//nolint:revive // file-length-limit, line-length-limit - section parsing is complex
package markdown

import (
	"bufio"
	"strings"

	bf "github.com/russross/blackfriday/v2"
)

// ExtractSectionContent extracts the content under a specific header.
// The headerName should match the header text exactly (case-sensitive).
// Content is extracted from after the header until the next header of
// the same or higher level.
//
//nolint:revive // cognitive-complexity - AST walking is inherently complex
func ExtractSectionContent(node *bf.Node, headerName string, level int) string {
	if node == nil {
		return ""
	}

	var content strings.Builder
	var inSection bool

	for child := node.FirstChild; child != nil; child = child.Next {
		if child.Type == bf.Heading {
			headingLevel := child.Level
			headingText := strings.TrimSpace(extractText(child))

			if !inSection {
				// Look for our target header
				if headingLevel == level && headingText == headerName {
					inSection = true

					continue
				}
			} else {
				// Check if we should stop (same or higher level header)
				if headingLevel <= level {
					break
				}
				// Include lower-level headers in content
				content.WriteString(renderNode(child))
			}

			continue
		}

		if inSection {
			content.WriteString(renderNode(child))
		}
	}

	return strings.TrimSpace(content.String())
}

// ExtractSectionContentAfterNode extracts content starting from a specific node
// until the next header of the same or higher level.
func ExtractSectionContentAfterNode(startNode *bf.Node, level int) string {
	if startNode == nil {
		return ""
	}

	var content strings.Builder

	// Start from the next sibling after the header
	for node := startNode.Next; node != nil; node = node.Next {
		if node.Type == bf.Heading {
			headingLevel := node.Level
			// Stop at same or higher level header
			if headingLevel <= level {
				break
			}
		}
		content.WriteString(renderNode(node))
	}

	return strings.TrimSpace(content.String())
}

// ExtractRequirements extracts all requirement blocks from the markdown AST.
// Requirements are identified by "### Requirement:" headers.
// Returns requirement blocks in document order.
//
//nolint:revive // cognitive-complexity - requirement parsing is complex
func ExtractRequirements(node *bf.Node) []RequirementBlock {
	requirements := make([]RequirementBlock, 0)
	if node == nil {
		return requirements
	}

	var currentReq *RequirementBlock
	var contentBuilder strings.Builder

	for child := node.FirstChild; child != nil; child = child.Next {
		if child.Type == bf.Heading {
			headingLevel := child.Level
			headingText := strings.TrimSpace(extractText(child))

			// Check for requirement header (### Requirement:)
			if headingLevel == 3 && strings.HasPrefix(headingText, "Requirement:") {
				// Save previous requirement
				if currentReq != nil {
					currentReq.Raw = strings.TrimRight(contentBuilder.String(), "\n") + "\n"
					currentReq.Scenarios = ExtractScenarioNames(currentReq.Raw)
					requirements = append(requirements, *currentReq)
				}

				// Start new requirement
				name := strings.TrimSpace(strings.TrimPrefix(headingText, "Requirement:"))
				headerLine := "### Requirement: " + name
				currentReq = &RequirementBlock{
					Name:       name,
					HeaderLine: headerLine,
				}
				contentBuilder.Reset()
				contentBuilder.WriteString(headerLine)
				contentBuilder.WriteString("\n")

				continue
			}

			// Check for H2 or H3 non-requirement header (ends current requirement)
			if headingLevel <= 3 && currentReq != nil {
				// Check if it's a non-requirement H3
				if headingLevel == 3 && !strings.HasPrefix(headingText, "Requirement:") {
					// This is another H3 that's not a requirement - end current
					currentReq.Raw = strings.TrimRight(contentBuilder.String(), "\n") + "\n"
					currentReq.Scenarios = ExtractScenarioNames(currentReq.Raw)
					requirements = append(requirements, *currentReq)
					currentReq = nil
					contentBuilder.Reset()

					continue
				}
				// H2 header ends current requirement
				if headingLevel == 2 {
					currentReq.Raw = strings.TrimRight(contentBuilder.String(), "\n") + "\n"
					currentReq.Scenarios = ExtractScenarioNames(currentReq.Raw)
					requirements = append(requirements, *currentReq)
					currentReq = nil
					contentBuilder.Reset()

					continue
				}
			}
		}

		// Append content to current requirement
		if currentReq != nil {
			contentBuilder.WriteString(renderNode(child))
		}
	}

	// Save final requirement
	if currentReq != nil {
		currentReq.Raw = strings.TrimRight(contentBuilder.String(), "\n") + "\n"
		currentReq.Scenarios = ExtractScenarioNames(currentReq.Raw)
		requirements = append(requirements, *currentReq)
	}

	return requirements
}

// ExtractRequirementsFromContent parses content string and extracts
// requirements. This is useful when you have section content as a string.
func ExtractRequirementsFromContent(content string) []RequirementBlock {
	node := Parse([]byte(content))

	return ExtractRequirements(node)
}

