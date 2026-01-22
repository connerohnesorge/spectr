package cmd

import (
	"encoding/json"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// TestJSONCValidation_RegexCharacters tests that JSONC validation
// correctly handles task descriptions containing regex pattern characters.
//
// This test ensures that regex metacharacters (.*  .+ .? ^$ [a-z] (foo|bar)
// \d+ \w+ \s+) are properly escaped during JSON marshalling and can be
// successfully round-tripped through the validation process.
//
// This corresponds to task 1.14 in the test-extreme-jsonc change proposal.
func TestJSONCValidation_RegexCharacters(
	t *testing.T,
) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "wildcard asterisk",
			description: "Match any: .*",
		},
		{
			name:        "one or more wildcard",
			description: "Match one or more: .+",
		},
		{
			name:        "optional wildcard",
			description: "Match optional: .?",
		},
		{
			name:        "anchors",
			description: "Line anchors: ^$ and ^start$ end^",
		},
		{
			name:        "character class",
			description: "Character class: [a-z] and [A-Z] and [0-9]",
		},
		{
			name:        "alternation",
			description: "Alternation: (foo|bar) or (yes|no|maybe)",
		},
		{
			name:        "digit pattern",
			description: "Digits: \\d+ and \\d* and \\d{2,4}",
		},
		{
			name:        "word pattern",
			description: "Words: \\w+ and \\w* and \\W+",
		},
		{
			name:        "whitespace pattern",
			description: "Whitespace: \\s+ and \\s* and \\S+",
		},
		{
			name:        "combined regex patterns",
			description: "Regex chars: .* .+ .? ^$ [a-z] (foo|bar) \\d+ \\w+ \\s+",
		},
		{
			name:        "complex regex email pattern",
			description: "Email regex: ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$",
		},
		{
			name:        "complex regex URL pattern",
			description: "URL regex: https?://[\\w.-]+(:[0-9]+)?(/[\\w/.-]*)?",
		},
		{
			name:        "lookahead and lookbehind",
			description: "Lookahead: (?=pattern) and lookbehind: (?<=pattern)",
		},
		{
			name:        "negated character class",
			description: "Negated: [^a-z] and [^0-9] and [^\\s]",
		},
		{
			name:        "quantifiers",
			description: "Quantifiers: {3} {2,5} {3,} * + ?",
		},
		{
			name:        "escape sequences",
			description: "Escapes: \\. \\* \\+ \\? \\[ \\] \\( \\) \\{ \\} \\^ \\$ \\|",
		},
		{
			name:        "boundary matchers",
			description: "Boundaries: \\b\\w+\\b and \\Bword\\B",
		},
		{
			name:        "backreferences",
			description: "Backrefs: (\\w+)\\s+\\1 and ([a-z])\\1",
		},
		{
			name:        "mixed regex and text",
			description: "Validate email with regex: ^[\\w.+-]+@[\\w.-]+\\.[a-z]{2,}$ before processing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the regex pattern in description
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.14",
						Section:     "Extreme Edge Cases",
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

			// Validate that the generated JSON can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf(
					"validateJSONCOutput failed for %s: %v",
					tt.name,
					err,
				)
			}

			// Verify round-trip: unmarshal and check that description is preserved
			var roundTrip parsers.TasksFile
			stripped := parsers.StripJSONComments(
				jsonData,
			)
			if err := json.Unmarshal(stripped, &roundTrip); err != nil {
				t.Fatalf(
					"round-trip unmarshal failed: %v",
					err,
				)
			}

			if len(roundTrip.Tasks) != 1 {
				t.Fatalf(
					"expected 1 task after round-trip, got %d",
					len(roundTrip.Tasks),
				)
			}

			if roundTrip.Tasks[0].Description != tt.description {
				t.Errorf(
					"round-trip failed to preserve regex description\noriginal: %q\nround-trip: %q",
					tt.description,
					roundTrip.Tasks[0].Description,
				)
			}

			// Verify all other fields are preserved
			if roundTrip.Tasks[0].ID != "1.14" {
				t.Errorf(
					"Task ID mismatch: got %q, want %q",
					roundTrip.Tasks[0].ID,
					"1.14",
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
		})
	}
}
