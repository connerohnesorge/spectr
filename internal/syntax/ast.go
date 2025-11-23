package syntax

// Node is the interface for all AST nodes.
type Node interface {
	Pos() int
}

// Document is the root node of the AST.
type Document struct {
	Nodes []Node
}

func (d *Document) Pos() int { return 0 }

// Header represents a markdown header.
type Header struct {
	Level int
	Text  string
	Raw   string
	Line  int
}

func (h *Header) Pos() int { return h.Line }

// Text represents a block of text.
type Text struct {
	Content string
	Raw     string
	Line    int
}

func (t *Text) Pos() int { return t.Line }

// CodeBlock represents a code block.
type CodeBlock struct {
	Language string
	Content  string
	Raw      string
	Line     int
}

func (c *CodeBlock) Pos() int { return c.Line }

// List represents a list item.
type List struct {
	Marker  string
	Content string
	Raw     string
	Line    int
}

func (l *List) Pos() int { return l.Line }
