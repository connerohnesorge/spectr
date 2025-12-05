package accept

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// =============================================================================
// CalculateSummary Tests
// =============================================================================

func TestCalculateSummary_Empty(t *testing.T) {
	var sections []Section
	summary := CalculateSummary(sections)

	if summary.Total != 0 {
		t.Errorf("Expected total 0, got %d", summary.Total)
	}
	if summary.Completed != 0 {
		t.Errorf("Expected completed 0, got %d", summary.Completed)
	}
}

func TestCalculateSummary_FlatTasks(t *testing.T) {
	sections := []Section{
		{
			Number: 1,
			Name:   "Test Section",
			Tasks: []Task{
				{ID: "1.1", Description: "Task 1", Completed: true},
				{ID: "1.2", Description: "Task 2", Completed: false},
				{ID: "1.3", Description: "Task 3", Completed: true},
			},
		},
	}
	summary := CalculateSummary(sections)

	if summary.Total != 3 {
		t.Errorf("Expected total 3, got %d", summary.Total)
	}
	if summary.Completed != 2 {
		t.Errorf("Expected completed 2, got %d", summary.Completed)
	}
}

func TestCalculateSummary_NestedTasks(t *testing.T) {
	sections := []Section{
		{
			Number: 1,
			Name:   "Nested Section",
			Tasks: []Task{
				{
					ID:          "1.1",
					Description: "Parent task",
					Completed:   true,
					Subtasks: []Task{
						{ID: "1.1.1", Description: "Child 1", Completed: true},
						{
							ID:          "1.1.2",
							Description: "Child 2",
							Completed:   false,
							Subtasks: []Task{
								{ID: "1.1.2.1", Description: "Grandchild", Completed: true},
							},
						},
					},
				},
			},
		},
	}
	summary := CalculateSummary(sections)

	// 1 parent + 2 children + 1 grandchild = 4 total
	if summary.Total != 4 {
		t.Errorf("Expected total 4, got %d", summary.Total)
	}
	// parent(1) + child1(1) + grandchild(1) = 3 completed
	if summary.Completed != 3 {
		t.Errorf("Expected completed 3, got %d", summary.Completed)
	}
}

func TestCalculateSummary_MixedCompleted(t *testing.T) {
	sections := []Section{
		{
			Number: 1,
			Name:   "Section 1",
			Tasks: []Task{
				{ID: "1.1", Description: "Done", Completed: true},
				{ID: "1.2", Description: "Not done", Completed: false},
			},
		},
		{
			Number: 2,
			Name:   "Section 2",
			Tasks: []Task{
				{ID: "2.1", Description: "Done", Completed: true},
				{ID: "2.2", Description: "Done", Completed: true},
				{ID: "2.3", Description: "Not done", Completed: false},
			},
		},
	}
	summary := CalculateSummary(sections)

	if summary.Total != 5 {
		t.Errorf("Expected total 5, got %d", summary.Total)
	}
	if summary.Completed != 3 {
		t.Errorf("Expected completed 3, got %d", summary.Completed)
	}
}

// =============================================================================
// countTasks Tests
// =============================================================================

func TestCountTasks_Empty(t *testing.T) {
	var tasks []Task
	total, completed := countTasks(tasks)

	if total != 0 {
		t.Errorf("Expected total 0, got %d", total)
	}
	if completed != 0 {
		t.Errorf("Expected completed 0, got %d", completed)
	}
}

func TestCountTasks_Simple(t *testing.T) {
	tasks := []Task{
		{ID: "1", Description: "Task 1", Completed: false},
		{ID: "2", Description: "Task 2", Completed: true},
		{ID: "3", Description: "Task 3", Completed: false},
	}
	total, completed := countTasks(tasks)

	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	if completed != 1 {
		t.Errorf("Expected completed 1, got %d", completed)
	}
}

