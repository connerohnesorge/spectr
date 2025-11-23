package syntax

import (
	"testing"
)

func TestLexer(t *testing.T) {
	input := `# Header
Some text
## Subheader
` + "```go" + `
fmt.Println("Hello")
` + "```" + `
- List item 1
* List item 2
`

	l := Lex(input)

	expected := []struct {
		Type  TokenType
		Value string
	}{
		{TokenHeader, "# Header\n"},
		{TokenText, "Some text\n"},
		{TokenHeader, "## Subheader\n"},
		{TokenCodeBlock, "```go\nfmt.Println(\"Hello\")\n```"},
		{TokenText, "\n"}, // Newline before list?
		{TokenList, "- List item 1\n"},
		{TokenList, "* List item 2\n"},
		{TokenEOF, ""},
	}

	for i, exp := range expected {
		tok := l.NextToken()
		if tok.Type != exp.Type {
			t.Errorf("test %d: expected type %v, got %v", i, exp.Type, tok.Type)
		}
		// Value matching might be tricky with exact whitespace, but let's try
		if tok.Type != TokenEOF && tok.Value != exp.Value {
			t.Errorf("test %d: expected value %q, got %q", i, exp.Value, tok.Value)
		}
	}
}
