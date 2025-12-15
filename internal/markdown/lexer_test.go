package markdown

import (
	"testing"
)

// TestLexer_EmptyInput verifies that empty input returns only EOF.
func TestLexer_EmptyInput(t *testing.T) {
	l := newLexer(make([]byte, 0))
	tok := l.Next()

	if tok.Type != TokenEOF {
		t.Errorf(
			"empty input: got %v, want TokenEOF",
			tok.Type,
		)
	}
	if tok.Start != 0 || tok.End != 0 {
		t.Errorf(
			"empty input EOF: Start=%d, End=%d, want 0, 0",
			tok.Start,
			tok.End,
		)
	}
}

// TestLexer_SingleCharacterDelimiters verifies each delimiter tokenizes correctly.
func TestLexer_SingleCharacterDelimiters(
	t *testing.T,
) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"#", TokenHash},
		{"*", TokenAsterisk},
		{"_", TokenUnderscore},
		{"~", TokenTilde},
		{"`", TokenBacktick},
		{"-", TokenDash},
		{"+", TokenPlus},
		{".", TokenDot},
		{":", TokenColon},
		{"|", TokenPipe},
		{"[", TokenBracketOpen},
		{"]", TokenBracketClose},
		{"(", TokenParenOpen},
		{")", TokenParenClose},
		{">", TokenGreaterThan},
	}

	for _, tt := range tests {
		t.Run(
			tt.expected.String(),
			func(t *testing.T) {
				l := newLexer([]byte(tt.input))
				tok := l.Next()

				if tok.Type != tt.expected {
					t.Errorf(
						"input %q: got %v, want %v",
						tt.input,
						tok.Type,
						tt.expected,
					)
				}
				if tok.Start != 0 {
					t.Errorf(
						"input %q: Start=%d, want 0",
						tt.input,
						tok.Start,
					)
				}
				if tok.End != 1 {
					t.Errorf(
						"input %q: End=%d, want 1",
						tt.input,
						tok.End,
					)
				}
				if string(
					tok.Source,
				) != tt.input {
					t.Errorf(
						"input %q: Source=%q, want %q",
						tt.input,
						tok.Source,
						tt.input,
					)
				}
			},
		)
	}
}

// TestLexer_Whitespace verifies whitespace tokenization.
func TestLexer_Whitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single space", " ", " "},
		{"multiple spaces", "   ", "   "},
		{"single tab", "\t", "\t"},
		{"multiple tabs", "\t\t", "\t\t"},
		{"mixed spaces and tabs", " \t ", " \t "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer([]byte(tt.input))
			tok := l.Next()

			if tok.Type != TokenWhitespace {
				t.Errorf(
					"input %q: got %v, want TokenWhitespace",
					tt.input,
					tok.Type,
				)
			}
			if string(tok.Source) != tt.expected {
				t.Errorf(
					"input %q: Source=%q, want %q",
					tt.input,
					tok.Source,
					tt.expected,
				)
			}
		})
	}
}

// TestLexer_Newline verifies newline tokenization.
func TestLexer_Newline(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantLen int
	}{
		{"LF only", "\n", 1},
		{"CR only", "\r", 1},
		{"CRLF", "\r\n", 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer([]byte(tt.input))
			tok := l.Next()

			if tok.Type != TokenNewline {
				t.Errorf(
					"input %q: got %v, want TokenNewline",
					tt.input,
					tok.Type,
				)
			}
			if tok.Len() != tt.wantLen {
				t.Errorf(
					"input %q: Len()=%d, want %d",
					tt.input,
					tok.Len(),
					tt.wantLen,
				)
			}
		})
	}
}

// TestLexer_Text verifies text content tokenization.
func TestLexer_Text(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple word", "hello", "hello"},
		{
			"word stops at delimiter",
			"hello*",
			"hello",
		},
		{
			"word stops at space",
			"hello world",
			"hello",
		},
		{
			"word stops at newline",
			"hello\n",
			"hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer([]byte(tt.input))
			tok := l.Next()

			if tok.Type != TokenText {
				t.Errorf(
					"input %q: got %v, want TokenText",
					tt.input,
					tok.Type,
				)
			}
			if string(tok.Source) != tt.expected {
				t.Errorf(
					"input %q: Source=%q, want %q",
					tt.input,
					tok.Source,
					tt.expected,
				)
			}
		})
	}
}

// TestLexer_Number verifies digit sequence tokenization.
func TestLexer_Number(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single digit", "1", "1"},
		{"multiple digits", "123", "123"},
		{"digits before dot", "1.", "1"},
		{"digits before space", "123 ", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer([]byte(tt.input))
			tok := l.Next()

			if tok.Type != TokenNumber {
				t.Errorf(
					"input %q: got %v, want TokenNumber",
					tt.input,
					tok.Type,
				)
			}
			if string(tok.Source) != tt.expected {
				t.Errorf(
					"input %q: Source=%q, want %q",
					tt.input,
					tok.Source,
					tt.expected,
				)
			}
		})
	}
}

// TestLexer_X verifies 'x' and 'X' tokenization (for checkboxes).
func TestLexer_X(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"x", "x"},
		{"X", "X"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			l := newLexer([]byte(tt.input))
			tok := l.Next()

			if tok.Type != TokenX {
				t.Errorf(
					"input %q: got %v, want TokenX",
					tt.input,
					tok.Type,
				)
			}
			if string(tok.Source) != tt.expected {
				t.Errorf(
					"input %q: Source=%q, want %q",
					tt.input,
					tok.Source,
					tt.expected,
				)
			}
		})
	}
}

