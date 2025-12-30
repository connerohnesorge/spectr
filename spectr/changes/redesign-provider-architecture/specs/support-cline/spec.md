# Delta Specification

## MODIFIED Requirements

### Requirement: Cline Provider Configuration

The provider SHALL be configured with these settings:

- ID: `cline`
- Name: `Cline`
- Priority: 7
- Config File: `CLINE.md`
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Cline provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `cline`, Name `Cline`, Priority 7
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.clinerules/commands/spectr/`
- **AND** it SHALL return a `ConfigFileInitializer` for `CLINE.md` with
  TemplateRef from TemplateManager
- **AND** it SHALL return a `SlashCommandsInitializer` for Markdown format slash
  commands

### Requirement: Cline Instruction File

The provider SHALL create and maintain a `CLINE.md` instruction file in the
project root.

#### Scenario: Instruction file creation

- **WHEN** `spectr init` runs with Cline provider selected
- **THEN** the ConfigFileInitializer creates `CLINE.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:start -->` and `<!--
  spectr:end -->` markers

#### Scenario: Instruction file updates

- **WHEN** `spectr init` runs in a project with Cline provider
- **THEN** the ConfigFileInitializer updates content between markers in
  `CLINE.md`
- **AND** preserves any user content outside the markers

### Requirement: Cline Slash Commands

The provider SHALL create slash commands in `.clinerules/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.clinerules/commands/spectr/`
  directory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.clinerules/commands/spectr/proposal.md`
- **AND** it SHALL create `.clinerules/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with YAML frontmatter
- **AND** frontmatter SHALL include `description` field
