# Parser Specification

## Requirements

### Requirement: Stateless Parse Function
The system SHALL provide a stateless Parse function that transforms source bytes into an immutable AST.

#### Scenario: Parse function signature
- **WHEN** `Parse(source []byte)` is called
- **THEN** it SHALL return `(Node, []ParseError)`
- **AND** the Node SHALL be the root `NodeDocument`
- **AND** errors SHALL be collected, not fatal (unless catastrophic)
- **AND** each call SHALL be independent (no shared state)

#### Scenario: Parse function concurrency
- **WHEN** multiple Parse calls run concurrently
- **THEN** they SHALL be safe without synchronization
- **AND** each call SHALL use its own internal state
- **AND** no global mutable state SHALL be accessed

#### Scenario: Internal pooling
- **WHEN** Parse executes
- **THEN** it SHALL use object pools internally for tokens and temporary allocations
- **AND** pools SHALL be accessed via sync.Pool (no explicit management)
- **AND** the caller SHALL NOT need to manage pool lifecycle

### Requirement: Tree-Sitter Style Incremental Parsing
The system SHALL support incremental reparsing where only changed portions of the document are re-parsed, reusing unchanged subtrees.

#### Scenario: Incremental parse with edit
- **WHEN** `ParseIncremental(oldTree Node, oldSource, newSource []byte)` is called
- **THEN** it SHALL compute the diff between oldSource and newSource
- **AND** it SHALL reparse only affected regions
- **AND** unchanged subtrees SHALL be reused (same Node pointers where hash matches)
- **AND** this is also a stateless function (no Parser struct)

#### Scenario: Edit region detection
- **WHEN** computing diff between old and new source
- **THEN** the parser SHALL identify: (startOffset, oldEndOffset, newEndOffset)
- **AND** it SHALL find the minimal changed region
- **AND** regions before and after the change MAY be reused

#### Scenario: Subtree reuse via hash
- **WHEN** reparsing a region
- **THEN** newly parsed nodes SHALL have hashes computed
- **AND** if a new node's hash matches an old node at same relative position
- **THEN** the old node SHALL be reused if its source offsets are still valid

#### Scenario: Offset adjustment for reused nodes
- **WHEN** reusing nodes after an edit
- **THEN** nodes before the edit SHALL have unchanged offsets
- **AND** nodes after the edit SHALL have offsets adjusted by (newLen - oldLen)
- **AND** the adjustment SHALL be done lazily or via offset transform

### Requirement: Diff-Based Edit Detection
The system SHALL accept old and new source text and compute the edit internally using an efficient diff algorithm.

#### Scenario: Diff algorithm selection
- **WHEN** computing diff between sources
- **THEN** the parser SHALL use a linear-time algorithm for common case (single edit point)
- **AND** it SHALL fall back to Myers diff for complex multi-edit cases
- **AND** the diff SHALL produce minimal edit regions

#### Scenario: Single edit optimization
- **WHEN** sources differ at only one contiguous region
- **THEN** the parser SHALL detect this in O(n) time via prefix/suffix matching
- **AND** this SHALL be the fast path for typical editing

#### Scenario: Multi-edit handling
- **WHEN** sources have multiple disjoint edit regions
- **THEN** the parser SHALL identify all changed regions
- **AND** each region SHALL be reparsed independently
- **AND** unchanged regions between edits SHALL be reused

### Requirement: Parser Error Handling
The system SHALL collect all parse errors and continue parsing to provide maximum feedback.

#### Scenario: Error collection mode
- **WHEN** the parser encounters an error
- **THEN** it SHALL add a `ParseError` to the error list
- **AND** it SHALL attempt error recovery
- **AND** parsing SHALL continue if possible

#### Scenario: Parse error structure
- **WHEN** a `ParseError` is created
- **THEN** it SHALL contain: `Offset int`, `Message string`, `Expected []TokenType`
- **AND** `Offset` SHALL be the byte offset where error occurred
- **AND** `Expected` SHALL list tokens that would have been valid

