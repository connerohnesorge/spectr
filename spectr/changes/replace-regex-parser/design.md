# Design: Lexer/Parser Architecture

## Two-Layer Architecture

This design establishes a clear separation between generic markdown parsing and Spectr-specific extraction logic through a two-layer architecture:

### Layer 1: Generic Markdown Parser (internal/mdparser)
A reusable markdown parser with **zero knowledge of Spectr semantics**. This layer provides:
- Lexical analysis (tokenization) of markdown syntax
- Syntax tree construction (AST) representing markdown structure
- No understanding of "Requirements", "Scenarios", or "Delta sections"

**Rationale:** Separating generic markdown parsing from domain logic provides:
- **Reusability**: The markdown parser can be used for any markdown processing task
- **Testability**: Generic parsing can be tested independently of Spectr semantics
- **Maintainability**: Changes to Spectr requirements don't affect the core parser
- **Separation of Concerns**: Syntax vs. semantics are clearly delineated

### Layer 2: Spectr Extractors (internal/parsers)
Domain-specific logic that **interprets markdown structure as Spectr elements**. This layer:
- Consumes the generic AST from Layer 1
- Applies Spectr business rules and validation
- Extracts Requirements, Scenarios, Delta sections
- Validates Spectr-specific conventions

**API Boundary:** Layer 2 depends on Layer 1, but Layer 1 has no knowledge of Layer 2. This ensures the markdown parser remains generic and reusable.

## Architecture

### Layer 1: Generic Markdown Parser (internal/mdparser)

The generic markdown parser converts raw markdown text into a structured Abstract Syntax Tree (AST) through two phases: lexing and parsing.

#### Lexer

The lexer performs tokenization, scanning input text and emitting a stream of `Token`s. The implementation follows the Go compiler's scanner pattern (`cmd/compile/internal/syntax/scanner.go`).

**Key Components:**
- `TokenType`: Enum for markdown token types (e.g., `TokenHeader`, `TokenText`, `TokenCodeBlock`, `TokenListItem`, `TokenEOF`)
  - **IMPORTANT**: Tokens are markdown-agnostic, not Spectr-specific
  - No `TokenRequirementHeader` or `TokenScenario` - these are Layer 2 concerns
