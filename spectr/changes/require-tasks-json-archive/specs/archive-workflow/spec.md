## REMOVED Requirements

### Requirement: Auto-Accept on Archive
**Reason**: Auto-accepting during archive bypasses the intentional workflow gate of running `spectr accept` before implementation. Archive should only archive completed changes, not handle format conversion.
**Migration**: Users must run `spectr accept <change-id>` before archiving. The archive command will display an actionable error message guiding users to do this.

## ADDED Requirements

### Requirement: Require Accepted Change for Archive
The system SHALL require that a change has been accepted (has `tasks.json`) before it can be archived.

#### Scenario: Archive accepted change
- **WHEN** archiving a change that has `tasks.json`
- **THEN** the system proceeds with the archive operation normally

#### Scenario: Archive unaccepted change with tasks.md
- **WHEN** archiving a change that has `tasks.md` but no `tasks.json`
- **THEN** the system displays an error: "No tasks.json found. Run `spectr accept <change-id>` to accept the change first."
- **AND** the system exits with non-zero status code
- **AND** no files are modified

#### Scenario: Archive change with no task files
- **WHEN** archiving a change that has neither `tasks.md` nor `tasks.json`
- **THEN** the system displays an error indicating no task file was found
- **AND** the system exits with non-zero status code
- **AND** no files are modified

## MODIFIED Requirements

### Requirement: Task Completion Checking
The system SHALL check task completion status and warn users before archiving. The system SHALL read from `tasks.json` (required).

#### Scenario: Display task status from JSON
- **WHEN** archiving a change with `tasks.json`
- **THEN** the system reads task status from JSON file
- **AND** displays task completion status (e.g., "3/5 complete")

#### Scenario: Warn on incomplete tasks
- **WHEN** a change has incomplete tasks
- **THEN** the system warns the user and requires confirmation to proceed (unless --yes flag is provided)

#### Scenario: Proceed with incomplete tasks when confirmed
- **WHEN** user confirms archiving despite incomplete tasks
- **THEN** the system proceeds with the archive operation