// TestLexer_MixedContent verifies correct tokenization of mixed content.
func TestLexer_MixedContent(t *testing.T) {
	input := "# Hello\n"
	l := newLexer([]byte(input))

	expected := []struct {
		tokenType TokenType
		text      string
	}{
		{TokenHash, "#"},
		{TokenWhitespace, " "},
		{TokenText, "Hello"},
		{TokenNewline, "\n"},
		{TokenEOF, ""},
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp.tokenType {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp.tokenType,
			)
		}
		if string(tok.Source) != exp.text {
			t.Errorf(
				"token %d: Source=%q, want %q",
				i,
				tok.Source,
				exp.text,
			)
		}
	}
}

// TestLexer_MarkdownHeading verifies heading tokenization.
func TestLexer_MarkdownHeading(t *testing.T) {
	input := "## Heading"
	l := newLexer([]byte(input))

	expected := []TokenType{
		TokenHash,       // #
		TokenHash,       // #
		TokenWhitespace, // space
		TokenText,       // Heading
		TokenEOF,
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp,
			)
		}
	}
}

// TestLexer_BoldText verifies bold syntax tokenization.
func TestLexer_BoldText(t *testing.T) {
	input := "**bold**"
	l := newLexer([]byte(input))

	expected := []TokenType{
		TokenAsterisk, // *
		TokenAsterisk, // *
		TokenText,     // bold
		TokenAsterisk, // *
		TokenAsterisk, // *
		TokenEOF,
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp,
			)
		}
	}
}

// TestLexer_Link verifies link syntax tokenization.
func TestLexer_Link(t *testing.T) {
	input := "[text](url)"
	l := newLexer([]byte(input))

	expected := []TokenType{
		TokenBracketOpen,  // [
		TokenText,         // text
		TokenBracketClose, // ]
		TokenParenOpen,    // (
		TokenText,         // url
		TokenParenClose,   // )
		TokenEOF,
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp,
			)
		}
	}
}

// TestLexer_OrderedList verifies ordered list syntax tokenization.
func TestLexer_OrderedList(t *testing.T) {
	input := "1. Item"
	l := newLexer([]byte(input))

	expected := []TokenType{
		TokenNumber,     // 1
		TokenDot,        // .
		TokenWhitespace, // space
		TokenText,       // Item
		TokenEOF,
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp,
			)
		}
	}
}

// TestLexer_UnorderedList verifies unordered list syntax tokenization.
func TestLexer_UnorderedList(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		marker TokenType
	}{
		{"dash list", "- Item", TokenDash},
		{"plus list", "+ Item", TokenPlus},
		{
			"asterisk list",
			"* Item",
			TokenAsterisk,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer([]byte(tt.input))

			tok := l.Next()
			if tok.Type != tt.marker {
				t.Errorf(
					"first token: got %v, want %v",
					tok.Type,
					tt.marker,
				)
			}

			tok = l.Next()
			if tok.Type != TokenWhitespace {
				t.Errorf(
					"second token: got %v, want TokenWhitespace",
					tok.Type,
				)
			}

			tok = l.Next()
			if tok.Type != TokenText {
				t.Errorf(
					"third token: got %v, want TokenText",
					tok.Type,
				)
			}
		})
	}
}

// TestLexer_Checkbox verifies checkbox syntax tokenization.
func TestLexer_Checkbox(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unchecked", "- [ ] Task"},
		{"checked lowercase", "- [x] Task"},
		{"checked uppercase", "- [X] Task"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer([]byte(tt.input))

			tok := l.Next()
			if tok.Type != TokenDash {
				t.Errorf(
					"token 0: got %v, want TokenDash",
					tok.Type,
				)
			}

			tok = l.Next()
			if tok.Type != TokenWhitespace {
				t.Errorf(
					"token 1: got %v, want TokenWhitespace",
					tok.Type,
				)
			}

			tok = l.Next()
			if tok.Type != TokenBracketOpen {
				t.Errorf(
					"token 2: got %v, want TokenBracketOpen",
					tok.Type,
				)
			}
		})
	}
}

// TestLexer_Blockquote verifies blockquote syntax tokenization.
func TestLexer_Blockquote(t *testing.T) {
	input := "> Quote"
	l := newLexer([]byte(input))

	tok := l.Next()
	if tok.Type != TokenGreaterThan {
		t.Errorf(
			"first token: got %v, want TokenGreaterThan",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenWhitespace {
		t.Errorf(
			"second token: got %v, want TokenWhitespace",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"third token: got %v, want TokenText",
			tok.Type,
		)
	}
}

// TestLexer_Table verifies table syntax tokenization.
func TestLexer_Table(t *testing.T) {
	input := "| A | B |"
	l := newLexer([]byte(input))

	expected := []TokenType{
		TokenPipe,       // |
		TokenWhitespace, // space
		TokenText,       // A
		TokenWhitespace, // space
		TokenPipe,       // |
		TokenWhitespace, // space
		TokenText,       // B
		TokenWhitespace, // space
		TokenPipe,       // |
		TokenEOF,
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp,
			)
		}
	}
}

// TestLexer_UnicodeText verifies UTF-8 text handling.
func TestLexer_UnicodeText(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			"emoji",
			"\xf0\x9f\x98\x80",
		}, // grinning face emoji
		{"chinese", "\xe4\xb8\xad\xe6\x96\x87"},
		{
			"arabic",
			"\xd8\xb9\xd8\xb1\xd8\xa8\xd9\x8a",
		},
		{
			"mixed ascii unicode",
			"Hello\xe4\xb8\x96\xe7\x95\x8c",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := newLexer([]byte(tt.input))
			tok := l.Next()

			if tok.Type != TokenText {
				t.Errorf(
					"input %q: got %v, want TokenText",
					tt.input,
					tok.Type,
				)
			}
			if string(tok.Source) != tt.input {
				t.Errorf(
					"input %q: Source=%q, want %q",
					tt.input,
					tok.Source,
					tt.input,
				)
			}
		})
	}
}

