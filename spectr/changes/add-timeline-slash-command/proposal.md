# Change: Add Timeline Slash Command for Implementation Planning

## Why

Teams implementing multiple interdependent changes need visibility into the
optimal implementation order and dependency graph. Currently, there's no way
to quickly see which changes can be parallelized, which must wait on others,
and what order maximizes progress. The `/spectr:timeline` slash command fills
this gap by analyzing all changes and generating a structured timeline with
dependency information and implementation notes.

## What Changes

- **ADDED**: New `/spectr:timeline` slash command for AI agents
  - Discovers all active changes in `spectr/changes/` (excluding archive)
  - Parses proposal metadata (frontmatter with `requires`/`enables`)
  - Analyzes dependency graph and detects circular dependencies
  - Generates `./spectr/timeline.json` with:
    - Dependency graph structure
    - Optimal implementation order
    - Change metadata (ID, title, task counts, risk notes)
    - Implementation notes and parallelization opportunities
  - Provides both machine-readable JSON and human-friendly formatting
- **ADDED**: Skill definition in `.agents/skills/spectr-timeline/` with SKILL.md
- **ADDED**: Timeline analysis in `timeline.json`:
  - Phase-based ordering (parallel vs sequential)
  - Dependency relationships with reasons
  - Critical path analysis
  - Implementation notes for coordination
- **MODIFIED**: `/spectr/AGENTS.md` to document the new command and expected timeline output

## Impact

- **Affected specs**: `slash-commands` (new command), `cli-interface` (informational)
- **Affected code**:
  - New `.agents/skills/spectr-timeline/SKILL.md` - command definition and instructions
  - New `spectr/timeline.json` output file - consumed by teams
  - Existing dependency resolution logic reused from chained-proposals feature
  - Optional: New `cmd/timeline.go` and `internal/timeline/` if implementing as CLI
    (not required for slash command version)
- **Breaking changes**: None - slash command is purely additive

## Scope Notes

This proposal focuses on the **slash command/skill implementation** that:

1. Reads existing proposal metadata
2. Analyzes dependencies
3. Produces `timeline.json` with human-friendly output

A future proposal could implement `/spectr timeline` as a CLI command if needed.
