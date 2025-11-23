package syntax

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// StateFn represents the state of the scanner as a function that returns the next state.
type StateFn func(*Lexer) StateFn

// Lexer holds the state of the scanner.
type Lexer struct {
	input  string     // the string being scanned
	start  int        // start position of this item
	pos    int        // current position in the input
	width  int        // width of last rune read from input
	tokens chan Token // channel of scanned items
	line   int        // current line number
}

// Lex creates a new lexer for the input string.
func Lex(input string) *Lexer {
	l := &Lexer{
		input:  input,
		tokens: make(chan Token),
		line:   1,
	}
	go l.run()
	return l
}

// run runs the state machine for the lexer.
func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// next returns the next rune in the input.
func (l *Lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return 0 // EOF
	}
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = w
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// peek returns the next rune in the input without consuming it.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
	if l.width == 1 && l.input[l.pos] == '\n' {
		l.line--
	}
}

// emit passes an item back to the client.
func (l *Lexer) emit(t TokenType) {
	l.tokens <- Token{t, l.input[l.start:l.pos], l.line}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *Lexer) ignore() {
	l.start = l.pos
}

// errorf returns an error token and terminates the scan.
func (l *Lexer) errorf(format string, args ...interface{}) StateFn {
	l.tokens <- Token{TokenError, fmt.Sprintf(format, args...), l.line}
	return nil
}

// NextToken returns the next token from the channel.
func (l *Lexer) NextToken() Token {
	return <-l.tokens
}

const (
	eof = 0
)

func lexText(l *Lexer) StateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], "```") {
			if l.pos > l.start {
				l.emit(TokenText)
			}
			return lexCodeBlock
		}
		if strings.HasPrefix(l.input[l.pos:], "#") {
			// Check if it's a header (must be at start of line or file)
			if l.pos == 0 || l.input[l.pos-1] == '\n' {
				if l.pos > l.start {
					l.emit(TokenText)
				}
				return lexHeader
			}
		}

		if strings.HasPrefix(l.input[l.pos:], "- ") || strings.HasPrefix(l.input[l.pos:], "* ") {
			if l.pos == 0 || l.input[l.pos-1] == '\n' {
				if l.pos > l.start {
					l.emit(TokenText)
				}
				return lexList
			}
		}

		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(TokenText)
	}
	l.emit(TokenEOF)
	return nil
}

func lexHeader(l *Lexer) StateFn {
	// Consume all #
	for l.peek() == '#' {
		l.next()
	}
	// Consume rest of line
	for {
		r := l.next()
		if r == '\n' || r == eof {
			break
		}
	}
	l.emit(TokenHeader)
	return lexText
}

func lexCodeBlock(l *Lexer) StateFn {
	l.pos += 3 // skip ```
	// Find end of code block
	for {
		if strings.HasPrefix(l.input[l.pos:], "```") {
			l.pos += 3
			l.emit(TokenCodeBlock)
			return lexText
		}
		if l.next() == eof {
			return l.errorf("unclosed code block")
		}
	}
}

func lexList(l *Lexer) StateFn {
	// Consume marker
	if strings.HasPrefix(l.input[l.pos:], "- ") {
		l.pos += 2
	} else if strings.HasPrefix(l.input[l.pos:], "* ") {
		l.pos += 2
	}

	// Consume rest of line
	for {
		r := l.next()
		if r == '\n' || r == eof {
			break
		}
	}
	l.emit(TokenList)
	return lexText
}
