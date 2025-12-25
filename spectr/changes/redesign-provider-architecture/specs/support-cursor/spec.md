## MODIFIED Requirements

### Requirement: Cursor Provider Configuration
The provider SHALL be configured with these settings:
- ID: `cursor`
- Name: `Cursor`
- Priority: 8
- Config File: (none - Cursor has no instruction file)
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Cursor provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `cursor`, Name `Cursor`, Priority 8
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's `Initializers(ctx, tm *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for `.cursorrules/commands/spectr/`
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash commands
- **AND** it SHALL NOT return a `ConfigFileInitializer` (Cursor has no instruction file)

### Requirement: Cursor Slash Commands
The provider SHALL create slash commands in `.cursorrules/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.cursorrules/commands/spectr/` directory

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it creates `.cursorrules/commands/spectr/proposal.md`
- **AND** it creates `.cursorrules/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field

