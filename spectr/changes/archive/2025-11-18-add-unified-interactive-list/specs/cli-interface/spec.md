# CLI Interface Spec Delta: Unified Interactive List Mode

## ADDED Requirements

### Requirement: Unified Item List Display

The system SHALL display changes and specifications together in a single
interactive table when invoked with appropriate flags, allowing users to browse
both item types simultaneously with clear visual differentiation.

#### Scenario: User opens unified interactive list

- **WHEN** the user runs `spectr list --interactive --all` from a directory with
  both changes and specs
- **THEN** a table appears showing both changes and specs rows
- **AND** each row indicates its type (change or spec)
- **AND** the table maintains correct ordering and alignment

#### Scenario: Unified list shows correct columns

- **WHEN** the unified interactive mode is active
- **THEN** the table displays: Type, ID, Title, and Type-Specific Details
  columns
- **AND** "Type-Specific Details" shows "Deltas/Tasks" for changes
- **AND** "Type-Specific Details" shows "Requirements" for specs

#### Scenario: User navigates mixed items

- **WHEN** the user navigates with arrow keys through a mixed list
- **THEN** the cursor moves smoothly between change and spec rows
- **AND** help text remains accurate and updated
- **AND** the selected row is clearly highlighted

### Requirement: Type-Aware Item Selection

The system SHALL track whether a selected item is a change or spec and provide
type-appropriate actions (e.g., edit only works for specs).

#### Scenario: Selecting a spec in unified mode

- **WHEN** the user presses Enter on a spec row
- **THEN** the spec ID is copied to clipboard
- **AND** a success message displays the ID and type indicator
- **AND** the user is returned to the interactive session or exited cleanly

#### Scenario: Selecting a change in unified mode

- **WHEN** the user presses Enter on a change row
- **THEN** the change ID is copied to clipboard
- **AND** a success message displays the ID and type indicator
- **AND** no edit action is attempted

#### Scenario: Edit action restricted to specs

- **WHEN** the user presses 'e' on a change row in unified mode
- **THEN** the action is ignored or a helpful message appears
- **AND** the interactive session continues

### Requirement: Backward-Compatible Separate Modes

The system SHALL maintain existing interactive modes for changes-only and
specs-only when `--all` flag is not provided.

#### Scenario: Changes-only mode still works

- **WHEN** the user runs `spectr list --interactive` without `--all`
- **THEN** only changes are displayed
- **AND** behavior is identical to the previous implementation
- **AND** edit functionality works as before for changes

#### Scenario: Specs-only mode still works

- **WHEN** the user runs `spectr list --specs --interactive` without `--all`
- **THEN** only specs are displayed
- **AND** behavior is identical to the previous implementation
- **AND** edit functionality works as before for specs

### Requirement: Enhanced List Command Flags

The system SHALL support new flag combinations to control listing behavior while
maintaining validation for mutually exclusive options.

#### Scenario: Flag validation for unified mode

- **WHEN** the user attempts `spectr list --interactive --all --json`
- **THEN** an error message is returned: "cannot use --interactive with --json"
- **AND** the command exits without running

#### Scenario: All flag with separate type flags

- **WHEN** the user provides `--all` with `--specs`
- **THEN** `--all` takes precedence and unified mode is used
- **AND** a warning may be shown (optional) about the redundant flag

#### Scenario: All flag in non-interactive mode

- **WHEN** the user runs `spectr list --all` without `--interactive`
- **THEN** both changes and specs are listed in text format
- **AND** each item shows its type in the output

## MODIFIED Requirements

### Requirement: Interactive List Mode

The interactive list mode in `spectr list` is extended to support unified
display of changes and specifications alongside existing separate modes.

#### Previous behavior

The system displays either changes OR specs in interactive mode based on the
`--specs` flag. Columns and behavior are specific to each item type.

#### New behavior

- When `--all` is provided with `--interactive`, both changes and specs are
  shown together with unified columns
- When neither `--all` nor `--specs` are provided, changes-only mode is default
  (backward compatible)
- When `--specs` is provided without `--all`, specs-only mode is used (backward
  compatible)
- Each item type is clearly labeled in the Type column (CHANGE or SPEC)
- Type-aware actions apply based on selected item (edit only for specs)

#### Scenario: Default behavior unchanged

- **WHEN** the user runs `spectr list --interactive`
- **THEN** the behavior is identical to before this change
- **AND** only changes are displayed
- **AND** columns show: ID, Title, Deltas, Tasks

#### Scenario: Unified mode opt-in

- **WHEN** the user explicitly uses `--all --interactive`
- **THEN** the new unified behavior is enabled
- **AND** users must opt-in to the new functionality
- **AND** columns show: Type, ID, Title, Details (context-aware)

#### Scenario: Unified mode displays both types

- **WHEN** unified mode is active
- **THEN** changes show Type="CHANGE" with delta and task counts
- **AND** specs show Type="SPEC" with requirement counts
- **AND** both types are navigable and selectable in the same table

#### Scenario: Type-specific actions in unified mode

- **WHEN** user presses 'e' on a change row in unified mode
- **THEN** the action is ignored (no edit for changes)
- **AND** help text does not show 'e' option
- **WHEN** user presses 'e' on a spec row in unified mode
- **THEN** the spec opens in the editor as usual
