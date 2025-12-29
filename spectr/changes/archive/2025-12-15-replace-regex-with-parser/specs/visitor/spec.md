# Delta Specification

## ADDED Requirements

### Requirement: Visitor Interface

The system SHALL define a Visitor interface that enables type-safe traversal of
AST nodes with double dispatch.

#### Scenario: Visitor interface definition

- **WHEN** the Visitor interface is defined
- **THEN** it SHALL have a method for each node type: `VisitDocument`,
  `VisitSection`, `VisitRequirement`, `VisitScenario`, etc.
- **AND** each method SHALL receive the specific node type as parameter
- **AND** each method SHALL return `error` to allow early termination

#### Scenario: Node Accept method

- **WHEN** `node.Accept(visitor Visitor)` is called
- **THEN** the node SHALL call the appropriate Visit method on the visitor
- **AND** it SHALL pass itself as the typed argument
- **AND** it SHALL return any error from the Visit method

#### Scenario: Default visitor implementation

- **WHEN** implementing a visitor
- **THEN** a `BaseVisitor` struct SHALL be provided with no-op defaults
- **AND** implementers SHALL embed `BaseVisitor` and override only needed
  methods
- **AND** `BaseVisitor` methods SHALL return nil (continue traversal)

### Requirement: Visitor Traversal Control

The system SHALL allow visitors to control traversal order and skip subtrees.

#### Scenario: Pre-order traversal

- **WHEN** default traversal is used
- **THEN** nodes SHALL be visited in pre-order (parent before children)
- **AND** children SHALL be visited left-to-right
- **AND** the visitor method is called before recursing into children

#### Scenario: Skip children via error

- **WHEN** a Visit method returns `SkipChildren` sentinel error
- **THEN** the traversal SHALL NOT visit that node's children
- **AND** traversal SHALL continue with the next sibling
- **AND** `SkipChildren` SHALL NOT be treated as a real error

#### Scenario: Stop traversal via error

- **WHEN** a Visit method returns a non-nil, non-SkipChildren error
- **THEN** traversal SHALL stop immediately
- **AND** the error SHALL be returned from the Walk function
- **AND** no further nodes SHALL be visited

### Requirement: Walk Function

The system SHALL provide a Walk function that orchestrates visitor traversal.

#### Scenario: Walk function signature

- **WHEN** calling the Walk function
- **THEN** the signature SHALL be `Walk(node *Node, visitor Visitor) error`
- **AND** it SHALL handle the recursion through children
- **AND** it SHALL respect SkipChildren and error returns

#### Scenario: Walk visits all descendants

- **WHEN** `Walk(root, visitor)` is called
- **THEN** every node in the tree SHALL have its Accept method called
- **AND** the order SHALL be consistent (pre-order depth-first)
- **AND** the visitor's Visit methods SHALL be called exactly once per node

#### Scenario: Walk with nil node

- **WHEN** `Walk(nil, visitor)` is called
- **THEN** it SHALL return nil without calling any Visit methods
- **AND** this SHALL allow safe handling of optional children

### Requirement: Typed Visitor Methods

The system SHALL define Visit methods for each concrete node type.

#### Scenario: Document visitor method

- **WHEN** visiting a `NodeDocument`
- **THEN** `VisitDocument(*NodeDocument) error` SHALL be called
- **AND** the method receives the typed node, not interface{}

#### Scenario: Section visitor method

- **WHEN** visiting a `NodeSection`
- **THEN** `VisitSection(*NodeSection) error` SHALL be called
- **AND** Level and Title fields SHALL be accessible

#### Scenario: Requirement visitor method

- **WHEN** visiting a `NodeRequirement`
- **THEN** `VisitRequirement(*NodeRequirement) error` SHALL be called
- **AND** Name field SHALL be accessible

#### Scenario: Scenario visitor method

- **WHEN** visiting a `NodeScenario`
- **THEN** `VisitScenario(*NodeScenario) error` SHALL be called
- **AND** Name field SHALL be accessible

#### Scenario: Inline element visitor methods

