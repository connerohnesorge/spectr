//nolint:revive // file-length-limit,argument-limit - validation logic requires comprehensive parameters
package validation

import (
	"errors"
	"fmt"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// validateAddedRequirements validates ADDED Requirements section
func validateAddedRequirements(
	addedContent, specPath string,
	lines []string,
	sectionLine int,
	fileAddedReqs map[string]bool,
	addedReqs map[string]string,
) []ValidationIssue {
	var issues []ValidationIssue
	requirements := ExtractRequirements(addedContent)

	if len(requirements) == 0 {
		issues = append(issues, ValidationIssue{
			Level: LevelError,
			Path:  specPath,
			Line:  sectionLine,
			Message: "ADDED Requirements section is empty " +
				"(no requirements found)",
		})

		return issues
	}

	for _, req := range requirements {
		normalized := NormalizeRequirementName(req.Name)
		reqPath := fmt.Sprintf(
			"%s: ADDED Requirement '%s'",
			specPath,
			req.Name,
		)
		reqLine := findRequirementLineInSection(lines, req.Name, sectionLine)

		// Check for SHALL/MUST
		if !ContainsShallOrMust(req.Content) {
			issues = append(issues, ValidationIssue{
				Level:   LevelError,
				Path:    reqPath,
				Line:    reqLine,
				Message: "ADDED requirement must contain SHALL or MUST",
			})
		}

		// Check for at least one scenario
		if len(req.Scenarios) == 0 {
			issues = append(issues, ValidationIssue{
				Level:   LevelError,
				Path:    reqPath,
				Line:    reqLine,
				Message: "ADDED requirement must have at least one scenario",
			})
		}

		// Check for duplicate within this file
		if fileAddedReqs[normalized] {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqLine,
				Message: fmt.Sprintf(
					"Duplicate requirement name in ADDED section: '%s'",
					req.Name,
				),
			})
		}
		fileAddedReqs[normalized] = true

		// Check for duplicate across files
		if existingPath, exists := addedReqs[normalized]; exists {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqLine,
				Message: fmt.Sprintf(
					"Requirement '%s' is ADDED in multiple files: "+
						"%s and %s",
					req.Name,
					existingPath,
					specPath,
				),
			})
		} else {
			addedReqs[normalized] = specPath
		}

		// Check for malformed scenarios
		if len(req.Scenarios) == 0 && hasMalformedScenarios(req.Content) {
			malformedLine := findMalformedScenarioLineInDelta(lines, reqLine)
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  malformedLine,
				Message: "Scenarios must use '#### Scenario:' format " +
					"(4 hashtags followed by 'Scenario:')",
			})
		}
	}

	return issues
}

// validateModifiedRequirements validates MODIFIED Requirements section
func validateModifiedRequirements(
	modifiedContent, specPath string,
	lines []string,
	sectionLine int,
	fileModifiedReqs map[string]bool,
	modifiedReqs map[string]string,
) []ValidationIssue {
	var issues []ValidationIssue
	requirements := ExtractRequirements(modifiedContent)

	if len(requirements) == 0 {
		issues = append(issues, ValidationIssue{
			Level: LevelError,
			Path:  specPath,
			Line:  sectionLine,
			Message: "MODIFIED Requirements section is empty " +
				"(no requirements found)",
		})

		return issues
	}

	for _, req := range requirements {
		normalized := NormalizeRequirementName(req.Name)
		reqPath := fmt.Sprintf(
			"%s: MODIFIED Requirement '%s'",
			specPath,
			req.Name,
		)
		reqLine := findRequirementLineInSection(lines, req.Name, sectionLine)

		// Check for SHALL/MUST
		if !ContainsShallOrMust(req.Content) {
			issues = append(issues, ValidationIssue{
				Level:   LevelError,
				Path:    reqPath,
				Line:    reqLine,
				Message: "MODIFIED requirement must contain SHALL or MUST",
			})
		}

		// Check for at least one scenario
		if len(req.Scenarios) == 0 {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqLine,
				Message: "MODIFIED requirement must have " +
					"at least one scenario",
			})
		}

		// Check for duplicate within this file
		if fileModifiedReqs[normalized] {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqLine,
				Message: fmt.Sprintf(
					"Duplicate requirement name in MODIFIED section: '%s'",
					req.Name,
				),
			})
		}
		fileModifiedReqs[normalized] = true

		// Check for duplicate across files
		if existingPath, exists := modifiedReqs[normalized]; exists {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqLine,
				Message: fmt.Sprintf(
					"Requirement '%s' is MODIFIED in multiple files: "+
						"%s and %s",
					req.Name,
					existingPath,
					specPath,
				),
			})
		} else {
			modifiedReqs[normalized] = specPath
		}

		// Check for malformed scenarios
		if len(req.Scenarios) == 0 && hasMalformedScenarios(req.Content) {
			malformedLine := findMalformedScenarioLineInDelta(lines, reqLine)
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  malformedLine,
				Message: "Scenarios must use '#### Scenario:' format " +
					"(4 hashtags followed by 'Scenario:')",
			})
		}
	}

	return issues
}

