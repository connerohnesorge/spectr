# Delta Specification

## MODIFIED Requirements

### Requirement: Windsurf Provider Configuration

The provider SHALL be configured with these settings:

- ID: `windsurf`
- Name: `Windsurf`
- Priority: 11
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Windsurf provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `windsurf`, Name `Windsurf`, Priority 11
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.windsurf/commands/spectr/`
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash
  commands
- **AND** it SHALL NOT return a `ConfigFileInitializer` (Windsurf has no
  instruction file)

#### Scenario: Provider priority ordering

- **WHEN** providers are registered
- **THEN** Windsurf SHALL have priority 11
- **AND** it SHALL be listed after Aider (priority 10)
- **AND** it SHALL be listed before Kilocode (priority 12)

### Requirement: Windsurf Slash Commands

The provider SHALL create slash commands in `.windsurf/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.windsurf/commands/spectr/`
  directory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.windsurf/commands/spectr/proposal.md`
- **AND** it SHALL create `.windsurf/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with YAML frontmatter
- **AND** frontmatter SHALL include a `description` field
