# Delta Specification

## ADDED Requirements

### Requirement: Immutable Node Structure

The system SHALL define immutable AST nodes that cannot be modified after
creation, using builder functions for transformations.

#### Scenario: Node fields are read-only

- **WHEN** a Node is created
- **THEN** all exported fields SHALL be read-only (no setter methods)
- **AND** the Children slice SHALL be a copy, not the original
- **AND** modifications SHALL require creating a new node via builder

#### Scenario: Node creation via constructor

- **WHEN** creating a new node
- **THEN** it SHALL be done via `NewNode(type NodeType, opts ...NodeOption)` or
  specific constructors
- **AND** the constructor SHALL compute and store the content hash
- **AND** the node SHALL be fully initialized before returning

#### Scenario: Node transformation via builder

- **WHEN** modifying a node (e.g., adding child)
- **THEN** the caller SHALL use `node.With(opts ...NodeOption)` to create a
  modified copy
- **AND** the original node SHALL remain unchanged
- **AND** unchanged children MAY be shared (structural sharing)

### Requirement: Node Content Hashing

The system SHALL compute and store a content hash for each node to enable
identity tracking and caching across incremental parses.

#### Scenario: Hash computation

- **WHEN** a node is created
- **THEN** its hash SHALL be computed from: NodeType + children hashes + text
  content hash
- **AND** the hash SHALL be stored in `node.Hash` field
- **AND** the hash type SHALL be `uint64` using a fast hash (e.g., xxhash, fnv)

#### Scenario: Hash equality implies content equality

- **WHEN** two nodes have the same hash
- **THEN** they SHALL have the same semantic content (with high probability)
- **AND** the hash SHALL be usable as a map key for caching

#### Scenario: Hash includes structure

- **WHEN** computing node hash
- **THEN** it SHALL include the recursive structure (child hashes)
- **AND** changing a deep child SHALL change all ancestor hashes
- **AND** this enables detecting subtree changes

### Requirement: Node Byte Slice Source Preservation

The system SHALL store the original source text as a byte slice view for
round-trip capability.

#### Scenario: Source slice storage

- **WHEN** a node is created
- **THEN** it SHALL store `Source []byte` as a slice into the original input
- **AND** `Source` SHALL span from the node's start to end byte offset
- **AND** `Source` SHALL NOT be a copy (zero-copy)

#### Scenario: Source enables round-trip

- **WHEN** reconstructing source from AST
- **THEN** leaf nodes SHALL return their `Source` directly
- **AND** formatting and whitespace SHALL be preserved exactly
- **AND** only modified subtrees need regeneration

#### Scenario: Source lifetime

- **WHEN** a node's `Source` slice is accessed
- **THEN** it SHALL remain valid as long as the original source buffer is
  retained
- **AND** the parser SHALL document this lifetime requirement

### Requirement: Node Type Hierarchy

The system SHALL define a node type hierarchy representing markdown document
structure.

#### Scenario: Document node as root

- **WHEN** parsing completes
- **THEN** the result SHALL be a `NodeDocument` containing all top-level nodes
- **AND** `NodeDocument` SHALL have no parent

#### Scenario: Block-level node types

- **WHEN** block elements are parsed
- **THEN** they SHALL be represented as:
  - `NodeSection` for H2 headers and their content
  - `NodeRequirement` for H3 Requirement: headers
  - `NodeScenario` for H4 Scenario: headers
  - `NodeParagraph` for paragraph content
  - `NodeList` for unordered/ordered lists
  - `NodeListItem` for individual list items
  - `NodeCodeBlock` for fenced code blocks
  - `NodeBlockquote` for blockquoted content

#### Scenario: Inline node types

- **WHEN** inline formatting is parsed
- **THEN** they SHALL be represented as:
  - `NodeText` for plain text content
  - `NodeStrong` for bold/strong emphasis
  - `NodeEmphasis` for italic emphasis
  - `NodeStrikethrough` for struck text
  - `NodeCode` for inline code
  - `NodeLink` for links (inline and reference)
  - `NodeWikilink` for wikilinks

#### Scenario: Flat children array for inline

