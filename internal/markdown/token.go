package markdown

// TokenType represents the type of a lexical token.
// Each delimiter character has its own type for fine-grained tokenization,
// enabling maximum flexibility for error recovery and precise error messages.
type TokenType uint8

const (
	// Structural tokens

	// TokenEOF signals end of input. Start == End == len(source).
	TokenEOF TokenType = iota
	// TokenNewline represents a line ending (\n or \r\n normalized).
	TokenNewline
	// TokenWhitespace represents contiguous ASCII spaces or tabs.
	TokenWhitespace
	// TokenText represents plain text content (not a delimiter).
	TokenText
	// TokenError represents invalid input with an error message.
	TokenError

	// Punctuation delimiter tokens - each character is its own token

	// TokenHash represents a single '#' character.
	TokenHash
	// TokenAsterisk represents a single '*' character.
	TokenAsterisk
	// TokenUnderscore represents a single '_' character.
	TokenUnderscore
	// TokenTilde represents a single '~' character.
	TokenTilde
	// TokenBacktick represents a single '`' character.
	TokenBacktick
	// TokenDash represents a single '-' character.
	TokenDash
	// TokenPlus represents a single '+' character.
	TokenPlus
	// TokenDot represents a single '.' character.
	TokenDot
	// TokenColon represents a single ':' character.
	TokenColon
	// TokenPipe represents a single '|' character.
	TokenPipe

	// Bracket tokens

	// TokenBracketOpen represents a '[' character.
	TokenBracketOpen
	// TokenBracketClose represents a ']' character.
	TokenBracketClose
	// TokenParenOpen represents a '(' character.
	TokenParenOpen
	// TokenParenClose represents a ')' character.
	TokenParenClose
	// TokenGreaterThan represents a '>' character.
	TokenGreaterThan

	// Special tokens

	// TokenNumber represents a sequence of digits (for ordered lists).
	TokenNumber
	// TokenX represents 'x' or 'X' in checkbox syntax.
	TokenX
)

const unknownTokenType = "Unknown"

// String returns a human-readable name for the token type.
// This is useful for debugging and error messages.
//
//nolint:revive // cyclomatic - switch cases are simple string returns
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenNewline:
		return "Newline"
	case TokenWhitespace:
		return "Whitespace"
	case TokenText:
		return "Text"
	case TokenError:
		return "Error"
	case TokenHash:
		return "Hash"
	case TokenAsterisk:
		return "Asterisk"
	case TokenUnderscore:
		return "Underscore"
	case TokenTilde:
		return "Tilde"
	case TokenBacktick:
		return "Backtick"
	case TokenDash:
		return "Dash"
	case TokenPlus:
		return "Plus"
	case TokenDot:
		return "Dot"
	case TokenColon:
		return "Colon"
	case TokenPipe:
		return "Pipe"
	case TokenBracketOpen:
		return "BracketOpen"
	case TokenBracketClose:
		return "BracketClose"
	case TokenParenOpen:
		return "ParenOpen"
	case TokenParenClose:
		return "ParenClose"
	case TokenGreaterThan:
		return "GreaterThan"
	case TokenNumber:
		return "Number"
	case TokenX:
		return "X"
	default:
		return unknownTokenType
	}
}

// Token represents a single lexical unit from the source text.
// It has position info and a zero-copy view into the original source.
type Token struct {
	// Type identifies the lexical category of this token.
	Type TokenType

	// Start is the byte offset from the beginning of the source.
	Start int

	// End is the byte offset past the last byte of this token (exclusive).
	// The token's content spans source[Start:End].
	End int

	// Source is a zero-copy slice view into the original source text.
	// It remains valid as long as the original source is retained.
	// For most tokens, this contains the token's text content.
	Source []byte

	// Message contains an error description for TokenError tokens.
	// For all other token types, this field is empty.
	Message string
}

// Len returns the byte length of the token.
func (t Token) Len() int {
	return t.End - t.Start
}

// Text returns the token's source content as a string.
// This creates a copy; use Source directly for zero-copy access.
func (t Token) Text() string {
	return string(t.Source)
}

// IsDelimiter returns true if the token is a punctuation delimiter.
func (t Token) IsDelimiter() bool {
	switch t.Type {
	case TokenHash,
		TokenAsterisk,
		TokenUnderscore,
		TokenTilde,
		TokenBacktick,
		TokenDash,
		TokenPlus,
		TokenDot,
		TokenColon,
		TokenPipe,
		TokenBracketOpen,
		TokenBracketClose,
		TokenParenOpen,
		TokenParenClose,
		TokenGreaterThan:
		return true
	case TokenEOF,
		TokenNewline,
		TokenWhitespace,
		TokenText,
		TokenError,
		TokenNumber,
		TokenX:
		return false
	default:
		return false
	}
}

// IsBracket returns true if the token is a bracket or parenthesis.
func (t Token) IsBracket() bool {
	switch t.Type {
	case TokenBracketOpen,
		TokenBracketClose,
		TokenParenOpen,
		TokenParenClose:
		return true
	case TokenEOF,
		TokenNewline,
		TokenWhitespace,
		TokenText,
		TokenError,
		TokenHash,
		TokenAsterisk,
		TokenUnderscore,
		TokenTilde,
		TokenBacktick,
		TokenDash,
		TokenPlus,
		TokenDot,
		TokenColon,
		TokenPipe,
		TokenGreaterThan,
		TokenNumber,
		TokenX:
		return false
	default:
		return false
	}
}

// IsStructural returns true if the token is structural (EOF, newline, etc).
func (t Token) IsStructural() bool {
	switch t.Type {
	case TokenEOF,
		TokenNewline,
		TokenWhitespace,
		TokenText,
		TokenError:
		return true
	case TokenHash,
		TokenAsterisk,
		TokenUnderscore,
		TokenTilde,
		TokenBacktick,
		TokenDash,
		TokenPlus,
		TokenDot,
		TokenColon,
		TokenPipe,
		TokenBracketOpen,
		TokenBracketClose,
		TokenParenOpen,
		TokenParenClose,
		TokenGreaterThan,
		TokenNumber,
		TokenX:
		return false
	default:
		return false
	}
}
