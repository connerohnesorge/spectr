package parser

import (
	"testing"
)

// TestNewLexer verifies that a new lexer is initialized correctly.
func TestNewLexer(t *testing.T) {
	input := "# Test Header\nSome text"
	l := NewLexer(input)

	if l.input != input {
		t.Errorf("expected input %q, got %q", input, l.input)
	}

	if l.line != 1 {
		t.Errorf("expected line 1, got %d", l.line)
	}

	if l.column != 1 {
		t.Errorf("expected column 1, got %d", l.column)
	}

	if l.start != 0 || l.pos != 0 {
		t.Errorf("expected start=0 pos=0, got start=%d pos=%d", l.start, l.pos)
	}

	if len(l.tokens) != 0 {
		t.Errorf("expected empty tokens, got %d tokens", len(l.tokens))
	}
}

// TestLexerNext verifies basic rune traversal and position tracking.
func TestLexerNext(t *testing.T) {
	l := NewLexer("abc\nde")

	// Read 'a'
	r := l.next()
	if r != 'a' {
		t.Errorf("expected 'a', got %c", r)
	}
	if l.line != 1 || l.column != 2 {
		t.Errorf("after 'a': expected 1:2, got %d:%d", l.line, l.column)
	}

	// Read 'b'
	r = l.next()
	if r != 'b' {
		t.Errorf("expected 'b', got %c", r)
	}
	if l.line != 1 || l.column != 3 {
		t.Errorf("after 'b': expected 1:3, got %d:%d", l.line, l.column)
	}

	// Read 'c'
	r = l.next()
	if r != 'c' {
		t.Errorf("expected 'c', got %c", r)
	}

	// Read '\n' - should increment line
	r = l.next()
	if r != '\n' {
		t.Errorf("expected newline, got %c", r)
	}
	if l.line != 2 || l.column != 1 {
		t.Errorf("after newline: expected 2:1, got %d:%d", l.line, l.column)
	}

	// Read 'd'
	r = l.next()
	if r != 'd' {
		t.Errorf("expected 'd', got %c", r)
	}
	if l.line != 2 || l.column != 2 {
		t.Errorf("after 'd': expected 2:2, got %d:%d", l.line, l.column)
	}

	// Read 'e'
	r = l.next()
	if r != 'e' {
		t.Errorf("expected 'e', got %c", r)
	}

	// Read past end - should return 0
	r = l.next()
	if r != 0 {
		t.Errorf("expected EOF (0), got %c", r)
	}
}

// TestLexerPeek verifies peek doesn't advance position.
func TestLexerPeek(t *testing.T) {
	l := NewLexer("ab")

	r := l.peek()
	if r != 'a' {
		t.Errorf("expected peek 'a', got %c", r)
	}

	// Position should not have changed
	if l.pos != 0 {
		t.Errorf("peek should not advance position, got pos=%d", l.pos)
	}

	// Next should still read 'a'
	r = l.next()
	if r != 'a' {
		t.Errorf("after peek, next should read 'a', got %c", r)
	}
}

// TestLexerBackup verifies backup restores previous position.
func TestLexerBackup(t *testing.T) {
	l := NewLexer("abc\nd")

	// Read 'a'
	l.next()
	if l.pos != 1 || l.line != 1 || l.column != 2 {
		t.Fatalf("after 'a': expected pos=1, 1:2, got pos=%d, %d:%d",
			l.pos, l.line, l.column)
	}

	// Backup
	l.backup()
	if l.pos != 0 || l.line != 1 || l.column != 1 {
		t.Errorf("after backup: expected pos=0, 1:1, got pos=%d, %d:%d",
			l.pos, l.line, l.column)
	}

	// Read 'a' again
	r := l.next()
	if r != 'a' {
		t.Errorf("after backup, next should read 'a', got %c", r)
	}

	// Test backup over newline
	l.next() // 'b'
	l.next() // 'c'
	l.next() // '\n'
	if l.line != 2 || l.column != 1 {
		t.Fatalf("after newline: expected 2:1, got %d:%d", l.line, l.column)
	}

	l.backup()
	if l.line != 1 {
		t.Errorf("backup over newline: expected line 1, got %d", l.line)
	}
}

