# Tasks

## 1. Core Infrastructure

- [ ] 1.1 Create `internal/ralph/` package directory structure
- [ ] 1.2 Define Task and TaskGraph types in `internal/ralph/graph.go`
- [ ] 1.3 Implement task graph parsing from tasks.jsonc files
- [ ] 1.4 Implement topological sort for dependency-aware execution order
- [ ] 1.5 Add parallel task detection (tasks with different prefixes)

## 2. Ralpher Interface

- [ ] 2.1 Define Ralpher interface in `internal/initialize/providers/ralpher.go`
- [ ] 2.2 Implement Ralpher for ClaudeProvider in `internal/initialize/providers/claude.go`
- [ ] 2.3 Add binary detection helper to check CLI availability
- [ ] 2.4 Write tests for Ralpher interface and Claude implementation

## 3. Prompt Generation

- [ ] 3.1 Create prompt template structure in `internal/ralph/prompt.go`
- [ ] 3.2 Implement change context loading (proposal.md, design.md, specs)
- [ ] 3.3 Implement prompt assembly with task details and context
- [ ] 3.4 Write tests for prompt generation with various context combinations

## 4. PTY Subprocess Management

- [ ] 4.1 Add creack/pty dependency for cross-platform PTY support
- [ ] 4.2 Implement PTY spawning and management in `internal/ralph/pty.go`
- [ ] 4.3 Handle PTY resize events from TUI
- [ ] 4.4 Implement graceful process termination on skip/abort
- [ ] 4.5 Write tests for PTY lifecycle management

## 5. Status File Watcher

- [ ] 5.1 Implement StatusWatcher in `internal/ralph/watcher.go`
- [ ] 5.2 Add polling loop with configurable interval (default 2s)
- [ ] 5.3 Detect status changes and emit events
- [ ] 5.4 Handle split file discovery (tasks-*.jsonc glob)
- [ ] 5.5 Write tests for status change detection

## 6. Session Persistence

- [ ] 6.1 Define SessionState struct in `internal/ralph/session.go`
- [ ] 6.2 Implement session save on interruption/quit
- [ ] 6.3 Implement session load and resume prompt
- [ ] 6.4 Implement session cleanup on completion
- [ ] 6.5 Write tests for session persistence round-trip

## 7. Orchestration Engine

- [ ] 7.1 Implement main orchestration loop in `internal/ralph/orchestrator.go`
- [ ] 7.2 Integrate task graph, prompt generation, PTY, and watcher
- [ ] 7.3 Implement retry logic with configurable maxRetries
- [ ] 7.4 Implement parallel execution for independent tasks
- [ ] 7.5 Handle user actions (retry, skip, abort, pause)
- [ ] 7.6 Write integration tests for orchestration scenarios

## 8. TUI Implementation

- [ ] 8.1 Create TUI model in `internal/ralph/tui.go` using Bubble Tea
- [ ] 8.2 Implement task list view with status indicators
- [ ] 8.3 Implement agent output pane with ANSI rendering
- [ ] 8.4 Implement keyboard controls (q, r, s, p)
- [ ] 8.5 Implement interactive task selection mode (--interactive)
- [ ] 8.6 Add help bar with available commands
- [ ] 8.7 Write TUI tests using teatest

## 9. CLI Command

- [ ] 9.1 Create `cmd/ralph.go` with Kong command struct
- [ ] 9.2 Implement change-id argument parsing and validation
- [ ] 9.3 Add --interactive flag for task selection mode
- [ ] 9.4 Add --max-retries flag (default 3)
- [ ] 9.5 Implement provider detection and Ralpher lookup
- [ ] 9.6 Wire up TUI and orchestrator
- [ ] 9.7 Write CLI integration tests

## 10. Documentation and Polish

- [ ] 10.1 Add ralph command to spectr help output
- [ ] 10.2 Update AGENTS.md with ralph usage instructions
- [ ] 10.3 Add error messages for common failure modes
- [ ] 10.4 Run full test suite and fix any failures
