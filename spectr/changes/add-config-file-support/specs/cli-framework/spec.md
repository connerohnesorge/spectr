# CLI Framework Specification Delta

## ADDED Requirements

### Requirement: Kong Resolver Support for Configuration
The CLI framework SHALL use Kong's resolver mechanism to provide configuration values from multiple sources (config files, environment variables) during parsing.

#### Scenario: Register custom resolvers
- **WHEN** initializing Kong in `main.go`
- **THEN** the system SHALL register custom resolvers using `kong.Resolvers()` option
- **AND** SHALL register environment variable resolver
- **AND** SHALL register user config file resolver
- **AND** resolvers SHALL be queried in precedence order

#### Scenario: Resolver provides flag values
- **WHEN** Kong parses flags and encounters a flag not provided on command line
- **THEN** Kong SHALL query each resolver in order
- **AND** SHALL use the first non-nil value returned
- **AND** SHALL fall back to struct field default if all resolvers return nil

#### Scenario: Resolver respects command hierarchy
- **WHEN** resolving flags for nested commands (e.g., `spectr validate --strict`)
- **THEN** the resolver SHALL check command-specific configuration first
- **AND** SHALL fall back to global configuration if command-specific not found

### Requirement: Environment Variable Tag Support
The CLI framework SHALL support `envar` struct tags on command fields to enable environment variable mapping.

#### Scenario: Define environment variable mapping
- **WHEN** a struct field defines `envar` tag like `envar:"SPECTR_STRICT"`
- **THEN** Kong SHALL automatically read from the `SPECTR_STRICT` environment variable
- **AND** SHALL use the environment value if CLI flag not provided
- **AND** SHALL respect precedence (CLI flag > env var)

#### Scenario: Environment variable for boolean flag
- **WHEN** a boolean flag has `envar` tag and env var is set to `true` or `false`
- **THEN** Kong SHALL parse the string value to boolean
- **AND** SHALL error on invalid boolean values

#### Scenario: Environment variable for string flag
- **WHEN** a string flag has `envar` tag and env var is set
- **THEN** Kong SHALL use the string value directly

#### Scenario: Environment variable for slice flag
- **WHEN** a slice flag has `envar` tag and env var is set
- **THEN** Kong SHALL parse comma-separated values into slice elements

### Requirement: Config Command Structure
The CLI framework SHALL include a `config` command with subcommands for managing configuration files.

#### Scenario: Config command registration
- **WHEN** the CLI is initialized
- **THEN** it SHALL include a ConfigCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr config`
- **AND** help text SHALL describe configuration management

#### Scenario: Config subcommands
- **WHEN** the ConfigCmd struct is defined
- **THEN** it SHALL have subcommand fields: Init, Show, Edit, Validate
- **AND** each SHALL be tagged with `cmd` and appropriate help text
- **AND** each SHALL implement a `Run() error` method

### Requirement: Resolver Error Handling
The CLI framework SHALL handle resolver errors gracefully and provide clear error messages.

#### Scenario: Config file parse error
- **WHEN** a resolver encounters a malformed config file
- **THEN** Kong SHALL exit with error before executing command
- **AND** error message SHALL include file path and parsing details

#### Scenario: Environment variable type mismatch
- **WHEN** an environment variable value cannot be parsed to the expected flag type
- **THEN** Kong SHALL exit with error describing the type mismatch
- **AND** error message SHALL include env var name and expected type

#### Scenario: Resolver initialization failure
- **WHEN** a resolver fails to initialize (e.g., config file permission error)
- **THEN** Kong SHALL exit with error before parsing
- **AND** error message SHALL describe the initialization failure

## MODIFIED Requirements

### Requirement: Struct-Based Command Definition
The CLI framework SHALL use Go struct types with struct tags to declaratively define command structure, subcommands, flags, arguments, and environment variable mappings.

#### Scenario: Root command definition
- **WHEN** the CLI is initialized
- **THEN** it SHALL use a root struct with subcommand fields tagged with `cmd` for command definitions
- **AND** each subcommand SHALL be a nested struct type with appropriate tags

#### Scenario: Subcommand registration
- **WHEN** a new subcommand is added to the CLI
- **THEN** it SHALL be defined as a struct field on the parent command struct
- **AND** it SHALL use `cmd` tag to indicate it is a subcommand
- **AND** it SHALL include a `help` tag describing the command purpose

#### Scenario: Struct fields with environment variable support
- **WHEN** a command struct field should support environment variables
- **THEN** it SHALL include an `envar` tag with the environment variable name
- **AND** the tag value SHALL follow `SPECTR_` prefix convention
- **AND** the field SHALL still support CLI flag and config file precedence
