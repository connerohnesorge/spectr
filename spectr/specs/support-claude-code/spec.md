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

#### Scenario: Provider identification

- **WHEN** the registry queries for Claude Code provider
- **THEN** it SHALL return provider with ID `claude-code`
- **AND** the provider priority is 1

#### Scenario: Provider metadata

- **WHEN** displaying provider options to users
- **THEN** the provider name is "Claude Code"
- **AND** it appears first in the list due to priority 1

### Requirement: Claude Code Instruction File

The provider SHALL create and maintain a `CLAUDE.md` instruction file in the
project root.

#### Scenario: Instruction file creation

- **WHEN** `spectr init` runs with Claude Code provider selected
- **THEN** the system creates `CLAUDE.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!--
  spectr:END -->` markers

#### Scenario: Instruction file updates

- **WHEN** `spectr init` runs in a project with Claude Code provider
- **THEN** the system updates content between markers in `CLAUDE.md`
- **AND** preserves any user content outside the markers

### Requirement: Claude Code Slash Commands

The provider SHALL create slash commands in `.claude/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** it creates `.claude/commands/spectr/` directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths

- **WHEN** the provider generates slash command files
- **THEN** it creates `.claude/commands/spectr/proposal.md`
- **AND** it creates `.claude/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

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
