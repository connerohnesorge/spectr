package parsers

import (
	"encoding/json"
	"testing"
)

func TestTask_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected string
	}{
		{
			name: "task without children",
			task: Task{
				ID:          "1.1",
				Section:     "Implementation",
				Description: "Create database schema",
				Status:      TaskStatusPending,
			},
			expected: `{"id":"1.1","section":"Implementation","description":"Create database schema","status":"pending"}`,
		},
		{
			name: "task with children reference",
			task: Task{
				ID:          "5",
				Section:     "Migrate Providers",
				Description: "Migrate all providers to new interface",
				Status:      TaskStatusInProgress,
				Children:    "$ref:specs/support-aider/tasks.jsonc",
			},
			expected: `{"id":"5","section":"Migrate Providers","description":"Migrate all providers to new interface","status":"in_progress","children":"$ref:specs/support-aider/tasks.jsonc"}`,
		},
		{
			name: "completed task",
			task: Task{
				ID:          "2.3",
				Section:     "Testing",
				Description: "Add unit tests for auth module",
				Status:      TaskStatusCompleted,
			},
			expected: `{"id":"2.3","section":"Testing","description":"Add unit tests for auth module","status":"completed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.task)
			if err != nil {
				t.Fatalf("failed to marshal task: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("marshaled task = %s, want %s", string(data), tt.expected)
			}

			// Verify unmarshaling works too
			var unmarshaled Task
			if err := json.Unmarshal([]byte(tt.expected), &unmarshaled); err != nil {
				t.Fatalf("failed to unmarshal task: %v", err)
			}

			if unmarshaled.ID != tt.task.ID {
				t.Errorf("unmarshaled ID = %s, want %s", unmarshaled.ID, tt.task.ID)
			}
			if unmarshaled.Section != tt.task.Section {
				t.Errorf("unmarshaled Section = %s, want %s", unmarshaled.Section, tt.task.Section)
			}
			if unmarshaled.Description != tt.task.Description {
				t.Errorf("unmarshaled Description = %s, want %s", unmarshaled.Description, tt.task.Description)
			}
			if unmarshaled.Status != tt.task.Status {
				t.Errorf("unmarshaled Status = %s, want %s", unmarshaled.Status, tt.task.Status)
			}
			if unmarshaled.Children != tt.task.Children {
				t.Errorf("unmarshaled Children = %s, want %s", unmarshaled.Children, tt.task.Children)
			}
		})
	}
}

func TestSummary_JSONSerialization(t *testing.T) {
	tests := []struct {
		name     string
		summary  Summary
		expected string
	}{
		{
			name: "empty summary",
			summary: Summary{
				Total:      0,
				Completed:  0,
				InProgress: 0,
				Pending:    0,
			},
			expected: `{"total":0,"completed":0,"in_progress":0,"pending":0}`,
		},
		{
			name: "partial progress",
			summary: Summary{
				Total:      10,
				Completed:  3,
				InProgress: 1,
				Pending:    6,
			},
			expected: `{"total":10,"completed":3,"in_progress":1,"pending":6}`,
		},
		{
			name: "all completed",
			summary: Summary{
				Total:      5,
				Completed:  5,
				InProgress: 0,
				Pending:    0,
			},
			expected: `{"total":5,"completed":5,"in_progress":0,"pending":0}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.summary)
			if err != nil {
				t.Fatalf("failed to marshal summary: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("marshaled summary = %s, want %s", string(data), tt.expected)
			}
		})
	}
}