// validateRemovedRequirements validates REMOVED Requirements section
func validateRemovedRequirements(
	removedContent, specPath string,
	lines []string,
	sectionLine int,
	fileRemovedReqs map[string]bool,
	removedReqs map[string]string,
) []ValidationIssue {
	var issues []ValidationIssue
	requirements := ExtractRequirements(removedContent)

	if len(requirements) == 0 {
		issues = append(issues, ValidationIssue{
			Level: LevelError,
			Path:  specPath,
			Line:  sectionLine,
			Message: "REMOVED Requirements section is empty " +
				"(no requirements found)",
		})

		return issues
	}

	for _, req := range requirements {
		normalized := NormalizeRequirementName(req.Name)
		reqPath := fmt.Sprintf(
			"%s: REMOVED Requirement '%s'",
			specPath,
			req.Name,
		)
		reqLine := findRequirementLineInSection(lines, req.Name, sectionLine)

		// Check for duplicate within this file
		if fileRemovedReqs[normalized] {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqLine,
				Message: fmt.Sprintf(
					"Duplicate requirement name in REMOVED section: '%s'",
					req.Name,
				),
			})
		}
		fileRemovedReqs[normalized] = true

		// Check for duplicate across files
		if existingPath, exists := removedReqs[normalized]; exists {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqLine,
				Message: fmt.Sprintf(
					"Requirement '%s' is REMOVED in multiple files: "+
						"%s and %s",
					req.Name,
					existingPath,
					specPath,
				),
			})
		} else {
			removedReqs[normalized] = specPath
		}
	}

	return issues
}

// validateRenamedRequirements validates RENAMED Requirements section
func validateRenamedRequirements(
	renamedContent, specPath string,
	lines []string,
	sectionLine int,
	fileRenamedFromReqs, fileRenamedToReqs map[string]bool,
	renamedFromReqs, renamedToReqs map[string]string,
) []ValidationIssue {
	var issues []ValidationIssue
	renames := parseRenamedRequirements(renamedContent)

	if len(renames) == 0 {
		issues = append(issues, ValidationIssue{
			Level: LevelError,
			Path:  specPath,
			Line:  sectionLine,
			Message: "RENAMED Requirements section is empty " +
				"(no rename pairs found)",
		})

		return issues
	}

	for _, rename := range renames {
		if rename.FromName == "" || rename.ToName == "" {
			// Try to find the line with malformed rename
			malformedLine := findMalformedRenamedLine(lines, sectionLine)
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  fmt.Sprintf("%s: RENAMED Requirements", specPath),
				Line:  malformedLine,
				Message: "Malformed RENAMED requirement " +
					"(expected format: '- FROM: ### Requirement: " +
					"OldName' followed by '- TO: ### Requirement: NewName')",
			})

			continue
		}

		normalizedFrom := NormalizeRequirementName(rename.FromName)
		normalizedTo := NormalizeRequirementName(rename.ToName)
		reqPath := fmt.Sprintf(
			"%s: RENAMED Requirement '%s' -> '%s'",
			specPath,
			rename.FromName,
			rename.ToName,
		)
		reqFromLine := findRenamedPairLine(lines, rename.FromName, sectionLine)
		reqToLine := findRenamedPairLine(lines, rename.ToName, reqFromLine)

		// Check for duplicate FROM names within this file
		if fileRenamedFromReqs[normalizedFrom] {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqFromLine,
				Message: fmt.Sprintf(
					"Duplicate FROM requirement name in "+
						"RENAMED section: '%s'",
					rename.FromName,
				),
			})
		}
		fileRenamedFromReqs[normalizedFrom] = true

		// Check for duplicate TO names within this file
		if fileRenamedToReqs[normalizedTo] {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqToLine,
				Message: fmt.Sprintf(
					"Duplicate TO requirement name in "+
						"RENAMED section: '%s'",
					rename.ToName,
				),
			})
		}
		fileRenamedToReqs[normalizedTo] = true

		// Check for duplicate FROM across files
		if existingPath, exists := renamedFromReqs[normalizedFrom]; exists {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqFromLine,
				Message: fmt.Sprintf(
					"Requirement '%s' is renamed (FROM) in "+
						"multiple files: %s and %s",
					rename.FromName,
					existingPath,
					specPath,
				),
			})
		} else {
			renamedFromReqs[normalizedFrom] = specPath
		}

		// Check for duplicate TO across files
		if existingPath, exists := renamedToReqs[normalizedTo]; exists {
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  reqPath,
				Line:  reqToLine,
				Message: fmt.Sprintf(
					"Requirement '%s' is renamed (TO) in "+
						"multiple files: %s and %s",
					rename.ToName,
					existingPath,
					specPath,
				),
			})
		} else {
			renamedToReqs[normalizedTo] = specPath
		}
	}

	return issues
}

