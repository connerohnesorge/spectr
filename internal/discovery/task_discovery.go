package discovery

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/connerohnesorge/spectr/internal/parsers"
	"github.com/connerohnesorge/spectr/internal/utils"
)

// TaskDiscovery handles finding and managing tasks in spectr changes
type TaskDiscovery struct {
	changeDir string
}

// NewTaskDiscovery creates a new TaskDiscovery instance
func NewTaskDiscovery(changeDir string) *TaskDiscovery {
	return &TaskDiscovery{
		changeDir: changeDir,
	}
}

// FindNextPendingTask finds the first pending task in the tasks.jsonc file
func (td *TaskDiscovery) FindNextPendingTask() (*parsers.Task, error) {
	tasksFile := filepath.Join(td.changeDir, "tasks.jsonc")

	// Read the tasks file
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %w", err)
	}

	// Parse JSONC (strip comments)
	jsonData := utils.StripJSONCComments(data)

	// Parse the tasks structure
	var tasksFileData parsers.TasksFile
	if err := json.Unmarshal(jsonData, &tasksFileData); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %w", err)
	}

	// Find the first pending task
	for _, task := range tasksFileData.Tasks {
		if task.Status == parsers.TaskStatusPending {
			return &task, nil
		}
	}

	return nil, errors.New("no pending tasks found")
}
