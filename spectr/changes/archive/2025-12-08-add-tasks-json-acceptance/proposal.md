# Change: Add Tasks JSON Acceptance Command

## Why

Inspired by Anthropic's research on [long-running
agents](https://www.anthropic.com/research/effective-harnesses-for-long-running-agents),
task lists in JSON format are significantly more stable for AI agents than
Markdown. Agents are less likely to accidentally modify or overwrite structured
JSON files compared to Markdown. This is because:

1. **Structural validation**: JSON parsers will immediately fail if the format
  is corrupted, whereas Markdown changes can silently corrupt task lists
2. **Atomic field updates**: Agents can update a single `status` field in JSON
  without risk of rewriting the entire task description
3. **Machine readability**: JSON provides unambiguous field boundaries,
  eliminating regex-based parsing that can fail on edge cases
4. **Drift prevention**: Once `tasks.md` is converted to `tasks.json`, the
  original is removed, ensuring a single source of truth

## What Changes

- Add a new `spectr accept <change-id>` command that converts `tasks.md` to
  `tasks.json` format
- The `accept` command validates the change before conversion (using existing
  validation)
- After successful conversion, `tasks.md` is removed to prevent drift between
  formats
- The `apply` slash command is updated to require agents to run `spectr accept`
  first before implementation begins
- All internal tooling (`parsers.CountTasks`, `view`, `list`, etc.) is updated
  to read from `tasks.json` when present, falling back to `tasks.md` for
  backward compatibility

## Impact

- Affected specs: `cli-framework`, `archive-workflow`
- Affected code:
  - `cmd/accept.go` (new)
  - `cmd/root.go` (add AcceptCmd)
  - `internal/parsers/parsers.go` (update CountTasks for JSON support)
  - `internal/archive/archiver.go` (update task counting)
  - `internal/list/lister.go` (update task counting)
  - `internal/view/dashboard.go` (update task counting)
  - `.claude/commands/spectr/apply.md` (add accept step)
  - Template files for slash commands
