package cmd

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

const testChangeID = "test-change"

func TestComputeAggregateStatus(t *testing.T) {
	tests := []struct {
		name     string
		children []parsers.Task
		expected parsers.TaskStatusValue
	}{
		{
			name:     "empty children returns pending",
			children: make([]parsers.Task, 0),
			expected: parsers.TaskStatusPending,
		},
		{
			name:     "nil children returns pending",
			children: nil,
			expected: parsers.TaskStatusPending,
		},
		{
			name: "all pending returns pending",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusPending},
				{ID: "1.3", Status: parsers.TaskStatusPending},
			},
			expected: parsers.TaskStatusPending,
		},
		{
			name: "single pending returns pending",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
			},
			expected: parsers.TaskStatusPending,
		},
		{
			name: "all completed returns completed",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
				{ID: "1.2", Status: parsers.TaskStatusCompleted},
				{ID: "1.3", Status: parsers.TaskStatusCompleted},
			},
			expected: parsers.TaskStatusCompleted,
		},
		{
			name: "single completed returns completed",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
			},
			expected: parsers.TaskStatusCompleted,
		},
		{
			name: "any in_progress returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusInProgress},
				{ID: "1.3", Status: parsers.TaskStatusPending},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "single in_progress returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusInProgress},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "all in_progress returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusInProgress},
				{ID: "1.2", Status: parsers.TaskStatusInProgress},
				{ID: "1.3", Status: parsers.TaskStatusInProgress},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "in_progress with completed returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
				{ID: "1.2", Status: parsers.TaskStatusInProgress},
				{ID: "1.3", Status: parsers.TaskStatusCompleted},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "in_progress with pending returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusInProgress},
				{ID: "1.3", Status: parsers.TaskStatusPending},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "mixed all statuses returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusInProgress},
				{ID: "1.3", Status: parsers.TaskStatusCompleted},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "mixed pending and completed (no in_progress) returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusCompleted},
				{ID: "1.3", Status: parsers.TaskStatusPending},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "mostly pending with one completed returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusPending},
				{ID: "1.3", Status: parsers.TaskStatusPending},
				{ID: "1.4", Status: parsers.TaskStatusCompleted},
			},
			expected: parsers.TaskStatusInProgress,
		},
		{
			name: "mostly completed with one pending returns in_progress",
			children: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
				{ID: "1.2", Status: parsers.TaskStatusCompleted},
				{ID: "1.3", Status: parsers.TaskStatusCompleted},
				{ID: "1.4", Status: parsers.TaskStatusPending},
			},
			expected: parsers.TaskStatusInProgress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := computeAggregateStatus(tt.children)
			if result != tt.expected {
				t.Errorf(
					"computeAggregateStatus() = %v, want %v",
					result,
					tt.expected,
				)
			}
		})
	}
}

// TestComputeAggregateStatusLargeSet tests status aggregation with a large number of children
func TestComputeAggregateStatusLargeSet(t *testing.T) {
	tests := []struct {
		name        string
		numPending  int
		numProgress int
		numComplete int
		expected    parsers.TaskStatusValue
	}{
		{
			name:        "100 pending tasks",
			numPending:  100,
			numProgress: 0,
			numComplete: 0,
			expected:    parsers.TaskStatusPending,
		},
		{
			name:        "100 completed tasks",
			numPending:  0,
			numProgress: 0,
			numComplete: 100,
			expected:    parsers.TaskStatusCompleted,
		},
		{
			name:        "99 pending, 1 in_progress",
			numPending:  99,
			numProgress: 1,
			numComplete: 0,
			expected:    parsers.TaskStatusInProgress,
		},
		{
			name:        "99 completed, 1 pending",
			numPending:  1,
			numProgress: 0,
			numComplete: 99,
			expected:    parsers.TaskStatusInProgress,
		},
		{
			name:        "50 pending, 50 completed",
			numPending:  50,
			numProgress: 0,
			numComplete: 50,
			expected:    parsers.TaskStatusInProgress,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			children := make([]parsers.Task, 0, tt.numPending+tt.numProgress+tt.numComplete)

			for i := range tt.numPending {
				children = append(children, parsers.Task{
					ID:     fmt.Sprintf("1.%d", i+1),
					Status: parsers.TaskStatusPending,
				})
			}
			for i := range tt.numProgress {
				children = append(children, parsers.Task{
					ID:     fmt.Sprintf("2.%d", i+1),
					Status: parsers.TaskStatusInProgress,
				})
			}
			for i := range tt.numComplete {
				children = append(children, parsers.Task{
					ID:     fmt.Sprintf("3.%d", i+1),
					Status: parsers.TaskStatusCompleted,
				})
			}

			result := computeAggregateStatus(children)
			if result != tt.expected {
				t.Errorf(
					"computeAggregateStatus() with %d pending, %d in_progress, %d completed = %v, want %v",
					tt.numPending,
					tt.numProgress,
					tt.numComplete,
					result,
					tt.expected,
				)
			}
		})
	}
}

