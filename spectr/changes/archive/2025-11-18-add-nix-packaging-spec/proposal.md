# Change: Add Nix Packaging Specification

## Why
Spectr currently lacks documented specifications for how the CLI is packaged and distributed via Nix flakes. The `packages.default` configuration exists but is not formally spec'd. This creates ambiguity about build requirements, distribution mechanisms, and packaging responsibilities. Formalizing this as a capability spec ensures the packaging workflow is explicit, testable, and maintainable.

## What Changes
- Introduces a new "Nix Packaging" capability in `specs/nix-packaging/`
- Documents the `packages.default` buildGoModule configuration
- Specifies requirements for building the CLI binary
- Defines distribution and release workflows
- Establishes conventions for version management

## Impact
- **Affected specs**: New capability (nix-packaging)
- **Affected code**: `flake.nix` (packages.default), main.go (version info), release workflows
- **Breaking changes**: None
- **Implementation**: No code changes required; spec-only documentation
