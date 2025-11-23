# Design: Markdown Lexer/Parser Architecture

## Context

Spectr parses markdown files to extract structured information (requirements, scenarios, sections, deltas). The current implementation uses regex-based line-by-line scanning across three files totaling ~1400 lines. This approach has proven fragile for edge cases involving code blocks, nested structures, and standard markdown features.

### Current Problems

1. **Code block blindness**: Regex matches `###` inside triple backticks as headers
2. **No context awareness**: Line-by-line scanning can't distinguish between different markdown contexts
3. **Brittle patterns**: Each new edge case requires careful regex tuning
4. **Maintenance burden**: Three separate parsing implementations with duplicated logic
5. **Poor error messages**: Can't report accurate line/column positions for malformed markdown

### Inspiration

The Go compiler's lexer/parser architecture (inspired by Rob Pike's lexer talk and used in `text/template`, `text/scanner`) provides a proven pattern:
- **Lexer**: Converts character stream → token stream
- **Parser**: Converts token stream → AST
- **Extractor**: Walks AST to extract domain structures

### Constraints

- Must handle all current markdown patterns (headings, lists, code blocks, scenarios)
- Must maintain backward compatibility with existing APIs
- Must be performant (no worse than 10% regression)
- Must follow Go idioms (interfaces, error handling, testing)
- Should not require full markdown spec compliance - only what Spectr needs

## Goals / Non-Goals

### Goals

- **Correct parsing** of markdown with code blocks, nested structures, and edge cases
- **Maintainable code** with clear separation of concerns (lexer, parser, extractor)
- **Better error messages** with accurate line/column positions
- **Test coverage** for all markdown edge cases that currently fail
- **API compatibility** - existing code using parsers package continues to work
- **Performance** - comparable or better than current regex approach

### Non-Goals

- **Full markdown spec compliance** - only implement what Spectr needs (headings, code blocks, lists, text)
- **Markdown rendering** - we're extracting structure, not rendering HTML
- **Markdown editing** - read-only parsing, no AST modification
- **External library** - custom implementation tailored to Spectr's needs
- **Streaming parsing** - file sizes are small (<1MB), load entire file into memory

## Architecture

### Three-Stage Pipeline

```
Markdown File → Lexer → Token Stream → Parser → AST → Extractor → Spectr Structures
```

#### Stage 1: Lexer (Tokenization)

**Purpose**: Convert character stream into token stream

**Token Types** (inspired by markdown spec, tailored to Spectr needs):
```go
type TokenType int

const (
    // Structure tokens
    HEADING_1       // # Title
    HEADING_2       // ## Section
    HEADING_3       // ### Requirement:
    HEADING_4       // #### Scenario:

    // Block tokens
    CODE_FENCE_START    // ```language
    CODE_FENCE_END      // ```
    INDENT_CODE_START   // 4+ spaces
    INDENT_CODE_END     // <4 spaces
    LIST_ITEM           // - item or 1. item
    BLOCKQUOTE          // > quote

    // Content tokens
    TEXT                // Regular text
    BLANK_LINE          // Empty line (structural significance)

    // Meta tokens
    EOF                 // End of file
    ERROR               // Lexical error
)

type Token struct {
    Type    TokenType
    Literal string    // Raw text
    Line    int       // Line number (1-indexed)
    Column  int       // Column number (1-indexed)
}
```

**Lexer State Machine**:
```
State: Normal
  - Read '#' → count hashes → emit HEADING_N
  - Read '`' → check for triple → emit CODE_FENCE_START
  - Read '-' or digit → check list pattern → emit LIST_ITEM
  - Read '>' → emit BLOCKQUOTE
  - Read '\n\n' → emit BLANK_LINE
  - Default → accumulate TEXT

State: InCodeFence
  - Read '```' → emit CODE_FENCE_END → return to Normal
  - Everything else → accumulate TEXT
  - Ignore all markdown syntax inside code fence

State: InIndentCode
  - Read line with <4 spaces → emit INDENT_CODE_END → return to Normal
  - Read line with 4+ spaces → accumulate TEXT
