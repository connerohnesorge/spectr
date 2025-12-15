# Printer Specification

## Requirements

### Requirement: Normalized Printer
The system SHALL provide a Printer that regenerates markdown source from AST nodes with consistent, minimal formatting.

#### Scenario: Print function signature
- **WHEN** printing an AST to markdown
- **THEN** the signature SHALL be `Print(node Node) []byte`
- **AND** it SHALL return normalized markdown as byte slice
- **AND** the output SHALL be valid markdown that parses to equivalent AST

#### Scenario: Print to writer
- **WHEN** printing to an io.Writer
- **THEN** the signature SHALL be `PrintTo(w io.Writer, node Node) error`
- **AND** it SHALL stream output without buffering entire result
- **AND** it SHALL return any write errors

### Requirement: Minimal Whitespace Formatting
The system SHALL use minimal whitespace in printed output for compact representation.

#### Scenario: Single space between inline elements
- **WHEN** printing inline content like text and emphasis
- **THEN** single space SHALL separate adjacent elements where needed
- **AND** no trailing whitespace SHALL appear on lines

#### Scenario: Single newline between blocks
- **WHEN** printing block elements like paragraphs and lists
- **THEN** single blank line SHALL separate distinct blocks
- **AND** consecutive headers SHALL have single blank line between them

#### Scenario: No leading/trailing blank lines
- **WHEN** printing a document
- **THEN** output SHALL NOT start with blank lines
- **AND** output SHALL NOT end with multiple newlines (one trailing newline allowed)

### Requirement: Header Printing
The system SHALL print headers with ATX style (hash prefix).

#### Scenario: Print section header
- **WHEN** printing a NodeSection with Level 2
- **THEN** output SHALL be `## {Title}\n`
- **AND** hash count SHALL match Level field

#### Scenario: Print requirement header
- **WHEN** printing a NodeRequirement
- **THEN** output SHALL be `### Requirement: {Name}\n`
- **AND** the "Requirement: " prefix SHALL always be present

#### Scenario: Print scenario header
- **WHEN** printing a NodeScenario
- **THEN** output SHALL be `#### Scenario: {Name}\n`
- **AND** the "Scenario: " prefix SHALL always be present

### Requirement: List Printing
The system SHALL print lists with consistent bullet style and indentation.

#### Scenario: Print unordered list
- **WHEN** printing a NodeList that is unordered
- **THEN** each item SHALL start with `- ` (dash space)
- **AND** nested content SHALL be indented by 2 spaces

#### Scenario: Print ordered list
- **WHEN** printing a NodeList that is ordered
- **THEN** each item SHALL start with `{n}. ` where n is 1-based index
- **AND** numbering SHALL always start at 1

#### Scenario: Print checkbox list item
- **WHEN** printing a NodeListItem with checkbox
- **THEN** unchecked SHALL print as `- [ ] {content}`
- **AND** checked SHALL print as `- [x] {content}`

#### Scenario: Print nested list
- **WHEN** printing a list item containing a nested list
- **THEN** nested list SHALL be indented by 2 spaces
- **AND** indentation SHALL accumulate for deeper nesting

### Requirement: Code Block Printing
The system SHALL print code blocks with consistent fence style.

#### Scenario: Print fenced code block
- **WHEN** printing a NodeCodeBlock
- **THEN** output SHALL use triple backticks as fence
- **AND** language identifier SHALL follow opening fence if present
- **AND** content SHALL be printed verbatim

#### Scenario: Code block content preservation
- **WHEN** printing code block content
- **THEN** internal newlines SHALL be preserved exactly
- **AND** no trailing whitespace normalization SHALL occur inside code

### Requirement: Inline Formatting Printing
The system SHALL print inline formatting with consistent delimiter style.

#### Scenario: Print strong emphasis
- **WHEN** printing a NodeStrong
- **THEN** output SHALL use `**{content}**` (asterisk style)
- **AND** nested content SHALL be printed recursively

#### Scenario: Print emphasis
- **WHEN** printing a NodeEmphasis
- **THEN** output SHALL use `*{content}*` (asterisk style)
- **AND** underscore style SHALL NOT be used

#### Scenario: Print strikethrough
- **WHEN** printing a NodeStrikethrough
- **THEN** output SHALL use `~~{content}~~`

#### Scenario: Print inline code
- **WHEN** printing a NodeCode
- **THEN** output SHALL use single backticks: `` `{content}` ``
- **AND** if content contains backtick, use double backticks with space padding

### Requirement: Link Printing
The system SHALL print links in inline style by default.

#### Scenario: Print inline link
- **WHEN** printing a NodeLink with URL
- **THEN** output SHALL be `[{text}]({url})`
- **AND** title SHALL be included if present: `[{text}]({url} "{title}")`

#### Scenario: Print wikilink
- **WHEN** printing a NodeWikilink
- **THEN** output SHALL be `[[{target}]]`
- **AND** display text SHALL be included if different: `[[{target}|{display}]]`
- **AND** anchor SHALL be included if present: `[[{target}#{anchor}]]`

### Requirement: Blockquote Printing
The system SHALL print blockquotes with consistent marker style.

#### Scenario: Print blockquote
- **WHEN** printing a NodeBlockquote
- **THEN** each line SHALL be prefixed with `> `
- **AND** nested blockquotes SHALL have multiple `> ` prefixes

#### Scenario: Print multi-paragraph blockquote
- **WHEN** printing blockquote with multiple paragraphs
- **THEN** blank lines within blockquote SHALL also have `>` prefix

### Requirement: Delta Section Printing
The system SHALL print Spectr delta sections with correct formatting.

#### Scenario: Print delta section header
- **WHEN** printing a NodeSection with DeltaType set
- **THEN** output SHALL be `## {DELTA_TYPE} Requirements\n`
- **AND** DELTA_TYPE SHALL be: ADDED, MODIFIED, REMOVED, or RENAMED

#### Scenario: Print WHEN/THEN/AND keywords
- **WHEN** printing a NodeListItem with Keyword field
- **THEN** keyword SHALL be printed bold: `- **WHEN** {content}`
- **AND** keyword SHALL be uppercase