// TestLexerEmit verifies token emission.
func TestLexerEmit(t *testing.T) {
	l := NewLexer("hello")

	// Consume "hel"
	l.next() // h
	l.next() // e
	l.next() // l

	l.emit(TokenText)

	if len(l.tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(l.tokens))
	}

	token := l.tokens[0]
	if token.Type != TokenText {
		t.Errorf("expected TokenText, got %s", token.Type)
	}
	if token.Value != "hel" {
		t.Errorf("expected value 'hel', got %q", token.Value)
	}
	if token.Pos.Line != 1 || token.Pos.Column != 1 {
		t.Errorf("expected position 1:1, got %d:%d",
			token.Pos.Line, token.Pos.Column)
	}

	// start should now be at position 3
	if l.start != 3 {
		t.Errorf("after emit, start should be 3, got %d", l.start)
	}
}

// TestLexerAccept verifies accept consumes valid runes.
func TestLexerAccept(t *testing.T) {
	l := NewLexer("abc123")

	// Accept 'a' from letters
	if !l.accept("abc") {
		t.Error("accept should consume 'a'")
	}
	if l.pos != 1 {
		t.Errorf("after accept: expected pos=1, got %d", l.pos)
	}

	// Try to accept number (should fail)
	if l.accept("123") {
		t.Error("accept should not consume 'b' when looking for numbers")
	}
	if l.pos != 1 {
		t.Errorf("failed accept should not advance, got pos=%d", l.pos)
	}
}

// TestLexerAcceptRun verifies acceptRun consumes multiple runes.
func TestLexerAcceptRun(t *testing.T) {
	l := NewLexer("###text")

	count := l.acceptRun("#")
	if count != 3 {
		t.Errorf("expected 3 '#' consumed, got %d", count)
	}
	if l.pos != 3 {
		t.Errorf("expected pos=3, got %d", l.pos)
	}

	// Next character should be 't'
	r := l.next()
	if r != 't' {
		t.Errorf("after acceptRun, expected 't', got %c", r)
	}
}

// TestLexerAtLineStart verifies line start detection.
func TestLexerAtLineStart(t *testing.T) {
	l := NewLexer("ab\ncd")

	if !l.atLineStart() {
		t.Error("should be at line start initially")
	}

	l.next() // 'a'
	if l.atLineStart() {
		t.Error("should not be at line start after reading 'a'")
	}

	l.next() // 'b'
	l.next() // '\n'
	if !l.atLineStart() {
		t.Error("should be at line start after newline")
	}

	l.next() // 'c'
	if l.atLineStart() {
		t.Error("should not be at line start after reading 'c'")
	}
}

// TestTokenTypeString verifies TokenType string representation.
func TestTokenTypeString(t *testing.T) {
	tests := []struct {
		typ      TokenType
		expected string
	}{
		{TokenEOF, "EOF"},
		{TokenText, "Text"},
		{TokenHeader, "Header"},
		{TokenCodeBlock, "CodeBlock"},
		{TokenListItem, "ListItem"},
		{TokenBlankLine, "BlankLine"},
		{TokenError, "Error"},
	}

	for _, tt := range tests {
		if got := tt.typ.String(); got != tt.expected {
			t.Errorf("TokenType(%d).String() = %q, want %q",
				tt.typ, got, tt.expected)
		}
	}
}

// TestPositionString verifies Position string representation.
func TestPositionString(t *testing.T) {
	pos := Position{Line: 42, Column: 17, Offset: 512}
	expected := "42:17"
	if got := pos.String(); got != expected {
		t.Errorf("Position.String() = %q, want %q", got, expected)
	}
}

// TestTokenString verifies Token string representation.
func TestTokenString(t *testing.T) {
	// Short token
	token := Token{
		Type:  TokenText,
		Value: "short",
		Pos:   Position{Line: 1, Column: 5, Offset: 4},
	}
	str := token.String()
	if str != "Text@1:5: short" {
		t.Errorf("Token.String() = %q, want %q", str, "Text@1:5: short")
	}

	// Long token (should truncate)
	longValue := "this is a very long token value that should be truncated"
	token = Token{
		Type:  TokenHeader,
		Value: longValue,
		Pos:   Position{Line: 2, Column: 1, Offset: 10},
	}
	str = token.String()
	if len(str) > 50 && !contains(str, "...") {
		t.Errorf("long Token.String() should truncate with '...': %q", str)
	}
}

// contains checks if s contains substr (helper for tests).
func contains(s, substr string) bool {
	return len(s) >= len(substr) && indexOf(s, substr) >= 0
}

// indexOf returns the index of substr in s, or -1 if not found.
func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}

