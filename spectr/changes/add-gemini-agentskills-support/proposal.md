# Add Gemini Agent Skills Support

## Summary

Extend the existing Gemini provider to support the [Agent Skills specification][agentskills],
enabling installation of Spectr's embedded skills
(`spectr-accept-wo-spectr-bin` and `spectr-validate-wo-spectr-bin`) in the
project's `.gemini/skills/` directory, alongside the existing TOML-based slash
commands.

## Motivation

**Problem:** The Gemini CLI provider currently only supports TOML-based slash
commands in `.gemini/commands/spectr/`. While Gemini CLI has native support for
Agent Skills (an experimental feature via `experimental.skills`), Spectr doesn't
leverage this capability to provide its task automation skills.

**Solution:** Extend the Gemini provider to:

1. Create `.gemini/skills/` directory for agent skills
2. Install `spectr-accept-wo-spectr-bin` and `spectr-validate-wo-spectr-bin`
   skills
3. Create `GEMINI.md` instruction file for workspace-wide guidance
4. Maintain existing TOML slash commands for backward compatibility

**Use Cases:**

1. **Task Acceptance:** Gemini users can invoke the accept skill to convert
   tasks.md to tasks.jsonc without the spectr binary
2. **Validation:** Gemini users can validate specs and changes using the
   validate skill in sandboxed environments
3. **Workspace Guidance:** GEMINI.md provides persistent context about Spectr
   workflows
4. **Progressive Disclosure:** Skills only load full instructions when activated,
   saving context tokens

## Scope

### In Scope

- Add `.gemini/skills/` directory initializer to Gemini provider
- Add `AgentSkillsInitializer` for `spectr-accept-wo-spectr-bin` skill
- Add `AgentSkillsInitializer` for `spectr-validate-wo-spectr-bin` skill
- Create `GEMINI.md` instruction file with Spectr guidance
- Maintain existing TOML slash commands in `.gemini/commands/spectr/`
- Skills installed in project `.gemini/skills/` (workspace scope only)

### Out of Scope

- User-level skills directory support (`~/.gemini/skills/`) - per decision to
  keep project-focused
- Modifying existing TOML slash command functionality
- Changes to skill content or templates (reuse existing embedded skills)
- Adding `allowed-tools` field to skill templates
- Extension skills support
- Custom skill creation or management

### Constraints

- **Project-focused only:** Skills go in `.gemini/skills/` in repo root, not
  home directory
- **Reuse existing skills:** Use the same embedded skills as Claude/Codex
  providers
- **No breaking changes:** Existing TOML slash commands remain functional
- **Follow established patterns:** Mirror Claude/Codex provider skill
  installation approach

## Impact Assessment

### User Impact

- **Low risk:** Purely additive feature, no changes to existing Gemini behavior
- **High value:** Gemini users gain access to task automation and validation
  without installing spectr binary
- **Better parity:** Aligns Gemini support with Claude Code and Codex providers

### Technical Impact

- **Modified file:** `internal/initialize/providers/gemini.go` (~15 lines added)
- **New template:** Create GEMINI.md template for instruction pointer
- **New spec:** `spectr/specs/support-gemini/spec.md` (merged from delta)
- **No new code patterns:** Reuses existing `AgentSkillsInitializer` and
  templates
- **Test coverage:** Existing provider tests should cover the change

### Dependencies

- **No new dependencies:** Uses existing `AgentSkillsInitializer` from
  provider-system
- **Embedded skills:** Depends on existing skill templates in
  `internal/initialize/templates/skills/`
- **Provider framework:** Uses established dual-filesystem initializer pattern

## Alternatives Considered

### Alternative 1: Install skills in home directory `~/.gemini/skills/`

**Rejected:** User explicitly chose project-focused approach. Global paths add
complexity and don't align with Spectr's per-project philosophy.

### Alternative 2: Remove TOML slash commands, use skills only

**Rejected:** User chose to keep both. Existing slash commands provide
immediate invocation while skills provide on-demand expertise. Different
use cases.

### Alternative 3: Add allowed-tools field to skill templates

**Rejected:** User chose not to modify templates. Existing skill templates
are already Agent Skills compliant and work without pre-approved tools.

### Alternative 4: Skip GEMINI.md instruction file

**Rejected:** User explicitly chose to create GEMINI.md for consistency with
other providers (CLAUDE.md, AGENTS.md).

## Success Criteria

1. **Functional:** `spectr init` with Gemini creates `.gemini/skills/` with both
   skills
2. **Instruction file:** `GEMINI.md` is created with Spectr guidance
3. **Backward compatible:** Existing TOML slash commands continue to work
4. **Validated:** `spectr validate add-gemini-agentskills-support` passes
5. **Tested:** Provider tests pass, skills are copied correctly
6. **Documented:** Spec deltas have WHEN/THEN scenarios for all requirements

## Related Work

- Mirrors Claude Code provider implementation
  (`internal/initialize/providers/claude.go`)
- Mirrors Codex provider implementation
  (`internal/initialize/providers/codex.go`)
- Uses existing `AgentSkillsInitializer` (provider-system spec)
- Builds on embedded skill templates (`internal/initialize/templates/skills/`)
- [Agent Skills standard][agentskills]

## References

- Current Gemini provider: `internal/initialize/providers/gemini.go`
- Claude Code provider: `internal/initialize/providers/claude.go` (lines 29-38)
- AgentSkillsInitializer: `internal/initialize/providers/agentskills.go`
- [Gemini CLI Agent Skills docs][gemini-cli]

[agentskills]: https://agentskills.io
[gemini-cli]: https://ai.google.dev/gemini-api/docs/gemini-cli
