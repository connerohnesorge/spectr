## ADDED Requirements

### Requirement: Shared TUI Component Library

The CLI SHALL use a shared `internal/tui` package for interactive TUI components, providing consistent styling, behavior, and composable building blocks across all interactive modes.

#### Scenario: TablePicker used for item selection

- **WHEN** any command needs an interactive table-based selection (list, archive, validation item picker)
- **THEN** the command SHALL use the `TablePicker` component from `internal/tui`
- **AND** the table SHALL use consistent styling from `tui.ApplyTableStyles()`
- **AND** navigation keys (↑/↓, j/k) SHALL work identically across all usages
- **AND** quit keys (q, Ctrl+C) SHALL work identically across all usages

#### Scenario: MenuPicker used for option selection

- **WHEN** any command needs an interactive menu selection (validation mode menu)
- **THEN** the command SHALL use the `MenuPicker` component from `internal/tui`
- **AND** the menu SHALL use consistent styling
- **AND** navigation and selection behavior SHALL match the TablePicker patterns

#### Scenario: Consistent string truncation

- **WHEN** any TUI component needs to truncate text for display
- **THEN** it SHALL use `tui.TruncateString()` with consistent ellipsis handling
- **AND** truncation SHALL add "..." suffix when text exceeds max length
- **AND** very short max lengths (≤3) SHALL truncate without ellipsis

#### Scenario: Consistent clipboard operations

- **WHEN** any TUI component needs to copy text to clipboard
- **THEN** it SHALL use `tui.CopyToClipboard()` from the shared package
- **AND** the function SHALL try native clipboard first
- **AND** the function SHALL fall back to OSC 52 for remote sessions

#### Scenario: Action registration pattern

- **WHEN** a command configures a TablePicker with custom actions
- **THEN** actions SHALL be registered via `WithAction(key, label, handler)`
- **AND** the help text SHALL automatically include all registered actions
- **AND** unregistered keys SHALL be ignored (no error)

#### Scenario: Domain logic remains in consuming packages

- **WHEN** the tui package is used by list or validation
- **THEN** domain-specific logic (archive workflow, validation execution) SHALL remain in consuming packages
- **AND** the tui package SHALL only provide UI primitives
- **AND** business logic SHALL not be coupled to the tui package

## MODIFIED Requirements

### Requirement: Table Visual Styling

The interactive table SHALL use clear visual styling to distinguish headers, selected rows, and borders, provided by the shared `internal/tui` package.

#### Scenario: Visual hierarchy in table

- **WHEN** interactive mode is displayed
- **THEN** column headers are visually distinct from data rows
- **AND** selected row has contrasting background/foreground colors
- **AND** table borders are visible and styled consistently
- **AND** table fits within terminal width gracefully
- **AND** styling SHALL be applied via `tui.ApplyTableStyles()`

#### Scenario: Consistent styling across commands

- **WHEN** user uses `spectr list -I`, `spectr archive`, or `spectr validate` interactive modes
- **THEN** all tables SHALL use identical styling
- **AND** colors, borders, and highlights SHALL match exactly
- **AND** the shared `tui.ApplyTableStyles()` function SHALL be the single source of truth

## MODIFIED Requirements

### Requirement: Interactive Archive Mode

The archive command SHALL provide an interactive table interface when no change ID argument is provided, displaying available changes in a navigable table format identical to the list command's interactive mode with project path information. The `-I`/`--interactive` flag has been removed as TUI is now the default behavior when no change ID is provided.

#### Scenario: User runs archive with no arguments

- **WHEN** user runs `spectr archive` with no change ID argument
- **THEN** an interactive table is displayed with columns: ID, Title, Deltas, Tasks
- **AND** the table supports arrow key navigation (↑/↓, j/k)
- **AND** the first row is selected by default
- **AND** the table uses the same visual styling as list -I
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
