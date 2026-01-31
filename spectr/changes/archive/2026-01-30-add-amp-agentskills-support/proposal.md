# Change: Add Amp Agent Skills Support

## Why

Amp (ampcode.com) is a production-grade AI coding assistant that uses the agent skills system for extensibility. While Amp is based on Claude Code's architecture, it has adopted `.agents/skills/` as the primary location for agent skills with `SKILL.md` frontmatter, deprecating custom commands in favor of user-invocable skills. Supporting Amp enables Spectr users to leverage Amp's skill discovery system and user-invocable skills (e.g., `/spectr:proposal`, `/spectr:apply`) for spec-driven development.

Amp represents a simplified, production-focused approach to agent skills:

- Uses `.agents/skills/` as the primary location (with compatibility for `.claude/skills/`)
- Requires YAML frontmatter in `SKILL.md` with `name` and `description` fields
- Supports user-invocable skills that agents can load on-demand
- Removed custom commands in favor of skills (as of January 29, 2026)

This change adds Amp as a first-class provider in Spectr.

## What Changes

- Create `AmpProvider` in `internal/initialize/providers/amp.go`
- Generate `.agents/skills/spectr-proposal/SKILL.md` and `.agents/skills/spectr-apply/SKILL.md` with Amp-compatible frontmatter
- Use `AgentSkillsInitializer` to copy embedded skill templates for `spectr-accept-wo-spectr-bin` and `spectr-validate-wo-spectr-bin`
- Register Amp provider with priority 15 (after Claude Code at 10, before Gemini at 20)
- Add Amp support spec documenting the skills-based architecture
- Create skill templates with proper frontmatter (`name`, `description`)

**Skills to be generated:**

- `.agents/skills/spectr-proposal/SKILL.md` - Create change proposals
- `.agents/skills/spectr-apply/SKILL.md` - Apply/accept proposals
- `.agents/skills/spectr-accept-wo-spectr-bin/SKILL.md` - Accept without binary
- `.agents/skills/spectr-validate-wo-spectr-bin/SKILL.md` - Validate without binary

## Impact

- Affected specs: `provider-system`, `agent-instructions`
- Affected code:
  - `internal/initialize/providers/amp.go` (new)
  - `internal/initialize/providers/registry.go` (add Amp registration)
  - `internal/domain/templates/skill-*.md.tmpl` (new skill templates)
  - `internal/initialize/templates.go` (load skill templates)
- Enables Amp users to invoke `/spectr:proposal` and `/spectr:apply` as user-invocable skills
- Provides embedded skills for validation and acceptance without the binary
- No breaking changes to existing providers
