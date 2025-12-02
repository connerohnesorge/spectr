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
			filePath := filepath.Join(tmpDir, "test.md")
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			title, err := ExtractTitle(filePath)
			if err != nil {
				t.Fatalf("ExtractTitle failed: %v", err)
			}
			if title != tt.expected {
				t.Errorf("Expected title %q, got %q", tt.expected, title)
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
		t.Errorf("Expected empty title, got %q", title)
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
			filePath := filepath.Join(tmpDir, "tasks.md")
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}

			status, err := CountTasks(filePath)
			if err != nil {
				t.Fatalf("CountTasks failed: %v", err)
			}
			if status.Total != tt.expectedTotal {
				t.Errorf("Expected total %d, got %d", tt.expectedTotal, status.Total)
			}
			if status.Completed != tt.expectedCompleted {
				t.Errorf("Expected completed %d, got %d", tt.expectedCompleted, status.Completed)
			}
		})
	}
}

func TestCountTasks_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.md")

	status, err := CountTasks(filePath)
	if err != nil {
		t.Fatalf("CountTasks should not error on missing file: %v", err)
	}
	if status.Total != 0 || status.Completed != 0 {
		t.Errorf("Expected zero status, got total=%d, completed=%d", status.Total, status.Completed)
	}
}

func TestCountDeltas(t *testing.T) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(tmpDir, "test-change")
	specsDir := filepath.Join(changeDir, "specs", "test-spec")
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
		t.Errorf("Expected 3 deltas, got %d", count)
	}
}

func TestCountDeltas_NoSpecs(t *testing.T) {
	tmpDir := t.TempDir()
	count, err := CountDeltas(tmpDir)
	if err != nil {
		t.Fatalf("CountDeltas should not error on missing specs: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 deltas, got %d", count)
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
		t.Fatalf("CountRequirements failed: %v", err)
	}
	if count != 3 {
		t.Errorf("Expected 3 requirements, got %d", count)
	}
}

func TestCountRequirements_NoRequirements(t *testing.T) {
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
		t.Fatalf("CountRequirements failed: %v", err)
	}
	if count != 0 {
		t.Errorf("Expected 0 requirements, got %d", count)
	}
}

func TestValidateTasksStructure_ValidFile(t *testing.T) {
	content := `## 1. Implementation
- [ ] 1.1 Create database schema
- [ ] 1.2 Implement API endpoint

## 2. Testing
- [ ] 2.1 Write unit tests
- [x] 2.2 Write integration tests
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure failed: %v", err)
	}

	// Should have 2 sections
	if len(result.Sections) != 2 {
		t.Errorf("Expected 2 sections, got %d", len(result.Sections))
	}

	// First section should be "Implementation" with 2 tasks
	if result.Sections[0].Number != 1 {
		t.Errorf("Expected section 1 number to be 1, got %d", result.Sections[0].Number)
	}
	if result.Sections[0].Name != "Implementation" {
		t.Errorf("Expected section 1 name to be 'Implementation', got %q", result.Sections[0].Name)
	}
	if result.Sections[0].TaskCount != 2 {
		t.Errorf("Expected section 1 to have 2 tasks, got %d", result.Sections[0].TaskCount)
	}

	// Second section should be "Testing" with 2 tasks
	if result.Sections[1].Number != 2 {
		t.Errorf("Expected section 2 number to be 2, got %d", result.Sections[1].Number)
	}
	if result.Sections[1].Name != "Testing" {
		t.Errorf("Expected section 2 name to be 'Testing', got %q", result.Sections[1].Name)
	}
	if result.Sections[1].TaskCount != 2 {
		t.Errorf("Expected section 2 to have 2 tasks, got %d", result.Sections[1].TaskCount)
	}

	// No orphaned tasks
	if result.OrphanedTasks != 0 {
		t.Errorf("Expected 0 orphaned tasks, got %d", result.OrphanedTasks)
	}

	// No empty sections
	if len(result.EmptySections) != 0 {
		t.Errorf(
			"Expected 0 empty sections, got %d: %v",
			len(result.EmptySections),
			result.EmptySections,
		)
	}

	// Sequential numbers
	if !result.SequentialNumbers {
		t.Error("Expected sequential numbers to be true")
	}

	// No gaps
	if len(result.NonSequentialGaps) != 0 {
		t.Errorf(
			"Expected 0 non-sequential gaps, got %d: %v",
			len(result.NonSequentialGaps),
			result.NonSequentialGaps,
		)
	}
}

func TestValidateTasksStructure_NoSections(t *testing.T) {
	content := `# Tasks

- [ ] Task 1
- [x] Task 2
- [ ] Task 3
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure failed: %v", err)
	}

	// No sections
	if len(result.Sections) != 0 {
		t.Errorf("Expected 0 sections, got %d", len(result.Sections))
	}

	// All tasks are orphaned
	if result.OrphanedTasks != 3 {
		t.Errorf("Expected 3 orphaned tasks, got %d", result.OrphanedTasks)
	}

	// Sequential is vacuously true when no sections
	if !result.SequentialNumbers {
		t.Error("Expected sequential numbers to be true (vacuously)")
	}

	// No gaps
	if len(result.NonSequentialGaps) != 0 {
		t.Errorf("Expected 0 non-sequential gaps, got %d", len(result.NonSequentialGaps))
	}
}

