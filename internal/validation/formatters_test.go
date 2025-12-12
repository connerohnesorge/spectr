package validation

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alecthomas/assert/v2"
)

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)

	return buf.String()
}

// TestPrintJSONReport_ValidReport tests printing a valid report as JSON
func TestPrintJSONReport_ValidReport(
	t *testing.T,
) {
	report := &ValidationReport{
		Valid: true,
		Summary: ValidationSummary{
			Errors:   0,
			Warnings: 0,
		},
		Issues: make([]ValidationIssue, 0),
	}

	output := captureOutput(func() {
		PrintJSONReport(report)
	})

	// Verify it's valid JSON
	var decoded ValidationReport
	err := json.Unmarshal(
		[]byte(output),
		&decoded,
	)
	assert.NoError(t, err)
	assert.True(t, decoded.Valid)
}

// TestPrintJSONReport_InvalidReport tests printing an invalid report as JSON
func TestPrintJSONReport_InvalidReport(
	t *testing.T,
) {
	report := &ValidationReport{
		Valid: false,
		Summary: ValidationSummary{
			Errors:   2,
			Warnings: 1,
		},
		Issues: []ValidationIssue{
			{
				Level:   "error",
				Path:    "spec.md",
				Message: "Missing required section",
			},
			{
				Level:   "error",
				Path:    "spec.md",
				Message: "Invalid format",
			},
			{
				Level:   "warning",
				Path:    "spec.md",
				Message: "Style issue",
			},
		},
	}

	output := captureOutput(func() {
		PrintJSONReport(report)
	})

	// Verify it's valid JSON
	var decoded ValidationReport
	err := json.Unmarshal(
		[]byte(output),
		&decoded,
	)
	assert.NoError(t, err)
	assert.False(t, decoded.Valid)
	assert.Equal(t, 2, decoded.Summary.Errors)
	assert.Equal(t, 1, decoded.Summary.Warnings)
	assert.Equal(t, 3, len(decoded.Issues))
}

// TestPrintJSONReport_NilReport tests handling of edge cases
func TestPrintJSONReport_NilReport(t *testing.T) {
	// This should handle gracefully
	output := captureOutput(func() {
		PrintJSONReport(nil)
	})

	// Should produce "null"
	assert.Contains(
		t,
		strings.TrimSpace(output),
		"null",
	)
}

// TestPrintHumanReport_ValidReport tests printing a valid report in human format
func TestPrintHumanReport_ValidReport(
	t *testing.T,
) {
	report := &ValidationReport{
		Valid: true,
		Summary: ValidationSummary{
			Errors:   0,
			Warnings: 0,
		},
		Issues: make([]ValidationIssue, 0),
	}

	output := captureOutput(func() {
		PrintHumanReport("test-spec", report)
	})

	assert.Contains(t, output, "✓")
	assert.Contains(t, output, "test-spec")
	assert.Contains(t, output, "valid")
}

// TestPrintHumanReport_InvalidReport tests printing an invalid report in human format
func TestPrintHumanReport_InvalidReport(
	t *testing.T,
) {
	report := &ValidationReport{
		Valid: false,
		Summary: ValidationSummary{
			Errors:   2,
			Warnings: 0,
		},
		Issues: []ValidationIssue{
			{
				Level:   "error",
				Path:    "spec.md",
				Message: "Missing required section",
			},
			{
				Level:   "error",
				Path:    "spec.md",
				Message: "Invalid format",
			},
		},
	}

	output := captureOutput(func() {
		PrintHumanReport("bad-spec", report)
	})

	assert.Contains(t, output, "✗")
	assert.Contains(t, output, "bad-spec")
	assert.Contains(t, output, "2 issue(s)")
	assert.Contains(
		t,
		output,
		"Missing required section",
	)
	assert.Contains(t, output, "Invalid format")
	assert.Contains(t, output, "error")
	assert.Contains(t, output, "spec.md")
}

