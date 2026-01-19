package cmd

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// Constants for test data
const (
	testSectionExtremeEdgeCases = "Extreme Edge Cases"
)

// Helper functions for JSONC validation tests

// createTaskWithDescription creates a test Task with a specific description
func createTaskWithDescription(
	id, section, description string,
	status parsers.TaskStatusValue,
) parsers.Task {
	return parsers.Task{
		ID:          id,
		Section:     section,
		Description: description,
		Status:      status,
	}
}

// marshalAndValidateRoundTrip marshals a Task to JSON, strips JSONC comments,
// unmarshals back, and validates that the description is preserved correctly.
// Returns an error if the round-trip fails or if the description doesn't match.
func marshalAndValidateRoundTrip(t *testing.T, task *parsers.Task) error {
	t.Helper()

	// Marshal the task to JSON
	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Errorf("Failed to marshal task: %v", err)

		return err
	}

	// Strip JSONC comments (simulating the comment stripping process)
	strippedData := parsers.StripJSONComments(jsonData)

	// Unmarshal back to a Task
	var unmarshaledTask parsers.Task
	if err := json.Unmarshal(strippedData, &unmarshaledTask); err != nil {
		t.Errorf("Failed to unmarshal task after comment stripping: %v", err)

		return err
	}

	// Verify the description matches
	if unmarshaledTask.Description != task.Description {
		t.Errorf(
			"Description mismatch after round-trip:\nOriginal: %q\nAfter:    %q",
			task.Description,
			unmarshaledTask.Description,
		)

		return err
	}

	// Verify all other fields match
	if unmarshaledTask.ID != task.ID {
		t.Errorf("ID mismatch: got %q, want %q", unmarshaledTask.ID, task.ID)
	}
	if unmarshaledTask.Section != task.Section {
		t.Errorf("Section mismatch: got %q, want %q", unmarshaledTask.Section, task.Section)
	}
	if unmarshaledTask.Status != task.Status {
		t.Errorf("Status mismatch: got %q, want %q", unmarshaledTask.Status, task.Status)
	}

	return nil
}

// createTasksFileWithTasks creates a TasksFile with the given tasks for testing
func createTasksFileWithTasks(tasks []parsers.Task) parsers.TasksFile {
	return parsers.TasksFile{
		Version: 1,
		Tasks:   tasks,
	}
}

// validateTasksFileSerialization validates that a TasksFile can be marshaled
// to JSON, have comments stripped, and then unmarshaled back with all task
// descriptions preserved.
func validateTasksFileSerialization(t *testing.T, tasksFile parsers.TasksFile) error {
	t.Helper()

	// Marshal the tasks file to JSON
	jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
	if err != nil {
		t.Errorf("Failed to marshal tasks file: %v", err)

		return err
	}

	// Strip JSONC comments
	strippedData := parsers.StripJSONComments(jsonData)

	// Unmarshal back to a TasksFile
	var unmarshaledFile parsers.TasksFile
	if err := json.Unmarshal(strippedData, &unmarshaledFile); err != nil {
		t.Errorf("Failed to unmarshal tasks file after comment stripping: %v", err)

		return err
	}

	// Verify version matches
	if unmarshaledFile.Version != tasksFile.Version {
		t.Errorf("Version mismatch: got %d, want %d", unmarshaledFile.Version, tasksFile.Version)
	}

	// Verify task count matches
	if len(unmarshaledFile.Tasks) != len(tasksFile.Tasks) {
		t.Errorf(
			"Task count mismatch: got %d, want %d",
			len(unmarshaledFile.Tasks),
			len(tasksFile.Tasks),
		)

		return nil
	}

	// Verify each task's description is preserved
	for i, originalTask := range tasksFile.Tasks {
		unmarshaledTask := unmarshaledFile.Tasks[i]
		if unmarshaledTask.Description != originalTask.Description {
			t.Errorf(
				"Task %d description mismatch:\nOriginal: %q\nAfter:    %q",
				i,
				originalTask.Description,
				unmarshaledTask.Description,
			)
		}
		if unmarshaledTask.ID != originalTask.ID {
			t.Errorf("Task %d ID mismatch: got %q, want %q", i, unmarshaledTask.ID, originalTask.ID)
		}
		if unmarshaledTask.Section != originalTask.Section {
			t.Errorf(
				"Task %d Section mismatch: got %q, want %q",
				i,
				unmarshaledTask.Section,
				originalTask.Section,
			)
		}
		if unmarshaledTask.Status != originalTask.Status {
			t.Errorf(
				"Task %d Status mismatch: got %q, want %q",
				i,
				unmarshaledTask.Status,
				originalTask.Status,
			)
		}
	}

	return nil
}

// TestHelperFunctions_Basic verifies the helper functions work correctly
func TestHelperFunctions_Basic(t *testing.T) {
	// Test createTaskWithDescription
	task := createTaskWithDescription(
		"1.1",
		"Test Section",
		"Test description",
		parsers.TaskStatusPending,
	)

	if task.ID != "1.1" {
		t.Errorf("Expected ID 1.1, got %s", task.ID)
	}
	if task.Description != "Test description" {
		t.Errorf("Expected description 'Test description', got %s", task.Description)
	}

	// Test marshalAndValidateRoundTrip
	if err := marshalAndValidateRoundTrip(t, &task); err != nil {
		t.Errorf("marshalAndValidateRoundTrip failed: %v", err)
	}

	// Test createTasksFileWithTasks
	tasks := []parsers.Task{task}
	tasksFile := createTasksFileWithTasks(tasks)

	if tasksFile.Version != 1 {
		t.Errorf("Expected version 1, got %d", tasksFile.Version)
	}
	if len(tasksFile.Tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasksFile.Tasks))
	}

	// Test validateTasksFileSerialization
	if err := validateTasksFileSerialization(t, tasksFile); err != nil {
		t.Errorf("validateTasksFileSerialization failed: %v", err)
	}
}

// TestJSONCValidation_SpecialCharacters tests that JSONC validation
// correctly handles special characters that require JSON escaping.
//
// This test ensures that task descriptions containing backslashes,
// quotes, newlines, tabs, and other control characters are properly
// escaped during JSON marshalling and can be successfully round-tripped
// through the validation process.
func TestJSONCValidation_SpecialCharacters(t *testing.T) {
	tests := []struct {
		name        string
		description string
		char        string // The special character being tested
	}{
		{
			name:        "backslash",
			description: "Task with backslash \\ character",
			char:        "\\",
		},
		{
			name:        "double quote",
			description: "Task with \"quoted\" text",
			char:        "\"",
		},
		{
			name:        "newline",
			description: "Task with\nnewline character",
			char:        "\n",
		},
		{
			name:        "tab",
			description: "Task with\ttab character",
			char:        "\t",
		},
		{
			name:        "carriage return",
			description: "Task with\rcarriage return",
			char:        "\r",
		},
		{
			name:        "backspace",
			description: "Task with\bbackspace character",
			char:        "\b",
		},
		{
			name:        "form feed",
			description: "Task with\fform feed character",
			char:        "\f",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with a task containing the special character
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1",
						Section:     "Test Section",
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
		})
	}
}

// TestJSONCValidation_Unicode verifies that JSONC validation preserves
// unicode characters and emoji correctly through round-trip conversion.
//
// This test ensures that task descriptions containing:
//   - Emoji (🚀, 🎉, 💻, 🔥)
//   - Chinese characters (你好世界)
//   - Arabic (مرحبا)
//   - Japanese (こんにちは)
//   - Special symbols (©, ®, ™, €)
//
// are properly encoded in JSON and can be successfully round-tripped
// through the validation process without losing data.
func TestJSONCValidation_Unicode(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "emoji rocket",
			description: "Deploy the application 🚀",
		},
		{
			name:        "emoji party",
			description: "Celebrate completion 🎉",
		},
		{
			name:        "emoji computer",
			description: "Write code on 💻",
		},
		{
			name:        "emoji fire",
			description: "This feature is 🔥",
		},
		{
			name:        "multiple emoji",
			description: "Deploy 🚀 and celebrate 🎉 with code 💻 that's 🔥",
		},
		{
			name:        "chinese characters",
			description: "你好世界 - Hello World in Chinese",
		},
		{
			name:        "arabic characters",
			description: "مرحبا - Hello in Arabic",
		},
		{
			name:        "japanese characters",
			description: "こんにちは - Hello in Japanese",
		},
		{
			name:        "copyright symbol",
			description: "Copyright © 2024 Company",
		},
		{
			name:        "registered trademark",
			description: "Product® is registered",
		},
		{
			name:        "trademark symbol",
			description: "Brand™ trademark",
		},
		{
			name:        "euro symbol",
			description: "Price: €99.99",
		},
		{
			name:        "mixed special symbols",
			description: "Symbols: © ® ™ € in one description",
		},
		{
			name:        "mixed unicode and ascii",
			description: "Update README.md with 你好 and add tests 🚀",
		},
		{
			name:        "all character types",
			description: "ASCII, 你好世界, مرحبا, こんにちは, © ® ™ €, 🚀🎉💻🔥",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the unicode description
			original := parsers.TasksFile{
				Version: 2,
				Tasks: []parsers.Task{
					{
						ID:          "1.1",
						Section:     "Test Section",
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(original, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent failed: %v", err)
			}

			// Validate that the generated JSONC can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf("validateJSONCOutput failed for %s: %v", tt.name, err)
			}

			// Verify round-trip is lossless:
			// Strip comments (though we don't have any in this test)
			stripped := parsers.StripJSONComments(jsonData)

			// Unmarshal back to TasksFile
			var result parsers.TasksFile
			if err := json.Unmarshal(stripped, &result); err != nil {
				t.Fatalf("round-trip unmarshal failed: %v", err)
			}

			// Verify we got exactly one task back
			if len(result.Tasks) != 1 {
				t.Fatalf("expected 1 task after round-trip, got %d", len(result.Tasks))
			}

			// Verify the description matches exactly (lossless round-trip)
			if result.Tasks[0].Description != original.Tasks[0].Description {
				t.Errorf(
					"round-trip lost data for unicode:\nOriginal: %q\nResult:   %q",
					original.Tasks[0].Description,
					result.Tasks[0].Description,
				)
			}

			// Verify all other fields match
			if result.Tasks[0].ID != original.Tasks[0].ID {
				t.Errorf("Task ID mismatch: got %q, want %q",
					result.Tasks[0].ID, original.Tasks[0].ID)
			}
			if result.Tasks[0].Section != original.Tasks[0].Section {
				t.Errorf("Task Section mismatch: got %q, want %q",
					result.Tasks[0].Section, original.Tasks[0].Section)
			}
			if result.Tasks[0].Status != original.Tasks[0].Status {
				t.Errorf("Task Status mismatch: got %q, want %q",
					result.Tasks[0].Status, original.Tasks[0].Status)
			}
			if result.Version != original.Version {
				t.Errorf("Version mismatch: got %d, want %d",
					result.Version, original.Version)
			}
		})
	}
}

