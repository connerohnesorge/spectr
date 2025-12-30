# Tokens Specification

## Requirements

### Requirement: Token Structure

The system SHALL define a Token struct that represents a single lexical unit
with byte-offset position and a view into the source text.

#### Scenario: Token contains byte offset position

- **WHEN** a Token is created during lexing
- **THEN** it SHALL contain a `Start` field with the byte offset from source
  start
- **AND** it SHALL contain an `End` field with the byte offset of the token end
  (exclusive)
- **AND** line and column SHALL be calculable on-demand from byte offset and
  source

#### Scenario: Token contains source view

- **WHEN** a Token references text content
- **THEN** it SHALL store a `[]byte` slice view into the original source
- **AND** it SHALL NOT copy the source text (zero-copy)
- **AND** the slice SHALL remain valid as long as the source is retained

#### Scenario: Token type identification

- **WHEN** a Token is examined
- **THEN** it SHALL have a `Type` field of type `TokenType`
- **AND** the type SHALL uniquely identify the lexical category

### Requirement: Fine-Grained Token Types

The system SHALL define fine-grained token types where each delimiter character
is a separate token for maximum error recovery flexibility.

#### Scenario: Emphasis delimiters as separate tokens

- **WHEN** lexing `**bold**`
- **THEN** the lexer SHALL emit: `TokenAsterisk`, `TokenAsterisk`,
  `TokenText("bold")`, `TokenAsterisk`, `TokenAsterisk`
- **AND** each asterisk SHALL be a separate token with its own position

#### Scenario: Underscore delimiters as separate tokens

- **WHEN** lexing `__underline__`
- **THEN** the lexer SHALL emit: `TokenUnderscore`, `TokenUnderscore`,
  `TokenText("underline")`, `TokenUnderscore`, `TokenUnderscore`

#### Scenario: Tilde delimiters for strikethrough

- **WHEN** lexing `~~struck~~`
- **THEN** the lexer SHALL emit: `TokenTilde`, `TokenTilde`,
  `TokenText("struck")`, `TokenTilde`, `TokenTilde`

#### Scenario: Backtick delimiters for code

- **WHEN** lexing `` `code` ``
- **THEN** the lexer SHALL emit: `TokenBacktick`, `TokenText("code")`,
  `TokenBacktick`

### Requirement: Block-Level Token Types

The system SHALL define token types for block-level markdown elements.

#### Scenario: Header hash tokens

- **WHEN** lexing `## Header`
- **THEN** the lexer SHALL emit: `TokenHash`, `TokenHash`, `TokenWhitespace`,
  `TokenText("Header")`
- **AND** each `#` SHALL be a separate `TokenHash`

#### Scenario: List bullet tokens

- **WHEN** lexing `- item`
- **THEN** the lexer SHALL emit: `TokenDash`, `TokenWhitespace`,
  `TokenText("item")`

#### Scenario: Ordered list tokens

- **WHEN** lexing `1. item`
- **THEN** the lexer SHALL emit: `TokenNumber("1")`, `TokenDot`,
  `TokenWhitespace`, `TokenText("item")`

#### Scenario: Checkbox tokens

- **WHEN** lexing `- [ ] task`
- **THEN** the lexer SHALL emit: `TokenDash`, `TokenWhitespace`,
  `TokenBracketOpen`, `TokenWhitespace`, `TokenBracketClose`, `TokenWhitespace`,
  `TokenText("task")`

#### Scenario: Checked checkbox tokens

- **WHEN** lexing `- [x] done`
- **THEN** the lexer SHALL emit: `TokenDash`, `TokenWhitespace`,
  `TokenBracketOpen`, `TokenX`, `TokenBracketClose`, `TokenWhitespace`,
  `TokenText("done")`

#### Scenario: Blockquote tokens

- **WHEN** lexing `> quoted`
- **THEN** the lexer SHALL emit: `TokenGreaterThan`, `TokenWhitespace`,
  `TokenText("quoted")`

#### Scenario: Code fence tokens

