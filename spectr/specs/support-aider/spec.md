# Support Aider Specification

## Purpose

Documents the Aider provider integration for Spectr.

## Requirements

### Requirement: Aider Provider Configuration

The provider SHALL be configured with these settings:

- ID: `aider`
- Name: `Aider`
- Priority: 11
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider identification

- **WHEN** the registry queries for Aider provider
- **THEN** it SHALL return provider with ID `aider`

#### Scenario: Provider metadata

- **WHEN** provider metadata is accessed
- **THEN** it SHALL have name `Aider`
- **AND** it SHALL have priority 11
- **AND** it SHALL use Markdown command format

### Requirement: No Instruction File

The Aider provider SHALL NOT create an instruction file.

#### Scenario: Config file check

- **WHEN** `HasConfigFile()` is called on Aider provider
- **THEN** it returns false

### Requirement: Aider Slash Commands

The provider SHALL create slash commands in `.aider/commands/spectr/` directory.

#### Scenario: Command paths

- **WHEN** the provider configures slash commands
- **THEN** it creates `.aider/commands/spectr/proposal.md`
- **AND** it creates `.aider/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field

### Requirement: Standard Frontmatter

The provider SHALL use standard frontmatter templates for all slash commands.

#### Scenario: Proposal command frontmatter

- **WHEN** proposal.md is created
- **THEN** it SHALL contain frontmatter with description: "Scaffold a new Spectr
  change and validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** apply.md is created
- **THEN** it SHALL contain frontmatter with description: "Implement an approved
  Spectr change and keep tasks in sync."

#### Scenario: Sync command frontmatter

- **WHEN** sync.md is created
- **THEN** it SHALL contain frontmatter with description: "Detect spec drift
  from code and update specs interactively."
