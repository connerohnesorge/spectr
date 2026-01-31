// Package cmd provides command-line interface implementations.
// This file implements the accept command which converts tasks.md to
// tasks.jsonc for machine-readable task tracking.
//
//nolint:revive // file-length-limit: accept command requires cohesive task parsing logic
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
	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
	"github.com/connerohnesorge/spectr/internal/validation"
)

// filePerm is the standard file permission for created files (rw-r--r--)
const filePerm = 0o644

// AcceptCmd represents the accept command for converting tasks.md to
// tasks.jsonc. This command parses the human-readable tasks.md file and
// produces a machine-readable tasks.jsonc file with structured task data.
// Both files are preserved after conversion: tasks.md remains as the
// human-readable source, while tasks.jsonc becomes the runtime source of truth.
// If tasks.jsonc is deleted, all commands automatically fall back to tasks.md.
//
// Note: A --sync-from-md flag is not needed because running accept again
// is idempotent and already handles re-generating tasks.jsonc from tasks.md.
type AcceptCmd struct {
	// ChangeID is the optional change identifier to process
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Convert tasks.md to tasks.jsonc (preserves tasks.md)"` //nolint:lll,revive // Kong struct tag exceeds line length
	// DryRun enables preview mode without writing files
	DryRun bool `                                        help:"Preview without writing"                              name:"dry-run"` //nolint:lll,revive // Kong struct tag with alignment padding

	// NoInteractive disables interactive prompts
	NoInteractive bool `help:"Disable prompts" name:"no-interactive"` //nolint:lll,revive // Kong struct tag exceeds line length
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

	// Load project configuration (optional)
	cfg, err := config.LoadConfig(projectRoot)
	if err != nil {
		return fmt.Errorf(
			"failed to load config: %w",
			err,
		)
	}

	// Validate the change before conversion
	if err = c.runValidation(changeDir); err != nil {
		return fmt.Errorf(
			"validation failed: %w",
			err,
		)
	}

	// Check proposal dependencies (chained proposals)
	// This is a hard fail - required dependencies must be archived
	if depErr := c.checkDependencies(projectRoot, changeID); depErr != nil {
		return depErr
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

	// Append configured tasks if present
	var appendCfg *config.AppendTasksConfig
	if cfg != nil && cfg.AppendTasks != nil &&
		cfg.AppendTasks.HasTasks() {
		appendCfg = cfg.AppendTasks
	}

	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.jsonc",
	)

	if c.DryRun {
		allTasks := tasks
		if appendCfg != nil && appendCfg.HasTasks() {
			appendedTasks := createAppendedTasks(tasks, appendCfg)
			allTasks = append(allTasks, appendedTasks...)
		}

		shouldSplit := shouldSplitTasksJSONC(allTasks)
		formatType := "flat (v1)"
		if shouldSplit {
			formatType = "hierarchical (v2)"
		}
		totalTasks := len(allTasks)

		fmt.Printf(
			"Would convert: %s\nwould write to: %s\nWould preserve: %s\nFound %d tasks\nFormat: %s\n", //nolint:lll,revive // Long format string for dry-run output
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			totalTasks,
			formatType,
		)

		return nil
	}

	// Check if we should split into hierarchical format
	// Split based on generated JSONC complexity, not tasks.md line count
	allTasks := tasks
	if appendCfg != nil && appendCfg.HasTasks() {
		appendedTasks := createAppendedTasks(tasks, appendCfg)
		allTasks = append(allTasks, appendedTasks...)
	}

	shouldSplit := shouldSplitTasksJSONC(allTasks)
	if shouldSplit {
		return writeAndCleanupHierarchical(
			changeID,
			changeDir,
			tasksMdPath,
			tasks,
			appendCfg,
		)
	}

	return writeAndCleanup(
		tasksMdPath,
		tasksJSONPath,
		tasks,
		appendCfg,
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

// writeAndCleanup writes the tasks.jsonc file and preserves tasks.md.
// We preserve tasks.md to avoid information loss (markdown formatting,
// comments, links) during the conversion to tasks.jsonc. Both files
// will coexist, with tasks.jsonc serving as the machine-readable format
// and tasks.md as the human-readable source of truth.
func writeAndCleanup(
	tasksMdPath, tasksJSONPath string,
	tasks []parsers.Task,
	appendCfg *config.AppendTasksConfig,
) error {
	if err := writeTasksJSONC(tasksJSONPath, tasks, appendCfg); err != nil {
		return fmt.Errorf(
			"failed to write tasks.jsonc: %w",
			err,
		)
	}

	// Preserve tasks.md to avoid information loss (formatting, comments, links)
	// Both tasks.md and tasks.jsonc now coexist after conversion

	totalTasks := len(tasks)
	if appendCfg != nil {
		totalTasks += len(appendCfg.Tasks)
	}

	fmt.Printf(
		"Converted %s -> %s\nPreserved %s\nWrote %d tasks\n",
		tasksMdPath,
		tasksJSONPath,
		tasksMdPath,
		totalTasks,
	)

	return nil
}

// writeAndCleanupHierarchical writes hierarchical v2 tasks.jsonc structure.
// Creates a root tasks.jsonc plus child tasks-{N}.jsonc files for each section.
// Preserves task status from existing files when re-running accept.
func writeAndCleanupHierarchical(
	changeID, changeDir, tasksMdPath string,
	tasks []parsers.Task,
	appendCfg *config.AppendTasksConfig,
) error {
	// Append configured tasks if present
	allTasks := tasks
	if appendCfg != nil && appendCfg.HasTasks() {
		appendedTasks := createAppendedTasks(tasks, appendCfg)
		allTasks = append(allTasks, appendedTasks...)
	}

	// Build status map from existing files
	statusMap := buildTaskStatusMap(changeDir)

	// Group tasks by section
	sections := groupTasksBySection(allTasks)

	// Write hierarchical structure
	if err := writeHierarchicalTasksJSONC(changeDir, changeID, sections, statusMap); err != nil {
		return fmt.Errorf("failed to write hierarchical tasks: %w", err)
	}

	// Count total tasks
	totalTasks := len(allTasks)

	fmt.Printf(
		"Converted %s -> hierarchical tasks.jsonc (v2)\nPreserved %s\nWrote %d tasks across %d sections\n",
		tasksMdPath,
		tasksMdPath,
		totalTasks,
		len(sections),
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

// sectionGroup represents a group of tasks under a common section
type sectionGroup struct {
	sectionNum  string // Section number (e.g., "1", "2")
	sectionName string // Section name (e.g., "Template Infrastructure")
	tasks       []parsers.Task
}

// groupTasksBySection groups tasks by their section number.
// Tasks are grouped by the first part of their ID (e.g., "1.1" -> "1").
// Tasks without sections or with non-standard IDs go into section "0".
func groupTasksBySection(tasks []parsers.Task) []sectionGroup {
	// Map from section number to tasks
	sectionMap := make(map[string]*sectionGroup)
	var sectionOrder []string // Track insertion order

	for _, task := range tasks {
		sectionNum := extractSectionNumber(task.ID)

		if _, exists := sectionMap[sectionNum]; !exists {
			// Create new section group
			sectionMap[sectionNum] = &sectionGroup{
				sectionNum:  sectionNum,
				sectionName: task.Section,
				tasks:       []parsers.Task{},
			}
			sectionOrder = append(sectionOrder, sectionNum)
		}

		sectionMap[sectionNum].tasks = append(sectionMap[sectionNum].tasks, task)
	}

	// Convert map to ordered slice
	result := make([]sectionGroup, 0, len(sectionMap))
	for _, num := range sectionOrder {
		result = append(result, *sectionMap[num])
	}

	return result
}

// extractSectionNumber extracts the section number from a task ID.
// Examples: "1.1" -> "1", "2.3" -> "2", "5" -> "0" (no subsection)
// Returns "0" for tasks without a clear section structure.
func extractSectionNumber(taskID string) string {
	parts := strings.Split(taskID, ".")
	if len(parts) >= 2 {
		// Has subsection - return first part
		return parts[0]
	}
	// No subsection - return "0" (preliminary tasks)
	return "0"
}

// shouldSplitTasksJSONC determines if tasks should be split into hierarchical format.
// Split occurs when the generated JSONC would be large (> 20 tasks or multiple sections).
func shouldSplitTasksJSONC(tasks []parsers.Task) bool {
	if len(tasks) <= 20 {
		return false
	}

	// Count sections (tasks with IDs like "1.1", "2.1" indicate multiple sections)
	sections := make(map[string]bool)
	for _, task := range tasks {
		sectionNum := extractSectionNumber(task.ID)
		if sectionNum != "0" {
			sections[sectionNum] = true
		}
	}

	// Split if > 20 tasks AND multiple sections
	return len(sections) > 1
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

// checkDependencies verifies that all required dependencies are archived.
// This is a hard check - if any required proposals are not archived,
// the accept command fails with a clear error message.
func (*AcceptCmd) checkDependencies(
	projectRoot, changeID string,
) error {
	err := validation.ValidateDependenciesForAccept(changeID, projectRoot)
	if err != nil {
		// Check if it's an UnmetDependenciesError for better formatting
		if unmetErr, ok := err.(*validation.UnmetDependenciesError); ok {
			fmt.Println("\nDependency check failed:")
			fmt.Printf(
				"  Change '%s' requires the following proposals to be archived:\n",
				unmetErr.ChangeID,
			)
			for _, dep := range unmetErr.Dependencies {
				fmt.Printf("    - %s\n", dep)
			}
			fmt.Println("\nArchive the required proposals first using 'spectr pr archive <id>'")
		}

		return err
	}

	return nil
}
