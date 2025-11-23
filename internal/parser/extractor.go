package parser

import (
	"fmt"
	"strings"
)

const (
	// endOfDocumentLine represents the end of document for range queries
	endOfDocumentLine = 999999
	// endOfDocumentOffset represents the end of document offset
	endOfDocumentOffset = 999999999
	// scenarioHeaderLevel is the header level for scenarios
	scenarioHeaderLevel = 4
	// newlineSeparator is used for joining content parts
	newlineSeparator = "\n"
)

// Requirement represents a parsed requirement with its scenarios.
//
// Designed to match RequirementBlock from internal/parsers
// with cleaner name/content separation.
type Requirement struct {
	Name      string     // Requirement name (from "### Requirement: <name>")
	Content   string     // Full content including scenarios
	Scenarios []Scenario // Parsed scenarios within this requirement
	Position  Position   // Position where the requirement starts
}

// Scenario represents a parsed scenario within a requirement.
type Scenario struct {
	Name     string   // Scenario name (from "#### Scenario: <name>")
	Content  string   // Scenario content (list items, paragraphs, etc.)
	Position Position // Position where the scenario starts
}

// Section represents a level-2 section in the document.
//
// This is used for extracting delta sections like "## ADDED Requirements"
// or "## MODIFIED Requirements".
type Section struct {
	Name     string   // Section name (text after ##)
	Content  string   // Full content of the section
	Position Position // Position where the section starts
}

// ExtractRequirements extracts all requirements from the document.
//
// It finds "### Requirement:" headers and collects content and scenarios
// for each requirement. Returns an error if scenarios are found outside
// of requirements (hierarchy violation).
func ExtractRequirements(doc *Document) ([]Requirement, error) {
	var requirements []Requirement

	// Find all requirement headers
	reqHeaders := FindHeaders(doc, func(h *Header) bool {
		text := strings.TrimSpace(h.Text)

		return h.Level == 3 && strings.HasPrefix(text, "Requirement:")
	})

	// Track all scenario positions to validate hierarchy
	scenarioHeaders := FindHeaders(doc, func(h *Header) bool {
		text := strings.TrimSpace(h.Text)

		return h.Level == 4 && strings.HasPrefix(text, "Scenario:")
	})

	// Validate that all scenarios are within requirements
	err := validateScenarioHierarchy(doc, reqHeaders, scenarioHeaders)
	if err != nil {
		return nil, err
	}

	// Extract each requirement
	for i, reqHeader := range reqHeaders {
		req, err := extractRequirement(doc, reqHeader, reqHeaders, i)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, req)
	}

	return requirements, nil
}

// ExtractSections extracts all level-2 sections from the document.
//
// This is used for finding delta operation sections like:
// - "## ADDED Requirements"
// - "## MODIFIED Requirements"
// - "## REMOVED Requirements"
func ExtractSections(doc *Document) ([]Section, error) {
	var sections []Section

	// Find all level-2 headers
	sectionHeaders := FindHeaders(doc, func(h *Header) bool {
		return h.Level == 2
	})

	// Extract content for each section
	for i, header := range sectionHeaders {
		// Determine the end position (start of next section or end of document)
		var endPos Position
		if i+1 < len(sectionHeaders) {
			endPos = sectionHeaders[i+1].Pos()
		} else {
			// Use a very large offset to get all remaining content
			endPos = Position{
				Line:   endOfDocumentLine,
				Column: 1,
				Offset: endOfDocumentOffset,
			}
		}

		// Extract nodes between this header and the next
		nodes := NodesBetween(doc, header.Pos(), endPos)
		content := extractSectionContent(nodes)

		sections = append(sections, Section{
			Name:     strings.TrimSpace(header.Text),
			Content:  content,
			Position: header.Pos(),
		})
	}

	return sections, nil
}

