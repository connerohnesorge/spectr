package taskexec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

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
func (su *StatusUpdater) UpdateTaskStatus(taskID string, status parsers.TaskStatusValue) error {
	tasksFile := filepath.Join(su.changeDir, "tasks.jsonc")

	// Read the tasks file
	data, err := os.ReadFile(tasksFile)
	if err != nil {
		return fmt.Errorf("failed to read tasks file: %w", err)
	}

	// Parse JSONC (strip comments)
	jsonData := stripJSONCComments(data)

	// Parse the tasks structure
	var tasksFileData parsers.TasksFile
	if err := json.Unmarshal(jsonData, &tasksFileData); err != nil {
		return fmt.Errorf("failed to parse tasks file: %w", err)
	}

	// Find and update the task
	taskFound := false
	for i := range tasksFileData.Tasks {
		if tasksFileData.Tasks[i].ID == taskID {
			tasksFileData.Tasks[i].Status = status
			taskFound = true
			break
		}
	}

	if !taskFound {
		return fmt.Errorf("task with ID %s not found", taskID)
	}

	// Marshal the updated data
	updatedJSON, err := json.MarshalIndent(tasksFileData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks data: %w", err)
	}

	// Write back to the file atomically
	tmpFile := tasksFile + ".tmp"
	if err := os.WriteFile(tmpFile, updatedJSON, 0644); err != nil {
		return fmt.Errorf("failed to write temporary file: %w", err)
	}

	if err := os.Rename(tmpFile, tasksFile); err != nil {
		os.Remove(tmpFile)
		return fmt.Errorf("failed to rename temporary file: %w", err)
	}

	return nil
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
