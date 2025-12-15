//nolint:revive // max-public-structs - node types intentionally public for AST API
package markdown

// NodeDocument is the root node of an AST.
// It contains all top-level block nodes from the parsed document.
type NodeDocument struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeDocument) Equal(other Node) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeDocument) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeSection represents an ATX-style header (H1-H6) and its content.
// For Spectr delta files, it may have a DeltaType indicating the change type.
type NodeSection struct {
	baseNode
	level     int    // Header level (1-6)
	title     []byte // Header text
	deltaType string // "ADDED", "MODIFIED", "REMOVED", "RENAMED", or ""
}

// Level returns the header level (1-6).
func (n *NodeSection) Level() int {
	return n.level
}

// Title returns the header text as a byte slice.
func (n *NodeSection) Title() []byte {
	return n.title
}

// DeltaType returns the delta type for Spectr delta sections.
// Returns one of "ADDED", "MODIFIED", "REMOVED", "RENAMED", or empty string.
func (n *NodeSection) DeltaType() string {
	return n.deltaType
}

// Equal performs deep structural comparison with another node.
func (n *NodeSection) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherSection, ok := other.(*NodeSection)
	if !ok {
		return false
	}
	if n.level != otherSection.level {
		return false
	}
	if n.deltaType != otherSection.deltaType {
		return false
	}
	if !bytesEqual(n.title, otherSection.title) {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeSection) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeRequirement represents a Spectr requirement header (### Requirement: Name).
type NodeRequirement struct {
	baseNode
	name string
}

// Name returns the requirement name extracted from the header.
func (n *NodeRequirement) Name() string {
	return n.name
}

// Equal performs deep structural comparison with another node.
func (n *NodeRequirement) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherReq, ok := other.(*NodeRequirement)
	if !ok {
		return false
	}
	if n.name != otherReq.name {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeRequirement) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeScenario represents a Spectr scenario header (#### Scenario: Name).
type NodeScenario struct {
	baseNode
	name string
}

// Name returns the scenario name extracted from the header.
func (n *NodeScenario) Name() string {
	return n.name
}

// Equal performs deep structural comparison with another node.
func (n *NodeScenario) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherScenario, ok := other.(*NodeScenario)
	if !ok {
		return false
	}
	if n.name != otherScenario.name {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeScenario) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeParagraph represents a paragraph of text.
// Its children are inline nodes (text, emphasis, links, etc.).
type NodeParagraph struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeParagraph) Equal(other Node) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeParagraph) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeList represents an unordered or ordered list.
// Its children are NodeListItem nodes.
type NodeList struct {
	baseNode
	ordered bool
}

// Ordered returns true if this is an ordered (numbered) list.
func (n *NodeList) Ordered() bool {
	return n.ordered
}

// Equal performs deep structural comparison with another node.
func (n *NodeList) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherList, ok := other.(*NodeList)
	if !ok {
		return false
	}
	if n.ordered != otherList.ordered {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeList) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeListItem represents a single list item.
// It may be a task item with a checkbox, or a Spectr bullet with WHEN/THEN/AND keyword.
type NodeListItem struct {
	baseNode
	checked *bool  // nil = no checkbox, *true = checked, *false = unchecked
	keyword string // "WHEN", "THEN", "AND", or ""
}

// Checked returns the checkbox state and whether a checkbox is present.
// If no checkbox, returns (false, false).
// If unchecked checkbox, returns (false, true).
// If checked checkbox, returns (true, true).
func (n *NodeListItem) Checked() (isChecked, hasCheckbox bool) {
	if n.checked == nil {
		return false, false
	}

	return *n.checked, true
}

// Keyword returns the Spectr bullet keyword (WHEN, THEN, AND) or empty string.
func (n *NodeListItem) Keyword() string {
	return n.keyword
}

