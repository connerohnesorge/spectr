//nolint:revive // line-length-limit - parsing logic prioritizes clarity
package validation

import (
	"bufio"
	"regexp"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parser"
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
	scanner := bufio.NewScanner(strings.NewReader(content))

	var currentSection string
	var currentContent strings.Builder
	sectionHeaderRegex := regexp.MustCompile(`^##\s+(.+)$`)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a section header (## header)
		matches := sectionHeaderRegex.FindStringSubmatch(line)
		if matches != nil {
			// Save previous section if exists
			if currentSection != "" {
				sections[currentSection] = strings.TrimSpace(
					currentContent.String(),
				)
			}

			// Start new section
			currentSection = strings.TrimSpace(matches[1])
			currentContent.Reset()
		} else if currentSection != "" {
			// Add line to current section content
			currentContent.WriteString(line)
			currentContent.WriteString("\n")
		}
	}

	// Save last section
	if currentSection != "" {
		sections[currentSection] = strings.TrimSpace(currentContent.String())
	}

	return sections
}

// ExtractRequirements returns all requirements found in content
// Looks for ### Requirement: headers
func ExtractRequirements(content string) []Requirement {
	// Initialize to empty slice instead of nil
	requirements := make([]Requirement, 0)
	scanner := bufio.NewScanner(strings.NewReader(content))

	requirementHeaderRegex := regexp.MustCompile(`^###\s+Requirement:\s*(.+)$`)
	var currentRequirement *Requirement
	var currentContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a requirement header
		matches := requirementHeaderRegex.FindStringSubmatch(line)
		if matches != nil {
			saveCurrentRequirement(
				currentRequirement,
				&currentContent,
				&requirements,
			)

			// Start new requirement
			currentRequirement = &Requirement{
				Name: strings.TrimSpace(matches[1]),
			}
			currentContent.Reset()

			continue
		}

		if currentRequirement == nil {
			continue
		}

		// Check if we should stop collecting
		if shouldStopRequirement(line) {
			closeRequirement(currentRequirement, &currentContent, &requirements)
			currentRequirement = nil

			continue
		}

		// Add line to current requirement content
		currentContent.WriteString(line)
		currentContent.WriteString(newline)
	}

	// Save last requirement
	saveCurrentRequirement(currentRequirement, &currentContent, &requirements)

	return requirements
}

// saveCurrentRequirement saves the current requirement if it exists
func saveCurrentRequirement(
	req *Requirement,
	content *strings.Builder,
	requirements *[]Requirement,
) {
	if req == nil {
		return
	}
	req.Content = strings.TrimSpace(content.String())
	req.Scenarios = ExtractScenarios(req.Content)
	*requirements = append(*requirements, *req)
}

// shouldStopRequirement checks if we should stop collecting requirement content
func shouldStopRequirement(line string) bool {
	// Stop if we hit another ### header (non-requirement)
	// But allow #### headers (scenarios) to pass through
	if strings.HasPrefix(line, "###") && !strings.HasPrefix(line, "####") {
		return true
	}

	// Stop if we hit a ## header (section boundary)
	// But make sure it's not a ### or #### header
	if strings.HasPrefix(line, "##") && !strings.HasPrefix(line, "###") {
		return true
	}

	return false
}

// closeRequirement finalizes and appends a requirement
func closeRequirement(
	req *Requirement,
	content *strings.Builder,
	requirements *[]Requirement,
) {
	req.Content = strings.TrimSpace(content.String())
	req.Scenarios = ExtractScenarios(req.Content)
	*requirements = append(*requirements, *req)
	content.Reset()
}

