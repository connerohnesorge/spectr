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
	"unicode"

	"github.com/connerohnesorge/spectr/internal/archive"
	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// matchSectionToCapability normalizes a section name to kebab-case for
// matching with delta spec directory names. It strips leading numbers,
// periods, and whitespace, then converts to kebab-case.
func matchSectionToCapability(sectionName string) string {
	// Strip leading numbers and punctuation (e.g., "5. Support Aider" -> "Support Aider")
	sectionName = strings.TrimLeftFunc(sectionName, func(r rune) bool {
		return unicode.IsDigit(r) || r == '.' || unicode.IsSpace(r)
	})

	// Convert to kebab-case
	var result strings.Builder
	for i, r := range sectionName {
		if r == '-' {
			// Preserve existing dashes
			if i > 0 && result.Len() > 0 && result.String()[result.Len()-1] != '-' {
				result.WriteRune('-')
			}
		} else if unicode.IsSpace(r) || r == '_' {
			if i > 0 && result.Len() > 0 && result.String()[result.Len()-1] != '-' {
				result.WriteRune('-')
			}
		} else if unicode.IsUpper(r) {
			if i > 0 && result.Len() > 0 && result.String()[result.Len()-1] != '-' {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result.WriteRune(r)
		}
	}

	return strings.ToLower(result.String())
}

// findMatchingDeltaSpec checks if a delta spec directory exists for the given
// capability name within the change directory and contains a spec.md file.
func findMatchingDeltaSpec(changeDir, capability string) bool {
	deltaSpecDir := filepath.Join(changeDir, "specs", capability)
	specPath := filepath.Join(deltaSpecDir, "spec.md")

	// Check if directory exists AND contains a spec.md file
	info, err := os.Stat(deltaSpecDir)
	if err != nil || !info.IsDir() {
		return false
	}

	_, err = os.Stat(specPath)
	return err == nil
}

// SectionTasks holds tasks for a specific section, used for hierarchical splitting.
type SectionTasks struct {
	SectionName string
	Tasks       []parsers.Task
	Capability  string // matched delta spec capability (empty if no match)
	HasChildren bool   // true if this section should have children ref
}

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
	DryRun bool `                                        help:"Preview without writing" name:"dry-run"` //nolint:lll,revive // Kong struct tag with alignment padding

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

	// Validate the change before conversion
	if err = c.runValidation(changeDir); err != nil {
		return fmt.Errorf(
			"validation failed: %w",
			err,
		)
	}

	// Check if any delta specs exist for hierarchical structure
	hasDeltaSpecs := c.hasDeltaSpecs(changeDir)

	var tasks []parsers.Task
	var childFiles map[string][]parsers.Task
	var hasHierarchy bool

	if hasDeltaSpecs {
		// Use hierarchical parsing with section tracking
		sections, err := parseTasksMdWithSections(tasksMdPath)
		if err != nil {
			return fmt.Errorf(
				"failed to parse tasks.md: %w",
				err,
			)
		}

		tasks, childFiles, hasHierarchy = splitTasksByCapability(changeDir, sections)

		// Safety check: if tasks.md has content but no valid tasks were found
		totalTasks := len(tasks)
		for _, childTasks := range childFiles {
			totalTasks += len(childTasks)
		}
		if totalTasks == 0 {
			info, statErr := os.Stat(tasksMdPath)
			if statErr == nil && info.Size() > 0 {
				return &specterrs.NoValidTasksError{
					TasksMdPath: tasksMdPath,
					FileSize:    info.Size(),
				}
			}
		}
	} else {
		// Use simple parsing (backwards compatible)
		tasks, err = parseTasksMd(tasksMdPath)
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
	}

	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.jsonc",
	)

	if c.DryRun {
		fmt.Printf(
			"Would convert: %s\nwould write to: %s\nWould preserve: %s\nFound %d tasks\n", //nolint:lll,revive // Long format string for dry-run output
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			len(tasks),
		)
		if hasHierarchy {
			fmt.Printf("Would create %d child task file(s)\n", len(childFiles))
			for cap, childTasks := range childFiles {
				fmt.Printf("  - specs/%s/tasks.jsonc: %d tasks\n", cap, len(childTasks))
			}
		}

		return nil
	}

	return c.writeHierarchicalTasks(
		tasksMdPath,
		tasksJSONPath,
		tasks,
		childFiles,
		hasHierarchy,
	)
}

// hasDeltaSpecs checks if the change directory has any delta spec directories.
func (c *AcceptCmd) hasDeltaSpecs(changeDir string) bool {
	specsDir := filepath.Join(changeDir, "specs")
	info, err := os.Stat(specsDir)
	if err != nil || !info.IsDir() {
		return false
	}

	entries, err := os.ReadDir(specsDir)
	if err != nil {
		return false
	}

	for _, entry := range entries {
		if entry.IsDir() {
			specPath := filepath.Join(specsDir, entry.Name(), "spec.md")
			if _, err := os.Stat(specPath); err == nil {
				return true
			}
		}
	}

	return false
}

