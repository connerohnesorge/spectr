//nolint:revive // line-length-limit - regex patterns and parsing logic need clarity
package parsers

import (
	"bufio"
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
// Uses AST-based parsing via markdown.ParseDocument for accurate extraction.
func ParseRequirements(
	filePath string,
) ([]RequirementBlock, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	doc, err := markdown.ParseDocument(content)
	if err != nil {
		return nil, err
	}

	// Get requirement names in document order
	names := doc.GetRequirementNames()
	requirements := make([]RequirementBlock, 0, len(names))

	for _, name := range names {
		req := doc.GetRequirement(name)
		if req == nil {
			continue
		}

		headerLine := "### Requirement: " + req.Name
		raw := headerLine + "\n" + req.Content

		requirements = append(requirements, RequirementBlock{
			HeaderLine: headerLine,
			Name:       req.Name,
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
	var scenarios []string

	scanner := bufio.NewScanner(
		strings.NewReader(requirementContent),
	)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if name, ok := markdown.MatchH4Scenario(line); ok {
			scenarios = append(
				scenarios,
				strings.TrimSpace(name),
			)
		}
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
