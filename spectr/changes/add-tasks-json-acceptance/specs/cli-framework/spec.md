## ADDED Requirements

### Requirement: Accept Command Structure
The CLI SHALL provide an `accept` command that converts tasks.md to tasks.json format when a change proposal is approved for implementation.

#### Scenario: Accept command registration
- **WHEN** the CLI is initialized
- **THEN** it SHALL include an AcceptCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr accept`
- **AND** help text SHALL describe task format conversion functionality

#### Scenario: Accept with change ID
- **WHEN** user runs `spectr accept <change-id>`
- **THEN** the system converts tasks.md to tasks.json for the specified change
- **AND** removes the original tasks.md file
- **AND** displays a success message with the generated file path

#### Scenario: Interactive accept selection
- **WHEN** user runs `spectr accept` without specifying a change ID
- **THEN** the system displays a list of active changes with tasks.md files
- **AND** prompts for selection

#### Scenario: Non-interactive accept with yes flag
- **WHEN** user runs `spectr accept <change-id> --yes`
- **THEN** the system converts without any confirmation prompts

### Requirement: Accept Command Flags
The accept command SHALL support flags for controlling behavior.

#### Scenario: Yes flag skips confirmation
- **WHEN** user provides the `-y` or `--yes` flag
- **THEN** the system skips the confirmation prompt before conversion

### Requirement: Accept Command Validation
The accept command SHALL validate preconditions before conversion.

#### Scenario: Change must exist
- **WHEN** user runs `spectr accept <change-id>` for a non-existent change
- **THEN** the system displays an error message and exits with non-zero code

#### Scenario: tasks.md must exist
- **WHEN** user runs `spectr accept <change-id>` and tasks.md does not exist
- **THEN** the system displays an error: "No tasks.md found for change '<change-id>'"

#### Scenario: tasks.json must not already exist
- **WHEN** user runs `spectr accept <change-id>` and tasks.json already exists
- **THEN** the system displays an error: "Change '<change-id>' has already been accepted"

### Requirement: tasks.json Schema
The accept command SHALL generate tasks.json with a structured schema that preserves task hierarchy and metadata.

#### Scenario: JSON structure fields
- **WHEN** tasks.json is generated
- **THEN** it SHALL include `version` field with value "1.0"
- **AND** it SHALL include `changeId` field with the change identifier
- **AND** it SHALL include `acceptedAt` field with ISO 8601 timestamp
- **AND** it SHALL include `sections` array with section objects
- **AND** it SHALL include `summary` object with total and completed counts

#### Scenario: Section object structure
- **WHEN** parsing a section from tasks.md
- **THEN** each section object SHALL include `name` (section title)
- **AND** SHALL include `number` (section number from header)
- **AND** SHALL include `tasks` array with task objects

#### Scenario: Task object structure
- **WHEN** parsing a task from tasks.md
- **THEN** each task object SHALL include `id` (e.g., "1.1", "2.3")
- **AND** SHALL include `description` (task text after ID, with indented detail lines appended)
- **AND** SHALL include `completed` boolean based on `[x]` vs `[ ]`
- **AND** SHALL include `subtasks` array (empty if no nested tasks, recursive for unlimited depth)

#### Scenario: Summary calculation
- **WHEN** generating tasks.json
- **THEN** summary.total SHALL equal total number of tasks across all sections (including all nested subtasks recursively)
- **AND** summary.completed SHALL equal number of tasks with `completed: true` (including nested subtasks)

### Requirement: Task ID Parsing
The accept command SHALL parse task IDs from the standard tasks.md format.

#### Scenario: Parse simple task ID
- **WHEN** parsing line `- [ ] 1.1 Create database schema`
- **THEN** task ID SHALL be "1.1"
- **AND** description SHALL be "Create database schema"
- **AND** completed SHALL be false

#### Scenario: Parse completed task
- **WHEN** parsing line `- [x] 2.3 Write unit tests`
- **THEN** task ID SHALL be "2.3"
- **AND** completed SHALL be true

#### Scenario: Parse nested subtask
- **WHEN** parsing line `- [ ] 1.1.1 Create users table`
- **THEN** task ID SHALL be "1.1.1"
- **AND** the task SHALL be added to subtasks of task "1.1"

#### Scenario: Parse deeply nested subtask
- **WHEN** parsing line `- [ ] 1.1.1.1 Add primary key constraint`
- **THEN** task ID SHALL be "1.1.1.1"
- **AND** the task SHALL be added to subtasks of task "1.1.1" (unlimited nesting depth)

#### Scenario: Handle task without ID
- **WHEN** parsing line `- [ ] Some task without numeric ID`
- **THEN** the system SHALL generate an auto-incremented ID within its section
- **AND** description SHALL be "Some task without numeric ID"

#### Scenario: Parse task with indented detail lines
- **WHEN** parsing a task followed by indented lines (2+ spaces or tab):
  ```
  - [ ] 1.1 Create database schema
    - Parse requirement headers
    - Extract requirement name and content
  ```
- **THEN** the indented lines SHALL be appended to the description with newlines
- **AND** the final description SHALL be "Create database schema\n- Parse requirement headers\n- Extract requirement name and content"

#### Scenario: Preserve completion status from tasks.md
- **WHEN** parsing a task marked `[x]` in tasks.md
- **THEN** the task SHALL have `completed: true` in tasks.json
- **AND** existing progress is preserved during acceptance

### Requirement: Section Header Parsing
The accept command SHALL parse section headers from tasks.md format.

#### Scenario: Parse numbered section header
- **WHEN** parsing line `## 1. Implementation`
- **THEN** section number SHALL be 1
- **AND** section name SHALL be "Implementation"

#### Scenario: Parse section header without number
- **WHEN** parsing line `## Validation`
- **THEN** section number SHALL be auto-assigned based on order
- **AND** section name SHALL be "Validation"

### Requirement: Accept Command Help Text
The accept command SHALL provide comprehensive help documentation.

#### Scenario: Command help display
- **WHEN** user invokes `spectr accept --help`
- **THEN** help text SHALL describe task conversion purpose
- **AND** SHALL explain that tasks.md is converted to tasks.json
- **AND** SHALL note that tasks.md is removed after conversion
- **AND** SHALL list available flags with descriptions
