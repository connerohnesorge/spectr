## MODIFIED Requirements

### Requirement: Codex Provider Configuration
The provider SHALL be configured with these settings:
- ID: `codex`
- Name: `Codex CLI`
- Priority: 9
- Config File: `AGENTS.md`
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Codex provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `codex`, Name `Codex CLI`, Priority 9
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers with global paths
- **WHEN** the provider's Initializers() method is called
- **THEN** it SHALL return a DirectoryInitializer for `~/.codex/prompts/` configured for global filesystem
- **AND** it SHALL return a ConfigFileInitializer for `AGENTS.md`
- **AND** it SHALL return a SlashCommandsInitializer for global slash commands configured for global filesystem

#### Scenario: Provider metadata
- **WHEN** provider is registered
- **THEN** the provider name is "Codex CLI"
- **AND** it appears after Cursor (priority 8) and before Aider (priority 10)

#### Scenario: Instruction file
- **WHEN** the provider returns initializers
- **THEN** it includes a ConfigFileInitializer for "AGENTS.md"

### Requirement: Codex Global Slash Commands
The provider SHALL create slash commands in the global `~/.codex/prompts/` directory.

#### Scenario: Global command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer configured for global filesystem SHALL create `~/.codex/prompts/` directory
- **AND** the directory is created in user's home directory via globalFs

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer configured for global filesystem executes
- **THEN** it creates `~/.codex/prompts/spectr-proposal.md`
- **AND** it creates `~/.codex/prompts/spectr-apply.md`

#### Scenario: Global path handling
- **WHEN** initializers configured for global filesystem execute
- **THEN** the executor provides both projectFs and globalFs filesystems
- **AND** the initializer uses globalFs based on its internal configuration
- **AND** paths work correctly regardless of current project directory

### Requirement: Codex Command Format
The provider SHALL use Markdown format with YAML frontmatter for slash commands.

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

