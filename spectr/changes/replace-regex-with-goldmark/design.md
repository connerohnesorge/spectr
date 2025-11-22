# Design: Goldmark AST-Based Parsing

## Context

Spectr currently uses regex-based line scanning for all markdown parsing:
- Title extraction: scans for first `# ` line
- Requirement parsing: regex for `### Requirement:`, accumulates lines until next `##`
- Scenario extraction: regex for `#### Scenario:`
- Delta sections: regex for `## ADDED|MODIFIED|REMOVED|RENAMED Requirements`
- Rename pairs: regex for `` - FROM: `### Requirement: X` `` format

This approach works for well-formed documents but fails on edge cases and provides poor error context.

**Stakeholders**: All users relying on validation, archiving, and listing functionality.

**Constraints**:
- Must maintain backward compatibility with existing spec format
- Must preserve all existing public APIs (zero breaking changes)
- Must handle all current test cases plus new edge cases
- Performance must not regress (ideally improve)

## Goals / Non-Goals

**Goals**:
1. Replace all regex parsing with goldmark AST traversal
2. Improve error reporting with precise line:column locations
3. Handle CommonMark edge cases (code blocks, nested structures)
4. Single-pass parsing for better performance
5. Extensible architecture for future spec format enhancements

**Non-Goals**:
- Changing spec file format or syntax
- Breaking existing public APIs
- Implementing custom markdown extensions (use standard CommonMark)
- Supporting non-CommonMark markdown flavors

## Decisions

### Decision 1: Use Goldmark as AST Provider

**Why Goldmark**:
- Official Go port of commonmark.js, used by Hugo and many Go projects
- Complete AST with source position tracking
- Extensible via renderer and parser plugins
- Excellent test coverage and active maintenance
- Zero dependencies beyond stdlib

**Alternatives Considered**:
- `github.com/gomarkdown/markdown`: Less precise AST, weaker source position tracking
- `github.com/russross/blackfriday`: Deprecated, no longer maintained
- Custom parser: High effort, would require extensive testing

**Decision**: Use `github.com/yuin/goldmark v1.7.0+`

### Decision 2: AST Walker Pattern

Implement a generic AST walker that visits nodes and executes callbacks:

```go
// AST walker signature
type NodeVisitor func(n ast.Node, entering bool) ast.WalkStatus

// Example usage
func findH1(doc ast.Node) string {
    var title string
    ast.Walk(doc, func(n ast.Node, entering bool) ast.WalkStatus {
        if heading, ok := n.(*ast.Heading); ok && entering && heading.Level == 1 {
            title = extractText(heading)
            return ast.WalkStop
        }
        return ast.WalkContinue
    })
    return title
}
```

**Why**: Goldmark's `ast.Walk` provides idiomatic traversal. Custom walkers add no value.

### Decision 3: Source Position Strategy

Enhance data structures with source position metadata:

```go
type RequirementBlock struct {
    HeaderLine string   // Preserved for compatibility
    Name       string
    Raw        string
    SourcePos  *SourcePosition  // NEW: AST-derived position
}

type SourcePosition struct {
    StartLine   int
    StartColumn int
    EndLine     int
    EndColumn   int
}
```

**Why**: Enables validation errors like `spec.md:42:5: Requirement missing scenario` instead of vague messages.

**Compatibility**: New field is optional, existing code unaffected.

### Decision 4: Delta Section Extraction

Use AST traversal to find `## ADDED Requirements` heading, then collect all nodes until next `##` heading:

```go
func extractSectionByHeading(doc ast.Node, targetHeading string) []ast.Node {
    var nodes []ast.Node
    inSection := false

    ast.Walk(doc, func(n ast.Node, entering bool) ast.WalkStatus {
        if h, ok := n.(*ast.Heading); ok && entering && h.Level == 2 {
            text := extractText(h)
            if strings.Contains(text, targetHeading) {
                inSection = true
                return ast.WalkContinue
            }
            if inSection {
                inSection = false
                return ast.WalkStop
            }
        }
        if inSection && entering {
            nodes = append(nodes, n)
        }
        return ast.WalkContinue
    })

    return nodes
}
```