func TestTasksFile_Version1(t *testing.T) {
	// Test backwards compatibility with version 1 format
	jsonData := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Implementation",
				"description": "Create database schema",
				"status": "pending"
			}
		]
	}`

	var tasksFile TasksFile
	if err := json.Unmarshal([]byte(jsonData), &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal version 1 tasks file: %v", err)
	}

	if tasksFile.Version != 1 {
		t.Errorf("Version = %d, want 1", tasksFile.Version)
	}

	if len(tasksFile.Tasks) != 1 {
		t.Errorf("len(Tasks) = %d, want 1", len(tasksFile.Tasks))
	}

	if tasksFile.Tasks[0].ID != "1.1" {
		t.Errorf("Tasks[0].ID = %s, want 1.1", tasksFile.Tasks[0].ID)
	}

	// Version 1 should not have Summary, Includes, or Parent
	if tasksFile.Summary != nil {
		t.Errorf("Summary should be nil for version 1")
	}
	if tasksFile.Includes != nil {
		t.Errorf("Includes should be nil for version 1")
	}
	if tasksFile.Parent != "" {
		t.Errorf("Parent should be empty for version 1")
	}
}

func TestTasksFile_Version2(t *testing.T) {
	// Test version 2 format with hierarchical structure
	jsonData := `{
		"version": 2,
		"summary": {
			"total": 10,
			"completed": 3,
			"in_progress": 1,
			"pending": 6
		},
		"tasks": [
			{
				"id": "1",
				"section": "Foundation",
				"description": "Create core interfaces",
				"status": "completed"
			},
			{
				"id": "5",
				"section": "Migrate Providers",
				"description": "Migrate all providers to new interface",
				"status": "in_progress",
				"children": "$ref:specs/support-aider/tasks.jsonc"
			}
		],
		"includes": ["specs/*/tasks.jsonc"]
	}`

	var tasksFile TasksFile
	if err := json.Unmarshal([]byte(jsonData), &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal version 2 tasks file: %v", err)
	}

	if tasksFile.Version != 2 {
		t.Errorf("Version = %d, want 2", tasksFile.Version)
	}

	if tasksFile.Summary == nil {
		t.Fatal("Summary should not be nil for version 2")
	}

	if tasksFile.Summary.Total != 10 {
		t.Errorf("Summary.Total = %d, want 10", tasksFile.Summary.Total)
	}

	if tasksFile.Summary.Completed != 3 {
		t.Errorf("Summary.Completed = %d, want 3", tasksFile.Summary.Completed)
	}

	if len(tasksFile.Includes) != 1 || tasksFile.Includes[0] != "specs/*/tasks.jsonc" {
		t.Errorf("Includes = %v, want [specs/*/tasks.jsonc]", tasksFile.Includes)
	}

	// Verify task with children reference
	migratedTask := tasksFile.Tasks[1]
	if migratedTask.Children != "$ref:specs/support-aider/tasks.jsonc" {
		t.Errorf("Tasks[1].Children = %s, want $ref:specs/support-aider/tasks.jsonc", migratedTask.Children)
	}
}

func TestTasksFile_ChildFile(t *testing.T) {
	// Test child file format (version 2 with parent)
	jsonData := `{
		"version": 2,
		"parent": "5",
		"tasks": [
			{
				"id": "5.1",
				"description": "Migrate aider.go to new Provider interface",
				"status": "pending"
			},
			{
				"id": "5.2",
				"description": "Add unit tests for Aider provider",
				"status": "pending"
			}
		]
	}`

	var tasksFile TasksFile
	if err := json.Unmarshal([]byte(jsonData), &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal child tasks file: %v", err)
	}

	if tasksFile.Version != 2 {
		t.Errorf("Version = %d, want 2", tasksFile.Version)
	}

	if tasksFile.Parent != "5" {
		t.Errorf("Parent = %s, want 5", tasksFile.Parent)
	}

	// Child files should not have Summary or Includes
	if tasksFile.Summary != nil {
		t.Errorf("Summary should be nil for child file")
	}
	if tasksFile.Includes != nil {
		t.Errorf("Includes should be nil for child file")
	}

	// Child tasks should not have Section or Children
	for _, task := range tasksFile.Tasks {
		if task.Section != "" {
			t.Errorf("Task %s Section should be empty for child file, got %s", task.ID, task.Section)
		}
		if task.Children != "" {
			t.Errorf("Task %s Children should be empty for child file, got %s", task.ID, task.Children)
		}
	}
}

func TestTasksFile_MarshalVersion2(t *testing.T) {
	// Test marshaling version 2 format
	tasksFile := TasksFile{
		Version: 2,
		Summary: &Summary{
			Total:      10,
			Completed:  3,
			InProgress: 1,
			Pending:    6,
		},
		Tasks: []Task{
			{
				ID:          "1",
				Section:     "Foundation",
				Description: "Create core interfaces",
				Status:      TaskStatusCompleted,
			},
			{
				ID:          "5",
				Section:     "Migrate Providers",
				Description: "Migrate all providers",
				Status:      TaskStatusInProgress,
				Children:    "$ref:specs/support-aider/tasks.jsonc",
			},
		},
		Includes: []string{"specs/*/tasks.jsonc"},
	}

	data, err := json.MarshalIndent(tasksFile, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal version 2 tasks file: %v", err)
	}

	// Verify the marshaled output contains expected fields
	jsonStr := string(data)
	if !contains(jsonStr, `"version": 2`) {
		t.Error("marshaled output should contain version 2")
	}
	if !contains(jsonStr, `"summary"`) {
		t.Error("marshaled output should contain summary")
	}
	if !contains(jsonStr, `"includes"`) {
		t.Error("marshaled output should contain includes")
	}
	if !contains(jsonStr, `"children"`) {
		t.Error("marshaled output should contain children field")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestTaskStatusValues(t *testing.T) {
	// Verify status values are correct
	if TaskStatusPending != "pending" {
		t.Errorf("TaskStatusPending = %s, want pending", TaskStatusPending)
	}
	if TaskStatusInProgress != "in_progress" {
		t.Errorf("TaskStatusInProgress = %s, want in_progress", TaskStatusInProgress)
	}
	if TaskStatusCompleted != "completed" {
		t.Errorf("TaskStatusCompleted = %s, want completed", TaskStatusCompleted)
	}
}
