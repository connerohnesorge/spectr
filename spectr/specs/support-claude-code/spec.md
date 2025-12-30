# Support Claude Code Specification

## Purpose

Documents the Claude Code provider integration for Spectr, enabling Spectr to
work seamlessly with Claude Code through instruction file management and slash
command generation.

## Requirements

### Requirement: Claude Code Provider Configuration

The provider SHALL be configured with these settings:

- ID: `claude-code`
- Name: `Claude Code`
- Priority: 1 (highest)
- Config File: `CLAUDE.md`
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Claude Code provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `claude-code`, Name `Claude Code`,
  Priority 1
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.claude/commands/spectr/`
- **AND** it SHALL return a `ConfigFileInitializer` for `CLAUDE.md` with
  TemplateRef from TemplateManager
- **AND** it SHALL return a `SlashCommandsInitializer` for slash commands in
  Markdown format

#### Scenario: Provider metadata

- **WHEN** the provider is registered
- **THEN** the provider name SHALL be "Claude Code"
- **AND** the provider priority SHALL be 1 (highest priority)
- **AND** the provider ID SHALL be "claude-code"

### Requirement: Claude Code Instruction File

The provider SHALL create and maintain a `CLAUDE.md` instruction file in the
project root.

#### Scenario: Instruction file creation

- **WHEN** `spectr init` runs with Claude Code provider selected
- **THEN** the ConfigFileInitializer creates `CLAUDE.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:start -->` and `<!--
  spectr:end -->` markers

#### Scenario: Instruction file updates

- **WHEN** `spectr init` runs in a project with Claude Code provider
- **THEN** the ConfigFileInitializer updates content between markers in
  `CLAUDE.md`
- **AND** preserves any user content outside the markers

### Requirement: Claude Code Slash Commands

The provider SHALL create slash commands in `.claude/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.claude/commands/spectr/`
  directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.claude/commands/spectr/proposal.md`
- **AND** it SHALL create `.claude/commands/spectr/apply.md`

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

#### Scenario: Sync command frontmatter

- **WHEN** generating the sync command file
- **THEN** the frontmatter description is "Detect spec drift from code and
  update specs interactively."

### Requirement: Command Path Construction

The provider SHALL construct command paths using the standard pattern.

#### Scenario: Path construction for Claude Code

- **WHEN** determining command file paths
- **THEN** it uses base directory `.claude/commands`
- **AND** appends `/spectr/` subdirectory
- **AND** appends command name with `.md` extension
- **AND** results in paths like `.claude/commands/spectr/proposal.md`

### Requirement: Claude Code Skills Directory

The provider SHALL create a `.claude/skills/` directory for agent skills.

#### Scenario: Skills directory creation

- **WHEN** the provider returns initializers
- **THEN** it SHALL include a `DirectoryInitializer` for `.claude/skills/`
- **AND** the directory SHALL be created before skills are installed

### Requirement: spectr-accept-wo-spectr-bin Skill

The provider SHALL install the `spectr-accept-wo-spectr-bin` skill for
accepting changes without the spectr binary.

#### Scenario: Skill installation path

- **WHEN** the provider returns initializers
- **THEN** it SHALL include an `AgentSkillsInitializer` for
  `spectr-accept-wo-spectr-bin`
- **AND** the skill SHALL be installed at
  `.claude/skills/spectr-accept-wo-spectr-bin/`

#### Scenario: Skill structure

- **WHEN** the skill is installed
- **THEN** it SHALL create `SKILL.md` with valid AgentSkills frontmatter
- **AND** it SHALL create `scripts/accept.sh` for task conversion
- **AND** `scripts/accept.sh` SHALL be executable

#### Scenario: SKILL.md content

- **WHEN** the `SKILL.md` file is created
- **THEN** the frontmatter `name` SHALL be `spectr-accept-wo-spectr-bin`
- **AND** the frontmatter `description` SHALL describe the skill's purpose
- **AND** the `compatibility` field SHALL note `jq` as a requirement
- **AND** the body SHALL contain usage instructions for the accept script

#### Scenario: accept.sh functionality

- **WHEN** `scripts/accept.sh` is executed with a change-id argument
- **THEN** it SHALL read `spectr/changes/<id>/tasks.md`
- **AND** it SHALL parse markdown task lists into structured format
- **AND** it SHALL output `spectr/changes/<id>/tasks.jsonc`
- **AND** the output SHALL be valid JSONC format
