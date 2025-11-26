# Support CoStrict Specification

## Purpose
Documents the CoStrict provider integration for Spectr.

## Requirements

### Requirement: CoStrict Provider Configuration
The provider SHALL be configured with these settings:
- ID: `costrict`
- Name: `CoStrict`
- Priority: 3
- Config File: `COSTRICT.md`
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for CoStrict provider
- **THEN** it SHALL return provider with ID `costrict`

### Requirement: CoStrict Instruction File
The provider SHALL create and maintain a `COSTRICT.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with CoStrict provider selected
- **THEN** the system creates `COSTRICT.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

### Requirement: CoStrict Slash Commands
The provider SHALL create slash commands in `.costrict/commands/spectr/` directory.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it creates `.costrict/commands/spectr/proposal.md`
- **AND** it creates `.costrict/commands/spectr/sync.md`
- **AND** it creates `.costrict/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
