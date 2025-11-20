# Archive Workflow Specification

## Purpose

This specification defines the workflow for archiving completed changes, including moving them to dated archive directories, validating change structure, parsing delta operations, merging deltas into base specs, and ensuring atomic spec updates.

## Requirements

### Requirement: Change Archive Directory Structure
The system SHALL archive completed changes to `spectr/changes/archive/YYYY-MM-DD-<change-id>/` where YYYY-MM-DD is the current date in ISO 8601 format.

#### Scenario: Archive with date prefix
- **WHEN** archiving a change on 2025-11-17
- **THEN** the change is moved to `spectr/changes/archive/2025-11-17-<change-id>/`

#### Scenario: Prevent duplicate archives
- **WHEN** an archive with the same name already exists
- **THEN** the system returns an error and does not overwrite the existing archive

### Requirement: Pre-Archive Validation
The system SHALL validate changes before archiving to ensure spec consistency.

#### Scenario: Validate proposal structure
- **WHEN** validating a change for archive
- **THEN** the system checks proposal.md structure and reports informational warnings

#### Scenario: Validate delta specs strictly
- **WHEN** validating delta specs
- **THEN** the system checks for required delta operations, scenario formatting, and blocks archive on errors

#### Scenario: Detect duplicate requirements within sections
- **WHEN** a delta spec has duplicate requirement names in the same section
- **THEN** the system returns a validation error with the duplicate requirement name

#### Scenario: Detect cross-section conflicts
- **WHEN** a requirement appears in multiple delta sections (e.g., ADDED and MODIFIED)
- **THEN** the system returns a validation error indicating the conflict

#### Scenario: Validate scenarios are properly formatted
- **WHEN** a requirement lacks properly formatted scenarios
- **THEN** the system returns an error requiring at least one `#### Scenario:` block per requirement

### Requirement: Task Completion Checking
The system SHALL check task completion status and warn users before archiving.

#### Scenario: Display task status
- **WHEN** archiving a change
- **THEN** the system displays task completion status (e.g., "3/5 complete")

#### Scenario: Warn on incomplete tasks
- **WHEN** a change has incomplete tasks
- **THEN** the system warns the user and requires confirmation to proceed (unless --yes flag is provided)

#### Scenario: Proceed with incomplete tasks when confirmed
- **WHEN** user confirms archiving despite incomplete tasks
- **THEN** the system proceeds with the archive operation

### Requirement: Delta Spec Discovery
The system SHALL find all delta specifications in a change directory for processing.

#### Scenario: Find delta specs in change
- **WHEN** preparing to archive a change
- **THEN** the system scans `spectr/changes/<id>/specs/*/spec.md` for delta specifications

#### Scenario: Identify corresponding main specs
- **WHEN** delta specs are found
- **THEN** the system maps them to `spectr/specs/*/spec.md` paths based on capability directory name

#### Scenario: Determine create vs update status
- **WHEN** mapping delta specs to main specs
- **THEN** the system checks if each main spec exists and marks it as "create" or "update"

### Requirement: Delta Operation Parsing
The system SHALL parse delta operations from change spec files following a strict format.

#### Scenario: Parse ADDED requirements
- **WHEN** a delta spec contains `## ADDED Requirements` section
- **THEN** the system extracts all requirement blocks with headers and scenarios

#### Scenario: Parse MODIFIED requirements
- **WHEN** a delta spec contains `## MODIFIED Requirements` section
- **THEN** the system extracts complete modified requirement blocks

#### Scenario: Parse REMOVED requirements
- **WHEN** a delta spec contains `## REMOVED Requirements` section
- **THEN** the system extracts requirement names to be removed

#### Scenario: Parse RENAMED requirements
- **WHEN** a delta spec contains `## RENAMED Requirements` section with FROM/TO pairs
- **THEN** the system extracts the old and new requirement names

#### Scenario: Require at least one delta operation
- **WHEN** a delta spec has no ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** the system returns an error indicating no delta operations were found

### Requirement: Delta Operation Application Order
The system SHALL apply delta operations in the order: RENAMED → REMOVED → MODIFIED → ADDED to ensure correct merging.

#### Scenario: Apply RENAMED first
- **WHEN** applying delta operations
- **THEN** RENAMED operations are applied first to update requirement names before other operations

#### Scenario: Apply REMOVED second
- **WHEN** applying delta operations
- **THEN** REMOVED operations are applied after RENAMED to delete requirements

#### Scenario: Apply MODIFIED third
- **WHEN** applying delta operations
- **THEN** MODIFIED operations are applied after REMOVED to update existing requirements

