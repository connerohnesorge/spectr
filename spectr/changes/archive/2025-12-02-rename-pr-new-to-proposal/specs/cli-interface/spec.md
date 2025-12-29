# Cli Interface Specification Delta

## MODIFIED Requirements

### Requirement: PR Proposal Subcommand

The `spectr pr proposal` subcommand SHALL create a pull request containing a Spectr change proposal for review, copying the change to an isolated git worktree without archiving. This command replaces the deprecated `spectr pr new` command.

The renaming from `new` to `proposal` aligns CLI terminology with the `/spectr:proposal` slash command naming convention, creating consistent vocabulary across CLI and IDE integrations.

#### Scenario: Proposal PR workflow execution

- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the system creates a temporary git worktree on branch `spectr/<change-id>`
- **AND** copies the change directory from source to worktree
- **AND** stages all changes in `spectr/` directory
- **AND** commits with structured message for proposal review
- **AND** pushes the branch to origin
- **AND** creates a PR using the detected platform's CLI
- **AND** cleans up the temporary worktree
- **AND** displays the PR URL on success

#### Scenario: Proposal PR does not archive

- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the original change remains in `spectr/changes/<change-id>/`
- **AND** the change is NOT moved to archive
- **AND** spec merging does NOT occur

#### Scenario: Proposal PR validates change first

- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the system runs validation on the change
- **AND** warnings are displayed if validation issues exist
- **AND** the PR workflow continues (validation does not block)

#### Scenario: User runs spectr pr without subcommand

- **WHEN** user runs `spectr pr` without a subcommand
- **THEN** help text is displayed showing available subcommands (archive, proposal)
- **AND** the command exits with code 0

#### Scenario: Unique prefix match for PR proposal command

- **WHEN** user runs `spectr pr proposal refactor`
- **AND** only one change ID starts with `refactor`
- **THEN** a resolution message is displayed
- **AND** the PR workflow proceeds with the resolved ID