- **WHEN** lexing a line starting with ``` or ~~~
- **THEN** the lexer SHALL emit `TokenBacktick` or `TokenTilde` for each
  character
- **AND** SHALL emit `TokenText` for the optional language identifier

### Requirement: Link and Wikilink Token Types

The system SHALL define token types for links and wikilinks with their component
parts.

#### Scenario: Inline link component tokens

- **WHEN** lexing `[text](url)`
- **THEN** the lexer SHALL emit: `TokenBracketOpen`, `TokenText("text")`,
  `TokenBracketClose`, `TokenParenOpen`, `TokenText("url")`, `TokenParenClose`

#### Scenario: Reference link component tokens

- **WHEN** lexing `[text][ref]`
- **THEN** the lexer SHALL emit: `TokenBracketOpen`, `TokenText("text")`,
  `TokenBracketClose`, `TokenBracketOpen`, `TokenText("ref")`,
  `TokenBracketClose`

#### Scenario: Wikilink tokens

- **WHEN** lexing `[[target]]`
- **THEN** the lexer SHALL emit: `TokenBracketOpen`, `TokenBracketOpen`,
  `TokenText("target")`, `TokenBracketClose`, `TokenBracketClose`

#### Scenario: Wikilink with display text

- **WHEN** lexing `[[target|display]]`
- **THEN** the lexer SHALL emit: `TokenBracketOpen`, `TokenBracketOpen`,
  `TokenText("target")`, `TokenPipe`, `TokenText("display")`,
  `TokenBracketClose`, `TokenBracketClose`

#### Scenario: Wikilink with anchor

- **WHEN** lexing `[[target#anchor]]`
- **THEN** the lexer SHALL emit: `TokenBracketOpen`, `TokenBracketOpen`,
  `TokenText("target")`, `TokenHash`, `TokenText("anchor")`,
  `TokenBracketClose`, `TokenBracketClose`

### Requirement: Whitespace and Structure Tokens

The system SHALL define tokens for whitespace, newlines, and structural
elements.

#### Scenario: Whitespace token

- **WHEN** lexing spaces or tabs
- **THEN** the lexer SHALL emit `TokenWhitespace` containing all contiguous
  whitespace
- **AND** the token SHALL preserve the original whitespace characters

#### Scenario: Newline token

- **WHEN** lexing a line ending
- **THEN** the lexer SHALL emit `TokenNewline`
- **AND** CRLF SHALL be normalized to a single `TokenNewline` but byte offsets
  SHALL reflect original positions

#### Scenario: End of file token

- **WHEN** lexing reaches end of input
- **THEN** the lexer SHALL emit `TokenEOF`
- **AND** `TokenEOF` SHALL have `Start == End == len(source)`

### Requirement: Error Token Type

The system SHALL define an error token type for invalid input that allows
continued lexing.

#### Scenario: Error token for invalid input

- **WHEN** the lexer encounters unrecognized input
- **THEN** it SHALL emit `TokenError` containing the invalid byte sequence
- **AND** it SHALL continue lexing from the next recognizable position
- **AND** the error token SHALL include a descriptive message

#### Scenario: Error token preserves position

- **WHEN** a `TokenError` is emitted
- **THEN** it SHALL have accurate `Start` and `End` byte offsets
- **AND** subsequent tokens SHALL have correct positions (not offset by error)

### Requirement: Token Type Enumeration

The system SHALL define a complete enumeration of all token types as a Go type.

#### Scenario: Token type is typed constant

- **WHEN** token types are defined
- **THEN** they SHALL be `const` values of type `TokenType`
- **AND** `TokenType` SHALL be based on `uint8` for compact storage
- **AND** token types SHALL be grouped by category with gaps for future
  additions

#### Scenario: Token type is stringable

- **WHEN** a `TokenType` is converted to string
- **THEN** it SHALL return a human-readable name (e.g., "TokenHash",
  "TokenText")
- **AND** the String() method SHALL be auto-generated or maintained consistently

#### Scenario: Complete token type list

- **WHEN** the token type enumeration is defined
- **THEN** it SHALL include at minimum:
  - Structural: `TokenEOF`, `TokenNewline`, `TokenWhitespace`, `TokenError`
  - Punctuation: `TokenHash`, `TokenAsterisk`, `TokenUnderscore`, `TokenTilde`,
    `TokenBacktick`
  - Brackets: `TokenBracketOpen`, `TokenBracketClose`, `TokenParenOpen`,
    `TokenParenClose`
  - List: `TokenDash`, `TokenPlus`, `TokenNumber`, `TokenDot`, `TokenX`
  - Special: `TokenGreaterThan`, `TokenPipe`, `TokenColon`
  - Content: `TokenText`
