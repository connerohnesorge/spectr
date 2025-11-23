package syntax

import (
	"fmt"
	"strings"
)

// Parser parses the input tokens into an AST.
type Parser struct {
	lexer *Lexer
	curr  Token
	peek  Token // We might need one lookahead, though the lexer is doing a lot of work.
}

// Parse parses the input string and returns the AST.
func Parse(input string) (*Document, error) {
	l := Lex(input)
	p := &Parser{
		lexer: l,
	}
	// Prime the pump
	p.next()
	return p.parseDocument()
}

func (p *Parser) next() {
	p.curr = p.lexer.NextToken()
}

func (p *Parser) parseDocument() (*Document, error) {
	doc := &Document{}
	for p.curr.Type != TokenEOF {
		if p.curr.Type == TokenError {
			return nil, fmt.Errorf("lexing error at line %d: %s", p.curr.Line, p.curr.Value)
		}

		var node Node
		var err error

		switch p.curr.Type {
		case TokenHeader:
			node, err = p.parseHeader()
		case TokenCodeBlock:
			node, err = p.parseCodeBlock()
		case TokenList:
			node, err = p.parseList()
		case TokenText:
			node = p.parseText()
		default:
			// Skip unknown or unexpected tokens for now, or error?
			// For robustness, maybe just treat as text or skip.
			p.next()
			continue
		}

		if err != nil {
			return nil, err
		}
		if node != nil {
			doc.Nodes = append(doc.Nodes, node)
		}
	}
	return doc, nil
}

func (p *Parser) parseHeader() (*Header, error) {
	// Value is like "## Header Text\n"
	line := p.curr.Line
	raw := p.curr.Value
	val := strings.TrimSpace(raw)
	level := 0
	for _, r := range val {
		if r == '#' {
			level++
		} else {
			break
		}
	}
	text := strings.TrimSpace(val[level:])
	p.next()
	return &Header{
		Level: level,
		Text:  text,
		Raw:   raw,
		Line:  line,
	}, nil
}

func (p *Parser) parseCodeBlock() (*CodeBlock, error) {
	// Value is like "```go\ncode\n```"
	line := p.curr.Line
	val := p.curr.Value

	// Strip fences
	lines := strings.Split(val, "\n")
	if len(lines) < 2 {
		// Should be caught by lexer, but just in case
		p.next()
		return &CodeBlock{Content: val, Line: line}, nil
	}

	// First line is ```lang
	lang := strings.TrimPrefix(strings.TrimSpace(lines[0]), "```")

	// Last line is ```
	// Content is everything in between
	content := ""
	if len(lines) > 2 {
		content = strings.Join(lines[1:len(lines)-1], "\n")
	}

	p.next()
	return &CodeBlock{
		Language: lang,
		Content:  content,
		Raw:      val,
		Line:     line,
	}, nil
}

func (p *Parser) parseText() *Text {
	t := &Text{
		Content: p.curr.Value,
		Raw:     p.curr.Value,
		Line:    p.curr.Line,
	}
	p.next()
	return t
}

func (p *Parser) parseList() (*List, error) {
	line := p.curr.Line
	raw := p.curr.Value
	val := strings.TrimSpace(raw)

	var marker string
	var content string
	if strings.HasPrefix(val, "- ") {
		marker = "-"
		content = strings.TrimPrefix(val, "- ")
	} else if strings.HasPrefix(val, "* ") {
		marker = "*"
		content = strings.TrimPrefix(val, "* ")
	}

	p.next()
	return &List{
		Marker:  marker,
		Content: content,
		Raw:     raw,
		Line:    line,
	}, nil
}
