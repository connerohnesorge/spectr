## ADDED Requirements

### Requirement: Tasks File Structure Validation
The validation system SHALL validate `tasks.md` files for structural correctness, ensuring tasks are organized under numbered section headers.

#### Scenario: Valid tasks.md with numbered sections
- **WHEN** a `tasks.md` file contains numbered section headers (`## 1.`, `## 2.`, etc.) with tasks under each
- **THEN** validation SHALL pass with no structure warnings
- **AND** the validation report SHALL confirm proper structure

#### Scenario: Tasks.md missing numbered sections
- **WHEN** a `tasks.md` file exists but contains no `## [number].` section headers
- **THEN** validation SHALL report a WARNING level issue
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include example of correct format: `## 1. Section Name`

#### Scenario: Tasks outside numbered sections
- **WHEN** task items (`- [ ]` or `- [x]`) appear before any numbered section header
- **THEN** validation SHALL report a WARNING for orphaned tasks
- **AND** the warning SHALL suggest grouping tasks under numbered sections

#### Scenario: Empty numbered section
- **WHEN** a numbered section header exists but contains no task items before the next section or end of file
- **THEN** validation SHALL report a WARNING for the empty section
- **AND** the warning SHALL identify which section number is empty

#### Scenario: Non-sequential section numbers
- **WHEN** section numbers skip values (e.g., `## 1.`, `## 3.`, missing `## 2.`)
- **THEN** validation SHALL report a WARNING about non-sequential numbering
- **AND** the warning SHALL identify the gap in sequence

#### Scenario: Valid section numbering
- **WHEN** sections are numbered sequentially starting from 1 (`## 1.`, `## 2.`, `## 3.`, etc.)
- **THEN** validation SHALL pass without sequence warnings

#### Scenario: Tasks.md file not found
- **WHEN** a change directory exists but `tasks.md` is missing
- **THEN** validation SHALL report a WARNING that tasks.md is recommended
- **AND** this SHALL NOT be an error (tasks.md remains optional for backward compatibility)
