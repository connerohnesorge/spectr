# Change: Add `spectr accept` command for tasks.md to tasks.json conversion

## Why

Inspired by Anthropic's research on [long-running agents](https://www.anthropic.com/research/effective-harnesses-for-long-running-agents), task lists in JSON format are significantly more stable for AI agents than Markdown. Agents are less likely to accidentally modify or overwrite structured JSON files compared to Markdown. This change introduces a `spectr accept` command that converts human-readable `tasks.md` to machine-stable `tasks.json` when a proposal is approved, then removes the original `tasks.md` to prevent drift. The `apply` slash command will require calling `spectr accept` first before implementation begins.

## What Changes

- **Add `spectr accept <change-id>` CLI command** - Converts tasks.md to tasks.json format
- **Define tasks.json schema** - Structured JSON format preserving sections, task IDs, descriptions, and completion status
- **Add task parser** - Parse tasks.md format into structured Task objects
  - Support unlimited recursive nesting (1.1.1.1...)
  - Append indented detail lines to task descriptions
  - Preserve existing `[x]` completion status
- **Remove tasks.md after conversion** - Delete original file to prevent drift between formats
- **Update apply slash command** - Instruct agent to run `spectr accept` if tasks.json missing
- **Archive existing tasks.md files** - Provide testdata fixtures from rich archive

## Impact

- Affected specs: `cli-framework`, `archive-workflow`
- Affected code:
  - `cmd/root.go` (new AcceptCmd)
  - `internal/accept/` (new package)
  - `internal/parsers/parsers.go` (enhanced task parsing)
  - `.claude/commands/spectr/apply.md` (updated workflow)
- Testdata: Use archived tasks.md files as test fixtures
- **Non-breaking**: tasks.md remains valid for proposal phase; conversion happens at acceptance
