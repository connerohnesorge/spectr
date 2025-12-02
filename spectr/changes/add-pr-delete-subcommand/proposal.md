# Change: Add `spectr pr delete` Subcommand

## Why

Users need a quick way to delete specs from the codebase via PR without manually removing directories. The current PR subcommands (`archive` and `new`) focus on creating or promoting changes, but there's no streamlined way to remove a spec when it's no longer needed.

A dedicated `spectr pr delete` (with `d` shorthand) subcommand provides:
- Single command to create a PR that removes an entire spec directory
- Complete isolation via git worktrees - user's working directory is never modified
- Supports partial change IDs consistent with other `spectr pr` subcommands
- Clear commit messages and PR bodies explaining the deletion

## What Changes

- **NEW**: Add `spectr pr delete <spec-id>` subcommand:
  - Alias: `spectr pr d <spec-id>`
  - Creates a worktree branch `spectr/delete-<spec-id>`
  - Removes the entire `spectr/specs/<spec-id>/` directory in the worktree
  - Stages deletion, commits with structured message, pushes, creates PR
  - Cleans up worktree automatically

- **NEW**: Add `PRDeleteCmd` struct in `cmd/pr.go`:
  - Shares common flags with other PR subcommands: `--base`, `--draft`, `--force`, `--dry-run`
  - Supports partial spec ID resolution (prefix and substring matching)

- **NEW**: Add `ModeDelete` constant and delete workflow in `internal/pr/`:
  - `executeDeleteInWorktree()` function to `rm -rf` the spec directory
  - Delete-specific commit and PR templates

## Impact

- **Affected specs**: `cli-interface` (new subcommand)
- **Affected code**:
  - `cmd/pr.go` - Add `Delete` subcommand to `PRCmd` struct
  - `internal/pr/workflow.go` - Add delete mode handling
  - `internal/pr/helpers.go` - Add `executeDeleteInWorktree()` function
  - `internal/pr/templates.go` - Add delete-specific templates
  - `internal/discovery/specs.go` - May need spec ID resolution (similar to change ID resolution)
- **User-visible changes**: New `spectr pr delete` and `spectr pr d` commands
- **Dependencies**: Same as existing `spectr pr` commands (git, platform CLIs)
