# Support Continue Specification

## Purpose

Documents the Continue provider integration for Spectr.

## Requirements

### Requirement: Continue Provider Configuration

The provider SHALL be configured with these settings:

- ID: `continue`
- Name: `Continue`
- Priority: 15
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider identification

- **WHEN** the registry queries for Continue provider
- **THEN** it SHALL return provider with ID `continue`

#### Scenario: Provider metadata

- **WHEN** the provider is queried for metadata
- **THEN** it SHALL return name `Continue`
- **AND** priority SHALL be 15

### Requirement: No Instruction File

The Continue provider SHALL NOT create an instruction file.

#### Scenario: Config file check

- **WHEN** `HasConfigFile()` is called on Continue provider
- **THEN** it SHALL return false

#### Scenario: Config file path

- **WHEN** the provider is queried for config file path
- **THEN** it SHALL return an empty string

### Requirement: Continue Slash Commands

The provider SHALL create slash commands in `.continue/commands/spectr/` directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** commands SHALL be placed in `.continue/commands/spectr/` directory

#### Scenario: Command paths

- **WHEN** the provider creates slash command files
- **THEN** it SHALL create `.continue/commands/spectr/proposal.md`
- **AND** it SHALL create `.continue/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** files SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field

### Requirement: Standard Frontmatter

The provider SHALL use standard frontmatter for all slash commands.

#### Scenario: Proposal command frontmatter

- **WHEN** proposal command is created
- **THEN** frontmatter SHALL include description: "Scaffold a new Spectr change and validate strictly."

#### Scenario: Sync command frontmatter

- **WHEN** sync command is created
- **THEN** frontmatter SHALL include description: "Detect spec drift from code and update specs interactively."

#### Scenario: Apply command frontmatter

- **WHEN** apply command is created
- **THEN** frontmatter SHALL include description: "Implement an approved Spectr change and keep tasks in sync."
