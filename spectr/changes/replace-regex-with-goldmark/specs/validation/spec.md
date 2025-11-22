## MODIFIED Requirements

### Requirement: Spec File Validation
The validation system SHALL validate spec files for structural correctness and adherence to Spectr conventions using goldmark AST-based parsing for robust markdown handling and precise error locations.

#### Scenario: Valid spec with all required sections
- **WHEN** a spec file contains Purpose and Requirements sections with properly formatted requirements and scenarios
- **THEN** validation SHALL pass with no errors
- **AND** the validation report SHALL indicate valid=true

#### Scenario: Missing Purpose section
- **WHEN** a spec file lacks a "## Purpose" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL indicate which section is missing with precise line:column location from AST
- **AND** the error message SHALL include remediation guidance showing correct format

#### Scenario: Missing Requirements section
- **WHEN** a spec file lacks a "## Requirements" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL provide example of correct structure with precise line location

#### Scenario: Requirement without scenarios
- **WHEN** a requirement exists without any "#### Scenario:" subsections
- **THEN** validation SHALL report a WARNING level issue with line:column of the requirement header
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include example scenario format

#### Scenario: Requirement missing SHALL or MUST
- **WHEN** a requirement text does not contain "SHALL" or "MUST" keywords
- **THEN** validation SHALL report a WARNING level issue with precise requirement location
- **AND** the message SHALL suggest using normative language

#### Scenario: Incorrect scenario format
- **WHEN** scenarios use formats other than "#### Scenario:" (e.g., bullets or bold text)
- **THEN** validation SHALL report an ERROR with exact line:column from AST
- **AND** the message SHALL show the correct "#### Scenario:" header format

#### Scenario: Code blocks containing fake headers
- **WHEN** a spec file contains code blocks with lines starting with `###` or `####`
- **THEN** validation SHALL NOT treat these as requirements or scenarios
- **AND** SHALL correctly parse using AST-based markdown parsing that distinguishes code from structure

#### Scenario: Nested lists in requirements
- **WHEN** a requirement contains nested bullet lists or numbered lists
- **THEN** validation SHALL correctly parse the requirement body using AST
- **AND** SHALL include nested content in the requirement block
- **AND** SHALL NOT falsely detect list items as new sections

### Requirement: Change Delta Validation
The validation system SHALL validate change delta specs for structural correctness and delta operation validity using goldmark AST parsing for accurate section extraction.

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
- **AND** the message SHALL indicate which sections are empty with line location from AST
- **AND** guidance SHALL explain requirement block format

#### Scenario: ADDED requirement without scenario
- **WHEN** an ADDED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR showing precise location
- **AND** the message SHALL indicate which requirement is missing scenarios

#### Scenario: MODIFIED requirement without scenario
- **WHEN** a MODIFIED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR showing requirement line:column
- **AND** the message SHALL require at least one scenario for MODIFIED requirements

#### Scenario: Duplicate requirement in same section
- **WHEN** two requirements with the same normalized name appear in the same delta section
- **THEN** validation SHALL fail with an ERROR showing both line locations from AST
- **AND** the message SHALL identify the duplicate requirement name with source positions

#### Scenario: Cross-section conflicts
- **WHEN** a requirement appears in both ADDED and MODIFIED sections
- **THEN** validation SHALL fail with an ERROR showing line locations of both occurrences
- **AND** the message SHALL indicate the conflicting requirement and sections

#### Scenario: RENAMED requirement validation
- **WHEN** a RENAMED section contains well-formed "FROM: X TO: Y" pairs
- **THEN** validation SHALL accept the renames
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Malformed markdown in delta sections
- **WHEN** delta sections contain malformed markdown (unclosed code blocks, broken links)
- **THEN** validation SHALL parse using goldmark AST which handles CommonMark edge cases
- **AND** SHALL extract requirements correctly despite formatting issues
- **AND** MAY report warnings for markdown quality issues

### Requirement: Helpful Error Messages
The validation system SHALL provide actionable error messages with remediation guidance and precise source locations derived from goldmark AST node positions.

#### Scenario: Error with remediation steps and source location
- **WHEN** validation fails due to missing required content
- **THEN** the error message SHALL explain what is wrong with precise line:column reference
- **AND** SHALL provide "Next steps" section with concrete actions
- **AND** SHALL include format examples when applicable
- **AND** location format SHALL be `path/to/file.md:42:5:` for easy editor navigation

#### Scenario: Ambiguous item name
- **WHEN** user provides an item name that matches both a change and a spec
- **THEN** validation SHALL report the ambiguity
- **AND** SHALL suggest using --type flag to disambiguate
- **AND** SHALL show available type options (change, spec)

#### Scenario: Item not found with suggestions
- **WHEN** user provides an item name that does not exist
- **THEN** validation SHALL report item not found
- **AND** SHALL provide nearest match suggestions based on string similarity
- **AND** SHALL limit suggestions to 5 most similar items