```

**Implementation Pattern** (inspired by `text/template/parse`):
```go
type Lexer struct {
    input   string    // The markdown being scanned
    start   int       // Start position of current token
    pos     int       // Current position
    width   int       // Width of last rune read
    line    int       // Current line number
    col     int       // Current column number
    state   stateFn   // Current state function
    tokens  chan Token // Channel of scanned tokens
}

type stateFn func(*Lexer) stateFn

// Example state function
func lexNormal(l *Lexer) stateFn {
    for {
        switch r := l.next(); {
        case r == '#':
            return lexHeading
        case r == '`':
            return lexBacktick
        case r == eof:
            l.emit(EOF)
            return nil
        default:
            return lexText
        }
    }
}
```

#### Stage 2: Parser (AST Construction)

**Purpose**: Convert token stream into abstract syntax tree

**AST Node Types**:
```go
type NodeType int

const (
    NodeDocument      // Root node
    NodeHeading       // # Header
    NodeCodeBlock     // ``` code ```
    NodeList          // - item
    NodeListItem      // Individual item
    NodeParagraph     // Text block
    NodeText          // Raw text
)

type Node interface {
    Type() NodeType
    Position() Position
    Children() []Node
}

type Position struct {
    Line   int
    Column int
}

type Document struct {
    Nodes []Node
}

type Heading struct {
    Level    int       // 1-6
    Text     string    // Header text
    Content  []Node    // Content until next heading of same/higher level
    Pos      Position
}

type CodeBlock struct {
    Language string    // Language identifier (if any)
    Code     string    // Code content
    Pos      Position
}
```

**Parser Logic**:
```go
type Parser struct {
    lexer   *Lexer
    current Token
    peek    Token
}

func (p *Parser) Parse() (*Document, error) {
    doc := &Document{Nodes: make([]Node, 0)}

    for p.current.Type != EOF {
        node := p.parseNode()
        if node != nil {
            doc.Nodes = append(doc.Nodes, node)
        }
    }

    return doc, nil
}

func (p *Parser) parseNode() Node {
    switch p.current.Type {
    case HEADING_1, HEADING_2, HEADING_3, HEADING_4:
        return p.parseHeading()
    case CODE_FENCE_START:
        return p.parseCodeBlock()
    case LIST_ITEM:
        return p.parseList()
    case TEXT:
        return p.parseParagraph()
    default:
        return nil
    }
}
```

**Context-Aware Parsing**:
- Track current heading level to determine content boundaries
- Inside code blocks, preserve everything verbatim
- Handle nested lists and code blocks within lists

#### Stage 3: Extractor (Domain Structure Extraction)

**Purpose**: Walk AST to extract Spectr-specific structures

**Extractor API** (matches current parsers package):
```go
// Extract requirements from a spec file
func ExtractRequirements(doc *Document) []RequirementBlock {
    extractor := &RequirementExtractor{}
    return extractor.Extract(doc)
}

// Extract sections (## headers → content map)
func ExtractSections(doc *Document) map[string]string {
    extractor := &SectionExtractor{}
    return extractor.Extract(doc)
}

// Extract delta operations from change spec
func ExtractDeltas(doc *Document) *DeltaPlan {
    extractor := &DeltaExtractor{}
    return extractor.Extract(doc)
}
```

**Extraction Logic**:
```go
type RequirementExtractor struct {
    inRequirementsSection bool
}

func (e *RequirementExtractor) Extract(doc *Document) []RequirementBlock {
    requirements := make([]RequirementBlock, 0)

    for _, node := range doc.Nodes {
        if h := asHeading(node); h != nil {
            if h.Level == 2 && containsIgnoreCase(h.Text, "Requirements") {
                e.inRequirementsSection = true
                continue
            }

            if h.Level == 3 && e.inRequirementsSection {
                if strings.HasPrefix(h.Text, "Requirement:") {
                    req := e.extractRequirement(h)
                    requirements = append(requirements, req)
                }
            }
        }
    }

    return requirements
}

