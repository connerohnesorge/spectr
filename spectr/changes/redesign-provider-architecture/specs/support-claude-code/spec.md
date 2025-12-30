# Delta Specification

## MODIFIED Requirements

### Requirement: Claude Code Provider Configuration

The provider SHALL be configured with these settings:

- ID: `claude-code`
- Name: `Claude Code`
- Priority: 1 (highest)
- Config File: `CLAUDE.md`
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Claude Code provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `claude-code`, Name `Claude Code`,
  Priority 1
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.claude/commands/spectr/`
- **AND** it SHALL return a `ConfigFileInitializer` for `CLAUDE.md` with
  TemplateRef from TemplateManager
- **AND** it SHALL return a `SlashCommandsInitializer` for slash commands in
  Markdown format

#### Scenario: Provider metadata

- **WHEN** the provider is registered
- **THEN** the provider name SHALL be "Claude Code"
- **AND** the provider priority SHALL be 1 (highest priority)
- **AND** the provider ID SHALL be "claude-code"

### Requirement: Claude Code Instruction File

The provider SHALL create and maintain a `CLAUDE.md` instruction file in the
project root.

#### Scenario: Instruction file creation

- **WHEN** `spectr init` runs with Claude Code provider selected
- **THEN** the ConfigFileInitializer creates `CLAUDE.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:start -->` and `<!--
  spectr:end -->` markers

#### Scenario: Instruction file updates

- **WHEN** `spectr init` runs in a project with Claude Code provider
- **THEN** the ConfigFileInitializer updates content between markers in
  `CLAUDE.md`
- **AND** preserves any user content outside the markers

### Requirement: Claude Code Slash Commands

The provider SHALL create slash commands in `.claude/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.claude/commands/spectr/`
  directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.claude/commands/spectr/proposal.md`
- **AND** it SHALL create `.claude/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** each file SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field
