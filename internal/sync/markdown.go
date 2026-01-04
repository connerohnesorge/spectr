package sync

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/markdown"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

// File permission for writing tasks.md
const filePermission = 0o644

// checkboxOffset is the offset from "- [" to the checkbox character
const checkboxOffset = 3

// SyncTasksToMarkdown updates tasks.md checkbox statuses from tasks.jsonc.
// Returns the number of tasks whose status was updated.
func SyncTasksToMarkdown(changeDir string) (int, error) {
	tasksJsoncPath := filepath.Join(changeDir, "tasks.jsonc")
	tasksMdPath := filepath.Join(changeDir, "tasks.md")

	// Skip if no tasks.jsonc (not yet accepted)
	if _, err := os.Stat(tasksJsoncPath); os.IsNotExist(err) {
		return 0, nil
	}

	// Read source of truth
	tasksFile, err := parsers.ReadTasksJson(tasksJsoncPath)
	if err != nil {
		return 0, fmt.Errorf("read tasks.jsonc: %w", err)
	}

	// Build ID -> status map
	statusMap := buildStatusMap(tasksFile.Tasks)

	// Check if tasks.md exists
	if _, err := os.Stat(tasksMdPath); os.IsNotExist(err) {
		// Skip sync if tasks.md doesn't exist (per spec: skip silently)
		return 0, nil
	}

	// Update existing tasks.md in-place
	return updateTasksMd(tasksMdPath, statusMap)
}

// buildStatusMap creates a map from task ID to checkbox character.
// pending/in_progress -> ' ' (unchecked), completed -> 'x' (checked)
func buildStatusMap(tasks []parsers.Task) map[string]rune {
	m := make(map[string]rune, len(tasks))
	for _, t := range tasks {
		if t.Status == parsers.TaskStatusCompleted {
			m[t.ID] = 'x'
		} else {
			m[t.ID] = ' ' // pending and in_progress both map to unchecked
		}
	}

	return m
}

// updateTasksMd reads tasks.md, updates checkbox statuses, writes back.
// Preserves all formatting, comments, and structure.
func updateTasksMd(path string, statusMap map[string]rune) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer func() { _ = file.Close() }()

	var lines []string
	var updated int
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		newLine, changed := updateTaskLine(line, statusMap)
		lines = append(lines, newLine)
		if changed {
			updated++
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	// Only write if changes were made
	if updated > 0 {
		content := strings.Join(lines, "\n")
		if !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		if err := os.WriteFile(path, []byte(content), filePermission); err != nil {
			return 0, err
		}
	}

	return updated, nil
}

// updateTaskLine updates a single line's checkbox if it's a task line.
// Returns the (possibly modified) line and whether it was changed.
func updateTaskLine(line string, statusMap map[string]rune) (string, bool) {
	match, ok := markdown.MatchFlexibleTask(line)
	if !ok {
		return line, false
	}

	taskID := match.Number
	if taskID == "" {
		return line, false
	}

	desiredStatus, exists := statusMap[taskID]
	if !exists {
		return line, false
	}

	// Compare checkbox states: both 'x' and 'X' count as checked
	currentIsChecked := match.Status == 'x' || match.Status == 'X'
	desiredIsChecked := desiredStatus == 'x' || desiredStatus == 'X'

	if currentIsChecked == desiredIsChecked {
		return line, false
	}

	// Update checkbox character
	// Find "- [" and update the character at checkboxOffset position
	idx := strings.Index(line, "- [")
	if idx == -1 {
		return line, false
	}
	checkboxIdx := idx + checkboxOffset

	newLine := line[:checkboxIdx] + string(desiredStatus) + line[checkboxIdx+1:]

	return newLine, true
}
