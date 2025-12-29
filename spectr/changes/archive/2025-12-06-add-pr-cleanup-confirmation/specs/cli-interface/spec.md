# Delta Specification

## ADDED Requirements

### Requirement: PR Proposal Local Change Cleanup Confirmation

After a successful `spectr pr proposal` command that creates a pull request, the
system SHALL prompt the user with a Bubbletea TUI confirmation menu asking
whether to remove the local change proposal directory from `spectr/changes/`.

This prompt helps users maintain a clean working directory by offering an
opportunity to remove proposals that have been submitted for review, while
defaulting to "No" for safety. The menu uses arrow key navigation and styled
rendering consistent with other spectr interactive modes.

#### Scenario: Successful PR proposal triggers cleanup prompt

- **WHEN** user runs `spectr pr proposal \<change-id\>` and PR creation succeeds
- **AND** the PR URL is displayed to the user
- **THEN** the system displays a Bubbletea TUI menu: "Remove local change
  proposal from spectr/changes/?"
- **AND** the menu shows two options: "No, keep it" and "Yes, remove it"
- **AND** the default selection is "No, keep it" (cursor starts on this option)
- **AND** the menu supports arrow key navigation (↑/↓) and Enter to confirm

#### Scenario: User confirms cleanup via TUI

- **WHEN** the cleanup TUI menu is displayed
- **AND** user navigates to "Yes, remove it" and presses Enter
- **THEN** the system removes the directory `spectr/changes/\<change-id\>/`
- **AND** displays confirmation: "Removed local change: \<change-id\>"

#### Scenario: User declines cleanup via TUI

- **WHEN** the cleanup TUI menu is displayed
- **AND** user presses Enter on the default "No, keep it" option
- **THEN** the system keeps the local change directory
- **AND** displays: "Local change kept: spectr/changes/\<change-id\>/"

#### Scenario: User cancels cleanup menu

- **WHEN** the cleanup TUI menu is displayed
- **AND** user presses 'q' or Ctrl+C
- **THEN** the system keeps the local change directory (same as declining)
- **AND** the command exits successfully

#### Scenario: Non-interactive mode skips prompt

- **WHEN** user runs `spectr pr proposal \<change-id\> --yes`
- **AND** PR creation succeeds
- **THEN** the cleanup prompt is NOT displayed
- **AND** the local change directory is kept (safe default)
- **AND** the command exits successfully

#### Scenario: Cleanup prompt only for proposal mode

- **WHEN** user runs `spectr pr archive \<change-id\>`
- **AND** PR creation succeeds
- **THEN** the cleanup prompt is NOT displayed
- **AND** the change is already moved to archive by the archive workflow

#### Scenario: PR creation fails

- **WHEN** user runs `spectr pr proposal \<change-id\>`
- **AND** PR creation fails at any step
- **THEN** the cleanup prompt is NOT displayed
- **AND** the local change directory remains intact

#### Scenario: Cleanup removal error handling

- **WHEN** the user confirms cleanup
- **AND** removal of the change directory fails (e.g., permission error)
- **THEN** display a warning: "Warning: Failed to remove change directory:
  `<error>`"
- **AND** the command still exits successfully (non-fatal error)