// TestLexerCurrentValue verifies currentValue returns the token being built.
func TestLexerCurrentValue(t *testing.T) {
	l := NewLexer("hello world")

	l.next() // h
	l.next() // e
	l.next() // l

	if got := l.currentValue(); got != "hel" {
		t.Errorf("currentValue() = %q, want %q", got, "hel")
	}

	l.next() // l
	l.next() // o

	if got := l.currentValue(); got != "hello" {
		t.Errorf("currentValue() = %q, want %q", got, "hello")
	}
}

// TestLexerIgnore verifies ignore advances start without emitting.
func TestLexerIgnore(t *testing.T) {
	l := NewLexer("   text")

	// Consume spaces
	l.next()
	l.next()
	l.next()

	// Ignore them
	l.ignore()

	if l.start != 3 {
		t.Errorf("after ignore, start should be 3, got %d", l.start)
	}

	if len(l.tokens) != 0 {
		t.Errorf("ignore should not emit tokens, got %d", len(l.tokens))
	}
}

// TestLexerEmitError verifies error token emission.
func TestLexerEmitError(t *testing.T) {
	l := NewLexer("test")

	l.next() // advance to position 1:2

	l.emitError("test error message")

	if len(l.tokens) != 1 {
		t.Fatalf("expected 1 token, got %d", len(l.tokens))
	}

	token := l.tokens[0]
	if token.Type != TokenError {
		t.Errorf("expected TokenError, got %s", token.Type)
	}
	if token.Value != "test error message" {
		t.Errorf("expected error message, got %q", token.Value)
	}
}

// TestLexerPeekString verifies peekString returns next n characters.
func TestLexerPeekString(t *testing.T) {
	l := NewLexer("```markdown")

	// Peek first 3 characters
	s := l.peekString(3)
	if s != "```" {
		t.Errorf("expected '```', got %q", s)
	}

	// Position should not change
	if l.pos != 0 {
		t.Errorf("peekString should not advance position, got pos=%d", l.pos)
	}

	// Peek beyond end
	l.pos = 10
	s = l.peekString(5)
	if s != "" {
		t.Errorf("peekString beyond end should return empty string, got %q", s)
	}
}

// TestLexText verifies basic text lexing.
func TestLexText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []TokenType
	}{
		{
			name:     "plain text",
			input:    "just some text",
			expected: []TokenType{TokenText, TokenEOF},
		},
		{
			name:     "text with newline",
			input:    "line one\nline two",
			expected: []TokenType{TokenText, TokenEOF},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []TokenType{TokenEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := l.Lex()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i] {
					t.Errorf("token %d: expected %s, got %s",
						i, tt.expected[i], tok.Type)
				}
			}
		})
	}
}

// TestLexHeader verifies header lexing.
func TestLexHeader(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedType  TokenType
		expectedValue string
	}{
		{
			name:          "single hash",
			input:         "# Header",
			expectedType:  TokenHeader,
			expectedValue: "# Header",
		},
		{
			name:          "double hash",
			input:         "## Header",
			expectedType:  TokenHeader,
			expectedValue: "## Header",
		},
		{
			name:          "triple hash",
			input:         "### Header",
			expectedType:  TokenHeader,
			expectedValue: "### Header",
		},
		{
			name:          "quad hash",
			input:         "#### Header",
			expectedType:  TokenHeader,
			expectedValue: "#### Header",
		},
		{
			name:          "header with no space",
			input:         "#Header",
			expectedType:  TokenHeader,
			expectedValue: "#Header",
		},
		{
			name:          "header with trailing space",
			input:         "# Header ",
			expectedType:  TokenHeader,
			expectedValue: "# Header ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := l.Lex()

			// Should have header token and EOF
			if len(tokens) < 2 {
				t.Fatalf("expected at least 2 tokens, got %d", len(tokens))
			}

			if tokens[0].Type != tt.expectedType {
				t.Errorf("expected %s, got %s", tt.expectedType, tokens[0].Type)
			}

			if tokens[0].Value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, tokens[0].Value)
			}
		})
	}
}

