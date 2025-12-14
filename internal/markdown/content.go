// Package markdown provides AST-based markdown parsing for spectr.
//
// This file contains content building functions that populate the Content
// fields of Section, Requirement, and Scenario structs after the AST walk.
// These functions compute line ranges and extract raw markdown text.
package markdown

import (
	"strings"
)

// Content building functions (used by ParseDocument)
//
// These functions are called after the AST walk to populate Content fields.
// They use line number information from headers to compute content ranges.

// buildContent builds section and requirement content from line ranges.
func buildContent(doc *Document) {
	buildSectionContent(doc)
	buildRequirementContent(doc)
	buildScenarioContent(doc)
}

// buildSectionContent populates Content for all sections.
func buildSectionContent(doc *Document) {
	for i, h := range doc.H2Headers {
		section := doc.Sections[h.Text]
		if section == nil {
			continue
		}

		endLine := len(doc.Lines)
		if i+1 < len(doc.H2Headers) {
			endLine = doc.H2Headers[i+1].Line - 1
		}

		section.EndLine = endLine
		startLine := section.StartLine - 1
		section.Content = buildContentFromLines(doc.Lines, startLine, endLine)
	}
}

// buildRequirementContent populates Content for all requirements.
func buildRequirementContent(doc *Document) {
	for i, h := range doc.H3Headers {
		name := parseRequirementName(h.Text)
		if name == "" {
			continue
		}

		req := doc.Requirements[name]
		if req == nil {
			continue
		}

		endLine := findEndLine(doc.H3Headers, i+1, len(doc.Lines))
		endLine = findEarlierBoundary(h.Line, doc.H2Headers, endLine)
		req.Content = buildContentFromLines(doc.Lines, h.Line, endLine)
	}
}

// buildScenarioContent populates Content for all scenarios.
func buildScenarioContent(doc *Document) {
	for i, h := range doc.H4Headers {
		name := parseScenarioName(h.Text)
		if name == "" {
			continue
		}

		scenario := doc.Scenarios[name]
		if scenario == nil {
			continue
		}

		endLine := findEndLine(doc.H4Headers, i+1, len(doc.Lines))
		endLine = findEarlierBoundary(h.Line, doc.H3Headers, endLine)
		endLine = findEarlierBoundary(h.Line, doc.H2Headers, endLine)
		scenario.Content = buildContentFromLines(doc.Lines, h.Line, endLine)
	}
}

// findEndLine returns the end line based on the next header.
func findEndLine(headers []Header, nextIdx, maxLine int) int {
	if nextIdx < len(headers) {
		return headers[nextIdx].Line - 1
	}

	return maxLine
}

// findEarlierBoundary finds if any header starts before endLine.
func findEarlierBoundary(startLine int, headers []Header, endLine int) int {
	for _, h := range headers {
		if h.Line > startLine && h.Line < endLine {
			return h.Line - 1
		}
	}

	return endLine
}

// buildContentFromLines builds content string from line range.
func buildContentFromLines(lines []string, startLine, endLine int) string {
	if startLine >= len(lines) {
		return ""
	}

	end := endLine
	if end > len(lines) {
		end = len(lines)
	}

	var buf strings.Builder

	for i := startLine; i < end; i++ {
		buf.WriteString(lines[i])
		if i < end-1 {
			buf.WriteString("\n")
		}
	}

	return buf.String()
}
