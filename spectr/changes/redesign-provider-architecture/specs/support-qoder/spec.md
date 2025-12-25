## MODIFIED Requirements

### Requirement: Qoder Provider Configuration
The provider SHALL be configured with these settings:
- ID: `qoder`
- Name: `Qoder`
- Priority: 4
- Config File: `QODER.md`
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Qoder provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `qoder`, Name `Qoder`, Priority 4
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's `Initializers(ctx, tm *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for `.qoder/commands/spectr/`
- **AND** it SHALL return a `ConfigFileInitializer` for `QODER.md` with TemplateRef from TemplateManager
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash commands

#### Scenario: Provider metadata
- **WHEN** the provider is registered
- **THEN** name SHALL be `Qoder`
- **AND** config file SHALL be `QODER.md`
- **AND** command format SHALL be Markdown

### Requirement: Qoder Instruction File
The provider SHALL create and maintain a `QODER.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with Qoder provider selected
- **THEN** the ConfigFileInitializer creates `QODER.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Instruction file update
- **WHEN** `spectr init` runs
- **THEN** the ConfigFileInitializer updates the Spectr instructions block in `QODER.md`
- **AND** preserves existing content outside the markers

### Requirement: Qoder Slash Commands
The provider SHALL create slash commands in `.qoder/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.qoder/commands/spectr/` directory
- **AND** all Spectr commands are placed under this directory

#### Scenario: Standard command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it creates `.qoder/commands/spectr/proposal.md`
- **AND** it creates `.qoder/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter at the top
- **AND** frontmatter includes a `description` field