- `Token`: Struct containing `Type`, `Value` (string), and `Position` (line/col for error reporting)
- `Lexer`: Struct holding the input, current position, and token emission mechanism
- `StateFn`: Function type `func(*Lexer) StateFn` (adapting Rob Pike's state machine approach)

**States:**
- `lexText`: Default state, consumes plain text
- `lexHeader`: Handles markdown headers (`#`, `##`, etc.)
- `lexCodeBlock`: Handles code fences (``````` ) and content within them
- `lexList`: Handles list items (`-`, `*`, numbered lists)

**Note:** The lexer produces markdown tokens without semantic interpretation. A `###` header is just a `TokenHeader` with level 3, not a "Requirement header".

#### Parser

The parser consumes the token stream and builds an Abstract Syntax Tree (AST), following the recursive descent pattern used in the Go compiler (`cmd/compile/internal/syntax/parser.go`).

**AST Nodes:**
- `Node`: Interface for all AST nodes
- `Document`: Root node containing all top-level elements
- `Header`: Represents a markdown header (level + text content)
- `Paragraph`: Represents a block of text
- `CodeBlock`: Represents a fenced code block (language + content)
- `List`: Represents a list (ordered/unordered + items)

**Emphasis:** These AST nodes represent **generic markdown constructs**, not Spectr semantics. A `Header` node doesn't know if it's a "Requirement" or just a section title. The AST is a faithful representation of the markdown document structure.

The parser uses recursive descent with token-based dispatching, ensuring robust error handling and structure recovery.

#### API

The public interface for the generic markdown parser:

```go
// Lexer API
func NewLexer(input string) *Lexer
func (l *Lexer) NextToken() Token

// Parser API
func Parse(input string) (*Document, error)

// AST traversal helpers
func (d *Document) Walk(visitor NodeVisitor)
func (h *Header) Level() int
func (h *Header) Text() string
```

**Contract:** The parser guarantees a well-formed AST representing valid markdown structure. It does NOT validate Spectr conventions.

### Layer 2: Spectr Extractors (internal/parsers)

Spectr extractors consume the generic markdown AST and apply domain-specific business rules to identify and validate Spectr elements.

**Core Responsibilities:**

1. **Requirement Extraction**
   - Find `Header` nodes matching pattern `### Requirement: [name]`
   - Validate that Requirements appear at the correct hierarchy level
   - Extract associated content (paragraphs, code blocks) as requirement body

2. **Scenario Extraction**
   - Find `Header` nodes matching pattern `#### Scenario: [name]`
   - Validate that Scenarios are children of Requirements
   - Extract scenario steps (list items with WHEN/THEN/GIVEN patterns)

3. **Delta Section Extraction**
   - Find `Header` nodes for Delta sections (`## ADDED Requirements`, etc.)
   - Validate delta structure and operation types
   - Extract requirements within each delta section

**Business Rules Applied:**
- Requirement headers must follow `### Requirement: [name]` pattern (case-sensitive)
- Scenarios must be `####` level headers (one level below Requirements)
- Scenarios must have at least one step
- Delta sections must contain valid operation types (ADDED, MODIFIED, REMOVED, RENAMED)
- MODIFIED and RENAMED requirements must reference existing requirements

**Example Logic:**
```go
func ExtractRequirements(doc *mdparser.Document) ([]RequirementBlock, error) {
    var requirements []RequirementBlock

    // Walk AST looking for Level 3 headers
    for _, node := range doc.Children() {
        header, ok := node.(*mdparser.Header)
        if !ok || header.Level() != 3 {
            continue
        }

        // Apply Spectr business rule: must start with "Requirement: "
        if !strings.HasPrefix(header.Text(), "Requirement: ") {
            continue
        }

        // Extract name and associated content
        name := strings.TrimPrefix(header.Text(), "Requirement: ")
        scenarios := extractScenariosUnderHeader(header)

        if len(scenarios) == 0 {
            return nil, fmt.Errorf("requirement %q has no scenarios", name)
        }

        requirements = append(requirements, RequirementBlock{
            Name: name,
            Scenarios: scenarios,
        })
    }

    return requirements, nil
}
```

#### API

The public interface for Spectr extractors:

```go
// Requirement extraction
func ExtractRequirements(doc *mdparser.Document) ([]RequirementBlock, error)
func ExtractScenarios(req *mdparser.Header) ([]Scenario, error)

// Delta extraction
func ExtractDeltaSpec(doc *mdparser.Document) (*DeltaPlan, error)
func ExtractDeltaSection(doc *mdparser.Document, operation string) ([]RequirementBlock, error)

// Validation helpers
func ValidateRequirementStructure(req RequirementBlock) error
func ValidateScenarioFormat(scenario Scenario) error
```

**Contract:** Extractors return domain objects (`RequirementBlock`, `Scenario`, `DeltaPlan`) or errors if Spectr business rules are violated.

## Code Sharing Architecture

The two-layer architecture eliminates parsing duplication across packages while maintaining clear separation of concerns.

### Current Duplication Problem

Currently, there are three separate parsing implementations:
1. **internal/parsers/**: Primary parsing for requirements and delta specs
2. **internal/validation/parser.go**: Duplicate implementations of ExtractSections, ExtractRequirements, ExtractScenarios
3. **internal/archive/spec_merger.go**: Additional regex-based parsing for spec merging

This duplication means:
- Bug fixes must be applied in multiple places
- Inconsistent behavior between parsing contexts
- Same brittleness issues affect all packages

### Unified Architecture

The new design establishes clear layers:

```
internal/mdparser/         ← NEW: Generic markdown parser (Layer 1)
├── lexer.go              → Tokenization (markdown-agnostic)
├── parser.go             → AST construction (markdown-agnostic)
└── ast.go                → AST node definitions

internal/parsers/          ← UPDATED: Spectr extractors (Layer 2)
├── extractor.go          → NEW: Core Spectr extraction logic
├── requirement_parser.go → Uses mdparser + extractor
└── delta_parser.go       → Uses mdparser + extractor

internal/validation/
├── parser.go             → Uses mdparser + extractor
└── change_rules.go       → Uses mdparser + extractor for RENAMED parsing

internal/archive/
└── spec_merger.go        → Uses mdparser + extractor for merging
```

### Integration Strategy

All call sites will use the same two-layer approach:

```go
// Step 1: Parse markdown into generic AST (Layer 1)
doc, err := mdparser.Parse(content)
if err != nil {
    return fmt.Errorf("markdown parsing failed: %w", err)
}

// Step 2: Extract Spectr elements with business rules (Layer 2)
requirements, err := parsers.ExtractRequirements(doc)
if err != nil {
    return fmt.Errorf("requirement extraction failed: %w", err)
}
```

### Benefits of Shared Architecture

- **Single Source of Truth**: One parser implementation for all contexts
- **Consistent Behavior**: Same parsing logic in validation, archiving, and extraction
- **Easier Maintenance**: Bug fixes and enhancements in one place
- **Better Testing**: Comprehensive tests on shared parser benefit all consumers
- **Reduced Code**: Eliminates ~300+ lines of duplicate parsing logic
- **Clear Separation**: Generic markdown parsing vs. Spectr semantics are distinct concerns

## Performance Considerations

### Expected Performance Profile

- **Lexer/Parser Performance**: May be 2-3x slower than regex for simple files
  - Regex: Single-pass pattern matching, minimal allocation
  - Lexer/Parser: Token stream + AST construction, more allocations

- **Trade-off Justification**: Correctness is the priority
  - The new architecture handles edge cases that regex cannot (e.g., code blocks containing requirement-like text)
  - Spectr processes small files (typically <1000 lines) where absolute performance is less critical
  - Robustness and maintainability outweigh raw speed

### Optimization Strategy

If benchmarking reveals unacceptable performance:

1. **Token Pooling**: Reuse token objects to reduce allocations
   ```go
   var tokenPool = sync.Pool{
       New: func() interface{} { return &Token{} },
   }
   ```

2. **Lazy AST Construction**: Build AST nodes on-demand during traversal
   - Don't materialize entire tree if only extracting specific sections
   - Stream-based processing for large files

3. **Caching**: Cache parsed ASTs for repeated operations
   - Archive workflow parses same spec files multiple times
   - Validation may parse files multiple times during interactive mode

4. **Profiling-Driven**: Use `pprof` to identify actual bottlenecks before optimizing

### Benchmark Suite

We will create benchmarks comparing:
- Regex-based parsing (current implementation)
- Lexer/parser implementation (new implementation)
- Various file sizes (100 lines, 1000 lines, 5000 lines)
- Edge cases (deep nesting, large code blocks)

**Decision Gate**: Only proceed with full migration if benchmarks show acceptable performance for typical Spectr use cases.

## Trade-offs

### Pros
- **Robustness**: Explicitly handles state, preventing misinterpretation of syntax inside code blocks
- **Precision**: Tracks line and column numbers for every token, enabling high-quality error messages
- **Separation of Concerns**: Lexing (syntax) is separated from Parsing (structure) and Extraction (semantics)
- **Reusability**: Generic markdown parser can be used for future features or separate tools
- **Maintainability**: Clear layering makes it obvious where to add new functionality

### Cons
- **Complexity**: More code than the current regex approach (~500 lines vs ~200 lines)
- **Performance**: Likely slower than raw regex for simple cases (mitigated by optimization if needed)
- **Migration Risk**: Requires updating multiple packages and validating behavior consistency

## Migration Strategy

The migration follows a phased approach with clear validation gates:

### Phase 1: Implement Core Parser (Tasks 1.1-1.7)
- Implement `internal/mdparser` package (lexer, parser, AST)
- Comprehensive unit tests for all markdown constructs
- Test edge cases: code blocks, nested lists, mixed content
- **Gate**: All markdown parsing tests pass

### Phase 2: Benchmark and Validate Performance (NEW)
- Create benchmark suite comparing regex vs. lexer/parser
- Test with real Spectr spec files from repository
- Profile allocation and CPU usage
- **Decision Gate**: If performance is unacceptable (>5x slower), implement optimizations
- **Exit Criteria**: Performance acceptable for typical use cases OR optimizations applied

### Phase 3: Implement Spectr Extractors (Tasks 2.1-2.3)
- Implement `internal/parsers/extractor.go` using mdparser
- Update `requirement_parser.go` and `delta_parser.go` to use extractor
- Verify existing `internal/parsers` tests pass
- **Gate**: All existing parser tests pass with new implementation

### Phase 4: Migrate Call Sites (Tasks 3.1-4.3)
- Replace `internal/validation/parser.go` implementations
- Update `internal/validation/change_rules.go` RENAMED parsing
- Replace `internal/archive/spec_merger.go` parsing logic
- **Gate**: All existing tests pass for validation and archive packages

### Phase 5: Comprehensive Validation (Tasks 5.1-5.4)
- Integration testing with full workflows (validate, archive, show)
- Edge case validation (markdown in code blocks, frontmatter, special characters)
- Regression testing against production spec files
- **Gate**: All integration tests pass, no behavior regressions

### Phase 6: Cleanup and Remove Old Code
- Remove old regex-based parsing functions
- Update documentation
- Archive this change proposal
- **Gate**: No references to old parsing code remain

**Rollback Plan**: If critical issues are discovered post-migration, the old regex code can be restored from git history. All phases maintain backward compatibility until Phase 6.
