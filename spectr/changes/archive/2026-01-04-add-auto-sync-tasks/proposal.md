# Change: Add Automatic tasks.jsonc → tasks.md Sync

## Why

AI agents and tools update `tasks.jsonc` during implementation (marking tasks
as `in_progress` or `completed`), but the human-readable `tasks.md` file
becomes stale. Users reviewing progress see outdated checkbox states in
`tasks.md`, creating confusion about actual task status. Manually keeping both
files in sync is tedious and error-prone.

## What Changes

- Add a **pre-command sync subroutine** that runs before every spectr
  subcommand, synchronizing task statuses from `tasks.jsonc` back to `tasks.md`
- Use **Kong's BeforeRun hook** pattern to implement this idiomatically
- Sync applies to **active changes only** (excludes archived changes)
- **Preserve markdown formatting**: update only `[ ]`/`[x]` markers, keeping
  comments, links, and structure intact
- Add **global `--no-sync` flag** to skip sync when needed
- **Silent by default**: no output unless errors occur (or `--verbose` passed)
- **tasks.jsonc is source of truth**: if task structure conflicts exist,
  tasks.jsonc wins

## Impact

- Affected specs: `cli-interface`
- Affected code:
  - `main.go` - Add Kong hook configuration
  - `cmd/root.go` - Add `--no-sync` global flag to CLI struct
  - New `internal/sync/` package - Sync logic implementation
- No breaking changes - behavior is additive and can be disabled
