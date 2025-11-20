package validation

// ValidationLevel represents the severity of a validation issue
type ValidationLevel string

const (
	// LevelError indicates a critical validation failure
	LevelError ValidationLevel = "ERROR"
	// LevelWarning indicates a non-critical issue that should be addressed
	LevelWarning ValidationLevel = "WARNING"
	// LevelInfo provides informational feedback
	LevelInfo ValidationLevel = "INFO"
)

// ValidationIssue represents a single validation problem or note
type ValidationIssue struct {
	Level   ValidationLevel `json:"level"`
	Path    string          `json:"path"`
	Line    int             `json:"line,omitempty"`
	Message string          `json:"message"`
}

// ValidationSummary provides aggregate counts of validation issues
type ValidationSummary struct {
	Errors   int `json:"errors"`
	Warnings int `json:"warnings"`
	Info     int `json:"info"`
}

// ValidationReport contains the complete validation results for an item
type ValidationReport struct {
	Valid   bool              `json:"valid"`
	Issues  []ValidationIssue `json:"issues"`
	Summary ValidationSummary `json:"summary"`
}

// NewValidationReport creates a new ValidationReport from a list of issues
func NewValidationReport(issuesParam []ValidationIssue) *ValidationReport {
	// Initialize issues slice to empty slice if nil
	issues := issuesParam
	if issues == nil {
		issues = make([]ValidationIssue, 0)
	}

	summary := ValidationSummary{}
	for _, issue := range issues {
		switch issue.Level {
		case LevelError:
			summary.Errors++
		case LevelWarning:
			summary.Warnings++
		case LevelInfo:
			summary.Info++
		}
	}

	// Report is valid if there are no errors
	valid := summary.Errors == 0

	return &ValidationReport{
		Valid:   valid,
		Issues:  issues,
		Summary: summary,
	}
}
