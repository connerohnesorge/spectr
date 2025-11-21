// Package validation provides validation rules for Spectr changes and specs.
//
//nolint:revive // file-length-limit,argument-limit,line-length-limit - validation logic requires comprehensive parameters
package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parsers"
)

// DeltaType represents the type of delta operation
type DeltaType string

const (
	DeltaAdded    DeltaType = "ADDED"
	DeltaModified DeltaType = "MODIFIED"
	DeltaRemoved  DeltaType = "REMOVED"
	DeltaRenamed  DeltaType = "RENAMED"
)

// ValidateChangeDeltaSpecs validates all delta spec files in a
// change directory. changeDir should be the path to a change directory
// (e.g., spectr/changes/add-feature). spectrRoot should be the path to the
// spectr/ directory (e.g., /path/to/project/spectr).
// Returns ValidationReport with all issues found, or error for issues
//
//nolint:revive // strictMode is intentional control flag
func ValidateChangeDeltaSpecs(
	changeDir string,
	spectrRoot string,
	strictMode bool,
) (*ValidationReport, error) {
	specsDir := filepath.Join(changeDir, "specs")

	// Check if specs directory exists
	info, err := os.Stat(specsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("specs directory not found: %s", specsDir)
		}

		return nil, fmt.Errorf("failed to access specs directory: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("specs path is not a directory: %s", specsDir)
	}

	// Find all spec.md files under specs/
	var specFiles []string
	err = filepath.Walk(specsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "spec.md" {
			specFiles = append(specFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk specs directory: %w", err)
	}

	if len(specFiles) == 0 {
		return nil, fmt.Errorf("no spec.md files found in specs directory: %s", specsDir)
	}

	// Track all issues across all spec files
	var allIssues []ValidationIssue

	// Track requirement names for duplicate/conflict detection across all files
	addedReqs := make(map[string]string)       // normalized name -> file path
	modifiedReqs := make(map[string]string)    // normalized name -> file path
	removedReqs := make(map[string]string)     // normalized name -> file path
	renamedFromReqs := make(map[string]string) // normalized FROM name -> file path
	renamedToReqs := make(map[string]string)   // normalized TO name -> file path

	// Track total delta count
	totalDeltas := 0

	// Process each spec file
	for _, specPath := range specFiles {
		fileIssues, deltaCount, err := validateSingleDeltaFile(
			specPath,
			addedReqs,
			modifiedReqs,
			removedReqs,
			renamedFromReqs,
			renamedToReqs,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to validate %s: %w", specPath, err)
		}

		allIssues = append(allIssues, fileIssues...)
		totalDeltas += deltaCount

		// Validate delta file against base spec
		baseSpecIssues, err := validateDeltaAgainstBaseSpec(specPath, specsDir, spectrRoot)
		if err != nil {
			return nil, fmt.Errorf("failed to validate %s against base spec: %w", specPath, err)
		}
		allIssues = append(allIssues, baseSpecIssues...)
	}

	// Check if there are no deltas at all
	if totalDeltas == 0 {
		allIssues = append(allIssues, ValidationIssue{
			Level: LevelError,
			Path:  specsDir,
			Line:  1, // Default to line 1 for missing deltas
			Message: "Change must have at least one delta " +
				"(ADDED, MODIFIED, REMOVED, or RENAMED requirement)",
		})
	}

	// Apply strict mode: convert warnings to errors
	if strictMode {
		for i := range allIssues {
			if allIssues[i].Level == LevelWarning {
				allIssues[i].Level = LevelError
			}
		}
	}

	return NewValidationReport(allIssues), nil
}

// validateSingleDeltaFile validates a single spec.md delta file
// Returns issues, delta count, and error
func validateSingleDeltaFile(
	specPath string,
	addedReqs, modifiedReqs, removedReqs, renamedFromReqs, renamedToReqs map[string]string,
) ([]ValidationIssue, int, error) {
	// Read file
	content, err := os.ReadFile(specPath)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to read file: %w", err)
	}

	contentStr := string(content)
	lines := strings.Split(contentStr, "\n")

	// Parse sections
	sections := ExtractSections(contentStr)
	var issues []ValidationIssue
	deltaCount := 0

	// Track requirement names within this file for duplicate detection
	fileAddedReqs := make(map[string]bool)
	fileModifiedReqs := make(map[string]bool)
	fileRemovedReqs := make(map[string]bool)
	fileRenamedFromReqs := make(map[string]bool)
	fileRenamedToReqs := make(map[string]bool)

	// Process ADDED Requirements
	if addedContent, hasAdded := sections["ADDED Requirements"]; hasAdded {
		deltaCount++
		addedLine := findDeltaSectionLine(lines, "ADDED Requirements")
		addedIssues := validateAddedRequirements(
			addedContent,
			specPath,
			lines,
			addedLine,
			fileAddedReqs,
			addedReqs,
		)
		issues = append(issues, addedIssues...)
	}

	// Process MODIFIED Requirements
	if modifiedContent, hasModified := sections["MODIFIED Requirements"]; hasModified {
		deltaCount++
		modifiedLine := findDeltaSectionLine(lines, "MODIFIED Requirements")
		modifiedIssues := validateModifiedRequirements(
			modifiedContent,
			specPath,
			lines,
			modifiedLine,
			fileModifiedReqs,
			modifiedReqs,
		)
		issues = append(issues, modifiedIssues...)
	}

	// Process REMOVED Requirements
	if removedContent, hasRemoved := sections["REMOVED Requirements"]; hasRemoved {
		deltaCount++
		removedLine := findDeltaSectionLine(lines, "REMOVED Requirements")
		removedIssues := validateRemovedRequirements(
			removedContent,
			specPath,
			lines,
			removedLine,
			fileRemovedReqs,
			removedReqs,
		)
		issues = append(issues, removedIssues...)
	}

	// Process RENAMED Requirements
	if renamedContent, hasRenamed := sections["RENAMED Requirements"]; hasRenamed {
		deltaCount++
		renamedLine := findDeltaSectionLine(lines, "RENAMED Requirements")
		renamedIssues := validateRenamedRequirements(
			renamedContent,
			specPath,
			lines,
			renamedLine,
			fileRenamedFromReqs,
			fileRenamedToReqs,
			renamedFromReqs,
			renamedToReqs,
		)
		issues = append(issues, renamedIssues...)
	}

	// Check for cross-section conflicts within this file
	for normalized := range fileAddedReqs {
		if fileModifiedReqs[normalized] {
			// Find actual requirement name for better error message
			reqName := findRequirementNameByNormalized(
				normalized,
				addedReqs,
				specPath,
			)
			// Find line number of the requirement in ADDED section
			addedLine := findDeltaSectionLine(lines, "ADDED Requirements")
			reqLine := findRequirementLineInSection(lines, reqName, addedLine)
			issues = append(issues, ValidationIssue{
				Level: LevelError,
				Path:  specPath,
				Line:  reqLine,
				Message: fmt.Sprintf(
					"Requirement '%s' appears in both ADDED and "+
						"MODIFIED sections",
					reqName,
				),
			})
		}
	}

	return issues, deltaCount, nil
}

// RenamedRequirement represents a FROM -> TO rename pair
type RenamedRequirement struct {
	FromName string
	ToName   string
}

// parseRenamedRequirements parses the RENAMED Requirements section
// Expected format:
// - FROM: ### Requirement: OldName
// - TO: ### Requirement: NewName
func parseRenamedRequirements(content string) []RenamedRequirement {
	var renames []RenamedRequirement
	lines := strings.Split(content, "\n")

	fromRegex := regexp.MustCompile(
		`^\s*-\s*FROM:\s*###\s*Requirement:\s*(.+)$`,
	)
	toRegex := regexp.MustCompile(
		`^\s*-\s*TO:\s*###\s*Requirement:\s*(.+)$`,
	)

	var currentFrom string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check for FROM line
		if matches := fromRegex.FindStringSubmatch(line); matches != nil {
			currentFrom = strings.TrimSpace(matches[1])

			continue
		}

		// Check for TO line
		matches := toRegex.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		toName := strings.TrimSpace(matches[1])

		// If we have a FROM, pair it with this TO
		if currentFrom != "" {
			renames = append(renames, RenamedRequirement{
				FromName: currentFrom,
				ToName:   toName,
			})
			currentFrom = "" // Reset for next pair
		} else {
			// TO without FROM - add as malformed
			renames = append(renames, RenamedRequirement{
				FromName: "",
				ToName:   toName,
			})
		}
	}

	// If we have a FROM without a TO, add as malformed
	if currentFrom != "" {
		renames = append(renames, RenamedRequirement{
			FromName: currentFrom,
			ToName:   "",
		})
	}

	return renames
}

