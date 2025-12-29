# CLI Interface Delta Spec

## ADDED Requirements

### Requirement: Interactive Archive Mode

The archive command SHALL provide an interactive table interface when no change ID argument is provided or when the `-I` or `--interactive` flag is used, displaying available changes in a navigable table format identical to the list command's interactive mode.

#### Scenario: User runs archive with no arguments

- **WHEN** user runs `spectr archive` with no change ID argument
- **THEN** an interactive table is displayed with columns: ID, Title, Deltas, Tasks
- **AND** the table supports arrow key navigation (↑/↓, j/k)
- **AND** the first row is selected by default
- **AND** the table uses the same visual styling as list -I

#### Scenario: User runs archive with -I flag

- **WHEN** user runs `spectr archive -I`
- **THEN** an interactive table is displayed even if other flags are present
- **AND** the behavior is identical to running archive with no arguments

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

### Requirement: Archive Interactive Table Display

The archive command's interactive table SHALL display the same information columns as the list command to help users make informed archiving decisions.

#### Scenario: Table columns match list command

- **WHEN** archive interactive mode is displayed
- **THEN** columns are: ID (30 chars), Title (40 chars), Deltas (10 chars), Tasks (15 chars)
- **AND** column widths match the list -I command exactly
- **AND** title text is truncated with ellipsis if longer than 38 characters
- **AND** task status shows format "completed/total" (e.g., "5/10")

#### Scenario: Visual styling consistency

- **WHEN** archive interactive table is displayed
- **THEN** the table uses identical styling to list -I
- **AND** column headers are visually distinct from data rows
- **AND** selected row has contrasting background/foreground colors
- **AND** table borders are visible and styled consistently
- **AND** help text shows navigation controls (↑/↓, j/k, enter, q)

### Requirement: Archive Selection Without Clipboard

The archive command's interactive mode SHALL NOT copy the selected change ID to the clipboard, unlike the list command, since the ID is immediately consumed by the archive workflow.

#### Scenario: Enter key captures selection

- **WHEN** user presses Enter on a selected change
- **THEN** the change ID is captured internally
- **AND** NO clipboard operation occurs
- **AND** NO "Copied: <id>" message is displayed
- **AND** the archive workflow proceeds immediately with the selected ID

#### Scenario: Workflow continuation

- **WHEN** a change is selected in interactive mode
- **THEN** the Archiver.Archive() method receives the selected change ID
- **AND** validation, task checking, and spec updates proceed as if the ID was provided as an argument
- **AND** all confirmation prompts and flags (--yes, --skip-specs) work normally
