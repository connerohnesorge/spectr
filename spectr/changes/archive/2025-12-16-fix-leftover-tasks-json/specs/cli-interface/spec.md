# Delta Specification

## MODIFIED Requirements

### Requirement: PR Remove Subcommand

The `spectr pr rm` subcommand SHALL create a pull request that removes a change
directory from the repository, using the same git worktree isolation as other PR
subcommands.

The command supports aliases `r` and `remove` for convenience.

#### Scenario: User runs spectr pr rm with change ID

- **WHEN** user runs `spectr pr rm \<change-id\>`
- **THEN** the system creates a temporary git worktree on branch
  `spectr/remove/\<change-id\>`
- **AND** removes the change directory from `spectr/changes/\<change-id\>` in
  the
  worktree
- **AND** stages the deletion
- **AND** commits with a structured message indicating removal
- **AND** pushes the branch to origin
- **AND** creates a PR using the detected platform's CLI
- **AND** cleans up the temporary worktree
- **AND** displays the PR URL on success
- **AND** removes the local change directory after successful PR creation

#### Scenario: User runs spectr pr rm without change ID

- **WHEN** user runs `spectr pr rm` without a change ID argument
- **THEN** an interactive table is displayed showing available changes
- **AND** user can navigate and select a change
- **AND** the remove workflow proceeds with the selected change ID

#### Scenario: User runs spectr pr r shorthand

- **WHEN** user runs `spectr pr r \<change-id\>`
- **THEN** the system executes the remove PR workflow identically to `spectr pr
  rm`
- **AND** all flags work with the alias

#### Scenario: Remove PR with draft flag

- **WHEN** user runs `spectr pr rm \<change-id\> --draft`
- **THEN** the PR is created as a draft PR on platforms that support it

#### Scenario: Remove PR with force flag

- **WHEN** user runs `spectr pr rm \<change-id\> --force`
- **AND** branch `spectr/remove/\<change-id\>` already exists on remote
- **THEN** the existing branch is deleted and recreated
- **AND** the PR workflow proceeds normally

#### Scenario: Remove branch conflict without force

- **WHEN** user runs `spectr pr rm \<change-id\>`
- **AND** branch `spectr/remove/\<change-id\>` already exists on remote
- **AND** `--force` flag is NOT provided
- **THEN** an error is displayed: "branch 'spectr/remove/\<change-id\>' already
  exists on remote; use --force to delete"
- **AND** the command exits with code 1

#### Scenario: Remove PR with dry-run flag

- **WHEN** user runs `spectr pr rm \<change-id\> --dry-run`
- **THEN** the system displays what would be done without executing
- **AND** no git operations are performed
- **AND** no PR is created
- **AND** no local cleanup is performed

#### Scenario: Remove PR with base branch flag

- **WHEN** user runs `spectr pr rm \<change-id\> --base develop`
- **THEN** the PR targets the `develop` branch instead of auto-detected
  main/master

#### Scenario: Change does not exist

- **WHEN** user runs `spectr pr rm \<change-id\>`
- **AND** the change does not exist in `spectr/changes/`
- **THEN** an error is displayed: "change '\<change-id\>' not found in
  spectr/changes/"
- **AND** the command exits with code 1

#### Scenario: Remove cleans up local change directory

- **WHEN** user runs `spectr pr rm \<change-id\>`
- **AND** PR creation succeeds
- **THEN** the system displays: "Cleaning up local change directory:
  spectr/changes/\<change-id\>/"
- **AND** the local change directory is removed including all files (tracked and
  untracked)

#### Scenario: Partial ID resolution for remove command

- **WHEN** user runs `spectr pr rm refactor`
- **AND** only one change ID starts with `refactor`
- **THEN** a resolution message is displayed
- **AND** the PR workflow proceeds with the resolved ID

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

#### Scenario: Cleanup for archive mode

- **WHEN** user runs `spectr pr archive \<change-id\>`
- **AND** PR creation succeeds
- **THEN** the system displays: "Cleaning up local change directory:
  spectr/changes/\<change-id\>/"
- **AND** the local change directory is removed
- **AND** the change is archived in the worktree (pulled when PR merges)

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
