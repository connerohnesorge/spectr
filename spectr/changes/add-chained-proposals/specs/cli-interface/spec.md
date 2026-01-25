## ADDED Requirements

### Requirement: Proposal Frontmatter Support

The system SHALL support YAML frontmatter in `proposal.md` files for declaring
metadata including dependencies.

#### Scenario: Parse valid frontmatter

- **WHEN** proposal.md contains YAML frontmatter block (`---` delimiters)
- **THEN** extract and parse frontmatter into ProposalMetadata
- **AND** make metadata available for validation and commands

#### Scenario: Handle missing frontmatter

- **WHEN** proposal.md has no frontmatter block
- **THEN** treat as valid with empty metadata (backward compatible)

#### Scenario: Handle malformed frontmatter

- **WHEN** proposal.md has invalid YAML in frontmatter
- **THEN** emit validation error with line number and details

### Requirement: Chained Proposal Dependencies

The system SHALL support `requires` and `enables` fields in proposal frontmatter
to declare dependencies between proposals.

#### Scenario: Declare required dependencies

- **WHEN** frontmatter contains `requires` list with id and optional reason
- **THEN** parse as list of Dependency objects
- **AND** validate that referenced change IDs exist or are archived

#### Scenario: Declare enabled proposals

- **WHEN** frontmatter contains `enables` list with id and optional reason
- **THEN** parse as informational metadata (not enforced)

#### Scenario: Self-reference detection

- **WHEN** proposal requires itself
- **THEN** emit validation error "Proposal cannot require itself"

### Requirement: Dependency Validation in Validate Command

The system SHALL check proposal dependencies during `spectr validate`.

#### Scenario: Warn on unmet dependencies

- **WHEN** `spectr validate <id>` runs on proposal with requires
- **AND** a required proposal is not archived
- **THEN** emit warning "Dependency '<id>' is not yet archived"

#### Scenario: Error on circular dependencies

- **WHEN** `spectr validate <id>` detects a dependency cycle
- **THEN** emit error "Circular dependency detected: A → B → ... → A"
- **AND** exit non-zero

#### Scenario: Pass validation when dependencies met

- **WHEN** all required proposals are archived
- **AND** no circular dependencies exist
- **THEN** validation passes (no dependency warnings/errors)

### Requirement: Dependency Enforcement in Accept Command

The system SHALL enforce `requires` dependencies before accepting a proposal.

#### Scenario: Block accept with unmet dependencies

- **WHEN** `spectr accept <id>` runs on proposal with requires
- **AND** any required proposal is not archived
- **THEN** exit non-zero with error listing unmet dependencies
- **AND** do not create tasks.jsonc

#### Scenario: Allow accept when dependencies met

- **WHEN** `spectr accept <id>` runs on proposal with requires
- **AND** all required proposals are archived
- **THEN** proceed with normal accept flow

#### Scenario: Accept proposal without dependencies

- **WHEN** `spectr accept <id>` runs on proposal without requires
- **THEN** proceed with normal accept flow (backward compatible)

### Requirement: Graph Command

The system SHALL provide a `spectr graph` command to visualize proposal
dependencies.

#### Scenario: Graph command registration

- **WHEN** CLI is initialized
- **THEN** register `spectr graph` command

#### Scenario: Display ASCII dependency tree

- **WHEN** `spectr graph [id]`
- **THEN** display ASCII tree showing requires/enables relationships
- **AND** indicate archived status with checkmark (✓) or pending (⧖)

#### Scenario: Display single proposal graph

- **WHEN** `spectr graph <id>`
- **THEN** show dependency tree rooted at specified proposal

#### Scenario: Display full dependency graph

- **WHEN** `spectr graph` (no args)
- **THEN** show all proposals and their relationships

#### Scenario: DOT format output

- **WHEN** `spectr graph --dot`
- **THEN** output Graphviz DOT format for external visualization

#### Scenario: JSON format output

- **WHEN** `spectr graph --json`
- **THEN** output JSON with nodes and edges arrays

### Requirement: Archive Status Detection

The system SHALL detect whether a change proposal has been archived.

#### Scenario: Detect archived change

- **WHEN** checking if change-id is archived
- **THEN** search `spectr/changes/archive/` for matching directory
- **AND** handle date prefix format (YYYY-MM-DD-<id>)

#### Scenario: Detect active change

- **WHEN** change-id exists in `spectr/changes/` (not archive)
- **THEN** report as not archived

#### Scenario: Handle unknown change ID

- **WHEN** change-id does not exist anywhere
- **THEN** report as unknown/not found
