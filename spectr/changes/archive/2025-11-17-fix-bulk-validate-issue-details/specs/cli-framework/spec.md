# Cli Framework Specification Delta

## MODIFIED Requirements

### Requirement: Validate Command Structure

The CLI SHALL provide a validate command for checking spec and change document
correctness.

#### Scenario: Validate command registration

- **WHEN** the CLI is initialized
- **THEN** it SHALL include a ValidateCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr validate`
- **AND** help text SHALL describe validation functionality

#### Scenario: Direct item validation invocation

- **WHEN** user invokes `spectr validate <item-name>`
- **THEN** the command SHALL validate the named item (change or spec)
- **AND** SHALL print validation results to stdout
- **AND** SHALL exit with code 0 for valid, 1 for invalid

#### Scenario: Bulk validation invocation

- **WHEN** user invokes `spectr validate --all`
- **THEN** the command SHALL validate all changes and specs
- **AND** SHALL print summary of results
- **AND** SHALL display full issue details for each failed item including level,
  path, and message
- **AND** SHALL exit with code 1 if any item fails validation

#### Scenario: Interactive validation invocation

- **WHEN** user invokes `spectr validate` without arguments in a TTY
- **THEN** the command SHALL prompt for what to validate
- **AND** SHALL execute the user's selection
