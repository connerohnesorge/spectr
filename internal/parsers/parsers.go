// Package parsers provides utilities for extracting and counting
// information from markdown specification files, including titles,
// tasks, deltas, and requirements.
package parsers

import (
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parser"
)

// ExtractTitle extracts the title from a markdown file by finding
// the first H1 heading and removing "Change:" or "Spec:" prefix if present
func ExtractTitle(filePath string) (string, error) {
	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Parse the document
	doc, err := parser.Parse(string(content))
	if err != nil {
		return "", err
	}

	// Find the first H1 heading
	headers := parser.FindHeaders(doc, func(h *parser.Header) bool {
		return h.Level == 1
	})

	if len(headers) == 0 {
		return "", nil
	}

	// Extract and clean the title
	title := strings.TrimSpace(headers[0].Text)

	// Remove "Change:" or "Spec:" prefix
	title = strings.TrimPrefix(title, "Change:")
	title = strings.TrimPrefix(title, "Spec:")
	title = strings.TrimSpace(title)

	return title, nil
}

// TaskStatus represents task completion status
type TaskStatus struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

//nolint:revive // complexity acceptable for task counting

// CountTasks counts tasks in tasks.md, identifying completed vs total
//
//nolint:revive // cognitive complexity acceptable for task counting logic
func CountTasks(filePath string) (TaskStatus, error) {
	status := TaskStatus{Total: 0, Completed: 0}

	// Read the file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		// Return zero status if file doesn't exist or can't be read
		return status, nil
	}

	// Parse the document
	doc, err := parser.Parse(string(content))
	if err != nil {
		// Return zero status on parse error
		return status, nil
	}

	// Helper to check if a line is a task and count it
	countTaskLine := func(line string) {
		line = strings.TrimSpace(line)
		// Remove leading list marker if present
		// (handles indented lists that became text)
		isListMarker := strings.HasPrefix(line, "-") ||
			strings.HasPrefix(line, "*")
		if isListMarker {
			line = strings.TrimSpace(line[1:])
		}
		// Check for checkbox pattern: [ ] or [x]
		if len(line) < 3 || !strings.HasPrefix(line, "[") {
			return
		}

		// Extract the character inside the checkbox
		checkbox := line[1:2]
		if checkbox != " " && checkbox != "x" && checkbox != "X" {
			return
		}

		status.Total++
		if strings.ToLower(checkbox) == "x" {
			status.Completed++
		}
	}

	// Walk the AST and count task list items
	parser.Walk(doc, func(n parser.Node) bool {
		switch node := n.(type) {
		case *parser.List:
			// Check each list item for task checkbox pattern
			for _, item := range node.Items {
				countTaskLine(item)
			}
		case *parser.Paragraph:
			// Check paragraph text for indented list items
			// (lexer treats indented lists as text)
			lines := strings.Split(node.Text, "\n")
			for _, line := range lines {
				trimmed := strings.TrimSpace(line)
				// Check if line looks like indented list item
				isListMarker := strings.HasPrefix(trimmed, "-") ||
					strings.HasPrefix(trimmed, "*")
				if !isListMarker {
					continue
				}

				countTaskLine(trimmed)
			}
		}

		return true
	})

	return status, nil
}

// CountDeltas counts the number of delta sections
// (ADDED, MODIFIED, REMOVED, RENAMED) in change spec files
func CountDeltas(changeDir string) (int, error) {
	count := 0
	specsDir := changeDir + "/specs"

	// Check if specs directory exists
	if _, err := os.Stat(specsDir); os.IsNotExist(err) {
		return 0, nil
	}

	// Walk through all spec files in the specs directory
	err := walkSpecFiles(specsDir, func(filePath string) error {
		// Read the file content
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Parse the document
		doc, err := parser.Parse(string(content))
		if err != nil {
			// Return 0 on parse error but don't fail the walk
			return nil
		}

		// Extract delta operations
		deltas, err := parser.ExtractDeltas(doc)
		if err != nil {
			// Return 0 on extraction error but don't fail the walk
			return nil
		}

		// Count total delta operations
		if len(deltas.Added) > 0 {
			count++
		}
		if len(deltas.Modified) > 0 {
			count++
		}
		if len(deltas.Removed) > 0 {
			count++
		}
		if len(deltas.Renamed) > 0 {
			count++
		}

		return nil
	})

	return count, err
}

// CountRequirements counts the number of requirements in a spec.md file
func CountRequirements(specPath string) (int, error) {
	// Read the file content
	content, err := os.ReadFile(specPath)
	if err != nil {
		return 0, err
	}

	// Parse the document
	doc, err := parser.Parse(string(content))
	if err != nil {
		// Return 0 on parse error
		return 0, nil
	}

	// Extract requirements
	requirements, err := parser.ExtractRequirements(doc)
	if err != nil {
		// Return 0 on extraction error
		return 0, nil
	}

	return len(requirements), nil
}

// walkSpecFiles walks through all spec.md files in a directory tree
func walkSpecFiles(root string, fn func(string) error) error {
	entries, err := os.ReadDir(root)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		path := root + "/" + entry.Name()
		if entry.IsDir() {
			err = walkSpecFiles(path, fn)
			if err != nil {
				return err
			}
		} else if entry.Name() == "spec.md" {
			err = fn(path)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
