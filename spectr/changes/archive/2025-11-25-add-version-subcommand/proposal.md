# Change: Add Version Subcommand

## Why

Users need to verify which version of spectr they have installed, especially
when troubleshooting issues or checking compatibility. The version information
must work correctly whether installed via goreleaser releases (which inject
version at build time via ldflags) or via nix flake (which uses the flake.nix
version attribute).

## What Changes

- Add `spectr version` subcommand to display version, commit, and build date
- Create `internal/version` package with version variables that can be set via
  ldflags
- Update `.goreleaser.yaml` to inject version information at build time
- Update `flake.nix` to inject version information via ldflags
- Support `--json` flag for machine-readable version output
- Support `--short` flag for version number only output

## Impact

- Affected specs: `specs/cli-framework/spec.md`
- Affected code:
  - `cmd/root.go` - add VersionCmd field
  - `cmd/version.go` - new file for version command
  - `internal/version/version.go` - new package for version variables
  - `.goreleaser.yaml` - add ldflags configuration
  - `flake.nix` - add ldflags to buildGoModule
