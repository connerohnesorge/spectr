// Package parser provides a robust lexer/parser architecture for parsing
// markdown specifications. It is modeled after the Go compiler's scanner
// and parser (cmd/compile/internal/syntax) combined with Rob Pike's state
// function approach from "Lexical Scanning in Go".
//
// The parser separates concerns into three layers:
//   - Lexer: Tokenizes input into a stream of tokens with position info
//   - Parser: Builds an Abstract Syntax Tree (AST) from tokens
//   - Extractor: Extracts Spectr-specific structures from the AST
//
// This architecture eliminates the brittleness of regex-based parsing
// and provides precise error reporting with line and column information.
//
//nolint:revive // file-length-limit - lexer infrastructure requires detail
package parser

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	// tokenStringMaxLen is max chars to show in Token.String()
	tokenStringMaxLen = 20
	// initialTokenCapacity is the initial capacity for token slice
	initialTokenCapacity = 64
)

// TokenType represents the type of a lexical token.
type TokenType int

const (
	// TokenEOF signals the end of input
	TokenEOF TokenType = iota

	// TokenText represents plain text content
	TokenText

	// TokenHeader represents a markdown header (# ## ### ####)
	TokenHeader

	// TokenCodeBlock represents a code fence (```) and its content
	TokenCodeBlock

	// TokenListItem represents a list item (- or *)
	TokenListItem

	// TokenBlankLine represents one or more blank lines
	TokenBlankLine

	// TokenError represents a lexical error
	TokenError
)

// String returns the string representation of a TokenType.
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenText:
		return "Text"
	case TokenHeader:
		return "Header"
	case TokenCodeBlock:
		return "CodeBlock"
	case TokenListItem:
		return "ListItem"
	case TokenBlankLine:
		return "BlankLine"
	case TokenError:
		return "Error"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}

// Position represents a position in the source text.
type Position struct {
	Line   int // Line number (1-based)
	Column int // Column number (1-based)
	Offset int // Byte offset (0-based)
}

// String returns a human-readable representation of a position.
func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Token represents a lexical token with its type, value, and position.
type Token struct {
	Type  TokenType // Type of token
	Value string    // Actual text content
	Pos   Position  // Position where token starts
}

// String returns a human-readable representation of a token.
func (t Token) String() string {
	if len(t.Value) > tokenStringMaxLen {
		return fmt.Sprintf(
			"%s@%s: %.20s...",
			t.Type,
			t.Pos,
			t.Value,
		)
	}

	return fmt.Sprintf("%s@%s: %s", t.Type, t.Pos, t.Value)
}

// StateFn represents a state in the lexer's state machine.
// It processes input and returns the next state function, or nil when done.
type StateFn func(*Lexer) StateFn

// Lexer tokenizes markdown input into a stream of tokens.
//
// The lexer uses a state machine approach where each state function
// processes a portion of input and returns the next state. This design
// is inspired by Rob Pike's "Lexical Scanning in Go" talk, adapted to
// fit the structural style of the Go compiler's scanner.
type Lexer struct {
	input  string  // The input string being scanned
	start  int     // Start position of current token
	pos    int     // Current position in input
	width  int     // Width of last rune read
	line   int     // Current line number (1-based)
	column int     // Current column number (1-based)
	tokens []Token // Emitted tokens
	state  StateFn // Current state function
}

// NewLexer creates a new lexer for the given input.
//
// The lexer starts at position 1:1 (line 1, column 1) following
// traditional text editor conventions.
func NewLexer(input string) *Lexer {
	return &Lexer{
		input:  input,
		start:  0,
		pos:    0,
		width:  0,
		line:   1,
		column: 1,
		tokens: make([]Token, 0, initialTokenCapacity),
		state:  nil, // Will be set to initial state function
	}
}

// next returns the next rune in the input and advances the position.
//
// Returns 0 (nul) if at end of input.
func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		l.width = 0

		return 0
	}

	r, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += width

	// Update line and column tracking
	if r == '\n' {
		l.line++
		l.column = 1
	} else {
		l.column++
	}

	return r
}

// peek returns the next rune without advancing the position.
//
// Returns 0 (nul) if at end of input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

// backup steps back one rune in the input.
//
// Can only be called once per call to next.
//
//nolint:revive // max-control-nesting - position tracking needs complexity
func (l *Lexer) backup() {
	if l.width == 0 {
		return
	}

	l.pos -= l.width

	// We need to check what rune we just backed over
	if l.pos >= 0 && l.pos < len(l.input) {
		r, _ := utf8.DecodeRuneInString(l.input[l.pos:])
		if r == '\n' {
			l.backupOverNewline()
		} else {
			l.column--
		}
	} else if l.pos == 0 {
		l.line = 1
		l.column = 1
	}

	l.width = 0
}

