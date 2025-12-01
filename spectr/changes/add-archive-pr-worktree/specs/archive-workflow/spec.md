# Archive Workflow Delta: Add PR Creation with Git Worktree Isolation

## ADDED Requirements

### Requirement: PR Creation Flag
The system SHALL provide a `--pr` flag on the `spectr archive` command that automates branch creation, commit, push, and pull request creation after a successful archive operation using git worktrees for isolation.

#### Scenario: Archive with PR flag creates pull request
- **WHEN** user runs `spectr archive <change-id> --pr`
- **THEN** the system completes the standard archive operation
- **AND** creates a git worktree on a new branch `archive-<change-id>`
- **AND** executes the archive operation within the worktree
- **AND** commits the changes with a structured message
- **AND** pushes the branch to origin
- **AND** creates a pull request using the detected platform CLI
- **AND** displays the PR URL on success

#### Scenario: PR flag requires explicit change ID
- **WHEN** user runs `spectr archive --pr` without a change ID
- **THEN** the system returns an error: "--pr flag requires explicit change ID"
- **AND** the archive operation does not proceed

#### Scenario: PR flag incompatible with interactive mode
- **WHEN** user runs `spectr archive --pr --interactive`
- **THEN** the system returns an error: "--pr cannot be used with --interactive"
- **AND** the archive operation does not proceed

#### Scenario: PR flag incompatible with no-validate
- **WHEN** user runs `spectr archive <change-id> --pr --no-validate`
- **THEN** the system returns an error: "--pr requires validation to ensure archive integrity"
- **AND** the archive operation does not proceed

### Requirement: Git Worktree Isolation
The system SHALL use git worktrees to isolate the PR creation workflow, ensuring the user's working directory is never modified by the `--pr` operation.

#### Scenario: Worktree created in temp directory
- **WHEN** the PR workflow begins
- **THEN** a git worktree is created in the system temp directory
- **AND** the worktree path includes a UUID to prevent conflicts
- **AND** the worktree is based on `origin/main` or `origin/master`

#### Scenario: Archive executes within worktree
- **WHEN** the worktree is created
- **THEN** the archive operation executes within the worktree
- **AND** all file modifications occur in the worktree, not the user's directory
- **AND** the user's working directory remains unchanged

#### Scenario: Worktree cleaned up after success
- **WHEN** the PR is created successfully
- **THEN** the worktree is removed automatically
- **AND** no temporary files remain in the temp directory

#### Scenario: Worktree cleaned up after failure
- **WHEN** any step in the PR workflow fails
- **THEN** the worktree is removed automatically
- **AND** an error message indicates what failed and how to recover

### Requirement: Git Hosting Platform Detection
The system SHALL detect the git hosting platform from the origin remote URL and use the appropriate CLI tool for PR creation.

#### Scenario: Detect GitHub from remote URL
- **WHEN** the origin remote URL contains "github.com"
- **THEN** the system identifies the platform as GitHub
- **AND** uses the `gh` CLI for PR creation

#### Scenario: Detect GitLab from remote URL
- **WHEN** the origin remote URL contains "gitlab.com" or "gitlab" in the hostname
- **THEN** the system identifies the platform as GitLab
- **AND** uses the `glab` CLI for MR creation

#### Scenario: Detect Gitea from remote URL
- **WHEN** the origin remote URL contains "gitea" or "forgejo" in the hostname
- **THEN** the system identifies the platform as Gitea
- **AND** uses the `tea` CLI for PR creation

#### Scenario: Detect Bitbucket from remote URL
- **WHEN** the origin remote URL contains "bitbucket.org" or "bitbucket" in the hostname
- **THEN** the system identifies the platform as Bitbucket
- **AND** provides a manual PR URL since no official CLI is supported

#### Scenario: Unknown platform returns error
- **WHEN** the origin remote URL does not match any known platform pattern
- **THEN** the system returns an error with the remote URL
- **AND** suggests manually creating the PR

### Requirement: CLI Tool Validation
The system SHALL validate that the required CLI tool is installed before attempting PR creation.

