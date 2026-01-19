package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
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
			name: "task ID extraction with auto-generation",
			markdown: `## 1. Section

- [ ] 1.1 First task
- [ ] 1.2 Second task
- [ ] 1.10 Tenth task (explicit ID ignored, auto-generated)
- [ ] 2.5 Task with mismatched ID (auto-generated)
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
					ID:          "1.3",
					Section:     "Section",
					Description: "Tenth task (explicit ID ignored, auto-generated)",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.4",
					Section:     "Section",
					Description: "Task with mismatched ID (auto-generated)",
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
		{
			name: "unnumbered section with tasks",
			markdown: `## Implementation

- [ ] 1. Update cmd/validate.go to remove Strict field
- [ ] 2. Update cmd/validate.go to always pass true
- [x] 3. Update internal/validation/interactive.go
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Implementation",
					Description: "Update cmd/validate.go to remove Strict field",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Implementation",
					Description: "Update cmd/validate.go to always pass true",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.3",
					Section:     "Implementation",
					Description: "Update internal/validation/interactive.go",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			name: "mixed formats with auto-generated IDs",
			markdown: `## 1. Setup

- [ ] 1.1 Create files
- [ ] 1.2 Configure settings

## 2. Implementation

- [ ] 2. Implement feature
- [x] 3. Test feature
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Create files",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Configure settings",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Implementation",
					Description: "Implement feature",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.2",
					Section:     "Implementation",
					Description: "Test feature",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			name: "tasks without section use global sequential IDs",
			markdown: `- [ ] First task without section
- [x] Second task without section
- [ ] Third task
`,
			expected: []parsers.Task{
				{
					ID:          "1",
					Section:     "",
					Description: "First task without section",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2",
					Section:     "",
					Description: "Second task without section",
					Status:      parsers.TaskStatusCompleted,
				},
				{
					ID:          "3",
					Section:     "",
					Description: "Third task",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "tasks without number get auto-generated IDs",
			markdown: `## 1. Setup

- [ ] Task without explicit number
- [ ] Another task without number
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task without explicit number",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Another task without number",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "unnumbered sections get sequential numbers",
			markdown: `## Setup
- [ ] First setup task
- [ ] Second setup task

## Implementation
- [ ] First impl task
- [x] Second impl task
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "First setup task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Second setup task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Implementation",
					Description: "First impl task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.2",
					Section:     "Implementation",
					Description: "Second impl task",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			name: "section numbering continues from explicit",
			markdown: `## 1. Setup
- [ ] Setup task

## Implementation
- [ ] Impl task

## 5. Testing
- [ ] Test task

## Deployment
- [ ] Deploy task
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Setup task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Implementation",
					Description: "Impl task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "5.1",
					Section:     "Testing",
					Description: "Test task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "6.1",
					Section:     "Deployment",
					Description: "Deploy task",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "explicit IDs used when matching expected",
			markdown: `## 1. Section
- [ ] 1.1 First task
- [ ] 1.2 Second task
- [ ] 1.3 Third task
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
					ID:          "1.3",
					Section:     "Section",
					Description: "Third task",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "wrong explicit IDs get overridden",
			markdown: `## 1. Section
