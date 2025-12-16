# Change: Fix PR Archive Validation Inconsistency

## Why

When running `spectr pr archive`, validation can fail with "ADDED requirement already exists in base spec" even though `spectr validate --all` passes for the same change. This creates a confusing user experience where local validation gives false confidence, only to fail during the PR creation workflow.

The root cause is that:
- `spectr validate` validates against **local** base specs in `spectr/specs/`
- `spectr pr archive` validates against **origin/main** base specs (via git worktree)

When a requirement has already been merged to main (e.g., from another PR) but doesn't exist locally, the discrepancy causes the archive to fail unexpectedly.

## What Changes

1. **Add `--base-branch` flag to validate command** - Allow validation to check against a specific branch's specs instead of only local specs
2. **Pre-flight validation in PR workflow** - Before creating the worktree, validate the change against the target branch's specs to fail fast
3. **Improved error messages** - When this specific validation failure occurs, explain the local vs remote discrepancy

## Impact

- Affected specs: `validation`
- Affected code:
  - `cmd/validate.go` - Add `--base-branch` flag
  - `internal/validation/change_rules.go` - Support base branch spec resolution
  - `internal/pr/workflow.go` - Add pre-flight validation
  - `internal/git/` - Add helper to read files from specific branches
