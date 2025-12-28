package validation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewValidator(t *testing.T) {
	// NewValidator no longer takes a strictMode parameter
	// Validation is always strict (warnings are converted to errors)
	v := NewValidator()
	if v == nil {
		t.Fatal(
			"NewValidator returned nil",
		)
	}
}

func TestValidator_ValidateSpec_ValidSpec(
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

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	// Validation is always strict now
	v := NewValidator()
	report, err := v.ValidateSpec(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpec returned error: %v",
			err,
		)
	}

	if !report.Valid {
		t.Error(
			"Expected Valid=true, got false",
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

	if report.Summary.Errors > 0 {
		t.Errorf(
			"Expected no errors, got %d",
			report.Summary.Errors,
		)
	}
}

func TestValidator_ValidateSpec_InvalidSpec(
	t *testing.T,
) {
	// Create an invalid spec file (missing Requirements section)
	content := `# Test Specification

This is just some content without proper sections.
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	// Validation is always strict now
	v := NewValidator()
	report, err := v.ValidateSpec(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpec returned error: %v",
			err,
		)
	}

	if report.Valid {
		t.Error(
			"Expected Valid=false, got true",
		)
	}

	if report.Summary.Errors < 1 {
		t.Errorf(
			"Expected at least 1 error, got %d",
			report.Summary.Errors,
		)
	}
}

func TestValidator_ValidateSpec_WarningsConvertedToErrors(
	t *testing.T,
) {
	// Create a spec with what would be warnings (missing scenarios)
	// but are now errors because validation is always strict
	content := `# Test Specification

## Requirements

### Requirement: User Authentication
The system SHALL provide user authentication functionality.
`

	tmpDir := t.TempDir()
	specPath := filepath.Join(tmpDir, "spec.md")
	err := os.WriteFile(
		specPath,
		[]byte(content),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create test file: %v",
			err,
		)
	}

	// Validation is always strict - warnings are converted to errors
	v := NewValidator()
	report, err := v.ValidateSpec(
		specPath,
	)
	if err != nil {
		t.Fatalf(
			"ValidateSpec returned error: %v",
			err,
		)
	}

	// Should be invalid (warnings become errors)
	if report.Valid {
		t.Error(
			"Expected Valid=false (warnings converted to errors), got true",
		)
	}

	// Warnings are converted to errors
	if report.Summary.Errors == 0 {
		t.Error(
			"Expected errors (converted from warnings), got none",
		)
	}

	if report.Summary.Warnings > 0 {
		t.Errorf(
			"Expected no warnings (should be converted to errors), got %d",
			report.Summary.Warnings,
		)
	}
}

func TestValidator_ValidateChange_ValidChange(
	t *testing.T,
) {
	// Create a valid change with delta specs
	changeDir := setupValidChangeDirectory(t)

	// Validation is always strict now
	v := NewValidator()
	report, err := v.ValidateChange(changeDir)
	if err != nil {
		t.Fatalf(
			"ValidateChange returned error: %v",
			err,
		)
	}

	if !report.Valid {
		t.Error(
			"Expected Valid=true, got false",
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

	if report.Summary.Errors > 0 {
		t.Errorf(
			"Expected no errors, got %d",
			report.Summary.Errors,
		)
	}
}

// setupValidChangeDirectory creates a valid change directory for testing
func setupValidChangeDirectory(
	t *testing.T,
) string {
	t.Helper()

	tmpDir := t.TempDir()
	changeDir := filepath.Join(
		tmpDir,
		"test-change",
	)
	specsDir := filepath.Join(
		changeDir,
		"specs",
		"auth",
	)
	err := os.MkdirAll(specsDir, 0o755)
	if err != nil {
		t.Fatalf(
			"Failed to create test directories: %v",
			err,
		)
	}

	deltaContent := `## ADDED Requirements

### Requirement: Two-Factor Authentication
The system SHALL require two-factor authentication for all users.

#### Scenario: OTP required
- **WHEN** user provides valid credentials
- **THEN** system prompts for OTP code

#### Scenario: OTP validation
- **WHEN** user provides valid OTP code
- **THEN** user is authenticated successfully
`

	specPath := filepath.Join(specsDir, "spec.md")
	err = os.WriteFile(
		specPath,
		[]byte(deltaContent),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create delta spec file: %v",
			err,
		)
	}

	return changeDir
}

func TestValidator_ValidateChange_InvalidChange(
	t *testing.T,
) {
	// Create an invalid change (missing deltas)
	tmpDir := t.TempDir()
	changeDir := filepath.Join(
		tmpDir,
		"test-change",
	)
	specsDir := filepath.Join(
		changeDir,
		"specs",
		"auth",
	)
	err := os.MkdirAll(specsDir, 0o755)
	if err != nil {
		t.Fatalf(
			"Failed to create test directories: %v",
			err,
		)
	}

	// Create spec.md with no delta sections
	deltaContent := `# Just a spec with no delta sections`
	specPath := filepath.Join(specsDir, "spec.md")
	err = os.WriteFile(
		specPath,
		[]byte(deltaContent),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create delta spec file: %v",
			err,
		)
	}

	// Validation is always strict now
	v := NewValidator()
	report, err := v.ValidateChange(
		changeDir,
	)
	if err != nil {
		t.Fatalf(
			"ValidateChange returned error: %v",
			err,
		)
	}

	if report.Valid {
		t.Error(
			"Expected Valid=false, got true",
		)
	}

	if report.Summary.Errors < 1 {
		t.Errorf(
			"Expected at least 1 error, got %d",
			report.Summary.Errors,
		)
	}
}

