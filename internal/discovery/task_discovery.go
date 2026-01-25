package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/utils"
)

// TaskDiscovery handles finding and managing tasks in spectr changes
type TaskDiscovery struct {
	changeDir string
	// Track visited files to detect circular references
	visited map[string]bool
}

// NewTaskDiscovery creates a new TaskDiscovery instance
func NewTaskDiscovery(changeDir string) *TaskDiscovery {
	return &TaskDiscovery{
		changeDir: changeDir,
		visited:   make(map[string]bool),
	}
}

// FindNextPendingTask finds the first pending task in the tasks.jsonc file
// Supports both v1 flat and v2 hierarchical formats with $ref resolution
func (td *TaskDiscovery) FindNextPendingTask() (*parsers.Task, error) {
	tasksFile := filepath.Join(td.changeDir, "tasks.jsonc")
	return td.findNextPendingTaskInFile(tasksFile)
}

// findNextPendingTaskInFile finds the first pending task in a specific file
func (td *TaskDiscovery) findNextPendingTaskInFile(filePath string) (*parsers.Task, error) {
	// Check for circular references
	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	if td.visited[absPath] {
		return nil, fmt.Errorf("circular reference detected: %s", absPath)
	}
	td.visited[absPath] = true

	// Read the tasks file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file %s: %w", filePath, err)
	}

	// Parse JSONC (strip comments)
	jsonData := utils.StripJSONCComments(data)

	// Parse the tasks structure
	var tasksFileData parsers.TasksFile
	if err := json.Unmarshal(jsonData, &tasksFileData); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file %s: %w", filePath, err)
	}

	// Find the first pending task, recursively checking children
	for _, task := range tasksFileData.Tasks {
		// If this task has children, check the child file first
		if task.Children != "" {
			childTask, err := td.resolveChildRef(task.Children, filepath.Dir(filePath))
			if err != nil {
				// If child resolution fails, log but continue to next task
				// This allows graceful degradation
				continue
			}
			if childTask != nil {
				return childTask, nil
			}
			// All children completed, continue to next task
			continue
		}

		// No children, check this task's status
		if task.Status == parsers.TaskStatusPending {
			return &task, nil
		}
	}

	return nil, errors.New("no pending tasks found")
}

// resolveChildRef resolves a $ref link to a child task file and finds the first pending task
func (td *TaskDiscovery) resolveChildRef(ref string, baseDir string) (*parsers.Task, error) {
	// Parse the $ref format: "$ref:specs/capability/tasks.jsonc"
	if !strings.HasPrefix(ref, "$ref:") {
		return nil, fmt.Errorf("invalid $ref format: %s (must start with $ref:)", ref)
	}

	// Extract the path after "$ref:"
	refPath := strings.TrimPrefix(ref, "$ref:")

	// Resolve relative to the base directory
	childPath := filepath.Join(baseDir, refPath)

	// Check if the child file exists
	if _, err := os.Stat(childPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("referenced file not found: %s", childPath)
	}

	// Recursively find the first pending task in the child file
	return td.findNextPendingTaskInFile(childPath)
}
