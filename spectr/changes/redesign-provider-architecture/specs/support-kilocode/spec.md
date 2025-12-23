## MODIFIED Requirements

### Requirement: Kilocode Provider Configuration
The provider SHALL be configured with these settings:
- ID: `kilocode`
- Name: `Kilocode`
- Priority: 14
- Config File: (none)
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Kilocode provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `kilocode`, Name `Kilocode`, Priority 14
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's Initializers() method is called
- **THEN** it SHALL return a DirectoryInitializer for `.kilocode/commands/spectr/`
- **AND** it SHALL return a SlashCommandsInitializer for Markdown format slash commands
- **AND** it SHALL NOT return a ConfigFileInitializer (Kilocode has no instruction file)

#### Scenario: Provider priority
- **WHEN** providers are sorted by priority
- **THEN** Kilocode SHALL have priority 14

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
- **THEN** it creates `.kilocode/commands/spectr/proposal.md`
- **AND** it creates `.kilocode/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field

