## ADDED Requirements

### Requirement: Help Toggle Hotkey

The interactive TUI modes SHALL hide hotkey hints by default and reveal them only when the user presses `?`, reducing visual clutter while maintaining discoverability.

#### Scenario: Default view shows minimal footer

- **WHEN** user enters any interactive TUI mode (list, archive, validate)
- **THEN** the footer displays only: item count, project path, and `?: help`
- **AND** the full hotkey reference is NOT shown
- **AND** navigation and all other hotkeys remain functional

#### Scenario: User presses '?' to reveal help

- **WHEN** user presses `?` while in interactive mode
- **THEN** the full hotkey reference is displayed in the footer area
- **AND** the reference includes all available hotkeys for the current mode
- **AND** the view updates immediately

#### Scenario: User dismisses help by pressing '?' again

- **WHEN** user presses `?` while help is visible
- **THEN** the help is hidden
- **AND** the minimal footer is restored

#### Scenario: Help auto-hides on navigation

- **WHEN** user presses a navigation key (↑/↓/j/k) while help is visible
- **THEN** the help is automatically hidden
- **AND** the navigation action is performed
- **AND** the minimal footer is restored

#### Scenario: Help content matches mode

- **WHEN** help is displayed in changes mode
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | q: quit`
- **WHEN** help is displayed in specs mode
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | q: quit`
- **WHEN** help is displayed in unified mode
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | t: filter | q: quit`

## MODIFIED Requirements

### Requirement: Interactive List Mode

The interactive list mode in `spectr list` is extended to support unified display of changes and specifications alongside existing separate modes.

#### Previous behavior

The system displays either changes OR specs in interactive mode based on the `--specs` flag. Columns and behavior are specific to each item type.

#### New behavior

- When `--all` is provided with `--interactive`, both changes and specs are shown together with unified columns
- When neither `--all` nor `--specs` are provided, changes-only mode is default (backward compatible)
- When `--specs` is provided without `--all`, specs-only mode is used (backward compatible)
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

#### Scenario: Help text uses minimal footer by default

- **WHEN** interactive mode is displayed in any mode (changes, specs, or unified)
- **THEN** the footer shows: item count, project path, and `?: help`
- **AND** the full hotkey reference is hidden until `?` is pressed

#### Scenario: Help text format for changes mode

- **WHEN** user presses `?` in changes mode (`spectr list -I`)
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | q: quit`
- **AND** pressing `?` again or navigating hides the help

#### Scenario: Help text format for specs mode

- **WHEN** user presses `?` in specs mode (`spectr list --specs -I`)
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | q: quit`
- **AND** archive hotkey is NOT shown (specs cannot be archived)
