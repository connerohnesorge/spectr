# Agent Instructions Specification

## Purpose

Defines how AI agent prompts should guide assistants to discover and interact
with spectr project structures, preferring direct file system access over CLI
commands for better reliability and context.

## Requirements

### Requirement: Direct File Access for Agents

Agent prompts SHALL instruct AI assistants to use direct file and directory
access methods (such as `ls spectr/changes/`, `ls spectr/specs/`, or file reads)
instead of CLI commands like `spectr list` to discover changes and
specifications.

#### Scenario: Agent discovering active changes

- **WHEN** an agent needs to find active changes in a project
- **THEN** the agent prompt SHALL instruct reading `spectr/changes/` directory
  directly
- **AND** SHALL NOT instruct running `spectr list`

#### Scenario: Agent discovering specifications

- **WHEN** an agent needs to find existing specifications in a project
- **THEN** the agent prompt SHALL instruct reading `spectr/specs/` directory
  directly
- **AND** SHALL NOT instruct running `spectr list --specs`

#### Scenario: Agent grounding proposal in current state

- **WHEN** an agent is creating a new change proposal
- **THEN** the agent prompt SHALL instruct reading `spectr/project.md` and
  exploring directories with `ls` or `rg`
- **AND** SHALL NOT require running `spectr list` commands

### Requirement: User Documentation Preserved

The `spectr list` command references SHALL remain in user-facing documentation
since formatted CLI output benefits human users.

#### Scenario: User-facing documentation unchanged

- **WHEN** a user reads README.md or docs/ content
- **THEN** they SHALL still see `spectr list` command examples and documentation
- **AND** the CLI command behavior SHALL remain unchanged

### Requirement: Delegation Context for Subagents

When orchestrators delegate implementation tasks to subagents or when agents
complete tasks from a change proposal, the instruction pointer SHALL include
guidance to provide change directory paths so subagents can reference the
authoritative specification.

#### Scenario: Orchestrator delegating task to coder subagent

- **WHEN** an orchestrator delegates a task from an active change proposal to a
  coder subagent
- **THEN** the instruction pointer SHALL guide the orchestrator to include the
  path to `<changes-dir>/<id>/proposal.md`
- **AND** SHALL guide inclusion of `<changes-dir>/<id>/tasks.jsonc` for task
  context
- **AND** SHALL guide inclusion of relevant delta spec paths
  `<changes-dir>/<id>/specs/<capability>/spec.md`

#### Scenario: Agent completing tasks from change proposal

- **WHEN** an agent is completing tasks defined in a change proposal
- **THEN** the instruction pointer SHALL instruct the agent to read the proposal
  and tasks files for authoritative context
- **AND** SHALL reference the change directory using template variables for
  dynamic paths

### Requirement: Incremental Task Status Updates

When agents complete tasks from a change proposal, task status SHALL be updated
immediately after each individual task is verified, not in batch at the end of
all work.

#### Scenario: Agent completing a single task

- **WHEN** an agent finishes implementing and verifying a single task from
  `tasks.jsonc`
- **THEN** the agent SHALL mark that task as `"completed"` immediately
- **AND** SHALL NOT wait until all tasks are done to update statuses

#### Scenario: Agent starting work on a task

- **WHEN** an agent begins work on a task from `tasks.jsonc`
- **THEN** the agent SHALL mark that task as `"in_progress"` before starting
  implementation
- **AND** SHALL NOT batch status transitions with other tasks

#### Scenario: Multiple tasks in sequence

- **WHEN** an agent is assigned multiple tasks to complete sequentially
- **THEN** each task's status SHALL be updated individually as it transitions
  through states
- **AND** the `tasks.jsonc` file SHALL reflect accurate progress at any point in
  time

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
