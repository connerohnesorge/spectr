# Change: Add AgentSkills Initializer

## Why

Currently, Spectr's `spectr accept` command requires the spectr binary to be
installed. When AI agents work in environments where the binary is unavailable
(e.g., sandboxed environments, CI pipelines, or fresh checkouts), they cannot
complete the accept workflow.

This change introduces an AgentSkills Initializer that:

1. Creates a generic mechanism to embed and scaffold entire skill directories
2. Provides a `spectr-accept-wo-spectr-bin` Claude skill with a bash script
   that converts `tasks.md` to `tasks.jsonc` using standard Unix tools
3. Enables agents to complete the accept workflow without the spectr binary

## What Changes

- **ADDED**: `AgentSkillsInitializer` - a generic initializer type that copies
  embedded skill directories to a target path
- **ADDED**: Embedded skill template at
  `internal/initialize/templates/skills/spectr-accept-wo-spectr-bin/`
  containing:
  - `SKILL.md` - AgentSkills-compliant skill definition
  - `scripts/accept.sh` - Bash script to convert tasks.md to tasks.jsonc
- **ADDED**: `TemplateManager.Skill(name)` method to retrieve embedded skill
  directories
- **MODIFIED**: `ClaudeProvider.Initializers()` to include the
  `AgentSkillsInitializer` for `.claude/skills/`

## Impact

- **Affected specs**: `provider-system`, `support-claude-code`
- **Affected code**:
  - `internal/initialize/providers/agentskills.go` - new initializer
  - `internal/initialize/templates.go` - skill accessor method
  - `internal/initialize/templates/skills/` - embedded skill directory
  - `internal/initialize/providers/claude.go` - add initializer instance
- **Breaking changes**: None - existing providers continue to work unchanged
