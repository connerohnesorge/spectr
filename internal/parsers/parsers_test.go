package parsers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "Change with prefix",
			content:  "# Change: Add Feature\n\nMore content",
			expected: "Add Feature",
		},
		{
			name:     "Spec with prefix",
			content:  "# Spec: Authentication\n\nMore content",
			expected: "Authentication",
		},
		{
			name:     "No prefix",
			content:  "# CLI Framework\n\nMore content",
			expected: "CLI Framework",
		},
		{
			name:     "Multiple headings",
			content:  "# First Heading\n## Second Heading\n# Third Heading",
			expected: "First Heading",
		},
		{
			name:     "Extra whitespace",
			content:  "#   Change:   Trim Whitespace   \n\nMore content",
			expected: "Trim Whitespace",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(
				tmpDir,
				"test.md",
			)
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			title, err := ExtractTitle(filePath)
			if err != nil {
				t.Fatalf(
					"ExtractTitle failed: %v",
					err,
				)
			}
			if title != tt.expected {
				t.Errorf(
					"Expected title %q, got %q",
					tt.expected,
					title,
				)
			}
		})
	}
}

func TestExtractTitle_NoHeading(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.md")
	content := "Some content without heading\n\nMore content"
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	title, err := ExtractTitle(filePath)
	if err != nil {
		t.Fatalf("ExtractTitle failed: %v", err)
	}
	if title != "" {
		t.Errorf(
			"Expected empty title, got %q",
			title,
		)
	}
}

func TestCountTasks(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedTotal     int
		expectedCompleted int
	}{
		{
			name: "Mixed tasks",
			content: `## Tasks
- [ ] Task 1
- [x] Task 2
- [ ] Task 3
- [X] Task 4`,
			expectedTotal:     4,
			expectedCompleted: 2,
		},
		{
			name: "All completed",
			content: `## Tasks
- [x] Task 1
- [X] Task 2`,
			expectedTotal:     2,
			expectedCompleted: 2,
		},
		{
			name: "All incomplete",
			content: `## Tasks
- [ ] Task 1
- [ ] Task 2`,
			expectedTotal:     2,
			expectedCompleted: 0,
		},
		{
			name: "With indentation",
			content: `## Tasks
  - [ ] Indented task 1
    - [x] Nested task 2`,
			expectedTotal:     2,
			expectedCompleted: 1,
		},
		{
			name: "Mixed content",
			content: `## Tasks
Some text
- [ ] Task 1
More text
- [x] Task 2
Not a task line`,
			expectedTotal:     2,
			expectedCompleted: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(
				tmpDir,
				"tasks.md",
			)
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			// CountTasks now takes a directory path, not a file path
			status, err := CountTasks(tmpDir)
			if err != nil {
				t.Fatalf(
					"CountTasks failed: %v",
					err,
				)
			}
			if status.Total != tt.expectedTotal {
				t.Errorf(
					"Expected total %d, got %d",
					tt.expectedTotal,
					status.Total,
				)
			}
			if status.Completed != tt.expectedCompleted {
				t.Errorf(
					"Expected completed %d, got %d",
					tt.expectedCompleted,
					status.Completed,
				)
			}
		})
	}
}

func TestCountTasks_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	// Empty directory with no tasks.json or tasks.md

	status, err := CountTasks(tmpDir)
	if err != nil {
		t.Fatalf(
			"CountTasks should not error on missing file: %v",
			err,
		)
	}
	if status.Total != 0 ||
		status.Completed != 0 ||
		status.InProgress != 0 {
		t.Errorf(
			"Expected zero status, got total=%d, completed=%d, inProgress=%d",
			status.Total,
			status.Completed,
			status.InProgress,
		)
	}
}

func TestCountDeltas(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(
		tmpDir,
		"test-change",
	)
	specsDir := filepath.Join(
		changeDir,
		"specs",
		"test-spec",
	)
	if err := os.MkdirAll(specsDir, 0o755); err != nil {
		t.Fatal(err)
	}

	specContent := `# Test Spec

## ADDED Requirements
### Requirement: New Feature

## MODIFIED Requirements
### Requirement: Updated Feature

## REMOVED Requirements
### Requirement: Old Feature
`
	specPath := filepath.Join(specsDir, "spec.md")
	if err := os.WriteFile(specPath, []byte(specContent), 0o644); err != nil {
		t.Fatal(err)
	}

	count, err := CountDeltas(changeDir)
	if err != nil {
		t.Fatalf("CountDeltas failed: %v", err)
	}
	if count != 3 {
		t.Errorf(
			"Expected 3 deltas, got %d",
			count,
		)
	}
}