// TestPrintHumanReport_WithWarnings tests printing report with warnings
func TestPrintHumanReport_WithWarnings(
	t *testing.T,
) {
	report := &ValidationReport{
		Valid: false,
		Summary: ValidationSummary{
			Errors:   1,
			Warnings: 1,
		},
		Issues: []ValidationIssue{
			{
				Level:   LevelError,
				Path:    "spec.md",
				Message: "Error message",
			},
			{
				Level:   LevelWarning,
				Path:    "spec.md",
				Message: "Warning message",
			},
		},
	}

	output := captureOutput(func() {
		PrintHumanReport("test-spec", report)
	})

	assert.Contains(t, output, "2 issue(s)")
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "[WARNING]")
	assert.Contains(t, output, "Error message")
	assert.Contains(t, output, "Warning message")
}

// TestPrintBulkJSONResults_ValidResults tests printing bulk results as JSON
func TestPrintBulkJSONResults_ValidResults(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "spec1",
			Type:  ItemTypeSpec,
			Valid: true,
			Report: &ValidationReport{
				Valid: true,
				Summary: ValidationSummary{
					Errors:   0,
					Warnings: 0,
				},
			},
		},
		{
			Name:  "spec2",
			Type:  ItemTypeSpec,
			Valid: true,
			Report: &ValidationReport{
				Valid: true,
				Summary: ValidationSummary{
					Errors:   0,
					Warnings: 0,
				},
			},
		},
	}

	output := captureOutput(func() {
		PrintBulkJSONResults(results)
	})

	// Verify it's valid JSON array
	var decoded []BulkResult
	err := json.Unmarshal(
		[]byte(output),
		&decoded,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(decoded))
	assert.Equal(t, "spec1", decoded[0].Name)
	assert.Equal(t, "spec2", decoded[1].Name)
	assert.True(t, decoded[0].Valid)
	assert.True(t, decoded[1].Valid)
}

// TestPrintBulkJSONResults_MixedResults tests printing mixed valid/invalid results
func TestPrintBulkJSONResults_MixedResults(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "valid-spec",
			Type:  ItemTypeSpec,
			Valid: true,
			Report: &ValidationReport{
				Valid: true,
			},
		},
		{
			Name:  "invalid-spec",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   "error",
						Path:    "spec.md",
						Message: "Error",
					},
				},
			},
		},
		{
			Name:  "error-spec",
			Type:  ItemTypeSpec,
			Valid: false,
			Error: "Failed to read file",
		},
	}

	output := captureOutput(func() {
		PrintBulkJSONResults(results)
	})

	var decoded []BulkResult
	err := json.Unmarshal(
		[]byte(output),
		&decoded,
	)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(decoded))
	assert.True(t, decoded[0].Valid)
	assert.False(t, decoded[1].Valid)
	assert.False(t, decoded[2].Valid)
	assert.Equal(
		t,
		"Failed to read file",
		decoded[2].Error,
	)
}

// TestPrintBulkJSONResults_EmptyResults tests printing empty results
func TestPrintBulkJSONResults_EmptyResults(
	t *testing.T,
) {
	results := make([]BulkResult, 0)

	output := captureOutput(func() {
		PrintBulkJSONResults(results)
	})

	var decoded []BulkResult
	err := json.Unmarshal(
		[]byte(output),
		&decoded,
	)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(decoded))
}

// TestPrintBulkHumanResults_AllValid tests printing all valid results
func TestPrintBulkHumanResults_AllValid(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "spec1",
			Type:  ItemTypeSpec,
			Valid: true,
		},
		{
			Name:  "spec2",
			Type:  ItemTypeSpec,
			Valid: true,
		},
		{
			Name:  "change1",
			Type:  ItemTypeChange,
			Valid: true,
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	assert.Contains(t, output, "✓ spec1 (spec)")
	assert.Contains(t, output, "✓ spec2 (spec)")
	assert.Contains(
		t,
		output,
		"✓ change1 (change)",
	)
	assert.Contains(
		t,
		output,
		"3 passed, 0 failed, 3 total",
	)
}

