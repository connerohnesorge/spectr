## ADDED Requirements

### Requirement: Accept Command Structure

The CLI SHALL provide an `accept` command that converts `tasks.md` to `tasks.json` format for stable agent manipulation during implementation.

#### Scenario: Accept command registration

- **WHEN** the CLI is initialized
- **THEN** it SHALL include an AcceptCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr accept`
- **AND** help text SHALL describe task format conversion functionality

#### Scenario: Accept with change ID

- **WHEN** user runs `spectr accept <change-id>`
- **THEN** the system validates the change exists in `spectr/changes/<change-id>/`
- **AND** the system parses `tasks.md` into structured format
- **AND** the system writes `tasks.json` with proper schema
- **AND** the system removes `tasks.md` to prevent drift

#### Scenario: Accept with validation

- **WHEN** user runs `spectr accept <change-id>`
- **THEN** the system validates the change before conversion
- **AND** the system blocks acceptance if validation fails
- **AND** the system displays validation errors

#### Scenario: Accept dry-run mode

- **WHEN** user runs `spectr accept <change-id> --dry-run`
- **THEN** the system displays what would be converted
- **AND** the system does NOT write tasks.json
- **AND** the system does NOT remove tasks.md

#### Scenario: Accept already accepted change

- **WHEN** user runs `spectr accept <change-id>` on a change that already has tasks.json
- **THEN** the system displays a message indicating change is already accepted
- **AND** the system exits with code 0 (success, idempotent)

#### Scenario: Accept change without tasks.md

- **WHEN** user runs `spectr accept <change-id>` on a change without tasks.md
- **THEN** the system displays an error indicating no tasks.md found
- **AND** the system exits with code 1

### Requirement: Tasks JSON Schema

The accept command SHALL generate `tasks.json` files conforming to a versioned schema with structured task objects.

#### Scenario: JSON file structure

- **WHEN** the accept command creates tasks.json
- **THEN** the file SHALL contain a root object with `version` and `tasks` fields
- **AND** `version` SHALL be integer 1 for this schema version
- **AND** `tasks` SHALL be an array of task objects

#### Scenario: Task object structure

- **WHEN** a task is serialized to JSON
- **THEN** it SHALL have `id` field containing the task identifier (e.g., "1.1")
- **AND** it SHALL have `section` field containing the section header (e.g., "Implementation")
- **AND** it SHALL have `description` field containing the full task text
- **AND** it SHALL have `status` field with value "pending", "in_progress", or "completed"

#### Scenario: Status value mapping from Markdown

- **WHEN** converting tasks.md to tasks.json
- **THEN** `- [ ]` SHALL map to status "pending"
- **AND** `- [x]` (case-insensitive) SHALL map to status "completed"

### Requirement: Accept Command Flags

The accept command SHALL support flags for controlling behavior.

#### Scenario: Dry-run flag

- **WHEN** user provides the `--dry-run` flag
- **THEN** the system previews the conversion without writing files
- **AND** displays the JSON that would be generated

#### Scenario: Interactive change selection

- **WHEN** user runs `spectr accept` without specifying a change ID
- **THEN** the system displays a list of active changes with tasks.md files
- **AND** prompts for selection using existing TUI components

## MODIFIED Requirements

### Requirement: Task Counting

The system SHALL count tasks in `tasks.md` files by identifying lines matching the pattern `- [ ]` or `- [x]` (case-insensitive), with completed tasks marked by `[x]`. The system SHALL prefer `tasks.json` when present.

#### Scenario: Count tasks from JSON

- **WHEN** the system counts tasks and `tasks.json` exists
- **THEN** it reads task status from the JSON file
- **AND** counts tasks by status field values
- **AND** reports `taskStatus` with total and completed counts

#### Scenario: Count completed and total tasks

- **WHEN** the system reads a `tasks.md` file with 3 tasks, 2 marked `[x]` and 1 marked `[ ]`
- **THEN** it reports `taskStatus` as `{ total: 3, completed: 2 }`

#### Scenario: Handle missing tasks file

- **WHEN** the system cannot find or read a `tasks.md` or `tasks.json` file for a change
- **THEN** it reports `taskStatus` as `{ total: 0, completed: 0 }`
- **AND** continues processing without error

#### Scenario: JSON takes precedence over Markdown

- **WHEN** both `tasks.json` and `tasks.md` exist (should not happen normally)
- **THEN** the system reads from `tasks.json`
- **AND** ignores `tasks.md`
