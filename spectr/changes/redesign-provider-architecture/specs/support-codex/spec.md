## MODIFIED Requirements

### Requirement: Codex Provider Configuration
The provider SHALL be configured with these settings:
- ID: `codex`
- Name: `Codex CLI`
- Priority: 10
- Config File: `AGENTS.md`
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Codex provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `codex`, Name `Codex CLI`, Priority 10
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers with global paths
- **WHEN** the provider's Initializers() method is called
- **THEN** it SHALL return a DirectoryInitializer for `~/.codex/prompts/` with IsGlobal() = true
- **AND** it SHALL return a ConfigFileInitializer for `AGENTS.md`
- **AND** it SHALL return a SlashCommandsInitializer for global slash commands with IsGlobal() = true

#### Scenario: Provider metadata
- **WHEN** provider is registered
- **THEN** the provider name is "Codex CLI"
- **AND** it appears after Cursor (priority 9) and before Aider (priority 11)

#### Scenario: Instruction file
- **WHEN** the provider returns initializers
- **THEN** it includes a ConfigFileInitializer for "AGENTS.md"

### Requirement: Codex Global Slash Commands
The provider SHALL create slash commands in the global `~/.codex/prompts/` directory.

#### Scenario: Global command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer with IsGlobal() = true SHALL create `~/.codex/prompts/spectr/` directory
- **AND** the directory is created in user's home directory via globalFs

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer with IsGlobal() = true executes
- **THEN** it creates `~/.codex/prompts/spectr-proposal.md`
- **AND** it creates `~/.codex/prompts/spectr-apply.md`

#### Scenario: Global path handling
- **WHEN** initializers with IsGlobal() = true execute
- **THEN** the executor provides the globalFs filesystem rooted at user's home directory
- **AND** paths work correctly regardless of current project directory

### Requirement: Codex Command Format
The provider SHALL use Markdown format with YAML frontmatter for slash commands.

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

