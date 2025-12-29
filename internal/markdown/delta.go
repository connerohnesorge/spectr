package markdown

import (
	"bytes"
	"strings"
)

// DeltaType represents the type of change in a delta specification.
type DeltaType string

// Delta type constants for categorizing requirement changes.
const (
	DeltaAdded    DeltaType = "ADDED"
	DeltaModified DeltaType = "MODIFIED"
	DeltaRemoved  DeltaType = "REMOVED"
	DeltaRenamed  DeltaType = "RENAMED"
)

// String length constants for parsing FROM/TO annotations.
const (
	lenFROM      = 4 // len("FROM")
	lenFROMColon = 5 // len("FROM:")
)

// Delta represents a parsed delta specification file.
// Delta files describe changes to requirements: additions, modifications,
// removals, and renames.
type Delta struct {
	// Root is the document node containing the entire AST.
	Root Node

	// Added contains requirements that are being added.
	Added []*Requirement

	// Modified contains requirements that are being modified.
	Modified []*Requirement

	// Removed contains the names of requirements being removed.
	// Only names are stored since the requirements no longer exist.
	Removed []string

	// Renamed contains rename operations with from/to pairs.
	Renamed []*RenamedRequirement

	// Errors contains all parse errors encountered during parsing.
	Errors []ParseError
}

// RenamedRequirement represents a requirement rename operation.
// It contains the old name and the new name of the requirement.
type RenamedRequirement struct {
	// From is the original requirement name.
	From string

	// To is the new requirement name.
	To string

	// Node is the AST node if parsed from a requirement header.
	// May be nil if parsed from a list item.
	Node *NodeRequirement
}

// ParseDelta parses markdown content as a delta specification and returns
// a Delta with categorized requirement changes.
func ParseDelta(
	content []byte,
) (*Delta, []ParseError) {
	root, errors := Parse(content)

	delta := &Delta{
		Root:     root,
		Added:    make([]*Requirement, 0),
		Modified: make([]*Requirement, 0),
		Removed:  make([]string, 0),
		Renamed:  make([]*RenamedRequirement, 0),
		Errors:   errors,
	}

	if root == nil {
		return delta, errors
	}

	// Extract delta sections using a visitor
	extractor := &deltaExtractor{
		delta:            delta,
		source:           content,
		currentDeltaType: "",
		currentSection:   "",
	}

	_ = Walk(root, extractor)

	return delta, errors
}

// deltaExtractor is a visitor that extracts delta-categorized requirements.
type deltaExtractor struct {
	BaseVisitor
	delta            *Delta
	source           []byte
	currentDeltaType DeltaType
	currentSection   string
	currentReq       *Requirement
}

// VisitSection tracks delta sections and updates the current delta type.
func (e *deltaExtractor) VisitSection(
	n *NodeSection,
) error {
	e.currentSection = string(n.Title())

	// Check if this is a delta section
	deltaType := n.DeltaType()
	if deltaType != "" {
		e.currentDeltaType = DeltaType(deltaType)
	} else {
		// Not a delta section, reset delta type
		e.currentDeltaType = ""
	}

	// Reset current requirement when entering a new section
	e.currentReq = nil

	return nil
}

// VisitRequirement extracts requirements and categorizes by delta type.
func (e *deltaExtractor) VisitRequirement(
	n *NodeRequirement,
) error {
	req := &Requirement{
		Name:      n.Name(),
		Section:   e.currentSection,
		Scenarios: make([]*Scenario, 0),
		Node:      n,
	}

	switch e.currentDeltaType {
	case DeltaAdded:
		e.delta.Added = append(e.delta.Added, req)
		e.currentReq = req
	case DeltaModified:
		e.delta.Modified = append(
			e.delta.Modified,
			req,
		)
		e.currentReq = req
	case DeltaRenamed:
		// For RENAMED sections, we need to look for FROM: annotation
		// The requirement name is the new name (To)
		renamed := &RenamedRequirement{
			To:   n.Name(),
			Node: n,
		}
		e.delta.Renamed = append(
			e.delta.Renamed,
			renamed,
		)
		e.currentReq = req
	case DeltaRemoved:
		// Requirements in REMOVED sections are handled via list items
		e.currentReq = req
	default:
		// Not in a delta section, ignore or treat as regular requirement
		e.currentReq = req
	}

	return nil
}

