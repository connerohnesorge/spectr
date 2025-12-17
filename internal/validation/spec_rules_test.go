package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateSpecFile_ValidSpec(
	t *testing.T,
) {
	// Create a valid spec file
	content := `# Test Specification

## Requirements

### Requirement: User Authentication
The system SHALL provide user authentication functionality.

#### Scenario: Successful login
- **WHEN** user provides valid credentials
- **THEN** user is authenticated and session is created

#### Scenario: Failed login
- **WHEN** user provides invalid credentials
- **THEN** authentication fails and error message is displayed
`

	// Write to temp file
	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	// Validate the spec
	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should be valid with no issues
	if !report.Valid {
		t.Errorf(
			"Expected valid report, got invalid with %d errors, %d warnings",
			report.Summary.Errors,
			report.Summary.Warnings,
		)
		for _, issue := range report.Issues {
			t.Logf(
				"  %s: %s - %s",
				issue.Level,
				issue.Path,
				issue.Message,
			)
		}
	}

	if len(report.Issues) != 0 {
		t.Errorf(
			"Expected 0 issues, got %d",
			len(report.Issues),
		)
	}
}

func TestValidateSpecFile_MissingRequirements(
	t *testing.T,
) {
	content := `# Test Specification

Some content here but no Requirements section.
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should be invalid due to missing Requirements
	if report.Valid {
		t.Error(
			"Expected invalid report due to missing Requirements section",
		)
	}

	if report.Summary.Errors != 1 {
		t.Errorf(
			"Expected 1 error, got %d",
			report.Summary.Errors,
		)
	}

	// Check that the error message is correct
	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"Requirements",
			) {
			found = true

			break
		}
	}
	if !found {
		t.Error(
			"Expected error about missing Requirements section",
		)
	}
}

func TestValidateSpecFile_RequirementWithoutShallOrMust(
	t *testing.T,
) {
	content := `# Test Specification

## Requirements

### Requirement: Some Feature
The system provides some feature.

#### Scenario: Test scenario
- **WHEN** something happens
- **THEN** something else happens
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should be invalid (warnings are now errors in always-strict mode)
	if report.Valid {
		t.Error(
			"Expected invalid report (warnings are errors in strict mode)",
		)
	}

	if report.Summary.Errors != 1 {
		t.Errorf(
			"Expected 1 error, got %d",
			report.Summary.Errors,
		)
	}

	// Check that the error message is correct
	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"SHALL or MUST",
			) {
			found = true

			break
		}
	}
	if !found {
		t.Error(
			"Expected error about missing SHALL or MUST",
		)
	}
}

func TestValidateSpecFile_RequirementWithoutScenarios(
	t *testing.T,
) {
	content := `# Test Specification

## Requirements

### Requirement: Some Feature
The system SHALL provide some feature.
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should be invalid (warnings are now errors in always-strict mode)
	if report.Valid {
		t.Error(
			"Expected invalid report (warnings are errors in strict mode)",
		)
	}

	if report.Summary.Errors != 1 {
		t.Errorf(
			"Expected 1 error, got %d",
			report.Summary.Errors,
		)
	}

	// Check that the error message is correct
	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"at least one scenario",
			) {
			found = true

			break
		}
	}
	if !found {
		t.Error(
			"Expected error about missing scenarios",
		)
	}
}

func TestValidateSpecFile_InvalidScenarioFormat(
	t *testing.T,
) {
	content := `# Test Specification

## Requirements

### Requirement: Some Feature
The system SHALL provide some feature.

##### Scenario: Test scenario
- **WHEN** something happens
- **THEN** something else happens
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should be invalid due to malformed scenario (##### instead of ####)
	if report.Valid {
		t.Error(
			"Expected invalid report due to malformed scenario format",
		)
	}

	if report.Summary.Errors == 0 {
		t.Errorf(
			"Expected at least 1 error, got %d",
			report.Summary.Errors,
		)
	}

	// Check that there's an error about scenario format
	foundFormatError := false
	foundMissingScenarioError := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"#### Scenario:",
			) {
			foundFormatError = true
		}
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"at least one scenario",
			) {
			foundMissingScenarioError = true
		}
	}

	if !foundFormatError {
		t.Error(
			"Expected error about incorrect scenario format",
		)
	}
	if !foundMissingScenarioError {
		t.Error(
			"Expected error about missing scenarios (since malformed ones don't count)",
		)
	}
}