// writeHierarchicalTasks writes root and child task files for hierarchical structure.
func (c *AcceptCmd) writeHierarchicalTasks(
	tasksMdPath, tasksJSONPath string,
	rootTasks []parsers.Task,
	childFiles map[string][]parsers.Task,
	hasHierarchy bool,
) error {
	if hasHierarchy {
		// Write root tasks.jsonc with hierarchical format
		includes := []string{"specs/*/tasks.jsonc"}
		if err := writeTasksJSONCHierarchical(tasksJSONPath, rootTasks, includes); err != nil {
			return fmt.Errorf("failed to write tasks.jsonc: %w", err)
		}

		// Write child task files
		for capability, childTasks := range childFiles {
			if len(childTasks) == 0 {
				continue
			}
			parentID := childTasks[0].ID
			if err := writeChildTasksJSONC(
				filepath.Dir(tasksJSONPath),
				capability,
				parentID,
				childTasks,
			); err != nil {
				return fmt.Errorf("failed to write child tasks.jsonc for %s: %w", capability, err)
			}
		}

		fmt.Printf(
			"Converted %s -> %s (hierarchical)\nPreserved %s\nWrote %d root tasks, %d child files\n",
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			len(rootTasks),
			len(childFiles),
		)
	} else {
		// Use simple write for backwards compatibility
		if err := writeTasksJSONC(tasksJSONPath, rootTasks); err != nil {
			return fmt.Errorf("failed to write tasks.jsonc: %w", err)
		}

		fmt.Printf(
			"Converted %s -> %s\nPreserved %s\nWrote %d tasks\n",
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			len(rootTasks),
		)
	}

	return nil
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
) error {
	if err := writeTasksJSONC(tasksJSONPath, tasks); err != nil {
		return fmt.Errorf(
			"failed to write tasks.jsonc: %w",
			err,
		)
	}

	// Preserve tasks.md to avoid information loss (formatting, comments, links)
	// Both tasks.md and tasks.jsonc now coexist after conversion

	fmt.Printf(
		"Converted %s -> %s\nPreserved %s\nWrote %d tasks\n",
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

// parseTasksMdWithSections parses tasks.md and returns tasks grouped by section.
// This is used for hierarchical task file splitting.
func parseTasksMdWithSections(
	path string,
) ([]SectionTasks, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var sections []SectionTasks
	state := &taskParseState{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check for section header (numbered or unnumbered)
		if name, number, ok := markdown.MatchAnySection(line); ok {
			state.handleSection(name, number)

			// Add new section to list
			sections = append(sections, SectionTasks{
				SectionName: state.sectionName,
			})
			continue
		}

		// Skip if no current section
		if len(sections) == 0 {
			continue
		}

		// Check for task line using flexible matching
		match, ok := markdown.MatchFlexibleTask(line)
		if !ok {
			continue
		}

		// Add task to current section
		task := state.createTask(match)
		lastIdx := len(sections) - 1
		sections[lastIdx].Tasks = append(sections[lastIdx].Tasks, task)
	}

	err = scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return sections, nil
}

// splitTasksByCapability partitions tasks into root and child files based on
// whether the section matches a delta spec directory.
func splitTasksByCapability(
	changeDir string,
	sections []SectionTasks,
) (rootTasks []parsers.Task, childFiles map[string][]parsers.Task, hasHierarchy bool) {
	childFiles = make(map[string][]parsers.Task)

	for _, section := range sections {
		if len(section.Tasks) == 0 {
			continue
		}

		capability := matchSectionToCapability(section.SectionName)

		// Check if this section matches a delta spec
		if findMatchingDeltaSpec(changeDir, capability) {
			hasHierarchy = true
			section.HasChildren = true

			// Add reference task to root
			refTask := parsers.Task{
				ID:          section.Tasks[0].ID, // Use first task's ID as parent ID
				Section:     section.SectionName,
				Description: section.Tasks[0].Description,
				Status:      section.Tasks[0].Status,
				Children:    fmt.Sprintf("$ref:specs/%s/tasks.jsonc", capability),
			}
			rootTasks = append(rootTasks, refTask)

			// Store child tasks (without section, keeping only ID and description)
			var childTasks []parsers.Task
			for _, task := range section.Tasks {
				childTasks = append(childTasks, parsers.Task{
					ID:          task.ID,
					Section:     "", // Child tasks don't need section
					Description: task.Description,
					Status:      task.Status,
					Children:    "", // Child tasks don't have children
				})
			}
			childFiles[capability] = childTasks
		} else {
			// No match - tasks stay in root
			section.HasChildren = false
			rootTasks = append(rootTasks, section.Tasks...)
		}
	}

	return rootTasks, childFiles, hasHierarchy
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