// findRequirementNameByNormalized finds the original requirement name
// given a normalized version
func findRequirementNameByNormalized(
	normalized string,
	_ map[string]string,
	_ string,
) string {
	// This is a helper to provide better error messages
	// In practice, we'd need to iterate through requirements again
	// For now, return the normalized version as it's close enough
	return normalized
}

// validateDeltaAgainstBaseSpec validates a delta file against the base spec
// Returns validation issues for pre-merge validation errors
func validateDeltaAgainstBaseSpec(
	deltaSpecPath string,
	specsDir string,
	spectrRoot string,
) ([]ValidationIssue, error) {
	// Extract capability name from delta spec path
	// Path structure: .../changes/<change-id>/specs/<capability>/spec.md
	// We want to extract <capability>
	relPath, err := filepath.Rel(specsDir, deltaSpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get relative path: %w", err)
	}

	// Extract capability (should be the directory name before spec.md)
	capability := filepath.Dir(relPath)
	if capability == "." || capability == "" {
		return nil, fmt.Errorf("invalid delta spec path structure: %s", deltaSpecPath)
	}

	// Construct base spec path
	baseSpecPath := filepath.Join(spectrRoot, "specs", capability, "spec.md")

	// Check if base spec exists
	_, err = os.Stat(baseSpecPath)
	baseExists := err == nil

	// Parse delta spec
	deltaPlan, err := parsers.ParseDeltaSpec(deltaSpecPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse delta spec: %w", err)
	}

	// Validate delta against base spec
	if err := ValidatePreMerge(baseSpecPath, deltaPlan, baseExists); err != nil {
		// Read delta file to find line number
		content, readErr := os.ReadFile(deltaSpecPath)
		lineNum := 1
		if readErr == nil {
			lines := strings.Split(string(content), "\n")
			lineNum = findPreMergeErrorLine(lines, err.Error(), deltaPlan)
		}

		return []ValidationIssue{
			{
				Level:   LevelError,
				Path:    deltaSpecPath,
				Line:    lineNum,
				Message: err.Error(),
			},
		}, nil
	}

	return nil, nil
}

