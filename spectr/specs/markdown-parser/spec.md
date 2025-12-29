# Markdown Parser Specification

## Requirements

### Requirement: Token-Based Lexer

The system SHALL provide a token-based lexer in `internal/markdown/lexer.go` that converts markdown source text into a stream of tokens with line and column positions.

#### Scenario: Lexer produces tokens for headers

- **WHEN** the lexer processes `## Section Name`
- **THEN** it SHALL produce a `TokenH2` token followed by text tokens
- **AND** each token SHALL include the line number and column offset
- **AND** the token value SHALL preserve the original text

#### Scenario: Lexer produces tokens for task checkboxes

- **WHEN** the lexer processes `- [ ] Task description`
- **THEN** it SHALL produce `TokenBullet`, `TokenCheckboxEmpty`, and text tokens
- **AND** when processing `- [x] Completed task`
- **THEN** it SHALL produce `TokenBullet`, `TokenCheckboxChecked`, and text tokens

#### Scenario: Lexer handles fenced code blocks

- **WHEN** the lexer encounters ``` or ~~~
- **THEN** it SHALL produce `TokenCodeFence` tokens
- **AND** content between fences SHALL be tokenized as `TokenText` without further parsing
- **AND** the language identifier after opening fence SHALL be captured

#### Scenario: Lexer handles inline formatting

- **WHEN** the lexer processes `**bold**` or `__bold__`
- **THEN** it SHALL produce `TokenBold` tokens wrapping the content
- **AND** when processing `*italic*` or `_italic_`
- **THEN** it SHALL produce `TokenItalic` tokens wrapping the content

#### Scenario: Lexer normalizes line endings

- **WHEN** the lexer processes content with CRLF line endings
- **THEN** it SHALL normalize them to LF internally
- **AND** line numbers SHALL still accurately reflect source positions

#### Scenario: Lexer tracks position accurately

- **WHEN** any token is produced
- **THEN** the token SHALL include the 1-based line number
- **AND** the token SHALL include the 1-based column number
- **AND** positions SHALL account for multi-byte UTF-8 characters

### Requirement: AST Parser

The system SHALL provide a parser in `internal/markdown/parser.go` that consumes tokens and produces an Abstract Syntax Tree (AST) with node types for document structure.

#### Scenario: Parser builds document structure

- **WHEN** the parser processes a stream of tokens
- **THEN** it SHALL produce a root `NodeDocument` containing child nodes
- **AND** H2 headers SHALL create `NodeSection` nodes
- **AND** content between sections SHALL be nested under the appropriate section

#### Scenario: Parser recognizes requirement headers

- **WHEN** the parser encounters `### Requirement: Name`
- **THEN** it SHALL create a `NodeRequirement` node with `Content = "Name"`
- **AND** subsequent content until the next H2/H3 SHALL be children of the requirement

#### Scenario: Parser recognizes scenario headers

- **WHEN** the parser encounters `#### Scenario: Description`
- **THEN** it SHALL create a `NodeScenario` node with `Content = "Description"`
- **AND** the scenario SHALL be nested under its parent requirement

#### Scenario: Parser handles nested lists

- **WHEN** the parser processes indented list items
- **THEN** it SHALL create nested `NodeList` and `NodeListItem` structures
- **AND** indentation level SHALL determine nesting depth

#### Scenario: Parser preserves source positions

- **WHEN** any AST node is created
- **THEN** the node SHALL include the line and column of its starting token
- **AND** positions SHALL be accessible for error reporting

### Requirement: Parse Error Reporting

The system SHALL provide structured parse errors with line numbers, column numbers, and contextual snippets.

#### Scenario: Error includes location

- **WHEN** a parse error occurs
- **THEN** the error SHALL include the 1-based line number
- **AND** the error SHALL include the 1-based column number
- **AND** the error message SHALL describe the problem

#### Scenario: Error includes context snippet

- **WHEN** a parse error is reported
- **THEN** the error SHALL include a snippet of the source text around the error
- **AND** the snippet SHALL highlight the error position with a caret (^)

#### Scenario: Multiple errors collected

- **WHEN** the parser encounters multiple errors
- **THEN** it SHALL collect all errors rather than stopping at the first
- **AND** parsing SHALL continue in a best-effort manner
- **AND** the result SHALL include both partial AST and error list

#### Scenario: Error formatting for display

- **WHEN** a ParseError is converted to string
- **THEN** the format SHALL be `line:column: message\n  snippet`
- **AND** the format SHALL be consistent across all error types

### Requirement: Spec File Parsing API

The system SHALL provide high-level API functions in `internal/markdown/api.go` for parsing Spectr spec files and extracting structured data.

#### Scenario: Parse complete spec file

- **WHEN** `ParseSpec(content)` is called with spec file content
- **THEN** it SHALL return a `*Spec` structure and any parse errors
- **AND** the Spec SHALL contain sections, requirements, and scenarios

#### Scenario: Extract sections from content

