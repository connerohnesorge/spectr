# Change: Add Validation Skill with Spectr Binary

## Why
Agents operating in environments where the `spectr` binary is available (e.g., local development, fully configured containers) should use the binary for validation instead of the bash-script fallback. The binary offers better performance, full feature parity (parallelism, cross-capability checks), and is the source of truth.

## What Changes
- Add a new agent skill `spectr-validate-w-spectr-bin` that wraps the `spectr validate` command.
- Update Claude Code provider to install this skill.
- Update Codex provider to install this skill.
- The new skill will reside in `.claude/skills/spectr-validate-w-spectr-bin` and `.codex/skills/spectr-validate-w-spectr-bin`.

## Impact
- **Affected Specs**: `support-claude-code`, `support-codex`
- **Affected Code**: `internal/initialize/providers/claude.go`, `internal/initialize/providers/codex.go`, `internal/initialize/templates/skills/`
