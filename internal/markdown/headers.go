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

// IsH2Header checks if a line starts with "## " (any H2 header).
func IsH2Header(line string) bool {
	trimmed := strings.TrimSpace(line)
	isH2 := strings.HasPrefix(trimmed, h2Prefix)
	isH3 := strings.HasPrefix(trimmed, h3Prefix)

	return isH2 && !isH3
}

// IsH3Header checks if a line starts with "### " (any H3 header).
func IsH3Header(line string) bool {
	trimmed := strings.TrimSpace(line)
	isH3 := strings.HasPrefix(trimmed, h3Prefix)
	isH4 := strings.HasPrefix(trimmed, h4Prefix)

	return isH3 && !isH4
}

// IsH4Header checks if a line starts with "#### " (any H4 header).
func IsH4Header(line string) bool {
	trimmed := strings.TrimSpace(line)
	if !strings.HasPrefix(trimmed, h4Prefix) {
		return false
	}
	// Must not be H5 or higher
	if strings.HasPrefix(trimmed, h5Prefix) {
		return false
	}
	// Must have space after ####
	rest := strings.TrimPrefix(trimmed, h4Prefix)

	return len(rest) > 0 && (rest[0] == ' ' || rest[0] == '\t')
}

// Convenience Functions
//
// These parse content and return extracted data.
// They use ParseDocument internally for AST-based extraction.

// FindRequirementHeaders finds all H3 requirement headers in content.
// Returns a slice of requirement names in document order.
//
// For multiple queries on the same content, prefer:
//
//	doc, _ := ParseDocument([]byte(content))
//	names := doc.GetRequirementNames()
func FindRequirementHeaders(content string) []string {
	doc, err := ParseDocument([]byte(content))
	if err != nil {
		return nil
	}

	return doc.GetRequirementNames()
}

// FindScenarioHeaders finds all H4 scenario headers in content.
// Returns a slice of scenario names.
func FindScenarioHeaders(content string) []string {
	doc, err := ParseDocument([]byte(content))
	if err != nil {
		return nil
	}

	var names []string
	for _, h := range doc.H4Headers {
		if name := parseScenarioName(h.Text); name != "" {
			names = append(names, name)
		}
	}

	return names
}
