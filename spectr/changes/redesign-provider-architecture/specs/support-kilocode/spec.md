# MODIFIED Requirements
## MODIFIED Requirements

### Requirement: Kilocode Provider Configuration

The provider SHALL be configured with these settings:

- ID: `kilocode`
- Name: `Kilocode`
- Priority: 12
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Kilocode provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `kilocode`, Name `Kilocode`, Priority 12
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for `.kilocode/commands/spectr/`
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash commands
- **AND** it SHALL NOT return a `ConfigFileInitializer` (Kilocode has no instruction file)

#### Scenario: Provider priority

- **WHEN** providers are sorted by priority
- **THEN** Kilocode SHALL have priority 12

#### Scenario: Command format check

- **WHEN** the provider is registered
- **THEN** it SHALL use Markdown format for slash commands

### Requirement: Kilocode Slash Commands

The provider SHALL create slash commands in `.kilocode/commands/spectr/` directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.kilocode/commands/spectr/` subdirectory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.kilocode/commands/spectr/proposal.md`
- **AND** it SHALL create `.kilocode/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with YAML frontmatter
- **AND** frontmatter SHALL include `description` field