// findDeltaSectionLine finds the line number where a delta section header appears
// Returns 1 if not found
func findDeltaSectionLine(lines []string, sectionName string) int {
	sectionHeader := "## " + sectionName
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), sectionHeader) {
			return i + 1 // Line numbers are 1-indexed
		}
	}

	return 1 // Default to line 1 if not found
}

// findRequirementLineInSection finds the line number where a requirement appears
// within a specific section (searches from startLine onwards)
// Returns startLine if not found
func findRequirementLineInSection(lines []string, reqName string, startLine int) int {
	reqHeader := "### Requirement: " + reqName
	// Start searching from startLine (convert to 0-indexed)
	searchStart := startLine - 1
	if searchStart < 0 {
		searchStart = 0
	}

	for i := searchStart; i < len(lines); i++ {
		if strings.HasPrefix(strings.TrimSpace(lines[i]), reqHeader) {
			return i + 1 // Line numbers are 1-indexed
		}
		// Stop if we hit another section
		if i > searchStart && strings.HasPrefix(strings.TrimSpace(lines[i]), "## ") {
			break
		}
	}

	return startLine // Default to section start if not found
}

// findRenamedPairLine finds the line number of a RENAMED requirement pair
// Searches for the FROM line
func findRenamedPairLine(lines []string, fromName string, startLine int) int {
	searchStart := startLine - 1
	if searchStart < 0 {
		searchStart = 0
	}

	for i := searchStart; i < len(lines); i++ {
		trimmed := strings.TrimSpace(lines[i])
		if i > searchStart && strings.HasPrefix(trimmed, "## ") {
			break
		}

		withoutBullet := strings.TrimSpace(strings.TrimPrefix(trimmed, "-"))
		withoutCode := strings.Trim(withoutBullet, "`")
		if strings.HasPrefix(
			withoutCode,
			"FROM: ### Requirement: "+fromName,
		) || strings.HasPrefix(
			withoutCode,
			"TO: ### Requirement: "+fromName,
		) {
			return i + 1 // Line numbers are 1-indexed
		}
	}

	return startLine // Default to section start if not found
}

