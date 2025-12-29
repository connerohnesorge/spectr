# Support Kilocode Specification

## Purpose

Documents the Kilocode provider integration for Spectr.

## Requirements

### Requirement: Kilocode Provider Configuration

The provider SHALL be configured with these settings:

- ID: `kilocode`
- Name: `Kilocode`
- Priority: 14
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider identification

- **WHEN** the registry queries for Kilocode provider
- **THEN** it SHALL return provider with ID `kilocode`

#### Scenario: Provider priority

- **WHEN** providers are sorted by priority
- **THEN** Kilocode SHALL have priority 14

#### Scenario: Command format check

- **WHEN** the provider is queried for command format
- **THEN** it SHALL return Markdown format

### Requirement: No Instruction File

The Kilocode provider SHALL NOT create an instruction file.

#### Scenario: Config file check

- **WHEN** `HasConfigFile()` is called on Kilocode provider
- **THEN** it returns false

#### Scenario: Config file path

- **WHEN** `ConfigFile()` is called on Kilocode provider
- **THEN** it returns empty string

### Requirement: Kilocode Slash Commands

The provider SHALL create slash commands in `.kilocode/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** commands are placed in `.kilocode/commands/spectr/` subdirectory

#### Scenario: Command paths

- **WHEN** the provider configures slash commands
- **THEN** it creates `.kilocode/commands/spectr/proposal.md`
- **AND** it creates `.kilocode/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field

#### Scenario: Proposal command frontmatter

- **WHEN** proposal command is created
- **THEN** frontmatter contains description "Scaffold a new Spectr change and
  validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** apply command is created
- **THEN** frontmatter contains description "Implement an approved Spectr change
  and keep tasks in sync."
