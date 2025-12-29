## MODIFIED Requirements

### Requirement: Task Completion Checking

The system SHALL check task completion status and warn users before archiving. The system SHALL read from `tasks.json` when present, falling back to `tasks.md`.

#### Scenario: Display task status from JSON

- **WHEN** archiving a change with `tasks.json`
- **THEN** the system reads task status from JSON file
- **AND** displays task completion status (e.g., "3/5 complete")

#### Scenario: Display task status from Markdown

- **WHEN** archiving a change with only `tasks.md`
- **THEN** the system reads task status from Markdown file
- **AND** displays task completion status (e.g., "3/5 complete")

#### Scenario: Warn on incomplete tasks

- **WHEN** a change has incomplete tasks
- **THEN** the system warns the user and requires confirmation to proceed (unless --yes flag is provided)

#### Scenario: Proceed with incomplete tasks when confirmed

- **WHEN** user confirms archiving despite incomplete tasks
- **THEN** the system proceeds with the archive operation

## ADDED Requirements

### Requirement: Auto-Accept on Archive

The system SHALL automatically convert `tasks.md` to `tasks.json` during archive if not already accepted, ensuring archived changes have stable task format.

#### Scenario: Archive triggers auto-accept

- **WHEN** archiving a change that has `tasks.md` but no `tasks.json`
- **THEN** the system displays a warning that auto-acceptance will occur
- **AND** the system converts `tasks.md` to `tasks.json` before archiving
- **AND** the system removes `tasks.md` after successful conversion

#### Scenario: Archive with existing tasks.json

- **WHEN** archiving a change that already has `tasks.json`
- **THEN** the system proceeds normally without conversion
- **AND** the archived change contains `tasks.json`

#### Scenario: Auto-accept failure blocks archive

- **WHEN** auto-acceptance fails during archive (e.g., invalid tasks.md format)
- **THEN** the system displays the conversion error
- **AND** the system aborts the archive operation
- **AND** no files are modified