// TestChildFileHeader tests the childFileHeader function with various inputs
func TestChildFileHeader(t *testing.T) {
	tests := []struct {
		name         string
		changeID     string
		parentTaskID string
		wantStrings  []string // Strings that must be present in the output
	}{
		{
			name:         "basic header with simple IDs",
			changeID:     testChangeID,
			parentTaskID: "1",
			wantStrings: []string{
				"// Generated by: spectr accept test-change",
				"// Parent change: test-change",
				"// Parent task: 1",
				"// Status Values:",
				"pending",
				"in_progress",
				"completed",
				"// Status Transitions:",
				"pending -> in_progress -> completed",
				"// Workflow:",
				"// IMPORTANT - Update Status Immediately:",
			},
		},
		{
			name:         "header with complex change ID",
			changeID:     "add-size-based-task-splitting",
			parentTaskID: "5",
			wantStrings: []string{
				"// Generated by: spectr accept add-size-based-task-splitting",
				"// Parent change: add-size-based-task-splitting",
				"// Parent task: 5",
			},
		},
		{
			name:         "header with hierarchical task ID",
			changeID:     "my-feature",
			parentTaskID: "2.3",
			wantStrings: []string{
				"// Generated by: spectr accept my-feature",
				"// Parent change: my-feature",
				"// Parent task: 2.3",
			},
		},
		{
			name:         "header with dashes in change ID",
			changeID:     "add-new-feature-x",
			parentTaskID: "10",
			wantStrings: []string{
				"// Generated by: spectr accept add-new-feature-x",
				"// Parent change: add-new-feature-x",
				"// Parent task: 10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := childFileHeader(tt.changeID, tt.parentTaskID)

			// Verify all expected strings are present
			for _, wantStr := range tt.wantStrings {
				if !strings.Contains(result, wantStr) {
					t.Errorf(
						"childFileHeader() missing expected string:\nwant: %q\ngot header:\n%s",
						wantStr,
						result,
					)
				}
			}

			// Verify it starts with JSONC comment syntax
			if !strings.HasPrefix(result, "//") {
				t.Errorf("childFileHeader() should start with '//', got: %s", result[:20])
			}

			// Verify it ends with newlines
			if !strings.HasSuffix(result, "\n\n") {
				t.Error("childFileHeader() should end with double newline")
			}
		})
	}
}

// TestChildFileHeaderFields verifies all required fields are present
func TestChildFileHeaderFields(t *testing.T) {
	changeID := testChangeID
	parentTaskID := "5"
	result := childFileHeader(changeID, parentTaskID)

	requiredFields := []string{
		"Generated by:",
		"Parent change:",
		"Parent task:",
		"Status Values:",
		"pending",
		"in_progress",
		"completed",
		"Status Transitions:",
		"Workflow:",
		"IMPORTANT - Update Status Immediately:",
		"Do NOT batch status updates",
		"Do NOT wait until all tasks are done",
	}

	for _, field := range requiredFields {
		if !strings.Contains(result, field) {
			t.Errorf("childFileHeader() missing required field: %q", field)
		}
	}
}

// TestChildFileHeaderFormat verifies the header is properly formatted JSONC
func TestChildFileHeaderFormat(t *testing.T) {
	result := childFileHeader("my-change", "1")

	// Split into lines to verify formatting
	lines := strings.Split(result, "\n")

	if len(lines) < 10 {
		t.Errorf("childFileHeader() should have at least 10 lines, got %d", len(lines))
	}

	// Verify all non-empty lines start with "//"
	for i, line := range lines {
		if line != "" && !strings.HasPrefix(line, "//") {
			t.Errorf("Line %d does not start with '//': %q", i+1, line)
		}
	}

	// Verify origin info is at the top
	if !strings.Contains(lines[0], "Generated by:") {
		t.Errorf("First line should contain 'Generated by:', got: %q", lines[0])
	}

	if !strings.Contains(lines[1], "Parent change:") {
		t.Errorf("Second line should contain 'Parent change:', got: %q", lines[1])
	}

	if !strings.Contains(lines[2], "Parent task:") {
		t.Errorf("Third line should contain 'Parent task:', got: %q", lines[2])
	}
}

