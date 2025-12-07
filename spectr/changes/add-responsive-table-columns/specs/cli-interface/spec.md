## ADDED Requirements

### Requirement: Responsive Table Column Layout
The interactive TUI table views SHALL detect terminal width and dynamically adjust column visibility and widths to ensure readable display across different screen sizes.

#### Scenario: Full width terminal displays all columns
- **WHEN** user runs `spectr list -I` on a terminal with 110+ columns
- **THEN** all columns are displayed at their default widths
- **AND** for changes view: ID (30), Title (40), Deltas (10), Tasks (15) are shown
- **AND** for specs view: ID (35), Title (45), Requirements (15) are shown
- **AND** for unified view: ID (30), Type (8), Title (40), Details (20) are shown

#### Scenario: Medium width terminal narrows Title column
- **WHEN** user runs `spectr list -I` on a terminal with 90-109 columns
- **THEN** the Title column width is reduced proportionally
- **AND** title truncation threshold is reduced to match narrower column
- **AND** all columns remain visible

#### Scenario: Narrow width terminal hides low-priority columns
- **WHEN** user runs `spectr list -I` on a terminal with 70-89 columns
- **THEN** the lowest-priority columns are hidden
- **AND** for changes view: Tasks column is hidden, Deltas may be narrowed
- **AND** for specs view: Requirements column width is reduced
- **AND** for unified view: Details column is hidden
- **AND** remaining columns are adjusted to fit available width

#### Scenario: Minimal width terminal shows essential columns only
- **WHEN** user runs `spectr list -I` on a terminal with fewer than 70 columns
- **THEN** only ID and Title columns are displayed
- **AND** title truncation is aggressive to fit available space
- **AND** help text indicates some columns are hidden

### Requirement: Dynamic Terminal Resize Handling
The interactive TUI SHALL respond to terminal resize events by recalculating and rebuilding the table layout without losing user state.

#### Scenario: Terminal resized wider during session
- **WHEN** user is in interactive mode and terminal width increases
- **THEN** table columns are recalculated for the new width
- **AND** previously hidden columns may become visible
- **AND** cursor position is preserved on the same item
- **AND** search filter state is preserved if active

#### Scenario: Terminal resized narrower during session
- **WHEN** user is in interactive mode and terminal width decreases
- **THEN** table columns are recalculated for the new width
- **AND** low-priority columns are hidden as needed
- **AND** cursor position is preserved on the same item
- **AND** the view does not overflow horizontally

#### Scenario: Resize does not interrupt search mode
- **WHEN** user is in search mode and terminal is resized
- **THEN** search input remains active
- **AND** filtered results are preserved
- **AND** table layout adapts to new width

### Requirement: Column Priority System
Each table view SHALL define column priorities to determine which columns are shown at each width breakpoint.

#### Scenario: Changes view column priorities
- **WHEN** calculating responsive columns for changes view
- **THEN** ID has highest priority (always shown)
- **AND** Title has second priority (always shown, width adjustable)
- **AND** Deltas has third priority (hidden below 80 columns)
- **AND** Tasks has lowest priority (hidden below 90 columns)

#### Scenario: Specs view column priorities
- **WHEN** calculating responsive columns for specs view
- **THEN** ID has highest priority (always shown)
- **AND** Title has second priority (always shown, width adjustable)
- **AND** Requirements has lowest priority (width reduced or hidden below 70 columns)

#### Scenario: Unified view column priorities
- **WHEN** calculating responsive columns for unified view
- **THEN** ID has highest priority (always shown)
- **AND** Type has second priority (always shown at fixed 8-character width)
- **AND** Title has third priority (width adjustable)
- **AND** Details has lowest priority (hidden below 90 columns)
