# Change: Add Delegation Guidance to Instruction Pointer Template

## Why
When orchestrator agents delegate implementation tasks to coder subagents, the subagents often lack sufficient context about the change proposal. Including the explicit path to the change directory (proposal, spec deltas, and tasks) enables subagents to reference the authoritative specification rather than relying on incomplete task descriptions.

## What Changes
- Update `instruction-pointer.md.tmpl` to include guidance about referencing change proposal paths when delegating tasks
- Use template variables (`{{ .ChangesDir }}`) for dynamic directory names
- Add clear instruction for providing `<change-dir>/proposal.md`, `<change-dir>/tasks.json`, and delta spec paths to subagents

## Impact
- Affected specs: `agent-instructions`
- Affected code: `internal/initialize/templates/spectr/instruction-pointer.md.tmpl`
