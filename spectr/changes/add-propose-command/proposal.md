# Change: Add `spectr propose` Command for Automated PR Creation

## Why

Creating change proposals is a core part of the Spectr workflow, but once a proposal is scaffolded and validated, users currently must manually create a git branch, commit, push, and open a PR. This multi-step process is repetitive and error-prone.

A `spectr propose <id>` command would streamline this workflow by automating the entire PR creation process in a single command. It detects the git hosting platform (GitHub, GitLab, Gitea) and uses the appropriate CLI tool (`gh`, `glab`, `tea`), reducing friction and ensuring consistent PR workflows across teams.

## What Changes

- **NEW**: Add `spectr propose <id>` command that:
  - Validates that `spectr/changes/<id>/` exists and is not yet committed to git
  - Creates a new git branch named `add-<id>` from the current branch
  - Stages only the `spectr/changes/<id>/` folder
  - Commits with a descriptive message including the change ID
  - Pushes the branch to the remote
  - Detects the git hosting platform (GitHub, GitLab, Gitea)
  - Invokes the appropriate PR CLI (`gh pr create`, `glab mr create`, `tea pr create`)
  - Displays the PR URL on success
  - **BREAKING**: None

## Impact

- **Affected specs**: `cli-interface`, `cli-framework`
- **Affected code**:
  - `cmd/` - Add `ProposeCmd` command handler
  - `internal/` - New `propose` package with git and PR operations
  - `main.go` - Register `ProposalCmd` in CLI
- **User-visible changes**: One new command
- **Dependencies**: Requires `git` and platform-specific PR CLI (`gh`, `glab`, or `tea`) installed
