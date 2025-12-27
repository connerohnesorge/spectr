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

#### Scenario: Provider returns initializers with home paths
- **WHEN** the provider's `Initializers(ctx context.Context, tm *TemplateManager)` method is called
- **THEN** it SHALL return a `HomeDirectoryInitializer` for `.codex/prompts/` (relative to home directory)
- **AND** it SHALL return a `ConfigFileInitializer` for `AGENTS.md` with TemplateRef from TemplateManager using `tm.Agents()`
- **AND** it SHALL return a `HomePrefixedSlashCommandsInitializer` with prefix `spectr-` for home slash commands in `.codex/prompts/`

#### Scenario: Provider metadata
- **WHEN** provider is registered
- **THEN** the provider name is "Codex CLI"
- **AND** it appears after Cursor (priority 8) and before Aider (priority 10)

#### Scenario: Instruction file
- **WHEN** the provider returns initializers
- **THEN** it includes a ConfigFileInitializer for "AGENTS.md"

#### Scenario: Create new instruction file
- **WHEN** `AGENTS.md` does not exist
- **THEN** the ConfigFileInitializer SHALL create it with instruction content between markers
- **AND** the markers SHALL be `<!-- spectr:start -->` and `<!-- spectr:end -->` (lowercase)

#### Scenario: Update existing instruction file
- **WHEN** `AGENTS.md` exists with spectr markers
- **THEN** the ConfigFileInitializer SHALL replace content between markers
- **AND** it SHALL preserve content outside markers
- **AND** the marker search SHALL be case-insensitive (matches both uppercase and lowercase)
- **AND** when writing, the system SHALL always use lowercase markers

### Requirement: Codex Global Slash Commands
The provider SHALL create slash commands in the home `~/.codex/prompts/` directory.

#### Scenario: Home command directory structure
- **WHEN** the provider returns initializers
- **THEN** `HomeDirectoryInitializer` SHALL create `~/.codex/prompts/` directory
- **AND** the directory is created in user's home directory via homeFs

#### Scenario: Command paths with prefix
- **WHEN** the `HomePrefixedSlashCommandsInitializer` executes
- **THEN** it SHALL create `.codex/prompts/spectr-proposal.md` in the home filesystem
- **AND** it SHALL create `.codex/prompts/spectr-apply.md` in the home filesystem

#### Scenario: Home path handling
- **WHEN** Home* initializers execute
- **THEN** the executor provides both projectFs and homeFs filesystems
- **AND** the Home* initializers use homeFs automatically (no configuration flag needed)
- **AND** paths work correctly regardless of current project directory

### Requirement: Codex Command Format
The provider SHALL use Markdown format with YAML frontmatter for slash commands.

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** each file SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field

