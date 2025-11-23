package parser

import (
	"fmt"
	"strings"
)

// Parser builds an Abstract Syntax Tree (AST) from a stream of tokens.
//
// The parser follows a simple recursive descent approach, converting
// each token type directly into the corresponding AST node. Since the
// lexer has already identified all markdown constructs, the parser's
// job is straightforward: transform tokens into nodes while preserving
// document order and position information.
//
// Design follows Go compiler patterns (cmd/compile/internal/syntax/parser.go)
// with a focus on clarity and error reporting.
type Parser struct {
	tokens []Token // Token stream from lexer
	pos    int     // Current position in token stream
}

// NewParser creates a new parser for the given token stream.
func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// Parse is the main entry point that lexes and parses input in one call.
//
// This convenience function combines lexing and parsing for common use cases.
// For more control, use NewLexer().Lex() followed by ParseTokens().
func Parse(input string) (*Document, error) {
	lexer := NewLexer(input)
	tokens := lexer.Lex()

	return ParseTokens(tokens)
}

// ParseTokens parses a token stream into an AST.
//
// This is the core parsing function. It converts a slice of tokens
// (produced by the lexer) into a Document node containing the full
// AST representation of the markdown.
//
// Returns an error if:
// - An error token is encountered
// - Unexpected token types are found
// - The token stream is malformed
func ParseTokens(tokens []Token) (*Document, error) {
	if len(tokens) == 0 {
		return NewDocument(), nil
	}

	p := NewParser(tokens)

	return p.parseDocument()
}

// parseDocument is the top-level parsing function.
//
// It processes the token stream sequentially, converting each token
// into the appropriate AST node type and adding it to the document's
// children.
func (p *Parser) parseDocument() (*Document, error) {
	doc := NewDocument()

	for !p.atEnd() {
		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		// Only add non-nil nodes (we might skip EOF)
		if node != nil {
			doc.Children = append(doc.Children, node)
		}
	}

	return doc, nil
}

// parseNode parses the next token into an AST node.
//
// This is the main dispatch function that routes each token type
// to its corresponding node constructor.
func (p *Parser) parseNode() (Node, error) {
	if p.atEnd() {
		return nil, nil
	}

	token := p.current()

	switch token.Type {
	case TokenEOF:
		// EOF token signals end - advance past it but don't create a node
		p.advance()

		return nil, nil

	case TokenHeader:
		return p.parseHeader()

	case TokenText:
		return p.parseParagraph()

	case TokenCodeBlock:
		return p.parseCodeBlock()

	case TokenListItem:
		return p.parseListItem()

	case TokenBlankLine:
		return p.parseBlankLine()

	case TokenError:
		// Error tokens indicate lexer failures
		return nil, fmt.Errorf("lexer error at %s: %s", token.Pos, token.Value)

	default:
		return nil, fmt.Errorf(
			"unexpected token type %s at %s",
			token.Type, token.Pos,
		)
	}
}

// parseHeader converts a TokenHeader into a Header node.
//
// The header token value contains the full header line including
// the # markers. We need to:
// 1. Count the # characters to determine level
// 2. Extract the text after the # markers
func (p *Parser) parseHeader() (*Header, error) {
	token := p.consume()

	// Count leading # characters
	level := 0
	text := token.Value

	for len(text) > 0 && text[0] == '#' {
		level++
		text = text[1:]
	}

	// Trim whitespace from the remaining text
	text = strings.TrimSpace(text)

	return NewHeader(level, text, token.Pos), nil
}

// parseParagraph converts a TokenText into a Paragraph node.
//
// Text tokens represent regular content that isn't part of any
// special markdown construct.
func (p *Parser) parseParagraph() (*Paragraph, error) {
	token := p.consume()

	return NewParagraph(token.Value, token.Pos), nil
}