#### Scenario: Error recovery strategy
- **WHEN** recovering from a parse error
- **THEN** the parser SHALL skip tokens until a synchronization point
- **AND** sync points SHALL include: newline after blank line, header, list marker
- **AND** an error node MAY be inserted to represent skipped content

#### Scenario: Maximum errors limit
- **WHEN** error count exceeds a threshold (default: 100)
- **THEN** parsing SHALL abort with a "too many errors" error
- **AND** partial AST up to that point SHALL still be returned

### Requirement: Block-Level Parsing
The system SHALL parse block-level elements following CommonMark-like precedence rules.

#### Scenario: Block element detection
- **WHEN** parsing at block level
- **THEN** the parser SHALL check for (in order): code fence, header, blockquote, list item, paragraph

#### Scenario: Header parsing
- **WHEN** a line starts with 1-6 `TokenHash` followed by `TokenWhitespace`
- **THEN** the parser SHALL create `NodeSection` with appropriate level
- **AND** remaining line content SHALL be parsed for inline formatting
- **AND** Spectr-specific headers (Requirement:, Scenario:) SHALL create specialized nodes

#### Scenario: List parsing with nesting
- **WHEN** lines start with `TokenDash`/`TokenPlus`/`TokenNumber+TokenDot`
- **THEN** the parser SHALL create `NodeList` containing `NodeListItem` children
- **AND** indentation SHALL determine nesting depth
- **AND** checkbox syntax SHALL create task items

#### Scenario: Code fence parsing
- **WHEN** a line starts with 3+ backticks or tildes
- **THEN** the parser SHALL create `NodeCodeBlock`
- **AND** content until matching fence SHALL be verbatim (no inline parsing)
- **AND** fence length and character SHALL be tracked for matching

#### Scenario: Paragraph parsing
- **WHEN** no other block element matches
- **THEN** consecutive non-blank lines SHALL form a `NodeParagraph`
- **AND** paragraph content SHALL be parsed for inline elements
- **AND** blank lines or block elements terminate the paragraph

### Requirement: Inline-Level Parsing
The system SHALL parse inline formatting within block elements, handling emphasis precedence correctly.

#### Scenario: Emphasis delimiter matching
- **WHEN** parsing emphasis (`*`, `_`)
- **THEN** the parser SHALL track a delimiter stack
- **AND** opening delimiters SHALL be matched with compatible closing delimiters
- **AND** CommonMark rules for left/right flanking SHALL be applied

#### Scenario: Emphasis nesting
- **WHEN** emphasis markers are nested like `***bold and italic***`
- **THEN** the parser SHALL create nested `NodeStrong` and `NodeEmphasis`
- **AND** the innermost emphasis type depends on delimiter arrangement

#### Scenario: Inline code parsing
- **WHEN** backticks are encountered
- **THEN** content until matching backtick sequence SHALL be `NodeCode`
- **AND** backtick count SHALL match (``` requires ```)
- **AND** content SHALL NOT be parsed for other inline elements

#### Scenario: Link parsing
- **WHEN** `[...]` is followed by `(...)`
- **THEN** the parser SHALL create `NodeLink`
- **AND** bracket content SHALL be parsed for inline formatting
- **AND** URL in parens SHALL NOT be parsed for inline formatting

#### Scenario: Wikilink parsing
- **WHEN** `[[...]]` is encountered
- **THEN** the parser SHALL create `NodeWikilink`
- **AND** content SHALL be split on `|` for target/display
- **AND** content SHALL be split on `#` for anchor

### Requirement: Reference Link Resolution
The system SHALL collect link definitions and resolve reference links in a two-pass approach.

#### Scenario: Link definition collection
- **WHEN** parsing block level
- **THEN** lines matching `[label]: url "title"` SHALL be collected
- **AND** definitions SHALL be stored in a map (label -> definition)
- **AND** labels SHALL be case-insensitive for matching

#### Scenario: Reference link resolution
- **WHEN** `[text][ref]` is encountered during inline parsing
- **THEN** the parser SHALL look up `ref` in the definitions map
- **AND** if found, `NodeLink` SHALL be created with resolved URL
- **AND** if not found, `ParseError` SHALL be added