// extractRequirement extracts a single requirement with its scenarios.
func extractRequirement(
	doc *Document,
	reqHeader *Header,
	allReqHeaders []*Header,
	currentIndex int,
) (Requirement, error) {
	// Extract requirement name
	name := extractRequirementName(reqHeader.Text)

	// Determine the end position for this requirement
	var endPos Position
	if currentIndex+1 < len(allReqHeaders) {
		// End at the next requirement
		endPos = allReqHeaders[currentIndex+1].Pos()
	} else {
		// Last requirement: find next level-2 header or end
		endPos = findNextLevel2Header(doc, reqHeader.Pos())
	}

	// Get all nodes between this requirement and the next
	nodes := NodesBetween(doc, reqHeader.Pos(), endPos)

	// Extract scenarios from the nodes
	scenarios, err := extractScenarios(nodes)
	if err != nil {
		return Requirement{}, err
	}

	// Extract full content (including scenarios)
	content := extractTextContent(nodes)

	return Requirement{
		Name:      name,
		Content:   content,
		Scenarios: scenarios,
		Position:  reqHeader.Pos(),
	}, nil
}

// findNextLevel2Header finds the next level-2 header after the given position.
// Returns a far-future position if no level-2 header is found.
func findNextLevel2Header(doc *Document, afterPos Position) Position {
	headers := FindHeaders(doc, func(h *Header) bool {
		return h.Level == 2 && h.Pos().Offset > afterPos.Offset
	})

	if len(headers) > 0 {
		return headers[0].Pos()
	}

	// No level-2 header found, return far-future position
	return Position{
		Line:   endOfDocumentLine,
		Column: 1,
		Offset: endOfDocumentOffset,
	}
}

// extractScenarios extracts all scenarios from a list of nodes.
func extractScenarios(nodes []Node) ([]Scenario, error) {
	var scenarios []Scenario
	var currentScenario *Scenario

	for _, node := range nodes {
		header, ok := node.(*Header)
		if !ok || header.Level != scenarioHeaderLevel {
			appendContentToScenario(currentScenario, node)

			continue
		}

		// Check if this is a scenario header
		if strings.HasPrefix(strings.TrimSpace(header.Text), "Scenario:") {
			currentScenario = saveAndStartNewScenario(
				&scenarios, currentScenario, header,
			)

			continue
		}

		// Non-scenario level-4 header ends current scenario
		if currentScenario != nil {
			scenarios = append(scenarios, *currentScenario)
			currentScenario = nil
		}
	}

	// Don't forget the last scenario
	if currentScenario != nil {
		scenarios = append(scenarios, *currentScenario)
	}

	return scenarios, nil
}

// saveAndStartNewScenario saves current scenario and starts a new one.
func saveAndStartNewScenario(
	scenarios *[]Scenario,
	current *Scenario,
	header *Header,
) *Scenario {
	if current != nil {
		*scenarios = append(*scenarios, *current)
	}

	name := extractScenarioName(header.Text)

	return &Scenario{
		Name:     name,
		Content:  "",
		Position: header.Pos(),
	}
}

// appendContentToScenario appends node text to current scenario if active.
func appendContentToScenario(scenario *Scenario, node Node) {
	if scenario == nil {
		return
	}

	text := nodeToText(node)
	if text == "" {
		return
	}

	if scenario.Content != "" {
		scenario.Content += newlineSeparator
	}
	scenario.Content += text
}

// validateScenarioHierarchy ensures all scenarios are within requirements.
func validateScenarioHierarchy(
	_ *Document,
	reqHeaders, scenarioHeaders []*Header,
) error {
	// If no scenarios, nothing to validate
	if len(scenarioHeaders) == 0 {
		return nil
	}

	// If there are scenarios but no requirements, that's an error
	if len(reqHeaders) == 0 {
		return fmt.Errorf(
			"found scenario at line %d but no requirements defined",
			scenarioHeaders[0].Pos().Line,
		)
	}

	// Check each scenario is within a requirement
	for _, scenario := range scenarioHeaders {
		if !isScenarioInRequirement(scenario, reqHeaders) {
			return fmt.Errorf(
				"scenario '%s' at line %d is not within a requirement",
				extractScenarioName(scenario.Text),
				scenario.Pos().Line,
			)
		}
	}

	return nil
}

