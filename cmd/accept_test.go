package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/config"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

// testProposalContent is a common proposal content used across multiple tests
const testProposalContent = `# Test Proposal

## Problem
Test problem

## Solution
Test solution
`

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
		{
			name: "multi-line task descriptions with sub-bullets",
			markdown: `## 3. Property-Based Testing

- [ ] 3.1 Create test infrastructure
- [ ] 3.2 Implement TestJSONCValidation_SpecialCharacters with test cases for:
  - Backslash ` + "`\\`" + `
  - Quote ` + "`\"`" + `
  - Newline ` + "`\\n`" + `
  - Tab ` + "`\\t`" + `
- [ ] 3.3 Implement TestJSONCValidation_Unicode with test cases for:
  - Emoji (🚀, 💻, ✅)
  - Non-ASCII characters (你好, مرحبا)
- [ ] 3.4 Simple task without sub-bullets
`,
			expected: []parsers.Task{
				{
					ID:          "3.1",
					Section:     "Property-Based Testing",
					Description: "Create test infrastructure",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "3.2",
					Section:     "Property-Based Testing",
					Description: "Implement TestJSONCValidation_SpecialCharacters with test cases for:\n  - Backslash `\\`\n  - Quote `\"`\n  - Newline `\\n`\n  - Tab `\\t`",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "3.3",
					Section:     "Property-Based Testing",
					Description: "Implement TestJSONCValidation_Unicode with test cases for:\n  - Emoji (🚀, 💻, ✅)\n  - Non-ASCII characters (你好, مرحبا)",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "3.4",
					Section:     "Property-Based Testing",
					Description: "Simple task without sub-bullets",
					Status:      parsers.TaskStatusPending,
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
			if err := writeTasksJSONC(tasksJSONPath, tt.tasks, nil, nil); err != nil {
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

	if err := writeTasksJSONC(tasksJSONPath, tasks, nil, nil); err != nil {
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

	err := writeTasksJSONC(tasksJsonPath, tasks, nil, nil)
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
	if err := os.WriteFile(proposalPath, []byte(testProposalContent), filePerm); err != nil {
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
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
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
	if err := os.WriteFile(proposalPath, []byte(testProposalContent), filePerm); err != nil {
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

	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
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

// TestWriteTasksJSONCWithAppendConfig verifies that append tasks are added correctly
func TestWriteTasksJSONCWithAppendConfig(t *testing.T) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(tmpDir, "tasks.jsonc")

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

	err := writeTasksJSONC(tasksJSONPath, existingTasks, appendCfg, nil)
	if err != nil {
		t.Fatalf("writeTasksJSONC() error = %v", err)
	}

	// Read and parse the output
	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	strippedData := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Should have 5 tasks (3 existing + 2 appended)
	if len(tasksFile.Tasks) != 5 {
		t.Errorf("expected 5 tasks, got %d", len(tasksFile.Tasks))
	}

	// Verify appended tasks have correct IDs (section 3, since max was 2)
	if tasksFile.Tasks[3].ID != "3.1" {
		t.Errorf("expected appended task ID '3.1', got '%s'", tasksFile.Tasks[3].ID)
	}
	if tasksFile.Tasks[4].ID != "3.2" {
		t.Errorf("expected appended task ID '3.2', got '%s'", tasksFile.Tasks[4].ID)
	}

	// Verify appended tasks have correct section
	if tasksFile.Tasks[3].Section != "Project Workflow" {
		t.Errorf("expected section 'Project Workflow', got '%s'", tasksFile.Tasks[3].Section)
	}

	// Verify appended tasks have pending status
	if tasksFile.Tasks[3].Status != parsers.TaskStatusPending {
		t.Errorf("expected status 'pending', got '%s'", tasksFile.Tasks[3].Status)
	}

	// Verify appended task descriptions
	if tasksFile.Tasks[3].Description != "Run linter and tests" {
		t.Error("wrong description for appended task")
	}
}

// TestAssignHierarchicalID tests the hierarchical ID generation function
func TestAssignHierarchicalID(t *testing.T) {
	tests := []struct {
		name        string
		existingID  string
		parentID    string
		childIndex  int
		expectedID  string
		description string
	}{
		{
			name:        "root task with no existing ID",
			existingID:  "",
			parentID:    "",
			childIndex:  1,
			expectedID:  "1",
			description: "generates simple numeric ID for root tasks",
		},
		{
			name:        "root task with existing ID",
			existingID:  "5",
			parentID:    "",
			childIndex:  1,
			expectedID:  "5",
			description: "preserves existing root task ID",
		},
		{
			name:        "child task without existing ID",
			existingID:  "",
			parentID:    "5",
			childIndex:  1,
			expectedID:  "5.1",
			description: "generates hierarchical ID for first child",
		},
		{
			name:        "child task with matching existing ID",
			existingID:  "5.2",
			parentID:    "5",
			childIndex:  2,
			expectedID:  "5.2",
			description: "preserves existing child ID when it matches parent",
		},
		{
			name:        "child task with non-matching existing ID",
			existingID:  "3.1",
			parentID:    "5",
			childIndex:  1,
			expectedID:  "5.1",
			description: "replaces existing ID when parent doesn't match",
		},
		{
			name:        "nested child task",
			existingID:  "",
			parentID:    "5.1",
			childIndex:  1,
			expectedID:  "5.1.1",
			description: "generates deeply nested hierarchical ID",
		},
		{
			name:        "nested child with existing ID",
			existingID:  "5.1.2",
			parentID:    "5.1",
			childIndex:  2,
			expectedID:  "5.1.2",
			description: "preserves existing nested child ID",
		},
		{
			name:        "multiple digit child index",
			existingID:  "",
			parentID:    "2",
			childIndex:  15,
			expectedID:  "2.15",
			description: "handles multi-digit child indices",
		},
		{
			name:        "root task with decimal existing ID",
			existingID:  "2.5",
			parentID:    "",
			childIndex:  1,
			expectedID:  "2.5",
			description: "preserves decimal root task ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := assignHierarchicalID(tt.existingID, tt.parentID, tt.childIndex)
			if result != tt.expectedID {
				t.Errorf(
					"assignHierarchicalID(%q, %q, %d) = %q, want %q\nDescription: %s",
					tt.existingID,
					tt.parentID,
					tt.childIndex,
					result,
					tt.expectedID,
					tt.description,
				)
			}
		})
	}
}

// TestValidateIDUniqueness tests the ID uniqueness validation function
func TestValidateIDUniqueness(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []parsers.Task
		wantError bool
		errorMsg  string
	}{
		{
			name: "all unique IDs",
			tasks: []parsers.Task{
				{ID: "1", Section: "A", Description: "Task 1", Status: parsers.TaskStatusPending},
				{ID: "2", Section: "A", Description: "Task 2", Status: parsers.TaskStatusPending},
				{ID: "3", Section: "B", Description: "Task 3", Status: parsers.TaskStatusPending},
			},
			wantError: false,
		},
		{
			name: "duplicate IDs",
			tasks: []parsers.Task{
				{ID: "1", Section: "A", Description: "Task 1", Status: parsers.TaskStatusPending},
				{ID: "2", Section: "A", Description: "Task 2", Status: parsers.TaskStatusPending},
				{
					ID:          "1",
					Section:     "B",
					Description: "Duplicate Task 1",
					Status:      parsers.TaskStatusPending,
				},
			},
			wantError: true,
			errorMsg:  "duplicate task IDs found: [1]",
		},
		{
			name: "multiple duplicate IDs",
			tasks: []parsers.Task{
				{ID: "1", Section: "A", Description: "Task 1", Status: parsers.TaskStatusPending},
				{ID: "2", Section: "A", Description: "Task 2", Status: parsers.TaskStatusPending},
				{
					ID:          "1",
					Section:     "B",
					Description: "Duplicate Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2",
					Section:     "C",
					Description: "Duplicate Task 2",
					Status:      parsers.TaskStatusPending,
				},
			},
			wantError: true,
			errorMsg:  "duplicate task IDs found:",
		},
		{
			name: "hierarchical IDs all unique",
			tasks: []parsers.Task{
				{ID: "1", Section: "A", Description: "Task 1", Status: parsers.TaskStatusPending},
				{
					ID:          "1.1",
					Section:     "A",
					Description: "Task 1.1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "A",
					Description: "Task 1.2",
					Status:      parsers.TaskStatusPending,
				},
				{ID: "2", Section: "B", Description: "Task 2", Status: parsers.TaskStatusPending},
				{
					ID:          "2.1",
					Section:     "B",
					Description: "Task 2.1",
					Status:      parsers.TaskStatusPending,
				},
			},
			wantError: false,
		},
		{
			name: "duplicate hierarchical IDs",
			tasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "A",
					Description: "Task 1.1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "A",
					Description: "Task 1.2",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.1",
					Section:     "B",
					Description: "Duplicate 1.1",
					Status:      parsers.TaskStatusPending,
				},
			},
			wantError: true,
			errorMsg:  "duplicate task IDs found: [1.1]",
		},
		{
			name:      "empty task list",
			tasks:     nil,
			wantError: false,
		},
		{
			name: "single task",
			tasks: []parsers.Task{
				{
					ID:          "1",
					Section:     "A",
					Description: "Single task",
					Status:      parsers.TaskStatusPending,
				},
			},
			wantError: false,
		},
		{
			name: "triple duplicate ID",
			tasks: []parsers.Task{
				{ID: "1", Section: "A", Description: "Task 1", Status: parsers.TaskStatusPending},
				{
					ID:          "1",
					Section:     "B",
					Description: "Duplicate 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1",
					Section:     "C",
					Description: "Duplicate 1 again",
					Status:      parsers.TaskStatusPending,
				},
			},
			wantError: true,
			errorMsg:  "duplicate task IDs found: [1]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIDUniqueness(tt.tasks)

			if tt.wantError {
				if err == nil {
					t.Error("validateIDUniqueness() expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("validateIDUniqueness() error = %q, want to contain %q", err.Error(), tt.errorMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateIDUniqueness() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestValidateIDUniquenessEdgeCases tests edge cases for ID validation
func TestValidateIDUniquenessEdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		tasks []parsers.Task
		want  error
	}{
		{
			name: "IDs that differ only in case (should be treated as different)",
			tasks: []parsers.Task{
				{ID: "a", Section: "A", Description: "Task a", Status: parsers.TaskStatusPending},
				{ID: "A", Section: "A", Description: "Task A", Status: parsers.TaskStatusPending},
			},
			want: nil, // Case-sensitive comparison
		},
		{
			name: "IDs with leading/trailing spaces (should be exact match)",
			tasks: []parsers.Task{
				{ID: "1", Section: "A", Description: "Task 1", Status: parsers.TaskStatusPending},
				{
					ID:          " 1",
					Section:     "A",
					Description: "Task with space",
					Status:      parsers.TaskStatusPending,
				},
			},
			want: nil, // Different IDs due to space
		},
		{
			name: "similar hierarchical IDs (1 vs 1.1)",
			tasks: []parsers.Task{
				{ID: "1", Section: "A", Description: "Task 1", Status: parsers.TaskStatusPending},
				{
					ID:          "1.1",
					Section:     "A",
					Description: "Task 1.1",
					Status:      parsers.TaskStatusPending,
				},
			},
			want: nil, // Different IDs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIDUniqueness(tt.tasks)
			if (err != nil) != (tt.want != nil) {
				t.Errorf("validateIDUniqueness() error = %v, want %v", err, tt.want)
			}
		})
	}
}

// TestWriteTasksJSONCWithDefaultSection verifies default section name
func TestWriteTasksJSONCWithDefaultSection(t *testing.T) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(tmpDir, "tasks.jsonc")

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

	err := writeTasksJSONC(tasksJSONPath, existingTasks, appendCfg, nil)
	if err != nil {
		t.Fatalf("writeTasksJSONC() error = %v", err)
	}

	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	strippedData := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Verify default section name is used
	if tasksFile.Tasks[1].Section != config.DefaultAppendTasksSection {
		t.Errorf("expected default section '%s', got '%s'",
			config.DefaultAppendTasksSection, tasksFile.Tasks[1].Section)
	}
}

// TestWriteTasksJSONCWithEmptyAppendTasks verifies no change when append tasks is empty
func TestWriteTasksJSONCWithEmptyAppendTasks(t *testing.T) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(tmpDir, "tasks.jsonc")

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

	err := writeTasksJSONC(tasksJSONPath, existingTasks, appendCfg, nil)
	if err != nil {
		t.Fatalf("writeTasksJSONC() error = %v", err)
	}

	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	strippedData := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	// Should have only the original task
	if len(tasksFile.Tasks) != 1 {
		t.Errorf("expected 1 task, got %d", len(tasksFile.Tasks))
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
				t.Errorf("findNextSectionNumber() = %d, want %d", got, tt.expected)
			}
		})
	}
}

