# Design: Amp Agent Skills Support

## Context

Amp is based on Claude Code but has adopted `.agents/skills/` as the standard location for agent skills. As of January 29, 2026, Amp removed custom commands entirely in favor of user-invocable skills.

## Key Architectural Decisions

### 1. Skills Over Commands

**Decision:** Generate agent skills (`.agents/skills/`) instead of slash commands.

**Rationale:**
- Amp has deprecated custom commands in favor of skills
- User-invocable skills allow both agent and user to invoke `/spectr:proposal`
- Skills support lazy-loading of instructions for better context efficiency
- Aligns with Amp's current architecture and roadmap

### 2. SKILL.md Frontmatter Format

**Decision:** Use Amp's minimal frontmatter format with `name` and `description`.

```yaml
---
name: spectr-proposal
description: Create a Spectr change proposal with delta specs and tasks
---
```

**Rationale:**
- Amp requires only `name` and `description` in frontmatter
- Simpler than Claude Code's extended frontmatter (no `allowed-tools`, `context`, etc.)
- Follows Amp's conventions for skill discovery

### 3. Skill Directory Structure

**Decision:** Generate skills in `.agents/skills/<skill-name>/SKILL.md`.

```
.agents/skills/
├── spectr-proposal/
│   └── SKILL.md
├── spectr-apply/
│   └── SKILL.md
├── spectr-accept-wo-spectr-bin/
│   ├── SKILL.md
│   └── scripts/accept.sh
└── spectr-validate-wo-spectr-bin/
    ├── SKILL.md
    └── scripts/validate.sh
```

**Rationale:**
- Matches Amp's expected structure
- Each skill in its own directory
- Scripts co-located with SKILL.md for easy bundling
- Supports resource files (scripts, templates) within skill directory

### 4. Template Strategy

**Decision:** Create new skill templates in `internal/domain/templates/`:
- `skill-proposal.md.tmpl`
- `skill-apply.md.tmpl`

**Rationale:**
- Keeps skill templates separate from slash command templates
- Allows different frontmatter and content structure
- Reuses the same instructional content as slash commands where applicable
- Template variables: `{{.BaseDir}}`, `{{.SpecsDir}}`, `{{.ChangesDir}}`

### 5. Provider Priority

**Decision:** Register Amp with priority 15 (Claude Code: 10, Gemini: 20).

**Rationale:**
- Amp is based on Claude Code but is a distinct tool
- Priority 15 places it logically between Claude Code and Gemini
- Ensures Amp initializers run after Claude Code when both are selected

### 6. Reuse AgentSkillsInitializer

**Decision:** Use existing `AgentSkillsInitializer` for embedded skills.

**Rationale:**
- Already handles recursive copying of skill directories
- Preserves file permissions (executable scripts)
- Idempotent execution
- No need for new initializer types

## Implementation Sequence

1. Create skill templates in `internal/domain/templates/`
2. Create `AmpProvider` in `internal/initialize/providers/amp.go`
3. Update `TemplateManager` to load skill templates
4. Register Amp in `RegisterAllProviders()`
5. Create Amp support spec in `spectr/specs/support-amp/spec.md`

## Testing Strategy

- Unit tests for `AmpProvider.Initializers()` output
- Integration tests verifying `.agents/skills/` directory creation
- Validate SKILL.md frontmatter parsing
- Test skill scripts with embedded `spectr-accept-wo-spectr-bin`
- Verify template variable substitution in skill content

## Alternatives Considered

### Alternative 1: Generate Custom Commands

**Rejected:** Amp removed custom commands on January 29, 2026. Would be immediately deprecated.

### Alternative 2: Use Claude Code Skills Format

**Rejected:** Amp uses simpler frontmatter without Claude-specific fields like `allowed-tools` or `context`.

### Alternative 3: Global Skills Only (~/.config/agents/skills/)

**Rejected:** Project-local skills (`.agents/skills/`) are the primary convention in Amp. Global skills are secondary.

## Migration Path

For users migrating from Claude Code to Amp:
1. Run `spectr init` and select Amp
2. Amp will automatically discover skills in `.agents/skills/`
3. Existing `.claude/skills/` remain for Claude Code compatibility
4. No manual migration required - both can coexist

## Future Considerations

- If Amp adds skill dependencies or advanced frontmatter, update templates
- Monitor Amp's skill API for changes in discovery or invocation
- Consider skill versioning if Amp introduces breaking changes
