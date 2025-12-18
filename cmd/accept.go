// Package cmd provides command-line interface implementations.
// This file implements the accept command which converts tasks.md to
// tasks.jsonc for machine-readable task tracking.
package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// filePerm is the standard file permission for created files (rw-r--r--)
const filePerm = 0644

// AcceptCmd represents the accept command for converting tasks.md to
// tasks.jsonc. This command parses the human-readable tasks.md file and
// produces a machine-readable tasks.jsonc file with structured task data.
type AcceptCmd struct {
	// ChangeID is the optional change identifier to process
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Change"` //nolint:lll,revive
	// DryRun enables preview mode without writing files
	DryRun bool `                                        help:"Preview without writing" name:"dry-run"` //nolint:lll,revive

	// NoInteractive disables interactive prompts
	NoInteractive bool `help:"Disable prompts" name:"no-interactive"` //nolint:lll,revive
}

// Run executes the accept command.
// It resolves the change ID, validates the change directory exists,
// and processes the tasks.md file to generate tasks.jsonc.
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

// processChange handles the conversion of tasks.md to tasks.jsonc.
// It validates that the change directory and tasks.md exist,
// validates the change, parses the markdown file, and writes the JSONC output.
func (c *AcceptCmd) processChange(
	projectRoot, changeID string,
) error {
	changeDir, tasksMdPath, err := resolveChangePaths(
		projectRoot,
		changeID,
	)
	if err != nil {
		return err
	}

	// Validate the change before conversion
	if err = c.runValidation(changeDir); err != nil {
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

	// Safety check: if tasks.md has content but no valid tasks were found
	if err := validateParsedTasks(tasks, tasksMdPath); err != nil {
		return err
	}

	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.jsonc",
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

	return writeAndCleanup(
		tasksMdPath,
		tasksJSONPath,
		tasks,
	)
}

// resolveChangePaths validates and returns the change directory and
// tasks.md path.
func resolveChangePaths(
	projectRoot, changeID string,
) (changeDir, tasksMdPath string, err error) {
	changeDir = filepath.Join(
		projectRoot,
		"spectr",
		"changes",
		changeID,
	)
	if _, err = os.Stat(changeDir); os.IsNotExist(
		err,
	) {
		return "", "", fmt.Errorf(
			"change directory not found: %s",
			changeDir,
		)
	}

	tasksMdPath = filepath.Join(
		changeDir,
		"tasks.md",
	)
	if _, err = os.Stat(tasksMdPath); os.IsNotExist(
		err,
	) {
		return "", "", fmt.Errorf(
			"tasks.md not found in change: %s",
			tasksMdPath,
		)
	}

	return changeDir, tasksMdPath, nil
}

// writeAndCleanup writes the tasks.jsonc file and removes tasks.md.
func writeAndCleanup(
	tasksMdPath, tasksJSONPath string,
	tasks []parsers.Task,
) error {
	if err := writeTasksJSONC(tasksJSONPath, tasks); err != nil {
		return fmt.Errorf(
			"failed to write tasks.jsonc: %w",
			err,
		)
	}

	// Remove tasks.md after successful tasks.jsonc creation
	if err := os.Remove(tasksMdPath); err != nil {
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
		return "", &specterrs.MissingChangeIDError{}
	}

	return selectChangeInteractive(projectRoot)
}

// taskParseState tracks state during tasks.md parsing.
type taskParseState struct {
	lastSectionNum   int    // Track for section numbering continuity
	sectionNum       string // Current section number as string
	sectionName      string // Current section name
	taskSeqInSection int    // Task sequence within section
	globalTaskSeq    int    // For tasks with no section
}

// handleSection processes a section header line and updates state.
func (s *taskParseState) handleSection(
	name, number string,
) {
	if number != "" {
		// Numbered section - use explicit number
		num, _ := strconv.Atoi(number)
		s.lastSectionNum = num
		s.sectionNum = number
	} else {
		// Unnumbered section - increment from last
		s.lastSectionNum++
		s.sectionNum = strconv.Itoa(s.lastSectionNum)
	}
	s.sectionName = strings.TrimSpace(name)
	s.taskSeqInSection = 0
}

// generateTaskID creates a task ID based on current state and match.
func (s *taskParseState) generateTaskID(
	matchNumber string,
) string {
	if s.sectionNum == "" {
		// No section context - global sequential
		s.globalTaskSeq++

		return strconv.Itoa(s.globalTaskSeq)
	}

	s.taskSeqInSection++
	expectedID := fmt.Sprintf(
		"%s.%d",
		s.sectionNum,
		s.taskSeqInSection,
	)

	if matchNumber != "" &&
		matchNumber == expectedID {
		return matchNumber // Explicit matches expected
	}

	return expectedID // Auto-generate or override
}

// createTask builds a Task from parsed match and current state.
func (s *taskParseState) createTask(
	match *markdown.FlexibleTaskMatch,
) parsers.Task {
	taskID := s.generateTaskID(match.Number)

	var status parsers.TaskStatusValue
	if match.Status == ' ' {
		status = parsers.TaskStatusPending
	} else {
		status = parsers.TaskStatusCompleted
	}

	return parsers.Task{
		ID:      taskID,
		Section: s.sectionName,
		Description: strings.TrimSpace(
			match.Content,
		),
		Status: status,
	}
}

// parseTasksMd parses tasks.md and returns a slice of Task structs.
// It extracts section headers (## lines), task IDs, descriptions, and status
// from the markdown structure. Supports flexible task formats:
//   - "- [ ] 1.1 Task" (decimal ID)
//   - "- [ ] 1. Task" (simple dot ID)
//   - "- [ ] 1 Task" (number only ID)
//   - "- [ ] Task" (no ID - auto-generated)
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
	state := &taskParseState{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for section header (numbered or unnumbered)
		if name, number, ok := markdown.MatchAnySection(line); ok {
			state.handleSection(name, number)

			continue
		}

		// Check for task line using flexible matching
		match, ok := markdown.MatchFlexibleTask(
			line,
		)
		if !ok {
			continue
		}

		tasks = append(
			tasks,
			state.createTask(match),
		)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf(
			"error reading file: %w",
			err,
		)
	}

	return tasks, nil
}

// validateParsedTasks checks if tasks.md has content but no valid tasks
// were found, which indicates a format mismatch.
func validateParsedTasks(
	tasks []parsers.Task,
	tasksMdPath string,
) error {
	if len(tasks) == 0 {
		info, statErr := os.Stat(tasksMdPath)
		if statErr == nil && info.Size() > 0 {
			return &specterrs.NoValidTasksError{
				TasksMdPath: tasksMdPath,
				FileSize:    info.Size(),
			}
		}
	}

	return nil
}
