# Lexer Specification

## Requirements

### Requirement: Lexer Structure
The system SHALL provide a Lexer struct that tokenizes markdown source text into a stream of fine-grained tokens with zero-copy byte slice views.

#### Scenario: Lexer initialization
- **WHEN** `NewLexer(source []byte)` is called
- **THEN** it SHALL return a `*Lexer` holding a reference to the source
- **AND** it SHALL NOT copy the source bytes
- **AND** it SHALL initialize position to byte offset 0

#### Scenario: Lexer source retention
- **WHEN** a Lexer is created
- **THEN** it SHALL retain the source slice for the lifetime of lexing
- **AND** all emitted tokens SHALL reference slices within this source
- **AND** the caller SHALL ensure source is not modified during lexing

### Requirement: Lexer Position Tracking
The system SHALL track position as byte offset only, with line/column calculable on demand.

#### Scenario: Byte offset tracking
- **WHEN** the lexer advances through source
- **THEN** it SHALL maintain a `pos` field with current byte offset
- **AND** `pos` SHALL be incremented by actual byte count (not rune count)

#### Scenario: Line and column calculation
- **WHEN** `LineCol(offset int)` is called on the lexer or source
- **THEN** it SHALL scan source from start to offset counting newlines
- **AND** it SHALL return `(line int, col int)` where line is 1-based
- **AND** column SHALL be byte offset from last newline (1-based)

#### Scenario: Cached line index for performance
- **WHEN** multiple `LineCol` calls are made
- **THEN** the system MAY cache a line-start-offsets index
- **AND** the index SHALL be built lazily on first position query
- **AND** binary search SHALL be used for O(log n) lookup

### Requirement: Lexer Token Emission
The system SHALL emit tokens one at a time via a Next() method that returns the next token.

#### Scenario: Next token retrieval
- **WHEN** `lexer.Next()` is called
- **THEN** it SHALL return the next `Token` from the source
- **AND** it SHALL advance the lexer position past the token
- **AND** repeated calls SHALL return subsequent tokens until EOF

#### Scenario: EOF handling
- **WHEN** `Next()` is called at end of source
- **THEN** it SHALL return `Token{Type: TokenEOF, Start: len(source), End: len(source)}`
- **AND** subsequent calls SHALL continue returning `TokenEOF`

#### Scenario: Peek without consuming
- **WHEN** `lexer.Peek()` is called
- **THEN** it SHALL return the next token WITHOUT advancing position
- **AND** subsequent `Next()` SHALL return the same token
- **AND** `Peek()` SHALL cache the token to avoid re-lexing

### Requirement: Lexer CRLF Normalization
The system SHALL normalize CRLF line endings while preserving accurate byte offsets.

#### Scenario: CRLF to LF normalization
- **WHEN** the lexer encounters `\r\n` (CRLF)
- **THEN** it SHALL emit a single `TokenNewline`
- **AND** the token `Start` SHALL point to the `\r`
- **AND** the token `End` SHALL point past the `\n` (Start + 2)

#### Scenario: Standalone CR handling
- **WHEN** the lexer encounters `\r` not followed by `\n`
- **THEN** it SHALL emit `TokenNewline` for the standalone `\r`
- **AND** `End` SHALL be `Start + 1`

#### Scenario: LF handling
- **WHEN** the lexer encounters `\n`
- **THEN** it SHALL emit `TokenNewline`
- **AND** `End` SHALL be `Start + 1`

### Requirement: Lexer Error Recovery
The system SHALL emit error tokens for invalid input and continue lexing to collect all errors.

#### Scenario: Error token emission
- **WHEN** the lexer encounters unrecognized input
- **THEN** it SHALL emit `Token{Type: TokenError, ...}` with the invalid bytes
- **AND** it SHALL set `Token.Message` to describe the error
- **AND** it SHALL advance past the invalid bytes

#### Scenario: Recovery after error
- **WHEN** an error token is emitted
- **THEN** the lexer SHALL attempt to resynchronize at the next recognizable boundary
- **AND** boundaries SHALL include: newline, whitespace after newline, known delimiter
- **AND** subsequent valid tokens SHALL have correct positions

#### Scenario: Invalid UTF-8 handling
- **WHEN** the lexer encounters invalid UTF-8 byte sequences
- **THEN** it SHALL emit `TokenError` containing the invalid bytes
- **AND** the message SHALL indicate "invalid UTF-8"
- **AND** lexing SHALL continue from the next byte

### Requirement: Lexer State Machine for Context
The system SHALL use a state machine to handle context-dependent tokenization.

#### Scenario: Code fence state
- **WHEN** the lexer encounters ``` or ~~~ at line start
- **THEN** it SHALL enter `StateFencedCode`
- **AND** while in this state, it SHALL emit only `TokenText` and `TokenNewline`
- **AND** it SHALL exit when matching fence is encountered at line start

#### Scenario: Code fence content preservation
- **WHEN** in `StateFencedCode` state
- **THEN** delimiter characters (`*`, `_`, `[`, etc.) SHALL be emitted as `TokenText`
- **AND** no special tokenization SHALL occur within fenced code

#### Scenario: Inline code state
- **WHEN** the lexer encounters a backtick sequence (1-n backticks)
- **THEN** it SHALL track the opening sequence length
- **AND** content until matching sequence SHALL be `TokenText`
- **AND** nested backticks of different length SHALL be literal

#### Scenario: Link URL state
- **WHEN** the lexer is inside `(...)` following `]`
- **THEN** it SHALL enter `StateLinkURL`
- **AND** special characters SHALL be treated as `TokenText` (URLs can contain `*`, `_`, etc.)

### Requirement: Lexer Unicode Handling
The system SHALL correctly handle Unicode text including multi-byte characters.

#### Scenario: Multi-byte UTF-8 characters
- **WHEN** lexing text containing multi-byte UTF-8 characters
- **THEN** byte offsets SHALL reflect actual byte positions
- **AND** `TokenText` content SHALL include complete UTF-8 sequences
- **AND** the lexer SHALL NOT split multi-byte characters

#### Scenario: Unicode whitespace
- **WHEN** lexing Unicode whitespace (e.g., NBSP, en-space)
- **THEN** only ASCII space (0x20) and tab (0x09) SHALL be `TokenWhitespace`
- **AND** other Unicode whitespace SHALL be `TokenText`
- **AND** this matches CommonMark behavior

#### Scenario: Unicode punctuation in emphasis
- **WHEN** determining if `*` or `_` can open/close emphasis
- **THEN** Unicode punctuation categories SHALL be considered
- **AND** behavior SHALL match CommonMark section 6.2 rules

### Requirement: Lexer Collect All Mode
The system SHALL provide a convenience method to collect all tokens at once.

#### Scenario: Collect all tokens
- **WHEN** `lexer.All()` is called
- **THEN** it SHALL return `[]Token` containing all tokens including `TokenEOF`
- **AND** the lexer position SHALL be at EOF after the call

#### Scenario: Collect with error aggregation
- **WHEN** `lexer.AllWithErrors()` is called
- **THEN** it SHALL return `([]Token, []LexError)`
- **AND** `LexError` SHALL contain position and message from error tokens
- **AND** the token slice SHALL still contain the `TokenError` entries

