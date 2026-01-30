# Tasks

## 1. Core Sync Implementation

- [x] 1.1 Create `internal/sync/sync.go` with `SyncTasksToMarkdown` function
- [x] 1.2 Implement tasks.jsonc parsing to read current task statuses
- [x] 1.3 Implement tasks.md parsing that preserves original
  structure/formatting
- [x] 1.4 Implement status marker replacement (`[ ]` ↔ `[x]`) in-place
- [x] 1.5 Handle edge case: tasks.md doesn't exist (skip sync gracefully)
- [x] 1.6 Handle edge case: tasks.jsonc doesn't exist (skip sync gracefully)
- [x] 1.7 Add unit tests for sync logic with table-driven test cases

## 2. Kong BeforeRun Hook Integration

- [x] 2.1 Add `--no-sync` global flag to `cmd/root.go` CLI struct
- [x] 2.2 Implement `BeforeRun` method on CLI struct for Kong hook
- [x] 2.3 In BeforeRun, discover all active changes with tasks.jsonc files
- [x] 2.4 Call sync for each discovered change before command executes
- [x] 2.5 Respect `--no-sync` flag to skip sync when set
- [ ] 2.6 Add integration test verifying sync runs before commands

## 3. Status Matching Logic

- [x] 3.1 Match tasks by ID between tasks.jsonc and tasks.md lines
- [x] 3.2 Handle flexible task ID formats (1.1, 1., 1, no-id)
- [x] 3.3 Map `pending` → `[ ]`, `in_progress` → `[ ]`, `completed` → `[x]`
- [x] 3.4 Preserve line content after status marker unchanged
- [x] 3.5 Add tests for ID matching edge cases

## 4. Output and Verbosity

- [x] 4.1 Add `--verbose` global flag to CLI struct
- [x] 4.2 Silent by default: no output on successful sync
- [x] 4.3 With `--verbose`: print "Synced N task statuses in change-id"
- [x] 4.4 Print errors to stderr if sync fails (but don't block command)

## 5. Discovery Integration

- [x] 5.1 Use `internal/discovery` to find active changes with tasks.jsonc
- [x] 5.2 Exclude `spectr/changes/archive/` from sync scope
- [x] 5.3 Handle case where spectr/ directory doesn't exist (not initialized)

## 6. Documentation

- [ ] 6.1 Update `spectr/AGENTS.md` with sync behavior documentation
- [x] 6.2 Add `--no-sync` flag to CLI help text
- [ ] 6.3 Document in `CLAUDE.md` that tasks.md is auto-synced

## 7. Final Validation

- [x] 7.1 Run `spectr validate add-auto-sync-tasks`
- [x] 7.2 Run full test suite: `nix develop -c tests`
- [x] 7.3 Run linter: `nix develop -c lint`
- [x] 7.4 Manual test: modify tasks.jsonc status, run any spectr command, verify
  tasks.md updated
