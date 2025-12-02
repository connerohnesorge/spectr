## ADDED Requirements

### Requirement: PR Command Structure
The system SHALL provide a `spectr pr` command with `archive` and `proposal` subcommands for creating pull requests from Spectr changes using git worktree isolation.

#### Scenario: User runs spectr pr without subcommand
- **WHEN** user runs `spectr pr` without a subcommand
- **THEN** help text is displayed showing available subcommands (archive, proposal)
- **AND** the command exits with code 0

#### Scenario: User runs spectr pr archive
- **WHEN** user runs `spectr pr archive <change-id>`
- **THEN** the system executes the archive PR workflow
- **AND** a PR is created with the archived change

#### Scenario: User runs spectr pr proposal
- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the system executes the proposal PR workflow
- **AND** a PR is created with the change proposal copied (not archived)

### Requirement: PR Archive Subcommand
The `spectr pr archive` subcommand SHALL create a pull request containing an archived Spectr change, executing the archive workflow in an isolated git worktree.

#### Scenario: Archive PR workflow execution
- **WHEN** user runs `spectr pr archive <change-id>`
- **THEN** the system creates a temporary git worktree on branch `spectr/<change-id>`
- **AND** executes `spectr archive <change-id> --yes` within the worktree
- **AND** stages all changes in `spectr/` directory
- **AND** commits with structured message including archive metadata
- **AND** pushes the branch to origin
- **AND** creates a PR using the detected platform's CLI
- **AND** cleans up the temporary worktree
- **AND** displays the PR URL on success

#### Scenario: Archive PR with skip-specs flag
- **WHEN** user runs `spectr pr archive <change-id> --skip-specs`
- **THEN** the `--skip-specs` flag is passed to the underlying archive command
- **AND** spec merging is skipped during the archive operation

#### Scenario: Archive PR preserves user working directory
- **WHEN** user runs `spectr pr archive <change-id>`
- **AND** user has uncommitted changes in their working directory
- **THEN** the user's working directory is NOT modified
- **AND** the archive operation executes only within the isolated worktree
- **AND** uncommitted changes are NOT included in the PR

### Requirement: PR Proposal Subcommand
The `spectr pr proposal` subcommand SHALL create a pull request containing a Spectr change proposal for review, copying the change to an isolated git worktree without archiving.

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

### Requirement: PR Common Flags
Both `spectr pr archive` and `spectr pr proposal` subcommands SHALL support common flags for controlling PR creation behavior.

#### Scenario: Base branch flag
- **WHEN** user provides `--base <branch>` flag
- **THEN** the PR targets the specified branch instead of auto-detected main/master

#### Scenario: Auto-detect base branch
- **WHEN** user does NOT provide `--base` flag
- **AND** `origin/main` exists
- **THEN** the PR targets `main`

#### Scenario: Fallback to master
- **WHEN** user does NOT provide `--base` flag
- **AND** `origin/main` does NOT exist
- **AND** `origin/master` exists
- **THEN** the PR targets `master`

#### Scenario: Draft PR flag
- **WHEN** user provides `--draft` flag
- **THEN** the PR is created as a draft PR on platforms that support it

#### Scenario: Force flag for existing branch
- **WHEN** user provides `--force` flag
- **AND** branch `spectr/<change-id>` already exists on remote
- **THEN** the existing branch is deleted and recreated
- **AND** the PR workflow proceeds normally

#### Scenario: Branch conflict without force
- **WHEN** branch `spectr/<change-id>` already exists on remote
- **AND** `--force` flag is NOT provided
- **THEN** an error is displayed: "Branch 'spectr/<change-id>' already exists. Use --force to overwrite."
- **AND** the command exits with code 1

#### Scenario: Dry run flag
- **WHEN** user provides `--dry-run` flag
- **THEN** the system displays what would be done without executing
- **AND** no git operations are performed
- **AND** no PR is created

### Requirement: Git Platform Detection
The system SHALL automatically detect the git hosting platform from the origin remote URL and use the appropriate CLI tool for PR creation.

#### Scenario: Detect GitHub platform
- **WHEN** origin remote URL contains `github.com`
- **THEN** platform is detected as GitHub
- **AND** `gh` CLI is used for PR creation

#### Scenario: Detect GitLab platform
- **WHEN** origin remote URL contains `gitlab.com` or `gitlab`
- **THEN** platform is detected as GitLab
- **AND** `glab` CLI is used for MR creation

#### Scenario: Detect Gitea platform
- **WHEN** origin remote URL contains `gitea` or `forgejo`
- **THEN** platform is detected as Gitea
- **AND** `tea` CLI is used for PR creation

#### Scenario: Detect Bitbucket platform
- **WHEN** origin remote URL contains `bitbucket.org` or `bitbucket`
- **THEN** platform is detected as Bitbucket
- **AND** manual PR URL is provided (no CLI automation)

#### Scenario: Unknown platform error
- **WHEN** origin remote URL does not match any known platform
- **THEN** an error is displayed with the detected URL
- **AND** guidance is provided for manual PR creation

#### Scenario: SSH URL format support
- **WHEN** origin remote uses SSH format (e.g., `git@github.com:org/repo.git`)
- **THEN** platform detection correctly identifies the host

