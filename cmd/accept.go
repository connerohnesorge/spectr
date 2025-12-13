// Package cmd provides command-line interface implementations.
// This file implements the accept command which converts tasks.md to
// tasks.json for machine-readable task tracking.
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// filePerm is the standard file permission for created files (rw-r--r--)
const filePerm = 0644

// Regex patterns for extracting task-specific data from markdown elements.
// These are ID extraction patterns, not structural parsing patterns.
// The markdown package handles structural parsing (headers, tasks).
var (
	// sectionIDPattern extracts section name from numbered H2 headers.
	// Example: "1. Core Accept Command" -> "Core Accept Command"
	sectionIDPattern = regexp.MustCompile(
		`^\d+\.\s+(.+)$`,
	)
	// taskIDPattern extracts task ID and description from task lines.
	// Example: "1.1 Create cmd/accept.go" -> ["1.1", "Create cmd/accept.go"]
	taskIDPattern = regexp.MustCompile(
		`^(\d+\.\d+)\s+(.+)$`,
	)
)

// AcceptCmd represents the accept command for converting tasks.md to
// tasks.json. This command parses the human-readable tasks.md file and
// produces a machine-readable tasks.json file with structured task data.
type AcceptCmd struct {
	// ChangeID is the optional change identifier to process
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change"` //nolint:lll,revive
	// DryRun enables preview mode without writing files
	DryRun bool `                                        help:"Preview without writing" name:"dry-run"` //nolint:lll,revive

	// NoInteractive disables interactive prompts
	NoInteractive bool `                                        help:"Disable prompts"         name:"no-interactive"` //nolint:lll,revive
}

// Run executes the accept command.
// It resolves the change ID, validates the change directory exists,
// and processes the tasks.md file to generate tasks.json.
func (c *AcceptCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf(
			"failed to get current directory: %w",
			err,
		)
	}

	changeID, err := c.resolveChangeID(
		projectRoot,
	)
	if err != nil {
		var userCancelledErr *specterrs.UserCancelledError
		if errors.As(err, &userCancelledErr) {
			return nil
		}

		return err
	}

	return c.processChange(projectRoot, changeID)
}

// processChange handles the conversion of tasks.md to tasks.json.
// It validates that the change directory and tasks.md exist,
// validates the change, parses the markdown file, and writes the JSON output.
func (c *AcceptCmd) processChange(
	projectRoot, changeID string,
) error {
	changeDir := filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		changeID,
	)
	_, err := os.Stat(changeDir)
	if os.IsNotExist(err) {
		return fmt.Errorf(
			"change directory not found: %s",
			changeDir,
		)
	}

	tasksMdPath := filepath.Join(
		changeDir,
		"tasks.md",
	)
	_, err = os.Stat(tasksMdPath)
	if os.IsNotExist(
		err,
	) {
		return fmt.Errorf(
			"tasks.md not found in change: %s",
			tasksMdPath,
		)
	}

	// Validate the change before conversion
	err = c.runValidation(changeDir)
	if err != nil {
		return fmt.Errorf(
			"validation failed: %w",
			err,
		)
	}

	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		return fmt.Errorf(
			"failed to parse tasks.md: %w",
			err,
		)
	}

	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.json",
	)

	if c.DryRun {
		fmt.Printf(
			"Would convert: %s\nwould write to: %s\nWould remove: %s\nFound %d tasks\n", //nolint:lll,revive
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			len(tasks),
		)

		return nil
	}

	err = writeTasksJSON(tasksJSONPath, tasks)
	if err != nil {
		return fmt.Errorf(
			"failed to write tasks.json: %w",
			err,
		)
	}

	// Remove tasks.md after successful tasks.json creation
	err = os.Remove(tasksMdPath)
	if err != nil {
		return fmt.Errorf(
			"failed to remove tasks.md: %w",
			err,
		)
	}

	fmt.Printf(
		"Converted %s -> %s\nRemoved %s\nWrote %d tasks\n",
		tasksMdPath,
		tasksJSONPath,
		tasksMdPath,
		len(tasks),
	)

	return nil
}

// runValidation validates the change before accepting
func (*AcceptCmd) runValidation(
	changeDir string,
) error {
	fmt.Println("Validating change...")

	report, err := archive.ValidatePreArchive(
		changeDir,
		true,
	)
	if err != nil {
		return err
	}

	if !report.Valid {
		fmt.Printf(
			"Validation failed: %d error(s), %d warning(s)\n",
			report.Summary.Errors,
			report.Summary.Warnings,
		)

		for _, issue := range report.Issues {
			fmt.Printf("  [%s] %s: %s\n",
				issue.Level,
				issue.Path,
				issue.Message,
			)
		}

		return &specterrs.ValidationRequiredError{
			Operation: "accepting",
		}
	}

	if report.Summary.Warnings > 0 {
		fmt.Printf(
			"Validation passed with %d warning(s)\n",
			report.Summary.Warnings,
		)
	} else {
		fmt.Println("Validation passed")
	}

	return nil
}

