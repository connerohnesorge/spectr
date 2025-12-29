# Change: Add Delegation/Completion Guidance to Instruction Pointer Template

## Why

When orchestrators delegate implementation tasks to coder subagents or
task-completing agents work through tasks, the agents often lack sufficient
context about the change proposal. Including explicit path references to the
change directory (proposal, spec deltas, and tasks) in the instruction pointer
enables subagents to reference the authoritative specification rather than
relying on incomplete task descriptions passed through delegation prompts.

## What Changes

- Update `instruction-pointer.md.tmpl` to include guidance about referencing
  change proposal paths when delegating or completing tasks
- Use template variables (`{{ .ChangesDir }}`) for dynamic directory names
- Add clear instruction for providing `<change-dir>/proposal.md`,
  `<change-dir>/tasks.json`, and delta spec paths to subagents

## Impact

- Affected specs: `agent-instructions`
- Affected code:
  `internal/initialize/templates/spectr/instruction-pointer.md.tmpl`