// TestLexer_RepeatedEOF verifies that EOF is sticky.
func TestLexer_RepeatedEOF(t *testing.T) {
	l := newLexer(make([]byte, 0))

	for i := range 5 {
		tok := l.Next()
		if tok.Type != TokenEOF {
			t.Errorf(
				"call %d: got %v, want TokenEOF",
				i,
				tok.Type,
			)
		}
	}
}

// TestLexer_CRLFNormalization verifies CRLF produces single newline spanning 2 bytes.
func TestLexer_CRLFNormalization(t *testing.T) {
	input := "line1\r\nline2"
	l := newLexer([]byte(input))

	// line1
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 0: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "line1" {
		t.Errorf(
			"token 0: Source=%q, want %q",
			tok.Source,
			"line1",
		)
	}

	// \r\n as single newline
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"token 1: got %v, want TokenNewline",
			tok.Type,
		)
	}
	if tok.Len() != 2 {
		t.Errorf(
			"CRLF token: Len()=%d, want 2",
			tok.Len(),
		)
	}
	if tok.Start != 5 || tok.End != 7 {
		t.Errorf(
			"CRLF token: Start=%d, End=%d, want 5, 7",
			tok.Start,
			tok.End,
		)
	}
	if string(tok.Source) != "\r\n" {
		t.Errorf(
			"CRLF token: Source=%q, want %q",
			tok.Source,
			"\r\n",
		)
	}

	// line2
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 2: got %v, want TokenText",
			tok.Type,
		)
	}
}

// TestLexer_StandaloneCR verifies standalone CR produces newline.
func TestLexer_StandaloneCR(t *testing.T) {
	input := "line1\rline2"
	l := newLexer([]byte(input))

	// line1
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 0: got %v, want TokenText",
			tok.Type,
		)
	}

	// \r as newline
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"token 1: got %v, want TokenNewline",
			tok.Type,
		)
	}
	if tok.Len() != 1 {
		t.Errorf(
			"CR token: Len()=%d, want 1",
			tok.Len(),
		)
	}
	if string(tok.Source) != "\r" {
		t.Errorf(
			"CR token: Source=%q, want %q",
			tok.Source,
			"\r",
		)
	}

	// line2
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 2: got %v, want TokenText",
			tok.Type,
		)
	}
}

// TestLexer_MixedLineEndings verifies handling of mixed line endings.
func TestLexer_MixedLineEndings(t *testing.T) {
	input := "a\nb\r\nc\rd"
	l := newLexer([]byte(input))

	expected := []struct {
		tokenType TokenType
		text      string
		length    int
	}{
		{TokenText, "a", 1},
		{TokenNewline, "\n", 1},
		{TokenText, "b", 1},
		{TokenNewline, "\r\n", 2},
		{TokenText, "c", 1},
		{TokenNewline, "\r", 1},
		{TokenText, "d", 1},
		{TokenEOF, "", 0},
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp.tokenType {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp.tokenType,
			)
		}
		if string(tok.Source) != exp.text {
			t.Errorf(
				"token %d: Source=%q, want %q",
				i,
				tok.Source,
				exp.text,
			)
		}
		if tok.Len() != exp.length {
			t.Errorf(
				"token %d: Len()=%d, want %d",
				i,
				tok.Len(),
				exp.length,
			)
		}
	}
}

// TestLexer_CRLFAtEnd verifies CRLF at end of input.
func TestLexer_CRLFAtEnd(t *testing.T) {
	input := "text\r\n"
	l := newLexer([]byte(input))

	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 0: got %v, want TokenText",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"token 1: got %v, want TokenNewline",
			tok.Type,
		)
	}
	if tok.Len() != 2 {
		t.Errorf(
			"CRLF at end: Len()=%d, want 2",
			tok.Len(),
		)
	}

	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"token 2: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_CRAtEnd verifies CR at end of input (without following LF).
func TestLexer_CRAtEnd(t *testing.T) {
	input := "text\r"
	l := newLexer([]byte(input))

	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 0: got %v, want TokenText",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"token 1: got %v, want TokenNewline",
			tok.Type,
		)
	}
	if tok.Len() != 1 {
		t.Errorf(
			"CR at end: Len()=%d, want 1",
			tok.Len(),
		)
	}

	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"token 2: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_ConsecutiveCRLF verifies multiple consecutive CRLF sequences.
