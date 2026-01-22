package discovery

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parsers"
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
	jsonData := stripJSONCComments(data)

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

	return nil, fmt.Errorf("no pending tasks found")
}

// stripJSONCComments removes single-line comments from JSONC content
func stripJSONCComments(data []byte) []byte {
	lines := strings.Split(string(data), "\n")
	var cleaned []string

	for _, line := range lines {
		// Skip comment lines (lines starting with // after whitespace)
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") {
			continue
		}
		// Remove inline comments
		if idx := strings.Index(line, "//"); idx != -1 {
			// Check if the // is inside quotes
			quoteCount := 0
			for i := 0; i < idx; i++ {
				if line[i] == '"' && (i == 0 || line[i-1] != '\\') {
					quoteCount++
				}
			}
			// Only strip comment if we're not inside quotes
			if quoteCount%2 == 0 {
				line = strings.TrimRight(line[:idx], " \t")
			}
		}
		// Only add non-empty lines
		if strings.TrimSpace(line) != "" {
			cleaned = append(cleaned, line)
		}
	}

	return []byte(strings.Join(cleaned, "\n"))
}