// TestCreateAppendedTasks verifies task creation from config
func TestCreateAppendedTasks(t *testing.T) {
	existingTasks := []parsers.Task{
		{ID: "1.1", Section: "Setup", Description: "Task 1", Status: parsers.TaskStatusPending},
		{ID: "2.1", Section: "Impl", Description: "Task 2", Status: parsers.TaskStatusPending},
	}

	cfg := &config.AppendTasksConfig{
		Section: "Workflow",
		Tasks:   []string{"Task A", "Task B", "Task C"},
	}

	tasks := createAppendedTasks(existingTasks, cfg)

	if len(tasks) != 3 {
		t.Fatalf("expected 3 tasks, got %d", len(tasks))
	}

	// Verify IDs start at section 3
	expectedIDs := []string{"3.1", "3.2", "3.3"}
	for i, task := range tasks {
		if task.ID != expectedIDs[i] {
			t.Errorf("task %d: expected ID '%s', got '%s'", i, expectedIDs[i], task.ID)
		}
		if task.Section != "Workflow" {
			t.Errorf("task %d: expected section 'Workflow', got '%s'", i, task.Section)
		}
		if task.Status != parsers.TaskStatusPending {
			t.Errorf("task %d: expected status 'pending', got '%s'", i, task.Status)
		}
	}

	if tasks[0].Description != "Task A" {
		t.Errorf("expected description 'Task A', got '%s'", tasks[0].Description)
	}
}

// TestCountLines verifies line counting for various file sizes
func TestCountLines(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedLines int
	}{
		{
			name:          "empty file",
			content:       "",
			expectedLines: 0,
		},
		{
			name:          "single line",
			content:       "single line",
			expectedLines: 1,
		},
		{
			name:          "single line with newline",
			content:       "single line\n",
			expectedLines: 1,
		},
		{
			name:          "multiple lines",
			content:       "line 1\nline 2\nline 3\n",
			expectedLines: 3,
		},
		{
			name:          "50 lines",
			content:       generateLines(50),
			expectedLines: 50,
		},
		{
			name:          "100 lines exactly",
			content:       generateLines(100),
			expectedLines: 100,
		},
		{
			name:          "101 lines",
			content:       generateLines(101),
			expectedLines: 101,
		},
		{
			name:          "150 lines",
			content:       generateLines(150),
			expectedLines: 150,
		},
		{
			name:          "lines with varying content",
			content:       "# Header\n\n- [ ] Task 1\n- [ ] Task 2\n\n## Section\n",
			expectedLines: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.md")

			if err := os.WriteFile(testFile, []byte(tt.content), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := countLines(testFile)
			if err != nil {
				t.Fatalf("countLines() error = %v", err)
			}

			if got != tt.expectedLines {
				t.Errorf("countLines() = %d, want %d", got, tt.expectedLines)
			}
		})
	}
}

// TestCountLinesFileNotFound verifies error handling for missing files
func TestCountLinesFileNotFound(t *testing.T) {
	_, err := countLines("/nonexistent/path/file.md")
	if err == nil {
		t.Error("countLines() expected error for nonexistent file, got nil")
	}
}

// TestShouldSplit verifies split detection threshold
func TestShouldSplit(t *testing.T) {
	tests := []struct {
		name        string
		lineCount   int
		shouldSplit bool
		description string
	}{
		{
			name:        "empty file",
			lineCount:   0,
			shouldSplit: false,
			description: "0 lines should not trigger split",
		},
		{
			name:        "small file",
			lineCount:   50,
			shouldSplit: false,
			description: "50 lines should not trigger split",
		},
		{
			name:        "exactly at threshold",
			lineCount:   100,
			shouldSplit: false,
			description: "100 lines exactly should not trigger split",
		},
		{
			name:        "one line over threshold",
			lineCount:   101,
			shouldSplit: true,
			description: "101 lines should trigger split",
		},
		{
			name:        "well over threshold",
			lineCount:   150,
			shouldSplit: true,
			description: "150 lines should trigger split",
		},
		{
			name:        "large file",
			lineCount:   500,
			shouldSplit: true,
			description: "500 lines should trigger split",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "tasks.md")

			// Generate file with specified line count
			content := generateLines(tt.lineCount)
			if err := os.WriteFile(testFile, []byte(content), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := shouldSplit(testFile)
			if err != nil {
				t.Fatalf("shouldSplit() error = %v", err)
			}

			if got != tt.shouldSplit {
				t.Errorf("shouldSplit() = %v, want %v (%s)", got, tt.shouldSplit, tt.description)
			}
		})
	}
}

// TestShouldSplitFileNotFound verifies error handling
func TestShouldSplitFileNotFound(t *testing.T) {
	_, err := shouldSplit("/nonexistent/path/tasks.md")
	if err == nil {
		t.Error("shouldSplit() expected error for nonexistent file, got nil")
	}
}

// generateLines creates a string with the specified number of lines.
// Each line contains "Line N" where N is the line number.
func generateLines(count int) string {
	if count == 0 {
		return ""
	}

	var lines []string
	for i := 1; i <= count; i++ {
		lines = append(lines, fmt.Sprintf("Line %d", i))
	}

	return strings.Join(lines, "\n") + "\n"
}

