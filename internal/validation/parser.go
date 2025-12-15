//nolint:revive // line-length-limit - parsing logic prioritizes clarity
package validation

import (
	"bufio"
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
func ExtractSections(
	content string,
) map[string]string {
	sections := make(map[string]string)
	scanner := bufio.NewScanner(
		strings.NewReader(content),
	)

	var currentSection string
	var currentContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a section header (## header)
		if sectionName, ok := markdown.MatchH2SectionHeader(line); ok {
			// Save previous section if exists
			if currentSection != "" {
				sections[currentSection] = strings.TrimSpace(
					currentContent.String(),
				)
			}

			// Start new section
			currentSection = strings.TrimSpace(
				sectionName,
			)
			currentContent.Reset()
		} else if currentSection != "" {
			// Add line to current section content
			currentContent.WriteString(line)
			currentContent.WriteString("\n")
		}
	}

	// Save last section
	if currentSection != "" {
		sections[currentSection] = strings.TrimSpace(
			currentContent.String(),
		)
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
	scanner := bufio.NewScanner(
		strings.NewReader(content),
	)

	var currentRequirement *Requirement
	var currentContent strings.Builder

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a requirement header
		if reqName, ok := markdown.MatchRequirementHeader(line); ok {
			saveCurrentRequirement(
				currentRequirement,
				&currentContent,
				&requirements,
			)

			// Start new requirement
			currentRequirement = &Requirement{
				Name: strings.TrimSpace(
					reqName,
				),
			}
			currentContent.Reset()

			continue
		}

		if currentRequirement == nil {
			continue
		}

		// Check if we should stop collecting
		if shouldStopRequirement(line) {
			closeRequirement(
				currentRequirement,
				&currentContent,
				&requirements,
			)
			currentRequirement = nil

			continue
		}

		// Add line to current requirement content
		currentContent.WriteString(line)
		currentContent.WriteString(newline)
	}

	// Save last requirement
	saveCurrentRequirement(
		currentRequirement,
		&currentContent,
		&requirements,
	)

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
	req.Content = strings.TrimSpace(
		content.String(),
	)
	req.Scenarios = ExtractScenarios(req.Content)
	*requirements = append(*requirements, *req)
}

// shouldStopRequirement checks if we should stop collecting requirement content
func shouldStopRequirement(line string) bool {
	// Stop if we hit another ### header (non-requirement)
	// But allow #### headers (scenarios) to pass through
	if strings.HasPrefix(line, "###") &&
		!strings.HasPrefix(line, "####") {
		return true
	}

	// Stop if we hit a ## header (section boundary)
	// But make sure it's not a ### or #### header
	if strings.HasPrefix(line, "##") &&
		!strings.HasPrefix(line, "###") {
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
	req.Content = strings.TrimSpace(
		content.String(),
	)
	req.Scenarios = ExtractScenarios(req.Content)
	*requirements = append(*requirements, *req)
	content.Reset()
}

// ExtractScenarios finds all #### Scenario: blocks in a requirement
func ExtractScenarios(
	requirementBlock string,
) []string {
	// Initialize to empty slice instead of nil
	scenarios := make([]string, 0)
	scanner := bufio.NewScanner(
		strings.NewReader(requirementBlock),
	)

	var currentScenario strings.Builder
	var inScenario bool

	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a scenario header (#### Scenario:)
		if _, ok := markdown.MatchScenarioHeader(line); ok {
			// Save previous scenario if exists
			if inScenario {
				scenarios = append(
					scenarios,
					strings.TrimSpace(
						currentScenario.String(),
					),
				)
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
			closeScenario(
				&currentScenario,
				&scenarios,
			)
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
			strings.TrimSpace(
				currentScenario.String(),
			),
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
func closeScenario(
	scenario *strings.Builder,
	scenarios *[]string,
) {
	*scenarios = append(
		*scenarios,
		strings.TrimSpace(scenario.String()),
	)
	scenario.Reset()
}

// ContainsShallOrMust checks if text contains SHALL or MUST (case-insensitive)
// Uses string-based matching instead of regex for better performance
func ContainsShallOrMust(text string) bool {
	// Convert to lowercase for case-insensitive comparison
	lower := strings.ToLower(text)

	// Check for "shall" or "must" as whole words
	// We need to ensure they are word boundaries (not part of another word)
	return containsWord(lower, "shall") ||
		containsWord(lower, "must")
}

// containsWord checks if a word exists as a whole word in the text
// (not as part of a larger word)
func containsWord(text, word string) bool {
	wordLen := len(word)
	textLen := len(text)

	idx := 0
	for {
		// Find the next occurrence of the word
		pos := strings.Index(text[idx:], word)
		if pos == -1 {
			return false
		}

		// Calculate absolute position
		absPos := idx + pos

		// Check if it's a word boundary at the start
		startOK := absPos == 0 ||
			!isWordChar(text[absPos-1])

		// Check if it's a word boundary at the end
		endPos := absPos + wordLen
		endOK := endPos >= textLen ||
			!isWordChar(text[endPos])

		if startOK && endOK {
			return true
		}

		// Move past this occurrence and continue searching
		idx = absPos + 1
		if idx >= textLen {
			return false
		}
	}
}

// isWordChar returns true if the byte is a word character (alphanumeric or underscore)
func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') ||
		(b >= 'A' && b <= 'Z') ||
		(b >= '0' && b <= '9') ||
		b == '_'
}

// NormalizeRequirementName normalizes requirement names for duplicate detection
// Trims whitespace, converts to lowercase, and removes extra spaces
// Uses string-based operations instead of regex for better performance
func NormalizeRequirementName(
	name string,
) string {
	// Trim leading/trailing whitespace
	normalized := strings.TrimSpace(name)

	// Convert to lowercase
	normalized = strings.ToLower(normalized)

	// Replace multiple spaces with single space using strings.Fields and Join
	// strings.Fields splits on any whitespace and removes empty strings
	fields := strings.Fields(normalized)
	normalized = strings.Join(fields, " ")

	return normalized
}