// TestRoundTripConversion_Version2Hierarchical tests round-trip JSON serialization
// for version 2 hierarchical task files (parent/child structure with $ref:tasks-*.jsonc)
func TestRoundTripConversion_Version2Hierarchical(t *testing.T) {
	// Create a root tasks file (version 2) with parent tasks that have children
	rootTasksFile := parsers.TasksFile{
		Version: 2,
		Tasks: []parsers.Task{
			{
				ID:          "1",
				Section:     "Foundation",
				Description: "Foundation tasks",
				Status:      parsers.TaskStatusPending,
				Children:    "$ref:tasks-1.jsonc",
			},
			{
				ID:          "2",
				Section:     "Implementation",
				Description: "Implementation tasks",
				Status:      parsers.TaskStatusInProgress,
				Children:    "$ref:tasks-2.jsonc",
			},
			{
				ID:          "3",
				Section:     "Testing",
				Description: "Testing tasks with special chars: `code`, **bold**, [link](url)",
				Status:      parsers.TaskStatusCompleted,
				Children:    "$ref:tasks-3.jsonc",
			},
		},
		Includes: []string{"tasks-*.jsonc"},
	}

	// Create child tasks files (each references parent)
	childTasksFile1 := parsers.TasksFile{
		Version: 2,
		Parent:  "1",
		Tasks: []parsers.Task{
			{
				ID:          "1.1",
				Section:     "Foundation",
				Description: "Set up project structure",
				Status:      parsers.TaskStatusCompleted,
			},
			{
				ID:          "1.2",
				Section:     "Foundation",
				Description: "Create initial configuration with backticks: `config.yaml`",
				Status:      parsers.TaskStatusCompleted,
			},
		},
	}

	childTasksFile2 := parsers.TasksFile{
		Version: 2,
		Parent:  "2",
		Tasks: []parsers.Task{
			{
				ID:          "2.1",
				Section:     "Implementation",
				Description: "Implement core logic // with comment-like text",
				Status:      parsers.TaskStatusInProgress,
			},
			{
				ID:          "2.2",
				Section:     "Implementation",
				Description: "Add validation /* block comment style */",
				Status:      parsers.TaskStatusPending,
			},
		},
	}

	childTasksFile3 := parsers.TasksFile{
		Version: 2,
		Parent:  "3",
		Tasks: []parsers.Task{
			{
				ID:          "3.1",
				Section:     "Testing",
				Description: "Write unit tests with \"quotes\" and 'apostrophes'",
				Status:      parsers.TaskStatusCompleted,
			},
		},
	}

	// Helper function to validate round-trip for a TasksFile
	validateRoundTrip := func(t *testing.T, tasksFile parsers.TasksFile) {
		t.Helper()

		// Marshal the tasks file to JSON
		jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
		if err != nil {
			t.Fatalf("Failed to marshal tasks file: %v", err)
		}

		// Strip JSONC comments (simulating the comment stripping process)
		strippedData := parsers.StripJSONComments(jsonData)

		// Unmarshal back to a TasksFile
		var unmarshaledFile parsers.TasksFile
		if err := json.Unmarshal(strippedData, &unmarshaledFile); err != nil {
			t.Fatalf("Failed to unmarshal tasks file after comment stripping: %v", err)
		}

		// Verify version matches
		if unmarshaledFile.Version != tasksFile.Version {
			t.Errorf(
				"Version mismatch: got %d, want %d",
				unmarshaledFile.Version,
				tasksFile.Version,
			)
		}

		// Verify parent field is preserved (if present)
		if unmarshaledFile.Parent != tasksFile.Parent {
			t.Errorf("Parent mismatch: got %q, want %q", unmarshaledFile.Parent, tasksFile.Parent)
		}

		// Verify includes field is preserved (if present)
		if len(unmarshaledFile.Includes) != len(tasksFile.Includes) {
			t.Errorf(
				"Includes count mismatch: got %d, want %d",
				len(unmarshaledFile.Includes),
				len(tasksFile.Includes),
			)
		}
		for i, include := range unmarshaledFile.Includes {
			if i < len(tasksFile.Includes) && include != tasksFile.Includes[i] {
				t.Errorf(
					"Includes[%d] mismatch: got %q, want %q",
					i,
					include,
					tasksFile.Includes[i],
				)
			}
		}

		// Verify task count matches
		if len(unmarshaledFile.Tasks) != len(tasksFile.Tasks) {
			t.Fatalf(
				"Task count mismatch: got %d, want %d",
				len(unmarshaledFile.Tasks),
				len(tasksFile.Tasks),
			)
		}

		// Verify each task's fields are preserved
		for i, originalTask := range tasksFile.Tasks {
			unmarshaledTask := unmarshaledFile.Tasks[i]
			if unmarshaledTask.Description != originalTask.Description {
				t.Errorf(
					"Task %d description mismatch:\nOriginal: %q\nAfter:    %q",
					i,
					originalTask.Description,
					unmarshaledTask.Description,
				)
			}
			if unmarshaledTask.ID != originalTask.ID {
				t.Errorf(
					"Task %d ID mismatch: got %q, want %q",
					i,
					unmarshaledTask.ID,
					originalTask.ID,
				)
			}
			if unmarshaledTask.Section != originalTask.Section {
				t.Errorf(
					"Task %d Section mismatch: got %q, want %q",
					i,
					unmarshaledTask.Section,
					originalTask.Section,
				)
			}
			if unmarshaledTask.Status != originalTask.Status {
				t.Errorf(
					"Task %d Status mismatch: got %q, want %q",
					i,
					unmarshaledTask.Status,
					originalTask.Status,
				)
			}
			if unmarshaledTask.Children != originalTask.Children {
				t.Errorf(
					"Task %d Children mismatch: got %q, want %q",
					i,
					unmarshaledTask.Children,
					originalTask.Children,
				)
			}
		}
	}

	// Test root tasks file round-trip
	t.Run("root_tasks_file", func(t *testing.T) {
		validateRoundTrip(t, rootTasksFile)
	})

	// Test child tasks file 1 round-trip (with parent field)
	t.Run("child_tasks_file_1", func(t *testing.T) {
		validateRoundTrip(t, childTasksFile1)
	})

	// Test child tasks file 2 round-trip (with comment-like text in descriptions)
	t.Run("child_tasks_file_2", func(t *testing.T) {
		validateRoundTrip(t, childTasksFile2)
	})

	// Test child tasks file 3 round-trip (with quotes in descriptions)
	t.Run("child_tasks_file_3", func(t *testing.T) {
		validateRoundTrip(t, childTasksFile3)
	})
}