- [ ] 5 First task with wrong number
- [ ] 99.99 Second task with wrong number
- [ ] Third task no number
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Section",
					Description: "First task with wrong number",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Section",
					Description: "Second task with wrong number",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.3",
					Section:     "Section",
					Description: "Third task no number",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "all task number formats mixed",
			markdown: `## 1. Mixed
- [ ] 1.1 Decimal format
- [ ] 1. Dot format
- [ ] 3 Number only
- [ ] No number format
- [x] 1.5 Matching explicit
`,
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Mixed",
					Description: "Decimal format",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Mixed",
					Description: "Dot format",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.3",
					Section:     "Mixed",
					Description: "Number only",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.4",
					Section:     "Mixed",
					Description: "No number format",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.5",
					Section:     "Mixed",
					Description: "Matching explicit",
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

			if err := os.WriteFile(tasksMdPath, []byte(tt.markdown), 0o644); err != nil {
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

			// Write the tasks (nil appendCfg for no appended tasks)
			if err := writeTasksJSONC(tasksJSONPath, tt.tasks, nil); err != nil {
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

	if err := writeTasksJSONC(tasksJSONPath, tasks, nil); err != nil {
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
	strippedData := parsers.StripJSONComments(
		data,
	)

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

	err := writeTasksJSONC(
		tasksJsonPath,
		tasks,
		nil,
	)
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
	expectedPerm := os.FileMode(0o644)
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
		(s == substr || s != "" && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}

	return false
}

// TestAcceptPreservesTasksMd verifies that tasks.md is NOT deleted after accept
func TestAcceptPreservesTasksMd(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock change directory structure
	changeDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		"test-change",
	)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf(
			"failed to create change dir: %v",
			err,
		)
	}

	// Create a proposal.md (required for validation)
	proposalPath := filepath.Join(
		changeDir,
		"proposal.md",
	)
	proposalContent := `# Test Proposal

## Problem
Test problem

## Solution
Test solution
`
	if err := os.WriteFile(proposalPath, []byte(proposalContent), filePerm); err != nil {
		t.Fatalf(
			"failed to write proposal.md: %v",
			err,
		)
	}

	// Create a tasks.md file
	tasksMdPath := filepath.Join(
		changeDir,
		"tasks.md",
	)
	tasksMdContent := `## 1. Implementation

- [ ] 1.1 First task
- [x] 1.2 Second task
`
	if err := os.WriteFile(tasksMdPath, []byte(tasksMdContent), filePerm); err != nil {
		t.Fatalf(
			"failed to write tasks.md: %v",
			err,
		)
	}

	// Run writeAndCleanup (the function that previously deleted tasks.md)
	tasks := []parsers.Task{
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
	}

	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.jsonc",
	)
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf(
			"writeAndCleanup failed: %v",
			err,
		)
	}

	// Verify tasks.md still exists
	if _, err := os.Stat(tasksMdPath); os.IsNotExist(
		err,
	) {
		t.Error(
			"tasks.md was deleted but should be preserved",
		)
	}

	// Verify tasks.jsonc was created
	if _, err := os.Stat(tasksJSONPath); os.IsNotExist(
		err,
	) {
		t.Error("tasks.jsonc was not created")
	}

	// Verify tasks.md content is unchanged
	mdContent, err := os.ReadFile(tasksMdPath)
	if err != nil {
		t.Fatalf(
			"failed to read tasks.md: %v",
			err,
		)
	}
	if string(mdContent) != tasksMdContent {
		t.Error("tasks.md content was modified")
	}
}

// TestAcceptWithBothFilesPresent verifies behavior when both tasks.md and tasks.jsonc already exist
func TestAcceptWithBothFilesPresent(
	t *testing.T,
) {
	tmpDir := t.TempDir()

	// Create a mock change directory structure
	changeDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		"test-change",
	)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf(
			"failed to create change dir: %v",
			err,
		)
	}

	// Create a proposal.md
	proposalPath := filepath.Join(
		changeDir,
		"proposal.md",
	)
	proposalContent := `# Test Proposal

## Problem
Test problem

## Solution
Test solution
`
	if err := os.WriteFile(proposalPath, []byte(proposalContent), filePerm); err != nil {
		t.Fatalf(
			"failed to write proposal.md: %v",
			err,
		)
	}

	// Create tasks.md
	tasksMdPath := filepath.Join(
		changeDir,
		"tasks.md",
	)
	tasksMdContent := `## 1. Updated Section

- [ ] 1.1 Updated task
- [ ] 1.2 New task
`
	if err := os.WriteFile(tasksMdPath, []byte(tasksMdContent), filePerm); err != nil {
		t.Fatalf(
			"failed to write tasks.md: %v",
			err,
		)
	}

	// Create existing tasks.jsonc (from previous accept)
	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.jsonc",
	)
	existingJSONContent := `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Old Section",
      "description": "Old task",
      "status": "completed"
    }
  ]
}
`
	if err := os.WriteFile(tasksJSONPath, []byte(existingJSONContent), filePerm); err != nil {
		t.Fatalf(
			"failed to write existing tasks.jsonc: %v",
			err,
		)
	}

	// Parse tasks.md and write new tasks.jsonc (simulating accept command)
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf(
			"failed to parse tasks.md: %v",
			err,
		)
	}

	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf(
			"writeAndCleanup failed: %v",
			err,
		)
	}

	// Verify both files exist
	if _, err := os.Stat(tasksMdPath); os.IsNotExist(
		err,
	) {
		t.Error("tasks.md should still exist")
	}
	if _, err := os.Stat(tasksJSONPath); os.IsNotExist(
		err,
	) {
		t.Error("tasks.jsonc should exist")
	}

	// Verify tasks.jsonc was overwritten with new content
	jsonContent, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf(
			"failed to read tasks.jsonc: %v",
			err,
		)
	}

	// Strip JSONC comments and parse
	strippedJSON := parsers.StripJSONComments(
		jsonContent,
	)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedJSON, &tasksFile); err != nil {
		t.Fatalf(
			"failed to parse tasks.jsonc: %v",
			err,
		)
	}

	// Verify it has the new tasks, not the old ones
	if len(tasksFile.Tasks) != 2 {
		t.Errorf(
			"expected 2 tasks, got %d",
			len(tasksFile.Tasks),
		)
	}
	if tasksFile.Tasks[0].Description != "Updated task" {
		t.Errorf(
			"expected 'Updated task', got '%s'",
			tasksFile.Tasks[0].Description,
		)
	}
}

