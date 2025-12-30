# Change: Add Provider-Specific Templates

## Why

Different AI coding tools (Claude Code, Crush, Cursor, etc.) have different
tool names, agent loops, and instruction styles. The current single-template
approach produces generic instructions that don't leverage provider-specific
capabilities or terminology (e.g., "View" vs "cat", "Edit" vs "write", agent
delegation patterns).

## What Changes

- Introduce per-provider template directories: `templates/{provider-id}/`
- Each provider can override any template: `AGENTS.md.tmpl`,
  `instruction-pointer.md.tmpl`, `slash-proposal.md.tmpl`,
  `slash-apply.md.tmpl`
- Generic templates remain as fallback for providers without custom
  templates
- `TemplateManager` resolves templates with provider-first lookup, then
  generic fallback

## Impact

- Affected specs: `cli-framework` (template resolution logic)
- Affected code:
  - `internal/initialize/templates.go` - template resolution logic
  - `internal/initialize/providers/provider.go` - pass provider ID to
    template renderer
  - `internal/initialize/templates/` - new provider subdirectories
