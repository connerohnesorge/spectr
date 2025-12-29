# CLI Interface Specification

## ADDED Requirements

### Requirement: Interactive List Mode

The list command SHALL provide an interactive table interface when the `-I` or `--interactive` flag is used, displaying items in a navigable table format.

#### Scenario: User launches interactive list for changes

- **WHEN** user runs `spectr list -I`
- **THEN** a table is displayed with columns: ID, Title, Deltas, Tasks
- **AND** the table supports arrow key navigation
- **AND** the first row is selected by default

#### Scenario: User launches interactive list for specs

- **WHEN** user runs `spectr list --specs -I`
- **THEN** a table is displayed with columns: ID, Title, Requirements
- **AND** the table supports arrow key navigation
- **AND** the first row is selected by default

#### Scenario: User navigates with keyboard

- **WHEN** user presses arrow keys or j/k
- **THEN** the selection moves up or down accordingly
- **AND** the selected row is visually highlighted

#### Scenario: Empty list in interactive mode

- **WHEN** user runs `spectr list -I` and no changes exist
- **THEN** display "No items found" message
- **AND** exit cleanly without entering interactive mode

### Requirement: Clipboard Copy on Selection

When a user presses Enter on a selected row in interactive mode, the item's ID SHALL be copied to the system clipboard.

#### Scenario: Copy change ID to clipboard

- **WHEN** user selects a change row and presses Enter
- **THEN** the change ID (kebab-case identifier) is copied to clipboard
- **AND** a success message is displayed (e.g., "Copied: add-archive-command")
- **AND** the interactive mode exits

#### Scenario: Copy spec ID to clipboard

- **WHEN** user selects a spec row and presses Enter
- **THEN** the spec ID is copied to clipboard
- **AND** a success message is displayed
- **AND** the interactive mode exits

#### Scenario: Clipboard failure handling

- **WHEN** clipboard operation fails
- **THEN** display error message to user
- **AND** do not exit interactive mode
- **AND** user can retry or quit manually

### Requirement: Interactive Mode Exit Controls

Users SHALL be able to exit interactive mode using standard quit commands.

#### Scenario: Quit with q key

- **WHEN** user presses 'q'
- **THEN** interactive mode exits
- **AND** no clipboard operation occurs
- **AND** command returns successfully

#### Scenario: Quit with Ctrl+C

- **WHEN** user presses Ctrl+C
- **THEN** interactive mode exits immediately
- **AND** no clipboard operation occurs
- **AND** command returns successfully

### Requirement: Table Visual Styling

The interactive table SHALL use clear visual styling to distinguish headers, selected rows, and borders.

#### Scenario: Visual hierarchy in table

- **WHEN** interactive mode is displayed
- **THEN** column headers are visually distinct from data rows
- **AND** selected row has contrasting background/foreground colors
- **AND** table borders are visible and styled consistently
- **AND** table fits within terminal width gracefully

### Requirement: Cross-Platform Clipboard Support

Clipboard operations SHALL work across Linux, macOS, and Windows platforms.

#### Scenario: Clipboard on Linux

- **WHEN** running on Linux
- **THEN** clipboard operations use X11 or Wayland clipboard APIs as appropriate
- **AND** fallback to OSC 52 escape sequences if desktop clipboard unavailable

#### Scenario: Clipboard on macOS

- **WHEN** running on macOS
- **THEN** clipboard operations use pbcopy or native clipboard APIs

#### Scenario: Clipboard on Windows

- **WHEN** running on Windows
- **THEN** clipboard operations use Windows clipboard APIs

#### Scenario: Clipboard in SSH/remote session

- **WHEN** running over SSH without X11 forwarding
- **THEN** use OSC 52 escape sequences to copy to local clipboard
- **AND** document this behavior for users
