# Change: Add spectr-validate-wo-spectr-bin AgentSkills Initializer

## Why

In sandboxed environments, CI pipelines, or fresh repository checkouts where the
spectr binary is not available, users currently have no way to validate their
specifications and change proposals. This creates friction in workflows where
the binary cannot be easily installed.

Following the pattern established by `spectr-accept-wo-spectr-bin`, we need a
validation skill that replicates the core validation logic using only bash and
standard tools (grep, sed, awk).

## What Changes

- Add new embedded skill `spectr-validate-wo-spectr-bin` under
  `internal/initialize/templates/skills/`
- Add `AgentSkillsInitializer` for the new skill in Claude Code provider
- Skill provides `scripts/validate.sh` implementing core validation checks:
  - Spec file validation (Requirements section, SHALL/MUST, scenarios)
  - Change delta validation (ADDED/MODIFIED/REMOVED/RENAMED sections)
  - Tasks.md validation (at least one task item)
  - Cross-file conflict detection (duplicates, missing base specs)

## Impact

- Affected specs: `support-claude-code`
- Affected code:
  - `internal/initialize/templates/skills/spectr-validate-wo-spectr-bin/` (new)
  - `internal/initialize/providers/claude.go` (add initializer)
