// Package taskexec provides functionality for executing and managing tasks in spectr changes.
package taskexec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/utils"
)

// filePerm is the permission mode for writing files
const filePerm = 0o644

// StatusUpdater handles updating task statuses in tasks.jsonc files
type StatusUpdater struct {
	changeDir string
}

// NewStatusUpdater creates a new StatusUpdater instance
func NewStatusUpdater(changeDir string) *StatusUpdater {
	return &StatusUpdater{
		changeDir: changeDir,
	}
}

// UpdateTaskStatus updates the status of a specific task
// Supports both v1 flat and v2 hierarchical formats
func (su *StatusUpdater) UpdateTaskStatus(taskID string, status parsers.TaskStatusValue) error {
	tasksFile := filepath.Join(su.changeDir, "tasks.jsonc")

	// Try to update in the root file first
	updated, err := su.updateTaskInFile(tasksFile, taskID, status)
	if err != nil {
		return err
	}

	if updated {
		// Task was updated in root file
		// For v2 hierarchical, check if we need to update parent status aggregation
		return su.updateParentStatusIfNeeded(tasksFile)
	}

	// Task not found in root file, search in child files (v2 hierarchical)
	return su.updateTaskInHierarchy(tasksFile, taskID, status)
}

// updateTaskInFile updates a task in a specific file
// Returns true if the task was found and updated, false otherwise
func (su *StatusUpdater) updateTaskInFile(
	filePath string,
	taskID string,
	status parsers.TaskStatusValue,
) (bool, error) {
	// Read the tasks file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return false, fmt.Errorf("failed to read tasks file %s: %w", filePath, err)
	}

	// Parse JSONC (strip comments)
	jsonData := utils.StripJSONCComments(data)

	// Parse the tasks structure
	var tasksFileData parsers.TasksFile
	if err := json.Unmarshal(jsonData, &tasksFileData); err != nil {
		return false, fmt.Errorf("failed to parse tasks file %s: %w", filePath, err)
	}

	// Find and update the task
	taskFound := false
	for i := range tasksFileData.Tasks {
		if tasksFileData.Tasks[i].ID != taskID {
			continue
		}

		tasksFileData.Tasks[i].Status = status
		taskFound = true

		break
	}

	if !taskFound {
		return false, nil
	}

	// Marshal the updated data
	updatedJSON, err := json.MarshalIndent(tasksFileData, "", "  ")
	if err != nil {
		return false, fmt.Errorf("failed to marshal tasks data: %w", err)
	}

	// Write back to the file atomically
	tmpFile := filePath + ".tmp"
	if err := os.WriteFile(tmpFile, updatedJSON, filePerm); err != nil {
		return false, fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tmpFile, filePath); err != nil {
		_ = os.Remove(tmpFile)

		return false, fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return true, nil
}

// updateTaskInHierarchy searches for a task in child files and updates it
func (su *StatusUpdater) updateTaskInHierarchy(
	rootFile string,
	taskID string,
	status parsers.TaskStatusValue,
) error {
	// Read the root tasks file to find child references
	data, err := os.ReadFile(rootFile)
	if err != nil {
		return fmt.Errorf("failed to read tasks file: %w", err)
	}

	// Parse JSONC (strip comments)
	jsonData := utils.StripJSONCComments(data)

	// Parse the tasks structure
	var tasksFileData parsers.TasksFile
	if err := json.Unmarshal(jsonData, &tasksFileData); err != nil {
		return fmt.Errorf("failed to parse tasks file: %w", err)
	}

	// Search through tasks with children
	for _, task := range tasksFileData.Tasks {
		if task.Children == "" {
			continue
		}

		// Resolve the child file path
		childPath, err := su.resolveChildPath(task.Children, filepath.Dir(rootFile))
		if err != nil {
			continue
		}

		// Try to update in the child file
		updated, err := su.updateTaskInFile(childPath, taskID, status)
		if err != nil {
			continue
		}

		if updated {
			// Task was updated in child file
			// Update parent status aggregation
			return su.updateParentStatusIfNeeded(childPath)
		}
	}

	return fmt.Errorf("task with ID %s not found in hierarchy", taskID)
}