// backupOverNewline handles backing up over a newline character.
func (l *Lexer) backupOverNewline() {
	// We're backing up over a newline, need to restore previous line
	l.line--
	// Scan backwards to find start of current (now previous) line
	prevLineStart := l.pos - 1
	for prevLineStart >= 0 {
		if l.input[prevLineStart] == '\n' {
			prevLineStart++

			break
		}
		prevLineStart--
	}
	if prevLineStart < 0 {
		prevLineStart = 0
	}
	l.column = l.pos - prevLineStart + 1
}

// emit creates a token from the current input slice and adds it
// to the token stream.
func (l *Lexer) emit(t TokenType) {
	token := Token{
		Type:  t,
		Value: l.input[l.start:l.pos],
		Pos: Position{
			Line:   l.lineAtStart(),
			Column: l.columnAtStart(),
			Offset: l.start,
		},
	}
	l.tokens = append(l.tokens, token)
	l.start = l.pos
}

// emitError creates an error token with the given message.
func (l *Lexer) emitError(message string) {
	token := Token{
		Type:  TokenError,
		Value: message,
		Pos: Position{
			Line:   l.line,
			Column: l.column,
			Offset: l.pos,
		},
	}
	l.tokens = append(l.tokens, token)
}

// ignore skips over the current input without emitting a token.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
//
// Returns true if a rune was consumed.
func (l *Lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()

	return false
}

// acceptRun consumes a run of runes from the valid set.
//
// Returns the number of runes consumed.
func (l *Lexer) acceptRun(valid string) int {
	count := 0
	for strings.ContainsRune(valid, l.next()) {
		count++
	}
	l.backup()

	return count
}

// acceptUntil consumes runes until encountering one from the stop set.
//
// Returns the number of runes consumed.
//
//nolint:unused // Will be used by state functions in task 1.2
func (l *Lexer) acceptUntil(stop string) int {
	count := 0
	for {
		r := l.next()
		if r == 0 || strings.ContainsRune(stop, r) {
			l.backup()

			break
		}
		count++
	}

	return count
}

// atLineStart returns true if the lexer is at the start of a line.
func (l *Lexer) atLineStart() bool {
	if l.pos == 0 {
		return true
	}
	if l.pos >= len(l.input) {
		return false
	}
	// Check if previous character was a newline
	if l.pos > 0 {
		r, _ := utf8.DecodeLastRuneInString(l.input[:l.pos])

		return r == '\n'
	}

	return false
}

// lineAtStart returns the line number where the current token started.
//
// This is calculated by counting newlines between start and position.
func (l *Lexer) lineAtStart() int {
	line := 1
	for i := range l.start {
		if l.input[i] == '\n' {
			line++
		}
	}

	return line
}

// columnAtStart returns the column number where the current token started.
//
// This is calculated by finding the last newline before start position.
func (l *Lexer) columnAtStart() int {
	column := 1
	for i := l.start - 1; i >= 0; i-- {
		if l.input[i] == '\n' {
			break
		}
		column++
	}

	return column
}

// currentValue returns the current token value being built.
func (l *Lexer) currentValue() string {
	return l.input[l.start:l.pos]
}

// Tokens returns all tokens emitted by the lexer.
//
// This method is typically called after running the lexer to completion.
func (l *Lexer) Tokens() []Token {
	return l.tokens
}

// run executes the state machine until it reaches a terminal state (nil).
//
// The initial state function must be set before calling run.
func (l *Lexer) run(initialState StateFn) {
	l.state = initialState
	for l.state != nil {
		l.state = l.state(l)
	}
}

// Lex tokenizes the input and returns all tokens.
//
// This is the main entry point for lexing. It runs the state machine
// starting from lexText and returns the complete token stream.
func (l *Lexer) Lex() []Token {
	l.run(lexText)

	return l.tokens
}