func (e *RequirementExtractor) extractRequirement(heading *Heading) RequirementBlock {
    name := strings.TrimPrefix(heading.Text, "Requirement:")
    name = strings.TrimSpace(name)

    // Collect scenarios from heading content
    scenarios := e.extractScenarios(heading.Content)

    // Reconstruct raw markdown
    raw := reconstructMarkdown(heading)

    return RequirementBlock{
        HeaderLine: heading.Text,
        Name:       name,
        Raw:        raw,
        Scenarios:  scenarios,
    }
}
```

### Error Handling

**Position-Aware Errors**:
```go
type ParseError struct {
    Message  string
    Position Position
    Context  string    // Surrounding text for context
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("%s at line %d, column %d\n%s",
        e.Message, e.Position.Line, e.Position.Column, e.Context)
}
```

**Error Examples**:
```
Unclosed code fence at line 42, column 1
   40: ### Requirement: Feature
   41: The system SHALL...
-> 42: ```go
   43: func Example() {
```

### Package Structure

```
internal/parsers/markdown/
├── token.go         # Token types and Token struct
├── lexer.go         # Lexer state machine
├── lexer_test.go    # Lexer unit tests
├── ast.go           # AST node types
├── parser.go        # Parser logic
├── parser_test.go   # Parser unit tests
├── extractor.go     # Spectr structure extraction
├── extractor_test.go # Extractor unit tests
└── doc.go           # Package documentation
```

**Public API** (in `internal/parsers/markdown`):
```go
// High-level API
func ParseMarkdown(content string) (*Document, error)
func ExtractRequirements(content string) ([]RequirementBlock, error)
func ExtractSections(content string) (map[string]string, error)
func ExtractDeltas(content string) (*DeltaPlan, error)

// Low-level API (for advanced use)
func NewLexer(input string) *Lexer
func NewParser(lexer *Lexer) *Parser
```

## Decisions

### Decision 1: Custom Implementation vs External Library

**Alternatives Considered**:

1. **goldmark** - Full markdown parser, AST-based
   - ✅ Pros: Battle-tested, compliant, extensible
   - ❌ Cons: Heavyweight (30+ packages), complex API, more than we need

2. **blackfriday** - Older markdown parser
   - ✅ Pros: Simpler than goldmark
   - ❌ Cons: Unmaintained, still more than needed

3. **text/template/parse** - Go's template lexer/parser
   - ✅ Pros: In stdlib, proven pattern
   - ❌ Cons: Template syntax, not markdown

4. **Custom implementation** (CHOSEN)
   - ✅ Pros: Tailored to Spectr needs, no dependencies, full control, learning opportunity
   - ❌ Cons: More initial work, need to handle edge cases

**Rationale**:
- Spectr needs a **subset** of markdown (headings, code blocks, lists, text)
- External libraries parse **everything** (emphasis, links, images, tables, etc.)
- Custom implementation: ~500 lines vs ~5000+ lines for goldmark
- Dependency avoidance aligns with project philosophy (zero external API dependencies)
- Educational value: team learns lexer/parser patterns

### Decision 2: Lexer State Machine vs Regex

**Choice**: State machine with state functions

**Rationale**:
- State functions naturally handle context (inside code block, inside list, etc.)
- Rob Pike's lexer pattern is proven in Go ecosystem (`text/template`, `go/scanner`)
- Easier to debug and test than complex regex
- Better error messages with position tracking
- Performance comparable to compiled regex

### Decision 3: AST Depth

**Choice**: Shallow AST focused on structure

**What we parse**:
- Headings (all levels)
- Code blocks (fenced and indented)
- Lists (unordered and ordered)
- Text paragraphs
- Blank lines (structural significance)

**What we ignore** (preserve as text):
- Emphasis (*italic*, **bold**)
- Links and images
- Inline code
- Tables
- HTML

**Rationale**:
- Requirements and scenarios contain inline markdown, but we don't need to parse it
- Simpler AST = simpler parser = fewer bugs
- Can expand later if needed

### Decision 4: Incremental Adoption

**Choice**: Implement new parser alongside old, migrate incrementally

**Migration Path**:
1. Implement `internal/parsers/markdown` with full test suite
2. Add integration tests comparing old vs new parser output
3. Update `internal/parsers` to use new parser internally
4. Run full test suite to ensure no regressions
5. Remove old regex-based code after validation

**Rationale**:
- De-risk the migration
- Easier to debug differences
- Can rollback if issues found

### Decision 5: Performance Targets

