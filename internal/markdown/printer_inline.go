package markdown

import "bytes"

// printParagraph prints a paragraph.
//
//nolint:revive // flag-parameter
func (p *printer) printParagraph(
	n *NodeParagraph,
	isFirst bool,
) {
	if !isFirst {
		p.writeBlankLine()
	}

	p.writeIndent()
	// Print inline children
	children := n.Children()
	for _, child := range children {
		p.printInline(child)
	}
	p.writeByte('\n')
}

// printInline prints an inline node without block-level formatting.
func (p *printer) printInline(node Node) {
	if node == nil {
		return
	}

	switch n := node.(type) {
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
	case *NodeWikilink:
		p.printWikilink(n)
	default:
		// Fallback: write source if available
		source := node.Source()
		if source != nil {
			p.write(source)
		}
	}
}

// printText prints plain text content.
func (p *printer) printText(n *NodeText) {
	p.write(n.Source())
}

// printStrong prints bold emphasis with **.
func (p *printer) printStrong(n *NodeStrong) {
	p.writeString("**")
	children := n.Children()
	if len(children) > 0 {
		for _, child := range children {
			p.printInline(child)
		}
	} else {
		// If no children, use source content (stripping delimiters if present)
		source := n.Source()
		if len(source) > 4 && source[0] == '*' && source[1] == '*' { //nolint:revive // add-constant
			// Strip ** delimiters
			p.write(source[2 : len(source)-2])
		} else {
			p.write(source)
		}
	}
	p.writeString("**")
}

// printEmphasis prints italic emphasis with *.
func (p *printer) printEmphasis(n *NodeEmphasis) {
	p.writeByte('*')
	children := n.Children()
	if len(children) > 0 {
		for _, child := range children {
			p.printInline(child)
		}
	} else {
		// If no children, use source content (stripping delimiters if present)
		source := n.Source()
		if len(source) > 2 && (source[0] == '*' || source[0] == '_') {
			// Strip delimiters
			p.write(source[1 : len(source)-1])
		} else {
			p.write(source)
		}
	}
	p.writeByte('*')
}

// printStrikethrough prints struck text with ~~.
func (p *printer) printStrikethrough(
	n *NodeStrikethrough,
) {
	p.writeString("~~")
	children := n.Children()
	if len(children) > 0 {
		for _, child := range children {
			p.printInline(child)
		}
	} else {
		// If no children, use source content (stripping delimiters if present)
		source := n.Source()
		if len(source) > 4 && source[0] == '~' && source[1] == '~' { //nolint:revive // add-constant
			// Strip ~~ delimiters
			p.write(source[2 : len(source)-2])
		} else {
			p.write(source)
		}
	}
	p.writeString("~~")
}

// printCode prints inline code with backticks.
func (p *printer) printCode(n *NodeCode) {
	source := n.Source()

	// Check if content contains backticks
	hasBacktick := bytes.Contains(
		source,
		[]byte{'`'},
	)

	if hasBacktick {
		// Use double backticks with space padding
		p.writeString("`` ")
		p.write(source)
		p.writeString(" ``")
	} else {
		p.writeByte('`')
		p.write(source)
		p.writeByte('`')
	}
}

// printLink prints an inline link [text](url).
func (p *printer) printLink(n *NodeLink) {
	p.writeByte('[')

	// Print link text from children
	children := n.Children()
	for _, child := range children {
		p.printInline(child)
	}

	p.writeString("](")
	p.write(n.URL())

	// Include title if present
	if title := n.Title(); len(title) > 0 {
		p.writeString(" \"")
		p.write(title)
		p.writeByte('"')
	}

	p.writeByte(')')
}

// printLinkDef prints a link definition [ref]: url "title".
//
//nolint:revive // flag-parameter
func (p *printer) printLinkDef(
	n *NodeLinkDef,
	isFirst bool,
) {
	if !isFirst {
		p.writeBlankLine()
	}

	p.writeIndent()
	p.writeByte('[')

	// Extract reference label from source or children
	children := n.Children()
	if len(children) > 0 {
		for _, child := range children {
			p.printInline(child)
		}
	} else {
		// Use first part of source as label
		source := n.Source()
		// Find the ]: to extract label
		idx := bytes.Index(source, []byte("]:"))
		if idx > 1 && source[0] == '[' {
			p.write(source[1:idx])
		}
	}

	p.writeString("]: ")
	p.write(n.URL())

	if title := n.Title(); len(title) > 0 {
		p.writeString(" \"")
		p.write(title)
		p.writeByte('"')
	}

	p.writeByte('\n')
}

// printWikilink prints a wikilink [[target|display#anchor]].
func (p *printer) printWikilink(n *NodeWikilink) {
	p.writeString("[[")
	p.write(n.Target())

	// Add anchor if present
	if anchor := n.Anchor(); len(anchor) > 0 {
		p.writeByte('#')
		p.write(anchor)
	}

	// Add display text if different from target
	if display := n.Display(); len(display) > 0 {
		p.writeByte('|')
		p.write(display)
	}

	p.writeString("]]")
}

// printBlockquote prints a blockquote.
//
//nolint:revive // flag-parameter
func (p *printer) printBlockquote(
	n *NodeBlockquote,
	isFirst bool,
) {
	if !isFirst {
		p.writeBlankLine()
	}

	children := n.Children()
	for _, child := range children {
		p.printBlockquoteChild(child)
	}
}
