//nolint:revive // line-length-limit - parsing logic prioritizes clarity
package parsers

import (
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/mdparser"
)

// RequirementBlock represents a requirement with its header and content
type RequirementBlock struct {
	HeaderLine string // "### Requirement: <name>"
	Name       string // Extracted requirement name
	Raw        string // Full block content (header + scenarios + body text)
}

// ParseRequirements parses all requirement blocks from a spec file.
//
// Returns a slice of RequirementBlock with their names and full content.
// Uses mdparser and ExtractRequirements to parse the file.
func ParseRequirements(filePath string) ([]RequirementBlock, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := mdparser.Parse(string(content))
	if err != nil {
		return nil, err
	}

	// ExtractRequirements validates that each requirement has scenarios
	// But ParseRequirements doesn't enforce this, so we catch and ignore that error
	requirements, err := ExtractRequirements(doc)
	if err != nil {
		// For backward compatibility, ParseRequirements allows requirements without scenarios
		// So we manually extract them without the validation
		requirements = extractRequirementsWithoutValidation(doc)
	}

	return requirements, nil
}

// extractRequirementsWithoutValidation extracts requirements without enforcing scenario validation
func extractRequirementsWithoutValidation(doc *mdparser.Document) []RequirementBlock {
	var requirements []RequirementBlock

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

		requirements = append(requirements, RequirementBlock{
			HeaderLine: "### Requirement: " + name,
			Name:       name,
			Raw:        raw,
		})
	}

	return requirements
}

// ParseScenarios extracts scenario blocks from requirement content.
//
// Returns a slice of scenario names found in the requirement.
// Uses mdparser to parse the content and extract scenario headers.
func ParseScenarios(requirementContent string) []string {
	doc, err := mdparser.Parse(requirementContent)
	if err != nil {
		return nil
	}

	var scenarios []string

	// Traverse all nodes looking for H4 headers with "Scenario: " prefix
	for _, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != 4 {
			continue
		}

		if strings.HasPrefix(header.Text, "Scenario: ") {
			name := strings.TrimPrefix(header.Text, "Scenario: ")
			name = strings.TrimSpace(name)
			scenarios = append(scenarios, name)
		}
	}

	return scenarios
}

// NormalizeRequirementName normalizes requirement names for matching.
//
// Trims whitespace and converts to lowercase for case-insensitive comparison
func NormalizeRequirementName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