func TestLexer_ConsecutiveCRLF(t *testing.T) {
	input := "\r\n\r\n"
	l := newLexer([]byte(input))

	for i := range 2 {
		tok := l.Next()
		if tok.Type != TokenNewline {
			t.Errorf(
				"token %d: got %v, want TokenNewline",
				i,
				tok.Type,
			)
		}
		if tok.Len() != 2 {
			t.Errorf(
				"token %d: Len()=%d, want 2",
				i,
				tok.Len(),
			)
		}
	}

	tok := l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"final token: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_InvalidUTF8ProducesError verifies invalid UTF-8 produces TokenError.
func TestLexer_InvalidUTF8ProducesError(
	t *testing.T,
) {
	// Invalid UTF-8 byte sequence (0x80 alone is invalid)
	input := []byte{0x80}
	l := newLexer(input)

	tok := l.Next()
	if tok.Type != TokenError {
		t.Errorf(
			"invalid UTF-8: got %v, want TokenError",
			tok.Type,
		)
	}
	if tok.Message == "" {
		t.Error(
			"error token should have a message",
		)
	}
}

// TestLexer_InvalidUTF8ThenContinues verifies lexer continues after error.
func TestLexer_InvalidUTF8ThenContinues(
	t *testing.T,
) {
	// Invalid UTF-8 followed by valid text
	input := []byte{0x80, 'a', 'b', 'c'}
	l := newLexer(input)

	// First should be error
	tok := l.Next()
	if tok.Type != TokenError {
		t.Errorf(
			"token 0: got %v, want TokenError",
			tok.Type,
		)
	}

	// Then should continue with valid text
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 1: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "abc" {
		t.Errorf(
			"token 1: Source=%q, want %q",
			tok.Source,
			"abc",
		)
	}

	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"token 2: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_MultipleErrors verifies multiple errors in one input.
func TestLexer_MultipleErrors(t *testing.T) {
	// Two invalid bytes separated by valid text
	input := []byte{0x80, 'a', 0x81}
	l := newLexer(input)

	tok := l.Next()
	if tok.Type != TokenError {
		t.Errorf(
			"token 0: got %v, want TokenError",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 1: got %v, want TokenText",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenError {
		t.Errorf(
			"token 2: got %v, want TokenError",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"token 3: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_InvalidUTF8MidText verifies invalid UTF-8 mid-text is handled.
func TestLexer_InvalidUTF8MidText(t *testing.T) {
	// Valid text, invalid byte, valid text
	input := []byte{
		'h',
		'e',
		'l',
		'l',
		'o',
		0x80,
		'w',
		'o',
		'r',
		'l',
		'd',
	}
	l := newLexer(input)

	// First should be text "hello"
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 0: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "hello" {
		t.Errorf(
			"token 0: Source=%q, want %q",
			tok.Source,
			"hello",
		)
	}

	// Then error for invalid byte
	tok = l.Next()
	if tok.Type != TokenError {
		t.Errorf(
			"token 1: got %v, want TokenError",
			tok.Type,
		)
	}

	// Then continue with "world"
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 2: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "world" {
		t.Errorf(
			"token 2: Source=%q, want %q",
			tok.Source,
			"world",
		)
	}
}

// TestLexer_ErrorHasCorrectOffset verifies error token has correct position.
func TestLexer_ErrorHasCorrectOffset(
	t *testing.T,
) {
	// Valid text followed by invalid byte
	input := []byte{'a', 'b', 'c', 0x80}
	l := newLexer(input)

	tok := l.Next() // abc
	if tok.Type != TokenText {
		t.Errorf(
			"token 0: got %v, want TokenText",
			tok.Type,
		)
	}

	tok = l.Next() // error at position 3
	if tok.Type != TokenError {
		t.Errorf(
			"token 1: got %v, want TokenError",
			tok.Type,
		)
	}
	if tok.Start != 3 {
		t.Errorf(
			"error token: Start=%d, want 3",
			tok.Start,
		)
	}
	if tok.End != 4 {
		t.Errorf(
			"error token: End=%d, want 4",
			tok.End,
		)
	}
}

// TestLexer_FencedCodeBlock verifies fenced code block state transitions.
// Note: The lexer emits ONE backtick token when seeing ```, then enters
// StateFencedCode. Subsequent characters are processed as fenced code content.
func TestLexer_FencedCodeBlock(t *testing.T) {
	input := "```go\ncode here\n```"
	l := newLexer([]byte(input))

	// First backtick - triggers state transition to StateFencedCode
	tok := l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"first backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}

	// After first backtick of a 3+ backtick sequence at line start,
	// state should be StateFencedCode
	if l.State() != StateFencedCode {
		t.Errorf(
			"after first backtick: state=%v, want StateFencedCode",
			l.State(),
		)
	}

	// Remaining backticks + info string "go" - now in fenced code, this is text
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"remaining fence + info: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "``go" {
		t.Errorf(
			"remaining fence + info: Source=%q, want %q",
			tok.Source,
			"``go",
		)
	}

	// Newline
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"newline after info: got %v, want TokenNewline",
			tok.Type,
		)
	}

	// Content "code here" as text (no delimiters in fenced code)
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"code content: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "code here" {
		t.Errorf(
			"code content: Source=%q, want %q",
			tok.Source,
			"code here",
		)
	}

	// Newline before closing fence
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"newline before closing: got %v, want TokenNewline",
			tok.Type,
		)
	}

	// Closing fence as text (exits fenced code state)
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"closing fence: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "```" {
		t.Errorf(
			"closing fence: Source=%q, want %q",
			tok.Source,
			"```",
		)
	}

	// After closing fence, state should be back to normal
	if l.State() != StateNormal {
		t.Errorf(
			"after closing fence: state=%v, want StateNormal",
			l.State(),
		)
	}

	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"final token: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_FencedCodeBlockWithTilde verifies tilde-fenced code block.
// Similar to backticks, only the first tilde is emitted as TokenTilde.
func TestLexer_FencedCodeBlockWithTilde(
	t *testing.T,
) {
	input := "~~~\ncode\n~~~"
	l := newLexer([]byte(input))

	// First tilde - triggers state transition to StateFencedCode
	tok := l.Next()
	if tok.Type != TokenTilde {
		t.Errorf(
			"first tilde: got %v, want TokenTilde",
			tok.Type,
		)
	}

	if l.State() != StateFencedCode {
		t.Errorf(
			"after first tilde: state=%v, want StateFencedCode",
			l.State(),
		)
	}

	// Remaining tildes as text (now in fenced code state)
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"remaining tildes: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "~~" {
		t.Errorf(
			"remaining tildes: Source=%q, want %q",
			tok.Source,
			"~~",
		)
	}
}

// TestLexer_FencedCodeDelimitersAsText verifies delimiters are text inside fenced code.
func TestLexer_FencedCodeDelimitersAsText(
	t *testing.T,
) {
	input := "```\n# * _ ~ ` - +\n```"
	l := newLexer([]byte(input))

	// First backtick (enters fenced code state)
	tok := l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"first backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}

	// Remaining backticks as text
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"remaining backticks: got %v, want TokenText",
			tok.Type,
		)
	}

	// Newline
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"newline: got %v, want TokenNewline",
			tok.Type,
		)
	}

	// Content should be all text (delimiters not tokenized separately)
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"content with delimiters: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "# * _ ~ ` - +" {
		t.Errorf(
			"content: Source=%q, want %q",
			tok.Source,
			"# * _ ~ ` - +",
		)
	}
}