// TestPrintBulkHumanResults_AllInvalid tests printing all invalid results
func TestPrintBulkHumanResults_AllInvalid(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "bad-spec",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   "error",
						Path:    "spec.md",
						Message: "Missing section",
					},
				},
			},
		},
		{
			Name:  "error-change",
			Type:  ItemTypeChange,
			Valid: false,
			Error: "File not found",
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	assert.Contains(
		t,
		output,
		"✗ bad-spec (spec)",
	)
	assert.Contains(t, output, "1 issue(s)")
	assert.Contains(t, output, "Missing section")
	assert.Contains(
		t,
		output,
		"✗ error-change (change): File not found",
	)
	// With failures, summary includes error/warning breakdown
	assert.Contains(
		t,
		output,
		"0 passed, 2 failed",
	)
	assert.Contains(t, output, "2 total")
}

// TestPrintBulkHumanResults_MixedResults tests printing mixed results
func TestPrintBulkHumanResults_MixedResults(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "good-spec",
			Type:  ItemTypeSpec,
			Valid: true,
		},
		{
			Name:  "bad-spec",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   "error",
						Path:    "spec.md",
						Message: "Error 1",
					},
					{
						Level:   "error",
						Path:    "spec.md",
						Message: "Error 2",
					},
				},
			},
		},
		{
			Name:  "good-change",
			Type:  ItemTypeChange,
			Valid: true,
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	assert.Contains(
		t,
		output,
		"✓ good-spec (spec)",
	)
	assert.Contains(
		t,
		output,
		"✗ bad-spec (spec)",
	)
	assert.Contains(t, output, "2 issue(s)")
	assert.Contains(
		t,
		output,
		"✓ good-change (change)",
	)
	// With failures, summary includes error/warning breakdown
	assert.Contains(
		t,
		output,
		"2 passed, 1 failed",
	)
	assert.Contains(t, output, "3 total")
}

// TestPrintBulkHumanResults_EmptyResults tests printing empty results
func TestPrintBulkHumanResults_EmptyResults(
	t *testing.T,
) {
	results := make([]BulkResult, 0)

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	assert.Contains(
		t,
		output,
		"0 passed, 0 failed, 0 total",
	)
}

// TestPrintBulkHumanResults_DetailedIssues tests that issues are printed with details
func TestPrintBulkHumanResults_DetailedIssues(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "spec-with-issues",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   LevelError,
						Path:    "spec.md:10",
						Message: "Missing purpose section",
					},
					{
						Level:   LevelWarning,
						Path:    "spec.md:25",
						Message: "Scenario could be more specific",
					},
				},
			},
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	assert.Contains(t, output, "2 issue(s)")
	assert.Contains(t, output, "[ERROR]")
	assert.Contains(t, output, "spec.md:10")
	assert.Contains(
		t,
		output,
		"Missing purpose section",
	)
	assert.Contains(t, output, "[WARNING]")
	assert.Contains(t, output, "spec.md:25")
	assert.Contains(
		t,
		output,
		"Scenario could be more specific",
	)
}

// TestPrintBulkHumanResults_ErrorOnly tests results with error but no report
func TestPrintBulkHumanResults_ErrorOnly(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "broken-spec",
			Type:  ItemTypeSpec,
			Valid: false,
			Error: "Failed to read file: permission denied",
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	assert.Contains(
		t,
		output,
		"✗ broken-spec (spec): Failed to read file: permission denied",
	)
	// With failures, summary includes error/warning breakdown
	assert.Contains(
		t,
		output,
		"0 passed, 1 failed",
	)
	assert.Contains(t, output, "1 total")
}

// TestBulkResult_JSONSerialization tests that BulkResult serializes correctly
func TestBulkResult_JSONSerialization(
	t *testing.T,
) {
	result := BulkResult{
		Name:  "test-spec",
		Type:  ItemTypeSpec,
		Valid: false,
		Report: &ValidationReport{
			Valid: false,
			Issues: []ValidationIssue{
				{
					Level:   "error",
					Path:    "spec.md",
					Message: "Test error",
				},
			},
		},
		Error: "",
	}

	data, err := json.Marshal(result)
	assert.NoError(t, err)

	var decoded BulkResult
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, "test-spec", decoded.Name)
	assert.Equal(t, ItemTypeSpec, decoded.Type)
	assert.False(t, decoded.Valid)
	assert.NotZero(t, decoded.Report)
	assert.Equal(t, 1, len(decoded.Report.Issues))
}

