## ADDED Requirements

### Requirement: Consolidated Regex Package

The system SHALL provide a dedicated `internal/regex/` package that consolidates all markdown parsing patterns used across the codebase, with files split by category and pre-compiled patterns at package initialization.

#### Scenario: Package structure by category

- **WHEN** the regex package is imported
- **THEN** it SHALL have separate files: `headers.go`, `tasks.go`, `renames.go`, `sections.go`
- **AND** each file SHALL contain semantically related patterns
- **AND** test files SHALL be co-located (`headers_test.go`, etc.)

#### Scenario: Pre-compiled patterns at init

- **WHEN** the regex package is imported
- **THEN** all patterns SHALL be compiled once at package initialization via package-level `var` with `regexp.MustCompile`
- **AND** subsequent pattern usage SHALL NOT require recompilation
- **AND** invalid patterns SHALL cause a panic at program startup (fail-fast)

#### Scenario: Header patterns available

- **WHEN** code needs to match markdown headers
- **THEN** it SHALL have access to `regex.H2SectionHeader` for `## Section Name`
- **AND** it SHALL have access to `regex.H2DeltaSection` for `## ADDED|MODIFIED|REMOVED|RENAMED Requirements`
- **AND** it SHALL have access to `regex.H2RequirementsSection` for exactly `## Requirements`
- **AND** it SHALL have access to `regex.H2NextSection` for finding section boundaries
- **AND** it SHALL have access to `regex.H3Requirement` for `### Requirement: Name`
- **AND** it SHALL have access to `regex.H3AnyHeader` for any `###` header
- **AND** it SHALL have access to `regex.H4Scenario` for `#### Scenario: Name`

#### Scenario: Task patterns available

- **WHEN** code needs to match task checkboxes
- **THEN** it SHALL have access to `regex.TaskCheckbox` for `- [ ]` and `- [x]` items
- **AND** it SHALL have access to `regex.NumberedTask` for `- [ ] 1.1 Description` format
- **AND** it SHALL have access to `regex.NumberedSection` for `## 1. Section Name` format

#### Scenario: Rename patterns available as separate variants

- **WHEN** code needs to match RENAMED section entries
- **THEN** it SHALL have access to `regex.RenamedFrom` for backtick-wrapped FROM lines
- **AND** it SHALL have access to `regex.RenamedTo` for backtick-wrapped TO lines
- **AND** it SHALL have access to `regex.RenamedFromAlt` for non-backtick FROM lines
- **AND** it SHALL have access to `regex.RenamedToAlt` for non-backtick TO lines
- **AND** both variants SHALL be exported separately (not combined into single pattern)

#### Scenario: Helper functions use (value, ok) return style

- **WHEN** code needs to extract values from pattern matches
- **THEN** `regex.MatchH3Requirement(line)` SHALL return `(name string, ok bool)`
- **AND** `regex.MatchH4Scenario(line)` SHALL return `(name string, ok bool)`
- **AND** `regex.MatchTaskCheckbox(line)` SHALL return `(state rune, ok bool)` with state normalized to lowercase
- **AND** `regex.MatchH2DeltaSection(line)` SHALL return `(deltaType string, ok bool)`
- **AND** `regex.MatchNumberedTask(line)` SHALL return `(checkbox, id, desc string, ok bool)`

#### Scenario: Section content extraction helpers

- **WHEN** code needs to extract content between markdown headers
- **THEN** `regex.FindSectionContent(content, sectionHeader)` SHALL return content from the specified H2 section to the next H2
- **AND** `regex.FindDeltaSectionContent(content, deltaType)` SHALL be a convenience wrapper for delta sections
- **AND** `regex.FindRequirementsSection(content)` SHALL extract the `## Requirements` section content
- **AND** empty string SHALL be returned if section not found

#### Scenario: Both patterns and helpers exported

- **WHEN** consuming code uses the regex package
- **THEN** raw `*regexp.Regexp` patterns SHALL be exported for advanced use cases
- **AND** helper functions SHALL be exported for common matching operations
- **AND** callers MAY use either patterns directly or helper functions

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
- **AND** detection SHALL use `regex.H4Scenario` pattern from the consolidated regex package

#### Scenario: Parsing uses consolidated regex package

- **WHEN** the validation system parses spec or delta files
- **THEN** it SHALL use patterns from `internal/regex/` package
- **AND** it SHALL NOT define local regex patterns for structural markdown elements
- **AND** behavior SHALL be identical to previous inline pattern implementation

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
- **THEN** validation SHALL accept the renames using `regex.RenamedFromAlt` and `regex.RenamedToAlt` patterns
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Delta parsing uses consolidated regex package

- **WHEN** the validation system parses delta spec files
- **THEN** it SHALL use patterns from `internal/regex/` package
- **AND** delta type detection SHALL use `regex.MatchH2DeltaSection`
- **AND** requirement extraction SHALL use `regex.MatchH3Requirement`
- **AND** scenario extraction SHALL use `regex.MatchH4Scenario`
- **AND** section content extraction SHALL use `regex.FindSectionContent` or `regex.FindDeltaSectionContent`
