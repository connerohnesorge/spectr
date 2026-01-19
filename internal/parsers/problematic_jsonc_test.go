package parsers

import (
	"encoding/json"
	"testing"
)

func TestStripJSONCommentsProblematicInputs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Task 1.1: Quotes and single quotes",
			input:    `{"description": "Normal task with \"quotes\" and 'single quotes'"}`,
			expected: `{"description": "Normal task with \"quotes\" and 'single quotes'"}`,
		},
		{
			name:     "Task 1.2: Backslashes",
			input:    `{"description": "Task with backslashes: C:\\Users\\test\\path and \\\\network\\share"}`,
			expected: `{"description": "Task with backslashes: C:\\Users\\test\\path and \\\\network\\share"}`,
		},
		{
			name:     "Task 1.3: Newlines",
			input:    `{"description": "Task with newlines: This is line 1\nThis is line 2"}`,
			expected: `{"description": "Task with newlines: This is line 1\nThis is line 2"}`,
		},
		{
			name:     "Task 1.4: JSON-like content",
			input:    `{"description": "Task with special JSON chars: { \"key\": \"value\" } and [array]"}`,
			expected: `{"description": "Task with special JSON chars: { \"key\": \"value\" } and [array]"}`,
		},
		{
			name:     "Task 1.5: Unicode and emojis",
			input:    `{"description": "Task with unicode: 🚀 Emoji test with 中文 Chinese chars and Ñoño"}`,
			expected: `{"description": "Task with unicode: 🚀 Emoji test with 中文 Chinese chars and Ñoño"}`,
		},
		{
			name:     "Task 1.6: Tabs",
			input:    `{"description": "Task with tabs:\tindented with tabs here"}`,
			expected: `{"description": "Task with tabs:\tindented with tabs here"}`,
		},
		{
			name:     "Task 1.7: Escape sequences",
			input:    `{"description": "Task with escape sequences: \\n \\t \\r \\b \\f \\\" \\\\ \\/ \\u0041"}`,
			expected: `{"description": "Task with escape sequences: \\n \\t \\r \\b \\f \\\" \\\\ \\/ \\u0041"}`,
		},
		{
			name:     "Task 1.8: Very long description",
			input:    `{"description": "Very long description that exceeds normal length limits to test if the system can handle descriptions that go on and on and on with lots of text that might cause buffer issues or memory problems or other edge cases that only appear with extremely long strings that contain multiple sentences and potentially hundreds of characters that need to be properly escaped and validated during the JSON marshaling process which should handle this gracefully without any errors or truncation or corruption of the data regardless of how verbose the task description becomes over time"}`,
			expected: `{"description": "Very long description that exceeds normal length limits to test if the system can handle descriptions that go on and on and on with lots of text that might cause buffer issues or memory problems or other edge cases that only appear with extremely long strings that contain multiple sentences and potentially hundreds of characters that need to be properly escaped and validated during the JSON marshaling process which should handle this gracefully without any errors or truncation or corruption of the data regardless of how verbose the task description becomes over time"}`,
		},
		{
			name:     "Task 1.9: Mixed special characters",
			input:    `{"description": "Mixed special chars: C:\\path\\to\\file.txt with \"quotes\" and {json} and \\u003chtml\\u003e and 50% discount"}`,
			expected: `{"description": "Mixed special chars: C:\\path\\to\\file.txt with \"quotes\" and {json} and \\u003chtml\\u003e and 50% discount"}`,
		},
		{
			name:     "Task 1.10: Backslash at end",
			input:    `{"description": "Backslash at end: path\\to\\directory\\"}`,
			expected: `{"description": "Backslash at end: path\\to\\directory\\"}`,
		},
		{
			name:     "Task 1.11: Control characters",
			input:    `{"description": "Control chars: ASCII control characters like \b (backspace) and \f (form feed)"}`,
			expected: `{"description": "Control chars: ASCII control characters like \b (backspace) and \f (form feed)"}`,
		},
		{
			name:     "Task 1.12: JSON-breaking syntax",
			input:    `{"description": "JSON edge case: Task with \"},{ which might break naive parsers"}`,
			expected: `{"description": "JSON edge case: Task with \"},{ which might break naive parsers"}`,
		},
		{
			name:     "Task 1.13: Unicode combining characters",
			input:    `{"description": "Unicode edge: Combining diacritics é café naïve Zürich"}`,
			expected: `{"description": "Unicode edge: Combining diacritics é café naïve Zürich"}`,
		},
		{
			name:     "Task 1.14: Math symbols",
			input:    `{"description": "Math symbols: ∑ ∫ √ ± × ÷ ≠ ≈ ≤ ≥"}`,
			expected: `{"description": "Math symbols: ∑ ∫ √ ± × ÷ ≠ ≈ ≤ ≥"}`,
		},
		{
			name:     "Task 1.15: Currency and symbols",
			input:    `{"description": "Currency and symbols: $100 €50 £30 ¥1000 © ® ™ § ¶"}`,
			expected: `{"description": "Currency and symbols: $100 €50 £30 ¥1000 © ® ™ § ¶"}`,
		},
		{
			name:     "Task 2.1: HTML-like content with unicode",
			input:    `{"description": "Output/return: \u003cpromise\u003e\"This is the last task: COMPLETE\"\u003c/promise\u003e\n"}`,
			expected: `{"description": "Output/return: \u003cpromise\u003e\"This is the last task: COMPLETE\"\u003c/promise\u003e\n"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that StripJSONComments preserves the content
			result := StripJSONComments([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("StripJSONComments() = %v, want %v", string(result), tt.expected)
			}

			// Test that the result is valid JSON
			var parsed map[string]any
			if err := json.Unmarshal(result, &parsed); err != nil {
				t.Errorf("Result is not valid JSON: %v\nResult: %s", err, string(result))
			}
		})
	}
}

func TestStripJSONCommentsWithCommentsAndEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "JSON with comments and quotes",
			input: `// Comment with "quotes"
{"description": "Task with \"quotes\" and 'single quotes'"}`,
			expected: `
{"description": "Task with \"quotes\" and 'single quotes'"}`,
		},
		{
			name: "JSON with comments and backslashes",
			input: `// Path: C:\test
{"description": "Task with backslashes: C:\\Users\\test\\path"}`,
			expected: `
{"description": "Task with backslashes: C:\\Users\\test\\path"}`,
		},
		{
			name: "JSON with multi-line comments and unicode",
			input: `/* Multi-line
   comment with 🚀 emoji */
{"description": "Task with unicode: 🚀 Emoji test"}`,
			expected: `
{"description": "Task with unicode: 🚀 Emoji test"}`,
		},
		{
			name: "JSON with comment-like content in strings",
			input: `// This is a comment
{"description": "This has // in the string"} // inline comment`,
			expected: `
{"description": "This has // in the string"} `,
		},
		{
			name: "Complex case with all edge cases",
			input: `// Header with \"quotes\" and C:\\path
/* Multi-line comment */
{"description": "Complex task with \"quotes\", C:\\path\\to\\file, 🚀 emoji, and // comment-like text"}`,
			expected: `

{"description": "Complex task with \"quotes\", C:\\path\\to\\file, 🚀 emoji, and // comment-like text"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := StripJSONComments([]byte(tt.input))
			if string(result) != tt.expected {
				t.Errorf("StripJSONComments() = %v, want %v", string(result), tt.expected)
			}

			// Verify the result is valid JSON
			var parsed map[string]any
			if err := json.Unmarshal(result, &parsed); err != nil {
				t.Errorf("Result is not valid JSON: %v\nResult: %s", err, string(result))
			}
		})
	}
}

// TestRoundTripJSONCMarshaling tests that JSONC can be marshaled and unmarshaled correctly
func TestRoundTripJSONCMarshaling(t *testing.T) {
	tests := []struct {
		name  string
		tasks []Task
	}{
		{
			name: "Tasks with problematic characters",
			tasks: []Task{
				{
					ID:          "1.1",
					Section:     "Edge Case Testing",
					Description: `Normal task with "quotes" and 'single quotes'`,
					Status:      "pending",
				},
				{
					ID:          "1.2",
					Section:     "Edge Case Testing",
					Description: `Task with backslashes: C:\Users\test\path and \\network\share`,
					Status:      "pending",
				},
				{
					ID:          "1.5",
					Section:     "Edge Case Testing",
					Description: `Task with unicode: 🚀 Emoji test with 中文 Chinese chars and Ñoño`,
					Status:      "pending",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a tasks file
			tasksFile := TasksFile{
				Version: 1,
				Tasks:   tt.tasks,
			}

			// Marshal to JSON
			jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal tasks: %v", err)
			}

			// Strip comments (should be no-op for freshly marshaled JSON)
			stripped := StripJSONComments(jsonData)

			// Parse back
			var parsed TasksFile
			if err := json.Unmarshal(stripped, &parsed); err != nil {
				t.Fatalf("Failed to unmarshal after stripping comments: %v", err)
			}

			// Verify round-trip integrity
			if parsed.Version != tasksFile.Version {
				t.Errorf("Version mismatch: got %d, want %d", parsed.Version, tasksFile.Version)
			}
			if len(parsed.Tasks) != len(tasksFile.Tasks) {
				t.Errorf(
					"Task count mismatch: got %d, want %d",
					len(parsed.Tasks),
					len(tasksFile.Tasks),
				)
			}
			for i, task := range tasksFile.Tasks {
				if i >= len(parsed.Tasks) {
					break
				}
				if parsed.Tasks[i].Description != task.Description {
					t.Errorf(
						"Task %d description mismatch:\ngot:  %q\nwant: %q",
						i,
						parsed.Tasks[i].Description,
						task.Description,
					)
				}
			}
		})
	}
}
