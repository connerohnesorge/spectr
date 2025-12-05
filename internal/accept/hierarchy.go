package accept

import "strings"

// organizeTaskHierarchy takes a flat list of tasks and organizes them
// into a hierarchy based on ID prefix matching
// (e.g., "1.1.1" is a subtask of "1.1").
func organizeTaskHierarchy(flatTasks []Task) []Task {
	if len(flatTasks) == 0 {
		return nil
	}

	// Create a map for quick lookup by ID
	taskMap := make(map[string]*Task)
	for i := range flatTasks {
		task := flatTasks[i]
		task.Subtasks = nil // Ensure clean slate
		taskMap[task.ID] = &task
	}

	// Result will hold only top-level tasks
	var result []Task

	// Process each task and assign to parent or root
	for i := range flatTasks {
		task := flatTasks[i]
		parentID := findParentID(task.ID)

		if parentID == "" {
			// Top-level task
			result = append(result, task)
		} else if parent, exists := taskMap[parentID]; exists {
			// Add as subtask to parent
			parent.Subtasks = append(parent.Subtasks, task)
		} else {
			// Parent doesn't exist, treat as top-level
			result = append(result, task)
		}
	}

	// Update result with modified tasks from taskMap that have subtasks
	return buildFinalHierarchy(result, taskMap)
}

// findParentID returns the parent task ID by removing the last segment.
// For "1.2.3", returns "1.2". For "1.2", returns "1". For "1", returns "".
func findParentID(taskID string) string {
	lastDot := strings.LastIndex(taskID, ".")
	if lastDot == -1 {
		return ""
	}

	return taskID[:lastDot]
}

// buildFinalHierarchy recursively builds the final task hierarchy.
func buildFinalHierarchy(tasks []Task, taskMap map[string]*Task) []Task {
	result := make([]Task, len(tasks))
	for i, task := range tasks {
		if mapped, exists := taskMap[task.ID]; exists {
			result[i] = Task{
				ID:          mapped.ID,
				Description: mapped.Description,
				Completed:   mapped.Completed,
				Subtasks:    buildFinalHierarchy(mapped.Subtasks, taskMap),
			}
		} else {
			result[i] = task
		}
	}

	return result
}