// resolveChildPath resolves a $ref link to an absolute file path
func (su *StatusUpdater) resolveChildPath(ref string, baseDir string) (string, error) {
	// Parse the $ref format: "$ref:specs/capability/tasks.jsonc"
	if !strings.HasPrefix(ref, "$ref:") {
		return "", fmt.Errorf("invalid $ref format: %s", ref)
	}

	// Extract the path after "$ref:"
	refPath := strings.TrimPrefix(ref, "$ref:")

	// Resolve relative to the base directory
	return filepath.Join(baseDir, refPath), nil
}

// updateParentStatusIfNeeded updates parent task status based on child completion
func (su *StatusUpdater) updateParentStatusIfNeeded(filePath string) error {
	// Read the child tasks file to get the parent ID
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read tasks file %s: %w", filePath, err)
	}

	// Parse JSONC (strip comments)
	jsonData := utils.StripJSONCComments(data)

	// Parse the tasks structure
	var childFileData parsers.TasksFile
	if err := json.Unmarshal(jsonData, &childFileData); err != nil {
		return fmt.Errorf("failed to parse tasks file %s: %w", filePath, err)
	}

	// If there's no parent field, this is a root file or v1 format - nothing to aggregate
	if childFileData.Parent == "" {
		return nil
	}

	// Compute the aggregated status of all child tasks
	aggregatedStatus := su.aggregateChildStatuses(childFileData.Tasks)

	// Update the parent task in the root file with the aggregated status
	rootFile := filepath.Join(su.changeDir, "tasks.jsonc")
	return su.updateParentTask(rootFile, childFileData.Parent, aggregatedStatus)
}

// aggregateChildStatuses computes the overall status based on child task statuses
// Rules:
// - If any child is pending, parent is in_progress (work started but not done)
// - If all children are completed, parent is completed
// - If all children are pending, parent is pending
func (su *StatusUpdater) aggregateChildStatuses(tasks []parsers.Task) parsers.TaskStatusValue {
	if len(tasks) == 0 {
		return parsers.TaskStatusPending
	}

	hasCompleted := false
	hasPending := false
	hasInProgress := false

	for _, task := range tasks {
		switch task.Status {
		case parsers.TaskStatusCompleted:
			hasCompleted = true
		case parsers.TaskStatusPending:
			hasPending = true
		case parsers.TaskStatusInProgress:
			hasInProgress = true
		}
	}

	// If all completed, parent is completed
	if hasCompleted && !hasPending && !hasInProgress {
		return parsers.TaskStatusCompleted
	}

	// If any in progress or mix of pending/completed, parent is in_progress
	if hasInProgress || (hasCompleted && hasPending) {
		return parsers.TaskStatusInProgress
	}

	// All pending
	return parsers.TaskStatusPending
}

// updateParentTask updates a parent task's status in the root file
func (su *StatusUpdater) updateParentTask(
	rootFile string,
	parentID string,
	status parsers.TaskStatusValue,
) error {
	// Read the root tasks file
	data, err := os.ReadFile(rootFile)
	if err != nil {
		return fmt.Errorf("failed to read root tasks file: %w", err)
	}

	// Parse JSONC (strip comments)
	jsonData := utils.StripJSONCComments(data)

	// Parse the tasks structure
	var rootFileData parsers.TasksFile
	if err := json.Unmarshal(jsonData, &rootFileData); err != nil {
		return fmt.Errorf("failed to parse root tasks file: %w", err)
	}

	// Find and update the parent task
	parentFound := false
	for i := range rootFileData.Tasks {
		if rootFileData.Tasks[i].ID != parentID {
			continue
		}

		rootFileData.Tasks[i].Status = status
		parentFound = true

		break
	}

	if !parentFound {
		// Parent task not found, but this is not a critical error
		// The parent might have been removed or the structure changed
		return nil
	}

	// Marshal the updated data
	updatedJSON, err := json.MarshalIndent(rootFileData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal root tasks data: %w", err)
	}

	// Write back to the file atomically
	tmpFile := rootFile + ".tmp"
	if err := os.WriteFile(tmpFile, updatedJSON, filePerm); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tmpFile, rootFile); err != nil {
		_ = os.Remove(tmpFile)

		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
}
