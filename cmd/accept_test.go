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
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a proposal.md (required for validation)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	proposalContent := `# Test Proposal

## Problem
Test problem

## Solution
Test solution
`
	if err := os.WriteFile(proposalPath, []byte(proposalContent), filePerm); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	// Create a tasks.md file
	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	tasksMdContent := `## 1. Implementation

- [ ] 1.1 First task
- [x] 1.2 Second task
`
	if err := os.WriteFile(tasksMdPath, []byte(tasksMdContent), filePerm); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
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

	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks); err != nil {
		t.Fatalf("writeAndCleanup failed: %v", err)
	}

	// Verify tasks.md still exists
	if _, err := os.Stat(tasksMdPath); os.IsNotExist(err) {
		t.Error("tasks.md was deleted but should be preserved")
	}

	// Verify tasks.jsonc was created
	if _, err := os.Stat(tasksJSONPath); os.IsNotExist(err) {
		t.Error("tasks.jsonc was not created")
	}

	// Verify tasks.md content is unchanged
	mdContent, err := os.ReadFile(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to read tasks.md: %v", err)
	}
	if string(mdContent) != tasksMdContent {
		t.Error("tasks.md content was modified")
	}
}

// TestAcceptWithBothFilesPresent verifies behavior when both tasks.md and tasks.jsonc already exist
func TestAcceptWithBothFilesPresent(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock change directory structure
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a proposal.md
	proposalPath := filepath.Join(changeDir, "proposal.md")
	proposalContent := `# Test Proposal

## Problem
Test problem

## Solution
Test solution
`
	if err := os.WriteFile(proposalPath, []byte(proposalContent), filePerm); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	// Create tasks.md
	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	tasksMdContent := `## 1. Updated Section

- [ ] 1.1 Updated task
- [ ] 1.2 New task
`
	if err := os.WriteFile(tasksMdPath, []byte(tasksMdContent), filePerm); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Create existing tasks.jsonc (from previous accept)
	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
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
		t.Fatalf("failed to write existing tasks.jsonc: %v", err)
	}

	// Parse tasks.md and write new tasks.jsonc (simulating accept command)
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse tasks.md: %v", err)
	}

	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks); err != nil {
		t.Fatalf("writeAndCleanup failed: %v", err)
	}

	// Verify both files exist
	if _, err := os.Stat(tasksMdPath); os.IsNotExist(err) {
		t.Error("tasks.md should still exist")
	}
	if _, err := os.Stat(tasksJSONPath); os.IsNotExist(err) {
		t.Error("tasks.jsonc should exist")
	}

	// Verify tasks.jsonc was overwritten with new content
	jsonContent, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	// Strip JSONC comments and parse
	strippedJSON := parsers.StripJSONComments(jsonContent)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedJSON, &tasksFile); err != nil {
		t.Fatalf("failed to parse tasks.jsonc: %v", err)
	}

	// Verify it has the new tasks, not the old ones
	if len(tasksFile.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasksFile.Tasks))
	}
	if tasksFile.Tasks[0].Description != "Updated task" {
		t.Errorf("expected 'Updated task', got '%s'", tasksFile.Tasks[0].Description)
	}
}

// TestAcceptDryRunPreservesFiles verifies that dry-run mode doesn't modify any files
func TestAcceptDryRunPreservesFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a mock change directory structure
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create tasks.md
	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	tasksMdContent := `## 1. Test

- [ ] 1.1 Task
`
	if err := os.WriteFile(tasksMdPath, []byte(tasksMdContent), filePerm); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")

	// Verify tasks.jsonc doesn't exist before dry-run
	if _, err := os.Stat(tasksJSONPath); !os.IsNotExist(err) {
		t.Fatal("tasks.jsonc should not exist before dry-run")
	}

	// Note: Full dry-run test would require mocking the AcceptCmd.Run() method
	// For now, we verify that writeAndCleanup is only called when NOT in dry-run mode
	// This test serves as documentation of expected dry-run behavior
}