// TestParseSections verifies section parsing with multiple sections
func TestParseSections(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		expected []Section
	}{
		{
			name: "multiple sections with tasks",
			markdown: `## 1. Setup
- [ ] 1.1 First setup task
- [ ] 1.2 Second setup task

## 2. Implementation
- [ ] 2.1 First impl task
- [x] 2.2 Second impl task

## 3. Testing
- [ ] 3.1 Test task
`,
			expected: []Section{
				{
					Name:      "Setup",
					Number:    "1",
					StartLine: 1,
					EndLine:   4, // Includes blank line before next section
					Tasks: []parsers.Task{
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
					},
				},
				{
					Name:      "Implementation",
					Number:    "2",
					StartLine: 5,
					EndLine:   8, // Includes blank line before next section
					Tasks: []parsers.Task{
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
					Name:      "Testing",
					Number:    "3",
					StartLine: 9,
					EndLine:   10, // Includes task line
					Tasks: []parsers.Task{
						{
							ID:          "3.1",
							Section:     "Testing",
							Description: "Test task",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name: "unnumbered sections",
			markdown: `## Setup
- [ ] Setup task

## Implementation
- [ ] Impl task
`,
			expected: []Section{
				{
					Name:      "Setup",
					Number:    "1",
					StartLine: 1,
					EndLine:   3, // Includes blank line
					Tasks: []parsers.Task{
						{
							ID:          "1.1",
							Section:     "Setup",
							Description: "Setup task",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
				{
					Name:      "Implementation",
					Number:    "2",
					StartLine: 4,
					EndLine:   5,
					Tasks: []parsers.Task{
						{
							ID:          "2.1",
							Section:     "Implementation",
							Description: "Impl task",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name: "sections with blank lines",
			markdown: `## 1. First Section

- [ ] 1.1 Task one


- [ ] 1.2 Task two

## 2. Second Section

- [ ] 2.1 Task three
`,
			expected: []Section{
				{
					Name:      "First Section",
					Number:    "1",
					StartLine: 1,
					EndLine:   7, // Includes all lines until next section
					Tasks: []parsers.Task{
						{
							ID:          "1.1",
							Section:     "First Section",
							Description: "Task one",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "1.2",
							Section:     "First Section",
							Description: "Task two",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
				{
					Name:      "Second Section",
					Number:    "2",
					StartLine: 8,
					EndLine:   10,
					Tasks: []parsers.Task{
						{
							ID:          "2.1",
							Section:     "Second Section",
							Description: "Task three",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name: "section with no tasks",
			markdown: `## 1. Empty Section

## 2. Section With Tasks
- [ ] 2.1 Task
`,
			expected: []Section{
				{
					Name:      "Empty Section",
					Number:    "1",
					StartLine: 1,
					EndLine:   2,
					Tasks:     make([]parsers.Task, 0),
				},
				{
					Name:      "Section With Tasks",
					Number:    "2",
					StartLine: 3,
					EndLine:   4,
					Tasks: []parsers.Task{
						{
							ID:          "2.1",
							Section:     "Section With Tasks",
							Description: "Task",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name:     "empty file",
			markdown: "",
			expected: make([]Section, 0),
		},
		{
			name: "only sections no tasks",
			markdown: `## 1. Section One

## 2. Section Two
`,
			expected: []Section{
				{
					Name:      "Section One",
					Number:    "1",
					StartLine: 1,
					EndLine:   2,
					Tasks:     make([]parsers.Task, 0),
				},
				{
					Name:      "Section Two",
					Number:    "2",
					StartLine: 3,
					EndLine:   3,
					Tasks:     make([]parsers.Task, 0),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "tasks.md")

			if err := os.WriteFile(testFile, []byte(tt.markdown), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := parseSections(testFile)
			if err != nil {
				t.Fatalf("parseSections() error = %v", err)
			}

			if len(got) != len(tt.expected) {
				t.Fatalf(
					"parseSections() returned %d sections, want %d",
					len(got),
					len(tt.expected),
				)
			}

			for i, section := range got {
				exp := tt.expected[i]
				if section.Name != exp.Name {
					t.Errorf("section %d: Name = %q, want %q", i, section.Name, exp.Name)
				}
				if section.Number != exp.Number {
					t.Errorf("section %d: Number = %q, want %q", i, section.Number, exp.Number)
				}
				if section.StartLine != exp.StartLine {
					t.Errorf(
						"section %d: StartLine = %d, want %d",
						i,
						section.StartLine,
						exp.StartLine,
					)
				}
				if section.EndLine != exp.EndLine {
					t.Errorf("section %d: EndLine = %d, want %d", i, section.EndLine, exp.EndLine)
				}
				if !reflect.DeepEqual(section.Tasks, exp.Tasks) {
					t.Errorf(
						"section %d: Tasks mismatch\ngot:  %+v\nwant: %+v",
						i,
						section.Tasks,
						exp.Tasks,
					)
				}
			}
		})
	}
}

// TestParseSectionsMalformedHeaders verifies handling of malformed section headers
func TestParseSectionsMalformedHeaders(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		description string
		expected    int // expected number of sections
	}{
		{
			name: "single hash header ignored",
			markdown: `# Not a section
- [ ] Task without section
## 1. Real Section
- [ ] 1.1 Task in section
`,
			description: "single # headers should not be treated as sections",
			expected:    1, // Only "Real Section"
		},
		{
			name: "triple hash header ignored",
			markdown: `### Not a section
- [ ] Task without section
## 1. Real Section
- [ ] 1.1 Task in section
`,
			description: "triple ### headers should not be treated as sections",
			expected:    1,
		},
		{
			name: "header without space after hash",
			markdown: `##NoSpace
- [ ] Task
## 1. Real Section
- [ ] 1.1 Another task
`,
			description: "## without space might not be recognized",
			expected:    1, // Depends on markdown parser - expecting only "Real Section"
		},
		{
			name: "mixed valid and invalid headers",
			markdown: `## 1. Valid Section
- [ ] 1.1 Task one
# Invalid Header
## 2. Another Valid Section
- [ ] 2.1 Task two
### Also Invalid
`,
			description: "only valid ## headers should create sections",
			expected:    2,
		},
		{
			name: "section header with no content after",
			markdown: `## 1. Section One
- [ ] 1.1 Task
## 2. Section Two
`,
			description: "section with no tasks should still be parsed",
			expected:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "tasks.md")

			if err := os.WriteFile(testFile, []byte(tt.markdown), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := parseSections(testFile)
			if err != nil {
				t.Fatalf("parseSections() error = %v", err)
			}

			if len(got) != tt.expected {
				t.Errorf("parseSections() returned %d sections, want %d (%s)",
					len(got), tt.expected, tt.description)
			}
		})
	}
}

// TestParseSectionsTasksWithoutSections verifies handling of tasks without sections
func TestParseSectionsTasksWithoutSections(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		description string
		expected    int // expected number of sections
	}{
		{
			name: "tasks before first section",
			markdown: `- [ ] Orphan task one
- [ ] Orphan task two

## 1. First Section
- [ ] 1.1 Task in section
`,
			description: "tasks before first section should not create a section",
			expected:    1, // Only "First Section"
		},
		{
			name: "only tasks no sections",
			markdown: `- [ ] Task one
- [ ] Task two
- [x] Task three
`,
			description: "file with only tasks and no sections",
			expected:    0,
		},
		{
			name: "tasks between sections",
			markdown: `## 1. Section One
- [ ] 1.1 Task

- [ ] Orphan task between sections

## 2. Section Two
- [ ] 2.1 Task
`,
			description: "tasks between sections belong to previous section",
			expected:    2,
		},
		{
			name:        "empty file",
			markdown:    "",
			description: "empty file should return no sections",
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "tasks.md")

			if err := os.WriteFile(testFile, []byte(tt.markdown), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			got, err := parseSections(testFile)
			if err != nil {
				t.Fatalf("parseSections() error = %v", err)
			}

			if len(got) != tt.expected {
				t.Errorf("parseSections() returned %d sections, want %d (%s)",
					len(got), tt.expected, tt.description)
			}
		})
	}
}

// TestExtractTasksForSection verifies task extraction by section name
func TestExtractTasksForSection(t *testing.T) {
	allTasks := []parsers.Task{
		{ID: "1.1", Section: "Setup", Description: "Task 1", Status: parsers.TaskStatusPending},
		{ID: "1.2", Section: "Setup", Description: "Task 2", Status: parsers.TaskStatusPending},
		{
			ID:          "2.1",
			Section:     "Implementation",
			Description: "Task 3",
			Status:      parsers.TaskStatusPending,
		},
		{
			ID:          "2.2",
			Section:     "Implementation",
			Description: "Task 4",
			Status:      parsers.TaskStatusCompleted,
		},
		{ID: "3.1", Section: "Testing", Description: "Task 5", Status: parsers.TaskStatusPending},
	}

	tests := []struct {
		name        string
		sectionName string
		expected    []parsers.Task
	}{
		{
			name:        "extract setup tasks",
			sectionName: "Setup",
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Task 2",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name:        "extract implementation tasks",
			sectionName: "Implementation",
			expected: []parsers.Task{
				{
					ID:          "2.1",
					Section:     "Implementation",
					Description: "Task 3",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.2",
					Section:     "Implementation",
					Description: "Task 4",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			name:        "extract testing tasks",
			sectionName: "Testing",
			expected: []parsers.Task{
				{
					ID:          "3.1",
					Section:     "Testing",
					Description: "Task 5",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name:        "nonexistent section",
			sectionName: "Documentation",
			expected:    make([]parsers.Task, 0),
		},
		{
			name:        "empty section name",
			sectionName: "",
			expected:    make([]parsers.Task, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractTasksForSection(allTasks, tt.sectionName)

			if len(got) != len(tt.expected) {
				t.Fatalf("extractTasksForSection() returned %d tasks, want %d",
					len(got), len(tt.expected))
			}

			// For empty slices, just checking length is sufficient
			if len(tt.expected) > 0 && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("extractTasksForSection() mismatch\ngot:  %+v\nwant: %+v",
					got, tt.expected)
			}
		})
	}
}

// TestSectionLineCount verifies line count calculation for sections
func TestSectionLineCount(t *testing.T) {
	tests := []struct {
		name     string
		section  Section
		expected int
	}{
		{
			name: "normal section",
			section: Section{
				StartLine: 1,
				EndLine:   10,
			},
			expected: 10,
		},
		{
			name: "single line section",
			section: Section{
				StartLine: 5,
				EndLine:   5,
			},
			expected: 1,
		},
		{
			name: "invalid section end before start",
			section: Section{
				StartLine: 10,
				EndLine:   5,
			},
			expected: 0,
		},
		{
			name: "large section",
			section: Section{
				StartLine: 1,
				EndLine:   150,
			},
			expected: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.section.LineCount()
			if got != tt.expected {
				t.Errorf("LineCount() = %d, want %d", got, tt.expected)
			}
		})
	}
}

// TestParseSectionsFileNotFound verifies error handling
func TestParseSectionsFileNotFound(t *testing.T) {
	_, err := parseSections("/nonexistent/path/tasks.md")
	if err == nil {
		t.Error("parseSections() expected error for nonexistent file, got nil")
	}
}

// TestShouldSplitSection verifies section splitting based on line count
func TestShouldSplitSection(t *testing.T) {
	tests := []struct {
		name        string
		section     Section
		shouldSplit bool
		description string
	}{
		{
			name: "small section under threshold",
			section: Section{
				Name:      "Small Section",
				StartLine: 1,
				EndLine:   50,
			},
			shouldSplit: false,
			description: "50 lines should not trigger section split",
		},
		{
			name: "section exactly at threshold",
			section: Section{
				Name:      "Boundary Section",
				StartLine: 1,
				EndLine:   100,
			},
			shouldSplit: false,
			description: "100 lines exactly should not trigger section split",
		},
		{
			name: "section one line over threshold",
			section: Section{
				Name:      "Just Over Threshold",
				StartLine: 1,
				EndLine:   101,
			},
			shouldSplit: true,
			description: "101 lines should trigger section split",
		},
		{
			name: "large section well over threshold",
			section: Section{
				Name:      "Large Section",
				StartLine: 1,
				EndLine:   150,
			},
			shouldSplit: true,
			description: "150 lines should trigger section split",
		},
		{
			name: "very large section",
			section: Section{
				Name:      "Very Large Section",
				StartLine: 10,
				EndLine:   310,
			},
			shouldSplit: true,
			description: "300+ lines should trigger section split",
		},
		{
			name: "empty section",
			section: Section{
				Name:      "Empty Section",
				StartLine: 1,
				EndLine:   1,
			},
			shouldSplit: false,
			description: "1 line should not trigger section split",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSplitSection(&tt.section)
			if got != tt.shouldSplit {
				t.Errorf("shouldSplitSection() = %v, want %v (%s)",
					got, tt.shouldSplit, tt.description)
			}
		})
	}
}

// TestExtractIDPrefix verifies ID prefix extraction for subsection grouping
func TestExtractIDPrefix(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected string
	}{
		{
			name:     "hierarchical ID with one dot",
			id:       "1.1",
			expected: "1",
		},
		{
			name:     "hierarchical ID with two dots",
			id:       "1.2.3",
			expected: "1.2",
		},
		{
			name:     "hierarchical ID with multiple dots",
			id:       "1.2.3.4.5",
			expected: "1.2.3.4",
		},
		{
			name:     "simple ID without dots",
			id:       "1",
			expected: "1",
		},
		{
			name:     "simple ID with higher number",
			id:       "5",
			expected: "5",
		},
		{
			name:     "ID with decimal format",
			id:       "2.1",
			expected: "2",
		},
		{
			name:     "ID with larger first part",
			id:       "10.5",
			expected: "10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractIDPrefix(tt.id)
			if got != tt.expected {
				t.Errorf("extractIDPrefix(%q) = %q, want %q", tt.id, got, tt.expected)
			}
		})
	}
}

// TestParseSubsections verifies subsection grouping by ID prefix
func TestParseSubsections(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []parsers.Task
		expected []SubsectionGroup
	}{
		{
			name: "tasks with common prefix",
			tasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Section",
					Description: "Task 1.1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Section",
					Description: "Task 1.2",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.3",
					Section:     "Section",
					Description: "Task 1.3",
					Status:      parsers.TaskStatusPending,
				},
			},
			expected: []SubsectionGroup{
				{
					Prefix: "1",
					Tasks: []parsers.Task{
						{
							ID:          "1.1",
							Section:     "Section",
							Description: "Task 1.1",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "1.2",
							Section:     "Section",
							Description: "Task 1.2",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "1.3",
							Section:     "Section",
							Description: "Task 1.3",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name: "tasks with multiple prefixes",
			tasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Section",
					Description: "Task 1.1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Section",
					Description: "Task 1.2",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Section",
					Description: "Task 2.1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.2",
					Section:     "Section",
					Description: "Task 2.2",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.3",
					Section:     "Section",
					Description: "Task 2.3",
					Status:      parsers.TaskStatusPending,
				},
			},
			expected: []SubsectionGroup{
				{
					Prefix: "1",
					Tasks: []parsers.Task{
						{
							ID:          "1.1",
							Section:     "Section",
							Description: "Task 1.1",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "1.2",
							Section:     "Section",
							Description: "Task 1.2",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
				{
					Prefix: "2",
					Tasks: []parsers.Task{
						{
							ID:          "2.1",
							Section:     "Section",
							Description: "Task 2.1",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "2.2",
							Section:     "Section",
							Description: "Task 2.2",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "2.3",
							Section:     "Section",
							Description: "Task 2.3",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name: "tasks with hierarchical IDs",
			tasks: []parsers.Task{
				{
					ID:          "1.1.1",
					Section:     "Section",
					Description: "Task 1.1.1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.1.2",
					Section:     "Section",
					Description: "Task 1.1.2",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2.1",
					Section:     "Section",
					Description: "Task 1.2.1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2.2",
					Section:     "Section",
					Description: "Task 1.2.2",
					Status:      parsers.TaskStatusPending,
				},
			},
			expected: []SubsectionGroup{
				{
					Prefix: "1.1",
					Tasks: []parsers.Task{
						{
							ID:          "1.1.1",
							Section:     "Section",
							Description: "Task 1.1.1",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "1.1.2",
							Section:     "Section",
							Description: "Task 1.1.2",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
				{
					Prefix: "1.2",
					Tasks: []parsers.Task{
						{
							ID:          "1.2.1",
							Section:     "Section",
							Description: "Task 1.2.1",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "1.2.2",
							Section:     "Section",
							Description: "Task 1.2.2",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name: "tasks with simple IDs no dots",
			tasks: []parsers.Task{
				{
					ID:          "1",
					Section:     "Section",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2",
					Section:     "Section",
					Description: "Task 2",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "3",
					Section:     "Section",
					Description: "Task 3",
					Status:      parsers.TaskStatusPending,
				},
			},
			expected: []SubsectionGroup{
				{
					Prefix: "1",
					Tasks: []parsers.Task{
						{
							ID:          "1",
							Section:     "Section",
							Description: "Task 1",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
				{
					Prefix: "2",
					Tasks: []parsers.Task{
						{
							ID:          "2",
							Section:     "Section",
							Description: "Task 2",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
				{
					Prefix: "3",
					Tasks: []parsers.Task{
						{
							ID:          "3",
							Section:     "Section",
							Description: "Task 3",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name: "single task",
			tasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Section",
					Description: "Task 1.1",
					Status:      parsers.TaskStatusPending,
				},
			},
			expected: []SubsectionGroup{
				{
					Prefix: "1",
					Tasks: []parsers.Task{
						{
							ID:          "1.1",
							Section:     "Section",
							Description: "Task 1.1",
							Status:      parsers.TaskStatusPending,
						},
					},
				},
			},
		},
		{
			name:     "empty task list",
			tasks:    make([]parsers.Task, 0),
			expected: make([]SubsectionGroup, 0),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseSubsections(tt.tasks)

			if len(got) != len(tt.expected) {
				t.Fatalf("parseSubsections() returned %d groups, want %d",
					len(got), len(tt.expected))
			}

			for i, group := range got {
				exp := tt.expected[i]
				if group.Prefix != exp.Prefix {
					t.Errorf("group %d: Prefix = %q, want %q", i, group.Prefix, exp.Prefix)
				}
				if !reflect.DeepEqual(group.Tasks, exp.Tasks) {
					t.Errorf("group %d: Tasks mismatch\ngot:  %+v\nwant: %+v",
						i, group.Tasks, exp.Tasks)
				}
			}
		})
	}
}

// TestParseSubsectionsFieldInitialization verifies that parseSubsections
// properly initializes all fields of SubsectionGroup structs.
// This test would have caught the bug where StartLine/EndLine fields
// were defined but never populated.
func TestParseSubsectionsFieldInitialization(t *testing.T) {
	tasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Test Section",
			Description: "Task 1",
			Status:      parsers.TaskStatusPending,
		},
		{
			ID:          "1.2",
			Section:     "Test Section",
			Description: "Task 2",
			Status:      parsers.TaskStatusPending,
		},
		{
			ID:          "2.1",
			Section:     "Test Section",
			Description: "Task 3",
			Status:      parsers.TaskStatusPending,
		},
	}

	groups := parseSubsections(tasks)

	if len(groups) != 2 {
		t.Fatalf("parseSubsections() returned %d groups, want 2", len(groups))
	}

	// Verify all fields are properly initialized for each group
	for i, group := range groups {
		// Prefix must be set (non-empty)
		if group.Prefix == "" {
			t.Errorf("group %d: Prefix is empty, should be initialized", i)
		}

		// Tasks must be initialized (non-nil) and contain tasks
		if group.Tasks == nil {
			t.Errorf("group %d: Tasks is nil, should be initialized", i)
		}
		if len(group.Tasks) == 0 {
			t.Errorf("group %d: Tasks is empty, should contain tasks", i)
		}

		// Verify Tasks slice contains actual task data
		for j, task := range group.Tasks {
			if task.ID == "" {
				t.Errorf("group %d, task %d: ID is empty", i, j)
			}
			if task.Description == "" {
				t.Errorf("group %d, task %d: Description is empty", i, j)
			}
		}
	}

	// Verify specific group contents
	if groups[0].Prefix != "1" {
		t.Errorf("group 0: Prefix = %q, want \"1\"", groups[0].Prefix)
	}
	if len(groups[0].Tasks) != 2 {
		t.Errorf("group 0: len(Tasks) = %d, want 2", len(groups[0].Tasks))
	}

	if groups[1].Prefix != "2" {
		t.Errorf("group 1: Prefix = %q, want \"2\"", groups[1].Prefix)
	}
	if len(groups[1].Tasks) != 1 {
		t.Errorf("group 1: len(Tasks) = %d, want 1", len(groups[1].Tasks))
	}
}

// TestLargeSectionSplitting verifies splitting large sections into multiple files
func TestLargeSectionSplitting(t *testing.T) {
	// Create a large section with multiple subsections
	var tasks []parsers.Task
	sectionName := "Large Implementation Section"

	// Create subsection 1 (tasks 1.1-1.30)
	for i := 1; i <= 30; i++ {
		tasks = append(tasks, parsers.Task{
			ID:          fmt.Sprintf("1.%d", i),
			Section:     sectionName,
			Description: fmt.Sprintf("Task 1.%d", i),
			Status:      parsers.TaskStatusPending,
		})
	}

	// Create subsection 2 (tasks 2.1-2.30)
	for i := 1; i <= 30; i++ {
		tasks = append(tasks, parsers.Task{
			ID:          fmt.Sprintf("2.%d", i),
			Section:     sectionName,
			Description: fmt.Sprintf("Task 2.%d", i),
			Status:      parsers.TaskStatusPending,
		})
	}

	// Create subsection 3 (tasks 3.1-3.30)
	for i := 1; i <= 30; i++ {
		tasks = append(tasks, parsers.Task{
			ID:          fmt.Sprintf("3.%d", i),
			Section:     sectionName,
			Description: fmt.Sprintf("Task 3.%d", i),
			Status:      parsers.TaskStatusPending,
		})
	}

	// Parse subsections
	groups := parseSubsections(tasks)

	// Verify we got 3 groups
	if len(groups) != 3 {
		t.Fatalf("parseSubsections() returned %d groups, want 3", len(groups))
	}

	// Verify each group has the correct prefix and task count
	expectedPrefixes := []string{"1", "2", "3"}
	for i, group := range groups {
		if group.Prefix != expectedPrefixes[i] {
			t.Errorf("group %d: Prefix = %q, want %q", i, group.Prefix, expectedPrefixes[i])
		}
		if len(group.Tasks) != 30 {
			t.Errorf("group %d: got %d tasks, want 30", i, len(group.Tasks))
		}
	}

	// Verify the first task in each group
	for i, group := range groups {
		firstTask := group.Tasks[0]
		expectedID := fmt.Sprintf("%s.1", expectedPrefixes[i])
		if firstTask.ID != expectedID {
			t.Errorf("group %d first task: ID = %q, want %q", i, firstTask.ID, expectedID)
		}
	}

	// Verify the last task in each group
	for i, group := range groups {
		lastTask := group.Tasks[len(group.Tasks)-1]
		expectedID := fmt.Sprintf("%s.30", expectedPrefixes[i])
		if lastTask.ID != expectedID {
			t.Errorf("group %d last task: ID = %q, want %q", i, lastTask.ID, expectedID)
		}
	}
}

// TestLoadExistingStatuses verifies loading statuses from existing tasks.jsonc files
func TestLoadExistingStatuses(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    func(t *testing.T, changeDir string)
		expectedMap   map[string]parsers.TaskStatusValue
		expectError   bool
		errorContains string
	}{
		{
			name: "no existing files",
			setupFiles: func(_ *testing.T, _ string) {
				// Don't create any files
			},
			expectedMap: make(map[string]parsers.TaskStatusValue),
			expectError: false,
		},
		{
			name: "version 1 flat file",
			setupFiles: func(t *testing.T, changeDir string) {
				tasksFile := parsers.TasksFile{
					Version: 1,
					Tasks: []parsers.Task{
						{
							ID:          "1.1",
							Section:     "Setup",
							Description: "Task 1",
							Status:      parsers.TaskStatusCompleted,
						},
						{
							ID:          "1.2",
							Section:     "Setup",
							Description: "Task 2",
							Status:      parsers.TaskStatusInProgress,
						},
						{
							ID:          "2.1",
							Section:     "Impl",
							Description: "Task 3",
							Status:      parsers.TaskStatusPending,
						},
					},
				}
				data, _ := json.MarshalIndent(tasksFile, "", "  ")
				err := os.WriteFile(filepath.Join(changeDir, "tasks.jsonc"), data, 0o644)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			expectedMap: map[string]parsers.TaskStatusValue{
				"1.1": parsers.TaskStatusCompleted,
				"1.2": parsers.TaskStatusInProgress,
				"2.1": parsers.TaskStatusPending,
			},
			expectError: false,
		},
		{
			name: "version 2 with child files",
			setupFiles: func(t *testing.T, changeDir string) {
				// Create root file
				rootFile := parsers.TasksFile{
					Version: 2,
					Tasks: []parsers.Task{
						{
							ID:          "1",
							Section:     "Setup",
							Description: "Setup section",
							Status:      parsers.TaskStatusInProgress,
							Children:    "$ref:tasks-1.jsonc",
						},
						{
							ID:          "2",
							Section:     "Impl",
							Description: "Impl section",
							Status:      parsers.TaskStatusPending,
							Children:    "$ref:tasks-2.jsonc",
						},
					},
					Includes: []string{"tasks-*.jsonc"},
				}
				data, _ := json.MarshalIndent(rootFile, "", "  ")
				err := os.WriteFile(filepath.Join(changeDir, "tasks.jsonc"), data, 0o644)
				if err != nil {
					t.Fatalf("failed to write root file: %v", err)
				}

				// Create child file 1
				child1File := parsers.TasksFile{
					Version: 2,
					Parent:  "1",
					Tasks: []parsers.Task{
						{
							ID:          "1.1",
							Section:     "Setup",
							Description: "Task 1.1",
							Status:      parsers.TaskStatusCompleted,
						},
						{
							ID:          "1.2",
							Section:     "Setup",
							Description: "Task 1.2",
							Status:      parsers.TaskStatusInProgress,
						},
					},
				}
				data, _ = json.MarshalIndent(child1File, "", "  ")
				err = os.WriteFile(filepath.Join(changeDir, "tasks-1.jsonc"), data, 0o644)
				if err != nil {
					t.Fatalf("failed to write child file 1: %v", err)
				}

				// Create child file 2
				child2File := parsers.TasksFile{
					Version: 2,
					Parent:  "2",
					Tasks: []parsers.Task{
						{
							ID:          "2.1",
							Section:     "Impl",
							Description: "Task 2.1",
							Status:      parsers.TaskStatusPending,
						},
						{
							ID:          "2.2",
							Section:     "Impl",
							Description: "Task 2.2",
							Status:      parsers.TaskStatusPending,
						},
					},
				}
				data, _ = json.MarshalIndent(child2File, "", "  ")
				err = os.WriteFile(filepath.Join(changeDir, "tasks-2.jsonc"), data, 0o644)
				if err != nil {
					t.Fatalf("failed to write child file 2: %v", err)
				}
			},
			expectedMap: map[string]parsers.TaskStatusValue{
				"1":   parsers.TaskStatusInProgress,
				"2":   parsers.TaskStatusPending,
				"1.1": parsers.TaskStatusCompleted,
				"1.2": parsers.TaskStatusInProgress,
				"2.1": parsers.TaskStatusPending,
				"2.2": parsers.TaskStatusPending,
			},
			expectError: false,
		},
		{
			name: "invalid JSON in root file",
			setupFiles: func(t *testing.T, changeDir string) {
				err := os.WriteFile(
					filepath.Join(changeDir, "tasks.jsonc"),
					[]byte("invalid json"),
					0o644,
				)
				if err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			},
			expectedMap:   nil,
			expectError:   true,
			errorContains: "failed to parse tasks.jsonc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
			if err := os.MkdirAll(changeDir, 0o755); err != nil {
				t.Fatalf("failed to create change dir: %v", err)
			}

			tt.setupFiles(t, changeDir)

			statusMap, err := loadExistingStatuses(changeDir)

			if tt.expectError {
				if err == nil {
					t.Error("loadExistingStatuses() expected error, got nil")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("loadExistingStatuses() error = %q, want to contain %q", err.Error(), tt.errorContains)
				}

				return
			}

			if err != nil {
				t.Fatalf("loadExistingStatuses() unexpected error = %v", err)
			}

			if !reflect.DeepEqual(statusMap, tt.expectedMap) {
				t.Errorf(
					"loadExistingStatuses() mismatch\ngot:  %+v\nwant: %+v",
					statusMap,
					tt.expectedMap,
				)
			}
		})
	}
}

// TestMergeTaskStatuses verifies status merging during regeneration
func TestMergeTaskStatuses(t *testing.T) {
	tests := []struct {
		name      string
		newTasks  []parsers.Task
		statusMap map[string]parsers.TaskStatusValue
		expected  []parsers.Task
	}{
		{
			name: "all IDs match - preserve statuses",
			newTasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Task 2",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Impl",
					Description: "Task 3",
					Status:      parsers.TaskStatusPending,
				},
			},
			statusMap: map[string]parsers.TaskStatusValue{
				"1.1": parsers.TaskStatusCompleted,
				"1.2": parsers.TaskStatusInProgress,
				"2.1": parsers.TaskStatusCompleted,
			},
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusCompleted,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Task 2",
					Status:      parsers.TaskStatusInProgress,
				},
				{
					ID:          "2.1",
					Section:     "Impl",
					Description: "Task 3",
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			name: "some IDs match - partial preservation",
			newTasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "New task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Impl",
					Description: "Task 3",
					Status:      parsers.TaskStatusPending,
				},
			},
			statusMap: map[string]parsers.TaskStatusValue{
				"1.1": parsers.TaskStatusCompleted,
				"2.1": parsers.TaskStatusInProgress,
			},
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusCompleted,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "New task",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "2.1",
					Section:     "Impl",
					Description: "Task 3",
					Status:      parsers.TaskStatusInProgress,
				},
			},
		},
		{
			name: "no IDs match - all pending",
			newTasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Task 2",
					Status:      parsers.TaskStatusPending,
				},
			},
			statusMap: map[string]parsers.TaskStatusValue{
				"3.1": parsers.TaskStatusCompleted,
				"3.2": parsers.TaskStatusInProgress,
			},
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Task 2",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name: "empty status map - all pending",
			newTasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Task 2",
					Status:      parsers.TaskStatusPending,
				},
			},
			statusMap: make(map[string]parsers.TaskStatusValue),
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
				{
					ID:          "1.2",
					Section:     "Setup",
					Description: "Task 2",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
		{
			name:     "empty new tasks",
			newTasks: nil,
			statusMap: map[string]parsers.TaskStatusValue{
				"1.1": parsers.TaskStatusCompleted,
			},
			expected: nil,
		},
		{
			name: "status change from completed to in_progress",
			newTasks: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusPending,
				},
			},
			statusMap: map[string]parsers.TaskStatusValue{
				"1.1": parsers.TaskStatusInProgress,
			},
			expected: []parsers.Task{
				{
					ID:          "1.1",
					Section:     "Setup",
					Description: "Task 1",
					Status:      parsers.TaskStatusInProgress,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mergeTaskStatuses(tt.newTasks, tt.statusMap)

			// Handle nil vs empty slice comparison
			if (got == nil && tt.expected == nil) || (len(got) == 0 && len(tt.expected) == 0) {
				return
			}

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("mergeTaskStatuses() mismatch\ngot:  %+v\nwant: %+v", got, tt.expected)
			}
		})
	}
}

// TestStatusPreservationIntegration tests end-to-end status preservation during regeneration
func TestStatusPreservationIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create proposal.md (required for validation)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(testProposalContent), filePerm); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	// Step 1: Create initial tasks.md
	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	initialTasksMd := `## 1. Setup
- [ ] 1.1 Create project structure
- [ ] 1.2 Setup configuration

## 2. Implementation
- [ ] 2.1 Implement feature A
- [ ] 2.2 Implement feature B
`
	if err := os.WriteFile(tasksMdPath, []byte(initialTasksMd), filePerm); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Step 2: First accept - generate initial tasks.jsonc
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse tasks.md: %v", err)
	}

	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("first writeAndCleanup failed: %v", err)
	}

	// Step 3: Manually update some task statuses (simulating work being done)
	jsonContent, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	strippedJSON := parsers.StripJSONComments(jsonContent)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedJSON, &tasksFile); err != nil {
		t.Fatalf("failed to parse tasks.jsonc: %v", err)
	}

	// Mark some tasks as completed/in_progress
	tasksFile.Tasks[0].Status = parsers.TaskStatusCompleted  // 1.1 completed
	tasksFile.Tasks[1].Status = parsers.TaskStatusInProgress // 1.2 in progress
	tasksFile.Tasks[2].Status = parsers.TaskStatusCompleted  // 2.1 completed

	// Write back the modified file
	updatedData, _ := json.MarshalIndent(tasksFile, "", "  ")
	output := tasksJSONHeader + string(updatedData)
	if err := os.WriteFile(tasksJSONPath, []byte(output), filePerm); err != nil {
		t.Fatalf("failed to write updated tasks.jsonc: %v", err)
	}

	// Step 4: Update tasks.md (add a new task, modify description of existing task)
	updatedTasksMd := `## 1. Setup
- [ ] 1.1 Create project structure
- [ ] 1.2 Setup configuration
- [ ] 1.3 New setup task

## 2. Implementation
- [ ] 2.1 Implement feature A (updated description)
- [ ] 2.2 Implement feature B
`
	if err := os.WriteFile(tasksMdPath, []byte(updatedTasksMd), filePerm); err != nil {
		t.Fatalf("failed to write updated tasks.md: %v", err)
	}

	// Step 5: Re-run accept (regeneration)
	newTasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse updated tasks.md: %v", err)
	}

	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, newTasks, nil); err != nil {
		t.Fatalf("second writeAndCleanup failed: %v", err)
	}

	// Step 6: Verify statuses were preserved
	jsonContent, err = os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read regenerated tasks.jsonc: %v", err)
	}

	strippedJSON = parsers.StripJSONComments(jsonContent)
	var finalTasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedJSON, &finalTasksFile); err != nil {
		t.Fatalf("failed to parse regenerated tasks.jsonc: %v", err)
	}

	// Verify we now have 5 tasks (3 + 2 new ones)
	if len(finalTasksFile.Tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(finalTasksFile.Tasks))
	}

	// Verify preserved statuses
	expectedStatuses := map[string]parsers.TaskStatusValue{
		"1.1": parsers.TaskStatusCompleted,  // preserved
		"1.2": parsers.TaskStatusInProgress, // preserved
		"1.3": parsers.TaskStatusPending,    // new task
		"2.1": parsers.TaskStatusCompleted,  // preserved
		"2.2": parsers.TaskStatusPending,    // was pending, stays pending
	}

	for _, task := range finalTasksFile.Tasks {
		expectedStatus, ok := expectedStatuses[task.ID]
		if !ok {
			t.Errorf("unexpected task ID: %s", task.ID)

			continue
		}
		if task.Status != expectedStatus {
			t.Errorf("task %s: status = %s, want %s", task.ID, task.Status, expectedStatus)
		}
	}

	// Verify description was updated for task 2.1
	for _, task := range finalTasksFile.Tasks {
		if task.ID != "2.1" {
			continue
		}

		if !strings.Contains(task.Description, "updated description") {
			t.Errorf("task 2.1 description was not updated: %s", task.Description)
		}
	}
}

// TestWriteTasksJSONCNilStatusMap verifies handling of nil status map
func TestWriteTasksJSONCNilStatusMap(t *testing.T) {
	tmpDir := t.TempDir()
	tasksJSONPath := filepath.Join(tmpDir, "tasks.jsonc")

	tasks := []parsers.Task{
		{ID: "1.1", Section: "Setup", Description: "Task 1", Status: parsers.TaskStatusPending},
		{ID: "1.2", Section: "Setup", Description: "Task 2", Status: parsers.TaskStatusPending},
	}

	// Pass nil status map (should not panic, should treat as empty map)
	err := writeTasksJSONC(tasksJSONPath, tasks, nil, nil)
	if err != nil {
		t.Fatalf("writeTasksJSONC() error = %v", err)
	}

	// Read and verify
	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	strippedData := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if len(tasksFile.Tasks) != 2 {
		t.Errorf("expected 2 tasks, got %d", len(tasksFile.Tasks))
	}

	// Verify all tasks have pending status (no preservation)
	for _, task := range tasksFile.Tasks {
		if task.Status != parsers.TaskStatusPending {
			t.Errorf(
				"task %s: status = %s, want %s",
				task.ID,
				task.Status,
				parsers.TaskStatusPending,
			)
		}
	}
}

// Helper: verify child file is valid version 2 with parent
func verifyChildFile(t *testing.T, childPath, childName string) {
	t.Helper()

	childData, err := os.ReadFile(childPath)
	if err != nil {
		t.Fatalf("failed to read %s: %v", childName, err)
	}

	strippedChild := parsers.StripJSONComments(childData)
	var childFile parsers.TasksFile
	if err := json.Unmarshal(strippedChild, &childFile); err != nil {
		t.Fatalf("failed to unmarshal %s: %v", childName, err)
	}

	if childFile.Version != 2 {
		t.Errorf("%s: expected version 2, got %d", childName, childFile.Version)
	}

	if childFile.Parent == "" {
		t.Errorf("%s: expected parent field, got empty", childName)
	}
}

// Helper: count child files in a directory
func countChildFiles(t *testing.T, changeDir string) int {
	t.Helper()

	entries, err := os.ReadDir(changeDir)
	if err != nil {
		t.Fatalf("failed to read change dir: %v", err)
	}

	count := 0
	for _, entry := range entries {
		if !strings.HasPrefix(entry.Name(), "tasks-") ||
			!strings.HasSuffix(entry.Name(), ".jsonc") {
			continue
		}

		count++
		childPath := filepath.Join(changeDir, entry.Name())
		verifyChildFile(t, childPath, entry.Name())
	}

	return count
}

// TestAcceptIntegration150LineSplit tests that a 150-line tasks.md generates split files
func TestAcceptIntegration150LineSplit(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a proposal.md (required for validation)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(testProposalContent), 0o644); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	// Create a tasks.md with >100 lines to trigger splitting
	// Generate 150 lines with 30 tasks across 3 sections
	var tasksMd strings.Builder
	tasksMd.WriteString("## 1. Foundation\n\n")
	for i := 1; i <= 10; i++ {
		tasksMd.WriteString(fmt.Sprintf("- [ ] 1.%d Task %d in Foundation section\n", i, i))
	}
	tasksMd.WriteString("\n## 2. Implementation\n\n")
	for i := 1; i <= 10; i++ {
		tasksMd.WriteString(fmt.Sprintf("- [ ] 2.%d Task %d in Implementation section\n", i, i))
	}
	tasksMd.WriteString("\n## 3. Testing\n\n")
	for i := 1; i <= 10; i++ {
		tasksMd.WriteString(fmt.Sprintf("- [ ] 3.%d Task %d in Testing section\n", i, i))
	}

	// Add padding to reach 150 lines
	for range 120 {
		tasksMd.WriteString("\n")
	}

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMd.String()), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Verify file is >100 lines
	lineCount, err := countLines(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to count lines: %v", err)
	}
	if lineCount <= 100 {
		t.Fatalf("expected >100 lines, got %d", lineCount)
	}

	// Parse tasks and write files
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse tasks.md: %v", err)
	}

	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("writeAndCleanup() error = %v", err)
	}

	// Verify root tasks.jsonc was created and is version 2
	rootData, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read root tasks.jsonc: %v", err)
	}

	strippedRoot := parsers.StripJSONComments(rootData)
	var rootFile parsers.TasksFile
	if err := json.Unmarshal(strippedRoot, &rootFile); err != nil {
		t.Fatalf("failed to unmarshal root JSON: %v", err)
	}

	if rootFile.Version != 2 {
		t.Errorf("expected version 2, got %d", rootFile.Version)
	}

	if len(rootFile.Includes) == 0 {
		t.Error("expected includes field, got none")
	}

	// Verify child files were created
	childCount := countChildFiles(t, changeDir)
	if childCount == 0 {
		t.Error("expected child files to be created, got none")
	}

	t.Logf("Successfully created %d child files", childCount)
}

// TestAcceptIntegration80LineFlatFile tests that an 80-line tasks.md generates a flat file
func TestAcceptIntegration80LineFlatFile(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a proposal.md (required for validation)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(testProposalContent), 0o644); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	// Create a tasks.md with ≤100 lines (should NOT trigger splitting)
	var tasksMd strings.Builder
	tasksMd.WriteString("## 1. Foundation\n\n")
	for i := 1; i <= 5; i++ {
		tasksMd.WriteString(fmt.Sprintf("- [ ] 1.%d Task %d in Foundation section\n", i, i))
	}
	tasksMd.WriteString("\n## 2. Implementation\n\n")
	for i := 1; i <= 5; i++ {
		tasksMd.WriteString(fmt.Sprintf("- [ ] 2.%d Task %d in Implementation section\n", i, i))
	}

	// Add padding to reach 80 lines (below threshold)
	for range 68 {
		tasksMd.WriteString("\n")
	}

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMd.String()), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Verify file is ≤100 lines
	lineCount, err := countLines(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to count lines: %v", err)
	}
	if lineCount > 100 {
		t.Fatalf("expected ≤100 lines, got %d", lineCount)
	}

	// Parse tasks and write files
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse tasks.md: %v", err)
	}

	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("writeAndCleanup() error = %v", err)
	}

	// Verify root tasks.jsonc was created and is version 1 (flat)
	rootData, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	strippedRoot := parsers.StripJSONComments(rootData)
	var rootFile parsers.TasksFile
	if err := json.Unmarshal(strippedRoot, &rootFile); err != nil {
		t.Fatalf("failed to unmarshal JSON: %v", err)
	}

	if rootFile.Version != 1 {
		t.Errorf("expected version 1 (flat), got %d", rootFile.Version)
	}

	// Verify NO child files were created
	entries, err := os.ReadDir(changeDir)
	if err != nil {
		t.Fatalf("failed to read change dir: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "tasks-") && strings.HasSuffix(entry.Name(), ".jsonc") {
			t.Errorf("unexpected child file created: %s", entry.Name())
		}
	}

	// Verify all tasks are in the root file
	if len(rootFile.Tasks) != 10 {
		t.Errorf("expected 10 tasks in flat file, got %d", len(rootFile.Tasks))
	}

	t.Log("Successfully created flat file (no splitting)")
}

// Helper: find first child file in directory
func findFirstChildFile(t *testing.T, changeDir string) string {
	t.Helper()

	entries, err := os.ReadDir(changeDir)
	if err != nil {
		t.Fatalf("failed to read change dir: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "tasks-") && strings.HasSuffix(entry.Name(), ".jsonc") {
			return filepath.Join(changeDir, entry.Name())
		}
	}

	t.Fatal("no child file found")

	return ""
}

// Helper: modify child file statuses
func modifyChildFileStatuses(t *testing.T, childPath string) {
	t.Helper()

	childData, err := os.ReadFile(childPath)
	if err != nil {
		t.Fatalf("failed to read child file: %v", err)
	}

	strippedChild := parsers.StripJSONComments(childData)
	var childFile parsers.TasksFile
	if err := json.Unmarshal(strippedChild, &childFile); err != nil {
		t.Fatalf("failed to unmarshal child JSON: %v", err)
	}

	// Mark first two tasks as completed
	if len(childFile.Tasks) >= 2 {
		childFile.Tasks[0].Status = parsers.TaskStatusCompleted
		childFile.Tasks[1].Status = parsers.TaskStatusInProgress
	}

	// Write modified child file back
	modifiedJSON, err := json.MarshalIndent(childFile, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal modified child: %v", err)
	}

	header := childFileHeader("test-change", childFile.Parent)
	modifiedContent := header + string(modifiedJSON)
	if err := os.WriteFile(childPath, []byte(modifiedContent), 0o644); err != nil {
		t.Fatalf("failed to write modified child: %v", err)
	}
}

// Helper: verify child file has preserved statuses
func verifyChildFileStatuses(t *testing.T, childPath string) {
	t.Helper()

	childData, err := os.ReadFile(childPath)
	if err != nil {
		t.Fatalf("failed to read child file after regeneration: %v", err)
	}

	strippedChild := parsers.StripJSONComments(childData)
	var regeneratedChild parsers.TasksFile
	if err := json.Unmarshal(strippedChild, &regeneratedChild); err != nil {
		t.Fatalf("failed to unmarshal regenerated child: %v", err)
	}

	if len(regeneratedChild.Tasks) < 2 {
		t.Fatalf("expected at least 2 tasks, got %d", len(regeneratedChild.Tasks))
	}

	// Verify first task is still completed
	if regeneratedChild.Tasks[0].Status != parsers.TaskStatusCompleted {
		t.Errorf("task 0: expected status %s, got %s",
			parsers.TaskStatusCompleted,
			regeneratedChild.Tasks[0].Status,
		)
	}

	// Verify second task is still in_progress
	if regeneratedChild.Tasks[1].Status != parsers.TaskStatusInProgress {
		t.Errorf("task 1: expected status %s, got %s",
			parsers.TaskStatusInProgress,
			regeneratedChild.Tasks[1].Status,
		)
	}
}

// TestAcceptIntegrationRegenerationPreservesStatuses tests that re-running accept preserves task statuses
func TestAcceptIntegrationRegenerationPreservesStatuses(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "spectr", "changes", "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a proposal.md (required for validation)
	proposalPath := filepath.Join(changeDir, "proposal.md")
	if err := os.WriteFile(proposalPath, []byte(testProposalContent), 0o644); err != nil {
		t.Fatalf("failed to write proposal.md: %v", err)
	}

	// Create initial tasks.md with >100 lines
	var tasksMd strings.Builder
	tasksMd.WriteString("## 1. Foundation\n\n")
	for i := 1; i <= 10; i++ {
		tasksMd.WriteString(fmt.Sprintf("- [ ] 1.%d Task %d in Foundation\n", i, i))
	}
	tasksMd.WriteString("\n## 2. Implementation\n\n")
	for i := 1; i <= 10; i++ {
		tasksMd.WriteString(fmt.Sprintf("- [ ] 2.%d Task %d in Implementation\n", i, i))
	}

	// Add padding to reach 150 lines
	for range 128 {
		tasksMd.WriteString("\n")
	}

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMd.String()), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// First generation
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse tasks.md: %v", err)
	}

	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("first writeAndCleanup() error = %v", err)
	}

	// Modify statuses in child files
	firstChildPath := findFirstChildFile(t, changeDir)
	modifyChildFileStatuses(t, firstChildPath)

	// Regenerate (re-run accept)
	tasks, err = parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("failed to parse tasks.md (second time): %v", err)
	}

	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("second writeAndCleanup() error = %v", err)
	}

	// Verify statuses were preserved
	verifyChildFileStatuses(t, firstChildPath)

	t.Log("Successfully preserved statuses across regeneration")
}

// TestEdgeCaseTasksMdNoSections tests handling of tasks.md with no sections.
// According to the design, files with no sections should generate a flat file
// regardless of size, since there's nothing to split on.
func TestEdgeCaseTasksMdNoSections(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create tasks.md with tasks but no section headers
	tasksMd := `- [ ] Task one without section
- [ ] Task two without section
- [ ] Task three without section
- [ ] Task four without section
- [ ] Task five without section
`

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMd), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Parse tasks
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("parseTasksMd() error = %v", err)
	}

	// Verify we got tasks
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	// Write the file - should generate flat file (version 1) since no sections
	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("writeAndCleanup() error = %v", err)
	}

	// Verify flat file was generated
	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	stripped := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(stripped, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal tasks.jsonc: %v", err)
	}

	// Verify version 1 (flat file)
	if tasksFile.Version != 1 {
		t.Errorf("expected version 1 (flat file), got version %d", tasksFile.Version)
	}

	// Verify no child files were created
	entries, err := os.ReadDir(changeDir)
	if err != nil {
		t.Fatalf("failed to read change dir: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "tasks-") && strings.HasSuffix(entry.Name(), ".jsonc") {
			t.Errorf("unexpected child file created: %s", entry.Name())
		}
	}

	t.Log("Successfully handled tasks.md with no sections (flat file)")
}

