## MODIFIED Requirements

### Requirement: Qwen Provider Configuration
The provider SHALL be configured with these settings:
- ID: `qwen`
- Name: `Qwen Code`
- Priority: 5
- Config File: `QWEN.md`
- Command Format: Markdown

#### Scenario: Provider registration
- **WHEN** the Qwen provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `qwen`, Name `Qwen Code`, Priority 5
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers
- **WHEN** the provider's `Initializers(ctx, tm *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for `.qwen/commands/spectr/`
- **AND** it SHALL return a `ConfigFileInitializer` for `QWEN.md` with TemplateRef from TemplateManager
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash commands

### Requirement: Qwen Instruction File
The provider SHALL create and maintain a `QWEN.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with Qwen provider selected
- **THEN** the ConfigFileInitializer creates `QWEN.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Instruction file updates
- **WHEN** `spectr init` runs in a project with Qwen provider
- **THEN** the ConfigFileInitializer updates content between markers in `QWEN.md`
- **AND** preserves any user content outside the markers

### Requirement: Qwen Slash Commands
The provider SHALL create slash commands in `.qwen/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.qwen/commands/spectr/` directory

#### Scenario: Command paths
- **WHEN** the SlashCommandsInitializer executes
- **THEN** it creates `.qwen/commands/spectr/proposal.md`
- **AND** it creates `.qwen/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field

