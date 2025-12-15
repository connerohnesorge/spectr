# Change: Replace Regex Parsing with Handbuilt Token-Based Parser

## Why
The current regex-based markdown parsing in `internal/regex/` is fragile, hard to extend, and provides limited error context. A handbuilt token-based lexer/parser will:
- Eliminate all regex patterns from the codebase for structural markdown parsing
- Provide precise line-based error reporting with column information
- Support a broader CommonMark subset for future extensibility
- Be easier to maintain, test, and debug than scattered regex patterns

## What Changes
- **NEW**: `internal/markdown/` package with token-based lexer and parser
- **NEW**: Lexer producing tokens (headers, lists, code blocks, emphasis, links, text)
- **NEW**: Parser consuming tokens to build structured AST
- **NEW**: Line/column error reporting with context snippets
- **REMOVED**: `internal/regex/` package entirely
- **MODIFIED**: All call sites updated to use new `internal/markdown/` API
- **MODIFIED**: `internal/parsers/` to use markdown package for parsing
- **MODIFIED**: `internal/validation/parser.go` to use markdown package
- **MODIFIED**: `internal/archive/spec_merger.go` to use markdown package
- **MODIFIED**: `cmd/accept.go` to use markdown package
- **MODIFIED**: `internal/git/platform.go` URL parsing (non-markdown regex retained or converted)

## Impact
- Affected specs: `validation` (regex consolidation requirements removed/modified)
- Affected code:
  - `internal/regex/` (deleted)
  - `internal/markdown/` (new)
  - `internal/parsers/requirement_parser.go`
  - `internal/parsers/delta_parser.go`
  - `internal/validation/parser.go`
  - `internal/validation/change_rules.go`
  - `internal/archive/spec_merger.go`
  - `cmd/accept.go`
- Breaking: API changes from regex helpers to markdown package functions

## Scope (CommonMark Subset + Extensions)
The parser will support a useful CommonMark subset plus Spectr extensions:
- **Headers**: H1-H6 (ATX style with `#`)
- **Lists**: Unordered (`-`, `*`, `+`), ordered (`1.`), task checkboxes (`- [ ]`, `- [x]`)
- **Code**: Fenced code blocks (``` and ~~~), inline code (`)
- **Emphasis**: Bold (`**`, `__`), italic (`*`, `_`), strikethrough (`~~`)
- **Links**: Inline links `[text](url)`, reference-style links `[text][ref]` with `[ref]: url` definitions
- **Wikilinks**: `[[spec-name]]` or `[[spec-name|display text]]` for linking to other specs/changes
- **Block elements**: Paragraphs, blockquotes (`>`)
- **Special Spectr patterns**: WHEN/THEN bullet points, requirement headers, scenario headers, delta sections

## Key Design Decisions
Based on user input:
1. **Full CommonMark subset** - Support headers, lists, code blocks, emphasis, links (not just minimal Spectr-specific patterns)
2. **Strict error mode** - Return line-based errors for malformed input; collect all errors rather than stopping at first
3. **Token-based lexer/parser** - Separate tokenization pass then parse tokens into AST
4. **Clean slate API** - New `internal/markdown/` package with improved API design (not drop-in replacement)

## Non-Goals
- Full CommonMark compliance (we support a useful subset, not 100% spec compliance)
- Tables (not needed for Spectr specs)
- HTML passthrough
- Setext-style headers (underlined with `===` or `---`)
- GFM extensions beyond task checkboxes and strikethrough
