# Delta Specification

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

#### Scenario: Help text uses condensed two-line format

- **WHEN** interactive mode is displayed in any mode (changes, specs, or
  unified)
- **THEN** the help text is formatted across two lines
- **AND** line 1 shows controls: navigation, actions, and item count
- **AND** line 2 shows project path
- **AND** navigation hint uses condensed format `↑/↓/j/k` instead of `↑/↓ or
  j/k`
- **AND** edit action uses short label `e: edit` instead of `e: edit proposal`
  or `e: edit spec`

#### Scenario: Help text format for changes mode

- **WHEN** user runs `spectr list -I` (changes mode)
- **THEN** line 1 shows: `↑/↓/j/k: nav | Enter: copy | e: edit | a: arch | q: quit`
- **AND** line 2 shows: `project: <path>`

#### Scenario: Help text format for specs mode

- **WHEN** user runs `spectr list --specs -I` (specs mode)
- **THEN** line 1 shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | q: quit`
- **AND** line 2 shows: `project: <path>`
- **AND** archive hotkey is NOT shown (specs cannot be archived)

#### Scenario: Help text format for unified mode

- **WHEN** user runs `spectr list --all -I` (unified mode)
- **THEN** line 1 shows: `↑/↓/j/k: nav | Enter: copy | e: edit | t: toggle | q: quit`
- **AND** line 2 shows: `project: <path>`
- **AND** archive hotkey is NOT shown in unified mode for simplicity

## ADDED Requirements

### Requirement: Archive Hotkey in Interactive Changes Mode

The interactive changes list mode SHALL provide an 'a' hotkey that archives the
currently selected change, invoking the same workflow as `spectr archive
<change-id>`.

#### Scenario: User presses 'a' to archive a change

- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses the 'a' key on a selected change
- **THEN** the interactive mode exits
- **AND** the archive workflow begins for the selected change ID
- **AND** validation, task checking, and spec updates proceed as if the ID was
  provided as an argument
- **AND** all confirmation prompts and flags work normally

#### Scenario: Archive hotkey not available in specs mode

- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses 'a' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'a: archive' option

#### Scenario: Archive hotkey not available in unified mode

- **WHEN** user is in unified interactive mode (`spectr list --all -I`)
- **AND** user presses 'a' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'a: archive' option
- **AND** this avoids confusion when a spec row is selected

#### Scenario: Archive workflow integration

- **WHEN** the archive hotkey triggers the archive workflow
- **THEN** the workflow uses the same code path as `spectr archive <id>`
- **AND** the selected change ID is passed to the archive workflow
- **AND** success or failure is reported after the workflow completes

#### Scenario: Help text shows archive hotkey in changes mode

- **WHEN** interactive changes mode is displayed
- **THEN** the help text includes `a: archive` in the controls line
- **AND** the hotkey appears after `e: edit` and before `q: quit`
