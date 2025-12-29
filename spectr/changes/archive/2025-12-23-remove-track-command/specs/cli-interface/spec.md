# Cli Interface Delta Spec

## REMOVED Requirements

### Requirement: Track Command

The CLI SHALL provide a `track` command that watches task status changes and automatically commits related changes.

#### Scenario: Track with change ID

- **WHEN** user runs `spectr track <change-id>`
- **THEN** the system watches tasks.json for the specified change
- **AND** displays current task status (X/Y completed)
- **AND** runs until all tasks are complete or interrupted

#### Scenario: Interactive track selection

- **WHEN** user runs `spectr track` without specifying a change ID
- **THEN** the system displays a list of active changes with tasks.json
- **AND** prompts for selection

#### Scenario: Auto-commit on task completion

- **WHEN** a task status changes to "completed" in tasks.json
- **THEN** the system detects modified files via git status
- **AND** stages all modified files except tasks.json, tasks.jsonc, tasks.md
- **AND** creates a commit with message format: `spectr(<change-id>): complete task <task-id>`
- **AND** includes footer: `[Automated by spectr track]`

#### Scenario: Auto-commit on task start

- **WHEN** a task status changes to "in_progress" in tasks.json
- **THEN** the system detects modified files via git status
- **AND** stages all modified files except tasks.json, tasks.jsonc, tasks.md
- **AND** creates a commit with message format: `spectr(<change-id>): start task <task-id>`
- **AND** includes footer: `[Automated by spectr track]`

#### Scenario: No files to commit warning

- **WHEN** a task status changes but no files have been modified (excluding task files)
- **THEN** the system prints a warning: "No files to commit for task <task-id>"
- **AND** continues watching for more task changes

#### Scenario: Git commit failure stops tracking

- **WHEN** a git commit operation fails (e.g., merge conflict, hook rejection)
- **THEN** the system displays the git error message
- **AND** stops tracking immediately
- **AND** exits with non-zero status code

#### Scenario: Graceful interruption

- **WHEN** user presses Ctrl+C during tracking
- **THEN** the system stops watching and exits cleanly
- **AND** displays "Tracking stopped" message

#### Scenario: All tasks already complete

- **WHEN** user runs `spectr track <change-id>` and all tasks are already completed
- **THEN** the system displays a message indicating all tasks are complete
- **AND** exits without starting the watch loop

### Requirement: Track Command Flags

The track command SHALL support flags for controlling behavior.

#### Scenario: No-interactive flag disables prompts

- **WHEN** user provides the `--no-interactive` flag
- **AND** no change-id is provided
- **THEN** the system displays usage error instead of prompting for selection

### Requirement: Track Command Binary Filtering

The track command SHALL support binary file filtering to prevent unintentional commits of binary files.

#### Scenario: Include-binaries flag enables binary file commits

- **WHEN** user provides the `--include-binaries` flag
- **THEN** the system includes binary files in automatic commits
- **AND** commits all modified files as before (excluding task files)

#### Scenario: Default binary file exclusion

- **WHEN** user runs `spectr track` without the `--include-binaries` flag
- **THEN** the system detects binary files using git diff --numstat
- **AND** excludes binary files from staging
- **AND** displays a warning listing skipped binary files
- **AND** continues to commit non-binary files normally

#### Scenario: Binary detection via git

- **WHEN** the system needs to determine if a file is binary
- **THEN** it SHALL use `git diff --numstat` to check for binary markers
- **AND** treat files with `-       -` output as binary

#### Scenario: Only binary files modified

- **WHEN** a task status changes and only binary files were modified (with no --include-binaries flag)
- **THEN** the system displays a warning: "No files to commit for task <task-id> (binary files excluded)"
- **AND** lists the skipped binary files
- **AND** continues watching for more task changes
