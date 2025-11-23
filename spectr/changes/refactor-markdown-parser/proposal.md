# Change: Refactor Markdown Parser to Lexer/Parser Architecture

## Why

Current markdown parsing relies on regex-based line-by-line scanning which is brittle for edge cases and difficult to maintain. This approach fails to properly handle:

- **Code blocks** containing markdown-like syntax (e.g., `###` inside triple backticks)
- **Nested structures** like lists containing code blocks or blockquotes with requirements
- **Standard markdown** features like indented code blocks, HTML comments, and escaped characters
- **Multi-line patterns** that don't fit the line-by-line model

The archived validate-command design.md acknowledged this limitation with "Add proper parser library if edge cases emerge during implementation" - those edge cases have emerged. A proper lexer/parser architecture, modeled after the Go compiler's text/scanner and go/scanner packages, will provide robust, maintainable markdown parsing that correctly handles the full markdown syntax while remaining performant.

## What Changes

- **NEW**: Create `internal/parsers/markdown` package with lexer/parser architecture
  - Lexer: Tokenizes markdown into meaningful tokens (HEADING, CODE_FENCE, TEXT, LIST_ITEM, etc.)
  - Parser: Builds abstract syntax tree (AST) from tokens
  - Extractor: Walks AST to extract Spectr structures (requirements, scenarios, sections)
- **MODIFY**: Update `internal/parsers` to use new lexer/parser instead of regex
  - Replace regex patterns with AST-based extraction
  - Maintain existing public APIs for backward compatibility
  - Improve error reporting with line/column information
- **MODIFY**: Update `internal/validation/parser.go` to leverage new parser
- **IMPROVE**: Add comprehensive test suite with markdown edge cases
  - Code blocks containing headers
  - Nested lists with code blocks
  - Indented code blocks
  - HTML comments and escaped characters
  - Mixed markdown features

## Impact

### Affected Capabilities
- **validation** - Uses parsing to extract and validate requirements (MODIFIED)
- **archive-workflow** - Uses parsing to merge delta specs into base specs (MODIFIED)
- **markdown-parsing** - New capability defining the lexer/parser architecture (ADDED)

### Affected Code
- `internal/parsers/parsers.go` - Refactored to use new parser
- `internal/parsers/requirement_parser.go` - Refactored to use new parser
- `internal/parsers/delta_parser.go` - Refactored to use new parser
- `internal/validation/parser.go` - Updated to use new parser
- `internal/archive/spec_merger.go` - May need updates for improved parsing
- **NEW**: `internal/parsers/markdown/` - Lexer, parser, AST, token definitions

### Breaking Changes
- **None** - Public APIs remain unchanged
- Improved parsing may reveal previously undetected malformed markdown (validation improvements, not breaking changes)

### Migration Path
- Existing specs and changes will be parsed more accurately
- Invalid markdown that previously slipped through may now be caught (this is a feature, not a bug)
- No action required from users - changes are internal to implementation

### Performance Impact
- **Expected**: Comparable or better performance due to single-pass lexing vs multiple regex scans
- **Measured**: Will benchmark before/after to ensure no regression
- Target: <10% performance change for typical spec files (<100 requirements)

### Testing Impact
- Expanded test coverage for markdown edge cases
- All existing tests must continue to pass
- New tests added for previously unsupported edge cases
