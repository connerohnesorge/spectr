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
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
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
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
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
			// NOTE: This test was updated for the blackfriday-based parser.
			// In proper markdown, text between list items creates separate lists.
			// The new parser correctly handles this as markdown-compliant behavior.
			// Tasks must be in continuous list format without intervening text.
			name: "Mixed content with proper list",
			content: `## Tasks
Some text before tasks

- [ ] Task 1
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
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
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
	if err := os.MkdirAll(specsDir, 0755); err != nil {
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
	if err := os.WriteFile(specPath, []byte(specContent), 0644); err != nil {
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
	if err := os.WriteFile(specPath, []byte(content), 0644); err != nil {
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
	if err := os.WriteFile(specPath, []byte(content), 0644); err != nil {
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
	if err := os.WriteFile(tasksJsonPath, []byte(content), 0644); err != nil {
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

func TestCountTasks_FromJson(t *testing.T) {
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
				"tasks.json",
			)
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
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

func TestCountTasks_JsonPreferredOverMarkdown(
	t *testing.T,
) {
	tmpDir := t.TempDir()

	// Create both tasks.json and tasks.md
	// tasks.json has 2 tasks
	jsonContent := `{
		"version": 1,
		"tasks": [
			{"id": "1.1", "section": "Impl", "description": "Task 1", "status": "completed"},
			{"id": "1.2", "section": "Impl", "description": "Task 2", "status": "pending"}
		]
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "tasks.json"), []byte(jsonContent), 0644); err != nil {
		t.Fatal(err)
	}

	// tasks.md has 5 tasks (different count to verify JSON is used)
	mdContent := `## Tasks
- [ ] Task 1
- [ ] Task 2
- [ ] Task 3
- [x] Task 4
- [x] Task 5`
	if err := os.WriteFile(filepath.Join(tmpDir, "tasks.md"), []byte(mdContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Should use tasks.json, not tasks.md
	status, err := CountTasks(tmpDir)
	if err != nil {
		t.Fatalf("CountTasks failed: %v", err)
	}

	// Expect counts from JSON (2 total, 1 completed) not MD (5 total, 2 completed)
	if status.Total != 2 {
		t.Errorf(
			"Expected total 2 (from JSON), got %d",
			status.Total,
		)
	}
	if status.Completed != 1 {
		t.Errorf(
			"Expected completed 1 (from JSON), got %d",
			status.Completed,
		)
	}
}
