//nolint:revive // line-length-limit - parsing logic need clarity
package parsers

import (
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/markdown"
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
//
//nolint:revive // function-length - parser is clearest as single function
func ParseRequirements(
	filePath string,
) ([]RequirementBlock, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := markdown.ParseDocument(content)
	if err != nil {
		// Return empty slice if content is empty or invalid
		return []RequirementBlock{}, nil
	}

	// Count requirement headers for pre-allocation
	reqCount := 0
	for _, header := range doc.Headers {
		if header.Level == 3 && strings.HasPrefix(header.Text, "Requirement:") {
			reqCount++
		}
	}

	requirements := make([]RequirementBlock, 0, reqCount)

	// Find all H3 headers that are requirements
	for _, header := range doc.Headers {
		if header.Level != 3 || !strings.HasPrefix(header.Text, "Requirement:") {
			continue
		}

		// Extract requirement name (everything after "Requirement:")
		name := strings.TrimPrefix(header.Text, "Requirement:")
		name = strings.TrimSpace(name)

		// Build the header line
		headerLine := "### " + header.Text

		// Get section content from the Sections map
		sectionContent := ""
		if section, ok := doc.Sections[header.Text]; ok {
			sectionContent = section.Content
		}

		// Build the raw content (header + content)
		raw := headerLine + "\n"
		if sectionContent != "" {
			raw += sectionContent + "\n"
		}

		requirements = append(requirements, RequirementBlock{
			HeaderLine: headerLine,
			Name:       name,
			Raw:        raw,
		})
	}

	return requirements, nil
}

// ParseScenarios extracts scenario blocks from requirement content.
//
// Returns a slice of scenario names found in the requirement.
func ParseScenarios(
	requirementContent string,
) []string {
	doc, err := markdown.ParseDocument([]byte(requirementContent))
	if err != nil {
		// Return empty slice if content is empty or invalid
		return []string{}
	}

	// Count scenario headers for pre-allocation
	scenarioCount := 0
	for _, header := range doc.Headers {
		if header.Level == 4 && strings.HasPrefix(header.Text, "Scenario:") {
			scenarioCount++
		}
	}

	scenarios := make([]string, 0, scenarioCount)

	// Find all H4 headers that are scenarios
	for _, header := range doc.Headers {
		if header.Level != 4 || !strings.HasPrefix(header.Text, "Scenario:") {
			continue
		}

		// Extract scenario name (everything after "Scenario:")
		name := strings.TrimPrefix(header.Text, "Scenario:")
		name = strings.TrimSpace(name)
		scenarios = append(scenarios, name)
	}

	return scenarios
}

// NormalizeRequirementName normalizes requirement names for matching.
//
// Trims whitespace and converts to lowercase for case-insensitive comparison
func NormalizeRequirementName(
	name string,
) string {
	return strings.ToLower(
		strings.TrimSpace(name),
	)
}