func TestCountDeltas_NoSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	count, err := CountDeltas(tmpDir)
	if err != nil {
		t.Fatalf(
			"CountDeltas should not error on missing specs: %v",
			err,
		)
	}
	if count != 0 {
		t.Errorf(
			"Expected 0 deltas, got %d",
			count,
		)
	}
}

func TestCountRequirements(t *testing.T) {
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")

	content := `# Test Spec

### Requirement: Feature 1
Description

### Requirement: Feature 2
Description

## Another Section

### Requirement: Feature 3
Description
`
	if err := os.WriteFile(specPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	count, err := CountRequirements(specPath)
	if err != nil {
		t.Fatalf(
			"CountRequirements failed: %v",
			err,
		)
	}
	if count != 3 {
		t.Errorf(
			"Expected 3 requirements, got %d",
			count,
		)
	}
}

func TestCountRequirements_NoRequirements(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")

	content := `# Test Spec

Some content without requirements
`
	if err := os.WriteFile(specPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	count, err := CountRequirements(specPath)
	if err != nil {
		t.Fatalf(
			"CountRequirements failed: %v",
			err,
		)
	}
	if count != 0 {
		t.Errorf(
			"Expected 0 requirements, got %d",
			count,
		)
	}
}

func TestReadTasksJson(t *testing.T) {
	tmpDir := t.TempDir()
	tasksJsonPath := filepath.Join(
		tmpDir,
		"tasks.json",
	)

	content := `{
	"version": 1,
	"tasks": [
		{"id": "1.1", "section": "Implementation", "description": "Task 1", "status": "completed"},
		{"id": "1.2", "section": "Implementation", "description": "Task 2", "status": "in_progress"},
		{"id": "1.3", "section": "Implementation", "description": "Task 3", "status": "pending"}
	]
}`
	if err := os.WriteFile(tasksJsonPath, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	tasksFile, err := ReadTasksJson(tasksJsonPath)
	if err != nil {
		t.Fatalf("ReadTasksJson failed: %v", err)
	}

	if tasksFile.Version != 1 {
		t.Errorf(
			"Expected version 1, got %d",
			tasksFile.Version,
		)
	}
	if len(tasksFile.Tasks) != 3 {
		t.Errorf(
			"Expected 3 tasks, got %d",
			len(tasksFile.Tasks),
		)
	}
	if tasksFile.Tasks[0].Status != TaskStatusCompleted {
		t.Errorf(
			"Expected first task status to be completed, got %s",
			tasksFile.Tasks[0].Status,
		)
	}
	if tasksFile.Tasks[1].Status != TaskStatusInProgress {
		t.Errorf(
			"Expected second task status to be in_progress, got %s",
			tasksFile.Tasks[1].Status,
		)
	}
	if tasksFile.Tasks[2].Status != TaskStatusPending {
		t.Errorf(
			"Expected third task status to be pending, got %s",
			tasksFile.Tasks[2].Status,
		)
	}
}

func TestCountTasks_FromJsonc(t *testing.T) {
	tests := []struct {
		name               string
		content            string
		expectedTotal      int
		expectedCompleted  int
		expectedInProgress int
	}{
		{
			name: "Mixed status tasks",
			content: `{
				"version": 1,
				"tasks": [
					{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "completed"},
					{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "in_progress"},
					{"id": "1.3", "section": "Impl", "description": "Task 3", "status": "pending"},
					{"id": "1.4", "section": "Impl", "description": "Task 4", "status": "completed"}
				]
			}`,
			expectedTotal:      4,
			expectedCompleted:  2,
			expectedInProgress: 1,
		},
		{
			name: "All completed",
			content: `{
				"version": 1,
				"tasks": [
					{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "completed"},
					{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "completed"}
				]
			}`,
			expectedTotal:      2,
			expectedCompleted:  2,
			expectedInProgress: 0,
		},
		{
			name: "All pending",
			content: `{
				"version": 1,
				"tasks": [
					{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "pending"},
					{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "pending"}
				]
			}`,
			expectedTotal:      2,
			expectedCompleted:  0,
			expectedInProgress: 0,
		},
		{
			name: "Empty tasks array",
			content: `{
				"version": 1,
				"tasks": []
			}`,
			expectedTotal:      0,
			expectedCompleted:  0,
			expectedInProgress: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(
				tmpDir,
				"tasks.jsonc",
			)
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			status, err := CountTasks(tmpDir)
			if err != nil {
				t.Fatalf(
					"CountTasks failed: %v",
					err,
				)
			}
			if status.Total != tt.expectedTotal {
				t.Errorf(
					"Expected total %d, got %d",
					tt.expectedTotal,
					status.Total,
				)
			}
			if status.Completed != tt.expectedCompleted {
				t.Errorf(
					"Expected completed %d, got %d",
					tt.expectedCompleted,
					status.Completed,
				)
			}
			if status.InProgress != tt.expectedInProgress {
				t.Errorf(
					"Expected in_progress %d, got %d",
					tt.expectedInProgress,
					status.InProgress,
				)
			}
		})
	}
}

func TestCountTasks_JsoncPreferredOverMarkdown(
	t *testing.T,
) {
	tmpDir := t.TempDir()

	// Create both tasks.jsonc and tasks.md
	// tasks.jsonc has 2 tasks
	jsoncContent := `// This is a JSONC file with comments
{
		"version": 1,
		"tasks": [
			{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "completed"},
			{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "pending"}
		]
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "tasks.jsonc"), []byte(jsoncContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// tasks.md has 5 tasks (different count to verify JSONC is used)
	mdContent := `## Tasks
- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
- [x] Task 4
- [x] Task 5`
	if err := os.WriteFile(filepath.Join(tmpDir, "tasks.md"), []byte(mdContent), 0o644); err != nil {
		t.Fatal(err)
	}

	// Should use tasks.jsonc, not tasks.md
	status, err := CountTasks(tmpDir)
	if err != nil {
		t.Fatalf("CountTasks failed: %v", err)
	}

	// Expect counts from JSONC (2 total, 1 completed) not MD (5 total, 2 completed)
	if status.Total != 2 {
		t.Errorf(
			"Expected total 2 (from JSONC), got %d",
			status.Total,
		)
	}
	if status.Completed != 1 {
		t.Errorf(
			"Expected completed 1 (from JSONC), got %d",
			status.Completed,
		)
	}
}

func TestCountTasks_IgnoresLegacyJson(
	t *testing.T,
) {
	// Test that tasks.json (legacy) is ignored in favor of tasks.md
	t.Run(
		"tasks.json ignored, falls back to tasks.md",
		func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create legacy tasks.json with 2 tasks
			jsonContent := `{
			"version": 1,
			"tasks": [
				{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "completed"},
				{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "pending"}
			]
		}`
			if err := os.WriteFile(filepath.Join(tmpDir, "tasks.json"), []byte(jsonContent), 0o644); err != nil {
				t.Fatal(err)
			}

			// Create tasks.md with 4 tasks
			mdContent := `## Tasks
- [ ] Task 1
- [ ] Task 2
- [x] Task 3
- [x] Task 4`
			if err := os.WriteFile(filepath.Join(tmpDir, "tasks.md"), []byte(mdContent), 0o644); err != nil {
				t.Fatal(err)
			}

			// Should use tasks.md, NOT tasks.json (legacy is ignored)
			status, err := CountTasks(tmpDir)
			if err != nil {
				t.Fatalf(
					"CountTasks failed: %v",
					err,
				)
			}

			// Expect counts from MD (4 total, 2 completed) not legacy JSON (2 total, 1 completed)
			if status.Total != 4 {
				t.Errorf(
					"Expected total 4 (from tasks.md), got %d - tasks.json should be ignored",
					status.Total,
				)
			}
			if status.Completed != 2 {
				t.Errorf(
					"Expected completed 2 (from tasks.md), got %d - tasks.json should be ignored",
					status.Completed,
				)
			}
		},
	)

	// Test that tasks.jsonc takes priority over both tasks.json and tasks.md
	t.Run(
		"tasks.jsonc wins over legacy tasks.json and tasks.md",
		func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create tasks.jsonc with 3 tasks (all completed)
			jsoncContent := `// JSONC format with comments
{
			"version": 1,
			"tasks": [
				{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "completed"},
				{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "completed"},
				{"id": "1.3", "section": "Impl", "description": "Task 3", "status": "completed"}
			]
		}`
			if err := os.WriteFile(filepath.Join(tmpDir, "tasks.jsonc"), []byte(jsoncContent), 0o644); err != nil {
				t.Fatal(err)
			}

			// Create legacy tasks.json with 2 tasks
			jsonContent := `{
			"version": 1,
			"tasks": [
				{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "pending"},
				{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "pending"}
			]
		}`
			if err := os.WriteFile(filepath.Join(tmpDir, "tasks.json"), []byte(jsonContent), 0o644); err != nil {
				t.Fatal(err)
			}

			// Create tasks.md with 5 tasks
			mdContent := `## Tasks
- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
- [ ] Task 4
- [ ] Task 5`
			if err := os.WriteFile(filepath.Join(tmpDir, "tasks.md"), []byte(mdContent), 0o644); err != nil {
				t.Fatal(err)
			}

			// Should use tasks.jsonc (3 total, 3 completed)
			status, err := CountTasks(tmpDir)
			if err != nil {
				t.Fatalf(
					"CountTasks failed: %v",
					err,
				)
			}

			if status.Total != 3 {
				t.Errorf(
					"Expected total 3 (from tasks.jsonc), got %d",
					status.Total,
				)
			}
			if status.Completed != 3 {
				t.Errorf(
					"Expected completed 3 (from tasks.jsonc), got %d",
					status.Completed,
				)
			}
		},
	)
}

func TestStripJSONComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "No comments",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "Single-line comment at start",
			input:    "// This is a comment\n{\"key\": \"value\"}",
			expected: "\n{\"key\": \"value\"}",
		},
		{
			name:     "Single-line comment inline",
			input:    "{\"key\": \"value\"} // comment",
			expected: "{\"key\": \"value\"} ",
		},
		{
			name:  "Multiple single-line comments",
			input: "// Header comment\n{\n\t// Comment before key\n\t\"key\": \"value\" // inline comment\n}",
			// Note: whitespace before comments is preserved, only the comment text is stripped
			expected: "\n{\n\t\n\t\"key\": \"value\" \n}",
		},
		{
			name:     "Multi-line comment",
			input:    "/* This is a\nmulti-line comment */\n{\"key\": \"value\"}",
			expected: "\n{\"key\": \"value\"}",
		},
		{
			name:     "Multi-line comment inline",
			input:    "{\"key\": /* comment */ \"value\"}",
			expected: "{\"key\":  \"value\"}",
		},
		{
			name:     "Comment-like content inside string preserved",
			input:    `{"url": "https://example.com", "comment": "This has // slashes"}`,
			expected: `{"url": "https://example.com", "comment": "This has // slashes"}`,
		},
		{
			name:     "Multi-line comment syntax inside string preserved",
			input:    `{"code": "/* not a comment */"}`,
			expected: `{"code": "/* not a comment */"}`,
		},
		{
			name:     "Escaped quotes in strings",
			input:    `{"message": "He said \"hello\"", "test": "value"} // comment`,
			expected: `{"message": "He said \"hello\"", "test": "value"} `,
		},
		{
			name:     "Complex escape sequences in strings",
			input:    `{"path": "C:\\path\\to\\file", "quote": "\""} // end`,
			expected: `{"path": "C:\\path\\to\\file", "quote": "\""} `,
		},
		{
			name:  "Mixed comments and data",
			input: "// File header\n/* Block comment explaining format */\n{\n\t\"version\": 1, // version number\n\t\"tasks\": [\n\t\t/* First task */\n\t\t{\"id\": \"1.1\", \"description\": \"Task with // in name\"}\n\t]\n}",
			// Note: whitespace before comments is preserved
			expected: "\n\n{\n\t\"version\": 1, \n\t\"tasks\": [\n\t\t\n\t\t{\"id\": \"1.1\", \"description\": \"Task with // in name\"}\n\t]\n}",
		},
		{
			name:     "Empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "Only comments",
			input:    "// Just a comment\n/* Another one */",
			expected: "\n",
		},
		{
			name:     "Nested-looking but not nested multi-line comments",
			input:    "/* outer /* inner */ after */ {\"key\": \"value\"}",
			expected: " after */ {\"key\": \"value\"}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripJSONComments(
				[]byte(tt.input),
			)
			if string(result) != tt.expected {
				t.Errorf(
					"StripJSONComments(%q)\n  got:      %q\n  expected: %q",
					tt.input,
					string(result),
					tt.expected,
				)
			}
		})
	}
}

