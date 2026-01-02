// Package cmd provides command-line interface implementations.
// This file implements the tasks command for viewing task status.
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/connerohnesorge/spectr/internal/discovery"
	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/specterrs"
)

// TasksCmd displays task status for a change, supporting hierarchical task files.
type TasksCmd struct {
	// ChangeID is the optional change identifier to display tasks for
	ChangeID string `arg:"" optional:"" predictor:"changeID" help:"Show tasks for a specific change"`
	// Flatten merges all tasks from root and child files into a single view
	Flatten bool `help:"Flatten all tasks into single view" name:"flatten"`
	// JSON outputs task data as JSON
	JSON bool `help:"Output as JSON" name:"json"`
	// NoInteractive disables interactive prompts
	NoInteractive bool `help:"Disable prompts" name:"no-interactive"`
}

// Run executes the tasks command.
// It resolves the change ID, reads the tasks.jsonc file(s),
// and displays the task status summary or flattened list.
func (c *TasksCmd) Run() error {
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

	return c.displayTasks(projectRoot, changeID)
}

// resolveChangeID resolves the change ID from argument or interactive selection.
func (c *TasksCmd) resolveChangeID(projectRoot string) (string, error) {
	if c.ChangeID != "" {
		normalizedID, _ := discovery.NormalizeItemPath(c.ChangeID)

		result, err := discovery.ResolveChangeID(normalizedID, projectRoot)
		if err != nil {
			return "", err
		}

		if result.PartialMatch {
			fmt.Printf("Resolved '%s' -> '%s'\n\n", c.ChangeID, result.ChangeID)
		}

		return result.ChangeID, nil
	}

	if c.NoInteractive {
		return "", &specterrs.MissingChangeIDError{}
	}

	return selectChangeInteractive(projectRoot)
}

// displayTasks reads and displays task information for the specified change.
func (c *TasksCmd) displayTasks(projectRoot, changeID string) error {
	changeDir := filepath.Join(projectRoot, "spectr", "changes", changeID)
	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")

	// Check if tasks.jsonc exists
	if _, err := os.Stat(tasksJSONPath); os.IsNotExist(err) {
		return fmt.Errorf("tasks.jsonc not found for change '%s'. Run 'spectr accept %s' first", changeID, changeID)
	}

	// Read root tasks file
	rootFile, err := readTasksFile(tasksJSONPath)
	if err != nil {
		return fmt.Errorf("failed to read tasks.jsonc: %w", err)
	}

	// Resolve child files if hierarchical
	var allTasks []parsers.Task
	var summary *parsers.Summary

	if rootFile.Version >= 2 && (len(rootFile.Includes) > 0 || hasChildrenReferences(rootFile.Tasks)) {
		// Resolve hierarchical structure
		allTasks, summary, err = c.resolveChildFiles(changeDir, rootFile)
		if err != nil {
			return fmt.Errorf("failed to resolve child files: %w", err)
		}
	} else {
		// Flat structure - use tasks directly from root file
		allTasks = rootFile.Tasks
		summary = computeTaskSummary(allTasks)
	}

	if c.JSON {
		return c.outputJSON(allTasks, summary, rootFile.Version >= 2)
	}

	if c.Flatten {
		return c.displayFlattened(allTasks, summary)
	}

	return c.displaySummary(allTasks, summary)
}

