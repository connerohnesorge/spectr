# Support Tabnine Specification

## Purpose
Documents the Tabnine provider integration for Spectr.

## Requirements

### Requirement: Tabnine Provider Configuration
The provider SHALL be configured with these settings:
- ID: `tabnine`
- Name: `Tabnine`
- Priority: 12
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Tabnine provider
- **THEN** it SHALL return provider with ID `tabnine`

### Requirement: No Instruction File
The Tabnine provider SHALL NOT create an instruction file.

#### Scenario: Config file check
- **WHEN** `HasConfigFile()` is called on Tabnine provider
- **THEN** it returns false

### Requirement: Tabnine Slash Commands
The provider SHALL create slash commands in `.tabnine/commands/spectr/` directory.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it creates `.tabnine/commands/spectr/proposal.md`
- **AND** it creates `.tabnine/commands/spectr/sync.md`
- **AND** it creates `.tabnine/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
