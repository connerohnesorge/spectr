# Delta Specification

## ADDED Requirements

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
