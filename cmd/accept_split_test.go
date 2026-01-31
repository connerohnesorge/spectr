package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// TestExtractSectionNumber verifies section number extraction from task IDs
func TestExtractSectionNumber(t *testing.T) {
	tests := []struct {
		name     string
		taskID   string
		expected string
	}{
		{
			name:     "standard subsection format",
			taskID:   "1.1",
			expected: "1",
		},
		{
			name:     "deeper nesting",
			taskID:   "2.3",
			expected: "2",
		},
		{
			name:     "single number without subsection",
			taskID:   "5",
			expected: "0",
		},
		{
			name:     "simple task without section",
			taskID:   "7",
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractSectionNumber(tt.taskID)
			if got != tt.expected {
				t.Errorf("extractSectionNumber(%s) = %s, want %s",
					tt.taskID, got, tt.expected)
			}
		})
	}
}

// TestGroupTasksBySection verifies task grouping by section
func TestGroupTasksBySection(t *testing.T) {
	tests := []struct {
		name     string
		tasks    []parsers.Task
		expected []sectionGroup
	}{
		{
			name: "single section",
			tasks: []parsers.Task{
				{ID: "1.1", Section: "Setup", Description: "Task 1"},
				{ID: "1.2", Section: "Setup", Description: "Task 2"},
			},
			expected: []sectionGroup{
				{
					sectionNum:  "1",
					sectionName: "Setup",
					tasks: []parsers.Task{
						{ID: "1.1", Section: "Setup", Description: "Task 1"},
						{ID: "1.2", Section: "Setup", Description: "Task 2"},
					},
				},
			},
		},
		{
			name: "multiple sections",
			tasks: []parsers.Task{
				{ID: "1.1", Section: "Setup", Description: "Task 1"},
				{ID: "2.1", Section: "Impl", Description: "Task 2"},
				{ID: "3.1", Section: "Test", Description: "Task 3"},
			},
			expected: []sectionGroup{
				{
					sectionNum:  "1",
					sectionName: "Setup",
					tasks: []parsers.Task{
						{ID: "1.1", Section: "Setup", Description: "Task 1"},
					},
				},
				{
					sectionNum:  "2",
					sectionName: "Impl",
					tasks: []parsers.Task{
						{ID: "2.1", Section: "Impl", Description: "Task 2"},
					},
				},
				{
					sectionNum:  "3",
					sectionName: "Test",
					tasks: []parsers.Task{
						{ID: "3.1", Section: "Test", Description: "Task 3"},
					},
				},
			},
		},
		{
			name: "tasks without sections",
			tasks: []parsers.Task{
				{ID: "1", Section: "", Description: "Task 1"},
				{ID: "2", Section: "", Description: "Task 2"},
			},
			expected: []sectionGroup{
				{
					sectionNum:  "0",
					sectionName: "",
					tasks: []parsers.Task{
						{ID: "1", Section: "", Description: "Task 1"},
						{ID: "2", Section: "", Description: "Task 2"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := groupTasksBySection(tt.tasks)

			if len(got) != len(tt.expected) {
				t.Fatalf("expected %d groups, got %d", len(tt.expected), len(got))
			}

			for i := range got {
				if got[i].sectionNum != tt.expected[i].sectionNum {
					t.Errorf("group %d: expected sectionNum %s, got %s",
						i, tt.expected[i].sectionNum, got[i].sectionNum)
				}
				if got[i].sectionName != tt.expected[i].sectionName {
					t.Errorf("group %d: expected sectionName %s, got %s",
						i, tt.expected[i].sectionName, got[i].sectionName)
				}
				if len(got[i].tasks) != len(tt.expected[i].tasks) {
					t.Errorf("group %d: expected %d tasks, got %d",
						i, len(tt.expected[i].tasks), len(got[i].tasks))
				}
			}
		})
	}
}

// TestShouldSplitTasksJSONC verifies split logic based on task complexity
func TestShouldSplitTasksJSONC(t *testing.T) {
	tests := []struct {
		name        string
		tasks       []parsers.Task
		shouldSplit bool
	}{
		{
			name: "19 tasks single section - no split",
			tasks: func() []parsers.Task {
				tasks := make([]parsers.Task, 19)
				for i := range 19 {
					tasks[i] = parsers.Task{ID: fmt.Sprintf("1.%d", i+1), Section: "Test"}
				}

				return tasks
			}(),
			shouldSplit: false,
		},
		{
			name: "20 tasks single section - no split",
			tasks: func() []parsers.Task {
				tasks := make([]parsers.Task, 20)
				for i := range 20 {
					tasks[i] = parsers.Task{ID: fmt.Sprintf("1.%d", i+1), Section: "Test"}
				}

				return tasks
			}(),
			shouldSplit: false,
		},
		{
			name: "21 tasks multiple sections - split",
			tasks: []parsers.Task{
				{ID: "1.1", Section: "Setup"},
				{ID: "1.2", Section: "Setup"},
				{ID: "1.3", Section: "Setup"},
				{ID: "1.4", Section: "Setup"},
				{ID: "1.5", Section: "Setup"},
				{ID: "1.6", Section: "Setup"},
				{ID: "1.7", Section: "Setup"},
				{ID: "1.8", Section: "Setup"},
				{ID: "1.9", Section: "Setup"},
				{ID: "1.10", Section: "Setup"},
				{ID: "2.1", Section: "Test"},
				{ID: "2.2", Section: "Test"},
				{ID: "2.3", Section: "Test"},
				{ID: "2.4", Section: "Test"},
				{ID: "2.5", Section: "Test"},
				{ID: "2.6", Section: "Test"},
				{ID: "2.7", Section: "Test"},
				{ID: "2.8", Section: "Test"},
				{ID: "2.9", Section: "Test"},
				{ID: "2.10", Section: "Test"},
				{ID: "2.11", Section: "Test"},
			},
			shouldSplit: true,
		},
		{
			name: "33 tasks 6 sections (amp case) - split",
			tasks: func() []parsers.Task {
				tasks := make([]parsers.Task, 33)
				// 6 sections with varying task counts
				sections := []int{5, 6, 6, 7, 5, 4} // totals 33
				taskIdx := 0
				for secNum := 1; secNum <= 6; secNum++ {
					for taskNum := 1; taskNum <= sections[secNum-1]; taskNum++ {
						tasks[taskIdx] = parsers.Task{
							ID:      fmt.Sprintf("%d.%d", secNum, taskNum),
							Section: fmt.Sprintf("Section %d", secNum),
						}
						taskIdx++
					}
				}

				return tasks
			}(),
			shouldSplit: true,
		},
		{
			name: "10 tasks no sections - no split",
			tasks: func() []parsers.Task {
				tasks := make([]parsers.Task, 10)
				for i := range 10 {
					tasks[i] = parsers.Task{ID: fmt.Sprintf("%d", i+1), Section: ""}
				}

				return tasks
			}(),
			shouldSplit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldSplitTasksJSONC(tt.tasks)
			if got != tt.shouldSplit {
				t.Errorf(
					"shouldSplitTasksJSONC() = %v, want %v (tasks: %d)",
					got,
					tt.shouldSplit,
					len(tt.tasks),
				)
			}
		})
	}
}

// TestAggregateSectionStatus verifies status aggregation logic
func TestAggregateSectionStatus(t *testing.T) {
	tests := []struct {
		name      string
		tasks     []parsers.Task
		statusMap map[string]parsers.TaskStatusValue
		expected  parsers.TaskStatusValue
	}{
		{
			name: "all completed",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
				{ID: "1.2", Status: parsers.TaskStatusCompleted},
			},
			statusMap: make(map[string]parsers.TaskStatusValue),
			expected:  parsers.TaskStatusCompleted,
		},
		{
			name: "any in progress",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusCompleted},
				{ID: "1.2", Status: parsers.TaskStatusInProgress},
				{ID: "1.3", Status: parsers.TaskStatusPending},
			},
			statusMap: make(map[string]parsers.TaskStatusValue),
			expected:  parsers.TaskStatusInProgress,
		},
		{
			name: "otherwise pending",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusPending},
			},
			statusMap: make(map[string]parsers.TaskStatusValue),
			expected:  parsers.TaskStatusPending,
		},
		{
			name: "status map overrides task status",
			tasks: []parsers.Task{
				{ID: "1.1", Status: parsers.TaskStatusPending},
				{ID: "1.2", Status: parsers.TaskStatusPending},
			},
			statusMap: map[string]parsers.TaskStatusValue{
				"1.1": parsers.TaskStatusCompleted,
				"1.2": parsers.TaskStatusCompleted,
			},
			expected: parsers.TaskStatusCompleted,
		},
		{
			name:      "empty tasks",
			tasks:     make([]parsers.Task, 0),
			statusMap: make(map[string]parsers.TaskStatusValue),
			expected:  parsers.TaskStatusPending,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := aggregateSectionStatus(tt.tasks, tt.statusMap)
			if got != tt.expected {
				t.Errorf("aggregateSectionStatus() = %s, want %s", got, tt.expected)
			}
		})
	}
}

