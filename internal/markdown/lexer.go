//nolint:revive // file-length-limit,line-length-limit: lexer needs comprehensive tokenization
package markdown

import (
	"unicode/utf8"
)

// LexerState represents the current context of the lexer.
// Different states affect how characters are tokenized.
type LexerState uint8

const (
	// StateNormal is the default state for normal markdown content.
	StateNormal LexerState = iota
	// StateFencedCode is active inside a fenced code block (``` or ~~~).
	// In this state, only TokenText and TokenNewline are emitted.
	StateFencedCode
	// StateInlineCode is active inside inline code (backticks).
	// Content is emitted as TokenText until matching backticks.
	StateInlineCode
	// StateLinkURL is active inside a link URL parentheses.
	// Special characters are treated as TokenText.
	StateLinkURL
)

// Lexer constants for parsing thresholds.
const (
	minFenceLength  = 3    // Minimum backticks/tildes for code fence
	utf8MultiByteHi = 0x80 // High bit set indicates multi-byte UTF-8
)

// LexError represents an error encountered during lexing.
// It captures position and message from TokenError tokens.
type LexError struct {
	Offset  int    // Byte offset where error occurred
	Message string // Error description
}

// Error implements the error interface.
func (e LexError) Error() string {
	return e.Message
}

// lexer tokenizes markdown source into a stream of fine-grained tokens.
// It is internal to the package - use Parse() for public API.
type lexer struct {
	source []byte // Original source (retained, not copied)
	pos    int    // Current byte offset in source

	// State machine for context-dependent tokenization
	state         LexerState
	fenceChar     byte // '`' or '~' for code fence matching
	fenceLen      int  // Length of opening fence (3+)
	backtickCount int  // Number of backticks for inline code matching

	// Peek caching to avoid re-lexing
	peeked  Token
	hasPeek bool
	atEOF   bool // Sticky EOF flag
}

// newLexer creates a new lexer for the given source.
// The source is retained by reference (not copied).
// Caller must ensure source is not modified during lexing.
func newLexer(source []byte) *lexer {
	return &lexer{
		source: source,
		pos:    0,
		state:  StateNormal,
	}
}

// Next returns the next token from the source and advances position.
// Repeated calls return subsequent tokens until TokenEOF.
// After EOF is returned, subsequent calls continue returning TokenEOF.
func (l *lexer) Next() Token {
	// Return cached peek token if available
	if l.hasPeek {
		l.hasPeek = false
		tok := l.peeked
		l.peeked = Token{} // Clear for GC

		return tok
	}

	return l.nextToken()
}

// Peek returns the next token WITHOUT advancing position.
// The token is cached, so subsequent Peek() or Next() returns the same token.
func (l *lexer) Peek() Token {
	if l.hasPeek {
		return l.peeked
	}

	l.peeked = l.nextToken()
	l.hasPeek = true

	return l.peeked
}

// All returns all tokens from the current position to EOF.
// The returned slice includes TokenEOF as the final element.
func (l *lexer) All() []Token {
	var tokens []Token
	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}

	return tokens
}

// AllWithErrors returns all tokens and a separate slice of lex errors.
// Error tokens remain in the token slice but are also extracted into the error slice.
func (l *lexer) AllWithErrors() ([]Token, []LexError) {
	var tokens []Token
	var errors []LexError

	for {
		tok := l.Next()
		tokens = append(tokens, tok)
		if tok.Type == TokenError {
			errors = append(errors, LexError{
				Offset:  tok.Start,
				Message: tok.Message,
			})
		}
		if tok.Type == TokenEOF {
			break
		}
	}

	return tokens, errors
}

// nextToken produces the next token from source.
// This is the core lexing logic.
func (l *lexer) nextToken() Token {
	// Handle EOF
	if l.pos >= len(l.source) || l.atEOF {
		l.atEOF = true

		return Token{
			Type:   TokenEOF,
			Start:  len(l.source),
			End:    len(l.source),
			Source: l.source[len(l.source):],
		}
	}

	// Dispatch based on current state
	switch l.state {
	case StateNormal:
		return l.lexNormal()
	case StateFencedCode:
		return l.lexFencedCodeContent()
	case StateInlineCode:
		return l.lexInlineCodeContent()
	case StateLinkURL:
		return l.lexLinkURLContent()
	default:
		return l.lexNormal()
	}
}

