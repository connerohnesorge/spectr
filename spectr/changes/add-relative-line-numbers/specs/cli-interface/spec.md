# Cli Interface Specification - Delta

## ADDED Requirements

### Requirement: Relative Line Numbers in Interactive List

The system SHALL support optional line number display in interactive list mode,
with relative and hybrid modes to assist vim-style count-prefix navigation.

#### Scenario: Toggle line number display with # key

- WHEN user presses `#` in interactive list mode
- THEN the system SHALL cycle through modes: off -> relative -> hybrid -> off
- AND the table display SHALL update immediately to reflect the new mode

#### Scenario: Relative line number mode display

- WHEN line number mode is `relative`
- THEN each row SHALL show its distance from the cursor row
- AND the cursor row SHALL display `0`
- AND rows above cursor SHALL show `1`, `2`, `3`... (ascending distance)
- AND rows below cursor SHALL show `1`, `2`, `3`... (ascending distance)

#### Scenario: Hybrid line number mode display

- WHEN line number mode is `hybrid`
- THEN the cursor row SHALL display its absolute position (1-indexed)
- AND all other rows SHALL show relative distance from cursor
- AND this matches Vim's `set number relativenumber` behavior

#### Scenario: Line number column styling

- WHEN line numbers are displayed
- THEN the line number column SHALL be dimmed (gray) for non-cursor rows
- AND the cursor row's line number SHALL be brighter or bold
- AND the column SHALL be right-aligned with consistent width

#### Scenario: Footer indicator for line number mode

- WHEN line number mode is not off
- THEN the footer SHALL include a mode indicator (`ln: rel` or `ln: hyb`)
- AND the indicator SHALL be removed when mode returns to off

#### Scenario: Default line number mode

- WHEN interactive list mode starts
- THEN line number mode SHALL default to `off`
- AND no line number column SHALL be displayed
- AND existing behavior SHALL be unchanged

#### Scenario: Help text includes line number toggle

- WHEN user presses `?` to view help
- THEN help text SHALL include `#: line numbers` in the key listing
