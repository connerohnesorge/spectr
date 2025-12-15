## Context
Spectr currently uses a consolidated `internal/regex/` package with pre-compiled regex patterns for parsing markdown spec files. While this works, regex-based parsing has limitations:
- Error messages lack precise location (line + column)
- Patterns are fragile when edge cases arise
- Adding new markdown features requires new regex patterns
- Testing regex patterns is cumbersome

The user requested eliminating all regex from the codebase and implementing a proper handbuilt parser with line-based error reporting.

**Stakeholders**: Spectr users, AI agents parsing spec files, CI/CD pipelines validating specs.

## User Decisions (from requirements gathering)

The following decisions were made based on user input:

### Core Architecture Decisions
| Question | User's Choice | Implication |
|----------|---------------|-------------|
| Markdown scope | **Full CommonMark subset** | Support headers, lists, code blocks, emphasis, links - not just minimal Spectr patterns |
| Error handling | **Strict with errors** | Return line-based errors for malformed input; collect all errors rather than stopping at first |
| Parser architecture | **Token-based lexer/parser** | Separate tokenization pass then parse tokens into AST; more complex but more flexible |
| API approach | **Clean slate API** | New `internal/markdown/` package with improved API; requires updating all call sites |

### Token and Lexer Decisions
| Question | User's Choice | Implication |
|----------|---------------|-------------|
| Token granularity | **Fine-grained** | Each delimiter is separate token (e.g., `*`, `*`, text, `*`, `*` for bold). Maximum flexibility for error recovery. |
| Lexer error handling | **Error tokens** | Emit `TokenError` with message, continue lexing. Allows collecting all errors in one pass. |
| Lexer API | **Internal only** | Lexer is implementation detail. Only expose parser API. Simpler public interface. |

### AST Decisions
| Question | User's Choice | Implication |
|----------|---------------|-------------|
| AST mutability | **Immutable with builders** | Nodes are read-only after creation. Use builder/transform functions to create modified copies. Thread-safe, predictable. |
| Position tracking | **Byte offset only** | Track only byte offsets. Line and column calculated on-demand from offset and source. Compact storage. |
| Inline content model | **Flat children array** | Paragraph has Children: [TextNode, StrongNode, TextNode]. Simple, matches CommonMark AST. |
| AST traversal | **Visitor pattern only** | Classic visitor pattern with Accept/Visit methods. Good for operations that vary by node type. |
| Node identity | **Content hash** | Hash of node content determines identity. Nodes with same content share identity. Good for caching. |
| Source preservation | **Store original text** | Each node stores its original source substring as byte slice view. Enables exact round-trip. |
| String representation | **Byte slice views** | `[]byte` slices into original source. Zero-copy but requires source lifetime management. |
| Node struct style | **Typed node structs** | Separate NodeSection, NodeRequirement, etc. implementing Node interface. Type-safe via assertions. |
| Type-specific fields | **Getter methods** | Private fields with Level(), URL() getters. Enforces immutability with verbose but safe API. |
| Parent pointers | **No parent pointers** | Nodes only reference children. Parent passed via visitor context. Simpler, truly immutable. |

### Parser Decisions
| Question | User's Choice | Implication |
|----------|---------------|-------------|
| Incremental parsing | **Tree-sitter style** | Full incremental with tree diffing. Maximum performance but significant complexity. |
| Edit granularity | **Full diff-based** | Accept old and new text, compute diff internally. Most flexible but slower for known edits. |
| Parser API | **Stateless Parse func** | No Parser struct, just Parse(source) function. Each call independent. Most concurrent. |
| Emphasis handling | **CommonMark strict** | Follow CommonMark spec exactly for delimiter matching. Handles edge cases like `*a _b* c_` correctly. |

### Additional Feature Decisions
| Question | User's Choice | Implication |
|----------|---------------|-------------|
| AST printer | **Normalized printer** | Regenerate markdown with consistent, minimal formatting. Useful for auto-formatting. |
| Position queries | **Interval tree index** | Build interval tree for O(log n) position queries. Extra memory but fast for repeated queries. |
| Index building | **Lazy on first query** | Build interval tree on first PositionQuery call. No overhead if never queried. |
| Transform API | **Visitor-based transforms** | TransformVisitor returns replacement nodes. Composable, functional style. |
| Transform signals | **Return (node, action)** | Return node and action enum (Keep/Replace/Delete). Explicit intent, clear semantics. |
| Query API | **Predicate-based Find** | Find(root, func(Node) bool) returns matching nodes. Flexible, Go-idiomatic. |
| Object pooling | **Full pooling** | Pool both tokens and nodes. Maximum performance but complex lifetime management. |
| Error location | **Offset only** | Error stores byte offset. Caller uses LineIndex to convert. Minimal data. |
| Streaming | **No streaming** | Always build full AST. Spectr files are small, streaming unnecessary complexity. |