// TestEdgeCaseTasksMdOnlyOneSection tests handling of tasks.md with only one section.
// According to the design, files with only one section don't need splitting since
// all tasks are already grouped together.
func TestEdgeCaseTasksMdOnlyOneSection(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create tasks.md with only one section
	tasksMd := `## 1. Single Section

- [ ] 1.1 Task one
- [ ] 1.2 Task two
- [ ] 1.3 Task three
- [ ] 1.4 Task four
- [ ] 1.5 Task five
`

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMd), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Parse tasks
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("parseTasksMd() error = %v", err)
	}

	// Verify we got tasks
	if len(tasks) != 5 {
		t.Fatalf("expected 5 tasks, got %d", len(tasks))
	}

	// Write the file - should generate flat file (version 1) since only one section
	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("writeAndCleanup() error = %v", err)
	}

	// Verify flat file was generated
	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	stripped := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(stripped, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal tasks.jsonc: %v", err)
	}

	// Verify version 1 (flat file) - single section doesn't need splitting
	if tasksFile.Version != 1 {
		t.Errorf("expected version 1 (flat file), got version %d", tasksFile.Version)
	}

	// Verify no child files were created
	entries, err := os.ReadDir(changeDir)
	if err != nil {
		t.Fatalf("failed to read change dir: %v", err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "tasks-") && strings.HasSuffix(entry.Name(), ".jsonc") {
			t.Errorf("unexpected child file created: %s", entry.Name())
		}
	}

	t.Log("Successfully handled tasks.md with only one section (flat file)")
}