// isScenarioInRequirement checks if a scenario falls within any
// requirement's range.
func isScenarioInRequirement(scenario *Header, reqHeaders []*Header) bool {
	scenarioOffset := scenario.Pos().Offset

	for i, req := range reqHeaders {
		reqStart := req.Pos().Offset

		// Determine the end of this requirement
		var reqEnd int
		if i+1 < len(reqHeaders) {
			reqEnd = reqHeaders[i+1].Pos().Offset
		} else {
			// Last requirement extends to end of document
			reqEnd = endOfDocumentOffset
		}

		// Check if scenario is in this requirement's range
		if scenarioOffset > reqStart && scenarioOffset < reqEnd {
			return true
		}
	}

	return false
}

// extractTextContent extracts text content from a list of nodes.
//
// This collects text from Paragraph and List nodes, but skips
// CodeBlock content to avoid extracting requirements that appear
// in code examples.
func extractTextContent(nodes []Node) string {
	var parts []string

	for _, node := range nodes {
		text := nodeToText(node)
		if text != "" {
			parts = append(parts, text)
		}
	}

	return strings.Join(parts, newlineSeparator)
}

// extractSectionContent extracts all content from a section including headers.
//
// This is different from extractTextContent as it includes requirement headers
// which are important for section content.
func extractSectionContent(nodes []Node) string {
	var parts []string

	for _, node := range nodes {
		switch n := node.(type) {
		case *Header:
			// Include headers in section content
			parts = append(parts, strings.TrimSpace(n.Text))
		case *Paragraph:
			parts = append(parts, strings.TrimSpace(n.Text))
		case *List:
			parts = append(parts, strings.Join(n.Items, newlineSeparator))
		case *CodeBlock:
			// Skip code blocks
			continue
		}
	}

	return strings.Join(parts, newlineSeparator)
}

// nodeToText converts a node to its text representation.
//
// Code blocks are intentionally skipped to avoid including
// requirements/scenarios that appear in examples.
func nodeToText(node Node) string {
	switch n := node.(type) {
	case *Paragraph:
		return strings.TrimSpace(n.Text)
	case *List:
		// Return all list items
		return strings.Join(n.Items, newlineSeparator)
	case *CodeBlock:
		// Explicitly skip code blocks
		return ""
	case *BlankLine:
		// Skip blank lines but preserve structure in extracted content
		return ""
	default:
		return ""
	}
}

// extractRequirementName extracts the requirement name from header text.
//
// Input: "Requirement: User Authentication"
// Output: "User Authentication"
func extractRequirementName(headerText string) string {
	text := strings.TrimSpace(headerText)
	if strings.HasPrefix(text, "Requirement:") {
		name := strings.TrimPrefix(text, "Requirement:")

		return strings.TrimSpace(name)
	}

	return text
}

// extractScenarioName extracts the scenario name from header text.
//
// Input: "Scenario: Valid credentials"
// Output: "Valid credentials"
func extractScenarioName(headerText string) string {
	text := strings.TrimSpace(headerText)
	if strings.HasPrefix(text, "Scenario:") {
		name := strings.TrimPrefix(text, "Scenario:")

		return strings.TrimSpace(name)
	}

	return text
}

