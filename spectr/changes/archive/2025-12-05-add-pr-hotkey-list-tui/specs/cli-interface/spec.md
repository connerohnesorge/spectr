# Delta Specification

## ADDED Requirements

### Requirement: PR Hotkey in Interactive Changes List Mode

The interactive changes list mode SHALL provide a `P` (Shift+P) hotkey that
exits the TUI and enters the `spectr pr` workflow for the selected change,
allowing users to create pull requests without manually copying the change ID.

#### Scenario: User presses Shift+P to enter PR mode

- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses the `P` key (Shift+P) on a selected change
- **THEN** the interactive mode exits
- **AND** the system enters PR mode for the selected change ID
- **AND** the user is prompted to select PR type (archive or proposal)

#### Scenario: PR hotkey not available in specs mode

- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses `P` key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show `P: pr` option

#### Scenario: PR hotkey not available in unified mode

- **WHEN** user is in unified interactive mode (`spectr list --all -I`)
- **AND** user presses `P` key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show `P: pr` option
- **AND** this avoids confusion when a spec row is selected

#### Scenario: Help text shows PR hotkey in changes mode

- **WHEN** user presses `?` in changes mode
- **THEN** the help text includes `P: pr` in the controls line
- **AND** the hotkey appears alongside other change-specific hotkeys (e, a)

#### Scenario: PR workflow integration

- **WHEN** the PR hotkey triggers the PR workflow
- **THEN** the workflow uses the same code path as `spectr pr`
- **AND** the selected change ID is passed to the PR workflow
- **AND** the user can select between archive and proposal modes

### Requirement: VHS Demo for PR Hotkey

The system SHALL provide a VHS tape demonstrating the Shift+P hotkey utility in
the interactive list TUI.

#### Scenario: Developer finds PR hotkey demo

- **WHEN** a developer reviews the VHS tape files in `assets/vhs/`
- **THEN** they SHALL find `pr-hotkey.tape` demonstrating the PR hotkey workflow

#### Scenario: User sees PR hotkey demo

- **WHEN** a user views the PR hotkey demo GIF
- **THEN** they SHALL see the interactive list being invoked
- **AND** they SHALL see the `P` key being pressed
- **AND** they SHALL see the PR mode being entered for the selected change
