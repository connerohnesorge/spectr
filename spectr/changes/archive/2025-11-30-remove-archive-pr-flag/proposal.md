# Change: Remove --pr flag from archive command

## Why

The `--pr` flag on `spectr archive` adds significant complexity (worktree management, platform detection, multi-CLI tool support) for a feature that duplicates functionality readily available through standard git workflows and AI assistants. Users can achieve the same result with `git add && git commit && gh pr create` which is more transparent and controllable.

## What Changes

- **REMOVED** `--pr` flag from `spectr archive` command
- **REMOVED** `internal/archive/pr.go` - PR creation orchestration
- **REMOVED** `internal/archive/pr_format.go` - commit/PR message formatting
- **REMOVED** `internal/archive/pr_test.go` - PR-related tests
- **REMOVED** `internal/git/` package entirely (only used for PR functionality):
  - `pr.go` - platform-specific PR creation (gh, glab, tea)
  - `platform.go` - git hosting platform detection
  - `platform_test.go` - platform detection tests
  - `operations.go` - worktree and branch utilities
- **REMOVED** `github.com/google/uuid` dependency (only used for branch name generation)
- **REMOVED** PR-related requirements from `archive-workflow` spec
- **REMOVED** PR flag requirement from `cli-interface` spec

## Impact

- Affected specs: `archive-workflow`, `cli-interface`
- Affected code: `internal/archive/`, `internal/git/`, `cmd/`, `go.mod`
- No breaking changes for users not using `--pr` flag
- Users who relied on `--pr` can use standard git commands instead