// TestValidateJSONCOutput_ErrorMessages tests that validation failures
// produce helpful error messages with context and suggestions.
func TestValidateJSONCOutput_ErrorMessages(t *testing.T) {
	tests := []struct {
		name             string
		invalidJSON      string
		expectInError    []string // Strings that should appear in the error message
		expectNotInError []string // Strings that should NOT appear in the error message
	}{
		{
			name:        "invalid character at known position",
			invalidJSON: `{"version": 1, "tasks": [{"id": "1", "description": "test\x"}]}`,
			expectInError: []string{
				"JSONC validation failed",
				"Problematic content near position",
				"Common causes:",
				"unescaped special characters",
				"bug in JSON escaping",
			},
		},
		{
			name:        "truncated JSON",
			invalidJSON: `{"version": 1, "tasks": [`,
			expectInError: []string{
				"JSONC validation failed",
				"Common causes:",
				"bug in JSON escaping",
			},
		},
		{
			name:        "unescaped backslash in middle of long content",
			invalidJSON: `{"version": 1, "tasks": [{"id": "1", "section": "Test", "description": "This is a long task description with an unescaped backslash \ in the middle that will cause parsing to fail", "status": "pending"}]}`,
			expectInError: []string{
				"JSONC validation failed",
				"Problematic content near position",
				"unescaped special characters",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate the invalid JSON
			err := validateJSONCOutput([]byte(tt.invalidJSON))

			// Should always return an error
			if err == nil {
				t.Fatal("Expected error but got none")
			}

			errMsg := err.Error()

			// Check that expected strings are present
			for _, expected := range tt.expectInError {
				if !strings.Contains(errMsg, expected) {
					t.Errorf(
						"Expected error message to contain %q, but it didn't.\nFull error: %s",
						expected,
						errMsg,
					)
				}
			}

			// Check that unexpected strings are not present
			for _, unexpected := range tt.expectNotInError {
				if strings.Contains(errMsg, unexpected) {
					t.Errorf(
						"Expected error message NOT to contain %q, but it did.\nFull error: %s",
						unexpected,
						errMsg,
					)
				}
			}

			// Verify the error is not empty and has reasonable length
			if len(errMsg) < 50 {
				t.Errorf("Error message seems too short: %q", errMsg)
			}
		})
	}
}

// TestBuildJSONCValidationError_ContextSize tests that the error context
// display properly handles edge cases like errors near the start or end of data.
func TestBuildJSONCValidationError_ContextSize(t *testing.T) {
	tests := []struct {
		name        string
		invalidJSON string
		description string
	}{
		{
			name:        "error at start",
			invalidJSON: `{invalid}`,
			description: "Error near the beginning should not panic",
		},
		{
			name:        "error at end",
			invalidJSON: strings.Repeat("a", 200) + `{"version": 1, "tasks": [}`,
			description: "Error near the end should not panic",
		},
		{
			name:        "very short invalid JSON",
			invalidJSON: `{`,
			description: "Very short JSON should not panic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic
			err := validateJSONCOutput([]byte(tt.invalidJSON))

			if err == nil {
				t.Fatal("Expected error but got none")
			}

			// Verify we got a reasonable error message
			errMsg := err.Error()
			if errMsg == "" {
				t.Error("Error message is empty")
			}

			if !strings.Contains(errMsg, "JSONC validation failed") {
				t.Errorf("Error message doesn't contain expected header: %s", errMsg)
			}
		})
	}
}

// verifyRoundTrip checks that a round-trip conversion preserves data
func verifyRoundTrip(t *testing.T, originalFile, roundTrippedFile *parsers.TasksFile) {
	// Verify version
	if roundTrippedFile.Version != originalFile.Version {
		t.Errorf(
			"Version mismatch: got %d, want %d",
			roundTrippedFile.Version,
			originalFile.Version,
		)
	}

	// Verify task count
	if len(roundTrippedFile.Tasks) != len(originalFile.Tasks) {
		t.Fatalf(
			"Task count mismatch: got %d, want %d",
			len(roundTrippedFile.Tasks),
			len(originalFile.Tasks),
		)
	}

	// Verify each task
	for i := range originalFile.Tasks {
		verifyTaskPreserved(t, i, &roundTrippedFile.Tasks[i], &originalFile.Tasks[i])
	}

	// Verify parent
	if roundTrippedFile.Parent != originalFile.Parent {
		t.Errorf("Parent mismatch: got %q, want %q", roundTrippedFile.Parent, originalFile.Parent)
	}

	// Verify includes
	verifyIncludesPreserved(t, originalFile.Includes, roundTrippedFile.Includes)
}

// verifyTaskPreserved checks that a task is preserved correctly
func verifyTaskPreserved(t *testing.T, idx int, roundTripped, original *parsers.Task) {
	if roundTripped.ID != original.ID {
		t.Errorf("Task %d ID mismatch: got %q, want %q", idx, roundTripped.ID, original.ID)
	}
	if roundTripped.Section != original.Section {
		t.Errorf(
			"Task %d Section mismatch: got %q, want %q",
			idx,
			roundTripped.Section,
			original.Section,
		)
	}
	if roundTripped.Description != original.Description {
		t.Errorf(
			"Task %d Description mismatch:\nOriginal: %q\nAfter:    %q",
			idx,
			original.Description,
			roundTripped.Description,
		)
	}
	if roundTripped.Status != original.Status {
		t.Errorf(
			"Task %d Status mismatch: got %q, want %q",
			idx,
			roundTripped.Status,
			original.Status,
		)
	}
	if roundTripped.Children != original.Children {
		t.Errorf(
			"Task %d Children mismatch: got %q, want %q",
			idx,
			roundTripped.Children,
			original.Children,
		)
	}
}

// verifyIncludesPreserved checks that includes are preserved
func verifyIncludesPreserved(t *testing.T, original, roundTripped []string) {
	if len(roundTripped) != len(original) {
		t.Errorf("Includes count mismatch: got %d, want %d", len(roundTripped), len(original))

		return
	}
	for i, inc := range original {
		if roundTripped[i] != inc {
			t.Errorf("Includes[%d] mismatch: got %q, want %q", i, roundTripped[i], inc)
		}
	}
}

// TestRoundTripConversion_RealWorldData tests round-trip conversion using
// actual archived tasks.jsonc files to ensure validation works with production data.
func TestRoundTripConversion_RealWorldData(t *testing.T) {
	// Get testdata directory and walk archived changes
	testDataDir := GetTestDataDir(t)
	archivedDir := filepath.Join(testDataDir, "integration", "changes", "archive")

	// Find all tasks.jsonc files in archived changes
	entries, err := os.ReadDir(archivedDir)
	if err != nil {
		t.Fatalf("Failed to read archived changes directory: %v", err)
	}

	var archivedFiles []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		// Check for tasks.jsonc in each archived change subdirectory
		tasksPath := filepath.Join(archivedDir, entry.Name(), "tasks.jsonc")
		if _, err := os.Stat(tasksPath); err == nil {
			archivedFiles = append(archivedFiles, tasksPath)
		}
	}

	if len(archivedFiles) == 0 {
		t.Skip(
			"No archived test files found in testdata/integration/changes/archive/",
		)
	}

	for _, filePath := range archivedFiles {
		t.Run(filepath.Base(filePath), func(t *testing.T) {
			// Read the archived tasks.jsonc file
			tasksFile, err := parsers.ReadTasksJson(filePath)
			if err != nil {
				t.Fatalf("Failed to read archived tasks.jsonc: %v", err)
			}

			// Marshal back to JSON
			jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal tasks file: %v", err)
			}

			// Strip JSONC comments (simulating the round-trip through JSONC processing)
			strippedData := parsers.StripJSONComments(jsonData)

			// Unmarshal back to a TasksFile
			var roundTrippedFile parsers.TasksFile
			if err := json.Unmarshal(strippedData, &roundTrippedFile); err != nil {
				t.Fatalf("Failed to unmarshal after round-trip: %v", err)
			}

			// Verify the round-trip is lossless
			verifyRoundTrip(t, tasksFile, &roundTrippedFile)
		})
	}
}

// TestJSONCValidation_FormatStrings tests that format string patterns commonly
// found in programming contexts are properly escaped and preserved in JSONC.
// This includes printf-style format specifiers (%s, %d, %x), shell-style variable
// substitutions (${var}, #{var}), and template variable patterns ({{var}}).
//
// These patterns are common in task descriptions related to code implementation,
// configuration files, and shell scripts. The test ensures they round-trip correctly
// through JSON marshaling/unmarshaling without data loss or corruption.
func TestJSONCValidation_FormatStrings(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "printf-style format specifiers",
			description: "Format string: %s %d %x %n ${var} #{var} {{var}}",
		},
		{
			name:        "single format specifiers",
			description: "Use %s for strings and %d for integers",
		},
		{
			name:        "hex and pointer formats",
			description: "Memory address: %p, Hex: %x, Octal: %o",
		},
		{
			name:        "dangerous %n specifier",
			description: "Warning: %n can write to memory in C printf",
		},
		{
			name:        "shell variable substitution",
			description: "Export PATH=${PATH}:/usr/local/bin",
		},
		{
			name:        "bash parameter expansion",
			description: "Use ${var:-default} for default values",
		},
		{
			name:        "ruby/perl hash variable",
			description: "Access hash with #{variable} interpolation",
		},
		{
			name:        "handlebars template",
			description: "Render with {{title}} and {{content}} placeholders",
		},
		{
			name:        "jinja2 template",
			description: "Template uses {{var}} for variables and {% for %} for loops",
		},
		{
			name:        "mixed format patterns",
			description: "Log: printf(%s, ${USER}) {{timestamp}}",
		},
		{
			name:        "nested braces",
			description: "Complex: {{outer {{inner}} }} and ${a${b}}",
		},
		{
			name:        "percent signs",
			description: "Complete: 100% done, use %% to escape percent",
		},
		{
			name:        "all specifiers combined",
			description: "Format: %s %d %x %n ${var} #{var} {{var}} %p %o %%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the format string description
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.13",
						Section:     testSectionExtremeEdgeCases,
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
			if roundTrip.Tasks[0].ID != tasksFile.Tasks[0].ID {
				t.Errorf("Task ID mismatch: got %q, want %q",
					roundTrip.Tasks[0].ID, tasksFile.Tasks[0].ID)
			}
			if roundTrip.Tasks[0].Section != tasksFile.Tasks[0].Section {
				t.Errorf("Task Section mismatch: got %q, want %q",
					roundTrip.Tasks[0].Section, tasksFile.Tasks[0].Section)
			}
			if roundTrip.Tasks[0].Status != tasksFile.Tasks[0].Status {
				t.Errorf("Task Status mismatch: got %q, want %q",
					roundTrip.Tasks[0].Status, tasksFile.Tasks[0].Status)
			}
		})
	}
}

