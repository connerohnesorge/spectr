## MODIFIED Requirements

### Requirement: OpenCode Provider Configuration
The provider SHALL be configured with these settings:
- ID: `opencode`
- Name: `OpenCode`
- Priority: 15
- Config File: None (OpenCode uses JSON config, instruction injection not supported)
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the OpenCode provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `opencode`, Name `OpenCode`, Priority 15
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's Initializers() method is called
- **THEN** it SHALL return a DirectoryInitializer for `.opencode/commands/spectr/`
- **AND** it SHALL return a SlashCommandsInitializer for Markdown format slash commands
- **AND** it SHALL NOT return a ConfigFileInitializer (OpenCode uses JSON config)

#### Scenario: Provider metadata
- **WHEN** provider is registered
- **THEN** the provider name is "OpenCode"
- **AND** it appears in the list ordered by priority

### Requirement: OpenCode Slash Commands
The provider SHALL create slash commands in `.opencode/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.opencode/commands/spectr/` directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it creates `.opencode/commands/spectr/proposal.md`
- **AND** it creates `.opencode/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field
