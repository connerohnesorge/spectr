## MODIFIED Requirements

### Requirement: CoStrict Provider Configuration
The provider SHALL be configured with these settings:
- ID: `costrict`
- Name: `CoStrict`
- Priority: 3
- Config File: `COSTRICT.md`
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the CoStrict provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `costrict`, Name `CoStrict`, Priority 3
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's Initializers() method is called
- **THEN** it SHALL return a DirectoryInitializer for `.costrict/commands/spectr/`
- **AND** it SHALL return a ConfigFileInitializer for `COSTRICT.md`
- **AND** it SHALL return a SlashCommandsInitializer for Markdown format slash commands

### Requirement: CoStrict Instruction File
The provider SHALL create and maintain a `COSTRICT.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with CoStrict provider selected
- **THEN** the ConfigFileInitializer creates `COSTRICT.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Instruction file updates
- **WHEN** `spectr init` runs in a project with CoStrict provider
- **THEN** the ConfigFileInitializer updates content between markers in `COSTRICT.md`
- **AND** preserves any user content outside the markers

### Requirement: CoStrict Slash Commands
The provider SHALL create slash commands in `.costrict/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.costrict/commands/spectr/` directory

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it creates `.costrict/commands/spectr/proposal.md`
- **AND** it creates `.costrict/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field

