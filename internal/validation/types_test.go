package validation

import (
	"encoding/json"
	"testing"
)

func TestValidationLevel_Constants(t *testing.T) {
	tests := []struct {
		name     string
		level    ValidationLevel
		expected string
	}{
		{"Error level", LevelError, "ERROR"},
		{"Warning level", LevelWarning, "WARNING"},
		{"Info level", LevelInfo, "INFO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.level) != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, string(tt.level))
			}
		})
	}
}

func TestValidationIssue_JSONMarshaling(t *testing.T) {
	issue := ValidationIssue{
		Level:   LevelError,
		Path:    "specs/auth/spec.md",
		Message: "Missing required section",
	}

	// Marshal to JSON
	data, err := json.Marshal(issue)
	if err != nil {
		t.Fatalf("Failed to marshal ValidationIssue: %v", err)
	}

	// Unmarshal back
	var unmarshaled ValidationIssue
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ValidationIssue: %v", err)
	}

	// Verify fields
	if unmarshaled.Level != issue.Level {
		t.Errorf("Level: expected %q, got %q", issue.Level, unmarshaled.Level)
	}
	if unmarshaled.Path != issue.Path {
		t.Errorf("Path: expected %q, got %q", issue.Path, unmarshaled.Path)
	}
	if unmarshaled.Message != issue.Message {
		t.Errorf("Message: expected %q, got %q", issue.Message, unmarshaled.Message)
	}
}

func TestValidationIssue_JSONStructure(t *testing.T) {
	issue := ValidationIssue{
		Level:   LevelWarning,
		Path:    "changes/add-feature/specs/auth/spec.md",
		Message: "Requirement should include SHALL or MUST",
	}

	data, err := json.Marshal(issue)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Check that JSON has expected fields
	expected := `{"level":"WARNING","path":"changes/add-feature/specs/auth/spec.md","message":"Requirement should include SHALL or MUST"}`
	if string(data) != expected {
		t.Errorf("JSON structure mismatch\nExpected: %s\nGot: %s", expected, string(data))
	}
}

func TestValidationSummary_JSONMarshaling(t *testing.T) {
	summary := ValidationSummary{
		Errors:   2,
		Warnings: 3,
		Info:     1,
	}

	// Marshal to JSON
	data, err := json.Marshal(summary)
	if err != nil {
		t.Fatalf("Failed to marshal ValidationSummary: %v", err)
	}

	// Unmarshal back
	var unmarshaled ValidationSummary
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ValidationSummary: %v", err)
	}

	// Verify fields
	if unmarshaled.Errors != summary.Errors {
		t.Errorf("Errors: expected %d, got %d", summary.Errors, unmarshaled.Errors)
	}
	if unmarshaled.Warnings != summary.Warnings {
		t.Errorf("Warnings: expected %d, got %d", summary.Warnings, unmarshaled.Warnings)
	}
	if unmarshaled.Info != summary.Info {
		t.Errorf("Info: expected %d, got %d", summary.Info, unmarshaled.Info)
	}
}

func TestValidationReport_JSONMarshaling(t *testing.T) {
	report := ValidationReport{
		Valid: false,
		Issues: []ValidationIssue{
			{
				Level:   LevelError,
				Path:    "specs/auth/spec.md",
				Message: "Missing Requirements section",
			},
			{
				Level:   LevelWarning,
				Path:    "specs/auth/spec.md",
				Message: "Requirement should have at least one scenario",
			},
		},
		Summary: ValidationSummary{
			Errors:   1,
			Warnings: 1,
			Info:     0,
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal ValidationReport: %v", err)
	}

	// Unmarshal back
	var unmarshaled ValidationReport
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal ValidationReport: %v", err)
	}

	// Verify fields
	if unmarshaled.Valid != report.Valid {
		t.Errorf("Valid: expected %v, got %v", report.Valid, unmarshaled.Valid)
	}
	if len(unmarshaled.Issues) != len(report.Issues) {
		t.Errorf("Issues count: expected %d, got %d", len(report.Issues), len(unmarshaled.Issues))
	}
	if unmarshaled.Summary.Errors != report.Summary.Errors {
		t.Errorf(
			"Summary.Errors: expected %d, got %d",
			report.Summary.Errors,
			unmarshaled.Summary.Errors,
		)
	}
	if unmarshaled.Summary.Warnings != report.Summary.Warnings {
		t.Errorf(
			"Summary.Warnings: expected %d, got %d",
			report.Summary.Warnings,
			unmarshaled.Summary.Warnings,
		)
	}
}

