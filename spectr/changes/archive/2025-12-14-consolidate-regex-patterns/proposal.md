# Change: Consolidate Regex-Based Markdown Parsing into Single Package

## Why

The codebase has 28+ regex patterns scattered across 6 files for parsing markdown structure:

- `internal/parsers/parsers.go` (3 patterns)
- `internal/parsers/delta_parser.go` (10 patterns)
- `internal/parsers/requirement_parser.go` (3 patterns)
- `internal/validation/parser.go` (5 patterns)
- `internal/archive/spec_merger.go` (5 patterns)
- `cmd/accept.go` (2 patterns)

This causes several issues:

1. **Duplication**: Same patterns (e.g., `^###\s+Requirement:`) appear in 4+ files
2. **Inconsistency**: Subtle variations exist (some use `(?m)`, some don't)
3. **No caching**: Patterns recompiled on every function call
4. **Testing burden**: Same behavior tested in multiple places
5. **Migration friction**: The upcoming `replace-regex-with-blackfriday` change must touch all 6 files

Consolidating regex patterns into `internal/regex/` will:

- Eliminate duplication immediately
- Pre-compile patterns once at package init
- Provide a single comparison point for blackfriday migration
- Improve testability with one canonical implementation

## What Changes

- **NEW**: `internal/regex/` package with pre-compiled patterns and matching functions
- **MODIFIED**: `internal/parsers/parsers.go` to use regex package
- **MODIFIED**: `internal/parsers/delta_parser.go` to use regex package
- **MODIFIED**: `internal/parsers/requirement_parser.go` to use regex package
- **MODIFIED**: `internal/validation/parser.go` to use regex package
- **MODIFIED**: `internal/archive/spec_merger.go` to use regex package
- **MODIFIED**: `cmd/accept.go` to use regex package
- **NOT affected**: `internal/git/platform.go` (URL parsing, not markdown)
- **NOT affected**: Simple utility patterns (`\s+`, `\n{3,}`, `shall|must`) remain inline

## Impact

- Affected specs: `validation` (parsing behavior documented there)
- Affected code:
  - `internal/parsers/*.go` - All parsers use consolidated package
  - `internal/validation/parser.go` - Uses consolidated package
  - `internal/archive/spec_merger.go` - Uses consolidated package
  - `cmd/accept.go` - Uses consolidated package

## Key Design Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Package name | `internal/regex/` | Clear purpose; not `markdown` to avoid confusion with future blackfriday package |
| Pattern compilation | Package-level `var` with `regexp.MustCompile` | Single compilation at init; panic on invalid patterns is acceptable |
| Function signatures | Return `(matches, found bool)` | Consistent API; no nil slice ambiguity |
| Section extraction | Keep line-by-line scanning in callers | Only patterns consolidated, not scanning logic |
| Utility patterns | Leave inline | `\s+`, `\n{3,}`, `shall|must` are simple and context-specific |

## Risk Assessment

- **Low risk**: Pure refactor; all existing tests must pass
- **Zero behavior change**: Output must be byte-for-byte identical
- **Mitigation**: Run full test suite after each file migration

## Relationship to Other Changes

This change is a **prerequisite** for `replace-regex-with-blackfriday`. After this change:

1. Blackfriday migration replaces `internal/regex/` with `internal/markdown/`
2. Comparison tests can compare `regex.Match*()` vs `markdown.Parse*()`
3. Cleaner diff: one package replacement instead of 6 file rewrites