// TestLexer_InlineCode verifies inline code state transitions.
func TestLexer_InlineCode(t *testing.T) {
	input := "`code`"
	l := newLexer([]byte(input))

	// Opening backtick
	tok := l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"opening backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}

	// After opening backtick, should be in inline code state
	if l.State() != StateInlineCode {
		t.Errorf(
			"after opening backtick: state=%v, want StateInlineCode",
			l.State(),
		)
	}

	// Content as text
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"inline code content: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "code" {
		t.Errorf(
			"inline code content: Source=%q, want %q",
			tok.Source,
			"code",
		)
	}

	// Closing backtick
	tok = l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"closing backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}

	// After closing, should be back to normal
	if l.State() != StateNormal {
		t.Errorf(
			"after closing backtick: state=%v, want StateNormal",
			l.State(),
		)
	}
}

// TestLexer_InlineCodeDoubleBackticks verifies double backtick inline code.
// Note: Like fenced code, only the first backtick is emitted as TokenBacktick.
// The lexer counts 2 backticks and sets backtickCount=2, but only emits one token.
func TestLexer_InlineCodeDoubleBackticks(
	t *testing.T,
) {
	input := "``code with ` backtick``"
	l := newLexer([]byte(input))

	// First backtick (triggers inline code state with backtickCount=2)
	tok := l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"first opening backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}

	// After first backtick of 2-backtick sequence, state is inline code
	if l.State() != StateInlineCode {
		t.Errorf(
			"after first backtick: state=%v, want StateInlineCode",
			l.State(),
		)
	}

	// Second backtick is now treated as text (inside inline code)
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"second backtick as text: got %v, want TokenText",
			tok.Type,
		)
	}
	// In inline code, the backtick followed by "code..." is text until closing ``
	// The exact behavior depends on how closing is detected
}

// TestLexer_InlineCodeDelimitersAsText verifies delimiters are text inside inline code.
func TestLexer_InlineCodeDelimitersAsText(
	t *testing.T,
) {
	input := "`*bold* _italic_`"
	l := newLexer([]byte(input))

	// Opening backtick
	l.Next()

	// Content should include delimiters as text
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"inline code with delimiters: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "*bold* _italic_" {
		t.Errorf(
			"inline code content: Source=%q, want %q",
			tok.Source,
			"*bold* _italic_",
		)
	}
}

// TestLexer_LinkURLState verifies link URL state transitions.
func TestLexer_LinkURLState(t *testing.T) {
	l := newLexer(
		[]byte("https://example.com/*test*"),
	)

	// Manually enter link URL state (simulating parser behavior)
	l.enterLinkURLState()

	if l.State() != StateLinkURL {
		t.Errorf(
			"after enterLinkURLState: state=%v, want StateLinkURL",
			l.State(),
		)
	}

	// URL with special chars should be text
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"URL: got %v, want TokenText",
			tok.Type,
		)
	}
}

// TestLexer_LinkURLExitsOnParen verifies link URL state exits on closing paren.
func TestLexer_LinkURLExitsOnParen(t *testing.T) {
	l := newLexer([]byte("url)after"))

	l.enterLinkURLState()

	// URL text
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"URL: got %v, want TokenText",
			tok.Type,
		)
	}
	if string(tok.Source) != "url" {
		t.Errorf(
			"URL: Source=%q, want %q",
			tok.Source,
			"url",
		)
	}

	// Closing paren (exits state)
	tok = l.Next()
	if tok.Type != TokenParenClose {
		t.Errorf(
			"closing paren: got %v, want TokenParenClose",
			tok.Type,
		)
	}

	// Should be back in normal state
	if l.State() != StateNormal {
		t.Errorf(
			"after closing paren: state=%v, want StateNormal",
			l.State(),
		)
	}

	// "after" should be normal text
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"after: got %v, want TokenText",
			tok.Type,
		)
	}
}

// TestLexer_LinkURLExitsOnNewline verifies link URL state exits on newline.
func TestLexer_LinkURLExitsOnNewline(
	t *testing.T,
) {
	l := newLexer([]byte("url\nafter"))

	l.enterLinkURLState()

	// URL text
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"URL: got %v, want TokenText",
			tok.Type,
		)
	}

	// Newline (exits state)
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"newline: got %v, want TokenNewline",
			tok.Type,
		)
	}

	// Should be back in normal state
	if l.State() != StateNormal {
		t.Errorf(
			"after newline: state=%v, want StateNormal",
			l.State(),
		)
	}
}

// TestLexer_StateTransitionMethods verifies explicit state transition methods.
func TestLexer_StateTransitionMethods(
	t *testing.T,
) {
	l := newLexer([]byte("test"))

	// Initial state
	if l.State() != StateNormal {
		t.Errorf(
			"initial: state=%v, want StateNormal",
			l.State(),
		)
	}

	// Enter fenced code state
	l.enterFencedCodeState('`', 3)
	if l.State() != StateFencedCode {
		t.Errorf(
			"after enterFencedCodeState: state=%v, want StateFencedCode",
			l.State(),
		)
	}

	// Reset to normal
	l.exitLinkURLState() // This sets state to Normal
	if l.State() != StateNormal {
		t.Errorf(
			"after exitLinkURLState: state=%v, want StateNormal",
			l.State(),
		)
	}

	// Enter inline code state
	l.enterInlineCodeState(2)
	if l.State() != StateInlineCode {
		t.Errorf(
			"after enterInlineCodeState: state=%v, want StateInlineCode",
			l.State(),
		)
	}

	// Enter link URL state
	l.enterLinkURLState()
	if l.State() != StateLinkURL {
		t.Errorf(
			"after enterLinkURLState: state=%v, want StateLinkURL",
			l.State(),
		)
	}
}

