//nolint:revive // line-length-limit - parsing logic prioritizes clarity
package validation

import (
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/markdown"
)

// Requirement represents a parsed requirement with its content and scenarios
type Requirement struct {
	Name      string
	Content   string
	Scenarios []string
}

// ExtractSections returns a map of section headers (## headers) to their content
// Example: "## Purpose" -> "This is the purpose..."
// Note: Section content includes everything until the next H2 header or EOF,
// including any nested H3, H4 headers and their content.
func ExtractSections(
	content string,
) map[string]string {
	sections := make(map[string]string)

	// Handle empty content
	if strings.TrimSpace(content) == "" {
		return sections
	}

	// Parse the document using the markdown package
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		// Return empty map for invalid content
		return sections
	}

	// Find H2 headers and extract content between them
	lines := strings.Split(content, "\n")

	// Collect H2 headers with their indices
	type h2Header struct {
		text string
		line int
	}
	var h2Headers []h2Header
	for _, header := range doc.Headers {
		if header.Level == 2 {
			h2Headers = append(h2Headers, h2Header{
				text: header.Text,
				line: header.Line,
			})
		}
	}

	// Extract content for each H2 section
	for i, h2 := range h2Headers {
		// Content starts on the line after the header
		startLine := h2.line // 1-indexed, so this is actually the next line (0-indexed)

		// Content ends at the next H2 header or EOF
		endLine := len(lines)
		if i+1 < len(h2Headers) {
			endLine = h2Headers[i+1].line - 1 // Line before the next H2
		}

		// Extract content lines
		var contentLines []string
		for lineNum := startLine; lineNum < endLine && lineNum < len(lines); lineNum++ {
			contentLines = append(contentLines, lines[lineNum])
		}

		sections[h2.text] = strings.TrimSpace(strings.Join(contentLines, newline))
	}

	return sections
}

// ExtractRequirements returns all requirements found in content
// Looks for ### Requirement: headers
func ExtractRequirements(
	content string,
) []Requirement {
	// Initialize to empty slice instead of nil
	requirements := make([]Requirement, 0)

	// Handle empty content
	if strings.TrimSpace(content) == "" {
		return requirements
	}

	// Parse the document using the markdown package
	doc, err := markdown.ParseDocument([]byte(content))
	if err != nil {
		// Return empty slice for invalid content
		return requirements
	}

	// Find all H3 headers that match "Requirement: ..."
	lines := strings.Split(content, "\n")

	for i, header := range doc.Headers {
		if header.Level != 3 {
			continue
		}

		// Check if header matches "Requirement: ..." pattern
		if !strings.HasPrefix(header.Text, "Requirement:") {
			continue
		}

		// Extract the requirement name
		name := strings.TrimSpace(
			strings.TrimPrefix(header.Text, "Requirement:"),
		)

		// Find the content between this header and the next stopping point
		reqContent := extractRequirementContent(
			doc,
			header,
			i,
			lines,
		)

		req := Requirement{
			Name:      name,
			Content:   reqContent,
			Scenarios: ExtractScenarios(reqContent),
		}
		requirements = append(requirements, req)
	}

	return requirements
}

// extractRequirementContent extracts content for a requirement header
func extractRequirementContent(
	doc *markdown.Document,
	header markdown.Header,
	headerIndex int,
	lines []string,
) string {
	// Content starts on the line after the header
	startLine := header.Line // 1-indexed line number, converted to 0-indexed below

	// Find the end line: next ## or ### header (but not #### which is scenario)
	endLine := len(lines)
	for j := headerIndex + 1; j < len(doc.Headers); j++ {
		nextHeader := doc.Headers[j]
		// Stop at ## (section) or ### (next requirement or other H3)
		if nextHeader.Level <= 3 {
			endLine = nextHeader.Line - 1

			break
		}
	}

	// Extract content lines (startLine is 1-indexed, so it's already the next line in 0-indexed)
	var contentLines []string
	for lineNum := startLine; lineNum < endLine && lineNum < len(lines); lineNum++ {
		contentLines = append(contentLines, lines[lineNum])
	}

	return strings.TrimSpace(strings.Join(contentLines, newline))
}

// ExtractScenarios finds all #### Scenario: blocks in a requirement
func ExtractScenarios(
	requirementBlock string,
) []string {
	// Initialize to empty slice instead of nil
	scenarios := make([]string, 0)

	// Handle empty content
	if strings.TrimSpace(requirementBlock) == "" {
		return scenarios
	}

	// Parse the document using the markdown package
	doc, err := markdown.ParseDocument([]byte(requirementBlock))
	if err != nil {
		// Return empty slice for invalid content
		return scenarios
	}

	lines := strings.Split(requirementBlock, "\n")

	// Find all H4 headers that match "Scenario: ..."
	for i, header := range doc.Headers {
		if header.Level != 4 {
			continue
		}

		// Check if header matches "Scenario: ..." pattern
		if !strings.HasPrefix(header.Text, "Scenario:") {
			continue
		}

		// Find the content for this scenario
		scenarioContent := extractScenarioContent(
			doc,
			header,
			i,
			lines,
		)

		scenarios = append(scenarios, scenarioContent)
	}

	return scenarios
}

// extractScenarioContent extracts content for a scenario header
func extractScenarioContent(
	doc *markdown.Document,
	header markdown.Header,
	headerIndex int,
	lines []string,
) string {
	// Scenario content includes the header line itself
	startLine := header.Line - 1 // Convert to 0-indexed

	// Find the end line: next ### or #### header
	endLine := len(lines)
	for j := headerIndex + 1; j < len(doc.Headers); j++ {
		nextHeader := doc.Headers[j]
		// Stop at ### (requirement) or #### (next scenario)
		if nextHeader.Level <= 4 {
			endLine = nextHeader.Line - 1

			break
		}
	}

	// Extract content lines (including the header line)
	var contentLines []string
	for lineNum := startLine; lineNum < endLine && lineNum < len(lines); lineNum++ {
		contentLines = append(contentLines, lines[lineNum])
	}

	return strings.TrimSpace(strings.Join(contentLines, newline))
}

// ContainsShallOrMust checks if text contains SHALL or MUST (case-insensitive)
func ContainsShallOrMust(text string) bool {
	shallMustRegex := regexp.MustCompile(
		`(?i)\b(shall|must)\b`,
	)

	return shallMustRegex.MatchString(text)
}

// NormalizeRequirementName normalizes requirement names for duplicate detection
// Trims whitespace, converts to lowercase, and removes extra spaces
func NormalizeRequirementName(
	name string,
) string {
	// Trim leading/trailing whitespace
	normalized := strings.TrimSpace(name)

	// Convert to lowercase
	normalized = strings.ToLower(normalized)

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized = spaceRegex.ReplaceAllString(
		normalized,
		" ",
	)

	return normalized
}