// ValidatePreMerge validates delta operations against base spec.
// It checks that:
// - ADDED requirements don't already exist in base spec
// - MODIFIED/REMOVED/RENAMED requirements DO exist in base spec
// - RENAMED TO requirements don't already exist (unless renaming to itself)
//
// If specExists is false, only ADDED operations are allowed.
//
//nolint:revive // specExists is a legitimate control parameter
func ValidatePreMerge(baseSpecPath string, deltaPlan *parsers.DeltaPlan, specExists bool) error {
	// If spec doesn't exist, only ADDED operations are allowed
	if !specExists {
		if len(deltaPlan.Modified) > 0 || len(deltaPlan.Removed) > 0 || len(deltaPlan.Renamed) > 0 {
			return errors.New(
				"target spec does not exist; only ADDED requirements are allowed for new specs",
			)
		}

		return nil
	}

	// Parse base spec to get existing requirements
	baseReqs, err := parsers.ParseRequirements(baseSpecPath)
	if err != nil {
		return fmt.Errorf("parse base spec: %w", err)
	}

	// Build map of existing requirement names (normalized)
	existing := make(map[string]bool)
	for _, req := range baseReqs {
		normalized := parsers.NormalizeRequirementName(req.Name)
		existing[normalized] = true
	}

	// Validate MODIFIED requirements exist in base
	for _, req := range deltaPlan.Modified {
		normalized := parsers.NormalizeRequirementName(req.Name)
		if !existing[normalized] {
			return fmt.Errorf("MODIFIED requirement %q does not exist in base spec", req.Name)
		}
	}

	// Validate REMOVED requirements exist in base
	for _, name := range deltaPlan.Removed {
		normalized := parsers.NormalizeRequirementName(name)
		if !existing[normalized] {
			return fmt.Errorf("REMOVED requirement %q does not exist in base spec", name)
		}
	}

	// Validate RENAMED FROM requirements exist in base
	for _, op := range deltaPlan.Renamed {
		fromNorm := parsers.NormalizeRequirementName(op.From)
		if !existing[fromNorm] {
			return fmt.Errorf("RENAMED FROM requirement %q does not exist in base spec", op.From)
		}

		// Check that TO name doesn't already exist (unless it's being renamed from something else)
		toNorm := parsers.NormalizeRequirementName(op.To)
		if existing[toNorm] && toNorm != fromNorm {
			return fmt.Errorf("RENAMED TO requirement %q already exists in base spec", op.To)
		}
	}

	// Validate ADDED requirements don't exist in base
	for _, req := range deltaPlan.Added {
		normalized := parsers.NormalizeRequirementName(req.Name)
		if existing[normalized] {
			return fmt.Errorf("ADDED requirement %q already exists in base spec", req.Name)
		}
	}

	return nil
}
