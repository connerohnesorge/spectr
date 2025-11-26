# Support Qoder Specification

## Purpose
Documents the Qoder provider integration for Spectr. Qoder is an AI coding assistant that uses `QODER.md` for configuration and `.qoder/commands/` for slash commands.

## Requirements

### Requirement: Qoder Provider Configuration
The provider SHALL be configured with these settings:
- ID: `qoder`
- Name: `Qoder`
- Priority: 4
- Config File: `QODER.md`
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Qoder provider
- **THEN** it SHALL return provider with ID `qoder`
- **AND** priority SHALL be 4

#### Scenario: Provider metadata
- **WHEN** the provider is queried for metadata
- **THEN** name SHALL be `Qoder`
- **AND** config file SHALL be `QODER.md`
- **AND** command format SHALL be Markdown

### Requirement: Qoder Instruction File
The provider SHALL create and maintain a `QODER.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with Qoder provider selected
- **THEN** the system creates `QODER.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Instruction file update
- **WHEN** `spectr update` runs
- **THEN** the system updates the Spectr instructions block in `QODER.md`
- **AND** preserves existing content outside the markers

### Requirement: Qoder Slash Commands
The provider SHALL create slash commands in `.qoder/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider configures slash commands
- **THEN** it creates `.qoder/commands/spectr/` directory
- **AND** all Spectr commands are placed under this directory

#### Scenario: Standard command paths
- **WHEN** the provider generates command file paths
- **THEN** it creates `.qoder/commands/spectr/proposal.md`
- **AND** it creates `.qoder/commands/spectr/sync.md`
- **AND** it creates `.qoder/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter at the top
- **AND** frontmatter includes a `description` field

#### Scenario: Proposal command frontmatter
- **WHEN** the proposal command file is created
- **THEN** frontmatter description SHALL be "Scaffold a new Spectr change and validate strictly."

#### Scenario: Apply command frontmatter
- **WHEN** the apply command file is created
- **THEN** frontmatter description SHALL be "Implement an approved Spectr change and keep tasks in sync."

#### Scenario: Sync command frontmatter
- **WHEN** the sync command file is created
- **THEN** frontmatter description SHALL be "Detect spec drift from code and update specs interactively."
