// Package markdown provides a token-based lexer and parser for CommonMark subset
// plus Spectr-specific extensions. It provides a robust, maintainable, and
// feature-rich approach for parsing Spectr specification files.
//
// # Architecture Overview
//
// The package uses a two-phase parsing architecture:
//
//  1. Lexer: Tokenizes input into fine-grained tokens (delimiters, text, whitespace)
//  2. Parser: Consumes tokens to build an immutable AST
//
// This separation provides clear concerns, testable components, and flexibility
// for error recovery.
//
//	Source []byte --> Lexer --> []Token --> Parser --> AST (Node tree)
//	                              |                      |
//	                              v                      v
//	                         TokenError            ParseError
//
// # Supported CommonMark Subset
//
// The parser supports a useful subset of CommonMark:
//
//   - Headers: ATX-style headers (H1-H6) using # prefix
//   - Lists: Unordered (-, *, +), ordered (1.), task checkboxes (- [ ], - [x])
//   - Code: Fenced code blocks (``` and ~~~), inline code (`)
//   - Emphasis: Bold (**, __), italic (*, _), strikethrough (~~)
//   - Links: Inline [text](url), reference [text][ref] with [ref]: url definitions
//   - Block elements: Paragraphs, blockquotes (>)
//
// # Spectr-Specific Extensions
//
// Beyond CommonMark, the parser recognizes Spectr patterns:
//
//   - Wikilinks: [[spec-name]], [[spec-name|display text]], [[target#anchor]]
//   - Requirement headers: ### Requirement: Name
//   - Scenario headers: #### Scenario: Description
//   - WHEN/THEN/AND bullets: - **WHEN** condition, - **THEN** result
//   - Delta sections: ## ADDED Requirements, ## MODIFIED Requirements, etc.
//
// # Key Types
//
// Token represents a lexical unit with position information:
//
//	type Token struct {
//	    Type    TokenType  // Token classification
//	    Start   int        // Byte offset from source start
//	    End     int        // Byte offset past last byte (exclusive)
//	    Source  []byte     // Zero-copy slice into original source
//	    Message string     // Error message (only for TokenError)
//	}
//
// Node is the interface for all AST nodes:
//
//	type Node interface {
//	    NodeType() NodeType        // Type classification
//	    Span() (start, end int)    // Byte offset range
//	    Hash() uint64              // Content hash for identity/caching
//	    Source() []byte            // Original source text
//	    Children() []Node          // Child nodes (immutable)
//	}
//
// Typed node structs implement the Node interface with type-specific getters:
//
//   - NodeDocument: Root document container
//   - NodeSection: Header sections with Level() and Title() getters
//   - NodeRequirement: Spectr requirement with Name() getter
//   - NodeScenario: Spectr scenario with Name() getter
//   - NodeParagraph: Paragraph container
//   - NodeList: List container with Ordered() getter
//   - NodeListItem: List item with Checked() and Keyword() getters
//   - NodeCodeBlock: Fenced code with Language() and Content() getters
//   - NodeBlockquote: Blockquote container
//   - NodeText: Plain text content
//   - NodeStrong: Bold emphasis
//   - NodeEmphasis: Italic emphasis
//   - NodeStrikethrough: Strikethrough text
//   - NodeCode: Inline code
//   - NodeLink: Link with URL() and Title() getters
//   - NodeWikilink: Wikilink with Target(), Display(), and Anchor() getters
//
// ParseError represents a parse error with location:
//
//	type ParseError struct {
//	    Offset   int           // Byte offset where error occurred
//	    Message  string        // Human-readable error description
//	    Expected []TokenType   // What tokens would have been valid
//	}
//
// # Key Functions
//
// Parse is the main entry point for parsing markdown:
//
//	func Parse(source []byte) (Node, []ParseError)
//
// Parse is stateless and safe for concurrent calls. It returns the root document
// node and any errors encountered. Even with errors, a partial AST is returned.
//
// ParseIncremental enables efficient reparsing after edits:
//
//	func ParseIncremental(oldTree Node, oldSource, newSource []byte) (Node, []ParseError)
//
// ParseIncremental computes the diff between old and new source, identifies
// affected regions, reparses only changed sections, and reuses unchanged subtrees
// via content hash matching. This provides tree-sitter style incremental parsing.
//
// # Usage Examples
//
// Basic parsing:
//
//	source := []byte("# Hello\n\nThis is a **test**.")
//	root, errs := markdown.Parse(source)
//	if len(errs) > 0 {
//	    for _, err := range errs {
//	        fmt.Printf("Error at offset %d: %s\n", err.Offset, err.Message)
//	    }
//	}
//	// Process AST...
//
// Using the visitor pattern:
//
//	type RequirementCollector struct {
//	    markdown.BaseVisitor
//	    Requirements []string
//	}
//
//	func (c *RequirementCollector) VisitRequirement(n *markdown.NodeRequirement) error {
//	    c.Requirements = append(c.Requirements, n.Name())
//	    return nil
//	}
//
//	collector := &RequirementCollector{}
//	markdown.Walk(root, collector)
//
// Using queries:
//
//	// Find all requirements
//	reqs := markdown.Find(root, markdown.IsType[*markdown.NodeRequirement]())
//
//	// Find requirement by name
//	req := markdown.FindFirst(root, markdown.And(
//	    markdown.IsType[*markdown.NodeRequirement](),
//	    markdown.HasName("My Requirement"),
//	))
//
// Using transforms:
//
//	// Rename a requirement
//	newRoot, err := markdown.Transform(root, markdown.RenameRequirement("OldName", "NewName"))
//
// Incremental parsing:
//
//	// Initial parse
//	root1, _ := markdown.Parse(source1)
//
//	// After edit, reparse incrementally
//	root2, _ := markdown.ParseIncremental(root1, source1, source2)
//
// # Error Handling
//
// The parser uses a collected error approach rather than fail-fast:
//
//   - All errors are collected during parsing
//   - Parsing continues via error recovery (skip to next sync point)
//   - A partial AST is always returned, even with errors
//   - Maximum 100 errors before aborting (configurable)
//
// Errors store byte offsets. Use LineIndex for line/column conversion:
//
//	idx := markdown.NewLineIndex(source)
//	for _, err := range errs {
//	    pos := idx.PositionAt(err.Offset)
//	    fmt.Printf("Line %d, Column %d: %s\n", pos.Line, pos.Column, err.Message)
//	}
//
// # Position Information
//
// All tokens and nodes track byte offsets into the original source. This is
// compact and efficient. For line/column display, use LineIndex:
//
//	idx := markdown.NewLineIndex(source)
//	line, col := idx.LineCol(offset)  // 1-based line numbers
//
// For frequent position queries on an AST, use PositionIndex which builds an
// interval tree for O(log n) lookups:
//
//	pidx := markdown.NewPositionIndex(root)
//	node := pidx.NodeAt(offset)         // Innermost node at offset
//	nodes := pidx.NodesAt(offset)       // All nodes containing offset
//	nodes := pidx.NodesInRange(10, 50)  // All nodes overlapping range
//
// # Thread Safety
//
// Parse and ParseIncremental are safe for concurrent calls (no shared state).
// AST nodes are immutable after creation. Object pools use sync.Pool internally.
// LineIndex and PositionIndex are safe for concurrent reads after construction.
//
// # Performance
//
// The package is optimized for performance:
//
//   - Zero-copy source views ([]byte slices into original input)
//   - Object pooling for tokens and nodes via sync.Pool
//   - Content hashing for efficient subtree comparison
//   - Lazy line index construction
//   - Interval tree for O(log n) position queries
//   - Incremental parsing reuses unchanged subtrees
//
// # Non-Goals
//
// This package intentionally does not support:
//
//   - Full CommonMark compliance (focused subset for Spectr)
//   - Tables
//   - HTML passthrough
//   - Setext-style headers (underlined with === or ---)
//   - GFM extensions beyond task checkboxes and strikethrough
//
//nolint:revive // line-length-limit: documentation lines exceed 80 chars for readability
package markdown
