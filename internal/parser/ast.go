package parser

//nolint:revive // file length acceptable for AST node definitions

import (
	"fmt"
	"strings"
)

// NodeType represents the type of an AST node.
type NodeType int

const (
	// NodeDocument represents the root document node
	NodeDocument NodeType = iota

	// NodeHeader represents a markdown header
	NodeHeader

	// NodeParagraph represents a block of text
	NodeParagraph

	// NodeCodeBlock represents a code fence
	NodeCodeBlock

	// NodeList represents a list structure
	NodeList

	// NodeBlankLine represents blank lines
	NodeBlankLine

	// maxDisplayLength is the maximum text length before truncation
	maxDisplayLength = 30
)

// String returns the string representation of a NodeType.
func (n NodeType) String() string {
	switch n {
	case NodeDocument:
		return "Document"
	case NodeHeader:
		return "Header"
	case NodeParagraph:
		return "Paragraph"
	case NodeCodeBlock:
		return "CodeBlock"
	case NodeList:
		return "List"
	case NodeBlankLine:
		return "BlankLine"
	default:
		return fmt.Sprintf("Unknown(%d)", n)
	}
}

// Node is the interface that all AST nodes must implement.
//
// All nodes have a type identifier and position information for
// error reporting and traversal.
type Node interface {
	// Type returns the type of this node
	Type() NodeType

	// Pos returns the position where this node starts
	Pos() Position

	// String returns a human-readable representation of the node
	String() string
}

// Document represents the root node of the parsed markdown document.
//
// It contains all top-level nodes in the document as children.
type Document struct {
	Children []Node   // Top-level nodes in the document
	position Position // Position is always 1:1 for document root
}

// Type returns NodeDocument.
func (*Document) Type() NodeType {
	return NodeDocument
}

// Pos returns the position where the document starts (always 1:1).
func (d *Document) Pos() Position {
	return d.position
}

// String returns a debug representation of the document.
func (d *Document) String() string {
	return fmt.Sprintf("Document{Children: %d nodes}", len(d.Children))
}

// NewDocument creates a new Document node.
func NewDocument() *Document {
	return &Document{
		Children: make([]Node, 0),
		position: Position{Line: 1, Column: 1, Offset: 0},
	}
}

// Header represents a markdown header (# ## ### etc).
//
// Headers have a level (1-6) and text content. They are leaf nodes
// that don't contain other nodes.
type Header struct {
	Level    int      // Header level (1-6)
	Text     string   // Header text (without # markers)
	position Position // Position where the header starts
}

// Type returns NodeHeader.
func (*Header) Type() NodeType {
	return NodeHeader
}

// Pos returns the position where the header starts.
func (h *Header) Pos() Position {
	return h.position
}

// String returns a debug representation of the header.
func (h *Header) String() string {
	if len(h.Text) > maxDisplayLength {
		return fmt.Sprintf(
			"Header{Level: %d, Text: %.30s... @%s}",
			h.Level, h.Text, h.position,
		)
	}

	return fmt.Sprintf(
		"Header{Level: %d, Text: %s @%s}",
		h.Level, h.Text, h.position,
	)
}

// NewHeader creates a new Header node.
func NewHeader(level int, text string, pos Position) *Header {
	return &Header{
		Level:    level,
		Text:     text,
		position: pos,
	}
}

// Paragraph represents a block of text content.
//
// Paragraphs are non-header, non-code text content. They are leaf nodes
// that don't contain other nodes.
type Paragraph struct {
	Text     string   // Paragraph text content
	position Position // Position where the paragraph starts
}

// Type returns NodeParagraph.
func (*Paragraph) Type() NodeType {
	return NodeParagraph
}

// Pos returns the position where the paragraph starts.
func (p *Paragraph) Pos() Position {
	return p.position
}

// String returns a debug representation of the paragraph.
func (p *Paragraph) String() string {
	text := strings.TrimSpace(p.Text)
	if len(text) > maxDisplayLength {
		return fmt.Sprintf("Paragraph{%.30s... @%s}", text, p.position)
	}

	return fmt.Sprintf("Paragraph{%s @%s}", text, p.position)
}

// NewParagraph creates a new Paragraph node.
func NewParagraph(text string, pos Position) *Paragraph {
	return &Paragraph{
		Text:     text,
		position: pos,
	}
}

// CodeBlock represents a code fence (```).
//
// Code blocks have an optional language specifier and content.
// They are leaf nodes that don't contain other nodes.
type CodeBlock struct {
	Language string   // Optional language specifier (e.g., "go", "python")
	Content  string   // Code content
	position Position // Position where the code block starts
}

// Type returns NodeCodeBlock.
func (*CodeBlock) Type() NodeType {
	return NodeCodeBlock
}