## Goals / Non-Goals

### Goals
- Replace all regex-based markdown parsing with a token-based lexer/parser
- Provide line and column numbers in all parse errors
- Support CommonMark subset: headers, lists, code blocks, emphasis, links
- Maintain or improve parse performance
- Retain existing functionality (validation, merging, requirement extraction)
- Clean, well-tested API in `internal/markdown/` package

### Non-Goals
- Full CommonMark compliance (we focus on what Spectr needs)
- Table support
- HTML passthrough
- Support for Setext headers (only ATX style)
- GFM extensions beyond task checkboxes and strikethrough

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                     internal/markdown/                          │
├─────────────────────────────────────────────────────────────────┤
│  ┌──────────┐    ┌──────────┐    ┌──────────┐    ┌──────────┐  │
│  │  lexer   │ -> │  tokens  │ -> │  parser  │ -> │   AST    │  │
│  └──────────┘    └──────────┘    └──────────┘    └──────────┘  │
│       │              │               │               │          │
│       v              v               v               v          │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │                    errors.go                              │  │
│  │  ParseError{Line, Column, Message, Snippet}               │  │
│  └──────────────────────────────────────────────────────────┘  │
├─────────────────────────────────────────────────────────────────┤
│                     Public API (api.go)                         │
│  - ParseSpec(content) -> *Spec, []ParseError                   │
│  - ParseDelta(content) -> *Delta, []ParseError                 │
│  - ExtractSections(content) -> map[string]Section              │
│  - ExtractRequirements(content) -> []Requirement               │
│  - FindSection(content, name) -> *Section, bool                │
└─────────────────────────────────────────────────────────────────┘
```

## Decisions

### Decision 1: Token-Based Lexer/Parser Architecture
**What**: Separate lexer (tokenizer) and parser passes.
**Why**:
- Clear separation of concerns
- Tokens are testable independently
- Parser logic is cleaner without character-level handling
- Easy to add new token types

**Alternatives considered**:
- Line-based state machine: Simpler but harder to extend for inline formatting
- Recursive descent without lexer: Mixes tokenization with parsing

### Decision 2: Fine-Grained Token Types
Token types are fine-grained with each delimiter as a separate token. See `specs/tokens/spec.md` for complete specification.

```go
type TokenType uint8

const (
    // Structural
    TokenEOF TokenType = iota
    TokenNewline
    TokenWhitespace
    TokenText
    TokenError  // Invalid input with error message

    // Punctuation delimiters (each character separate)
    TokenHash        // #
    TokenAsterisk    // *
    TokenUnderscore  // _
    TokenTilde       // ~
    TokenBacktick    // `
    TokenDash        // -
    TokenPlus        // +
    TokenDot         // .
    TokenColon       // :
    TokenPipe        // |

    // Brackets
    TokenBracketOpen   // [
    TokenBracketClose  // ]
    TokenParenOpen     // (
    TokenParenClose    // )
    TokenGreaterThan   // >

    // Special
    TokenNumber  // Digit sequence for ordered lists
    TokenX       // x or X in checkboxes
)

// Token with byte offset positions and zero-copy source view
type Token struct {
    Type    TokenType
    Start   int    // Byte offset from source start
    End     int    // Byte offset past last byte (exclusive)
    Source  []byte // Slice view into original source (zero-copy)
    Message string // Error message (only for TokenError)
}
```

### Decision 3: Immutable AST Node Types
Nodes are immutable with content hashing for identity. See `specs/ast/spec.md` for complete specification.

```go
type NodeType uint8

const (
    NodeDocument NodeType = iota
    NodeSection           // H2 section with Level
    NodeRequirement       // ### Requirement: with Name
    NodeScenario          // #### Scenario: with Name
    NodeParagraph
    NodeList
    NodeListItem
    NodeCodeBlock         // With Language, Content
    NodeBlockquote
    NodeCode              // Inline code
    NodeEmphasis          // Italic
    NodeStrong            // Bold
    NodeStrikethrough
    NodeLink              // With URL, Title
    NodeLinkDef           // [ref]: url definition
    NodeWikilink          // [[target|display#anchor]]
    NodeText
)

