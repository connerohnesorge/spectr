//nolint:revive // file-length-limit: node types require comprehensive definitions
package markdown

import (
	"hash/fnv"
	"strconv"
)

// NodeType represents the type of an AST node.
// Each type corresponds to a different markdown construct.
type NodeType uint8

const (
	// Block-level node types

	// NodeTypeDocument is the root node containing all top-level nodes.
	NodeTypeDocument NodeType = iota
	// NodeTypeSection represents an H2+ header and its content.
	NodeTypeSection
	// NodeTypeRequirement represents a ### Requirement: header.
	NodeTypeRequirement
	// NodeTypeScenario represents a #### Scenario: header.
	NodeTypeScenario
	// NodeTypeParagraph represents paragraph content.
	NodeTypeParagraph
	// NodeTypeList represents an unordered or ordered list.
	NodeTypeList
	// NodeTypeListItem represents a single list item.
	NodeTypeListItem
	// NodeTypeCodeBlock represents a fenced code block.
	NodeTypeCodeBlock
	// NodeTypeBlockquote represents blockquoted content.
	NodeTypeBlockquote

	// Inline node types

	// NodeTypeText represents plain text content.
	NodeTypeText
	// NodeTypeStrong represents bold/strong emphasis (**text** or __text__).
	NodeTypeStrong
	// NodeTypeEmphasis represents italic emphasis (*text* or _text_).
	NodeTypeEmphasis
	// NodeTypeStrikethrough represents struck text (~~text~~).
	NodeTypeStrikethrough
	// NodeTypeCode represents inline code (`code`).
	NodeTypeCode
	// NodeTypeLink represents a link [text](url) or [text][ref].
	NodeTypeLink
	// NodeTypeLinkDef represents a link definition [ref]: url.
	NodeTypeLinkDef
	// NodeTypeWikilink represents a wikilink [[target|display#anchor]].
	NodeTypeWikilink
)

// String returns a human-readable name for the node type.
func (t NodeType) String() string {
	switch t {
	case NodeTypeDocument:
		return "Document"
	case NodeTypeSection:
		return "Section"
	case NodeTypeRequirement:
		return "Requirement"
	case NodeTypeScenario:
		return "Scenario"
	case NodeTypeParagraph:
		return "Paragraph"
	case NodeTypeList:
		return "List"
	case NodeTypeListItem:
		return "ListItem"
	case NodeTypeCodeBlock:
		return "CodeBlock"
	case NodeTypeBlockquote:
		return "Blockquote"
	case NodeTypeText:
		return "Text"
	case NodeTypeStrong:
		return "Strong"
	case NodeTypeEmphasis:
		return "Emphasis"
	case NodeTypeStrikethrough:
		return "Strikethrough"
	case NodeTypeCode:
		return "Code"
	case NodeTypeLink:
		return "Link"
	case NodeTypeLinkDef:
		return "LinkDef"
	case NodeTypeWikilink:
		return "Wikilink"
	default:
		return "Unknown"
	}
}

// Node is the interface implemented by all AST nodes.
// All nodes are immutable after creation - modifications require
// creating new nodes via builders.
type Node interface {
	// NodeType returns the type classification of this node.
	NodeType() NodeType

	// Span returns the byte offset range (start, end) of this node.
	// Start is inclusive, end is exclusive.
	Span() (start, end int)

	// Hash returns a content hash for identity tracking and caching.
	// Nodes with the same hash have the same semantic content.
	Hash() uint64

	// Source returns a zero-copy byte slice view into the original source.
	// This remains valid as long as the original source buffer is retained.
	Source() []byte

	// Children returns an immutable copy of child nodes.
	// Modifications to the returned slice do not affect the node.
	Children() []Node

	// Equal performs deep structural comparison with another node.
	Equal(other Node) bool
}

// baseNode contains the common fields shared by all node types.
// It is embedded in each concrete node struct.
type baseNode struct {
	nodeType NodeType
	hash     uint64
	start    int
	end      int
	source   []byte
	children []Node
}

// NodeType returns the type classification of this node.
func (n *baseNode) NodeType() NodeType {
	return n.nodeType
}

// Span returns the byte offset range (start, end) of this node.
func (n *baseNode) Span() (start, end int) {
	return n.start, n.end
}

// Hash returns the content hash for identity tracking and caching.
func (n *baseNode) Hash() uint64 {
	return n.hash
}

// Source returns a zero-copy byte slice view into the original source.
func (n *baseNode) Source() []byte {
	return n.source
}

// Children returns an immutable copy of child nodes.
func (n *baseNode) Children() []Node {
	if n.children == nil {
		return nil
	}
	// Return a copy to prevent external modification
	result := make([]Node, len(n.children))
	copy(result, n.children)

	return result
}

