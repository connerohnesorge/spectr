## ADDED Requirements

### Requirement: AST-Based Markdown Parsing
The validation system SHALL use blackfriday AST-based parsing instead of regex patterns for extracting markdown structure from spec and delta files.

#### Scenario: Parse spec file with AST
- **WHEN** the validation system parses a spec.md file
- **THEN** it SHALL use blackfriday to build an AST representation
- **AND** it SHALL extract headers by walking AST nodes of type Heading
- **AND** it SHALL extract content by collecting text between header boundaries
- **AND** results SHALL be equivalent to previous regex-based extraction

#### Scenario: Parse delta spec with AST
- **WHEN** the validation system parses a delta spec file
- **THEN** it SHALL identify ADDED/MODIFIED/REMOVED/RENAMED sections via H2 heading text
- **AND** it SHALL extract requirement blocks via H3 headings starting with "Requirement:"
- **AND** it SHALL extract scenarios via H4 headings starting with "Scenario:"

#### Scenario: Parse task checkboxes with AST
- **WHEN** the validation system parses a tasks.md file
- **THEN** it SHALL identify task items via List nodes containing checkbox patterns
- **AND** it SHALL determine completion status from checkbox character (space vs x/X)
- **AND** results SHALL match previous `^\s*-\s*\[([xX ])\]` regex behavior

### Requirement: Markdown Package Architecture
The system SHALL provide a dedicated `internal/markdown/` package that encapsulates all blackfriday AST operations, exposing a clean API for spec parsing needs.

#### Scenario: Package exposes typed extraction functions
- **WHEN** code imports the markdown package
- **THEN** it SHALL have access to `ParseDocument(content []byte) *Document`
- **AND** it SHALL have access to `ExtractHeaders(doc *Document, level int) []Header`
- **AND** it SHALL have access to `ExtractSections(doc *Document) map[string]Section`
- **AND** it SHALL have access to `ExtractTasks(doc *Document) []Task`

#### Scenario: Types provide source location
- **WHEN** a Header, Section, or Task is extracted
- **THEN** the type SHALL include line number information from the AST
- **AND** error messages SHALL be able to reference specific line numbers

#### Scenario: Package hides blackfriday internals
- **WHEN** consuming code uses the markdown package
- **THEN** it SHALL NOT need to import blackfriday directly
- **AND** all blackfriday types SHALL be wrapped in package-specific types
- **AND** changes to blackfriday version SHALL only affect internal/markdown/

## MODIFIED Requirements

### Requirement: Spec File Validation
The validation system SHALL validate spec files for structural correctness and adherence to Spectr conventions.

#### Scenario: Valid spec with all required sections
- **WHEN** a spec file contains a Requirements section with properly formatted requirements and scenarios
- **THEN** validation SHALL pass with no errors
- **AND** the validation report SHALL indicate valid=true

#### Scenario: Missing Requirements section
- **WHEN** a spec file lacks a "## Requirements" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL provide example of correct structure

#### Scenario: Requirement without scenarios
- **WHEN** a requirement exists without any "#### Scenario:" subsections
- **THEN** validation SHALL report a WARNING level issue
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include example scenario format

#### Scenario: Requirement missing SHALL or MUST
- **WHEN** a requirement text does not contain "SHALL" or "MUST" keywords
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL suggest using normative language

#### Scenario: Incorrect scenario format
- **WHEN** scenarios use formats other than "#### Scenario:" (e.g., bullets or bold text)
- **THEN** validation SHALL report an ERROR
- **AND** the message SHALL show the correct "#### Scenario:" header format
- **AND** detection SHALL use AST heading level 4 with "Scenario:" prefix

### Requirement: Change Delta Validation
The validation system SHALL validate change delta specs for structural correctness and delta operation validity.

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