// Immutable node with content hash and source preservation
type Node struct {
    Type     NodeType
    Hash     uint64    // Content hash for identity/caching
    Start    int       // Byte offset from source start
    End      int       // Byte offset past last byte
    Source   []byte    // Original source slice (zero-copy)
    Children []*Node   // Immutable child slice

    // Type-specific fields accessed via methods:
    // - Level() int           (Section)
    // - Title() []byte        (Section)
    // - Name() string         (Requirement, Scenario)
    // - Language() []byte     (CodeBlock)
    // - URL() []byte          (Link)
    // - Target() []byte       (Wikilink)
    // - Display() []byte      (Wikilink)
    // - Anchor() []byte       (Wikilink)
}

// Line and column calculated on demand
func (n *Node) Position(idx *LineIndex) Position {
    return idx.PositionAt(n.Start)
}
```

### Decision 4: Error Structure
```go
type ParseError struct {
    Offset   int           // Byte offset where error occurred
    Message  string        // Human-readable error description
    Expected []TokenType   // What tokens would have been valid
}

// Line/column calculated on demand from offset
func (e ParseError) Position(idx *LineIndex) Position {
    return idx.PositionAt(e.Offset)
}

func (e ParseError) Error() string {
    return fmt.Sprintf("offset %d: %s", e.Offset, e.Message)
}
```

### Decision 5: Collected Errors with Recovery
**What**: Parser collects all errors and continues parsing via error recovery.
**Why**: Better UX - users see all problems at once, not one at a time.

```go
type ParseResult struct {
    Root   *Node         // Root document node (may be partial on errors)
    Errors []ParseError  // All errors encountered
    State  *ParseState   // State for incremental reparsing
}

func Parse(source []byte) ParseResult {
    // Collects errors, continues parsing where possible
    // On error, skips to next synchronization point
}
```

### Decision 6: Tree-Sitter Style Incremental Parsing
**What**: Full incremental parsing with diff-based edit detection and subtree reuse.
**Why**: Maximum performance for editor integrations and repeated validation.

See `specs/parser/spec.md` for complete specification.

```go
// Incremental parse with automatic diff computation
func ParseIncremental(oldTree *Node, oldSource, newSource []byte) ParseResult {
    // 1. Compute diff between oldSource and newSource
    // 2. Identify affected regions
    // 3. Reparse only changed regions
    // 4. Reuse unchanged subtrees (matched by content hash)
    // 5. Adjust offsets for nodes after edit point
}

// Edit region from diff
type EditRegion struct {
    StartOffset    int  // Where edit begins
    OldEndOffset   int  // Where old content ended
    NewEndOffset   int  // Where new content ends
}
```

### Decision 7: Visitor Pattern for AST Traversal
**What**: Classic visitor pattern with double dispatch.
**Why**: Type-safe, extensible for different operations.

See `specs/visitor/spec.md` for complete specification.

```go
type Visitor interface {
    VisitDocument(*NodeDocument) error
    VisitSection(*NodeSection) error
    VisitRequirement(*NodeRequirement) error
    VisitScenario(*NodeScenario) error
    VisitParagraph(*NodeParagraph) error
    // ... methods for all node types
}

// Walk traverses the AST calling visitor methods
func Walk(node *Node, v Visitor) error

