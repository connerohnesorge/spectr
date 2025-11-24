// Package mdparser provides a robust markdown parser built on a
// lexer/parser architecture.
//
// This package implements a two-phase parsing approach:
//  1. Lexer: Tokenizes markdown input into a stream of tokens
//     (headers, text, code blocks, lists)
//  2. Parser: Builds an Abstract Syntax Tree (AST) from the token
//     stream
//
// Design Philosophy:
//
// The parser is intentionally generic and does not encode
// Spectr-specific semantics.
// It provides a clean AST representation of markdown structure that
// can be traversed and interpreted by higher-level extractors
// (see internal/parsers/extractor.go).
//
// This separation allows the parser to be:
//   - Reusable across different markdown parsing needs
//   - Testable independent of business logic
//   - Maintainable without regex brittleness
//
// Supported Markdown Elements:
//   - Headers (H1-H6): ## Header Text
//   - Paragraphs: Regular text blocks
//   - Code Blocks: Fenced code blocks with optional language
//     identifiers (```go)
//   - Lists: Both ordered (1. item) and unordered (- item) lists
//   - Blank Lines: Preserved for structure reconstruction
//
// Limitations:
//   - No inline formatting (bold, italic, links) - these are
//     preserved as text
//   - No tables, blockquotes, or horizontal rules
//   - No nested list support (items are flat)
//
// Example Usage:
//
//	doc, err := mdparser.Parse(markdownContent)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Traverse the AST
//	for _, node := range doc.Children {
//	    if header, ok := node.(*mdparser.Header); ok {
//	        fmt.Printf("Found header: %s (level %d)\n",
//	            header.Text, header.Level)
//	    }
//	}
//
// Architecture:
//
// The lexer uses a state machine pattern where each state function
// returns the next state function to execute. This provides clean
// separation of concerns and makes it easy to handle
// context-dependent parsing (e.g., code blocks).
//
// The parser builds the AST in a single pass, using a two-token
// lookahead to make parsing decisions. Error recovery is minimal -
// the parser returns an error on the first malformed structure
// encountered.
package mdparser

// Node is the interface implemented by all AST nodes.
// This represents generic markdown structure without Spectr-specific semantics.
type Node interface {
	// Pos returns the starting position of this node in the source
	Pos() Position
	// End returns the ending position of this node in the source
	End() Position
	// String returns a string representation of the node (for debugging)
	String() string
}

// Position represents a location in the source file.
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column number
	Offset int // 0-based byte offset
}

// Document is the root node of the AST, containing all top-level elements.
type Document struct {
	StartPos Position
	EndPos   Position
	Children []Node
}

func (d *Document) Pos() Position { return d.StartPos }
func (d *Document) End() Position { return d.EndPos }
func (*Document) String() string {
	return "Document"
}

// Header represents a markdown header (e.g., "## Section Title").
type Header struct {
	StartPos Position
	EndPos   Position
	Level    int    // 1-6 for #, ##, ###, etc.
	Text     string // Header text without the # symbols
}

func (h *Header) Pos() Position { return h.StartPos }
func (h *Header) End() Position { return h.EndPos }
func (*Header) String() string {
	return "Header"
}

// Paragraph represents a block of text.
type Paragraph struct {
	StartPos Position
	EndPos   Position
	Lines    []string // Lines of text in this paragraph
}

func (p *Paragraph) Pos() Position { return p.StartPos }
func (p *Paragraph) End() Position { return p.EndPos }
func (*Paragraph) String() string {
	return "Paragraph"
}

// CodeBlock represents a fenced code block.
type CodeBlock struct {
	StartPos Position
	EndPos   Position
	Language string   // Language identifier (e.g., "go", "bash")
	Lines    []string // Lines of code content
}

func (c *CodeBlock) Pos() Position { return c.StartPos }
func (c *CodeBlock) End() Position { return c.EndPos }
func (*CodeBlock) String() string {
	return "CodeBlock"
}

// List represents a list (ordered or unordered).
type List struct {
	StartPos Position
	EndPos   Position
	Ordered  bool
	Items    []*ListItem
}

func (l *List) Pos() Position { return l.StartPos }
func (l *List) End() Position { return l.EndPos }
func (*List) String() string {
	return "List"
}

// ListItem represents a single item in a list.
type ListItem struct {
	StartPos Position
	EndPos   Position
	Text     string // Item text without the bullet/number
	Marker   string // Original list marker (e.g., "- ", "1. ")
	Children []Node // Nested content (paragraphs, sublists, etc.)
}

func (li *ListItem) Pos() Position { return li.StartPos }
func (li *ListItem) End() Position { return li.EndPos }
func (*ListItem) String() string {
	return "ListItem"
}

// BlankLine represents one or more blank lines (for structure preservation).
type BlankLine struct {
	StartPos Position
	EndPos   Position
	Count    int // Number of consecutive blank lines
}

func (b *BlankLine) Pos() Position { return b.StartPos }
func (b *BlankLine) End() Position { return b.EndPos }
func (*BlankLine) String() string {
	return "BlankLine"
}
