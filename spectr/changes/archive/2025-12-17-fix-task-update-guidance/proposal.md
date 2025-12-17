# Change: Fix Task Update Guidance for Incremental Updates

## Why

Current prompts and JSONC header have contradictory guidance about when to update task statuses. Some say "after all work is done" while others say "after each task". This leads to batch updates instead of incremental tracking, defeating the purpose of real-time progress visibility.

## What Changes

- Update `tasksJSONHeader` in `cmd/accept_writer.go` to explicitly state "update IMMEDIATELY after each task"
- Fix `.claude/commands/spectr/apply.md` to remove "after all work is done" language
- Fix `spectr/AGENTS.md` to remove "ensure every item is finished before updating" contradiction
- Fix `.claude/agents/coder.md` to reference `tasks.jsonc` (not `tasks.md`) and add incremental update guidance

## Impact

- Affected specs: agent-instructions (workflow guidance)
- Affected code: cmd/accept_writer.go, .claude/commands/spectr/apply.md, spectr/AGENTS.md, .claude/agents/coder.md
