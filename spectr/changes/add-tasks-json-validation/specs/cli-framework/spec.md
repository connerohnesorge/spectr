## ADDED Requirements

### Requirement: Tasks JSON Validation
The validate command SHALL validate `tasks.json` files when present in a change directory, ensuring JSON is well-formed and all task status values conform to the `TaskStatusValue` enum.

#### Scenario: Valid tasks.json passes validation
- **WHEN** `spectr validate <change-id>` is run
- **AND** the change directory contains a valid `tasks.json` file
- **AND** all task status values are one of: `pending`, `in_progress`, `completed`
- **THEN** validation SHALL pass without tasks.json-related errors

#### Scenario: Invalid status value reports error
- **WHEN** `spectr validate <change-id>` is run
- **AND** the change directory contains a `tasks.json` file
- **AND** a task has an invalid status value (e.g., `"done"`, `"unknown"`)
- **THEN** validation SHALL report an error for each invalid task
- **AND** the error message SHALL include the task ID
- **AND** the error message SHALL specify the invalid value and list valid options

#### Scenario: Malformed JSON reports error
- **WHEN** `spectr validate <change-id>` is run
- **AND** the change directory contains a `tasks.json` file that is not valid JSON
- **THEN** validation SHALL report a JSON parsing error
- **AND** validation SHALL fail

#### Scenario: Missing tasks.json is not an error
- **WHEN** `spectr validate <change-id>` is run
- **AND** the change directory does not contain a `tasks.json` file
- **THEN** validation SHALL NOT report an error for missing tasks.json
- **AND** validation SHALL continue checking other files