// TestEdgeCaseTasksMdExactly100Lines tests handling of tasks.md with exactly 100 lines.
// According to the design, the threshold is ">" 100, so exactly 100 lines should use flat format.
func TestEdgeCaseTasksMdExactly100Lines(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a tasks.md file with exactly 100 lines
	lines := []string{
		"## 1. Section One",
		"",
	}
	for i := 1; i <= 48; i++ {
		lines = append(lines, fmt.Sprintf("- [ ] 1.%d Task %d", i, i))
	}
	// Total so far: 1 (header) + 1 (blank) + 48 (tasks) = 50 lines

	lines = append(lines,
		"",
		"## 2. Section Two",
		"",
	)
	for i := 1; i <= 47; i++ {
		lines = append(lines, fmt.Sprintf("- [ ] 2.%d Task %d", i, i))
	}
	// Total: 50 + 1 (blank) + 1 (header) + 1 (blank) + 47 (tasks) = 100 lines

	tasksMd := strings.Join(lines, "\n")

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMd), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Verify line count
	lineCount, err := countLines(tasksMdPath)
	if err != nil {
		t.Fatalf("countLines() error = %v", err)
	}
	if lineCount != 100 {
		t.Fatalf("expected exactly 100 lines, got %d", lineCount)
	}

	// Check if splitting is triggered
	split, err := shouldSplit(tasksMdPath)
	if err != nil {
		t.Fatalf("shouldSplit() error = %v", err)
	}

	// Exactly 100 lines should NOT trigger splitting (threshold is > 100)
	if split {
		t.Error("shouldSplit() = true for 100 lines, want false (threshold is > 100)")
	}

	// Parse tasks
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("parseTasksMd() error = %v", err)
	}

	// Write the file - should generate flat file (version 1)
	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("writeAndCleanup() error = %v", err)
	}

	// Verify flat file was generated
	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	stripped := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(stripped, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal tasks.jsonc: %v", err)
	}

	// Verify version 1 (flat file)
	if tasksFile.Version != 1 {
		t.Errorf(
			"expected version 1 (flat file) for exactly 100 lines, got version %d",
			tasksFile.Version,
		)
	}

	t.Log("Successfully handled tasks.md with exactly 100 lines (flat file)")
}

