// Package parsers provides extraction logic for Spectr-specific
// elements from generic markdown AST structures.
//
// Architecture:
//
// This package implements the second layer of the parsing architecture:
//  1. internal/mdparser: Generic markdown parser (AST builder)
//  2. internal/parsers: Spectr-specific extractors (business logic)
//
// The extractors in this package traverse the generic markdown AST
// and identify Spectr conventions (Requirements, Scenarios, Delta
// operations).
//
// Key Functions:
//   - ExtractRequirements: Find all "### Requirement: Name" blocks
//   - ExtractScenarios: Find all "#### Scenario: Name" blocks
//   - ExtractDeltaSections: Find delta operations (ADDED, MODIFIED,
//     REMOVED)
//   - ExtractRenamedRequirements: Parse RENAMED section with FROM/TO
//     pairs
//
// Design Principles:
//
// These extractors are pure functions that operate on immutable AST nodes.
// They do not modify the AST or maintain state. This makes them easy to
// test and reason about.
//
// The extractors handle edge cases that regex-based parsing cannot:
//   - Requirements inside code blocks (ignored)
//   - Malformed hierarchy (proper error reporting)
//   - Nested structures (correct scoping)
//
// Example Usage:
//
//	doc, _ := mdparser.Parse(content)
//	requirements, err := parsers.ExtractRequirements(doc)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	for _, req := range requirements {
//	    fmt.Printf("Requirement: %s\n", req.Name)
//	    scenarios, _ := parsers.ParseScenarios(req.Raw)
//	    fmt.Printf("  Scenarios: %d\n", len(scenarios))
//	}
//
//nolint:revive // file length justified by comprehensive extractor logic
package parsers

import (
	"fmt"
	"strings"

	"github.com/connerohnesorge/spectr/internal/mdparser"
)

// Scenario represents a test scenario within a requirement.
type Scenario struct {
	Name  string   // Scenario name
	Steps []string // Scenario steps (WHEN, THEN, GIVEN)
}

// RenamedRequirement represents a requirement rename operation.
type RenamedRequirement struct {
	From string // Old requirement name
	To   string // New requirement name
}