// Sentinel to skip children without stopping
var SkipChildren = errors.New("skip children")
```

### Decision 8: Package Structure
```
internal/markdown/
├── doc.go              # Package documentation
├── token.go            # Token types and Token struct
├── lexer.go            # Lexer implementation (internal)
├── lexer_test.go
├── node.go             # Node interface and typed node structs
├── node_types.go       # NodeSection, NodeRequirement, etc. implementations
├── node_test.go
├── parser.go           # Stateless Parse function
├── parser_test.go
├── incremental.go      # ParseIncremental function
├── incremental_test.go
├── errors.go           # ParseError type (offset only)
├── lineindex.go        # LineIndex for line/col calculation
├── lineindex_test.go
├── visitor.go          # Visitor interface and Walk function
├── visitor_test.go
├── transform.go        # TransformVisitor and Transform function
├── transform_test.go
├── query.go            # Find, FindFirst, predicate combinators
├── query_test.go
├── index.go            # PositionIndex with interval tree
├── index_test.go
├── printer.go          # Normalized markdown printer
├── printer_test.go
├── pool.go             # Token and node object pools
├── pool_test.go
├── api.go              # Public API functions
├── api_test.go
├── spec.go             # Spec-specific parsing (requirements, scenarios)
├── spec_test.go
├── delta.go            # Delta-specific parsing (ADDED, MODIFIED, etc.)
├── delta_test.go
├── wikilink.go         # Wikilink resolution
├── wikilink_test.go
├── compat.go           # Compatibility helpers for migration
└── compat_test.go
```

### Decision 9: Migration Strategy
1. Create `internal/markdown/` package alongside `internal/regex/`
2. Implement lexer with comprehensive tests
3. Implement parser with comprehensive tests
4. Create compatibility layer that matches existing function signatures
5. Update call sites one package at a time
6. Delete `internal/regex/` once all call sites migrated
7. Update validation spec to remove regex-related requirements

## Lexer State Machine

The lexer operates primarily on a line-by-line basis for block elements, with inline tokenization for emphasis, code, and links.

### Block-Level State Machine
```
┌─────────┐     '#'      ┌──────────┐
│  START  │ ──────────── │  HEADER  │ ─── count '#' chars ─── emit TokenH1..H6
└─────────┘              └──────────┘
     │
     │ '-', '*', '+'     ┌──────────┐
     ├─────────────────- │   LIST   │ ─── check for checkbox [ ]/[x]
     │                   └──────────┘
     │
     │ '`' (triple)      ┌───────────┐
     ├─────────────────- │ CODE_FENCE│ ─── capture until matching fence
     │                   └───────────┘
     │
     │ '>'               ┌────────────┐
     ├─────────────────- │ BLOCKQUOTE │
     │                   └────────────┘
     │
     │ digit + '.'       ┌─────────────┐
     ├─────────────────- │ ORDERED_LIST│
     │                   └─────────────┘
     │
     │ other             ┌──────────┐
     └─────────────────- │ PARAGRAPH│ ─── inline tokenization
                         └──────────┘
```

### Inline Tokenization
For text within paragraphs, list items, and headers:
```
┌──────────┐    '**' or '__'    ┌────────┐
│  INLINE  │ ────────────────── │  BOLD  │
└──────────┘                    └────────┘
     │
     │ '*' or '_'        ┌──────────┐
     ├─────────────────- │  ITALIC  │
     │                   └──────────┘
     │
     │ '~~'              ┌───────────────┐
     ├─────────────────- │ STRIKETHROUGH │
     │                   └───────────────┘
     │
     │ '`'               ┌─────────────┐
     ├─────────────────- │ INLINE_CODE │
     │                   └─────────────┘
     │
     │ '['               ┌──────────┐
     └─────────────────- │   LINK   │ ─┬─ followed by ](url) → inline link
                         └──────────┘  └─ followed by ][ref] → reference link
