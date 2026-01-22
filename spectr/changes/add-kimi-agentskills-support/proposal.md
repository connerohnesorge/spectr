# Add Kimi Agent Skills Support

## Summary

This proposal adds support for Kimi CLI's agent skills and AGENTS.md directory,
following the same pattern as existing Claude Code and Codex providers. This
enables Kimi users to leverage Spectr's agent skills for spec-driven development
workflows.

## Motivation

Kimi is a popular AI coding assistant that would benefit from Spectr's agent
skills support. Currently, Spectr supports Claude Code and Codex providers with
agent skills that enable:

- Converting tasks.md to tasks.jsonc without requiring spectr binary
- Validating specs without requiring spectr binary

Adding Kimi support extends Spectr's reach to Kimi users and maintains
consistency across AI coding assistants.

## Scope

### In Scope

- Create Kimi provider implementation
- Support `.kimi/skills` directory creation
- Support `.kimi/commands` directory creation
- Create `AGENTS.md` instruction file
- Implement `spectr-proposal` and `spectr-apply` slash commands
- Install agent skills: `spectr-accept-wo-spectr-bin` and
  `spectr-validate-wo-spectr-bin`

### Out of Scope

- Home directory configuration (e.g., `~/.kimi/prompts`)
- Additional Kimi-specific skills beyond the core Spectr skills
- Integration with Kimi's proprietary features

## Impact Assessment

### Positive Impact

- Kimi users gain access to Spectr's agent skills
- Consistent experience across AI coding assistants
- No breaking changes to existing providers
- Follows established patterns from Claude/Codex implementations

### Risk Assessment

- Low risk: Follows proven implementation pattern
- Minimal maintenance overhead
- No external dependencies added

## Success Criteria

1. `spectr init` creates `.kimi/skills` and `.kimi/commands` directories for Kimi
   projects
2. `AGENTS.md` file is created with Kimi-specific instructions
3. Agent skills are properly installed and executable
4. Slash commands `spectr-proposal` and `spectr-apply` work correctly
5. Validation passes: `spectr validate add-kimi-agentskills-support`