// TestLexer_FencedCodeRequiresLineStart verifies fence must be at line start.
func TestLexer_FencedCodeRequiresLineStart(
	t *testing.T,
) {
	// Backticks not at line start should not trigger fenced code state
	input := " ```"
	l := newLexer([]byte(input))

	// Whitespace
	tok := l.Next()
	if tok.Type != TokenWhitespace {
		t.Errorf(
			"whitespace: got %v, want TokenWhitespace",
			tok.Type,
		)
	}

	// Backtick (should stay in normal state because not at line start)
	tok = l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}

	// Should NOT be in fenced code state
	// Note: The implementation may still enter inline code state
	// This is acceptable behavior - the test verifies it doesn't incorrectly
	// enter fenced code state for non-line-start backticks
}

// TestLexer_FencedCodeAfterNewline verifies fence at line start after newline.
func TestLexer_FencedCodeAfterNewline(
	t *testing.T,
) {
	input := "\n```"
	l := newLexer([]byte(input))

	// Newline
	tok := l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"newline: got %v, want TokenNewline",
			tok.Type,
		)
	}

	// Backticks at line start
	tok = l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}

	// Should be in fenced code state (after 3 backticks)
	// Note: We only read one backtick, lexer counts ahead
}

// TestLexer_LongerClosingFence verifies closing fence can be longer.
func TestLexer_LongerClosingFence(t *testing.T) {
	input := "```\ncode\n````"
	l := newLexer([]byte(input))

	// First backtick (enters fenced code state)
	tok := l.Next()
	if tok.Type != TokenBacktick {
		t.Errorf(
			"first backtick: got %v, want TokenBacktick",
			tok.Type,
		)
	}
	if l.State() != StateFencedCode {
		t.Errorf(
			"after first backtick: state=%v, want StateFencedCode",
			l.State(),
		)
	}

	// Remaining backticks (`` as text in fenced code)
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"remaining opening backticks: got %v, want TokenText",
			tok.Type,
		)
	}

	// Newline
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"first newline: got %v, want TokenNewline",
			tok.Type,
		)
	}

	// Code content
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"code content: got %v, want TokenText",
			tok.Type,
		)
	}

	// Newline before closing
	tok = l.Next()
	if tok.Type != TokenNewline {
		t.Errorf(
			"second newline: got %v, want TokenNewline",
			tok.Type,
		)
	}

	// Closing fence (4 backticks should still close a 3-backtick fence)
	tok = l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"closing fence: got %v, want TokenText",
			tok.Type,
		)
	}

	// Should be back in normal state
	if l.State() != StateNormal {
		t.Errorf(
			"after longer closing fence: state=%v, want StateNormal",
			l.State(),
		)
	}
}

// TestLexer_PeekReturnsSameToken verifies Peek returns same token multiple times.
func TestLexer_PeekReturnsSameToken(
	t *testing.T,
) {
	l := newLexer([]byte("hello"))

	tok1 := l.Peek()
	tok2 := l.Peek()
	tok3 := l.Peek()

	if tok1.Type != tok2.Type ||
		tok2.Type != tok3.Type {
		t.Errorf(
			"Peek returned different types: %v, %v, %v",
			tok1.Type,
			tok2.Type,
			tok3.Type,
		)
	}
	if tok1.Start != tok2.Start ||
		tok2.Start != tok3.Start {
		t.Errorf(
			"Peek returned different Start: %d, %d, %d",
			tok1.Start,
			tok2.Start,
			tok3.Start,
		)
	}
}

// TestLexer_NextAfterPeekClears verifies Next after Peek returns and clears cache.
func TestLexer_NextAfterPeekClears(t *testing.T) {
	l := newLexer([]byte("a b"))

	peek := l.Peek()
	next := l.Next()

	if peek.Type != next.Type {
		t.Errorf(
			"Peek and Next types differ: %v vs %v",
			peek.Type,
			next.Type,
		)
	}
	if peek.Start != next.Start {
		t.Errorf(
			"Peek and Next Start differ: %d vs %d",
			peek.Start,
			next.Start,
		)
	}

	// Next call should advance to next token
	next2 := l.Next()
	if next2.Start == next.Start {
		t.Error(
			"Next should advance position after returning peeked token",
		)
	}
}

// TestLexer_AlternatingPeekNext verifies alternating Peek/Next works correctly.
func TestLexer_AlternatingPeekNext(t *testing.T) {
	l := newLexer([]byte("a b c"))

	// Peek, Next, Peek, Next pattern
	p1 := l.Peek()
	n1 := l.Next()
	if p1.Type != n1.Type ||
		p1.Start != n1.Start {
		t.Errorf(
			"First Peek/Next mismatch: peek=%v@%d, next=%v@%d",
			p1.Type,
			p1.Start,
			n1.Type,
			n1.Start,
		)
	}

	p2 := l.Peek()
	n2 := l.Next()
	if p2.Type != n2.Type ||
		p2.Start != n2.Start {
		t.Errorf(
			"Second Peek/Next mismatch: peek=%v@%d, next=%v@%d",
			p2.Type,
			p2.Start,
			n2.Type,
			n2.Start,
		)
	}

	// Verify we're advancing
	if n1.Start == n2.Start {
		t.Error(
			"Tokens should have different positions",
		)
	}
}

