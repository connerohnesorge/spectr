package mdparser

import (
	"fmt"
	"strings"
)

// Parser builds an AST from a token stream.
//
// The parser uses a two-token lookahead (current and peek) to make
// parsing decisions. This allows it to handle ambiguous cases and
// provide better error messages.
type Parser struct {
	lexer   *Lexer
	current Token
	peek    Token
}

// Parse is the main entry point for parsing markdown content.
//
// It creates a lexer and parser, then builds a complete Document AST
// from the input string. The parser makes a single pass through the
// token stream.
//
// Parameters:
//   - input: Raw markdown content as a string
//
// Returns:
//   - *Document: Root AST node containing all parsed elements
//   - error: Parse error if the input is malformed
//
// Example:
//
//	doc, err := mdparser.Parse("## Header\n\nParagraph text")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Parsed %d nodes\n", len(doc.Children))
func Parse(input string) (*Document, error) {
	lexer := NewLexer(input)
	parser := &Parser{lexer: lexer}

	// Initialize current and peek tokens
	parser.advance()
	parser.advance()

	return parser.parseDocument()
}

// advance moves to the next token.
func (p *Parser) advance() {
	p.current = p.peek
	p.peek = p.lexer.NextToken()
}

// parseDocument parses the entire document.
func (p *Parser) parseDocument() (*Document, error) {
	doc := &Document{
		StartPos: Position{Line: 1, Column: 1, Offset: 0},
		Children: nil,
	}

	for p.current.Type != TokenEOF {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}
		if node != nil {
			doc.Children = append(doc.Children, node)
		}
	}

	doc.EndPos = p.current.Pos

	return doc, nil
}

// parseNode parses a single node based on the current token.
func (p *Parser) parseNode() (Node, error) {
	switch p.current.Type {
	case TokenHeader:
		return p.parseHeader()
	case TokenCodeFence:
		return p.parseCodeBlock()
	case TokenListItem:
		return p.parseList()
	case TokenText:
		return p.parseParagraph()
	case TokenBlankLine:
		return p.parseBlankLine()
	case TokenError:
		return nil, fmt.Errorf(
			"lexer error at line %d: %s",
			p.current.Pos.Line,
			p.current.Value,
		)
	case TokenEOF, TokenCodeContent:
		// Skip EOF and code content tokens at top level
		p.advance()

		return nil, nil
	default:
		// Skip unknown tokens
		p.advance()

		return nil, nil
	}
}

// parseHeader parses a header token into a Header node.
func (p *Parser) parseHeader() (*Header, error) {
	token := p.current
	p.advance()

	// Extract level and text from token value
	level := 0
	text := token.Value
	for i := 0; i < len(text) && text[i] == '#'; i++ {
		level++
	}

	// Remove # symbols and trim whitespace
	text = strings.TrimSpace(text[level:])

	return &Header{
		StartPos: token.Pos,
		EndPos:   token.Pos,
		Level:    level,
		Text:     text,
	}, nil
}

// parseParagraph parses consecutive text tokens into a Paragraph node.
func (p *Parser) parseParagraph() (*Paragraph, error) {
	startPos := p.current.Pos
	var lines []string

	// Collect consecutive text lines
	for p.current.Type == TokenText {
		lines = append(lines, p.current.Value)
		p.advance()
	}

	endPos := p.current.Pos

	return &Paragraph{
		StartPos: startPos,
		EndPos:   endPos,
		Lines:    lines,
	}, nil
}

// parseCodeBlock parses a code block (from opening fence to closing fence).
func (p *Parser) parseCodeBlock() (*CodeBlock, error) {
	startPos := p.current.Pos
	openingFence := p.current.Value
	p.advance()

	// Extract language from opening fence
	language := strings.TrimSpace(strings.TrimPrefix(openingFence, "```"))

	var lines []string

	// Collect code content until closing fence
	for p.current.Type != TokenEOF {
		if p.current.Type == TokenCodeFence {
			// Found closing fence
			endPos := p.current.Pos
			p.advance()

			return &CodeBlock{
				StartPos: startPos,
				EndPos:   endPos,
				Language: language,
				Lines:    lines,
			}, nil
		}

		if p.current.Type == TokenCodeContent {
			lines = append(lines, p.current.Value)
		}
		p.advance()
	}

	// Missing closing fence - return what we have
	return &CodeBlock{
		StartPos: startPos,
		EndPos:   p.current.Pos,
		Language: language,
		Lines:    lines,
	}, nil
}

// parseList parses consecutive list items into a List node.
func (p *Parser) parseList() (*List, error) {
	startPos := p.current.Pos

	// Determine if ordered or unordered based on first item
	firstItem := p.current.Value
	ordered := len(firstItem) > 0 &&
		firstItem[0] >= '0' &&
		firstItem[0] <= '9'

	var items []*ListItem

	// Collect consecutive list items
	for p.current.Type == TokenListItem {
		item, err := p.parseListItem()
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	endPos := p.current.Pos

	return &List{
		StartPos: startPos,
		EndPos:   endPos,
		Ordered:  ordered,
		Items:    items,
	}, nil
}

// parseListItem parses a single list item.
func (p *Parser) parseListItem() (*ListItem, error) {
	startPos := p.current.Pos
	text := p.current.Value

	// Remove bullet/number prefix
	text = strings.TrimSpace(text)
	if strings.HasPrefix(text, "- ") || strings.HasPrefix(text, "* ") {
		text = text[2:]
	} else if len(text) > 0 && text[0] >= '0' && text[0] <= '9' {
		// Remove number and period
		dotIdx := strings.Index(text, ". ")
		if dotIdx >= 0 {
			text = text[dotIdx+2:]
		}
	}

	p.advance()

	// TODO: Handle nested content (paragraphs, sublists, etc.)
	// For now, just store the text

	return &ListItem{
		StartPos: startPos,
		EndPos:   p.current.Pos,
		Text:     text,
		Children: nil,
	}, nil
}

// parseBlankLine parses a blank line token into a BlankLine node.
func (p *Parser) parseBlankLine() (*BlankLine, error) {
	startPos := p.current.Pos

	// Count consecutive blank lines
	count := 0
	for p.current.Type == TokenBlankLine {
		count++
		p.advance()
	}

	endPos := p.current.Pos

	return &BlankLine{
		StartPos: startPos,
		EndPos:   endPos,
		Count:    count,
	}, nil
}
