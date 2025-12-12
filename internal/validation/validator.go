package validation

import "path/filepath"

// Validator is the main orchestrator for validation operations.
// It coordinates validation of specs and changes using the underlying
// rule functions.
type Validator struct {
	strictMode bool
}

// NewValidator creates a new Validator with the specified strict mode.
// When strictMode is true, warnings are treated as errors in the
// validation reports.
func NewValidator(strictMode bool) *Validator {
	return &Validator{
		strictMode: strictMode,
	}
}

// ValidateSpec validates a specification file at the given path.
// This is a wrapper around ValidateSpecFile that applies the
// validator's strictMode setting.
// Returns a ValidationReport with all issues found, or an error for
// filesystem issues.
func (v *Validator) ValidateSpec(
	path string,
) (*ValidationReport, error) {
	// Delegate to the spec validation rule function
	return ValidateSpecFile(path, v.strictMode)
}

// ValidateChange validates all delta spec files in a change directory.
// This is a wrapper around ValidateChangeDeltaSpecs that applies the
// validator's strictMode setting.
// changeDir should be the path to a change directory
// (e.g., spectr/changes/add-feature).
// Returns a ValidationReport with all issues found, or an error for
// filesystem issues.
func (v *Validator) ValidateChange(
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
		v.strictMode,
	)
}

// CreateReport creates a ValidationReport from a list of issues.
// This is a helper method that applies strict mode logic to the report.
// When strictMode is enabled, warnings in the issues are already
// converted to errors by the underlying validation functions before
// reaching this point.
func (*Validator) CreateReport(
	issues []ValidationIssue,
) *ValidationReport {
	// Use the standard report creation function
	// Note: strictMode conversion happens in the underlying validation
	// functions, not here, to ensure consistency across all validation
	// paths
	return NewValidationReport(issues)
}
