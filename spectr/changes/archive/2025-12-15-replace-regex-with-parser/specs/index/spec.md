# Delta Specification

## ADDED Requirements

### Requirement: Position Index Structure

The system SHALL provide a PositionIndex that enables O(log n) queries for
finding nodes at a given source position.

#### Scenario: Index creation

- **WHEN** `NewPositionIndex(root Node)` is called
- **THEN** it SHALL return a `*PositionIndex` ready for queries
- **AND** the index SHALL be built lazily on first query (not at construction)

#### Scenario: Lazy index building

- **WHEN** the first position query is made
- **THEN** the interval tree SHALL be built from the AST
- **AND** subsequent queries SHALL reuse the built index
- **AND** index construction SHALL be O(n log n) where n is node count

#### Scenario: Index storage

- **WHEN** an index is built
- **THEN** it SHALL use an interval tree data structure
- **AND** each node's [Start, End) range SHALL be stored
- **AND** the tree SHALL support overlapping intervals (parent contains
  children)

### Requirement: Position Query Operations

The system SHALL provide query methods for finding nodes at a source position.

#### Scenario: Find node at offset

- **WHEN** `index.NodeAt(offset int)` is called
- **THEN** it SHALL return the innermost (most specific) node containing that
  offset
- **AND** query time SHALL be O(log n)
- **AND** offset outside document range SHALL return nil

#### Scenario: Find all nodes at offset

- **WHEN** `index.NodesAt(offset int)` is called
- **THEN** it SHALL return all nodes containing that offset (from root to leaf)
- **AND** nodes SHALL be ordered from outermost (root) to innermost (leaf)
- **AND** query time SHALL be O(log n + k) where k is result count

#### Scenario: Find nodes in range

- **WHEN** `index.NodesInRange(start, end int)` is called
- **THEN** it SHALL return all nodes overlapping the given range
- **AND** a node overlaps if its range intersects [start, end)
- **AND** partial overlaps SHALL be included

### Requirement: Interval Tree Implementation

The system SHALL implement an efficient interval tree for position indexing.

#### Scenario: Interval tree structure

- **WHEN** the interval tree is built
- **THEN** it SHALL be a balanced binary search tree (e.g., augmented AVL or
  red-black)
- **AND** each node SHALL store: interval [start, end), max endpoint in subtree,
  AST node reference

#### Scenario: Interval tree invariants

- **WHEN** the tree is queried
- **THEN** the binary search property SHALL hold on interval start points
- **AND** the max endpoint augmentation SHALL enable efficient pruning
- **AND** the tree SHALL remain balanced after construction

#### Scenario: Query algorithm

- **WHEN** querying for nodes at offset
- **THEN** the algorithm SHALL traverse from root, pruning subtrees where max <
  offset
- **AND** it SHALL check if each visited node's interval contains the offset
- **AND** it SHALL collect all containing intervals

### Requirement: Index Invalidation

The system SHALL handle index invalidation when the AST changes.

#### Scenario: Index tied to AST version

- **WHEN** an index is created for an AST
- **THEN** it SHALL store the AST root's hash
- **AND** queries on a modified AST SHALL detect the mismatch

#### Scenario: Stale index detection

- **WHEN** the AST is modified (new root created)
- **THEN** the old index SHALL NOT be valid for the new AST
- **AND** attempting to use stale index SHALL return error or rebuild

#### Scenario: Explicit rebuild

- **WHEN** `index.Rebuild(root Node)` is called
- **THEN** it SHALL discard old tree and build new index
- **AND** subsequent queries SHALL use the new index

### Requirement: Line Index Integration

The system SHALL integrate with LineIndex for line/column calculations.

#### Scenario: Combined position and line lookup

- **WHEN** `index.PositionAt(offset int)` is called
- **THEN** it SHALL return `Position{Line, Column, Offset}`
- **AND** it SHALL use internal LineIndex for line/column calculation

#### Scenario: LineIndex caching

- **WHEN** multiple line/column queries are made
- **THEN** the LineIndex SHALL be computed once and cached
- **AND** it SHALL be built lazily on first position query

#### Scenario: Node position convenience

- **WHEN** `index.NodePosition(node Node)` is called
- **THEN** it SHALL return `Position` for the node's start offset
- **AND** this SHALL be equivalent to `PositionAt(node.Start())`

### Requirement: Range Query Utilities

The system SHALL provide utility functions for common range queries.

#### Scenario: Find enclosing section

- **WHEN** `index.EnclosingSection(offset int)` is called
- **THEN** it SHALL return the NodeSection containing that offset
- **AND** it SHALL return nil if offset is not within any section

#### Scenario: Find enclosing requirement

- **WHEN** `index.EnclosingRequirement(offset int)` is called
- **THEN** it SHALL return the NodeRequirement containing that offset
- **AND** it SHALL return nil if offset is not within any requirement

#### Scenario: Sibling navigation

- **WHEN** `index.NextSibling(node Node)` is called
- **THEN** it SHALL return the next sibling node at the same level
- **AND** it SHALL return nil if node is last sibling
- **AND** this requires parent context from the index

### Requirement: Index Memory Efficiency

The system SHALL optimize memory usage for the position index.

#### Scenario: Node references not copies

- **WHEN** the index stores nodes
- **THEN** it SHALL store references/pointers to original AST nodes
- **AND** it SHALL NOT copy node data
- **AND** index memory overhead SHALL be O(n) pointers plus tree structure

#### Scenario: Compact interval representation

- **WHEN** intervals are stored
- **THEN** start and end SHALL be stored as int (not Position structs)
- **AND** line/column SHALL be computed on demand, not stored
