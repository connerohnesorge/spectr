# Change: Clean Up Leftover tasks.json After PR Operations

## Why
When `spectr pr archive` or `spectr pr rm` commands complete, the local change directory in the user's working directory is not cleaned up. This causes `tasks.json` files (and potentially other untracked files) to remain as orphans after the PR is merged and pulled. The user is left with stale files that should have been removed or archived.

## What Changes
- **Add local cleanup after PR rm**: After successfully creating a PR to remove a change, delete the local change directory to prevent orphan files
- **Add local cleanup after PR archive**: After successfully creating a PR to archive a change, delete the local change directory (the archived version will be pulled when the PR merges)
- **Add warning before cleanup**: Display a warning message before cleaning up local files, informing the user what will be removed
- **Add `--no-cleanup` flag**: Allow users to skip local cleanup if they want to preserve local changes

## Impact
- Affected specs: `archive-workflow`, `cli-interface`
- Affected code: `internal/pr/workflow.go`, `cmd/pr.go`
- Breaking changes: None (new behavior is additive, with opt-out flag)
