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
	// Children is an optional reference to a child tasks file (e.g., "$ref:tasks-1.jsonc")
	// Used in version 2 hierarchical task files
	Children string `json:"children,omitempty"`
}

// TasksFile represents the root structure of a tasks.json file
type TasksFile struct {
	Version int    `json:"version"`
	Tasks   []Task `json:"tasks"`
	// Parent is the parent task ID for child task files (version 2)
	Parent string `json:"parent,omitempty"`
	// Includes is a list of child task file patterns (version 2 root files)
	Includes []string `json:"includes,omitempty"`
}