// ExtractRequirements extracts all requirements from a markdown document.
//
// Requirements are identified by H3 headers matching "### Requirement: [name]".
// This function traverses the AST to find requirement headers and collects
// all content (scenarios, paragraphs) that belongs to each requirement.
//
// Validation:
//   - Each requirement MUST have at least one scenario (#### Scenario: Name)
//   - Requirements inside code blocks are ignored (not extracted)
//
// Parameters:
//   - doc: Parsed markdown document (from mdparser.Parse)
//
// Returns:
//   - []RequirementBlock: Slice of requirements with names and raw content
//   - error: Validation error if a requirement has no scenarios
//
// Example:
//
//	doc, _ := mdparser.Parse(specContent)
//	requirements, err := ExtractRequirements(doc)
//	if err != nil {
//	    log.Fatalf("Invalid spec: %v", err)
//	}
//	for _, req := range requirements {
//	    fmt.Printf("Found requirement: %s\n", req.Name)
//	}
func ExtractRequirements(doc *mdparser.Document) ([]RequirementBlock, error) {
	var requirements []RequirementBlock

	for i, node := range doc.Children {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != RequirementHeaderLevel {
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
		siblings := getSiblingsUntilNextHeader(
			doc.Children, i+1, RequirementHeaderLevel,
		)

		// Extract scenarios under this requirement
		scenarios, err := ExtractScenarios(header, siblings)
		if err != nil {
			return nil, fmt.Errorf("requirement %q: %w", name, err)
		}

		// Validate at least one scenario
		if len(scenarios) == 0 {
			return nil, fmt.Errorf("requirement %q has no scenarios", name)
		}

		// Build raw content
		raw := buildRawContent(header, siblings)

		requirements = append(requirements, RequirementBlock{
			HeaderLine: fmt.Sprintf("### Requirement: %s", name),
			Name:       name,
			Raw:        raw,
		})
	}

	return requirements, nil
}

// ExtractScenarios extracts scenarios from nodes following a
// requirement header.
//
// Scenarios are identified by H4 headers matching
// "#### Scenario: [name]".
// This function collects scenario steps (WHEN/THEN/GIVEN) from list
// items following each scenario header.
//
// Parameters:
//   - reqHeader: The requirement header (for error context)
//   - siblings: AST nodes following the requirement header
//
// Returns:
//   - []Scenario: Slice of scenarios with names and steps
//   - error: Parse error if scenario structure is invalid
//
// Example:
//
//	header := &mdparser.Header{Level: 3, Text: "Requirement: Login"}
//	siblings := getSiblingsUntilNextHeader(doc.Children, i+1, 3)
//	scenarios, err := ExtractScenarios(header, siblings)
//	for _, scenario := range scenarios {
//	    fmt.Printf("Scenario: %s (%d steps)\n",
//	        scenario.Name, len(scenario.Steps))
//	}
func ExtractScenarios(
	_ *mdparser.Header,
	siblings []mdparser.Node,
) ([]Scenario, error) {
	var scenarios []Scenario

	for i, node := range siblings {
		header, ok := node.(*mdparser.Header)
		if !ok {
			continue
		}

		// Stop if we hit H3 or higher (end of requirement section)
		if header.Level <= RequirementHeaderLevel {
			break
		}

		// Only process H4 headers for scenarios
		if header.Level != ScenarioHeaderLevel {
			continue
		}

		// Check if this is a scenario header
		if !strings.HasPrefix(header.Text, "Scenario: ") {
			continue
		}

		// Extract scenario name
		name := strings.TrimPrefix(header.Text, "Scenario: ")
		name = strings.TrimSpace(name)

		// Extract steps from siblings until next scenario or end
		steps := extractScenarioSteps(siblings, i+1)

		scenarios = append(scenarios, Scenario{
			Name:  name,
			Steps: steps,
		})
	}

	return scenarios, nil
}

// ExtractDeltaSections extracts all delta sections from a markdown
// document.
//
// Delta sections are H2 headers like:
//   - ## ADDED Requirements
//   - ## MODIFIED Requirements
//   - ## REMOVED Requirements
//
// This function is used to parse delta spec files in
// changes/[name]/specs/.
//
// Parameters:
//   - doc: Parsed markdown document
//
// Returns:
//   - map[string][]RequirementBlock: Map of operation type to
//     requirements
//   - error: Parse error if delta structure is invalid
//
// Example:
//
//	doc, _ := mdparser.Parse(deltaContent)
//	deltas, err := ExtractDeltaSections(doc)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("ADDED: %d requirements\n", len(deltas["ADDED"]))
//	fmt.Printf("MODIFIED: %d\n", len(deltas["MODIFIED"]))
func ExtractDeltaSections(
	doc *mdparser.Document,
) (map[string][]RequirementBlock, error) {
	deltaOps := []string{"ADDED", "MODIFIED", "REMOVED", "RENAMED"}
	result := make(map[string][]RequirementBlock)

	for _, op := range deltaOps {
		// Find section header
		sectionIdx := findDeltaSectionHeader(doc.Children, op)
		if sectionIdx == -1 {
			continue // Section doesn't exist
		}

		// Get siblings until next H2
		siblings := getSiblingsUntilNextHeader(doc.Children, sectionIdx+1, 2)

		// For RENAMED, handle specially (no requirements, just FROM/TO pairs)
		if op == "RENAMED" {
			continue
		}

		// Extract requirements from section
		reqs, err := extractRequirementsFromSection(siblings)
		if err != nil {
			return nil, fmt.Errorf("%s section: %w", op, err)
		}

		result[op] = reqs
	}

	return result, nil
}

// ExtractRenamedRequirements extracts RENAMED requirements from a
// markdown document.
//
// The RENAMED section contains list items with FROM and TO patterns:
//   - FROM: `### Requirement: Old Name`
//   - TO: `### Requirement: New Name`
//
// This function handles malformed pairs gracefully:
//   - FROM without TO: Included with empty To field
//   - TO without FROM: Included with empty From field
//
// Parameters:
//   - doc: Parsed markdown document
//
// Returns:
//   - []RenamedRequirement: Slice of rename operations
//   - error: Always nil (malformed pairs are included in result)
//
// Example:
//
//	doc, _ := mdparser.Parse(deltaContent)
//	renamed, _ := ExtractRenamedRequirements(doc)
//	for _, r := range renamed {
//	    fmt.Printf("Rename: %s -> %s\n", r.From, r.To)
//	}
func ExtractRenamedRequirements(
	doc *mdparser.Document,
) ([]RenamedRequirement, error) {
	// Find RENAMED section header
	sectionIdx := findDeltaSectionHeader(doc.Children, "RENAMED")
	if sectionIdx == -1 {
		return nil, nil // No RENAMED section
	}

	// Get siblings until next H2
	siblings := getSiblingsUntilNextHeader(doc.Children, sectionIdx+1, 2)

	var renamed []RenamedRequirement
	var currentFrom string

	for _, node := range siblings {
		list, ok := node.(*mdparser.List)
		if !ok {
			continue
		}

		currentFrom = processRenamedListItems(list.Items, &renamed, currentFrom)
	}

	// If we have a FROM without a TO at the end, add as malformed
	if currentFrom != "" {
		renamed = append(renamed, RenamedRequirement{
			From: currentFrom,
			To:   "",
		})
	}

	return renamed, nil
}

// processRenamedListItems processes list items in a RENAMED section
// and returns the updated currentFrom value.
//
//nolint:revive // currentFrom parameter is intentionally modified
func processRenamedListItems(
	items []*mdparser.ListItem,
	renamed *[]RenamedRequirement,
	currentFrom string,
) string {
	for _, item := range items {
		text := strings.TrimSpace(item.Text)

		if strings.HasPrefix(text, "FROM: ") {
			currentFrom = handleFromPattern(text, renamed, currentFrom)

			continue
		}

		if strings.HasPrefix(text, "TO: ") {
			currentFrom = handleToPattern(text, renamed, currentFrom)
		}
	}

	return currentFrom
}

// handleFromPattern processes a FROM pattern in a RENAMED section.
func handleFromPattern(
	text string,
	renamed *[]RenamedRequirement,
	currentFrom string,
) string {
	// Save any unpaired FROM from before
	if currentFrom != "" {
		*renamed = append(*renamed, RenamedRequirement{
			From: currentFrom,
			To:   "", // Unpaired FROM
		})
	}

	from := extractRequirementNameFromBacktick(text, "FROM: ")
	if from != "" {
		return from
	}

	return ""
}

// handleToPattern processes a TO pattern in a RENAMED section.
func handleToPattern(
	text string,
	renamed *[]RenamedRequirement,
	currentFrom string,
) string {
	to := extractRequirementNameFromBacktick(text, "TO: ")
	if to == "" {
		return currentFrom
	}

	if currentFrom != "" {
		// Paired FROM and TO
		*renamed = append(*renamed, RenamedRequirement{
			From: currentFrom,
			To:   to,
		})

		return "" // Reset for next pair
	}

	// TO without FROM - add as malformed
	*renamed = append(*renamed, RenamedRequirement{
		From: "",
		To:   to,
	})

	return ""
}

// Helper functions

// getSiblingsUntilNextHeader returns nodes from start index until
// next header of given level or higher.
func getSiblingsUntilNextHeader(
	nodes []mdparser.Node,
	startIdx, maxLevel int,
) []mdparser.Node {
	if startIdx >= len(nodes) {
		return nil
	}

	var siblings []mdparser.Node
	for i := startIdx; i < len(nodes); i++ {
		node := nodes[i]

		// Check if we've hit a header of the specified level or higher
		header, ok := node.(*mdparser.Header)
		if ok && header.Level <= maxLevel {
			break
		}

		siblings = append(siblings, node)
	}

	return siblings
}

// findDeltaSectionHeader finds the index of a delta section header
// (e.g., "## ADDED Requirements").
// Returns -1 if not found.
func findDeltaSectionHeader(nodes []mdparser.Node, operation string) int {
	expectedText := fmt.Sprintf("%s Requirements", operation)

	for i, node := range nodes {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != 2 {
			continue
		}

		if strings.TrimSpace(header.Text) == expectedText {
			return i
		}
	}

	return -1
}

// extractRequirementsFromSection extracts requirements from a section's nodes.
// This handles requirements within delta sections (ADDED, MODIFIED, REMOVED).
func extractRequirementsFromSection(
	nodes []mdparser.Node,
) ([]RequirementBlock, error) {
	var requirements []RequirementBlock

	for i, node := range nodes {
		header, ok := node.(*mdparser.Header)
		if !ok || header.Level != RequirementHeaderLevel {
			continue
		}

		// Check if this is a requirement header
		if !strings.HasPrefix(header.Text, "Requirement: ") {
			continue
		}

		// Extract requirement name
		name := strings.TrimPrefix(header.Text, "Requirement: ")
		name = strings.TrimSpace(name)

		// Get siblings until next H3 or higher
		siblings := getSiblingsUntilNextHeader(
			nodes, i+1, RequirementHeaderLevel,
		)

		// Build raw content
		raw := buildRawContent(header, siblings)

		requirements = append(requirements, RequirementBlock{
			HeaderLine: fmt.Sprintf("### Requirement: %s", name),
			Name:       name,
			Raw:        raw,
		})
	}

	return requirements, nil
}

// buildRawContent constructs the raw markdown text from a header and
// its siblings.
func buildRawContent(
	header *mdparser.Header,
	siblings []mdparser.Node,
) string {
	var sb strings.Builder

	// Add header line
	reqName := header.Text[len("Requirement: "):]
	sb.WriteString(fmt.Sprintf("### Requirement: %s\n", reqName))

	// Add sibling content
	for _, node := range siblings {
		renderNode(&sb, node)
	}

	return sb.String()
}

// renderNode renders a single AST node to a string builder.
func renderNode(sb *strings.Builder, node mdparser.Node) {
	switch n := node.(type) {
	case *mdparser.Header:
		renderHeader(sb, n)
	case *mdparser.Paragraph:
		renderParagraph(sb, n)
	case *mdparser.CodeBlock:
		renderCodeBlock(sb, n)
	case *mdparser.List:
		renderList(sb, n)
	case *mdparser.BlankLine:
		renderBlankLine(sb, n)
	}
}

// renderHeader renders a header node.
func renderHeader(sb *strings.Builder, h *mdparser.Header) {
	sb.WriteString(strings.Repeat("#", h.Level))
	sb.WriteString(" ")
	sb.WriteString(h.Text)
	sb.WriteString("\n")
}

// renderParagraph renders a paragraph node.
func renderParagraph(sb *strings.Builder, p *mdparser.Paragraph) {
	for _, line := range p.Lines {
		sb.WriteString(line)
		sb.WriteString("\n")
	}
}

// renderCodeBlock renders a code block node.
//
//nolint:revive // "\n" repetition is clearer than a constant here
func renderCodeBlock(sb *strings.Builder, cb *mdparser.CodeBlock) {
	sb.WriteString("```")
	sb.WriteString(cb.Language)
	sb.WriteString("\n")

	for _, line := range cb.Lines {
		sb.WriteString(line)
		sb.WriteString("\n")
	}

	sb.WriteString("```\n")
}

// renderList renders a list node.
func renderList(sb *strings.Builder, l *mdparser.List) {
	for _, item := range l.Items {
		if l.Ordered {
			sb.WriteString("1. ")
		} else {
			sb.WriteString("- ")
		}

		sb.WriteString(item.Text)
		sb.WriteString("\n")
	}
}

// renderBlankLine renders a blank line node.
func renderBlankLine(sb *strings.Builder, bl *mdparser.BlankLine) {
	for range bl.Count {
		sb.WriteString("\n")
	}
}

// extractScenarioSteps extracts list items that follow a scenario
// header.
// These are typically WHEN/THEN/GIVEN steps.
//
//nolint:revive // early-return would make logic less clear here
func extractScenarioSteps(nodes []mdparser.Node, startIdx int) []string {
	var steps []string

	for i := startIdx; i < len(nodes); i++ {
		node := nodes[i]

		// Stop if we hit another header
		if header, ok := node.(*mdparser.Header); ok {
			// Stop at H4 or higher (next scenario or requirement)
			if header.Level <= ScenarioHeaderLevel {
				break
			}
		}

		// Extract list items
		if list, ok := node.(*mdparser.List); ok {
			for _, item := range list.Items {
				text := strings.TrimSpace(item.Text)
				// Only include items that look like scenario steps
				if strings.Contains(text, "**WHEN**") ||
					strings.Contains(text, "**THEN**") ||
					strings.Contains(text, "**GIVEN**") {
					steps = append(steps, text)
				}
			}
		}
	}

	return steps
}

// extractRequirementNameFromBacktick extracts requirement name from
// text.
// Handles both backtick and non-backtick formats:
// - "FROM: `### Requirement: Old Name`" -> "Old Name"
// - "FROM: ### Requirement: Old Name" -> "Old Name"
//
//nolint:revive // text parameter is intentionally modified for clarity
func extractRequirementNameFromBacktick(text, prefix string) string {
	// Remove prefix
	text = strings.TrimPrefix(text, prefix)
	text = strings.TrimSpace(text)

	// Extract from backticks if present (optional)
	if strings.HasPrefix(text, "`") && strings.HasSuffix(text, "`") {
		text = strings.TrimPrefix(text, "`")
		text = strings.TrimSuffix(text, "`")
		text = strings.TrimSpace(text)
	}

	// Remove "### Requirement: " prefix
	text = strings.TrimPrefix(text, "###")
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "Requirement:")
	text = strings.TrimSpace(text)

	return text
}