func TestValidationReport_CompleteJSONRoundTrip(t *testing.T) {
	original := ValidationReport{
		Valid: true,
		Issues: []ValidationIssue{
			{Level: LevelInfo, Path: "specs/test/spec.md", Message: "All good"},
		},
		Summary: ValidationSummary{
			Errors:   0,
			Warnings: 0,
			Info:     1,
		},
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal
	var restored ValidationReport
	err = json.Unmarshal(data, &restored)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Deep comparison
	if restored.Valid != original.Valid {
		t.Error("Valid field mismatch")
	}
	if len(restored.Issues) != len(original.Issues) {
		t.Fatalf(
			"Issues length mismatch: expected %d, got %d",
			len(original.Issues),
			len(restored.Issues),
		)
	}
	if restored.Issues[0].Level != original.Issues[0].Level {
		t.Error("Issue level mismatch")
	}
	if restored.Issues[0].Path != original.Issues[0].Path {
		t.Error("Issue path mismatch")
	}
	if restored.Issues[0].Message != original.Issues[0].Message {
		t.Error("Issue message mismatch")
	}
	if restored.Summary.Errors != original.Summary.Errors {
		t.Error("Summary errors mismatch")
	}
	if restored.Summary.Warnings != original.Summary.Warnings {
		t.Error("Summary warnings mismatch")
	}
	if restored.Summary.Info != original.Summary.Info {
		t.Error("Summary info mismatch")
	}
}

func TestNewValidationReport_EmptyIssues(t *testing.T) {
	report := NewValidationReport(make([]ValidationIssue, 0))

	if !report.Valid {
		t.Error("Expected Valid=true for no issues")
	}
	if len(report.Issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(report.Issues))
	}
	if report.Summary.Errors != 0 {
		t.Errorf("Expected 0 errors, got %d", report.Summary.Errors)
	}
	if report.Summary.Warnings != 0 {
		t.Errorf("Expected 0 warnings, got %d", report.Summary.Warnings)
	}
	if report.Summary.Info != 0 {
		t.Errorf("Expected 0 info, got %d", report.Summary.Info)
	}
}

func TestNewValidationReport_WithErrors(t *testing.T) {
	issues := []ValidationIssue{
		{Level: LevelError, Path: "test.md", Message: "Error 1"},
		{Level: LevelError, Path: "test.md", Message: "Error 2"},
		{Level: LevelWarning, Path: "test.md", Message: "Warning 1"},
	}

	report := NewValidationReport(issues)

	if report.Valid {
		t.Error("Expected Valid=false when errors are present")
	}
	if report.Summary.Errors != 2 {
		t.Errorf("Expected 2 errors, got %d", report.Summary.Errors)
	}
	if report.Summary.Warnings != 1 {
		t.Errorf("Expected 1 warning, got %d", report.Summary.Warnings)
	}
	if report.Summary.Info != 0 {
		t.Errorf("Expected 0 info, got %d", report.Summary.Info)
	}
}

func TestNewValidationReport_OnlyWarnings(t *testing.T) {
	issues := []ValidationIssue{
		{Level: LevelWarning, Path: "test.md", Message: "Warning 1"},
		{Level: LevelWarning, Path: "test.md", Message: "Warning 2"},
		{Level: LevelInfo, Path: "test.md", Message: "Info 1"},
	}

	report := NewValidationReport(issues)

	if !report.Valid {
		t.Error("Expected Valid=true when only warnings and info are present")
	}
	if report.Summary.Errors != 0 {
		t.Errorf("Expected 0 errors, got %d", report.Summary.Errors)
	}
	if report.Summary.Warnings != 2 {
		t.Errorf("Expected 2 warnings, got %d", report.Summary.Warnings)
	}
	if report.Summary.Info != 1 {
		t.Errorf("Expected 1 info, got %d", report.Summary.Info)
	}
}

func TestNewValidationReport_MixedIssues(t *testing.T) {
	issues := []ValidationIssue{
		{Level: LevelError, Path: "test.md", Message: "Error 1"},
		{Level: LevelWarning, Path: "test.md", Message: "Warning 1"},
		{Level: LevelWarning, Path: "test.md", Message: "Warning 2"},
		{Level: LevelInfo, Path: "test.md", Message: "Info 1"},
		{Level: LevelInfo, Path: "test.md", Message: "Info 2"},
		{Level: LevelInfo, Path: "test.md", Message: "Info 3"},
	}

	report := NewValidationReport(issues)

	if report.Valid {
		t.Error("Expected Valid=false when errors are present")
	}
	if report.Summary.Errors != 1 {
		t.Errorf("Expected 1 error, got %d", report.Summary.Errors)
	}
	if report.Summary.Warnings != 2 {
		t.Errorf("Expected 2 warnings, got %d", report.Summary.Warnings)
	}
	if report.Summary.Info != 3 {
		t.Errorf("Expected 3 info, got %d", report.Summary.Info)
	}
	if len(report.Issues) != 6 {
		t.Errorf("Expected 6 total issues, got %d", len(report.Issues))
	}
}

func TestValidationReport_NilIssues(t *testing.T) {
	report := NewValidationReport(nil)

	if !report.Valid {
		t.Error("Expected Valid=true for nil issues")
	}
	if report.Issues == nil {
		t.Error("Expected Issues to be non-nil slice")
	}
	if len(report.Issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(report.Issues))
	}
}

func TestValidationReport_JSONWithNilIssues(t *testing.T) {
	report := &ValidationReport{
		Valid:   true,
		Issues:  nil,
		Summary: ValidationSummary{},
	}

	data, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	var unmarshaled ValidationReport
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// JSON null array unmarshals to nil slice, which is fine
	if !unmarshaled.Valid {
		t.Error("Valid field should be true")
	}
}
