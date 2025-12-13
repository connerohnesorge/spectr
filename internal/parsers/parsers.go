// Package parsers provides utilities for extracting and counting
// information from markdown specification files, including titles,
// tasks, deltas, and requirements.
package parsers

import (
	"bufio"
	"encoding/json"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/markdown"
)

// Markdown header levels
const (
	headerLevelH2 = 2
	headerLevelH3 = 3
)

// ExtractTitle extracts the title from a markdown file by finding
// the first H1 heading and removing "Change:" or "Spec:" prefix if present
func ExtractTitle(
	filePath string,
) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Look for H1 heading (# Title)
		if !strings.HasPrefix(line, "# ") {
			continue
		}
		title := strings.TrimPrefix(line, "# ")
		title = strings.TrimSpace(title)

		// Remove "Change:" or "Spec:" prefix
		title = strings.TrimPrefix(
			title,
			"Change:",
		)
		title = strings.TrimPrefix(title, "Spec:")
		title = strings.TrimSpace(title)

		return title, nil
	}

	return "", scanner.Err()
}

// TaskStatus represents task completion status
type TaskStatus struct {
	Total      int `json:"total"`
	Completed  int `json:"completed"`
	InProgress int `json:"in_progress"`
}

// ReadTasksJson reads and parses a tasks.json file
func ReadTasksJson(
	filePath string,
) (*TasksFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tasksFile TasksFile
	if err := json.Unmarshal(data, &tasksFile); err != nil {
		return nil, err
	}

	return &tasksFile, nil
}

// CountTasks counts tasks in a change directory, checking tasks.json first
// and falling back to tasks.md if tasks.json doesn't exist
func CountTasks(
	changeDir string,
) (TaskStatus, error) {
	// First, try to read tasks.json
	tasksJsonPath := changeDir + "/tasks.json"
	if _, err := os.Stat(tasksJsonPath); err == nil {
		return countTasksFromJson(tasksJsonPath)
	}

	// Fall back to tasks.md
	tasksMdPath := changeDir + "/tasks.md"

	return countTasksFromMarkdown(tasksMdPath)
}

// countTasksFromJson counts tasks from a tasks.json file
func countTasksFromJson(
	filePath string,
) (TaskStatus, error) {
	status := TaskStatus{
		Total:      0,
		Completed:  0,
		InProgress: 0,
	}

	tasksFile, err := ReadTasksJson(filePath)
	if err != nil {
		return status, err
	}

	status.Total = len(tasksFile.Tasks)
	for _, task := range tasksFile.Tasks {
		switch task.Status {
		case TaskStatusCompleted:
			status.Completed++
		case TaskStatusInProgress:
			status.InProgress++
		case TaskStatusPending:
			// Pending tasks are counted in Total, no separate counter needed
		}
	}

	return status, nil
}

// countTasksFromMarkdown counts tasks from a tasks.md file
func countTasksFromMarkdown(
	filePath string,
) (TaskStatus, error) {
	status := TaskStatus{
		Total:      0,
		Completed:  0,
		InProgress: 0,
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		// Return zero status if file doesn't exist or can't be read
		return status, nil
	}

	doc, err := markdown.ParseDocument(content)
	if err != nil {
		// Return zero status if content is empty or invalid
		return status, nil
	}

	// Count all tasks recursively (including nested children)
	countTasksRecursive(doc.Tasks, &status)

	return status, nil
}

// countTasksRecursive counts tasks and their children recursively
func countTasksRecursive(tasks []markdown.Task, status *TaskStatus) {
	for _, task := range tasks {
		status.Total++
		if task.Checked {
			status.Completed++
		}
		// Count nested tasks
		countTasksRecursive(task.Children, status)
	}
}

// CountDeltas counts the number of delta sections
// (ADDED, MODIFIED, REMOVED, RENAMED) in change spec files
func CountDeltas(changeDir string) (int, error) {
	count := 0
	specsDir := changeDir + "/specs"

	// Check if specs directory exists
	if _, err := os.Stat(specsDir); os.IsNotExist(
		err,
	) {
		return 0, nil
	}

	// Walk through all spec files in the specs directory
	err := walkSpecFiles(
		specsDir,
		func(filePath string) error {
			content, err := os.ReadFile(filePath)
			if err != nil {
				return err
			}

			doc, err := markdown.ParseDocument(content)
			if err != nil {
				// Skip files that can't be parsed (empty, binary, etc.)
				return nil
			}

			// Look for H2 headers that match delta patterns
			for _, header := range doc.Headers {
				if header.Level != 2 {
					continue
				}
				// Check if header matches delta pattern:
				// "ADDED Requirements", "MODIFIED Requirements", etc.
				if isDeltaHeader(header.Text) {
					count++
				}
			}

			return nil
		},
	)

	return count, err
}

// isDeltaHeader checks if a header text matches a delta section pattern
func isDeltaHeader(text string) bool {
	deltaTypes := []string{"ADDED", "MODIFIED", "REMOVED", "RENAMED"}
	for _, deltaType := range deltaTypes {
		if strings.HasPrefix(text, deltaType+" Requirements") {
			return true
		}
	}

	return false
}

// CountRequirements counts the number of requirements in a spec.md file
func CountRequirements(
	specPath string,
) (int, error) {
	content, err := os.ReadFile(specPath)
	if err != nil {
		return 0, err
	}

	doc, err := markdown.ParseDocument(content)
	if err != nil {
		// Return 0 if content is empty or invalid
		return 0, nil
	}

	count := 0
	// Look for H3 headers that start with "Requirement:"
	for _, header := range doc.Headers {
		isH3 := header.Level == headerLevelH3
		isReq := strings.HasPrefix(header.Text, "Requirement:")
		if isH3 && isReq {
			count++
		}
	}

	return count, nil
}

// walkSpecFiles walks through all spec.md files in a directory tree
func walkSpecFiles(
	root string,
	fn func(string) error,
) error {
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
