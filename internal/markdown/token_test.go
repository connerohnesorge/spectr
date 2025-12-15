package markdown

import (
	"testing"
)

// TestTokenType_String verifies that all TokenType values have meaningful String() output.
func TestTokenType_String(t *testing.T) {
	tests := []struct {
		tokenType TokenType
		expected  string
	}{
		// Structural tokens
		{TokenEOF, "EOF"},
		{TokenNewline, "Newline"},
		{TokenWhitespace, "Whitespace"},
		{TokenText, "Text"},
		{TokenError, "Error"},

		// Punctuation delimiter tokens
		{TokenHash, "Hash"},
		{TokenAsterisk, "Asterisk"},
		{TokenUnderscore, "Underscore"},
		{TokenTilde, "Tilde"},
		{TokenBacktick, "Backtick"},
		{TokenDash, "Dash"},
		{TokenPlus, "Plus"},
		{TokenDot, "Dot"},
		{TokenColon, "Colon"},
		{TokenPipe, "Pipe"},

		// Bracket tokens
		{TokenBracketOpen, "BracketOpen"},
		{TokenBracketClose, "BracketClose"},
		{TokenParenOpen, "ParenOpen"},
		{TokenParenClose, "ParenClose"},
		{TokenGreaterThan, "GreaterThan"},

		// Special tokens
		{TokenNumber, "Number"},
		{TokenX, "X"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			got := tt.tokenType.String()
			if got != tt.expected {
				t.Errorf(
					"TokenType(%d).String() = %q, want %q",
					tt.tokenType,
					got,
					tt.expected,
				)
			}
		})
	}
}

// TestTokenType_String_NotEmpty verifies all token types return non-empty strings.
func TestTokenType_String_NotEmpty(t *testing.T) {
	// Test all defined token types (0 through TokenX)
	for i := TokenType(0); i <= TokenX; i++ {
		s := i.String()
		if s == "" {
			t.Errorf(
				"TokenType(%d).String() returned empty string",
				i,
			)
		}
		if s == "Unknown" {
			t.Errorf(
				"TokenType(%d).String() returned 'Unknown' for a defined token type",
				i,
			)
		}
	}
}

// TestTokenType_String_Unknown verifies that undefined token types return "Unknown".
func TestTokenType_String_Unknown(t *testing.T) {
	// Test an undefined token type (beyond the last defined one)
	undefinedType := TokenX + 1
	got := undefinedType.String()
	if got != "Unknown" {
		t.Errorf(
			"TokenType(%d).String() = %q, want %q",
			undefinedType,
			got,
			"Unknown",
		)
	}

	// Test a much larger undefined value
	largeType := TokenType(255)
	got = largeType.String()
	if got != "Unknown" {
		t.Errorf(
			"TokenType(%d).String() = %q, want %q",
			largeType,
			got,
			"Unknown",
		)
	}
}

// TestToken_FieldAccess verifies Token struct fields are accessible.
func TestToken_FieldAccess(t *testing.T) {
	source := []byte("hello world")
	tok := Token{
		Type:    TokenText,
		Start:   0,
		End:     5,
		Source:  source[:5],
		Message: "",
	}

	if tok.Type != TokenText {
		t.Errorf(
			"Token.Type = %v, want TokenText",
			tok.Type,
		)
	}
	if tok.Start != 0 {
		t.Errorf(
			"Token.Start = %d, want 0",
			tok.Start,
		)
	}
	if tok.End != 5 {
		t.Errorf(
			"Token.End = %d, want 5",
			tok.End,
		)
	}
	if string(tok.Source) != "hello" {
		t.Errorf(
			"Token.Source = %q, want %q",
			tok.Source,
			"hello",
		)
	}
	if tok.Message != "" {
		t.Errorf(
			"Token.Message = %q, want empty string",
			tok.Message,
		)
	}
}

// TestToken_Len verifies the Len() method.
func TestToken_Len(t *testing.T) {
	tests := []struct {
		name     string
		start    int
		end      int
		expected int
	}{
		{"empty token", 0, 0, 0},
		{"single char", 0, 1, 1},
		{"multi char", 0, 5, 5},
		{"offset token", 10, 15, 5},
		{"large token", 100, 200, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Token{
				Start: tt.start,
				End:   tt.end,
			}
			if got := tok.Len(); got != tt.expected {
				t.Errorf(
					"Token{Start: %d, End: %d}.Len() = %d, want %d",
					tt.start,
					tt.end,
					got,
					tt.expected,
				)
			}
		})
	}
}

