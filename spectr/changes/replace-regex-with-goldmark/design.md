# Design: Replace Regex Parsing with Goldmark AST

## Context

Current parsing implementation uses line-by-line regex matching in three core files:
- `parsers.go`: Title extraction, task/delta/requirement counting
- `requirement_parser.go`: Requirement block extraction with scenarios
- `delta_parser.go`: Delta section parsing (ADDED/MODIFIED/REMOVED/RENAMED)

**Limitations of regex approach:**
- Fragile handling of code blocks (can match patterns inside fenced code)
- Poor handling of nested markdown structures
- Difficult to maintain and extend (complex regex patterns)
- No AST representation for richer analysis
- Line-by-line scanning doesn't understand document structure

**Why goldmark:**
- AST-based parsing provides proper CommonMark compliance
- Extensible architecture for custom walkers
- Robust handling of edge cases (code blocks, nested elements)
- Industry-standard library (used by Hugo, gokrazy, and other major projects)
- Source mapping allows line number tracking from byte offsets

**Scope:** Complete migration of all three parser files to goldmark-based implementation.

## Goals / Non-Goals

### Goals
- Replace all regex parsing with goldmark AST traversal
- Improve robustness for edge cases (code blocks, nested markdown, inline code)
- Make parsing more maintainable and extensible
- Preserve exact markdown text for archive merging workflow
- Maintain line number tracking for validation error messages
- Keep all existing tests passing

### Non-Goals
- Not changing spec file formats or conventions
- Not adding new markdown features or extensions
- Not optimizing for performance over correctness
- Not introducing external dependencies beyond goldmark

## Decisions

### AST Architecture

**Decision:** Use goldmark parser with custom AST walkers

**Approach:**
- Parse markdown to goldmark AST once per file
- Create reusable walker utilities for common patterns (find headings, extract sections)
- Build domain-specific parsers on top of walkers

**Rationale:**
- Single parse pass is more efficient than multiple regex scans
- AST walkers are composable and testable
- Separates markdown parsing from domain logic

### Line Number Tracking

**Decision:** Convert goldmark segment offsets to line/column/byte location using source text

**Implementation:**
```go
type SourceLocation struct {
    LineNumber int
    Column     int
    ByteOffset int
}

func segmentToLineNumber(source []byte, segment text.Segment) int {
    // Count newlines in source[0:segment.Start]
    return bytes.Count(source[0:segment.Start], []byte("\n")) + 1
}

func segmentToColumn(source []byte, segment text.Segment) int {
    // Find the last newline before segment.Start, then offset into the current line
    lastNewline := bytes.LastIndex(source[:segment.Start], []byte("\n"))
    return segment.Start - lastNewline
}

func segmentToLocation(source []byte, segment text.Segment) SourceLocation {
    return SourceLocation{
        LineNumber: segmentToLineNumber(source, segment),
        Column:     segmentToColumn(source, segment),
        ByteOffset: segment.Start,
    }
}
```

**Rationale:**
- Goldmark provides `text.Segment` with byte offsets
- Validation errors need human-readable line numbers
- One-time conversion cost is acceptable for error reporting
- Column and byte offsets improve precision for error pointers in validators

**Trade-off:** Small performance cost to calculate line numbers, but only needed for error cases.

### Custom Delta Section Parsing

**Decision:** Use AST walker + domain-specific logic (not goldmark extensions)

**Approach:**
1. Walk AST to find H2 nodes matching `## (ADDED|MODIFIED|REMOVED|RENAMED) Requirements`
2. Collect H3 requirement nodes under each delta section until next H2
3. Extract text content and build RequirementBlock structures

**Rationale:**
- Delta sections are domain-specific, not general markdown features
- Goldmark extensions add complexity without clear benefit
- AST walker with pattern matching is straightforward

**Alternative considered:** Custom goldmark parser extension - rejected as over-engineered.

### Rename Format Handling

**Decision:** Parse with AST walker, then post-process list items for FROM/TO

**Implementation:**
1. Find RENAMED section via H2 heading
2. Walk list items under that section
3. Post-process text content to extract FROM/TO pairs from inline code

**Format:**
```markdown
## RENAMED Requirements
- FROM: `### Requirement: Old Name`
- TO: `### Requirement: New Name`
```

**Rationale:**
- AST gives us list structure (avoids false matches in code blocks)
- Post-processing text content is simpler than custom AST nodes
- More robust than regex for nested structures

### Data Structure Changes

**BREAKING CHANGES** (internal only - no external consumers):

```go
// Before
type RequirementBlock struct {
    HeaderLine string
    Name       string
    Raw        string
}

