# Delta Specification

## ADDED Requirements

### Requirement: Provider Search in Init Wizard

The initialization wizard's tool selection step SHALL provide a `/` hotkey that
activates a text search mode, allowing users to filter the displayed provider
list by typing a search query that matches against provider names.

#### Scenario: User presses '/' to enter search mode

- **WHEN** user is on the tool selection step of the init wizard (`StepSelect`)
- **AND** user presses the '/' key
- **THEN** search mode is activated
- **AND** a text input field is displayed below the provider list
- **AND** the cursor is placed in the text input field
- **AND** the user can type a search query

#### Scenario: Search filters providers in real-time

- **WHEN** search mode is active
- **AND** user types characters into the search input
- **THEN** the provider list is filtered in real-time
- **AND** only providers whose name contains the search query (case-insensitive)
  are displayed
- **AND** the cursor moves to the first matching provider if current selection
  is filtered out

#### Scenario: Search with no matches shows empty list

- **WHEN** search mode is active
- **AND** user types a query that matches no providers
- **THEN** the provider list displays no items
- **AND** a message indicates no matches found (e.g., "No providers match
  'xyz'")

#### Scenario: User presses Escape to exit search mode

- **WHEN** search mode is active
- **AND** user presses the Escape key
- **THEN** search mode is deactivated
- **AND** the search query is cleared
- **AND** all providers are displayed again in the list
- **AND** the text input field is hidden

#### Scenario: Selection preserved during filtering

- **WHEN** search mode is active
- **AND** user has previously selected providers (checked checkboxes)
- **AND** user types a query that filters out some selected providers
- **THEN** the selection state of filtered-out providers is preserved
- **AND** when search is cleared, previously selected providers remain selected

#### Scenario: Navigation works while searching

- **WHEN** search mode is active
- **AND** filtered results are displayed
- **THEN** arrow key navigation (up/down, j/k) moves through filtered rows
- **AND** space key toggles selection on the currently highlighted filtered
  provider
- **AND** Enter key proceeds to the Review step with all selections (including
  filtered-out ones)

#### Scenario: Help text shows search hotkey

- **WHEN** the tool selection step is displayed and search mode is NOT active
- **THEN** the help text includes '/: search' in the controls line
- **AND** the search hotkey is shown alongside existing controls (navigate,
  toggle, all, none, enter, quit)

#### Scenario: Search mode visual indicator

- **WHEN** search mode is active
- **THEN** the search input field is visually distinct (styled text input)
- **AND** the current search query is visible in the input field
- **AND** the help text updates to show 'Esc: exit search' instead of '/:
  search'

#### Scenario: Cursor adjustment on filter change

- **WHEN** search mode is active
- **AND** the user types additional characters that reduce the filtered list
- **AND** the current cursor position is beyond the new list length
- **THEN** the cursor is adjusted to the last valid position in the filtered
  list
- **AND** the cursor does not go out of bounds
