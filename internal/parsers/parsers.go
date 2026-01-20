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

// StripJSONComments removes JSONC-style comments from JSON content.
// Handles single-line (//) and multi-line (/* */) comments.
// Comments inside strings are preserved (not stripped).
func StripJSONComments(data []byte) []byte {
	result := make([]byte, 0, len(data))
	i := 0

	for i < len(data) {
		switch {
		case data[i] == '"':
			i = copyJSONString(data, i, &result)
		case isLineComment(data, i):
			i = skipLineComment(data, i)
		case isBlockComment(data, i):
			i = skipBlockComment(data, i)
		default:
			result = append(result, data[i])
			i++
		}
	}

	return result
}

func isLineComment(data []byte, i int) bool {
	return i+1 < len(data) && data[i] == '/' &&
		data[i+1] == '/'
}

func isBlockComment(data []byte, i int) bool {
	return i+1 < len(data) && data[i] == '/' &&
		data[i+1] == '*'
}

func skipLineComment(data []byte, start int) int {
	pos := start + 2
	for pos < len(data) && data[pos] != '\n' {
		pos++
	}

	return pos
}

func skipBlockComment(
	data []byte,
	start int,
) int {
	pos := start + 2
	for pos+1 < len(data) {
		if data[pos] == '*' &&
			data[pos+1] == '/' {
			return pos + 2
		}
		pos++
	}

	return pos
}

func copyJSONString(
	data []byte,
	start int,
	result *[]byte,
) int {
	*result = append(*result, data[start])
	pos := start + 1

	for pos < len(data) {
		if data[pos] == '\\' &&
			pos+1 < len(data) {
			*result = append(
				*result,
				data[pos],
				data[pos+1],
			)
			pos += 2

			continue
		}
		*result = append(*result, data[pos])
		if data[pos] == '"' {
			return pos + 1
		}
		pos++
	}

	return pos
}

// ReadTasksJson reads and parses a tasks.json file.
// Supports JSONC format with single-line and multi-line comments, and trailing commas.
func ReadTasksJson(
	filePath string,
) (*TasksFile, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Convert JSONC to standard JSON (handles comments AND trailing commas)
	data = JSONCToJSON(data)

	var tasksFile TasksFile
	if err := json.Unmarshal(data, &tasksFile); err != nil {
		return nil, err
	}

	return &tasksFile, nil
}

// CountTasks counts tasks in a change directory, checking tasks.jsonc first
// and falling back to tasks.md if tasks.jsonc doesn't exist.
// NOTE: Legacy tasks.json files are silently ignored (hard break).
func CountTasks(
	changeDir string,
) (TaskStatus, error) {
	// First, try to read tasks.jsonc (new format with comment support)
	tasksJsoncPath := changeDir + "/tasks.jsonc"
	if _, err := os.Stat(tasksJsoncPath); err == nil {
		return countTasksFromJson(tasksJsoncPath)
	}

	// Fall back to tasks.md (for not-yet-accepted changes)
	// Note: tasks.json is NOT checked - legacy files are ignored
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

	file, err := os.Open(filePath)
	if err != nil {
		// Return zero status if file doesn't exist or can't be read
		return status, nil
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		state, ok := markdown.MatchTaskCheckbox(
			line,
		)
		if !ok {
			continue
		}
		status.Total++
		if markdown.IsTaskChecked(state) {
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
	if _, err := os.Stat(specsDir); os.IsNotExist(
		err,
	) {
		return 0, nil
	}

	// Walk through all spec files in the specs directory
	err := walkSpecFiles(
		specsDir,
		func(filePath string) error {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer func() { _ = file.Close() }()

			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := strings.TrimSpace(
					scanner.Text(),
				)
				if _, ok := markdown.MatchH2DeltaSection(line); ok {
					count++
				}
			}

			return scanner.Err()
		},
	)

	return count, err
}

// CountRequirements counts the number of requirements in a spec.md file
func CountRequirements(
	specPath string,
) (int, error) {
	file, err := os.Open(specPath)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	count := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if _, ok := markdown.MatchRequirementHeader(line); ok {
			count++
		}
	}

	return count, scanner.Err()
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