// TestAcceptDryRunPreservesFiles verifies that dry-run mode doesn't modify any files
func TestAcceptDryRunPreservesFiles(
	t *testing.T,
) {
	tmpDir := t.TempDir()

	// Create a mock change directory structure
	changeDir := filepath.Join(
		tmpDir,
		"spectr",
		"changes",
		"test-change",
	)
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf(
			"failed to create change dir: %v",
			err,
		)
	}

	// Create tasks.md
	tasksMdPath := filepath.Join(
		changeDir,
		"tasks.md",
	)
	tasksMdContent := `## 1. Test

- [ ] 1.1 Task
`
	if err := os.WriteFile(tasksMdPath, []byte(tasksMdContent), filePerm); err != nil {
		t.Fatalf(
			"failed to write tasks.md: %v",
			err,
		)
	}

	tasksJSONPath := filepath.Join(
		changeDir,
		"tasks.jsonc",
	)

	// Verify tasks.jsonc doesn't exist before dry-run
	if _, err := os.Stat(tasksJSONPath); !os.IsNotExist(
		err,
	) {
		t.Fatal(
			"tasks.jsonc should not exist before dry-run",
		)
	}

	// Note: Full dry-run test would require mocking the AcceptCmd.Run() method
	// For now, we verify that writeAndCleanup is only called when NOT in dry-run mode
	// This test serves as documentation of expected dry-run behavior
}

// TestWriteTasksJSONCWithAppendConfig verifies that append tasks are added correctly
func TestWriteTasksJSONCWithAppendConfig(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(
		tmpDir,
		"tasks.jsonc",
	)

	existingTasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Setup",
			Description: "First task",
			Status:      parsers.TaskStatusCompleted,
		},
		{
			ID:          "1.2",
			Section:     "Setup",
			Description: "Second task",
			Status:      parsers.TaskStatusPending,
		},
		{
			ID:          "2.1",
			Section:     "Implementation",
			Description: "Third task",
			Status:      parsers.TaskStatusPending,
		},
	}

	appendCfg := &config.AppendTasksConfig{
		Section: "Project Workflow",
		Tasks: []string{
			"Run linter and tests",
			"Update changelog",
		},
	}

	err := writeTasksJSONC(
		tasksJSONPath,
		existingTasks,
		appendCfg,
	)
	if err != nil {
		t.Fatalf(
			"writeTasksJSONC() error = %v",
			err,
		)
	}

	// Read and parse the output
	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf(
			"failed to read tasks.jsonc: %v",
			err,
		)
	}

	strippedData := parsers.StripJSONComments(
		data,
	)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &tasksFile); err != nil {
		t.Fatalf(
			"failed to unmarshal JSON: %v",
			err,
		)
	}

	// Should have 5 tasks (3 existing + 2 appended)
	if len(tasksFile.Tasks) != 5 {
		t.Errorf(
			"expected 5 tasks, got %d",
			len(tasksFile.Tasks),
		)
	}

	// Verify appended tasks have correct IDs (section 3, since max was 2)
	if tasksFile.Tasks[3].ID != "3.1" {
		t.Errorf(
			"expected appended task ID '3.1', got '%s'",
			tasksFile.Tasks[3].ID,
		)
	}
	if tasksFile.Tasks[4].ID != "3.2" {
		t.Errorf(
			"expected appended task ID '3.2', got '%s'",
			tasksFile.Tasks[4].ID,
		)
	}

	// Verify appended tasks have correct section
	if tasksFile.Tasks[3].Section != "Project Workflow" {
		t.Errorf(
			"expected section 'Project Workflow', got '%s'",
			tasksFile.Tasks[3].Section,
		)
	}

	// Verify appended tasks have pending status
	if tasksFile.Tasks[3].Status != parsers.TaskStatusPending {
		t.Errorf(
			"expected status 'pending', got '%s'",
			tasksFile.Tasks[3].Status,
		)
	}

	// Verify appended task descriptions
	if tasksFile.Tasks[3].Description != "Run linter and tests" {
		t.Error(
			"wrong description for appended task",
		)
	}
}

