# Nix Packaging Specification

## ADDED Requirements

### Requirement: CLI Binary Build via Nix
The system SHALL build the spectr CLI binary using Nix flakes with a `packages.default` configuration that uses `buildGoModule` to compile the Go source code into an executable binary.

#### Scenario: Build default package
- **WHEN** user runs `nix build` in the project root
- **THEN** a spectr CLI binary is produced at `result/bin/spectr`

#### Scenario: Build with specific Go version
- **WHEN** the flake.nix specifies Go 1.25.0 as the compiler
- **THEN** the built binary uses Go 1.25.0 runtime

### Requirement: Vendor Hash Configuration
The system SHALL specify a `vendorHash` in `packages.default` that ensures reproducible builds by pinning Go module dependencies to a known state.

#### Scenario: Reproducible builds
- **WHEN** `nix build` is executed multiple times with the same source
- **THEN** the output binary hash remains identical

### Requirement: Package Metadata
The system SHALL include package metadata (pname, version, description, homepage, license, maintainers) in the flake outputs for distribution and discoverability.

#### Scenario: Package metadata presence
- **WHEN** the flake is evaluated
- **THEN** pname="spectr", version follows semantic versioning, and license is Apache 2.0

#### Scenario: Homepage and license information
- **WHEN** the package is published to Nixpkgs or documentation systems
- **THEN** homepage points to the authoritative repository and license is explicitly stated

### Requirement: Development Shell Integration
The system SHALL provide a development shell via `devShells.default` that includes Go toolchain, linting, testing, and formatting tools required for spectr development.

#### Scenario: Enter development environment
- **WHEN** developer runs `nix develop`
- **THEN** all required build and development tools are available in PATH (Go, golangci-lint, gotestsum, etc.)

#### Scenario: Use live reload during development
- **WHEN** developer is in the development shell
- **THEN** `air` command is available for live reloading code changes

### Requirement: Source Code Inclusion
The system SHALL ensure the flake correctly specifies the project source (`src = self`) so that all Go source files, go.mod, and go.sum are included in the build context.

#### Scenario: All source files included
- **WHEN** `nix build` executes
- **THEN** all .go files and module metadata are available to the Go compiler

### Requirement: Output Structure
The system SHALL produce a standard Nix package output with the binary executable in the expected location within the derivation.

#### Scenario: Binary in correct location
- **WHEN** build completes successfully
- **THEN** the spectr binary is located at `$out/bin/spectr` and is executable

### Requirement: Cross-Platform Support
The system SHALL support building on multiple platforms (x86_64-linux, aarch64-linux, x86_64-darwin, aarch64-darwin) through proper Nix flake configuration.

#### Scenario: Build on supported platforms
- **WHEN** the flake is evaluated on aarch64-darwin (Apple Silicon)
- **THEN** the build produces a native aarch64-darwin binary