func TestValidator_ValidateChange_MissingSpecsDir(
	t *testing.T,
) {
	tmpDir := t.TempDir()
	changeDir := filepath.Join(
		tmpDir,
		"test-change",
	)
	// Don't create specs directory

	v := NewValidator()
	_, err := v.ValidateChange(changeDir)

	if err == nil {
		t.Fatal(
			"Expected error for missing specs directory, got nil",
		)
	}
}

func TestValidator_ValidateSpec_NonexistentFile(
	t *testing.T,
) {
	v := NewValidator()
	_, err := v.ValidateSpec(
		"/nonexistent/path/spec.md",
	)

	if err == nil {
		t.Fatal(
			"Expected error for nonexistent file, got nil",
		)
	}
}

func TestValidator_CreateReport(t *testing.T) {
	// CreateReport is a helper that just creates a report from issues
	// Warning-to-error conversion happens in the validation functions,
	// not in CreateReport
	tests := []struct {
		name       string
		issues     []ValidationIssue
		wantValid  bool
		wantErrors int
	}{
		{
			name: "no issues",
			issues: make(
				[]ValidationIssue,
				0,
			),
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "errors present",
			issues: []ValidationIssue{
				{
					Level:   LevelError,
					Path:    "test.md",
					Message: "Error 1",
				},
				{
					Level:   LevelError,
					Path:    "test.md",
					Message: "Error 2",
				},
			},
			wantValid:  false,
			wantErrors: 2,
		},
		{
			name: "warnings (not converted by CreateReport)",
			issues: []ValidationIssue{
				{
					Level:   LevelWarning,
					Path:    "test.md",
					Message: "Warning 1",
				},
			},
			wantValid:  true,
			wantErrors: 0,
		},
		{
			name: "mixed issues",
			issues: []ValidationIssue{
				{
					Level:   LevelError,
					Path:    "test.md",
					Message: "Error",
				},
				{
					Level:   LevelWarning,
					Path:    "test.md",
					Message: "Warning",
				},
				{
					Level:   LevelInfo,
					Path:    "test.md",
					Message: "Info",
				},
			},
			wantValid:  false,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewValidator()
			report := v.CreateReport(tt.issues)

			if report == nil {
				t.Fatal(
					"CreateReport returned nil",
				)
			}

			if report.Valid != tt.wantValid {
				t.Errorf(
					"Expected Valid=%v, got %v",
					tt.wantValid,
					report.Valid,
				)
			}

			if report.Summary.Errors != tt.wantErrors {
				t.Errorf(
					"Expected %d errors, got %d",
					tt.wantErrors,
					report.Summary.Errors,
				)
			}

			if len(
				report.Issues,
			) != len(
				tt.issues,
			) {
				t.Errorf(
					"Expected %d issues, got %d",
					len(tt.issues),
					len(report.Issues),
				)
			}
		})
	}
}

