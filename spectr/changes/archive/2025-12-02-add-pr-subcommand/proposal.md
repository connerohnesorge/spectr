# Change: Add `spectr pr` Subcommand with Git Worktree Isolation

## Why

Users need a streamlined way to create pull requests for Spectr changes without
polluting their working directory. The current workflow requires manual steps:
checkout new branch, archive/copy change, stage, commit, push, create PR. This
is error-prone and tedious.

A dedicated `spectr pr` subcommand with `archive` and `new` variants provides:

- Complete isolation via git worktrees - user's working directory is never
  modified
- Automatic platform detection (GitHub, GitLab, Gitea, Forgejo, Bitbucket)
- Structured commit messages and PR bodies with change metadata
- Single command to go from completed change to PR under review

This is distinct from a `--pr` flag on archive because it provides a dedicated
namespace for PR operations with multiple subcommands (`archive` for completed
changes, `proposal` for in-progress changes).

## What Changes

- **NEW**: Add `spectr pr` top-level command with two subcommands:
  - `spectr pr archive <change-id>` - Archive a change and create PR (MUST ALSO
    SUPPORT partial ids similar to `archive`)
  - `spectr pr proposal <change-id>` - Copy change to PR branch without
    archiving

- **NEW**: Both subcommands share common PR workflow:
  1. Detect git hosting platform from `origin` remote URL
  2. Validate required CLI tool is installed (`gh`, `glab`, `tea`)
  3. Create temporary git worktree on branch `spectr/<change-id>`
  4. Execute change operation within worktree:
     - `archive`: Run `spectr archive <change-id> --yes`
     - `proposal`: Copy `spectr/changes/<change-id>/` to worktree
  5. Stage the entire `spectr/` directory (`git add spectr/`)
  6. Commit with structured message including change metadata
  7. Push worktree branch to origin
  8. Create PR/MR from worktree using platform-appropriate CLI
  9. Clean up temporary worktree
  10. Display PR URL on success

- **NEW**: Add `internal/git/` package with:
  - Platform detection from remote URLs
  - Worktree management (create, execute in, cleanup)
  - PR CLI abstraction (gh, glab, tea, bitbucket manual URL)

- **NEW**: Add `internal/pr/` package with:
  - PR workflow orchestration
  - Commit/PR message templating
  - Change metadata extraction

## Impact

- **Affected specs**: `cli-interface` (new command)
- **Affected code**:
  - `cmd/root.go` - Add `PR` command to CLI struct
  - `cmd/pr.go` (NEW) - PR command with archive/new subcommands
  - `internal/git/` (NEW) - Git operations package
  - `internal/pr/` (NEW) - PR workflow package
- **User-visible changes**: New top-level `spectr pr` command with subcommands
- **Dependencies**:
  - Requires `git` (>= 2.5 for worktree support) available in PATH
  - Requires platform-specific PR CLI for PR creation:
    - GitHub: `gh` CLI
    - GitLab: `glab` CLI
    - Gitea/Forgejo: `tea` CLI
    - Bitbucket: Manual PR URL provided (no official CLI)