func TestReadTasksJsonWithComments(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedVersion   int
		expectedTaskCount int
		expectedFirstID   string
		expectedFirstDesc string
	}{
		{
			name: "JSONC with header comment",
			content: `// tasks.jsonc - Spectr task tracking file
// Status values: pending, in_progress, completed
{
	"version": 1,
	"tasks": [
		{"id": "1.1", "section": "Setup", "description": "Initialize project", "status": "completed"}
	]
}`,
			expectedVersion:   1,
			expectedTaskCount: 1,
			expectedFirstID:   "1.1",
			expectedFirstDesc: "Initialize project",
		},
		{
			name: "JSONC with inline comments",
			content: `{
	"version": 1, // Schema version
	"tasks": [
		{
			"id": "2.1",       // Task identifier
			"section": "Impl", // Category
			"description": "Build feature", // What to do
			"status": "pending" // Current state
		}
	]
}`,
			expectedVersion:   1,
			expectedTaskCount: 1,
			expectedFirstID:   "2.1",
			expectedFirstDesc: "Build feature",
		},
		{
			name: "JSONC with block comments",
			content: `/*
 * This is the tasks file for the project.
 * It tracks implementation progress.
 */
{
	"version": 1,
	"tasks": [
		/* First task in the list */
		{"id": "3.1", "section": "Testing", "description": "Write tests", "status": "in_progress"}
	]
}`,
			expectedVersion:   1,
			expectedTaskCount: 1,
			expectedFirstID:   "3.1",
			expectedFirstDesc: "Write tests",
		},
		{
			name: "Multiple tasks with comments",
			content: `// Task file
{
	"version": 1,
	"tasks": [
		// Setup tasks
		{"id": "1.1", "section": "Setup", "description": "First task", "status": "completed"},
		{"id": "1.2", "section": "Setup", "description": "Second task", "status": "pending"},
		// Implementation tasks
		{"id": "2.1", "section": "Impl", "description": "Third task", "status": "in_progress"}
	]
}`,
			expectedVersion:   1,
			expectedTaskCount: 3,
			expectedFirstID:   "1.1",
			expectedFirstDesc: "First task",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(
				tmpDir,
				"tasks.jsonc",
			)
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			tasksFile, err := ReadTasksJson(
				filePath,
			)
			if err != nil {
				t.Fatalf(
					"ReadTasksJson failed: %v",
					err,
				)
			}

			if tasksFile.Version != tt.expectedVersion {
				t.Errorf(
					"Expected version %d, got %d",
					tt.expectedVersion,
					tasksFile.Version,
				)
			}
			if len(
				tasksFile.Tasks,
			) != tt.expectedTaskCount {
				t.Errorf(
					"Expected %d tasks, got %d",
					tt.expectedTaskCount,
					len(tasksFile.Tasks),
				)
			}
			if len(tasksFile.Tasks) == 0 {
				return
			}
			if tasksFile.Tasks[0].ID != tt.expectedFirstID {
				t.Errorf(
					"Expected first task ID %q, got %q",
					tt.expectedFirstID,
					tasksFile.Tasks[0].ID,
				)
			}
			if tasksFile.Tasks[0].Description != tt.expectedFirstDesc {
				t.Errorf(
					"Expected first task description %q, got %q",
					tt.expectedFirstDesc,
					tasksFile.Tasks[0].Description,
				)
			}
		})
	}
}