// lexText is the default state that consumes plain text.
//
// It detects the start of special markdown constructs:
// - Headers (# at line start)
// - Code blocks (``` at line start)
// - List items (- or * at line start)
// - Blank lines (empty lines)
//
// When it encounters these constructs, it emits any accumulated text
// and transitions to the appropriate state.
func lexText(l *Lexer) StateFn {
	for {
		// Check for EOF
		if l.peek() == 0 {
			// Emit any remaining text
			if l.pos > l.start {
				l.emit(TokenText)
			}
			l.emit(TokenEOF)

			return nil
		}

		// Only check for special constructs at line start
		if l.atLineStart() {
			// Check for blank line
			if l.peek() == '\n' {
				// Emit any accumulated text first
				if l.pos > l.start {
					l.emit(TokenText)
				}

				return lexBlankLine
			}

			// Check for code block
			if l.peek() == '`' && l.peekString(3) == "```" {
				// Emit any accumulated text first
				if l.pos > l.start {
					l.emit(TokenText)
				}

				return lexCodeBlock
			}

			// Check for header
			if l.peek() == '#' {
				// Emit any accumulated text first
				if l.pos > l.start {
					l.emit(TokenText)
				}

				return lexHeader
			}

			// Check for list item
			r := l.peek()
			if r == '-' || r == '*' {
				// Need to check if it's followed by space (proper list item)
				// Use peekString to look at marker + next char
				twoChars := l.peekString(2)
				if len(twoChars) == 2 {
					nextR := rune(twoChars[1])
					if nextR == ' ' || nextR == '\t' {
						// Emit any accumulated text first
						if l.pos > l.start {
							l.emit(TokenText)
						}

						return lexList
					}
				}
			}
		}

		// Consume regular text
		l.next()
	}
}

// lexHeader handles markdown headers.
//
// It counts the # characters (1-6 supported), captures the header text
// until end of line, and emits a TokenHeader. The token value includes
// the # characters and the text.
func lexHeader(l *Lexer) StateFn {
	// Consume # characters
	hashCount := l.acceptRun("#")

	// Headers must have at least one # and at most 6
	if hashCount == 0 || hashCount > 6 {
		l.emitError(fmt.Sprintf("invalid header: %d hash characters", hashCount))

		return lexText
	}

	// Consume optional space after #
	l.accept(" \t")

	// Consume until end of line
	l.acceptUntil("\n")

	// Emit the header token (includes # and text)
	l.emit(TokenHeader)

	// Consume the newline if present
	if l.peek() == '\n' {
		l.next()
		l.ignore()
	}

	return lexText
}

// lexCodeBlock handles code fences (```).
//
// This is critically important: it must consume everything between the
// opening and closing ``` markers without treating content as markdown.
// This prevents false positives when markdown syntax appears in code examples.
//
// The token value includes the opening fence, optional language specifier,
// all content, and the closing fence.
func lexCodeBlock(l *Lexer) StateFn {
	// Consume opening ```
	l.acceptRun("`")

	// Consume optional language specifier (rest of line)
	l.acceptUntil("\n")

	// Consume the newline after opening fence
	if l.peek() == '\n' {
		l.next()
	}

	// Now consume everything until closing ```
	for {
		r := l.next()

		if r == 0 {
			// EOF without closing fence - emit what we have
			l.emit(TokenCodeBlock)
			l.emitError("unclosed code block")

			return nil
		}

		// Check for closing fence (``` at line start)
		if r == '\n' && l.peek() == '`' && l.peekString(3) == "```" {
			// Consume the closing ```
			l.acceptRun("`")

			// Consume rest of line (some markdown allows text after closing fence)
			l.acceptUntil("\n")

			// Emit the code block (without trailing newline)
			l.emit(TokenCodeBlock)

			// Consume the newline after the fence but don't include it in token
			if l.peek() == '\n' {
				l.next()
				l.ignore()
			}

			return lexText
		}
	}
}

// lexList handles list items (- or *).
//
// It consumes the list marker, the space after it, and the content
// until end of line. The token value includes the marker and content.
func lexList(l *Lexer) StateFn {
	// Consume list marker (- or *)
	if !l.accept("-*") {
		l.emitError("expected list marker")

		return lexText
	}

	// Consume required space/tab after marker
	if !l.accept(" \t") {
		l.emitError("list marker must be followed by space")

		return lexText
	}

	// Consume rest of line
	l.acceptUntil("\n")

	l.emit(TokenListItem)

	// Consume the newline if present
	if l.peek() == '\n' {
		l.next()
		l.ignore()
	}

	return lexText
}

// lexBlankLine handles blank lines.
//
// It consumes one or more consecutive blank lines and emits a single
// TokenBlankLine. This helps preserve document structure.
func lexBlankLine(l *Lexer) StateFn {
	// Consume consecutive blank lines
	for l.peek() == '\n' {
		l.next()
	}

	l.emit(TokenBlankLine)

	return lexText
}

// peekString returns the next n characters without advancing position.
//
// Returns empty string if fewer than n characters remain.
func (l *Lexer) peekString(n int) string {
	if l.pos+n > len(l.input) {
		return ""
	}

	return l.input[l.pos : l.pos+n]
}