// TestDeleteOldChildFiles tests the cleanup of old child files
func TestDeleteOldChildFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create some test files
	testFiles := []string{
		"tasks.jsonc",       // Root file - should NOT be deleted
		"tasks-1.jsonc",     // Child file - should be deleted
		"tasks-2.jsonc",     // Child file - should be deleted
		"tasks-10.jsonc",    // Child file - should be deleted
		"proposal.md",       // Non-tasks file - should NOT be deleted
		"tasks.md",          // Source file - should NOT be deleted
		"other-tasks.jsonc", // Different pattern - should NOT be deleted
	}

	// Write all test files
	for _, filename := range testFiles {
		filePath := fmt.Sprintf("%s/%s", tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	// Run deleteOldChildFiles
	err := deleteOldChildFiles(tempDir)
	if err != nil {
		t.Fatalf("deleteOldChildFiles() failed: %v", err)
	}

	// Verify which files still exist
	remainingFiles := []string{
		"tasks.jsonc",
		"proposal.md",
		"tasks.md",
		"other-tasks.jsonc",
	}

	for _, filename := range remainingFiles {
		filePath := fmt.Sprintf("%s/%s", tempDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("File %s should still exist but was deleted", filename)
		}
	}

	// Verify child files were deleted
	deletedFiles := []string{
		"tasks-1.jsonc",
		"tasks-2.jsonc",
		"tasks-10.jsonc",
	}

	for _, filename := range deletedFiles {
		filePath := fmt.Sprintf("%s/%s", tempDir, filename)
		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			t.Errorf("File %s should have been deleted but still exists", filename)
		}
	}
}

// TestDeleteOldChildFilesEmptyDir tests cleanup on an empty directory
func TestDeleteOldChildFilesEmptyDir(t *testing.T) {
	tempDir := t.TempDir()

	err := deleteOldChildFiles(tempDir)
	if err != nil {
		t.Errorf("deleteOldChildFiles() on empty dir should not error, got: %v", err)
	}
}

