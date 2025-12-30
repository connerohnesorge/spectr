# Support Antigravity Specification

## Purpose

Documents the Antigravity provider integration for Spectr.

## Requirements

### Requirement: Antigravity Provider Configuration

The provider SHALL be configured with these settings:

- ID: `antigravity`
- Name: `Antigravity`
- Priority: 6
- Config File: `AGENTS.md`
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Antigravity provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `antigravity`, Name `Antigravity`,
  Priority 6
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for `.agent/workflows/`
- **AND** it SHALL return a `ConfigFileInitializer` for `AGENTS.md` with
  TemplateRef from TemplateManager
- **AND** it SHALL return a `PrefixedSlashCommandsInitializer` with prefix
  `spectr-` for Markdown format slash commands in `.agent/workflows/`

#### Scenario: Configuration file location

- **WHEN** the provider is initialized
- **THEN** ConfigFileInitializer SHALL target `AGENTS.md`
- **AND** command format SHALL be Markdown

### Requirement: Antigravity Instruction File

The provider SHALL create and maintain an `AGENTS.md` instruction file in the
project root.

#### Scenario: Instruction file creation

- **WHEN** `spectr init` runs with Antigravity provider selected
- **THEN** the ConfigFileInitializer creates `AGENTS.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:start -->` and `<!--
  spectr:end -->` markers

#### Scenario: Instruction file updates

- **WHEN** `spectr init` runs for Antigravity provider
- **THEN** the ConfigFileInitializer updates content between `<!-- spectr:start
  -->` and `<!-- spectr:end -->` markers
- **AND** preserves content outside the markers

### Requirement: Antigravity Slash Commands

The provider SHALL create slash commands in `.agent/workflows/` directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.agent/workflows/` directory (not
  `.agent/commands/`)
- **AND** all Spectr commands reside in `.agent/workflows/` subdirectory

#### Scenario: Command file paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.agent/workflows/spectr-proposal.md`
- **AND** it SHALL create `.agent/workflows/spectr-apply.md`

#### Scenario: Command file format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** each file SHALL include YAML frontmatter at the top
- **AND** frontmatter SHALL include a `description` field

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
