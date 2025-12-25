## MODIFIED Requirements

### Requirement: Continue Provider Configuration
The provider SHALL be configured with these settings:
- ID: `continue`
- Name: `Continue`
- Priority: 13
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Continue provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `continue`, Name `Continue`, Priority 13
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's `Initializers(ctx, tm *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for `.continue/commands/spectr/`
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash commands
- **AND** it SHALL NOT return a `ConfigFileInitializer` (Continue has no instruction file)

#### Scenario: Provider metadata
- **WHEN** the provider is registered
- **THEN** it SHALL have name `Continue`
- **AND** priority SHALL be 13

### Requirement: Continue Slash Commands
The provider SHALL create slash commands in `.continue/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.continue/commands/spectr/` directory

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.continue/commands/spectr/proposal.md`
- **AND** it SHALL create `.continue/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** files SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field

