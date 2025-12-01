# Change: Add `--pr` Flag to Archive Command with Git Worktree Isolation

## Why

The previous `--pr` flag implementation was removed (2025-11-30-remove-archive-pr-flag) due to complexity concerns around mixing archive operations with the user's working directory state. Users who have uncommitted changes or are mid-development risk polluting their archive commits with unrelated work.

This proposal re-introduces the `--pr` flag using **git worktrees** to provide complete isolation. The archive operation executes in a fresh worktree on a new branch, ensuring:
- The user's working directory is never modified
- No risk of including uncommitted changes in the PR
- Clean, atomic commits containing only the archive results
- Safer rollback if any step fails

## What Changes

- **NEW**: Add `--pr` flag to `spectr archive` command that:
  1. Detects the git hosting platform (GitHub, GitLab, Gitea, Forgejo, Bitbucket) from remote URL
  2. Validates the required CLI tool is installed (`gh`, `glab`, `tea`, or standard git for Bitbucket)
  3. Creates a temporary git worktree on a new branch `archive-<change-id>`
  4. Executes `spectr archive <change-id> --yes` within the worktree
  5. Stages the entire `spectr/` directory (`git add spectr/`)
  6. Commits with a structured message including change ID and operation summary
  7. Pushes the branch to origin
  8. Creates a PR/MR using the platform-appropriate CLI
  9. Cleans up the temporary worktree
  10. Displays the PR URL on success

- **NEW**: Add `internal/git/` package with:
  - Platform detection from remote URLs
  - Worktree management (create, execute, cleanup)
  - PR CLI abstraction (gh, glab, tea)

## Impact

- **Affected specs**: `archive-workflow`
- **Affected code**:
  - `internal/archive/cmd.go` - Add `PR` flag field to `ArchiveCmd` struct
  - `internal/archive/archiver.go` - Add PR workflow orchestration after successful archive
  - `internal/git/` (NEW) - Git worktree operations, platform detection, PR CLI invocation
- **User-visible changes**: One new optional flag on existing command
- **Dependencies**:
  - Requires `git` (>= 2.5 for worktree support) available in PATH
  - Requires platform-specific PR CLI for PR creation:
    - GitHub: `gh` CLI
    - GitLab: `glab` CLI
    - Gitea/Forgejo: `tea` CLI
    - Bitbucket: Manual PR URL provided (no official CLI)