// TestBulkResult_JSONOmitEmpty tests that empty fields are omitted in JSON
func TestBulkResult_JSONOmitEmpty(t *testing.T) {
	// Result with error but no report
	result1 := BulkResult{
		Name:  "test",
		Type:  ItemTypeSpec,
		Valid: false,
		Error: "test error",
	}

	data1, err := json.Marshal(result1)
	assert.NoError(t, err)
	// Report should not be in JSON
	assert.NotContains(t, string(data1), "report")
	assert.Contains(
		t,
		string(data1),
		"test error",
	)

	// Result with report but no error
	result2 := BulkResult{
		Name:  "test",
		Type:  ItemTypeSpec,
		Valid: true,
		Report: &ValidationReport{
			Valid: true,
		},
	}

	data2, err := json.Marshal(result2)
	assert.NoError(t, err)
	// Error should not be in JSON (omitempty)
	assert.NotContains(
		t,
		string(data2),
		"\"error\"",
	)
	assert.Contains(t, string(data2), "report")
}

// TestPrintHumanReport_SpecialCharacters tests handling of special characters
func TestPrintHumanReport_SpecialCharacters(
	t *testing.T,
) {
	report := &ValidationReport{
		Valid: false,
		Issues: []ValidationIssue{
			{
				Level:   "error",
				Path:    "spec.md",
				Message: "Message with \"quotes\" and 'apostrophes'",
			},
		},
	}

	output := captureOutput(func() {
		PrintHumanReport("test-spec", report)
	})

	assert.Contains(t, output, "quotes")
	assert.Contains(t, output, "apostrophes")
}

// TestPrintBulkHumanResults_LargeNumberOfResults tests with many results
func TestPrintBulkHumanResults_LargeNumberOfResults(
	t *testing.T,
) {
	results := make([]BulkResult, 100)
	for i := range 100 {
		results[i] = BulkResult{
			Name: "spec-" + string(
				rune('a'+i%26),
			),
			Type:  ItemTypeSpec,
			Valid: i%2 == 0, // Half valid, half invalid
		}
		if !results[i].Valid {
			results[i].Report = &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   "error",
						Path:    "spec.md",
						Message: "Error",
					},
				},
			}
		}
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	// Should include error/warning breakdown when there are failures
	assert.Contains(
		t,
		output,
		"50 passed, 50 failed",
	)
	assert.Contains(t, output, "100 total")
}

// TestToRelativePath tests the path conversion helper
func TestToRelativePath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute path with spectr prefix",
			input:    "/home/user/project/spectr/changes/foo/spec.md",
			expected: "changes/foo/spec.md",
		},
		{
			name:     "absolute path with nested spectr",
			input:    "/home/user/spectr-project/spectr/specs/auth/spec.md",
			expected: "specs/auth/spec.md",
		},
		{
			name:     "path without spectr prefix",
			input:    "/home/user/project/other/file.md",
			expected: "file.md",
		},
		{
			name:     "simple file name",
			input:    "spec.md",
			expected: "spec.md",
		},
		{
			name:     "empty path",
			input:    "",
			expected: ".",
		},
		{
			name:     "path ending with spectr",
			input:    "/home/user/spectr/foo.md",
			expected: "foo.md",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := ToRelativePath(tc.input)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// TestPrintBulkHumanResults_BlankLinesBetweenFailedItems tests visual separation
func TestPrintBulkHumanResults_BlankLinesBetweenFailedItems(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "fail1",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   LevelError,
						Path:    "spec.md",
						Message: "Error 1",
					},
				},
			},
		},
		{
			Name:  "fail2",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   LevelError,
						Path:    "spec.md",
						Message: "Error 2",
					},
				},
			},
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	// Should have blank line between failed items (but not before first)
	lines := strings.Split(output, "\n")
	foundBlankBetweenFails := false
	inFirstFail := false
	for _, line := range lines {
		if strings.Contains(line, "✗ fail1") {
			inFirstFail = true
		}
		if inFirstFail && line == "" {
			// Check if next non-empty contains fail2
			foundBlankBetweenFails = true
		}
		if foundBlankBetweenFails &&
			strings.Contains(line, "✗ fail2") {
			break
		}
	}
	assert.True(
		t,
		foundBlankBetweenFails,
		"Expected blank line between failed items",
	)
}

