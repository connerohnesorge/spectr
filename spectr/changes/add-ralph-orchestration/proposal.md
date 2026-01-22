# Change: Add Ralph Task Orchestration TUI

## Why

Currently, implementing changes from a tasks.jsonc file requires manual work:
developers must read each task, copy context, invoke their AI CLI, and track
progress. This is error-prone and tedious for large change proposals with many
tasks.

Ralph automates this workflow by orchestrating agent CLI sessions for each task,
injecting full change context, streaming live output, and tracking completion
via status file polling.

## What Changes

- **NEW**: `spectr ralph <change-id>` command with interactive TUI
- **NEW**: `Ralpher` interface on Provider for CLI invocation configuration
- **NEW**: `internal/ralph/` package for orchestration logic
- **MODIFIED**: Provider interface gains optional `Ralpher` implementation
- Dependency-aware task execution with parallel independent tasks
- Live PTY output streaming in Bubble Tea TUI
- Session persistence for resume support
- Auto-retry with configurable limit on task failures

## Impact

- Affected specs: provider-system (new), ralph-orchestration (new)
- Affected code:
  - `cmd/ralph.go` - New CLI command
  - `internal/ralph/` - Orchestration engine, TUI, session state
  - `internal/initialize/providers/` - Ralpher interface addition
  - `internal/initialize/providers/claude.go` - Initial Ralpher implementation
