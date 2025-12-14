# Support Opencode Specification

## Purpose
This document describes the requirements for the OpenCode provider.

## Requirements

### Requirement: OpenCode Provider Configuration
The provider SHALL be configured with these settings:
- ID: `opencode`
- Name: `OpenCode`
- Priority: 16 (after Continue)
- Config File: None (OpenCode uses JSON config, instruction injection not supported)
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for OpenCode provider
- **THEN** it SHALL return provider with ID `opencode`
- **AND** the provider priority is 16

#### Scenario: Provider metadata
- **WHEN** displaying provider options to users
- **THEN** the provider name is "OpenCode"
- **AND** it appears in the list ordered by priority

### Requirement: OpenCode Slash Commands
The provider SHALL create slash commands in `.opencode/command/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider configures slash commands
- **THEN** it creates `.opencode/command/spectr/` directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths
- **WHEN** the provider generates slash command files
- **THEN** it creates `.opencode/command/spectr/proposal.md`
- **AND** it creates `.opencode/command/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

### Requirement: Standard Frontmatter
The provider SHALL use standard frontmatter templates for each command type.

#### Scenario: Proposal command frontmatter
- **WHEN** generating the proposal command file
- **THEN** the frontmatter description is "Scaffold a new Spectr change and validate strictly."

#### Scenario: Apply command frontmatter
- **WHEN** generating the apply command file
- **THEN** the frontmatter description is "Implement an approved Spectr change and keep tasks in sync."

### Requirement: Command Path Construction
The provider SHALL construct command paths using the standard pattern.

#### Scenario: Path construction for OpenCode
- **WHEN** determining command file paths
- **THEN** it uses base directory `.opencode/command`
- **AND** appends `/spectr/` subdirectory
- **AND** appends command name with `.md` extension
- **AND** results in paths like `.opencode/command/spectr/proposal.md`

### Requirement: No Instruction File
The provider SHALL NOT create an instruction file since OpenCode uses JSON configuration.

#### Scenario: Config file check
- **WHEN** checking if provider has config file
- **THEN** it returns false
- **AND** no instruction file is created during initialization

