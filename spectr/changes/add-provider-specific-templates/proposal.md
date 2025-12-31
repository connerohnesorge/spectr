# Proposal: Provider-Specific Template System

## Why

Different AI coding tools have distinct capabilities, workflows, and instruction styles that the current single-template approach cannot leverage:

- **Claude Code** has an orchestrator pattern with coder/tester/stuck subagents and 200k context window
- **Codex** uses global prompts directory (`~/.codex/prompts/`) and AgentSkills in home directory
- **Gemini** uses TOML format instead of Markdown
- **Other tools** have varying tool names ("View" vs "cat"), agent loops, and delegation patterns

Currently, all 16 providers share the same 4 generic templates (slash-proposal.md.tmpl, slash-apply.md.tmpl, and TOML variants). This produces instructions that:
- Don't mention provider-specific features (like Claude's subagent delegation)
- Can't reference provider-specific file locations (like Codex's `~/.codex/prompts/`)
- Miss opportunities to optimize for each tool's strengths (like Claude's 200k context)

**User Pain Point:** When AI assistants read these generic slash commands, they don't receive guidance on how to use the specific tool's features effectively. Claude Code users don't learn about the orchestrator pattern, Codex users don't know about global prompts, etc.

## What

Introduce a **provider-specific template system with partial overrides** that allows each AI tool to customize specific sections of slash command templates while inheriting shared content from generic templates.

### Core Mechanism

Use Go's native `{{define}}` blocks to enable section-based composition:

1. **Generic templates** define sections (guardrails, steps, reference) that all providers share
2. **Provider templates** override only the sections they want to customize (e.g., Claude overrides "guardrails")
3. **Composition logic** merges provider overrides into generic templates at initialization
4. **Automatic fallback** uses generic sections when provider doesn't provide overrides

### Scope

- **Included:** Only slash command templates (slash-proposal.md.tmpl, slash-apply.md.tmpl)
- **Initial providers:** claude-code, codex (as proof-of-concept)
- **Future expansion:** gemini (TOML format), remaining 14 providers as needed
- **Not included:** Initialization templates (project.md, AGENTS.md, instruction-pointer.md) remain generic

### Example Customization

**Generic guardrails** (all providers):
```markdown
- Favor straightforward, minimal implementations first
- Keep changes tightly scoped to the requested outcome
- Refer to `spectr/AGENTS.md` and `spectr/project.md` for conventions
```

**Claude Code override** (adds provider-specific guidance):
```markdown
- **Use the orchestrator pattern**: Delegate to coder/tester/stuck subagents
- **Leverage extended context**: Use Claude's 200k context window
- **Reference skills**: Check `.claude/skills/` for available AgentSkills
- Favor straightforward, minimal implementations first
- Keep changes tightly scoped to the requested outcome
- Refer to `spectr/AGENTS.md` and `spectr/project.md` for conventions
```

## Impact

### Benefits

1. **Better Tool Utilization**
   - Claude Code users see orchestrator pattern guidance
   - Codex users learn about `~/.codex/prompts/` and global AgentSkills
   - Each tool's instructions leverage its unique capabilities

2. **Backward Compatibility**
   - All existing providers continue working without changes
   - Providers can opt-in to custom templates at any time
   - Gradual migration path (no big-bang deployment)

3. **Maintainability**
   - Shared content (90%+ of template) lives in generic templates
   - Provider overrides are small (typically 1-2 sections)
   - No full template duplication (avoids content drift)

4. **Scalability**
   - Easy to add new providers with custom templates
   - Template composition uses Go's native `text/template` (no external dependencies)
   - Performance overhead is negligible (composition happens once at init)

### Risks & Mitigations

**Risk 1: Content Drift**
- **Concern:** Provider templates could diverge significantly from generic templates
- **Mitigation:** Establish guidelines for what to customize (tool-specific) vs preserve (shared best practices)
- **Mitigation:** Code review process for all template changes

**Risk 2: Complexity**
- **Concern:** Template composition adds complexity to TemplateManager
- **Mitigation:** Use Go's native `{{define}}` blocks (well-documented, battle-tested)
- **Mitigation:** Comprehensive unit and integration tests
- **Mitigation:** Clear documentation in design.md

**Risk 3: Breaking Changes**
- **Concern:** Refactoring templates could break existing providers
- **Mitigation:** Phased migration with backward compatibility at every step
- **Mitigation:** Existing `SlashCommand()` method unchanged (calls new method internally)
- **Mitigation:** End-to-end validation before each phase

### User-Facing Changes

**For end users (running `spectr init`):**
- No changes to command interface
- `spectr init --provider=claude-code` produces Claude-specific slash commands (if implemented)
- `spectr init --provider=aider` produces generic slash commands (no customization yet)
- Behavior is transparent (users just see better instructions)

