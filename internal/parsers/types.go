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
	// Children is a reference to child tasks file (version 2+)
	// Format: "$ref:specs/<capability>/tasks.jsonc"
	Children string `json:"children,omitempty"`
}

// Summary represents aggregated task status counts for hierarchical tasks
type Summary struct {
	// Total is the total number of tasks across all files
	Total int `json:"total"`
	// Completed is the number of tasks with status "completed"
	Completed int `json:"completed"`
	// InProgress is the number of tasks with status "in_progress"
	InProgress int `json:"in_progress"`
	// Pending is the number of tasks with status "pending"
	Pending int `json:"pending"`
}

// TasksFile represents the root structure of a tasks.json file
type TasksFile struct {
	// Version is the schema version (1 = flat, 2 = hierarchical)
	Version int `json:"version"`
	// Tasks is the list of tasks
	Tasks []Task `json:"tasks"`
	// Summary is optional aggregated counts for hierarchical tasks (version 2+)
	Summary *Summary `json:"summary,omitempty"`
	// Includes is glob patterns for auto-discovery of child task files (version 2+)
	Includes []string `json:"includes,omitempty"`
	// Parent is the parent task ID for child task files (version 2+)
	Parent string `json:"parent,omitempty"`
}

// TasksFileVersion1 represents the version 1 schema for backwards compatibility
type TasksFileVersion1 struct {
	Version int    `json:"version"`
	Tasks   []Task `json:"tasks"`
}
