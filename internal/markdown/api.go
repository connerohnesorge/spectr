//nolint:revive // file-length-limit: API surface requires comprehensive methods
package markdown

import (
	"strings"
)

// Spec represents a fully parsed specification file.
// It provides convenient access to sections, requirements, and parse errors.
type Spec struct {
	// Root is the document node containing the entire AST.
	Root Node

	// Sections contains all sections indexed by their name (case-preserved).
	// Use FindSection for case-insensitive lookup.
	Sections map[string]*Section

	// Requirements contains all requirements found in the document.
	Requirements []*Requirement

	// Errors contains all parse errors encountered during parsing.
	Errors []ParseError
}

// Section represents a document section (header and its content).
type Section struct {
	// Name is the section title text.
	Name string

	// Level is the header level (1-6).
	Level int

	// Start is the byte offset where the section starts.
	Start int

	// End is the byte offset where the section ends.
	End int

	// Content is the raw source content of the section.
	Content []byte

	// Node is the underlying AST node for this section.
	Node *NodeSection
}

// Requirement represents a parsed requirement.
type Requirement struct {
	// Name is the requirement name (text after "Requirement:").
	Name string

	// Section is the name of the parent section containing this requirement.
	Section string

	// Scenarios contains all scenarios within this requirement.
	Scenarios []*Scenario

	// Node is the underlying AST node for this requirement.
	Node *NodeRequirement
}

// Scenario represents a parsed scenario.
type Scenario struct {
	// Name is the scenario name (text after "Scenario:").
	Name string

	// Node is the underlying AST node for this scenario.
	Node *NodeScenario
}

// Wikilink represents an extracted wikilink with parsed components.
type Wikilink struct {
	// Target is the link target (e.g., "validation" or "changes/my-change").
	Target string

	// Display is the optional display text (empty if not specified).
	Display string

	// Anchor is the optional anchor/fragment (empty if not specified).
	Anchor string

	// Start is the byte offset where the wikilink starts.
	Start int

	// End is the byte offset where the wikilink ends.
	End int

	// Node is the underlying AST node for this wikilink.
	Node *NodeWikilink
}

// ParseSpec parses markdown content and returns a Spec with structured access
// to sections, requirements, and scenarios.
func ParseSpec(
	content []byte,
) (*Spec, []ParseError) {
	root, errors := Parse(content)

	spec := &Spec{
		Root:         root,
		Sections:     make(map[string]*Section),
		Requirements: make([]*Requirement, 0),
		Errors:       errors,
	}

	if root == nil {
		return spec, errors
	}

	// Extract sections and requirements using a visitor
	extractor := &specExtractor{
		spec:           spec,
		source:         content,
		currentSection: "",
	}

	_ = Walk(root, extractor)

	return spec, errors
}

// specExtractor is a visitor that extracts sections and requirements.
// It tracks the current requirement and associates scenarios with it.
type specExtractor struct {
	BaseVisitor
	spec           *Spec
	source         []byte
	currentSection string
	currentReq     *Requirement
}

// VisitSection extracts section information and resets current requirement.
func (e *specExtractor) VisitSection(
	n *NodeSection,
) error {
	name := string(n.Title())
	start, end := n.Span()

	section := &Section{
		Name:    name,
		Level:   n.Level(),
		Start:   start,
		End:     end,
		Content: e.source[start:end],
		Node:    n,
	}

	e.spec.Sections[name] = section
	e.currentSection = name
	// Reset current requirement when entering a new section
	e.currentReq = nil

	return nil
}

// VisitRequirement extracts requirement information.
func (e *specExtractor) VisitRequirement(
	n *NodeRequirement,
) error {
	req := &Requirement{
		Name:      n.Name(),
		Section:   e.currentSection,
		Scenarios: make([]*Scenario, 0),
		Node:      n,
	}

	e.spec.Requirements = append(
		e.spec.Requirements,
		req,
	)
	e.currentReq = req

	return nil
}