// ExtractScenarios extracts scenario blocks from requirement content.
// Returns full scenario blocks (header + content).
func ExtractScenarios(content string) []string {
	scenarios := make([]string, 0)

	// Parse the content
	node := Parse([]byte(content))
	if node == nil {
		return scenarios
	}

	var currentScenario strings.Builder
	var inScenario bool

	for child := node.FirstChild; child != nil; child = child.Next {
		if child.Type == bf.Heading {
			headingLevel := child.Level
			headingText := strings.TrimSpace(extractText(child))

			// Check for scenario header (#### Scenario:)
			if headingLevel == 4 && strings.HasPrefix(headingText, "Scenario:") {
				// Save previous scenario
				if inScenario {
					scenarios = append(scenarios, strings.TrimSpace(currentScenario.String()))
				}

				// Start new scenario
				currentScenario.Reset()
				currentScenario.WriteString(renderNode(child))
				inScenario = true

				continue
			}

			// Higher level header ends scenario
			if headingLevel <= 4 && inScenario {
				scenarios = append(scenarios, strings.TrimSpace(currentScenario.String()))
				inScenario = false
				currentScenario.Reset()

				// Check if this is a new scenario
				if headingLevel == 4 && strings.HasPrefix(headingText, "Scenario:") {
					currentScenario.WriteString(renderNode(child))
					inScenario = true
				}

				continue
			}
		}

		// Append content to current scenario
		if inScenario {
			currentScenario.WriteString(renderNode(child))
		}
	}

	// Save final scenario
	if inScenario {
		scenarios = append(scenarios, strings.TrimSpace(currentScenario.String()))
	}

	return scenarios
}

// ExtractScenarioNames extracts scenario names from requirement content.
// Returns just the scenario name strings (text after "Scenario:").
func ExtractScenarioNames(content string) []string {
	names := make([]string, 0)

	// Parse the content
	node := Parse([]byte(content))
	if node == nil {
		return names
	}

	node.Walk(func(n *bf.Node, entering bool) bf.WalkStatus {
		if !entering || n.Type != bf.Heading {
			return bf.GoToNext
		}

		if n.Level != 4 {
			return bf.GoToNext
		}

		headingText := strings.TrimSpace(extractText(n))
		if strings.HasPrefix(headingText, "Scenario:") {
			scenarioPrefix := "Scenario:"
			name := strings.TrimSpace(strings.TrimPrefix(headingText, scenarioPrefix))
			names = append(names, name)
		}

		return bf.GoToNext
	})

	return names
}

// ExtractDeltaSectionContent extracts the content of a delta section.
// sectionType should be one of: ADDED, MODIFIED, REMOVED, RENAMED
func ExtractDeltaSectionContent(node *bf.Node, sectionType string) string {
	sectionType = strings.ToUpper(strings.TrimSpace(sectionType))

	// Find the delta section header
	headerNode := FindDeltaSection(node, sectionType)
	if headerNode == nil {
		return ""
	}

	// Extract content after this header until next H2
	return ExtractSectionContentAfterNode(headerNode, 2)
}

// ExtractRequirementsSection extracts the content under "## Requirements" header.
func ExtractRequirementsSection(node *bf.Node) string {
	return ExtractSectionContent(node, "Requirements", 2)
}

// FindRequirementsHeader finds the "## Requirements" header node.
func FindRequirementsHeader(node *bf.Node) *bf.Node {
	if node == nil {
		return nil
	}

	var result *bf.Node
	node.Walk(func(n *bf.Node, entering bool) bf.WalkStatus {
		if !entering || n.Type != bf.Heading {
			return bf.GoToNext
		}

		if n.Level != 2 {
			return bf.GoToNext
		}

		headerText := strings.TrimSpace(extractText(n))
		if headerText == "Requirements" {
			result = n

			return bf.Terminate
		}

		return bf.GoToNext
	})

	return result
}

// ExtractOrderedRequirementNames extracts requirement names in document order.
// This is useful for preserving requirement order during spec merging.
func ExtractOrderedRequirementNames(content string) []string {
	names := make([]string, 0)

	// Simple line-by-line extraction for ordering (faster than full AST)
	scanner := bufio.NewScanner(strings.NewReader(content))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "### Requirement:") {
			name := strings.TrimSpace(strings.TrimPrefix(line, "### Requirement:"))
			names = append(names, name)
		}
	}

	return names
}

// SplitSpec splits a spec into preamble, requirements section, and after sections.
// Returns (preamble, requirementsContent, afterContent).
func SplitSpec(content string) (preamble, requirements, after string) {
	node := Parse([]byte(content))
	if node == nil {
		return content, "", ""
	}

	// Find the ## Requirements header
	reqHeader := FindRequirementsHeader(node)
	if reqHeader == nil {
		return content, "", ""
	}

	var preambleBuilder strings.Builder
	var requirementsBuilder strings.Builder
	var afterBuilder strings.Builder

	state := 0 // 0=preamble, 1=requirements, 2=after

	for child := node.FirstChild; child != nil; child = child.Next {
		if child.Type == bf.Heading && child.Level == 2 {
			headerText := strings.TrimSpace(extractText(child))

			if headerText == "Requirements" {
				state = 1
				preambleBuilder.WriteString(renderNode(child))

				continue
			}

			if state == 1 {
				// Hit next H2 after Requirements
				state = 2
			}
		}

		switch state {
		case 0:
			preambleBuilder.WriteString(renderNode(child))
		case 1:
			requirementsBuilder.WriteString(renderNode(child))
		case 2:
			afterBuilder.WriteString(renderNode(child))
		}
	}

	return preambleBuilder.String(), requirementsBuilder.String(), afterBuilder.String()
}
