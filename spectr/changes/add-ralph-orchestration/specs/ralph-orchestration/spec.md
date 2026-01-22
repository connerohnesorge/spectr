# Ralph Orchestration Changes

## ADDED Requirements

### Requirement: Ralph CLI Command

The system SHALL provide a `spectr ralph <change-id>` command that orchestrates
agent CLI sessions for tasks in the specified change's tasks.jsonc files.

#### Scenario: Start orchestration for change

- WHEN user runs `spectr ralph add-feature-x`
- THEN the system SHALL locate `spectr/changes/add-feature-x/tasks.jsonc`
- AND parse all tasks including split files (tasks-*.jsonc)
- AND launch the TUI showing task list and agent output panes

#### Scenario: Change not found

- WHEN user runs `spectr ralph nonexistent-change`
- THEN the system SHALL exit with error "change 'nonexistent-change' not found"
- AND suggest running `spectr list` to see available changes

#### Scenario: No tasks.jsonc exists

- WHEN user runs `spectr ralph add-feature-x` but tasks.jsonc does not exist
- THEN the system SHALL exit with error "tasks.jsonc not found for change 'add-feature-x'"
- AND suggest running `spectr accept add-feature-x` first

### Requirement: Interactive Task Selection

The system SHALL support both automatic and interactive task selection modes.

#### Scenario: Default run-all mode

- WHEN user runs `spectr ralph <change-id>` without flags
- THEN the system SHALL automatically execute all pending tasks in dependency order

#### Scenario: Interactive selection mode

- WHEN user runs `spectr ralph <change-id> --interactive`
- THEN the TUI SHALL display a task selector before starting orchestration
- AND user can toggle tasks on/off with space key
- AND user presses enter to start with selected tasks only

### Requirement: Dependency-Aware Execution

The system SHALL execute tasks in dependency order based on hierarchical task IDs,
parallelizing independent tasks when possible.

#### Scenario: Sequential dependent tasks

- WHEN tasks 1.1, 1.2, 1.3 exist with shared prefix
- THEN the system SHALL execute 1.1 before 1.2 before 1.3
- AND wait for each to complete before starting the next

#### Scenario: Parallel independent tasks

- WHEN tasks 1.1 and 2.1 exist with different prefixes
- THEN the system MAY execute them in parallel
- AND track progress for both simultaneously in TUI

#### Scenario: Child task dependencies

- WHEN a root task has children referenced via `$ref:tasks-N.jsonc`
- THEN the system SHALL execute all child tasks before marking root complete
- AND child tasks follow same dependency rules within their file

### Requirement: Full Context Injection

The system SHALL inject complete change context into each task's prompt.

#### Scenario: Prompt includes task details

- WHEN generating prompt for task 1.3
- THEN the prompt SHALL include task ID, section, and description

#### Scenario: Prompt includes proposal

- WHEN generating prompt for any task
- THEN the prompt SHALL include contents of proposal.md

#### Scenario: Prompt includes design if exists

- WHEN design.md exists in the change directory
- THEN the prompt SHALL include its contents
- AND clearly label it as design context

#### Scenario: Prompt includes delta specs

- WHEN delta specs exist in `spectr/changes/<change-id>/specs/`
- THEN the prompt SHALL include all spec.md contents
- AND organize by capability name

### Requirement: Live Agent Output Streaming

The system SHALL stream agent CLI output in real-time within the TUI.

#### Scenario: Display agent output

- WHEN agent CLI produces output
- THEN the TUI SHALL display it in the agent output pane
- AND scroll automatically to show latest output

#### Scenario: Handle ANSI escape codes

- WHEN agent output contains ANSI color codes
- THEN the TUI SHALL render colors correctly
- AND handle cursor movement sequences appropriately

#### Scenario: Output pane scrollback

- WHEN agent output exceeds visible area
- THEN the TUI SHALL maintain scrollback buffer
- AND allow user to scroll up to see previous output

### Requirement: Status File Polling

The system SHALL detect task completion by polling tasks.jsonc for status changes.

#### Scenario: Detect completion

- WHEN task status changes from "in_progress" to "completed" in tasks.jsonc
- THEN the system SHALL recognize task as complete
- AND update TUI to show completion
- AND proceed to next task

#### Scenario: Poll interval

- WHEN monitoring task status
- THEN the system SHALL poll every 2 seconds
- AND minimize filesystem overhead

#### Scenario: Handle external status changes

- WHEN user manually edits tasks.jsonc status
- THEN the system SHALL detect and respect the change
- AND update orchestration state accordingly

### Requirement: Session Persistence

The system SHALL persist orchestration state to support resume after interruption.

#### Scenario: Save session state

- WHEN orchestration is interrupted (quit, crash, Ctrl+C)
- THEN the system SHALL save state to `.ralph-session.json`
- AND include completed tasks, failed tasks, retry counts

#### Scenario: Resume session

- WHEN user runs `spectr ralph <change-id>` and session file exists
- THEN the system SHALL prompt "Resume previous session? [Y/n]"
- AND if yes, continue from last state
- AND if no, start fresh (delete session file)

#### Scenario: Session cleanup

- WHEN all tasks complete successfully
- THEN the system SHALL delete `.ralph-session.json`
- AND report final summary

### Requirement: Error Handling with Retries

The system SHALL automatically retry failed tasks with configurable limits.

#### Scenario: Automatic retry

- WHEN agent exits with non-zero code
- THEN the system SHALL retry the task
- AND increment retry counter

#### Scenario: Retry limit reached

- WHEN retry count reaches maxRetries (default 3)
- THEN the system SHALL pause orchestration
- AND prompt user with options: retry, skip, abort

#### Scenario: Skip task

- WHEN user selects "skip" for a failed task
- THEN the system SHALL mark task as skipped
- AND continue to next task if independent
- AND skip dependent tasks with warning

#### Scenario: Abort orchestration

- WHEN user selects "abort"
- THEN the system SHALL save session state
- AND exit gracefully with summary of progress

### Requirement: TUI Keyboard Controls

The system SHALL provide keyboard controls for orchestration management.

#### Scenario: Quit orchestration

- WHEN user presses 'q' or Ctrl+C
- THEN the system SHALL confirm quit if task in progress
- AND save session state before exiting

#### Scenario: Retry current task

- WHEN user presses 'r' while task in progress or failed
- THEN the system SHALL terminate current attempt if running
- AND restart the task from beginning

#### Scenario: Skip current task

- WHEN user presses 's' while task in progress
- THEN the system SHALL terminate current task
- AND mark as skipped
- AND proceed to next task

#### Scenario: Pause orchestration

- WHEN user presses 'p'
- THEN the system SHALL pause after current task completes
- AND display "Paused - press 'p' to resume"
