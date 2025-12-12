// Package cmd provides command-line interface implementations.
// This file implements the accept command which converts tasks.md to
// tasks.json for machine-readable task tracking.
package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// filePerm is the standard file permission for created files (rw-r--r--)
const filePerm = 0644

// Regex patterns for parsing tasks.md - compiled once at package level
var (
	sectionPattern = regexp.MustCompile(
		`^##\s+\d+\.\s+(.+)$`,
	)
	taskPattern = regexp.MustCompile(
		`^-\s+\[([ xX])\]\s+(\d+\.\d+)\s+(.+)$`,
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

// parseTaskFromMatch creates a Task from regex match results.
func parseTaskFromMatch(
	matches []string,
	section string,
) parsers.Task {
	checkbox := matches[1]
	taskID := matches[2]
	description := strings.TrimSpace(matches[3])

	var status parsers.TaskStatusValue
	if checkbox == " " {
		status = parsers.TaskStatusPending
	} else {
		status = parsers.TaskStatusCompleted
	}

	return parsers.Task{
		ID:          taskID,
		Section:     section,
		Description: description,
		Status:      status,
	}
}

// parseTasksMd parses tasks.md and returns a slice of Task structs.
//
// It extracts section headers (## lines), task IDs, descriptions, and status
// from the markdown structure.
func parseTasksMd(
	path string,
) ([]parsers.Task, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to open file: %w",
			err,
		)
	}
	defer func() { _ = file.Close() }()

	var tasks []parsers.Task
	var currentSection string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for section header (e.g., "## 1. Core Accept Command")
		if matches := sectionPattern.FindStringSubmatch(line); matches != nil {
			currentSection = strings.TrimSpace(
				matches[1],
			)

			continue
		}

		// Check for task line (e.g., "- [ ] 1.1 Create `cmd/accept.go`...")
		matches := taskPattern.FindStringSubmatch(
			line,
		)
		if matches == nil {
			continue
		}

		tasks = append(
			tasks,
			parseTaskFromMatch(
				matches,
				currentSection,
			),
		)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf(
			"error reading file: %w",
			err,
		)
	}

	return tasks, nil
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
