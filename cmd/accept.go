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

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// filePerm is the standard file permission for created files (rw-r--r--)
const filePerm = 0644

// AcceptCmd represents the accept command for converting tasks.md to
// tasks.json. This command parses the human-readable tasks.md file and
// produces a machine-readable tasks.json file with structured task data.
type AcceptCmd struct {
	// ChangeID is the optional change identifier to process
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change"`
	// DryRun enables preview mode without writing files
	DryRun bool `name:"dry-run" help:"Preview without writing"`
	// NoInteractive disables interactive prompts
	NoInteractive bool `name:"no-interactive" help:"Disable prompts"`
}

// Run executes the accept command.
// It resolves the change ID, validates the change directory exists,
// and processes the tasks.md file to generate tasks.json.
func (c *AcceptCmd) Run() error {
	projectRoot, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}

	changeID, err := c.resolveChangeID(projectRoot)
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
func (c *AcceptCmd) processChange(projectRoot, changeID string) error {
	changeDir := filepath.Join(projectRoot, "spectr", "changes", changeID)
	if _, err := os.Stat(changeDir); os.IsNotExist(err) {
		return fmt.Errorf("change directory not found: %s", changeDir)
	}

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if _, err := os.Stat(tasksMdPath); os.IsNotExist(err) {
		return fmt.Errorf("tasks.md not found in change: %s", tasksMdPath)
	}

	// Validate the change before conversion
	if err := c.runValidation(changeDir); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		return fmt.Errorf("failed to parse tasks.md: %w", err)
	}

	tasksJSONPath := filepath.Join(changeDir, "tasks.json")

	if c.DryRun {
		fmt.Printf("Would convert: %s\n", tasksMdPath)
		fmt.Printf("Would write to: %s\n", tasksJSONPath)
		fmt.Printf("Would remove: %s\n", tasksMdPath)
		fmt.Printf("Found %d tasks\n", len(tasks))

		return nil
	}

	if err := writeTasksJSON(tasksJSONPath, tasks); err != nil {
		return fmt.Errorf("failed to write tasks.json: %w", err)
	}

	// Remove tasks.md after successful tasks.json creation
	if err := os.Remove(tasksMdPath); err != nil {
		return fmt.Errorf("failed to remove tasks.md: %w", err)
	}

	fmt.Printf("Converted %s -> %s\n", tasksMdPath, tasksJSONPath)
	fmt.Printf("Removed %s\n", tasksMdPath)
	fmt.Printf("Wrote %d tasks\n", len(tasks))

	return nil
}

// runValidation validates the change before accepting
func (*AcceptCmd) runValidation(changeDir string) error {
	fmt.Println("Validating change...")

	report, err := archive.ValidatePreArchive(changeDir, true)
	if err != nil {
		return err
	}

	if !report.Valid {
		fmt.Printf("Validation failed: %d error(s), %d warning(s)\n",
			report.Summary.Errors, report.Summary.Warnings)

		for _, issue := range report.Issues {
			fmt.Printf("  [%s] %s: %s\n",
				issue.Level,
				issue.Path,
				issue.Message,
			)
		}

		return &specterrs.ValidationRequiredError{Operation: "accepting"}
	}

	if report.Summary.Warnings > 0 {
		fmt.Printf("Validation passed with %d warning(s)\n",
			report.Summary.Warnings)
	} else {
		fmt.Println("Validation passed")
	}

	return nil
}

// resolveChangeID resolves the change ID from argument or interactive
// selection. If a change ID is provided, it uses partial matching to
// resolve the full ID. Otherwise, it prompts for interactive selection
// (unless NoInteractive is set).
func (c *AcceptCmd) resolveChangeID(projectRoot string) (string, error) {
	if c.ChangeID != "" {
		result, err := discovery.ResolveChangeID(c.ChangeID, projectRoot)
		if err != nil {
			return "", err
		}

		if result.PartialMatch {
			fmt.Printf("Resolved '%s' -> '%s'\n\n", c.ChangeID, result.ChangeID)
		}

		return result.ChangeID, nil
	}

	if c.NoInteractive {
		return "", errors.New(
			"usage: spectr accept <change-id> [flags]\n" +
				"       spectr accept <change-id> --dry-run",
		)
	}

	return selectChangeInteractive(projectRoot)
}

// parseTasksMd parses tasks.md and returns a slice of Task structs.
// It uses the markdown package to extract section headers, task IDs,
// descriptions, and status from the markdown structure.
func parseTasksMd(path string) ([]parsers.Task, error) {
	tasksWithIDs, err := markdown.ExtractTasksWithIDsFromFile(path)
	if err != nil {
		return nil, err
	}

	// Return nil for empty results to maintain backwards compatibility
	if len(tasksWithIDs) == 0 {
		return nil, nil
	}

	tasks := make([]parsers.Task, len(tasksWithIDs))
	for i, t := range tasksWithIDs {
		var status parsers.TaskStatusValue
		if t.Checked {
			status = parsers.TaskStatusCompleted
		} else {
			status = parsers.TaskStatusPending
		}

		tasks[i] = parsers.Task{
			ID:          t.ID,
			Section:     t.Section,
			Description: t.Description,
			Status:      status,
		}
	}

	return tasks, nil
}

// writeTasksJSON writes tasks to a tasks.json file.
// The output follows the TasksFile structure defined in parsers/types.go.
func writeTasksJSON(path string, tasks []parsers.Task) error {
	tasksFile := parsers.TasksFile{
		Version: 1,
		Tasks:   tasks,
	}

	jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks to JSON: %w", err)
	}

	if err := os.WriteFile(path, jsonData, filePerm); err != nil {
		return fmt.Errorf("failed to write tasks.json: %w", err)
	}

	return nil
}
