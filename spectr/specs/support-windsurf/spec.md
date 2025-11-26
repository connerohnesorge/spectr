# Support Windsurf Specification

## Purpose
Documents the Windsurf provider integration for Spectr.

## Requirements

### Requirement: Windsurf Provider Configuration
The provider SHALL be configured with these settings:
- ID: `windsurf`
- Name: `Windsurf`
- Priority: 13
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Windsurf provider
- **THEN** it SHALL return provider with ID `windsurf`

#### Scenario: Provider priority ordering
- **WHEN** providers are listed in priority order
- **THEN** Windsurf SHALL appear with priority 13
- **AND** it SHALL be listed after Tabnine (priority 12)
- **AND** it SHALL be listed before Kilocode (priority 14)

### Requirement: No Instruction File
The Windsurf provider SHALL NOT create an instruction file.

#### Scenario: Config file check
- **WHEN** `HasConfigFile()` is called on Windsurf provider
- **THEN** it SHALL return false

#### Scenario: Config file path
- **WHEN** the provider is initialized
- **THEN** the `configFile` field SHALL be an empty string

### Requirement: Windsurf Slash Commands
The provider SHALL create slash commands in `.windsurf/commands/spectr/` directory.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it SHALL create `.windsurf/commands/spectr/proposal.md`
- **AND** it SHALL create `.windsurf/commands/spectr/sync.md`
- **AND** it SHALL create `.windsurf/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with YAML frontmatter
- **AND** frontmatter SHALL include a `description` field

#### Scenario: Proposal command frontmatter
- **WHEN** the proposal command is created
- **THEN** the frontmatter description SHALL be "Scaffold a new Spectr change and validate strictly."

#### Scenario: Sync command frontmatter
- **WHEN** the sync command is created
- **THEN** the frontmatter description SHALL be "Detect spec drift from code and update specs interactively."

#### Scenario: Apply command frontmatter
- **WHEN** the apply command is created
- **THEN** the frontmatter description SHALL be "Implement an approved Spectr change and keep tasks in sync."
