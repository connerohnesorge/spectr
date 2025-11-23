// Package parsers provides utilities for extracting and counting
// information from markdown specification files, including titles,
// tasks, deltas, and requirements.
package parsers

import (
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/mdparser"
)

// ExtractTitle extracts the title from a markdown file by finding
// the first H1 heading and removing "Change:" or "Spec:" prefix if present
func ExtractTitle(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return "", err
	}

	// Find first H1 header
	for _, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != 1 {
			continue
		}

		title := strings.TrimSpace(header.Text)

		// Remove "Change:" or "Spec:" prefix
		title = strings.TrimPrefix(title, "Change:")
		title = strings.TrimPrefix(title, "Spec:")
		title = strings.TrimSpace(title)

		return title, nil
	}

	return "", nil
}

// TaskStatus represents task completion status
type TaskStatus struct {
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

// CountTasks counts tasks in tasks.md, identifying completed vs total
func CountTasks(filePath string) (TaskStatus, error) {
	status := TaskStatus{Total: 0, Completed: 0}

	content, err := os.ReadFile(filePath)
	if err != nil {
		// Return zero status if file doesn't exist or can't be read
		return status, nil
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return status, nil
	}

	// Traverse AST looking for list items with checkboxes
	for _, node := range doc.Children {
		switch n := node.(type) {
		case *mdparser.List:
			processListCheckboxes(n.Items, &status)
		case *mdparser.Paragraph:
			processParagraphCheckboxes(n.Lines, &status)
		}
	}

	return status, nil
}

// processListCheckboxes processes checkboxes in list items.
func processListCheckboxes(items []*mdparser.ListItem, status *TaskStatus) {
	for _, item := range items {
		updateStatusFromCheckbox(item.Text, status)
	}
}

// processParagraphCheckboxes processes checkboxes in paragraph lines.
// These handle indented list items that the lexer treats as paragraphs.
func processParagraphCheckboxes(lines []string, status *TaskStatus) {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this looks like a checkbox list item
		hasDash := strings.HasPrefix(trimmed, "- [")
		hasAsterisk := strings.HasPrefix(trimmed, "* [")
		if !hasDash && !hasAsterisk {
			continue
		}

		// Extract the checkbox part: skip "- " or "* "
		if len(trimmed) <= 2 {
			continue
		}

		updateStatusFromCheckbox(trimmed[2:], status)
	}
}

// updateStatusFromCheckbox checks for checkbox pattern and updates
// status.
//
//nolint:revive // text parameter is intentionally modified for clarity
func updateStatusFromCheckbox(text string, status *TaskStatus) {
	text = strings.TrimSpace(text)
	if !strings.HasPrefix(text, "[") {
		return
	}

	// Extract checkbox marker
	if len(text) < MinCheckboxLength {
		return
	}

	marker := text[1]
	if marker != ' ' && marker != 'x' && marker != 'X' {
		return
	}

	// This is a task
	status.Total++
	if marker == 'x' || marker == 'X' {
		status.Completed++
	}
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
		deltaCount, err := countDeltasInFile(filePath)
		if err != nil {
			return err
		}

		count += deltaCount

		return nil
	})

	return count, err
}

// countDeltasInFile counts delta sections in a single file.
func countDeltasInFile(filePath string) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return 0, err
	}

	count := 0
	deltaOps := []string{"ADDED", "MODIFIED", "REMOVED", "RENAMED"}

	for _, node := range doc.Children {
		if isDeltaSectionHeader(node, deltaOps) {
			count++
		}
	}

	return count, nil
}

// isDeltaSectionHeader checks if a node is a delta section header.
func isDeltaSectionHeader(node mdparser.Node, deltaOps []string) bool {
	header, ok := node.(*mdparser.Header)
	if !ok || header.Level != 2 {
		return false
	}

	headerText := strings.TrimSpace(header.Text)
	for _, op := range deltaOps {
		expectedText := op + " Requirements"
		if headerText == expectedText {
			return true
		}
	}

	return false
}

// CountRequirements counts the number of requirements in a spec.md file
func CountRequirements(specPath string) (int, error) {
	content, err := os.ReadFile(specPath)
	if err != nil {
		return 0, err
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return 0, err
	}

	count := 0

	// Traverse AST looking for H3 headers with "Requirement:" prefix
	for _, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != RequirementHeaderLevel {
			continue
		}

		if strings.HasPrefix(header.Text, "Requirement:") {
			count++
		}
	}

	return count, nil
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
