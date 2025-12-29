## ADDED Requirements

### Requirement: Pre-made Example Projects for VHS Demos

The system SHALL provide pre-made spectr project examples in the `examples/` directory that VHS tape files use for demonstrations, keeping demos focused on spectr commands rather than file creation boilerplate.

#### Scenario: Developer creates clean demo

- **WHEN** a VHS tape file needs a spectr project for demonstration
- **THEN** it SHALL copy from a pre-made example in `examples/` directory
- **AND** the demo output SHALL focus on spectr commands, not `cat` heredocs creating files

#### Scenario: Developer maintains example project

- **WHEN** a change to demo content is needed
- **THEN** the developer SHALL edit the pre-made example in `examples/` directory
- **AND** the change will automatically apply to any tape using that example

### Requirement: Batch GIF Generation Command

The system SHALL provide a `generate-gif` command (via Nix flake) to generate all demo GIFs in one command, supporting both full regeneration and single-demo regeneration.

#### Scenario: Developer regenerates all GIFs

- **WHEN** a developer runs `generate-gif` in the nix develop shell
- **THEN** all VHS tape files SHALL be processed
- **AND** GIFs SHALL be output to `assets/gifs/` directory

#### Scenario: Developer regenerates single GIF

- **WHEN** a developer runs `generate-gif <demo-name>`
- **THEN** only the specified demo's GIF SHALL be regenerated

#### Scenario: Developer gets command usage help

- **WHEN** a developer runs `generate-gif --help`
- **THEN** they SHALL see available demo names and usage instructions
