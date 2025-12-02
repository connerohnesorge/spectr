# Change: Remove spectr list command from agent prompts

## Why
AI agents can directly read directories (`spectr/changes/`, `spectr/specs/`) using `ls`, file reads, or `rg` without needing the `spectr list` CLI command. Instructing agents to use `spectr list` is unnecessary and adds cognitive overhead when agents already have more flexible tools at their disposal.

## What Changes
- Remove all references to `spectr list` and `spectr list --specs` from agent workflow prompts
- Replace with direct file/directory access instructions where context is needed
- Update templates that generate these prompts for downstream projects
- Keep `spectr list` references in user-facing documentation (README.md, docs/) since humans benefit from the formatted output

## Impact
- Affected specs: agent-instructions (new capability documenting agent prompt conventions)
- Affected code:
  - `spectr/AGENTS.md`
  - `.agent/workflows/spectr-proposal.md`
  - `.agent/workflows/spectr-sync.md`
  - `.agent/workflows/spectr-apply.md`
  - `.claude/commands/spectr/proposal.md`
  - `.claude/commands/spectr/sync.md`
  - `.claude/commands/spectr/apply.md`
  - `.gemini/commands/spectr/proposal.toml`
  - `.gemini/commands/spectr/sync.toml`
  - `.gemini/commands/spectr/apply.toml`
  - `internal/initialize/templates/spectr/AGENTS.md.tmpl`
  - `internal/initialize/templates/tools/slash-apply.md.tmpl`
  - `internal/initialize/templates/tools/slash-sync.md.tmpl`