**Why**: More robust than regex, correctly handles code blocks and nested structures.

### Decision 5: Rename Parsing Strategy

RENAMED sections use non-standard format (bullet lists with inline code):

```markdown
## RENAMED Requirements
- FROM: `### Requirement: Old Name`
- TO: `### Requirement: New Name`
```

**Strategy**: Use AST to extract list items, then apply targeted regex on list item text:

```go
func parseRenamedSection(doc ast.Node) []RenameOp {
    sectionNodes := extractSectionByHeading(doc, "RENAMED Requirements")

    var renamed []RenameOp
    var currentFrom string

    for _, node := range sectionNodes {
        if list, ok := node.(*ast.List); ok {
            // Iterate list items
            for item := list.FirstChild(); item != nil; item = item.NextSibling() {
                text := extractText(item)
                if from := parseFromLine(text); from != "" {
                    currentFrom = from
                } else if to := parseToLine(text); to != "" && currentFrom != "" {
                    renamed = append(renamed, RenameOp{From: currentFrom, To: to})
                    currentFrom = ""
                }
            }
        }
    }

    return renamed
}

func parseFromLine(text string) string {
    // Targeted regex on clean text (no markdown noise)
    re := regexp.MustCompile(`^FROM:\s*` + "`" + `###\s+Requirement:\s*(.+?)` + "`" + `$`)
    matches := re.FindStringSubmatch(strings.TrimSpace(text))
    if len(matches) > 1 {
        return matches[1]
    }
    return ""
}
```

**Why**: Hybrid approach - AST for structure, regex for simple pattern extraction within clean text.

### Decision 6: Line Number Mapping

Goldmark provides `ast.Node.Lines()` which returns `text.Segments` with byte offsets. Convert to line:column:

```go
func getSourcePosition(node ast.Node, source []byte) *SourcePosition {
    if node.Lines().Len() == 0 {
        return nil
    }

    seg := node.Lines().At(0)
    startLine := bytes.Count(source[:seg.Start], []byte("\n")) + 1
    startCol := seg.Start - bytes.LastIndex(source[:seg.Start], []byte("\n"))

    lastSeg := node.Lines().At(node.Lines().Len() - 1)
    endLine := bytes.Count(source[:lastSeg.Stop], []byte("\n")) + 1
    endCol := lastSeg.Stop - bytes.LastIndex(source[:lastSeg.Stop], []byte("\n"))

    return &SourcePosition{
        StartLine:   startLine,
        StartColumn: startCol,
        EndLine:     endLine,
        EndColumn:   endCol,
    }
}
```

**Why**: Accurate source positions for error reporting, consistent with standard editor navigation.

### Decision 7: Backward Compatibility Strategy

**Public API Preservation**:
- Keep function signatures unchanged (add `*AST` suffix to new implementations initially)
- Run parallel tests: `regex_output == ast_output` for all existing test cases
- After validation, replace old implementations, remove `AST` suffix
- Deprecation: Remove old regex functions entirely (not public API, internal only)

**Data Structure Compatibility**:
- Add optional fields (`SourcePos *SourcePosition`) that default to nil
- Existing consumers ignore new fields
- New consumers can opt-in to enhanced metadata

### Decision 8: Migration Phases

**Phase 1: Foundation** (tasks 1.x)
- Add goldmark dependency
- Create AST utility functions
- Establish testing patterns

**Phase 2: Simple Replacements** (tasks 2.x)
- Replace title extraction, simple counters
- Validate equivalence with existing tests

**Phase 3: Complex Parsing** (tasks 3.x-5.x)
- Migrate requirement and delta parsing
- Add source position tracking
- Handle rename edge cases

**Phase 4: Integration** (tasks 6.x-7.x)
- Update all consumers
- Remove old implementations
- Comprehensive validation

**Why Phased**: Reduces risk, allows incremental validation, easier to debug issues.

## Risks / Trade-offs

### Risk 1: Performance Regression
**Mitigation**: Benchmark before/after. Goldmark's single-pass parsing should be faster than multiple regex scans, but measure to confirm.

### Risk 2: Breaking Changes in Goldmark
**Mitigation**: Pin to specific version (`v1.7.x`), test upgrades before adopting.

### Risk 3: Edge Cases Not Covered by Tests
**Mitigation**: Add comprehensive edge case tests (malformed markdown, unusual whitespace, code blocks with fake headers).

### Risk 4: Migration Bugs
**Mitigation**: Parallel testing (regex vs AST), phased rollout, extensive integration testing before removing old code.

### Trade-off: Added Dependency
**Accept**: Goldmark is well-maintained, widely used, and critical for correct markdown parsing. Benefits outweigh dependency cost.

## Migration Plan

### Step 1: Add Goldmark Dependency
```bash
go get github.com/yuin/goldmark@v1.7.8
```

### Step 2: Implement AST Utilities
Create `internal/parsers/ast_utils.go` with:
- `parseMarkdownToAST(filePath string) (ast.Node, []byte, error)`
- `extractTextContent(node ast.Node) string`
- `findHeadingByLevel(doc ast.Node, level int) *ast.Heading`
- `extractSectionByHeading(doc ast.Node, heading string) []ast.Node`
- `getSourcePosition(node ast.Node, source []byte) *SourcePosition`

### Step 3: Implement New Parsers
For each regex parser, create AST equivalent:
- `ExtractTitleAST()` -> replaces `ExtractTitle()`
- `ParseRequirementsAST()` -> replaces `ParseRequirements()`
- `ParseDeltaSpecAST()` -> replaces `ParseDeltaSpec()`

### Step 4: Parallel Testing
```go
func TestExtractTitle_Equivalence(t *testing.T) {
    tests := []string{"testdata/spec1.md", "testdata/spec2.md", ...}
    for _, path := range tests {
        regexTitle, _ := ExtractTitle(path)
        astTitle, _ := ExtractTitleAST(path)
        assert.Equal(t, regexTitle, astTitle)
    }
}
```

### Step 5: Integration Testing
Run full test suite:
```bash
go test ./...
go run main.go validate --all --strict
```

### Step 6: Switchover
Replace old function calls with new implementations:
```bash
# Example
sed -i 's/ExtractTitle(/ExtractTitleAST(/g' internal/**/*.go
```

### Step 7: Cleanup
Remove old regex-based implementations, rename `*AST` functions to original names.

### Rollback Plan
If issues arise:
1. Git revert to pre-migration commit
2. Incremental rollback: Keep new functions, revert consumer updates
3. Debug in isolation using parallel tests

## Open Questions

1. **Should we expose AST in public API for advanced use cases?**
   - **Leaning No**: Keep AST as implementation detail, expose only parsed data structures
   - Can revisit if consumers need direct AST access

2. **How to handle non-CommonMark extensions (if needed in future)?**
   - **Use goldmark extensions**: Well-supported pattern
   - Example: Custom renderer for spec-specific syntax

3. **Should source positions be included in JSON output?**
   - **Leaning Yes**: Helpful for tooling (editors, LSPs)
   - Make it optional via flag (`--include-positions`)

## Success Criteria

1. ✅ All existing tests pass with AST-based parsing
2. ✅ New edge case tests pass (code blocks, nested structures)
3. ✅ `spectr validate --all --strict` passes on entire project
4. ✅ Performance benchmarks show no regression (ideally improvement)
5. ✅ Error messages include precise line:column references
6. ✅ Zero breaking changes to public APIs
7. ✅ All regex-based parsing code removed from codebase