func TestCountTasks_WithSubtasks(t *testing.T) {
	tasks := []Task{
		{
			ID:          "1",
			Description: "Parent",
			Completed:   true,
			Subtasks: []Task{
				{ID: "1.1", Description: "Child 1", Completed: false},
				{ID: "1.2", Description: "Child 2", Completed: true},
			},
		},
	}
	total, completed := countTasks(tasks)

	// 1 parent + 2 subtasks = 3 total
	if total != 3 {
		t.Errorf("Expected total 3, got %d", total)
	}
	// parent(1) + child2(1) = 2 completed
	if completed != 2 {
		t.Errorf("Expected completed 2, got %d", completed)
	}
}

// =============================================================================
// WriteTasksJSON Tests
// =============================================================================

func TestWriteTasksJSON_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	sections := []Section{
		{Number: 1, Name: "Test", Tasks: nil},
	}

	err := WriteTasksJSON(filePath, "test-change", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		t.Error("Expected file to be created")
	}
}

func TestWriteTasksJSON_ValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	sections := []Section{
		{
			Number: 1,
			Name:   "Implementation",
			Tasks: []Task{
				{ID: "1.1", Description: "Task 1", Completed: true},
			},
		},
	}

	err := WriteTasksJSON(filePath, "test-change", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Errorf("Output is not valid JSON: %v", err)
	}
}

