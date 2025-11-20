## ADDED Requirements

### Requirement: Propose Command for Automated PR Creation

The `spectr propose <change-id>` command SHALL automate the workflow of creating a pull request from a newly scaffolded change proposal. It SHALL validate that the change exists and is uncommitted, create a dedicated git branch, stage only the change folder, commit the changes, push to the remote, and invoke the appropriate PR creation tool based on the git hosting platform.

#### Scenario: User creates PR for a new proposal

- **WHEN** user runs `spectr propose my-feature`
- **AND** the change proposal exists in `spectr/changes/my-feature/`
- **AND** the folder is not yet committed to git
- **AND** a git repository with an `origin` remote is configured
- **THEN** a new branch named `add-my-feature` is created
- **AND** only the `spectr/changes/my-feature/` folder is staged and committed
- **AND** the branch is pushed to the origin remote
- **AND** the appropriate PR CLI tool is invoked (gh, glab, or tea)
- **AND** the PR is created with title "Propose: my-feature"
- **AND** the PR body includes the purpose and changes from the proposal
- **AND** the PR URL is displayed to the user

#### Scenario: Git platform auto-detection for GitHub

- **WHEN** the `origin` remote URL contains `github.com`
- **THEN** the `gh pr create` command is used
- **AND** the PR is created on GitHub

#### Scenario: Git platform auto-detection for GitLab

- **WHEN** the `origin` remote URL contains `gitlab.com` or is a self-hosted GitLab instance
- **THEN** the `glab mr create` command is used
- **AND** the MR (merge request) is created on GitLab

#### Scenario: Git platform auto-detection for Gitea

- **WHEN** the `origin` remote URL contains `gitea` or `forgejo`
- **THEN** the `tea pr create` command is used
- **AND** the PR is created on Gitea or Forgejo

#### Scenario: Change folder does not exist

- **WHEN** user runs `spectr propose non-existent`
- **AND** the folder `spectr/changes/non-existent/` does not exist
- **THEN** an error is displayed: "Change proposal 'non-existent' not found in spectr/changes/"
- **AND** no git operations occur
- **AND** the command exits with error code 1

#### Scenario: Change folder is already committed

- **WHEN** user runs `spectr propose already-tracked`
- **AND** the folder `spectr/changes/already-tracked/` is already tracked by git
- **THEN** an error is displayed: "Change 'already-tracked' is already committed. Use git/gh/glab/tea directly to create a PR."
- **AND** no branch is created
- **AND** the command exits with error code 1

#### Scenario: Not in a git repository

- **WHEN** user runs `spectr propose my-feature`
- **AND** the current directory is not inside a git repository
- **THEN** an error is displayed: "Not in a git repository. Initialize git with 'git init'."
- **AND** the command exits with error code 1

#### Scenario: Origin remote not configured

- **WHEN** user runs `spectr propose my-feature`
- **AND** the change folder exists and is uncommitted
- **AND** the git repository has no `origin` remote
- **THEN** an error is displayed: "No 'origin' remote configured. Run 'git remote add origin <url>' first."
- **AND** no branch is created
- **AND** the command exits with error code 1

#### Scenario: Git hosting platform not detected

- **WHEN** user runs `spectr propose my-feature`
- **AND** the origin remote URL does not match GitHub, GitLab, or Gitea patterns
- **THEN** an error is displayed: "Could not detect git hosting platform. Remote URL: [url]. Please create PR manually using gh, glab, or tea."
- **AND** the branch is created and pushed, but PR creation is not attempted
- **AND** the command exits with error code 1

#### Scenario: PR CLI tool not installed

- **WHEN** user runs `spectr propose my-feature` for a GitHub repository
- **AND** the `gh` CLI tool is not installed
- **THEN** an error is displayed: "gh not found. Install from https://github.com/cli/cli"
- **AND** the branch is created and pushed, but the PR is not created
- **AND** the command exits with error code 1

#### Scenario: Unrelated uncommitted changes are not included

- **WHEN** user has uncommitted changes in other directories (not in `spectr/changes/my-feature/`)
- **AND** runs `spectr propose my-feature`
- **THEN** only the `spectr/changes/my-feature/` folder is staged and committed
- **AND** unrelated uncommitted changes remain in the working tree

#### Scenario: PR creation succeeds and URL is displayed

- **WHEN** the PR is successfully created
- **THEN** a success message is displayed with the PR URL
- **AND** the message format is "PR created: https://github.com/owner/repo/pull/123"
- **AND** the command exits with code 0
