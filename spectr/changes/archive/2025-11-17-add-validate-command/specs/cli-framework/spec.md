# CLI Framework Specification Delta

## ADDED Requirements

### Requirement: Validate Command Structure

The CLI SHALL provide a validate command for checking spec and change document correctness.

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
- **AND** SHALL exit with code 1 if any item fails validation

#### Scenario: Interactive validation invocation

- **WHEN** user invokes `spectr validate` without arguments in a TTY
- **THEN** the command SHALL prompt for what to validate
- **AND** SHALL execute the user's selection

### Requirement: Validate Command Flags

The validate command SHALL support flags for controlling validation behavior and output format.

#### Scenario: Strict mode flag

- **WHEN** user provides `--strict` flag
- **THEN** validation SHALL treat warnings as errors
- **AND** exit code SHALL be 1 if warnings exist
- **AND** validation report SHALL reflect strict mode in valid field

#### Scenario: JSON output flag

- **WHEN** user provides `--json` flag
- **THEN** output SHALL be formatted as JSON
- **AND** SHALL include items, summary, and version fields
- **AND** SHALL be parseable by standard JSON tools

#### Scenario: Type disambiguation flag

- **WHEN** user provides `--type change` or `--type spec`
- **THEN** the command SHALL treat the item as the specified type
- **AND** SHALL skip type auto-detection
- **AND** SHALL error if item does not exist as that type

#### Scenario: All items flag

- **WHEN** user provides `--all` flag
- **THEN** the command SHALL validate all changes and all specs
- **AND** SHALL run in bulk validation mode

#### Scenario: Changes only flag

- **WHEN** user provides `--changes` flag
- **THEN** the command SHALL validate all changes only
- **AND** SHALL skip specs

#### Scenario: Specs only flag

- **WHEN** user provides `--specs` flag
- **THEN** the command SHALL validate all specs only
- **AND** SHALL skip changes

#### Scenario: Non-interactive flag

- **WHEN** user provides `--no-interactive` flag
- **THEN** the command SHALL not prompt for input
- **AND** SHALL print usage hint if no item specified
- **AND** SHALL exit with code 1

### Requirement: Validate Command Help Text

The validate command SHALL provide comprehensive help documentation.

#### Scenario: Command help display

- **WHEN** user invokes `spectr validate --help`
- **THEN** help text SHALL describe validation purpose
- **AND** SHALL list all available flags with descriptions
- **AND** SHALL show usage examples for common scenarios
- **AND** SHALL indicate optional vs required arguments

### Requirement: Positional Argument Support for Item Name

The validate command SHALL accept an optional positional argument for the item to validate.

#### Scenario: Optional item name argument

- **WHEN** validate command is defined
- **THEN** it SHALL have an ItemName field tagged with `arg:"" optional:""`
- **AND** the field type SHALL be pointer to string or string with zero value check
- **AND** omitting the argument SHALL be valid (triggers interactive or bulk mode)

#### Scenario: Item name provided

- **WHEN** user provides item name as positional argument
- **THEN** the command SHALL validate that specific item
- **AND** SHALL auto-detect whether it's a change or spec
- **AND** SHALL respect --type flag if provided for disambiguation