// computeHash computes a content hash for a node using FNV-1a.
// The hash includes: NodeType, children hashes, and source content.
func computeHash(
	nodeType NodeType,
	children []Node,
	source []byte,
) uint64 {
	h := fnv.New64a()

	// Hash the node type
	h.Write([]byte{byte(nodeType)})

	// Hash children hashes
	for _, child := range children {
		childHash := child.Hash()
		// Write hash bytes in big-endian order
		h.Write([]byte{
			byte(childHash >> 56),
			byte(childHash >> 48),
			byte(childHash >> 40),
			byte(childHash >> 32),
			byte(childHash >> 24),
			byte(childHash >> 16),
			byte(childHash >> 8),
			byte(childHash),
		})
	}

	// Hash source content
	h.Write(source)

	return h.Sum64()
}

// computeHashWithExtra computes a hash that includes extra type-specific data.
// Used for nodes with fields beyond the common baseNode fields.
func computeHashWithExtra(
	nodeType NodeType,
	children []Node,
	source []byte,
	extra []byte,
) uint64 {
	h := fnv.New64a()

	// Hash the node type
	h.Write([]byte{byte(nodeType)})

	// Hash children hashes
	for _, child := range children {
		childHash := child.Hash()
		h.Write([]byte{
			byte(childHash >> 56),
			byte(childHash >> 48),
			byte(childHash >> 40),
			byte(childHash >> 32),
			byte(childHash >> 24),
			byte(childHash >> 16),
			byte(childHash >> 8),
			byte(childHash),
		})
	}

	// Hash source content
	h.Write(source)

	// Hash extra type-specific data
	h.Write(extra)

	return h.Sum64()
}

// equalNodes performs deep structural comparison of two nodes.
func equalNodes(a, b Node) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Quick check: same type
	if a.NodeType() != b.NodeType() {
		return false
	}

	// Quick check: hash comparison (fast path)
	if a.Hash() != b.Hash() {
		return false
	}

	// Deep check: compare source
	aSource := a.Source()
	bSource := b.Source()
	if len(aSource) != len(bSource) {
		return false
	}
	for i := range aSource {
		if aSource[i] != bSource[i] {
			return false
		}
	}

	// Deep check: compare children recursively
	aChildren := a.Children()
	bChildren := b.Children()
	if len(aChildren) != len(bChildren) {
		return false
	}
	for i := range aChildren {
		if !equalNodes(
			aChildren[i],
			bChildren[i],
		) {
			return false
		}
	}

	return true
}

// NodeBuilder provides a fluent API for constructing AST nodes.
// It supports building new nodes and transforming existing ones.
type NodeBuilder struct {
	nodeType NodeType
	start    int
	end      int
	source   []byte
	children []Node

	// Type-specific fields
	level     int    // for Section
	title     []byte // for Section
	deltaType string // for Section
	name      string // for Requirement, Scenario
	language  []byte // for CodeBlock
	content   []byte // for CodeBlock
	ordered   bool   // for List
	checked   *bool  // for ListItem
	keyword   string // for ListItem
	url       []byte // for Link
	linkTitle []byte // for Link
	target    []byte // for Wikilink
	display   []byte // for Wikilink
	anchor    []byte // for Wikilink
}

// NewNodeBuilder creates a new builder for the specified node type.
func NewNodeBuilder(
	nodeType NodeType,
) *NodeBuilder {
	return &NodeBuilder{
		nodeType: nodeType,
	}
}

// WithStart sets the start byte offset.
func (b *NodeBuilder) WithStart(
	start int,
) *NodeBuilder {
	b.start = start

	return b
}

// WithEnd sets the end byte offset.
func (b *NodeBuilder) WithEnd(
	end int,
) *NodeBuilder {
	b.end = end

	return b
}

// WithSource sets the source byte slice.
func (b *NodeBuilder) WithSource(
	source []byte,
) *NodeBuilder {
	b.source = source

	return b
}

// WithChildren sets the children nodes.
func (b *NodeBuilder) WithChildren(
	children []Node,
) *NodeBuilder {
	// Make a defensive copy
	if children != nil {
		b.children = make([]Node, len(children))
		copy(b.children, children)
	} else {
		b.children = nil
	}

	return b
}

// WithLevel sets the header level (for Section nodes).
func (b *NodeBuilder) WithLevel(
	level int,
) *NodeBuilder {
	b.level = level

	return b
}

// WithTitle sets the title (for Section nodes).
func (b *NodeBuilder) WithTitle(
	title []byte,
) *NodeBuilder {
	b.title = title

	return b
}

// WithDeltaType sets the delta type (for Section nodes).
func (b *NodeBuilder) WithDeltaType(
	deltaType string,
) *NodeBuilder {
	b.deltaType = deltaType

	return b
}