- **WHEN** a paragraph contains mixed content like `Hello **world**!`
- **THEN** NodeParagraph.Children SHALL be: [NodeText("Hello "),
  NodeStrong([NodeText("world")]), NodeText("!")]
- **AND** the array SHALL be flat at each level (no unnecessary nesting)

### Requirement: Node Position Information

The system SHALL track byte offsets for each node, with line/column calculable
on demand.

#### Scenario: Byte offset storage

- **WHEN** a node is created
- **THEN** it SHALL have `Start` and `End` fields (byte offsets)
- **AND** `Start` SHALL be the offset of the first byte of the node
- **AND** `End` SHALL be the offset past the last byte (exclusive)

#### Scenario: Span calculation

- **WHEN** `node.Span()` is called
- **THEN** it SHALL return `(start, end int)` byte offsets
- **AND** `end - start` SHALL equal `len(node.Source)`

#### Scenario: Position calculation on demand

- **WHEN** `node.Position(lineIndex *LineIndex)` is called
- **THEN** it SHALL return `Position{Line, Column, Offset}` for the start
- **AND** line/column calculation SHALL use the provided line index
- **AND** if lineIndex is nil, it SHALL be computed from source

### Requirement: Node Type-Specific Fields

The system SHALL provide type-specific data for nodes beyond the common fields.

#### Scenario: Section node fields

- **WHEN** a `NodeSection` is created
- **THEN** it SHALL have `Level int` (1-6 for H1-H6)
- **AND** it SHALL have `Title []byte` for the header text

#### Scenario: Requirement node fields

- **WHEN** a `NodeRequirement` is created
- **THEN** it SHALL have `Name string` extracted from "### Requirement: Name"
- **AND** it SHALL have scenarios as children

#### Scenario: Code block fields

- **WHEN** a `NodeCodeBlock` is created
- **THEN** it SHALL have `Language []byte` for the info string (may be nil)
- **AND** it SHALL have `Content []byte` for the code content (without fences)

#### Scenario: Link node fields

- **WHEN** a `NodeLink` is created
- **THEN** it SHALL have `URL []byte` for the link destination
- **AND** it SHALL have `Title []byte` for the optional title (may be nil)
- **AND** children SHALL be the link text nodes

#### Scenario: Wikilink node fields

- **WHEN** a `NodeWikilink` is created
- **THEN** it SHALL have `Target []byte` for the link target
- **AND** it SHALL have `Display []byte` for optional display text (may be nil)
- **AND** it SHALL have `Anchor []byte` for optional anchor (may be nil)

### Requirement: Node Builder API

The system SHALL provide a builder API for constructing and transforming nodes.

#### Scenario: Builder for new nodes

- **WHEN** constructing a node programmatically
- **THEN** `NewNodeBuilder(NodeType)` SHALL return a builder
- **AND** builder SHALL have methods: `WithChildren()`, `WithSource()`,
  `WithStart()`, `WithEnd()`
- **AND** `Build()` SHALL return the immutable node

#### Scenario: Builder for transformations

- **WHEN** transforming an existing node
- **THEN** `node.ToBuilder()` SHALL return a builder pre-populated with node
  data
- **AND** modifications SHALL be made via builder methods
- **AND** `Build()` SHALL return a new node, original unchanged

#### Scenario: Builder validates consistency

- **WHEN** `Build()` is called
- **THEN** it SHALL validate that Start <= End
- **AND** it SHALL validate children are properly nested
- **AND** it SHALL compute the content hash

### Requirement: Node Equality and Comparison

The system SHALL support node equality comparison via content hash and
structural comparison.

#### Scenario: Hash-based equality

- **WHEN** comparing nodes with `node1.Hash == node2.Hash`
- **THEN** equal hashes SHALL indicate same content (with collision probability)
- **AND** this SHALL be the fast path for equality checking

#### Scenario: Deep equality

- **WHEN** `node1.Equal(node2)` is called
- **THEN** it SHALL perform deep structural comparison
- **AND** it SHALL compare: Type, children recursively, source content
- **AND** it SHALL return true only if semantically identical

### Requirement: Typed Node Structs

The system SHALL implement separate typed structs for each node type, all
implementing a common Node interface.

#### Scenario: Node interface definition

