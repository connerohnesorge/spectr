# Change: Improve README Installation Documentation

## Why

The README's Installation section is missing the primary installation method documented in `docs/src/content/docs/getting-started/installation.md`. Users arriving at the GitHub repository should find the easiest installation path (pre-built binaries from GitHub Releases) first, but currently the README only shows Nix Flakes and building from source.

The official documentation describes GitHub Releases as "the easiest way to install Spectr" but this option is not mentioned in the README at all. This creates friction for users who don't use Nix and don't want to build from source.

## What Changes

1. **Add GitHub Releases section** as the primary (first) installation method in the README
   - Include curl commands for Linux, macOS, and Windows
   - List all available platforms (Linux x86_64/arm64, macOS Intel/Apple Silicon, Windows x86_64/arm64)

2. **Reorder installation methods** to prioritize ease of use:
   - GitHub Releases (easiest, first)
   - Nix Flakes (existing, second)
   - Building from Source (existing, last)

3. **Update Table of Contents** to reflect the new installation subsections

## Impact

- **Affected files**: `README.md`
- **Affected specs**: `documentation` (alignment with existing requirement for installation instructions)
- **User impact**: Improved onboarding experience for new users who prefer binary downloads
- **Risk**: Low - additive documentation change only

## Out of Scope

- NixOS Configuration and Development Shell examples (these are appropriate for the detailed docs site, not the README)
- Changes to the actual installation process or release artifacts
