# Validation Specification - Dependency Syntax Delta

## MODIFIED Requirements

### Requirement: Change Delta Validation
The validation system SHALL validate change delta specs for structural correctness, delta operation validity, and dependency reference validity.

#### Scenario: Valid change with deltas
- **WHEN** a change directory contains specs with proper ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** validation SHALL pass with no errors
- **AND** each delta requirement SHALL be counted toward the total

#### Scenario: Change with no deltas
- **WHEN** a change directory has no specs/ subdirectory or no delta sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL explain that at least one delta is required
- **AND** remediation guidance SHALL explain the delta header format

#### Scenario: Delta sections present but empty
- **WHEN** delta sections exist (## ADDED Requirements) but contain no requirement entries
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
- **AND** the message SHALL require at least one scenario for MODIFIED requirements

#### Scenario: Duplicate requirement in same section
- **WHEN** two requirements with the same normalized name appear in the same delta section
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL identify the duplicate requirement name

#### Scenario: Cross-section conflicts
- **WHEN** a requirement appears in both ADDED and MODIFIED sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate the conflicting requirement and sections

#### Scenario: RENAMED requirement validation
- **WHEN** a RENAMED section contains well-formed "FROM: X TO: Y" pairs
- **THEN** validation SHALL accept the renames
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Valid dependency declarations
- **WHEN** a change proposal.md contains @depends(change-id) or @requires(spec:capability-id) syntax
- **THEN** validation SHALL parse dependency references
- **AND** SHALL verify referenced changes exist in spectr/changes/ or spectr/changes/archive/
- **AND** SHALL verify referenced specs exist in spectr/specs/
- **AND** SHALL pass if all dependencies are valid

#### Scenario: Missing change dependency in development mode
- **WHEN** a change proposal.md contains @depends(change-id) and the referenced change does not exist
- **THEN** validation SHALL generate a WARNING level issue in normal mode
- **AND** the warning message SHALL indicate the missing change ID
- **AND** SHALL suggest checking the change ID spelling or creating the dependency first

#### Scenario: Missing change dependency in strict mode
- **WHEN** validation runs in strict mode (used during archive) and a change proposal.md contains @depends(change-id) for a non-existent change
- **THEN** validation SHALL generate an ERROR level issue
- **AND** valid SHALL be false
- **AND** the error message SHALL indicate the missing dependency must be resolved before archiving

#### Scenario: Missing spec requirement in development mode
- **WHEN** a change proposal.md contains @requires(spec:capability-id) and the referenced spec does not exist
- **THEN** validation SHALL generate a WARNING level issue in normal mode
- **AND** the warning message SHALL indicate the missing spec ID
- **AND** SHALL suggest checking the spec ID spelling or implementing the capability first

#### Scenario: Missing spec requirement in strict mode
- **WHEN** validation runs in strict mode and a change proposal.md contains @requires(spec:capability-id) for a non-existent spec
- **THEN** validation SHALL generate an ERROR level issue
- **AND** valid SHALL be false
- **AND** the error message SHALL indicate the missing requirement must be resolved before archiving

#### Scenario: Self-referencing dependency
- **WHEN** a change proposal.md contains @depends(same-change-id) referencing itself
- **THEN** validation SHALL generate an ERROR level issue
- **AND** the message SHALL indicate that self-references are not allowed
- **AND** SHALL show the line number where the self-reference appears

#### Scenario: Duplicate dependency declarations
- **WHEN** a change proposal.md contains the same @depends() or @requires() declaration multiple times
- **THEN** validation SHALL generate a WARNING level issue
- **AND** the message SHALL indicate the duplicate dependency
- **AND** SHALL suggest removing redundant declarations

#### Scenario: Malformed dependency syntax
- **WHEN** a change proposal.md contains malformed @depends() or @requires() syntax (missing parentheses, invalid characters)
- **THEN** validation SHALL generate a WARNING level issue
- **AND** the message SHALL show the correct syntax format
- **AND** SHALL provide examples of valid dependency declarations

#### Scenario: Dependency validation with line numbers
- **WHEN** validation detects dependency issues
- **THEN** error messages SHALL include the file path and line number where the dependency appears
- **AND** SHALL format as "proposal.md:12: Missing dependency @depends(foo)"
- **AND** SHALL help developers locate and fix issues quickly

## ADDED Requirements

### Requirement: Dependency Reference Parsing
The validation system SHALL parse dependency references from change proposal files using inline syntax.

#### Scenario: Parse @depends syntax
- **WHEN** a proposal.md file contains `@depends(change-id)` inline in text
- **THEN** the parser SHALL extract the change-id
- **AND** SHALL record the dependency type as CHANGE
- **AND** SHALL capture the line number where the reference appears

#### Scenario: Parse @requires syntax
- **WHEN** a proposal.md file contains `@requires(spec:capability-id)` inline in text
- **THEN** the parser SHALL extract the capability-id
- **AND** SHALL record the dependency type as SPEC
- **AND** SHALL capture the line number where the reference appears

#### Scenario: Parse multiple dependencies
- **WHEN** a proposal.md file contains multiple @depends() or @requires() declarations
- **THEN** the parser SHALL extract all dependency references
- **AND** SHALL preserve order of appearance
- **AND** SHALL track line numbers for each reference

#### Scenario: Handle whitespace variations
- **WHEN** dependency syntax contains extra whitespace (e.g., `@depends( change-id )`)
- **THEN** the parser SHALL normalize whitespace and extract the correct ID
- **AND** SHALL not generate errors for whitespace variations

#### Scenario: Ignore dependencies in code blocks
- **WHEN** a proposal.md file contains @depends() or @requires() within fenced code blocks (```)
- **THEN** the parser SHOULD ignore these references as examples, not actual dependencies
- **AND** SHALL only parse dependencies in prose sections

### Requirement: Dependency Validation Rules
The validation system SHALL validate that all declared dependencies reference existing changes or specs.

#### Scenario: Validate change dependency exists
- **WHEN** validating a @depends(change-id) reference
- **THEN** the validator SHALL check if the change exists in spectr/changes/
- **AND** SHALL also check spectr/changes/archive/ for archived changes
- **AND** SHALL pass if found in either location

#### Scenario: Validate spec requirement exists
- **WHEN** validating a @requires(spec:capability-id) reference
- **THEN** the validator SHALL check if the spec exists in spectr/specs/capability-id/spec.md
- **AND** SHALL pass if the spec file exists and is readable

#### Scenario: Check for self-references
- **WHEN** validating dependencies for a change
- **THEN** the validator SHALL compare each @depends() reference against the current change ID
- **AND** SHALL generate an ERROR if a change depends on itself

#### Scenario: Check for duplicate declarations
- **WHEN** validating dependencies for a change
- **THEN** the validator SHALL track all unique dependency references
- **AND** SHALL generate a WARNING if the same dependency is declared multiple times
- **AND** SHALL report the line numbers of duplicate declarations

#### Scenario: Respect strict mode flag
- **WHEN** dependency validation runs with strictMode=true
- **THEN** missing dependencies SHALL generate ERROR level issues instead of WARNING
- **AND** valid SHALL be set to false if any dependencies are missing
- **AND** exit code SHALL be non-zero for missing dependencies

#### Scenario: Provide actionable error messages
- **WHEN** dependency validation fails
- **THEN** error messages SHALL include the file path and line number
- **AND** SHALL explain what is wrong (e.g., "Change 'foo' not found")
- **AND** SHALL provide remediation guidance (e.g., "Check spelling or create the dependent change first")
- **AND** SHALL show examples of correct syntax if syntax is malformed