// VisitScenario adds scenarios to the current requirement.
func (e *specExtractor) VisitScenario(
	n *NodeScenario,
) error {
	// Add to current requirement if one exists
	if e.currentReq != nil {
		e.currentReq.Scenarios = append(
			e.currentReq.Scenarios,
			&Scenario{
				Name: n.Name(),
				Node: n,
			},
		)
	}

	return nil
}

// ExtractSections parses content and returns a map of section name to Section.
// The keys are the exact section titles (case-preserved).
func ExtractSections(
	content []byte,
) map[string]*Section {
	root, _ := Parse(content)
	if root == nil {
		return make(map[string]*Section)
	}

	sections := make(map[string]*Section)

	extractor := &sectionExtractor{
		sections: sections,
		source:   content,
	}

	_ = Walk(root, extractor)

	return sections
}

// sectionExtractor is a visitor that extracts only sections.
type sectionExtractor struct {
	BaseVisitor
	sections map[string]*Section
	source   []byte
}

// VisitSection extracts section information.
func (e *sectionExtractor) VisitSection(
	n *NodeSection,
) error {
	name := string(n.Title())
	start, end := n.Span()

	e.sections[name] = &Section{
		Name:    name,
		Level:   n.Level(),
		Start:   start,
		End:     end,
		Content: e.source[start:end],
		Node:    n,
	}

	return nil
}

// ExtractRequirements parses content and returns all requirements.
func ExtractRequirements(
	content []byte,
) []*Requirement {
	root, _ := Parse(content)
	if root == nil {
		return []*Requirement{}
	}

	extractor := &requirementExtractor{
		requirements:   make([]*Requirement, 0),
		currentSection: "",
	}

	_ = Walk(root, extractor)

	return extractor.requirements
}

// requirementExtractor is a visitor that extracts requirements.
// It tracks scenarios that follow each requirement.
type requirementExtractor struct {
	BaseVisitor
	requirements   []*Requirement
	currentSection string
	currentReq     *Requirement
}

// VisitSection tracks current section and resets current requirement.
func (e *requirementExtractor) VisitSection(
	n *NodeSection,
) error {
	e.currentSection = string(n.Title())
	e.currentReq = nil

	return nil
}

// VisitRequirement extracts requirement and sets it as current.
func (e *requirementExtractor) VisitRequirement(
	n *NodeRequirement,
) error {
	req := &Requirement{
		Name:      n.Name(),
		Section:   e.currentSection,
		Scenarios: make([]*Scenario, 0),
		Node:      n,
	}

	e.requirements = append(e.requirements, req)
	e.currentReq = req

	return nil
}

// VisitScenario adds scenarios to the current requirement.
func (e *requirementExtractor) VisitScenario(
	n *NodeScenario,
) error {
	if e.currentReq != nil {
		e.currentReq.Scenarios = append(
			e.currentReq.Scenarios,
			&Scenario{
				Name: n.Name(),
				Node: n,
			},
		)
	}

	return nil
}

// FindSection performs case-insensitive lookup for a section by name.
// Returns the section and true if found, nil and false otherwise.
func FindSection(
	content []byte,
	name string,
) (*Section, bool) {
	root, _ := Parse(content)
	if root == nil {
		return nil, false
	}

	finder := &sectionFinder{
		targetName: strings.ToLower(name),
		source:     content,
	}

	_ = Walk(root, finder)

	if finder.found != nil {
		return finder.found, true
	}

	return nil, false
}

// sectionFinder is a visitor that finds a section by name (case-insensitive).
type sectionFinder struct {
	BaseVisitor
	targetName string
	source     []byte
	found      *Section
}

// VisitSection checks if this section matches the target name.
func (f *sectionFinder) VisitSection(
	n *NodeSection,
) error {
	name := string(n.Title())
	if strings.ToLower(name) == f.targetName {
		start, end := n.Span()
		f.found = &Section{
			Name:    name,
			Level:   n.Level(),
			Start:   start,
			End:     end,
			Content: f.source[start:end],
			Node:    n,
		}
		// Stop traversal once found
		return SkipChildren
	}

	return nil
}

