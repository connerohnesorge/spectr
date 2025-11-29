## ADDED Requirements

### Requirement: Delete Hotkey for Specs in Interactive Mode

The interactive specs list mode and unified mode SHALL provide a 'd' hotkey that deletes the entire spec folder for the currently selected spec row, with a confirmation prompt before performing the destructive action.

#### Scenario: User presses 'd' to delete a spec in specs mode

- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses the 'd' key on a selected spec
- **THEN** a confirmation prompt is displayed: "Delete spec '<spec-id>'? This will remove the entire folder. (y/N)"
- **AND** the TUI waits for user input

#### Scenario: User confirms spec deletion

- **WHEN** the confirmation prompt is displayed
- **AND** user presses 'y' or 'Y'
- **THEN** the entire `spectr/specs/<spec-id>/` folder is deleted using `os.RemoveAll()`
- **AND** a success message is displayed: "Deleted: <spec-id>"
- **AND** the row is removed from the table
- **AND** the table is refreshed with updated row count
- **AND** if the deleted row was the last row, the cursor moves to the new last row

#### Scenario: User cancels spec deletion

- **WHEN** the confirmation prompt is displayed
- **AND** user presses 'n', 'N', Escape, or any other key
- **THEN** the deletion is cancelled
- **AND** a message is displayed: "Cancelled"
- **AND** the TUI returns to the interactive list view
- **AND** the same row remains selected

#### Scenario: Delete in unified mode on spec row

- **WHEN** user is in unified interactive mode (`spectr list --all -I`)
- **AND** user presses 'd' on a row where Type is "SPEC"
- **THEN** the confirmation and deletion flow proceeds as in specs mode
- **AND** the behavior is identical to the specs-only mode

#### Scenario: Delete ignored for change rows in unified mode

- **WHEN** user is in unified interactive mode (`spectr list --all -I`)
- **AND** user presses 'd' on a row where Type is "CHANGE"
- **THEN** the key press is ignored (no action taken)
- **AND** an optional error message may be displayed: "Cannot delete changes; use archive instead"
- **AND** the TUI remains in interactive mode

#### Scenario: Delete hotkey not available for changes-only mode

- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses 'd' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'd: delete' option

#### Scenario: Spec folder does not exist

- **WHEN** user presses 'd' to delete a spec
- **AND** the spec folder at `spectr/specs/<spec-id>/` does not exist
- **THEN** an error message is displayed: "Spec folder not found: <path>"
- **AND** the TUI remains in interactive mode
- **AND** the user can continue navigating or quit

#### Scenario: Deletion fails with filesystem error

- **WHEN** user confirms deletion
- **AND** `os.RemoveAll()` fails (e.g., permission denied, read-only filesystem)
- **THEN** an error message is displayed with the underlying error details
- **AND** the TUI remains in interactive mode
- **AND** the row is NOT removed from the table

#### Scenario: Help text shows delete hotkey in specs mode

- **WHEN** user presses '?' in interactive specs mode
- **THEN** the help text includes 'd: delete' in the controls line
- **AND** the hotkey appears in the list of available actions

#### Scenario: Help text shows delete hotkey in unified mode

- **WHEN** user presses '?' in unified interactive mode
- **THEN** the help text includes 'd: delete (specs)' in the controls line
- **AND** the hotkey description clarifies it applies only to spec rows