**Targets**:
- **Throughput**: Parse 100KB markdown in <10ms (current: ~8ms with regex)
- **Memory**: <1MB allocation for typical spec file
- **Regression**: <10% slower than current implementation
- **Scale**: Linear performance up to 1000 requirements per file

**Measurement**:
- Benchmark suite comparing old vs new parser
- Profile memory allocations
- Test with large generated spec files (100, 500, 1000 requirements)

## Alternatives Considered

### Alternative 1: Fix Regex Patterns

**What**: Add special cases to regex for code blocks

**Why Not**:
- Regex for "not inside code fence" is extremely complex
- Doesn't solve indented code blocks
- Doesn't solve nested structures
- Perpetuates maintenance burden

### Alternative 2: Pre-process to Remove Code Blocks

**What**: Strip code blocks before regex parsing

**Why Not**:
- Loses position information for errors
- Can't handle scenarios inside code blocks (valid use case)
- Fragile - what if code block markers are malformed?

### Alternative 3: Use Goldmark + Custom Renderer

**What**: Parse with goldmark, walk AST with custom renderer

**Why Not**:
- Still heavyweight dependency (30+ packages)
- Learning curve for goldmark's AST model
- Overkill for our needs (we don't need 90% of markdown features)

## Risks / Trade-offs

### Risk: Implementation Complexity

**Risk**: Custom parser might have subtle bugs

**Mitigation**:
- Comprehensive test suite with edge cases
- Fuzzing with generated markdown
- Compare output with current parser (should match for all existing specs)
- Gradual rollout - implement, test, migrate

### Risk: Performance Regression

**Risk**: New parser might be slower than regex

**Mitigation**:
- Benchmark early and often
- Profile hot paths
- Optimize lexer state machine if needed
- Set hard performance requirement: <10% regression

### Risk: Maintenance Burden

**Risk**: Custom code requires ongoing maintenance

**Counter-argument**:
- Current regex code also requires maintenance
- Lexer/parser is more maintainable than complex regex
- Clear architecture makes changes easier
- Comprehensive tests catch regressions

### Trade-off: Code Volume

**Trade-off**: More lines of code than current implementation

**Justification**:
- Better separation of concerns (lexer, parser, extractor)
- More test coverage
- Better documentation
- Improved maintainability despite higher line count
- Estimate: +400 lines (lexer: 200, parser: 150, extractor: 50)

### Trade-off: Initial Development Time

**Trade-off**: Takes longer than fixing regex

**Justification**:
- Pays off in reduced future maintenance
- Handles current and future edge cases
- Better foundation for potential future features (markdown linting, formatting)

## Implementation Plan

### Phase 1: Lexer Foundation (Estimated: 2-3 focused sessions)

1. Define token types (`token.go`)
2. Implement lexer state machine (`lexer.go`)
   - Normal state (headings, text)
   - Code fence state
   - List state
3. Write lexer tests (`lexer_test.go`)
   - Token stream correctness
   - Position tracking
   - Edge cases (nested code blocks, etc.)

### Phase 2: Parser & AST (Estimated: 2-3 focused sessions)

1. Define AST node types (`ast.go`)
2. Implement parser (`parser.go`)
   - Document parsing
   - Heading hierarchy
   - Code block handling
   - List parsing
3. Write parser tests (`parser_test.go`)
   - AST structure correctness
   - Round-trip testing (parse → reconstruct)

### Phase 3: Extractor (Estimated: 1-2 focused sessions)

1. Implement requirement extraction (`extractor.go`)
2. Implement section extraction
3. Implement delta extraction
4. Write extractor tests (`extractor_test.go`)
   - Match current parser output
   - Handle edge cases

### Phase 4: Integration (Estimated: 1-2 focused sessions)

1. Update `internal/parsers` to use new parser
2. Update `internal/validation/parser.go`
3. Run full test suite
4. Fix any regressions
5. Add integration tests

### Phase 5: Validation & Cleanup (Estimated: 1 focused session)

1. Performance benchmarks
2. Memory profiling
3. Remove old regex code
4. Update documentation
5. Final testing

**Total Estimate**: 7-11 focused coding sessions

## Testing Strategy

### Unit Tests

**Lexer Tests** (`lexer_test.go`):
```go
func TestLexer_Headings(t *testing.T) { ... }
func TestLexer_CodeFence(t *testing.T) { ... }
func TestLexer_CodeFenceWithHeaderInside(t *testing.T) { ... }
func TestLexer_Lists(t *testing.T) { ... }
func TestLexer_PositionTracking(t *testing.T) { ... }
```

**Parser Tests** (`parser_test.go`):
```go
func TestParser_Document(t *testing.T) { ... }
func TestParser_HeadingHierarchy(t *testing.T) { ... }
func TestParser_CodeBlockDoesNotInterfere(t *testing.T) { ... }
func TestParser_NestedLists(t *testing.T) { ... }
```

**Extractor Tests** (`extractor_test.go`):
```go
func TestExtractRequirements(t *testing.T) { ... }
func TestExtractRequirements_WithCodeBlocks(t *testing.T) { ... }
func TestExtractSections(t *testing.T) { ... }
func TestExtractDeltas(t *testing.T) { ... }
```

### Integration Tests

**Compatibility Tests**:
```go
// Ensure new parser matches old parser output for all existing specs
func TestBackwardCompatibility(t *testing.T) {
    specs := discoverAllSpecs()
    for _, spec := range specs {
        oldOutput := oldParser.Parse(spec)
        newOutput := newParser.Parse(spec)
        assert.Equal(t, oldOutput, newOutput)
    }
}
```

**Edge Case Tests**:
```go
func TestMarkdownEdgeCases(t *testing.T) {
    tests := []struct {
        name     string
        markdown string
        expected Result
    }{
        {
            name: "Code block with requirement header inside",
            markdown: `
## Requirements

### Requirement: Feature

\`\`\`markdown
### Requirement: Not a Real Requirement
This is inside a code block
\`\`\`

#### Scenario: Test
- **WHEN** ...
`,
            expected: Result{
                Requirements: 1,  // Only the real one
                Scenarios: 1,
            },
        },
        // More edge cases...
    }
}
```

### Benchmark Tests

```go
func BenchmarkLexer(b *testing.B) { ... }
func BenchmarkParser(b *testing.B) { ... }
func BenchmarkOldParser(b *testing.B) { ... }
func BenchmarkNewParser(b *testing.B) { ... }

