# Delta Specification

## ADDED Requirements

### Requirement: Tasks File Validation

The validation system SHALL validate that tasks.md files contain at least one
task item when present in a change directory.

#### Scenario: tasks.md with valid tasks

- **WHEN** a change directory contains a tasks.md file with task items (`- [ ]`
  or `- [x]`)
- **THEN** validation SHALL pass without task-related errors

#### Scenario: tasks.md exists but is empty

- **WHEN** a change directory contains a tasks.md file with no task items
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL indicate that no tasks were found
- **AND** the error SHALL include the path to the tasks.md file

#### Scenario: tasks.md with only section headers

- **WHEN** a change directory contains a tasks.md file with only section headers
  and no task items
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL indicate the expected task format (`- [ ]` or
  `- [x]`)

#### Scenario: tasks.md does not exist

- **WHEN** a change directory does not contain a tasks.md file
- **THEN** validation SHALL NOT report an error for the missing tasks.md
- **AND** validation SHALL proceed with other checks normally