// TestJSONCValidation_EdgeCases tests edge cases that might cause json.Marshal issues.
// These tests validate that unusual but valid input strings are correctly handled
// during JSON marshaling and unmarshaling, with no data loss or corruption.
func TestJSONCValidation_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		description string
		shouldPass  bool // true if we expect successful round-trip, false if we expect error
	}{
		{
			name:        "empty string",
			description: "",
			shouldPass:  true,
		},
		{
			name:        "single space",
			description: " ",
			shouldPass:  true,
		},
		{
			name:        "multiple spaces",
			description: "   ",
			shouldPass:  true,
		},
		{
			name:        "tab only",
			description: "\t",
			shouldPass:  true,
		},
		{
			name:        "multiple tabs",
			description: "\t\t",
			shouldPass:  true,
		},
		{
			name:        "newline only",
			description: "\n",
			shouldPass:  true,
		},
		{
			name:        "multiple newlines",
			description: "\n\n",
			shouldPass:  true,
		},
		{
			name:        "carriage return only",
			description: "\r",
			shouldPass:  true,
		},
		{
			name:        "CRLF line ending",
			description: "\r\n",
			shouldPass:  true,
		},
		{
			name:        "mixed whitespace",
			description: " \t\n ",
			shouldPass:  true,
		},
		{
			name:        "whitespace sandwich",
			description: "  text  ",
			shouldPass:  true,
		},
		{
			name:        "very long description 10KB",
			description: generateLongString(10 * 1024),
			shouldPass:  true,
		},
		{
			name:        "very long description 50KB",
			description: generateLongString(50 * 1024),
			shouldPass:  true,
		},
		{
			name:        "very long description 100KB",
			description: generateLongString(100 * 1024),
			shouldPass:  true,
		},
		{
			name:        "mixed special characters",
			description: "Task\nwith\ttabs\\and\"quotes",
			shouldPass:  true,
		},
		{
			name:        "escaped backslash and quote",
			description: "Path: C:\\Program Files\\App\\file.txt with \"quotes\"",
			shouldPass:  true,
		},
		{
			name:        "all JSON escape characters",
			description: "\\\"\n\r\t\b\f",
			shouldPass:  true,
		},
		{
			name:        "control character null (ASCII 0)",
			description: "text\x00text",
			shouldPass:  true,
		},
		{
			name:        "control character bell (ASCII 7)",
			description: "text\x07text",
			shouldPass:  true,
		},
		{
			name:        "control character backspace (ASCII 8)",
			description: "text\btext",
			shouldPass:  true,
		},
		{
			name:        "control character vertical tab (ASCII 11)",
			description: "text\vtext",
			shouldPass:  true,
		},
		{
			name:        "control character form feed (ASCII 12)",
			description: "text\ftext",
			shouldPass:  true,
		},
		{
			name:        "control character escape (ASCII 27)",
			description: "text\x1btext",
			shouldPass:  true,
		},
		{
			name:        "multiple control characters",
			description: "\x01\x02\x03\x04\x05\x06",
			shouldPass:  true,
		},
		{
			name:        "text with embedded nulls",
			description: "before\x00middle\x00after",
			shouldPass:  true,
		},
		{
			name:        "all printable ASCII",
			description: generatePrintableASCII(),
			shouldPass:  true,
		},
		{
			name:        "repeated backslashes",
			description: "\\\\\\\\\\\\\\\\",
			shouldPass:  true,
		},
		{
			name:        "repeated quotes",
			description: "\"\"\"\"\"\"\"\"",
			shouldPass:  true,
		},
		{
			name:        "alternating backslash and quote",
			description: "\\\"\\\"\\\"\\\"",
			shouldPass:  true,
		},
		{
			name:        "newlines between text",
			description: "line1\nline2\nline3\nline4\nline5",
			shouldPass:  true,
		},
		{
			name:        "unicode BOM",
			description: "\uFEFFtext with BOM",
			shouldPass:  true,
		},
		{
			name:        "zero-width characters",
			description: "text\u200Bwith\u200Czero\u200Dwidth\uFEFFchars",
			shouldPass:  true,
		},
		{
			name:        "right-to-left override",
			description: "text\u202Eright-to-left",
			shouldPass:  true,
		},
		{
			name:        "combining characters",
			description: "e\u0301\u0302\u0303", // e with multiple diacritics
			shouldPass:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the edge case description
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.1",
						Section:     "Edge Case Testing",
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
			if err != nil {
				if tt.shouldPass {
					t.Errorf(
						"Expected successful marshal for %q, but got error: %v",
						tt.name,
						err,
					)
				}

				return
			}

			// Validate that the generated JSON can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				if tt.shouldPass {
					t.Errorf("validateJSONCOutput failed for %s: %v", tt.name, err)
				}

				return
			}

			// Verify round-trip: unmarshal and check that description is preserved
			var roundTrip parsers.TasksFile
			stripped := parsers.StripJSONComments(jsonData)
			if err := json.Unmarshal(stripped, &roundTrip); err != nil {
				if tt.shouldPass {
					t.Errorf("round-trip unmarshal failed for %s: %v", tt.name, err)
				}

				return
			}

			if len(roundTrip.Tasks) != 1 {
				t.Errorf("expected 1 task after round-trip, got %d", len(roundTrip.Tasks))

				return
			}

			if roundTrip.Tasks[0].Description != tt.description {
				if tt.shouldPass {
					t.Errorf(
						"round-trip failed to preserve description for %s\noriginal: %q\nround-trip: %q",
						tt.name,
						tt.description,
						roundTrip.Tasks[0].Description,
					)
				}
			} else if !tt.shouldPass {
				t.Errorf(
					"Expected failure for %q, but round-trip succeeded",
					tt.name,
				)
			}
		})
	}
}

// generateLongString creates a string of the specified length filled with
// repeating patterns of text to simulate very long task descriptions.
func generateLongString(length int) string {
	if length <= 0 {
		return ""
	}

	pattern := "This is a test task description with some special chars: \\n\\t\"quotes\" and more. "
	patternLen := len(pattern)

	// Calculate how many full patterns we need
	fullPatterns := length / patternLen
	remainder := length % patternLen

	// Build the string
	result := ""
	for range fullPatterns {
		result += pattern
	}

	// Add the remainder
	if remainder > 0 {
		result += pattern[:remainder]
	}

	return result
}

// generatePrintableASCII generates a string containing all printable ASCII
// characters (space through tilde, ASCII 32-126).
func generatePrintableASCII() string {
	result := ""
	for i := 32; i <= 126; i++ {
		result += string(rune(i))
	}

	return result
}