// TestBuildTaskStatusMap verifies status preservation from existing files
func TestBuildTaskStatusMap(t *testing.T) {
	tmpDir := t.TempDir()

	// Create v1 tasks.jsonc
	tasksJSON := `{
  "version": 1,
  "tasks": [
    {
      "id": "1.1",
      "section": "Setup",
      "description": "Task 1",
      "status": "completed"
    },
    {
      "id": "1.2",
      "section": "Setup",
      "description": "Task 2",
      "status": "in_progress"
    }
  ]
}`
	tasksPath := filepath.Join(tmpDir, "tasks.jsonc")
	if err := os.WriteFile(tasksPath, []byte(tasksJSON), 0o644); err != nil {
		t.Fatalf("failed to write tasks.jsonc: %v", err)
	}

	// Build status map
	statusMap := buildTaskStatusMap(tmpDir)

	// Verify status extraction
	if statusMap["1.1"] != parsers.TaskStatusCompleted {
		t.Errorf("expected 1.1 to be completed, got %s", statusMap["1.1"])
	}
	if statusMap["1.2"] != parsers.TaskStatusInProgress {
		t.Errorf("expected 1.2 to be in_progress, got %s", statusMap["1.2"])
	}
}

// TestBuildTaskStatusMapV2 verifies status preservation from v2 hierarchical files
func TestBuildTaskStatusMapV2(t *testing.T) {
	tmpDir := t.TempDir()

	// Create v2 root tasks.jsonc
	rootJSON := `{
  "version": 2,
  "tasks": [
    {
      "id": "1",
      "section": "Setup",
      "description": "Setup tasks",
      "status": "in_progress",
      "children": "$ref:tasks-1.jsonc"
    }
  ],
  "includes": ["tasks-*.jsonc"]
}`
	rootPath := filepath.Join(tmpDir, "tasks.jsonc")
	if err := os.WriteFile(rootPath, []byte(rootJSON), 0o644); err != nil {
		t.Fatalf("failed to write root tasks.jsonc: %v", err)
	}

	// Create v2 child tasks-1.jsonc
	childJSON := `{
  "version": 2,
  "parent": "1",
  "tasks": [
    {
      "id": "1.1",
      "section": "Setup",
      "description": "Task 1",
      "status": "completed"
    },
    {
      "id": "1.2",
      "section": "Setup",
      "description": "Task 2",
      "status": "pending"
    }
  ]
}`
	childPath := filepath.Join(tmpDir, "tasks-1.jsonc")
	if err := os.WriteFile(childPath, []byte(childJSON), 0o644); err != nil {
		t.Fatalf("failed to write child tasks.jsonc: %v", err)
	}

	// Build status map
	statusMap := buildTaskStatusMap(tmpDir)

	// Verify status extraction from both files
	if statusMap["1"] != parsers.TaskStatusInProgress {
		t.Errorf("expected task 1 to be in_progress, got %s", statusMap["1"])
	}
	if statusMap["1.1"] != parsers.TaskStatusCompleted {
		t.Errorf("expected 1.1 to be completed, got %s", statusMap["1.1"])
	}
	if statusMap["1.2"] != parsers.TaskStatusPending {
		t.Errorf("expected 1.2 to be pending, got %s", statusMap["1.2"])
	}
}

