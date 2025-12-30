# Support Opencode Specification

## Requirements

### Requirement: OpenCode Provider Configuration

The provider SHALL be configured with these settings:

- ID: `opencode`
- Name: `OpenCode`
- Priority: 15
- Config File: None (OpenCode uses JSON config, instruction injection not
  supported)
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the OpenCode provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `opencode`, Name `OpenCode`, Priority 15
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.opencode/commands/spectr/`
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash
  commands
- **AND** it SHALL NOT return a `ConfigFileInitializer` (OpenCode uses JSON
  config)

#### Scenario: Provider metadata

- **WHEN** provider is registered
- **THEN** the provider name is "OpenCode"
- **AND** it appears in the list ordered by priority

### Requirement: OpenCode Slash Commands

The provider SHALL create slash commands in `.opencode/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.opencode/commands/spectr/`
  directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.opencode/commands/spectr/proposal.md`
- **AND** it SHALL create `.opencode/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** each file SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field

### Requirement: Standard Frontmatter

The provider SHALL use standard frontmatter templates for each command type.

#### Scenario: Proposal command frontmatter

- **WHEN** generating the proposal command file
- **THEN** the frontmatter description is "Scaffold a new Spectr change and
  validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** generating the apply command file
- **THEN** the frontmatter description is "Implement an approved Spectr change
  and keep tasks in sync."

### Requirement: Command Path Construction

The provider SHALL construct command paths using the standard pattern.

#### Scenario: Path construction for OpenCode

- **WHEN** determining command file paths
- **THEN** it uses base directory `.opencode/command`
- **AND** appends `/spectr/` subdirectory
- **AND** appends command name with `.md` extension
- **AND** results in paths like `.opencode/command/spectr/proposal.md`

### Requirement: No Instruction File

The provider SHALL NOT create an instruction file since OpenCode uses JSON
configuration.

#### Scenario: Config file check

- **WHEN** checking if provider has config file
- **THEN** it returns false
- **AND** no instruction file is created during initialization
