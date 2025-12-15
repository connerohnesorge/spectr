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
- **AND** detection SHALL use `markdown.MatchScenarioHeader` from the markdown parser package

#### Scenario: Parsing uses markdown parser package
- **WHEN** the validation system parses spec or delta files
- **THEN** it SHALL use functions from `internal/markdown/` package
- **AND** it SHALL NOT import `internal/regex/` package or define local regex patterns
- **AND** behavior SHALL be identical to previous implementation

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
- **THEN** validation SHALL accept the renames using `markdown.MatchRenamedFrom` and `markdown.MatchRenamedTo` functions
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Delta parsing uses markdown parser package
- **WHEN** the validation system parses delta spec files
- **THEN** it SHALL use functions from `internal/markdown/` package
- **AND** delta type detection SHALL use `markdown.MatchDeltaSection`
- **AND** requirement extraction SHALL use `markdown.MatchRequirementHeader`
- **AND** scenario extraction SHALL use `markdown.MatchScenarioHeader`
- **AND** section content extraction SHALL use `markdown.FindSection` or `markdown.FindDeltaSection`

## ADDED Requirements

### Requirement: Wikilink Validation
The validation system SHALL validate wikilinks in spec files to ensure all targets exist and are resolvable.

#### Scenario: Valid wikilink to spec
- **WHEN** a spec file contains `[[validation]]`
- **AND** `spectr/specs/validation/spec.md` exists
- **THEN** validation SHALL pass with no wikilink errors

#### Scenario: Valid wikilink to change
- **WHEN** a spec file contains `[[changes/my-change]]`
- **AND** `spectr/changes/my-change/proposal.md` exists
- **THEN** validation SHALL pass with no wikilink errors

#### Scenario: Invalid wikilink target
- **WHEN** a spec file contains `[[nonexistent-spec]]`
- **AND** no matching spec or change exists
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL indicate the unresolved wikilink target
- **AND** in strict mode validation SHALL fail (valid=false)

#### Scenario: Wikilink with invalid anchor
- **WHEN** a spec file contains `[[validation#Nonexistent Requirement]]`
- **AND** the target spec exists but the anchor does not
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL indicate the unresolved anchor

