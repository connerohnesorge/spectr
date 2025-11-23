# Config Management Specification

## Purpose

This specification defines configuration file support for Spectr, including YAML config file loading from user directories, environment variable support, precedence resolution, and configuration management commands for initializing, viewing, editing, and validating user preferences.

## ADDED Requirements

### Requirement: User Configuration File Loading
The system SHALL load user configuration from a YAML file located at `~/.config/spectr/config.yaml` on Linux/macOS, respecting the XDG Base Directory Specification.

#### Scenario: Load existing config file
- **WHEN** the CLI starts and `~/.config/spectr/config.yaml` exists
- **THEN** the system SHALL read and parse the YAML file
- **AND** SHALL merge configuration values with command defaults

#### Scenario: Missing config file
- **WHEN** the CLI starts and `~/.config/spectr/config.yaml` does not exist
- **THEN** the system SHALL continue with default values
- **AND** SHALL NOT error or warn about missing config file

#### Scenario: XDG_CONFIG_HOME override
- **WHEN** the `XDG_CONFIG_HOME` environment variable is set
- **THEN** the system SHALL use `$XDG_CONFIG_HOME/spectr/config.yaml` instead of `~/.config/spectr/config.yaml`

#### Scenario: Invalid YAML syntax
- **WHEN** the config file contains invalid YAML syntax
- **THEN** the system SHALL exit with an error message including line number and syntax details
- **AND** SHALL NOT proceed with command execution

#### Scenario: Unknown configuration keys
- **WHEN** the config file contains keys that do not correspond to any CLI flags
- **THEN** the system SHALL warn about unknown keys
- **AND** SHALL continue execution with valid keys

### Requirement: Environment Variable Support
The system SHALL support environment variables prefixed with `SPECTR_` that map to CLI flags, with hyphens converted to underscores and names uppercased.

#### Scenario: Boolean flag from environment variable
- **WHEN** the environment variable `SPECTR_STRICT=true` is set
- **THEN** the `--strict` flag SHALL be enabled for all commands
- **AND** SHALL behave identically to passing `--strict` on the command line

#### Scenario: String flag from environment variable
- **WHEN** the environment variable `SPECTR_TYPE=spec` is set
- **THEN** the `--type` flag SHALL use the value `spec`

#### Scenario: Invalid boolean value
- **WHEN** an environment variable like `SPECTR_STRICT=invalid` contains a non-boolean value
- **THEN** the system SHALL exit with an error indicating the invalid value and expected type

#### Scenario: Command-specific environment variables
- **WHEN** environment variables like `SPECTR_VALIDATE_STRICT=true` are set
- **THEN** the system SHALL apply them only to the `validate` command
- **AND** SHALL NOT affect other commands

### Requirement: Configuration Precedence
The system SHALL resolve configuration values using precedence order: CLI flags > Environment variables > User config file > Hard-coded defaults.

#### Scenario: CLI flag overrides environment variable
- **WHEN** `SPECTR_STRICT=false` is set and user runs `spectr validate --strict`
- **THEN** the `--strict` flag SHALL be enabled (CLI flag wins)

#### Scenario: Environment variable overrides config file
- **WHEN** config file sets `strict: false` and `SPECTR_STRICT=true` is set
- **THEN** the `--strict` flag SHALL be enabled (env var wins)

#### Scenario: Config file overrides defaults
- **WHEN** config file sets `strict: true` and no CLI flag or env var is provided
- **THEN** the `--strict` flag SHALL be enabled (config file wins)

#### Scenario: Default used when no config provided
- **WHEN** no config file, env var, or CLI flag sets a value
- **THEN** the system SHALL use the hard-coded default value

### Requirement: Configuration File Schema
The configuration file SHALL use YAML format with structure mirroring the CLI command hierarchy, supporting both global and command-specific settings.

#### Scenario: Global configuration settings
- **WHEN** config file includes top-level keys like `json: true`
- **THEN** these settings SHALL apply to all commands that support those flags

#### Scenario: Command-specific configuration
- **WHEN** config file includes nested structure like `validate:\n  strict: true`
- **THEN** the `strict` setting SHALL apply only to the `validate` command
- **AND** SHALL NOT affect other commands

#### Scenario: Command-specific overrides global
- **WHEN** config file has both global `json: false` and command-specific `validate:\n  json: true`
- **THEN** the `validate` command SHALL use `json: true`
- **AND** other commands SHALL use `json: false`