// NormalizeRequirementName normalizes requirement names for comparison.
//
// This matches the behavior in internal/parsers/requirement_parser.go
// for compatibility with existing code.
func NormalizeRequirementName(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

// DeltaSpec represents all delta operations in a change specification.
//
// This structure is compatible with parsers.DeltaPlan but uses the new parser.
type DeltaSpec struct {
	Added    []Requirement        // Requirements being added
	Modified []Requirement        // Requirements being modified
	Removed  []Requirement        // Requirements being removed
	Renamed  []RenamedRequirement // Requirements being renamed
}

// RenamedRequirement represents a requirement rename operation.
type RenamedRequirement struct {
	From string // Original requirement name
	To   string // New requirement name
}

// ExtractDeltas extracts all delta operations from a document.
//
// It finds level-2 sections with delta operation headers (ADDED, MODIFIED,
// REMOVED, RENAMED) and extracts requirements within each section.
//
// Returns an error if the document structure is invalid (e.g., scenarios
// outside of requirements within delta sections).
func ExtractDeltas(doc *Document) (*DeltaSpec, error) {
	delta := &DeltaSpec{
		Added:    make([]Requirement, 0),
		Modified: make([]Requirement, 0),
		Removed:  make([]Requirement, 0),
		Renamed:  make([]RenamedRequirement, 0),
	}

	// Extract all level-2 sections
	sections, err := ExtractSections(doc)
	if err != nil {
		return nil, err
	}

	// Process each delta section
	for _, section := range sections {
		err := processDeltaSection(doc, section, delta)
		if err != nil {
			return nil, err
		}
	}

	return delta, nil
}

// processDeltaSection processes a single delta section.
func processDeltaSection(
	doc *Document,
	section Section,
	delta *DeltaSpec,
) error {
	sectionName := strings.TrimSpace(section.Name)

	// Check for ADDED Requirements
	if matchesDeltaSection(sectionName, "ADDED") {
		reqs, err := extractRequirementsFromSection(doc, section)
		if err != nil {
			return fmt.Errorf(
				"extracting ADDED requirements: %w", err,
			)
		}
		delta.Added = append(delta.Added, reqs...)

		return nil
	}

	// Check for MODIFIED Requirements
	if matchesDeltaSection(sectionName, "MODIFIED") {
		reqs, err := extractRequirementsFromSection(doc, section)
		if err != nil {
			return fmt.Errorf(
				"extracting MODIFIED requirements: %w", err,
			)
		}
		delta.Modified = append(delta.Modified, reqs...)

		return nil
	}

	// Check for REMOVED Requirements
	if matchesDeltaSection(sectionName, "REMOVED") {
		reqs, err := extractRequirementsFromSection(doc, section)
		if err != nil {
			return fmt.Errorf(
				"extracting REMOVED requirements: %w", err,
			)
		}
		delta.Removed = append(delta.Removed, reqs...)

		return nil
	}

	// Check for RENAMED Requirements
	if !matchesDeltaSection(sectionName, "RENAMED") {
		return nil
	}

	renamed, err := extractRenamedFromSection(doc, section)
	if err != nil {
		return fmt.Errorf(
			"extracting RENAMED requirements: %w", err,
		)
	}
	delta.Renamed = append(delta.Renamed, renamed...)

	return nil
}

// matchesDeltaSection checks if a section name matches a delta operation.
//
// Matches variations like "ADDED Requirements", "Added Requirements", etc.
func matchesDeltaSection(sectionName, operation string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(sectionName))
	expected := fmt.Sprintf("%s REQUIREMENTS", strings.ToUpper(operation))

	return normalized == expected
}

//nolint:revive // complexity acceptable for extraction logic