func TestCountTasksFromJson_Fixture(
	t *testing.T,
) {
	// Test with real fixture - 16 pending tasks
	status, err := countTasksFromJson(
		"testdata/tasks_fixture.jsonc",
	)
	if err != nil {
		t.Fatalf(
			"countTasksFromJson failed: %v",
			err,
		)
	}

	// Fixture has 16 tasks, all pending
	if status.Total != 16 {
		t.Errorf(
			"Expected Total=16, got %d",
			status.Total,
		)
	}
	if status.Completed != 0 {
		t.Errorf(
			"Expected Completed=0, got %d",
			status.Completed,
		)
	}
	if status.InProgress != 0 {
		t.Errorf(
			"Expected InProgress=0, got %d",
			status.InProgress,
		)
	}
}

func TestCountTasksFromJson_AllScenarios(
	t *testing.T,
) {
	tests := []struct {
		name           string
		content        string
		wantTotal      int
		wantCompleted  int
		wantInProgress int
	}{
		{
			name: "all pending",
			content: `{"version":1,"tasks":[
				{"id":"1","section":"A","description":"T1","status":"pending"},
				{"id":"2","section":"A","description":"T2","status":"pending"},
				{"id":"3","section":"A","description":"T3","status":"pending"}
			]}`,
			wantTotal:      3,
			wantCompleted:  0,
			wantInProgress: 0,
		},
		{
			name: "all completed",
			content: `{"version":1,"tasks":[
				{"id":"1","section":"A","description":"T1","status":"completed"},
				{"id":"2","section":"A","description":"T2","status":"completed"}
			]}`,
			wantTotal:      2,
			wantCompleted:  2,
			wantInProgress: 0,
		},
		{
			name: "mixed states",
			content: `{"version":1,"tasks":[
				{"id":"1","section":"A","description":"T1","status":"completed"},
				{"id":"2","section":"A","description":"T2","status":"in_progress"},
				{"id":"3","section":"A","description":"T3","status":"pending"},
				{"id":"4","section":"A","description":"T4","status":"completed"}
			]}`,
			wantTotal:      4,
			wantCompleted:  2,
			wantInProgress: 1,
		},
		{
			name:           "empty tasks",
			content:        `{"version":1,"tasks":[]}`,
			wantTotal:      0,
			wantCompleted:  0,
			wantInProgress: 0,
		},
		{
			name: "only in_progress",
			content: `{"version":1,"tasks":[
				{"id":"1","section":"A","description":"T1","status":"in_progress"},
				{"id":"2","section":"A","description":"T2","status":"in_progress"}
			]}`,
			wantTotal:      2,
			wantCompleted:  0,
			wantInProgress: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(
				tmpDir,
				"tasks.jsonc",
			)
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			status, err := countTasksFromJson(
				filePath,
			)
			if err != nil {
				t.Fatalf(
					"countTasksFromJson failed: %v",
					err,
				)
			}

			if status.Total != tt.wantTotal {
				t.Errorf(
					"Total = %d, want %d",
					status.Total,
					tt.wantTotal,
				)
			}
			if status.Completed != tt.wantCompleted {
				t.Errorf(
					"Completed = %d, want %d",
					status.Completed,
					tt.wantCompleted,
				)
			}
			if status.InProgress != tt.wantInProgress {
				t.Errorf(
					"InProgress = %d, want %d",
					status.InProgress,
					tt.wantInProgress,
				)
			}
		})
	}
}

