package validation

import (
	"fmt"
	"os"
	"strings"
)

// ValidateSpecFile validates a spec file according to Spectr rules
// Returns a ValidationReport containing all issues found, or an error
// for filesystem issues
//
//nolint:revive // strictMode is intentional control flag
func ValidateSpecFile(path string, strictMode bool) (*ValidationReport, error) {
	// Read the file
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file: %w", err)
	}

	// Parse sections
	sections := ExtractSections(string(content))
	issues := make([]ValidationIssue, 0)

	// Rule 1: Check for ## Purpose section (ERROR if missing)
	purposeContent, hasPurpose := sections["Purpose"]
	if !hasPurpose {
		issues = append(issues, ValidationIssue{
			Level:   LevelError,
			Path:    path,
			Message: "Missing required '## Purpose' section",
		})
	}

	// Rule 2: Check for ## Requirements section (ERROR if missing)
	requirementsContent, hasRequirements := sections["Requirements"]
	if !hasRequirements {
		issues = append(issues, ValidationIssue{
			Level:   LevelError,
			Path:    path,
			Message: "Missing required '## Requirements' section",
		})
	}

	// Rule 3: Check Purpose section length (WARNING if < 50 chars)
	if hasPurpose && len(purposeContent) < 50 {
		issues = append(issues, ValidationIssue{
			Level: LevelWarning,
			Path:  path,
			Message: fmt.Sprintf(
				"Purpose section is too short "+
					"(%d characters, minimum 50 recommended)",
				len(purposeContent),
			),
		})
	}

	// Rule 4-7: Validate requirements (only if Requirements section exists)
	if hasRequirements {
		requirements := ExtractRequirements(requirementsContent)

		for _, req := range requirements {
			reqPath := fmt.Sprintf("%s: Requirement '%s'", path, req.Name)

			// Rule 4: Check for SHALL or MUST (WARNING if missing)
			if !ContainsShallOrMust(req.Content) {
				issues = append(issues, ValidationIssue{
					Level: LevelWarning,
					Path:  reqPath,
					Message: "Requirement should contain SHALL or " +
						"MUST to indicate normative requirement",
				})
			}

			// Rule 5: Check for at least one scenario (WARNING)
			if len(req.Scenarios) == 0 {
				issues = append(issues, ValidationIssue{
					Level: LevelWarning,
					Path:  reqPath,
					Message: "Requirement should have " +
						"at least one scenario",
				})
			}

			// Rule 6: Check scenario format (ERROR if wrong format)
			// This is implicitly handled by ExtractScenarios - if
			// there's content that looks like scenarios but doesn't
			// match #### Scenario: format, they won't be extracted
			// We need to explicitly check for malformed scenarios
			if len(req.Scenarios) == 0 &&
				hasMalformedScenarios(req.Content) {
				issues = append(issues, ValidationIssue{
					Level: LevelError,
					Path:  reqPath,
					Message: "Scenarios must use '#### Scenario:' " +
						"format (4 hashtags followed by 'Scenario:')",
				})
			}
		}
	}

	// Apply strict mode: convert warnings to errors
	if strictMode {
		for i := range issues {
			if issues[i].Level == LevelWarning {
				issues[i].Level = LevelError
			}
		}
	}

	// Create and return the validation report
	report := NewValidationReport(issues)

	return report, nil
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
	if containsPattern(content, "##### Scenario:") {
		return true
	}

	// Check for ###### Scenario: (6 hashtags - wrong)
	if containsPattern(content, "###### Scenario:") {
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
func containsPattern(content, pattern string) bool {
	return len(content) > 0 && len(pattern) > 0 &&
		strings.Contains(content, pattern)
}
