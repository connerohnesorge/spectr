# Delta: Documentation Specification

## ADDED Requirements

### Requirement: Visual CLI Demonstrations
The system SHALL provide visual demonstrations of core CLI workflows using VHS-generated GIF recordings to help users quickly understand Spectr's capabilities.

#### Scenario: User sees initialization demo
- **WHEN** a user reads the Quick Start section of the README
- **THEN** they SHALL see a GIF demonstrating the `spectr init` command and resulting directory structure

#### Scenario: User sees validation demo
- **WHEN** a user reads about validation in the documentation
- **THEN** they SHALL see a GIF showing validation errors and how to fix them

#### Scenario: User sees complete workflow demo
- **WHEN** a user visits the getting-started guide
- **THEN** they SHALL see a GIF demonstrating the complete workflow from proposal to archive

### Requirement: Reproducible Demo Source Files
The system SHALL maintain VHS tape files as version-controlled source for all demo GIFs to enable easy regeneration when the CLI changes.

#### Scenario: Developer regenerates outdated GIF
- **WHEN** a developer updates a CLI command
- **THEN** they SHALL be able to run the corresponding VHS tape file to regenerate an accurate GIF

#### Scenario: Developer creates new demo
- **WHEN** a developer wants to add a new demo
- **THEN** they SHALL find existing tape files as examples in `assets/vhs/` directory

#### Scenario: Contributor finds demo standards
- **WHEN** a contributor reads the development documentation
- **THEN** they SHALL find guidelines for VHS tape configuration (theme, size, typing speed)

### Requirement: Demo Asset Organization
The system SHALL organize demo assets with clear separation between source files (VHS tapes) and generated outputs (GIFs).

#### Scenario: Developer locates tape source
- **WHEN** a developer needs to update a demo
- **THEN** they SHALL find VHS tape files in `assets/vhs/` directory

#### Scenario: Documentation references generated GIF
- **WHEN** the README or docs site needs to embed a demo
- **THEN** they SHALL reference GIF files from `assets/gifs/` directory

#### Scenario: Developer regenerates all demos
- **WHEN** a developer runs the regeneration command
- **THEN** all GIFs SHALL be generated from their corresponding tape files and placed in `assets/gifs/`

### Requirement: Core Workflow Coverage
The system SHALL provide demo GIFs covering all essential Spectr workflows: initialization, listing, validation, and archiving.

#### Scenario: User learns initialization
- **WHEN** a user views the init demo GIF
- **THEN** they SHALL see `spectr init` being run and the resulting `spectr/` directory structure

#### Scenario: User learns listing
- **WHEN** a user views the list demo GIF
- **THEN** they SHALL see both `spectr list` (changes) and `spectr list --specs` (specifications) commands

#### Scenario: User learns validation
- **WHEN** a user views the validate demo GIF
- **THEN** they SHALL see `spectr validate` catching an error, the error being fixed, and validation passing

#### Scenario: User learns archiving
- **WHEN** a user views the archive demo GIF
- **THEN** they SHALL see `spectr archive` merging deltas into specs and moving the change to the archive directory

#### Scenario: User sees end-to-end workflow
- **WHEN** a user views the workflow demo GIF
- **THEN** they SHALL see the complete three-stage workflow from creating a change through archiving it