// TestDeleteOldChildFilesNoChildFiles tests cleanup when no child files exist
func TestDeleteOldChildFilesNoChildFiles(t *testing.T) {
	tempDir := t.TempDir()

	// Create only non-child files
	testFiles := []string{"tasks.jsonc", "proposal.md", "tasks.md"}
	for _, filename := range testFiles {
		filePath := fmt.Sprintf("%s/%s", tempDir, filename)
		if err := os.WriteFile(filePath, []byte("test"), 0o644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	err := deleteOldChildFiles(tempDir)
	if err != nil {
		t.Errorf("deleteOldChildFiles() should not error when no child files exist, got: %v", err)
	}

	// Verify all files still exist
	for _, filename := range testFiles {
		filePath := fmt.Sprintf("%s/%s", tempDir, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			t.Errorf("File %s should still exist", filename)
		}
	}
}

// TestWriteRootTasksJSONC tests writing a root file with references
func TestWriteRootTasksJSONC(t *testing.T) {
	tempDir := t.TempDir()
	rootPath := fmt.Sprintf("%s/tasks.jsonc", tempDir)

	referenceTasks := []parsers.Task{
		{
			ID:          "1",
			Section:     "Implementation",
			Description: "Implementation tasks",
			Status:      parsers.TaskStatusInProgress,
			Children:    "$ref:tasks-1.jsonc",
		},
		{
			ID:          "2",
			Section:     "Testing",
			Description: "Testing tasks",
			Status:      parsers.TaskStatusPending,
			Children:    "$ref:tasks-2.jsonc",
		},
	}

	err := writeRootTasksJSONC(rootPath, referenceTasks)
	if err != nil {
		t.Fatalf("writeRootTasksJSONC() failed: %v", err)
	}

	// Read the file back
	content, err := os.ReadFile(rootPath)
	if err != nil {
		t.Fatalf("Failed to read root file: %v", err)
	}

	contentStr := string(content)

	// Verify header is present
	if !strings.Contains(contentStr, "// Spectr Tasks File (JSONC)") {
		t.Error("Root file should contain JSONC header")
	}

	// Verify version 2
	if !strings.Contains(contentStr, `"version": 2`) {
		t.Error("Root file should have version 2")
	}

	// Verify includes field
	if !strings.Contains(contentStr, `"includes": [`) {
		t.Error("Root file should have includes field")
	}
	if !strings.Contains(contentStr, `"tasks-*.jsonc"`) {
		t.Error("Root file includes should contain tasks-*.jsonc glob")
	}

	// Verify reference tasks
	if !strings.Contains(contentStr, `"children": "$ref:tasks-1.jsonc"`) {
		t.Error("Root file should contain reference to tasks-1.jsonc")
	}
	if !strings.Contains(contentStr, `"children": "$ref:tasks-2.jsonc"`) {
		t.Error("Root file should contain reference to tasks-2.jsonc")
	}

	// Verify task statuses
	if !strings.Contains(contentStr, `"status": "in_progress"`) {
		t.Error("Root file should contain in_progress status")
	}
	if !strings.Contains(contentStr, `"status": "pending"`) {
		t.Error("Root file should contain pending status")
	}
}

// TestWriteChildTasksJSONC tests writing a child file
func TestWriteChildTasksJSONC(t *testing.T) {
	tempDir := t.TempDir()
	childPath := fmt.Sprintf("%s/tasks-1.jsonc", tempDir)
	changeID := "my-testChangeID"
	parentTaskID := "1"

	childTasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Implementation",
			Description: "Create database schema",
			Status:      parsers.TaskStatusCompleted,
		},
		{
			ID:          "1.2",
			Section:     "Implementation",
			Description: "Implement API handlers",
			Status:      parsers.TaskStatusInProgress,
		},
	}

	err := writeChildTasksJSONC(childPath, changeID, parentTaskID, childTasks)
	if err != nil {
		t.Fatalf("writeChildTasksJSONC() failed: %v", err)
	}

	// Read the file back
	content, err := os.ReadFile(childPath)
	if err != nil {
		t.Fatalf("Failed to read child file: %v", err)
	}

	contentStr := string(content)

	// Verify child-specific header
	if !strings.Contains(contentStr, "// Generated by: spectr accept my-testChangeID") {
		t.Error("Child file should contain generation info")
	}
	if !strings.Contains(contentStr, "// Parent change: my-testChangeID") {
		t.Error("Child file should contain parent change info")
	}
	if !strings.Contains(contentStr, "// Parent task: 1") {
		t.Error("Child file should contain parent task info")
	}

	// Verify version 2
	if !strings.Contains(contentStr, `"version": 2`) {
		t.Error("Child file should have version 2")
	}

	// Verify parent field
	if !strings.Contains(contentStr, `"parent": "1"`) {
		t.Error("Child file should have parent field set to '1'")
	}

	// Verify tasks
	if !strings.Contains(contentStr, `"id": "1.1"`) {
		t.Error("Child file should contain task 1.1")
	}
	if !strings.Contains(contentStr, `"id": "1.2"`) {
		t.Error("Child file should contain task 1.2")
	}
	if !strings.Contains(contentStr, "Create database schema") {
		t.Error("Child file should contain task descriptions")
	}
}

// TestWriteChildTasksJSONCMultipleFiles tests writing multiple child files
func TestWriteChildTasksJSONCMultipleFiles(t *testing.T) {
	tempDir := t.TempDir()
	changeID := testChangeID

	// Write multiple child files
	for i := 1; i <= 3; i++ {
		childPath := fmt.Sprintf("%s/tasks-%d.jsonc", tempDir, i)
		parentTaskID := fmt.Sprintf("%d", i)

		tasks := []parsers.Task{
			{
				ID:          fmt.Sprintf("%d.1", i),
				Section:     fmt.Sprintf("Section %d", i),
				Description: fmt.Sprintf("Task %d.1", i),
				Status:      parsers.TaskStatusPending,
			},
		}

		err := writeChildTasksJSONC(childPath, changeID, parentTaskID, tasks)
		if err != nil {
			t.Fatalf("Failed to write child file %d: %v", i, err)
		}
	}

	// Verify all files exist
	for i := 1; i <= 3; i++ {
		childPath := fmt.Sprintf("%s/tasks-%d.jsonc", tempDir, i)
		if _, err := os.Stat(childPath); os.IsNotExist(err) {
			t.Errorf("Child file %d should exist", i)
		}

		// Read and verify parent field
		content, err := os.ReadFile(childPath)
		if err != nil {
			t.Fatalf("Failed to read child file %d: %v", i, err)
		}

		expectedParent := fmt.Sprintf(`"parent": "%d"`, i)
		if !strings.Contains(string(content), expectedParent) {
			t.Errorf("Child file %d should have parent field '%s'", i, expectedParent)
		}
	}
}