func TestWriteTasksJSON_CorrectFields(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	sections := []Section{
		{
			Number: 1,
			Name:   "Test Section",
			Tasks: []Task{
				{ID: "1.1", Description: "Task", Completed: true},
				{ID: "1.2", Description: "Task 2", Completed: false},
			},
		},
	}

	err := WriteTasksJSON(filePath, "my-change-id", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result TasksJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Check version
	if result.Version != "1.0" {
		t.Errorf("Expected version '1.0', got %q", result.Version)
	}

	// Check changeId
	if result.ChangeID != "my-change-id" {
		t.Errorf("Expected changeId 'my-change-id', got %q", result.ChangeID)
	}

	// Check acceptedAt is valid RFC3339
	_, err = time.Parse(time.RFC3339, result.AcceptedAt)
	if err != nil {
		t.Errorf("AcceptedAt is not valid RFC3339: %q, error: %v", result.AcceptedAt, err)
	}

	// Check sections
	if len(result.Sections) != 1 {
		t.Errorf("Expected 1 section, got %d", len(result.Sections))
	}
	if result.Sections[0].Name != "Test Section" {
		t.Errorf("Expected section name 'Test Section', got %q", result.Sections[0].Name)
	}

	// Check summary
	if result.Summary.Total != 2 {
		t.Errorf("Expected summary total 2, got %d", result.Summary.Total)
	}
	if result.Summary.Completed != 1 {
		t.Errorf("Expected summary completed 1, got %d", result.Summary.Completed)
	}
}

func TestWriteTasksJSON_PrettyPrinted(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	sections := []Section{
		{Number: 1, Name: "Test", Tasks: nil},
	}

	err := WriteTasksJSON(filePath, "test", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	content := string(data)

	// Check for 2-space indentation by looking for indented fields
	if !strings.Contains(content, "  \"version\"") {
		t.Error("Expected 2-space indentation for top-level fields")
	}

	// Verify it's not minified (contains newlines between fields)
	lines := strings.Split(content, "\n")
	if len(lines) < 5 {
		t.Error("Expected pretty-printed output with multiple lines")
	}
}

func TestWriteTasksJSON_TrailingNewline(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	sections := []Section{
		{Number: 1, Name: "Test", Tasks: nil},
	}

	err := WriteTasksJSON(filePath, "test", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("File is empty")
	}

	if data[len(data)-1] != '\n' {
		t.Error("Expected file to end with newline")
	}
}

func TestWriteTasksJSON_AtomicWrite(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	sections := []Section{
		{Number: 1, Name: "Test", Tasks: nil},
	}

	err := WriteTasksJSON(filePath, "test", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	// Check that no temp files remain
	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for _, entry := range entries {
		if strings.Contains(entry.Name(), ".tmp") {
			t.Errorf("Temp file was not cleaned up: %s", entry.Name())
		}
	}

	// Should only have the final tasks.json file
	if len(entries) != 1 {
		t.Errorf("Expected 1 file in directory, got %d", len(entries))
	}
	if entries[0].Name() != "tasks.json" {
		t.Errorf("Expected file 'tasks.json', got %q", entries[0].Name())
	}
}

// =============================================================================
// Integration Tests with Fixtures
// =============================================================================

func TestWriteTasksJSON_SimpleFixture(t *testing.T) {
	// Parse the simple.md fixture
	fixtureDir := filepath.Join("testdata", "simple.md")
	sections, err := ParseTasksFile(fixtureDir)
	if err != nil {
		t.Fatalf("Failed to parse simple.md: %v", err)
	}

	// Write to JSON
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "tasks.json")

	err = WriteTasksJSON(outputPath, "add-naming-philosophy-note", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	// Read and verify
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	var result TasksJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify structure
	if result.Version != "1.0" {
		t.Errorf("Expected version '1.0', got %q", result.Version)
	}
	if result.ChangeID != "add-naming-philosophy-note" {
		t.Errorf("Expected changeId 'add-naming-philosophy-note', got %q", result.ChangeID)
	}

	// simple.md has 3 sections
	if len(result.Sections) != 3 {
		t.Errorf("Expected 3 sections, got %d", len(result.Sections))
	}

	// Verify section names
	expectedSections := []string{"Documentation Updates", "Validation", "Final Review"}
	for i, expected := range expectedSections {
		if i >= len(result.Sections) {
			break
		}
		if result.Sections[i].Name != expected {
			t.Errorf("Section %d: expected name %q, got %q", i+1, expected, result.Sections[i].Name)
		}
	}

	// All tasks in simple.md are completed
	if result.Summary.Total != result.Summary.Completed {
		t.Errorf(
			"Expected all tasks completed, got %d/%d",
			result.Summary.Completed,
			result.Summary.Total,
		)
	}
}

func TestWriteTasksJSON_ComplexFixture(t *testing.T) {
	// Parse the complex.md fixture
	fixtureDir := filepath.Join("testdata", "complex.md")
	sections, err := ParseTasksFile(fixtureDir)
	if err != nil {
		t.Fatalf("Failed to parse complex.md: %v", err)
	}

	// Write to JSON
	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "tasks.json")

	err = WriteTasksJSON(outputPath, "complex-refactor", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	// Read and verify
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("Failed to read output: %v", err)
	}

	var result TasksJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// complex.md has 8 sections
	if len(result.Sections) != 8 {
		t.Errorf("Expected 8 sections, got %d", len(result.Sections))
	}

	// Verify summary counts are reasonable
	// complex.md has many tasks, including nested subtasks
	if result.Summary.Total == 0 {
		t.Error("Expected non-zero total tasks")
	}

	// Verify that we have both completed and incomplete tasks
	// (complex.md has some tasks still marked incomplete)
	if result.Summary.Completed == 0 {
		t.Error("Expected some completed tasks")
	}

	// Verify first section has tasks
	if len(result.Sections) > 0 {
		section1 := result.Sections[0]
		if section1.Name != "Foundation: Create New Abstractions" {
			t.Errorf(
				"Expected first section name 'Foundation: Create New Abstractions', got %q",
				section1.Name,
			)
		}
		if len(section1.Tasks) == 0 {
			t.Error("Expected first section to have tasks")
		}
	}

	// Verify JSON structure is valid by checking a few known fields
	if result.Version != "1.0" {
		t.Errorf("Expected version '1.0', got %q", result.Version)
	}
	if result.ChangeID != "complex-refactor" {
		t.Errorf("Expected changeId 'complex-refactor', got %q", result.ChangeID)
	}
}

// =============================================================================
// Error Handling Tests
// =============================================================================

func TestWriteTasksJSON_InvalidPath(t *testing.T) {
	sections := []Section{
		{Number: 1, Name: "Test", Tasks: nil},
	}

	// Try to write to a non-existent directory
	err := WriteTasksJSON("/nonexistent/directory/tasks.json", "test", sections)
	if err == nil {
		t.Error("Expected error when writing to invalid path")
	}
}

func TestWriteTasksJSON_PermissionDenied(t *testing.T) {
	// Skip on systems where we can't test permissions properly
	if os.Getuid() == 0 {
		t.Skip("Skipping permission test when running as root")
	}

	tmpDir := t.TempDir()

	// Create a read-only directory
	readOnlyDir := filepath.Join(tmpDir, "readonly")
	if err := os.Mkdir(readOnlyDir, 0444); err != nil {
		t.Fatalf("Failed to create read-only directory: %v", err)
	}
	// Ensure cleanup
	defer func() { _ = os.Chmod(readOnlyDir, 0755) }()

	sections := []Section{
		{Number: 1, Name: "Test", Tasks: nil},
	}

	err := WriteTasksJSON(filepath.Join(readOnlyDir, "tasks.json"), "test", sections)
	if err == nil {
		t.Error("Expected error when writing to read-only directory")
	}
}

// =============================================================================
// Edge Cases
// =============================================================================

func TestWriteTasksJSON_EmptySections(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	var sections []Section

	err := WriteTasksJSON(filePath, "empty-change", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result TasksJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if len(result.Sections) != 0 {
		t.Errorf("Expected 0 sections, got %d", len(result.Sections))
	}
	if result.Summary.Total != 0 {
		t.Errorf("Expected summary total 0, got %d", result.Summary.Total)
	}
}

func TestWriteTasksJSON_SpecialCharacters(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	sections := []Section{
		{
			Number: 1,
			Name:   "Test \"quotes\" and <brackets>",
			Tasks: []Task{
				{
					ID:          "1.1",
					Description: "Task with special chars: \"quotes\", <angle>, &amp;",
					Completed:   false,
				},
			},
		},
	}

	err := WriteTasksJSON(filePath, "special-chars", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result TasksJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal JSON with special characters: %v", err)
	}

	// Verify special characters are preserved
	if result.Sections[0].Name != "Test \"quotes\" and <brackets>" {
		t.Errorf("Special characters not preserved in section name: %q", result.Sections[0].Name)
	}
}

func TestWriteTasksJSON_LongTaskDescription(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	longDesc := strings.Repeat("This is a very long task description. ", 100)

	sections := []Section{
		{
			Number: 1,
			Name:   "Long Description Test",
			Tasks: []Task{
				{ID: "1.1", Description: longDesc, Completed: false},
			},
		},
	}

	err := WriteTasksJSON(filePath, "long-desc", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result TasksJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	if result.Sections[0].Tasks[0].Description != longDesc {
		t.Error("Long description was not preserved correctly")
	}
}

func TestWriteTasksJSON_DeeplyNestedTasks(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "tasks.json")

	// Create a deeply nested task structure
	sections := []Section{
		{
			Number: 1,
			Name:   "Deep Nesting",
			Tasks: []Task{
				{
					ID:          "1",
					Description: "Level 1",
					Completed:   true,
					Subtasks: []Task{
						{
							ID:          "1.1",
							Description: "Level 2",
							Completed:   true,
							Subtasks: []Task{
								{
									ID:          "1.1.1",
									Description: "Level 3",
									Completed:   false,
									Subtasks: []Task{
										{
											ID:          "1.1.1.1",
											Description: "Level 4",
											Completed:   true,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	err := WriteTasksJSON(filePath, "deep-nest", sections)
	if err != nil {
		t.Fatalf("WriteTasksJSON failed: %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result TasksJSON
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify total count: 4 tasks total
	if result.Summary.Total != 4 {
		t.Errorf("Expected 4 total tasks, got %d", result.Summary.Total)
	}
	// Verify completed count: 3 completed (levels 1, 2, 4)
	if result.Summary.Completed != 3 {
		t.Errorf("Expected 3 completed tasks, got %d", result.Summary.Completed)
	}

	// Verify nesting structure
	level1 := result.Sections[0].Tasks[0]
	if level1.ID != "1" {
		t.Errorf("Expected level 1 ID '1', got %q", level1.ID)
	}
	if len(level1.Subtasks) != 1 {
		t.Fatalf("Expected 1 subtask at level 1, got %d", len(level1.Subtasks))
	}

	level2 := level1.Subtasks[0]
	if level2.ID != "1.1" {
		t.Errorf("Expected level 2 ID '1.1', got %q", level2.ID)
	}

	level3 := level2.Subtasks[0]
	if level3.ID != "1.1.1" {
		t.Errorf("Expected level 3 ID '1.1.1', got %q", level3.ID)
	}

	level4 := level3.Subtasks[0]
	if level4.ID != "1.1.1.1" {
		t.Errorf("Expected level 4 ID '1.1.1.1', got %q", level4.ID)
	}
}

// =============================================================================
// Table-Driven Tests
// =============================================================================

func TestCalculateSummary_TableDriven(t *testing.T) {
	tests := []struct {
		name              string
		sections          []Section
		expectedTotal     int
		expectedCompleted int
	}{
		{
			name:              "Empty sections",
			sections:          nil,
			expectedTotal:     0,
			expectedCompleted: 0,
		},
		{
			name: "Single section with tasks",
			sections: []Section{
				{
					Number: 1,
					Name:   "Test",
					Tasks: []Task{
						{ID: "1.1", Completed: true},
						{ID: "1.2", Completed: false},
					},
				},
			},
			expectedTotal:     2,
			expectedCompleted: 1,
		},
		{
			name: "Multiple sections",
			sections: []Section{
				{
					Number: 1,
					Name:   "Section 1",
					Tasks: []Task{
						{ID: "1.1", Completed: true},
					},
				},
				{
					Number: 2,
					Name:   "Section 2",
					Tasks: []Task{
						{ID: "2.1", Completed: true},
						{ID: "2.2", Completed: true},
					},
				},
			},
			expectedTotal:     3,
			expectedCompleted: 3,
		},
		{
			name: "Section with nested tasks",
			sections: []Section{
				{
					Number: 1,
					Name:   "Nested",
					Tasks: []Task{
						{
							ID:        "1.1",
							Completed: false,
							Subtasks: []Task{
								{ID: "1.1.1", Completed: true},
								{ID: "1.1.2", Completed: false},
							},
						},
					},
				},
			},
			expectedTotal:     3,
			expectedCompleted: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := CalculateSummary(tt.sections)
			if summary.Total != tt.expectedTotal {
				t.Errorf("Total: expected %d, got %d", tt.expectedTotal, summary.Total)
			}
			if summary.Completed != tt.expectedCompleted {
				t.Errorf("Completed: expected %d, got %d", tt.expectedCompleted, summary.Completed)
			}
		})
	}
}

func TestCountTasks_TableDriven(t *testing.T) {
	tests := []struct {
		name              string
		tasks             []Task
		expectedTotal     int
		expectedCompleted int
	}{
		{
			name:              "Empty tasks",
			tasks:             nil,
			expectedTotal:     0,
			expectedCompleted: 0,
		},
		{
			name: "All completed",
			tasks: []Task{
				{ID: "1", Completed: true},
				{ID: "2", Completed: true},
			},
			expectedTotal:     2,
			expectedCompleted: 2,
		},
		{
			name: "None completed",
			tasks: []Task{
				{ID: "1", Completed: false},
				{ID: "2", Completed: false},
			},
			expectedTotal:     2,
			expectedCompleted: 0,
		},
		{
			name: "With deep subtasks",
			tasks: []Task{
				{
					ID:        "1",
					Completed: true,
					Subtasks: []Task{
						{
							ID:        "1.1",
							Completed: false,
							Subtasks: []Task{
								{ID: "1.1.1", Completed: true},
							},
						},
					},
				},
			},
			expectedTotal:     3,
			expectedCompleted: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			total, completed := countTasks(tt.tasks)
			if total != tt.expectedTotal {
				t.Errorf("Total: expected %d, got %d", tt.expectedTotal, total)
			}
			if completed != tt.expectedCompleted {
				t.Errorf("Completed: expected %d, got %d", tt.expectedCompleted, completed)
			}
		})
	}
}