- **WHEN** visiting inline elements
- **THEN** there SHALL be: `VisitText`, `VisitStrong`, `VisitEmphasis`,
  `VisitStrikethrough`, `VisitCode`, `VisitLink`, `VisitWikilink`
- **AND** each SHALL receive its typed node

#### Scenario: Block element visitor methods

- **WHEN** visiting block elements
- **THEN** there SHALL be: `VisitParagraph`, `VisitList`, `VisitListItem`,
  `VisitCodeBlock`, `VisitBlockquote`
- **AND** each SHALL receive its typed node

### Requirement: Visitor Enter/Leave Pattern

The system SHALL support an optional Enter/Leave pattern for visitors needing
both pre and post visitation.

#### Scenario: EnterLeave visitor interface

- **WHEN** an `EnterLeaveVisitor` interface is defined
- **THEN** it SHALL have `EnterX` and `LeaveX` methods for each node type
- **AND** `EnterX` SHALL be called before visiting children
- **AND** `LeaveX` SHALL be called after visiting children

#### Scenario: WalkEnterLeave function

- **WHEN** `WalkEnterLeave(node, visitor EnterLeaveVisitor)` is called
- **THEN** for each node: EnterX, recurse children, LeaveX
- **AND** LeaveX SHALL be called even if children returned SkipChildren
- **AND** LeaveX SHALL NOT be called if EnterX returns a real error

#### Scenario: BaseEnterLeaveVisitor

- **WHEN** implementing an EnterLeaveVisitor
- **THEN** `BaseEnterLeaveVisitor` SHALL provide no-op defaults
- **AND** both Enter and Leave methods SHALL default to nil return

### Requirement: Visitor Utilities

The system SHALL provide utility functions for common visitor patterns.

#### Scenario: Collect nodes by type

- **WHEN** `CollectByType[T Node](root *Node) []*T` is called
- **THEN** it SHALL return all nodes of type T in the tree
- **AND** the result SHALL be in traversal order
- **AND** this SHALL use the visitor internally

#### Scenario: Find first node by predicate

- **WHEN** `FindFirst(root *Node, pred func(*Node) bool) *Node` is called
- **THEN** it SHALL return the first node where pred returns true
- **AND** it SHALL stop traversal after finding the match
- **AND** it SHALL return nil if no match found

#### Scenario: Count nodes

- **WHEN** `Count(root *Node) int` is called
- **THEN** it SHALL return the total number of nodes in the tree
- **AND** this SHALL include the root node

#### Scenario: Depth calculation

- **WHEN** `Depth(root *Node) int` is called
- **THEN** it SHALL return the maximum depth of the tree
- **AND** root depth SHALL be 1
- **AND** an empty document SHALL have depth 1

### Requirement: Visitor Thread Safety

The system SHALL document thread safety guarantees for visitors.

#### Scenario: Visitor instance safety

- **WHEN** a visitor is used
- **THEN** a single visitor instance SHALL NOT be used concurrently
- **AND** each Walk call SHALL use its own visitor instance OR external
  synchronization

#### Scenario: AST immutability enables concurrent reading

- **WHEN** multiple visitors traverse the same AST
- **THEN** concurrent read-only traversal SHALL be safe
- **AND** no locking SHALL be required for concurrent Walk calls on same tree
- **AND** this is enabled by AST immutability

### Requirement: Visitor Error Handling

The system SHALL define clear error handling semantics for visitors.

#### Scenario: SkipChildren sentinel

- **WHEN** `SkipChildren` error is defined
- **THEN** it SHALL be a package-level sentinel error: `var SkipChildren =
  errors.New("skip children")`
- **AND** checking SHALL be via `errors.Is(err, SkipChildren)`

#### Scenario: Error wrapping

- **WHEN** a visitor returns a wrapped error
- **THEN** Walk SHALL return the wrapped error unchanged
- **AND** callers MAY use `errors.Unwrap` to access underlying error
- **AND** source position SHOULD be included in error wrapping

#### Scenario: Panic recovery

- **WHEN** a visitor method panics
- **THEN** Walk SHALL NOT recover the panic
- **AND** the panic SHALL propagate to the caller
- **AND** this allows debugging panic locations
