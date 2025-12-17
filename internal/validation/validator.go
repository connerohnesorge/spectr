package validation

import "path/filepath"

// Validator is the main orchestrator for validation operations.
// It coordinates validation of specs and changes using the underlying
// rule functions.
type Validator struct{}

// NewValidator creates a new Validator.
// Validation always treats warnings as errors (strict mode).
func NewValidator() *Validator {
	return &Validator{}
}

// ValidateSpec validates a specification file at the given path.
// This is a wrapper around ValidateSpecFile that always validates strictly.
// Returns a ValidationReport with all issues found, or an error for
// filesystem issues.
func (*Validator) ValidateSpec(
	path string,
) (*ValidationReport, error) {
	// Delegate to the spec validation rule function
	return ValidateSpecFile(path)
}

// ValidateChange validates all delta spec files in a change directory.
// This is a wrapper around ValidateChangeDeltaSpecs that always validates
// strictly. changeDir should be the path to a change directory
// (e.g., spectr/changes/add-feature).
// Returns a ValidationReport with all issues found, or an error for
// filesystem issues.
func (*Validator) ValidateChange(
	changeDir string,
) (*ValidationReport, error) {
	// Derive spectrRoot from changeDir
	// changeDir format: /path/to/project/spectr/changes/<change-id>
	// spectrRoot should be: /path/to/project/spectr
	spectrRoot := filepath.Dir(
		filepath.Dir(changeDir),
	)

	// Delegate to the change validation rule function
	return ValidateChangeDeltaSpecs(
		changeDir,
		spectrRoot,
	)
}

// CreateReport creates a ValidationReport from a list of issues.
// This is a helper method for creating validation reports.
// Warnings are converted to errors by the underlying validation functions
// before reaching this point (validation is always strict).
func (*Validator) CreateReport(
	issues []ValidationIssue,
) *ValidationReport {
	// Use the standard report creation function
	// Note: Warning-to-error conversion happens in the underlying validation
	// functions, not here, to ensure consistency across all validation paths
	return NewValidationReport(issues)
}
