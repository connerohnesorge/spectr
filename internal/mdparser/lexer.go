package mdparser

import (
	"strings"
	"unicode"
)

// TokenType represents the type of a token.
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenHeader
	TokenText
	TokenCodeFence
	TokenCodeContent
	TokenListItem
	TokenBlankLine
	TokenError
)

const maxHeaderLevel = 6

func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenHeader:
		return "Header"
	case TokenText:
		return "Text"
	case TokenCodeFence:
		return "CodeFence"
	case TokenCodeContent:
		return "CodeContent"
	case TokenListItem:
		return "ListItem"
	case TokenBlankLine:
		return "BlankLine"
	case TokenError:
		return "Error"
	default:
		return "Unknown"
	}
}

// Token represents a lexical token.
type Token struct {
	Type  TokenType
	Value string
	Pos   Position
}

// Lexer performs lexical analysis of markdown input using a state machine.
//
// The lexer tokenizes markdown line-by-line, maintaining state to handle
// context-dependent structures like code blocks. It tracks position information
// (line, column, offset) for each token to enable accurate error reporting.
//
// State Management:
//   - inCodeBlock: Tracks whether we're inside a fenced code block
//   - State functions return the next state function to execute
//
// The lexer buffers tokens internally and returns them via NextToken().
type Lexer struct {
	input       string
	start       int // start position of current token
	pos         int // current position in input
	line        int // current line (1-based)
	col         int // current column (1-based)
	tokens      []Token
	inCodeBlock bool // state: are we inside a code block?
}

// stateFn represents a state in the lexer state machine.
// Each state function processes input and returns the next state to execute.
// Returning nil signals the end of lexing.
type stateFn func(*Lexer) stateFn

// NewLexer creates a new lexer for the given input.
//
// The lexer is initialized with line and column positions starting at 1.
// Call NextToken() repeatedly to get tokens from the input.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input: input,
		line:  1,
		col:   1,
	}
}

// NextToken returns the next token from the input.
//
// The lexer runs its state machine lazily, generating tokens on demand.
// When the input is exhausted, it returns TokenEOF indefinitely.
//
// This method is safe to call repeatedly - once EOF is reached, all
// subsequent calls will return EOF tokens.
func (l *Lexer) NextToken() Token {
	if len(l.tokens) == 0 {
		// Run the state machine to generate tokens
		for state := lexStart; state != nil; {
			state = state(l)
			if len(l.tokens) == 0 {
				continue
			}
			token := l.tokens[0]
			l.tokens = l.tokens[1:]

			return token
		}

		// Return EOF if no more tokens
		return Token{Type: TokenEOF, Pos: l.currentPos()}
	}

	token := l.tokens[0]
	l.tokens = l.tokens[1:]

	return token
}

// Tokenize returns all tokens from the input.
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF || tok.Type == TokenError {
			break
		}
	}

	return tokens
}

// peek returns the next character without consuming it.
func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}

	return rune(l.input[l.pos])
}

// next consumes and returns the next character.
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r := rune(l.input[l.pos])
	l.pos++
	if r == '\n' {
		l.line++
		l.col = 1
	} else {
		l.col++
	}

	return r
}

// emit creates a token of the given type and adds it to the token stream.
func (l *Lexer) emit(t TokenType) {
	value := l.input[l.start:l.pos]

	// Calculate the starting column by counting back from current position
	startCol := l.col
	for i := l.pos - 1; i >= l.start && i >= 0; i-- {
		if l.input[i] == '\n' {
			break
		}
		startCol--
	}
	if startCol < 1 {
		startCol = 1
	}

	token := Token{
		Type:  t,
		Value: value,
		Pos:   Position{Line: l.line, Column: startCol, Offset: l.start},
	}
	l.tokens = append(l.tokens, token)
	l.start = l.pos
}

// ignore skips over the pending input.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// currentPos returns the current position.
func (l *Lexer) currentPos() Position {
	return Position{Line: l.line, Column: l.col, Offset: l.pos}
}

// atStartOfLine checks if we're at the start of a line.
func (l *Lexer) atStartOfLine() bool {
	if l.pos == 0 {
		return true
	}
	if l.pos > 0 && l.input[l.pos-1] == '\n' {
		return true
	}

	return false
}

// peekLine returns the rest of the current line without consuming it.
func (l *Lexer) peekLine() string {
	end := l.pos
	for end < len(l.input) && l.input[end] != '\n' {
		end++
	}

	return l.input[l.pos:end]
}

// skipToEndOfLine consumes characters until the end of the line.
func (l *Lexer) skipToEndOfLine() {
	for l.peek() != '\n' && l.peek() != 0 {
		l.next()
	}
}

// checkOrderedListItem checks if the line starts with an ordered list item.
func (*Lexer) checkOrderedListItem(line string) bool {
	if len(line) == 0 || !unicode.IsDigit(rune(line[0])) {
		return false
	}

	dotIdx := strings.Index(line, ". ")
	if dotIdx <= 0 {
		return false
	}

	for _, r := range line[:dotIdx] {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
