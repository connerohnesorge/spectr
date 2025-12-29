# Cli Interface Specification Delta

## ADDED Requirements

### Requirement: Editor Hotkey in Interactive Specs List

The interactive specs list mode SHALL provide an 'e' hotkey that opens the
selected spec file in the user's configured editor.

#### Scenario: User presses 'e' to edit a spec

- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses the 'e' key on a selected spec
- **THEN** the file `spectr/specs/\<spec-id\>/spec.md` is opened in the editor
  specified by $EDITOR environment variable
- **AND** the TUI waits for the editor to close
- **AND** the TUI remains active after the editor closes
- **AND** the same row remains selected

#### Scenario: User edits spec and returns to TUI

- **WHEN** user presses 'e' to open a spec
- **AND** makes changes in the editor and saves
- **AND** closes the editor
- **THEN** the TUI returns to the interactive list view
- **AND** the user can continue navigating or edit another spec
- **AND** the user can quit with 'q' or Ctrl+C as normal

#### Scenario: EDITOR environment variable not set

- **WHEN** user presses 'e' to edit a spec
- **AND** $EDITOR environment variable is not set
- **THEN** display an error message "EDITOR environment variable not set"
- **AND** the TUI remains in interactive mode
- **AND** the user can continue navigating or quit

#### Scenario: Spec file does not exist

- **WHEN** user presses 'e' to edit a spec
- **AND** the spec file at `spectr/specs/\<spec-id\>/spec.md` does not exist
- **THEN** display an error message "Spec file not found: \<path\>"
- **AND** the TUI remains in interactive mode
- **AND** the user can continue navigating or quit

#### Scenario: Editor launch fails

- **WHEN** user presses 'e' to edit a spec
- **AND** the editor process fails to launch (e.g., editor binary not found,
  permission error)
- **THEN** display an error message with the underlying error details
- **AND** the TUI remains in interactive mode
- **AND** the user can retry or quit

#### Scenario: Help text shows editor hotkey

- **WHEN** interactive specs mode is displayed
- **THEN** the help text includes "e: edit spec" or similar guidance
- **AND** the help text shows all available keys including navigation, enter, e,
  and quit keys

### Requirement: Editor Hotkey Scope

The 'e' hotkey for opening files in $EDITOR SHALL only be available in specs
list mode, not in changes list mode.

#### Scenario: Editor hotkey not available for changes

- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses 'e' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'e: edit' option
- **AND** only standard navigation and clipboard actions are available

#### Scenario: Rationale for specs-only scope

- **WHEN** user reviews this specification
- **THEN** they understand that changes have multiple files (proposal.md,
  tasks.md, design.md, delta specs)
- **AND** pressing 'e' on a change would be ambiguous (which file to open?)
- **AND** specs have a single canonical file (spec.md) making 'e' unambiguous
- **AND** this design decision can be revisited in a future change if multi-file
  editing is needed
