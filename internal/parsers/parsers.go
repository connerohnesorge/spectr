// Package parsers provides utilities for extracting and counting
// information from markdown specification files, including titles,
// tasks, deltas, and requirements.
package parsers

import (
	"bufio"
	"errors"
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
	Total     int `json:"total"`
	Completed int `json:"completed"`
}

// CountTasks counts tasks in tasks.md, identifying completed vs total
func CountTasks(filePath string) (TaskStatus, error) {
	status := TaskStatus{Total: 0, Completed: 0}

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

// TaskSection represents a numbered section in a tasks.md file
type TaskSection struct {
	Number    int    `json:"number"`     // The section number (1, 2, 3, etc.)
	Name      string `json:"name"`       // The section header text
	TaskCount int    `json:"task_count"` // Number of tasks in this section
	Line      int    `json:"line"`       // Line number where section starts
}

// TasksStructureResult contains validation results for a tasks.md file
type TasksStructureResult struct {
	// Sections is the list of numbered sections found
	Sections []TaskSection `json:"sections"`
	// OrphanedTasks is the count of tasks not under any numbered section
	OrphanedTasks int `json:"orphaned_tasks"`
	// EmptySections contains names of sections that have no tasks
	EmptySections []string `json:"empty_sections"`
	// SequentialNumbers indicates whether section numbers are sequential
	SequentialNumbers bool `json:"sequential_numbers"`
	// NonSequentialGaps lists which numbers are missing if not sequential
	NonSequentialGaps []int `json:"non_sequential_gaps"`
}

// ValidateTasksStructure parses a tasks.md file and validates its structure
func ValidateTasksStructure(filePath string) (*TasksStructureResult, error) {
	result := newTasksStructureResult()

	file, err := os.Open(filePath)
	if err != nil {
		// Only ignore missing files; propagate other I/O errors
		if errors.Is(err, os.ErrNotExist) {
			return result, nil
		}

		return nil, err
	}
	defer func() { _ = file.Close() }()

	parseTasksFile(file, result)

	if err := file.Close(); err != nil {
		return result, err
	}

	finalizeTasksResult(result)

	return result, nil
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