// TestToken_Text verifies the Text() method.
func TestToken_Text(t *testing.T) {
	tests := []struct {
		name     string
		source   []byte
		expected string
	}{
		{"simple text", []byte("hello"), "hello"},
		{"empty source", make([]byte, 0), ""},
		{
			"unicode text",
			[]byte("hello, world"),
			"hello, world",
		},
		{
			"special chars",
			[]byte("**bold**"),
			"**bold**",
		},
		{"newline", []byte("\n"), "\n"},
		{"whitespace", []byte("   "), "   "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Token{Source: tt.source}
			if got := tok.Text(); got != tt.expected {
				t.Errorf(
					"Token{Source: %q}.Text() = %q, want %q",
					tt.source,
					got,
					tt.expected,
				)
			}
		})
	}
}

// TestToken_Text_NilSource verifies Text() handles nil Source gracefully.
func TestToken_Text_NilSource(t *testing.T) {
	tok := Token{
		Type:   TokenText,
		Start:  0,
		End:    5,
		Source: nil,
	}

	got := tok.Text()
	if got != "" {
		t.Errorf(
			"Token with nil Source.Text() = %q, want empty string",
			got,
		)
	}
}

// TestToken_IsDelimiter verifies the IsDelimiter() method.
func TestToken_IsDelimiter(t *testing.T) {
	delimiterTypes := []TokenType{
		TokenHash, TokenAsterisk, TokenUnderscore, TokenTilde,
		TokenBacktick, TokenDash, TokenPlus, TokenDot, TokenColon,
		TokenPipe, TokenBracketOpen, TokenBracketClose,
		TokenParenOpen, TokenParenClose, TokenGreaterThan,
	}

	nonDelimiterTypes := []TokenType{
		TokenEOF, TokenNewline, TokenWhitespace, TokenText, TokenError,
		TokenNumber, TokenX,
	}

	for _, tt := range delimiterTypes {
		t.Run(
			tt.String()+"_is_delimiter",
			func(t *testing.T) {
				tok := Token{Type: tt}
				if !tok.IsDelimiter() {
					t.Errorf(
						"Token{Type: %v}.IsDelimiter() = false, want true",
						tt,
					)
				}
			},
		)
	}

	for _, tt := range nonDelimiterTypes {
		t.Run(
			tt.String()+"_not_delimiter",
			func(t *testing.T) {
				tok := Token{Type: tt}
				if tok.IsDelimiter() {
					t.Errorf(
						"Token{Type: %v}.IsDelimiter() = true, want false",
						tt,
					)
				}
			},
		)
	}
}

// TestToken_IsBracket verifies the IsBracket() method.
func TestToken_IsBracket(t *testing.T) {
	bracketTypes := []TokenType{
		TokenBracketOpen, TokenBracketClose, TokenParenOpen, TokenParenClose,
	}

	nonBracketTypes := []TokenType{
		TokenEOF, TokenNewline, TokenWhitespace, TokenText, TokenError,
		TokenHash, TokenAsterisk, TokenUnderscore, TokenTilde,
		TokenBacktick, TokenDash, TokenPlus, TokenDot, TokenColon,
		TokenPipe, TokenGreaterThan, TokenNumber, TokenX,
	}

	for _, tt := range bracketTypes {
		t.Run(
			tt.String()+"_is_bracket",
			func(t *testing.T) {
				tok := Token{Type: tt}
				if !tok.IsBracket() {
					t.Errorf(
						"Token{Type: %v}.IsBracket() = false, want true",
						tt,
					)
				}
			},
		)
	}

	for _, tt := range nonBracketTypes {
		t.Run(
			tt.String()+"_not_bracket",
			func(t *testing.T) {
				tok := Token{Type: tt}
				if tok.IsBracket() {
					t.Errorf(
						"Token{Type: %v}.IsBracket() = true, want false",
						tt,
					)
				}
			},
		)
	}
}

// TestToken_IsStructural verifies the IsStructural() method.
func TestToken_IsStructural(t *testing.T) {
	structuralTypes := []TokenType{
		TokenEOF, TokenNewline, TokenWhitespace, TokenText, TokenError,
	}

	nonStructuralTypes := []TokenType{
		TokenHash, TokenAsterisk, TokenUnderscore, TokenTilde,
		TokenBacktick, TokenDash, TokenPlus, TokenDot, TokenColon,
		TokenPipe, TokenBracketOpen, TokenBracketClose,
		TokenParenOpen, TokenParenClose, TokenGreaterThan,
		TokenNumber, TokenX,
	}

	for _, tt := range structuralTypes {
		t.Run(
			tt.String()+"_is_structural",
			func(t *testing.T) {
				tok := Token{Type: tt}
				if !tok.IsStructural() {
					t.Errorf(
						"Token{Type: %v}.IsStructural() = false, want true",
						tt,
					)
				}
			},
		)
	}

	for _, tt := range nonStructuralTypes {
		t.Run(
			tt.String()+"_not_structural",
			func(t *testing.T) {
				tok := Token{Type: tt}
				if tok.IsStructural() {
					t.Errorf(
						"Token{Type: %v}.IsStructural() = true, want false",
						tt,
					)
				}
			},
		)
	}
}