func TestValidateTasksStructure_OrphanedTasks(t *testing.T) {
	content := `# Tasks

- [ ] Orphan task 1
- [ ] Orphan task 2

## 1. Implementation
- [ ] 1.1 Create schema
- [x] 1.2 Implement API

## 2. Testing
- [ ] 2.1 Write tests
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure failed: %v", err)
	}

	// Should have 2 sections
	if len(result.Sections) != 2 {
		t.Errorf("Expected 2 sections, got %d", len(result.Sections))
	}

	// 2 orphaned tasks before the first section
	if result.OrphanedTasks != 2 {
		t.Errorf("Expected 2 orphaned tasks, got %d", result.OrphanedTasks)
	}

	// First section has 2 tasks
	if result.Sections[0].TaskCount != 2 {
		t.Errorf("Expected section 1 to have 2 tasks, got %d", result.Sections[0].TaskCount)
	}

	// Second section has 1 task
	if result.Sections[1].TaskCount != 1 {
		t.Errorf("Expected section 2 to have 1 task, got %d", result.Sections[1].TaskCount)
	}

	// Sequential numbers
	if !result.SequentialNumbers {
		t.Error("Expected sequential numbers to be true")
	}
}

func TestValidateTasksStructure_EmptySections(t *testing.T) {
	content := `## 1. Planning
Some description text here.

## 2. Implementation
- [ ] 2.1 Create schema
- [x] 2.2 Implement API

## 3. Documentation
No tasks in this section either.

## 4. Testing
- [ ] 4.1 Write tests
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure failed: %v", err)
	}

	// Should have 4 sections
	if len(result.Sections) != 4 {
		t.Errorf("Expected 4 sections, got %d", len(result.Sections))
	}

	// 2 empty sections: Planning and Documentation
	if len(result.EmptySections) != 2 {
		t.Errorf(
			"Expected 2 empty sections, got %d: %v",
			len(result.EmptySections),
			result.EmptySections,
		)
	}

	// Verify the empty section names
	expectedEmpty := map[string]bool{"Planning": true, "Documentation": true}
	for _, name := range result.EmptySections {
		if !expectedEmpty[name] {
			t.Errorf("Unexpected empty section: %q", name)
		}
	}

	// Sequential numbers
	if !result.SequentialNumbers {
		t.Error("Expected sequential numbers to be true")
	}
}