- **WHEN** the Node interface is defined
- **THEN** it SHALL include common methods: `NodeType() NodeType`, `Span() (int,
  int)`, `Hash() uint64`, `Source() []byte`, `Children() []Node`
- **AND** all concrete node types SHALL implement this interface

#### Scenario: Concrete node structs

- **WHEN** node types are implemented
- **THEN** each SHALL be a separate struct: `NodeDocument`, `NodeSection`,
  `NodeRequirement`, `NodeScenario`, `NodeParagraph`, `NodeList`,
  `NodeListItem`, `NodeCodeBlock`, `NodeBlockquote`, `NodeText`, `NodeStrong`,
  `NodeEmphasis`, `NodeStrikethrough`, `NodeCode`, `NodeLink`, `NodeWikilink`
- **AND** type assertions SHALL be used to access type-specific data

#### Scenario: Type-safe access via assertion

- **WHEN** accessing type-specific data
- **THEN** callers SHALL use type assertion: `if section, ok :=
  node.(*NodeSection); ok { ... }`
- **AND** this provides compile-time type safety for type-specific operations

### Requirement: Getter Methods for Type-Specific Data

The system SHALL use private fields with getter methods for type-specific data
to enforce immutability.

#### Scenario: Private fields with getters

- **WHEN** type-specific data is accessed
- **THEN** fields SHALL be private (lowercase): `level`, `title`, `url`, etc.
- **AND** getter methods SHALL provide read access: `Level() int`, `Title()
  []byte`, etc.
- **AND** no setter methods SHALL be provided

#### Scenario: NodeSection getters

- **WHEN** accessing NodeSection data
- **THEN** `section.Level()` SHALL return the header level (1-6)
- **AND** `section.Title()` SHALL return the header text as []byte
- **AND** `section.DeltaType()` SHALL return the delta type if present (ADDED,
  MODIFIED, etc.)

#### Scenario: NodeRequirement getters

- **WHEN** accessing NodeRequirement data
- **THEN** `req.Name()` SHALL return the requirement name as string
- **AND** scenarios SHALL be accessible via `req.Children()`

#### Scenario: NodeScenario getters

- **WHEN** accessing NodeScenario data
- **THEN** `scenario.Name()` SHALL return the scenario name as string

#### Scenario: NodeCodeBlock getters

- **WHEN** accessing NodeCodeBlock data
- **THEN** `block.Language()` SHALL return the language identifier as []byte
  (may be nil)
- **AND** `block.Content()` SHALL return the code content as []byte (without
  fences)

#### Scenario: NodeLink getters

- **WHEN** accessing NodeLink data
- **THEN** `link.URL()` SHALL return the link destination as []byte
- **AND** `link.Title()` SHALL return the optional title as []byte (may be nil)

#### Scenario: NodeWikilink getters

- **WHEN** accessing NodeWikilink data
- **THEN** `wikilink.Target()` SHALL return the link target as []byte
- **AND** `wikilink.Display()` SHALL return optional display text as []byte (may
  be nil)
- **AND** `wikilink.Anchor()` SHALL return optional anchor as []byte (may be
  nil)

#### Scenario: NodeListItem getters

- **WHEN** accessing NodeListItem data
- **THEN** `item.Checked()` SHALL return (bool, bool) for checkbox state and
  presence
- **AND** `item.Keyword()` SHALL return optional WHEN/THEN/AND keyword as string

### Requirement: No Parent Pointers

The system SHALL NOT store parent pointers in nodes to maintain simplicity and
true immutability.

#### Scenario: Children-only references

- **WHEN** a node is created
- **THEN** it SHALL only store references to children
- **AND** it SHALL NOT store a reference to its parent
- **AND** this avoids cycles and simplifies immutability

#### Scenario: Parent access via visitor context

- **WHEN** traversing the AST with a visitor
- **THEN** parent information SHALL be available via VisitorContext
- **AND** the context SHALL maintain the path from root to current node
- **AND** callers SHALL use `ctx.Parent()` to access the parent during
  visitation

#### Scenario: Upward navigation without parent pointers

- **WHEN** upward navigation is needed outside of visitation
- **THEN** callers SHALL use the PositionIndex or re-traverse from root
- **AND** this trade-off favors simplicity over convenience
