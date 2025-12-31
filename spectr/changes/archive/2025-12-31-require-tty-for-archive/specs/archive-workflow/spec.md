## ADDED Requirements

### Requirement: TTY Requirement for Archive Command

The system SHALL require a TTY (terminal) environment when executing the
`spectr archive` command to ensure only humans can perform archive operations.

#### Scenario: Archive in TTY environment succeeds

- **WHEN** `spectr archive` is executed in a terminal with TTY
- **THEN** the system proceeds with normal archive workflow

#### Scenario: Archive in non-TTY environment fails

- **WHEN** `spectr archive` is executed in a non-TTY environment (e.g., piped
  input, automated script, CI/CD)
- **THEN** the system returns an error indicating TTY is required
- **AND** no archive operation is performed

#### Scenario: Error message is clear

- **WHEN** archive fails due to missing TTY
- **THEN** the error message clearly states that archive requires an interactive
  terminal
- **AND** suggests running the command directly in a terminal

#### Scenario: TTY check occurs before any operations

- **WHEN** `spectr archive` is executed
- **THEN** the TTY check is performed before validation, task checking, or any
  other archive operations
- **AND** no files are modified if TTY check fails
