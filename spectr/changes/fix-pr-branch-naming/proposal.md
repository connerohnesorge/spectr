# Change: Fix PR Branch Naming by Subcommand Type

## Why

Currently, both `spectr pr archive` and `spectr pr new` create branches with the same naming pattern `spectr/<change-id>`. This makes it difficult to distinguish between branches created for archived changes versus proposal reviews when viewing git history or branch lists. The branch name should convey the purpose of the PR at a glance.

## What Changes

- `spectr pr archive` creates branches named `spectr/archive/<change-id>`
- `spectr pr new` creates branches named `spectr/proposal/<change-id>`
- Update branch existence checks to use mode-specific prefix
- Update existing spec scenarios to reflect new branch naming

## Impact

- Affected specs: cli-interface (adds PR Branch Naming Convention requirement)
- Affected code: `internal/pr/workflow.go:99` (branch name generation in `prepareWorkflowContext`)
- Related changes: Supersedes branch naming in `add-pr-subcommand` change which uses `spectr/<change-id>` format
- **BREAKING**: Users with existing remote branches named `spectr/<change-id>` will need to use `--force` to recreate them with the new naming convention, or delete them manually
