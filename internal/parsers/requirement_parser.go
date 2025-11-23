//nolint:revive // line-length-limit - parsing logic needs clarity
package parsers

import (
	"fmt"
	"os"
	"strings"

	"github.com/connerohnesorge/spectr/internal/parser"
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
	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Parse with new parser
	doc, err := parser.Parse(string(content))
	if err != nil {
		return nil, fmt.Errorf("failed to parse spec file: %w", err)
	}

	// Extract requirements using new parser
	reqs, err := parser.ExtractRequirements(doc)
	if err != nil {
		return nil, fmt.Errorf("failed to extract requirements: %w", err)
	}

	// Convert to RequirementBlock format for backward compatibility
	var blocks []RequirementBlock
	for _, req := range reqs {
		blocks = append(blocks, requirementToBlock(req))
	}

	return blocks, nil
}

// ParseScenarios extracts scenario blocks from requirement content.
//
// Returns a slice of scenario names found in the requirement.
func ParseScenarios(requirementContent string) []string {
	// Parse the content
	doc, err := parser.Parse(requirementContent)
	if err != nil {
		// Return empty slice on parse error (maintain backward compatibility)
		return nil
	}

	// Extract requirements (which include scenarios)
	reqs, err := parser.ExtractRequirements(doc)
	if err != nil {
		// Return empty slice on extraction error
		return nil
	}

	// Collect all scenario names from all requirements
	var scenarios []string
	for _, req := range reqs {
		for _, scenario := range req.Scenarios {
			scenarios = append(scenarios, scenario.Name)
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

// requirementToBlock converts a parser.Requirement to a RequirementBlock.
//
// This adapter function maintains backward compatibility with existing code
// that expects the RequirementBlock structure. It reconstructs the raw
// markdown including scenario headers from the parsed data.
func requirementToBlock(req parser.Requirement) RequirementBlock {
	// Reconstruct the header line in the expected format
	headerLine := fmt.Sprintf("### Requirement: %s", req.Name)

	// Build raw content: header + content + scenarios
	var rawBuilder strings.Builder
	rawBuilder.WriteString(headerLine)
	rawBuilder.WriteString("\n")
	rawBuilder.WriteString(req.Content)

	// Append scenario headers and content
	for _, scenario := range req.Scenarios {
		rawBuilder.WriteString("\n\n#### Scenario: ")
		rawBuilder.WriteString(scenario.Name)
		rawBuilder.WriteString("\n")
		rawBuilder.WriteString(scenario.Content)
	}

	return RequirementBlock{
		HeaderLine: headerLine,
		Name:       req.Name,
		Raw:        rawBuilder.String(),
	}
}
