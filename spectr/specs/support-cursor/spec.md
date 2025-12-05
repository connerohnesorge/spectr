# Support Cursor Specification

## Purpose
Documents the Cursor provider integration for Spectr.

## Requirements

### Requirement: Cursor Provider Configuration
The provider SHALL be configured with these settings:
- ID: `cursor`
- Name: `Cursor`
- Priority: 9
- Config File: (none - Cursor has no instruction file)
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Cursor provider
- **THEN** it SHALL return provider with ID `cursor`

### Requirement: No Instruction File
The Cursor provider SHALL NOT create an instruction file.

#### Scenario: Config file check
- **WHEN** `HasConfigFile()` is called on Cursor provider
- **THEN** it returns false

### Requirement: Cursor Slash Commands
The provider SHALL create slash commands in `.cursorrules/commands/spectr/` directory.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it creates `.cursorrules/commands/spectr/proposal.md`
- **AND** it creates `.cursorrules/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
