# Change: Add Version Subcommand

## Why

Users and automation systems need a standard way to determine which version of Spectr is installed for troubleshooting, compatibility verification, and dependency management. Currently, there is no built-in command to display version information, forcing users to rely on external methods like checking installation artifacts or git tags.

A version command provides:
- **Bug reporting**: Users can accurately report which version exhibits issues
- **Compatibility checking**: Scripts can verify minimum version requirements
- **Feature availability**: Users can determine if their version supports specific features
- **CI/CD integration**: Automated pipelines can validate tool versions
- **Debugging**: Developers can quickly identify the build being used (commit, build date, platform)

The file-based approach using Go's `embed` package is chosen because it works seamlessly with both Nix builds and standard Go builds. Nix's hermetic build environment doesn't easily support dynamic ldflags injection, while the embed approach allows Nix to generate a VERSION file during the build phase that gets embedded into the binary. This provides a unified solution that works across all build systems without special-casing.

## What Changes

- Add `version` subcommand to the CLI that displays version information
- Create a `VERSION` file at the project root containing version metadata in JSON format with fields:
  - `version`: Version number (from git tags)
  - `commit`: Git commit hash (short)
  - `date`: Build date (ISO 8601 format)
- Use Go's `embed` package to embed the VERSION file into the binary at compile time
- Read and parse version information from the embedded file at runtime
- Display comprehensive build information including:
  - Version number (from VERSION file)
  - Git commit hash (from VERSION file)
  - Build date (from VERSION file)
  - Go version used for compilation (from runtime)
  - OS and architecture (from runtime)
- Support `--short` flag to output only the version number (for scripting)
- Support `--json` flag for machine-readable output
- Support both build systems:
  - **Nix**: `flake.nix` generates VERSION file during build phase
  - **GoReleaser**: Build hooks generate VERSION file before compilation

## Impact

- **Affected specs**: `cli-interface` (new version command requirement)
- **Affected code**:
  - `VERSION`: New file at project root containing version metadata (JSON format)
  - `cmd/version.go`: New file implementing VersionCmd with embed package
  - `cmd/root.go`: Add version subcommand registration
  - `.goreleaser.yaml`: Add before hooks to generate VERSION file during build
  - `flake.nix`: Add VERSION file generation during Nix build phase
- **Breaking changes**: None
- **Dependencies**: No new dependencies required (embed is standard library)
- **Backward compatibility**: Fully backward compatible, purely additive change