// VisitScenario adds scenarios to the current requirement.
func (e *deltaExtractor) VisitScenario(
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

// VisitList processes list items for REMOVED requirements and RENAMED pairs.
//
//nolint:revive // unused-receiver: interface requires method
func (*deltaExtractor) VisitList(
	*NodeList,
) error {
	return nil // Continue to visit children
}

// VisitListItem handles list items in delta sections.
func (e *deltaExtractor) VisitListItem(
	n *NodeListItem,
) error {
	// Get the source text of the list item
	itemText := strings.TrimSpace(
		string(n.Source()),
	)

	switch e.currentDeltaType {
	case DeltaRemoved:
		// List items in REMOVED sections are requirement names
		// Extract the name (after the bullet point marker)
		name := extractListItemText(itemText)
		if name != "" {
			e.delta.Removed = append(
				e.delta.Removed,
				name,
			)
		}

	case DeltaRenamed:
		// Check for FROM: OldName TO: NewName format in list items
		from, to := parseRenamedListItem(itemText)
		if from != "" && to != "" {
			e.delta.Renamed = append(
				e.delta.Renamed,
				&RenamedRequirement{
					From: from,
					To:   to,
				},
			)
		}
	case DeltaAdded, DeltaModified:
		// List items in ADDED/MODIFIED sections are handled differently
	default:
		// Not in a delta section
	}

	return nil
}

// VisitParagraph handles paragraphs for FROM: annotations in RENAMED sections.
func (e *deltaExtractor) VisitParagraph(
	n *NodeParagraph,
) error {
	if e.currentDeltaType != DeltaRenamed {
		return nil
	}

	// Check if paragraph contains FROM: annotation
	paragraphText := string(n.Source())
	fromName := parseFromAnnotation(paragraphText)

	if fromName != "" &&
		len(e.delta.Renamed) > 0 {
		// Associate with the most recent renamed requirement
		lastRenamed := e.delta.Renamed[len(e.delta.Renamed)-1]
		if lastRenamed.From == "" {
			lastRenamed.From = fromName
		}
	}

	return nil
}

// FindDeltaSection returns the content of a specific delta section type.
// It searches for sections matching the given delta type
// (ADDED, MODIFIED, REMOVED, RENAMED).
// Returns empty string if no matching section is found.
func FindDeltaSection(
	content []byte,
	deltaType DeltaType,
) string {
	root, _ := Parse(content)
	if root == nil {
		return ""
	}

	finder := &deltaSectionFinder{
		targetType: deltaType,
		source:     content,
	}

	_ = Walk(root, finder)

	if finder.found != nil {
		return string(finder.found.Content)
	}

	return ""
}

// deltaSectionFinder is a visitor that finds a delta section by type.
type deltaSectionFinder struct {
	BaseVisitor
	targetType DeltaType
	source     []byte
	found      *Section
}

// VisitSection checks if this section matches the target delta type.
func (f *deltaSectionFinder) VisitSection(
	n *NodeSection,
) error {
	if n.DeltaType() == string(f.targetType) {
		start, end := n.Span()
		f.found = &Section{
			Name:    string(n.Title()),
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

// FindAllDeltaSections returns all delta sections in the content.
// Returns a map from delta type to the section content.
func FindAllDeltaSections(
	content []byte,
) map[DeltaType]*Section {
	root, _ := Parse(content)
	if root == nil {
		return make(map[DeltaType]*Section)
	}

	finder := &allDeltaSectionsFinder{
		sections: make(map[DeltaType]*Section),
		source:   content,
	}

	_ = Walk(root, finder)

	return finder.sections
}

// allDeltaSectionsFinder is a visitor that collects all delta sections.
type allDeltaSectionsFinder struct {
	BaseVisitor
	sections map[DeltaType]*Section
	source   []byte
}

// VisitSection collects delta sections.
func (f *allDeltaSectionsFinder) VisitSection(
	n *NodeSection,
) error {
	deltaType := n.DeltaType()
	if deltaType != "" {
		start, end := n.Span()
		f.sections[DeltaType(deltaType)] = &Section{
			Name:    string(n.Title()),
			Level:   n.Level(),
			Start:   start,
			End:     end,
			Content: f.source[start:end],
			Node:    n,
		}
	}

	return nil
}

// parseRenamedListItem parses a list item for FROM: OldName TO: NewName format.
// Returns (from, to) names, or empty strings if not found.
func parseRenamedListItem(
	text string,
) (from, to string) {
	// Normalize and look for patterns:
	// - FROM: OldName TO: NewName
	// - from: OldName to: NewName
	// - FROM OldName TO NewName (without colons)
	normalizedText := strings.TrimSpace(
		text,
	) //nolint:revive // modifies-parameter
	normalizedText = extractListItemText(
		normalizedText,
	) //nolint:revive // modifies-parameter

	upper := strings.ToUpper(normalizedText)

	// Try "FROM: ... TO: ..." format
	fromIdx := strings.Index(upper, "FROM:")
	if fromIdx == -1 {
		// Try without colon
		fromIdx = strings.Index(upper, "FROM ")
	}

	if fromIdx == -1 {
		return "", ""
	}

	// Find the TO part
	toIdx := strings.Index(upper, "TO:")
	if toIdx == -1 {
		toIdx = strings.Index(upper, " TO ")
		if toIdx != -1 {
			toIdx++ // Adjust for space
		}
	}

	if toIdx == -1 || toIdx <= fromIdx {
		return "", ""
	}

	// Extract FROM value
	fromStart := fromIdx + lenFROM
	if fromIdx < len(normalizedText) &&
		normalizedText[fromStart] == ':' {
		fromStart++
	}
	fromValue := strings.TrimSpace(
		normalizedText[fromStart:toIdx],
	)

	// Extract TO value
	toStart := toIdx + 2 // len("TO")
	if toStart < len(normalizedText) &&
		normalizedText[toStart] == ':' {
		toStart++
	}
	toValue := strings.TrimSpace(
		normalizedText[toStart:],
	)

	return fromValue, toValue
}

// parseFromAnnotation extracts the FROM: OldName annotation from text.
// Returns the old name, or empty string if not found.
func parseFromAnnotation(text string) string {
	// Look for FROM: followed by a name
	upper := strings.ToUpper(text)

	fromIdx := strings.Index(upper, "FROM:")
	if fromIdx == -1 {
		return ""
	}

	// Find the value after FROM:
	start := fromIdx + lenFROMColon

	// Find the end of the line or next annotation
	end := len(text)
	if newlineIdx := strings.Index(text[start:], "\n"); newlineIdx != -1 {
		end = start + newlineIdx
	}

	value := strings.TrimSpace(text[start:end])

	// Remove any trailing punctuation or markup
	value = strings.TrimRight(value, ".,;:)")

	return value
}

// extractListItemText removes the bullet marker from a list item and
// returns the content.
//
//nolint:revive // line-length-limit
func extractListItemText(text string) string {
	result := strings.TrimSpace(
		text,
	) //nolint:revive // modifies-parameter

	// Skip leading bullet characters (-, *, +)
	if result != "" &&
		(result[0] == '-' || result[0] == '*' || result[0] == '+') {
		result = strings.TrimSpace(result[1:])
	} else {
		// Check for ordered list (number followed by dot)
		for i, ch := range result {
			if ch >= '0' && ch <= '9' {
				continue
			}
			if ch == '.' {
				result = strings.TrimSpace(result[i+1:])

				break
			}

			break
		}
	}

	// Handle checkboxes [ ] or [x]
	if strings.HasPrefix(result, "[ ]") ||
		strings.HasPrefix(result, "[x]") ||
		strings.HasPrefix(result, "[X]") {
		result = strings.TrimSpace(result[3:])
	}

	return result
}

// GetDeltaRequirementNames returns all requirement names affected by the delta.
// This includes added, modified, removed, and renamed requirements.
func GetDeltaRequirementNames(
	delta *Delta,
) []string {
	names := make([]string, 0)

	for _, req := range delta.Added {
		names = append(names, req.Name)
	}
	for _, req := range delta.Modified {
		names = append(names, req.Name)
	}
	names = append(names, delta.Removed...)
	for _, renamed := range delta.Renamed {
		if renamed.From != "" {
			names = append(names, renamed.From)
		}
		if renamed.To != "" {
			names = append(names, renamed.To)
		}
	}

	return names
}

// HasDeltaSection checks if the content contains a specific delta section type.
func HasDeltaSection(
	content []byte,
	deltaType DeltaType,
) bool {
	return FindDeltaSection(
		content,
		deltaType,
	) != ""
}

// CountDeltaChanges returns the total number of changes in a delta.
func CountDeltaChanges(delta *Delta) int {
	return len(
		delta.Added,
	) + len(
		delta.Modified,
	) + len(
		delta.Removed,
	) + len(
		delta.Renamed,
	)
}

// ValidateRenamed checks that all renamed requirements have both From and To
// values. Returns the names of any incomplete rename operations.
func ValidateRenamed(delta *Delta) []string {
	incomplete := make([]string, 0)
	for _, renamed := range delta.Renamed {
		if renamed.From == "" {
			incomplete = append(
				incomplete,
				"missing FROM for: "+renamed.To,
			)
		}
		if renamed.To == "" {
			incomplete = append(
				incomplete,
				"missing TO for: "+renamed.From,
			)
		}
	}

	return incomplete
}

// ExtractDeltaContent extracts content from a specific delta section type.
// This is useful when you need the raw markdown content of a delta section.
func ExtractDeltaContent(
	content []byte,
	deltaType DeltaType,
) []byte {
	sectionContent := FindDeltaSection(
		content,
		deltaType,
	)
	if sectionContent == "" {
		return nil
	}

	return []byte(sectionContent)
}

// MergeDelta merges delta changes into the merged delta.
// This is useful for combining multiple delta files.
func MergeDelta(merged, delta *Delta) {
	merged.Added = append(
		merged.Added,
		delta.Added...)
	merged.Modified = append(
		merged.Modified,
		delta.Modified...)
	merged.Removed = append(
		merged.Removed,
		delta.Removed...)
	merged.Renamed = append(
		merged.Renamed,
		delta.Renamed...)
	merged.Errors = append(
		merged.Errors,
		delta.Errors...)
}

// isDeltaSectionHeader checks if a header title indicates a delta section.
func isDeltaSectionHeader(
	title string,
) (DeltaType, bool) {
	upper := strings.ToUpper(
		strings.TrimSpace(title),
	)

	if strings.HasPrefix(upper, "ADDED") {
		return DeltaAdded, true
	}
	if strings.HasPrefix(upper, "MODIFIED") {
		return DeltaModified, true
	}
	if strings.HasPrefix(upper, "REMOVED") {
		return DeltaRemoved, true
	}
	if strings.HasPrefix(upper, "RENAMED") {
		return DeltaRenamed, true
	}

	return "", false
}

// FindRenamedPairs extracts all FROM/TO pairs from the content.
// This is a convenience function for getting just the rename mappings.
func FindRenamedPairs(
	content []byte,
) map[string]string {
	delta, _ := ParseDelta(content)
	pairs := make(map[string]string)

	for _, renamed := range delta.Renamed {
		if renamed.From != "" &&
			renamed.To != "" {
			pairs[renamed.From] = renamed.To
		}
	}

	return pairs
}

// FindAddedRequirements returns the names of all added requirements.
func FindAddedRequirements(
	content []byte,
) []string {
	delta, _ := ParseDelta(content)
	names := make([]string, len(delta.Added))
	for i, req := range delta.Added {
		names[i] = req.Name
	}

	return names
}

// FindModifiedRequirements returns the names of all modified requirements.
func FindModifiedRequirements(
	content []byte,
) []string {
	delta, _ := ParseDelta(content)
	names := make([]string, len(delta.Modified))
	for i, req := range delta.Modified {
		names[i] = req.Name
	}

	return names
}

// FindRemovedRequirements returns the names of all removed requirements.
func FindRemovedRequirements(
	content []byte,
) []string {
	delta, _ := ParseDelta(content)

	return delta.Removed
}

// GetDeltaSummary returns a human-readable summary of delta changes.
//
//nolint:revive // function-length - straightforward iteration over delta fields
func GetDeltaSummary(delta *Delta) string {
	var buf bytes.Buffer

	if len(delta.Added) > 0 {
		buf.WriteString("Added: ")
		for i, req := range delta.Added {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(req.Name)
		}
		buf.WriteString("\n")
	}

	if len(delta.Modified) > 0 {
		buf.WriteString("Modified: ")
		for i, req := range delta.Modified {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(req.Name)
		}
		buf.WriteString("\n")
	}

	if len(delta.Removed) > 0 {
		buf.WriteString("Removed: ")
		for i, name := range delta.Removed {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(name)
		}
		buf.WriteString("\n")
	}

	if len(delta.Renamed) > 0 {
		buf.WriteString("Renamed: ")
		for i, renamed := range delta.Renamed {
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(renamed.From)
			buf.WriteString(" -> ")
			buf.WriteString(renamed.To)
		}
		buf.WriteString("\n")
	}

	return buf.String()
}
