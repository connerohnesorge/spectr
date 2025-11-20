# Change: Add `--pr` Flag to Archive Command for Automated PR Creation

## Why

After successfully archiving a change, users currently must manually create a git branch, stage the archived files and updated specs, commit, push, and open a pull request to merge these changes back to the main branch. This multi-step process is repetitive, error-prone, and interrupts the natural workflow of completing a change.

A `--pr` flag on `spectr archive` would streamline this workflow by automating the entire PR creation process after a successful archive operation. Similar to the `spectr propose` command pattern, it would detect the git hosting platform (GitHub, GitLab, Gitea) and use the appropriate CLI tool (`gh`, `glab`, `tea`) to create the PR, ensuring consistent workflows across teams.

## What Changes

- **NEW**: Add `--pr` flag to `spectr archive` command that:
  - Only activates after successful completion of the archive operation (validation, spec merging, directory move)
  - Creates a new git branch named `archive-<change-id>` from the current branch
  - Stages the archived directory (`spectr/changes/archive/YYYY-MM-DD-<change-id>/`) and all updated spec files
  - Commits with a descriptive message including the change ID and spec operation summary
  - Pushes the branch to the remote repository
  - Detects the git hosting platform (GitHub, GitLab, Gitea) from the remote URL
  - Invokes the appropriate PR CLI tool (`gh pr create`, `glab mr create`, `tea pr create`)
  - Displays the PR URL on success
  - **BREAKING**: None (new optional flag)

## Impact

- **Affected specs**: `archive-workflow`, `cli-interface`
- **Affected code**:
  - `cmd/archive.go` - Add `PR` flag field to `ArchiveCmd` struct
  - `internal/archive/archiver.go` - Add PR workflow logic after successful archive
  - `internal/archive/` - New git operations module (or reuse from propose package if it exists)
- **User-visible changes**: One new optional flag on existing command
- **Dependencies**:
  - Requires `git` available in PATH
  - Requires platform-specific PR CLI (`gh`, `glab`, or `tea`) for PR creation
  - Should coordinate with or reuse git detection/PR logic from `spectr propose` command (from add-propose-command)
