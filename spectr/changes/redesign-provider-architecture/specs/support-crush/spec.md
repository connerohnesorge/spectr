## MODIFIED Requirements

### Requirement: Crush Provider Configuration
The provider SHALL be configured with these settings:
- ID: `crush`
- Name: `Crush`
- Priority: 16 (after existing providers)
- Config File: `CRUSH.md`
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Crush provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `crush`, Name `Crush`, Priority 16
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's Initializers() method is called
- **THEN** it SHALL return a DirectoryInitializer for `.crush/commands/spectr/`
- **AND** it SHALL return a ConfigFileInitializer for `CRUSH.md`
- **AND** it SHALL return a SlashCommandsInitializer for Markdown format slash commands

#### Scenario: Provider metadata
- **WHEN** provider is registered
- **THEN** the provider name is "Crush"
- **AND** it appears in the provider list after Continue (priority 15)

### Requirement: Crush Instruction File
The provider SHALL create and maintain a `CRUSH.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with Crush provider selected
- **THEN** the ConfigFileInitializer creates `CRUSH.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Instruction file updates
- **WHEN** `spectr init` runs in a project with Crush provider
- **THEN** the ConfigFileInitializer updates content between markers in `CRUSH.md`
- **AND** preserves any user content outside the markers

### Requirement: Crush Slash Commands
The provider SHALL create slash commands in `.crush/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.crush/commands/spectr/` directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it creates `.crush/commands/spectr/proposal.md`
- **AND** it creates `.crush/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

