package ralph

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseTaskGraph(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    map[string]string
		wantTaskCount int
		wantRootCount int
		wantChildren  map[string]int // parent ID -> child count
		wantErr       bool
	}{
		{
			name: "single file with simple tasks",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{
							"id": "1.1",
							"section": "Core Infrastructure",
							"description": "Create ralph package directory",
							"status": "pending"
						},
						{
							"id": "1.2",
							"section": "Core Infrastructure",
							"description": "Define Task and TaskGraph types",
							"status": "pending"
						},
						{
							"id": "2.1",
							"section": "Ralpher Interface",
							"description": "Define Ralpher interface",
							"status": "pending"
						}
					]
				}`,
			},
			wantTaskCount: 3,
			wantRootCount: 3,
			wantChildren: map[string]int{
				"1": 2, // 1.1, 1.2
				"2": 1, // 2.1
			},
		},
		{
			name: "multiple files with tasks",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{
							"id": "1.1",
							"section": "Core",
							"description": "Task 1.1",
							"status": "pending"
						}
					]
				}`,
				"tasks-2.jsonc": `{
					"version": 1,
					"tasks": [
						{
							"id": "2.1",
							"section": "Section 2",
							"description": "Task 2.1",
							"status": "pending"
						}
					]
				}`,
			},
			wantTaskCount: 2,
			wantRootCount: 2,
			wantChildren: map[string]int{
				"1": 1,
				"2": 1,
			},
		},
		{
			name: "nested task hierarchy",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{
							"id": "1.1",
							"section": "Core",
							"description": "Task 1.1",
							"status": "pending"
						},
						{
							"id": "1.1.1",
							"section": "Core",
							"description": "Task 1.1.1",
							"status": "pending"
						},
						{
							"id": "1.1.2",
							"section": "Core",
							"description": "Task 1.1.2",
							"status": "pending"
						},
						{
							"id": "1.2",
							"section": "Core",
							"description": "Task 1.2",
							"status": "pending"
						}
					]
				}`,
			},
			wantTaskCount: 4,
			wantRootCount: 2,
			wantChildren: map[string]int{
				"1":   2, // 1.1, 1.2
				"1.1": 2, // 1.1.1, 1.1.2
			},
		},
		{
			name: "with JSONC comments",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					// This is a comment
					"version": 1,
					/* Multi-line
					   comment */
					"tasks": [
						{
							"id": "1.1",
							"section": "Core Infrastructure",
							"description": "Create package", // inline comment
							"status": "pending"
						}
					]
				}`,
			},
			wantTaskCount: 1,
			wantRootCount: 1,
			wantChildren: map[string]int{
				"1": 1,
			},
		},
		{
			name:          "no tasks files",
			setupFiles:    make(map[string]string),
			wantTaskCount: 0,
			wantRootCount: 0,
			wantChildren:  nil,
			wantErr:       true,
		},
		{
			name: "invalid JSON",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{
							"id": "1.1",
							"section": "Core",
							invalid syntax here
						}
					]
				}`,
			},
			wantTaskCount: 0,
			wantRootCount: 0,
			wantChildren:  nil,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir := t.TempDir()

			// Setup test files
			for filename, content := range tt.setupFiles {
				filePath := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}

			// Parse the task graph
			graph, err := ParseTaskGraph(tmpDir)

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseTaskGraph() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			// Verify task count
			if len(graph.Tasks) != tt.wantTaskCount {
				t.Errorf("got %d tasks, want %d", len(graph.Tasks), tt.wantTaskCount)
			}

			// Verify root count
			if len(graph.Roots) != tt.wantRootCount {
				t.Errorf("got %d roots, want %d", len(graph.Roots), tt.wantRootCount)
			}

			// Verify children counts
			for parentID, wantCount := range tt.wantChildren {
				gotCount := len(graph.Children[parentID])
				if gotCount != wantCount {
					t.Errorf("parent %s: got %d children, want %d", parentID, gotCount, wantCount)
				}
			}
		})
	}
}

func TestGetParentID(t *testing.T) {
	tests := []struct {
		name   string
		taskID string
		want   string
	}{
		{
			name:   "root task",
			taskID: "1",
			want:   "",
		},
		{
			name:   "first level child",
			taskID: "1.1",
			want:   "1",
		},
		{
			name:   "second level child",
			taskID: "1.1.1",
			want:   "1.1",
		},
		{
			name:   "third level child",
			taskID: "1.2.3.4",
			want:   "1.2.3",
		},
		{
			name:   "different root",
			taskID: "2.1",
			want:   "2",
		},
		{
			name:   "another root",
			taskID: "10",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getParentID(tt.taskID)
			if got != tt.want {
				t.Errorf("getParentID(%s) = %s, want %s", tt.taskID, got, tt.want)
			}
		})
	}
}

func TestTaskGraphRelationships(t *testing.T) {
	// Create a more complex task graph for relationship testing
	tmpDir := t.TempDir()

	content := `{
		"version": 1,
		"tasks": [
			{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"},
			{"id": "1.2", "section": "Core", "description": "Task 1.2", "status": "pending"},
			{"id": "1.2.1", "section": "Core", "description": "Task 1.2.1", "status": "pending"},
			{"id": "1.2.2", "section": "Core", "description": "Task 1.2.2", "status": "pending"},
			{"id": "2.1", "section": "Other", "description": "Task 2.1", "status": "pending"},
			{"id": "2.1.1", "section": "Other", "description": "Task 2.1.1", "status": "pending"}
		]
	}`

	filePath := filepath.Join(tmpDir, "tasks.jsonc")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	graph, err := ParseTaskGraph(tmpDir)
	if err != nil {
		t.Fatalf("ParseTaskGraph() error = %v", err)
	}

	// Test specific relationships
	t.Run("verify parent-child relationships", func(t *testing.T) {
		tests := []struct {
			parent   string
			expected []string
		}{
			{"1", []string{"1.1", "1.2"}},
			{"1.2", []string{"1.2.1", "1.2.2"}},
			{"2", []string{"2.1"}},
			{"2.1", []string{"2.1.1"}},
		}

		for _, tt := range tests {
			children := graph.Children[tt.parent]
			if len(children) != len(tt.expected) {
				t.Errorf(
					"parent %s: got %d children, want %d",
					tt.parent,
					len(children),
					len(tt.expected),
				)

				continue
			}

			// Check that all expected children are present
			childMap := make(map[string]bool)
			for _, child := range children {
				childMap[child] = true
			}

			for _, expected := range tt.expected {
				if !childMap[expected] {
					t.Errorf("parent %s: missing expected child %s", tt.parent, expected)
				}
			}
		}
	})

	t.Run("verify roots", func(t *testing.T) {
		expectedRoots := map[string]bool{
			"1.1": true,
			"1.2": true,
			"2.1": true,
		}

		if len(graph.Roots) != len(expectedRoots) {
			t.Errorf("got %d roots, want %d", len(graph.Roots), len(expectedRoots))
		}

		for _, root := range graph.Roots {
			if !expectedRoots[root] {
				t.Errorf("unexpected root: %s", root)
			}
		}
	})

	t.Run("verify all tasks exist", func(t *testing.T) {
		expectedTasks := []string{"1.1", "1.2", "1.2.1", "1.2.2", "2.1", "2.1.1"}
		for _, taskID := range expectedTasks {
			task, exists := graph.Tasks[taskID]
			if !exists {
				t.Errorf("task %s not found in graph", taskID)

				continue
			}
			if task.ID != taskID {
				t.Errorf("task ID mismatch: got %s, want %s", task.ID, taskID)
			}
		}
	})
}

func TestParseTaskGraphPreservesTaskData(t *testing.T) {
	// Test that all task fields are correctly preserved
	tmpDir := t.TempDir()

	content := `{
		"version": 1,
		"tasks": [
			{
				"id": "1.1",
				"section": "Test Section",
				"description": "Test description with special chars: \"quotes\" and \\backslashes\\",
				"status": "in_progress",
				"children": "$ref:tasks-2.jsonc"
			}
		]
	}`

	filePath := filepath.Join(tmpDir, "tasks.jsonc")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	graph, err := ParseTaskGraph(tmpDir)
	if err != nil {
		t.Fatalf("ParseTaskGraph() error = %v", err)
	}

	task := graph.Tasks["1.1"]
	if task == nil {
		t.Fatal("task 1.1 not found")
	}

	// Verify all fields
	const expectedID = "1.1"
	if task.ID != expectedID {
		t.Errorf("ID = %s, want %s", task.ID, expectedID)
	}
	if task.Section != "Test Section" {
		t.Errorf("Section = %s, want Test Section", task.Section)
	}
	expectedDesc := `Test description with special chars: "quotes" and \backslashes\`
	if task.Description != expectedDesc {
		t.Errorf("Description = %s, want %s", task.Description, expectedDesc)
	}
	if task.Status != "in_progress" {
		t.Errorf("Status = %s, want in_progress", task.Status)
	}
	if task.Children != "$ref:tasks-2.jsonc" {
		t.Errorf("Children = %s, want $ref:tasks-2.jsonc", task.Children)
	}
}

//nolint:revive // Test function complexity is acceptable for comprehensive test coverage
func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    map[string]string
		wantStages    [][]string // expected stages (order matters within stages is flexible)
		wantErr       bool
		validateStage func(t *testing.T, stages [][]string) // custom validation
	}{
		{
			name: "simple linear dependencies",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"},
						{"id": "1.2", "section": "Core", "description": "Task 1.2", "status": "pending"},
						{"id": "1.3", "section": "Core", "description": "Task 1.3", "status": "pending"}
					]
				}`,
			},
			validateStage: func(t *testing.T, stages [][]string) {
				if len(stages) != 3 {
					t.Errorf("expected 3 stages, got %d", len(stages))

					return
				}
				// Each stage should have exactly 1 task (linear)
				for i, stage := range stages {
					if len(stage) != 1 {
						t.Errorf("stage %d: expected 1 task, got %d", i, len(stage))
					}
				}
				// Verify execution order: 1.1 -> 1.2 -> 1.3
				if stages[0][0] != testTaskIDOne || stages[1][0] != "1.2" || stages[2][0] != "1.3" {
					t.Errorf("incorrect order: got %v", stages)
				}
			},
		},
		{
			name: "parallel execution - independent roots",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"},
						{"id": "1.2", "section": "Core", "description": "Task 1.2", "status": "pending"},
						{"id": "2.1", "section": "Other", "description": "Task 2.1", "status": "pending"},
						{"id": "2.2", "section": "Other", "description": "Task 2.2", "status": "pending"}
					]
				}`,
			},
			validateStage: func(t *testing.T, stages [][]string) {
				if len(stages) != 2 {
					t.Errorf("expected 2 stages, got %d", len(stages))

					return
				}
				// Stage 0: 1.1 and 2.1 can run in parallel
				if len(stages[0]) != 2 {
					t.Errorf("stage 0: expected 2 tasks, got %d", len(stages[0]))

					return
				}
				stage0Map := make(map[string]bool)
				for _, id := range stages[0] {
					stage0Map[id] = true
				}
				if !stage0Map["1.1"] || !stage0Map["2.1"] {
					t.Errorf("stage 0: expected [1.1, 2.1], got %v", stages[0])
				}

				// Stage 1: 1.2 and 2.2 can run in parallel
				if len(stages[1]) != 2 {
					t.Errorf("stage 1: expected 2 tasks, got %d", len(stages[1]))

					return
				}
				stage1Map := make(map[string]bool)
				for _, id := range stages[1] {
					stage1Map[id] = true
				}
				if !stage1Map["1.2"] || !stage1Map["2.2"] {
					t.Errorf("stage 1: expected [1.2, 2.2], got %v", stages[1])
				}
			},
		},
		{
			name: "multiple independent trees",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{"id": "1.1", "section": "Tree1", "description": "Task 1.1", "status": "pending"},
						{"id": "1.2", "section": "Tree1", "description": "Task 1.2", "status": "pending"},
						{"id": "2.1", "section": "Tree2", "description": "Task 2.1", "status": "pending"},
						{"id": "2.2", "section": "Tree2", "description": "Task 2.2", "status": "pending"},
						{"id": "3.1", "section": "Tree3", "description": "Task 3.1", "status": "pending"}
					]
				}`,
			},
			validateStage: func(t *testing.T, stages [][]string) {
				if len(stages) != 2 {
					t.Errorf("expected 2 stages, got %d", len(stages))

					return
				}
				// Stage 0: all .1 tasks can run in parallel
				if len(stages[0]) != 3 {
					t.Errorf("stage 0: expected 3 tasks, got %d", len(stages[0]))

					return
				}
				stage0Map := make(map[string]bool)
				for _, id := range stages[0] {
					stage0Map[id] = true
				}
				if !stage0Map["1.1"] || !stage0Map["2.1"] || !stage0Map["3.1"] {
					t.Errorf("stage 0: expected [1.1, 2.1, 3.1], got %v", stages[0])
				}

				// Stage 1: 1.2 and 2.2 can run in parallel
				if len(stages[1]) != 2 {
					t.Errorf("stage 1: expected 2 tasks, got %d", len(stages[1]))

					return
				}
				stage1Map := make(map[string]bool)
				for _, id := range stages[1] {
					stage1Map[id] = true
				}
				if !stage1Map["1.2"] || !stage1Map["2.2"] {
					t.Errorf("stage 1: expected [1.2, 2.2], got %v", stages[1])
				}
			},
		},
		{
			name: "nested hierarchy with parent dependencies",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{"id": "1", "section": "Root", "description": "Task 1", "status": "pending"},
						{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"},
						{"id": "1.2", "section": "Core", "description": "Task 1.2", "status": "pending"},
						{"id": "1.2.1", "section": "Core", "description": "Task 1.2.1", "status": "pending"},
						{"id": "1.2.2", "section": "Core", "description": "Task 1.2.2", "status": "pending"}
					]
				}`,
			},
			validateStage: func(t *testing.T, stages [][]string) {
				// Stage 0: task 1 (root)
				// Stage 1: task 1.1
				// Stage 2: task 1.2
				// Stage 3: task 1.2.1
				// Stage 4: task 1.2.2
				if len(stages) != 5 {
					t.Errorf("expected 5 stages, got %d: %v", len(stages), stages)

					return
				}

				expectedOrder := []string{"1", "1.1", "1.2", "1.2.1", "1.2.2"}
				for i, expected := range expectedOrder {
					if len(stages[i]) != 1 {
						t.Errorf("stage %d: expected 1 task, got %d", i, len(stages[i]))

						continue
					}
					if stages[i][0] != expected {
						t.Errorf("stage %d: expected %s, got %s", i, expected, stages[i][0])
					}
				}
			},
		},
		{
			name: "complex mixed hierarchy",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"},
						{"id": "1.1.1", "section": "Core", "description": "Task 1.1.1", "status": "pending"},
						{"id": "1.1.2", "section": "Core", "description": "Task 1.1.2", "status": "pending"},
						{"id": "1.2", "section": "Core", "description": "Task 1.2", "status": "pending"},
						{"id": "2.1", "section": "Other", "description": "Task 2.1", "status": "pending"},
						{"id": "2.1.1", "section": "Other", "description": "Task 2.1.1", "status": "pending"}
					]
				}`,
			},
			validateStage: func(t *testing.T, stages [][]string) {
				// Stage 0: 1.1 and 2.1 (independent roots)
				// Stage 1: 1.1.1 (child of 1.1), 1.2 (sibling of 1.1), and 2.1.1 (child of 2.1) can run in parallel
				// Stage 2: 1.1.2 (sequential sibling of 1.1.1)
				if len(stages) != 3 {
					t.Errorf("expected 3 stages, got %d: %v", len(stages), stages)

					return
				}

				// Stage 0: should have 1.1 and 2.1
				if len(stages[0]) != 2 {
					t.Errorf("stage 0: expected 2 tasks, got %d", len(stages[0]))

					return
				}
				stage0Map := make(map[string]bool)
				for _, id := range stages[0] {
					stage0Map[id] = true
				}
				if !stage0Map["1.1"] || !stage0Map["2.1"] {
					t.Errorf("stage 0: expected [1.1, 2.1], got %v", stages[0])
				}

				// Stage 1: should have 1.1.1, 1.2, and 2.1.1
				// 1.2 depends on 1.1 (previous sibling), which completed in stage 0
				// 1.1.1 depends on 1.1 (parent), which completed in stage 0
				// 2.1.1 depends on 2.1 (parent), which completed in stage 0
				if len(stages[1]) != 3 {
					t.Errorf("stage 1: expected 3 tasks, got %d: %v", len(stages[1]), stages[1])

					return
				}
				stage1Map := make(map[string]bool)
				for _, id := range stages[1] {
					stage1Map[id] = true
				}
				if !stage1Map["1.1.1"] || !stage1Map["1.2"] || !stage1Map["2.1.1"] {
					t.Errorf("stage 1: expected [1.1.1, 1.2, 2.1.1], got %v", stages[1])
				}

				// Stage 2: should have 1.1.2
				if len(stages[2]) != 1 || stages[2][0] != "1.1.2" {
					t.Errorf("stage 2: expected [1.1.2], got %v", stages[2])
				}
			},
		},
		{
			name: "single task",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"}
					]
				}`,
			},
			validateStage: func(t *testing.T, stages [][]string) {
				if len(stages) != 1 {
					t.Errorf("expected 1 stage, got %d", len(stages))

					return
				}
				if len(stages[0]) != 1 || stages[0][0] != testTaskIDOne {
					t.Errorf("expected [[1.1]], got %v", stages)
				}
			},
		},
		{
			name: "gaps in numbering",
			setupFiles: map[string]string{
				"tasks.jsonc": `{
					"version": 1,
					"tasks": [
						{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"},
						{"id": "1.3", "section": "Core", "description": "Task 1.3", "status": "pending"},
						{"id": "1.5", "section": "Core", "description": "Task 1.5", "status": "pending"}
					]
				}`,
			},
			validateStage: func(t *testing.T, stages [][]string) {
				// With gaps, tasks should run in parallel since there's no previous sibling
				if len(stages) != 1 {
					t.Errorf("expected 1 stage (all parallel), got %d: %v", len(stages), stages)

					return
				}
				if len(stages[0]) != 3 {
					t.Errorf("stage 0: expected 3 tasks, got %d", len(stages[0]))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory for test
			tmpDir := t.TempDir()

			// Setup test files
			for filename, content := range tt.setupFiles {
				filePath := filepath.Join(tmpDir, filename)
				if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write test file: %v", err)
				}
			}

			// Parse the task graph
			graph, err := ParseTaskGraph(tmpDir)
			if err != nil {
				t.Fatalf("ParseTaskGraph() error = %v", err)
			}

			// Perform topological sort
			stages, err := graph.TopologicalSort()

			// Check error expectation
			if (err != nil) != tt.wantErr {
				t.Errorf("TopologicalSort() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if tt.wantErr {
				return
			}

			// Use custom validation if provided
			if tt.validateStage != nil {
				tt.validateStage(t, stages)
			}
		})
	}
}

func TestGetRootPrefix(t *testing.T) {
	graph := &TaskGraph{}

	tests := []struct {
		name   string
		taskID string
		want   string
	}{
		{
			name:   "single component root",
			taskID: "1",
			want:   "1",
		},
		{
			name:   "two component task",
			taskID: "1.2",
			want:   "1",
		},
		{
			name:   "three component task",
			taskID: "1.2.3",
			want:   "1",
		},
		{
			name:   "deeply nested task",
			taskID: "1.2.3.4.5",
			want:   "1",
		},
		{
			name:   "different root - task 2",
			taskID: "2.1",
			want:   "2",
		},
		{
			name:   "different root - task 10",
			taskID: "10.5.3",
			want:   "10",
		},
		{
			name:   "empty string",
			taskID: "",
			want:   "",
		},
		{
			name:   "root only - single digit",
			taskID: "5",
			want:   "5",
		},
		{
			name:   "root only - double digit",
			taskID: "99",
			want:   "99",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graph.GetRootPrefix(tt.taskID)
			if got != tt.want {
				t.Errorf("GetRootPrefix(%q) = %q, want %q", tt.taskID, got, tt.want)
			}
		})
	}
}

func TestCanRunInParallel(t *testing.T) {
	graph := &TaskGraph{}

	tests := []struct {
		name    string
		taskID1 string
		taskID2 string
		want    bool
		reason  string
	}{
		// Different roots - can run in parallel
		{
			name:    "different roots - simple",
			taskID1: "1.1",
			taskID2: "2.1",
			want:    true,
			reason:  "different root prefixes (1 vs 2)",
		},
		{
			name:    "different roots - same depth",
			taskID1: "1.2.3",
			taskID2: "2.1.1",
			want:    true,
			reason:  "different root prefixes (1 vs 2)",
		},
		{
			name:    "different roots - different depths",
			taskID1: "1.1",
			taskID2: "2.1.1.1",
			want:    true,
			reason:  "different root prefixes (1 vs 2)",
		},
		{
			name:    "different roots - root level",
			taskID1: "1",
			taskID2: "2",
			want:    true,
			reason:  "different root prefixes",
		},
		{
			name:    "different roots - high numbers",
			taskID1: "10.1",
			taskID2: "20.1",
			want:    true,
			reason:  "different root prefixes (10 vs 20)",
		},

		// Same root - cannot run in parallel
		{
			name:    "same root - sequential siblings",
			taskID1: "1.1",
			taskID2: "1.2",
			want:    false,
			reason:  "same root, sequential siblings",
		},
		{
			name:    "same root - non-adjacent siblings",
			taskID1: "1.1",
			taskID2: "1.5",
			want:    false,
			reason:  "same root, siblings in same tree",
		},
		{
			name:    "same root - nested siblings",
			taskID1: "1.2.1",
			taskID2: "1.2.3",
			want:    false,
			reason:  "same root, siblings in same subtree",
		},

		// Parent-child relationships - cannot run in parallel
		{
			name:    "parent-child - direct",
			taskID1: "1",
			taskID2: "1.1",
			want:    false,
			reason:  "parent-child relationship",
		},
		{
			name:    "child-parent - reversed order",
			taskID1: "1.1",
			taskID2: "1",
			want:    false,
			reason:  "child-parent relationship",
		},
		{
			name:    "parent-grandchild",
			taskID1: "1",
			taskID2: "1.1.1",
			want:    false,
			reason:  "parent-grandchild relationship",
		},
		{
			name:    "parent-child - nested",
			taskID1: "1.2",
			taskID2: "1.2.1",
			want:    false,
			reason:  "parent-child relationship in subtree",
		},
		{
			name:    "parent-deeply-nested-child",
			taskID1: "1.2",
			taskID2: "1.2.3.4.5",
			want:    false,
			reason:  "parent-descendant relationship",
		},

		// Edge cases
		{
			name:    "same task",
			taskID1: "1.1",
			taskID2: "1.1",
			want:    false,
			reason:  "same task cannot run in parallel with itself",
		},
		{
			name:    "empty first task",
			taskID1: "",
			taskID2: "1.1",
			want:    false,
			reason:  "empty task ID",
		},
		{
			name:    "empty second task",
			taskID1: "1.1",
			taskID2: "",
			want:    false,
			reason:  "empty task ID",
		},
		{
			name:    "both empty",
			taskID1: "",
			taskID2: "",
			want:    false,
			reason:  "both task IDs empty",
		},

		// Similar looking IDs but different roots
		{
			name:    "looks similar but different roots",
			taskID1: "1.10",
			taskID2: "11.0",
			want:    true,
			reason:  "different root prefixes (1 vs 11)",
		},
		{
			name:    "single vs multi-digit roots",
			taskID1: "1.1",
			taskID2: "10.1",
			want:    true,
			reason:  "different root prefixes (1 vs 10)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := graph.CanRunInParallel(tt.taskID1, tt.taskID2)
			if got != tt.want {
				t.Errorf("CanRunInParallel(%q, %q) = %v, want %v\nReason: %s",
					tt.taskID1, tt.taskID2, got, tt.want, tt.reason)
			}
		})
	}
}

func TestCanRunInParallel_Symmetry(t *testing.T) {
	// Test that CanRunInParallel is symmetric: f(a,b) == f(b,a)
	graph := &TaskGraph{}

	testPairs := []struct {
		taskID1 string
		taskID2 string
	}{
		{"1.1", "2.1"},
		{"1", "1.1"},
		{"1.1", "1.2"},
		{"1.2.1", "1.2.2"},
		{"1.1.1", "2.3.4"},
	}

	for _, pair := range testPairs {
		t.Run(pair.taskID1+"_vs_"+pair.taskID2, func(t *testing.T) {
			result1 := graph.CanRunInParallel(pair.taskID1, pair.taskID2)
			result2 := graph.CanRunInParallel(pair.taskID2, pair.taskID1)

			if result1 != result2 {
				t.Errorf("CanRunInParallel is not symmetric: f(%q, %q) = %v, but f(%q, %q) = %v",
					pair.taskID1, pair.taskID2, result1,
					pair.taskID2, pair.taskID1, result2)
			}
		})
	}
}

func TestCanRunInParallel_WithTopologicalSort(t *testing.T) {
	// Test that tasks identified as parallel are actually in the same stage
	// of the topological sort
	tmpDir := t.TempDir()

	content := `{
		"version": 1,
		"tasks": [
			{"id": "1.1", "section": "Core", "description": "Task 1.1", "status": "pending"},
			{"id": "1.2", "section": "Core", "description": "Task 1.2", "status": "pending"},
			{"id": "2.1", "section": "Other", "description": "Task 2.1", "status": "pending"},
			{"id": "2.2", "section": "Other", "description": "Task 2.2", "status": "pending"},
			{"id": "3.1", "section": "Third", "description": "Task 3.1", "status": "pending"}
		]
	}`

	filePath := filepath.Join(tmpDir, "tasks.jsonc")
	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	graph, err := ParseTaskGraph(tmpDir)
	if err != nil {
		t.Fatalf("ParseTaskGraph() error = %v", err)
	}

	stages, err := graph.TopologicalSort()
	if err != nil {
		t.Fatalf("TopologicalSort() error = %v", err)
	}

	// Check that tasks in the same stage can run in parallel
	for stageIdx, stage := range stages {
		// Check all pairs within a stage
		for i := range stage {
			for j := i + 1; j < len(stage); j++ {
				taskID1 := stage[i]
				taskID2 := stage[j]

				canParallel := graph.CanRunInParallel(taskID1, taskID2)
				if !canParallel {
					t.Errorf(
						"Stage %d: tasks %q and %q are in the same stage but CanRunInParallel returned false",
						stageIdx,
						taskID1,
						taskID2,
					)
				}
			}
		}
	}

	// Verify specific expected parallel relationships
	expectedParallel := [][2]string{
		{"1.1", "2.1"}, // Different roots, stage 0
		{"1.1", "3.1"}, // Different roots, stage 0
		{"2.1", "3.1"}, // Different roots, stage 0
		{"1.2", "2.2"}, // Different roots, stage 1
	}

	for _, pair := range expectedParallel {
		if !graph.CanRunInParallel(pair[0], pair[1]) {
			t.Errorf("Expected %q and %q to be parallel, but got false", pair[0], pair[1])
		}
	}

	// Verify specific expected non-parallel relationships
	expectedNonParallel := [][2]string{
		{"1.1", "1.2"}, // Same root, sequential
		{"2.1", "2.2"}, // Same root, sequential
	}

	for _, pair := range expectedNonParallel {
		if graph.CanRunInParallel(pair[0], pair[1]) {
			t.Errorf("Expected %q and %q to NOT be parallel, but got true", pair[0], pair[1])
		}
	}
}

func TestGetDependencies(t *testing.T) {
	tests := []struct {
		name      string
		tasks     map[string]*Task
		taskID    string
		wantDeps  []string // can be in any order
		wantCount int      // expected number of dependencies
	}{
		{
			name: "first child with parent",
			tasks: map[string]*Task{
				"1":   {ID: "1"},
				"1.1": {ID: "1.1"},
			},
			taskID:    "1.1",
			wantDeps:  []string{"1"},
			wantCount: 1,
		},
		{
			name: "second child depends on parent and previous sibling",
			tasks: map[string]*Task{
				"1":   {ID: "1"},
				"1.1": {ID: "1.1"},
				"1.2": {ID: "1.2"},
			},
			taskID:    "1.2",
			wantDeps:  []string{"1", "1.1"},
			wantCount: 2,
		},
		{
			name: "nested child with parent",
			tasks: map[string]*Task{
				"1.1":   {ID: "1.1"},
				"1.1.1": {ID: "1.1.1"},
			},
			taskID:    "1.1.1",
			wantDeps:  []string{"1.1"},
			wantCount: 1,
		},
		{
			name: "child without parent in graph (orphan root)",
			tasks: map[string]*Task{
				"1.1": {ID: "1.1"},
			},
			taskID:    "1.1",
			wantDeps:  nil,
			wantCount: 0,
		},
		{
			name: "root task has no dependencies",
			tasks: map[string]*Task{
				"1": {ID: "1"},
			},
			taskID:    "1",
			wantDeps:  nil,
			wantCount: 0,
		},
		{
			name: "third sibling depends on parent and previous sibling",
			tasks: map[string]*Task{
				"1":   {ID: "1"},
				"1.1": {ID: "1.1"},
				"1.2": {ID: "1.2"},
				"1.3": {ID: "1.3"},
			},
			taskID:    "1.3",
			wantDeps:  []string{"1", "1.2"},
			wantCount: 2,
		},
		{
			name: "gap in sibling numbering (no previous sibling)",
			tasks: map[string]*Task{
				"1":   {ID: "1"},
				"1.1": {ID: "1.1"},
				"1.3": {ID: "1.3"}, // 1.2 missing
			},
			taskID:    "1.3",
			wantDeps:  []string{"1"}, // only parent, no previous sibling
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			graph := &TaskGraph{
				Tasks:    tt.tasks,
				Children: make(map[string][]string),
			}

			deps := graph.getDependencies(tt.taskID)

			// Check count
			if len(deps) != tt.wantCount {
				t.Errorf("got %d dependencies, want %d: %v", len(deps), tt.wantCount, deps)

				return
			}

			// Check all expected dependencies are present
			depsMap := make(map[string]bool)
			for _, dep := range deps {
				depsMap[dep] = true
			}

			for _, expected := range tt.wantDeps {
				if !depsMap[expected] {
					t.Errorf("missing expected dependency: %s (got %v)", expected, deps)
				}
			}
		})
	}
}
