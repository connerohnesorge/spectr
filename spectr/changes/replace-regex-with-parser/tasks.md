## 1. Token Implementation (specs/tokens/spec.md)
- [ ] 1.1 Create `internal/markdown/doc.go` with package documentation
- [ ] 1.2 Create `internal/markdown/token.go` with fine-grained TokenType enum
- [ ] 1.3 Implement structural tokens: TokenEOF, TokenNewline, TokenWhitespace, TokenText, TokenError
- [ ] 1.4 Implement punctuation delimiter tokens: TokenHash, TokenAsterisk, TokenUnderscore, TokenTilde, TokenBacktick, TokenDash, TokenPlus, TokenDot, TokenColon, TokenPipe
- [ ] 1.5 Implement bracket tokens: TokenBracketOpen, TokenBracketClose, TokenParenOpen, TokenParenClose, TokenGreaterThan
- [ ] 1.6 Implement special tokens: TokenNumber, TokenX (for checkboxes)
- [ ] 1.7 Implement Token struct with Type, Start, End, Source ([]byte), Message fields
- [ ] 1.8 Implement TokenType.String() method for debugging
- [ ] 1.9 Create `internal/markdown/token_test.go` with token type tests

## 2. Lexer Implementation (specs/lexer/spec.md)
- [ ] 2.1 Create `internal/markdown/lexer.go` with Lexer struct holding source []byte reference
- [ ] 2.2 Implement byte offset tracking (pos field incremented by byte count)
- [ ] 2.3 Implement Next() method returning next Token and advancing position
- [ ] 2.4 Implement Peek() method returning next Token WITHOUT advancing (with caching)
- [ ] 2.5 Implement TokenEOF handling at end of source
- [ ] 2.6 Implement CRLF normalization: \r\n → single TokenNewline with End = Start + 2
- [ ] 2.7 Implement standalone \r handling as TokenNewline
- [ ] 2.8 Implement TokenError emission for unrecognized input with Message field
- [ ] 2.9 Implement error recovery: resync at newline, whitespace, or known delimiter
- [ ] 2.10 Implement invalid UTF-8 handling: emit TokenError, continue from next byte
- [ ] 2.11 Implement state machine for code fence context (StateFencedCode)
- [ ] 2.12 Implement state machine for inline code context (backtick sequence tracking)
- [ ] 2.13 Implement state machine for link URL context (StateLinkURL)
- [ ] 2.14 Implement multi-byte UTF-8 character handling (complete sequences in TokenText)
- [ ] 2.15 Implement Unicode whitespace handling (only ASCII space/tab as TokenWhitespace)
- [ ] 2.16 Implement All() method returning []Token including TokenEOF
- [ ] 2.17 Implement AllWithErrors() returning ([]Token, []LexError)
- [ ] 2.18 Create `internal/markdown/lexer_test.go` with comprehensive tests
- [ ] 2.19 Add CRLF normalization tests
- [ ] 2.20 Add error token recovery tests
- [ ] 2.21 Add state machine transition tests

## 3. AST Node Implementation (specs/ast/spec.md)
- [ ] 3.1 Create `internal/markdown/node.go` with Node interface definition
- [ ] 3.2 Define Node interface: NodeType(), Span(), Hash(), Source(), Children()
- [ ] 3.3 Create `internal/markdown/node_types.go` with typed node structs
- [ ] 3.4 Implement NodeDocument struct with private fields and getters
- [ ] 3.5 Implement NodeSection struct with level, title fields and Level(), Title(), DeltaType() getters
- [ ] 3.6 Implement NodeRequirement struct with name field and Name() getter
- [ ] 3.7 Implement NodeScenario struct with name field and Name() getter
- [ ] 3.8 Implement NodeParagraph struct
- [ ] 3.9 Implement NodeList struct with ordered field and Ordered() getter
- [ ] 3.10 Implement NodeListItem struct with checked, keyword fields and Checked(), Keyword() getters
- [ ] 3.11 Implement NodeCodeBlock struct with language, content fields and Language(), Content() getters
- [ ] 3.12 Implement NodeBlockquote struct
- [ ] 3.13 Implement NodeText struct
- [ ] 3.14 Implement NodeStrong struct
- [ ] 3.15 Implement NodeEmphasis struct
- [ ] 3.16 Implement NodeStrikethrough struct
- [ ] 3.17 Implement NodeCode struct
- [ ] 3.18 Implement NodeLink struct with url, title fields and URL(), Title() getters
- [ ] 3.19 Implement NodeWikilink struct with target, display, anchor fields and getters
- [ ] 3.20 Implement content hash computation (uint64 using fnv or xxhash)
- [ ] 3.21 Implement hash from: NodeType + children hashes + text content hash
- [ ] 3.22 Implement Source []byte field as zero-copy slice into original input
- [ ] 3.23 Implement Start/End byte offset fields
- [ ] 3.24 Implement Children []Node as immutable copy
- [ ] 3.25 Implement NodeBuilder for programmatic node construction
- [ ] 3.26 Implement node.ToBuilder() for transformation pipelines
- [ ] 3.27 Implement builder validation: Start <= End, proper nesting
- [ ] 3.28 Implement node.Span() returning (start, end int) byte offsets
- [ ] 3.29 Implement node.Equal(other) for deep structural comparison
- [ ] 3.30 Create `internal/markdown/node_test.go` with immutability and hashing tests