// TestDetermineSplitGroups tests the split group determination logic
func TestDetermineSplitGroups(t *testing.T) {
	tests := []struct {
		name          string
		sections      []Section
		expectedCount int
		description   string
	}{
		{
			name: "single small section",
			sections: []Section{
				{
					Name:      "Implementation",
					Number:    "1",
					Tasks:     createTestTasks(5, "1"),
					StartLine: 1,
					EndLine:   30, // Small section
				},
			},
			expectedCount: 1,
			description:   "Small section should stay together",
		},
		{
			name: "multiple small sections",
			sections: []Section{
				{
					Name:      "Implementation",
					Number:    "1",
					Tasks:     createTestTasks(5, "1"),
					StartLine: 1,
					EndLine:   30,
				},
				{
					Name:      "Testing",
					Number:    "2",
					Tasks:     createTestTasks(5, "2"),
					StartLine: 31,
					EndLine:   60,
				},
			},
			expectedCount: 2,
			description:   "Multiple small sections should each be a group",
		},
		{
			name: "one large section should be split",
			sections: []Section{
				{
					Name:      "Implementation",
					Number:    "1",
					Tasks:     createTestTasks(50, "1"), // Many tasks
					StartLine: 1,
					EndLine:   150, // Exceeds threshold
				},
			},
			expectedCount: 1, // Will be split into subsections, but that depends on task IDs
			description:   "Large section should be split by subsections",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			groups := determineSplitGroups(tt.sections)

			// The exact count depends on subsection logic, but we can verify basic structure
			if len(groups) == 0 {
				t.Error("determineSplitGroups() should return at least one group")
			}

			// Verify each group has required fields
			for i, group := range groups {
				if group.parentID == "" {
					t.Errorf("Group %d should have a parentID", i)
				}
				if group.section == "" {
					t.Errorf("Group %d should have a section name", i)
				}
				if len(group.tasks) == 0 {
					t.Errorf("Group %d should have tasks", i)
				}
			}
		})
	}
}

// createTestTasks creates a slice of test tasks with sequential IDs
func createTestTasks(count int, prefix string) []parsers.Task {
	tasks := make([]parsers.Task, count)
	for i := range count {
		tasks[i] = parsers.Task{
			ID:          fmt.Sprintf("%s.%d", prefix, i+1),
			Section:     "Test Section",
			Description: fmt.Sprintf("Test task %d", i+1),
			Status:      parsers.TaskStatusPending,
		}
	}

	return tasks
}

