package cmd

import (
	"encoding/json"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// TestTask110_JSONCCommentInjection verifies that pathological JSONC comment
// injection patterns in task descriptions are properly escaped and can be
// round-tripped without data corruption.
//
// This test specifically covers task 1.10 from the test-extreme-jsonc change proposal:
// "JSONC comment injection: // comment */ or /* comment */ in description"
//
// The test ensures that:
//  1. Pathological comment injection patterns are preserved as-is (not stripped)
//  2. JSON marshaling properly escapes these strings
//  3. JSONC parsing with StripJSONComments does NOT remove comment syntax inside strings
//  4. Round-trip conversion (marshal → strip comments → unmarshal) is lossless
func TestTask110_JSONCCommentInjection(
	t *testing.T,
) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "task 1.10 exact pattern",
			description: "// comment */ or /* comment */ in description",
		},
		{
			name:        "comment termination injection attempt",
			description: "Text // */ should not terminate anything",
		},
		{
			name:        "block comment start in single-line",
			description: "// /* mixing comment types",
		},
		{
			name:        "nested block comment attempt",
			description: "/* outer /* inner */ outer still open",
		},
		{
			name:        "alternating mixed patterns",
			description: "// /* // /* // */ */ // */",
		},
		{
			name:        "triple slash with block",
			description: "/// comment /* block */ more",
		},
		{
			name:        "documentation style bombardment",
			description: "/***** */ /* */ // /* */ /*****/",
		},
		{
			name:        "pathological slash sequence",
			description: "////////// /* */ //////////",
		},
		{
			name:        "comment with escapes",
			description: "// comment\\nwith\\t/* block */\\rescapes",
		},
		{
			name:        "JSON-like with comments",
			description: "{\"key\": \"value\"} // /* */ looks like JSON with comment",
		},
		{
			name:        "URL with comment patterns",
			description: "http://example.com/* not a comment *///path",
		},
		{
			name:        "backslash before comment",
			description: "C:\\\\path // /* Windows */ style",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the pathological comment pattern
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.10",
						Section:     testSectionExtremeEdgeCases,
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(
				tasksFile,
				"",
				"  ",
			)
			if err != nil {
				t.Fatalf(
					"json.MarshalIndent failed: %v",
					err,
				)
			}

			// Validate that the generated JSON is syntactically correct
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf(
					"validateJSONCOutput failed for %s: %v",
					tt.name,
					err,
				)
			}

			// Verify round-trip: strip comments and unmarshal
			stripped := parsers.StripJSONComments(
				jsonData,
			)
			var roundTrip parsers.TasksFile
			if err := json.Unmarshal(stripped, &roundTrip); err != nil {
				t.Fatalf(
					"round-trip unmarshal failed: %v\nOriginal JSON:\n%s\nStripped JSON:\n%s",
					err,
					string(jsonData),
					string(stripped),
				)
			}

			// Verify we got exactly one task back
			if len(roundTrip.Tasks) != 1 {
				t.Fatalf(
					"expected 1 task after round-trip, got %d",
					len(roundTrip.Tasks),
				)
			}

			// CRITICAL: Verify the description is EXACTLY preserved
			// Pathological comment patterns should NOT be stripped - they're part of the data
			if roundTrip.Tasks[0].Description != tt.description {
				t.Errorf(
					"round-trip failed to preserve pathological comment description:\nOriginal:    %q\nRound-trip:  %q\nOriginal JSON:\n%s\nStripped JSON:\n%s",
					tt.description,
					roundTrip.Tasks[0].Description,
					string(jsonData),
					string(stripped),
				)
			}

			// Verify all other fields match
			if roundTrip.Tasks[0].ID != "1.10" {
				t.Errorf(
					"Task ID mismatch: got %q, want %q",
					roundTrip.Tasks[0].ID,
					"1.10",
				)
			}
			if roundTrip.Tasks[0].Section != testSectionExtremeEdgeCases {
				t.Errorf(
					"Task Section mismatch: got %q, want %q",
					roundTrip.Tasks[0].Section,
					testSectionExtremeEdgeCases,
				)
			}
			if roundTrip.Tasks[0].Status != parsers.TaskStatusPending {
				t.Errorf(
					"Task Status mismatch: got %q, want %q",
					roundTrip.Tasks[0].Status,
					parsers.TaskStatusPending,
				)
			}
			if roundTrip.Version != 1 {
				t.Errorf(
					"Version mismatch: got %d, want %d",
					roundTrip.Version,
					1,
				)
			}
		})
	}
}
