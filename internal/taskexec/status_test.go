package taskexec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestUpdateTaskStatus(t *testing.T) {
	tests := []struct {
		name           string
		initialContent string
		taskID         string
		newStatus      parsers.TaskStatusValue
		wantContent    string
		wantErr        bool
	}{
		{
			name: "update task from pending to in_progress",
			initialContent: `{
				"version": 1,
				"tasks": [
					{
						"id": "1.1",
						"section": "Test",
						"description": "First task",
						"status": "pending"
					}
				]
			}`,
			taskID:    "1.1",
			newStatus: parsers.TaskStatusInProgress,
			wantContent: `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Test",
      "description": "First task",
      "status": "in_progress"
    }
  ]
}`,
			wantErr: false,
		},
		{
			name: "update task from in_progress to completed",
			initialContent: `{
				"version": 1,
				"tasks": [
					{
						"id": "1.1",
						"section": "Test",
						"description": "First task",
						"status": "in_progress"
					}
				]
			}`,
			taskID:    "1.1",
			newStatus: parsers.TaskStatusCompleted,
			wantContent: `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Test",
      "description": "First task",
      "status": "completed"
    }
  ]
}`,
			wantErr: false,
		},
		{
			name: "update task with comments",
			initialContent: `{
				// Version comment
				"version": 1,
				"tasks": [
					{
						"id": "1.1",
						"section": "Test",
						"description": "First task",
						"status": "pending" // Initial status
					}
				]
			}`,
			taskID:    "1.1",
			newStatus: parsers.TaskStatusCompleted,
			wantContent: `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Test",
      "description": "First task",
      "status": "completed"
    }
  ]
}`,
			wantErr: false,
		},
		{
			name: "task not found",
			initialContent: `{
				"version": 1,
				"tasks": [
					{
						"id": "1.1",
						"section": "Test",
						"description": "First task",
						"status": "pending"
					}
				]
			}`,
			taskID:      "2.1",
			newStatus:   parsers.TaskStatusCompleted,
			wantContent: "",
			wantErr:     true,
		},
		{
			name: "multiple tasks",
			initialContent: `{
				"version": 1,
				"tasks": [
					{
						"id": "1.1",
						"section": "Test",
						"description": "First task",
						"status": "completed"
					},
					{
						"id": "1.2",
						"section": "Test",
						"description": "Second task",
						"status": "pending"
					},
					{
						"id": "1.3",
						"section": "Test",
						"description": "Third task",
						"status": "in_progress"
					}
				]
			}`,
			taskID:    "1.2",
			newStatus: parsers.TaskStatusInProgress,
			wantContent: `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Test",
      "description": "First task",
      "status": "completed"
    },
    {
      "id": "1.2",
      "section": "Test",
      "description": "Second task",
      "status": "in_progress"
    },
    {
      "id": "1.3",
      "section": "Test",
      "description": "Third task",
      "status": "in_progress"
    }
  ]
}`,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir := t.TempDir()

			// Write initial tasks file
			tasksFile := filepath.Join(tempDir, "tasks.jsonc")
			if err := os.WriteFile(tasksFile, []byte(tt.initialContent), 0o644); err != nil {
				t.Fatalf("Failed to write tasks file: %v", err)
			}

			// Create status updater
			su := NewStatusUpdater(tempDir)

			// Update task status
			err := su.UpdateTaskStatus(tt.taskID, tt.newStatus)

			// Check results
			if tt.wantErr {
				if err == nil {
					t.Error("Expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)

				return
			}

			// Read the updated file
			updatedContent, err := os.ReadFile(tasksFile)
			if err != nil {
				t.Fatalf("Failed to read updated file: %v", err)
			}

			if string(updatedContent) != tt.wantContent {
				t.Errorf(
					"File content mismatch.\nExpected:\n%s\nGot:\n%s",
					tt.wantContent,
					string(updatedContent),
				)
			}
		})
	}
}

func TestUpdateTaskStatusAtomic(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Write initial tasks file
	initialContent := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test",
				"description": "First task",
				"status": "pending"
			}
		]
	}`

	tasksFile := filepath.Join(tempDir, "tasks.jsonc")
	if err := os.WriteFile(tasksFile, []byte(initialContent), 0o644); err != nil {
		t.Fatalf("Failed to write tasks file: %v", err)
	}

	// Create status updater
	su := NewStatusUpdater(tempDir)

	// Update task status
	if err := su.UpdateTaskStatus("1.1", parsers.TaskStatusCompleted); err != nil {
		t.Fatalf("Failed to update task status: %v", err)
	}

	// Verify the file was updated
	updatedContent, err := os.ReadFile(tasksFile)
	if err != nil {
		t.Fatalf("Failed to read updated file: %v", err)
	}

	if !contains(string(updatedContent), `"status": "completed"`) {
		t.Error("Task status was not updated correctly")
	}

	// Verify no temporary file remains
	if _, err := os.Stat(tasksFile + ".tmp"); !os.IsNotExist(err) {
		t.Error("Temporary file was not cleaned up")
	}
}

func contains(s, substr string) bool {
	if len(s) < len(substr) {
		return false
	}
	if s == substr {
		return true
	}
	if s == "" {
		return false
	}

	return s[:len(substr)] == substr || contains(s[1:], substr)
}