```

### Link Definition Handling
Reference-style links require a two-pass approach:
1. **First pass**: Collect all link definitions `[ref]: url "optional title"` at block level
2. **Second pass**: Resolve `[text][ref]` references using collected definitions

Link definitions:
- Must appear at start of line (no leading whitespace except for continuation)
- Format: `[label]: destination "optional title"`
- Labels are case-insensitive for matching
- Can span multiple lines with indentation
- Are removed from final output (not rendered)

### Wikilink Handling (Spectr Extension)
Wikilinks provide a convenient way to link between specs and changes:

**Syntax**:
- `[[spec-name]]` - link to spec, display text is spec name
- `[[spec-name|Display Text]]` - link to spec with custom display text
- `[[changes/my-change]]` - link to a change
- `[[validation#Requirement: Spec File Validation]]` - link to specific requirement

**Resolution**:
- Parser produces `NodeWikilink` with `Target` and optional `Display` fields
- Resolution is deferred to rendering/validation layer (not parser's job)
- Resolution rules:
  1. Check `spectr/specs/{target}/spec.md`
  2. Check `spectr/changes/{target}/proposal.md`
  3. If contains `#`, treat as anchor within the target
- Unresolved wikilinks should be flagged by validation

**Token structure**:
```go
type WikilinkToken struct {
    Target  string  // e.g., "validation" or "changes/my-change"
    Display string  // optional display text after |
    Anchor  string  // optional anchor after # (e.g., "Requirement: Name")
    Line    int
    Column  int
}
```

### Emphasis Disambiguation Rules
Following CommonMark spec section 6.2-6.4:
1. **Delimiter run**: sequence of `*` or `_` not preceded/followed by same char
2. **Left-flanking**: not followed by whitespace, and either not followed by punctuation OR preceded by whitespace/punctuation
3. **Right-flanking**: not preceded by whitespace, and either not preceded by punctuation OR followed by whitespace/punctuation
4. **`*` can open/close emphasis** if left/right-flanking respectively
5. **`_` intraword**: `_` cannot open/close if surrounded by alphanumerics (e.g., `foo_bar_baz` is not emphasis)

Implementation will track delimiter stack for proper matching.

## Spectr-Specific Parsing

The parser will recognize Spectr-specific patterns:

### Requirement Headers
```markdown
### Requirement: Name Here
```
Parsed as `NodeRequirement` with `Content = "Name Here"`

### Scenario Headers
```markdown
#### Scenario: Description Here
```
Parsed as `NodeScenario` with `Content = "Description Here"`

### WHEN/THEN/AND Bullets
```markdown
- **WHEN** condition
- **THEN** result
- **AND** additional
```
Parsed as `NodeTaskItem` with special `Keyword` field.

### Delta Sections
```markdown
## ADDED Requirements
## MODIFIED Requirements
## REMOVED Requirements
## RENAMED Requirements
```
Recognized as `NodeSection` with `DeltaType` field.

## Non-Markdown Regex (Retained)

The following regex patterns in `internal/git/platform.go` are NOT markdown-related and will be retained or converted separately:
- SSH URL pattern: `^(?:[\w-]+@)?([^:]+):(.+)$`
- HTTPS URL pattern: `^(?:https?|ssh|git)://(?:[\w-]+@)?([^/]+)/(.+)$`

These could optionally be converted to string parsing functions for consistency, but are out of scope for the markdown parser change.

## Risks / Trade-offs

### Risk: Performance Regression
**Mitigation**: Benchmark lexer and parser against regex implementation. Token-based parsing is typically faster than regex for structured content.

### Risk: Edge Cases in Lexer
**Mitigation**: Comprehensive test suite with edge cases from CommonMark spec test suite. Fuzzing tests.

### Risk: Breaking Existing Functionality
**Mitigation**:
- Maintain parallel implementations during migration
- Extensive integration tests comparing old vs new output
- Gradual migration with feature flags if needed

### Trade-off: Increased Code Size
Handbuilt parser is more code than regex patterns, but:
- Code is more readable and maintainable
- Tests are clearer
- Errors are better

## Migration Plan

### Phase 1: Core Package (Est. complexity: Medium)
1. Create `internal/markdown/token.go` with token types
2. Create `internal/markdown/lexer.go` with lexer
3. Create `internal/markdown/node.go` with AST types
4. Create `internal/markdown/parser.go` with parser
5. Comprehensive unit tests for each component

### Phase 2: API and Spec Parsing (Est. complexity: Medium)
1. Create `internal/markdown/api.go` with public API
2. Create `internal/markdown/spec.go` for Spectr-specific parsing
3. Create `internal/markdown/delta.go` for delta parsing
4. Integration tests against real spec files

### Phase 3: Migration (Est. complexity: Low-Medium per file)
1. Update `internal/parsers/requirement_parser.go`
2. Update `internal/parsers/delta_parser.go`
3. Update `internal/validation/parser.go`
4. Update `internal/archive/spec_merger.go`
5. Update `cmd/accept.go`
6. Update `internal/validation/change_rules.go`

### Phase 4: Cleanup
1. Delete `internal/regex/` package
2. Update validation spec (remove regex-related requirements)
3. Update documentation

### Rollback
If issues arise:
- Revert migration commits
- `internal/regex/` remains functional until fully migrated
- No parallel runtime overhead (one or the other, not both)

## Resolved Design Questions

1. **Reference-style links?** `[text][ref]` with `[ref]: url`
   - **Decision**: Yes - support reference-style links for flexibility. Parser will collect link definitions during first pass and resolve references during AST construction.

2. **Windows line endings (CRLF)?**
   - **Decision**: Yes - lexer will normalize CRLF to LF during tokenization while preserving accurate line numbers

3. **Lexer visibility?**
   - **Decision**: Internal only - expose only the parser API; lexer is implementation detail

4. **`internal/git/platform.go` regex?**
   - **Decision**: Out of scope - those regex patterns are for URL parsing, not markdown. Address in separate change if desired.

5. **Emphasis parsing edge cases?**
   - **Decision**: Follow CommonMark rules for emphasis: delimiter runs, left/right flanking, intraword emphasis. Implementation will handle `*foo*bar*` correctly.

6. **Error recovery strategy?**
   - **Decision**: On error, skip to next recognizable structure (next header, next list item, etc.) and continue parsing. This maximizes useful output even from malformed input.