// parseCodeBlock converts a TokenCodeBlock into a CodeBlock node.
//
// The code block token value contains the full code fence including
// opening fence, language specifier, content, and closing fence.
// We need to:
// 1. Extract the language specifier (first line after ```)
// 2. Extract the content (everything between fences)
func (p *Parser) parseCodeBlock() (*CodeBlock, error) {
	token := p.consume()

	// Split the token value into lines
	lines := strings.Split(token.Value, "\n")
	if len(lines) == 0 {
		return NewCodeBlock("", "", token.Pos), nil
	}

	// First line contains opening fence and optional language
	firstLine := lines[0]

	// Extract language by removing leading backticks
	language := strings.TrimLeft(firstLine, "`")
	language = strings.TrimSpace(language)

	// Find the content between fences
	// Start from line 1 (skip opening fence)
	// End before the last line if it's a closing fence
	contentStart := 1
	contentEnd := len(lines)

	// Check if last line is a closing fence
	if contentEnd > 1 {
		lastLine := strings.TrimSpace(lines[contentEnd-1])
		if strings.HasPrefix(lastLine, "```") {
			contentEnd--
		}
	}

	// Extract content lines
	var contentLines []string
	if contentEnd > contentStart {
		contentLines = lines[contentStart:contentEnd]
	}

	content := strings.Join(contentLines, "\n")

	return NewCodeBlock(language, content, token.Pos), nil
}

// parseListItem converts a TokenListItem into a List node.
//
// Each token represents a single list item. The token value contains
// the list marker (- or *) followed by the item text.
func (p *Parser) parseListItem() (*List, error) {
	token := p.consume()

	// Remove the list marker (first character) and leading whitespace
	text := token.Value
	if len(text) > 0 && (text[0] == '-' || text[0] == '*') {
		text = text[1:]
	}
	text = strings.TrimSpace(text)

	return NewList(text, token.Pos), nil
}

// parseBlankLine converts a TokenBlankLine into a BlankLine node.
//
// The token value contains the blank line(s). We count the newlines
// to determine how many consecutive blank lines there are.
func (p *Parser) parseBlankLine() (*BlankLine, error) {
	token := p.consume()

	// Count newlines in the token value
	count := strings.Count(token.Value, "\n")
	if count == 0 {
		count = 1 // At minimum, one blank line
	}

	return NewBlankLine(count, token.Pos), nil
}

// Helper methods for navigating the token stream

// current returns the token at the current position without advancing.
//
// Returns a sentinel EOF token if at end of stream.
func (p *Parser) current() Token {
	if p.atEnd() {
		// Return a sentinel EOF token
		return Token{
			Type: TokenEOF,
			Pos:  Position{Line: 0, Column: 0, Offset: 0},
		}
	}

	return p.tokens[p.pos]
}

// peek returns the next token without advancing position.
//
// Returns a sentinel EOF token if at end of stream.
func (p *Parser) peek() Token {
	if p.pos+1 >= len(p.tokens) {
		return Token{
			Type: TokenEOF,
			Pos:  Position{Line: 0, Column: 0, Offset: 0},
		}
	}

	return p.tokens[p.pos+1]
}

// advance moves to the next token.
func (p *Parser) advance() {
	if !p.atEnd() {
		p.pos++
	}
}

// consume returns the current token and advances to the next.
func (p *Parser) consume() Token {
	token := p.current()
	p.advance()

	return token
}

// expect checks if the current token matches the expected type.
//
// If it matches, the token is consumed and returned.
// If it doesn't match, an error is returned.
//
// This is useful for enforcing expected token sequences.
func (p *Parser) expect(expected TokenType) (Token, error) {
	token := p.current()
	if token.Type != expected {
		return token, fmt.Errorf(
			"expected %s but got %s at %s",
			expected,
			token.Type,
			token.Pos,
		)
	}
	p.advance()

	return token, nil
}

// atEnd returns true if the parser has consumed all tokens.
func (p *Parser) atEnd() bool {
	return p.pos >= len(p.tokens)
}