#### Scenario: HTTPS URL format support
- **WHEN** origin remote uses HTTPS format (e.g., `https://github.com/org/repo.git`)
- **THEN** platform detection correctly identifies the host

### Requirement: Platform CLI Availability
The system SHALL verify that the required platform CLI tool is installed and authenticated before attempting PR creation.

#### Scenario: CLI not installed
- **WHEN** the required CLI tool (gh/glab/tea) is not installed
- **THEN** an error is displayed: "<tool> CLI is required for <platform> PR creation. Install: <install-url>"
- **AND** the command exits with code 1

#### Scenario: CLI not authenticated
- **WHEN** the required CLI tool is installed but not authenticated
- **THEN** an error is displayed with authentication instructions
- **AND** the command exits with code 1

### Requirement: Git Worktree Isolation
The PR commands SHALL use git worktrees to provide complete isolation from the user's working directory.

#### Scenario: Worktree created in temp directory
- **WHEN** PR workflow starts
- **THEN** a worktree is created in the system temp directory
- **AND** the worktree path includes a UUID to prevent conflicts

#### Scenario: Worktree based on origin branch
- **WHEN** worktree is created
- **THEN** it is based on the remote base branch (origin/main or origin/master)
- **AND** it does NOT include any local uncommitted changes

#### Scenario: Worktree cleanup on success
- **WHEN** PR workflow completes successfully
- **THEN** the temporary worktree is removed
- **AND** no temporary files remain

#### Scenario: Worktree cleanup on failure
- **WHEN** PR workflow fails at any stage
- **THEN** the temporary worktree is still removed
- **AND** an appropriate error message is displayed

#### Scenario: Git version requirement
- **WHEN** git version is less than 2.5
- **THEN** an error is displayed: "Git >= 2.5 required for worktree support. Current version: <version>"
- **AND** the command exits with code 1

### Requirement: PR Commit Message Format
The system SHALL generate structured commit messages that follow conventional commit format and include relevant metadata.

#### Scenario: Archive commit message format
- **WHEN** `spectr pr archive` commits changes
- **THEN** the commit message includes:
  - Title: `spectr(archive): <change-id>`
  - Archive location path
  - Spec operation counts (added, modified, removed, renamed)
  - Attribution: "Generated by: spectr pr archive"

#### Scenario: Proposal commit message format
- **WHEN** `spectr pr proposal` commits changes
- **THEN** the commit message includes:
  - Title: `spectr(proposal): <change-id>`
  - Proposal location path
  - Attribution: "Generated by: spectr pr proposal"

### Requirement: PR Body Content
The system SHALL generate PR body content that helps reviewers understand the change.

#### Scenario: Archive PR body content
- **WHEN** PR is created for archive
- **THEN** the PR body includes:
  - Summary section with change ID and archive location
  - Spec updates table with operation counts
  - List of updated capabilities
  - Review checklist

#### Scenario: Proposal PR body content
- **WHEN** PR is created for proposal
- **THEN** the PR body includes:
  - Summary section with change ID and location
  - List of included files (proposal.md, tasks.md, specs/)
  - Review checklist

### Requirement: PR Branch Naming
The system SHALL use a consistent branch naming convention for PR branches.

#### Scenario: Branch name format
- **WHEN** PR workflow creates a branch
- **THEN** the branch is named `spectr/<change-id>`

#### Scenario: Branch name with special characters
- **WHEN** change ID contains only valid kebab-case characters
- **THEN** the branch name is valid for git

### Requirement: PR Error Handling
The system SHALL provide clear error messages and guidance when PR creation fails.

#### Scenario: Not in git repository
- **WHEN** user runs `spectr pr` from outside a git repository
- **THEN** an error is displayed: "Not in a git repository"
- **AND** the command exits with code 1

#### Scenario: No origin remote
- **WHEN** user runs `spectr pr` and no origin remote exists
- **THEN** an error is displayed: "No 'origin' remote configured"
- **AND** guidance is provided to add a remote

#### Scenario: Change does not exist
- **WHEN** user runs `spectr pr <subcommand> <change-id>`
- **AND** the change does not exist
- **THEN** an error is displayed: "Change '<change-id>' not found"
- **AND** the command exits with code 1

#### Scenario: Push failure
- **WHEN** git push fails (e.g., network error)
- **THEN** an error is displayed with the git error message
- **AND** guidance is provided for manual recovery
- **AND** worktree is still cleaned up

#### Scenario: PR creation failure with pushed branch
- **WHEN** push succeeds but PR creation fails
- **THEN** an error is displayed with the PR CLI error
- **AND** a message indicates: "Branch was pushed. Create PR manually or retry."
- **AND** worktree is still cleaned up

### Requirement: Partial Change ID Resolution for PR Commands
The `spectr pr` subcommands SHALL support intelligent partial ID matching consistent with the archive command's resolution algorithm.

#### Scenario: Exact ID match for PR commands
- **WHEN** user runs `spectr pr archive exact-change-id`
- **AND** a change with ID `exact-change-id` exists
- **THEN** the PR workflow proceeds with that change

#### Scenario: Unique prefix match for PR commands
- **WHEN** user runs `spectr pr proposal refactor`
- **AND** only one change ID starts with `refactor`
- **THEN** a resolution message is displayed
- **AND** the PR workflow proceeds with the resolved ID
