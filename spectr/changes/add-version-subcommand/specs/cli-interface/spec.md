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

### Requirement: Build-Time Version Injection
The version command SHALL support version metadata injection at build time via Go linker flags (ldflags), integrated with GoReleaser for release builds.

#### Scenario: GoReleaser injects version metadata
- **WHEN** Spectr is built via GoReleaser
- **THEN** the version variable is set to the git tag (e.g., "v1.2.3")
- **AND** the commit variable is set to the short git commit hash
- **AND** the date variable is set to the build date in ISO 8601 format
- **AND** all values are displayed in the version output

#### Scenario: Manual build with ldflags
- **WHEN** developer builds with `go build -ldflags "-X main.version=v1.0.0"`
- **THEN** the version command displays the injected version number
- **AND** other fields display their default or injected values

### Requirement: Development Build Handling
The version command SHALL gracefully handle cases where version metadata is not injected, displaying "dev" or appropriate defaults.

#### Scenario: Build without ldflags injection
- **WHEN** user builds Spectr with `go build` without ldflags
- **THEN** the version command displays "dev" as the version
- **AND** commit displays "unknown" or empty string
- **AND** date displays "unknown" or empty string
- **AND** goVersion displays the runtime Go version
- **AND** os and arch display the current platform values

#### Scenario: Development build identification
- **WHEN** user runs version command on a development build
- **THEN** the output clearly indicates this is a development build
- **AND** the user understands this is not an official release
- **AND** the command still exits successfully with code 0

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
