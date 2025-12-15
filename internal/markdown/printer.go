package markdown

import (
	"bytes"
	"io"
	"strings"
)

// Print renders the AST node to normalized markdown.
// It returns the rendered markdown as a byte slice.
// The output uses consistent formatting (ATX headers, dash bullets)
// with minimal whitespace.
func Print(node Node) []byte {
	if node == nil {
		return nil
	}
	var buf bytes.Buffer
	_ = PrintTo(&buf, node)

	return buf.Bytes()
}

// PrintTo renders the AST node to the provided io.Writer.
// It streams output without buffering the entire result.
// Returns any write errors encountered.
func PrintTo(w io.Writer, node Node) error {
	if node == nil {
		return nil
	}
	p := &printer{
		w:         w,
		indent:    0,
		listDepth: 0,
		ordered: make(
			[]int,
			0,
			8, //nolint:revive // add-constant
		), // Stack of ordered list counters
	}
	p.printNode(node, false)

	return p.err
}

// printer maintains state during markdown rendering.
type printer struct {
	w          io.Writer
	indent     int   // Current indentation level (in spaces)
	listDepth  int   // Depth of nested lists
	ordered    []int // Stack of ordered list counters (1-based)
	err        error // Accumulated error
	needsBlank bool  // Whether next block needs preceding blank line
}

// write writes bytes to the output, tracking errors.
func (p *printer) write(b []byte) {
	if p.err != nil {
		return
	}
	_, p.err = p.w.Write(b)
}

// writeString writes a string to the output.
func (p *printer) writeString(s string) {
	if p.err != nil {
		return
	}
	_, p.err = io.WriteString(p.w, s)
}

// writeByte writes a single byte to the output.
func (p *printer) writeByte(b byte) {
	if p.err != nil {
		return
	}
	_, p.err = p.w.Write([]byte{b})
}

// writeIndent writes the current indentation.
func (p *printer) writeIndent() {
	if p.indent > 0 {
		p.writeString(
			strings.Repeat(" ", p.indent),
		)
	}
}

// writeBlankLine writes a blank line if needed before a block element.
func (p *printer) writeBlankLine() {
	if p.needsBlank {
		p.writeByte('\n')
	}
	p.needsBlank = true
}

// printNode dispatches to the appropriate print method based on node type.
func (p *printer) printNode(
	node Node,
	isFirst bool,
) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *NodeDocument:
		p.printDocument(n)
	case *NodeSection:
		p.printSection(n, isFirst)
	case *NodeRequirement:
		p.printRequirement(n, isFirst)
	case *NodeScenario:
		p.printScenario(n, isFirst)
	case *NodeParagraph:
		p.printParagraph(n, isFirst)
	case *NodeList:
		p.printList(n, isFirst)
	case *NodeListItem:
		p.printListItem(n)
	case *NodeCodeBlock:
		p.printCodeBlock(n, isFirst)
	case *NodeBlockquote:
		p.printBlockquote(n, isFirst)
	case *NodeText:
		p.printText(n)
	case *NodeStrong:
		p.printStrong(n)
	case *NodeEmphasis:
		p.printEmphasis(n)
	case *NodeStrikethrough:
		p.printStrikethrough(n)
	case *NodeCode:
		p.printCode(n)
	case *NodeLink:
		p.printLink(n)
	case *NodeLinkDef:
		p.printLinkDef(n, isFirst)
	case *NodeWikilink:
		p.printWikilink(n)
	default:
		// For unknown node types, try to print children
		children := node.Children()
		for i, child := range children {
			p.printNode(child, i == 0)
		}
	}
}

// printDocument prints the document root node.
func (p *printer) printDocument(n *NodeDocument) {
	p.needsBlank = false // Don't start with blank line
	children := n.Children()
	for i, child := range children {
		p.printNode(child, i == 0)
	}
	// Ensure single trailing newline (if there's content)
	if len(children) > 0 {
		p.writeByte('\n')
	}
}