// ExtractScenarios finds all #### Scenario: blocks in a requirement
func ExtractScenarios(requirementBlock string) []string {
	// Initialize to empty slice instead of nil
	scenarios := make([]string, 0)
	scanner := bufio.NewScanner(strings.NewReader(requirementBlock))

	scenarioHeaderRegex := regexp.MustCompile(
		`^####\s+Scenario:\s*(.+)$`,
	)
	var currentScenario strings.Builder
	var inScenario bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a scenario header (#### Scenario:)
		matches := scenarioHeaderRegex.FindStringSubmatch(line)
		if matches != nil {
			// Save previous scenario if exists
			if inScenario {
				scenarios = append(scenarios,
					strings.TrimSpace(currentScenario.String()))
			}

			// Start new scenario with the header line
			currentScenario.Reset()
			currentScenario.WriteString(line)
			currentScenario.WriteString(newline)
			inScenario = true

			continue
		}

		// Process lines when we're inside a scenario
		if !inScenario {
			continue
		}

		// Check if we should stop collecting (hit header boundary)
		if shouldStopScenario(line) {
			closeScenario(&currentScenario, &scenarios)
			inScenario = false

			continue
		}

		// Add line to current scenario
		currentScenario.WriteString(line)
		currentScenario.WriteString(newline)
	}

	// Save last scenario
	if inScenario {
		scenarios = append(
			scenarios,
			strings.TrimSpace(currentScenario.String()),
		)
	}

	return scenarios
}

// shouldStopScenario checks if we should stop collecting scenario content
func shouldStopScenario(line string) bool {
	// Stop if we hit another #### header (next scenario or other)
	if strings.HasPrefix(line, "####") {
		return true
	}

	// Stop if we hit a ### header (next requirement)
	if strings.HasPrefix(line, "###") {
		return true
	}

	return false
}

// closeScenario finalizes and appends a scenario
func closeScenario(scenario *strings.Builder, scenarios *[]string) {
	*scenarios = append(*scenarios, strings.TrimSpace(scenario.String()))
	scenario.Reset()
}

// ContainsShallOrMust checks if text contains SHALL or MUST (case-insensitive)
// Uses the parser to skip code blocks and only check actual content
func ContainsShallOrMust(text string) bool {
	// Parse the document using the new parser
	doc, err := parser.Parse(text)
	if err != nil {
		// If parsing fails, fall back to simple string search
		// This maintains backward compatibility
		upper := strings.ToUpper(text)

		return containsWord(upper, "SHALL") || containsWord(upper, "MUST")
	}

	// Walk the AST and check only content nodes (skip code blocks)
	found := false
	parser.Walk(doc, func(node parser.Node) bool {
		switch n := node.(type) {
		case *parser.Paragraph:
			// Check paragraph text for SHALL/MUST
			upper := strings.ToUpper(n.Text)
			if containsWord(upper, "SHALL") || containsWord(upper, "MUST") {
				found = true

				return false // stop walking
			}
		case *parser.List:
			// Check list items for SHALL/MUST
			for _, item := range n.Items {
				upper := strings.ToUpper(item)
				if containsWord(upper, "SHALL") || containsWord(upper, "MUST") {
					found = true

					return false // stop walking
				}
			}
		case *parser.Header:
			// Also check headers (requirements can be stated in headers)
			upper := strings.ToUpper(n.Text)
			if containsWord(upper, "SHALL") || containsWord(upper, "MUST") {
				found = true

				return false // stop walking
			}
		}

		return true // continue walking
	})

	return found
}

// containsWord checks if a word exists as a whole word in text
// Both text and word should be uppercase for case-insensitive matching
func containsWord(text, word string) bool {
	// Simple word boundary check: the word must be surrounded by non-letter characters
	// or be at the start/end of the string
	start := 0
	for {
		idx := strings.Index(text[start:], word)
		if idx == -1 {
			return false
		}
		idx += start

		// Check left boundary
		leftOK := idx == 0 || !isLetter(rune(text[idx-1]))

		// Check right boundary
		rightIdx := idx + len(word)
		rightOK := rightIdx >= len(text) || !isLetter(rune(text[rightIdx]))

		if leftOK && rightOK {
			return true
		}

		// Continue searching after this match
		start = idx + 1
		if start >= len(text) {
			return false
		}
	}
}

// isLetter checks if a rune is a letter (simplified ASCII version)
func isLetter(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')
}

// NormalizeRequirementName normalizes requirement names for duplicate detection
// Trims whitespace, converts to lowercase, and removes extra spaces
func NormalizeRequirementName(name string) string {
	// Convert to lowercase first
	normalized := strings.ToLower(name)

	// Split on any whitespace (spaces, tabs, newlines) and filter out empty strings
	// Then rejoin with single spaces - this handles all whitespace normalization
	fields := strings.Fields(normalized)

	return strings.Join(fields, " ")
}
