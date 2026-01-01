# Markdown Package

Custom two-phase lexer+parser for CommonMark subset + Spectr extensions. NOT goldmark.

## OVERVIEW
Two-phase architecture: tokenizes input → builds immutable AST. Separation enables testing, error recovery, and incremental reparsing (tree-sitter-style).

## STRUCTURE
```
internal/markdown/
├── lexer.go              # Tokenizer (delimiters, text, whitespace)
├── parser.go             # AST builder from tokens
├── node.go              # AST node types and interfaces
├── visitor.go           # Visitor pattern support
├── query.go             # AST query utilities
├── transform.go         # AST modification utilities
├── incremental.go        # Diff-based incremental parsing
├── compat.go            # Legacy compatibility layer
├── delta.go             # Delta spec parsing
├── wikilink.go          # Wikilink parsing [[target|text]]
├── lineindex.go         # Line/column conversion
├── positionindex.go     # Interval tree for O(log n) queries
└── *_test.go            # Comprehensive test coverage
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Parse markdown | Parse() in api.go | Main entry point |
| Incremental reparse | ParseIncremental() | Reuses unchanged subtrees |
| Find nodes | Find(), FindFirst() | Query utilities |
| Visit nodes | Walk() with Visitor | Visitor pattern |
| Transform AST | Transform() | Apply modifications |
| Position info | LineIndex, PositionIndex | Line/col conversion |

## CONVENTIONS
- **Zero-copy source**: Tokens store []byte slices into original input
- **Immutable AST**: Nodes immutable after creation, safe for concurrent reads
- **Content hashing**: Hash() on nodes enables subtree comparison
- **Collected errors**: Parser continues past errors, returns up to 100
- **Thread-safe**: Parse() and ParseIncremental() safe for concurrent calls

## UNIQUE TO THIS PACKAGE
- **Spectr extensions**: Wikilinks [[target]], requirement headers `### Requirement:`, scenario `#### Scenario:`, WHEN/THEN bullets
- **Delta operations**: Recognizes `## ADDED|MODIFIED|REMOVED|RENAMED Requirements`
- **Incremental parsing**: ParseIncremental() computes source diff, reparses only changed sections, reuses unchanged subtrees via hash matching

## KEY TYPES
- **Token**: Lexical unit with Type, Start/End offsets, Source slice, Message (errors)
- **Node interface**: NodeType(), Span(), Hash(), Source(), Children()
- **NodeDocument/Section/Requirement/Scenario/...**: Typed node implementations with specific getters

## COMMON PATTERNS
```go
// Basic parse
root, errs := markdown.Parse(source)

// Visitor pattern
type Collector struct { markdown.BaseVisitor }
func (c *Collector) VisitRequirement(n *markdown.NodeRequirement) error { ... }
markdown.Walk(root, &Collector)

// Query
reqs := markdown.Find(root, markdown.IsType[*markdown.NodeRequirement]())

// Transform
newRoot := markdown.Transform(root, markdown.RenameRequirement("Old", "New"))
```

## PERFORMANCE
- Object pooling for tokens/nodes (sync.Pool)
- Lazy line index construction
- O(log n) position queries via interval tree
- Incremental parsing avoids full reparse on edits