// lexNormal handles tokenization in the default state.
//
//nolint:revive // function-length - lexer state machine requires many branches
func (l *lexer) lexNormal() Token {
	start := l.pos
	b := l.source[l.pos]

	// Handle newlines (with CRLF normalization)
	if b == '\n' {
		l.pos++

		return Token{
			Type:   TokenNewline,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	}
	if b == '\r' {
		l.pos++
		// Check for CRLF
		if l.pos < len(l.source) &&
			l.source[l.pos] == '\n' {
			l.pos++
		}

		return Token{
			Type:   TokenNewline,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	}

	// Handle ASCII whitespace (space and tab only)
	if b == ' ' || b == '\t' {
		return l.lexWhitespace()
	}

	// Handle single-character delimiters
	switch b {
	case '#':
		l.pos++

		return Token{
			Type:   TokenHash,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '*':
		l.pos++

		return Token{
			Type:   TokenAsterisk,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '_':
		l.pos++

		return Token{
			Type:   TokenUnderscore,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '~':
		return l.lexTilde()
	case '`':
		return l.lexBacktick()
	case '-':
		l.pos++

		return Token{
			Type:   TokenDash,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '+':
		l.pos++

		return Token{
			Type:   TokenPlus,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '.':
		l.pos++

		return Token{
			Type:   TokenDot,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case ':':
		l.pos++

		return Token{
			Type:   TokenColon,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '|':
		l.pos++

		return Token{
			Type:   TokenPipe,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '[':
		l.pos++

		return Token{
			Type:   TokenBracketOpen,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case ']':
		l.pos++

		return Token{
			Type:   TokenBracketClose,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '(':
		l.pos++

		return Token{
			Type:   TokenParenOpen,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case ')':
		l.pos++

		return Token{
			Type:   TokenParenClose,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	case '>':
		l.pos++

		return Token{
			Type:   TokenGreaterThan,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	}

	// Handle digits (for ordered lists)
	if b >= '0' && b <= '9' {
		return l.lexNumber()
	}

	// Handle 'x' or 'X' for checkboxes
	if b == 'x' || b == 'X' {
		l.pos++

		return Token{
			Type:   TokenX,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	}

	// Handle text (everything else)
	return l.lexText()
}

// lexWhitespace consumes contiguous ASCII whitespace (space and tab).
func (l *lexer) lexWhitespace() Token {
	start := l.pos
	for l.pos < len(l.source) {
		b := l.source[l.pos]
		if b != ' ' && b != '\t' {
			break
		}
		l.pos++
	}

	return Token{
		Type:   TokenWhitespace,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// lexNumber consumes a sequence of digits.
func (l *lexer) lexNumber() Token {
	start := l.pos
	for l.pos < len(l.source) {
		b := l.source[l.pos]
		if b < '0' || b > '9' {
			break
		}
		l.pos++
	}

	return Token{
		Type:   TokenNumber,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// lexBacktick handles backtick characters.
// It counts consecutive backticks and may transition to inline code or fenced code state.
func (l *lexer) lexBacktick() Token {
	start := l.pos
	count := 0

	// Count consecutive backticks
	for l.pos < len(l.source) && l.source[l.pos] == '`' {
		l.pos++
		count++
	}

	// Check if this is a code fence (3+ backticks at line start)
	// A line start is either position 0 or immediately after a newline
	atLineStart := start == 0 ||
		(start > 0 && (l.source[start-1] == '\n' || l.source[start-1] == '\r'))

	if count >= minFenceLength && atLineStart {
		// Enter fenced code state
		l.state = StateFencedCode
		l.fenceChar = '`'
		l.fenceLen = count
	} else if count > 0 && l.state == StateNormal {
		// Enter inline code state
		l.state = StateInlineCode
		l.backtickCount = count
	}

	// Return individual backtick tokens
	// Reset position to after first backtick only
	l.pos = start + 1

	return Token{
		Type:   TokenBacktick,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// lexTilde handles tilde characters.
// It counts consecutive tildes and may transition to fenced code state.
func (l *lexer) lexTilde() Token {
	start := l.pos
	count := 0

	// Count consecutive tildes
	for l.pos < len(l.source) && l.source[l.pos] == '~' {
		l.pos++
		count++
	}

	// Check if this is a code fence (3+ tildes at line start)
	// A line start is either position 0 or immediately after a newline
	atLineStart := start == 0 ||
		(start > 0 && (l.source[start-1] == '\n' || l.source[start-1] == '\r'))

	if count >= minFenceLength && atLineStart {
		// Enter fenced code state
		l.state = StateFencedCode
		l.fenceChar = '~'
		l.fenceLen = count
	}

	// Return individual tilde tokens
	// Reset position to after first tilde only
	l.pos = start + 1

	return Token{
		Type:   TokenTilde,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// lexText consumes text content (non-delimiter, non-whitespace characters).
// It handles multi-byte UTF-8 sequences correctly.
func (l *lexer) lexText() Token {
	start := l.pos

	for l.pos < len(l.source) {
		b := l.source[l.pos]

		// Stop at delimiters, whitespace, and newlines
		if isDelimiterByte(b) || b == ' ' ||
			b == '\t' ||
			b == '\n' ||
			b == '\r' {
			break
		}

		// Handle UTF-8 multi-byte sequences
		if b >= utf8MultiByteHi {
			r, size := utf8.DecodeRune(
				l.source[l.pos:],
			)
			if r == utf8.RuneError && size == 1 {
				// Invalid UTF-8: emit what we have so far as text, then handle error
				if l.pos > start {
					return Token{
						Type:   TokenText,
						Start:  start,
						End:    l.pos,
						Source: l.source[start:l.pos],
					}
				}
				// Emit error for invalid byte
				l.pos++

				return Token{
					Type:    TokenError,
					Start:   start,
					End:     l.pos,
					Source:  l.source[start:l.pos],
					Message: "invalid UTF-8 byte sequence",
				}
			}
			l.pos += size

			continue
		}

		l.pos++
	}

	if l.pos == start {
		// No text consumed - should not happen in normal flow
		// Safety: advance by one byte and emit error
		l.pos++

		return Token{
			Type:    TokenError,
			Start:   start,
			End:     l.pos,
			Source:  l.source[start:l.pos],
			Message: "unexpected character",
		}
	}

	return Token{
		Type:   TokenText,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// lexFencedCodeContent handles tokenization inside a fenced code block.
// Only TokenText and TokenNewline are emitted until the closing fence.
//
//nolint:revive // function-length: complex state machine for code fence parsing
func (l *lexer) lexFencedCodeContent() Token {
	start := l.pos

	// Check for newline first
	if l.pos < len(l.source) {
		b := l.source[l.pos]
		if b == '\n' {
			l.pos++

			return Token{
				Type:   TokenNewline,
				Start:  start,
				End:    l.pos,
				Source: l.source[start:l.pos],
			}
		}
		if b == '\r' {
			l.pos++
			if l.pos < len(l.source) &&
				l.source[l.pos] == '\n' {
				l.pos++
			}

			return Token{
				Type:   TokenNewline,
				Start:  start,
				End:    l.pos,
				Source: l.source[start:l.pos],
			}
		}
	}

	// Check if we're at line start and this might be a closing fence
	atLineStart := start == 0 ||
		(start > 0 && (l.source[start-1] == '\n' || l.source[start-1] == '\r'))

	if atLineStart && l.pos < len(l.source) &&
		l.source[l.pos] == l.fenceChar {
		// Count fence characters
		fenceStart := l.pos
		count := 0
		for l.pos < len(l.source) && l.source[l.pos] == l.fenceChar {
			l.pos++
			count++
		}

		// Check if this closes the fence (same or more characters)
		if count >= l.fenceLen {
			// Verify rest of line is blank or whitespace
			valid := true
			checkPos := l.pos
			for checkPos < len(l.source) {
				c := l.source[checkPos]
				if c == '\n' || c == '\r' {
					break
				}
				if c != ' ' && c != '\t' {
					valid = false

					break
				}
				checkPos++
			}

			if valid {
				// Exit fenced code state
				l.state = StateNormal
				l.fenceChar = 0
				l.fenceLen = 0

				// Return the fence as text (closing fence)
				return Token{
					Type:   TokenText,
					Start:  fenceStart,
					End:    l.pos,
					Source: l.source[fenceStart:l.pos],
				}
			}
		}

		// Not a closing fence - reset and treat as content
		l.pos = fenceStart
	}

	// Consume text until newline
	for l.pos < len(l.source) {
		b := l.source[l.pos]
		if b == '\n' || b == '\r' {
			break
		}
		l.pos++
	}

	if l.pos == start {
		// Nothing consumed - shouldn't happen, safety
		return l.nextToken()
	}

	return Token{
		Type:   TokenText,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// lexInlineCodeContent handles tokenization inside inline code.
// Content is emitted as TokenText until matching backtick sequence.
//
//nolint:revive // function-length: complex state machine for inline code parsing
func (l *lexer) lexInlineCodeContent() Token {
	start := l.pos

	// Check for closing backticks
	if l.pos < len(l.source) &&
		l.source[l.pos] == '`' {
		// Count consecutive backticks
		count := 0
		for l.pos < len(l.source) && l.source[l.pos] == '`' {
			l.pos++
			count++
		}

		if count == l.backtickCount {
			// Exit inline code state
			l.state = StateNormal
			l.backtickCount = 0

			// Emit individual backtick tokens
			// Reset to emit just the first backtick
			l.pos = start + 1

			return Token{
				Type:   TokenBacktick,
				Start:  start,
				End:    l.pos,
				Source: l.source[start:l.pos],
			}
		}

		// Not matching - treat as text content
		return Token{
			Type:   TokenText,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	}

	// Handle newline
	if l.pos < len(l.source) {
		b := l.source[l.pos]
		if b == '\n' {
			l.pos++

			return Token{
				Type:   TokenNewline,
				Start:  start,
				End:    l.pos,
				Source: l.source[start:l.pos],
			}
		}
		if b == '\r' {
			l.pos++
			if l.pos < len(l.source) &&
				l.source[l.pos] == '\n' {
				l.pos++
			}

			return Token{
				Type:   TokenNewline,
				Start:  start,
				End:    l.pos,
				Source: l.source[start:l.pos],
			}
		}
	}

	// Consume text until backtick or newline
	for l.pos < len(l.source) {
		b := l.source[l.pos]
		if b == '`' || b == '\n' || b == '\r' {
			break
		}
		l.pos++
	}

	if l.pos == start {
		// Safety: should not happen
		l.pos++

		return Token{
			Type:    TokenError,
			Start:   start,
			End:     l.pos,
			Source:  l.source[start:l.pos],
			Message: "unexpected state in inline code",
		}
	}

	return Token{
		Type:   TokenText,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// lexLinkURLContent handles tokenization inside link URL parentheses.
// Special characters are treated as TokenText (URLs can contain *, _, etc.).
//
//nolint:revive // function-length: complex state machine for URL parsing
func (l *lexer) lexLinkURLContent() Token {
	start := l.pos

	// Check for closing parenthesis
	if l.pos < len(l.source) &&
		l.source[l.pos] == ')' {
		l.pos++
		l.state = StateNormal

		return Token{
			Type:   TokenParenClose,
			Start:  start,
			End:    l.pos,
			Source: l.source[start:l.pos],
		}
	}

	// Handle newline (URLs shouldn't span lines, exit state)
	if l.pos < len(l.source) {
		b := l.source[l.pos]
		if b == '\n' {
			l.pos++
			l.state = StateNormal

			return Token{
				Type:   TokenNewline,
				Start:  start,
				End:    l.pos,
				Source: l.source[start:l.pos],
			}
		}
		if b == '\r' {
			l.pos++
			if l.pos < len(l.source) &&
				l.source[l.pos] == '\n' {
				l.pos++
			}
			l.state = StateNormal

			return Token{
				Type:   TokenNewline,
				Start:  start,
				End:    l.pos,
				Source: l.source[start:l.pos],
			}
		}
	}

	// Handle whitespace (for separating URL from title)
	if l.pos < len(l.source) {
		b := l.source[l.pos]
		if b == ' ' || b == '\t' {
			return l.lexWhitespace()
		}
	}

	// Consume URL text until ) or whitespace or newline
	for l.pos < len(l.source) {
		b := l.source[l.pos]
		if b == ')' || b == ' ' || b == '\t' ||
			b == '\n' ||
			b == '\r' {
			break
		}
		l.pos++
	}

	if l.pos == start {
		// Nothing consumed - shouldn't happen
		l.pos++

		return Token{
			Type:    TokenError,
			Start:   start,
			End:     l.pos,
			Source:  l.source[start:l.pos],
			Message: "unexpected state in link URL",
		}
	}

	return Token{
		Type:   TokenText,
		Start:  start,
		End:    l.pos,
		Source: l.source[start:l.pos],
	}
}

// enterLinkURLState transitions the lexer to StateLinkURL.
// This should be called by the parser when it detects ]( sequence.
func (l *lexer) enterLinkURLState() {
	l.state = StateLinkURL
}

// exitLinkURLState returns the lexer to StateNormal.
func (l *lexer) exitLinkURLState() {
	l.state = StateNormal
}

// enterFencedCodeState transitions to StateFencedCode.
// char should be '`' or '~', length is the fence length (3+).
func (l *lexer) enterFencedCodeState(
	char byte,
	length int,
) {
	l.state = StateFencedCode
	l.fenceChar = char
	l.fenceLen = length
}

// enterInlineCodeState transitions to StateInlineCode.
// count is the number of opening backticks.
func (l *lexer) enterInlineCodeState(count int) {
	l.state = StateInlineCode
	l.backtickCount = count
}

// State returns the current lexer state.
func (l *lexer) State() LexerState {
	return l.state
}

// Pos returns the current byte offset in the source.
func (l *lexer) Pos() int {
	return l.pos
}

// isDelimiterByte returns true if the byte is a markdown delimiter character.
func isDelimiterByte(b byte) bool {
	switch b {
	case '#',
		'*',
		'_',
		'~',
		'`',
		'-',
		'+',
		'.',
		':',
		'|',
		'[',
		']',
		'(',
		')',
		'>':
		return true
	default:
		return false
	}
}
