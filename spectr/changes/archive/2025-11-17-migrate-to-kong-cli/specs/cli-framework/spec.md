# CLI Framework Specification

## ADDED Requirements

### Requirement: Struct-Based Command Definition

The CLI framework SHALL use Go struct types with struct tags to declaratively define command structure, subcommands, flags, and arguments.

#### Scenario: Root command definition

- **WHEN** the CLI is initialized
- **THEN** it SHALL use a root struct with subcommand fields tagged with `cmd` for command definitions
- **AND** each subcommand SHALL be a nested struct type with appropriate tags

#### Scenario: Subcommand registration

- **WHEN** a new subcommand is added to the CLI
- **THEN** it SHALL be defined as a struct field on the parent command struct
- **AND** it SHALL use `cmd` tag to indicate it is a subcommand
- **AND** it SHALL include a `help` tag describing the command purpose

### Requirement: Declarative Flag Definition

The CLI framework SHALL define flags using struct fields with Kong struct tags instead of imperative flag registration.

#### Scenario: String flag definition

- **WHEN** a command requires a string flag
- **THEN** it SHALL be defined as a struct field with `name` tag for the flag name
- **AND** it MAY include `short` tag for single-character shorthand
- **AND** it SHALL include `help` tag describing the flag purpose
- **AND** it MAY include `default` tag for default values

#### Scenario: Boolean flag definition

- **WHEN** a command requires a boolean flag
- **THEN** it SHALL be defined as a bool struct field with appropriate tags
- **AND** the flag SHALL default to false unless explicitly set

#### Scenario: Slice flag definition

- **WHEN** a command requires a multi-value flag
- **THEN** it SHALL be defined as a slice type struct field
- **AND** it SHALL support comma-separated values or repeated flag usage

### Requirement: Positional Argument Support

The CLI framework SHALL support positional arguments using struct fields tagged with `arg`.

#### Scenario: Optional positional argument

- **WHEN** a command accepts an optional positional argument
- **THEN** it SHALL be defined with `arg` and `optional` tags
- **AND** the field SHALL be a pointer type or have a zero value for "not provided"

#### Scenario: Required positional argument

- **WHEN** a command requires a positional argument
- **THEN** it SHALL be defined with `arg` tag without `optional`
- **AND** parsing SHALL fail if the argument is not provided

### Requirement: Automatic Method Dispatch

The CLI framework SHALL automatically invoke the appropriate command's Run method after parsing.

#### Scenario: Command execution

- **WHEN** a command is successfully parsed
- **THEN** the framework SHALL call the command struct's `Run() error` method
- **AND** it SHALL pass any configured context values to the Run method
- **AND** it SHALL handle the returned error appropriately

### Requirement: Built-in Help Generation

The CLI framework SHALL automatically generate help text from struct tags and types.

#### Scenario: Root help display

- **WHEN** the CLI is invoked with `--help` or no arguments
- **THEN** it SHALL display a list of available subcommands
- **AND** it SHALL show descriptions from `help` tags
- **AND** it SHALL indicate required vs optional arguments

#### Scenario: Subcommand help display

- **WHEN** a subcommand is invoked with `--help`
- **THEN** it SHALL display the command description
- **AND** it SHALL list all flags with their types and help text
- **AND** it SHALL show positional argument requirements

### Requirement: Error Handling and Exit Codes

The CLI framework SHALL provide appropriate error messages and exit codes for parsing and execution failures.

#### Scenario: Parse error handling

- **WHEN** invalid flags or arguments are provided
- **THEN** it SHALL display an error message
- **AND** it SHALL show usage information
- **AND** it SHALL exit with non-zero status code

#### Scenario: Execution error handling

- **WHEN** a command's Run method returns an error
- **THEN** it SHALL display the error message
- **AND** it SHALL exit with non-zero status code

### Requirement: Backward-Compatible CLI Interface

The CLI framework SHALL maintain the same command syntax and flag names as the previous implementation.

#### Scenario: Init command compatibility

- **WHEN** users invoke `spectr init` with existing flag combinations
- **THEN** the behavior SHALL be identical to the previous Cobra-based implementation
- **AND** all flag names SHALL remain unchanged
- **AND** short flag aliases SHALL remain unchanged
- **AND** positional argument handling SHALL remain unchanged

#### Scenario: Help text accessibility

- **WHEN** users invoke `spectr --help` or `spectr init --help`
- **THEN** help information SHALL be displayed (format may differ from Cobra)
- **AND** all commands and flags SHALL be documented