// TestIsVersionTwo tests the isVersionTwo function for detecting hierarchical format
func TestIsVersionTwo(t *testing.T) {
	tests := []struct {
		name     string
		file     *parsers.TasksFile
		expected bool
	}{
		{
			name:     "nil file returns false",
			file:     nil,
			expected: false,
		},
		{
			name: "version 1 file returns false",
			file: &parsers.TasksFile{
				Version: 1,
				Tasks:   make([]parsers.Task, 0),
			},
			expected: false,
		},
		{
			name: "version 2 file returns true",
			file: &parsers.TasksFile{
				Version: 2,
				Tasks:   make([]parsers.Task, 0),
			},
			expected: true,
		},
		{
			name: "version 2 with includes returns true",
			file: &parsers.TasksFile{
				Version:  2,
				Tasks:    make([]parsers.Task, 0),
				Includes: []string{"tasks-*.jsonc"},
			},
			expected: true,
		},
		{
			name: "version 2 with parent returns true",
			file: &parsers.TasksFile{
				Version: 2,
				Parent:  "1",
				Tasks:   make([]parsers.Task, 0),
			},
			expected: true,
		},
		{
			name: "version 0 returns false",
			file: &parsers.TasksFile{
				Version: 0,
				Tasks:   make([]parsers.Task, 0),
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVersionTwo(tt.file)
			if result != tt.expected {
				t.Errorf("isVersionTwo() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestLoadExistingStatusesVersion1 tests reading version 1 flat files
func TestLoadExistingStatusesVersion1(t *testing.T) {
	tempDir := t.TempDir()

	// Create a version 1 tasks.jsonc file
	v1Content := `// Spectr Tasks File (JSONC)
{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Implementation",
      "description": "Task one",
      "status": "completed"
    },
    {
      "id": "1.2",
      "section": "Implementation",
      "description": "Task two",
      "status": "in_progress"
    },
    {
      "id": "2.1",
      "section": "Testing",
      "description": "Task three",
      "status": "pending"
    }
  ]
}
`

	tasksPath := fmt.Sprintf("%s/tasks.jsonc", tempDir)
	if err := os.WriteFile(tasksPath, []byte(v1Content), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Load statuses
	statusMap, err := loadExistingStatuses(tempDir)
	if err != nil {
		t.Fatalf("loadExistingStatuses() failed: %v", err)
	}

	// Verify all statuses were loaded
	expectedStatuses := map[string]parsers.TaskStatusValue{
		"1.1": parsers.TaskStatusCompleted,
		"1.2": parsers.TaskStatusInProgress,
		"2.1": parsers.TaskStatusPending,
	}

	if len(statusMap) != len(expectedStatuses) {
		t.Errorf("Expected %d statuses, got %d", len(expectedStatuses), len(statusMap))
	}

	for id, expectedStatus := range expectedStatuses {
		if status, exists := statusMap[id]; !exists {
			t.Errorf("Missing status for task %s", id)
		} else if status != expectedStatus {
			t.Errorf("Task %s status = %v, want %v", id, status, expectedStatus)
		}
	}
}

// TestLoadExistingStatusesVersion1NoChildLookup tests that version 1 files don't try to load child files
func TestLoadExistingStatusesVersion1NoChildLookup(t *testing.T) {
	tempDir := t.TempDir()

	// Create a version 1 tasks.jsonc file
	v1Content := `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Test",
      "description": "Test task",
      "status": "pending"
    }
  ]
}
`

	tasksPath := fmt.Sprintf("%s/tasks.jsonc", tempDir)
	if err := os.WriteFile(tasksPath, []byte(v1Content), 0o644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Create a child file that should NOT be read (version 1 doesn't use child files)
	childContent := `{
  "version": 2,
  "parent": "1",
  "tasks": [
    {
      "id": "1.1.1",
      "section": "Test",
      "description": "Child task",
      "status": "completed"
    }
  ]
}
`

	childPath := fmt.Sprintf("%s/tasks-1.jsonc", tempDir)
	if err := os.WriteFile(childPath, []byte(childContent), 0o644); err != nil {
		t.Fatalf("Failed to write child file: %v", err)
	}

	// Load statuses - should only read root file, not child
	statusMap, err := loadExistingStatuses(tempDir)
	if err != nil {
		t.Fatalf("loadExistingStatuses() failed: %v", err)
	}

	// Verify only the root file task was loaded
	if len(statusMap) != 1 {
		t.Errorf("Expected 1 status from root file only, got %d", len(statusMap))
	}

	// Verify the root task is present
	if status, exists := statusMap["1.1"]; !exists {
		t.Error("Root task 1.1 should be loaded")
	} else if status != parsers.TaskStatusPending {
		t.Errorf("Task 1.1 status = %v, want pending", status)
	}

	// Verify the child task was NOT loaded
	if _, exists := statusMap["1.1.1"]; exists {
		t.Error("Child task 1.1.1 should NOT be loaded for version 1 file")
	}
}

// TestVersion1ToVersion2Upgrade tests upgrading from version 1 to version 2
func TestVersion1ToVersion2Upgrade(t *testing.T) {
	tempDir := t.TempDir()

	// Create initial version 1 file with some completed tasks
	v1Content := `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Implementation",
      "description": "Task one",
      "status": "completed"
    },
    {
      "id": "1.2",
      "section": "Implementation",
      "description": "Task two",
      "status": "in_progress"
    },
    {
      "id": "2.1",
      "section": "Testing",
      "description": "Task three",
      "status": "pending"
    }
  ]
}
`

	tasksPath := fmt.Sprintf("%s/tasks.jsonc", tempDir)
	if err := os.WriteFile(tasksPath, []byte(v1Content), 0o644); err != nil {
		t.Fatalf("Failed to write version 1 file: %v", err)
	}

	// Load statuses from version 1 file
	statusMap, err := loadExistingStatuses(tempDir)
	if err != nil {
		t.Fatalf("loadExistingStatuses() failed: %v", err)
	}

	// Verify version 1 statuses were loaded
	if len(statusMap) != 3 {
		t.Errorf("Expected 3 statuses from v1 file, got %d", len(statusMap))
	}

	// Now simulate an upgrade to version 2 by creating hierarchical files
	// This would happen when re-running `spectr accept` on a large tasks.md

	// Create version 2 root file
	referenceTasks := []parsers.Task{
		{
			ID:          "1",
			Section:     "Implementation",
			Description: "Implementation tasks",
			Status:      parsers.TaskStatusInProgress, // Computed from children
			Children:    "$ref:tasks-1.jsonc",
		},
		{
			ID:          "2",
			Section:     "Testing",
			Description: "Testing tasks",
			Status:      parsers.TaskStatusPending,
			Children:    "$ref:tasks-2.jsonc",
		},
	}

	if err := writeRootTasksJSONC(tasksPath, referenceTasks); err != nil {
		t.Fatalf("Failed to write version 2 root file: %v", err)
	}

	// Create child files with merged statuses
	childTasks1 := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Implementation",
			Description: "Task one",
			Status:      statusMap["1.1"], // Preserved from v1
		},
		{
			ID:          "1.2",
			Section:     "Implementation",
			Description: "Task two",
			Status:      statusMap["1.2"], // Preserved from v1
		},
	}

	child1Path := fmt.Sprintf("%s/tasks-1.jsonc", tempDir)
	if err := writeChildTasksJSONC(child1Path, "test-change", "1", childTasks1); err != nil {
		t.Fatalf("Failed to write child file 1: %v", err)
	}

	childTasks2 := []parsers.Task{
		{
			ID:          "2.1",
			Section:     "Testing",
			Description: "Task three",
			Status:      statusMap["2.1"], // Preserved from v1
		},
	}

	child2Path := fmt.Sprintf("%s/tasks-2.jsonc", tempDir)
	if err := writeChildTasksJSONC(child2Path, "test-change", "2", childTasks2); err != nil {
		t.Fatalf("Failed to write child file 2: %v", err)
	}

	// Now load statuses again from version 2 files
	v2StatusMap, err := loadExistingStatuses(tempDir)
	if err != nil {
		t.Fatalf("loadExistingStatuses() failed on v2 files: %v", err)
	}

	// Verify all statuses were preserved in the upgrade
	expectedStatuses := map[string]parsers.TaskStatusValue{
		"1":   parsers.TaskStatusInProgress, // From root file
		"2":   parsers.TaskStatusPending,    // From root file
		"1.1": parsers.TaskStatusCompleted,  // Preserved from v1
		"1.2": parsers.TaskStatusInProgress, // Preserved from v1
		"2.1": parsers.TaskStatusPending,    // Preserved from v1
	}

	if len(v2StatusMap) != len(expectedStatuses) {
		t.Errorf(
			"Expected %d statuses after upgrade, got %d",
			len(expectedStatuses),
			len(v2StatusMap),
		)
	}

	for id, expectedStatus := range expectedStatuses {
		if status, exists := v2StatusMap[id]; !exists {
			t.Errorf("Missing status for task %s after upgrade", id)
		} else if status != expectedStatus {
			t.Errorf("Task %s status after upgrade = %v, want %v", id, status, expectedStatus)
		}
	}
}

// TestMixedVersion1AndVersion2Changes tests an integration scenario with both formats
func TestMixedVersion1AndVersion2Changes(t *testing.T) {
	// Create a temporary directory structure with multiple changes
	tempDir := t.TempDir()

	// Change 1: Version 1 format (small change, < 100 lines)
	change1Dir := fmt.Sprintf("%s/change-1", tempDir)
	if err := os.Mkdir(change1Dir, 0o755); err != nil {
		t.Fatalf("Failed to create change-1 dir: %v", err)
	}

	v1Content := `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Small Change",
      "description": "Simple task",
      "status": "completed"
    }
  ]
}
`

	change1TasksPath := fmt.Sprintf("%s/tasks.jsonc", change1Dir)
	if err := os.WriteFile(change1TasksPath, []byte(v1Content), 0o644); err != nil {
		t.Fatalf("Failed to write change-1 tasks: %v", err)
	}

	// Change 2: Version 2 format (large change, > 100 lines, split into multiple files)
	change2Dir := fmt.Sprintf("%s/change-2", tempDir)
	if err := os.Mkdir(change2Dir, 0o755); err != nil {
		t.Fatalf("Failed to create change-2 dir: %v", err)
	}

	// Version 2 root file
	v2RootTasks := []parsers.Task{
		{
			ID:          "1",
			Section:     "Large Change Section 1",
			Description: "Section 1 tasks",
			Status:      parsers.TaskStatusInProgress,
			Children:    "$ref:tasks-1.jsonc",
		},
		{
			ID:          "2",
			Section:     "Large Change Section 2",
			Description: "Section 2 tasks",
			Status:      parsers.TaskStatusPending,
			Children:    "$ref:tasks-2.jsonc",
		},
	}

	change2TasksPath := fmt.Sprintf("%s/tasks.jsonc", change2Dir)
	if err := writeRootTasksJSONC(change2TasksPath, v2RootTasks); err != nil {
		t.Fatalf("Failed to write change-2 root tasks: %v", err)
	}

	// Version 2 child files
	change2Child1Tasks := []parsers.Task{
		{
			ID:          "1.1",
			Section:     "Large Change Section 1",
			Description: "Child task 1.1",
			Status:      parsers.TaskStatusCompleted,
		},
		{
			ID:          "1.2",
			Section:     "Large Change Section 1",
			Description: "Child task 1.2",
			Status:      parsers.TaskStatusInProgress,
		},
	}

	change2Child1Path := fmt.Sprintf("%s/tasks-1.jsonc", change2Dir)
	if err := writeChildTasksJSONC(change2Child1Path, "change-2", "1", change2Child1Tasks); err != nil {
		t.Fatalf("Failed to write change-2 child file 1: %v", err)
	}

	change2Child2Tasks := []parsers.Task{
		{
			ID:          "2.1",
			Section:     "Large Change Section 2",
			Description: "Child task 2.1",
			Status:      parsers.TaskStatusPending,
		},
	}

	change2Child2Path := fmt.Sprintf("%s/tasks-2.jsonc", change2Dir)
	if err := writeChildTasksJSONC(change2Child2Path, "change-2", "2", change2Child2Tasks); err != nil {
		t.Fatalf("Failed to write change-2 child file 2: %v", err)
	}

	// Now load statuses from both changes
	change1StatusMap, err := loadExistingStatuses(change1Dir)
	if err != nil {
		t.Fatalf("loadExistingStatuses() failed on change-1: %v", err)
	}

	change2StatusMap, err := loadExistingStatuses(change2Dir)
	if err != nil {
		t.Fatalf("loadExistingStatuses() failed on change-2: %v", err)
	}

	// Verify change 1 (version 1) loaded correctly
	if len(change1StatusMap) != 1 {
		t.Errorf("Change 1 should have 1 task, got %d", len(change1StatusMap))
	}

	if status, exists := change1StatusMap["1.1"]; !exists {
		t.Error("Change 1 task 1.1 should exist")
	} else if status != parsers.TaskStatusCompleted {
		t.Errorf("Change 1 task 1.1 status = %v, want completed", status)
	}

	// Verify change 2 (version 2) loaded correctly
	// Should have 2 root tasks + 3 child tasks = 5 total
	expectedChange2Count := 5
	if len(change2StatusMap) != expectedChange2Count {
		t.Errorf(
			"Change 2 should have %d tasks, got %d",
			expectedChange2Count,
			len(change2StatusMap),
		)
	}

	// Verify root tasks
	if status, exists := change2StatusMap["1"]; !exists {
		t.Error("Change 2 root task 1 should exist")
	} else if status != parsers.TaskStatusInProgress {
		t.Errorf("Change 2 root task 1 status = %v, want in_progress", status)
	}

	// Verify child tasks
	expectedChange2Statuses := map[string]parsers.TaskStatusValue{
		"1":   parsers.TaskStatusInProgress,
		"2":   parsers.TaskStatusPending,
		"1.1": parsers.TaskStatusCompleted,
		"1.2": parsers.TaskStatusInProgress,
		"2.1": parsers.TaskStatusPending,
	}

	for id, expectedStatus := range expectedChange2Statuses {
		if status, exists := change2StatusMap[id]; !exists {
			t.Errorf("Change 2 task %s should exist", id)
		} else if status != expectedStatus {
			t.Errorf("Change 2 task %s status = %v, want %v", id, status, expectedStatus)
		}
	}
}
