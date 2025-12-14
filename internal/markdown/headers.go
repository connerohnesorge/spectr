package markdown

import (
	"strings"
)

// Header prefix constants
const (
	h2Prefix = "## "
	h3Prefix = "### "
	h4Prefix = "####"
	h5Prefix = "#####"
)

// Line Matchers
//
// These functions check individual lines for header patterns.
// They are designed for efficient line-by-line iteration:
//
//	for _, line := range lines {
//	    if name, ok := markdown.MatchH3Requirement(line); ok {
//	        // handle requirement
//	    }
//	}
//
// For bulk extraction, use ParseDocument + GetRequirementNames instead.

// MatchH2SectionHeader checks if a line is an H2 section header
// and extracts the name. Returns the section name and true if matched.
func MatchH2SectionHeader(line string) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, h2Prefix) {
		return "", false
	}

	// Must not be H3 or higher
	if strings.HasPrefix(trimmed, h3Prefix) {
		return "", false
	}

	name = strings.TrimSpace(strings.TrimPrefix(trimmed, h2Prefix))

	return name, name != ""
}

// MatchH2DeltaSection checks if a line is a delta section header.
// Returns the delta type (ADDED, MODIFIED, REMOVED, RENAMED) and true.
func MatchH2DeltaSection(line string) (deltaType string, ok bool) {
	name, ok := MatchH2SectionHeader(line)
	if !ok {
		return "", false
	}

	for _, dt := range deltaTypes {
		if name == dt+" Requirements" {
			return dt, true
		}
	}

	return "", false
}

// MatchH3Requirement checks if a line is a requirement header.
// Returns the requirement name and true if matched.
func MatchH3Requirement(line string) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, h3Prefix) {
		return "", false
	}

	// Must not be H4 or higher
	if strings.HasPrefix(trimmed, h4Prefix) {
		return "", false
	}

	rest := strings.TrimPrefix(trimmed, h3Prefix)
	if !strings.HasPrefix(rest, "Requirement:") {
		return "", false
	}

	name = strings.TrimSpace(strings.TrimPrefix(rest, "Requirement:"))

	return name, name != ""
}

// MatchH4Scenario checks if a line is a scenario header and extracts the name.
// Returns the scenario name and true if matched.
func MatchH4Scenario(line string) (name string, ok bool) {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, h4Prefix) {
		return "", false
	}

	// Must not be H5 or higher
	if strings.HasPrefix(trimmed, h5Prefix) {
		return "", false
	}

	// Get rest after ####, handling variable whitespace
	rest := strings.TrimPrefix(trimmed, h4Prefix)
	rest = strings.TrimLeft(rest, " \t")

	// Check for Scenario: prefix
	if !strings.HasPrefix(rest, "Scenario:") {
		return "", false
	}

	name = strings.TrimSpace(strings.TrimPrefix(rest, "Scenario:"))

	return name, name != ""
}