// findPreMergeErrorLine finds the line number related to a pre-merge validation error
func findPreMergeErrorLine(lines []string, errMsg string, deltaPlan *parsers.DeltaPlan) int {
	// Extract requirement name from error message
	// Error formats:
	// - "MODIFIED requirement %q does not exist in base spec"
	// - "REMOVED requirement %q does not exist in base spec"
	// - "RENAMED FROM requirement %q does not exist in base spec"
	// - "RENAMED TO requirement %q already exists in base spec"
	// - "ADDED requirement %q already exists in base spec"

	sectionName, reqName := extractSectionAndReqName(errMsg)

	if reqName == "" {
		return 1 // Can't find requirement name, default to line 1
	}

	sectionLine := findDeltaSectionLine(lines, sectionName)
	if strings.Contains(sectionName, "RENAMED") {
		return findRenamedPairLine(lines, reqName, sectionLine)
	}

	return findRequirementLineInSection(lines, reqName, sectionLine)
}

// extractSectionAndReqName extracts the section name and requirement name from an error message
func extractSectionAndReqName(errMsg string) (string, string) {
	var sectionName, reqName string

	// Determine section name based on error message
	switch {
	case strings.Contains(errMsg, "MODIFIED requirement"):
		sectionName = "MODIFIED Requirements"
	case strings.Contains(errMsg, "REMOVED requirement"):
		sectionName = "REMOVED Requirements"
	case strings.Contains(errMsg, "RENAMED FROM requirement"):
		sectionName = "RENAMED Requirements"
	case strings.Contains(errMsg, "RENAMED TO requirement"):
		sectionName = "RENAMED Requirements"
	case strings.Contains(errMsg, "ADDED requirement"):
		sectionName = "ADDED Requirements"
	default:
		return sectionName, reqName
	}

	// Extract requirement name from quotes
	if start := strings.Index(errMsg, "\""); start != -1 {
		if end := strings.Index(errMsg[start+1:], "\""); end != -1 {
			reqName = errMsg[start+1 : start+1+end]
		}
	}

	return sectionName, reqName
}

// findMalformedScenarioLineInDelta finds the line number of a malformed scenario in delta file
// Searches from reqLine onwards
// Returns reqLine if not found
func findMalformedScenarioLineInDelta(lines []string, reqLine int) int {
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
		// Stop if we hit another requirement or section
		if i > searchStart {
			if strings.HasPrefix(line, "### Requirement:") {
				break
			}
			if strings.HasPrefix(line, "## ") {
				break
			}
		}
	}

	return reqLine // Default to requirement line if not found
}

// findMalformedRenamedLine finds the line number of a malformed RENAMED requirement
// Searches from sectionLine onwards for FROM or TO without its pair
func findMalformedRenamedLine(lines []string, sectionLine int) int {
	searchStart := sectionLine - 1
	if searchStart < 0 {
		searchStart = 0
	}

	for i := searchStart; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		// Check for FROM or TO lines
		if strings.Contains(line, "FROM:") || strings.Contains(line, "TO:") {
			return i + 1 // Line numbers are 1-indexed
		}
		// Stop if we hit another section
		if i > searchStart && strings.HasPrefix(line, "## ") {
			break
		}
	}

	return sectionLine // Default to section start if not found
}
