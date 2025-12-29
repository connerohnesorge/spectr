# Delta Specification

## ADDED Requirements

### Requirement: Search Hotkey in Interactive Lists

The interactive list modes SHALL provide a '/' hotkey that activates a text
search mode, allowing users to filter the displayed list by typing a search
query that matches against item IDs and titles.

#### Scenario: User presses '/' to enter search mode

- **WHEN** user is in any interactive list mode (changes, specs, or unified)
- **AND** user presses the '/' key
- **THEN** search mode is activated
- **AND** a text input field is displayed below or above the table
- **AND** the cursor is placed in the text input field
- **AND** the user can type a search query

#### Scenario: Search filters rows in real-time

- **WHEN** search mode is active
- **AND** user types characters into the search input
- **THEN** the table rows are filtered in real-time
- **AND** only rows where ID or title contains the search query
  (case-insensitive) are displayed
- **AND** the first matching row is automatically selected

#### Scenario: Search with no matches shows empty table

- **WHEN** search mode is active
- **AND** user types a query that matches no items
- **THEN** the table displays no rows
- **AND** a message indicates no matches found

#### Scenario: User presses Escape to exit search mode

- **WHEN** search mode is active
- **AND** user presses the Escape key
- **THEN** search mode is deactivated
- **AND** the search query is cleared
- **AND** all items are displayed again in the table
- **AND** the text input field is hidden

#### Scenario: User presses '/' again to clear search

- **WHEN** search mode is active
- **AND** the search query is not empty
- **AND** user presses '/' key
- **THEN** the search input gains focus (normal text input behavior)

- **WHEN** search mode is active
- **AND** the search query is empty
- **AND** user presses '/' key
- **THEN** search mode is deactivated
- **AND** all items are displayed again

#### Scenario: Navigation works while searching

- **WHEN** search mode is active
- **AND** filtered results are displayed
- **THEN** arrow key navigation (up/down, j/k) moves through filtered rows
- **AND** Enter key copies the selected filtered item's ID
- **AND** other hotkeys (e, a, t) work on the selected filtered item

#### Scenario: Help text shows search hotkey

- **WHEN** interactive mode is displayed in any mode
- **THEN** the help text includes '/: search' in the controls line
- **AND** the search hotkey is shown for all modes (changes, specs, unified)

#### Scenario: Search mode visual indicator

- **WHEN** search mode is active
- **THEN** the search input field is visually distinct
- **AND** the current search query is visible
- **AND** the help text updates to show 'Esc: exit search'