## 4. Line Index Implementation
- [ ] 4.1 Create `internal/markdown/lineindex.go` with LineIndex struct
- [ ] 4.2 Implement lazy line-start-offsets index construction
- [ ] 4.3 Implement binary search for O(log n) LineCol lookup
- [ ] 4.4 Implement LineCol(offset) returning (line, col int) with 1-based line numbers
- [ ] 4.5 Implement PositionAt(offset) returning Position struct
- [ ] 4.6 Create `internal/markdown/lineindex_test.go` with position tests

## 5. Parser Implementation (specs/parser/spec.md)
- [ ] 5.1 Create `internal/markdown/parser.go` with stateless Parse function
- [ ] 5.2 Implement Parse(source []byte) returning (Node, []ParseError)
- [ ] 5.3 Ensure Parse is safe for concurrent calls (no shared state)
- [ ] 5.4 Implement internal object pooling via sync.Pool
- [ ] 5.5 Implement ParseError struct with Offset, Message, Expected []TokenType
- [ ] 5.6 Implement error collection mode: add ParseError, attempt recovery, continue
- [ ] 5.7 Implement error recovery: skip to sync points (blank line, header, list marker)
- [ ] 5.8 Implement maximum errors limit (default: 100) with abort
- [ ] 5.9 Implement block-level parsing: code fence, header, blockquote, list item, paragraph
- [ ] 5.10 Implement header parsing: 1-6 TokenHash + TokenWhitespace → NodeSection
- [ ] 5.11 Implement Spectr Requirement: and Scenario: header detection
- [ ] 5.12 Implement list parsing with nesting via indentation
- [ ] 5.13 Implement checkbox syntax parsing for task items
- [ ] 5.14 Implement code fence parsing: 3+ backticks/tildes, verbatim content
- [ ] 5.15 Implement paragraph parsing: consecutive non-blank lines
- [ ] 5.16 Implement inline parsing with delimiter stack for emphasis
- [ ] 5.17 Implement CommonMark strict left/right flanking rules
- [ ] 5.18 Implement underscore intraword restriction (foo_bar_baz not emphasis)
- [ ] 5.19 Implement overlapping delimiter resolution per CommonMark spec
- [ ] 5.20 Implement inline code parsing with backtick sequence matching
- [ ] 5.21 Implement link parsing: [text](url) and [text][ref]
- [ ] 5.22 Implement wikilink parsing: [[target|display#anchor]]
- [ ] 5.23 Implement link definition collection (first pass): [label]: url "title"
- [ ] 5.24 Implement reference link resolution (second pass) with case-insensitive labels
- [ ] 5.25 Implement shortcut reference links: [text][] and [text]
- [ ] 5.26 Implement delta section detection: ADDED|MODIFIED|REMOVED|RENAMED Requirements
- [ ] 5.27 Implement WHEN/THEN/AND keyword extraction in list items
- [ ] 5.28 Create `internal/markdown/parser_test.go` with comprehensive tests
- [ ] 5.29 Add CommonMark emphasis edge case tests

## 6. Incremental Parsing (specs/parser/spec.md)
- [ ] 6.1 Create `internal/markdown/incremental.go` with ParseIncremental function
- [ ] 6.2 Implement ParseIncremental(oldTree, oldSource, newSource) returning (Node, []ParseError)
- [ ] 6.3 Implement diff computation between oldSource and newSource
- [ ] 6.4 Implement single edit optimization: O(n) prefix/suffix matching
- [ ] 6.5 Implement multi-edit fallback: Myers diff for complex cases
- [ ] 6.6 Implement EditRegion struct: StartOffset, OldEndOffset, NewEndOffset
- [ ] 6.7 Implement affected region identification
- [ ] 6.8 Implement subtree reuse via content hash matching
- [ ] 6.9 Implement offset adjustment for nodes after edit point
- [ ] 6.10 Implement link definition and line index reuse
- [ ] 6.11 Create `internal/markdown/incremental_test.go` with incremental parsing tests

## 7. Visitor Implementation (specs/visitor/spec.md)
- [ ] 7.1 Create `internal/markdown/visitor.go` with Visitor interface
- [ ] 7.2 Implement Visit methods for all node types via interface
- [ ] 7.3 Implement Walk(node, visitor) for pre-order depth-first traversal
- [ ] 7.4 Implement SkipChildren sentinel error for skipping subtrees
- [ ] 7.5 Implement Walk with nil node returning nil (safe handling)
- [ ] 7.6 Implement BaseVisitor with no-op defaults for embedding
- [ ] 7.7 Implement VisitorContext for passing parent information during traversal
- [ ] 7.8 Implement ctx.Parent() to access parent during visitation
- [ ] 7.9 Implement EnterLeaveVisitor interface with EnterX/LeaveX methods
- [ ] 7.10 Implement WalkEnterLeave(node, visitor) with pre/post visitation
- [ ] 7.11 Implement BaseEnterLeaveVisitor with no-op defaults
- [ ] 7.12 Create `internal/markdown/visitor_test.go` with traversal tests

## 8. Transform Implementation (specs/transform/spec.md)
- [ ] 8.1 Create `internal/markdown/transform.go` with TransformAction type
- [ ] 8.2 Implement TransformAction enum: ActionKeep, ActionReplace, ActionDelete
- [ ] 8.3 Implement TransformVisitor interface with TransformX methods
- [ ] 8.4 Implement Transform(root, visitor) returning (Node, error)
- [ ] 8.5 Implement post-order traversal for transforms (children before parent)
- [ ] 8.6 Implement ActionDelete: remove node from parent's children
- [ ] 8.7 Implement ActionReplace: substitute returned node
- [ ] 8.8 Implement BaseTransformVisitor with defaults returning (original, ActionKeep, nil)
- [ ] 8.9 Implement Compose(t1, t2) for sequential transform composition
- [ ] 8.10 Implement Pipeline(transforms...) for transform chains
- [ ] 8.11 Implement When(pred, transform) for conditional transforms
- [ ] 8.12 Implement Map(f func(Node) Node) utility transform
- [ ] 8.13 Implement Filter(pred) utility transform
- [ ] 8.14 Implement RenameRequirement(oldName, newName) helper
- [ ] 8.15 Implement AddScenario(reqName, scenario) helper
- [ ] 8.16 Implement RemoveRequirement(name) helper
- [ ] 8.17 Create `internal/markdown/transform_test.go` with transform tests

## 9. Query Implementation (specs/query/spec.md)
- [ ] 9.1 Create `internal/markdown/query.go` with Find function
- [ ] 9.2 Implement Find(root, pred func(Node) bool) returning []Node
- [ ] 9.3 Implement FindFirst(root, pred) returning Node (nil if none)
- [ ] 9.4 Implement FindByType[T Node](root) returning []*T
- [ ] 9.5 Implement FindFirstByType[T Node](root) returning *T
- [ ] 9.6 Implement And(p1, p2) predicate combinator
- [ ] 9.7 Implement Or(p1, p2) predicate combinator
- [ ] 9.8 Implement Not(p) predicate combinator
- [ ] 9.9 Implement All(preds...) predicate combinator
- [ ] 9.10 Implement Any(preds...) predicate combinator
- [ ] 9.11 Implement IsType[T]() predicate factory
- [ ] 9.12 Implement HasName(name) predicate factory
- [ ] 9.13 Implement InRange(start, end) predicate factory
- [ ] 9.14 Implement HasChild(pred) predicate factory
- [ ] 9.15 Implement HasDescendant(pred) predicate factory
- [ ] 9.16 Implement Count(root, pred) returning int without slice allocation
- [ ] 9.17 Implement Exists(root, pred) returning bool with short-circuit
- [ ] 9.18 Create `internal/markdown/query_test.go` with query tests

## 10. Position Index Implementation (specs/index/spec.md)
- [ ] 10.1 Create `internal/markdown/index.go` with PositionIndex struct
- [ ] 10.2 Implement NewPositionIndex(root) constructor
- [ ] 10.3 Implement lazy index building on first query
- [ ] 10.4 Implement interval tree data structure (augmented BST)
- [ ] 10.5 Implement index.NodeAt(offset) returning innermost node
- [ ] 10.6 Implement index.NodesAt(offset) returning all nodes containing offset
- [ ] 10.7 Implement index.NodesInRange(start, end) returning overlapping nodes
- [ ] 10.8 Implement O(log n) query algorithm with subtree pruning
- [ ] 10.9 Implement index.Rebuild(root) for explicit rebuilding
- [ ] 10.10 Implement stale index detection via AST root hash
- [ ] 10.11 Integrate LineIndex for line/column calculations
- [ ] 10.12 Implement index.PositionAt(offset) returning Position
- [ ] 10.13 Implement index.EnclosingSection(offset) utility
- [ ] 10.14 Implement index.EnclosingRequirement(offset) utility
- [ ] 10.15 Create `internal/markdown/index_test.go` with index tests

## 11. Printer Implementation (specs/printer/spec.md)
- [ ] 11.1 Create `internal/markdown/printer.go` with Print function
- [ ] 11.2 Implement Print(node Node) returning []byte
- [ ] 11.3 Implement PrintTo(w io.Writer, node Node) returning error
- [ ] 11.4 Implement minimal whitespace formatting (single space/newline)
- [ ] 11.5 Implement header printing with ATX style (# prefix)
- [ ] 11.6 Implement Requirement: and Scenario: header formatting
- [ ] 11.7 Implement unordered list printing with dash bullets
- [ ] 11.8 Implement ordered list printing with 1-based numbering
- [ ] 11.9 Implement checkbox list item printing
- [ ] 11.10 Implement nested list indentation (2 spaces)
- [ ] 11.11 Implement fenced code block printing with triple backticks
- [ ] 11.12 Implement strong emphasis printing with **
- [ ] 11.13 Implement emphasis printing with *
- [ ] 11.14 Implement strikethrough printing with ~~
- [ ] 11.15 Implement inline code printing with backticks
- [ ] 11.16 Implement inline link printing: [text](url)
- [ ] 11.17 Implement wikilink printing: [[target|display]]
- [ ] 11.18 Implement blockquote printing with > prefix
- [ ] 11.19 Implement delta section header printing
- [ ] 11.20 Implement WHEN/THEN/AND keyword printing in bold
- [ ] 11.21 Create `internal/markdown/printer_test.go` with printer tests

## 12. Pool Implementation (specs/pool/spec.md)
- [ ] 12.1 Create `internal/markdown/pool.go` with pool types
- [ ] 12.2 Implement token pool using sync.Pool
- [ ] 12.3 Implement pool.GetToken() and pool.PutToken()
- [ ] 12.4 Implement typed node pools for each node type
- [ ] 12.5 Implement pool.GetSection(), pool.GetRequirement(), etc.
- [ ] 12.6 Implement pool.PutNode(n Node) with type routing
- [ ] 12.7 Implement children slice pool with size buckets
- [ ] 12.8 Implement pool.GetChildren(capacity) and pool.PutChildren()
- [ ] 12.9 Implement field clearing before return to pool
- [ ] 12.10 Implement optional pool statistics tracking
- [ ] 12.11 Create `internal/markdown/pool_test.go` with pool tests

## 13. High-Level API Implementation
- [ ] 13.1 Create `internal/markdown/api.go` with public API functions
- [ ] 13.2 Implement ParseSpec(content) (*Spec, []ParseError)
- [ ] 13.3 Implement ExtractSections(content) map[string]Section
- [ ] 13.4 Implement ExtractRequirements(content) []Requirement
- [ ] 13.5 Implement FindSection(content, name) (*Section, bool)
- [ ] 13.6 Implement ExtractWikilinks(content) []Wikilink using visitor
- [ ] 13.7 Create `internal/markdown/api_test.go` with integration tests

## 14. Wikilink Resolution
- [ ] 14.1 Create `internal/markdown/wikilink.go` for wikilink types and resolution
- [ ] 14.2 Implement ResolveWikilink(target, projectRoot) (path string, exists bool)
- [ ] 14.3 Implement resolution for spec targets (spectr/specs/{target}/spec.md)
- [ ] 14.4 Implement resolution for change targets (spectr/changes/{target}/proposal.md)
- [ ] 14.5 Implement anchor resolution within target files
- [ ] 14.6 Create `internal/markdown/wikilink_test.go` with resolution tests

## 15. Delta Parsing Implementation
- [ ] 15.1 Create `internal/markdown/delta.go` for delta-specific parsing
- [ ] 15.2 Implement ParseDelta(content) (*Delta, []ParseError)
- [ ] 15.3 Implement FindDeltaSection(content, deltaType) string
- [ ] 15.4 Implement RENAMED FROM/TO pair parsing
- [ ] 15.5 Create `internal/markdown/delta_test.go` with delta tests

## 16. Compatibility Helpers
- [ ] 16.1 Create `internal/markdown/compat.go` for old API compatibility
- [ ] 16.2 Implement MatchRequirementHeader(line) (string, bool)
- [ ] 16.3 Implement MatchScenarioHeader(line) (string, bool)
- [ ] 16.4 Implement IsH2Header(line) bool
- [ ] 16.5 Implement IsH3Header(line) bool
- [ ] 16.6 Implement MatchTaskCheckbox(line) (rune, bool)
- [ ] 16.7 Implement MatchNumberedTask(line) (*NumberedTaskMatch, bool)
- [ ] 16.8 Implement MatchNumberedSection(line) (string, bool)
- [ ] 16.9 Create `internal/markdown/compat_test.go` with compatibility tests

## 17. Migration - Parsers Package
- [ ] 17.1 Update `internal/parsers/requirement_parser.go` to use markdown package
- [ ] 17.2 Update `internal/parsers/delta_parser.go` to use markdown package
- [ ] 17.3 Run existing parsers tests to verify compatibility
- [ ] 17.4 Remove regex import from parsers package

## 18. Migration - Validation Package
- [ ] 18.1 Update `internal/validation/parser.go` to use markdown package
- [ ] 18.2 Remove inline regex patterns from ContainsShallOrMust
- [ ] 18.3 Remove inline regex from NormalizeRequirementName
- [ ] 18.4 Update `internal/validation/change_rules.go` if needed
- [ ] 18.5 Add wikilink validation (check targets exist)
- [ ] 18.6 Run existing validation tests to verify compatibility

## 19. Migration - Archive Package
- [ ] 19.1 Update `internal/archive/spec_merger.go` to use markdown package
- [ ] 19.2 Replace multiNewline regex with string-based normalization
- [ ] 19.3 Run existing archive tests to verify compatibility

## 20. Migration - Cmd Package
- [ ] 20.1 Update `cmd/accept.go` to use markdown package
- [ ] 20.2 Run accept command tests to verify compatibility

## 21. Performance and Benchmarking
- [ ] 21.1 Create `internal/markdown/benchmark_test.go` with benchmark tests
- [ ] 21.2 Benchmark lexer tokenization throughput
- [ ] 21.3 Benchmark parser AST construction
- [ ] 21.4 Benchmark incremental parsing with small edits
- [ ] 21.5 Benchmark position index queries
- [ ] 21.6 Benchmark object pool hit rates
- [ ] 21.7 Compare performance with old regex implementation
- [ ] 21.8 Optimize hot paths identified by benchmarks

## 22. Cleanup and Documentation
- [ ] 22.1 Delete `internal/regex/` package entirely
- [ ] 22.2 Update `spectr/specs/validation/spec.md` to remove regex requirements
- [ ] 22.3 Run `go mod tidy` to clean dependencies
- [ ] 22.4 Run full test suite (`go test ./...`)
- [ ] 22.5 Run linter (`golangci-lint run`)
- [ ] 22.6 Update any documentation referencing regex package
- [ ] 22.7 Add package documentation with usage examples
