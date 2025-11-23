//nolint:revive // line-length-limit - parsing logic prioritizes clarity
package validation

import (
	"strings"

	"github.com/connerohnesorge/spectr/internal/syntax"
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
	doc, err := syntax.Parse(content)
	if err != nil {
		return make(map[string]string)
	}

	sections := make(map[string]string)
	var currentSection string
	var currentContent strings.Builder

	for _, node := range doc.Nodes {
		if h, ok := node.(*syntax.Header); ok && h.Level == 2 {
			if currentSection != "" {
				sections[currentSection] = strings.TrimSpace(currentContent.String())
			}
			currentSection = strings.TrimSpace(h.Text)
			currentContent.Reset()
		} else if currentSection != "" {
			currentContent.WriteString(node.RawString())
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
	doc, err := syntax.Parse(content)
	if err != nil {
		return []Requirement{}
	}

	var requirements []Requirement
	var currentReq *Requirement
	var currentContent strings.Builder

	for _, node := range doc.Nodes {
		if h, ok := node.(*syntax.Header); ok {
			if h.Level == 3 && strings.HasPrefix(h.Text, "Requirement:") {
				// Save previous requirement
				if currentReq != nil {
					currentReq.Content = strings.TrimSpace(currentContent.String())
					currentReq.Scenarios = ExtractScenarios(currentReq.Content)
					requirements = append(requirements, *currentReq)
				}

				// Start new requirement
				name := strings.TrimSpace(strings.TrimPrefix(h.Text, "Requirement:"))
				currentReq = &Requirement{
					Name: name,
				}
				currentContent.Reset()
				continue
			} else if h.Level == 2 || h.Level == 3 {
				// Section boundary (##) or other Level 3 header ends requirement
				if currentReq != nil {
					currentReq.Content = strings.TrimSpace(currentContent.String())
					currentReq.Scenarios = ExtractScenarios(currentReq.Content)
					requirements = append(requirements, *currentReq)
					currentReq = nil
				}
				// Don't continue, as this header might be start of something else (but we only care about Requirements here)
				// Actually, if it's Level 2, it's a section header, we ignore it in this context (we are parsing requirements from a block).
				// If it's another Level 3, it's ignored unless it's a Requirement.
				continue
			}
		}

		if currentReq != nil {
			currentContent.WriteString(node.RawString())
		}
	}

	// Save last requirement
	if currentReq != nil {
		currentReq.Content = strings.TrimSpace(currentContent.String())
		currentReq.Scenarios = ExtractScenarios(currentReq.Content)
		requirements = append(requirements, *currentReq)
	}

	return requirements
}

// ExtractScenarios finds all #### Scenario: blocks in a requirement
func ExtractScenarios(requirementBlock string) []string {
	doc, err := syntax.Parse(requirementBlock)
	if err != nil {
		return []string{}
	}

	var scenarios []string
	var currentScenario strings.Builder
	var inScenario bool

	for _, node := range doc.Nodes {
		if h, ok := node.(*syntax.Header); ok {
			if h.Level == 4 && strings.HasPrefix(h.Text, "Scenario:") {
				// Save previous scenario
				if inScenario {
					scenarios = append(scenarios, strings.TrimSpace(currentScenario.String()))
				}
				// Start new scenario
				inScenario = true
				currentScenario.Reset()
				currentScenario.WriteString(h.Raw) // Include header
				continue
			} else if h.Level <= 3 {
				// Requirement or Section boundary ends scenario
				if inScenario {
					scenarios = append(scenarios, strings.TrimSpace(currentScenario.String()))
					inScenario = false
				}
				continue
			} else if h.Level == 4 {
				// Another Level 4 header (not Scenario) ends scenario
				if inScenario {
					scenarios = append(scenarios, strings.TrimSpace(currentScenario.String()))
					inScenario = false
				}
				continue
			}
		}

		if inScenario {
			currentScenario.WriteString(node.RawString())
		}
	}

	// Save last scenario
	if inScenario {
		scenarios = append(scenarios, strings.TrimSpace(currentScenario.String()))
	}

	return scenarios
}

// ContainsShallOrMust checks if text contains SHALL or MUST (case-insensitive)
func ContainsShallOrMust(text string) bool {
	f := func(c rune) bool {
		return !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9'))
	}
	words := strings.FieldsFunc(text, f)
	for _, w := range words {
		upperW := strings.ToUpper(w)
		if upperW == "SHALL" || upperW == "MUST" {
			return true
		}
	}
	return false
}

// NormalizeRequirementName normalizes requirement names for duplicate detection
// Trims whitespace, converts to lowercase, and removes extra spaces
func NormalizeRequirementName(name string) string {
	return strings.Join(strings.Fields(strings.ToLower(name)), " ")
}