**For provider authors (future):**
- Can create provider-specific templates in `internal/initialize/templates/providers/{id}/`
- Simple opt-in: update `Initializers()` to call `tm.ProviderSlashCommand("{id}", cmd)`
- Clear documentation and examples (claude-code, codex as reference)

### Technical Impact

**Files Modified:**
- `internal/domain/template.go` - Add `ProviderTemplate` field, composition logic
- `internal/domain/templates/slash-proposal.md.tmpl` - Refactor with `{{define}}` blocks
- `internal/domain/templates/slash-apply.md.tmpl` - Refactor with `{{define}}` blocks
- `internal/initialize/templates.go` - Parse provider templates, add provider-aware methods
- `internal/initialize/providers/provider.go` - Extend `TemplateManager` interface
- `internal/initialize/providers/claude.go` - Use `ProviderSlashCommand("claude-code", cmd)`
- `internal/initialize/providers/codex.go` - Use `ProviderSlashCommand("codex", cmd)`

**Files Created:**
- `internal/initialize/templates/providers/claude-code/slash-proposal.md.tmpl`
- `internal/initialize/templates/providers/claude-code/slash-apply.md.tmpl`
- `internal/initialize/templates/providers/codex/slash-proposal.md.tmpl`
- `internal/initialize/templates/providers/codex/slash-apply.md.tmpl`
- `internal/initialize/templates/providers/gemini/` (placeholder directory for future)

**Testing Impact:**
- Add unit tests for template composition (`internal/domain/template_test.go`)
- Add integration tests for provider template parsing (`internal/initialize/templates_test.go`)
- Add end-to-end validation (run `spectr init` for each provider)

### Capabilities Affected

- **provider-system** - Core changes to template resolution and composition
- **cli-interface** - New `ProviderSlashCommand` and `ProviderTOMLSlashCommand` methods in TemplateManager interface

### Success Metrics

- ✅ `spectr init --provider=claude-code` produces slash commands mentioning "orchestrator pattern" and "coder/tester subagents"
- ✅ `spectr init --provider=codex` produces slash commands mentioning "`~/.codex/prompts/`" and "global AgentSkills"
- ✅ `spectr init --provider=aider` produces generic slash commands (no customization)
- ✅ All existing providers continue working without changes
- ✅ `spectr validate add-provider-specific-templates` passes with zero errors
- ✅ Unit tests achieve >90% coverage for new composition logic
- ✅ Integration tests verify all 16 providers can initialize successfully

## Alternatives Considered

### Alternative 1: Full Template Duplication

Create complete copies of templates for each provider.

**Pros:**
- Simple to understand (one file = one provider)
- Maximum flexibility (can change anything)

**Cons:**
- Massive duplication (16 providers × 4 templates = 64 template files)
- Content drift inevitable (generic improvements won't propagate)
- Maintenance nightmare (update one thing, update 64 files)

**Verdict:** Rejected due to maintenance burden

### Alternative 2: String-Based Section Markers

Use HTML comments to mark sections: `<!-- SECTION:guardrails -->...<!-- /SECTION:guardrails -->`

**Pros:**
- Works with any template engine
- Easy to visually identify sections

**Cons:**
- No compile-time validation
- Manual string manipulation (error-prone)
- Doesn't integrate with Go's template engine
- Custom parsing logic required

**Verdict:** Rejected due to lack of type safety and integration

### Alternative 3: Separate Files Per Section

Store each section in its own file: `guardrails.md.tmpl`, `steps.md.tmpl`, `reference.md.tmpl`

**Pros:**
- Very clear what each section contains
- Easy to override (replace one file)

**Cons:**
- File proliferation (16 providers × 3 sections × 4 templates = 192 files)
- Custom assembly logic required
- Hard to see full template at a glance
- Breaks existing template structure

**Verdict:** Rejected due to complexity and file proliferation

### Alternative 4: Configuration-Based Overrides

Store overrides in YAML/TOML config: `overrides.yaml` with section content

**Pros:**
- Centralized configuration
- Could support runtime overrides

**Cons:**
- Mixes template content (Markdown) with config (YAML/TOML)
- Harder to edit (no syntax highlighting for Markdown in YAML strings)
- Additional dependency (YAML/TOML parser)
- Doesn't leverage Go's template engine

**Verdict:** Rejected due to poor developer experience

### Selected Approach: Go Template `{{define}}` Blocks

**Why this wins:**
- Native Go feature (no dependencies)
- Compile-time validation (template parse errors)
- Integrates with existing `text/template` usage
- Clean composition via `Clone()` + `AddParseTree()`
- Minimal file changes (only wraps existing content in `{{define}}`)

## Next Steps

1. ✅ Create this proposal
2. Create delta specs for `provider-system` and `cli-interface` capabilities
3. Create `tasks.md` with implementation checklist
4. Run `spectr validate add-provider-specific-templates` and fix all issues
5. Get proposal approved
6. Implement per tasks.md (in separate work session)
