# Change: Replace Regex-Based Markdown Parsing with Blackfriday AST

## Why

The current markdown parsing implementation uses 28+ regex patterns scattered across 6 files (`internal/parsers/*.go`, `internal/validation/parser.go`, `internal/archive/spec_merger.go`, `cmd/accept.go`). This approach has several issues:

1. **Fragility**: Regex patterns are duplicated and inconsistent (e.g., two different RENAMED parsing patterns exist)
2. **No caching**: Patterns are recompiled on every function call
3. **Limited structure awareness**: Regex cannot understand markdown nesting or context
4. **Maintenance burden**: Adding new spec formats requires updating multiple regex patterns

Replacing with blackfriday's AST-based parsing provides:
- Single source of truth for markdown structure
- Proper handling of nested elements
- Better error messages with line numbers
- Easier extensibility for future spec formats
- Single parse of markdown, no need to reparse then entire markdown document

## What Changes

- **NEW**: `internal/markdown/` package with AST-based parsing using blackfriday v2
- **MODIFIED**: `internal/parsers/` to use new markdown package instead of regex
- **MODIFIED**: `internal/validation/parser.go` to use new markdown package
- **MODIFIED**: `internal/archive/spec_merger.go` to use new markdown package
- **MODIFIED**: `cmd/accept.go` to use new markdown package for task parsing
- **MODIFIED**: `go.mod` to add `github.com/russross/blackfriday/v2` dependency

## Impact

- Affected specs: `validation` (parsing behavior documented there)
- Affected code:
  - `internal/parsers/parsers.go` (task counting, delta counting, requirement counting)
  - `internal/parsers/delta_parser.go` (delta section extraction)
  - `internal/parsers/requirement_parser.go` (requirement block parsing)
  - `internal/validation/parser.go` (section, requirement, scenario extraction)
  - `internal/archive/spec_merger.go` (spec reconstruction)
  - `cmd/accept.go` (task markdown parsing)
- **NOT affected**: `internal/git/platform.go` (URL parsing stays regex)
- **NOT affected**: Simple utility patterns (`\s+`, `\n{3,}`, `shall|must`)

## Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Malformed markdown | Accept blackfriday's interpretation | Trust the library; no need for regex equivalence on edge cases |
| AST caching | Parse once, never reparse | Architecture enforces single parse per file per operation |
| Line numbers | Hard requirement | Types MUST include line numbers for error messages |
| Comparison tests | Keep permanently | Maintain as regression suite for future changes |
| API shape | All-in-one `ParseDocument()` | Single parse returns everything; callers use what they need |
| Package structure | Multiple files | parser.go, headers.go, sections.go, tasks.go, types.go |
| AST exposure | Strictly hide internals | No access to raw blackfriday nodes |
| Section content | Raw markdown text | `Section.Content` is string with original markdown |
| Task content | Full line text | `Task.Line` contains complete original line |
| Parse errors | Return error | `ParseDocument` returns `(nil, error)` on invalid input |
| Nested tasks | Preserve hierarchy | `Task` has `Children []Task` for nested items |
| Error types | Use specterrs | Add markdown-specific errors to `internal/specterrs/` |

## Risk Assessment

- **Low risk**: All existing tests continue to pass
- **Medium risk**: Subtle parsing differences in edge cases (malformed markdown)
- **Mitigation**: Comprehensive test coverage comparing old vs new parser output; comparison tests kept permanently
