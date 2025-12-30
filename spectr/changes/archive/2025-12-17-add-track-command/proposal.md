# Change: Add Track Command for Automatic Git Commits

## Why

During change implementation, developers manually commit after completing each
task. This is tedious and can lead to inconsistent commit history. An automated
tracking command that watches task status changes and auto-commits related
changes would streamline the workflow.

## What Changes

- New `spectr track [change-id]` command watches tasks.json for status updates
- When a task status changes to "in_progress" or "completed", automatically
  stages and commits modified files
- Excludes task files (tasks.json, tasks.jsonc, tasks.md) from staging
- Commit message format: `spectr(<change-id>): start|complete task <task-id>`
  with `[Automated by spectr track]` footer
- Warns user when no files to commit (only task file changed)
- Stops tracking immediately if git commit fails (e.g., merge conflict, hook
  rejection)
- Runs in foreground until all tasks complete or user interrupts (Ctrl+C)
- Interactive change selection if no change-id provided

## Impact

- Affected specs: cli-interface
- Affected code:
  - `cmd/track.go` - New command definition
  - `cmd/root.go` - Add Track field to CLI struct
  - `internal/track/` - New package for tracking logic
  - `internal/specterrs/` - Add track-specific error types
  - `go.mod` - Add fsnotify dependency
