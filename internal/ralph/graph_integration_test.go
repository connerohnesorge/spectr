package ralph

import (
	"path/filepath"
	"testing"
)

// TestParseTaskGraphWithRealTasksFile tests parsing the actual tasks.jsonc
// from the add-ralph-orchestration change to ensure it works with real data.
func TestParseTaskGraphWithRealTasksFile(t *testing.T) {
	// Get the path to the actual tasks.jsonc file
	changeDir := filepath.Join("..", "..", "spectr", "changes", "add-ralph-orchestration")

	graph, err := ParseTaskGraph(changeDir)
	if err != nil {
		t.Fatalf("Failed to parse real tasks.jsonc: %v", err)
	}

	// Verify we got some tasks
	if len(graph.Tasks) == 0 {
		t.Fatal("Expected tasks to be parsed, got 0")
	}

	// Verify we have some roots
	if len(graph.Roots) == 0 {
		t.Fatal("Expected root tasks, got 0")
	}

	// Verify specific task exists (task 1.3 which we're currently working on)
	task13, exists := graph.Tasks["1.3"]
	if !exists {
		t.Fatal("Task 1.3 should exist in the graph")
	}

	// Verify task 1.3 has expected properties
	if task13.Section != "Core Infrastructure" {
		t.Errorf("Task 1.3 section = %s, want Core Infrastructure", task13.Section)
	}

	if task13.Status != "completed" {
		t.Errorf("Task 1.3 status = %s, want completed", task13.Status)
	}

	// Verify parent-child relationships exist
	// Task 1 should have children (1.1, 1.2, 1.3, 1.4, 1.5)
	children1 := graph.Children["1"]
	if len(children1) < 5 {
		t.Errorf("Task 1 should have at least 5 children, got %d", len(children1))
	}

	// Verify task 1.3 is a child of virtual parent "1"
	found := false
	for _, child := range children1 {
		if child == "1.3" {
			found = true

			break
		}
	}
	if !found {
		t.Error("Task 1.3 should be a child of task 1")
	}

	t.Logf("Successfully parsed %d tasks with %d roots", len(graph.Tasks), len(graph.Roots))
}
