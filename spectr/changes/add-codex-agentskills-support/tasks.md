# Implementation Tasks

## Tasks

- [ ] Read current Codex provider implementation (`internal/initialize/providers/codex.go`)
- [ ] Add `DirectoryInitializer` for `.codex/skills/` to Codex provider's `Initializers()` method
- [ ] Add `AgentSkillsInitializer` for `spectr-accept-wo-spectr-bin` skill to Codex provider
- [ ] Add `AgentSkillsInitializer` for `spectr-validate-wo-spectr-bin` skill to Codex provider
- [ ] Verify initializer ordering (directories before skills)
- [ ] Run `nix develop -c 'lint'` to verify code quality
- [ ] Run `nix develop -c 'tests'` to verify all tests pass
- [ ] Run `spectr init` manually to test Codex provider creates `.codex/skills/` with both skills
- [ ] Verify `SKILL.md` files exist in both skill directories
- [ ] Verify `scripts/accept.sh` and `scripts/validate.sh` are executable
- [ ] Run `.codex/skills/spectr-validate-wo-spectr-bin/scripts/validate.sh --all` to verify validate skill works
- [ ] Update support-codex spec in `spectr/specs/support-codex/spec.md` by merging delta requirements
- [ ] Archive this change proposal to `spectr/changes/archive/YYYY-MM-DD-add-codex-agentskills-support/`

## Validation Checklist

- [ ] All spec delta requirements have at least one scenario with WHEN/THEN format
- [ ] No duplicate requirement IDs across all capabilities
- [ ] Tasks list contains actionable, verifiable items
- [ ] `spectr validate add-codex-agentskills-support` passes without errors
- [ ] Implementation matches spec requirements exactly
- [ ] Code follows Go conventions and passes linting
- [ ] All tests pass with adequate coverage
- [ ] Manual testing confirms skills are installed and functional
