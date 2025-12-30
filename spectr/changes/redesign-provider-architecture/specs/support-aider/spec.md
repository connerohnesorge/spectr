# Delta Specification

## MODIFIED Requirements

### Requirement: Aider Provider Configuration

The provider SHALL be configured with these settings:

- ID: `aider`
- Name: `Aider`
- Priority: 10
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Aider provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `aider`, Name `Aider`, Priority 10
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.aider/commands/spectr/`
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash
  commands
- **AND** it SHALL NOT return a `ConfigFileInitializer` (Aider has no
  instruction file)

### Requirement: Aider Slash Commands

The provider SHALL create slash commands in `.aider/commands/spectr/` directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.aider/commands/spectr/` directory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.aider/commands/spectr/proposal.md`
- **AND** it SHALL create `.aider/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with YAML frontmatter
- **AND** frontmatter SHALL include `description` field
