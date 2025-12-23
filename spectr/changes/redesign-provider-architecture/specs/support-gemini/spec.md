## MODIFIED Requirements

### Requirement: Gemini Provider Configuration
The provider SHALL be configured with these settings:
- ID: `gemini`
- Name: `Gemini CLI`
- Priority: 2
- Config File: (none - Gemini has no instruction file)
- Command Format: TOML

#### Scenario: Provider registration
- **WHEN** the Gemini provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `gemini`, Name `Gemini CLI`, Priority 2
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's Initializers() method is called
- **THEN** it SHALL return a DirectoryInitializer for `.gemini/commands/spectr/`
- **AND** it SHALL return a SlashCommandsInitializer for TOML format slash commands
- **AND** it SHALL NOT return a ConfigFileInitializer (Gemini has no instruction file)

### Requirement: Gemini Slash Commands
The provider SHALL create slash commands in `.gemini/commands/spectr/` directory using TOML format.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.gemini/commands/spectr/` directory

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it creates `.gemini/commands/spectr/proposal.toml`
- **AND** it creates `.gemini/commands/spectr/apply.toml`

#### Scenario: TOML command format
- **WHEN** slash command files are created by SlashCommandsInitializer
- **THEN** they use TOML format with `.toml` extension
- **AND** include `description` field with command description
- **AND** include `prompt` field with command content