func TestReadTasksJsonWithTrailingCommas(t *testing.T) {
	tests := []struct {
		name              string
		content           string
		expectedVersion   int
		expectedTaskCount int
		wantError         bool
	}{
		{
			name: "Trailing comma in array",
			content: `{
	"version": 1,
	"tasks": [
		{"id": "1.1", "section": "Setup", "description": "Task 1", "status": "pending"},
		{"id": "1.2", "section": "Setup", "description": "Task 2", "status": "completed"},
	]
}`,
			expectedVersion:   1,
			expectedTaskCount: 2,
			wantError:         false,
		},
		{
			name: "Trailing comma in object",
			content: `{
	"version": 1,
	"tasks": [
		{"id": "1.1", "section": "Setup", "description": "Task 1", "status": "pending",}
	],
}`,
			expectedVersion:   1,
			expectedTaskCount: 1,
			wantError:         false,
		},
		{
			name: "Multiple trailing commas with comments",
			content: `// Task file with trailing commas
{
	"version": 1, // version number
	"tasks": [
		{"id": "1.1", "section": "Setup", "description": "Task 1", "status": "pending",}, // first task
		{"id": "1.2", "section": "Setup", "description": "Task 2", "status": "completed",}, // second task
	], // end tasks
}`,
			expectedVersion:   1,
			expectedTaskCount: 2,
			wantError:         false,
		},
		{
			name: "Nested trailing commas",
			content: `{
	"version": 1,
	"tasks": [
		{
			"id": "1.1",
			"section": "Setup",
			"description": "Task with trailing comma",
			"status": "pending",
		},
	],
}`,
			expectedVersion:   1,
			expectedTaskCount: 1,
			wantError:         false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "tasks.jsonc")
			if err := os.WriteFile(filePath, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			tasksFile, err := ReadTasksJson(filePath)
			if tt.wantError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tt.wantError && err != nil {
				t.Fatalf("ReadTasksJson failed: %v", err)
			}
			if tt.wantError {
				return
			}

			if tasksFile.Version != tt.expectedVersion {
				t.Errorf("Expected version %d, got %d", tt.expectedVersion, tasksFile.Version)
			}
			if len(tasksFile.Tasks) != tt.expectedTaskCount {
				t.Errorf("Expected %d tasks, got %d", tt.expectedTaskCount, len(tasksFile.Tasks))
			}
		})
	}
}
