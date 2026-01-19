# Delta Spec: CLI Interface - JSONC Validation

## ADDED Requirements

### Requirement: JSONC Output Validation

The system SHALL validate all generated JSONC output to ensure it can be successfully parsed.

#### Scenario: Validate JSONC after generation
- WHEN `writeTasksJSONC`, `writeRootTasksJSONC`, or `writeChildTasksJSONC` completes
- THEN verify the generated JSONC can be parsed via `StripJSONComments` + `json.Unmarshal`
- AND return error if parsing fails

#### Scenario: Validate round-trip conversion
- WHEN converting Task → JSONC → Task
- THEN the resulting Task MUST be identical to the original
- AND no data loss occurs in any field

#### Scenario: Handle special characters in task descriptions
- WHEN task descriptions contain backslashes, quotes, newlines, tabs, or unicode
- THEN json.Marshal MUST produce valid escape sequences
- AND the JSONC MUST parse correctly when read back

#### Scenario: Validate with external JSONC parser
- WHEN JSONC validation is enabled
- THEN use an external JSONC parser library to verify format correctness
- AND report detailed errors if validation fails

### Requirement: Property-Based JSONC Testing

The system SHALL test JSONC generation with randomized inputs to find edge cases.

#### Scenario: Test special character combinations
- WHEN running property-based tests
- THEN generate task descriptions with special characters: `\`, `"`, `\n`, `\t`, `\r`, `\b`, `\f`
- AND verify all combinations produce valid JSONC

#### Scenario: Test unicode and emoji
- WHEN task descriptions contain unicode or emoji
- THEN JSONC output MUST preserve the characters correctly
- AND round-trip conversion MUST be lossless

#### Scenario: Test JSON meta-characters
- WHEN task descriptions contain `{`, `}`, `[`, `]`, `:`, `,`
- THEN json.Marshal MUST escape them correctly
- AND parsing MUST not confuse them with JSON structure

#### Scenario: Test JSONC-like comments in descriptions
- WHEN task descriptions contain `//` or `/* */`
- THEN json.Marshal MUST escape them correctly
- AND `StripJSONComments` MUST not remove them from the description

### Requirement: JSONC Validation Error Reporting

The system SHALL provide detailed error messages when JSONC validation fails.

#### Scenario: Report invalid escape sequences
- WHEN JSONC validation fails due to invalid escape sequence
- THEN error MUST include: task ID, description excerpt, character position
- AND suggest the correct escape sequence

#### Scenario: Report round-trip data loss
- WHEN round-trip validation detects data loss
- THEN error MUST include: task ID, original value, parsed value
- AND identify the specific field that lost data

#### Scenario: Fail fast on write
- WHEN JSONC validation fails
- THEN do not write the file to disk
- AND return error immediately to prevent corruption

## MODIFIED Requirements

### Requirement: Accept Command Structure

The system SHALL convert `tasks.md` to `tasks.jsonc` with validation.

#### Scenario: Accept command registration
- WHEN initialized
- THEN register `spectr accept`

#### Scenario: Accept with change ID
- WHEN `spectr accept <id>`
- THEN validate change, parse tasks.md, generate JSONC
- AND validate JSONC output before writing
- AND preserve tasks.md

#### Scenario: Accept with validation
- WHEN running
- THEN abort if change validation fails
- OR abort if JSONC validation fails

#### Scenario: JSONC validation failure
- WHEN JSONC validation fails during accept
- THEN display detailed error with task ID and description
- AND suggest fixes for common issues (escape sequences, special chars)
- AND exit with non-zero status

#### Scenario: Interactive change selection
- WHEN no change ID provided
- THEN prompt for change selection
- AND validate selected change before conversion