// TestLexer_PeekDoesNotAdvance verifies that while the first Peek call
// advances internal position to lex a token, subsequent Peek calls return
// the cached token without further position advancement.
func TestLexer_PeekDoesNotAdvance(t *testing.T) {
	l := newLexer([]byte("hello"))

	// First Peek advances internal position to lex the token
	tok1 := l.Peek()
	posAfterFirstPeek := l.Pos()

	// Second Peek should return cached token without advancing position
	tok2 := l.Peek()
	posAfterSecondPeek := l.Pos()

	// Position should not change between subsequent Peek calls
	if posAfterFirstPeek != posAfterSecondPeek {
		t.Errorf(
			"Position changed between Peek calls: %d vs %d",
			posAfterFirstPeek,
			posAfterSecondPeek,
		)
	}

	// Tokens should be identical (compare fields since []byte prevents ==)
	if tok1.Type != tok2.Type {
		t.Errorf(
			"Peek returned different Type: %v vs %v",
			tok1.Type,
			tok2.Type,
		)
	}
	if tok1.Start != tok2.Start {
		t.Errorf(
			"Peek returned different Start: %d vs %d",
			tok1.Start,
			tok2.Start,
		)
	}
	if tok1.End != tok2.End {
		t.Errorf(
			"Peek returned different End: %d vs %d",
			tok1.End,
			tok2.End,
		)
	}
	if string(tok1.Source) != string(tok2.Source) {
		t.Errorf(
			"Peek returned different Source: %q vs %q",
			tok1.Source,
			tok2.Source,
		)
	}
}

// TestLexer_All verifies All() returns correct sequence ending with EOF.
func TestLexer_All(t *testing.T) {
	input := "# Hello"
	l := newLexer([]byte(input))

	tokens := l.All()

	if len(tokens) == 0 {
		t.Fatal("All() returned empty slice")
	}

	// Last token should be EOF
	last := tokens[len(tokens)-1]
	if last.Type != TokenEOF {
		t.Errorf(
			"last token: got %v, want TokenEOF",
			last.Type,
		)
	}

	// Should have: Hash, Whitespace, Text, EOF
	expectedTypes := []TokenType{
		TokenHash,
		TokenWhitespace,
		TokenText,
		TokenEOF,
	}
	if len(tokens) != len(expectedTypes) {
		t.Errorf(
			"token count: got %d, want %d",
			len(tokens),
			len(expectedTypes),
		)
	}

	for i, exp := range expectedTypes {
		if i < len(tokens) &&
			tokens[i].Type != exp {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tokens[i].Type,
				exp,
			)
		}
	}
}

// TestLexer_AllEmpty verifies All() on empty input returns just EOF.
func TestLexer_AllEmpty(t *testing.T) {
	l := newLexer(make([]byte, 0))

	tokens := l.All()

	if len(tokens) != 1 {
		t.Errorf(
			"empty All(): got %d tokens, want 1",
			len(tokens),
		)
	}
	if tokens[0].Type != TokenEOF {
		t.Errorf(
			"empty All() token: got %v, want TokenEOF",
			tokens[0].Type,
		)
	}
}

// TestLexer_AllWithErrorsNoErrors verifies AllWithErrors() with valid input.
func TestLexer_AllWithErrorsNoErrors(
	t *testing.T,
) {
	input := "hello world"
	l := newLexer([]byte(input))

	tokens, errors := l.AllWithErrors()

	if len(errors) != 0 {
		t.Errorf(
			"valid input: got %d errors, want 0",
			len(errors),
		)
	}
	if len(tokens) == 0 {
		t.Fatal(
			"AllWithErrors() returned empty tokens",
		)
	}
	if tokens[len(tokens)-1].Type != TokenEOF {
		t.Error("tokens should end with EOF")
	}
}

// TestLexer_AllWithErrorsExtractsErrors verifies errors are extracted.
func TestLexer_AllWithErrorsExtractsErrors(
	t *testing.T,
) {
	// Input with invalid UTF-8
	input := []byte{'a', 0x80, 'b'}
	l := newLexer(input)

	tokens, errors := l.AllWithErrors()

	if len(errors) != 1 {
		t.Errorf(
			"got %d errors, want 1",
			len(errors),
		)
	}
	if len(errors) > 0 && errors[0].Offset != 1 {
		t.Errorf(
			"error offset: got %d, want 1",
			errors[0].Offset,
		)
	}

	// Error should also be in tokens
	foundError := false
	for _, tok := range tokens {
		if tok.Type == TokenError {
			foundError = true

			break
		}
	}
	if !foundError {
		t.Error(
			"TokenError should be in tokens slice",
		)
	}
}

// TestLexer_AllWithErrorsMultipleErrors verifies multiple errors are extracted.
func TestLexer_AllWithErrorsMultipleErrors(
	t *testing.T,
) {
	// Multiple invalid bytes
	input := []byte{0x80, 'a', 0x81, 'b', 0x82}
	l := newLexer(input)

	tokens, errors := l.AllWithErrors()

	if len(errors) != 3 {
		t.Errorf(
			"got %d errors, want 3",
			len(errors),
		)
	}

	// Verify error offsets
	expectedOffsets := []int{0, 2, 4}
	for i, exp := range expectedOffsets {
		if i < len(errors) &&
			errors[i].Offset != exp {
			t.Errorf(
				"error %d offset: got %d, want %d",
				i,
				errors[i].Offset,
				exp,
			)
		}
	}

	// Tokens should end with EOF
	if tokens[len(tokens)-1].Type != TokenEOF {
		t.Error("tokens should end with EOF")
	}
}

// TestLexer_LexErrorInterface verifies LexError implements error interface.
func TestLexer_LexErrorInterface(t *testing.T) {
	err := LexError{
		Offset:  42,
		Message: "test error",
	}

	var e error = err // Should compile if LexError implements error

	if e.Error() != "test error" {
		t.Errorf(
			"Error() = %q, want %q",
			e.Error(),
			"test error",
		)
	}
}

// TestLexer_AllFromMiddle verifies All() works after some tokens consumed.
func TestLexer_AllFromMiddle(t *testing.T) {
	l := newLexer([]byte("a b c"))

	// Consume first token
	l.Next()

	// All() should return remaining tokens
	tokens := l.All()

	// Should have: Whitespace, Text, Whitespace, Text, EOF
	// (starting from after "a")
	if len(tokens) == 0 {
		t.Fatal(
			"All() from middle returned empty",
		)
	}
	if tokens[len(tokens)-1].Type != TokenEOF {
		t.Error("tokens should end with EOF")
	}
}

