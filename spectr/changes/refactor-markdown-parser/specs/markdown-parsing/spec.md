# Markdown Parsing Capability

## ADDED Requirements

### Requirement: Markdown Tokenization

The markdown parser SHALL tokenize markdown content into a stream of typed tokens with position information.

#### Scenario: Heading tokenization

- **WHEN** markdown contains heading syntax (e.g., `# Title`, `## Section`, `### Requirement:`)
- **THEN** the lexer SHALL emit HEADING_N tokens with appropriate level (1-6)
- **AND** SHALL capture the heading text without the hash symbols
- **AND** SHALL record accurate line and column position for each token

#### Scenario: Code fence tokenization

- **WHEN** markdown contains triple backtick code fences with optional language identifier
- **THEN** the lexer SHALL emit CODE_FENCE_START and CODE_FENCE_END tokens
- **AND** SHALL preserve all content between fences verbatim
- **AND** SHALL NOT tokenize markdown syntax inside code fences
- **AND** SHALL handle nested backticks and markdown-like syntax within code blocks

#### Scenario: List tokenization

- **WHEN** markdown contains list items (unordered `-` or ordered `1.`)
- **THEN** the lexer SHALL emit LIST_ITEM tokens
- **AND** SHALL handle indented nested lists
- **AND** SHALL preserve list item content

#### Scenario: Position tracking

- **WHEN** tokenizing markdown content
- **THEN** each token SHALL include line number (1-indexed)
- **AND** each token SHALL include column number (1-indexed)
- **AND** position information SHALL be accurate for error reporting

### Requirement: AST Construction

The markdown parser SHALL construct an abstract syntax tree representing the document structure.

#### Scenario: Document parsing

- **WHEN** parsing a token stream
- **THEN** the parser SHALL create a Document root node
- **AND** SHALL build child nodes for each structural element (headings, code blocks, lists, paragraphs)
- **AND** SHALL preserve parent-child relationships

#### Scenario: Heading hierarchy

- **WHEN** parsing headings of different levels
- **THEN** the parser SHALL determine content boundaries based on heading levels
- **AND** SHALL treat content after a heading as children until the next heading of equal or higher level
- **AND** SHALL properly nest lower-level headings within higher-level headings

#### Scenario: Code block preservation

- **WHEN** parsing code blocks
- **THEN** the parser SHALL preserve all content verbatim
- **AND** SHALL NOT interpret markdown syntax within code blocks
- **AND** SHALL distinguish between fenced code blocks and indented code blocks
- **AND** SHALL preserve language identifiers for syntax highlighting hints

#### Scenario: Nested structures

- **WHEN** parsing nested markdown structures (lists with code blocks, lists with nested lists)
- **THEN** the parser SHALL correctly represent nesting in the AST
- **AND** SHALL maintain proper parent-child relationships

### Requirement: Spectr Structure Extraction

The markdown parser SHALL extract Spectr-specific structures from the AST.

#### Scenario: Requirement extraction

- **WHEN** extracting requirements from a spec file AST
- **THEN** the extractor SHALL identify `### Requirement:` headings within `## Requirements` sections
- **AND** SHALL extract requirement name from heading text
- **AND** SHALL collect all content until the next requirement or section boundary
- **AND** SHALL NOT extract requirement headers found inside code blocks

#### Scenario: Scenario extraction

- **WHEN** extracting scenarios from requirement content
- **THEN** the extractor SHALL identify `#### Scenario:` headings
- **AND** SHALL extract scenario name and content
- **AND** SHALL include WHEN/THEN clauses in scenario content
- **AND** SHALL NOT extract scenario headers found inside code blocks

#### Scenario: Section extraction

- **WHEN** extracting sections from a spec file
- **THEN** the extractor SHALL identify `## Section Name` headings
- **AND** SHALL map section names to their content
- **AND** SHALL collect content until the next `##` heading
- **AND** SHALL preserve formatting within sections

#### Scenario: Delta operation extraction

- **WHEN** extracting delta operations from a change spec
- **THEN** the extractor SHALL identify `## ADDED Requirements`, `## MODIFIED Requirements`, `## REMOVED Requirements`, and `## RENAMED Requirements` sections
- **AND** SHALL extract requirement blocks within each section
- **AND** SHALL parse FROM/TO pairs in RENAMED sections
- **AND** SHALL return a structured DeltaPlan with all operations

### Requirement: Context-Aware Parsing

The markdown parser SHALL distinguish between different markdown contexts to avoid false matches.

#### Scenario: Ignoring markdown syntax in code blocks

