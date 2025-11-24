//nolint:revive // line-length-limit - parsing logic prioritizes clarity
package validation

import (
	"strings"

	"github.com/connerohnesorge/spectr/internal/mdparser"
	"github.com/connerohnesorge/spectr/internal/parsers"
)

// Requirement represents a parsed requirement with its content and scenarios
type Requirement struct {
	Name      string
	Content   string
	Scenarios []string
}

// ExtractSections returns a map of section headers (## headers) to their content
// Example: "## Purpose" -> "This is the purpose..."
func ExtractSections(content string) map[string]string {
	sections := make(map[string]string)

	doc, err := mdparser.Parse(content)
	if err != nil {
		// Return empty map on parse error (graceful degradation)
		return sections
	}

	var currentSection string
	var currentContent strings.Builder

	for _, node := range doc.Children {
		// Check if this is an H2 header (section boundary)
		if header, ok := node.(*mdparser.Header); ok && header.Level == 2 {
			// Save previous section if exists
			if currentSection != "" {
				sections[currentSection] = strings.TrimSpace(currentContent.String())
			}

			// Start new section
			currentSection = strings.TrimSpace(header.Text)
			currentContent.Reset()

			continue
		}

		// Add node content to current section if we're in one
		if currentSection != "" {
			nodeContent := renderNodeToString(node)
			currentContent.WriteString(nodeContent)
		}
	}

	// Save last section
	if currentSection != "" {
		sections[currentSection] = strings.TrimSpace(currentContent.String())
	}

	return sections
}