// ExtractWikilinks parses content and returns all wikilinks found.
// Uses a visitor to traverse the AST and collect all NodeWikilink nodes.
func ExtractWikilinks(
	content []byte,
) []*Wikilink {
	root, _ := Parse(content)
	if root == nil {
		return []*Wikilink{}
	}

	extractor := &wikilinkExtractor{
		wikilinks: make([]*Wikilink, 0),
	}

	_ = Walk(root, extractor)

	return extractor.wikilinks
}

// wikilinkExtractor is a visitor that collects all wikilinks.
type wikilinkExtractor struct {
	BaseVisitor
	wikilinks []*Wikilink
}

// VisitWikilink extracts wikilink information.
func (e *wikilinkExtractor) VisitWikilink(
	n *NodeWikilink,
) error {
	start, end := n.Span()

	wikilink := &Wikilink{
		Target:  string(n.Target()),
		Display: string(n.Display()),
		Anchor:  string(n.Anchor()),
		Start:   start,
		End:     end,
		Node:    n,
	}

	e.wikilinks = append(e.wikilinks, wikilink)

	return nil
}

// FindRequirement finds a requirement by name (case-insensitive).
// Returns the requirement and true if found, nil and false otherwise.
func FindRequirement(
	content []byte,
	name string,
) (*Requirement, bool) {
	requirements := ExtractRequirements(content)
	targetName := strings.ToLower(name)

	for _, req := range requirements {
		if strings.ToLower(
			req.Name,
		) == targetName {
			return req, true
		}
	}

	return nil, false
}

// FindScenario finds a scenario by name (case-insensitive).
// Returns the scenario and true if found, nil and false otherwise.
func FindScenario(
	content []byte,
	name string,
) (*Scenario, bool) {
	root, _ := Parse(content)
	if root == nil {
		return nil, false
	}

	targetName := strings.ToLower(name)

	finder := &scenarioFinder{
		targetName: targetName,
	}

	_ = Walk(root, finder)

	if finder.found != nil {
		return finder.found, true
	}

	return nil, false
}

// scenarioFinder is a visitor that finds a scenario by name.
type scenarioFinder struct {
	BaseVisitor
	targetName string
	found      *Scenario
}

// VisitScenario checks if this scenario matches.
func (f *scenarioFinder) VisitScenario(
	n *NodeScenario,
) error {
	if strings.ToLower(n.Name()) == f.targetName {
		f.found = &Scenario{
			Name: n.Name(),
			Node: n,
		}
		// Stop traversal
		return SkipChildren
	}

	return nil
}

// GetSectionContent returns the text content of a section by name.
// Returns empty string if section not found.
func GetSectionContent(
	content []byte,
	sectionName string,
) string {
	section, found := FindSection(
		content,
		sectionName,
	)
	if !found {
		return ""
	}

	return string(section.Content)
}

// HasRequirement checks if a requirement with the given name exists.
func HasRequirement(
	content []byte,
	name string,
) bool {
	_, found := FindRequirement(content, name)

	return found
}

// HasSection checks if a section with the given name exists (case-insensitive).
func HasSection(
	content []byte,
	name string,
) bool {
	_, found := FindSection(content, name)

	return found
}

// GetRequirementNames returns the names of all requirements in the content.
func GetRequirementNames(
	content []byte,
) []string {
	requirements := ExtractRequirements(content)
	names := make([]string, len(requirements))
	for i, req := range requirements {
		names[i] = req.Name
	}

	return names
}

// GetSectionNames returns the names of all sections in the content.
func GetSectionNames(content []byte) []string {
	sections := ExtractSections(content)
	names := make([]string, 0, len(sections))
	for name := range sections {
		names = append(names, name)
	}

	return names
}

// CountRequirements returns the total number of requirements in the content.
func CountRequirements(content []byte) int {
	return len(ExtractRequirements(content))
}

// CountScenarios returns the total number of scenarios in the content.
func CountScenarios(content []byte) int {
	root, _ := Parse(content)
	if root == nil {
		return 0
	}

	return Count(root, IsType[*NodeScenario]())
}

// CountWikilinks returns the total number of wikilinks in the content.
func CountWikilinks(content []byte) int {
	root, _ := Parse(content)
	if root == nil {
		return 0
	}

	return Count(root, IsType[*NodeWikilink]())
}
