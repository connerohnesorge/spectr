package markdown

import (
	"bytes"
	"strings"
)

// printSection prints a section header with ATX style.
//
//nolint:revive // flag-parameter
func (p *printer) printSection(
	n *NodeSection,
	isFirst bool,
) {
	if !isFirst {
		p.writeBlankLine()
	}

	// Write hash prefix based on level
	level := n.Level()
	if level < 1 {
		level = 1
	}
	if level > 6 { //nolint:revive // add-constant
		level = 6 //nolint:revive // add-constant
	}
	p.writeString(strings.Repeat("#", level))
	p.writeByte(' ')

	// For delta sections, format as "DELTA_TYPE Requirements"
	if deltaType := n.DeltaType(); deltaType != "" {
		p.writeString(deltaType)
		p.writeString(" Requirements")
	} else {
		// Write title
		p.write(n.Title())
	}
	p.writeByte('\n')

	// Print children (content under the header)
	children := n.Children()
	for i, child := range children {
		p.printNode(child, i == 0)
	}
}

// printRequirement prints a requirement header (### Requirement: Name).
//
//nolint:revive // flag-parameter
func (p *printer) printRequirement(
	n *NodeRequirement,
	isFirst bool,
) {
	if !isFirst {
		p.writeBlankLine()
	}

	p.writeString("### Requirement: ")
	p.writeString(n.Name())
	p.writeByte('\n')

	// Print children (scenarios, paragraphs, etc.)
	children := n.Children()
	for i, child := range children {
		p.printNode(child, i == 0)
	}
}

// printScenario prints a scenario header (#### Scenario: Name).
//
//nolint:revive // flag-parameter
func (p *printer) printScenario(
	n *NodeScenario,
	isFirst bool,
) {
	if !isFirst {
		p.writeBlankLine()
	}

	p.writeString("#### Scenario: ")
	p.writeString(n.Name())
	p.writeByte('\n')

	// Print children
	children := n.Children()
	for i, child := range children {
		p.printNode(child, i == 0)
	}
}

// printList prints an unordered or ordered list.
//
//nolint:revive // flag-parameter
func (p *printer) printList(
	n *NodeList,
	isFirst bool,
) {
	if !isFirst && p.listDepth == 0 {
		p.writeBlankLine()
	}

	p.listDepth++
	if n.Ordered() {
		// Push a new counter for this ordered list
		p.ordered = append(p.ordered, 1)
	}

	children := n.Children()
	for _, child := range children {
		if item, ok := child.(*NodeListItem); ok {
			p.printListItem(item)
			if n.Ordered() {
				// Increment counter
				p.ordered[len(p.ordered)-1]++
			}
		} else {
			p.printNode(child, false)
		}
	}

	if n.Ordered() {
		// Pop the counter
		p.ordered = p.ordered[:len(p.ordered)-1]
	}
	p.listDepth--
}

// printListItem prints a single list item.
//
//nolint:revive // function-length - list item formatting requires multiple steps
func (p *printer) printListItem(n *NodeListItem) {
	p.writeIndent()

	// Determine bullet style
	if len(p.ordered) > 0 &&
		p.ordered[len(p.ordered)-1] > 0 {
		// Ordered list
		num := p.ordered[len(p.ordered)-1]
		p.writeString(itoa(num))
		p.writeString(". ")
	} else {
		// Unordered list
		p.writeString("- ")
	}

	// Handle checkbox
	checked, hasCheckbox := n.Checked()
	if hasCheckbox {
		if checked {
			p.writeString("[x] ")
		} else {
			p.writeString("[ ] ")
		}
	}

	// Handle WHEN/THEN/AND keywords
	keyword := n.Keyword()
	if keyword != "" {
		p.writeString("**")
		p.writeString(strings.ToUpper(keyword))
		p.writeString("** ")
	}

	// Print children inline
	children := n.Children()
	hasNestedList := false
	for i, child := range children {
		if _, isList := child.(*NodeList); isList {
			hasNestedList = true
			if i > 0 {
				// Add newline before nested list
				p.writeByte('\n')
			}
			// Increase indent for nested list
			oldIndent := p.indent
			p.indent += 2
			p.printList(child.(*NodeList), true)
			p.indent = oldIndent
		} else {
			p.printInline(child)
		}
	}

	if !hasNestedList || len(children) == 0 {
		p.writeByte('\n')
	}
}

// printCodeBlock prints a fenced code block.
//
//nolint:revive // flag-parameter
func (p *printer) printCodeBlock(
	n *NodeCodeBlock,
	isFirst bool,
) {
	if !isFirst {
		p.writeBlankLine()
	}

	p.writeIndent()
	p.writeString("```")
	if lang := n.Language(); len(lang) > 0 {
		p.write(lang)
	}
	p.writeByte('\n')

	// Print content verbatim, preserving internal newlines
	content := n.Content()
	if len(content) > 0 {
		// Add indentation to each line of content
		if p.indent > 0 {
			lines := bytes.Split(
				content,
				[]byte{'\n'},
			)
			for i, line := range lines {
				if i > 0 {
					p.writeByte('\n')
				}
				p.writeIndent()
				p.write(line)
			}
		} else {
			p.write(content)
		}
		// Ensure content ends with newline
		if content[len(content)-1] != '\n' {
			p.writeByte('\n')
		}
	}

	p.writeIndent()
	p.writeString("```\n")
}

// printBlockquoteChild prints a child of a blockquote with > prefix.
//
//nolint:revive // function-length - blockquote child formatting handles multiple node types
func (p *printer) printBlockquoteChild(
	node Node,
) {
	if node == nil {
		return
	}

	switch n := node.(type) {
	case *NodeParagraph:
		p.writeIndent()
		p.writeString("> ")
		children := n.Children()
		for _, child := range children {
			p.printInline(child)
		}
		p.writeByte('\n')
	case *NodeBlockquote:
		// Nested blockquote - print with additional >
		children := n.Children()
		for _, child := range children {
			p.writeIndent()
			p.writeString("> ")
			p.printBlockquoteChild(child)
		}
	case *NodeList:
		// List in blockquote
		children := n.Children()
		for i, child := range children {
			item, ok := child.(*NodeListItem)
			if !ok {
				continue
			}
			p.writeIndent()
			p.writeString("> ")
			if n.Ordered() {
				p.writeString(itoa(i + 1))
				p.writeString(". ")
			} else {
				p.writeString("- ")
			}
			itemChildren := item.Children()
			for _, ic := range itemChildren {
				p.printInline(ic)
			}
			p.writeByte('\n')
		}
	default:
		// For other types, prefix each line with >
		p.writeIndent()
		p.writeString("> ")
		source := node.Source()
		if source != nil {
			// Handle multi-line source
			lines := bytes.Split(source, []byte{'\n'})
			for i, line := range lines {
				if i > 0 {
					p.writeByte('\n')
					p.writeIndent()
					p.writeString("> ")
				}
				p.write(line)
			}
		}
		p.writeByte('\n')
	}
}