// TestEdgeCaseTasksMd101Lines tests handling of tasks.md with 101 lines.
// According to the design, 101 lines should trigger splitting.
func TestEdgeCaseTasksMd101Lines(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a tasks.md file with exactly 101 lines
	lines := []string{
		"## 1. Section One",
		"",
	}
	for i := 1; i <= 48; i++ {
		lines = append(lines, fmt.Sprintf("- [ ] 1.%d Task %d", i, i))
	}
	// Total so far: 1 (header) + 1 (blank) + 48 (tasks) = 50 lines

	lines = append(lines,
		"",
		"## 2. Section Two",
		"",
	)
	for i := 1; i <= 48; i++ {
		lines = append(lines, fmt.Sprintf("- [ ] 2.%d Task %d", i, i))
	}
	// Total: 50 + 1 (blank) + 1 (header) + 1 (blank) + 48 (tasks) = 101 lines

	tasksMd := strings.Join(lines, "\n")

	tasksMdPath := filepath.Join(changeDir, "tasks.md")
	if err := os.WriteFile(tasksMdPath, []byte(tasksMd), 0o644); err != nil {
		t.Fatalf("failed to write tasks.md: %v", err)
	}

	// Verify line count
	lineCount, err := countLines(tasksMdPath)
	if err != nil {
		t.Fatalf("countLines() error = %v", err)
	}
	if lineCount != 101 {
		t.Fatalf("expected exactly 101 lines, got %d", lineCount)
	}

	// Check if splitting is triggered
	split, err := shouldSplit(tasksMdPath)
	if err != nil {
		t.Fatalf("shouldSplit() error = %v", err)
	}

	// 101 lines should trigger splitting
	if !split {
		t.Error("shouldSplit() = false for 101 lines, want true (threshold is > 100)")
	}

	// Parse tasks
	tasks, err := parseTasksMd(tasksMdPath)
	if err != nil {
		t.Fatalf("parseTasksMd() error = %v", err)
	}

	// Write the file - should generate split files (version 2)
	tasksJSONPath := filepath.Join(changeDir, "tasks.jsonc")
	if err := writeAndCleanup(tasksMdPath, tasksJSONPath, tasks, nil); err != nil {
		t.Fatalf("writeAndCleanup() error = %v", err)
	}

	// Verify root file was generated
	data, err := os.ReadFile(tasksJSONPath)
	if err != nil {
		t.Fatalf("failed to read tasks.jsonc: %v", err)
	}

	stripped := parsers.StripJSONComments(data)
	var tasksFile parsers.TasksFile
	if err := json.Unmarshal(stripped, &tasksFile); err != nil {
		t.Fatalf("failed to unmarshal tasks.jsonc: %v", err)
	}

	// Verify version 2 (split files)
	if tasksFile.Version != 2 {
		t.Errorf(
			"expected version 2 (split files) for 101 lines, got version %d",
			tasksFile.Version,
		)
	}

	// Verify child files were created
	entries, err := os.ReadDir(changeDir)
	if err != nil {
		t.Fatalf("failed to read change dir: %v", err)
	}

	childFileCount := 0
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), "tasks-") && strings.HasSuffix(entry.Name(), ".jsonc") {
			childFileCount++
		}
	}

	// Should have at least one child file
	if childFileCount == 0 {
		t.Error("expected child files to be created for 101 lines, got none")
	}

	t.Logf(
		"Successfully handled tasks.md with 101 lines (split into %d child files)",
		childFileCount,
	)
}