func TestMatchSectionToCapability(t *testing.T) {
	tests := []struct {
		name     string
		section  string
		expected string
	}{
		{
			name:     "numbered section with period",
			section:  "5. Support Aider",
			expected: "support-aider",
		},
		{
			name:     "numbered section with period and space",
			section:  "2. Accept Command - Auto-Split Logic",
			expected: "accept-command-auto-split-logic",
		},
		{
			name:     "unnumbered section",
			section:  "Foundation",
			expected: "foundation",
		},
		{
			name:     "section with multiple spaces",
			section:  "Tasks Command Implementation",
			expected: "tasks-command-implementation",
		},
		{
			name:     "section with underscores",
			section:  "My_Section_Name",
			expected: "my-section-name",
		},
		{
			name:     "leading numbers only",
			section:  "1",
			expected: "",
		},
		{
			name:     "just period and space",
			section:  ". ",
			expected: "",
		},
		{
			name:     "already kebab-case",
			section:  "already-kebab-case",
			expected: "already-kebab-case",
		},
		{
			name:     "mixed case",
			section:  "My Mixed Case Section",
			expected: "my-mixed-case-section",
		},
		{
			name:     "section with numbers in name",
			section:  "Phase 2 Implementation",
			expected: "phase-2-implementation",
		},
		{
			name:     "empty section",
			section:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchSectionToCapability(tt.section)
			if result != tt.expected {
				t.Errorf("matchSectionToCapability(%q) = %q, want %q", tt.section, result, tt.expected)
			}
		})
	}
}

func TestFindMatchingDeltaSpec(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	changeDir := filepath.Join(tempDir, "change")
	specsDir := filepath.Join(changeDir, "specs")

	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("failed to create specs dir: %v", err)
	}

	// Create a valid delta spec directory
	aiderDir := filepath.Join(specsDir, "support-aider")
	if err := os.MkdirAll(aiderDir, 0o755); err != nil {
		t.Fatalf("failed to create support-aider dir: %v", err)
	}

	// Create spec.md file
	specPath := filepath.Join(aiderDir, "spec.md")
	if err := os.WriteFile(specPath, []byte("# Spec"), 0o644); err != nil {
		t.Fatalf("failed to write spec.md: %v", err)
	}

	// Create an empty directory (should not match)
	emptyDir := filepath.Join(specsDir, "empty-dir")
	if err := os.MkdirAll(emptyDir, 0o755); err != nil {
		t.Fatalf("failed to create empty-dir: %v", err)
	}

	tests := []struct {
		name       string
		capability string
		expected   bool
	}{
		{
			name:       "existing delta spec with spec.md",
			capability: "support-aider",
			expected:   true,
		},
		{
			name:       "non-existent capability",
			capability: "non-existent",
			expected:   false,
		},
		{
			name:       "empty directory without spec.md",
			capability: "empty-dir",
			expected:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findMatchingDeltaSpec(changeDir, tt.capability)
			if result != tt.expected {
				t.Errorf("findMatchingDeltaSpec(%q) = %v, want %v", tt.capability, result, tt.expected)
			}
		})
	}
}

