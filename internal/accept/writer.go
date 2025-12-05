package accept

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// WriteTasksJSON writes a tasks.json file from parsed Section data.
// It uses atomic write (write to temp file, then rename) for safety.
func WriteTasksJSON(filePath, changeID string, sections []Section) error {
	tasksJSON := TasksJSON{
		Version:    "1.0",
		ChangeID:   changeID,
		AcceptedAt: time.Now().UTC().Format(time.RFC3339),
		Sections:   sections,
		Summary:    CalculateSummary(sections),
	}

	data, err := json.MarshalIndent(tasksJSON, "", "  ")
	if err != nil {
		return err
	}

	// Ensure the data ends with a newline
	data = append(data, '\n')

	// Write atomically: create temp file in same directory, then rename
	dir := filepath.Dir(filePath)
	tempFile, err := os.CreateTemp(dir, "tasks-*.json.tmp")
	if err != nil {
		return err
	}
	tempPath := tempFile.Name()

	// Ensure cleanup on error
	defer func() {
		// If tempPath still exists, we failed - clean up
		if _, statErr := os.Stat(tempPath); statErr == nil {
			_ = os.Remove(tempPath)
		}
	}()

	if _, err := tempFile.Write(data); err != nil {
		_ = tempFile.Close()

		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	// Atomic rename
	return os.Rename(tempPath, filePath)
}

// CalculateSummary recursively counts all tasks and completed tasks
// across sections.
func CalculateSummary(sections []Section) Summary {
	var total, completed int
	for _, section := range sections {
		t, c := countTasks(section.Tasks)
		total += t
		completed += c
	}

	return Summary{
		Total:     total,
		Completed: completed,
	}
}

// countTasks recursively counts tasks and their subtasks.
func countTasks(tasks []Task) (total, completed int) {
	for _, task := range tasks {
		total++
		if task.Completed {
			completed++
		}
		// Recursively count subtasks
		subTotal, subCompleted := countTasks(task.Subtasks)
		total += subTotal
		completed += subCompleted
	}

	return total, completed
}
