# Proposal: Provider-Specific Template Overrides for AI Tools

## Why

Different AI coding tools have distinct characteristics that generic templates cannot effectively address:

- **Claude Code**: Orchestrator pattern with subagent delegation (coder/tester/stuck), 200k context window
- **Codex**: Global prompts directory (`~/.codex/prompts/`), home directory-based skill management
- **Gemini**: TOML format, different API and workflow patterns
- **Other tools**: Varying tool names, agent loop patterns, and best practices

Currently, all 16 providers use identical generic templates. This approach:
- Misses tool-specific features (Claude's subagent pattern, Codex's global prompts)
- Produces instructions that don't leverage each tool's strengths
- Can't reference tool-specific file locations or capabilities
- Creates friction: AI assistants don't learn how to effectively use the specific tool

**Example Gap**: When Claude Code users initialize Spectr, they receive generic instructions with no mention of the orchestrator pattern, coder/tester subagents, or how to delegate tasks effectively.

## What

Introduce a **provider-specific template override system** that allows tools to customize sections of slash command templates while inheriting shared content from generic templates.

### Core Mechanism

1. **Generic templates** define section blocks (guardrails, steps, reference) using Go's `{{define}}` directive
2. **Provider templates** override specific sections with tool-specific guidance
3. **Composition merges** base template with provider overrides at initialization
4. **Automatic fallback** uses base sections when provider doesn't override

### Example

**Generic Guardrails** (all providers):
```markdown
- Favor straightforward, minimal implementations
- Keep changes tightly scoped to the requested outcome
- Refer to spectr/AGENTS.md and spectr/project.md for conventions
```

**Claude Code Override**:
```markdown
- **Use the orchestrator pattern**: Delegate to coder/tester/stuck subagents
- **Leverage 200k context**: Use Claude's extended window for comprehensive analysis
- Favor straightforward, minimal implementations
- Keep changes tightly scoped to the requested outcome
- Refer to spectr/AGENTS.md and spectr/project.md for conventions
```

### Scope

- **Templates**: Only slash commands (slash-proposal.md.tmpl, slash-apply.md.tmpl)
- **Initial Providers**: claude-code, codex (proof-of-concept)
- **Format**: Markdown and TOML templates supported equally
- **Not Included**: Initialization templates (project.md, AGENTS.md) remain generic

## Impact

### User-Facing Changes

**For end users** (`spectr init`):
- No command interface changes
- Tool-specific slash commands produced transparently
- Better instructions for each tool (invisible improvement)

**For provider authors** (future):
- Simple opt-in: call `tm.ProviderSlashCommand("{id}", cmd)` instead of `tm.SlashCommand(cmd)`
- Create `internal/initialize/templates/providers/{id}/` directory with override templates
- Clear patterns from claude-code and codex examples

### Technical Impact

**Files Modified** (7):
- `internal/domain/template.go` - TemplateRef composition support
- `internal/domain/templates/slash-proposal.md.tmpl` - Add {{define}} blocks
- `internal/domain/templates/slash-apply.md.tmpl` - Add {{define}} blocks
- `internal/initialize/templates.go` - Parse provider templates, add provider-aware methods
- `internal/initialize/providers/provider.go` - Extend TemplateManager interface
- `internal/initialize/providers/claude.go` - Opt-in to provider templates
- `internal/initialize/providers/codex.go` - Opt-in to provider templates

**Files Created** (5):
- `internal/initialize/templates/providers/claude-code/slash-proposal.md.tmpl`
- `internal/initialize/templates/providers/claude-code/slash-apply.md.tmpl`
- `internal/initialize/templates/providers/codex/slash-proposal.md.tmpl`
- `internal/initialize/templates/providers/codex/slash-apply.md.tmpl`
- `internal/initialize/templates/providers/gemini/.gitkeep` (placeholder)

### Benefits

1. **Tool Optimization**: Each AI tool gets instructions tailored to its capabilities
2. **Content Reuse**: 90% of template content shared (avoid duplication)
3. **Gradual Migration**: Providers opt-in, no forced changes
4. **Scalability**: Easy to add new providers with custom templates
5. **Backward Compatible**: All existing providers work unchanged

### Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| Content drift | Establish guidelines: customize tool-specific, preserve shared best practices |
| Template complexity | Use Go's native `{{define}}` (well-documented, no custom parsing) |
| Breaking changes | Phased migration with fallback at each step |
| Parse-time validation errors | Catch during initialization, fail fast with clear error messages |

## Alternatives Considered

### Alternative 1: Full Template Duplication
Create complete template copies for each provider.
- ❌ 16 providers × 4 templates = 64 files
- ❌ Content drift inevitable
- ❌ Maintenance nightmare

### Alternative 2: String-Based Section Markers
Use HTML comments: `<!-- SECTION:guardrails -->...<!-- /SECTION:guardrails -->`
- ❌ No compile-time validation
- ❌ Manual string manipulation (error-prone)
- ❌ Doesn't integrate with Go's template engine

### Alternative 3: Separate Files Per Section
Store sections individually: `guardrails.md.tmpl`, `steps.md.tmpl`
- ❌ File proliferation (192+ files)
- ❌ Hard to see full template at a glance
- ❌ Complex assembly logic

### Selected Approach: Go Template `{{define}}` Blocks ✓
- ✓ Native Go feature (no dependencies)
- ✓ Compile-time template validation
- ✓ Works with existing `text/template` usage
- ✓ Clean composition via `Clone()` + `AddParseTree()`
- ✓ Minimal changes to existing templates (just add wrappers)

## Implementation Overview

### Phase 1: Foundation
- Refactor generic templates with `{{define}}` blocks
- Add composition support to TemplateRef
- Update TemplateManager to parse provider directories

### Phase 2: Provider Templates
- Create claude-code overrides (guardrails section)
- Create codex overrides (steps and reference sections)
- Test composition produces expected output

### Phase 3: Provider Migration
- Update ClaudeProvider to use `ProviderSlashCommand()`
- Update CodexProvider to use `ProviderSlashCommand()`
- Verify backward compatibility

### Phase 4: Validation
- Run full test suite
- End-to-end validation with all providers
- Confirm spectr validate passes

## Success Criteria

- ✓ Generic templates refactored with sections
- ✓ Provider-specific templates override targeted sections
- ✓ Template composition works correctly
- ✓ All 16 providers continue working
- ✓ Claude Code uses custom templates
- ✓ Codex uses custom templates
- ✓ `spectr init` produces provider-specific slash commands
- ✓ Providers without custom templates fall back to generic
- ✓ Validation passes (zero errors)
- ✓ All tests pass

## Next Steps

1. ✓ Create this proposal
2. Create delta specs for provider-system
3. Create detailed design.md with code examples
4. Create tasks.md with implementation checklist
5. Run spectr validate
6. Get proposal approved
