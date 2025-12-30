# Delta Specification

## ADDED Requirements

### Requirement: Item Name Path Normalization

Commands accepting item names (validate, archive, accept) SHALL normalize path
arguments to extract the item ID and infer the item type from the path
structure.

#### Scenario: Path with spectr/changes prefix

- **WHEN** user runs a command with argument `spectr/changes/my-change`
- **THEN** the system SHALL extract `my-change` as the item ID
- **AND** SHALL infer the item type as "change"

#### Scenario: Path with spectr/changes prefix and trailing content

- **WHEN** user runs a command with argument
  `spectr/changes/my-change/specs/foo/spec.md`
- **THEN** the system SHALL extract `my-change` as the item ID
- **AND** SHALL infer the item type as "change"

#### Scenario: Path with spectr/specs prefix

- **WHEN** user runs a command with argument `spectr/specs/my-spec`
- **THEN** the system SHALL extract `my-spec` as the item ID
- **AND** SHALL infer the item type as "spec"

#### Scenario: Path with spectr/specs prefix and spec.md file

- **WHEN** user runs a command with argument `spectr/specs/my-spec/spec.md`
- **THEN** the system SHALL extract `my-spec` as the item ID
- **AND** SHALL infer the item type as "spec"

#### Scenario: Simple ID without path prefix

- **WHEN** user runs a command with argument `my-change`
- **THEN** the system SHALL use `my-change` as-is for lookup
- **AND** SHALL use existing auto-detection logic for item type

#### Scenario: Absolute path normalization

- **WHEN** user runs a command with argument
  `/home/user/project/spectr/changes/my-change`
- **THEN** the system SHALL extract `my-change` as the item ID
- **AND** SHALL infer the item type as "change"

#### Scenario: Inferred type precedence

- **WHEN** user provides a path argument that contains `spectr/changes/` or
  `spectr/specs/`
- **THEN** the inferred type from path SHALL be used for validation
- **AND** SHALL NOT trigger "exists as both change and spec" ambiguity errors