// TestPrintBulkHumanResults_GroupedIssuesByFile tests file grouping
func TestPrintBulkHumanResults_GroupedIssuesByFile(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "multi-file-issues",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   LevelError,
						Path:    "/home/user/spectr/changes/foo/spec.md",
						Message: "Error 1",
					},
					{
						Level:   LevelError,
						Path:    "/home/user/spectr/changes/foo/spec.md",
						Message: "Error 2",
					},
					{
						Level:   LevelWarning,
						Path:    "/home/user/spectr/changes/foo/tasks.md",
						Message: "Warning 1",
					},
				},
			},
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	// Should use relative paths
	assert.Contains(
		t,
		output,
		"changes/foo/spec.md",
	)
	assert.Contains(
		t,
		output,
		"changes/foo/tasks.md",
	)
	// Should not contain full absolute paths
	assert.NotContains(
		t,
		output,
		"/home/user/spectr/",
	)
}

// TestPrintBulkHumanResults_EnhancedSummaryWithBreakdown tests error/warning counts
func TestPrintBulkHumanResults_EnhancedSummaryWithBreakdown(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "valid-spec",
			Type:  ItemTypeSpec,
			Valid: true,
		},
		{
			Name:  "mixed-issues",
			Type:  ItemTypeSpec,
			Valid: false,
			Report: &ValidationReport{
				Valid: false,
				Issues: []ValidationIssue{
					{
						Level:   LevelError,
						Path:    "spec.md",
						Message: "Error 1",
					},
					{
						Level:   LevelError,
						Path:    "spec.md",
						Message: "Error 2",
					},
					{
						Level:   LevelWarning,
						Path:    "spec.md",
						Message: "Warning 1",
					},
				},
			},
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	// Should show error/warning breakdown when failures exist
	assert.Contains(
		t,
		output,
		"1 passed, 1 failed",
	)
	assert.Contains(t, output, "2 errors")
	assert.Contains(t, output, "1 warnings")
	assert.Contains(t, output, "2 total")
}

// TestPrintBulkHumanResults_NoBreakdownWhenAllPass tests summary without breakdown
func TestPrintBulkHumanResults_NoBreakdownWhenAllPass(
	t *testing.T,
) {
	results := []BulkResult{
		{
			Name:  "spec1",
			Type:  ItemTypeSpec,
			Valid: true,
		},
		{
			Name:  "spec2",
			Type:  ItemTypeSpec,
			Valid: true,
		},
	}

	output := captureOutput(func() {
		PrintBulkHumanResults(results)
	})

	// Should not include error/warning breakdown when all pass
	assert.Contains(
		t,
		output,
		"2 passed, 0 failed, 2 total",
	)
	assert.NotContains(t, output, "errors")
	assert.NotContains(t, output, "warnings")
}

// TestFormatLevel tests the formatLevel helper (without TTY)
func TestFormatLevel(t *testing.T) {
	// Note: In test environment, isTTY() returns false, so no colors
	errorLabel := formatLevel(LevelError)
	assert.Contains(t, errorLabel, "[ERROR]")

	warningLabel := formatLevel(LevelWarning)
	assert.Contains(t, warningLabel, "[WARNING]")

	infoLabel := formatLevel(LevelInfo)
	assert.Contains(t, infoLabel, "[INFO]")
}
