package cmd

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// TestJSONCValidation_NestedQuotes tests that JSONC validation correctly
// handles nested quotes in task descriptions. This test validates the extreme
// edge case of task 1.1 from test-extreme-jsonc change proposal.
//
// Test case: "He said \"Hello\" and she replied \"Hi there\""
//
// This ensures that:
//  1. The description with nested quotes can be marshaled to valid JSON
//  2. The JSON can be parsed back successfully
//  3. Round-trip conversion preserves the exact string (lossless)
func TestJSONCValidation_NestedQuotes(t *testing.T) {
	// This is the exact description from task 1.1 in test-extreme-jsonc
	description := "Nested quotes: \"He said \\\"Hello\\\" and she replied \\\"Hi there\\\"\""

	// Create a task with nested quotes in the description
	task := parsers.Task{
		ID:          "1.1",
		Section:     "Extreme Edge Cases",
		Description: description,
		Status:      parsers.TaskStatusPending,
	}

	// Create a TasksFile containing the task
	tasksFile := parsers.TasksFile{
		Version: 1,
		Tasks:   []parsers.Task{task},
	}

	// Marshal to JSON (this is what writeTasksJSONC does internally)
	jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent() failed for nested quotes: %v", err)
	}

	// Validate that the generated JSONC can be parsed
	if err := validateJSONCOutput(jsonData); err != nil {
		t.Errorf("validateJSONCOutput() failed for nested quotes: %v", err)
	}

	// Verify round-trip is lossless:
	// Strip JSONC comments (though we don't have any in this test)
	stripped := parsers.StripJSONComments(jsonData)

	// Unmarshal back to TasksFile
	var roundTrip parsers.TasksFile
	if err := json.Unmarshal(stripped, &roundTrip); err != nil {
		t.Fatalf(
			"round-trip unmarshal failed for nested quotes: %v\nJSON:\n%s",
			err,
			string(jsonData),
		)
	}

	// Verify we got exactly one task back
	if len(roundTrip.Tasks) != 1 {
		t.Fatalf("expected 1 task after round-trip, got %d", len(roundTrip.Tasks))
	}

	// Verify the description matches exactly (lossless round-trip)
	if roundTrip.Tasks[0].Description != description {
		t.Errorf(
			"round-trip lost data for nested quotes:\nOriginal: %q\nResult:   %q",
			description,
			roundTrip.Tasks[0].Description,
		)
	}

	// Verify all other fields match
	if roundTrip.Tasks[0].ID != task.ID {
		t.Errorf("Task ID mismatch: got %q, want %q", roundTrip.Tasks[0].ID, task.ID)
	}
	if roundTrip.Tasks[0].Section != task.Section {
		t.Errorf("Task Section mismatch: got %q, want %q", roundTrip.Tasks[0].Section, task.Section)
	}
	if roundTrip.Tasks[0].Status != task.Status {
		t.Errorf("Task Status mismatch: got %q, want %q", roundTrip.Tasks[0].Status, task.Status)
	}
	if roundTrip.Version != tasksFile.Version {
		t.Errorf("Version mismatch: got %d, want %d", roundTrip.Version, tasksFile.Version)
	}

	// Additional verification: Check that the JSON contains properly escaped quotes
	jsonStr := string(jsonData)

	// The JSON should contain the escaped version of the description
	// In JSON, backslashes are escaped as \\ and quotes as \"
	// So "He said \"Hello\"" becomes "He said \\\"Hello\\\"" in the JSON string
	if !strings.Contains(jsonStr, "Nested quotes") {
		t.Error("JSON output missing 'Nested quotes' text")
	}

	// Verify the JSON is valid by checking it contains the "description" field
	if !strings.Contains(jsonStr, "\"description\"") {
		t.Error("JSON output missing 'description' field")
	}
}