// TestJSONCValidation_CommentLikeStrings verifies that task descriptions
// containing JSONC comment syntax (// and /* */) are preserved correctly
// when marshaled to JSON and parsed back via StripJSONComments.
//
// This is critical because StripJSONComments should only remove actual
// comments, not comment-like strings that appear inside JSON string values.
func TestJSONCValidation_CommentLikeStrings(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "single-line comment syntax",
			description: "Task with // comment",
		},
		{
			name:        "block comment start",
			description: "Task with /* comment",
		},
		{
			name:        "block comment end",
			description: "Task with */ comment",
		},
		{
			name:        "full block comment",
			description: "Task with /* block */ comment",
		},
		{
			name:        "URL with slashes",
			description: "http://example.com",
		},
		{
			name:        "URL with path",
			description: "https://example.com/path/to/resource",
		},
		{
			name:        "mixed comment syntax at start",
			description: "// Update the documentation",
		},
		{
			name:        "mixed comment syntax at end",
			description: "Update the documentation //",
		},
		{
			name:        "block comment mid-sentence",
			description: "Update /* inline comment */ documentation",
		},
		{
			name:        "multiple slashes",
			description: "Path: ////multiple////slashes////",
		},
		{
			name:        "file path with backslashes and comment",
			description: "C:\\Program Files\\App\\file.txt // Windows path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a task with the test description
			task := parsers.Task{
				ID:          "1.1",
				Section:     "Test Section",
				Description: tt.description,
				Status:      parsers.TaskStatusPending,
			}

			// Create a TasksFile with the task
			tasksFile := parsers.TasksFile{
				Version: 2,
				Tasks:   []parsers.Task{task},
			}

			// Marshal to JSON (this is what writeTasksJSONC does)
			jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent() error = %v", err)
			}

			// Verify JSON is valid by unmarshaling it directly
			var directUnmarshal parsers.TasksFile
			if err := json.Unmarshal(jsonData, &directUnmarshal); err != nil {
				t.Fatalf("json.Unmarshal() error = %v, JSON:\n%s", err, string(jsonData))
			}

			// Verify direct unmarshal preserves description
			if directUnmarshal.Tasks[0].Description != tt.description {
				t.Errorf(
					"Direct unmarshal failed to preserve description:\nwant: %q\ngot:  %q",
					tt.description,
					directUnmarshal.Tasks[0].Description,
				)
			}

			// Now test with StripJSONComments (simulating JSONC parsing)
			// This is the critical test: StripJSONComments should NOT remove
			// comment-like syntax when it's inside a quoted string
			strippedData := parsers.StripJSONComments(jsonData)

			var roundTripUnmarshal parsers.TasksFile
			if err := json.Unmarshal(strippedData, &roundTripUnmarshal); err != nil {
				t.Fatalf(
					"json.Unmarshal() after StripJSONComments error = %v\nOriginal JSON:\n%s\nStripped JSON:\n%s",
					err,
					string(jsonData),
					string(strippedData),
				)
			}

			// Verify the description is preserved exactly after round-trip
			if roundTripUnmarshal.Tasks[0].Description != tt.description {
				t.Errorf(
					"Round-trip conversion failed to preserve description:\nwant: %q\ngot:  %q\nOriginal JSON:\n%s\nStripped JSON:\n%s",
					tt.description,
					roundTripUnmarshal.Tasks[0].Description,
					string(jsonData),
					string(strippedData),
				)
			}

			// Verify all other fields are preserved
			if roundTripUnmarshal.Tasks[0].ID != task.ID {
				t.Errorf(
					"Round-trip failed to preserve ID: want %q, got %q",
					task.ID,
					roundTripUnmarshal.Tasks[0].ID,
				)
			}
			if roundTripUnmarshal.Tasks[0].Section != task.Section {
				t.Errorf(
					"Round-trip failed to preserve Section: want %q, got %q",
					task.Section,
					roundTripUnmarshal.Tasks[0].Section,
				)
			}
			if roundTripUnmarshal.Tasks[0].Status != task.Status {
				t.Errorf(
					"Round-trip failed to preserve Status: want %q, got %q",
					task.Status,
					roundTripUnmarshal.Tasks[0].Status,
				)
			}
			if roundTripUnmarshal.Version != tasksFile.Version {
				t.Errorf(
					"Round-trip failed to preserve Version: want %d, got %d",
					tasksFile.Version,
					roundTripUnmarshal.Version,
				)
			}
		})
	}
}

// TestRoundTripConversion_AllFields tests that all Task fields are preserved
// during JSON marshal/unmarshal round-trip conversion. This ensures that the
// JSONC serialization process correctly handles all field types including
// ID, Section, Description, Status, and Children.
func TestRoundTripConversion_AllFields(t *testing.T) {
	tests := []struct {
		name      string
		tasksFile parsers.TasksFile
	}{
		{
			name: "all task fields with pending status",
			tasksFile: parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.1",
						Section:     "Implementation",
						Description: "Implement feature X with special chars: \\ \" \n \t",
						Status:      parsers.TaskStatusPending,
					},
				},
			},
		},
		{
			name: "all task fields with in_progress status",
			tasksFile: parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "2.1",
						Section:     "Testing",
						Description: "Test feature Y with unicode: é ñ 中",
						Status:      parsers.TaskStatusInProgress,
					},
				},
			},
		},
		{
			name: "all task fields with completed status",
			tasksFile: parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "3.1",
						Section:     "Documentation",
						Description: "Document feature Z",
						Status:      parsers.TaskStatusCompleted,
					},
				},
			},
		},
		{
			name: "task with children reference",
			tasksFile: parsers.TasksFile{
				Version: 2,
				Tasks: []parsers.Task{
					{
						ID:          "4",
						Section:     "Round-Trip Testing",
						Description: "Comprehensive round-trip tests",
						Status:      parsers.TaskStatusInProgress,
						Children:    "$ref:tasks-4.jsonc",
					},
				},
			},
		},
		{
			name: "multiple tasks with various statuses",
			tasksFile: parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1",
						Section:     "Setup",
						Description: "Initialize project",
						Status:      parsers.TaskStatusCompleted,
					},
					{
						ID:          "2",
						Section:     "Development",
						Description: "Implement core features",
						Status:      parsers.TaskStatusInProgress,
					},
					{
						ID:          "3",
						Section:     "Testing",
						Description: "Write comprehensive tests",
						Status:      parsers.TaskStatusPending,
					},
				},
			},
		},
		{
			name: "version 2 root file with parent and includes",
			tasksFile: parsers.TasksFile{
				Version:  2,
				Parent:   "root",
				Includes: []string{"tasks-1.jsonc", "tasks-2.jsonc", "tasks-3.jsonc"},
				Tasks: []parsers.Task{
					{
						ID:          "1",
						Section:     "Phase 1",
						Description: "First phase tasks",
						Status:      parsers.TaskStatusCompleted,
						Children:    "$ref:tasks-1.jsonc",
					},
				},
			},
		},
		{
			name: "version 2 child file with parent",
			tasksFile: parsers.TasksFile{
				Version: 2,
				Parent:  "4",
				Tasks: []parsers.Task{
					{
						ID:          "4.1",
						Section:     "Round-Trip Testing",
						Description: "Test all fields",
						Status:      parsers.TaskStatusPending,
					},
					{
						ID:          "4.2",
						Section:     "Round-Trip Testing",
						Description: "Test hierarchical format",
						Status:      parsers.TaskStatusPending,
					},
				},
			},
		},
		{
			name: "empty includes array",
			tasksFile: parsers.TasksFile{
				Version:  2,
				Includes: nil,
				Tasks: []parsers.Task{
					{
						ID:          "1",
						Section:     "Solo",
						Description: "Single task",
						Status:      parsers.TaskStatusPending,
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			jsonData, err := json.MarshalIndent(tt.tasksFile, "", "  ")
			if err != nil {
				t.Fatalf("Failed to marshal tasks file: %v", err)
			}

			// Strip JSONC comments (simulating the comment stripping process)
			strippedData := parsers.StripJSONComments(jsonData)

			// Unmarshal back to TasksFile
			var roundTripped parsers.TasksFile
			if err := json.Unmarshal(strippedData, &roundTripped); err != nil {
				t.Fatalf("Failed to unmarshal after round-trip: %v", err)
			}

			// Deep equality check on Version
			if roundTripped.Version != tt.tasksFile.Version {
				t.Errorf(
					"Version mismatch: got %d, want %d",
					roundTripped.Version,
					tt.tasksFile.Version,
				)
			}

			// Deep equality check on Parent
			if roundTripped.Parent != tt.tasksFile.Parent {
				t.Errorf(
					"Parent mismatch: got %q, want %q",
					roundTripped.Parent,
					tt.tasksFile.Parent,
				)
			}

			// Deep equality check on Includes
			if len(roundTripped.Includes) != len(tt.tasksFile.Includes) {
				t.Errorf(
					"Includes length mismatch: got %d, want %d",
					len(roundTripped.Includes),
					len(tt.tasksFile.Includes),
				)
			} else {
				for i := range tt.tasksFile.Includes {
					if roundTripped.Includes[i] != tt.tasksFile.Includes[i] {
						t.Errorf(
							"Includes[%d] mismatch: got %q, want %q",
							i,
							roundTripped.Includes[i],
							tt.tasksFile.Includes[i],
						)
					}
				}
			}

			// Deep equality check on Tasks array
			if len(roundTripped.Tasks) != len(tt.tasksFile.Tasks) {
				t.Fatalf(
					"Tasks length mismatch: got %d, want %d",
					len(roundTripped.Tasks),
					len(tt.tasksFile.Tasks),
				)
			}

			// Verify each task field individually
			for i, originalTask := range tt.tasksFile.Tasks {
				rtTask := roundTripped.Tasks[i]

				if rtTask.ID != originalTask.ID {
					t.Errorf("Task[%d].ID mismatch: got %q, want %q", i, rtTask.ID, originalTask.ID)
				}

				if rtTask.Section != originalTask.Section {
					t.Errorf(
						"Task[%d].Section mismatch: got %q, want %q",
						i,
						rtTask.Section,
						originalTask.Section,
					)
				}

				if rtTask.Description != originalTask.Description {
					t.Errorf(
						"Task[%d].Description mismatch:\n  got:  %q\n  want: %q",
						i,
						rtTask.Description,
						originalTask.Description,
					)
				}

				if rtTask.Status != originalTask.Status {
					t.Errorf(
						"Task[%d].Status mismatch: got %q, want %q",
						i,
						rtTask.Status,
						originalTask.Status,
					)
				}

				if rtTask.Children != originalTask.Children {
					t.Errorf(
						"Task[%d].Children mismatch: got %q, want %q",
						i,
						rtTask.Children,
						originalTask.Children,
					)
				}
			}
		})
	}
}