// Equal performs deep structural comparison with another node.
func (n *NodeListItem) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherItem, ok := other.(*NodeListItem)
	if !ok {
		return false
	}
	// Compare checked state
	if (n.checked == nil) != (otherItem.checked == nil) {
		return false
	}
	if n.checked != nil &&
		*n.checked != *otherItem.checked {
		return false
	}
	if n.keyword != otherItem.keyword {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeListItem) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeCodeBlock represents a fenced code block (``` or ~~~).
type NodeCodeBlock struct {
	baseNode
	language []byte // Language identifier (may be nil)
	content  []byte // Code content (without fences)
}

// Language returns the language identifier as a byte slice.
// Returns nil if no language was specified.
func (n *NodeCodeBlock) Language() []byte {
	return n.language
}

// Content returns the code content (without fence lines) as a byte slice.
func (n *NodeCodeBlock) Content() []byte {
	return n.content
}

// Equal performs deep structural comparison with another node.
func (n *NodeCodeBlock) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherBlock, ok := other.(*NodeCodeBlock)
	if !ok {
		return false
	}
	if !bytesEqual(
		n.language,
		otherBlock.language,
	) {
		return false
	}
	if !bytesEqual(
		n.content,
		otherBlock.content,
	) {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeCodeBlock) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeBlockquote represents blockquoted content (lines starting with >).
type NodeBlockquote struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeBlockquote) Equal(other Node) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeBlockquote) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeText represents plain text content.
type NodeText struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeText) Equal(other Node) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeText) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// Text returns the text content as a string.
// This creates a copy; use Source() for zero-copy access.
func (n *NodeText) Text() string {
	return string(n.source)
}

// NodeStrong represents bold/strong emphasis (**text** or __text__).
type NodeStrong struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeStrong) Equal(other Node) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeStrong) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeEmphasis represents italic emphasis (*text* or _text_).
type NodeEmphasis struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeEmphasis) Equal(other Node) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeEmphasis) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeStrikethrough represents struck text (~~text~~).
type NodeStrikethrough struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeStrikethrough) Equal(
	other Node,
) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeStrikethrough) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeCode represents inline code (`code`).
type NodeCode struct {
	baseNode
}

// Equal performs deep structural comparison with another node.
func (n *NodeCode) Equal(other Node) bool {
	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeCode) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// Code returns the code content as a string.
// This creates a copy; use Source() for zero-copy access.
func (n *NodeCode) Code() string {
	return string(n.source)
}

// NodeLink represents a link [text](url "title") or [text][ref].
type NodeLink struct {
	baseNode
	url   []byte
	title []byte
}

// URL returns the link destination as a byte slice.
func (n *NodeLink) URL() []byte {
	return n.url
}

// Title returns the optional link title as a byte slice.
// Returns nil if no title was specified.
func (n *NodeLink) Title() []byte {
	return n.title
}

// Equal performs deep structural comparison with another node.
func (n *NodeLink) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherLink, ok := other.(*NodeLink)
	if !ok {
		return false
	}
	if !bytesEqual(n.url, otherLink.url) {
		return false
	}
	if !bytesEqual(n.title, otherLink.title) {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeLink) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeLinkDef represents a link definition [ref]: url "title".
// Link definitions are collected during parsing and used to resolve reference links.
type NodeLinkDef struct {
	baseNode
	url   []byte
	title []byte
}

// URL returns the link destination as a byte slice.
func (n *NodeLinkDef) URL() []byte {
	return n.url
}

// Title returns the optional link title as a byte slice.
// Returns nil if no title was specified.
func (n *NodeLinkDef) Title() []byte {
	return n.title
}

// Equal performs deep structural comparison with another node.
func (n *NodeLinkDef) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherDef, ok := other.(*NodeLinkDef)
	if !ok {
		return false
	}
	if !bytesEqual(n.url, otherDef.url) {
		return false
	}
	if !bytesEqual(n.title, otherDef.title) {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeLinkDef) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// NodeWikilink represents a wikilink [[target|display#anchor]].
type NodeWikilink struct {
	baseNode
	target  []byte // Link target (e.g., "validation" or "changes/my-change")
	display []byte // Optional display text (may be nil)
	anchor  []byte // Optional anchor (may be nil)
}

// Target returns the link target as a byte slice.
func (n *NodeWikilink) Target() []byte {
	return n.target
}

// Display returns the optional display text as a byte slice.
// Returns nil if no display text was specified (defaults to target).
func (n *NodeWikilink) Display() []byte {
	return n.display
}

// Anchor returns the optional anchor as a byte slice.
// Returns nil if no anchor was specified.
func (n *NodeWikilink) Anchor() []byte {
	return n.anchor
}

// Equal performs deep structural comparison with another node.
func (n *NodeWikilink) Equal(other Node) bool {
	if other == nil {
		return false
	}
	otherWikilink, ok := other.(*NodeWikilink)
	if !ok {
		return false
	}
	if !bytesEqual(
		n.target,
		otherWikilink.target,
	) {
		return false
	}
	if !bytesEqual(
		n.display,
		otherWikilink.display,
	) {
		return false
	}
	if !bytesEqual(
		n.anchor,
		otherWikilink.anchor,
	) {
		return false
	}

	return equalNodes(n, other)
}

// ToBuilder creates a builder pre-populated with this node's data.
func (n *NodeWikilink) ToBuilder() *NodeBuilder {
	return nodeToBuilder(n)
}

// bytesEqual compares two byte slices for equality.
// Handles nil slices correctly.
func bytesEqual(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