### Requirement: Config Init Command
The system SHALL provide `spectr config init` command to scaffold a default configuration file with comments and examples.

#### Scenario: Initialize config file
- **WHEN** user runs `spectr config init`
- **THEN** the system SHALL create `~/.config/spectr/` directory if it does not exist
- **AND** SHALL write `config.yaml` with commented examples for all available settings
- **AND** SHALL exit with success message showing the created file path

#### Scenario: Config file already exists
- **WHEN** user runs `spectr config init` and `config.yaml` already exists
- **THEN** the system SHALL prompt for confirmation before overwriting
- **AND** SHALL preserve existing file if user declines

#### Scenario: Non-interactive init with force flag
- **WHEN** user runs `spectr config init --force`
- **THEN** the system SHALL overwrite existing config file without prompting

#### Scenario: Template includes all commands
- **WHEN** the config template is generated
- **THEN** it SHALL include commented examples for all commands (init, list, validate, archive, view)
- **AND** SHALL document all available flags for each command
- **AND** SHALL include explanatory comments for precedence and syntax

### Requirement: Config Show Command
The system SHALL provide `spectr config show` command to display the merged effective configuration with source annotations.

#### Scenario: Show merged configuration
- **WHEN** user runs `spectr config show`
- **THEN** the system SHALL display all configuration values currently in effect
- **AND** SHALL indicate the source of each value (default, config file, env var, or CLI flag)
- **AND** SHALL format output as human-readable text by default

#### Scenario: Show with JSON output
- **WHEN** user runs `spectr config show --json`
- **THEN** the system SHALL output configuration as JSON
- **AND** SHALL include fields: key, value, source, command (if command-specific)

#### Scenario: Show with no config file
- **WHEN** user runs `spectr config show` and no config file exists
- **THEN** the system SHALL display default values
- **AND** SHALL indicate source as "default" for all values

#### Scenario: Show filtered by command
- **WHEN** user runs `spectr config show --command validate`
- **THEN** the system SHALL display only configuration applicable to the `validate` command
- **AND** SHALL include both global and command-specific settings

### Requirement: Config Edit Command
The system SHALL provide `spectr config edit` command to open the config file in the user's default editor.

#### Scenario: Edit existing config file
- **WHEN** user runs `spectr config edit`
- **THEN** the system SHALL open `~/.config/spectr/config.yaml` in the editor specified by `$EDITOR`
- **AND** SHALL wait for the editor to close before exiting

#### Scenario: Edit with missing config file
- **WHEN** user runs `spectr config edit` and no config file exists
- **THEN** the system SHALL prompt to initialize the config file first
- **AND** SHALL run `spectr config init` if user confirms

#### Scenario: EDITOR not set
- **WHEN** user runs `spectr config edit` and `$EDITOR` is not set
- **THEN** the system SHALL fall back to `vi` on Unix-like systems
- **AND** SHALL exit with error on Windows if no suitable editor found

#### Scenario: Edit with custom editor
- **WHEN** user runs `spectr config edit --editor nano`
- **THEN** the system SHALL use `nano` instead of `$EDITOR`

### Requirement: Config Validate Command
The system SHALL provide `spectr config validate` command to check configuration file syntax and schema correctness.

#### Scenario: Validate correct config file
- **WHEN** user runs `spectr config validate`
- **THEN** the system SHALL parse the config file
- **AND** SHALL validate all keys against known CLI flags
- **AND** SHALL validate all values match expected types
- **AND** SHALL exit with success message if valid

#### Scenario: Validate with syntax errors
- **WHEN** user runs `spectr config validate` on a file with YAML syntax errors
- **THEN** the system SHALL display error message with line number
- **AND** SHALL exit with non-zero status code

#### Scenario: Validate with unknown keys
- **WHEN** user runs `spectr config validate` on a file with unknown keys
- **THEN** the system SHALL warn about each unknown key
- **AND** SHALL exit with non-zero status code if `--strict` flag is provided
- **AND** SHALL exit with zero status code (warning only) if `--strict` is not provided

#### Scenario: Validate with type mismatches
- **WHEN** user runs `spectr config validate` on a file with wrong value types
- **THEN** the system SHALL report each type mismatch with expected and actual types
- **AND** SHALL exit with non-zero status code

