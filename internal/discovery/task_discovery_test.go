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

func TestFindNextPendingTaskHierarchical(t *testing.T) {
	tests := []struct {
		name       string
		setupFiles func(t *testing.T, tempDir string)
		wantTaskID string
		wantErr    bool
	}{
		{
			name: "v2 hierarchical - finds pending task in child file",
			setupFiles: func(t *testing.T, tempDir string) {
				// Root tasks.jsonc
				rootTasks := `{
					"version": 2,
					"tasks": [
						{
							"id": "1",
							"section": "Foundation",
							"description": "Setup core",
							"status": "completed"
						},
						{
							"id": "2",
							"section": "Features",
							"description": "Implement features",
							"status": "in_progress",
							"children": "$ref:specs/features/tasks.jsonc"
						}
					]
				}`
				os.WriteFile(filepath.Join(tempDir, "tasks.jsonc"), []byte(rootTasks), 0o644)

				// Create child directory
				os.MkdirAll(filepath.Join(tempDir, "specs", "features"), 0o755)

				// Child tasks.jsonc
				childTasks := `{
					"version": 2,
					"parent": "2",
					"tasks": [
						{
							"id": "2.1",
							"description": "First feature task",
							"status": "completed"
						},
						{
							"id": "2.2",
							"description": "Second feature task",
							"status": "pending"
						}
					]
				}`
				os.WriteFile(
					filepath.Join(tempDir, "specs", "features", "tasks.jsonc"),
					[]byte(childTasks),
					0o644,
				)
			},
			wantTaskID: "2.2",
			wantErr:    false,
		},
		{
			name: "v2 hierarchical - skips completed parent with all children completed",
			setupFiles: func(t *testing.T, tempDir string) {
				// Root tasks.jsonc
				rootTasks := `{
					"version": 2,
					"tasks": [
						{
							"id": "1",
							"section": "Features",
							"description": "All done",
							"status": "completed",
							"children": "$ref:specs/features/tasks.jsonc"
						},
						{
							"id": "2",
							"section": "Next",
							"description": "Next task",
							"status": "pending"
						}
					]
				}`
				os.WriteFile(filepath.Join(tempDir, "tasks.jsonc"), []byte(rootTasks), 0o644)

				// Create child directory
				os.MkdirAll(filepath.Join(tempDir, "specs", "features"), 0o755)

				// Child tasks.jsonc (all completed)
				childTasks := `{
					"version": 2,
					"parent": "1",
					"tasks": [
						{
							"id": "1.1",
							"description": "Done",
							"status": "completed"
						}
					]
				}`
				os.WriteFile(
					filepath.Join(tempDir, "specs", "features", "tasks.jsonc"),
					[]byte(childTasks),
					0o644,
				)
			},
			wantTaskID: "2",
			wantErr:    false,
		},
		{
			name: "v2 hierarchical - handles missing child file gracefully",
			setupFiles: func(t *testing.T, tempDir string) {
				// Root tasks.jsonc with $ref to non-existent file
				rootTasks := `{
					"version": 2,
					"tasks": [
						{
							"id": "1",
							"section": "Features",
							"description": "Missing child",
							"status": "in_progress",
							"children": "$ref:specs/missing/tasks.jsonc"
						},
						{
							"id": "2",
							"section": "Next",
							"description": "Next task",
							"status": "pending"
						}
					]
				}`
				os.WriteFile(filepath.Join(tempDir, "tasks.jsonc"), []byte(rootTasks), 0o644)
			},
			wantTaskID: "2",
			wantErr:    false,
		},
		{
			name: "v2 hierarchical - detects circular references",
			setupFiles: func(t *testing.T, tempDir string) {
				// Root tasks.jsonc
				rootTasks := `{
					"version": 2,
					"tasks": [
						{
							"id": "1",
							"section": "Circular",
							"description": "Circular ref",
							"status": "in_progress",
							"children": "$ref:specs/a/tasks.jsonc"
						}
					]
				}`
				os.WriteFile(filepath.Join(tempDir, "tasks.jsonc"), []byte(rootTasks), 0o644)

				// Create child directories
				os.MkdirAll(filepath.Join(tempDir, "specs", "a"), 0o755)
				os.MkdirAll(filepath.Join(tempDir, "specs", "b"), 0o755)

				// Child A references Child B
				childA := `{
					"version": 2,
					"parent": "1",
					"tasks": [
						{
							"id": "1.1",
							"description": "Ref to B",
							"status": "in_progress",
							"children": "$ref:../b/tasks.jsonc"
						}
					]
				}`
				os.WriteFile(
					filepath.Join(tempDir, "specs", "a", "tasks.jsonc"),
					[]byte(childA),
					0o644,
				)

				// Child B references back to root (circular)
				childB := `{
					"version": 2,
					"parent": "1.1",
					"tasks": [
						{
							"id": "1.1.1",
							"description": "Ref back to root",
							"status": "in_progress",
							"children": "$ref:../../tasks.jsonc"
						}
					]
				}`
				os.WriteFile(
					filepath.Join(tempDir, "specs", "b", "tasks.jsonc"),
					[]byte(childB),
					0o644,
				)
			},
			wantTaskID: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary directory
			tempDir := t.TempDir()

			// Setup test files
			tt.setupFiles(t, tempDir)

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