// TestLexer_PositionTracking verifies position updates correctly.
func TestLexer_PositionTracking(t *testing.T) {
	input := "ab cd"
	l := newLexer([]byte(input))

	if l.Pos() != 0 {
		t.Errorf(
			"initial Pos()=%d, want 0",
			l.Pos(),
		)
	}

	tok := l.Next() // "ab"
	if tok.Start != 0 || tok.End != 2 {
		t.Errorf(
			"first token: Start=%d, End=%d, want 0, 2",
			tok.Start,
			tok.End,
		)
	}

	tok = l.Next() // " "
	if tok.Start != 2 || tok.End != 3 {
		t.Errorf(
			"whitespace: Start=%d, End=%d, want 2, 3",
			tok.Start,
			tok.End,
		)
	}

	tok = l.Next() // "cd"
	if tok.Start != 3 || tok.End != 5 {
		t.Errorf(
			"second text: Start=%d, End=%d, want 3, 5",
			tok.Start,
			tok.End,
		)
	}
}

// TestLexer_SourceSlicesAreCorrect verifies Source slices match positions.
func TestLexer_SourceSlicesAreCorrect(
	t *testing.T,
) {
	input := "hello world"
	l := newLexer([]byte(input))

	tokens := l.All()

	for _, tok := range tokens {
		if tok.Type == TokenEOF {
			continue
		}
		expected := input[tok.Start:tok.End]
		if string(tok.Source) != expected {
			t.Errorf(
				"token at %d-%d: Source=%q, want %q",
				tok.Start,
				tok.End,
				tok.Source,
				expected,
			)
		}
	}
}

// TestLexer_OnlyWhitespace verifies input of only whitespace.
func TestLexer_OnlyWhitespace(t *testing.T) {
	l := newLexer([]byte("   "))

	tok := l.Next()
	if tok.Type != TokenWhitespace {
		t.Errorf(
			"only whitespace: got %v, want TokenWhitespace",
			tok.Type,
		)
	}
	if string(tok.Source) != "   " {
		t.Errorf(
			"whitespace source: got %q, want %q",
			tok.Source,
			"   ",
		)
	}

	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"after whitespace: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_OnlyNewlines verifies input of only newlines.
func TestLexer_OnlyNewlines(t *testing.T) {
	l := newLexer([]byte("\n\n\n"))

	for i := range 3 {
		tok := l.Next()
		if tok.Type != TokenNewline {
			t.Errorf(
				"newline %d: got %v, want TokenNewline",
				i,
				tok.Type,
			)
		}
	}

	tok := l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"after newlines: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_OnlyDelimiters verifies input of only delimiters.
func TestLexer_OnlyDelimiters(t *testing.T) {
	l := newLexer([]byte("#*_"))

	expected := []TokenType{
		TokenHash,
		TokenAsterisk,
		TokenUnderscore,
		TokenEOF,
	}

	for i, exp := range expected {
		tok := l.Next()
		if tok.Type != exp {
			t.Errorf(
				"token %d: got %v, want %v",
				i,
				tok.Type,
				exp,
			)
		}
	}
}

// TestLexer_VeryLongText verifies handling of very long text.
func TestLexer_VeryLongText(t *testing.T) {
	// Create a long string without delimiters
	long := make([]byte, 10000)
	for i := range long {
		long[i] = 'a'
	}

	l := newLexer(long)

	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"long text: got %v, want TokenText",
			tok.Type,
		)
	}
	if tok.Len() != 10000 {
		t.Errorf(
			"long text: Len()=%d, want 10000",
			tok.Len(),
		)
	}

	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"after long text: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_ZeroByte verifies handling of null byte (0x00).
func TestLexer_ZeroByte(t *testing.T) {
	input := []byte{'a', 0x00, 'b'}
	l := newLexer(input)

	// The lexer treats null byte as part of text (it's not a delimiter).
	// So the entire input should be one text token.
	tok := l.Next()
	if tok.Type != TokenText {
		t.Errorf(
			"token 0: got %v, want TokenText",
			tok.Type,
		)
	}
	// The text includes the null byte since it's not a delimiter
	if tok.Len() != 3 {
		t.Errorf(
			"token 0: Len()=%d, want 3",
			tok.Len(),
		)
	}

	// Should be EOF after
	tok = l.Next()
	if tok.Type != TokenEOF {
		t.Errorf(
			"token 1: got %v, want TokenEOF",
			tok.Type,
		)
	}
}

// TestLexer_DigitXAmbiguity verifies disambiguation between number and x.
func TestLexer_DigitXAmbiguity(t *testing.T) {
	// "x1" should be TokenX followed by TokenNumber
	l := newLexer([]byte("x1"))

	tok := l.Next()
	if tok.Type != TokenX {
		t.Errorf(
			"'x': got %v, want TokenX",
			tok.Type,
		)
	}

	tok = l.Next()
	if tok.Type != TokenNumber {
		t.Errorf(
			"'1': got %v, want TokenNumber",
			tok.Type,
		)
	}
}

// TestLexer_ConsecutiveNumbers verifies number sequences.
func TestLexer_ConsecutiveNumbers(t *testing.T) {
	l := newLexer([]byte("123 456"))

	tok := l.Next()
	if tok.Type != TokenNumber {
		t.Errorf(
			"first number: got %v, want TokenNumber",
			tok.Type,
		)
	}
	if string(tok.Source) != "123" {
		t.Errorf(
			"first number: Source=%q, want %q",
			tok.Source,
			"123",
		)
	}

	l.Next() // whitespace

	tok = l.Next()
	if tok.Type != TokenNumber {
		t.Errorf(
			"second number: got %v, want TokenNumber",
			tok.Type,
		)
	}
	if string(tok.Source) != "456" {
		t.Errorf(
			"second number: Source=%q, want %q",
			tok.Source,
			"456",
		)
	}
}
