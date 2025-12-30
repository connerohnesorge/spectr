# Delta Specification

## ADDED Requirements

### Requirement: Cross-Capability Requirement Name Independence

The validation system SHALL treat requirement names as scoped to their
capability, allowing the same requirement name to exist in different capability
specs.

#### Scenario: Same name in different capability deltas

- **WHEN** a change modifies requirements with the same name in different
  capability specs
- **THEN** validation SHALL pass without duplicate name errors
- **AND** each requirement SHALL be matched to its own capability's base spec
- **AND** example: `support-aider::No Instruction File` is distinct from
  `support-cursor::No Instruction File`

#### Scenario: Same name REMOVED in multiple capabilities

- **WHEN** a change removes a requirement named "X" from capability A
- **AND** removes a requirement named "X" from capability B
- **THEN** validation SHALL NOT report "Requirement 'X' is REMOVED in multiple
  files"
- **AND** each removal SHALL be validated against its respective capability spec

#### Scenario: Same name MODIFIED in multiple capabilities

- **WHEN** a change modifies a requirement named "X" in capability A
- **AND** modifies a requirement named "X" in capability B
- **THEN** validation SHALL NOT report "Requirement 'X' is MODIFIED in multiple
  files"
- **AND** each modification SHALL be validated against its respective capability
  spec

#### Scenario: Same name ADDED in multiple capabilities

- **WHEN** a change adds a requirement named "X" to capability A (new spec)
- **AND** adds a requirement named "X" to capability B (new spec)
- **THEN** validation SHALL NOT report duplicate requirement errors
- **AND** each addition SHALL be validated independently

## MODIFIED Requirements

### Requirement: Change Delta Validation

The validation system SHALL validate change delta specs for structural
correctness and delta operation validity.

#### Scenario: Valid change with deltas

- **WHEN** a change directory contains specs with proper
  ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** validation SHALL pass with no errors
- **AND** each delta requirement SHALL be counted toward the total

#### Scenario: Change with no deltas

- **WHEN** a change directory has no specs/ subdirectory or no delta sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL explain that at least one delta is required
- **AND** remediation guidance SHALL explain the delta header format

#### Scenario: Delta sections present but empty

- **WHEN** delta sections exist (## ADDED Requirements) but contain no
  requirement entries
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which sections are empty
- **AND** guidance SHALL explain requirement block format

#### Scenario: ADDED requirement without scenario

- **WHEN** an ADDED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which requirement is missing scenarios

#### Scenario: MODIFIED requirement without scenario

- **WHEN** a MODIFIED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL require at least one scenario for MODIFIED
  requirements

#### Scenario: Duplicate requirement in same section within file

- **WHEN** two requirements with the same normalized name appear in the same
  delta section of the same file
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL identify the duplicate requirement name

#### Scenario: Cross-section conflicts within file

- **WHEN** a requirement appears in both ADDED and MODIFIED sections of the same
  file
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate the conflicting requirement and sections

#### Scenario: RENAMED requirement validation

- **WHEN** a RENAMED section contains well-formed "FROM: X TO: Y" pairs
- **THEN** validation SHALL accept the renames using `markdown.MatchRenamedFrom`
  and `markdown.MatchRenamedTo` functions
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Delta parsing uses markdown parser package

- **WHEN** the validation system parses delta spec files
- **THEN** it SHALL use functions from `internal/markdown/` package
- **AND** delta type detection SHALL use `markdown.MatchDeltaSection`
- **AND** requirement extraction SHALL use `markdown.MatchRequirementHeader`
- **AND** scenario extraction SHALL use `markdown.MatchScenarioHeader`
- **AND** section content extraction SHALL use `markdown.FindSection` or
  `markdown.FindDeltaSection`

#### Scenario: Cross-capability same-named requirements allowed

- **WHEN** multiple delta specs in different capability directories contain
  requirements with the same name
- **THEN** validation SHALL NOT report duplicate errors
- **AND** uniqueness SHALL be scoped to (capability, requirement_name) tuples
