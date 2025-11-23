package syntax

// Node is the interface for all AST nodes.
type Node interface {
	Pos() int
	RawString() string
}

// Document is the root node of the AST.
type Document struct {
	Nodes []Node
}

func (d *Document) Pos() int          { return 0 }
func (d *Document) RawString() string { return "" }

// Header represents a markdown header.
type Header struct {
	Level int
	Text  string
	Raw   string
	Line  int
}

func (h *Header) Pos() int          { return h.Line }
func (h *Header) RawString() string { return h.Raw }

// Text represents a block of text.
type Text struct {
	Content string
	Raw     string
	Line    int
}

func (t *Text) Pos() int          { return t.Line }
func (t *Text) RawString() string { return t.Raw }

// CodeBlock represents a code block.
type CodeBlock struct {
	Language string
	Content  string
	Raw      string
	Line     int
}

func (c *CodeBlock) Pos() int          { return c.Line }
func (c *CodeBlock) RawString() string { return c.Raw }

// List represents a list item.
type List struct {
	Marker  string
	Content string
	Raw     string
	Line    int
}

func (l *List) Pos() int          { return l.Line }
func (l *List) RawString() string { return l.Raw }