// WithName sets the name (for Requirement and Scenario nodes).
func (b *NodeBuilder) WithName(
	name string,
) *NodeBuilder {
	b.name = name

	return b
}

// WithLanguage sets the language (for CodeBlock nodes).
func (b *NodeBuilder) WithLanguage(
	language []byte,
) *NodeBuilder {
	b.language = language

	return b
}

// WithContent sets the code content (for CodeBlock nodes).
func (b *NodeBuilder) WithContent(
	content []byte,
) *NodeBuilder {
	b.content = content

	return b
}

// WithOrdered sets whether the list is ordered (for List nodes).
func (b *NodeBuilder) WithOrdered(
	ordered bool,
) *NodeBuilder {
	b.ordered = ordered

	return b
}

// WithChecked sets the checkbox state (for ListItem nodes).
// Pass nil for no checkbox, pointer to bool for checkbox state.
func (b *NodeBuilder) WithChecked(
	checked *bool,
) *NodeBuilder {
	b.checked = checked

	return b
}

// WithKeyword sets the keyword (for ListItem nodes).
func (b *NodeBuilder) WithKeyword(
	keyword string,
) *NodeBuilder {
	b.keyword = keyword

	return b
}

// WithURL sets the URL (for Link nodes).
func (b *NodeBuilder) WithURL(
	url []byte,
) *NodeBuilder {
	b.url = url

	return b
}

// WithLinkTitle sets the link title (for Link nodes).
func (b *NodeBuilder) WithLinkTitle(
	title []byte,
) *NodeBuilder {
	b.linkTitle = title

	return b
}

// WithTarget sets the target (for Wikilink nodes).
func (b *NodeBuilder) WithTarget(
	target []byte,
) *NodeBuilder {
	b.target = target

	return b
}

// WithDisplay sets the display text (for Wikilink nodes).
func (b *NodeBuilder) WithDisplay(
	display []byte,
) *NodeBuilder {
	b.display = display

	return b
}

// WithAnchor sets the anchor (for Wikilink nodes).
func (b *NodeBuilder) WithAnchor(
	anchor []byte,
) *NodeBuilder {
	b.anchor = anchor

	return b
}

// Validate checks that the builder state is valid.
// Returns an error if validation fails.
func (b *NodeBuilder) Validate() error {
	// Validate Start <= End
	if b.start > b.end {
		return &BuilderValidationError{
			Field:   "Start/End",
			Message: "Start must be less than or equal to End",
		}
	}

	// Validate children nesting (children must be within parent span)
	for i, child := range b.children {
		if child == nil {
			return &BuilderValidationError{
				Field:   "Children",
				Message: "Child node cannot be nil",
			}
		}
		childStart, childEnd := child.Span()
		if childStart < b.start ||
			childEnd > b.end {
			return &BuilderValidationError{
				Field:   "Children",
				Message: "Child node span must be within parent span",
				Index:   i,
			}
		}
	}

	// Type-specific validation
	if b.nodeType == NodeTypeSection {
		if b.level < 1 || b.level > 6 {
			return &BuilderValidationError{
				Field:   "Level",
				Message: "Section level must be between 1 and 6",
			}
		}
	}

	return nil
}

