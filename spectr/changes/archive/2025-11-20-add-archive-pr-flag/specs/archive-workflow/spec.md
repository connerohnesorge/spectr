## ADDED Requirements

### Requirement: Archive PR Automation Flag

The system SHALL provide a `--pr` flag on the `spectr archive` command that automatically creates a pull request after successful archive completion, including branch creation, committing archived files and updated specs, pushing to remote, and invoking the appropriate platform PR CLI tool.

#### Scenario: Archive with PR flag creates branch and PR

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** the archive operation completes successfully
- **AND** a git repository with an `origin` remote is configured
- **THEN** a new branch named `archive-my-feature` is created
- **AND** the archived directory and updated specs are staged and committed
- **AND** the branch is pushed to the origin remote
- **AND** the appropriate PR CLI tool is invoked based on platform detection
- **AND** the PR URL is displayed to the user

#### Scenario: Archive fails, no git operations occur

- **WHEN** user runs `spectr archive invalid-change --pr`
- **AND** the archive operation fails validation
- **THEN** no git branch is created
- **AND** no files are committed
- **AND** no PR is created
- **AND** the command exits with error code 1

#### Scenario: PR flag compatible with other archive flags

- **WHEN** user runs `spectr archive my-feature --pr --yes --skip-specs`
- **THEN** the archive operation skips confirmation prompts
- **AND** spec updates are skipped
- **AND** PR is created after successful archive
- **AND** the PR body notes that spec updates were skipped

### Requirement: Archive PR Branch Naming

The system SHALL create archive PR branches with the naming convention `archive-<change-id>` to clearly indicate the branch purpose and maintain consistency with change proposal branch naming.

#### Scenario: Branch name follows convention

- **WHEN** user archives a change named `user-authentication` with `--pr` flag
- **THEN** the created branch is named `archive-user-authentication`
- **AND** the branch is created from the current branch

#### Scenario: Branch name conflict handling

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** a branch named `archive-my-feature` already exists
- **THEN** an error is displayed: "Branch 'archive-my-feature' already exists"
- **AND** the archive operation completes successfully
- **AND** no new branch is created
- **AND** the command exits with error code 1

### Requirement: Archive PR Commit Strategy

The system SHALL commit all archive-related changes atomically in a single commit, including the archived directory, removal of the original change directory, and all updated spec files.

#### Scenario: Commit includes archived directory and updated specs

- **WHEN** archiving a change with `--pr` flag
- **AND** spec updates are not skipped
- **THEN** the commit includes the new archived directory at `spectr/changes/archive/YYYY-MM-DD-<change-id>/`
- **AND** the commit includes all updated files in `spectr/specs/`
- **AND** the removal of the original `spectr/changes/<change-id>/` directory is detected by git

#### Scenario: Commit when specs are skipped

- **WHEN** user runs `spectr archive my-feature --pr --skip-specs`
- **THEN** the commit includes only the archived directory
- **AND** the commit includes the removal of the original change directory
- **AND** no spec files are included in the commit

#### Scenario: Commit message format includes operation summary

- **WHEN** a commit is created for archive with PR
- **THEN** the commit message starts with "Archive: <change-id>"
- **AND** the body includes the archive location
- **AND** the body includes spec operation counts (added/modified/removed/renamed)
- **AND** the message ends with "Change-Id: <change-id>" trailer

### Requirement: Archive PR Platform Detection

The system SHALL detect the git hosting platform from the origin remote URL and invoke the appropriate PR creation CLI tool (gh for GitHub, glab for GitLab, tea for Gitea/Forgejo).

#### Scenario: GitHub platform detection

- **WHEN** the origin remote URL contains `github.com`
- **AND** user runs archive with `--pr` flag
- **THEN** the `gh pr create` command is used to create the PR
- **AND** the PR is created on GitHub

#### Scenario: GitLab platform detection

- **WHEN** the origin remote URL contains `gitlab.com` or matches a GitLab instance
- **AND** user runs archive with `--pr` flag
- **THEN** the `glab mr create` command is used to create the merge request
- **AND** the MR is created on GitLab

