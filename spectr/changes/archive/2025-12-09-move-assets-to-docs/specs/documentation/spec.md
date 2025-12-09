## MODIFIED Requirements

### Requirement: Reproducible Demo Source Files
The system SHALL maintain VHS tape files as version-controlled source for all demo GIFs to enable easy regeneration when the CLI changes.

#### Scenario: Developer regenerates outdated GIF
- **WHEN** a developer updates a CLI command
- **THEN** they SHALL be able to run the corresponding VHS tape file to regenerate an accurate GIF

#### Scenario: Developer creates new demo
- **WHEN** a developer wants to add a new demo
- **THEN** they SHALL find existing tape files as examples in `docs/src/assets/vhs/` directory

#### Scenario: Contributor finds demo standards
- **WHEN** a contributor reads the development documentation
- **THEN** they SHALL find guidelines for VHS tape configuration (theme, size, typing speed)

### Requirement: Demo Asset Organization
The system SHALL organize demo assets with clear separation between source files (VHS tapes) and generated outputs (GIFs).

#### Scenario: Developer locates tape source
- **WHEN** a developer needs to update a demo
- **THEN** they SHALL find VHS tape files in `docs/src/assets/vhs/` directory

#### Scenario: Documentation references generated GIF
- **WHEN** the README or docs site needs to embed a demo
- **THEN** they SHALL reference GIF files from `docs/src/assets/gifs/` directory

#### Scenario: Developer regenerates all demos
- **WHEN** a developer runs the regeneration command
- **THEN** all GIFs SHALL be generated from their corresponding tape files and placed in `docs/src/assets/gifs/`

### Requirement: Batch GIF Generation Command
The system SHALL provide a `generate-gif` command (via Nix flake) to generate all demo GIFs in one command, supporting both full regeneration and single-demo regeneration.

#### Scenario: Developer regenerates all GIFs
- **WHEN** a developer runs `generate-gif` in the nix develop shell
- **THEN** all VHS tape files SHALL be processed
- **AND** GIFs SHALL be output to `docs/src/assets/gifs/` directory

#### Scenario: Developer regenerates single GIF
- **WHEN** a developer runs `generate-gif <demo-name>`
- **THEN** only the specified demo's GIF SHALL be regenerated

#### Scenario: Developer gets command usage help
- **WHEN** a developer runs `generate-gif --help`
- **THEN** they SHALL see available demo names and usage instructions
