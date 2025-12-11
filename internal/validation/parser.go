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
func ExtractSections(content string) map[string]string {
	node := markdown.Parse([]byte(content))

	return markdown.ExtractH2Sections(node)
}

// ExtractRequirements returns all requirements found in content
// Looks for ### Requirement: headers
func ExtractRequirements(content string) []Requirement {
	mdReqs := markdown.ExtractRequirementsFromContent(content)

	// Convert markdown.RequirementBlock to validation.Requirement
	reqs := make([]Requirement, len(mdReqs))
	for i, mdReq := range mdReqs {
		// Extract content without the header line for compatibility
		contentWithoutHeader := mdReq.Raw
		if strings.HasPrefix(contentWithoutHeader, mdReq.HeaderLine) {
			contentWithoutHeader = strings.TrimPrefix(contentWithoutHeader, mdReq.HeaderLine)
			contentWithoutHeader = strings.TrimPrefix(contentWithoutHeader, "\n")
		}

		// Extract scenarios as full blocks (header + content)
		scenarios := markdown.ExtractScenarios(mdReq.Raw)

		reqs[i] = Requirement{
			Name:      mdReq.Name,
			Content:   strings.TrimSpace(contentWithoutHeader),
			Scenarios: scenarios,
		}
	}

	return reqs
}

// ExtractScenarios finds all #### Scenario: blocks in a requirement
func ExtractScenarios(requirementBlock string) []string {
	return markdown.ExtractScenarios(requirementBlock)
}

// ContainsShallOrMust checks if text contains SHALL or MUST (case-insensitive)
func ContainsShallOrMust(text string) bool {
	shallMustRegex := regexp.MustCompile(`(?i)\b(shall|must)\b`)

	return shallMustRegex.MatchString(text)
}

// NormalizeRequirementName normalizes requirement names for duplicate detection
// Trims whitespace, converts to lowercase, and removes extra spaces
func NormalizeRequirementName(name string) string {
	// Trim leading/trailing whitespace
	normalized := strings.TrimSpace(name)

	// Convert to lowercase
	normalized = strings.ToLower(normalized)

	// Replace multiple spaces with single space
	spaceRegex := regexp.MustCompile(`\s+`)
	normalized = spaceRegex.ReplaceAllString(normalized, " ")

	return normalized
}