// TestEdgeCaseMalformedTaskIDs tests graceful handling of malformed task IDs.
// The system should either auto-generate valid IDs or handle malformed IDs without crashing.
func TestEdgeCaseMalformedTaskIDs(t *testing.T) {
	tests := []struct {
		name        string
		markdown    string
		expectError bool
		description string
	}{
		{
			name: "task with invalid characters in ID",
			markdown: `## 1. Test Section

- [ ] 1.a Task with letter in ID
- [ ] 1.2 Normal task
`,
			expectError: false, // Should handle gracefully by auto-generating or preserving
			description: "Task IDs with invalid characters should be handled gracefully",
		},
		{
			name: "task with spaces in ID",
			markdown: `## 1. Test Section

- [ ] 1. 1 Task with space in ID
- [ ] 1.2 Normal task
`,
			expectError: false, // Parser should handle this
			description: "Task IDs with spaces should be handled",
		},
		{
			name: "task with empty ID",
			markdown: `## 1. Test Section

- [ ] Task with no ID
- [ ] 1.1 Normal task
`,
			expectError: false, // Should auto-generate ID
			description: "Tasks without IDs should get auto-generated IDs",
		},
		{
			name: "task with very long ID",
			markdown: `## 1. Test Section

- [ ] 1.1.1.1.1.1.1.1 Task with deeply nested ID
- [ ] 1.2 Normal task
`,
			expectError: false, // Should preserve or handle gracefully
			description: "Tasks with deeply nested IDs should be preserved",
		},
		{
			name: "duplicate task IDs",
			markdown: `## 1. Test Section

- [ ] 1.1 First task
- [ ] 1.1 Duplicate ID task
- [ ] 1.2 Normal task
`,
			expectError: false, // Parsing should work, validation catches duplicates
			description: "Duplicate task IDs should be parsed (validation catches them)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "tasks.md")

			if err := os.WriteFile(testFile, []byte(tt.markdown), 0o644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Parse the tasks
			tasks, err := parseTasksMd(testFile)
			if tt.expectError && err == nil {
				t.Errorf("expected error but got none for: %s", tt.description)
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v for: %s", err, tt.description)
			}

			// Verify we got some tasks (even if IDs are malformed)
			if !tt.expectError && len(tasks) == 0 {
				t.Error("expected tasks to be parsed despite malformed IDs, got 0 tasks")
			}

			// Try to write the tasks to verify no crashes
			if !tt.expectError && len(tasks) > 0 {
				tasksJSONPath := filepath.Join(tmpDir, "tasks.jsonc")
				err := writeTasksJSONC(tasksJSONPath, tasks, nil, nil)
				if err != nil {
					// Writing may fail for some edge cases, but shouldn't crash
					t.Logf("Writing failed (expected for some edge cases): %v", err)
				}
			}

			t.Logf("Successfully handled: %s", tt.description)
		})
	}
}

