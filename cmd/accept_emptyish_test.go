package cmd

import (
	"encoding/json"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// TestJSONCValidation_EmptyishStrings tests that JSONC validation correctly
// handles empty strings and quote-like character combinations that appear
// "empty-ish" or contain various quote styles.
//
// This validates task 1.12 from the test-extreme-jsonc change proposal:
// "Empty-ish: \"\" and ” and “"
//
// These edge cases are important because they can confuse parsers, especially:
// - Empty strings (zero-length)
// - Strings containing only quote characters
// - Mixed quote styles within a single string
func TestJSONCValidation_EmptyishStrings(
	t *testing.T,
) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "truly empty string",
			description: "",
		},
		{
			name:        "double quotes only",
			description: "\"\"",
		},
		{
			name:        "single quotes only",
			description: "''",
		},
		{
			name:        "backticks only",
			description: "``",
		},
		{
			name:        "all three quote types together",
			description: "\"\" and '' and ``",
		},
		{
			name:        "nested empty quotes",
			description: "\"\"\"\"",
		},
		{
			name:        "single quotes with text",
			description: "Task with '' empty quotes",
		},
		{
			name:        "double quotes with text",
			description: "Task with \"\" empty quotes",
		},
		{
			name:        "backticks with text",
			description: "Task with `` empty backticks",
		},
		{
			name:        "mixed quotes at start",
			description: "\"\"''`` mixed empty quotes",
		},
		{
			name:        "mixed quotes at end",
			description: "mixed empty quotes \"\"''``",
		},
		{
			name:        "quotes separated by spaces",
			description: "\"\" '' ``",
		},
		{
			name:        "quotes with no spaces",
			description: "\"\"''``",
		},
		{
			name:        "multiple double quote pairs",
			description: "\"\" \"\" \"\" \"\"",
		},
		{
			name:        "multiple single quote pairs",
			description: "'' '' '' ''",
		},
		{
			name:        "multiple backtick pairs",
			description: "`` `` `` ``",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the empty-ish string description
			original := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.12",
						Section:     "Extreme Edge Cases",
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(
				original,
				"",
				"  ",
			)
			if err != nil {
				t.Fatalf(
					"json.MarshalIndent failed: %v",
					err,
				)
			}

			// Validate that the generated JSONC can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf(
					"validateJSONCOutput failed for %s: %v",
					tt.name,
					err,
				)
			}

			// Verify round-trip is lossless:
			// Strip comments (though we don't have any in this test)
			stripped := parsers.StripJSONComments(
				jsonData,
			)

			// Unmarshal back to TasksFile
			var result parsers.TasksFile
			if err := json.Unmarshal(stripped, &result); err != nil {
				t.Fatalf(
					"round-trip unmarshal failed: %v",
					err,
				)
			}

			// Verify we got exactly one task back
			if len(result.Tasks) != 1 {
				t.Fatalf(
					"expected 1 task after round-trip, got %d",
					len(result.Tasks),
				)
			}

			// Verify the description matches exactly (lossless round-trip)
			if result.Tasks[0].Description != original.Tasks[0].Description {
				t.Errorf(
					"round-trip lost data for empty-ish string:\nOriginal: %q\nResult:   %q",
					original.Tasks[0].Description,
					result.Tasks[0].Description,
				)
			}

			// Verify all other fields match
			if result.Tasks[0].ID != original.Tasks[0].ID {
				t.Errorf(
					"Task ID mismatch: got %q, want %q",
					result.Tasks[0].ID,
					original.Tasks[0].ID,
				)
			}
			if result.Tasks[0].Section != original.Tasks[0].Section {
				t.Errorf(
					"Task Section mismatch: got %q, want %q",
					result.Tasks[0].Section,
					original.Tasks[0].Section,
				)
			}
			if result.Tasks[0].Status != original.Tasks[0].Status {
				t.Errorf(
					"Task Status mismatch: got %q, want %q",
					result.Tasks[0].Status,
					original.Tasks[0].Status,
				)
			}
			if result.Version != original.Version {
				t.Errorf(
					"Version mismatch: got %d, want %d",
					result.Version,
					original.Version,
				)
			}
		})
	}
}
