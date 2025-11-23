# Validation Delta Spec

## MODIFIED Requirements

### Requirement: Spec File Validation

The validation system SHALL validate spec files for structural correctness using a lexer/parser architecture that correctly handles markdown edge cases.

#### Scenario: Valid spec with all required sections

- **WHEN** a spec file contains Purpose and Requirements sections with properly formatted requirements and scenarios
- **THEN** validation SHALL pass with no errors using AST-based parsing
- **AND** the validation report SHALL indicate valid=true
- **AND** SHALL correctly parse requirements even with code blocks present

#### Scenario: Missing Purpose section

- **WHEN** a spec file lacks a "## Purpose" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL indicate which section is missing with accurate line/column position
- **AND** the error message SHALL include remediation guidance showing correct format

#### Scenario: Missing Requirements section

- **WHEN** a spec file lacks a "## Requirements" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error SHALL include position information from the AST
- **AND** the error message SHALL provide example of correct structure

#### Scenario: Requirement without scenarios

- **WHEN** a requirement exists without any "#### Scenario:" subsections
- **THEN** validation SHALL report a WARNING level issue
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include line and column position of the requirement
- **AND** the warning SHALL include example scenario format

#### Scenario: Requirement missing SHALL or MUST

- **WHEN** a requirement text does not contain "SHALL" or "MUST" keywords
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL suggest using normative language
- **AND** SHALL include position information for the requirement

#### Scenario: Incorrect scenario format

- **WHEN** scenarios use formats other than "#### Scenario:" (e.g., bullets or bold text)
- **THEN** validation SHALL report an ERROR
- **AND** the message SHALL show the correct "#### Scenario:" header format
- **AND** SHALL include line and column of the incorrect format

#### Scenario: Code blocks do not interfere with validation

- **WHEN** a spec contains code blocks with requirement or scenario syntax inside
- **THEN** validation SHALL NOT treat code block content as actual requirements or scenarios
- **AND** SHALL only validate actual markdown structure outside code blocks
- **AND** SHALL correctly count requirements excluding those in code blocks

### Requirement: Change Delta Validation

The validation system SHALL validate change delta specs using AST-based parsing that correctly handles markdown edge cases.

#### Scenario: Valid change with deltas

- **WHEN** a change directory contains specs with proper ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** validation SHALL pass with no errors using the new parser
- **AND** each delta requirement SHALL be counted toward the total
- **AND** SHALL correctly parse deltas even with code blocks in requirements

#### Scenario: Change with no deltas

- **WHEN** a change directory has no specs/ subdirectory or no delta sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL explain that at least one delta is required
- **AND** remediation guidance SHALL explain the delta header format

#### Scenario: Delta sections present but empty

- **WHEN** delta sections exist (## ADDED Requirements) but contain no requirement entries
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which sections are empty with position information
- **AND** guidance SHALL explain requirement block format

#### Scenario: ADDED requirement without scenario

- **WHEN** an ADDED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which requirement is missing scenarios
- **AND** SHALL include line and column position

#### Scenario: MODIFIED requirement without scenario

- **WHEN** a MODIFIED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL require at least one scenario for MODIFIED requirements
- **AND** SHALL include position information

#### Scenario: Duplicate requirement in same section

- **WHEN** two requirements with the same normalized name appear in the same delta section
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL identify the duplicate requirement name
- **AND** SHALL include positions of both occurrences

#### Scenario: Cross-section conflicts

- **WHEN** a requirement appears in both ADDED and MODIFIED sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate the conflicting requirement and sections
- **AND** SHALL include position information for both occurrences

#### Scenario: RENAMED requirement validation

- **WHEN** a RENAMED section contains well-formed "FROM: X TO: Y" pairs
- **THEN** validation SHALL accept the renames using AST-based parsing
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Code blocks in delta specs

- **WHEN** a delta spec contains code blocks with requirement examples
- **THEN** validation SHALL NOT extract requirement headers from inside code blocks
- **AND** SHALL only validate actual delta operations outside code blocks
- **AND** SHALL correctly report the count of actual delta requirements

### Requirement: Helpful Error Messages

The validation system SHALL provide actionable error messages with position information using data from the markdown parser.

#### Scenario: Error with remediation steps and position

- **WHEN** validation fails due to missing required content
- **THEN** the error message SHALL explain what is wrong
- **AND** SHALL include line and column number from the AST
- **AND** SHALL provide "Next steps" section with concrete actions
- **AND** SHALL include format examples when applicable
- **AND** MAY include context showing surrounding markdown

#### Scenario: Parse error with position

- **WHEN** validation encounters malformed markdown (e.g., unclosed code fence)
- **THEN** the error message SHALL include line and column of the error
- **AND** SHALL show context (surrounding lines) for debugging
- **AND** SHALL provide actionable suggestion for fixing the error

#### Scenario: Item not found with suggestions

- **WHEN** user provides an item name that does not exist
- **THEN** validation SHALL report item not found
- **AND** SHALL provide nearest match suggestions based on string similarity
- **AND** SHALL limit suggestions to 5 most similar items