#### Scenario: Apply ADDED last
- **WHEN** applying delta operations
- **THEN** ADDED operations are applied last to insert new requirements

### Requirement: Requirement Name Normalization
The system SHALL normalize requirement names by trimming whitespace and using case-insensitive matching to prevent duplicates.

#### Scenario: Normalize whitespace
- **WHEN** comparing requirement names
- **THEN** the system trims leading and trailing whitespace

#### Scenario: Case-insensitive matching
- **WHEN** comparing requirement names
- **THEN** the system uses lowercase comparison to match requirements

#### Scenario: Prevent duplicates due to formatting
- **WHEN** requirements differ only in whitespace or case
- **THEN** the system treats them as the same requirement

### Requirement: Spec Merging Algorithm
The system SHALL merge delta operations into base specs while preserving requirement ordering.

#### Scenario: Preserve original requirement order
- **WHEN** merging modified or renamed requirements
- **THEN** the system maintains their original position in the requirements list

#### Scenario: Append new requirements at end
- **WHEN** adding new requirements
- **THEN** the system appends them after all existing requirements

#### Scenario: Rebuild requirements section
- **WHEN** merging is complete
- **THEN** the system reconstructs the spec with proper markdown structure and spacing

### Requirement: New Spec Creation
The system SHALL create new spec files when a delta spec has no corresponding main spec.

#### Scenario: Generate spec skeleton for new specs
- **WHEN** creating a new spec file
- **THEN** the system generates a skeleton with title, purpose placeholder, and requirements section

#### Scenario: Restrict operations for new specs
- **WHEN** creating a new spec
- **THEN** only ADDED operations are allowed (MODIFIED/REMOVED/RENAMED return an error)

#### Scenario: Error on invalid operations for new specs
- **WHEN** a delta spec for a new capability contains MODIFIED/REMOVED/RENAMED operations
- **THEN** the system returns an error indicating only ADDED is allowed for new specs

### Requirement: Post-Merge Validation
The system SHALL validate rebuilt specs before writing them to ensure correctness.

#### Scenario: Validate rebuilt spec structure
- **WHEN** a spec has been rebuilt from delta operations
- **THEN** the system validates the spec structure before writing

#### Scenario: Prevent writing invalid specs
- **WHEN** post-merge validation fails
- **THEN** the system aborts the archive operation without writing any specs

#### Scenario: Display validation errors
- **WHEN** post-merge validation fails
- **THEN** the system displays detailed validation errors with file paths and issue descriptions

### Requirement: Atomic Spec Updates
The system SHALL prepare all spec updates before writing to ensure all-or-nothing consistency.

#### Scenario: Prepare all updates first
- **WHEN** processing multiple delta specs
- **THEN** the system builds and validates all updated specs before writing any files

#### Scenario: Abort on first validation error
- **WHEN** any spec fails validation during preparation
- **THEN** the system aborts without writing any spec files

#### Scenario: Write all specs after validation
- **WHEN** all specs pass validation
- **THEN** the system writes all updated specs to their main spec locations

### Requirement: Spec Update Display
The system SHALL display spec update operations with clear formatting and counts.

#### Scenario: Display specs to be updated
- **WHEN** spec updates are prepared
- **THEN** the system lists each spec with "create" or "update" status

#### Scenario: Display operation counts per spec
- **WHEN** applying spec updates
- **THEN** the system displays counts: "+ N added, ~ N modified, - N removed, → N renamed" per spec

#### Scenario: Display total operation counts
- **WHEN** all spec updates are complete
- **THEN** the system displays totals across all specs

### Requirement: Spec Update Confirmation
The system SHALL prompt for confirmation before applying spec updates unless --yes flag is provided.

#### Scenario: Prompt for spec update confirmation
- **WHEN** spec updates are ready to apply
- **THEN** the system prompts "Proceed with spec updates?" and waits for user response

#### Scenario: Skip prompt with yes flag
- **WHEN** --yes flag is provided
- **THEN** the system proceeds with spec updates without prompting

#### Scenario: Skip updates when user declines
- **WHEN** user declines spec update prompt
- **THEN** the system skips spec updates but continues with archive operation

### Requirement: Archive Success Reporting
The system SHALL display clear success messages after archiving.

#### Scenario: Display archive success
- **WHEN** archive operation completes successfully
- **THEN** the system displays "Change '<change-id>' archived as '<archive-name>'"

#### Scenario: Display spec update success
- **WHEN** spec updates are applied
- **THEN** the system displays "Specs updated successfully" after showing operation counts

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
