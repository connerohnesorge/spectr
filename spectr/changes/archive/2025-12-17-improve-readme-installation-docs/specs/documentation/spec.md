## MODIFIED Requirements

### Requirement: Comprehensive README with Multiple Sections

The system SHALL provide a comprehensive README.md file that serves both end users and developers, including installation instructions, usage guide, command reference, architecture overview, and contribution guidelines.

#### Scenario: User finds installation instructions

- **WHEN** a new user visits the repository
- **THEN** they SHALL find clear instructions for installing via pre-built binaries from GitHub Releases, Nix Flakes, or building from source
- **AND** the GitHub Releases method SHALL be documented first as the easiest installation path
- **AND** all available platforms SHALL be listed (Linux x86_64/arm64, macOS x86_64/arm64, Windows x86_64/arm64)

#### Scenario: Developer understands architecture

- **WHEN** a developer reads the README
- **THEN** they SHALL find an architecture overview explaining the clean separation of concerns and package structure

#### Scenario: Contributor knows how to contribute

- **WHEN** someone wants to contribute
- **THEN** they SHALL find guidelines for code style, testing, commit conventions, and PR process
