## MODIFIED Requirements

### Requirement: Interactive Validation Mode

The validation system SHALL support interactive selection when invoked without arguments in a TTY, using a bubbletea-based TUI with menu-driven navigation and item picker.

#### Scenario: Interactive mode main menu

- **WHEN** validate command is invoked without arguments in an interactive terminal
- **THEN** it SHALL display a menu with options: "Validate All", "Validate All Changes", "Validate All Specs", "Pick Specific Item", "Quit"
- **AND** user SHALL be able to navigate options using arrow keys or j/k
- **AND** user SHALL select an option by pressing Enter
- **AND** selected option SHALL be executed immediately

#### Scenario: Pick specific item with search

- **WHEN** user selects "Pick Specific Item" from the main menu
- **THEN** a searchable list of all changes and specs SHALL be displayed
- **AND** items SHALL be sorted alphabetically with type indicator (change/spec)
- **AND** user SHALL navigate the list with arrow keys or j/k
- **AND** pressing Enter on an item SHALL validate that specific item
- **AND** pressing q or Ctrl+C SHALL return to main menu

#### Scenario: Non-interactive environment detection

- **WHEN** validate command is invoked without arguments in non-interactive environment (CI/CD)
- **THEN** it SHALL print usage hints for non-interactive invocation
- **AND** SHALL exit with code 1
- **AND** SHALL NOT hang waiting for input

#### Scenario: Interactive validation execution

- **WHEN** user selects a validation option in interactive mode
- **THEN** validation SHALL execute using existing validation logic
- **AND** results SHALL be displayed in human-readable format (not JSON)
- **AND** user SHALL see success/failure summary
- **AND** for failures, detailed issues SHALL be shown
- **AND** user SHALL be returned to main menu after viewing results (or exit on quit)

#### Scenario: Consistent styling with other TUIs

- **WHEN** interactive validation TUI is displayed
- **THEN** it SHALL use lipgloss styling consistent with internal/list/interactive.go
- **AND** it SHALL use the same color scheme and formatting patterns
- **AND** help text SHALL be displayed showing available key bindings
- **AND** selected items SHALL be highlighted with cursor style

## ADDED Requirements

### Requirement: Helper Functions in Internal Package

The validation system SHALL organize helper functions in internal/validation/ package following clean architecture patterns, with cmd/ serving as a thin command layer.

#### Scenario: Helper functions accessible to internal packages

- **WHEN** validation logic needs to determine item types or format results
- **THEN** helper functions SHALL be available in internal/validation/helpers.go
- **AND** item collection functions SHALL be in internal/validation/items.go
- **AND** formatting functions SHALL be in internal/validation/formatters.go
- **AND** all helpers SHALL be unit tested in their respective test files

#### Scenario: Type determination logic reusable

- **WHEN** any validation component needs to determine if an item is a change or spec
- **THEN** it SHALL use DetermineItemType() from internal/validation/helpers.go
- **AND** the function SHALL return itemTypeInfo with isChange, isSpec, and itemType fields
- **AND** it SHALL handle ambiguous cases (item exists as both change and spec)
- **AND** it SHALL respect explicit --type flag when provided

#### Scenario: Validation item collection reusable

- **WHEN** bulk validation needs to collect items to validate
- **THEN** it SHALL use GetAllItems(), GetChangeItems(), or GetSpecItems() from internal/validation/items.go
- **AND** functions SHALL return []ValidationItem with name, itemType, and path
- **AND** functions SHALL handle missing directories gracefully
- **AND** functions SHALL leverage internal/discovery for ID enumeration

#### Scenario: Result formatting separated from command logic

- **WHEN** validation results need to be displayed
- **THEN** formatting functions SHALL be in internal/validation/formatters.go
- **AND** FormatJSONReport() and FormatHumanReport() SHALL handle single item results
- **AND** FormatBulkJSONResults() and FormatBulkHumanResults() SHALL handle multiple items
- **AND** formatters SHALL accept report data and return formatted strings
- **AND** formatters SHALL not directly write to stdout (return strings instead)

### Requirement: Interactive TUI Architecture

The validation interactive mode SHALL follow the bubbletea model-update-view pattern used in other project TUIs, ensuring consistency and maintainability.

#### Scenario: Bubbletea model structure

- **WHEN** interactive validation TUI is initialized
- **THEN** it SHALL define a model struct implementing tea.Model interface
- **AND** model SHALL contain state for current screen, selected option, validation results
- **AND** model SHALL have Init() method returning initial commands
- **AND** model SHALL have Update(msg tea.Msg) method handling events
- **AND** model SHALL have View() string method rendering current state

#### Scenario: Key binding handling

- **WHEN** user presses keys in interactive validation TUI
- **THEN** arrow up/down and j/k SHALL navigate menu items
- **AND** Enter SHALL select the highlighted option
- **AND** q or Ctrl+C SHALL quit the TUI
- **AND** Esc SHALL go back to previous screen (when in item picker)
- **AND** unrecognized keys SHALL be ignored

#### Scenario: Integration with cmd layer

- **WHEN** cmd/validate.go needs to launch interactive mode
- **THEN** it SHALL call RunInteractiveValidation() from internal/validation/interactive.go
- **AND** function SHALL accept projectPath, validator, and JSON flag
- **AND** function SHALL return error on failure
- **AND** function SHALL handle TTY detection internally
- **AND** cmd layer SHALL only handle program initialization and error printing
