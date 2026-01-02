# Conclaude Hooks Specification

## Purpose

Documents the hook system for conclaude configuration, enabling Claude Code
users to extend behavior through lifecycle hooks including the new postToolUse
hook for observing tool results.

## ADDED Requirements

### Requirement: postToolUse Hook Configuration

The system SHALL support a `postToolUse` configuration section for hooking into
tool use results after they complete.

#### Scenario: Basic postToolUse configuration

- **WHEN** a `.conclaude.yaml` contains a `postToolUse` section
- **THEN** it SHALL accept a `commands` array
- **AND** each command entry SHALL support `tool` filter (optional)
- **AND** each command entry SHALL support `run` command string (required)
- **AND** each command entry SHALL support `showStdout` boolean (optional,
  default false)
- **AND** each command entry SHALL support `showStderr` boolean (optional,
  default false)
- **AND** each command entry SHALL support `timeout` in seconds (optional)

#### Scenario: postToolUse with tool filter

- **WHEN** a postToolUse command specifies a `tool` filter
- **THEN** the hook SHALL only execute for matching tool names
- **AND** exact matches SHALL be supported (e.g., `AskUserQuestion`)
- **AND** glob patterns SHALL be supported (e.g., `Ask*`, `*Search*`)
- **AND** omitting `tool` SHALL match all tools

#### Scenario: postToolUse without tool filter

- **WHEN** a postToolUse command omits the `tool` filter
- **THEN** the hook SHALL execute for every tool use result
- **AND** this enables universal logging of all tool activity

### Requirement: postToolUse Environment Variables

The system SHALL pass tool use data to hook commands via environment variables.

#### Scenario: Tool name environment variable

- **WHEN** a postToolUse hook command executes
- **THEN** `CONCLAUDE_TOOL_NAME` SHALL contain the tool name (e.g.,
  `AskUserQuestion`)
- **AND** the value SHALL be the exact tool name as used by Claude

#### Scenario: Tool input environment variable

- **WHEN** a postToolUse hook command executes
- **THEN** `CONCLAUDE_TOOL_INPUT` SHALL contain the tool input as JSON
- **AND** the JSON SHALL be compact (no pretty-printing)
- **AND** special characters SHALL be properly escaped

#### Scenario: Tool output environment variable

- **WHEN** a postToolUse hook command executes
- **THEN** `CONCLAUDE_TOOL_OUTPUT` SHALL contain the tool output as JSON
- **AND** for text outputs, it SHALL be JSON-encoded string
- **AND** for structured outputs, it SHALL be the JSON object

#### Scenario: Timestamp environment variable

- **WHEN** a postToolUse hook command executes
- **THEN** `CONCLAUDE_TOOL_TIMESTAMP` SHALL contain ISO 8601 timestamp
- **AND** the timestamp SHALL be in UTC (e.g., `2025-01-02T15:04:05Z`)

### Requirement: postToolUse Read-Only Semantics

The postToolUse hook SHALL be read-only with no ability to modify tool output or
inject messages.

#### Scenario: Hook output is discarded

- **WHEN** a postToolUse hook command produces stdout
- **THEN** the output SHALL NOT be injected into the conversation
- **AND** the output SHALL be discarded (unless showStdout is true for logging)
- **AND** Claude SHALL NOT see any hook output

#### Scenario: Hook cannot modify tool output

- **WHEN** a postToolUse hook executes
- **THEN** the original tool output SHALL be preserved
- **AND** no transformation of tool output SHALL occur
- **AND** the hook observes but cannot alter the tool result

#### Scenario: Hook failure does not block Claude

- **WHEN** a postToolUse hook command fails (non-zero exit)
- **THEN** the failure SHALL be logged (if showStderr enabled)
- **AND** Claude SHALL continue normally with the original tool output
- **AND** hook failures SHALL NOT interrupt the conversation

### Requirement: postToolUse Hook Execution

The system SHALL execute postToolUse hooks after each tool use completes.

#### Scenario: Sequential hook execution

- **WHEN** multiple postToolUse commands match a tool
- **THEN** they SHALL execute sequentially in configuration order
- **AND** each hook SHALL receive the same environment variables
- **AND** later hooks SHALL NOT depend on earlier hook outputs

#### Scenario: Async execution option

- **WHEN** a postToolUse command specifies `async: true`
- **THEN** the hook MAY execute asynchronously (fire-and-forget)
- **AND** Claude SHALL NOT wait for async hooks to complete
- **AND** async hooks SHALL still receive all environment variables

#### Scenario: Timeout handling

- **WHEN** a postToolUse hook exceeds its timeout
- **THEN** the hook process SHALL be terminated
- **AND** an error SHALL be logged (if showStderr enabled)
- **AND** Claude SHALL continue with the next hook or tool use

### Requirement: postToolUse Configuration Schema

The system SHALL validate postToolUse configuration against a schema.

#### Scenario: Valid postToolUse configuration

- **WHEN** validating a postToolUse section
- **THEN** `commands` SHALL be an array (required if section present)
- **AND** each command SHALL have `run` string (required)
- **AND** each command MAY have `tool` string (optional filter)
- **AND** each command MAY have `showStdout` boolean (optional)
- **AND** each command MAY have `showStderr` boolean (optional)
- **AND** each command MAY have `timeout` number (optional, seconds)
- **AND** each command MAY have `async` boolean (optional, default false)
- **AND** each command MAY have `enabled` boolean (optional, default true)

#### Scenario: Invalid postToolUse configuration

- **WHEN** a postToolUse command lacks `run` field
- **THEN** validation SHALL fail with clear error message
- **AND** the error SHALL indicate the missing field

### Requirement: AskUserQuestion Q&A Logging Example

The system SHALL document an example for logging AskUserQuestion interactions.

#### Scenario: Q&A append logging configuration

- **WHEN** a user wants to log all Q&A interactions
- **THEN** the configuration SHALL support:

```yaml
postToolUse:
  commands:
    - tool: "AskUserQuestion"
      run: "echo \"## $(date -Iseconds)\" >> .claude/qa-log.md && echo \"\" >> .claude/qa-log.md && echo \"**Q:** $CONCLAUDE_TOOL_INPUT\" >> .claude/qa-log.md && echo \"**A:** $CONCLAUDE_TOOL_OUTPUT\" >> .claude/qa-log.md && echo \"\" >> .claude/qa-log.md"
```

- **AND** this SHALL append Q&A pairs to `.claude/qa-log.md`
- **AND** each entry SHALL include timestamp, question, and answer

#### Scenario: Dedicated logging script

- **WHEN** complex logging logic is needed
- **THEN** users SHALL be able to reference external scripts:

```yaml
postToolUse:
  commands:
    - tool: "AskUserQuestion"
      run: ".claude/scripts/log-qa.sh"
```

- **AND** the script SHALL receive data via environment variables
- **AND** the script SHALL handle JSON parsing and formatting