// readTasksFile reads and parses a tasks.jsonc file.
func readTasksFile(path string) (*parsers.TasksFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Strip JSONC comments before parsing
	cleanData := parsers.StripJSONComments(data)

	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(cleanData, &tasksFile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &tasksFile, nil
}

// hasChildrenReferences checks if any tasks have children references.
func hasChildrenReferences(tasks []parsers.Task) bool {
	for _, task := range tasks {
		if task.Children != "" {
			return true
		}
	}
	return false
}

// resolveChildFiles resolves child task files and merges all tasks.
func (c *TasksCmd) resolveChildFiles(changeDir string, rootFile *parsers.TasksFile) ([]parsers.Task, *parsers.Summary, error) {
	allTasks := make([]parsers.Task, 0, len(rootFile.Tasks))
	processedFiles := make(map[string]bool)

	// First, add root tasks and track which child files to process
	for _, task := range rootFile.Tasks {
		allTasks = append(allTasks, task)

		if task.Children != "" {
			childPath := extractChildPath(task.Children, changeDir)
			if childPath != "" && !processedFiles[childPath] {
				childTasks, err := c.loadChildFile(childPath, task.ID, processedFiles)
				if err != nil {
					return nil, nil, err
				}
				allTasks = append(allTasks, childTasks...)
				processedFiles[childPath] = true
			}
		}
	}

	// Process includes patterns
	for _, pattern := range rootFile.Includes {
		matches := findMatchingFiles(changeDir, pattern)
		for _, match := range matches {
			if processedFiles[match] {
				continue
			}
			childTasks, err := c.loadChildFile(match, "", processedFiles)
			if err != nil {
				return nil, nil, err
			}
			allTasks = append(allTasks, childTasks...)
			processedFiles[match] = true
		}
	}

	summary := computeTaskSummary(allTasks)
	return allTasks, summary, nil
}

// extractChildPath extracts the file path from a children reference.
func extractChildPath(childrenRef, changeDir string) string {
	// Handle $ref:specs/<capability>/tasks.jsonc format
	if strings.HasPrefix(childrenRef, "$ref:") {
		refPath := strings.TrimPrefix(childrenRef, "$ref:")
		// Make relative to changeDir if needed
		if !filepath.IsAbs(refPath) {
			if strings.HasPrefix(refPath, "specs/") {
				return filepath.Join(changeDir, refPath)
			}
			return filepath.Join(changeDir, refPath)
		}
		return refPath
	}
	return ""
}

// findMatchingFiles finds files matching a glob pattern.
func findMatchingFiles(changeDir, pattern string) []string {
	// Handle specs/*/tasks.jsonc pattern
	if strings.HasPrefix(pattern, "specs/") && strings.HasSuffix(pattern, "tasks.jsonc") {
		specsDir := filepath.Join(changeDir, "specs")
		entries, err := os.ReadDir(specsDir)
		if err != nil {
			return nil
		}

		var matches []string
		for _, entry := range entries {
			if entry.IsDir() {
				tasksPath := filepath.Join(specsDir, entry.Name(), "tasks.jsonc")
				if _, err := os.Stat(tasksPath); err == nil {
					matches = append(matches, tasksPath)
				}
			}
		}
		return matches
	}
	return nil
}

// loadChildFile loads a child task file and returns its tasks.
func (c *TasksCmd) loadChildFile(path, parentID string, processedFiles map[string]bool) ([]parsers.Task, error) {
	if processedFiles[path] {
		return nil, nil
	}

	childFile, err := readTasksFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read child file %s: %w", path, err)
	}

	var tasks []parsers.Task
	for _, task := range childFile.Tasks {
		// If parentID is set, prefix task IDs
		if parentID != "" && !strings.HasPrefix(task.ID, parentID+".") {
			task.ID = parentID + "." + task.ID
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// computeTaskSummary computes summary counts from a list of tasks.
func computeTaskSummary(tasks []parsers.Task) *parsers.Summary {
	summary := &parsers.Summary{
		Total: len(tasks),
	}

	for _, task := range tasks {
		switch task.Status {
		case parsers.TaskStatusCompleted:
			summary.Completed++
		case parsers.TaskStatusInProgress:
			summary.InProgress++
		case parsers.TaskStatusPending:
			summary.Pending++
		}
	}

	return summary
}

// displaySummary shows section-by-section progress summary.
func (c *TasksCmd) displaySummary(tasks []parsers.Task, summary *parsers.Summary) error {
	// Group tasks by section
	sectionTasks := make(map[string][]parsers.Task)
	for _, task := range tasks {
		section := task.Section
		if section == "" {
			section = "Uncategorized"
		}
		sectionTasks[section] = append(sectionTasks[section], task)
	}

	// Sort sections
	var sections []string
	for section := range sectionTasks {
		sections = append(sections, section)
	}
	sort.Strings(sections)

	// Display header
	fmt.Println("Tasks")

	// Display per-section progress
	for _, section := range sections {
		tasks := sectionTasks[section]
		completed := countStatus(tasks, parsers.TaskStatusCompleted)
		inProgress := countStatus(tasks, parsers.TaskStatusInProgress)
		total := len(tasks)

		percent := float64(completed) / float64(total) * 100
		fmt.Printf("%s: %d/%d completed (%.0f%%)", section, completed, total, percent)
		if inProgress > 0 {
			fmt.Printf(", %d in progress", inProgress)
		}
		fmt.Println()
	}

	// Display total progress
	if summary != nil {
		if summary.Total > 0 {
			percent := float64(summary.Completed) / float64(summary.Total) * 100
			fmt.Printf("\nTotal: %d/%d completed (%.0f%%)\n", summary.Completed, summary.Total, percent)
		} else {
			fmt.Println("\nNo tasks found")
		}
	}

	return nil
}

// displayFlattened shows all tasks in a flattened list.
func (c *TasksCmd) displayFlattened(tasks []parsers.Task, summary *parsers.Summary) error {
	// Sort by ID
	sortedTasks := make([]parsers.Task, len(tasks))
	copy(sortedTasks, tasks)
	sort.Slice(sortedTasks, func(i, j int) bool {
		return sortedTasks[i].ID < sortedTasks[j].ID
	})

	fmt.Println("All Tasks:")
	fmt.Println(strings.Repeat("-", 60))

	for _, task := range sortedTasks {
		statusIcon := statusToIcon(task.Status)
		section := task.Section
		if section == "" {
			section = "-"
		}
		fmt.Printf("[%s] %s | %s | %s\n", statusIcon, task.ID, section, task.Description)
	}

	fmt.Println(strings.Repeat("-", 60))

	if summary != nil {
		fmt.Printf("Total: %d tasks (%d completed, %d in progress, %d pending)\n",
			summary.Total, summary.Completed, summary.InProgress, summary.Pending)
	}

	return nil
}

// outputJSON outputs task data as JSON.
func (c *TasksCmd) outputJSON(tasks []parsers.Task, summary *parsers.Summary, isHierarchical bool) error {
	type output struct {
		Hierarchical bool             `json:"hierarchical"`
		Summary      *parsers.Summary `json:"summary,omitempty"`
		Tasks        []parsers.Task   `json:"tasks"`
	}

	out := output{
		Hierarchical: isHierarchical,
		Summary:      summary,
		Tasks:        tasks,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(out)
}

// countStatus counts tasks with a specific status.
func countStatus(tasks []parsers.Task, status parsers.TaskStatusValue) int {
	count := 0
	for _, task := range tasks {
		if task.Status == status {
			count++
		}
	}
	return count
}

// statusToIcon converts a status to a display icon.
func statusToIcon(status parsers.TaskStatusValue) string {
	switch status {
	case parsers.TaskStatusCompleted:
		return "✓"
	case parsers.TaskStatusInProgress:
		return "▶"
	case parsers.TaskStatusPending:
		return "○"
	default:
		return "?"
	}
}