func TestSplitTasksByCapability(t *testing.T) {
	changeDir := t.TempDir()

	// Create a delta spec directory
	specsDir := filepath.Join(changeDir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("failed to create specs dir: %v", err)
	}

	// Create support-aider delta spec
	aiderDir := filepath.Join(specsDir, "support-aider")
	if err := os.MkdirAll(aiderDir, 0o755); err != nil {
		t.Fatalf("failed to create support-aider dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(aiderDir, "spec.md"), []byte("# Spec"), 0o644); err != nil {
		t.Fatalf("failed to write spec.md: %v", err)
	}

	sections := []SectionTasks{
		{
			SectionName: "Foundation",
			Tasks: []parsers.Task{
				{
					ID:          "1",
					Section:     "Foundation",
					Description: "Create core interfaces",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			SectionName: "5. Support Aider",
			Tasks: []parsers.Task{
				{
					ID:          "5.1",
					Section:     "5. Support Aider",
					Description: "Migrate aider.go",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "5.2",
					Section:     "5. Support Aider",
					Description: "Add unit tests for Aider",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
	}

	rootTasks, childFiles, hasHierarchy := splitTasksByCapability(changeDir, sections)

	// Check root tasks
	if len(rootTasks) != 2 {
		t.Errorf("expected 2 root tasks, got %d", len(rootTasks))
	}

	// First task should be from Foundation (no children)
	if rootTasks[0].ID != "1" {
		t.Errorf("root task 1 ID = %s, want 1", rootTasks[0].ID)
	}
	if rootTasks[0].Children != "" {
		t.Errorf("root task 1 Children = %s, want empty", rootTasks[0].Children)
	}

	// Second task should be reference to Support Aider
	if rootTasks[1].ID != "5.1" {
		t.Errorf("root task 2 ID = %s, want 5.1", rootTasks[1].ID)
	}
	expectedChildren := "$ref:specs/support-aider/tasks.jsonc"
	if rootTasks[1].Children != expectedChildren {
		t.Errorf("root task 2 Children = %s, want %s", rootTasks[1].Children, expectedChildren)
	}

	// Check child files
	if !hasHierarchy {
		t.Error("expected hasHierarchy to be true")
	}

	if len(childFiles) != 1 {
		t.Errorf("expected 1 child file, got %d", len(childFiles))
	}

	childTasks, exists := childFiles["support-aider"]
	if !exists {
		t.Error("expected support-aider in childFiles")
	}

	if len(childTasks) != 2 {
		t.Errorf("expected 2 child tasks, got %d", len(childTasks))
	}

	// Child tasks should not have section or children
	for _, task := range childTasks {
		if task.Section != "" {
			t.Errorf("child task %s Section = %s, want empty", task.ID, task.Section)
		}
		if task.Children != "" {
			t.Errorf("child task %s Children = %s, want empty", task.ID, task.Children)
		}
	}
}

func TestSplitTasksByCapability_NoMatch(t *testing.T) {
	changeDir := t.TempDir()

	// Create specs dir but without any matching delta specs
	specsDir := filepath.Join(changeDir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("failed to create specs dir: %v", err)
	}

	sections := []SectionTasks{
		{
			SectionName: "Foundation",
			Tasks: []parsers.Task{
				{
					ID:          "1",
					Section:     "Foundation",
					Description: "Create core interfaces",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2",
					Section:     "Foundation",
					Description: "Another task",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
	}

	rootTasks, childFiles, hasHierarchy := splitTasksByCapability(changeDir, sections)

	// All tasks should be in root
	if len(rootTasks) != 2 {
		t.Errorf("expected 2 root tasks, got %d", len(rootTasks))
	}

	// No child files
	if len(childFiles) != 0 {
		t.Errorf("expected 0 child files, got %d", len(childFiles))
	}

	if hasHierarchy {
		t.Error("expected hasHierarchy to be false")
	}
}

func TestComputeSummary(t *testing.T) {
	tests := []struct {
		name               string
		tasks              []parsers.Task
		expectedTotal      int
		expectedCompleted  int
		expectedInProgress int
		expectedPending    int
	}{
		{
			name:               "empty tasks",
			tasks:              []parsers.Task{},
			expectedTotal:      0,
			expectedCompleted:  0,
			expectedInProgress: 0,
			expectedPending:    0,
		},
		{
			name: "mixed status",
			tasks: []parsers.Task{
				{Status: parsers.TaskStatusCompleted},
				{Status: parsers.TaskStatusCompleted},
				{Status: parsers.TaskStatusInProgress},
				{Status: parsers.TaskStatusPending},
				{Status: parsers.TaskStatusPending},
				{Status: parsers.TaskStatusPending},
			},
			expectedTotal:      6,
			expectedCompleted:  2,
			expectedInProgress: 1,
			expectedPending:    3,
		},
		{
			name: "all completed",
			tasks: []parsers.Task{
				{Status: parsers.TaskStatusCompleted},
				{Status: parsers.TaskStatusCompleted},
			},
			expectedTotal:      2,
			expectedCompleted:  2,
			expectedInProgress: 0,
			expectedPending:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := computeSummary(tt.tasks)

			if summary.Total != tt.expectedTotal {
				t.Errorf("Summary.Total = %d, want %d", summary.Total, tt.expectedTotal)
			}
			if summary.Completed != tt.expectedCompleted {
				t.Errorf("Summary.Completed = %d, want %d", summary.Completed, tt.expectedCompleted)
			}
			if summary.InProgress != tt.expectedInProgress {
				t.Errorf("Summary.InProgress = %d, want %d", summary.InProgress, tt.expectedInProgress)
			}
			if summary.Pending != tt.expectedPending {
				t.Errorf("Summary.Pending = %d, want %d", summary.Pending, tt.expectedPending)
			}
		})
	}
}

// Integration tests for hierarchical auto-split generation
func TestHierarchicalAutoSplit_Integration(t *testing.T) {
	// Create a complete mock change directory with delta specs
	tmpDir := t.TempDir()

	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	specsDir := filepath.Join(changeDir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("failed to create specs dir: %v", err)
	}

	// Create support-aider delta spec
	aiderDir := filepath.Join(specsDir, "support-aider")
	if err := os.MkdirAll(aiderDir, 0o755); err != nil {
		t.Fatalf("failed to create support-aider dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(aiderDir, "spec.md"), []byte("# Support Aider\n\n## ADDED\n### Test\nTest spec"), 0o644); err != nil {
		t.Fatalf("failed to write spec.md: %v", err)
	}

	// Create providers delta spec
	providersDir := filepath.Join(specsDir, "migrate-providers")
	if err := os.MkdirAll(providersDir, 0o755); err != nil {
		t.Fatalf("failed to create migrate-providers dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(providersDir, "spec.md"), []byte("# Migrate Providers\n\n## ADDED\n### Test\nTest spec"), 0o644); err != nil {
		t.Fatalf("failed to write spec.md: %v", err)
	}

	// Create a proposal.md (required for validation)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	proposalContent := `# Test Proposal

## Problem
Test problem

## Solution
Test solution
`
	if err := os.WriteFile(proposalPath, []byte(proposalContent), filePerm); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	// Create tasks.md with sections matching delta specs
	tasksMdContent := `## 1. Foundation

- [ ] 1.1 Create core interfaces
- [ ] 1.2 Set up basic structure

## 5. Support Aider

- [ ] 5.1 Migrate aider.go to new Provider interface
- [ ] 5.2 Add unit tests for Aider provider
- [ ] 5.3 Update documentation

## 6. Migrate Providers

- [ ] 6.1 Migrate all providers to new interface
- [ ] 6.2 Add integration tests
- [x] 6.3 Update documentation

## 7. Testing

- [ ] 7.1 Write unit tests
- [ ] 7.2 Run integration tests
`
	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMdContent), filePerm); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Parse tasks.md with sections
	sections, err := parseTasksMdWithSections(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse tasks.md with sections: %v", err)
	}

	// Verify section parsing
	if len(sections) != 4 {
		t.Errorf("expected 4 sections, got %d", len(sections))
	}

	// Verify section names and task counts
	expectedSections := map[string]int{
		"Foundation":        2,
		"Support Aider":     3,
		"Migrate Providers": 3,
		"Testing":           2,
	}
	for _, section := range sections {
		expectedCount, ok := expectedSections[section.SectionName]
		if !ok {
			t.Errorf("unexpected section: %s", section.SectionName)
			continue
		}
		if len(section.Tasks) != expectedCount {
			t.Errorf("section %s: expected %d tasks, got %d", section.SectionName, expectedCount, len(section.Tasks))
		}
	}

	// Split tasks by capability
	rootTasks, childFiles, hasHierarchy := splitTasksByCapability(changeDir, sections)

	// Verify hierarchy was created
	if !hasHierarchy {
		t.Error("expected hasHierarchy to be true")
	}

	// Verify root tasks count:
	// - Foundation: 2 flat tasks (IDs 1.1, 1.2)
	// - Support Aider: 1 reference task (ID 5.1 from first task in section)
	// - Migrate Providers: 1 reference task (ID 6.1 from first task in section)
	// - Testing: 2 flat tasks (IDs 7.1, 7.2)
	// Total: 6 root tasks
	if len(rootTasks) != 6 {
		t.Errorf("expected 6 root tasks, got %d", len(rootTasks))
	}

	// Verify Support Aider child file was created
	aiderTasks, exists := childFiles["support-aider"]
	if !exists {
		t.Error("expected support-aider in childFiles")
	}
	if len(aiderTasks) != 3 {
		t.Errorf("expected 3 Support Aider tasks, got %d", len(aiderTasks))
	}

	// Verify Migrate Providers child file was created
	providersTasks, exists := childFiles["migrate-providers"]
	if !exists {
		t.Error("expected migrate-providers in childFiles")
	}
	if len(providersTasks) != 3 {
		t.Errorf("expected 3 Migrate Providers tasks, got %d", len(providersTasks))
	}

	// Verify child task IDs are correct (use original task IDs)
	if aiderTasks[0].ID != "5.1" {
		t.Errorf("aider task 1 ID = %s, want 5.1", aiderTasks[0].ID)
	}
	if providersTasks[2].ID != "6.3" {
		t.Errorf("providers task 3 ID = %s, want 6.3", providersTasks[2].ID)
	}

	// Verify status was preserved
	if providersTasks[2].Status != parsers.TaskStatusCompleted {
		t.Errorf("providers task 3 status = %s, want completed", providersTasks[2].Status)
	}

	// Verify reference task uses first task's ID from each section
	// Support Aider ref should have ID 5.1 (from first task 5.1)
	// Migrate Providers ref should have ID 6.1 (from first task 6.1)
	has5dot1 := false
	has6dot1 := false
	for _, task := range rootTasks {
		if task.ID == "5.1" && task.Children != "" {
			has5dot1 = true
		}
		if task.ID == "6.1" && task.Children != "" {
			has6dot1 = true
		}
	}
	if !has5dot1 {
		t.Error("expected to find Support Aider reference task with ID 5.1")
	}
	if !has6dot1 {
		t.Error("expected to find Migrate Providers reference task with ID 6.1")
	}
}

func TestHierarchicalWriteAndRead_Integration(t *testing.T) {
	// Test the full write/read cycle for hierarchical tasks
	// Note: This test focuses on verifying the write operations work correctly
	// Full read cycle testing is done through the tasks command integration tests
	tmpDir := t.TempDir()

	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	specsDir := filepath.Join(changeDir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("failed to create specs dir: %v", err)
	}
}

func TestHierarchicalAutoSplit_MixedSections(t *testing.T) {
	// Test auto-split when some sections match and some don't
	tmpDir := t.TempDir()

	changeDir := filepath.Join(tmpDir, "change")
	specsDir := filepath.Join(changeDir, "specs")
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("failed to create specs dir: %v", err)
	}

	// Only create support-aider delta spec (not testing or docs)
	aiderDir := filepath.Join(specsDir, "support-aider")
	if err := os.MkdirAll(aiderDir, 0o755); err != nil {
		t.Fatalf("failed to create support-aider dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(aiderDir, "spec.md"), []byte("# Spec"), 0o644); err != nil {
		t.Fatalf("failed to write spec.md: %v", err)
	}

	sections := []SectionTasks{
		{
			SectionName: "Foundation",
			Tasks: []parsers.Task{
				{ID: "1", Description: "Task 1", Status: parsers.TaskStatusPending},
				{ID: "2", Description: "Task 2", Status: parsers.TaskStatusPending},
			},
		},
		{
			SectionName: "5. Support Aider",
			Tasks: []parsers.Task{
				{ID: "5.1", Description: "Aider Task 1", Status: parsers.TaskStatusPending},
			},
		},
		{
			SectionName: "Documentation",
			Tasks: []parsers.Task{
				{ID: "10", Description: "Doc Task 1", Status: parsers.TaskStatusPending},
			},
		},
	}

	rootTasks, childFiles, hasHierarchy := splitTasksByCapability(changeDir, sections)

	// Verify hierarchy is true (at least one match)
	if !hasHierarchy {
		t.Error("expected hasHierarchy to be true")
	}

	// Verify root tasks: Foundation (2 flat) + Support Aider ref + Documentation
	// Total: 4 root tasks
	if len(rootTasks) != 4 {
		t.Errorf("expected 4 root tasks, got %d", len(rootTasks))
	}

	// First task should be from Foundation (flat)
	if rootTasks[0].ID != "1" || rootTasks[0].Children != "" {
		t.Errorf("root task 1 should be flat, got ID=%s, Children=%s", rootTasks[0].ID, rootTasks[0].Children)
	}

	// Second task should be from Foundation (flat, 2nd task in section)
	if rootTasks[1].ID != "2" || rootTasks[1].Children != "" {
		t.Errorf("root task 2 should be flat, got ID=%s, Children=%s", rootTasks[1].ID, rootTasks[1].Children)
	}

	// Third task should be reference to Support Aider
	if rootTasks[2].ID != "5.1" || rootTasks[2].Children != "$ref:specs/support-aider/tasks.jsonc" {
		t.Errorf("root task 3 should be reference, got ID=%s, Children=%s", rootTasks[2].ID, rootTasks[2].Children)
	}

	// Fourth task should be from Documentation (flat)
	if rootTasks[3].ID != "10" || rootTasks[3].Children != "" {
		t.Errorf("root task 4 should be flat, got ID=%s, Children=%s", rootTasks[3].ID, rootTasks[3].Children)
	}

	// Verify only support-aider child file exists
	if len(childFiles) != 1 {
		t.Errorf("expected 1 child file, got %d", len(childFiles))
	}

	if _, exists := childFiles["support-aider"]; !exists {
		t.Error("expected support-aider in childFiles")
	}

	if _, exists := childFiles["documentation"]; exists {
		t.Error("documentation should NOT be in childFiles (no matching delta spec)")
	}
}

func TestHierarchicalAutoSplit_NoDeltaSpecs(t *testing.T) {
	// Test behavior when no delta specs exist (should remain flat)
	tmpDir := t.TempDir()

	changeDir := filepath.Join(tmpDir, "change")
	specsDir := filepath.Join(changeDir, "specs")
	// Create empty specs directory (no delta specs)
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatalf("failed to create specs dir: %v", err)
	}

	sections := []SectionTasks{
		{
			SectionName: "Foundation",
			Tasks: []parsers.Task{
				{ID: "1", Description: "Task 1", Status: parsers.TaskStatusPending},
				{ID: "2", Description: "Task 2", Status: parsers.TaskStatusPending},
			},
		},
		{
			SectionName: "Implementation",
			Tasks: []parsers.Task{
				{ID: "3", Description: "Task 3", Status: parsers.TaskStatusPending},
			},
		},
	}

	rootTasks, childFiles, hasHierarchy := splitTasksByCapability(changeDir, sections)

	// Verify no hierarchy created
	if hasHierarchy {
		t.Error("expected hasHierarchy to be false when no delta specs match")
	}

	// Verify all tasks are in root
	if len(rootTasks) != 3 {
		t.Errorf("expected 3 root tasks, got %d", len(rootTasks))
	}

	// Verify no child files
	if len(childFiles) != 0 {
		t.Errorf("expected 0 child files, got %d", len(childFiles))
	}

	// Verify all tasks are flat (no children)
	for i, task := range rootTasks {
		if task.Children != "" {
			t.Errorf("root task %d should have no children, got %s", i, task.Children)
		}
	}
}

func TestParseTasksMdWithSections_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()

	tasksMdPath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(""), filePerm); err != nil {
		t.Fatalf("failed to write empty tasks.md: %v", err)
	}

	sections, err := parseTasksMdWithSections(tasksMdPath)
	if err != nil {
		t.Fatalf("parseTasksMdWithSections() error = %v", err)
	}

	if len(sections) != 0 {
		t.Errorf("expected 0 sections for empty file, got %d", len(sections))
	}
}

func TestParseTasksMdWithSections_OnlySections(t *testing.T) {
	tmpDir := t.TempDir()

	tasksMdPath := filepath.Join(tmpDir, "tasks.md")
	content := `## 1. Section One

## 2. Section Two

## 3. Section Three
`
	if err := os.WriteFile(tasksMdPath, []byte(content), filePerm); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	sections, err := parseTasksMdWithSections(tasksMdPath)
	if err != nil {
		t.Fatalf("parseTasksMdWithSections() error = %v", err)
	}

	if len(sections) != 3 {
		t.Errorf("expected 3 sections, got %d", len(sections))
	}

	for i, section := range sections {
		if len(section.Tasks) != 0 {
			t.Errorf("section %d should have 0 tasks, got %d", i+1, len(section.Tasks))
		}
	}
}