#### Scenario: Gitea platform detection

- **WHEN** the origin remote URL contains `gitea` or `forgejo`
- **AND** user runs archive with `--pr` flag
- **THEN** the `tea pr create` command is used to create the PR
- **AND** the PR is created on Gitea or Forgejo

#### Scenario: Platform detection fails

- **WHEN** the origin remote URL does not match any known platform
- **AND** user runs archive with `--pr` flag
- **THEN** an error is displayed with the remote URL
- **AND** the message guides user to create PR manually with gh, glab, or tea
- **AND** the archive completes successfully but PR is not created
- **AND** the branch is created and pushed

### Requirement: Archive PR Title and Body

The system SHALL generate a PR with a descriptive title and body that summarizes the archive operation, spec updates, and provides review guidance.

#### Scenario: PR title follows convention

- **WHEN** a PR is created for archiving change `my-feature`
- **THEN** the PR title is "Archive: my-feature"

#### Scenario: PR body includes archive summary

- **WHEN** a PR is created after archive
- **THEN** the PR body includes "Archived completed change: `<change-id>`"
- **AND** the body includes the archive location path
- **AND** the body includes spec operation counts
- **AND** the body lists updated capabilities
- **AND** the body includes review notes section
- **AND** the body footer notes "Generated by `spectr archive --pr`"

#### Scenario: PR body when specs skipped

- **WHEN** a PR is created with `--skip-specs` flag
- **THEN** the PR body includes "Spec updates skipped (--skip-specs flag used)"
- **AND** the spec operation counts section is omitted
- **AND** the updated capabilities list is omitted

### Requirement: Archive PR Error Handling

The system SHALL handle git operation errors gracefully, providing clear error messages and leaving the archive in a valid state even when PR creation fails.

#### Scenario: Not in git repository

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** the current directory is not in a git repository
- **THEN** an error is displayed: "Not in a git repository. Initialize git with 'git init'."
- **AND** the archive operation completes successfully
- **AND** changes remain uncommitted
- **AND** the command exits with error code 1

#### Scenario: Origin remote not configured

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** the git repository has no `origin` remote
- **THEN** an error is displayed: "No 'origin' remote configured. Run 'git remote add origin <url>' first."
- **AND** the archive operation completes successfully
- **AND** no branch is created
- **AND** the command exits with error code 1

#### Scenario: PR CLI tool not installed

- **WHEN** user runs `spectr archive my-feature --pr` for a GitHub repository
- **AND** the `gh` CLI tool is not installed
- **THEN** an error is displayed: "gh not found. Install from https://github.com/cli/cli"
- **AND** the archive operation completes successfully
- **AND** the branch is created and pushed
- **AND** user can manually create PR using gh after installation
- **AND** the command exits with error code 1

#### Scenario: Push fails due to network error

- **WHEN** pushing the branch fails due to network error
- **THEN** an error is displayed with the git error message
- **AND** the archive operation is complete
- **AND** the branch and commit exist locally
- **AND** user can manually retry push with `git push`
- **AND** the command exits with error code 1

#### Scenario: PR creation fails

- **WHEN** the PR CLI tool fails to create the PR
- **THEN** an error is displayed with the tool output
- **AND** the archive operation is complete
- **AND** the branch is created, committed, and pushed
- **AND** user can manually create PR via web UI or CLI
- **AND** the command exits with error code 1

### Requirement: Archive PR Success Reporting

The system SHALL display the PR URL after successful PR creation to provide immediate feedback and enable quick access to the created pull request.

#### Scenario: PR created successfully

- **WHEN** archive with `--pr` flag completes successfully
- **AND** the PR is created
- **THEN** a success message displays the PR URL
- **AND** the message format is "PR created: https://github.com/owner/repo/pull/123" (or equivalent for GitLab/Gitea)
- **AND** the command exits with code 0

#### Scenario: Display success after archive confirmation

- **WHEN** archive completes and PR is created
- **THEN** the PR URL is displayed after the "Successfully archived" message
- **AND** both success messages are clearly visible to the user
