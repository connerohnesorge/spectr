# Delta Specification

## ADDED Requirements

### Requirement: Version Command Structure

The CLI SHALL provide a `version` command that displays version information
including version number, git commit hash, and build date.

#### Scenario: Version command registration

- **WHEN** the CLI is initialized
- **THEN** it SHALL include a VersionCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr version`
- **AND** help text SHALL describe version display functionality

#### Scenario: Version command invocation

- **WHEN** user runs `spectr version` without flags
- **THEN** the system displays version in format: `spectr version {version}
  (commit: {commit}, built: {date})`
- **AND** version SHALL be the semantic version (e.g., `0.1.0` or `dev`)
- **AND** commit SHALL be the git commit hash (short or full) or `unknown`
- **AND** date SHALL be the build date in ISO 8601 format or `unknown`

#### Scenario: Version command with short flag

- **WHEN** user runs `spectr version --short`
- **THEN** the system displays only the version number (e.g., `0.1.0`)
- **AND** no other information is displayed

#### Scenario: Version command with JSON flag

- **WHEN** user runs `spectr version --json`
- **THEN** the system outputs version data as JSON
- **AND** JSON SHALL include fields: `version`, `commit`, `date`
- **AND** SHALL be parseable by standard JSON tools

### Requirement: Version Variable Injection

The version information SHALL be injectable at build time via Go ldflags,
supporting both goreleaser releases and nix flake builds.

#### Scenario: Goreleaser version injection

- **WHEN** goreleaser builds the binary
- **THEN** version SHALL be set from git tag via ldflags
- **AND** commit SHALL be set from git commit hash via ldflags
- **AND** date SHALL be set from build timestamp via ldflags

#### Scenario: Nix flake version injection

- **WHEN** nix builds the binary via flake.nix
- **THEN** version SHALL be set from the flake package version attribute via
  ldflags
- **AND** commit and date MAY be `unknown` if not available in nix build context

#### Scenario: Development build defaults

- **WHEN** binary is built without ldflags (e.g., `go build`)
- **THEN** version SHALL default to `dev`
- **AND** commit SHALL default to `unknown`
- **AND** date SHALL default to `unknown`

### Requirement: Version Package Location

The version variables SHALL be defined in a dedicated `internal/version` package
for clean separation and easy ldflags targeting.

#### Scenario: Package structure

- **WHEN** the version package is imported
- **THEN** it SHALL expose `Version`, `Commit`, and `Date` string variables
- **AND** variables SHALL have default values for development builds
- **AND** the ldflags path SHALL be
  `github.com/connerohnesorge/spectr/internal/version`