// Comparison
func BenchmarkComparison(b *testing.B) {
    b.Run("OldParser", func(b *testing.B) { ... })
    b.Run("NewParser", func(b *testing.B) { ... })
}
```

## Migration Path

### Step 1: No User Impact

- Implementation happens in `internal/parsers/markdown`
- Existing code unchanged
- All tests pass

### Step 2: Internal Migration

- Update `internal/parsers` to call new parser
- Keep old code as fallback
- Feature flag to enable/disable new parser
- Run both parsers, compare results

### Step 3: Validation

- Run on all specs in repository
- Ensure identical output (or document improvements)
- Performance benchmarks confirm <10% regression
- Memory profiling shows acceptable allocation

### Step 4: Rollout

- Enable new parser by default
- Monitor for issues
- Keep old code for one release cycle
- Remove old code after validation period

### Step 5: Leverage

- Use improved error messages
- Add new features (markdown linting, formatting)
- Expand to handle more markdown features if needed

## Success Criteria

1. **Correctness**: All existing specs parse identically (or with documented improvements)
2. **Edge cases**: Code blocks with markdown syntax parse correctly
3. **Performance**: <10% regression on benchmark suite
4. **Tests**: >90% code coverage on new parser code
5. **API compatibility**: No changes to public parsers package API
6. **Error messages**: Position-aware errors with context
7. **Documentation**: Comprehensive godoc and package documentation

## Open Questions

1. **Fuzzing**: Should we add fuzzing tests to catch edge cases?
   - **Answer**: Yes, use Go's built-in fuzzing for lexer/parser

2. **Incremental parsing**: Do we need to support parsing partial documents?
   - **Answer**: No, files are small enough to parse completely

3. **Streaming**: Should lexer support streaming input?
   - **Answer**: No, load entire file for simplicity

4. **AST mutation**: Should AST be mutable for future editing features?
   - **Answer**: No, read-only for now, can add later if needed

5. **Markdown dialect**: Which markdown flavor should we target?
   - **Answer**: CommonMark subset - only features Spectr uses