// TestEdgeCaseMissingParentTaskInChildFile tests validation of child files
// that reference non-existent parent tasks.
func TestEdgeCaseMissingParentTaskInChildFile(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	if err := os.MkdirAll(changeDir, 0o755); err != nil {
		t.Fatalf("failed to create change dir: %v", err)
	}

	// Create a root tasks.jsonc with version 2
	rootTasks := parsers.TasksFile{
		Version: 2,
		Tasks: []parsers.Task{
			{
				ID:          "1",
				Section:     "Section One",
				Description: "First section tasks",
				Status:      parsers.TaskStatusPending,
				Children:    "$ref:tasks-1.jsonc",
			},
			// Note: No task with ID "2" exists
		},
		Includes: []string{"tasks-*.jsonc"},
	}

	rootPath := filepath.Join(changeDir, "tasks.jsonc")
	rootData, err := json.MarshalIndent(rootTasks, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal root tasks: %v", err)
	}

	if err := os.WriteFile(rootPath, rootData, 0o644); err != nil {
		t.Fatalf("failed to write root tasks.jsonc: %v", err)
	}

	// Create a child file that references non-existent parent "2"
	childTasks := parsers.TasksFile{
		Version: 2,
		Parent:  "2", // This parent doesn't exist in root file
		Tasks: []parsers.Task{
			{
				ID:          "2.1",
				Section:     "Section Two",
				Description: "Task in non-existent parent",
				Status:      parsers.TaskStatusPending,
			},
		},
	}

	childPath := filepath.Join(changeDir, "tasks-2.jsonc")
	childData, err := json.MarshalIndent(childTasks, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal child tasks: %v", err)
	}

	if err := os.WriteFile(childPath, childData, 0o644); err != nil {
		t.Fatalf("failed to write child tasks file: %v", err)
	}

	// Try to load existing statuses - this should handle missing parent gracefully
	// The loadExistingStatuses function should not crash even if child references missing parent
	statusMap, err := loadExistingStatuses(changeDir)
	if err != nil {
		// Some implementations may return an error for orphaned child files
		t.Logf("loadExistingStatuses returned error (acceptable): %v", err)
	} else {
		// If no error, verify statuses were still loaded
		if len(statusMap) == 0 {
			t.Log("No statuses loaded (child file ignored due to missing parent)")
		} else {
			t.Logf("Statuses loaded despite missing parent: %d entries", len(statusMap))
		}
	}

	// The key is that the system should not crash - it should either:
	// 1. Skip the orphaned child file
	// 2. Return an error that can be handled
	// 3. Load the child tasks anyway (orphaned but not broken)
	t.Log("Successfully handled missing parent task in child file without crashing")
}

// TestEdgeCaseEmptyTasksFile tests handling of empty tasks.md file.
func TestEdgeCaseEmptyTasksFile(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "tasks.md")

	// Create empty file
	if err := os.WriteFile(testFile, []byte(""), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Parse the empty file
	tasks, err := parseTasksMd(testFile)
	if err != nil {
		t.Errorf("parseTasksMd() should handle empty file without error, got: %v", err)
	}

	// Should return empty task list
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks from empty file, got %d", len(tasks))
	}

	// Verify shouldSplit returns false for empty file
	split, err := shouldSplit(testFile)
	if err != nil {
		t.Fatalf("shouldSplit() error = %v", err)
	}

	if split {
		t.Error("shouldSplit() = true for empty file, want false")
	}

	t.Log("Successfully handled empty tasks.md file")
}

// TestEdgeCaseTasksFileWithOnlyComments tests handling of tasks.md with only comments.
func TestEdgeCaseTasksFileWithOnlyComments(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "tasks.md")

	markdown := `<!-- This is a comment -->
<!-- Another comment -->
<!-- TODO: Add tasks here -->
`

	if err := os.WriteFile(testFile, []byte(markdown), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Parse the file
	tasks, err := parseTasksMd(testFile)
	if err != nil {
		t.Errorf("parseTasksMd() should handle comments-only file, got error: %v", err)
	}

	// Should return empty task list
	if len(tasks) != 0 {
		t.Errorf("expected 0 tasks from comments-only file, got %d", len(tasks))
	}

	t.Log("Successfully handled tasks.md with only comments")
}

// TestJSONCValidation_JSONMetaCharacters tests that task descriptions containing
// JSON structural characters are properly escaped and survive round-trip validation.
func TestJSONCValidation_JSONMetaCharacters(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "opening brace",
			description: "Task with { opening brace",
		},
		{
			name:        "closing brace",
			description: "Task with } closing brace",
		},
		{
			name:        "opening bracket",
			description: "Task with [ opening bracket",
		},
		{
			name:        "closing bracket",
			description: "Task with ] closing bracket",
		},
		{
			name:        "colon",
			description: "Task with : colon separator",
		},
		{
			name:        "comma",
			description: "Task with , comma separator",
		},
		{
			name:        "json object",
			description: `Task with {"key": "value"} JSON object`,
		},
		{
			name:        "json array",
			description: "Task with [1, 2, 3] JSON array",
		},
		{
			name:        "nested json",
			description: `Task with {"items": [1, 2], "name": "test"} nested structure`,
		},
		{
			name:        "multiple braces",
			description: "Task with {{nested}} and {separate} braces",
		},
		{
			name:        "mixed json chars",
			description: `Handle: {"a": [1, 2], "b": {"c": "d"}}`,
		},
		{
			name:        "colon in middle",
			description: "Configure setting: value pairs properly",
		},
		{
			name:        "brackets and braces",
			description: "Array[{object}] syntax support",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a task with the test description
			task := parsers.Task{
				ID:          "1.1",
				Section:     "Test Section",
				Description: tt.description,
				Status:      parsers.TaskStatusPending,
			}

			// Step 1: Verify json.Marshal properly escapes the description
			marshalled, err := json.Marshal(task)
			if err != nil {
				t.Fatalf("json.Marshal failed: %v", err)
			}

			// Step 2: Verify the marshalled JSON is valid and can be unmarshalled
			var unmarshalled parsers.Task
			if err := json.Unmarshal(marshalled, &unmarshalled); err != nil {
				t.Fatalf("json.Unmarshal failed: %v\nMarshalled JSON: %s", err, string(marshalled))
			}

			// Step 3: Verify the description survived the round-trip unchanged
			if unmarshalled.Description != tt.description {
				t.Errorf(
					"Description changed during round-trip:\nOriginal: %q\nGot:      %q",
					tt.description,
					unmarshalled.Description,
				)
			}

			// Step 4: Test with TasksFile structure (as used in actual code)
			tasksFile := parsers.TasksFile{
				Version: 2,
				Tasks:   []parsers.Task{task},
			}

			// Marshal the full TasksFile structure with indentation
			fullMarshalled, err := json.MarshalIndent(tasksFile, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent failed: %v", err)
			}

			// Step 5: Verify the full structure can be parsed back
			var fullUnmarshalled parsers.TasksFile
			if err := json.Unmarshal(fullMarshalled, &fullUnmarshalled); err != nil {
				t.Fatalf(
					"Failed to unmarshal full TasksFile: %v\nJSON: %s",
					err,
					string(fullMarshalled),
				)
			}

			// Step 6: Verify task description is still intact
			if len(fullUnmarshalled.Tasks) != 1 {
				t.Fatalf("Expected 1 task, got %d", len(fullUnmarshalled.Tasks))
			}

			if fullUnmarshalled.Tasks[0].Description != tt.description {
				t.Errorf(
					"Description changed in full structure round-trip:\nOriginal: %q\nGot:      %q",
					tt.description,
					fullUnmarshalled.Tasks[0].Description,
				)
			}

			t.Logf("Successfully validated JSON meta characters in description: %q", tt.description)
		})
	}
}