// renderNodeToString renders an AST node back to markdown text
func renderNodeToString(node mdparser.Node) string {
	var sb strings.Builder

	switch n := node.(type) {
	case *mdparser.Header:
		sb.WriteString(strings.Repeat("#", n.Level))
		sb.WriteString(" ")
		sb.WriteString(n.Text)
		sb.WriteString("\n")

	case *mdparser.Paragraph:
		for _, line := range n.Lines {
			sb.WriteString(line)
			sb.WriteString("\n")
		}

	case *mdparser.CodeBlock:
		sb.WriteString("```")
		sb.WriteString(n.Language)
		sb.WriteString("\n")
		for _, line := range n.Lines {
			sb.WriteString(line)
			sb.WriteString("\n")
		}
		sb.WriteString("```\n")

	case *mdparser.List:
		for _, item := range n.Items {
			if n.Ordered {
				sb.WriteString("1. ")
			} else {
				sb.WriteString("- ")
			}
			sb.WriteString(item.Text)
			sb.WriteString("\n")
		}

	case *mdparser.BlankLine:
		for range n.Count {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// ExtractRequirements returns all requirements found in content
// Looks for ### Requirement: headers
// Uses the shared parsers.ExtractRequirements function but allows requirements without scenarios
func ExtractRequirements(content string) []Requirement {
	// Initialize to empty slice instead of nil
	requirements := make([]Requirement, 0)

	doc, err := mdparser.Parse(content)
	if err != nil {
		// Return empty slice on parse error (graceful degradation)
		return requirements
	}

	// Extract requirements manually from the document
	// We don't use parsers.ParseRequirements because it's file-based
	// and we need content-based extraction here
	parsedReqs := extractRequirementsManually(doc)

	// Convert from parsers.RequirementBlock to validation.Requirement
	for _, req := range parsedReqs {
		// Extract scenarios from the raw content
		scenarios := ExtractScenarios(req.Raw)

		requirements = append(requirements, Requirement{
			Name:      req.Name,
			Content:   strings.TrimSpace(req.Raw[len(req.HeaderLine)+1:]), // Skip header line
			Scenarios: scenarios,
		})
	}

	return requirements
}

// extractRequirementsManually extracts requirements without enforcing scenario validation
// This is used for backward compatibility with validation code
func extractRequirementsManually(doc *mdparser.Document) []parsers.RequirementBlock {
	var requirements []parsers.RequirementBlock

	for i, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != 3 {
			continue
		}

		// Check if this is a requirement header
		if !strings.HasPrefix(header.Text, "Requirement: ") {
			continue
		}

		// Extract requirement name
		name := strings.TrimPrefix(header.Text, "Requirement: ")
		name = strings.TrimSpace(name)

		// Get siblings (nodes after this header until next H3 or H2)
		siblings := getSiblingsUntilNextHeader(doc.Children, i+1, 3)

		// Build raw content
		raw := buildRawContent(header, siblings)

		requirements = append(requirements, parsers.RequirementBlock{
			HeaderLine: "### Requirement: " + name,
			Name:       name,
			Raw:        raw,
		})
	}

	return requirements
}

// Helper functions copied from parsers/extractor.go

// getSiblingsUntilNextHeader returns nodes from start index until next header of given level or higher.
func getSiblingsUntilNextHeader(nodes []mdparser.Node, startIdx, maxLevel int) []mdparser.Node {
	if startIdx >= len(nodes) {
		return nil
	}

	var siblings []mdparser.Node
	for i := startIdx; i < len(nodes); i++ {
		node := nodes[i]

		// Check if we've hit a header of the specified level or higher
		if header, ok := node.(*mdparser.Header); ok && header.Level <= maxLevel {
			break
		}

		siblings = append(siblings, node)
	}

	return siblings
}

// buildRawContent constructs the raw markdown text from a header and its siblings.
func buildRawContent(header *mdparser.Header, siblings []mdparser.Node) string {
	var sb strings.Builder

	// Add header line
	sb.WriteString("### Requirement: ")
	sb.WriteString(strings.TrimPrefix(header.Text, "Requirement: "))
	sb.WriteString("\n")

	// Add sibling content
	for _, node := range siblings {
		switch n := node.(type) {
		case *mdparser.Header:
			sb.WriteString(strings.Repeat("#", n.Level))
			sb.WriteString(" ")
			sb.WriteString(n.Text)
			sb.WriteString("\n")

		case *mdparser.Paragraph:
			for _, line := range n.Lines {
				sb.WriteString(line)
				sb.WriteString("\n")
			}

		case *mdparser.CodeBlock:
			sb.WriteString("```")
			sb.WriteString(n.Language)
			sb.WriteString("\n")
			for _, line := range n.Lines {
				sb.WriteString(line)
				sb.WriteString("\n")
			}
			sb.WriteString("```\n")

		case *mdparser.List:
			for _, item := range n.Items {
				if n.Ordered {
					sb.WriteString("1. ")
				} else {
					sb.WriteString("- ")
				}
				sb.WriteString(item.Text)
				sb.WriteString("\n")
			}

		case *mdparser.BlankLine:
			for range n.Count {
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

// ExtractScenarios finds all #### Scenario: blocks in a requirement
// Uses mdparser to parse scenario structure
func ExtractScenarios(requirementBlock string) []string {
	// Initialize to empty slice instead of nil
	scenarios := make([]string, 0)

	doc, err := mdparser.Parse(requirementBlock)
	if err != nil {
		// Return empty slice on parse error (graceful degradation)
		return scenarios
	}

	// Traverse the AST to find scenario headers
	for i, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != 4 {
			continue
		}

		// Check if this is a scenario header
		if !strings.HasPrefix(header.Text, "Scenario: ") {
			continue
		}

		// Get siblings until next H4 or higher
		siblings := getSiblingsUntilNextHeader(doc.Children, i+1, 4)

		// Build full scenario text (header + content)
		var sb strings.Builder
		sb.WriteString("#### Scenario: ")
		sb.WriteString(strings.TrimPrefix(header.Text, "Scenario: "))
		sb.WriteString("\n")

		// Add sibling content
		for _, sibling := range siblings {
			sb.WriteString(renderNodeToString(sibling))
		}

		scenarios = append(scenarios, strings.TrimSpace(sb.String()))
	}

	return scenarios
}

// ContainsShallOrMust checks if text contains SHALL or MUST (case-insensitive)
// Uses word boundary checking without regex
func ContainsShallOrMust(text string) bool {
	textLower := strings.ToLower(text)

	// Check for "shall" with word boundaries
	if containsWord(textLower, "shall") {
		return true
	}

	// Check for "must" with word boundaries
	if containsWord(textLower, "must") {
		return true
	}

	return false
}

// containsWord checks if text contains word with word boundaries
// A word boundary is defined as space, punctuation, or string start/end
func containsWord(text, word string) bool {
	idx := strings.Index(text, word)
	for idx != -1 {
		// Check if this is a word boundary (not part of a larger word)
		wordLen := len(word)

		// Check character before
		beforeOK := idx == 0 || !isAlphanumeric(rune(text[idx-1]))

		// Check character after
		afterOK := idx+wordLen >= len(text) || !isAlphanumeric(rune(text[idx+wordLen]))

		if beforeOK && afterOK {
			return true
		}

		// Continue searching
		nextSearchStart := idx + 1
		nextIdx := strings.Index(text[nextSearchStart:], word)
		if nextIdx == -1 {
			break
		}

		// Convert relative index to absolute index within the original text
		idx = nextSearchStart + nextIdx
	}

	return false
}

// isAlphanumeric checks if a character is alphanumeric
func isAlphanumeric(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')
}

// NormalizeRequirementName normalizes requirement names for duplicate detection
// Trims whitespace, converts to lowercase, and removes extra spaces
func NormalizeRequirementName(name string) string {
	// Trim leading/trailing whitespace
	normalized := strings.TrimSpace(name)

	// Convert to lowercase
	normalized = strings.ToLower(normalized)

	// Replace multiple spaces with single space (without regex)
	normalized = collapseSpaces(normalized)

	return normalized
}

// collapseSpaces replaces multiple consecutive spaces with a single space
func collapseSpaces(s string) string {
	var result strings.Builder
	prevSpace := false

	for _, r := range s {
		isSpace := r == ' ' || r == '\t' || r == '\n' || r == '\r'

		if isSpace {
			if !prevSpace {
				result.WriteRune(' ')
				prevSpace = true
			}
		} else {
			result.WriteRune(r)
			prevSpace = false
		}
	}

	return result.String()
}