// TestJSONCValidation_HTMLInjection tests that JSONC validation correctly
// handles HTML and XSS injection attempts in task descriptions.
//
// This test ensures that potentially malicious HTML tags, script tags, and
// HTML comments are properly escaped during JSON marshaling and can be
// successfully round-tripped through the validation process without losing data.
//
// The test covers:
//   - <script> tags with XSS payloads
//   - HTML comments (<!-- -->)
//   - Various HTML tags (<div>, <img>, etc.)
//   - Mixed HTML and text content
//   - HTML entities and special characters
//
// This is critical for security and data integrity when task descriptions
// might contain code examples, documentation snippets, or malicious input.
func TestJSONCValidation_HTMLInjection(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "simple script tag",
			description: "<script>alert('xss')</script>",
		},
		{
			name:        "script tag with double quotes",
			description: "<script>alert(\"xss\")</script>",
		},
		{
			name:        "html comment",
			description: "<!-- comment -->",
		},
		{
			name:        "html comment with text",
			description: "Text before <!-- comment --> text after",
		},
		{
			name:        "script and comment together",
			description: "<script>alert('xss')</script> and <!-- comment -->",
		},
		{
			name:        "img tag with onerror",
			description: "<img src=x onerror=alert('xss')>",
		},
		{
			name:        "div tag",
			description: "<div>content</div>",
		},
		{
			name:        "nested html tags",
			description: "<div><span>nested</span></div>",
		},
		{
			name:        "html with attributes",
			description: "<a href=\"javascript:alert('xss')\">click</a>",
		},
		{
			name:        "multiple script tags",
			description: "<script>alert('1')</script><script>alert('2')</script>",
		},
		{
			name:        "html entities",
			description: "&lt;script&gt;alert('xss')&lt;/script&gt;",
		},
		{
			name:        "mixed html and special chars",
			description: "<script>alert(\"test\\nwith\\ttabs\")</script>",
		},
		{
			name:        "html comment multiline",
			description: "<!-- This is a\nmultiline\ncomment -->",
		},
		{
			name:        "style tag with css",
			description: "<style>body { background: red; }</style>",
		},
		{
			name:        "iframe injection",
			description: "<iframe src=\"javascript:alert('xss')\"></iframe>",
		},
		{
			name:        "html with unicode",
			description: "<script>alert('🚀');</script>",
		},
		{
			name:        "xss with quotes",
			description: "<script>alert('He said \"Hello\"')</script>",
		},
		{
			name:        "svg xss vector",
			description: "<svg onload=alert('xss')>",
		},
		{
			name:        "data uri xss",
			description: "<img src=\"data:text/html,<script>alert('xss')</script>\">",
		},
		{
			name:        "all html special chars",
			description: "<>&\"'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the HTML injection description
			original := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.5",
						Section:     testSectionExtremeEdgeCases,
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(original, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent failed: %v", err)
			}

			// Validate that the generated JSONC can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf("validateJSONCOutput failed for %s: %v", tt.name, err)
			}

			// Verify round-trip is lossless:
			// Strip comments (though we don't have any in this test)
			stripped := parsers.StripJSONComments(jsonData)

			// Unmarshal back to TasksFile
			var result parsers.TasksFile
			if err := json.Unmarshal(stripped, &result); err != nil {
				t.Fatalf("round-trip unmarshal failed: %v", err)
			}

			// Verify we got exactly one task back
			if len(result.Tasks) != 1 {
				t.Fatalf("expected 1 task after round-trip, got %d", len(result.Tasks))
			}

			// Verify the description matches exactly (lossless round-trip)
			if result.Tasks[0].Description != original.Tasks[0].Description {
				t.Errorf(
					"round-trip lost data for HTML injection:\nOriginal: %q\nResult:   %q",
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

// TestJSONCValidation_QuoteBombardment verifies that task descriptions
// containing pathological numbers of consecutive quotes are properly escaped
// and can be successfully round-tripped through JSON marshaling.
//
// This test validates the edge case of "quote bombardment" - many consecutive
// quote characters that could potentially break JSON escaping if not handled
// correctly. Each quote must be properly escaped as \" in the JSON output.
//
// This test covers task 1.7 from test-extreme-jsonc change proposal:
// "Quote bombardment: \"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\"\""
func TestJSONCValidation_QuoteBombardment(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "10 consecutive quotes",
			description: `""""""""""`,
		},
		{
			name:        "20 consecutive quotes",
			description: `""""""""""""""""""""`,
		},
		{
			name:        "38 consecutive quotes (from task 1.7)",
			description: `""""""""""""""""""""""""""""""""""""""`,
		},
		{
			name:        "50 consecutive quotes",
			description: strings.Repeat(`"`, 50),
		},
		{
			name:        "100 consecutive quotes",
			description: strings.Repeat(`"`, 100),
		},
		{
			name:        "quotes with text prefix",
			description: `Quote bombardment: """"""""""""""""""""""""""""""""""""""""""`,
		},
		{
			name:        "quotes with text suffix",
			description: `"""""""""""""""""""""""""""""""""""""" in description`,
		},
		{
			name:        "quotes with text in middle",
			description: `"""""""""""" text """""""""""`,
		},
		{
			name:        "alternating quotes and spaces",
			description: `" " " " " " " " " "`,
		},
		{
			name:        "quotes with newlines",
			description: "\"\"\"\"\"\n\"\"\"\"\"\n\"\"\"\"\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the quote bombardment description
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.7",
						Section:     testSectionExtremeEdgeCases,
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
					"round-trip failed to preserve quote bombardment\noriginal length: %d quotes\nround-trip length: %d quotes\noriginal: %q\nround-trip: %q",
					strings.Count(tt.description, `"`),
					strings.Count(roundTrip.Tasks[0].Description, `"`),
					tt.description,
					roundTrip.Tasks[0].Description,
				)
			}

			// Verify all other fields are preserved
			if roundTrip.Tasks[0].ID != "1.7" {
				t.Errorf("Task ID mismatch: got %q, want %q", roundTrip.Tasks[0].ID, "1.7")
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

			// Additional validation: count escaped quotes in JSON
			jsonStr := string(jsonData)
			// In JSON, each quote in the description should be escaped as \"
			// So we should find the escaped version in the JSON string
			if !strings.Contains(jsonStr, `\"`) {
				t.Error("JSON output does not contain escaped quotes as expected")
			}
		})
	}
}

// TestJSONCValidation_PathTraversal verifies that JSONC validation preserves
// path traversal strings correctly through round-trip conversion.
//
// This test ensures that task descriptions containing path traversal patterns:
//   - Unix path traversal: ../../../etc/passwd
//   - Windows path traversal: ..\..\..\..\windows\system32
//
// are properly escaped in JSON and can be successfully round-tripped
// through the validation process without losing data.
func TestJSONCValidation_PathTraversal(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "unix path traversal simple",
			description: "../../../etc/passwd",
		},
		{
			name:        "unix path traversal deep",
			description: "../../../../../../../../etc/passwd",
		},
		{
			name:        "windows path traversal simple",
			description: "..\\..\\..\\..\\windows\\system32",
		},
		{
			name:        "windows path traversal deep",
			description: "..\\..\\..\\..\\..\\..\\..\\..\\windows\\system32\\config\\sam",
		},
		{
			name:        "mixed path traversal",
			description: "../../../etc/passwd and ..\\..\\..\\..\\windows\\system32",
		},
		{
			name:        "path traversal with file",
			description: "Read file: ../../../etc/passwd or ..\\..\\..\\..\\windows\\system32\\drivers\\etc\\hosts",
		},
		{
			name:        "url encoded path traversal",
			description: "%2e%2e%2f%2e%2e%2f%2e%2e%2fetc%2fpasswd",
		},
		{
			name:        "double encoded path traversal",
			description: "%252e%252e%252f%252e%252e%252f",
		},
		{
			name:        "path traversal with null byte",
			description: "../../../etc/passwd\x00.jpg",
		},
		{
			name:        "path traversal with unicode",
			description: "..\\..\\..\\..\\windows\\system32\\驱动程序",
		},
		{
			name:        "forward and backward slash mix",
			description: "..\\../..\\../etc/passwd",
		},
		{
			name:        "multiple consecutive dots",
			description: "....//....//etc/passwd",
		},
		{
			name:        "path with spaces",
			description: ".. / .. / .. / etc / passwd",
		},
		{
			name:        "absolute and relative paths",
			description: "/etc/passwd and ../../../etc/shadow and C:\\Windows\\System32",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the path traversal description
			original := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.6",
						Section:     testSectionExtremeEdgeCases,
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(original, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent failed: %v", err)
			}

			// Validate that the generated JSONC can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf("validateJSONCOutput failed for %s: %v", tt.name, err)
			}

			// Verify round-trip is lossless:
			// Strip comments (though we don't have any in this test)
			stripped := parsers.StripJSONComments(jsonData)

			// Unmarshal back to TasksFile
			var result parsers.TasksFile
			if err := json.Unmarshal(stripped, &result); err != nil {
				t.Fatalf("round-trip unmarshal failed: %v", err)
			}

			// Verify we got exactly one task back
			if len(result.Tasks) != 1 {
				t.Fatalf("expected 1 task after round-trip, got %d", len(result.Tasks))
			}

			// Verify the description matches exactly (lossless round-trip)
			if result.Tasks[0].Description != original.Tasks[0].Description {
				t.Errorf(
					"round-trip lost data for path traversal:\nOriginal: %q\nResult:   %q",
					original.Tasks[0].Description,
					result.Tasks[0].Description,
				)
			}

			// Verify all other fields match
			if result.Tasks[0].ID != original.Tasks[0].ID {
				t.Errorf("Task ID mismatch: got %q, want %q",
					result.Tasks[0].ID, original.Tasks[0].ID)
			}
			if result.Tasks[0].Section != original.Tasks[0].Section {
				t.Errorf("Task Section mismatch: got %q, want %q",
					result.Tasks[0].Section, original.Tasks[0].Section)
			}
			if result.Tasks[0].Status != original.Tasks[0].Status {
				t.Errorf("Task Status mismatch: got %q, want %q",
					result.Tasks[0].Status, original.Tasks[0].Status)
			}
			if result.Version != original.Version {
				t.Errorf("Version mismatch: got %d, want %d",
					result.Version, original.Version)
			}
		})
	}
}

