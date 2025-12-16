package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

func TestParseTasksMd(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected []parsers.Task
	}{
		{
			name: "basic parsing with sections and tasks",
			markdown: `## 1. Core Accept Command

- [ ] 1.1 Create cmd/accept.go with AcceptCmd struct following Kong patterns
- [x] 1.2 Add AcceptCmd to CLI struct in cmd/root.go
- [ ] 1.3 Implement Run() method that validates change exists

## 2. JSON Schema and Types

- [x] 2.1 Define TasksFile struct with version field and tasks array
- [ ] 2.2 Define Task struct with id, section, description, and status fields
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Core Accept Command",
					Description: "Create cmd/accept.go with AcceptCmd struct following Kong patterns",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Core Accept Command",
					Description: "Add AcceptCmd to CLI struct in cmd/root.go",
					Status:      parsers.TaskStatusCompleted,
				},
				{
					ID:          "1.3",
					Section:     "Core Accept Command",
					Description: "Implement Run() method that validates change exists",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "JSON Schema and Types",
					Description: "Define TasksFile struct with version field and tasks array",
					Status:      parsers.TaskStatusCompleted,
				},
				{
					ID:          "2.2",
					Section:     "JSON Schema and Types",
					Description: "Define Task struct with id, section, description, and status fields",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "completed vs pending tasks",
			markdown: `## 1. Testing

- [ ] 1.1 This task is pending
- [x] 1.2 This task is completed
- [X] 1.3 This task is also completed with uppercase X
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Testing",
					Description: "This task is pending",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Testing",
					Description: "This task is completed",
					Status:      parsers.TaskStatusCompleted,
				},
				{
					ID:          "1.3",
					Section:     "Testing",
					Description: "This task is also completed with uppercase X",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			name: "section extraction from headers",
			markdown: `## 1. First Section

- [ ] 1.1 Task in first section

## 2. Second Section

- [ ] 2.1 Task in second section

## 3. Third Section With Spaces

- [ ] 3.1 Task in third section
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "First Section",
					Description: "Task in first section",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Second Section",
					Description: "Task in second section",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "3.1",
					Section:     "Third Section With Spaces",
					Description: "Task in third section",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "task ID extraction",
			markdown: `## 1. Section

- [ ] 1.1 First task
- [ ] 1.2 Second task
- [ ] 1.10 Tenth task
- [ ] 2.5 Task in different section numbering
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Section",
					Description: "First task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Section",
					Description: "Second task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.10",
					Section:     "Section",
					Description: "Tenth task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.5",
					Section:     "Section",
					Description: "Task in different section numbering",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name:     "empty file",
			markdown: "",
			expected: nil,
		},
		{
			name: "file with only sections no tasks",
			markdown: `## 1. Section One

## 2. Section Two
`,
			expected: nil,
		},
		{
			name: "task with backticks and special characters",
			markdown: `## 1. Implementation

- [ ] 1.1 Create ` + "`cmd/accept.go`" + ` with ` + "`AcceptCmd`" + ` struct
- [x] 1.2 Add function that returns ` + "`*CLI`" + `
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Implementation",
					Description: "Create `cmd/accept.go` with `AcceptCmd` struct",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Implementation",
					Description: "Add function that returns `*CLI`",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp file with markdown content
			tmpDir := t.TempDir()
			tasksMdPath := filepath.Join(
				tmpDir,
				"tasks.md",
			)

			if err := os.WriteFile(tasksMdPath, []byte(tt.markdown), 0644); err != nil {
				t.Fatalf(
					"failed to write test file: %v",
					err,
				)
			}

			// Parse the file
			got, err := parseTasksMd(tasksMdPath)
			if err != nil {
				t.Fatalf(
					"parseTasksMd() error = %v",
					err,
				)
			}

			// Compare results
			if !reflect.DeepEqual(
				got,
				tt.expected,
			) {
				t.Errorf(
					"parseTasksMd() mismatch\ngot:  %+v\nwant: %+v",
					got,
					tt.expected,
				)
			}
		})
	}
}

func TestParseTasksMdFileNotFound(t *testing.T) {
	_, err := parseTasksMd(
		"/nonexistent/path/tasks.md",
	)
	if err == nil {
		t.Error(
			"parseTasksMd() expected error for nonexistent file, got nil",
		)
	}
}

func TestWriteTasksJson(t *testing.T) {
	tests := []struct {
		name  string
		tasks []parsers.Task
	}{
		{
			name: "basic writing with proper JSON structure",
			tasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Implementation",
					Description: "First task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Implementation",
					Description: "Second task",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			name:  "empty tasks array",
			tasks: make([]parsers.Task, 0),
		},
		{
			name: "single task",
			tasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Testing",
					Description: "Only task",
					Status:      parsers.TaskStatusInProgress,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tasksJSONPath := filepath.Join(
				tmpDir,
				"tasks.jsonc",
			)

			// Write the tasks
			if err := writeTasksJSONC(tasksJSONPath, tt.tasks); err != nil {
				t.Fatalf(
					"writeTasksJSONC() error = %v",
					err,
				)
			}

			// Read back and verify
			data, err := os.ReadFile(
				tasksJSONPath,
			)
			if err != nil {
				t.Fatalf(
					"failed to read written file: %v",
					err,
				)
			}

			// Strip JSONC comments before unmarshalling
			data = parsers.StripJSONComments(data)

			var tasksFile parsers.TasksFile
			if err := json.Unmarshal(data, &tasksFile); err != nil {
				t.Fatalf(
					"failed to unmarshal JSON: %v",
					err,
				)
			}

			// Verify version is 1
			if tasksFile.Version != 1 {
				t.Errorf(
					"version = %d, want 1",
					tasksFile.Version,
				)
			}

			// Verify tasks match
			if !reflect.DeepEqual(
				tasksFile.Tasks,
				tt.tasks,
			) {
				t.Errorf(
					"tasks mismatch\ngot:  %+v\nwant: %+v",
					tasksFile.Tasks,
					tt.tasks,
				)
			}
		})
	}
}

func TestWriteTasksJsonIndentation(t *testing.T) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(
		tmpDir,
		"tasks.jsonc",
	)

	tasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Test",
			Description: "Task",
			Status:      parsers.TaskStatusPending,
		},
	}

	if err := writeTasksJSONC(tasksJSONPath, tasks); err != nil {
		t.Fatalf(
			"writeTasksJSONC() error = %v",
			err,
		)
	}

	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read file: %v", err)
	}

	content := string(data)

	// Verify indentation uses 2 spaces (not tabs, not 4 spaces)
	expectedIndent := "  \"version\""
	if !contains(content, expectedIndent) {
		t.Error(
			"JSON indentation incorrect, expected 2-space indent",
		)
	}

	// Strip JSONC comments before unmarshalling
	strippedData := parsers.StripJSONComments(data)

	// Verify it's valid JSON with proper structure
	var parsed map[string]any
	if err := json.Unmarshal(strippedData, &parsed); err != nil {
		t.Errorf(
			"output is not valid JSONC: %v",
			err,
		)
	}

	// Verify required top-level keys exist
	if _, ok := parsed["version"]; !ok {
		t.Error("JSON missing 'version' key")
	}
	if _, ok := parsed["tasks"]; !ok {
		t.Error("JSON missing 'tasks' key")
	}
}

func TestWriteTasksJsonFilePermissions(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	tasksJsonPath := filepath.Join(
		tmpDir,
		"tasks.jsonc",
	)

	tasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Test",
			Description: "Task",
			Status:      parsers.TaskStatusPending,
		},
	}

	err := writeTasksJSONC(tasksJsonPath, tasks)
	if err != nil {
		t.Fatalf(
			"writeTasksJSONC() error = %v",
			err,
		)
	}

	info, err := os.Stat(tasksJsonPath)
	if err != nil {
		t.Fatalf("failed to stat file: %v", err)
	}

	// Verify file permissions are 0644
	expectedPerm := os.FileMode(0644)
	if info.Mode().Perm() != expectedPerm {
		t.Errorf(
			"file permissions = %o, want %o",
			info.Mode().Perm(),
			expectedPerm,
		)
	}
}

func TestAcceptCmdStructure(t *testing.T) {
	cmd := &AcceptCmd{}
	val := reflect.ValueOf(cmd).Elem()

	// Check ChangeID field exists
	changeIDField := val.FieldByName("ChangeID")
	if !changeIDField.IsValid() {
		t.Error(
			"AcceptCmd does not have ChangeID field",
		)
	}

	// Check DryRun field exists
	dryRunField := val.FieldByName("DryRun")
	if !dryRunField.IsValid() {
		t.Error(
			"AcceptCmd does not have DryRun field",
		)
	}

	// Check NoInteractive field exists
	noInteractiveField := val.FieldByName(
		"NoInteractive",
	)
	if !noInteractiveField.IsValid() {
		t.Error(
			"AcceptCmd does not have NoInteractive field",
		)
	}
}

func TestCLIHasAcceptCommand(t *testing.T) {
	cli := &CLI{}
	val := reflect.ValueOf(cli).Elem()
	acceptField := val.FieldByName("Accept")

	if !acceptField.IsValid() {
		t.Fatal(
			"CLI struct does not have Accept field",
		)
	}

	// Check the type
	if acceptField.Type().Name() != "AcceptCmd" {
		t.Errorf(
			"Accept field type: got %s, want AcceptCmd",
			acceptField.Type().Name(),
		)
	}
}

// contains is a helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}
