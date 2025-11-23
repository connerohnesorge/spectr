# Change: Add Version Subcommand

## Why

Users and automation systems need a standard way to determine which version of Spectr is installed for troubleshooting, compatibility verification, and dependency management. Currently, there is no built-in command to display version information, forcing users to rely on external methods like checking installation artifacts or git tags.

A version command provides:
- **Bug reporting**: Users can accurately report which version exhibits issues
- **Compatibility checking**: Scripts can verify minimum version requirements
- **Feature availability**: Users can determine if their version supports specific features
- **CI/CD integration**: Automated pipelines can validate tool versions
- **Debugging**: Developers can quickly identify the build being used (commit, build date, platform)

## What Changes

- Add `version` subcommand to the CLI that displays version information
- Support build-time version injection via GoReleaser's ldflags mechanism
- Display comprehensive build information including:
  - Version number (from git tags)
  - Git commit hash (short)
  - Build date (ISO 8601 format)
  - Go version used for compilation
  - OS and architecture
- Support `--short` flag to output only the version number (for scripting)
- Support `--json` flag for machine-readable output
- Integrate with GoReleaser to automatically inject version metadata during builds

## Impact

- **Affected specs**: `cli-interface` (new version command requirement)
- **Affected code**:
  - `cmd/root.go`: Add Version field to CLI struct
  - `cmd/version.go`: New file implementing VersionCmd
  - `main.go`: Add version variables for ldflags injection
  - `.goreleaser.yaml`: Add ldflags configuration for version metadata
- **Breaking changes**: None
- **Dependencies**: No new dependencies required
- **Backward compatibility**: Fully backward compatible, purely additive change