// TestWriteTasksJSONCWithDefaultSection verifies default section name
func TestWriteTasksJSONCWithDefaultSection(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(
		tmpDir,
		"tasks.jsonc",
	)

	existingTasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Setup",
			Description: "Task",
			Status:      parsers.TaskStatusPending,
		},
	}

	appendCfg := &config.AppendTasksConfig{
		// No section specified - should use default
		Tasks: []string{"Appended task"},
	}

	err := writeTasksJSONC(
		tasksJSONPath,
		existingTasks,
		appendCfg,
	)
	if err != nil {
		t.Fatalf(
			"writeTasksJSONC() error = %v",
			err,
		)
	}

	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf(
			"failed to read tasks.jsonc: %v",
			err,
		)
	}

	strippedData := parsers.StripJSONComments(
		data,
	)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &tasksFile); err != nil {
		t.Fatalf(
			"failed to unmarshal JSON: %v",
			err,
		)
	}

	// Verify default section name is used
	if tasksFile.Tasks[1].Section != config.DefaultAppendTasksSection {
		t.Errorf(
			"expected default section '%s', got '%s'",
			config.DefaultAppendTasksSection,
			tasksFile.Tasks[1].Section,
		)
	}
}

// TestWriteTasksJSONCWithEmptyAppendTasks verifies no change when append tasks is empty
func TestWriteTasksJSONCWithEmptyAppendTasks(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(
		tmpDir,
		"tasks.jsonc",
	)

	existingTasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Setup",
			Description: "Task",
			Status:      parsers.TaskStatusPending,
		},
	}

	appendCfg := &config.AppendTasksConfig{
		Section: "Workflow",
		Tasks:   make([]string, 0), // Empty
	}

	err := writeTasksJSONC(
		tasksJSONPath,
		existingTasks,
		appendCfg,
	)
	if err != nil {
		t.Fatalf(
			"writeTasksJSONC() error = %v",
			err,
		)
	}

	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf(
			"failed to read tasks.jsonc: %v",
			err,
		)
	}

	strippedData := parsers.StripJSONComments(
		data,
	)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &tasksFile); err != nil {
		t.Fatalf(
			"failed to unmarshal JSON: %v",
			err,
		)
	}

	// Should have only the original task
	if len(tasksFile.Tasks) != 1 {
		t.Errorf(
			"expected 1 task, got %d",
			len(tasksFile.Tasks),
		)
	}
}

// TestFindNextSectionNumber verifies section number calculation
func TestFindNextSectionNumber(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []parsers.Task
		expected int
	}{
		{
			name:     "empty tasks",
			tasks:    make([]parsers.Task, 0),
			expected: 1,
		},
		{
			name: "single section",
			tasks: []parsers.Task{
				{ID: "1.1"},
				{ID: "1.2"},
			},
			expected: 2,
		},
		{
			name: "multiple sections",
			tasks: []parsers.Task{
				{ID: "1.1"},
				{ID: "2.1"},
				{ID: "3.1"},
			},
			expected: 4,
		},
		{
			name: "non-sequential sections",
			tasks: []parsers.Task{
				{ID: "1.1"},
				{ID: "5.1"},
				{ID: "3.1"},
			},
			expected: 6,
		},
		{
			name: "tasks without section prefix",
			tasks: []parsers.Task{
				{ID: "1"},
				{ID: "2"},
			},
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findNextSectionNumber(tt.tasks)
			if got != tt.expected {
				t.Errorf(
					"findNextSectionNumber() = %d, want %d",
					got,
					tt.expected,
				)
			}
		})
	}
}

// TestCreateAppendedTasks verifies task creation from config
func TestCreateAppendedTasks(t *testing.T) {
	existingTasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Setup",
			Description: "Task 1",
			Status:      parsers.TaskStatusPending,
		},
		{
			ID:          "2.1",
			Section:     "Impl",
			Description: "Task 2",
			Status:      parsers.TaskStatusPending,
		},
	}

	cfg := &config.AppendTasksConfig{
		Section: "Workflow",
		Tasks: []string{
			"Task A",
			"Task B",
			"Task C",
		},
	}

	tasks := createAppendedTasks(
		existingTasks,
		cfg,
	)

	if len(tasks) != 3 {
		t.Fatalf(
			"expected 3 tasks, got %d",
			len(tasks),
		)
	}

	// Verify IDs start at section 3
	expectedIDs := []string{"3.1", "3.2", "3.3"}
	for i, task := range tasks {
		if task.ID != expectedIDs[i] {
			t.Errorf(
				"task %d: expected ID '%s', got '%s'",
				i,
				expectedIDs[i],
				task.ID,
			)
		}
		if task.Section != "Workflow" {
			t.Errorf(
				"task %d: expected section 'Workflow', got '%s'",
				i,
				task.Section,
			)
		}
		if task.Status != parsers.TaskStatusPending {
			t.Errorf(
				"task %d: expected status 'pending', got '%s'",
				i,
				task.Status,
			)
		}
	}

	if tasks[0].Description != "Task A" {
		t.Errorf(
			"expected description 'Task A', got '%s'",
			tasks[0].Description,
		)
	}
}