// TestApplyStatusPreservation verifies status application to tasks
func TestApplyStatusPreservation(t *testing.T) {
	tasks := []parsers.Task{
		{ID: "1.1", Status: parsers.TaskStatusPending},
		{ID: "1.2", Status: parsers.TaskStatusPending},
		{ID: "1.3", Status: parsers.TaskStatusPending},
	}

	statusMap := map[string]parsers.TaskStatusValue{
		"1.1": parsers.TaskStatusCompleted,
		"1.2": parsers.TaskStatusInProgress,
		// 1.3 not in map - should remain pending
	}

	applyStatusPreservation(tasks, statusMap)

	if tasks[0].Status != parsers.TaskStatusCompleted {
		t.Errorf("expected 1.1 to be completed, got %s", tasks[0].Status)
	}
	if tasks[1].Status != parsers.TaskStatusInProgress {
		t.Errorf("expected 1.2 to be in_progress, got %s", tasks[1].Status)
	}
	if tasks[2].Status != parsers.TaskStatusPending {
		t.Errorf("expected 1.3 to remain pending, got %s", tasks[2].Status)
	}
}

// TestWriteHierarchicalTasksJSONC verifies v2 hierarchical file generation
func TestWriteHierarchicalTasksJSONC(t *testing.T) {
	tmpDir := t.TempDir()

	sections := []sectionGroup{
		{
			sectionNum:  "1",
			sectionName: "Setup",
			tasks: []parsers.Task{
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
					Status:      parsers.TaskStatusCompleted,
				},
			},
		},
		{
			sectionNum:  "2",
			sectionName: "Implementation",
			tasks: []parsers.Task{
				{
					ID:          "2.1",
					Section:     "Implementation",
					Description: "Task 3",
					Status:      parsers.TaskStatusPending,
				},
			},
		},
	}

	statusMap := make(map[string]parsers.TaskStatusValue)

	err := writeHierarchicalTasksJSONC(tmpDir, "test-change", sections, statusMap)
	if err != nil {
		t.Fatalf("writeHierarchicalTasksJSONC() error = %v", err)
	}

	// Verify root file exists and has correct structure
	rootPath := filepath.Join(tmpDir, "tasks.jsonc")
	rootData, err := os.ReadFile(rootPath)
	if err != nil {
		t.Fatalf("failed to read root tasks.jsonc: %v", err)
	}

	rootJSON := parsers.StripJSONComments(rootData)
	var rootFile parsers.TasksFile
	if err := json.Unmarshal(rootJSON, &rootFile); err != nil {
		t.Fatalf("failed to parse root tasks.jsonc: %v", err)
	}

	// Verify v2 format
	if rootFile.Version != 2 {
		t.Errorf("expected version 2, got %d", rootFile.Version)
	}

	// Verify includes
	if len(rootFile.Includes) != 1 || rootFile.Includes[0] != "tasks-*.jsonc" {
		t.Errorf("expected includes ['tasks-*.jsonc'], got %v", rootFile.Includes)
	}

	// Verify root tasks have children references
	if len(rootFile.Tasks) != 2 {
		t.Fatalf("expected 2 root tasks, got %d", len(rootFile.Tasks))
	}

	if rootFile.Tasks[0].Children != "$ref:tasks-1.jsonc" {
		t.Errorf("expected children ref 'tasks-1.jsonc', got %s", rootFile.Tasks[0].Children)
	}

	// Verify child file exists
	child1Path := filepath.Join(tmpDir, "tasks-1.jsonc")
	child1Data, err := os.ReadFile(child1Path)
	if err != nil {
		t.Fatalf("failed to read tasks-1.jsonc: %v", err)
	}

	child1JSON := parsers.StripJSONComments(child1Data)
	var child1File parsers.TasksFile
	if err := json.Unmarshal(child1JSON, &child1File); err != nil {
		t.Fatalf("failed to parse tasks-1.jsonc: %v", err)
	}

	// Verify child file structure
	if child1File.Version != 2 {
		t.Errorf("expected child version 2, got %d", child1File.Version)
	}
	if child1File.Parent != "1" {
		t.Errorf("expected parent '1', got %s", child1File.Parent)
	}
	if len(child1File.Tasks) != 2 {
		t.Errorf("expected 2 child tasks, got %d", len(child1File.Tasks))
	}
}

// TestStripJSONCComments verifies comment removal
func TestStripJSONCComments(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "single line comment",
			input:    `{"key": "value"} // comment`,
			expected: `{"key": "value"} `,
		},
		{
			name:     "multi-line comment",
			input:    `{"key": /* comment */ "value"}`,
			expected: `{"key":  "value"}`,
		},
		{
			name:     "no comments",
			input:    `{"key": "value"}`,
			expected: `{"key": "value"}`,
		},
		{
			name:     "comment with slashes in string",
			input:    `{"url": "http://example.com"}`,
			expected: `{"url": "http://example.com"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(stripJSONCComments([]byte(tt.input)))
			if got != tt.expected {
				t.Errorf("stripJSONCComments() = %q, want %q", got, tt.expected)
			}
		})
	}
}