// Build creates the immutable node from the builder state.
// It validates the builder and computes the content hash.
// Returns nil if validation fails (caller should check Validate() first for error details).
//
//nolint:revive // function-length: builder pattern requires type dispatch
func (b *NodeBuilder) Build() Node {
	if err := b.Validate(); err != nil {
		return nil
	}

	// Make defensive copy of children
	var children []Node
	if b.children != nil {
		children = make([]Node, len(b.children))
		copy(children, b.children)
	}

	base := baseNode{
		nodeType: b.nodeType,
		start:    b.start,
		end:      b.end,
		source:   b.source,
		children: children,
	}

	switch b.nodeType {
	case NodeTypeDocument:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeDocument{baseNode: base}

	case NodeTypeSection:
		extra := make(
			[]byte,
			0,
			1+len(b.title)+len(b.deltaType),
		)
		extra = append(extra, byte(b.level))
		extra = append(extra, b.title...)
		extra = append(
			extra,
			[]byte(b.deltaType)...)
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			extra,
		)

		return &NodeSection{
			baseNode:  base,
			level:     b.level,
			title:     b.title,
			deltaType: b.deltaType,
		}

	case NodeTypeRequirement:
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			[]byte(b.name),
		)

		return &NodeRequirement{
			baseNode: base,
			name:     b.name,
		}

	case NodeTypeScenario:
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			[]byte(b.name),
		)

		return &NodeScenario{
			baseNode: base,
			name:     b.name,
		}

	case NodeTypeParagraph:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeParagraph{baseNode: base}

	case NodeTypeList:
		extra := []byte{0}
		if b.ordered {
			extra[0] = 1
		}
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			extra,
		)

		return &NodeList{
			baseNode: base,
			ordered:  b.ordered,
		}

	case NodeTypeListItem:
		extra := make([]byte, 0, 2+len(b.keyword))
		// Encode checked state: 0 = no checkbox, 1 = unchecked, 2 = checked
		switch {
		case b.checked == nil:
			extra = append(extra, 0)
		case *b.checked:
			extra = append(extra, 2)
		default:
			extra = append(extra, 1)
		}
		extra = append(
			extra,
			[]byte(b.keyword)...)
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			extra,
		)

		return &NodeListItem{
			baseNode: base,
			checked:  b.checked,
			keyword:  b.keyword,
		}

	case NodeTypeCodeBlock:
		extra := make(
			[]byte,
			0,
			len(b.language)+1+len(b.content),
		)
		extra = append(extra, b.language...)
		extra = append(extra, 0) // separator
		extra = append(extra, b.content...)
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			extra,
		)

		return &NodeCodeBlock{
			baseNode: base,
			language: b.language,
			content:  b.content,
		}

	case NodeTypeBlockquote:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeBlockquote{baseNode: base}

	case NodeTypeText:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeText{baseNode: base}

	case NodeTypeStrong:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeStrong{baseNode: base}

	case NodeTypeEmphasis:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeEmphasis{baseNode: base}

	case NodeTypeStrikethrough:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeStrikethrough{baseNode: base}

	case NodeTypeCode:
		base.hash = computeHash(
			b.nodeType,
			children,
			b.source,
		)

		return &NodeCode{baseNode: base}

	case NodeTypeLink:
		extra := make(
			[]byte,
			0,
			len(b.url)+1+len(b.linkTitle),
		)
		extra = append(extra, b.url...)
		extra = append(extra, 0) // separator
		extra = append(extra, b.linkTitle...)
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			extra,
		)

		return &NodeLink{
			baseNode: base,
			url:      b.url,
			title:    b.linkTitle,
		}

	case NodeTypeLinkDef:
		extra := make(
			[]byte,
			0,
			len(b.url)+1+len(b.linkTitle),
		)
		extra = append(extra, b.url...)
		extra = append(extra, 0)
		extra = append(extra, b.linkTitle...)
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			extra,
		)

		return &NodeLinkDef{
			baseNode: base,
			url:      b.url,
			title:    b.linkTitle,
		}

	case NodeTypeWikilink:
		extra := make(
			[]byte,
			0,
			len(
				b.target,
			)+1+len(
				b.display,
			)+1+len(
				b.anchor,
			),
		)
		extra = append(extra, b.target...)
		extra = append(extra, 0)
		extra = append(extra, b.display...)
		extra = append(extra, 0)
		extra = append(extra, b.anchor...)
		base.hash = computeHashWithExtra(
			b.nodeType,
			children,
			b.source,
			extra,
		)

		return &NodeWikilink{
			baseNode: base,
			target:   b.target,
			display:  b.display,
			anchor:   b.anchor,
		}

	default:
		return nil
	}
}

// BuilderValidationError represents a validation error in the builder.
type BuilderValidationError struct {
	Field   string
	Message string
	Index   int // For array fields, the index of the problematic element
}

func (e *BuilderValidationError) Error() string {
	if e.Index > 0 {
		return e.Field + "[" + strconv.Itoa(
			e.Index,
		) + "]: " + e.Message
	}

	return e.Field + ": " + e.Message
}

// nodeToBuilder creates a builder from an existing node for transformation.
// This is implemented by each concrete node type.
func nodeToBuilder(n Node) *NodeBuilder {
	if n == nil {
		return nil
	}

	start, end := n.Span()
	b := &NodeBuilder{
		nodeType: n.NodeType(),
		start:    start,
		end:      end,
		source:   n.Source(),
		children: n.Children(),
	}

	// Copy type-specific fields
	switch node := n.(type) {
	case *NodeSection:
		b.level = node.level
		b.title = node.title
		b.deltaType = node.deltaType
	case *NodeRequirement:
		b.name = node.name
	case *NodeScenario:
		b.name = node.name
	case *NodeList:
		b.ordered = node.ordered
	case *NodeListItem:
		b.checked = node.checked
		b.keyword = node.keyword
	case *NodeCodeBlock:
		b.language = node.language
		b.content = node.content
	case *NodeLink:
		b.url = node.url
		b.linkTitle = node.title
	case *NodeLinkDef:
		b.url = node.url
		b.linkTitle = node.title
	case *NodeWikilink:
		b.target = node.target
		b.display = node.display
		b.anchor = node.anchor
	}

	return b
}
