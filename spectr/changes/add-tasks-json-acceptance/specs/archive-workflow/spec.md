## MODIFIED Requirements

### Requirement: Task Completion Checking
The system SHALL check task completion status from either tasks.json or tasks.md, preferring the JSON format when present.

#### Scenario: Check tasks from tasks.json
- **WHEN** archiving a change that has tasks.json
- **THEN** the system reads completion status from the JSON file
- **AND** uses summary.total and summary.completed for display

#### Scenario: Fall back to tasks.md
- **WHEN** archiving a change that has tasks.md but no tasks.json
- **THEN** the system reads completion status from the Markdown file
- **AND** counts `[x]` and `[ ]` markers as before

#### Scenario: Display task status
- **WHEN** archiving a change
- **THEN** the system displays task completion status (e.g., "3/5 complete")

#### Scenario: Warn on incomplete tasks
- **WHEN** a change has incomplete tasks
- **THEN** the system warns the user and requires confirmation to proceed (unless --yes flag is provided)

#### Scenario: Proceed with incomplete tasks when confirmed
- **WHEN** user confirms archiving despite incomplete tasks
- **THEN** the system proceeds with the archive operation

## ADDED Requirements

### Requirement: Task Format Detection
The system SHALL automatically detect whether a change uses tasks.json or tasks.md format.

#### Scenario: Prefer tasks.json when both exist
- **WHEN** a change directory contains both tasks.json and tasks.md
- **THEN** the system SHALL use tasks.json as the source of truth
- **AND** SHALL log a warning about the duplicate files

#### Scenario: Use tasks.md when no JSON exists
- **WHEN** a change directory contains only tasks.md
- **THEN** the system SHALL use tasks.md for task tracking

#### Scenario: Handle missing task files
- **WHEN** a change directory has neither tasks.json nor tasks.md
- **THEN** the system SHALL report zero tasks
- **AND** SHALL continue without error

### Requirement: tasks.json Parsing
The system SHALL parse tasks.json files to extract completion metrics.

#### Scenario: Parse summary from JSON
- **WHEN** reading a tasks.json file
- **THEN** the system extracts total and completed from the summary object

#### Scenario: Handle malformed JSON
- **WHEN** tasks.json exists but is not valid JSON
- **THEN** the system displays an error with the file path and parse error
- **AND** blocks archive until the file is fixed

#### Scenario: Handle missing summary
- **WHEN** tasks.json exists but lacks a summary object
- **THEN** the system calculates totals by iterating through sections and tasks (including nested subtasks recursively)