// TestJSONCValidation_AllPrintableASCII tests that all printable ASCII characters
// (ASCII 32-126) are correctly handled in JSONC output. This is the specific test
// for task 1.15 from the test-extreme-jsonc change proposal.
//
// The test verifies that the exact string from task 1.15:
// !"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\]^_`abcdefghijklmnopqrstuvwxyz{|}~
// can be marshaled to JSON, parsed back, and round-trip correctly without data loss.
func TestJSONCValidation_AllPrintableASCII(t *testing.T) {
	// The exact description from task 1.15
	task115Description := "All printable ASCII: !\"#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~"

	// Also test with the generated version to ensure consistency
	generatedASCII := generatePrintableASCII()

	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "task 1.15 exact description",
			description: task115Description,
		},
		{
			name:        "generated printable ASCII",
			description: generatedASCII,
		},
		{
			name:        "printable ASCII with prefix",
			description: "Testing all chars: " + generatedASCII,
		},
		{
			name:        "printable ASCII with suffix",
			description: generatedASCII + " <- all printable ASCII chars",
		},
		{
			name:        "printable ASCII in middle of text",
			description: "Start -> " + generatedASCII + " <- End",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the test description
			original := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.15",
						Section:     testSectionExtremeEdgeCases,
						Description: tt.description,
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(original, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent failed: %v", err)
			}

			// Validate that the generated JSONC can be parsed back
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf("validateJSONCOutput failed: %v", err)
			}

			// Verify round-trip: Strip comments and unmarshal
			stripped := parsers.StripJSONComments(jsonData)

			var result parsers.TasksFile
			if err := json.Unmarshal(stripped, &result); err != nil {
				t.Fatalf("round-trip unmarshal failed: %v", err)
			}

			// Verify we got exactly one task back
			if len(result.Tasks) != 1 {
				t.Fatalf("expected 1 task after round-trip, got %d", len(result.Tasks))
			}

			// Verify the description matches exactly (lossless round-trip)
			if result.Tasks[0].Description != original.Tasks[0].Description {
				t.Errorf(
					"round-trip lost data:\nOriginal: %q\nResult:   %q",
					original.Tasks[0].Description,
					result.Tasks[0].Description,
				)
			}

			// Verify all other fields match
			if result.Tasks[0].ID != original.Tasks[0].ID {
				t.Errorf("Task ID mismatch: got %q, want %q",
					result.Tasks[0].ID, original.Tasks[0].ID)
			}
			if result.Tasks[0].Section != original.Tasks[0].Section {
				t.Errorf("Task Section mismatch: got %q, want %q",
					result.Tasks[0].Section, original.Tasks[0].Section)
			}
			if result.Tasks[0].Status != original.Tasks[0].Status {
				t.Errorf("Task Status mismatch: got %q, want %q",
					result.Tasks[0].Status, original.Tasks[0].Status)
			}
			if result.Version != original.Version {
				t.Errorf("Version mismatch: got %d, want %d",
					result.Version, original.Version)
			}
		})
	}
}

// TestJSONCValidation_JSONInjectionAttempt tests that task descriptions
// containing JSON injection attempts are properly escaped and do not corrupt
// the generated JSONC structure.
//
// A malicious or malformed task description like:
//
//	"},{"id":"injected","status":"hacked
//
// should be treated as a simple string value, not as a JSON structure break.
// This test verifies that:
//  1. The description is properly JSON-escaped when marshaled
//  2. The generated JSONC remains valid and parseable
//  3. Round-trip conversion preserves the exact injection attempt string
//  4. The JSONC structure is not corrupted (no extra tasks injected)
//
// This is a security-critical test for any system that processes user-provided
// text and serializes it to JSON.
func TestJSONCValidation_JSONInjectionAttempt(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "basic JSON injection - object break",
			description: "\"},{\"id\":\"injected\",\"status\":\"hacked",
		},
		{
			name:        "JSON injection with array break",
			description: "\"}],\"tasks\":[{\"id\":\"injected\",\"status\":\"hacked",
		},
		{
			name:        "JSON injection with nested structures",
			description: "\"},\"extra\":{\"nested\":{\"id\":\"injected\"}},\"more\":{",
		},
		{
			name:        "JSON injection closing and opening",
			description: "\"}},{\"id\":\"1.99\",\"description\":\"injected task\",\"status\":\"completed\",\"section\":\"Injected\"}",
		},
		{
			name:        "JSON key injection attempt",
			description: "\",\"admin\":true,\"injected\":\"",
		},
		{
			name:        "JSON array injection",
			description: "\"],[\"injected\",\"array\",\"elements\"],[\"",
		},
		{
			name:        "Mixed quotes and braces",
			description: "\"},\"key\":\"value\",\"another\":\"{\"nested\":\"object\"}",
		},
		{
			name:        "Escaped quote injection",
			description: "\\\"},\\\"injection\\\":\\\"attempt",
		},
		{
			name:        "Unicode escape injection",
			description: "\\u0022},{\\u0022id\\u0022:\\u0022injected",
		},
		{
			name:        "Complex nested injection",
			description: "\"},\"tasks\":[{\"id\":\"evil\",\"children\":\"$ref:evil.jsonc\"}],\"includes\":[\"evil.jsonc\"],\"version\":999,\"data\":{",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the injection attempt as a task description
			original := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.1",
						Section:     "Security Testing",
						Description: "Normal task before injection",
						Status:      parsers.TaskStatusPending,
					},
					{
						ID:          "1.2",
						Section:     "Security Testing",
						Description: tt.description, // The injection attempt
						Status:      parsers.TaskStatusPending,
					},
					{
						ID:          "1.3",
						Section:     "Security Testing",
						Description: "Normal task after injection",
						Status:      parsers.TaskStatusPending,
					},
				},
			}

			// Marshal to JSON (this is what writeTasksJSONC does internally)
			jsonData, err := json.MarshalIndent(original, "", "  ")
			if err != nil {
				t.Fatalf("json.MarshalIndent() error = %v", err)
			}

			// CRITICAL: Validate that the generated JSON is valid and parseable
			// If the injection succeeded, this would fail or produce wrong structure
			if err := validateJSONCOutput(jsonData); err != nil {
				t.Errorf("validateJSONCOutput failed for %s: %v\nGenerated JSON:\n%s",
					tt.name, err, string(jsonData))
			}

			// Verify round-trip: unmarshal and check structure integrity
			stripped := parsers.StripJSONComments(jsonData)
			var roundTrip parsers.TasksFile
			if err := json.Unmarshal(stripped, &roundTrip); err != nil {
				t.Fatalf("round-trip unmarshal failed: %v\nJSON:\n%s", err, string(jsonData))
			}

			// CRITICAL: Verify that we still have exactly 3 tasks (not injected extras)
			if len(roundTrip.Tasks) != 3 {
				t.Errorf(
					"JSON injection corrupted structure: expected 3 tasks, got %d\nJSON:\n%s",
					len(roundTrip.Tasks),
					string(jsonData),
				)
			}

			// Verify that the injection attempt is preserved as a safe string
			if roundTrip.Tasks[1].Description != tt.description {
				t.Errorf(
					"round-trip failed to preserve injection attempt as safe string:\nOriginal: %q\nResult:   %q",
					tt.description,
					roundTrip.Tasks[1].Description,
				)
			}

			// Verify that other task fields are not affected by the injection
			if roundTrip.Tasks[0].Description != "Normal task before injection" {
				t.Errorf("Task before injection was corrupted: %q", roundTrip.Tasks[0].Description)
			}
			if len(roundTrip.Tasks) >= 3 &&
				roundTrip.Tasks[2].Description != "Normal task after injection" {
				t.Errorf("Task after injection was corrupted: %q", roundTrip.Tasks[2].Description)
			}

			// Verify task IDs are preserved (not replaced by injected IDs)
			if roundTrip.Tasks[1].ID != "1.2" {
				t.Errorf("Task ID corrupted: expected %q, got %q", "1.2", roundTrip.Tasks[1].ID)
			}

			// Verify task statuses are preserved (not changed to "hacked" or "completed")
			if roundTrip.Tasks[1].Status != parsers.TaskStatusPending {
				t.Errorf(
					"Task status corrupted: expected %q, got %q",
					parsers.TaskStatusPending,
					roundTrip.Tasks[1].Status,
				)
			}

			// Verify version is preserved (not changed by injection)
			if roundTrip.Version != 1 {
				t.Errorf("Version corrupted: expected 1, got %d", roundTrip.Version)
			}
		})
	}
}