func TestValidateTasksStructure_NonSequentialNumbers(t *testing.T) {
	content := `## 1. First Section
- [ ] Task 1

## 3. Third Section
- [ ] Task 3

## 5. Fifth Section
- [ ] Task 5
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure failed: %v", err)
	}

	// Should have 3 sections
	if len(result.Sections) != 3 {
		t.Errorf("Expected 3 sections, got %d", len(result.Sections))
	}

	// Non-sequential numbers
	if result.SequentialNumbers {
		t.Error("Expected sequential numbers to be false")
	}

	// Gaps should be 2 and 4
	if len(result.NonSequentialGaps) != 2 {
		t.Errorf(
			"Expected 2 gaps, got %d: %v",
			len(result.NonSequentialGaps),
			result.NonSequentialGaps,
		)
	}

	expectedGaps := map[int]bool{2: true, 4: true}
	for _, gap := range result.NonSequentialGaps {
		if !expectedGaps[gap] {
			t.Errorf("Unexpected gap: %d", gap)
		}
	}
}

func TestValidateTasksStructure_SingleSection(t *testing.T) {
	content := `## 1. Only Section
- [ ] 1.1 Task one
- [x] 1.2 Task two
- [ ] 1.3 Task three
`
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure failed: %v", err)
	}

	// Should have 1 section
	if len(result.Sections) != 1 {
		t.Errorf("Expected 1 section, got %d", len(result.Sections))
	}

	// Section should have 3 tasks
	if result.Sections[0].TaskCount != 3 {
		t.Errorf("Expected section to have 3 tasks, got %d", result.Sections[0].TaskCount)
	}

	// No orphaned tasks
	if result.OrphanedTasks != 0 {
		t.Errorf("Expected 0 orphaned tasks, got %d", result.OrphanedTasks)
	}

	// Sequential (only one section starting at 1)
	if !result.SequentialNumbers {
		t.Error("Expected sequential numbers to be true")
	}

	// No empty sections
	if len(result.EmptySections) != 0 {
		t.Errorf("Expected 0 empty sections, got %d", len(result.EmptySections))
	}
}

func TestValidateTasksStructure_NonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.md")

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure should not error on missing file: %v", err)
	}

	// Should return empty result
	if len(result.Sections) != 0 {
		t.Errorf("Expected 0 sections for nonexistent file, got %d", len(result.Sections))
	}
	if result.OrphanedTasks != 0 {
		t.Errorf("Expected 0 orphaned tasks for nonexistent file, got %d", result.OrphanedTasks)
	}
	if len(result.EmptySections) != 0 {
		t.Errorf(
			"Expected 0 empty sections for nonexistent file, got %d",
			len(result.EmptySections),
		)
	}
	if !result.SequentialNumbers {
		t.Error("Expected sequential numbers to be true for nonexistent file")
	}
	if len(result.NonSequentialGaps) != 0 {
		t.Errorf("Expected 0 gaps for nonexistent file, got %d", len(result.NonSequentialGaps))
	}
}

func TestValidateTasksStructure_EmptyFile(t *testing.T) {
	content := ""
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.md")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := ValidateTasksStructure(filePath)
	if err != nil {
		t.Fatalf("ValidateTasksStructure failed: %v", err)
	}

	// Should return empty result
	if len(result.Sections) != 0 {
		t.Errorf("Expected 0 sections for empty file, got %d", len(result.Sections))
	}
	if result.OrphanedTasks != 0 {
		t.Errorf("Expected 0 orphaned tasks for empty file, got %d", result.OrphanedTasks)
	}
	if len(result.EmptySections) != 0 {
		t.Errorf("Expected 0 empty sections for empty file, got %d", len(result.EmptySections))
	}
	if !result.SequentialNumbers {
		t.Error("Expected sequential numbers to be true for empty file")
	}
	if len(result.NonSequentialGaps) != 0 {
		t.Errorf("Expected 0 gaps for empty file, got %d", len(result.NonSequentialGaps))
	}
}