#### Scenario: CLI tool not installed
- **WHEN** the detected platform requires a CLI tool that is not installed
- **THEN** the system returns an error identifying the missing tool
- **AND** provides the installation URL for the tool

#### Scenario: CLI tool installed and authenticated
- **WHEN** the required CLI tool is installed
- **THEN** the system proceeds with PR creation
- **AND** uses the tool's default authentication

### Requirement: Branch Naming Convention
The system SHALL create branches with a consistent naming convention for archive PRs.

#### Scenario: Branch name from change ID
- **WHEN** creating a branch for archive PR
- **THEN** the branch is named `archive-<change-id>`
- **AND** the branch is based on `origin/main` or `origin/master`

#### Scenario: Branch already exists remotely
- **WHEN** a branch with name `archive-<change-id>` already exists on origin
- **THEN** the system returns an error
- **AND** suggests deleting the existing branch or using a different approach

### Requirement: Commit Message Format
The system SHALL generate structured commit messages for archive PR commits.

#### Scenario: Commit message includes archive metadata
- **WHEN** committing archive changes
- **THEN** the commit message uses format `archive(<change-id>): Archive completed change`
- **AND** the body includes the archive location
- **AND** the body includes spec operation counts (added, modified, removed, renamed)
- **AND** the body includes attribution to `spectr archive --pr`

#### Scenario: Commit message with skip-specs
- **WHEN** `--skip-specs` flag was used with `--pr`
- **THEN** the commit message omits the spec operation counts
- **AND** the body notes that spec updates were skipped

### Requirement: PR Content Generation
The system SHALL generate structured PR titles and bodies with archive summary information.

#### Scenario: PR title format
- **WHEN** creating a PR
- **THEN** the title is `archive(<change-id>): Archive completed change`

#### Scenario: PR body includes summary
- **WHEN** creating a PR
- **THEN** the body includes the change ID and archive location
- **AND** the body includes spec operation counts in a table
- **AND** the body lists updated capability names
- **AND** the body includes a review checklist
- **AND** the body includes attribution footer

### Requirement: Draft PR Support
The system SHALL support creating draft PRs via a `--draft` flag.

#### Scenario: Draft PR creation
- **WHEN** user runs `spectr archive <change-id> --pr --draft`
- **THEN** the PR is created as a draft (where supported by platform)
- **AND** the PR is not marked as ready for review

#### Scenario: Draft flag with unsupported platform
- **WHEN** `--draft` is used with Gitea
- **THEN** the system proceeds without draft flag
- **AND** displays a warning that draft PRs are not supported

### Requirement: PR Workflow Error Recovery
The system SHALL provide clear error messages and recovery guidance when PR creation fails.

#### Scenario: Not in git repository
- **WHEN** `--pr` is used outside a git repository
- **THEN** the system returns error: "Not in a git repository"
- **AND** suggests running `git init`

#### Scenario: No origin remote
- **WHEN** `--pr` is used without an origin remote configured
- **THEN** the system returns error: "No 'origin' remote configured"
- **AND** suggests running `git remote add origin <url>`

#### Scenario: Push fails
- **WHEN** pushing the branch fails
- **THEN** the system returns the git error message
- **AND** notes that the archive was successful but not pushed
- **AND** suggests manually pushing and creating the PR

#### Scenario: PR creation fails
- **WHEN** PR CLI invocation fails
- **THEN** the system returns the CLI error output
- **AND** notes that the branch was pushed successfully
- **AND** provides the manual PR creation URL

### Requirement: Base Branch Detection
The system SHALL automatically detect the default base branch for the repository.

#### Scenario: Main branch detection
- **WHEN** determining the base branch
- **THEN** the system checks for `origin/main` first
- **AND** falls back to `origin/master` if main doesn't exist

#### Scenario: Neither main nor master exists
- **WHEN** neither `origin/main` nor `origin/master` exists
- **THEN** the system returns an error
- **AND** suggests checking the remote configuration