// TestToken_ErrorMessage verifies TokenError's Message field usage.
func TestToken_ErrorMessage(t *testing.T) {
	tests := []struct {
		name    string
		message string
	}{
		{"simple error", "unexpected character"},
		{"empty message", ""},
		{
			"detailed error",
			"invalid UTF-8 sequence at byte offset 42",
		},
		{
			"unicode message",
			"unexpected character: '!'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tok := Token{
				Type:    TokenError,
				Start:   0,
				End:     1,
				Source:  []byte("?"),
				Message: tt.message,
			}

			if tok.Message != tt.message {
				t.Errorf(
					"Token.Message = %q, want %q",
					tok.Message,
					tt.message,
				)
			}
			if tok.Type != TokenError {
				t.Errorf(
					"Token.Type = %v, want TokenError",
					tok.Type,
				)
			}
		})
	}
}

// TestToken_ZeroCopySource verifies that Source provides zero-copy access.
func TestToken_ZeroCopySource(t *testing.T) {
	original := []byte("hello world")
	tok := Token{
		Type:   TokenText,
		Start:  0,
		End:    5,
		Source: original[:5],
	}

	// Verify Source points to the same underlying memory
	// by checking that modifications to original affect Source
	if &tok.Source[0] != &original[0] {
		t.Error(
			"Token.Source does not share memory with original slice",
		)
	}

	// Verify Text() creates a copy (modifications don't affect it)
	text := tok.Text()
	original[0] = 'H'

	if tok.Source[0] != 'H' {
		t.Error(
			"Token.Source should reflect changes to original",
		)
	}
	if text != "hello" {
		t.Error(
			"Token.Text() should return a copy, not be affected by original changes",
		)
	}
}

// TestToken_EOF verifies EOF token characteristics.
func TestToken_EOF(t *testing.T) {
	sourceLen := 100
	tok := Token{
		Type:   TokenEOF,
		Start:  sourceLen,
		End:    sourceLen,
		Source: nil,
	}

	if tok.Len() != 0 {
		t.Errorf(
			"EOF token Len() = %d, want 0",
			tok.Len(),
		)
	}
	if tok.Text() != "" {
		t.Errorf(
			"EOF token Text() = %q, want empty string",
			tok.Text(),
		)
	}
	if !tok.IsStructural() {
		t.Error("EOF token should be structural")
	}
	if tok.IsDelimiter() {
		t.Error(
			"EOF token should not be a delimiter",
		)
	}
	if tok.IsBracket() {
		t.Error(
			"EOF token should not be a bracket",
		)
	}
}

// TestToken_Completeness verifies all token types are covered by helper methods.
func TestToken_Completeness(t *testing.T) {
	// Every token type should be classified by exactly one of:
	// - IsStructural
	// - IsDelimiter (and some are also IsBracket)
	// - Special (TokenNumber, TokenX)

	allTypes := []TokenType{
		TokenEOF, TokenNewline, TokenWhitespace, TokenText, TokenError,
		TokenHash, TokenAsterisk, TokenUnderscore, TokenTilde,
		TokenBacktick, TokenDash, TokenPlus, TokenDot, TokenColon,
		TokenPipe, TokenBracketOpen, TokenBracketClose,
		TokenParenOpen, TokenParenClose, TokenGreaterThan,
		TokenNumber, TokenX,
	}

	for _, tt := range allTypes {
		tok := Token{Type: tt}
		isStructural := tok.IsStructural()
		isDelimiter := tok.IsDelimiter()
		isSpecial := tt == TokenNumber ||
			tt == TokenX

		// Each token should be exactly one category (except brackets are also delimiters)
		categories := 0
		if isStructural {
			categories++
		}
		if isDelimiter {
			categories++
		}
		if isSpecial {
			categories++
		}

		if categories == 0 {
			t.Errorf(
				"Token type %v is not categorized by any helper method",
				tt,
			)
		}
		if categories <= 1 || tok.IsBracket() {
			// Brackets are both delimiters and have their own IsBracket method, which is fine
			continue
		}

		if isStructural && isDelimiter {
			t.Errorf(
				"Token type %v is incorrectly categorized as both structural and delimiter",
				tt,
			)
		}
	}
}
