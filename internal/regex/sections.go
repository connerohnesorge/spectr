package regex

import (
	"fmt"
	"regexp"
)

// FindSectionContent extracts content between a specific H2 section
// and the next H2. The sectionHeader parameter is the exact text after
// "## " (e.g., "Requirements"). Returns empty string if section not found.
//
// Example:
//
//	content := "## Purpose\nSome text\n## Requirements\nReq content"
//	FindSectionContent(content, "Requirements") // returns "\nReq content"
func FindSectionContent(content, sectionHeader string) string {
	patternStr := fmt.Sprintf(
		`(?m)^##\s+%s\s*$`,
		regexp.QuoteMeta(sectionHeader),
	)
	pattern := regexp.MustCompile(patternStr)
	matches := pattern.FindStringIndex(content)
	if matches == nil {
		return ""
	}

	sectionStart := matches[1]
	nextMatches := H2NextSection.FindStringIndex(content[sectionStart:])

	if nextMatches != nil {
		return content[sectionStart : sectionStart+nextMatches[0]]
	}

	return content[sectionStart:]
}

// FindDeltaSectionContent extracts content from a delta section
// (ADDED, MODIFIED, etc.). This is a convenience wrapper around
// FindSectionContent for delta specs.
//
// Example:
//
//	FindDeltaSectionContent(content, "ADDED") // finds "## ADDED Reqs"
func FindDeltaSectionContent(content, deltaType string) string {
	return FindSectionContent(content, deltaType+" Requirements")
}

// FindRequirementsSection extracts the "## Requirements" section content.
// This is a convenience wrapper around FindSectionContent for spec files.
func FindRequirementsSection(content string) string {
	return FindSectionContent(content, "Requirements")
}

// FindAllH3Requirements finds all requirement headers in content and
// returns their names. Uses multiline mode to find all occurrences of
// "### Requirement: Name".
func FindAllH3Requirements(content string) []string {
	pattern := regexp.MustCompile(`(?m)^###\s+Requirement:\s*(.+)$`)
	matches := pattern.FindAllStringSubmatch(content, -1)

	names := make([]string, 0, len(matches))
	for _, match := range matches {
		if len(match) > 1 {
			names = append(names, match[1])
		}
	}

	return names
}

// FindSectionIndex returns the start and end indices of a section
// header. Returns nil if the section is not found.
func FindSectionIndex(content, sectionHeader string) []int {
	patternStr := fmt.Sprintf(
		`(?m)^##\s+%s\s*$`,
		regexp.QuoteMeta(sectionHeader),
	)
	pattern := regexp.MustCompile(patternStr)

	return pattern.FindStringIndex(content)
}