// extractRequirementsFromSection extracts requirements from a delta
// section.
//
// Creates temporary doc with section content and extracts requirements.
//
//nolint:revive // cognitive complexity acceptable for extraction logic
func extractRequirementsFromSection(
	doc *Document,
	section Section,
) ([]Requirement, error) {
	// Find the start and end positions for this section
	startPos := section.Position

	// Find all level-2 headers to determine section boundaries
	allSections, _ := ExtractSections(doc)
	var endPos Position

	// Find the next section after this one
	foundCurrent := false
	for _, s := range allSections {
		if s.Position.Offset == startPos.Offset {
			foundCurrent = true

			continue
		}
		if foundCurrent {
			endPos = s.Position

			break
		}
	}

	// If no next section, use end of document
	if endPos.Offset == 0 {
		endPos = Position{
			Line:   endOfDocumentLine,
			Column: 1,
			Offset: endOfDocumentOffset,
		}
	}

	// Get all requirement headers in this range
	reqHeaders := FindHeaders(doc, func(h *Header) bool {
		if h.Level != 3 {
			return false
		}
		if !strings.HasPrefix(strings.TrimSpace(h.Text), "Requirement:") {
			return false
		}
		// Check if this header is within the section range
		inRange := h.Pos().Offset > startPos.Offset &&
			h.Pos().Offset < endPos.Offset

		return inRange
	})

	// Get scenario headers in this range for validation
	scenarioHeaders := FindHeaders(doc, func(h *Header) bool {
		if h.Level != 4 {
			return false
		}
		if !strings.HasPrefix(strings.TrimSpace(h.Text), "Scenario:") {
			return false
		}
		inRange := h.Pos().Offset > startPos.Offset &&
			h.Pos().Offset < endPos.Offset

		return inRange
	})

	// Validate that all scenarios are within requirements
	err := validateScenarioHierarchy(doc, reqHeaders, scenarioHeaders)
	if err != nil {
		return nil, err
	}

	// Extract each requirement
	var requirements []Requirement
	for i, reqHeader := range reqHeaders {
		req, err := extractRequirement(doc, reqHeader, reqHeaders, i)
		if err != nil {
			return nil, err
		}
		requirements = append(requirements, req)
	}

	return requirements, nil
}

//nolint:revive // function length acceptable for parsing logic

// extractRenamedFromSection extracts RENAMED requirements from a section.
//
// Parses FROM/TO list item pairs in the format:
// - FROM: `### Requirement: Old Name`
// - TO: `### Requirement: New Name`
//
//nolint:revive // cognitive complexity acceptable for parsing RENAMED sections
func extractRenamedFromSection(doc *Document, section Section) ([]RenamedRequirement, error) {
	// Find the start and end positions for this section
	startPos := section.Position

	// Find section boundaries
	allSections, _ := ExtractSections(doc)
	var endPos Position

	foundCurrent := false
	for _, s := range allSections {
		if s.Position.Offset == startPos.Offset {
			foundCurrent = true

			continue
		}
		if foundCurrent {
			endPos = s.Position

			break
		}
	}

	if endPos.Offset == 0 {
		endPos = Position{
			Line:   endOfDocumentLine,
			Column: 1,
			Offset: endOfDocumentOffset,
		}
	}

	// Get all nodes in this section
	nodes := NodesBetween(doc, startPos, endPos)

	// Parse FROM/TO pairs from list items
	var renamed []RenamedRequirement
	var currentFrom string

	for _, node := range nodes {
		list, ok := node.(*List)
		if !ok {
			continue
		}

		for _, item := range list.Items {
			item = strings.TrimSpace(item)

			// Check for FROM line
			if strings.HasPrefix(item, "FROM:") {
				currentFrom = parseRenameLine(item, "FROM:")

				continue
			}

			// Check for TO line
			if currentFrom == "" || !strings.HasPrefix(item, "TO:") {
				continue
			}

			to := parseRenameLine(item, "TO:")
			if to != "" {
				renamed = append(renamed, RenamedRequirement{
					From: currentFrom,
					To:   to,
				})
				currentFrom = ""
			}
		}
	}

	return renamed, nil
}

// parseRenameLine parses a FROM or TO line and extracts the requirement name.
//
// Input: "FROM: `### Requirement: User Authentication`"
// Output: "User Authentication"
func parseRenameLine(line, prefix string) string {
	// Remove the prefix (FROM: or TO:)
	cleaned := strings.TrimPrefix(line, prefix)
	cleaned = strings.TrimSpace(cleaned)

	// Remove backticks if present
	cleaned = strings.Trim(cleaned, "`")
	cleaned = strings.TrimSpace(cleaned)

	// Remove the "### Requirement:" prefix
	cleaned = strings.TrimPrefix(cleaned, "###")
	cleaned = strings.TrimSpace(cleaned)
	cleaned = strings.TrimPrefix(cleaned, "Requirement:")
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}

//nolint:revive // file-length acceptable for extractor implementation
