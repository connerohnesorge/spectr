// Package parsers provides utilities for extracting and counting
// information from markdown specification files, including titles,
// tasks, deltas, and requirements.
package parsers

import (
	"os"
	"strings"

	"github.com/yuin/goldmark/ast"
	extast "github.com/yuin/goldmark/extension/ast"
)

// ExtractTitle extracts the title from a markdown file by finding
// the first H1 heading and removing "Change:" or "Spec:" prefix if present
func ExtractTitle(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Parse markdown into AST
	doc := ParseMarkdown(content)

	// Find first H1 heading
	var title string
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// Check if this is an H1 heading
		if heading, ok := n.(*ast.Heading); ok && heading.Level == 1 {
			title = ExtractTextContent(heading)
			return ast.WalkStop, nil
		}

		return ast.WalkContinue, nil
	})

	if title == "" {
		return "", nil
	}

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

// CountTasks counts tasks in tasks.md, identifying completed vs total
func CountTasks(filePath string) (TaskStatus, error) {
	status := TaskStatus{Total: 0, Completed: 0}

	content, err := os.ReadFile(filePath)
	if err != nil {
		// Return zero status if file doesn't exist or can't be read
		return status, nil
	}

	// Parse markdown into AST (with GFM extension for task lists)
	doc := ParseMarkdown(content)

	// Walk the AST looking for TaskCheckBox nodes (from GFM extension)
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// Check if this is a TaskCheckBox node
		if taskBox, ok := n.(*extast.TaskCheckBox); ok {
			status.Total++
			if taskBox.IsChecked {
				status.Completed++
			}
		}

		return ast.WalkContinue, nil
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
		content, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		// Parse markdown into AST
		doc := ParseMarkdown(content)

		// Walk the AST looking for H2 delta section headers
		ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
			if !entering {
				return ast.WalkContinue, nil
			}

			// Check if this is an H2 heading
			if heading, ok := n.(*ast.Heading); ok && heading.Level == 2 {
				headingText := ExtractTextContent(heading)
				headingText = strings.TrimSpace(headingText)

				// Check if it matches delta operation pattern
				if strings.HasSuffix(headingText, " Requirements") {
					prefix := strings.TrimSuffix(headingText, " Requirements")
					prefix = strings.TrimSpace(prefix)
					if prefix == "ADDED" || prefix == "MODIFIED" ||
					   prefix == "REMOVED" || prefix == "RENAMED" {
						count++
					}
				}
			}

			return ast.WalkContinue, nil
		})

		return nil
	})

	return count, err
}

// CountRequirements counts the number of requirements in a spec.md file
func CountRequirements(specPath string) (int, error) {
	content, err := os.ReadFile(specPath)
	if err != nil {
		return 0, err
	}

	// Parse markdown into AST
	doc := ParseMarkdown(content)

	count := 0

	// Walk the AST looking for H3 requirement headers
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}

		// Check if this is an H3 heading
		if heading, ok := n.(*ast.Heading); ok && heading.Level == 3 {
			headingText := ExtractTextContent(heading)
			headingText = strings.TrimSpace(headingText)

			// Check if it starts with "Requirement:"
			if strings.HasPrefix(headingText, "Requirement:") {
				count++
			}
		}

		return ast.WalkContinue, nil
	})

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
