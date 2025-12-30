## Context

Spectr provides AI agents with structured workflows through slash commands and
instruction files. However, the `spectr accept` command requires the spectr
binary, which may not be available in all environments.

The AgentSkills specification (https://agentskills.io) defines a portable
format for AI assistant capabilities. By creating a skill that replicates
`spectr accept` functionality using standard Unix tools, agents can work
effectively even without the spectr binary.

**Constraints:**
- Must follow AgentSkills specification for skill structure
- Must use `embed.FS` for packaging skill files
- Must be generic enough to support future skills
- Accept script must work with standard Unix tools (bash, sed, awk, jq)

## Goals / Non-Goals

**Goals:**
- Create a reusable `AgentSkillsInitializer` that can embed any skill directory
- Provide a working `spectr-accept-wo-spectr-bin` skill for Claude Code
- Convert tasks.md to tasks.jsonc without requiring the spectr binary
- Follow existing initializer patterns for consistency

**Non-Goals:**
- Full `spectr accept` feature parity (validation, hierarchical tasks)
- Supporting skills for non-Claude providers in this change
- Runtime skill discovery or loading

## Decisions

### Decision: Generic AgentSkillsInitializer

Create a generic initializer that takes a skill name and target directory,
copying all files from the embedded skill template.

**Rationale:** This allows future skills to be added without code changes to
the initializer itself - just add new embedded directories.

**Alternatives considered:**
- Skill-specific initializer: Would require new code for each skill
- Template-based generation: Overkill for static skill files

### Decision: Skill Directory Structure

Embed skills under `internal/initialize/templates/skills/<skill-name>/`:

```
templates/skills/
└── spectr-accept-wo-spectr-bin/
    ├── SKILL.md
    └── scripts/
        └── accept.sh
```

**Rationale:** Keeps all templates together; mirrors the target structure.

### Decision: Bash-based Accept Script

Use bash with jq for JSON generation rather than pure sed/awk.

**Rationale:**
- jq is widely available and handles JSON properly
- Avoids complex escaping issues with pure text manipulation
- Produces valid JSONC output reliably

## Risks / Trade-offs

- **jq dependency**: Script requires jq. Mitigation: Document requirement in
  SKILL.md compatibility field.
- **Incomplete parsing**: Bash script may not handle all tasks.md edge cases.
  Mitigation: Document limitations; recommend full spectr binary for complex
  cases.
- **Maintenance burden**: Two implementations of accept logic. Mitigation:
  The script is deliberately minimal - just task conversion.

## Open Questions

- Should we add a `--no-validate` flag to the script for speed?
- Should we support other providers (Gemini, etc.) with similar skills?
