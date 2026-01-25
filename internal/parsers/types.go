// Package parsers provides utilities for extracting and counting
// information from markdown specification files, including titles,
// tasks, deltas, and requirements.
//
// This file contains JSON schema types for the tasks.json file format.
package parsers

// TaskStatusValue represents the status of a task in tasks.json
type TaskStatusValue string

// Task status constants for JSON serialization
const (
	TaskStatusPending    TaskStatusValue = "pending"
	TaskStatusInProgress TaskStatusValue = "in_progress"
	TaskStatusCompleted  TaskStatusValue = "completed"
)

// Task represents a single task in tasks.json
type Task struct {
	// ID is the task identifier, e.g., "1.1", "2.3"
	ID string `json:"id"`
	// Section is the task section, e.g., "Implementation", "Testing"
	Section string `json:"section"`
	// Description is the full task description text
	Description string `json:"description"`
	// Status is one of pending, in_progress, completed
	Status TaskStatusValue `json:"status"`
	// Children is a $ref to a child task file (v2 hierarchical format)
	// Format: "$ref:specs/capability/tasks.jsonc"
	Children string `json:"children,omitempty"`
}

// TaskSummary represents task completion statistics
type TaskSummary struct {
	Total      int `json:"total"`
	Completed  int `json:"completed"`
	InProgress int `json:"in_progress"`
	Pending    int `json:"pending"`
}

// TasksFile represents the root structure of a tasks.json file
type TasksFile struct {
	Version int    `json:"version"`
	Tasks   []Task `json:"tasks"`
	// Summary provides quick overview (v2 format only)
	Summary *TaskSummary `json:"summary,omitempty"`
	// Includes is a list of glob patterns for child task files (v2 format only)
	Includes []string `json:"includes,omitempty"`
	// Parent is the parent task ID (used in child task files, v2 format only)
	Parent string `json:"parent,omitempty"`
}
