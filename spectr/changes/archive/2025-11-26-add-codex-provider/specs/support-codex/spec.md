# Support Codex Specification

## Purpose
Documents the Codex CLI provider integration for Spectr, enabling Spectr to work with OpenAI's Codex CLI through global prompt file management and slash command generation.

## ADDED Requirements

### Requirement: Codex Provider Configuration
The provider SHALL be configured with these settings:
- ID: `codex`
- Name: `Codex CLI`
- Priority: 10
- Config File: `AGENTS.md`
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for Codex provider
- **THEN** it SHALL return provider with ID `codex`
- **AND** the provider priority is 10

#### Scenario: Provider metadata
- **WHEN** displaying provider options to users
- **THEN** the provider name is "Codex CLI"
- **AND** it appears after Cursor (priority 9) and before Aider (priority 11)

#### Scenario: Instruction file
- **WHEN** checking if Codex provider has a config file
- **THEN** HasConfigFile() returns true
- **AND** ConfigFile() returns "AGENTS.md"

### Requirement: Codex Global Slash Commands
The provider SHALL create slash commands in the global `~/.codex/prompts/` directory.

#### Scenario: Global command directory structure
- **WHEN** the provider configures slash commands
- **THEN** it creates `~/.codex/prompts/spectr/` directory
- **AND** the directory is created in user's home directory, not project directory

#### Scenario: Command paths
- **WHEN** the provider generates slash command files
- **THEN** it creates `~/.codex/prompts/spectr-proposal.md`
- **AND** it creates `~/.codex/prompts/spectr-sync.md`
- **AND** it creates `~/.codex/prompts/spectr-apply.md`

#### Scenario: Global path expansion
- **WHEN** resolving command paths
- **THEN** the `~` prefix is expanded to user's home directory
- **AND** paths work correctly regardless of current project directory

### Requirement: Codex Command Format
The provider SHALL use Markdown format with YAML frontmatter for slash commands.

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

#### Scenario: Proposal command frontmatter
- **WHEN** generating the proposal command file
- **THEN** the frontmatter description is "Scaffold a new Spectr change and validate strictly."

#### Scenario: Apply command frontmatter
- **WHEN** generating the apply command file
- **THEN** the frontmatter description is "Implement an approved Spectr change and keep tasks in sync."

#### Scenario: Sync command frontmatter
- **WHEN** generating the sync command file
- **THEN** the frontmatter description is "Detect spec drift from code and update specs interactively."

### Requirement: Global Path Support in Provider Framework
The provider framework SHALL support global paths (starting with `~/` or `/`) in addition to project-relative paths.

#### Scenario: Global path detection
- **WHEN** a command path starts with `~/` or `/`
- **THEN** the system treats it as a global path
- **AND** does not prepend the project directory

#### Scenario: Home directory expansion
- **WHEN** a path starts with `~/`
- **THEN** the system expands `~` to the user's home directory
- **AND** uses `os.UserHomeDir()` for cross-platform compatibility

#### Scenario: IsConfigured with global paths
- **WHEN** checking if a provider with global paths is configured
- **THEN** the system checks the expanded absolute path
- **AND** does not look in the project directory

### Requirement: Codex Command Invocation
Users SHALL invoke Spectr commands in Codex using the `/spectr-<command>` pattern.

#### Scenario: Invoking proposal command
- **WHEN** user types `/spectr-proposal` in Codex
- **THEN** Codex loads and executes the proposal prompt

#### Scenario: Invoking sync command
- **WHEN** user types `/spectr-sync` in Codex
- **THEN** Codex loads and executes the sync prompt

#### Scenario: Invoking apply command
- **WHEN** user types `/spectr-apply` in Codex
- **THEN** Codex loads and executes the apply prompt