// Pos returns the position where the code block starts.
func (c *CodeBlock) Pos() Position {
	return c.position
}

// String returns a debug representation of the code block.
func (c *CodeBlock) String() string {
	content := strings.TrimSpace(c.Content)
	if len(content) > maxDisplayLength {
		content = content[:maxDisplayLength] + "..."
	}
	if c.Language != "" {
		return fmt.Sprintf(
			"CodeBlock{Lang: %s, Content: %s @%s}",
			c.Language, content, c.position,
		)
	}

	return fmt.Sprintf(
		"CodeBlock{Content: %s @%s}",
		content, c.position,
	)
}

// NewCodeBlock creates a new CodeBlock node.
func NewCodeBlock(language, content string, pos Position) *CodeBlock {
	return &CodeBlock{
		Language: language,
		Content:  content,
		position: pos,
	}
}

// List represents a list structure.
//
// List nodes contain the text of list items. In the current design,
// each list item is a separate List node (one node per item).
// This matches the lexer's TokenListItem behavior.
type List struct {
	// List item texts (for grouped lists, typically 1 item per node)
	Items    []string
	position Position // Position where the list starts
}

// Type returns NodeList.
func (*List) Type() NodeType {
	return NodeList
}

// Pos returns the position where the list starts.
func (l *List) Pos() Position {
	return l.position
}

// String returns a debug representation of the list.
func (l *List) String() string {
	if len(l.Items) == 0 {
		return fmt.Sprintf("List{Empty @%s}", l.position)
	}
	if len(l.Items) == 1 {
		item := l.Items[0]
		if len(item) > maxDisplayLength {
			item = item[:maxDisplayLength] + "..."
		}

		return fmt.Sprintf("List{%s @%s}", item, l.position)
	}

	return fmt.Sprintf("List{Items: %d @%s}", len(l.Items), l.position)
}

// NewList creates a new List node with a single item.
func NewList(item string, pos Position) *List {
	return &List{
		Items:    []string{item},
		position: pos,
	}
}

// BlankLine represents one or more blank lines.
//
// Blank lines help preserve document structure and are useful for
// determining paragraph boundaries.
type BlankLine struct {
	Count    int      // Number of consecutive blank lines
	position Position // Position where the blank lines start
}

// Type returns NodeBlankLine.
func (*BlankLine) Type() NodeType {
	return NodeBlankLine
}

// Pos returns the position where the blank lines start.
func (b *BlankLine) Pos() Position {
	return b.position
}

// String returns a debug representation of the blank line.
func (b *BlankLine) String() string {
	if b.Count == 1 {
		return fmt.Sprintf("BlankLine @%s", b.position)
	}

	return fmt.Sprintf("BlankLine{Count: %d @%s}", b.Count, b.position)
}

// NewBlankLine creates a new BlankLine node.
func NewBlankLine(count int, pos Position) *BlankLine {
	return &BlankLine{
		Count:    count,
		position: pos,
	}
}

// Walk traverses the AST and calls the visitor function for each node.
//
// The visitor function is called with each node in depth-first order.
// If the visitor returns false, traversal is stopped immediately.
//
// This is useful for the Extractor to search for specific patterns
// in the document structure.
func Walk(node Node, visitor func(Node) bool) bool {
	if node == nil {
		return true
	}

	// Visit the current node
	if !visitor(node) {
		return false
	}

	// Recurse into children if this is a Document
	if doc, ok := node.(*Document); ok {
		for _, child := range doc.Children {
			if !Walk(child, visitor) {
				return false
			}
		}
	}

	return true
}

// FindHeaders returns all Header nodes matching the predicate.
//
// This is a convenience for finding specific headers.
func FindHeaders(node Node, predicate func(*Header) bool) []*Header {
	var headers []*Header

	Walk(node, func(n Node) bool {
		if h, ok := n.(*Header); ok {
			if predicate(h) {
				headers = append(headers, h)
			}
		}

		return true
	})

	return headers
}

//nolint:revive // file-length acceptable for AST definitions

// NodesBetween returns all nodes between startPos and endPos in the document.
//
// This is useful for extracting content between two headers.
// Nodes are returned in document order.
func NodesBetween(doc *Document, startPos, endPos Position) []Node {
	var nodes []Node
	inRange := false

	Walk(doc, func(n Node) bool {
		pos := n.Pos()

		// Start collecting after startPos
		if !inRange && (pos.Offset >= startPos.Offset) {
			inRange = true
		}

		// Stop before endPos
		if inRange && (pos.Offset >= endPos.Offset) {
			return false
		}

		// Collect nodes in range (excluding the start node itself)
		if inRange && pos.Offset > startPos.Offset {
			nodes = append(nodes, n)
		}

		return true
	})

	return nodes
}

//nolint:revive // file length acceptable for AST node definitions
