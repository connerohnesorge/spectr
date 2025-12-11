//nolint:revive // line-length-limit - parsing logic needs clarity
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
func ParseRequirements(filePath string) ([]RequirementBlock, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	mdRequirements := markdown.ExtractRequirementsFromContent(string(content))

	// Convert markdown.RequirementBlock to parsers.RequirementBlock
	requirements := make([]RequirementBlock, len(mdRequirements))
	for i, mdReq := range mdRequirements {
		requirements[i] = RequirementBlock{
			HeaderLine: mdReq.HeaderLine,
			Name:       mdReq.Name,
			Raw:        mdReq.Raw,
		}
	}

	return requirements, nil
}

// ParseScenarios extracts scenario blocks from requirement content.
//
// Returns a slice of scenario names found in the requirement.
func ParseScenarios(requirementContent string) []string {
	return markdown.ExtractScenarioNames(requirementContent)
}

// NormalizeRequirementName normalizes requirement names for matching.
//
// Trims whitespace and converts to lowercase for case-insensitive comparison
func NormalizeRequirementName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}