func TestValidator_IntegrationWithMultipleCapabilities(
	t *testing.T,
) {
	// Test a change that affects multiple capabilities
	// Create proper directory structure: spectr/changes/<change-id>/
	projectRoot := t.TempDir()
	spectrRoot := filepath.Join(
		projectRoot,
		"spectr",
	)
	changesRoot := filepath.Join(
		spectrRoot,
		"changes",
	)
	specsRoot := filepath.Join(
		spectrRoot,
		"specs",
	)
	changeDir := filepath.Join(
		changesRoot,
		"multi-capability-change",
	)

	// Create auth capability delta
	authDir := filepath.Join(
		changeDir,
		"specs",
		"auth",
	)
	err := os.MkdirAll(authDir, 0o755)
	if err != nil {
		t.Fatalf(
			"Failed to create auth directory: %v",
			err,
		)
	}

	authContent := `## ADDED Requirements

### Requirement: OAuth Support
The system SHALL support OAuth 2.0 authentication.

#### Scenario: OAuth login
- **WHEN** user clicks OAuth login
- **THEN** user is redirected to OAuth provider
`

	err = os.WriteFile(
		filepath.Join(authDir, "spec.md"),
		[]byte(authContent),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create auth spec: %v",
			err,
		)
	}

	// Create notifications capability delta
	notifDir := filepath.Join(
		changeDir,
		"specs",
		"notifications",
	)
	err = os.MkdirAll(notifDir, 0o755)
	if err != nil {
		t.Fatalf(
			"Failed to create notifications directory: %v",
			err,
		)
	}

	notifContent := `## MODIFIED Requirements

### Requirement: Email Notifications
The system SHALL send email notifications for authentication events.

#### Scenario: Login notification
- **WHEN** user logs in successfully
- **THEN** email notification is sent
`

	err = os.WriteFile(
		filepath.Join(notifDir, "spec.md"),
		[]byte(notifContent),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create notifications spec: %v",
			err,
		)
	}

	// Create base spec for notifications (MODIFIED requires existing spec)
	notifBaseDir := filepath.Join(
		specsRoot,
		"notifications",
	)
	err = os.MkdirAll(notifBaseDir, 0o755)
	if err != nil {
		t.Fatalf(
			"Failed to create notifications base directory: %v",
			err,
		)
	}

	notifBaseContent := `## Requirements

### Requirement: Email Notifications
The system SHALL send email notifications.

#### Scenario: Basic notification
- **WHEN** event occurs
- **THEN** email is sent
`

	err = os.WriteFile(
		filepath.Join(notifBaseDir, "spec.md"),
		[]byte(notifBaseContent),
		0o644,
	)
	if err != nil {
		t.Fatalf(
			"Failed to create notifications base spec: %v",
			err,
		)
	}

	v := NewValidator()
	report, err := v.ValidateChange(changeDir)
	if err != nil {
		t.Fatalf(
			"ValidateChange returned error: %v",
			err,
		)
	}

	if report.Valid {
		return
	}
	t.Error(
		"Expected valid report for multi-capability change",
	)
	for _, issue := range report.Issues {
		t.Logf(
			"  %s: %s - %s",
			issue.Level,
			issue.Path,
			issue.Message,
		)
	}

	if report.Summary.Errors > 0 {
		t.Errorf(
			"Expected no errors, got %d",
			report.Summary.Errors,
		)
	}
}
