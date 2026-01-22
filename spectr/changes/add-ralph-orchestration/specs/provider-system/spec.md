# Provider System Changes

## ADDED Requirements

### Requirement: Ralpher Interface

The system SHALL define a Ralpher interface that providers can optionally
implement to support task orchestration via `spectr ralph`.

#### Scenario: Provider implements Ralpher

- WHEN a provider implements the Ralpher interface
- THEN it SHALL be usable with `spectr ralph`
- AND provide InvokeTask and Binary methods

#### Scenario: Provider does not implement Ralpher

- WHEN a provider does not implement Ralpher
- THEN `spectr ralph` SHALL skip that provider
- AND display "Provider X does not support ralph orchestration"

#### Scenario: InvokeTask returns configured command

- WHEN InvokeTask is called with task and prompt
- THEN it SHALL return an exec.Cmd configured for PTY attachment
- AND the command SHALL be ready to execute immediately

#### Scenario: Binary returns CLI name

- WHEN Binary is called
- THEN it SHALL return the CLI binary name (e.g., "claude")
- AND the name SHALL be used for display and binary detection

### Requirement: Claude Code Ralpher Implementation

The system SHALL implement Ralpher for Claude Code provider as the initial
supported provider.

#### Scenario: Claude Code InvokeTask

- WHEN InvokeTask is called for Claude Code
- THEN it SHALL return `exec.Command("claude", "--print", "--dangerously-skip-permissions")`
- AND configure stdin to receive the prompt
- AND set appropriate environment variables

#### Scenario: Claude Code Binary

- WHEN Binary is called for Claude Code
- THEN it SHALL return "claude"

#### Scenario: Detect Claude Code availability

- WHEN initializing ralph for Claude Code
- THEN the system SHALL check if "claude" binary exists in PATH
- AND return error "claude CLI not found" if missing
