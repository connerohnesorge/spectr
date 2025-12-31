# Add Agent Skills Discovery Support

## Summary

Extend Spectr's provider system with agent skills discovery capabilities, following the [Agent Skills standard](https://agentskills.io). This adds the ability to list embedded agent skills, parse SKILL.md frontmatter metadata, check installation status, and integrate with the CLI via `spectr list --skills`.

## Motivation

**Problem:** Spectr currently supports installing agent skills (via `AgentSkillsInitializer`) but provides no way to discover what skills are available, view their metadata, or check if they're installed. Users cannot easily see which skills exist or understand what dependencies they require.

**Solution:** Add discovery APIs and CLI integration to make agent skills transparent and discoverable. Users can run `spectr list --skills` to see all embedded skills, their descriptions, compatibility requirements, and installation status.

**Use Cases:**
1. **Discovery:** "What agent skills does Spectr provide?"
2. **Metadata:** "What dependencies does the accept skill require?"
3. **Status:** "Is the validate skill already installed in my project?"
4. **Documentation:** Skills become self-documenting via frontmatter

## Scope

### In Scope
- Parse SKILL.md frontmatter (YAML format with name, description, compatibility)
- List all embedded agent skills from `internal/initialize/templates/skills/`
- Check if a skill is installed in `.claude/skills/<skill-name>/`
- CLI command: `spectr list --skills` with text, long, and JSON output formats
- Extend provider-system spec with discovery requirements

### Out of Scope
- Installing/removing skills via CLI (future enhancement)
- Discovering skills from other sources (installed, custom, global paths)
- Interactive skill management (future enhancement)
- Skill version tracking or updates

### Constraints
- **Repository-focused only**: Only check `.claude/skills/` in repo root, no global paths
- **Embedded skills only**: Only discover skills embedded in Spectr binary
- **Read-only operations**: No state changes, purely informational
- **Extend existing spec**: Add to provider-system spec, not a new capability

## Impact Assessment

### User Impact
- **Low risk, high value**: Purely additive feature, no breaking changes
- **Improved UX**: Skills become discoverable and self-documenting
- **Better onboarding**: New users can explore available skills

### Technical Impact
- **New package**: `internal/discovery` for skill parsing/listing (~100 lines)
- **Extended packages**: `internal/list` and `cmd/list` gain skill support (~105 lines)
- **Test coverage**: ~200 lines of new tests
- **Total code**: ~455 lines (well under 500-line simplicity target)

### Dependencies
- **No new external dependencies**: Uses existing `gopkg.in/yaml.v3`, `afero`, `io/fs`
- **Existing patterns**: Follows established list/discovery architecture

## Alternatives Considered

### Alternative 1: Separate `spectr skills` top-level command
**Rejected:** Inconsistent with existing patterns. Spectr uses `list --specs`, `list --all`, so `list --skills` maintains consistency.

### Alternative 2: Discover installed skills from all sources
**Deferred:** Adds complexity (scanning multiple providers, handling conflicts). Start simple with embedded-only, extend later if needed.

### Alternative 3: Add install/remove commands now
**Deferred:** Focus on discovery first. Write operations require more design (error handling, validation, undo). Can add in future iteration.

## Success Criteria

1. **Functional**: `spectr list --skills` outputs all embedded skills
2. **Validated**: `spectr validate add-agent-skills-discovery` passes
3. **Tested**: All tests pass with >80% coverage for new code
4. **Documented**: Spec deltas have WHEN/THEN scenarios for all requirements
5. **Simple**: Implementation stays under 500 lines total

## Related Work

- Extends existing `AgentSkillsInitializer` (provider-system spec, line 714)
- Builds on `TemplateManager.SkillFS()` for embedded skill access
- Follows patterns from `internal/list` package (ListChanges, ListSpecs)
- Uses YAML frontmatter like Hugo/Jekyll/OpenAI Codex skills

## References

- [Agent Skills Standard](https://agentskills.io)
- [OpenAI Codex Skills Documentation](https://platform.openai.com/docs/guides/codex/skills)
- Existing spec: `spectr/specs/provider-system/spec.md` (lines 714-821)
- Existing implementation: `internal/initialize/providers/agentskills.go`
