package validation

import (
	"strings"
	"testing"
)

// TestValidateChangeDeltaSpecs_AddedRequirementAlreadyExistsInBase tests that
// adding a requirement that already exists in the base spec fails validation
func TestValidateChangeDeltaSpecs_AddedRequirementAlreadyExistsInBase(
	t *testing.T,
) {
	changeDir, spectrRoot := createChangeDir(
		t,
		map[string]string{
			"auth/spec.md": `## ADDED Requirements

### Requirement: User Authentication
The system SHALL provide user authentication.

#### Scenario: Login
- **WHEN** user logs in
- **THEN** user is authenticated
`,
		},
	)

	// Create base spec with the same requirement
	createBaseSpec(
		t,
		spectrRoot,
		"auth",
		`## Requirements

### Requirement: User Authentication
The system SHALL provide existing authentication.

#### Scenario: Existing login
- **WHEN** user logs in
- **THEN** user is authenticated
`,
	)

	report, err := ValidateChangeDeltaSpecs(
		changeDir,
		spectrRoot,
	)
	if err != nil {
		t.Fatalf(
			"ValidateChangeDeltaSpecs returned error: %v",
			err,
		)
	}

	if report.Valid {
		t.Error(
			"Expected invalid report - ADDED requirement already exists in base spec",
		)
	}

	// Look for the specific error message
	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"already exists in base spec",
			) {
			found = true

			break
		}
	}

	if found {
		return
	}

	t.Error(
		"Expected error about requirement already existing in base spec",
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

// TestValidateChangeDeltaSpecs_ModifiedRequirementDoesNotExistInBase tests that
// modifying a requirement that doesn't exist in the base spec fails validation
func TestValidateChangeDeltaSpecs_ModifiedRequirementDoesNotExistInBase(
	t *testing.T,
) {
	changeDir, spectrRoot := createChangeDir(
		t,
		map[string]string{
			"auth/spec.md": `## MODIFIED Requirements

### Requirement: User Authentication
The system SHALL provide enhanced authentication.

#### Scenario: Enhanced login
- **WHEN** user logs in
- **THEN** user gets enhanced auth
`,
		},
	)

	// Create base spec WITHOUT the requirement
	createBaseSpec(
		t,
		spectrRoot,
		"auth",
		`## Requirements

### Requirement: Other Feature
The system SHALL provide other feature.

#### Scenario: Other
- **WHEN** other feature used
- **THEN** it works
`,
	)

	report, err := ValidateChangeDeltaSpecs(
		changeDir,
		spectrRoot,
	)
	if err != nil {
		t.Fatalf(
			"ValidateChangeDeltaSpecs returned error: %v",
			err,
		)
	}

	if report.Valid {
		t.Error(
			"Expected invalid report - MODIFIED requirement does not exist in base spec",
		)
	}

	// Look for the specific error message
	found := false
	for _, issue := range report.Issues {
		if issue.Level == LevelError &&
			strings.Contains(
				issue.Message,
				"does not exist in base spec",
			) {
			found = true

			break
		}
	}

	if found {
		return
	}

	t.Error(
		"Expected error about requirement not existing in base spec",
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

// TestValidateChangeDeltaSpecs_ValidBaseSpecValidation tests that
// valid ADDED and MODIFIED operations pass validation
func TestValidateChangeDeltaSpecs_ValidBaseSpecValidation(
	t *testing.T,
) {
	changeDir, spectrRoot := createChangeDir(
		t,
		map[string]string{
			"auth/spec.md": `## ADDED Requirements

### Requirement: New Feature
The system SHALL provide new feature.

#### Scenario: New
- **WHEN** new feature used
- **THEN** it works

## MODIFIED Requirements

### Requirement: Existing Feature
The system SHALL provide enhanced existing feature.

#### Scenario: Enhanced
- **WHEN** existing feature used
- **THEN** it works better
`,
		},
	)

	// Create base spec with only Existing Feature (not New Feature)
	createBaseSpec(
		t,
		spectrRoot,
		"auth",
		`## Requirements

### Requirement: Existing Feature
The system SHALL provide existing feature.

#### Scenario: Basic
- **WHEN** existing feature used
- **THEN** it works
`,
	)

	report, err := ValidateChangeDeltaSpecs(
		changeDir,
		spectrRoot,
	)
	if err != nil {
		t.Fatalf(
			"ValidateChangeDeltaSpecs returned error: %v",
			err,
		)
	}

	if report.Valid {
		return
	}

	t.Error(
		"Expected valid report - ADDED requirement is new, MODIFIED requirement exists",
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
