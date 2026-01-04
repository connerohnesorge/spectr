# Agent Instructions Delta

## ADDED Requirements

### Requirement: Slash Command Frontmatter for Agentic Discovery

Slash command templates SHALL include YAML frontmatter with metadata fields that
enable automatic discovery and invocation by AI coding assistants.

#### Scenario: Claude Code discovers slash commands

- **WHEN** Claude Code loads a project with spectr slash commands
- **THEN** the `SlashCommand` tool SHALL be able to discover `/spectr:proposal`
  and `/spectr:apply` commands
- **AND** the commands SHALL appear in `/help` output with descriptions

#### Scenario: OpenCode routes commands to agents

- **WHEN** OpenCode loads a project with spectr slash commands
- **THEN** the command system SHALL read the `agent` frontmatter field
- **AND** SHALL route command execution to the specified agent when present

#### Scenario: Frontmatter uses superset approach

- **WHEN** spectr generates slash command files
- **THEN** the frontmatter SHALL include fields for multiple AI tools
- **AND** each tool SHALL ignore fields it does not recognize
- **AND** the following fields SHALL be supported:
  - `description`: Brief command description (Claude Code, OpenCode)
  - `allowed-tools`: Permitted tool list (Claude Code)
  - `agent`: Target agent for routing (OpenCode)
  - `model`: Optional model override (Claude Code, OpenCode)
  - `subtask`: Force subagent invocation (OpenCode)

### Requirement: Frontmatter Field Values

Slash command frontmatter fields SHALL have appropriate default values that
enable useful agentic behavior without requiring user customization.

#### Scenario: Proposal command frontmatter defaults

- **WHEN** spectr generates the proposal slash command
- **THEN** the `description` field SHALL be "Create a Spectr change proposal"
- **AND** the `allowed-tools` field SHALL include read and search tools
- **AND** the `agent` field SHALL be unset or null (inherit default)

#### Scenario: Apply command frontmatter defaults

- **WHEN** spectr generates the apply slash command
- **THEN** the `description` field SHALL be "Apply a Spectr change proposal"
- **AND** the `allowed-tools` field SHALL include Bash for running spectr accept
- **AND** the `agent` field SHALL be unset or null (inherit default)
