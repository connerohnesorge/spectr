# Design

## Implementation Details

### Ralpher Interface

```go
// Ralpher is an optional interface that providers can implement
// to support task orchestration via spectr ralph.
type Ralpher interface {
    // InvokeTask creates an exec.Cmd configured to run the agent CLI
    // with the given task context. The prompt is passed according to
    // the provider's input method (stdin, file, argument).
    //
    // Parameters:
    //   - ctx: Context for cancellation
    //   - task: The task being executed (id, description, section)
    //   - prompt: Full prompt content with injected context
    //
    // Returns exec.Cmd ready for PTY attachment, or error if
    // the provider cannot be invoked (binary not found, etc.)
    InvokeTask(ctx context.Context, task Task, prompt string) (*exec.Cmd, error)

    // Binary returns the CLI binary name for display and detection.
    // Example: "claude", "gemini", "cursor"
    Binary() string
}
```

### Task Dependency Graph

Tasks are parsed from tasks*.jsonc files and organized into a dependency graph:

```go
type TaskGraph struct {
    Tasks    map[string]*Task  // task ID -> task
    Children map[string][]string // parent ID -> child IDs
    Roots    []string          // tasks with no dependencies
}

// Task represents a single task from tasks.jsonc
type Task struct {
    ID          string
    Section     string
    Description string
    Status      string // pending, in_progress, completed
    Children    string // "$ref:tasks-N.jsonc" or empty
}
```

Execution order:

1. Parse all tasks*.jsonc files into unified graph
2. Identify root tasks (no parent prefix in ID)
3. Execute tasks in topological order
4. Parallelize tasks that share no prefix (independent sections)

### Prompt Template Structure

```markdown
# Task: {task.ID} - {task.Section}

## Task Description
{task.Description}

## Change Context

### Proposal
{contents of proposal.md}

### Design (if exists)
{contents of design.md}

### Relevant Specs
{contents of delta specs from spectr/changes/<id>/specs/}

## Instructions
Complete this task and update the task status in tasks.jsonc to "completed"
when done. If blocked, set status to "in_progress" and describe the blocker.
```

### TUI Layout (Bubble Tea)

```text
┌─────────────────────────────────────────────────────────────┐
│ spectr ralph: add-feature-x                    [3/12 tasks] │
├─────────────────────────────────────────────────────────────┤
│ Tasks                                                       │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ ✓ 1.1 Set up project structure                         │ │
│ │ ✓ 1.2 Create database schema                           │ │
│ │ ▶ 1.3 Implement API endpoint              [in_progress] │ │
│ │ ○ 1.4 Add frontend component                  [pending] │ │
│ │ ○ 2.1 Write unit tests                        [pending] │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ Agent Output (task 1.3)                                     │
│ ┌─────────────────────────────────────────────────────────┐ │
│ │ Reading file: internal/api/handler.go                   │ │
│ │ I'll implement the GET /users endpoint...               │ │
│ │ Writing to internal/api/handler.go...                   │ │
│ │ █                                                       │ │
│ └─────────────────────────────────────────────────────────┘ │
├─────────────────────────────────────────────────────────────┤
│ [q] quit  [r] retry  [s] skip  [p] pause  [i] interactive   │
└─────────────────────────────────────────────────────────────┘
```

### Session State Persistence

```go
// SessionState persists orchestration progress for resume
type SessionState struct {
    ChangeID      string            `json:"change_id"`
    StartedAt     time.Time         `json:"started_at"`
    LastUpdated   time.Time         `json:"last_updated"`
    CompletedIDs  []string          `json:"completed_ids"`
    FailedIDs     []string          `json:"failed_ids"`
    RetryCount    map[string]int    `json:"retry_count"`
    CurrentTaskID string            `json:"current_task_id,omitempty"`
}
```

Stored at: `spectr/changes/<id>/.ralph-session.json`

### Status File Polling

The orchestrator polls tasks.jsonc every 2 seconds to detect status changes:

```go
type StatusWatcher struct {
    paths    []string      // tasks.jsonc, tasks-1.jsonc, etc.
    interval time.Duration
    onChange func(taskID string, newStatus string)
}
```

Task completion detected when:

1. Agent process exits with code 0, AND
2. Task status in tasks.jsonc changed to "completed"

If agent exits but status unchanged, treat as failure and retry.

### Error Handling

Retry logic:

1. On task failure (non-zero exit or timeout), increment retry count
2. If retries < maxRetries (default 3), re-invoke task
3. If retries exhausted, pause orchestration and prompt user:
   - Retry: Reset retry count and try again
   - Skip: Mark task skipped, continue to next
   - Abort: Save session state and exit

### Package Structure

```text
internal/ralph/
├── orchestrator.go    # Main orchestration loop
├── graph.go           # Task dependency graph
├── prompt.go          # Prompt template generation
├── session.go         # Session state persistence
├── watcher.go         # Status file polling
├── tui.go             # Bubble Tea TUI model
├── tui_views.go       # TUI view rendering
└── pty.go             # PTY subprocess management
```

## Context

This feature addresses GitHub issue #347 "Ralph Design TUI". The name "ralph"
is a project-specific term for the orchestration loop that continuously feeds
tasks to agent CLIs.

## Goals / Non-Goals

**Goals:**

- Automate task execution from tasks.jsonc files
- Provide live visibility into agent progress
- Support resume after interruption
- Handle failures gracefully with retries

**Non-Goals:**

- Support all providers immediately (Claude Code only initially)
- Replace manual task execution (ralph is optional)
- Modify how tasks.jsonc is generated (accept command unchanged)

## Decisions

1. **PTY over pipes**: Using PTY gives full terminal emulation, supporting
   agent CLIs that use colors, progress bars, or interactive prompts.

2. **Status polling over IPC**: Polling tasks.jsonc is simpler than custom
   IPC and works with any agent that can edit files.

3. **In-memory prompts**: Avoids temp file cleanup complexity and prevents
   prompt files from cluttering the workspace.

4. **Bubble Tea TUI**: Consistent with existing spectr TUI components and
   provides rich terminal UI capabilities.

## Risks / Trade-offs

- **PTY complexity**: Cross-platform PTY handling can be tricky. Mitigate
  with creack/pty library (already used in ecosystem).

- **Polling latency**: 2-second poll interval means up to 2s delay in
  detecting completion. Acceptable for human-observed orchestration.

- **Single provider**: Initial Claude-only support limits immediate utility.
  Mitigate by designing Ralpher interface for easy provider addition.

## Open Questions

None - all key decisions made during requirements gathering.