#### Scenario: Shortcut reference links
- **WHEN** `[text][]` or `[text]` (without second brackets) is encountered
- **THEN** the parser SHALL use `text` as the reference label
- **AND** resolution SHALL proceed as with explicit reference

### Requirement: Parser State for Incremental
The system SHALL maintain parse state to enable efficient incremental updates.

#### Scenario: Parse state structure
- **WHEN** a parse completes
- **THEN** the result SHALL include the AST and a `ParseState` for incremental use
- **AND** `ParseState` SHALL contain: link definitions, line index, source reference

#### Scenario: State reuse for incremental
- **WHEN** `ParseIncremental` is called
- **THEN** unchanged link definitions SHALL be reused
- **AND** line index for unchanged regions SHALL be reused
- **AND** new definitions in changed regions SHALL be merged

### Requirement: Spectr-Specific Parsing
The system SHALL recognize and parse Spectr-specific markdown patterns.

#### Scenario: Requirement header parsing
- **WHEN** a line matches `### Requirement: {name}`
- **THEN** the parser SHALL create `NodeRequirement` with extracted name
- **AND** the node SHALL be a child of the enclosing section

#### Scenario: Scenario header parsing
- **WHEN** a line matches `#### Scenario: {name}`
- **THEN** the parser SHALL create `NodeScenario` with extracted name
- **AND** the node SHALL be a child of the enclosing requirement

#### Scenario: Delta section parsing
- **WHEN** a line matches `## ADDED|MODIFIED|REMOVED|RENAMED Requirements`
- **THEN** the parser SHALL create `NodeSection` with `DeltaType` field set
- **AND** the section SHALL be flagged as a delta section

#### Scenario: WHEN/THEN/AND pattern parsing
- **WHEN** a list item contains `**WHEN**`, `**THEN**`, or `**AND**` at start
- **THEN** the parser SHALL set a `Keyword` field on the list item
- **AND** the keyword SHALL be extracted for semantic use

### Requirement: CommonMark Strict Emphasis
The system SHALL follow CommonMark specification strictly for emphasis parsing, including edge cases.

#### Scenario: Left-flanking delimiter run
- **WHEN** a delimiter run is evaluated for opening emphasis
- **THEN** it SHALL be left-flanking if: (1) not followed by whitespace, AND (2) not followed by punctuation OR preceded by whitespace/punctuation
- **AND** this rule SHALL be applied per CommonMark spec section 6.2

#### Scenario: Right-flanking delimiter run
- **WHEN** a delimiter run is evaluated for closing emphasis
- **THEN** it SHALL be right-flanking if: (1) not preceded by whitespace, AND (2) not preceded by punctuation OR followed by whitespace/punctuation
- **AND** this rule SHALL be applied per CommonMark spec section 6.2

#### Scenario: Underscore intraword restriction
- **WHEN** `_` delimiters are surrounded by alphanumeric characters
- **THEN** they SHALL NOT open or close emphasis
- **AND** `foo_bar_baz` SHALL be parsed as literal text, not emphasis
- **AND** `*` delimiters do NOT have this restriction

#### Scenario: Overlapping delimiter resolution
- **WHEN** delimiters overlap like `*a _b* c_`
- **THEN** the parser SHALL follow CommonMark's "process emphasis" algorithm
- **AND** the result SHALL match CommonMark spec exactly
- **AND** `*a _b* c_` SHALL parse as: emphasis("a _b") + text(" c_")

#### Scenario: Delimiter stack processing
- **WHEN** multiple potential emphasis delimiters exist
- **THEN** the parser SHALL use a delimiter stack
- **AND** it SHALL process bottom-up, matching compatible openers with closers
- **AND** closer must match opener type (* matches *, _ matches _)

#### Scenario: Triple delimiter handling
- **WHEN** `***text***` is encountered
- **THEN** it SHALL create nested strong+emphasis nodes
- **AND** the outer node type depends on delimiter arrangement
- **AND** common patterns SHALL be handled efficiently

