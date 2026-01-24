package discovery

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/utils"
)

func TestFindNextPendingTask(t *testing.T) {
	tests := []struct {
		name         string
		tasksContent string
		wantTaskID   string
		wantErr      bool
	}{
		{
			name: "finds first pending task",
			tasksContent: `{
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
						"status": "pending"
					}
				]
			}`,
			wantTaskID: "1.2",
			wantErr:    false,
		},
		{
			name: "finds task with comments",
			tasksContent: `{
				// Version of the tasks file
				"version": 1,
				"tasks": [
					// This is a comment
					{
						"id": "1.1",
						"section": "Test",
						"description": "First task",
						"status": "pending" // This task is pending
					}
				]
			}`,
			wantTaskID: "1.1",
			wantErr:    false,
		},
		{
			name: "no pending tasks",
			tasksContent: `{
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
					}
				]
			}`,
			wantTaskID: "",
			wantErr:    true,
		},
		{
			name: "invalid json",
			tasksContent: `{
				"version": 1,
				"tasks": [
					{
						"id": "1.1",
						"section": "Test",
						"description": "First task",
						"status": "pending",
					}
				]
			}`,
			wantTaskID: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir := t.TempDir()

			// Write test tasks file
			tasksFile := filepath.Join(tempDir, "tasks.jsonc")
			if err := os.WriteFile(tasksFile, []byte(tt.tasksContent), 0o644); err != nil {
				t.Fatalf("Failed to write tasks file: %v", err)
			}

			// Create task discovery
			td := NewTaskDiscovery(tempDir)

			// Find next pending task
			task, err := td.FindNextPendingTask()

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

			if task.ID != tt.wantTaskID {
				t.Errorf("Expected task ID %s, got %s", tt.wantTaskID, task.ID)
			}
		})
	}
}

func TestStripJSONCComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "removes single line comments",
			input: `// This is a comment
{
	"version": 1 // Inline comment
}`,
			expected: `{
	"version": 1
}`,
		},
		{
			name: "preserves comments in strings",
			input: `{
	"description": "This contains // not a comment"
}`,
			expected: `{
	"description": "This contains // not a comment"
}`,
		},
		{
			name: "handles empty lines and spaces",
			input: `// Comment

{
	// Another comment
	"test": "value"
}`,
			expected: `{
	"test": "value"
}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(utils.StripJSONCComments([]byte(tt.input)))
			if result != tt.expected {
				t.Errorf("Expected:\n%s\nGot:\n%s", tt.expected, result)
			}
		})
	}
}
