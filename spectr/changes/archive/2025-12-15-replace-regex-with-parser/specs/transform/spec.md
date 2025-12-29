## ADDED Requirements

### Requirement: Transform Action Type

The system SHALL define an action type to signal transform intentions explicitly.

#### Scenario: TransformAction enumeration

- **WHEN** transform actions are defined
- **THEN** there SHALL be: `ActionKeep`, `ActionReplace`, `ActionDelete`
- **AND** `ActionKeep` means return original node unchanged
- **AND** `ActionReplace` means use returned node as replacement
- **AND** `ActionDelete` means remove node from parent's children

#### Scenario: Action type definition

- **WHEN** the action type is implemented
- **THEN** it SHALL be `type TransformAction uint8`
- **AND** constants SHALL be defined for each action

### Requirement: Transform Visitor Interface

The system SHALL define a TransformVisitor interface for AST rewriting via visitor pattern.

#### Scenario: Transform visitor signature

- **WHEN** the TransformVisitor interface is defined
- **THEN** it SHALL have a method for each node type
- **AND** each method signature SHALL be `TransformX(*NodeX) (Node, TransformAction, error)`
- **AND** the returned Node SHALL be the replacement (only used when action is Replace)

#### Scenario: Transform visitor methods

- **WHEN** implementing TransformVisitor
- **THEN** there SHALL be: `TransformDocument`, `TransformSection`, `TransformRequirement`, `TransformScenario`, `TransformParagraph`, `TransformList`, `TransformListItem`, etc.
- **AND** each receives typed node and returns (Node, action, error)

#### Scenario: Base transform visitor

- **WHEN** implementing a transform
- **THEN** `BaseTransformVisitor` SHALL provide defaults that return (original, ActionKeep, nil)
- **AND** implementers SHALL embed BaseTransformVisitor and override only needed methods

### Requirement: Transform Function

The system SHALL provide a Transform function that applies a TransformVisitor to an AST.

#### Scenario: Transform function signature

- **WHEN** transforming an AST
- **THEN** the signature SHALL be `Transform(root Node, visitor TransformVisitor) (Node, error)`
- **AND** it SHALL return the transformed root node
- **AND** original AST SHALL remain unchanged (immutability)

#### Scenario: Transform traversal order

- **WHEN** Transform traverses the AST
- **THEN** it SHALL use post-order traversal (children before parent)
- **AND** this allows transforms to see already-transformed children
- **AND** parent transform sees the results of child transforms

#### Scenario: Transform with deletion

- **WHEN** a transform method returns ActionDelete
- **THEN** the node SHALL be removed from its parent's children
- **AND** sibling nodes SHALL be preserved
- **AND** deleting root node SHALL return nil root

#### Scenario: Transform with replacement

- **WHEN** a transform method returns (newNode, ActionReplace, nil)
- **THEN** the original node SHALL be replaced by newNode in parent's children
- **AND** newNode's children SHALL already be transformed (post-order)

### Requirement: Transform Composition

The system SHALL support composing multiple transforms.

#### Scenario: Sequential transform composition

- **WHEN** `Compose(t1, t2 TransformVisitor)` is called
- **THEN** it SHALL return a TransformVisitor that applies t1 then t2
- **AND** t2 receives the output of t1 for each node

#### Scenario: Transform pipeline

- **WHEN** `Pipeline(transforms ...TransformVisitor)` is called
- **THEN** it SHALL return a TransformVisitor applying all transforms in order
- **AND** each transform sees the result of all previous transforms

#### Scenario: Conditional transform

- **WHEN** `When(pred func(Node) bool, t TransformVisitor)` is called
- **THEN** it SHALL return a TransformVisitor that applies t only when pred returns true
- **AND** nodes not matching pred SHALL pass through unchanged

### Requirement: Common Transform Utilities

The system SHALL provide utility transforms for common operations.

#### Scenario: Map transform

- **WHEN** `Map(f func(Node) Node)` is called
- **THEN** it SHALL return a TransformVisitor that applies f to every node
- **AND** f returning the same node SHALL be treated as ActionKeep

#### Scenario: Filter transform

- **WHEN** `Filter(pred func(Node) bool)` is called
- **THEN** it SHALL return a TransformVisitor that deletes nodes where pred returns false
- **AND** nodes matching pred SHALL be kept

#### Scenario: Replace by type transform

- **WHEN** `ReplaceType[T Node](f func(*T) Node)` is called
- **THEN** it SHALL return a TransformVisitor that applies f to nodes of type T
- **AND** other node types SHALL pass through unchanged

### Requirement: Transform Error Handling

The system SHALL define clear error handling for transforms.

#### Scenario: Transform error propagation

- **WHEN** a transform method returns a non-nil error
- **THEN** Transform SHALL stop immediately
- **AND** the error SHALL be returned from Transform
- **AND** partial results SHALL NOT be returned

#### Scenario: Skip subtree in transform

- **WHEN** a transform wants to skip a node's children
- **THEN** it SHALL return (original, ActionKeep, nil) in pre-transform hook
- **AND** children SHALL NOT be visited for transformation

#### Scenario: Transform with context

- **WHEN** transforms need to share state
- **THEN** TransformVisitor implementations MAY have fields for context
- **AND** a single TransformVisitor instance SHALL NOT be used concurrently

### Requirement: Transform Result Validation

The system SHALL validate transform results for consistency.

#### Scenario: Replaced node type compatibility

- **WHEN** a transform replaces a node
- **THEN** the replacement SHALL be type-compatible with the original's position
- **AND** replacing NodeSection with NodeText in section position SHALL be an error

#### Scenario: Children consistency after transform

- **WHEN** Transform completes
- **THEN** all parent-child relationships SHALL be consistent
- **AND** child Start/End offsets MAY be invalid (transforms may create synthetic nodes)

#### Scenario: Hash recomputation

- **WHEN** a node is replaced or modified via transform
- **THEN** the new node's hash SHALL be computed
- **AND** all ancestor hashes SHALL be recomputed (since children changed)

### Requirement: Specialized Transform Helpers

The system SHALL provide helpers for common Spectr transformations.

#### Scenario: Rename requirement transform

- **WHEN** `RenameRequirement(oldName, newName string)` is called
- **THEN** it SHALL return a TransformVisitor that renames matching requirements
- **AND** only requirements with Name() == oldName SHALL be affected

#### Scenario: Add scenario transform

- **WHEN** `AddScenario(reqName string, scenario *NodeScenario)` is called
- **THEN** it SHALL return a TransformVisitor that adds scenario to matching requirement
- **AND** scenario SHALL be appended to requirement's children

#### Scenario: Remove requirement transform

- **WHEN** `RemoveRequirement(name string)` is called
- **THEN** it SHALL return a TransformVisitor that deletes matching requirements
- **AND** action SHALL be ActionDelete for matching requirements