func TestValidateSpecFile_AlwaysStrict(
	t *testing.T,
) {
	// This spec has what would have been a warning (missing SHALL/MUST)
	// but is now always treated as an error
	content := `# Test Specification

## Requirements

### Requirement: Some Feature
The system provides some feature.

#### Scenario: Test scenario
- **WHEN** something happens
- **THEN** something else happens
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	// Validate (always strict)
	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should be invalid (warnings are always converted to errors)
	if report.Valid {
		t.Error(
			"Expected invalid report (warnings become errors)",
		)
	}

	// Warnings are converted to errors
	// 1 error: missing SHALL/MUST in requirement
	if report.Summary.Errors != 1 {
		t.Errorf(
			"Expected 1 error, got %d",
			report.Summary.Errors,
		)
	}

	if report.Summary.Warnings != 0 {
		t.Errorf(
			"Expected 0 warnings (all converted to errors), got %d",
			report.Summary.Warnings,
		)
	}
}

func TestValidateSpecFile_FileNotFound(
	t *testing.T,
) {
	nonexistentPath := "/tmp/nonexistent-spec-file-12345.md"

	_, err := ValidateSpecFile(
		nonexistentPath,
	)
	if err == nil {
		t.Error(
			"Expected error for nonexistent file, got nil",
		)
	}

	if !strings.Contains(
		err.Error(),
		"failed to read spec file",
	) {
		t.Errorf(
			"Expected file read error, got: %v",
			err,
		)
	}
}

func TestValidateSpecFile_MultipleIssues(
	t *testing.T,
) {
	content := `# Test Specification

## Requirements

### Requirement: Feature One
The system provides feature one.

### Requirement: Feature Two
The system SHALL provide feature two.

#### Scenario: Test scenario
- **WHEN** something happens
- **THEN** something else happens

### Requirement: Feature Three
The system MUST provide feature three.
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should be invalid (warnings are now errors in always-strict mode)
	if report.Valid {
		t.Error(
			"Expected invalid report (warnings are errors)",
		)
	}

	// Expect:
	// - 3 ERRORs (warnings converted to errors):
	//   - Feature One (no SHALL/MUST)
	//   - Feature One (no scenarios)
	//   - Feature Three (no scenarios)
	// Feature Two has SHALL and has scenario, so no issues
	expectedErrors := 3
	expectedWarnings := 0

	if report.Summary.Errors != expectedErrors {
		t.Errorf(
			"Expected %d errors, got %d",
			expectedErrors,
			report.Summary.Errors,
		)
	}

	if report.Summary.Warnings != expectedWarnings {
		t.Errorf(
			"Expected %d warnings, got %d",
			expectedWarnings,
			report.Summary.Warnings,
		)
	}

	if len(
		report.Issues,
	) != expectedErrors+expectedWarnings {
		t.Errorf(
			"Expected %d total issues, got %d",
			expectedErrors+expectedWarnings,
			len(report.Issues),
		)
	}
}

func TestValidateSpecFile_BoldScenarioFormat(
	t *testing.T,
) {
	content := `# Test Specification

## Requirements

### Requirement: Some Feature
The system SHALL provide some feature.

**Scenario: Test scenario**
- **WHEN** something happens
- **THEN** something else happens
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should have error about malformed scenario
	if report.Valid {
		t.Error(
			"Expected invalid report due to malformed scenario format",
		)
	}

	foundFormatError := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"#### Scenario:",
			) {
			foundFormatError = true

			break
		}
	}

	if !foundFormatError {
		t.Error(
			"Expected error about incorrect scenario format (bold instead of header)",
		)
	}
}

func TestValidateSpecFile_BulletScenarioFormat(
	t *testing.T,
) {
	content := `# Test Specification

## Requirements

### Requirement: Some Feature
The system SHALL provide some feature.

- **Scenario: Test scenario**
  - **WHEN** something happens
  - **THEN** something else happens
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	report, err := ValidateSpecFile(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpecFile returned error: %v",
			err,
		)
	}

	// Should have error about malformed scenario
	if report.Valid {
		t.Error(
			"Expected invalid report due to malformed scenario format",
		)
	}

	foundFormatError := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"#### Scenario:",
			) {
			foundFormatError = true

			break
		}
	}

	if !foundFormatError {
		t.Error(
			"Expected error about incorrect scenario format (bullet instead of header)",
		)
	}
}
