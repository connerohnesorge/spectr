# Support Antigravity Specification

## Purpose
Documents the Antigravity provider integration for Spectr.

## Requirements

### Requirement: Antigravity Provider Configuration
The provider SHALL be configured with these settings:
- ID: `antigravity`
- Name: `Antigravity`
- Priority: 7
- Config File: `AGENTS.md`
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Antigravity provider
- **THEN** it SHALL return provider with ID `antigravity`
- **AND** priority SHALL be 7

#### Scenario: Configuration file location
- **WHEN** the provider is initialized
- **THEN** config file SHALL be `AGENTS.md`
- **AND** command format SHALL be Markdown

### Requirement: Antigravity Instruction File
The provider SHALL create and maintain an `AGENTS.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with Antigravity provider selected
- **THEN** the system creates `AGENTS.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Instruction file updates
- **WHEN** `spectr init` runs for Antigravity provider
- **THEN** the system updates content between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers
- **AND** preserves content outside the markers

### Requirement: Antigravity Slash Commands
The provider SHALL create slash commands in `.agent/workflows/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider configures slash commands
- **THEN** it uses `.agent/workflows/` as base directory (not `.agent/commands/`)
- **AND** all Spectr commands reside in `.agent/workflows/` subdirectory

#### Scenario: Command file paths
- **WHEN** the provider creates slash command files
- **THEN** it creates `.agent/workflows/spectr-proposal.md`
- **AND** it creates `.agent/workflows/spectr-sync.md`
- **AND** it creates `.agent/workflows/spectr-apply.md`

#### Scenario: Command file format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter at the top
- **AND** frontmatter includes a `description` field

### Requirement: Standard Frontmatter
The provider SHALL use standard frontmatter templates for each command type.

#### Scenario: Proposal command frontmatter
- **WHEN** creating the proposal command file
- **THEN** it SHALL include frontmatter:
  ```yaml
  ---
  description: Scaffold a new Spectr change and validate strictly.
  ---
  ```

#### Scenario: Apply command frontmatter
- **WHEN** creating the apply command file
- **THEN** it SHALL include frontmatter:
  ```yaml
  ---
  description: Implement an approved Spectr change and keep tasks in sync.
  ---
  ```

#### Scenario: Sync command frontmatter
- **WHEN** creating the sync command file
- **THEN** it SHALL include frontmatter:
  ```yaml
  ---
  description: Detect spec drift from code and update specs interactively.
  ---
  ```
