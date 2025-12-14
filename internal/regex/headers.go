package regex

import "regexp"

// Header patterns - all pre-compiled at package init
var (
	// H2SectionHeader matches "## Section Name" and captures the section name.
	// Used for extracting section headers from spec files.
	H2SectionHeader = regexp.MustCompile(`^##\s+(.+)$`)

	// H2DeltaSection matches "## ADDED|MODIFIED|REMOVED|RENAMED Requirements"
	// and captures the delta type (ADDED, MODIFIED, REMOVED, or RENAMED).
	H2DeltaSection = regexp.MustCompile(
		`^##\s+(ADDED|MODIFIED|REMOVED|RENAMED)\s+Requirements\s*$`,
	)

	// H2RequirementsSection matches exactly "## Requirements" section header.
	// Uses multiline mode for searching within multi-line content.
	H2RequirementsSection = regexp.MustCompile(`(?m)^##\s+Requirements\s*$`)

	// H2NextSection matches any "## " header for finding section boundaries.
	// Uses multiline mode for searching within multi-line content.
	H2NextSection = regexp.MustCompile(`(?m)^##\s+`)

	// H2Pattern matches any H2 header (## followed by space).
	// Used for detecting section boundaries without capturing content.
	H2Pattern = regexp.MustCompile(`^##\s+`)

	// H3Requirement matches "### Requirement: Name" and captures the name.
	H3Requirement = regexp.MustCompile(`^###\s+Requirement:\s*(.+)$`)

	// H3AnyHeader matches any "### " header for section boundary detection.
	H3AnyHeader = regexp.MustCompile(`^###\s+`)

	// H4Scenario matches "#### Scenario: Name" and captures the name.
	H4Scenario = regexp.MustCompile(`^####\s+Scenario:\s*(.+)$`)
)

// MatchH2SectionHeader checks if a line is an H2 section header
// and extracts the name. Returns the section name (raw capture, untrimmed) and
// true if matched, or empty string and false otherwise.
func MatchH2SectionHeader(line string) (name string, ok bool) {
	matches := H2SectionHeader.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// MatchH2DeltaSection checks if a line is a delta section header.
// Returns the delta type (ADDED, MODIFIED, REMOVED, or RENAMED) and
// true if matched, or empty string and false otherwise.
func MatchH2DeltaSection(line string) (deltaType string, ok bool) {
	matches := H2DeltaSection.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// MatchH3Requirement checks if a line is a requirement header and
// extracts the name. Returns the requirement name (raw capture, untrimmed) and
// true if matched, or empty string and false otherwise.
// Callers should use strings.TrimSpace() on the name if trimming is needed.
func MatchH3Requirement(line string) (name string, ok bool) {
	matches := H3Requirement.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// MatchH4Scenario checks if a line is a scenario header and extracts the name.
// Returns the scenario name (raw capture, untrimmed) and true if matched,
// or empty string and false otherwise.
// Callers should use strings.TrimSpace() on the name if trimming is needed.
func MatchH4Scenario(line string) (name string, ok bool) {
	matches := H4Scenario.FindStringSubmatch(line)
	if len(matches) < 2 {
		return "", false
	}

	return matches[1], true
}

// IsH2Header checks if a line starts with "## " (any H2 header).
func IsH2Header(line string) bool {
	return H2Pattern.MatchString(line)
}

// IsH3Header checks if a line starts with "### " (any H3 header).
func IsH3Header(line string) bool {
	return H3AnyHeader.MatchString(line)
}