// TestLexCodeBlock verifies code block lexing.
func TestLexCodeBlock(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple code block",
			input:    "```\ncode\n```",
			expected: "```\ncode\n```",
		},
		{
			name:     "code block with language",
			input:    "```go\nfunc main() {}\n```",
			expected: "```go\nfunc main() {}\n```",
		},
		{
			name:     "code block with markdown inside",
			input:    "```\n# Not a header\n## Also not a header\n- Not a list\n```",
			expected: "```\n# Not a header\n## Also not a header\n- Not a list\n```",
		},
		{
			name:     "critical: requirement in code block",
			input:    "```go\n### Requirement: This is in a code block\n```\n### Requirement: This is real",
			expected: "```go\n### Requirement: This is in a code block\n```",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := l.Lex()

			// Find the code block token
			var codeBlockToken *Token
			for i := range tokens {
				if tokens[i].Type == TokenCodeBlock {
					codeBlockToken = &tokens[i]

					break
				}
			}

			if codeBlockToken == nil {
				t.Fatal("expected to find TokenCodeBlock")
			}

			if codeBlockToken.Value != tt.expected {
				t.Errorf("expected value %q, got %q", tt.expected, codeBlockToken.Value)
			}
		})
	}
}

// TestLexCodeBlockCritical verifies the critical use case.
//
// This is the most important test: it ensures that markdown syntax
// inside code blocks is NOT tokenized as headers.
func TestLexCodeBlockCritical(t *testing.T) {
	input := "```go\n### Requirement: This is inside a code block\n```\n### Requirement: This is a real requirement"

	l := NewLexer(input)
	tokens := l.Lex()

	// Verify token sequence
	expectedTypes := []TokenType{
		TokenCodeBlock, // The code block containing fake requirement
		TokenHeader,    // The real requirement
		TokenEOF,
	}

	if len(tokens) != len(expectedTypes) {
		t.Fatalf("expected %d tokens, got %d", len(expectedTypes), len(tokens))
	}

	for i, expectedType := range expectedTypes {
		if tokens[i].Type != expectedType {
			t.Errorf("token %d: expected %s, got %s (value: %q)",
				i, expectedType, tokens[i].Type, tokens[i].Value)
		}
	}

	// Verify the code block contains the fake requirement
	codeBlockValue := tokens[0].Value
	if !contains(codeBlockValue, "### Requirement: This is inside a code block") {
		t.Errorf("code block should contain fake requirement, got: %q", codeBlockValue)
	}

	// Verify the header is the real requirement
	headerValue := tokens[1].Value
	if !contains(headerValue, "### Requirement: This is a real requirement") {
		t.Errorf("header should be real requirement, got: %q", headerValue)
	}
}

// TestLexList verifies list item lexing.
func TestLexList(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedType  TokenType
		expectedValue string
	}{
		{
			name:          "dash list",
			input:         "- Item one",
			expectedType:  TokenListItem,
			expectedValue: "- Item one",
		},
		{
			name:          "asterisk list",
			input:         "* Item two",
			expectedType:  TokenListItem,
			expectedValue: "* Item two",
		},
		{
			name:          "list with tab",
			input:         "-\tTabbed item",
			expectedType:  TokenListItem,
			expectedValue: "-\tTabbed item",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := l.Lex()

			if len(tokens) < 2 {
				t.Fatalf("expected at least 2 tokens, got %d", len(tokens))
			}

			if tokens[0].Type != tt.expectedType {
				t.Errorf("expected %s, got %s", tt.expectedType, tokens[0].Type)
			}

			if tokens[0].Value != tt.expectedValue {
				t.Errorf("expected value %q, got %q", tt.expectedValue, tokens[0].Value)
			}
		})
	}
}

// TestLexBlankLine verifies blank line lexing.
func TestLexBlankLine(t *testing.T) {
	tests := []struct {
		name  string
		input string
		count int // number of TokenBlankLine expected
	}{
		{
			name:  "single blank line",
			input: "\n",
			count: 1,
		},
		{
			name:  "multiple blank lines",
			input: "\n\n\n",
			count: 1, // Should be combined into one
		},
		{
			name:  "blank lines between text",
			input: "text\n\nmore text",
			count: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewLexer(tt.input)
			tokens := l.Lex()

			blankLineCount := 0
			for _, tok := range tokens {
				if tok.Type == TokenBlankLine {
					blankLineCount++
				}
			}

			if blankLineCount != tt.count {
				t.Errorf("expected %d blank line tokens, got %d", tt.count, blankLineCount)
			}
		})
	}
}