// After
type RequirementBlock struct {
    HeaderLine string
    Name       string
    Raw        string          // Preserved for archive merging
    ASTNode    ast.Node        // NEW: Reference to AST node
    Location   SourceLocation  // NEW: Line/column for errors
}

type SourceLocation struct {
    LineNumber int
    Column     int
    ByteOffset int
}
```

**Backward compatibility:**
- `Raw` field preserved with exact markdown text
- Archive merging uses `Raw` content (no change to spec_merger.go logic)
- Public validation API unchanged

**Impact:** All internal consumers must be updated simultaneously (archive, validation, list, view).

## Migration Plan

### Phase 1: Foundation (Tasks 1.1-1.3)
**Deliverables:**
- Add `github.com/yuin/goldmark` to go.mod
- Create `internal/parsers/ast_utils.go` with utilities:
  - `ParseMarkdown(content []byte) ast.Node`
  - `FindHeading(node ast.Node, level int, text string) ast.Node`
  - `SegmentToLocation(source []byte, segment text.Segment) SourceLocation` (line, column, byte offset helpers)
  - `ExtractTextContent(node ast.Node) string`

**Testing:** Unit tests for each utility function.

### Phase 2: Simple Parsers (Tasks 2.1-2.6)
**Deliverables:**
- Migrate `ExtractTitle()` to AST-based implementation
- Migrate `CountRequirements()`, `CountDeltas()`, `CountTasks()` to AST

**Testing:**
- Run parallel tests (regex vs goldmark) to validate equivalence
- Use existing test fixtures in parsers_test.go

### Phase 3: Requirements (Tasks 3.1-3.5)
**Deliverables:**
- Migrate `ParseRequirements()` to AST-based (update RequirementBlock)
- Migrate `ParseScenarios()` to AST-based

**Testing:**
- requirement_parser_test.go must pass
- Verify line numbers in RequirementBlock.Location

### Phase 4: Delta Sections (Tasks 4.1-4.8)
**Deliverables:**
- Implement AST-based delta section detection
- Parse ADDED/MODIFIED/REMOVED sections
- Parse RENAMED section with list post-processing

**Testing:**
- delta_parser_test.go must pass
- Add tests for edge cases (code blocks with "## ADDED" inside)

### Phase 5: Integration (Tasks 5.1-5.10)
**Deliverables:**
- Update consumers: archive/spec_merger.go, validation/*.go, list/lister.go, view/dashboard.go
- Update all test files
- Performance benchmarking (before/after comparison)
- Remove old regex code

**Testing:**
- Full integration test suite
- Validate no regressions in archive merging behavior

## Risks / Trade-offs

### Risk: Performance Regression
**Details:** AST parsing may be slower than simple regex for small files

**Mitigation:**
- Benchmark on real spec files (typical size 100-500 lines)
- Profile hot paths
- Cache AST if same file parsed multiple times

**Trade-off:** Accept slight slowdown for correctness and maintainability

### Risk: Line Number Calculation Complexity
**Details:** Converting byte offsets to line numbers requires counting newlines

**Mitigation:**
- Create utility function with comprehensive tests
- Document edge cases (CRLF vs LF, empty files)
- Cache line number calculations where needed

**Trade-off:** Small runtime cost for better error messages

### Risk: Breaking Internal APIs
**Details:** RequirementBlock and DeltaPlan structures change

**Mitigation:**
- All consumers are internal (no external packages)
- Update all consumers in same change
- Use compiler to catch breaking changes

**Trade-off:** One-time migration effort vs long-term maintainability

### Risk: Archive Merging Behavior Change
**Details:** Spec merging depends on exact text matching

**Mitigation:**
- Preserve `Raw` field with exact markdown content
- Add integration tests for archive workflow
- Validate merged specs match expected format

**Trade-off:** Must maintain two representations (AST + Raw text)

### Risk: Rename Format Parsing Complexity
**Details:** FROM/TO parsing with inline code backticks is tricky

**Mitigation:**
- Use AST to find list structure (avoids false positives)
- Post-process text content for extraction
- Add comprehensive test cases with edge cases

**Trade-off:** Two-stage parsing (AST + text) vs complex regex

## Open Questions

1. **Should we add goldmark extensions for custom node types?**
   - Current decision: No, use post-processing
   - Revisit if we need more complex custom syntax

2. **Do we need feature flags for gradual rollout?**
   - Current decision: No, all-at-once migration
   - All consumers are internal, can coordinate updates

3. **Should line number calculation be cached?**
   - Current decision: Calculate on-demand for errors
   - Profile first, optimize if needed

4. **How to handle malformed markdown that goldmark rejects?**
   - Current decision: Return parse errors to user
   - Goldmark is permissive (CommonMark spec), unlikely to fail on valid files
