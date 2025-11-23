## ADDED Requirements

### Requirement: Version Command
The CLI SHALL provide a `version` subcommand that displays version and build information to help users identify which version of Spectr they are running.

#### Scenario: User runs version command
- **WHEN** user runs `spectr version`
- **THEN** the system displays version information in human-readable format
- **AND** output includes version number, git commit, build date, Go version, OS, and architecture
- **AND** the command exits with code 0

#### Scenario: Help text available for version command
- **WHEN** user runs `spectr version --help`
- **THEN** the system displays help text explaining the version command and its flags
- **AND** the help text lists `--short` and `--json` flags

### Requirement: Version Short Output
The version command SHALL support a `--short` flag that outputs only the version number for use in scripts and automation.

#### Scenario: User requests short version output
- **WHEN** user runs `spectr version --short`
- **THEN** the system outputs only the version number (e.g., "v1.2.3" or "dev")
- **AND** no other information is displayed
- **AND** output is suitable for parsing in shell scripts

#### Scenario: Short output with piping
- **WHEN** user runs `spectr version --short` and pipes to another command
- **THEN** the output contains only the version string with no formatting or colors
- **AND** the output is a single line with no trailing whitespace except newline

### Requirement: Version JSON Output
The version command SHALL support a `--json` flag that outputs version information in JSON format for machine consumption.

#### Scenario: User requests JSON version output
- **WHEN** user runs `spectr version --json`
- **THEN** the system outputs valid JSON containing all version fields
- **AND** JSON includes fields: version, commit, date, goVersion, os, arch
- **AND** JSON is properly formatted and parseable

#### Scenario: JSON output structure validation
- **WHEN** user parses the JSON output from `spectr version --json`
- **THEN** each field is a string value
- **AND** the version field contains the version number
- **AND** the commit field contains the git commit hash (short form, 7 chars)
- **AND** the date field contains the build date in ISO 8601 format
- **AND** the goVersion field contains the Go version (e.g., "go1.25.0")
- **AND** the os and arch fields contain runtime.GOOS and runtime.GOARCH values

### Requirement: Version Command Integration
The version command SHALL be integrated into the root CLI structure following the established Kong command pattern used by other subcommands.

#### Scenario: Version command in CLI structure
- **WHEN** developer examines the CLI struct in `cmd/root.go`
- **THEN** a Version field of type VersionCmd exists
- **AND** the field has appropriate Kong tags for help text
- **AND** the field follows the same pattern as other commands (Init, List, Validate, etc.)

#### Scenario: Version command implementation pattern
- **WHEN** developer examines `cmd/version.go`
- **THEN** VersionCmd struct has Short and JSON boolean fields with Kong tags
- **AND** VersionCmd has a Run() error method
- **AND** the file includes comprehensive doc comments
- **AND** the code follows the same patterns as other cmd files

### Requirement: Version Metadata from Embedded File
The version command SHALL read version metadata from an embedded VERSION file in JSON format using Go's embed package, supporting both Nix builds and GoReleaser builds.

#### Scenario: VERSION file embedded and parsed correctly
- **WHEN** Spectr is built via any build system (GoReleaser, Nix, or go build)
- **THEN** the VERSION file is embedded in the binary at compile time
- **AND** the version command reads and parses the JSON file at runtime
- **AND** the JSON contains "version", "commit", and "date" fields as strings
- **AND** all values from the file are displayed in the version output

#### Scenario: Custom VERSION file support
- **WHEN** developer creates a custom VERSION file with specific values before build
- **THEN** `go build` embeds the custom file into the binary
- **AND** the version command displays the custom version, commit, and date values
- **AND** the JSON format is validated during parsing

#### Scenario: VERSION file format validation
- **WHEN** the embedded VERSION file is parsed
- **THEN** the system validates it contains valid JSON
- **AND** the system checks for required fields: version, commit, date
- **AND** if parsing fails, the system falls back to development defaults
- **AND** the command still exits successfully with code 0

### Requirement: Development Build Handling
The version command SHALL gracefully handle cases where the VERSION file is missing or contains empty values, displaying "dev" or appropriate defaults.

#### Scenario: Build without VERSION file
- **WHEN** user builds Spectr with `go build` without a VERSION file
- **THEN** the version command displays "dev" as the version
- **AND** commit displays "unknown"
- **AND** date displays "unknown"
- **AND** goVersion displays the runtime Go version
- **AND** os and arch display the current platform values

#### Scenario: VERSION file with empty values
- **WHEN** the VERSION file exists but contains empty strings for fields
- **THEN** empty version field defaults to "dev"
- **AND** empty commit field defaults to "unknown"
- **AND** empty date field defaults to "unknown"
- **AND** the version command still operates normally

#### Scenario: Development build identification
- **WHEN** user runs version command on a development build
- **THEN** the output clearly indicates this is a development build (shows "dev")
- **AND** the user understands this is not an official release
- **AND** the command still exits successfully with code 0
