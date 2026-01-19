# Add Codex AgentSkills Support

## Summary

Add agent skills support to the Codex provider, enabling installation of embedded
Spectr skills (`spectr-accept-wo-spectr-bin` and `spectr-validate-wo-spectr-bin`)
in the repository's `.codex/skills/` directory.

## Motivation

**Problem:** Codex CLI users cannot access Spectr's embedded agent skills. The
Codex provider currently only supports slash commands in the global
`~/.codex/prompts/` directory and the `AGENTS.md` instruction file, but lacks the
ability to install agent skills for task automation.

**Solution:** Extend the Codex provider to install agent skills in the
repository's `.codex/skills/` directory, following the same pattern as the Claude
Code provider. This enables Codex users to leverage skills like
`spectr-accept-wo-spectr-bin` (for converting tasks.md to tasks.jsonc) and
`spectr-validate-wo-spectr-bin` (for validating specs without the binary).

**Use Cases:**

1. **Task Acceptance:** Codex users can run the accept skill to convert markdown
   tasks to JSONC
2. **Validation:** Codex users can validate specs and changes using the validate
   skill
3. **Consistency:** All AI coding assistants get access to the same Spectr
   automation capabilities
4. **Repository-Focused:** Skills are installed per-repository, not globally

## Scope

### In Scope

- Add `.codex/skills/` directory initializer to Codex provider
- Add `AgentSkillsInitializer` for `spectr-accept-wo-spectr-bin` skill
- Add `AgentSkillsInitializer` for `spectr-validate-wo-spectr-bin` skill
- Update support-codex spec with new requirements for skills directory and skill
  installation
- Skills installed in repository `.codex/skills/` (project filesystem, NOT home)

### Out of Scope

- Global skills directory support (explicitly excluded per user request)
- Modifying existing Codex slash commands or AGENTS.md functionality
- Changes to skill content or templates (reuse existing embedded skills)
- Discovery or listing functionality (already covered by add-agent-skills-discovery)
- Custom skill creation or management

### Constraints

- **Repository-focused only**: Skills go in `.codex/skills/` in repo root, not
  home directory
- **Reuse existing skills**: Use the same embedded skills as Claude Code provider
- **No breaking changes**: Existing Codex functionality remains unchanged
- **Follow established patterns**: Mirror Claude Code provider's skill
  installation approach

## Impact Assessment

### User Impact

- **Low risk**: Purely additive feature, no changes to existing Codex behavior
- **High value**: Codex users gain access to task automation and validation
  without installing spectr binary
- **Better parity**: Reduces feature gap between Codex and Claude Code support

### Technical Impact

- **Modified file**: `internal/initialize/providers/codex.go` (~10 lines added)
- **Updated spec**: `spectr/specs/support-codex/spec.md` (new requirements section)
- **No new code**: Reuses existing `AgentSkillsInitializer` and skill templates
- **Test coverage**: Existing provider tests should cover the change

### Dependencies

- **No new dependencies**: Uses existing AgentSkillsInitializer from
  provider-system
- **Embedded skills**: Depends on existing skill templates in
  `internal/initialize/templates/skills/`
- **Provider framework**: Uses established dual-filesystem initializer pattern

## Alternatives Considered

### Alternative 1: Install skills in home directory `~/.codex/skills/`

**Rejected:** User explicitly requested repository-focused approach only. Global
paths add complexity and don't align with spectr's per-project philosophy.

### Alternative 2: Use different skill names for Codex

**Rejected:** Skills are provider-agnostic. Using different names would fragment
the ecosystem and confuse users.

### Alternative 3: Add Codex-specific skills

**Deferred:** Current embedded skills work for all providers. If Codex needs
custom skills in the future, we can add them then.

## Success Criteria

1. **Functional**: `spectr init` with Codex creates `.codex/skills/` with both
   skills
2. **Validated**: `spectr validate add-codex-agentskills-support` passes
3. **Tested**: Existing provider tests pass, skills are copied correctly
4. **Documented**: Spec deltas have WHEN/THEN scenarios for all requirements
5. **Simple**: Implementation adds <20 lines of code

## Related Work

- Mirrors Claude Code provider implementation
  (`internal/initialize/providers/claude.go`)
- Uses existing `AgentSkillsInitializer` (provider-system spec, line 714)
- Builds on embedded skill templates (`internal/initialize/templates/skills/`)
- Follows established Codex patterns (home directory for prompts, project for
  config)

## References

- Claude Code provider: `internal/initialize/providers/claude.go` (lines 29-38)
- AgentSkillsInitializer: `internal/initialize/providers/agentskills.go`
- Provider system spec: `spectr/specs/provider-system/spec.md` (lines 714-821)
- Support Codex spec: `spectr/specs/support-codex/spec.md`
- Agent Skills standard: <https://agentskills.io>
