package markdown

import (
	"strings"
)

// Section Extraction Functions
//
// These functions parse content and extract section data.
// They use ParseDocument internally for AST-based extraction.
//
// For multiple queries on the same content, prefer using ParseDocument once:
//
//	doc, _ := ParseDocument([]byte(content))
//	reqs := doc.GetSectionContent("Requirements")
//	added := doc.GetSectionContent("ADDED Requirements")

// FindSectionContent extracts content between a specific H2 section
// and the next H2. The sectionHeader parameter is the exact text after
// "## " (e.g., "Requirements"). Returns empty string if section not found.
func FindSectionContent(content, sectionHeader string) string {
	doc, err := ParseDocument([]byte(content))
	if err != nil {
		return ""
	}

	return doc.GetSectionContent(sectionHeader)
}

// FindDeltaSectionContent extracts content from a delta section
// (ADDED, MODIFIED, etc.). This is a convenience wrapper around
// FindSectionContent for delta specs.
func FindDeltaSectionContent(content, deltaType string) string {
	return FindSectionContent(content, deltaType+" Requirements")
}

// FindRequirementsSection extracts the "## Requirements" section content.
// This is a convenience wrapper around FindSectionContent for spec files.
func FindRequirementsSection(content string) string {
	return FindSectionContent(content, "Requirements")
}

// FindAllH3Requirements finds all requirement headers in content and
// returns their names.
func FindAllH3Requirements(content string) []string {
	return FindRequirementHeaders(content)
}

// FindSectionIndex returns the start and end indices of a section header.
// Returns nil if the section is not found.
func FindSectionIndex(content, sectionHeader string) []int {
	targetHeader := "## " + sectionHeader
	lines := strings.Split(content, "\n")
	pos := 0

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == targetHeader {
			return []int{pos, pos + len(line)}
		}

		pos += len(line) + 1
	}

	return nil
}

// ExtractRequirementBlock extracts a requirement block (header + content)
// from content starting at a given line. Returns block content and end line.
func ExtractRequirementBlock(lines []string, startLine int) (string, int) {
	if startLine >= len(lines) {
		return "", startLine
	}

	var result strings.Builder
	result.WriteString(lines[startLine])
	result.WriteString("\n")

	endLine := startLine + 1

	for endLine < len(lines) {
		line := lines[endLine]
		trimmed := strings.TrimSpace(line)

		// Stop at next requirement header
		if _, ok := MatchH3Requirement(trimmed); ok {
			break
		}

		// Stop at next H2 section
		if IsH2Header(trimmed) {
			break
		}

		// Stop at non-requirement H3 header
		if IsH3Header(trimmed) && !strings.Contains(trimmed, "Requirement:") {
			break
		}

		result.WriteString(line)
		result.WriteString("\n")
		endLine++
	}

	return result.String(), endLine
}

// SplitIntoSections splits content into sections based on H2 headers.
// Returns a map of section name to section content.
func SplitIntoSections(content string) map[string]string {
	doc, err := ParseDocument([]byte(content))
	if err != nil {
		return make(map[string]string)
	}

	result := make(map[string]string, len(doc.Sections))
	for name, section := range doc.Sections {
		result[name] = strings.TrimSpace(section.Content)
	}

	return result
}
