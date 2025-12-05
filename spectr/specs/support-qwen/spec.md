# Support Qwen Specification

## Purpose
Documents the Qwen Code provider integration for Spectr.

## Requirements

### Requirement: Qwen Provider Configuration
The provider SHALL be configured with these settings:
- ID: `qwen`
- Name: `Qwen Code`
- Priority: 6
- Config File: `QWEN.md`
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Qwen provider
- **THEN** it SHALL return provider with ID `qwen`

### Requirement: Qwen Instruction File
The provider SHALL create and maintain a `QWEN.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with Qwen provider selected
- **THEN** the system creates `QWEN.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

### Requirement: Qwen Slash Commands
The provider SHALL create slash commands in `.qwen/commands/spectr/` directory.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it creates `.qwen/commands/spectr/proposal.md`
- **AND** it creates `.qwen/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