- **WHEN** parsing a spec containing a code block with markdown syntax inside
- **THEN** headers inside code blocks SHALL NOT be parsed as document headers
- **AND** requirement syntax inside code blocks SHALL NOT be extracted as requirements
- **AND** scenario syntax inside code blocks SHALL NOT be extracted as scenarios
- **AND** the code block content SHALL be preserved exactly as written

#### Scenario: Distinguishing code fence types

- **WHEN** parsing fenced code blocks with different delimiters (triple backtick, triple tilde)
- **THEN** the parser SHALL recognize both types
- **AND** SHALL match opening and closing delimiters correctly
- **AND** SHALL handle language identifiers after opening delimiter

#### Scenario: Handling indented code blocks

- **WHEN** parsing indented code blocks (4 or more spaces)
- **THEN** the parser SHALL recognize them as code blocks
- **AND** SHALL preserve content verbatim
- **AND** SHALL NOT interpret markdown syntax within indented code

### Requirement: Error Handling and Reporting

The markdown parser SHALL provide detailed error messages with position information for malformed markdown.

#### Scenario: Unclosed code fence error

- **WHEN** parsing markdown with an unclosed code fence (opening ``` without closing ```)
- **THEN** the parser SHALL report a ParseError
- **AND** the error SHALL include line and column of the unclosed fence
- **AND** the error message SHALL include context (surrounding lines)

#### Scenario: Position-aware error messages

- **WHEN** any parsing error occurs
- **THEN** the error SHALL include line number (1-indexed)
- **AND** the error SHALL include column number (1-indexed)
- **AND** the error SHALL include relevant context for debugging
- **AND** the error message SHALL be actionable (suggest fix)

### Requirement: Performance and Scalability

The markdown parser SHALL parse spec files efficiently without performance regression from the previous regex-based implementation.

#### Scenario: Parsing performance target

- **WHEN** parsing a 100KB markdown file
- **THEN** parsing SHALL complete in less than 10 milliseconds
- **AND** performance SHALL be within 10% of the previous regex implementation
- **AND** performance SHALL scale linearly with file size

#### Scenario: Memory efficiency

- **WHEN** parsing a typical spec file (less than 100 requirements)
- **THEN** total memory allocation SHALL be less than 1MB
- **AND** SHALL not leak memory over multiple parse operations
- **AND** SHALL reuse allocations where possible

#### Scenario: Large document handling

- **WHEN** parsing a large spec file (1000 requirements)
- **THEN** parsing SHALL complete successfully
- **AND** performance SHALL remain linear (no exponential slowdown)
- **AND** memory usage SHALL scale linearly with document size

### Requirement: API Compatibility

The markdown parser SHALL maintain backward compatibility with existing parser APIs in the internal/parsers package.

#### Scenario: Public API preservation

- **WHEN** integrating the new parser into internal/parsers
- **THEN** all existing public function signatures SHALL remain unchanged
- **AND** function return types SHALL match previous implementation
- **AND** existing code using the parsers package SHALL continue to work without modification

#### Scenario: Output equivalence

- **WHEN** parsing existing spec files with the new parser
- **THEN** extracted requirements SHALL match previous parser output
- **AND** extracted scenarios SHALL match previous parser output
- **AND** extracted sections SHALL match previous parser output
- **AND** any differences SHALL be documented as intentional improvements (e.g., correctly handling previously broken edge cases)

### Requirement: Lexer State Machine

The markdown lexer SHALL use a state function pattern for tokenization.

#### Scenario: State function transitions

- **WHEN** the lexer is tokenizing markdown
- **THEN** the lexer SHALL maintain a current state function
- **AND** each state function SHALL return the next state function
- **AND** SHALL transition between states based on input characters
- **AND** SHALL emit tokens at appropriate boundaries

#### Scenario: State: Normal

- **WHEN** lexer is in Normal state
- **THEN** SHALL recognize heading start (`#` characters)
- **AND** SHALL recognize code fence start (triple backtick)
- **AND** SHALL recognize list item start (`-` or digit followed by `.`)
- **AND** SHALL accumulate regular text
- **AND** SHALL transition to appropriate state for each syntax element

#### Scenario: State: InCodeFence

- **WHEN** lexer enters InCodeFence state (after triple backtick)
- **THEN** SHALL ignore all markdown syntax until closing triple backtick
- **AND** SHALL accumulate all content as TEXT token
- **AND** SHALL emit CODE_FENCE_END when closing delimiter found
- **AND** SHALL return to Normal state after closing delimiter

#### Scenario: EOF handling

- **WHEN** lexer reaches end of input
- **THEN** SHALL emit EOF token
- **AND** SHALL finalize any incomplete tokens
- **AND** SHALL report error for unclosed structures (code fences, etc.)
