# Support Cline Specification

## Purpose
Documents the Cline provider integration for Spectr.

## Requirements

### Requirement: Cline Provider Configuration
The provider SHALL be configured with these settings:
- ID: `cline`
- Name: `Cline`
- Priority: 8
- Config File: `CLINE.md`
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Cline provider
- **THEN** it SHALL return provider with ID `cline`

### Requirement: Cline Instruction File
The provider SHALL create and maintain a `CLINE.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with Cline provider selected
- **THEN** the system creates `CLINE.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

### Requirement: Cline Slash Commands
The provider SHALL create slash commands in `.clinerules/commands/spectr/` directory.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it creates `.clinerules/commands/spectr/proposal.md`
- **AND** it creates `.clinerules/commands/spectr/sync.md`
- **AND** it creates `.clinerules/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
