# Add Slash Command Frontmatter

## Summary

Add YAML frontmatter to spectr's slash command templates (proposal.md, apply.md) to enable agentic discovery by Claude Code and OpenCode. This allows AI tools to automatically discover and invoke spectr commands based on context.

## Motivation

Claude Code's `SlashCommand` tool and OpenCode's command system can automatically discover and invoke slash commands that have proper frontmatter metadata. Currently, spectr's slash command templates lack frontmatter, meaning:

- Claude Code cannot automatically invoke `/spectr:proposal` when a user mentions "create a proposal"
- OpenCode cannot route commands to specific agents or models
- Commands lack discoverability metadata (description, allowed-tools)

Adding frontmatter enables agentic workflows where the AI tool proactively suggests or invokes spectr commands.

## Scope

### In Scope

- Add YAML frontmatter to `slash-proposal.md.tmpl` and `slash-apply.md.tmpl`
- Use superset approach: include fields for both Claude Code and OpenCode
- Update TOML templates to ensure parity where applicable
- No changes to spectr CLI parsing (frontmatter is consumed by AI tools only)

### Out of Scope

- spectr CLI parsing frontmatter (AI tools handle this)
- New slash commands (only updating existing templates)
- Provider-specific template variants (single template per command)

## Frontmatter Fields

### Claude Code Fields
- `description`: Brief description shown in `/help` and used by `SlashCommand` tool
- `allowed-tools`: Tools the command can use (e.g., `Bash(spectr:*), Read, Glob, Grep`)
- `model`: Optional model override

### OpenCode Fields
- `description`: Brief description shown in command list
- `agent`: Optional agent to route to (e.g., `plan`)
- `model`: Optional model override
- `subtask`: Boolean to force subagent invocation

### Superset Approach

Each template includes all fields. Each tool ignores unknown fields:

```yaml
---
description: Create a Spectr change proposal
allowed-tools: Read, Glob, Grep, Write, Edit, Bash(spectr:*)
agent: plan
model: null
subtask: false
---
```

## Impact Assessment

### Positive Impact
- Claude Code users get automatic command discovery
- OpenCode users can route commands to appropriate agents
- No breaking changes (frontmatter is additive)
- Better developer experience with agentic workflows

### Risk Assessment
- Low risk: Frontmatter is parsed by AI tools, not spectr CLI
- No behavior change for existing users
- Template rendering unchanged (Go templates preserve frontmatter)

## Success Criteria

1. `spectr init` generates slash command files with proper frontmatter
2. Claude Code's `SlashCommand` tool can discover `/spectr:proposal` and `/spectr:apply`
3. OpenCode command system shows descriptions and respects agent routing
4. `spectr validate add-slash-command-frontmatter` passes