#### Scenario: Validate missing config file
- **WHEN** user runs `spectr config validate` and no config file exists
- **THEN** the system SHALL exit with error message indicating file not found
- **AND** SHALL suggest running `spectr config init`

### Requirement: Kong Resolver Integration
The system SHALL integrate with Kong's resolver mechanism to provide configuration values during CLI parsing.

#### Scenario: Resolver loads config before parsing
- **WHEN** Kong begins parsing command-line arguments
- **THEN** the custom resolver SHALL already have loaded the user config file
- **AND** SHALL provide values when Kong queries for flag defaults

#### Scenario: Resolver returns nil for unset values
- **WHEN** Kong queries for a flag value that is not in the config file or env vars
- **THEN** the resolver SHALL return nil to allow Kong to continue to next resolver or default

#### Scenario: Resolver chains multiple sources
- **WHEN** multiple resolvers are configured (env var resolver, file resolver)
- **THEN** Kong SHALL query them in order: env var resolver first, file resolver second
- **AND** SHALL use the first non-nil value returned

#### Scenario: Resolver handles nested command paths
- **WHEN** resolving a flag for a subcommand like `spectr validate --strict`
- **THEN** the resolver SHALL check command-specific config first (`validate.strict`)
- **AND** SHALL fall back to global config if command-specific not found

### Requirement: Configuration Type Coercion
The system SHALL automatically coerce configuration values from YAML to the appropriate Go types expected by Kong struct fields.

#### Scenario: Boolean value coercion
- **WHEN** config file has `strict: true` as YAML boolean
- **THEN** the system SHALL convert to Go `bool` type

#### Scenario: String value coercion
- **WHEN** config file has `type: "spec"` as YAML string
- **THEN** the system SHALL convert to Go `string` type

#### Scenario: String slice coercion
- **WHEN** config file has `tools: [git, make]` as YAML array
- **THEN** the system SHALL convert to Go `[]string` type

#### Scenario: Invalid type coercion
- **WHEN** config file has a value that cannot be coerced to the expected type
- **THEN** the system SHALL exit with error describing the type mismatch
- **AND** SHALL indicate the config key and expected type

### Requirement: Configuration File Path Discovery
The system SHALL support multiple methods for locating the user configuration file.

#### Scenario: Default XDG location
- **WHEN** no environment variables or flags override the path
- **THEN** the system SHALL use `~/.config/spectr/config.yaml` on Unix-like systems

#### Scenario: XDG_CONFIG_HOME environment variable
- **WHEN** `XDG_CONFIG_HOME` is set to `/custom/path`
- **THEN** the system SHALL use `/custom/path/spectr/config.yaml`

#### Scenario: Home directory expansion
- **WHEN** resolving config path
- **THEN** the system SHALL expand `~` to the user's home directory
- **AND** SHALL handle missing home directory gracefully with error message

#### Scenario: Config path override flag
- **WHEN** user runs `spectr config show --config /custom/config.yaml`
- **THEN** the system SHALL use the specified path instead of default location

### Requirement: Error Handling for Configuration Loading
The system SHALL provide clear error messages when configuration loading fails.

#### Scenario: Unreadable config file
- **WHEN** config file exists but is not readable (permissions issue)
- **THEN** the system SHALL exit with error indicating permission problem
- **AND** SHALL display the attempted file path

#### Scenario: Malformed YAML
- **WHEN** config file contains malformed YAML
- **THEN** the system SHALL exit with error showing line number and column
- **AND** SHALL include the YAML parsing error message

#### Scenario: Config directory creation failure
- **WHEN** `spectr config init` cannot create `~/.config/spectr/` directory
- **THEN** the system SHALL exit with error indicating permission or disk space issue

### Requirement: Backward Compatibility
The system SHALL maintain 100% backward compatibility with existing CLI usage patterns.

#### Scenario: No config file present
- **WHEN** user has no config file and uses only CLI flags
- **THEN** behavior SHALL be identical to previous Spectr versions
- **AND** no warnings or messages about missing config shall appear

#### Scenario: Empty config file
- **WHEN** config file exists but is empty or contains only comments
- **THEN** the system SHALL use default values for all flags
- **AND** SHALL NOT error on empty config

#### Scenario: Existing scripts unchanged
- **WHEN** users run existing shell scripts with explicit CLI flags
- **THEN** all scripts SHALL continue working without modification
- **AND** config files SHALL NOT override explicit CLI flags
