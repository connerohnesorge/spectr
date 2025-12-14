# Design: Blackfriday AST-Based Markdown Parsing

## Context

Spectr parses markdown specification files to extract:
- H2 sections (`## Requirements`, `## ADDED Requirements`, etc.)
- H3 requirement headers (`### Requirement: Name`)
- H4 scenario headers (`#### Scenario: Name`)
- Task checkboxes (`- [ ]`, `- [x]`)
- Requirement content blocks (everything between headers)

Following the `consolidate-regex-patterns` change, all markdown-related regex patterns are now consolidated in `internal/regex/`:
- `headers.go`: H2, H3, H4 header patterns and matchers
- `tasks.go`: Task checkbox and numbered task patterns
- `renames.go`: RENAMED section FROM/TO patterns
- `sections.go`: Section content extraction helpers

While consolidation addressed duplication and pre-compilation, the regex approach still lacks structural understanding of markdown.

## Goals

- Single markdown parsing implementation in `internal/markdown/`
- Type-safe AST traversal instead of string matching
- Preserve all existing parsing behavior
- Improve error messages with source locations

## Non-Goals

- Changing spec file format
- Adding new markdown features
- Replacing non-markdown regex (git URLs, whitespace normalization)

## Decisions

### Decision 1: Use blackfriday v2 with AST walking

**What**: Use `github.com/russross/blackfriday/v2` with custom AST walker.

**Why**: 
- Mature, widely-used library
- Provides raw AST access via `blackfriday.Node`
- Fast parsing without HTML rendering overhead
- User's explicit preference

**Alternatives considered**:
- goldmark: More compliant but heavier, overkill for our header-focused parsing
- Custom parser: Too much effort for well-solved problem

### Decision 2: Create `internal/markdown/` package

**What**: New package with focused API for spec parsing needs.

**Why**:
- Clean separation from business logic
- Testable in isolation
- Single place to handle blackfriday specifics

**Structure**:
```
internal/markdown/
├── parser.go      # Core AST parsing, ParseDocument() function
├── headers.go     # Header extraction logic
├── sections.go    # Section content extraction logic
├── tasks.go       # Task checkbox parsing logic
└── types.go       # Shared types (Document, Header, Section, Task)
```

**Types** (in types.go):
```go
// Document represents a parsed markdown document with all extracted elements.
type Document struct {
    Headers  []Header           // All headers H1-H4
    Sections map[string]Section // Sections keyed by header text
    Tasks    []Task             // All task checkboxes (hierarchical)
}

// Header represents a markdown heading with source location.
type Header struct {
    Level int    // 1-6
    Text  string // Header text content
    Line  int    // Source line number (1-indexed)
}

// Section represents content between headers.
type Section struct {
    Header  Header // The section's header
    Content string // Raw markdown text between this header and next
}

// Task represents a checkbox item with hierarchy support.
type Task struct {
    Line     string // Full original line text
    Checked  bool   // true if [x] or [X]
    LineNum  int    // Source line number (1-indexed)
    Children []Task // Nested task items
}
```

### Decision 3: Preserve existing public APIs

**What**: Keep function signatures in `parsers`, `validation`, `archive` unchanged.

**Why**:
- Minimize blast radius
- No changes to cmd layer
- Existing tests remain valid

**Implementation**: Internal functions call new markdown package, transform results to existing types.

### Decision 4: All-in-one parsing API

**What**: Single `ParseDocument()` function that parses once and returns everything.

```go
// ParseDocument parses markdown content and extracts all structural elements.
// Returns error for invalid input (empty content, binary data).
// Blackfriday internals are never exposed; all data is in package-defined types.
func ParseDocument(content []byte) (*Document, error)
```

**Why**:
- Parse once, never reparse - architecture enforces this
- Callers use only what they need from the returned Document
- No need for multiple parsing functions or caching logic
- Blackfriday internals strictly hidden from callers

**Internal implementation** uses recursive walker with callback pattern:
```go
// Internal only - not exported
type nodeVisitor func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus
```

### Decision 5: Error handling with specterrs

**What**: Add markdown-specific error types to `internal/specterrs/markdown.go`.

```go
// MarkdownParseError indicates markdown content failed to parse.
type MarkdownParseError struct {
    Path string // File path if known, empty otherwise
    Line int    // Line number if known, 0 otherwise
    Err  error  // Underlying error
}

func (e *MarkdownParseError) Error() string {
    if e.Line > 0 {
        return fmt.Sprintf("failed to parse markdown %s at line %d: %v", e.Path, e.Line, e.Err)
    }
    if e.Path != "" {
        return fmt.Sprintf("failed to parse markdown %s: %v", e.Path, e.Err)
    }
    return fmt.Sprintf("failed to parse markdown: %v", e.Err)
}

func (e *MarkdownParseError) Unwrap() error { return e.Err }

// EmptyContentError indicates empty or whitespace-only content was provided.
type EmptyContentError struct {
    Path string
}

func (e *EmptyContentError) Error() string {
    if e.Path != "" {
        return fmt.Sprintf("markdown file is empty: %s", e.Path)
    }
    return "markdown content is empty"
}

// BinaryContentError indicates binary (non-text) content was provided.
type BinaryContentError struct {
    Path string
}

func (e *BinaryContentError) Error() string {
    if e.Path != "" {
        return fmt.Sprintf("file appears to be binary, not markdown: %s", e.Path)
    }
    return "content appears to be binary, not markdown"
}
```

**Why**:
- Consistent with existing specterrs pattern
- Structured errors allow programmatic handling
- Line numbers included where available

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Parsing differences | Medium - edge cases may break | Comprehensive comparison tests |
| Performance regression | Low - AST is faster than repeated regex | Benchmark before/after |
| Blackfriday quirks | Low - mature library | Pin to specific version |

## Migration Plan

1. Add blackfriday dependency to go.mod
2. Add markdown error types to `internal/specterrs/markdown.go`
3. Implement `internal/markdown/` package with types and ParseDocument()
4. Write unit tests for markdown package
5. Create comparison tests with regex patterns **embedded in test file** (copied from `internal/regex/`), comparing regex vs AST output - **keep permanently as regression suite**
6. Update consumers one file at a time to use markdown package:
   - `internal/parsers/parsers.go`
   - `internal/parsers/requirement_parser.go`
   - `internal/parsers/delta_parser.go`
   - `internal/validation/parser.go`
   - `internal/archive/spec_merger.go`
   - `cmd/accept.go`
7. Remove `internal/regex/` package entirely after all consumers migrated (comparison tests retain embedded patterns)
8. Run full test suite

**Rollback**: Revert to commit before change; no database or external state involved.

## Resolved Questions

All design questions have been resolved:

| Question | Resolution |
|----------|------------|
| Malformed markdown handling | Accept blackfriday's interpretation |
| AST caching strategy | Parse once, never reparse |
| Line numbers | Hard requirement - all types include line numbers |
| Comparison tests | Keep permanently as regression suite |
| API shape | All-in-one ParseDocument() |
| Package structure | Multiple files as specified |
| AST exposure | Strictly hide blackfriday internals |
| Section content format | Raw markdown text |
| Task content format | Full line text preserved |
| Parse error handling | Return (nil, error) using specterrs types |
| Nested tasks | Preserve hierarchy with Children []Task |
