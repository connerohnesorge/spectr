# Design: Blackfriday AST-Based Markdown Parsing

## Context

Spectr parses markdown specification files to extract:
- H2 sections (`## Requirements`, `## ADDED Requirements`, etc.)
- H3 requirement headers (`### Requirement: Name`)
- H4 scenario headers (`#### Scenario: Name`)
- Task checkboxes (`- [ ]`, `- [x]`)
- Requirement content blocks (everything between headers)

Current implementation uses line-by-line scanning with regex patterns. This works but is fragile and duplicated across multiple packages.

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
├── parser.go      # Core AST parsing and document structure
├── headers.go     # Header extraction (H1-H4)
├── sections.go    # Section content extraction
├── tasks.go       # Task checkbox parsing
└── types.go       # Shared types (Header, Section, Task)
```

### Decision 3: Preserve existing public APIs

**What**: Keep function signatures in `parsers`, `validation`, `archive` unchanged.

**Why**:
- Minimize blast radius
- No changes to cmd layer
- Existing tests remain valid

**Implementation**: Internal functions call new markdown package, transform results to existing types.

### Decision 4: AST walking pattern

**What**: Use recursive walker with callback pattern for node processing.

```go
type NodeVisitor func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus

func Walk(doc *blackfriday.Node, visitor NodeVisitor) {
    doc.Walk(visitor)
}
```

**Why**:
- Matches blackfriday's native API
- Allows selective processing (skip subtrees)
- Clean handling of enter/exit events for nested content

## Risks / Trade-offs

| Risk | Impact | Mitigation |
|------|--------|------------|
| Parsing differences | Medium - edge cases may break | Comprehensive comparison tests |
| Performance regression | Low - AST is faster than repeated regex | Benchmark before/after |
| Blackfriday quirks | Low - mature library | Pin to specific version |

## Migration Plan

1. Add blackfriday dependency
2. Implement `internal/markdown/` package with tests
3. Create comparison tests (regex vs AST output)
4. Replace internal implementations one file at a time
5. Remove unused regex patterns
6. Run full test suite

**Rollback**: Revert to commit before change; no database or external state involved.

## Open Questions

None - straightforward refactoring with clear scope.
