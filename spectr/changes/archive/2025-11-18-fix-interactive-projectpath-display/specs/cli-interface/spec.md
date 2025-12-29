## MODIFIED Requirements

### Requirement: Interactive List Mode

The list command SHALL provide an interactive table interface when the `-I` or `--interactive` flag is used, displaying items in a navigable table format with project path information.

#### Scenario: User launches interactive list for changes

- **WHEN** user runs `spectr list -I`
- **THEN** a table is displayed with columns: ID, Title, Deltas, Tasks
- **AND** the table supports arrow key navigation
- **AND** the first row is selected by default
- **AND** the project path is displayed in the interface

#### Scenario: User launches interactive list for specs

- **WHEN** user runs `spectr list --specs -I`
- **THEN** a table is displayed with columns: ID, Title, Requirements
- **AND** the table supports arrow key navigation
- **AND** the first row is selected by default
- **AND** the project path is displayed in the interface

#### Scenario: User navigates with keyboard

- **WHEN** user presses arrow keys or j/k
- **THEN** the selection moves up or down accordingly
- **AND** the selected row is visually highlighted

#### Scenario: Empty list in interactive mode

- **WHEN** user runs `spectr list -I` and no changes exist
- **THEN** display "No items found" message
- **AND** exit cleanly without entering interactive mode

### Requirement: Interactive Archive Mode

The archive command SHALL provide an interactive table interface when no change ID argument is provided or when the `-I` or `--interactive` flag is used, displaying available changes in a navigable table format identical to the list command's interactive mode with project path information.

#### Scenario: User runs archive with no arguments

- **WHEN** user runs `spectr archive` with no change ID argument
- **THEN** an interactive table is displayed with columns: ID, Title, Deltas, Tasks
- **AND** the table supports arrow key navigation (↑/↓, j/k)
- **AND** the first row is selected by default
- **AND** the table uses the same visual styling as list -I
- **AND** the project path is displayed in the interface

#### Scenario: User runs archive with -I flag

- **WHEN** user runs `spectr archive -I`
- **THEN** an interactive table is displayed even if other flags are present
- **AND** the behavior is identical to running archive with no arguments
- **AND** the project path is displayed in the interface

#### Scenario: User selects change for archiving

- **WHEN** user presses Enter on a selected row in archive interactive mode
- **THEN** the change ID is captured (not copied to clipboard)
- **AND** the interactive mode exits
- **AND** the archive workflow proceeds with the selected change ID
- **AND** validation, task checking, and spec updates proceed as normal

#### Scenario: User cancels archive selection

- **WHEN** user presses 'q' or Ctrl+C in archive interactive mode
- **THEN** interactive mode exits
- **AND** archive command returns successfully without archiving anything
- **AND** a "Cancelled" message is displayed

#### Scenario: No changes available for archiving

- **WHEN** user runs `spectr archive` and no changes exist in changes/ directory
- **THEN** display "No changes available to archive" message
- **AND** exit cleanly without entering interactive mode
- **AND** command returns successfully

#### Scenario: Archive with explicit change ID bypasses interactive mode

- **WHEN** user runs `spectr archive <change-id>`
- **THEN** interactive mode is NOT triggered
- **AND** archive proceeds directly with the specified change ID
- **AND** behavior is unchanged from current implementation

## ADDED Requirements

### Requirement: Project Path Display in Interactive Mode

The interactive table interfaces SHALL display the project root path to provide users with context about which project they are working with.

#### Scenario: Project path shown in changes interactive mode

- **WHEN** user runs `spectr list -I` for changes
- **THEN** the project root path is displayed in the help text or table header
- **AND** the path is the absolute path to the project directory

#### Scenario: Project path shown in specs interactive mode

- **WHEN** user runs `spectr list --specs -I`
- **THEN** the project root path is displayed in the help text or table header
- **AND** the path is the absolute path to the project directory

#### Scenario: Project path shown in archive interactive mode

- **WHEN** user runs `spectr archive` without arguments
- **THEN** the project root path is displayed in the help text or table header
- **AND** the path is the absolute path to the project directory

#### Scenario: Project path properly initialized for changes

- **WHEN** `RunInteractiveChanges()` is invoked
- **THEN** the `projectPath` parameter is passed from the calling command
- **AND** the `projectPath` field on `interactiveModel` is set during initialization

#### Scenario: Project path properly initialized for archive

- **WHEN** `RunInteractiveArchive()` is invoked
- **THEN** the `projectPath` parameter is passed from the calling command
- **AND** the `projectPath` field on `interactiveModel` is set during initialization
