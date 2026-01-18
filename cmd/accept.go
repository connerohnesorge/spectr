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
)

// filePerm is the standard file permission for created files (rw-r--r--)
const filePerm = 0o644

// splitThreshold is the line count threshold for splitting tasks.jsonc files.
// Files with more than this many lines will be split into multiple files.
const splitThreshold = 100

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
	if cfg != nil && cfg.AppendTasks != nil && cfg.AppendTasks.HasTasks() {
		appendCfg = cfg.AppendTasks
	}

	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.jsonc",
	)

	if c.DryRun {
		totalTasks := len(tasks)
		if appendCfg != nil {
			totalTasks += len(appendCfg.Tasks)
		}
		fmt.Printf(
			"Would convert: %s\nwould write to: %s\nWould preserve: %s\nFound %d tasks\n", //nolint:lll,revive // Long format string for dry-run output
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			totalTasks,
		)

		return nil
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
//
// During regeneration (re-running `spectr accept`), this function loads
// existing task statuses and preserves them for tasks whose IDs match.
//
// If the tasks.md file exceeds the split threshold (100 lines), this function
// will automatically split it into multiple tasks.jsonc files for better
// agent readability.
func writeAndCleanup(
	tasksMdPath, tasksJSONPath string,
	tasks []parsers.Task,
	appendCfg *config.AppendTasksConfig,
) error {
	// Load existing statuses from the change directory
	changeDir := filepath.Dir(tasksJSONPath)
	statusMap, err := loadExistingStatuses(changeDir)
	if err != nil {
		return fmt.Errorf(
			"failed to load existing statuses: %w",
			err,
		)
	}

	// Check if we should split the file
	split, err := shouldSplit(tasksMdPath)
	if err != nil {
		return fmt.Errorf(
			"failed to check split threshold: %w",
			err,
		)
	}

	totalTasks := len(tasks)
	if appendCfg != nil {
		totalTasks += len(appendCfg.Tasks)
	}

	if split {
		// Route to hierarchical writing (version 2)
		sections, err := parseSections(tasksMdPath)
		if err != nil {
			return fmt.Errorf(
				"failed to parse sections: %w",
				err,
			)
		}

		// Extract change ID from path for child file headers
		changeID := filepath.Base(changeDir)

		if err := writeHierarchicalTasks(changeDir, changeID, sections, statusMap); err != nil {
			return fmt.Errorf(
				"failed to write hierarchical tasks: %w",
				err,
			)
		}

		fmt.Printf(
			"Converted %s -> %s (split into multiple files)\nPreserved %s\nWrote %d tasks\n",
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			totalTasks,
		)
	} else {
		// Route to flat file writing (version 1)
		if err := writeTasksJSONC(tasksJSONPath, tasks, appendCfg, statusMap); err != nil {
			return fmt.Errorf(
				"failed to write tasks.jsonc: %w",
				err,
			)
		}

		fmt.Printf(
			"Converted %s -> %s\nPreserved %s\nWrote %d tasks\n",
			tasksMdPath,
			tasksJSONPath,
			tasksMdPath,
			totalTasks,
		)
	}

	// Preserve tasks.md to avoid information loss (formatting, comments, links)
	// Both tasks.md and tasks.jsonc now coexist after conversion

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

// countLines counts the number of lines in a file.
// Returns the line count or an error if the file cannot be read.
func countLines(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	lineCount := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return 0, fmt.Errorf("error reading file: %w", err)
	}

	return lineCount, nil
}

// shouldSplit determines whether a tasks.md file should be split into
// multiple tasks.jsonc files based on the line count threshold.
// Returns true if the file exceeds splitThreshold (100 lines).
func shouldSplit(path string) (bool, error) {
	lineCount, err := countLines(path)
	if err != nil {
		return false, err
	}

	return lineCount > splitThreshold, nil
}

// Section represents a section in tasks.md with its name, tasks, and line range.
// Sections are identified by headers like "## 1. Section Name" or "## Section Name".
// The line range tracks where the section starts and ends in the file for size calculations.
type Section struct {
	Name      string         // Section name (e.g., "Implementation", "Testing")
	Number    string         // Section number (e.g., "1", "2") - empty for unnumbered sections
	Tasks     []parsers.Task // Tasks belonging to this section
	StartLine int            // Line number where section starts (1-indexed)
	EndLine   int            // Line number where section ends (1-indexed)
}

// LineCount returns the number of lines in this section.
func (s *Section) LineCount() int {
	if s.EndLine < s.StartLine {
		return 0
	}

	return s.EndLine - s.StartLine + 1
}

// parseSections extracts sections from tasks.md file.
// It identifies sections by "## N. Section Name" or "## Section Name" headers,
// tracks their line ranges, and returns a slice of Section structs.
// Tasks are extracted and associated with their containing section.
func parseSections(path string) ([]Section, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer func() { _ = file.Close() }()

	var sections []Section
	var currentSection *Section
	state := &taskParseState{}
	lineNum := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Check for section header
		if name, number, ok := markdown.MatchAnySection(line); ok {
			// Close previous section
			if currentSection != nil {
				currentSection.EndLine = lineNum - 1
				sections = append(sections, *currentSection)
			}

			// Start new section
			state.handleSection(name, number)
			currentSection = &Section{
				Name:      state.sectionName,
				Number:    state.sectionNum,
				Tasks:     make([]parsers.Task, 0),
				StartLine: lineNum,
			}

			continue
		}

		// Check for task line
		if match, ok := markdown.MatchFlexibleTask(line); ok {
			task := state.createTask(match)

			// Add task to current section if we have one
			if currentSection != nil {
				currentSection.Tasks = append(currentSection.Tasks, task)
			}
		}
	}

	// Close final section
	if currentSection != nil {
		currentSection.EndLine = lineNum
		sections = append(sections, *currentSection)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return sections, nil
}

// extractTasksForSection returns all tasks that belong to the specified section.
// It matches tasks by their section name and returns them in the order they appear.
func extractTasksForSection(allTasks []parsers.Task, sectionName string) []parsers.Task {
	var tasks []parsers.Task
	for _, task := range allTasks {
		if task.Section == sectionName {
			tasks = append(tasks, task)
		}
	}

	return tasks
}

// shouldSplitSection determines whether a section should be split into subsections
// based on its line count. Returns true if the section exceeds splitThreshold (100 lines).
func shouldSplitSection(section *Section) bool {
	return section.LineCount() > splitThreshold
}

// SubsectionGroup represents a group of tasks that share a common ID prefix.
// For example, tasks "1.1", "1.2", "1.3" share the prefix "1" and would be
// grouped together. This allows large sections to be split into smaller,
// more manageable chunks while preserving logical groupings.
type SubsectionGroup struct {
	Prefix    string         // ID prefix shared by all tasks (e.g., "1", "2.1")
	Tasks     []parsers.Task // Tasks with this prefix
	StartLine int            // Line number where subsection starts (1-indexed)
	EndLine   int            // Line number where subsection ends (1-indexed)
}

// LineCount returns the number of lines in this subsection group.
func (sg *SubsectionGroup) LineCount() int {
	if sg.EndLine < sg.StartLine {
		return 0
	}

	return sg.EndLine - sg.StartLine + 1
}

// parseSubsections groups tasks by their ID prefix to create subsections.
// For example, tasks "1.1", "1.2", "1.3" all share prefix "1" and will be
// grouped together. This is used when a section is too large and needs to
// be split into smaller chunks.
//
// The function extracts the prefix from each task's ID (everything before the
// last dot) and groups tasks with the same prefix together.
//
// Example:
//
//	Tasks: ["1.1", "1.2", "2.1", "2.2", "2.3"]
//	Groups: [["1.1", "1.2"], ["2.1", "2.2", "2.3"]]
func parseSubsections(tasks []parsers.Task) []SubsectionGroup {
	// Map to track subsection groups by prefix
	groupMap := make(map[string]*SubsectionGroup)
	var orderedPrefixes []string

	for _, task := range tasks {
		// Extract prefix from task ID (everything before the last dot)
		prefix := extractIDPrefix(task.ID)

		// Create new group if it doesn't exist
		if _, exists := groupMap[prefix]; !exists {
			groupMap[prefix] = &SubsectionGroup{
				Prefix: prefix,
				Tasks:  make([]parsers.Task, 0),
			}
			orderedPrefixes = append(orderedPrefixes, prefix)
		}

		// Add task to the group
		groupMap[prefix].Tasks = append(groupMap[prefix].Tasks, task)
	}

	// Convert map to ordered slice
	groups := make([]SubsectionGroup, 0, len(orderedPrefixes))
	for _, prefix := range orderedPrefixes {
		groups = append(groups, *groupMap[prefix])
	}

	return groups
}

// extractIDPrefix returns the prefix of a task ID.
// For hierarchical IDs like "1.2.3", it returns everything before the last dot ("1.2").
// For simple IDs like "1", it returns the ID itself.
// For IDs with one dot like "1.1", it returns the part before the dot ("1").
func extractIDPrefix(id string) string {
	// Find the last dot in the ID
	lastDot := strings.LastIndex(id, ".")
	if lastDot == -1 {
		// No dot found - this is a simple ID like "1"
		return id
	}

	// Return everything before the last dot
	return id[:lastDot]
}

// assignHierarchicalID generates a hierarchical ID for a child task based on the parent ID.
// It uses dot notation to create hierarchical relationships:
//   - Root tasks: "1", "2", "5" (no parent)
//   - Child tasks: "5.1", "5.2" (parent "5")
//   - Nested children: "5.1.1", "5.1.2" (parent "5.1")
//
// The function preserves existing IDs from tasks.md when the task already has an ID.
// For new tasks or when generating IDs for split files, it prepends the parent ID.
//
// Parameters:
//   - existingID: The task's current ID from tasks.md (may be empty)
//   - parentID: The parent task ID (empty string for root tasks)
//   - childIndex: The sequential index of this child within the parent (1-based)
//
// Returns the hierarchical ID (e.g., "5.1", "5.2", etc.)
func assignHierarchicalID(existingID, parentID string, childIndex int) string {
	// If no parent, this is a root task - use existing ID or generate simple ID
	if parentID == "" {
		if existingID != "" {
			return existingID
		}

		return strconv.Itoa(childIndex)
	}

	// If existing ID already has the parent prefix, preserve it
	// For example, if parentID is "5" and existingID is "5.1", keep "5.1"
	if existingID != "" && strings.HasPrefix(existingID, parentID+".") {
		return existingID
	}

	// Generate hierarchical ID: parent.child
	// For example, parent "5" + child 1 = "5.1"
	return fmt.Sprintf("%s.%d", parentID, childIndex)
}

// validateIDUniqueness checks that all task IDs are unique within a set of tasks.
// Duplicate IDs can cause confusion and errors when tracking task status.
//
// Returns an error if any duplicate IDs are found, listing all duplicates.
// Returns nil if all IDs are unique.
func validateIDUniqueness(tasks []parsers.Task) error {
	seen := make(map[string]bool)
	duplicates := make([]string, 0)

	for _, task := range tasks {
		if seen[task.ID] {
			// Only add to duplicates list once
			if !stringSliceContains(duplicates, task.ID) {
				duplicates = append(duplicates, task.ID)
			}
		} else {
			seen[task.ID] = true
		}
	}

	if len(duplicates) > 0 {
		return fmt.Errorf("duplicate task IDs found: %v", duplicates)
	}

	return nil
}

// stringSliceContains checks if a string slice contains a specific string.
func stringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}

	return false
}