- **WHEN** `ExtractSections(content)` is called
- **THEN** it SHALL return a map of section name to `Section` struct
- **AND** each Section SHALL contain its content and child requirements

#### Scenario: Extract requirements from content

- **WHEN** `ExtractRequirements(content)` is called
- **THEN** it SHALL return a slice of `Requirement` structs
- **AND** each Requirement SHALL include name, content, and scenarios

#### Scenario: Find specific section

- **WHEN** `FindSection(content, "Requirements")` is called
- **THEN** it SHALL return the Section and `true` if found
- **AND** it SHALL return `nil, false` if the section does not exist

### Requirement: Delta File Parsing API

The system SHALL provide API functions for parsing delta spec files with ADDED, MODIFIED, REMOVED, and RENAMED sections.

#### Scenario: Parse delta spec file

- **WHEN** `ParseDelta(content)` is called with delta file content
- **THEN** it SHALL return a `*Delta` structure with categorized requirements
- **AND** Added, Modified, Removed, and Renamed fields SHALL be populated

#### Scenario: Extract delta section content

- **WHEN** `FindDeltaSection(content, "ADDED")` is called
- **THEN** it SHALL return the content between `## ADDED Requirements` and the next H2
- **AND** it SHALL return empty string if section not found

#### Scenario: Parse RENAMED entries

- **WHEN** a RENAMED section contains FROM/TO pairs
- **THEN** the parser SHALL extract both backtick-wrapped and plain formats
- **AND** each rename SHALL have From and To fields populated

### Requirement: Wikilink Parsing

The system SHALL support wikilink syntax for linking between specs and changes with optional display text and anchors.

#### Scenario: Parse simple wikilink

- **WHEN** the lexer processes `[[validation]]`
- **THEN** it SHALL produce a `TokenWikilink` with `Target = "validation"`
- **AND** `Display` SHALL be empty (defaults to target)

#### Scenario: Parse wikilink with display text

- **WHEN** the lexer processes `[[validation|Validation Spec]]`
- **THEN** it SHALL produce a `TokenWikilink` with `Target = "validation"` and `Display = "Validation Spec"`

#### Scenario: Parse wikilink with anchor

- **WHEN** the lexer processes `[[validation#Requirement: Spec File Validation]]`
- **THEN** it SHALL produce a `TokenWikilink` with `Target = "validation"` and `Anchor = "Requirement: Spec File Validation"`

#### Scenario: Parse wikilink to change

- **WHEN** the lexer processes `[[changes/replace-regex-with-parser]]`
- **THEN** it SHALL produce a `TokenWikilink` with `Target = "changes/replace-regex-with-parser"`

#### Scenario: Wikilink resolution to spec

- **WHEN** `ResolveWikilink("validation", projectRoot)` is called
- **THEN** it SHALL check for `spectr/specs/validation/spec.md`
- **AND** return the path and `true` if the file exists

#### Scenario: Wikilink resolution to change

- **WHEN** `ResolveWikilink("changes/my-change", projectRoot)` is called
- **THEN** it SHALL check for `spectr/changes/my-change/proposal.md`
- **AND** return the path and `true` if the file exists

### Requirement: Reference-Style Link Parsing

The system SHALL support reference-style links with link definitions collected in a first pass and resolved during parsing.

#### Scenario: Parse link definition

- **WHEN** the lexer processes `[ref]: https://example.com "Title"`
- **THEN** it SHALL produce a `TokenLinkDef` with label, URL, and optional title

#### Scenario: Resolve reference link

- **WHEN** the parser encounters `[text][ref]` and a definition `[ref]: url` exists
- **THEN** it SHALL create a `NodeLink` with the resolved URL
- **AND** the display text SHALL be "text"

#### Scenario: Unresolved reference link

- **WHEN** the parser encounters `[text][ref]` without a matching definition
- **THEN** it SHALL create a parse error indicating undefined reference
- **AND** the error SHALL include the line number and reference label

### Requirement: Compatibility Helpers

The system SHALL provide helper functions that maintain compatibility with common usage patterns from the old regex-based API.

#### Scenario: Match H3 requirement header

- **WHEN** `MatchRequirementHeader(line)` is called
- **THEN** it SHALL return `(name string, ok bool)` matching the old API
- **AND** it SHALL use the lexer/parser internally, not regex

#### Scenario: Match H4 scenario header

- **WHEN** `MatchScenarioHeader(line)` is called
- **THEN** it SHALL return `(name string, ok bool)` matching the old API
- **AND** it SHALL use the lexer/parser internally, not regex

#### Scenario: Check if line is H2 header

- **WHEN** `IsH2Header(line)` is called
- **THEN** it SHALL return `true` if the line starts with `##`
- **AND** it SHALL use simple string prefix check, not regex

#### Scenario: Match task checkbox

- **WHEN** `MatchTaskCheckbox(line)` is called
- **THEN** it SHALL return `(state rune, ok bool)` with 'x', 'X', or ' '
- **AND** it SHALL use the lexer internally, not regex
