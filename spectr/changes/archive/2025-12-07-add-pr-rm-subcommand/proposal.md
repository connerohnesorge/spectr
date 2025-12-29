# Change: Add PR Remove Subcommand

## Why

When a change proposal is abandoned, stale, or superseded, there is no convenient way to create a PR that cleanly removes it from the repository. Currently, users must manually:

1. Delete the change directory locally
2. Stage and commit the deletion
3. Create a branch and push
4. Create a PR through the platform CLI or web interface

A `spectr pr rm <change_id>` command would automate this workflow, providing the same git worktree isolation and platform detection as the existing `spectr pr archive` and `spectr pr proposal` commands.

## What Changes

- Add `spectr pr rm` (aliases: `r`, `remove`) subcommand to remove a change via PR
- Uses the same worktree-based isolation as other PR subcommands
- Branch naming follows existing pattern: `spectr/remove/<change-id>`
- Generates appropriate commit message and PR body for removal context
- Supports all common PR flags: `--base`, `--draft`, `--force`, `--dry-run`
- Interactive selection when no change ID is provided (consistent with other PR subcommands)
- Validates that the change exists before proceeding

## Impact

- Affected specs: `cli-interface`
- Affected code: `cmd/pr.go`, `cmd/pr_helpers.go`, `internal/pr/workflow.go`, `internal/pr/templates.go`, `internal/pr/helpers.go`
