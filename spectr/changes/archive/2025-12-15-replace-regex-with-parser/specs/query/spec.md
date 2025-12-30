# Delta Specification

## ADDED Requirements

### Requirement: Predicate-Based Find

The system SHALL provide a Find function for locating nodes matching a
predicate.

#### Scenario: Find function signature

- **WHEN** finding nodes by predicate
- **THEN** the signature SHALL be `Find(root Node, pred func(Node) bool) []Node`
- **AND** it SHALL return all nodes where pred returns true
- **AND** results SHALL be in pre-order traversal order

#### Scenario: Find traverses entire tree

- **WHEN** Find is called
- **THEN** it SHALL visit every node in the AST
- **AND** it SHALL collect all nodes where pred returns true
- **AND** traversal SHALL be depth-first pre-order

#### Scenario: Find returns empty for no matches

- **WHEN** no nodes match the predicate
- **THEN** Find SHALL return an empty slice (not nil)
- **AND** the slice SHALL have length 0

### Requirement: FindFirst Function

The system SHALL provide a FindFirst function for locating the first matching
node.

#### Scenario: FindFirst function signature

- **WHEN** finding first matching node
- **THEN** the signature SHALL be `FindFirst(root Node, pred func(Node) bool)
  Node`
- **AND** it SHALL return the first node where pred returns true
- **AND** it SHALL return nil if no match found

#### Scenario: FindFirst short-circuits

- **WHEN** a match is found
- **THEN** FindFirst SHALL stop traversal immediately
- **AND** remaining nodes SHALL NOT be visited
- **AND** this enables efficient early termination

### Requirement: Type-Safe Find Functions

The system SHALL provide generic Find functions for type-safe queries.

#### Scenario: FindByType generic function

- **WHEN** `FindByType[T Node](root Node) []*T` is called
- **THEN** it SHALL return all nodes of type T
- **AND** results SHALL be properly typed (no interface{} casting needed)
- **AND** results SHALL be in traversal order

#### Scenario: FindFirstByType generic function

- **WHEN** `FindFirstByType[T Node](root Node) *T` is called
- **THEN** it SHALL return the first node of type T
- **AND** it SHALL return nil if no node of type T exists

#### Scenario: Type parameter constraint

- **WHEN** T is specified
- **THEN** T SHALL be constrained to types implementing Node interface
- **AND** compile-time error SHALL occur for invalid type parameters

### Requirement: Predicate Combinators

The system SHALL provide combinators for building complex predicates.

#### Scenario: And combinator

- **WHEN** `And(p1, p2 func(Node) bool)` is called
- **THEN** it SHALL return a predicate that is true when both p1 AND p2 are true
- **AND** short-circuit evaluation SHALL apply (p2 not called if p1 is false)

#### Scenario: Or combinator

- **WHEN** `Or(p1, p2 func(Node) bool)` is called
- **THEN** it SHALL return a predicate that is true when p1 OR p2 is true
- **AND** short-circuit evaluation SHALL apply (p2 not called if p1 is true)

#### Scenario: Not combinator

- **WHEN** `Not(p func(Node) bool)` is called
- **THEN** it SHALL return a predicate that negates p
- **AND** `Not(p)(node)` equals `!p(node)`

#### Scenario: All combinator

- **WHEN** `All(preds ...func(Node) bool)` is called
- **THEN** it SHALL return a predicate true when all preds are true
- **AND** it SHALL short-circuit on first false

#### Scenario: Any combinator

- **WHEN** `Any(preds ...func(Node) bool)` is called
- **THEN** it SHALL return a predicate true when any pred is true
- **AND** it SHALL short-circuit on first true

### Requirement: Common Predicate Factories

The system SHALL provide factory functions for common predicates.

#### Scenario: IsType predicate

- **WHEN** `IsType[T Node]()` is called
- **THEN** it SHALL return a predicate that matches nodes of type T
- **AND** type checking SHALL use Go type assertion

#### Scenario: HasName predicate

- **WHEN** `HasName(name string)` is called
- **THEN** it SHALL return a predicate matching nodes with Name() == name
- **AND** it SHALL work for NodeRequirement, NodeScenario, etc.

#### Scenario: InRange predicate

- **WHEN** `InRange(start, end int)` is called
- **THEN** it SHALL return a predicate matching nodes within [start, end)
- **AND** a node is in range if its span overlaps the given range

#### Scenario: HasChild predicate

- **WHEN** `HasChild(pred func(Node) bool)` is called
- **THEN** it SHALL return a predicate true if any direct child matches pred
- **AND** only immediate children SHALL be checked (not descendants)

#### Scenario: HasDescendant predicate

- **WHEN** `HasDescendant(pred func(Node) bool)` is called
- **THEN** it SHALL return a predicate true if any descendant matches pred
- **AND** all descendants SHALL be checked recursively

### Requirement: Query Result Operations

The system SHALL provide operations on query results.

#### Scenario: Map over results

- **WHEN** results need to be transformed
- **THEN** standard Go slice operations SHALL be used
- **AND** no special Map function is needed (use range loop)

#### Scenario: Count matching nodes

- **WHEN** `Count(root Node, pred func(Node) bool)` is called
- **THEN** it SHALL return the number of matching nodes
- **AND** it SHALL NOT allocate a slice (just count)

#### Scenario: Exists check

- **WHEN** `Exists(root Node, pred func(Node) bool)` is called
- **THEN** it SHALL return true if any node matches
- **AND** it SHALL short-circuit on first match

### Requirement: Ancestor Queries

The system SHALL provide queries for ancestor relationships using visitor
context.

#### Scenario: Ancestors function

- **WHEN** `Ancestors(ctx VisitorContext)` is called during visitation
- **THEN** it SHALL return ancestor nodes from parent to root
- **AND** the slice SHALL be ordered [parent, grandparent, ..., root]

#### Scenario: Parent function

- **WHEN** `Parent(ctx VisitorContext)` is called during visitation
- **THEN** it SHALL return the immediate parent node
- **AND** it SHALL return nil for root node

#### Scenario: Depth function

- **WHEN** `Depth(ctx VisitorContext)` is called during visitation
- **THEN** it SHALL return the depth of current node
- **AND** root depth SHALL be 0
- **AND** root's children depth SHALL be 1

### Requirement: Sibling Queries

The system SHALL provide queries for sibling relationships.

#### Scenario: Siblings function

- **WHEN** `Siblings(ctx VisitorContext)` is called during visitation
- **THEN** it SHALL return all sibling nodes (excluding current node)
- **AND** siblings SHALL be in document order

#### Scenario: PreviousSibling function

- **WHEN** `PreviousSibling(ctx VisitorContext)` is called
- **THEN** it SHALL return the preceding sibling node
- **AND** it SHALL return nil if current is first child

#### Scenario: NextSibling function

- **WHEN** `NextSibling(ctx VisitorContext)` is called
- **THEN** it SHALL return the following sibling node
- **AND** it SHALL return nil if current is last child