// resolveChangeID resolves the change ID from argument or interactive
// selection. If a change ID is provided, it uses partial matching to
// resolve the full ID. Otherwise, it prompts for interactive selection
// (unless NoInteractive is set).
func (c *AcceptCmd) resolveChangeID(
	projectRoot string,
) (string, error) {
	if c.ChangeID != "" {
		// Normalize path to extract change ID
		normalizedID, _ := discovery.NormalizeItemPath(
			c.ChangeID,
		)

		result, err := discovery.ResolveChangeID(
			normalizedID,
			projectRoot,
		)
		if err != nil {
			return "", err
		}

		if result.PartialMatch {
			fmt.Printf(
				"Resolved '%s' -> '%s'\n\n",
				c.ChangeID,
				result.ChangeID,
			)
		}

		return result.ChangeID, nil
	}

	if c.NoInteractive {
		// TODO: Define error type for this?
		return "", errors.New(
			"usage: spectr accept <change-id> [flags]\n" +
				"       spectr accept <change-id> --dry-run",
		)
	}

	return selectChangeInteractive(projectRoot)
}

// parseTaskFromMarkdown creates a parsers.Task from a markdown.Task.
// It extracts the task ID and description from the task line text.
// Returns nil if the task line doesn't match the expected format.
func parseTaskFromMarkdown(
	mdTask markdown.Task,
	section string,
) *parsers.Task {
	// Extract the text after the checkbox from the full line
	// Line format: "- [ ] 1.1 Description" or "- [x] 1.1 Description"
	line := strings.TrimSpace(mdTask.Line)

	// Remove the list marker and checkbox
	// Find the checkbox end position
	checkboxEnd := strings.Index(line, "] ")
	if checkboxEnd == -1 {
		return nil
	}
	textAfterCheckbox := strings.TrimSpace(line[checkboxEnd+2:])

	// Extract task ID and description using the ID pattern
	matches := taskIDPattern.FindStringSubmatch(textAfterCheckbox)
	if matches == nil {
		return nil
	}

	taskID := matches[1]
	description := strings.TrimSpace(matches[2])

	var status parsers.TaskStatusValue
	if mdTask.Checked {
		status = parsers.TaskStatusCompleted
	} else {
		status = parsers.TaskStatusPending
	}

	return &parsers.Task{
		ID:          taskID,
		Section:     section,
		Description: description,
		Status:      status,
	}
}

// parseTasksMd parses tasks.md and returns a slice of Task structs.
// It uses the markdown package for structural parsing (headers, tasks),
// then extracts task IDs and descriptions from the parsed elements.
func parseTasksMd(path string) ([]parsers.Task, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	doc, err := markdown.ParseDocument(content)
	if err != nil {
		// Empty file is valid - just return no tasks
		var emptyErr *specterrs.EmptyContentError
		if errors.As(err, &emptyErr) {
			return nil, nil
		}

		return nil, fmt.Errorf("failed to parse markdown: %w", err)
	}

	// Build section map: line number -> section name (numbered H2 headers)
	sections := make(map[int]string)
	for _, h := range doc.Headers {
		if h.Level != 2 {
			continue
		}

		if m := sectionIDPattern.FindStringSubmatch(h.Text); m != nil {
			sections[h.Line] = strings.TrimSpace(m[1])
		}
	}

	var tasks []parsers.Task
	var processTasks func(mdTasks []markdown.Task)
	processTasks = func(mdTasks []markdown.Task) {
		for _, mdTask := range mdTasks {
			section := findSection(mdTask.LineNum, sections)
			if task := parseTaskFromMarkdown(mdTask, section); task != nil {
				tasks = append(tasks, *task)
			}

			processTasks(mdTask.Children)
		}
	}
	processTasks(doc.Tasks)

	return tasks, nil
}

// findSection finds the section name for a task at a given line number.
func findSection(taskLine int, sections map[int]string) string {
	var result string
	var maxLine int

	for line, section := range sections {
		if line < taskLine && line > maxLine {
			maxLine, result = line, section
		}
	}

	return result
}

// writeTasksJSON writes tasks to a tasks.json file.
// The output follows the TasksFile structure defined in parsers/types.go.
func writeTasksJSON(
	path string,
	tasks []parsers.Task,
) error {
	tasksFile := parsers.TasksFile{
		Version: 1,
		Tasks:   tasks,
	}

	jsonData, err := json.MarshalIndent(
		tasksFile,
		"",
		"  ",
	)
	if err != nil {
		return fmt.Errorf(
			"failed to marshal tasks to JSON: %w",
			err,
		)
	}

	if err := os.WriteFile(path, jsonData, filePerm); err != nil {
		return fmt.Errorf(
			"failed to write tasks.json: %w",
			err,
		)
	}

	return nil
}
