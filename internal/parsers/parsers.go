// Package parsers provides utilities for extracting and counting
// information from markdown specification files, including titles,
// tasks, deltas, and requirements.
package parsers

import (
	"bufio"
	"encoding/json"
	"os"
	"regexp"
	"strings"
)

// ExtractTitle extracts the title from a markdown file by finding
// the first H1 heading and removing "Change:" or "Spec:" prefix if present
func ExtractTitle(filePath string) (string, error) {
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
		title = strings.TrimPrefix(title, "Change:")
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
func ReadTasksJson(filePath string) (*TasksFile, error) {
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
func CountTasks(changeDir string) (TaskStatus, error) {
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
func countTasksFromJson(filePath string) (TaskStatus, error) {
	status := TaskStatus{Total: 0, Completed: 0, InProgress: 0}

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
func countTasksFromMarkdown(filePath string) (TaskStatus, error) {
	status := TaskStatus{Total: 0, Completed: 0, InProgress: 0}

	file, err := os.Open(filePath)
	if err != nil {
		// Return zero status if file doesn't exist or can't be read
		return status, nil
	}
	defer func() { _ = file.Close() }()

	// Regex to match task lines: - [ ] or - [x] (case-insensitive)
	taskPattern := regexp.MustCompile(`^\s*-\s*\[([xX ])\]`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches := taskPattern.FindStringSubmatch(line)
		if len(matches) <= 1 {
			continue
		}
		status.Total++
		marker := strings.ToLower(strings.TrimSpace(matches[1]))
		if marker == "x" {
			status.Completed++
		}
	}

	return status, scanner.Err()
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
		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()

		// Match delta section headers
		deltaPattern := regexp.MustCompile(
			`^##\s+(ADDED|MODIFIED|REMOVED|RENAMED)\s+Requirements`,
		)

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if deltaPattern.MatchString(line) {
				count++
			}
		}

		return scanner.Err()
	})

	return count, err
}

// CountRequirements counts the number of requirements in a spec.md file
func CountRequirements(specPath string) (int, error) {
	file, err := os.Open(specPath)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	count := 0
	reqPattern := regexp.MustCompile(`^###\s+Requirement:`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if reqPattern.MatchString(line) {
			count++
		}
	}

	return count, scanner.Err()
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