// TestJSONCValidation_LiteralNewlineChars tests that task descriptions containing
// literal backslash-n sequences (not actual newline characters) are properly
// escaped in JSONC output and can be successfully round-tripped.
//
// This test covers task 1.11 from test-extreme-jsonc change proposal:
// "Literal newline char test: Line1\nLine2\nLine3 (with literal \n not actual newline)"
//
// The distinction is critical:
//   - Actual newline character: '\n' (single byte, ASCII 10)
//   - Literal backslash-n: '\\' + 'n' (two characters: backslash and letter n)
//
// When marshaling to JSON:
//   - Actual newline '\n' becomes "\\n" in JSON
//   - Literal backslash-n '\\n' becomes "\\\\n" in JSON (backslash escaped, then n)
func TestJSONCValidation_LiteralNewlineChars(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "literal backslash-n sequence",
			description: "Line1\\nLine2\\nLine3",
		},
		{
			name:        "literal backslash-t sequence",
			description: "Col1\\tCol2\\tCol3",
		},
		{
			name:        "literal backslash-r sequence",
			description: "Line1\\rLine2\\rLine3",
		},
		{
			name:        "mixed literal escape sequences",
			description: "Text\\nWith\\tMixed\\rLiteral\\bEscapes",
		},
		{
			name:        "literal backslash followed by various chars",
			description: "\\a\\b\\c\\d\\e\\f\\g\\h\\i\\j\\k\\l\\m\\n\\o\\p\\q\\r\\s\\t\\u\\v\\w\\x\\y\\z",
		},
		{
			name:        "literal backslash-n at start",
			description: "\\nStarts with literal backslash-n",
		},
		{
			name:        "literal backslash-n at end",
			description: "Ends with literal backslash-n\\n",
		},
		{
			name:        "multiple consecutive literal backslash-n",
			description: "Multiple\\n\\n\\nLiteral\\n\\n\\nNewlines",
		},
		{
			name:        "literal vs actual newline comparison",
			description: "Literal: \\n vs Actual: \n (mixed)",
		},
		{
			name:        "literal backslash-backslash-n (triple escape)",
			description: "Triple escape: \\\\nText",
		},
		{
			name:        "path with literal backslash-n",
			description: "C:\\new\\folder\\test.txt",
		},
		{
			name:        "regex with literal backslash-n",
			description: "Regex pattern: .*\\n.* matches lines with newline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the literal escape sequence
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.11",
						Section:     testSectionExtremeEdgeCases,
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

			// Validate that the generated JSONC can be parsed back
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

			// Verify the description is preserved EXACTLY (this is the critical test)
			if roundTrip.Tasks[0].Description != tt.description {
				t.Errorf(
					"round-trip failed to preserve literal escape sequences\noriginal: %q\nround-trip: %q\nJSON output:\n%s",
					tt.description,
					roundTrip.Tasks[0].Description,
					string(jsonData),
				)
			}

			// Additional verification: check the actual JSON string representation
			// For "Line1\\nLine2" input, JSON should contain "Line1\\\\nLine2"
			// (the backslash is escaped as \\, making it \\\\ in the JSON literal)
			jsonStr := string(jsonData)

			// Verify that literal backslash sequences are properly escaped
			// If input contains \n (backslash + n), JSON must contain \\n (escaped backslash + n)
			if !strings.Contains(tt.description, "\\n") {
				return
			}

			// The JSON representation should have the backslash escaped
			if !strings.Contains(jsonStr, "\\\\n") {
				t.Errorf(
					"JSON output does not properly escape literal \\n sequence\nInput: %q\nJSON:\n%s",
					tt.description,
					jsonStr,
				)
			}
		})
	}
}

// TestJSONCValidation_MixedBombardment tests pathological mixed quote and backslash patterns.
// This test ensures that task descriptions containing extreme combinations of quotes
// and backslashes (as seen in task 1.9 of test-extreme-jsonc) are properly escaped
// during JSON marshalling and can be successfully round-tripped through the validation process.
func TestJSONCValidation_MixedBombardment(t *testing.T) {
	tests := []struct {
		name        string
		description string
	}{
		{
			name:        "mixed bombardment from task 1.9",
			description: "Mixed bombardment: \\\"\\\"\\\"\\\"\\\\ \\\"\\\"\\\"\\\"\\\\ \\\"\\\"\\\"\\\"\\\\ \\\"\\\"\\\"\\\"\\\\",
		},
		{
			name:        "alternating quote backslash",
			description: "\\\"\\\"\\\"\\\"\\\"\\\"\\\"\\\"",
		},
		{
			name:        "quote backslash pairs",
			description: "\\\"\\\\\\\"\\\\\\\"\\\\\\\"\\\\",
		},
		{
			name:        "backslash then quotes",
			description: "\\\\\\\\\\\"\\\"\\\"\\\"",
		},
		{
			name:        "quotes then backslashes",
			description: "\\\"\\\"\\\"\\\"\\\\\\\\\\\\\\\\",
		},
		{
			name:        "mixed with text",
			description: "Start \\\"\\\"\\\"\\\"\\\\ middle \\\"\\\"\\\"\\\"\\\\ end",
		},
		{
			name:        "triple quote backslash pattern",
			description: "\\\"\\\"\\\"\\\\\\\\\\\"\\\"\\\"\\\\\\\\\\\"\\\"\\\"\\\\\\\\",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a TasksFile with the pathological description
			tasksFile := parsers.TasksFile{
				Version: 1,
				Tasks: []parsers.Task{
					{
						ID:          "1.9",
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
				t.Fatalf("round-trip unmarshal failed: %v\nJSON:\n%s", err, string(jsonData))
			}

			if len(roundTrip.Tasks) != 1 {
				t.Fatalf("expected 1 task after round-trip, got %d", len(roundTrip.Tasks))
			}

			if roundTrip.Tasks[0].Description != tt.description {
				t.Errorf(
					"round-trip failed to preserve description\noriginal: %q\nround-trip: %q\nJSON:\n%s",
					tt.description,
					roundTrip.Tasks[0].Description,
					string(jsonData),
				)
			}

			// Verify all other fields are preserved
			if roundTrip.Tasks[0].ID != "1.9" {
				t.Errorf("Task ID mismatch: got %q, want %q", roundTrip.Tasks[0].ID, "1.9")
			}
			if roundTrip.Tasks[0].Section != "Extreme Edge Cases" {
				t.Errorf(
					"Task Section mismatch: got %q, want %q",
					roundTrip.Tasks[0].Section,
					"Extreme Edge Cases",
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

// TestFinalCompletion validates the final task (2.1) from the test-extreme-jsonc
// change proposal. This is a comprehensive integration test that verifies the
// specific pathological edge case: HTML-like tags with escaped quotes and newlines.
//
// Task 2.1 requirement: The system SHALL validate JSONC output with pathological
// edge case inputs.
//
// Task 2.1 scenario: GIVEN a task description with pathological special characters
// WHEN the system generates JSONC THEN the output SHALL be parseable and round-trip
// correctly.
func TestFinalCompletion(t *testing.T) {
	// This is the exact description from task 2.1 in tasks.jsonc
	description := "Output/return: <promise>\"This is the last task: COMPLETE\"</promise>\n"

	// Create a task with the pathological description
	task := parsers.Task{
		ID:          "2.1",
		Section:     "99.99.99 Final section (after all tasks)",
		Description: description,
		Status:      parsers.TaskStatusPending,
	}

	// Create a TasksFile with the task
	tasksFile := parsers.TasksFile{
		Version: 1,
		Tasks:   []parsers.Task{task},
	}

	// Marshal to JSON (this is what writeTasksJSONC does)
	jsonData, err := json.MarshalIndent(tasksFile, "", "  ")
	if err != nil {
		t.Fatalf("json.MarshalIndent() failed: %v", err)
	}

	// Validate that the generated JSONC is valid JSON
	if err := validateJSONCOutput(jsonData); err != nil {
		t.Fatalf("validateJSONCOutput() failed: %v", err)
	}

	// Verify round-trip: unmarshal and check that description is preserved
	var roundTrip parsers.TasksFile
	stripped := parsers.StripJSONComments(jsonData)
	if err := json.Unmarshal(stripped, &roundTrip); err != nil {
		t.Fatalf("round-trip unmarshal failed: %v", err)
	}

	// Verify we got exactly one task back
	if len(roundTrip.Tasks) != 1 {
		t.Fatalf("expected 1 task after round-trip, got %d", len(roundTrip.Tasks))
	}

	// Verify the description matches exactly (lossless round-trip)
	if roundTrip.Tasks[0].Description != description {
		t.Errorf(
			"round-trip failed to preserve description\noriginal: %q\nround-trip: %q",
			description,
			roundTrip.Tasks[0].Description,
		)
	}

	// Verify all fields are preserved
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

	// Extra verification: Ensure specific characters are present
	if !strings.Contains(roundTrip.Tasks[0].Description, "<promise>") {
		t.Error("Missing <promise> tag in round-tripped description")
	}
	if !strings.Contains(roundTrip.Tasks[0].Description, "</promise>") {
		t.Error("Missing </promise> tag in round-tripped description")
	}
	if !strings.Contains(roundTrip.Tasks[0].Description, "\"This is the last task: COMPLETE\"") {
		t.Error("Missing quoted text in round-tripped description")
	}
	if !strings.HasSuffix(roundTrip.Tasks[0].Description, "\n") {
		t.Error("Missing trailing newline in round-tripped description")
	}

	// SUCCESS: This is the last task - COMPLETE
	t.Log(
		"Final validation test COMPLETE: HTML-like tags with pathological special characters are properly escaped and round-trip correctly",
	)
}
