package syntax

import "fmt"

// TokenType identifies the type of lexical tokens.
type TokenType int

const (
	TokenError TokenType = iota
	TokenEOF
	TokenText
	TokenHeader    // #, ##, etc.
	TokenCodeBlock // ``` ... ```
	TokenList      // -, *, 1.
)

// Token represents a lexical token.
type Token struct {
	Type  TokenType
	Value string
	Line  int
}

func (t Token) String() string {
	switch t.Type {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return t.Value
	}
	if len(t.Value) > 10 {
		return fmt.Sprintf("%d:%q...", t.Type, t.Value[:10])
	}
	return fmt.Sprintf("%d:%q", t.Type, t.Value)
}
