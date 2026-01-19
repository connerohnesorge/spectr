package cmd

import (
	"encoding/json"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// TestJSONCValidation_UnicodeEdgeCases tests pathological Unicode edge cases including
// zero-width characters, RTL marks, and other tricky Unicode scenarios.
//
// This test ensures that task descriptions containing:
//   - Zero-width characters (ZWSP, ZWNJ, ZWJ)
//   - Right-to-left override marks (RLO, LRO, PDF)
//   - Bidirectional text marks
//   - Null bytes (if any)
//   - Combining characters
//
// are properly encoded in JSON and can be successfully round-tripped
// through the validation process without losing data.
func TestJSONCValidation_UnicodeEdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "zero-width space",
			description: "text\u200Bwith\u200Bzero\u200Bwidth\u200Bspaces",
		},
		{
			name:        "zero-width non-joiner",
			description: "text\u200Cwith\u200Czero\u200Cwidth\u200Cnon-joiner",
		},
		{
			name:        "zero-width joiner",
			description: "text\u200Dwith\u200Dzero\u200Dwidth\u200Djoiner",
		},
		{
			name:        "all three zero-width chars",
			description: "combined\u200B\u200C\u200Dzero-width characters",
		},
		{
			name:        "right-to-left override",
			description: "text\u202Eright-to-left\u202Cafter",
		},
		{
			name:        "left-to-right override",
			description: "text\u202Dleft-to-right\u202Cafter",
		},
		{
			name:        "right-to-left mark",
			description: "text\u200Fwith RTL mark",
		},
		{
			name:        "left-to-right mark",
			description: "text\u200Ewith LTR mark",
		},
		{
			name:        "bidirectional isolate markers",
			description: "text\u2066isolated\u2069text",
		},
		{
			name:        "mixed RTL and LTR marks",
			description: "English\u200Fעברית\u200EEnglish again",
		},
		{
			name:        "Arabic with RTL marks",
			description: "\u200Fمرحبا\u200F بالعالم",
		},
		{
			name:        "Hebrew with RTL marks",
			description: "\u200Fשלום\u200F עולם",
		},
		{
			name:        "null byte in middle",
			description: "text\x00with\x00null\x00bytes",
		},
		{
			name:        "null byte at start",
			description: "\x00starts with null",
		},
		{
			name:        "null byte at end",
			description: "ends with null\x00",
		},
		{
			name:        "multiple consecutive null bytes",
			description: "multiple\x00\x00\x00nulls",
		},
		{
			name:        "BOM byte order mark",
			description: "\uFEFFtext with BOM at start",
		},
		{
			name:        "BOM in middle",
			description: "text\uFEFFwith\uFEFFBOM\uFEFFin\uFEFFmiddle",
		},
		{
			name:        "combining diacritical marks",
			description: "e\u0301\u0302\u0303\u0304\u0305", // e with multiple combining marks
		},
		{
			name:        "combining marks on multiple chars",
			description: "a\u0301b\u0302c\u0303d\u0304e\u0305",
		},
		{
			name:        "all directional marks combined",
			description: "\u200E\u200F\u202A\u202B\u202C\u202D\u202E directional marks",
		},
		{
			name:        "soft hyphen",
			description: "super\u00ADcalifragilistic",
		},
		{
			name:        "non-breaking space",
			description: "non\u00A0breaking\u00A0space",
		},
		{
			name:        "narrow no-break space",
			description: "narrow\u202Fno-break\u202Fspace",
		},
		{
			name:        "word joiner",
			description: "word\u2060joiner\u2060test",
		},
		{
			name:        "zero-width no-break space (deprecated BOM)",
			description: "text\uFEFFwith\uFEFFzero-width\uFEFFno-break",
		},
		{
			name:        "variation selectors",
			description: "text\uFE0Ewith\uFE0Fvariation\uFE0Eselectors",
		},
		{
			name:        "interlinear annotation",
			description: "text\uFFF9annotation\uFFFAseparator\uFFFBterminator",
		},
		{
			name:        "mixed everything pathological",
			description: "\u200B\u200C\u200D\u200E\u200F\u202A\u202B\u202C\u202D\u202E\uFEFF\x00all\x00mixed\u0301\u0302",
		},
		{
			name:        "replacement character",
			description: "invalid\uFFFDcharacter",
		},
		{
			name:        "object replacement character",
			description: "object\uFFFCreplacement",
		},
		{
			name:        "line separator",
			description: "line\u2028separator",
		},
		{
			name:        "paragraph separator",
			description: "paragraph\u2029separator",
		},
		{
			name:        "hangul filler",
			description: "hangul\u3164filler",
		},
		{
			name:        "ideographic space",
			description: "ideographic\u3000space",
		},
		{
			name:        "mongolian vowel separator",
			description: "mongolian\u180Evowel\u180Eseparator",
		},
		{
			name:        "actual task description from spec",
			description: "Unicode edge: null bytes (if any), zero-width chars: \u200B\u200C\u200D, RTL marks: \u200F\u200E",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the Unicode edge case description
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.4",
						Section:     "Extreme Edge Cases",
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent failed: %v", err)
			}

			// Validate that the generated JSON can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf("validateJSONCOutput failed for %s: %v", tt.name, err)
			}

			// Verify round-trip: unmarshal and check that description is preserved
			var roundTrip parsers.TasksFile
			stripped := parsers.StripJSONComments(jsonData)
			if err := json.Unmarshal(stripped, &roundTrip); err != nil {
				t.Fatalf("round-trip unmarshal failed: %v", err)
			}

			if len(roundTrip.Tasks) != 1 {
				t.Fatalf("expected 1 task after round-trip, got %d", len(roundTrip.Tasks))
			}

			if roundTrip.Tasks[0].Description != tt.description {
				t.Errorf(
					"round-trip failed to preserve description\noriginal: %q\nround-trip: %q",
					tt.description,
					roundTrip.Tasks[0].Description,
				)
			}

			// Verify all other fields are preserved
			if roundTrip.Tasks[0].ID != "1.4" {
				t.Errorf("Task ID mismatch: got %q, want %q",
					roundTrip.Tasks[0].ID, "1.4")
			}
			if roundTrip.Tasks[0].Section != "Extreme Edge Cases" {
				t.Errorf("Task Section mismatch: got %q, want %q",
					roundTrip.Tasks[0].Section, "Extreme Edge Cases")
			}
			if roundTrip.Tasks[0].Status != parsers.TaskStatusPending {
				t.Errorf("Task Status mismatch: got %q, want %q",
					roundTrip.Tasks[0].Status, parsers.TaskStatusPending)
			}
		})
	}
}
