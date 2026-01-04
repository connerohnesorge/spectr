# Tasks

## 1. Core Sync Implementation

- [ ] 1.1 Create `internal/sync/sync.go` with `SyncTasksToMarkdown` function
- [ ] 1.2 Implement tasks.jsonc parsing to read current task statuses
- [ ] 1.3 Implement tasks.md parsing that preserves original structure/formatting
- [ ] 1.4 Implement status marker replacement (`[ ]` ↔ `[x]`) in-place
- [ ] 1.5 Handle edge case: tasks.md doesn't exist (skip sync gracefully)
- [ ] 1.6 Handle edge case: tasks.jsonc doesn't exist (skip sync gracefully)
- [ ] 1.7 Add unit tests for sync logic with table-driven test cases

## 2. Kong BeforeRun Hook Integration

- [ ] 2.1 Add `--no-sync` global flag to `cmd/root.go` CLI struct
- [ ] 2.2 Implement `BeforeRun` method on CLI struct for Kong hook
- [ ] 2.3 In BeforeRun, discover all active changes with tasks.jsonc files
- [ ] 2.4 Call sync for each discovered change before command executes
- [ ] 2.5 Respect `--no-sync` flag to skip sync when set
- [ ] 2.6 Add integration test verifying sync runs before commands

## 3. Status Matching Logic

- [ ] 3.1 Match tasks by ID between tasks.jsonc and tasks.md lines
- [ ] 3.2 Handle flexible task ID formats (1.1, 1., 1, no-id)
- [ ] 3.3 Map `pending` → `[ ]`, `in_progress` → `[ ]`, `completed` → `[x]`
- [ ] 3.4 Preserve line content after status marker unchanged
- [ ] 3.5 Add tests for ID matching edge cases

## 4. Output and Verbosity

- [ ] 4.1 Add `--verbose` global flag to CLI struct
- [ ] 4.2 Silent by default: no output on successful sync
- [ ] 4.3 With `--verbose`: print "Synced N task statuses in change-id"
- [ ] 4.4 Print errors to stderr if sync fails (but don't block command)

## 5. Discovery Integration

- [ ] 5.1 Use `internal/discovery` to find active changes with tasks.jsonc
- [ ] 5.2 Exclude `spectr/changes/archive/` from sync scope
- [ ] 5.3 Handle case where spectr/ directory doesn't exist (not initialized)

## 6. Documentation

- [ ] 6.1 Update `spectr/AGENTS.md` with sync behavior documentation
- [ ] 6.2 Add `--no-sync` flag to CLI help text
- [ ] 6.3 Document in `CLAUDE.md` that tasks.md is auto-synced

## 7. Final Validation

- [ ] 7.1 Run `spectr validate add-auto-sync-tasks`
- [ ] 7.2 Run full test suite: `nix develop -c tests`
- [ ] 7.3 Run linter: `nix develop -c lint`
- [ ] 7.4 Manual test: modify tasks.jsonc status, run any spectr command, verify tasks.md updated
