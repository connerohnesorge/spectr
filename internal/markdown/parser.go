package markdown

import (
	"bytes"
	"strings"
	"unicode/utf8"

	"github.com/connerohnesorge/spectr/internal/specterrs"
	"github.com/russross/blackfriday/v2"
)

// Parser constants
const (
	tabWidth        = 4 // Number of spaces per tab for indentation
	minCheckboxLen  = 5 // Minimum length for "- [ ]"
	checkboxCharIdx = 3 // Index of checkbox character in "- [x]"
	closeBracketIdx = 4 // Index of closing bracket in "- [ ]"

	// Header levels
	headerLevelH2 = 2
	headerLevelH3 = 3
	headerLevelH4 = 4
)

// Delta section types
var deltaTypes = []string{"ADDED", "MODIFIED", "REMOVED", "RENAMED"}

// ParseDocument parses markdown content using blackfriday AST and extracts
// all structural elements in a single pass. Returns error for invalid input.
//
// The returned Document contains indexed maps for O(1) lookups:
//   - Sections by name
//   - Requirements by name
//   - Scenarios by name
//   - Headers by level
func ParseDocument(content []byte) (*Document, error) {
	// Validate input
	if len(bytes.TrimSpace(content)) == 0 {
		return nil, &specterrs.EmptyContentError{}
	}

	if !utf8.Valid(content) || bytes.Contains(content, []byte{0}) {
		return nil, &specterrs.BinaryContentError{}
	}

	lines := strings.Split(string(content), "\n")

	// Parse AST with blackfriday
	extensions := blackfriday.CommonExtensions |
		blackfriday.NoEmptyLineBeforeBlock
	parser := blackfriday.New(blackfriday.WithExtensions(extensions))
	ast := parser.Parse(content)

	doc := &Document{
		Content:      content,
		Lines:        lines,
		Headers:      make([]Header, 0),
		Sections:     make(map[string]*Section),
		Requirements: make(map[string]*Requirement),
		Scenarios:    make(map[string]*Scenario),
		H2Headers:    make([]Header, 0),
		H3Headers:    make([]Header, 0),
		H4Headers:    make([]Header, 0),
		Tasks:        make([]Task, 0),
	}

	// Single AST walk extracts all structure
	walkAndExtract(ast, doc)

	// Build section and requirement content from line ranges
	buildContent(doc)

	// Extract tasks from lines (blackfriday doesn't preserve checkbox state)
	extractTasks(doc)

	return doc, nil
}

// walkState holds the current parsing state during AST walk.
type walkState struct {
	currentSection *Section
	currentReq     *Requirement
}

// walkAndExtract performs a single AST walk to extract all headers, sections,
// requirements, and scenarios.
func walkAndExtract(root *blackfriday.Node, doc *Document) {
	state := &walkState{}

	root.Walk(func(n *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if !entering || n.Type != blackfriday.Heading {
			return blackfriday.GoToNext
		}

		handleHeading(n, doc, state)

		return blackfriday.GoToNext
	})
}

// handleHeading processes a heading node and updates document and state.
func handleHeading(n *blackfriday.Node, doc *Document, state *walkState) {
	text := extractNodeText(n)
	lineNum := findHeaderLine(doc.Lines, n.Level, text)

	header := Header{Level: n.Level, Text: text, Line: lineNum}
	doc.Headers = append(doc.Headers, header)

	switch n.Level {
	case headerLevelH2:
		handleH2(doc, state, header)
	case headerLevelH3:
		handleH3(doc, state, header)
	case headerLevelH4:
		handleH4(doc, state, header)
	}
}

// handleH2 processes an H2 section header.
func handleH2(doc *Document, state *walkState, header Header) {
	doc.H2Headers = append(doc.H2Headers, header)
	state.currentSection = &Section{
		Name:      header.Text,
		Header:    header,
		StartLine: header.Line + 1,
	}

	if dt := parseDeltaType(header.Text); dt != "" {
		state.currentSection.IsDelta = true
		state.currentSection.DeltaType = dt
	}

	doc.Sections[header.Text] = state.currentSection
	state.currentReq = nil
}

// handleH3 processes an H3 requirement header.
func handleH3(doc *Document, state *walkState, header Header) {
	doc.H3Headers = append(doc.H3Headers, header)

	name := parseRequirementName(header.Text)
	if name == "" {
		return
	}

	sectionName := ""
	if state.currentSection != nil {
		sectionName = state.currentSection.Name
	}

	state.currentReq = &Requirement{
		Name:    name,
		Section: sectionName,
		Header:  header,
	}
	doc.Requirements[name] = state.currentReq
}

// handleH4 processes an H4 scenario header.
func handleH4(doc *Document, state *walkState, header Header) {
	doc.H4Headers = append(doc.H4Headers, header)

	name := parseScenarioName(header.Text)
	if name == "" || state.currentReq == nil {
		return
	}

	scenario := &Scenario{
		Name:        name,
		Requirement: state.currentReq.Name,
		Header:      header,
	}
	doc.Scenarios[name] = scenario
	state.currentReq.Scenarios = append(state.currentReq.Scenarios, scenario)
}

// extractNodeText walks node children to build text content.
func extractNodeText(n *blackfriday.Node) string {
	var buf strings.Builder

	n.Walk(func(child *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering && child.Literal != nil {
			buf.Write(child.Literal)
		}

		return blackfriday.GoToNext
	})

	return buf.String()
}

// findHeaderLine finds the line number of a header by searching lines.
func findHeaderLine(lines []string, level int, text string) int {
	prefix := strings.Repeat("#", level) + " "

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, prefix) {
			continue
		}

		lineText := strings.TrimSpace(strings.TrimPrefix(trimmed, prefix))
		if lineText == text {
			return i + 1 // 1-indexed
		}
	}

	return 0
}

// parseDeltaType extracts delta type from "ADDED Requirements" etc.
func parseDeltaType(text string) string {
	for _, dt := range deltaTypes {
		if text == dt+" Requirements" {
			return dt
		}
	}

	return ""
}

// parseRequirementName extracts name from "Requirement: Foo".
func parseRequirementName(text string) string {
	const prefix = "Requirement: "
	if strings.HasPrefix(text, prefix) {
		return strings.TrimSpace(text[len(prefix):])
	}

	return ""
}

// parseScenarioName extracts name from "Scenario: Foo".
func parseScenarioName(text string) string {
	const prefix = "Scenario: "
	if strings.HasPrefix(text, prefix) {
		return strings.TrimSpace(text[len(prefix):])
	}

	return ""
}
