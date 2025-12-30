# Delta Specification

## ADDED Requirements

### Requirement: Stdout Output Mode for Interactive List

The `spectr list` command SHALL support a `--stdout` flag that, when combined
with interactive mode (`-I`), outputs the selected item ID to stdout instead of
copying it to the system clipboard.

#### Scenario: User runs list with --stdout and -I flags

- **WHEN** user runs `spectr list -I --stdout`
- **AND** user navigates to a row and presses Enter
- **THEN** the selected ID is printed to stdout (just the ID, no formatting)
- **AND** no clipboard operation is performed
- **AND** the command exits with code 0

#### Scenario: Stdout mode with changes

- **WHEN** user runs `spectr list -I --stdout` (changes mode)
- **AND** user selects a change and presses Enter
- **THEN** only the change ID is printed to stdout (e.g., `add-feature`)
- **AND** no "Copied:" prefix or other formatting is included

#### Scenario: Stdout mode with specs

- **WHEN** user runs `spectr list --specs -I --stdout`
- **AND** user selects a spec and presses Enter
- **THEN** only the spec ID is printed to stdout (e.g., `cli-interface`)
- **AND** no "Copied:" prefix or other formatting is included

#### Scenario: Stdout mode with unified view

- **WHEN** user runs `spectr list --all -I --stdout`
- **AND** user selects an item and presses Enter
- **THEN** only the item ID is printed to stdout
- **AND** no "Copied:" prefix or other formatting is included

#### Scenario: Stdout flag requires interactive mode

- **WHEN** user runs `spectr list --stdout` without `-I`
- **THEN** an error is displayed: "cannot use --stdout without --interactive
  (-I)"
- **AND** the command exits with code 1

#### Scenario: Stdout flag mutually exclusive with JSON

- **WHEN** user runs `spectr list -I --stdout --json`
- **THEN** an error is displayed: "cannot use --stdout with --json"
- **AND** the command exits with code 1

#### Scenario: Stdout mode enables piping

- **WHEN** user runs `spectr list -I --stdout | xargs echo`
- **AND** user selects an item
- **THEN** the pipeline receives the clean ID string
- **AND** the downstream command processes the ID correctly

#### Scenario: User quits without selection in stdout mode

- **WHEN** user runs `spectr list -I --stdout`
- **AND** user presses 'q' or Ctrl+C without selecting
- **THEN** nothing is printed to stdout
- **AND** the command exits with code 0

#### Scenario: Stdout mode help text

- **WHEN** user runs `spectr list --help`
- **THEN** the help text shows `--stdout` flag
- **AND** the description explains it outputs to stdout instead of clipboard
- **AND** the help indicates it requires `-I` flag
