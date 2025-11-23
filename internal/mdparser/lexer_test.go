package mdparser

import (
	"testing"
)

func TestLexer_Headers(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "single hash header",
			input: "# Title",
			expected: []Token{
				{Type: TokenHeader, Value: "# Title"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "multiple hash header",
			input: "### Requirement: Login",
			expected: []Token{
				{Type: TokenHeader, Value: "### Requirement: Login"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "header with newline",
			input: "## Section\nNext line",
			expected: []Token{
				{Type: TokenHeader, Value: "## Section"},
				{Type: TokenText, Value: "Next line"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "multiple headers",
			input: "# H1\n## H2\n### H3",
			expected: []Token{
				{Type: TokenHeader, Value: "# H1"},
				{Type: TokenHeader, Value: "## H2"},
				{Type: TokenHeader, Value: "### H3"},
				{Type: TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d: expected type %v, got %v", i, tt.expected[i].Type, tok.Type)
				}
				if tok.Type != TokenEOF && tok.Value != tt.expected[i].Value {
					t.Errorf(
						"token %d: expected value %q, got %q",
						i,
						tt.expected[i].Value,
						tok.Value,
					)
				}
			}
		})
	}
}

func TestLexer_CodeBlocks(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "simple code block",
			input: "```\ncode\n```",
			expected: []Token{
				{Type: TokenCodeFence, Value: "```"},
				{Type: TokenCodeContent, Value: "code"},
				{Type: TokenCodeFence, Value: "```"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "code block with language",
			input: "```go\nfunc main() {}\n```",
			expected: []Token{
				{Type: TokenCodeFence, Value: "```go"},
				{Type: TokenCodeContent, Value: "func main() {}"},
				{Type: TokenCodeFence, Value: "```"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "code block with markdown syntax inside",
			input: "```\n### Not a header\n- Not a list\n```",
			expected: []Token{
				{Type: TokenCodeFence, Value: "```"},
				{Type: TokenCodeContent, Value: "### Not a header"},
				{Type: TokenCodeContent, Value: "- Not a list"},
				{Type: TokenCodeFence, Value: "```"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "code block with requirement syntax",
			input: "```markdown\n### Requirement: Should not parse\n#### Scenario: Should not parse\n```",
			expected: []Token{
				{Type: TokenCodeFence, Value: "```markdown"},
				{Type: TokenCodeContent, Value: "### Requirement: Should not parse"},
				{Type: TokenCodeContent, Value: "#### Scenario: Should not parse"},
				{Type: TokenCodeFence, Value: "```"},
				{Type: TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d: expected type %v, got %v", i, tt.expected[i].Type, tok.Type)
				}
				if tok.Type != TokenEOF && tok.Value != tt.expected[i].Value {
					t.Errorf(
						"token %d: expected value %q, got %q",
						i,
						tt.expected[i].Value,
						tok.Value,
					)
				}
			}
		})
	}
}

func TestLexer_Lists(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "unordered list",
			input: "- Item 1\n- Item 2",
			expected: []Token{
				{Type: TokenListItem, Value: "- Item 1"},
				{Type: TokenListItem, Value: "- Item 2"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "ordered list",
			input: "1. First\n2. Second",
			expected: []Token{
				{Type: TokenListItem, Value: "1. First"},
				{Type: TokenListItem, Value: "2. Second"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "list with asterisk",
			input: "* Item 1\n* Item 2",
			expected: []Token{
				{Type: TokenListItem, Value: "* Item 1"},
				{Type: TokenListItem, Value: "* Item 2"},
				{Type: TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d: expected type %v, got %v", i, tt.expected[i].Type, tok.Type)
				}
				if tok.Type != TokenEOF && tok.Value != tt.expected[i].Value {
					t.Errorf(
						"token %d: expected value %q, got %q",
						i,
						tt.expected[i].Value,
						tok.Value,
					)
				}
			}
		})
	}
}

func TestLexer_BlankLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "single blank line",
			input: "line 1\n\nline 2",
			expected: []Token{
				{Type: TokenText, Value: "line 1"},
				{Type: TokenBlankLine},
				{Type: TokenText, Value: "line 2"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "multiple blank lines",
			input: "line 1\n\n\nline 2",
			expected: []Token{
				{Type: TokenText, Value: "line 1"},
				{Type: TokenBlankLine},
				{Type: TokenBlankLine},
				{Type: TokenText, Value: "line 2"},
				{Type: TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d: expected type %v, got %v", i, tt.expected[i].Type, tok.Type)
				}
			}
		})
	}
}

func TestLexer_StateTransitions(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []TokenType
	}{
		{
			name:  "text to header",
			input: "Some text\n# Header",
			want:  []TokenType{TokenText, TokenHeader, TokenEOF},
		},
		{
			name:  "header to code block",
			input: "## Title\n```\ncode\n```",
			want: []TokenType{
				TokenHeader,
				TokenCodeFence,
				TokenCodeContent,
				TokenCodeFence,
				TokenEOF,
			},
		},
		{
			name:  "code block to list",
			input: "```\ncode\n```\n- Item",
			want: []TokenType{
				TokenCodeFence,
				TokenCodeContent,
				TokenCodeFence,
				TokenListItem,
				TokenEOF,
			},
		},
		{
			name:  "list to blank to text",
			input: "- Item\n\nText",
			want:  []TokenType{TokenListItem, TokenBlankLine, TokenText, TokenEOF},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != len(tt.want) {
				t.Fatalf("expected %d tokens, got %d", len(tt.want), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tt.want[i] {
					t.Errorf("token %d: expected type %v, got %v", i, tt.want[i], tok.Type)
				}
			}
		})
	}
}

func TestLexer_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []Token
	}{
		{
			name:  "empty input",
			input: "",
			expected: []Token{
				{Type: TokenEOF},
			},
		},
		{
			name:  "only whitespace",
			input: "   \n\t\n  ",
			expected: []Token{
				{Type: TokenBlankLine},
				{Type: TokenBlankLine},
				{Type: TokenBlankLine},
				{Type: TokenEOF},
			},
		},
		{
			name:  "unclosed code block",
			input: "```\ncode without closing fence",
			expected: []Token{
				{Type: TokenCodeFence, Value: "```"},
				{Type: TokenCodeContent, Value: "code without closing fence"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "hash not at line start",
			input: "text # not a header",
			expected: []Token{
				{Type: TokenText, Value: "text # not a header"},
				{Type: TokenEOF},
			},
		},
		{
			name:  "dash not followed by space",
			input: "-not a list",
			expected: []Token{
				{Type: TokenText, Value: "-not a list"},
				{Type: TokenEOF},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := NewLexer(tt.input)
			tokens := lexer.Tokenize()

			if len(tokens) != len(tt.expected) {
				t.Fatalf("expected %d tokens, got %d", len(tt.expected), len(tokens))
			}

			for i, tok := range tokens {
				if tok.Type != tt.expected[i].Type {
					t.Errorf("token %d: expected type %v, got %v", i, tt.expected[i].Type, tok.Type)
				}
			}
		})
	}
}

func TestLexer_Position(t *testing.T) {
	input := "# Header\nText line\n\n```\ncode\n```"
	lexer := NewLexer(input)
	tokens := lexer.Tokenize()

	// Verify positions are being tracked
	for i, tok := range tokens {
		if tok.Type == TokenEOF {
			continue
		}
		if tok.Pos.Line < 1 {
			t.Errorf("token %d: invalid line number %d", i, tok.Pos.Line)
		}
		if tok.Pos.Column < 1 {
			t.Errorf("token %d: invalid column number %d", i, tok.Pos.Column)
		}
	}
}

func TestLexer_ComplexDocument(t *testing.T) {
	input := `# Main Title

## Section 1

Some paragraph text here.

- List item 1
- List item 2

## Section 2

` + "```go" + `
func main() {
    // ### Not a header
    // - Not a list
}
` + "```" + `

Final text.
`

	lexer := NewLexer(input)
	tokens := lexer.Tokenize()

	// Verify we got reasonable tokens
	if len(tokens) < 5 {
		t.Fatalf("expected at least 5 tokens for complex document, got %d", len(tokens))
	}

	// Verify code content is not parsed as headers/lists
	foundCodeContent := false
	for _, tok := range tokens {
		if tok.Type == TokenCodeContent {
			foundCodeContent = true
			if tok.Value == "    // ### Not a header" {
				// Good - this is inside code block
				continue
			}
		}
		// Make sure we don't have headers inside code blocks
		if foundCodeContent && tok.Type == TokenHeader && tok.Value == "### Not a header" {
			t.Error("found header token inside code block - state machine failed")
		}
	}
}
