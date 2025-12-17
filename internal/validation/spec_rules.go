package validation

import (
	"fmt"
	"os"
	"strings"
)

// ValidateSpecFile validates a spec file according to Spectr rules
// Returns a ValidationReport containing all issues found, or an error
// for filesystem issues
// Note: Always applies strict validation (warnings are converted to errors)
func ValidateSpecFile(
	path string,
) (*ValidationReport, error) {
	// Read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to read spec file: %w",
			err,
		)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// Parse sections
	sections := ExtractSections(contentStr)
	issues := make([]ValidationIssue, 0)

	// Rule 1: Check for ## Requirements section (ERROR if missing)
	requirementsContent, hasRequirements := sections["Requirements"]
	if !hasRequirements {
		issues = append(issues, ValidationIssue{
			Level:   LevelError,
			Path:    path,
			Line:    1, // Missing section defaults to line 1
			Message: "Missing required '## Requirements' section",
		})
	}

	// Rule 2-5: Validate requirements (only if Requirements section exists)
	if hasRequirements {
		reqIssues := validateRequirements(
			path,
			requirementsContent,
			lines,
		)
		issues = append(issues, reqIssues...)
	}

	// Always convert warnings to errors (strict validation)
	convertWarningsToErrors(issues)

	// Create and return the validation report
	return NewValidationReport(issues), nil
}

// validateRequirements validates all requirements in a spec file
// Returns a slice of validation issues found
func validateRequirements(
	path, requirementsContent string,
	lines []string,
) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	requirements := ExtractRequirements(requirementsContent)
	requirementsLine := findSectionLine(lines, "Requirements")

	for _, req := range requirements {
		reqIssues := validateSingleRequirement(
			path,
			req,
			lines,
			requirementsLine,
		)
		issues = append(issues, reqIssues...)
	}

	return issues
}

// validateSingleRequirement validates a single requirement
// Returns a slice of validation issues found
func validateSingleRequirement(
	path string,
	req Requirement,
	lines []string,
	requirementsLine int,
) []ValidationIssue {
	issues := make([]ValidationIssue, 0)
	reqPath := fmt.Sprintf("%s: Requirement '%s'", path, req.Name)
	reqLine := findRequirementLine(lines, req.Name, requirementsLine)

	// Rule 2: Check for SHALL or MUST (WARNING if missing)
	if !ContainsShallOrMust(req.Content) {
		issues = append(issues, ValidationIssue{
			Level: LevelWarning,
			Path:  reqPath,
			Line:  reqLine,
			Message: "Requirement should contain SHALL or " +
				"MUST to indicate normative requirement",
		})
	}

	// Rule 3: Check for at least one scenario (WARNING)
	if len(req.Scenarios) == 0 {
		issues = append(issues, ValidationIssue{
			Level:   LevelWarning,
			Path:    reqPath,
			Line:    reqLine,
			Message: "Requirement should have at least one scenario",
		})
	}

	// Rule 4: Check scenario format (ERROR if wrong format)
	if len(req.Scenarios) == 0 && hasMalformedScenarios(req.Content) {
		malformedLine := findMalformedScenarioLine(lines, reqLine)
		issues = append(issues, ValidationIssue{
			Level: LevelError,
			Path:  reqPath,
			Line:  malformedLine,
			Message: "Scenarios must use '#### Scenario:' " +
				"format (4 hashtags followed by 'Scenario:')",
		})
	}

	return issues
}

// convertWarningsToErrors converts all warnings to errors in-place
func convertWarningsToErrors(issues []ValidationIssue) {
	for i := range issues {
		if issues[i].Level == LevelWarning {
			issues[i].Level = LevelError
		}
	}
}

// hasMalformedScenarios detects if content has scenario-like text that
// doesn't match proper format
func hasMalformedScenarios(content string) bool {
	// Look for common malformations:
	// - "**Scenario:" (bold instead of header)
	// - "### Scenario:" (3 hashtags instead of 4)
	// - "##### Scenario:" (5+ hashtags)
	// - "###### Scenario:" (6 hashtags)
	// - "- **Scenario:" (bullet point)
	// - "Scenario:" at start of line without hashtags

	// Simple heuristic: if content contains "Scenario:" but
	// ExtractScenarios found none, and it's not just in regular prose
	// (would need more context to be certain)
	// For now, we'll check for common markdown scenario patterns
	// that are wrong

	// Check for ### Scenario: (3 hashtags - wrong)
	if containsPattern(content, "### Scenario:") {
		return true
	}

	// Check for ##### Scenario: (5 hashtags - wrong)
	if containsPattern(
		content,
		"##### Scenario:",
	) {
		return true
	}

	// Check for ###### Scenario: (6 hashtags - wrong)
	if containsPattern(
		content,
		"###### Scenario:",
	) {
		return true
	}

	// Check for **Scenario: (bold - wrong)
	if containsPattern(content, "**Scenario:") {
		return true
	}

	// Check for - **Scenario: (bullet + bold - wrong)
	if containsPattern(content, "- **Scenario:") {
		return true
	}

	return false
}

// containsPattern checks if content contains the given pattern
func containsPattern(
	content, pattern string,
) bool {
	return len(content) > 0 && len(pattern) > 0 &&
		strings.Contains(content, pattern)
}

// findSectionLine finds the line number where a section header appears
// Returns 1 if not found
func findSectionLine(
	lines []string,
	sectionName string,
) int {
	sectionHeader := "## " + sectionName
	for i, line := range lines {
		if strings.HasPrefix(
			strings.TrimSpace(line),
			sectionHeader,
		) {
			return i + 1 // Line numbers are 1-indexed
		}
	}

	return 1 // Default to line 1 if not found
}

// findRequirementLine finds the line number where a requirement appears
// Searches from startLine onwards
// Returns startLine if not found
func findRequirementLine(
	lines []string,
	reqName string,
	startLine int,
) int {
	reqHeader := "### Requirement: " + reqName
	// Start searching from startLine (convert to 0-indexed)
	searchStart := startLine - 1
	if searchStart < 0 {
		searchStart = 0
	}

	for i := searchStart; i < len(lines); i++ {
		if strings.Contains(lines[i], reqHeader) {
			return i + 1 // Line numbers are 1-indexed
		}
	}

	return startLine // Default to section start if not found
}

// findMalformedScenarioLine finds the line number of a malformed scenario
// Searches from reqLine onwards
// Returns reqLine if not found
func findMalformedScenarioLine(
	lines []string,
	reqLine int,
) int {
	// Start searching from reqLine (convert to 0-indexed)
	searchStart := reqLine - 1
	if searchStart < 0 {
		searchStart = 0
	}

	// Look for common malformations
	malformedPatterns := []string{
		"### Scenario:",
		"##### Scenario:",
		"###### Scenario:",
		"**Scenario:",
		"- **Scenario:",
	}

	for i := searchStart; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		for _, pattern := range malformedPatterns {
			if strings.Contains(line, pattern) {
				return i + 1 // Line numbers are 1-indexed
			}
		}
	}

	return reqLine // Default to requirement line if not found
}