// TestLexMixedContent verifies lexing of mixed markdown content.
func TestLexMixedContent(t *testing.T) {
	input := `# Header

Some text here.

## Subheader

- List item 1
- List item 2

` + "```go\ncode\n```" + `

More text.`

	l := NewLexer(input)
	tokens := l.Lex()

	// Verify we got diverse token types
	types := make(map[TokenType]bool)
	for _, tok := range tokens {
		types[tok.Type] = true
	}

	expectedTypes := []TokenType{
		TokenHeader,
		TokenText,
		TokenBlankLine,
		TokenListItem,
		TokenCodeBlock,
		TokenEOF,
	}

	for _, expectedType := range expectedTypes {
		if !types[expectedType] {
			t.Errorf("expected to find %s token", expectedType)
		}
	}
}

// TestLexStateTransitions verifies correct transitions between states.
func TestLexStateTransitions(t *testing.T) {
	input := "text\n# header\nmore text\n```\ncode\n```\nfinal text"

	l := NewLexer(input)
	tokens := l.Lex()

	expectedSequence := []TokenType{
		TokenText,      // "text\n"
		TokenHeader,    // "# header"
		TokenText,      // "\nmore text\n"
		TokenCodeBlock, // "```\ncode\n```"
		TokenText,      // "\nfinal text"
		TokenEOF,
	}

	if len(tokens) != len(expectedSequence) {
		t.Fatalf("expected %d tokens, got %d", len(expectedSequence), len(tokens))
	}

	for i, expected := range expectedSequence {
		if tokens[i].Type != expected {
			t.Errorf("token %d: expected %s, got %s (value: %q)",
				i, expected, tokens[i].Type, tokens[i].Value)
		}
	}
}

// TestLexPositionTracking verifies accurate line and column tracking.
func TestLexPositionTracking(t *testing.T) {
	input := "line1\n# Header\nline3"

	l := NewLexer(input)
	tokens := l.Lex()

	// First token (text) should be at 1:1
	if tokens[0].Pos.Line != 1 || tokens[0].Pos.Column != 1 {
		t.Errorf("first token: expected 1:1, got %d:%d",
			tokens[0].Pos.Line, tokens[0].Pos.Column)
	}

	// Header should be at line 2
	var headerToken *Token
	for i := range tokens {
		if tokens[i].Type == TokenHeader {
			headerToken = &tokens[i]

			break
		}
	}

	if headerToken == nil {
		t.Fatal("expected to find header token")
	}

	if headerToken.Pos.Line != 2 {
		t.Errorf("header: expected line 2, got %d", headerToken.Pos.Line)
	}
}

// TestLexUnclosedCodeBlock verifies handling of unclosed code blocks.
func TestLexUnclosedCodeBlock(t *testing.T) {
	input := "```\ncode without closing fence"

	l := NewLexer(input)
	tokens := l.Lex()

	// Should get code block token and error token
	hasCodeBlock := false
	hasError := false

	for _, tok := range tokens {
		if tok.Type == TokenCodeBlock {
			hasCodeBlock = true
		}
		if tok.Type == TokenError {
			hasError = true
		}
	}

	if !hasCodeBlock {
		t.Error("expected TokenCodeBlock for unclosed fence")
	}

	if !hasError {
		t.Error("expected TokenError for unclosed fence")
	}
}

// TestLexNotListItem verifies that - without space is not a list.
func TestLexNotListItem(t *testing.T) {
	input := "-notalist"

	l := NewLexer(input)
	tokens := l.Lex()

	// Should be treated as text, not list item
	if len(tokens) < 1 {
		t.Fatal("expected at least 1 token")
	}

	if tokens[0].Type != TokenText {
		t.Errorf("expected TokenText, got %s", tokens[0].Type)
	}
}

// TestLexHeaderNotAtLineStart verifies # in middle of line is not header.
func TestLexHeaderNotAtLineStart(t *testing.T) {
	input := "This has a # in the middle"

	l := NewLexer(input)
	tokens := l.Lex()

	// Should be all text, no header
	for _, tok := range tokens {
		if tok.Type == TokenHeader {
			t.Error("# in middle of line should not be a header")
		}
	}
}

// TestLexMultipleHeaders verifies multiple headers are tokenized correctly.
func TestLexMultipleHeaders(t *testing.T) {
	input := "# First\n## Second\n### Third\n#### Fourth"

	l := NewLexer(input)
	tokens := l.Lex()

	headerCount := 0
	for _, tok := range tokens {
		if tok.Type == TokenHeader {
			headerCount++
		}
	}

	if headerCount != 4 {
		t.Errorf("expected 4 headers, got %d", headerCount)
	}
}
