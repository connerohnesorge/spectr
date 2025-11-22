//nolint:revive // line-length-limit - regex patterns and parsing logic need clarity
package parsers

import (
	"bufio"
	"os"
	"regexp"
	"strings"
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
func ParseRequirements(filePath string) ([]RequirementBlock, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var requirements []RequirementBlock
	var currentReq *RequirementBlock
	reqPattern := regexp.MustCompile(`^###\s+Requirement:\s*(.+)$`)
	h2Pattern := regexp.MustCompile(`^##\s+`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Check if this is a new requirement header
		matches := reqPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			// Save previous requirement if exists
			if currentReq != nil {
				requirements = append(requirements, *currentReq)
			}

			// Start new requirement
			name := strings.TrimSpace(matches[1])
			currentReq = &RequirementBlock{
				HeaderLine: line,
				Name:       name,
				Raw:        line + "\n",
			}

			continue
		}

		// Check if we hit a new section (## header) - ends current requirement
		if h2Pattern.MatchString(line) {
			if currentReq != nil {
				requirements = append(requirements, *currentReq)
				currentReq = nil
			}

			continue
		}

		// Append line to current requirement if we're in one
		if currentReq != nil {
			currentReq.Raw += line + "\n"
		}
	}

	// Don't forget the last requirement
	if currentReq != nil {
		requirements = append(requirements, *currentReq)
	}

	return requirements, scanner.Err()
}

// ParseScenarios extracts scenario blocks from requirement content.
//
// Returns a slice of scenario names found in the requirement.
func ParseScenarios(requirementContent string) []string {
	var scenarios []string
	scenarioPattern := regexp.MustCompile(`^####\s+Scenario:\s*(.+)$`)

	scanner := bufio.NewScanner(strings.NewReader(requirementContent))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		matches := scenarioPattern.FindStringSubmatch(line)
		if len(matches) > 1 {
			scenarios = append(
				scenarios,
				strings.TrimSpace(matches[1]),
			)
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
