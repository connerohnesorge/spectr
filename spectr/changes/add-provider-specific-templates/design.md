# Design: Provider-Specific Template System with Partial Overrides

## Problem Statement

Different AI coding tools (Claude Code, Crush, Cursor, Codex, Gemini, etc.) have different:
- Tool names and capabilities (e.g., "View" vs "cat", "Edit" vs "write")
- Agent loop patterns (orchestrator pattern, direct implementation, etc.)
- Instruction styles (200k context management, skill delegation, etc.)
- File formats (Markdown vs TOML)

The current single-template approach produces generic instructions that don't leverage provider-specific capabilities or terminology. All 16 providers share the same 4 core templates (slash-proposal.md.tmpl, slash-apply.md.tmpl, and TOML variants), missing opportunities for targeted optimization.

## Design Goals

1. **Partial Overrides**: Providers can override specific sections (e.g., "Guardrails") without duplicating entire templates
2. **Inheritance**: Unspecified sections automatically fall back to generic template
3. **Backward Compatibility**: Existing providers continue to work without changes
4. **Type Safety**: Compile-time validation of section names via Go templates
5. **Maintainability**: Clear composition logic, easy to debug
6. **Performance**: Template parsing happens once at initialization (no runtime overhead)

## Architectural Approach

### Template Composition with Go's `{{define}}` Blocks

Use Go's native `text/template` package features for section-based composition:

**Base Template Structure:**
```markdown
{{define "guardrails"}}
## Guardrails
- Generic best practices...
{{end}}

{{define "steps"}}
## Steps
1. Generic workflow steps...
{{end}}

{{define "reference"}}
## Reference
- Generic documentation pointers...
{{end}}

{{define "main"}}
# Proposal Creation Guide
{{template "guardrails" .}}
{{template "steps" .}}
{{template "reference" .}}
{{end}}
```

**Provider Override (only override what changes):**
```markdown
{{define "guardrails"}}
## Guardrails (Claude Code Edition)
- Use the orchestrator pattern: delegate to coder/tester/stuck subagents
- Leverage Claude's 200k context window
- [... rest of generic guardrails ...]
{{end}}
```

### Why This Approach?

**Advantages:**
- Native Go feature (no external dependencies)
- Compile-time template validation
- Clean composition via `template.Clone()` + `AddParseTree()`
- Last-wins semantics (provider overrides automatically replace base sections)
- Works with Go's template execution model

**Rejected Alternatives:**
1. **String-based section markers** (`<!-- SECTION:name -->`) - No compile-time validation, manual string manipulation
2. **Separate files per section** - File proliferation (12+ files per provider), complex assembly
3. **Full template duplication** - Maintenance burden, content drift between providers

## Directory Structure

```
internal/
  domain/
    templates/
      slash-proposal.md.tmpl          # Generic base (refactored with {{define}} blocks)
      slash-apply.md.tmpl              # Generic base (refactored with {{define}} blocks)
      slash-proposal.toml.tmpl         # Generic TOML (unchanged for now)
      slash-apply.toml.tmpl            # Generic TOML (unchanged for now)

  initialize/
    templates/
      providers/                       # NEW: Provider-specific overrides
        claude-code/
          slash-proposal.md.tmpl       # Override "guardrails" section
          slash-apply.md.tmpl           # Override "guardrails" section
        gemini/
          slash-proposal.toml.tmpl     # Override TOML sections (future)
          slash-apply.toml.tmpl         # Override TOML sections (future)
        codex/
          slash-proposal.md.tmpl       # Override "steps" section
          slash-apply.md.tmpl           # Override "steps" section
```

**Rationale:**
- Generic templates stay in `internal/domain/templates/` (no breaking changes to embed directive)
- Provider overrides in `internal/initialize/templates/providers/` (mirrors skill templates pattern)
- Flat provider directories (no nested format subdirectories - format encoded in extension)
- Consistent naming across providers (no provider prefixes, directory provides namespace)

## Code Changes

### 1. Update `domain.TemplateRef` (internal/domain/template.go)

Add optional `ProviderTemplate` field for composition:

```go
// TemplateRef is a type-safe reference to a parsed template.
// Supports provider-specific overrides via template composition.
type TemplateRef struct {
    Name             string             // template file name (e.g., "slash-proposal.md.tmpl")
    Template         *template.Template // base template with all sections
    ProviderTemplate *template.Template // provider overrides (nil if none)
}

// Render executes the template with the given context.
// If ProviderTemplate is present, composes base + provider before rendering.
func (tr TemplateRef) Render(ctx *TemplateContext) (string, error) {
    // No provider overrides - use base template directly
    if tr.ProviderTemplate == nil {
        var buf bytes.Buffer
        if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
            return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
        }
        return buf.String(), nil
    }

    // Compose base + provider overrides
    composed, err := tr.composeTemplate()
    if err != nil {
        return "", fmt.Errorf("failed to compose template %s: %w", tr.Name, err)
    }

    // Execute "main" template (the composition entry point)
    var buf bytes.Buffer
    if err := composed.ExecuteTemplate(&buf, "main", ctx); err != nil {
        return "", fmt.Errorf("failed to render composed template: %w", err)
    }
    return buf.String(), nil
}

// composeTemplate merges base template with provider overrides.
// Provider sections override base sections (last-wins semantics).
func (tr TemplateRef) composeTemplate() (*template.Template, error) {
    // Clone base template (preserves all base sections)
    composed := template.Must(tr.Template.Clone())

    // Merge provider overrides - last-wins for duplicate {{define}} names
    for _, t := range tr.ProviderTemplate.Templates() {
        if _, err := composed.AddParseTree(t.Name(), t.Tree); err != nil {
            return nil, fmt.Errorf("failed to merge provider template %s: %w", t.Name(), err)
        }
    }

    return composed, nil
}
```

**Key Design Decisions:**
- `ProviderTemplate` is optional (nil = no composition needed)
- Composition happens lazily in `Render()` (not during TemplateManager construction)
- Uses `template.Clone()` to avoid mutating base template
- `AddParseTree()` provides last-wins semantics automatically
- Execute "main" template (the composition root defined in base)

### 2. Update `TemplateManager` (internal/initialize/templates.go)

Add provider template parsing and provider-aware methods:

```go
type TemplateManager struct {
    templates         *template.Template                 // Existing: base + domain merged
    providerTemplates map[string]*template.Template      // NEW: provider-id -> overrides
}

func NewTemplateManager() (*TemplateManager, error) {
    // Parse main templates (existing logic)
    mainTmpl, err := template.ParseFS(templateFS, "templates/**/*.tmpl")
    if err != nil {
        return nil, fmt.Errorf("failed to parse main templates: %w", err)
    }

    // Parse and merge domain templates (existing logic)
    domainTmpl, err := template.ParseFS(domain.TemplateFS, "templates/*.tmpl")
    if err != nil {
        return nil, fmt.Errorf("failed to parse domain templates: %w", err)
    }

    for _, t := range domainTmpl.Templates() {
        if _, err := mainTmpl.AddParseTree(t.Name(), t.Tree); err != nil {
            return nil, fmt.Errorf("failed to merge template %s: %w", t.Name(), err)
        }
    }

    // NEW: Parse provider-specific templates
    providerTmpls := make(map[string]*template.Template)

    // Check if providers directory exists
    entries, err := fs.ReadDir(templateFS, "templates/providers")
    if err != nil {
        // No providers directory yet - that's OK (backward compatibility)
        return &TemplateManager{
            templates:         mainTmpl,
            providerTemplates: providerTmpls,
        }, nil
    }

    // Discover and parse each provider directory
    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }

        providerID := entry.Name()
        pattern := fmt.Sprintf("templates/providers/%s/*.tmpl", providerID)

        providerTmpl, err := template.ParseFS(templateFS, pattern)
        if err != nil {
            // Provider directory exists but has no templates - skip
            continue
        }

        if providerTmpl != nil {
            providerTmpls[providerID] = providerTmpl
        }
    }

    return &TemplateManager{
        templates:         mainTmpl,
        providerTemplates: providerTmpls,
    }, nil
}

// ProviderSlashCommand returns a Markdown template reference for a specific provider.
// Falls back to generic template if provider has no custom templates.
func (tm *TemplateManager) ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef {
    names := map[domain.SlashCommand]string{
        domain.SlashProposal: "slash-proposal.md.tmpl",
        domain.SlashApply:    "slash-apply.md.tmpl",
    }

    return domain.TemplateRef{
        Name:             names[cmd],
        Template:         tm.templates,                    // Base template
        ProviderTemplate: tm.providerTemplates[providerID], // nil if provider has no overrides
    }
}

// ProviderTOMLSlashCommand returns a TOML template reference for a specific provider.
// Falls back to generic template if provider has no custom templates.
func (tm *TemplateManager) ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef {
    names := map[domain.SlashCommand]string{
        domain.SlashProposal: "slash-proposal.toml.tmpl",
        domain.SlashApply:    "slash-apply.toml.tmpl",
    }

    return domain.TemplateRef{
        Name:             names[cmd],
        Template:         tm.templates,
        ProviderTemplate: tm.providerTemplates[providerID],
    }
}

// SlashCommand (existing - backward compatible)
// Returns generic template with no provider overrides
func (tm *TemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
    return tm.ProviderSlashCommand("", cmd) // Empty provider ID = no overrides
}

// TOMLSlashCommand (existing - backward compatible)
// Returns generic TOML template with no provider overrides
func (tm *TemplateManager) TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
    return tm.ProviderTOMLSlashCommand("", cmd)
}
```

**Key Design Decisions:**
- `providerTemplates` map stores all parsed provider templates (parsed once at init)
- Graceful degradation: missing providers directory doesn't break initialization
- Provider ID lookup returns nil if provider has no custom templates (automatic fallback)
- Existing `SlashCommand()` and `TOMLSlashCommand()` call new methods with empty provider ID
- No changes to existing callers required (backward compatibility)

### 3. Update `TemplateManager` Interface (internal/initialize/providers/provider.go)

Extend interface with provider-aware methods:

```go
type TemplateManager interface {
    // Existing methods (unchanged)
    InstructionPointer() domain.TemplateRef
    Agents() domain.TemplateRef
    SlashCommand(cmd domain.SlashCommand) domain.TemplateRef
    TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef
    SkillFS(skillName string) (fs.FS, error)

    // NEW: Provider-aware template methods
    ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef
    ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef
}
```

**Key Design Decisions:**
- Interface extension (not breaking change - adds methods, doesn't remove)
- Old methods remain for backward compatibility
- New methods follow same naming pattern with "Provider" prefix
- ProviderID is string (matches Registration.ID)

### 4. Update Provider Implementations

Providers opt-in by using new methods:

**Claude Code (internal/initialize/providers/claude.go):**
```go
func (*ClaudeProvider) Initializers(_ context.Context, tm TemplateManager) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".claude/commands/spectr"),
        NewDirectoryInitializer(".claude/skills"),
        NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer()),
        NewSlashCommandsInitializer(".claude/commands/spectr",
            map[domain.SlashCommand]domain.TemplateRef{
                domain.SlashProposal: tm.ProviderSlashCommand("claude-code", domain.SlashProposal),
                domain.SlashApply:    tm.ProviderSlashCommand("claude-code", domain.SlashApply),
            },
        ),
        // ... AgentSkills initializers unchanged ...
    }
}
```

**Codex (internal/initialize/providers/codex.go):**
```go
func (*CodexProvider) Initializers(_ context.Context, tm TemplateManager) []Initializer {
    return []Initializer{
        // ... existing directory/config initializers ...
        NewHomeSlashCommandsInitializer("spectr",
            map[domain.SlashCommand]domain.TemplateRef{
                domain.SlashProposal: tm.ProviderSlashCommand("codex", domain.SlashProposal),
                domain.SlashApply:    tm.ProviderSlashCommand("codex", domain.SlashApply),
            },
        ),
        // ... AgentSkills initializers unchanged ...
    }
}
```

**Gemini (internal/initialize/providers/gemini.go):**
```go
func (*GeminiProvider) Initializers(_ context.Context, tm TemplateManager) []Initializer {
    return []Initializer{
        // ... existing initializers ...
        NewTOMLSlashCommandsInitializer(".gemini/commands/spectr",
            map[domain.SlashCommand]domain.TemplateRef{
                domain.SlashProposal: tm.ProviderTOMLSlashCommand("gemini", domain.SlashProposal),
                domain.SlashApply:    tm.ProviderTOMLSlashCommand("gemini", domain.SlashApply),
            },
        ),
    }
}
```

**Other Providers:**
- No changes required (continue using `tm.SlashCommand()` or `tm.TOMLSlashCommand()`)
- Automatically use generic templates
- Can migrate to provider-specific templates at any time

## Template Section Definitions

### Generic Markdown Templates

Refactor `internal/domain/templates/slash-proposal.md.tmpl`:

```markdown
{{define "guardrails"}}
## Guardrails

- Favor straightforward, minimal implementations first and add complexity
  only when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `{{ .AgentsFile }}` and `{{ .ProjectFile }}` (located inside the
  `{{ .BaseDir }}/` directory—run `ls {{ .BaseDir }}`) if you need additional
  Spectr conventions or clarifications.
- Identify any vague or ambiguous details and ask the necessary follow-up
  questions before editing files.

Note: You are not implementing yet, you are fully planning and creating the change proposal using spectr.
{{end}}

{{define "steps"}}
## Steps

1. Review `{{ .ProjectFile }}`, read `{{ .SpecsDir }}/` and
   `{{ .ChangesDir }}/` directories, and inspect related code or docs (e.g.,
   via `rg`/`ls`) to ground the proposal in current behaviour; note any gaps
   that require clarification.
2. Choose a unique verb-led `change-id` and scaffold `proposal.md`,
   `tasks.md`, and `design.md` (when needed) under
   `{{ .ChangesDir }}/<id>/`.
3. Map the change into concrete capabilities or requirements, breaking
   multi-scope efforts into distinct spec deltas with clear relationships and
   sequencing.
4. Capture architectural reasoning in `design.md` when the solution spans
   multiple systems, introduces new patterns, or demands trade-off discussion
   before committing to specs.
5. Draft spec deltas in `{{ .ChangesDir }}/<id>/specs/<capability>/spec.md`
   (one folder per capability) using `## ADDED|MODIFIED|REMOVED Requirements`
   with at least one `#### Scenario:` per requirement and cross-reference
   related capabilities when relevant.
6. Draft `tasks.md` as an ordered list of small, verifiable work items that
   deliver user-visible progress, include validation (tests, tooling), and
   highlight dependencies or parallelizable work.
7. Validate with `spectr validate <id>` and resolve every issue before
   sharing the proposal.
{{end}}

{{define "reference"}}
## Reference

- Read delta specs directly at
  `{{ .ChangesDir }}/<id>/specs/<capability>/spec.md` when validation fails.
- Read existing specs at `{{ .SpecsDir }}/<capability>/spec.md` to understand
  current state.
- Search existing requirements with `rg -n "Requirement:|Scenario:"
  {{ .SpecsDir }}` before writing new ones.
- Explore the codebase with `rg <keyword>`, `ls`, or direct file reads so
  proposals align with current implementation realities.
{{end}}

{{define "main"}}
# Proposal Creation Guide
{{template "guardrails" .}}
{{template "steps" .}}
{{template "reference" .}}
{{end}}
```

Similar refactoring for `slash-apply.md.tmpl`:

```markdown
{{define "guardrails"}}
## Guardrails

- Favor straightforward, minimal implementations first and add complexity only
  when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `{{ .AgentsFile }}` and `{{ .BaseDir }}/project.md` (located inside
  the `{{ .BaseDir }}/` directory—run `ls {{ .BaseDir }}`) if you need
  additional Spectr conventions or clarifications.
{{end}}

{{define "steps"}}
## Steps

Track these steps as TODOs and complete them one by one.

1. Run `spectr accept <id>` to convert `tasks.md` to `tasks.jsonc` format for
   stable task tracking.
2. Read `{{ .ChangesDir }}/<id>/proposal.md`, `design.md` (if present), and
   `tasks.jsonc` to confirm scope and acceptance criteria.
3. Work through tasks sequentially, keeping edits minimal and focused on the
   requested change. Update the task status in `tasks.jsonc` after verifying
   the work.
4. Confirm completion before updating statuses—make sure every item in
   `tasks.jsonc` is finished.
5. Verify/Update all task status in `tasks.jsonc` after all work is done. Tasks
   have status values: `pending`, `in_progress`, `completed`.
6. Read `{{ .ChangesDir }}/` and `{{ .SpecsDir }}/` directories when additional
   context is required.
{{end}}

{{define "reference"}}
## Reference

- Read `{{ .ChangesDir }}/<id>/proposal.md` for proposal details.
- Read `{{ .ChangesDir }}/<id>/specs/<capability>/spec.md` for delta specs.
{{end}}

{{define "main"}}
# Guardrails
{{template "guardrails" .}}
{{template "steps" .}}
{{template "reference" .}}
{{end}}
```

### Provider-Specific Templates

**Claude Code** (`internal/initialize/templates/providers/claude-code/slash-proposal.md.tmpl`):
```markdown
{{define "guardrails"}}
## Guardrails (Claude Code Edition)

- **Use the orchestrator pattern**: You are Claude Code with 200k context. Delegate coding tasks to the `coder` subagent, testing to the `tester` subagent, and escalations to the `stuck` subagent.
- **Leverage extended context**: Take advantage of Claude's 200k context window for comprehensive proposals and full codebase understanding.
- **Reference skills**: Use AgentSkills available at `.claude/skills/` when applicable (e.g., spectr-validate-wo-spectr-bin).
- Favor straightforward, minimal implementations first and add complexity
  only when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `{{ .AgentsFile }}` and `{{ .ProjectFile }}` (located inside the
  `{{ .BaseDir }}/` directory—run `ls {{ .BaseDir }}`) if you need additional
  Spectr conventions or clarifications.
- Identify any vague or ambiguous details and ask the necessary follow-up
  questions before editing files.

Note: You are not implementing yet, you are fully planning and creating the change proposal using spectr.
{{end}}
```

**Claude Code** (`internal/initialize/templates/providers/claude-code/slash-apply.md.tmpl`):
```markdown
{{define "guardrails"}}
## Guardrails (Claude Code Edition)

- **Delegate implementation**: Use the `coder` subagent for implementation tasks, `tester` for validation, `stuck` for escalations.
- **Verify subagent work**: After coder completes a task, run `git diff --stat` and read files with significant changes (>10 lines modified).
- **Test every change**: Use the `tester` subagent after EVERY coder completion to verify the implementation.
- Favor straightforward, minimal implementations first and add complexity only
  when it is requested or clearly required.
- Keep changes tightly scoped to the requested outcome.
- Refer to `{{ .AgentsFile }}` and `{{ .BaseDir }}/project.md` (located inside
  the `{{ .BaseDir }}/` directory—run `ls {{ .BaseDir }}`) if you need
  additional Spectr conventions or clarifications.
{{end}}
```

**Codex** (`internal/initialize/templates/providers/codex/slash-proposal.md.tmpl`):
```markdown
{{define "steps"}}
## Steps

1. Review `{{ .ProjectFile }}`, read `{{ .SpecsDir }}/` and
   `{{ .ChangesDir }}/` directories, and inspect related code or docs (e.g.,
   via `rg`/`ls`) to ground the proposal in current behaviour.
2. **Check global prompts**: Review `~/.codex/prompts/` for existing Spectr-related prompts before creating new ones.
3. Choose a unique verb-led `change-id` and scaffold `proposal.md`,
   `tasks.md`, and `design.md` (when needed) under
   `{{ .ChangesDir }}/<id>/`.
4. Map the change into concrete capabilities or requirements, breaking
   multi-scope efforts into distinct spec deltas with clear relationships and
   sequencing.
5. Capture architectural reasoning in `design.md` when the solution spans
   multiple systems, introduces new patterns, or demands trade-off discussion
   before committing to specs.
6. Draft spec deltas in `{{ .ChangesDir }}/<id>/specs/<capability>/spec.md`
   (one folder per capability) using `## ADDED|MODIFIED|REMOVED Requirements`
   with at least one `#### Scenario:` per requirement.
7. **Reference AgentSkills**: If applicable, reference skills available in `~/.codex/skills/`.
8. Draft `tasks.md` as an ordered list of small, verifiable work items.
9. Validate with `spectr validate <id>` and resolve every issue before
   sharing the proposal.
{{end}}
```

## Backward Compatibility Strategy

### Phase 1: Add new infrastructure (no breaking changes)
- Add `ProviderTemplate` field to `TemplateRef` (optional, defaults to nil)
- Refactor generic templates to use `{{define}}` blocks
- Update `TemplateRef.Render()` to handle composition
- Add provider template parsing to `NewTemplateManager()`
- Extend `TemplateManager` interface with new methods

**Impact:** Zero - all existing code continues working

### Phase 2: Add provider-aware methods (backward compatible)
- Implement `ProviderSlashCommand()` and `ProviderTOMLSlashCommand()`
- Keep existing `SlashCommand()` and `TOMLSlashCommand()` unchanged
- Existing methods call new methods with empty provider ID

**Impact:** Zero - providers using old methods get generic templates

### Phase 3: Create provider templates (opt-in)
- Create `templates/providers/claude-code/` with overrides
- Create `templates/providers/codex/` with overrides
- Create `templates/providers/gemini/` (placeholder for future)

**Impact:** Zero until providers opt-in to using them

### Phase 4: Migrate providers (opt-in)
- Update ClaudeProvider, CodexProvider, GeminiProvider to use new methods
- Other 13 providers unchanged

**Impact:** Only affects 3 providers that explicitly opt-in

## Composition Algorithm Details

### Template Merging Process

1. **Base Template Parsing** (NewTemplateManager):
   ```go
   // All templates parsed into single template.Template
   // Template names: "slash-proposal.md.tmpl", "slash-apply.md.tmpl", etc.
   // Sections defined: "guardrails", "steps", "reference", "main"
   ```

2. **Provider Template Parsing** (NewTemplateManager):
   ```go
   // Each provider directory parsed into separate template.Template
   // Provider template names: same as base ("guardrails", "steps", etc.)
   // Stored in map: providerTemplates["claude-code"] = parsed template
   ```

3. **Template Composition** (TemplateRef.Render):
   ```go
   // Clone base template (preserves all base sections)
   composed := template.Must(tr.Template.Clone())

   // Merge provider overrides (last-wins for duplicate section names)
   for _, t := range tr.ProviderTemplate.Templates() {
       composed.AddParseTree(t.Name(), t.Tree) // Overwrites if exists
   }

   // Result: composed template with provider sections replacing base sections
   ```

4. **Template Execution**:
   ```go
   // Execute "main" template (the composition root)
   composed.ExecuteTemplate(&buf, "main", ctx)

   // "main" calls {{template "guardrails" .}}, {{template "steps" .}}, etc.
   // Each {{template}} directive resolves to either:
   // - Provider section (if defined in provider template)
   // - Base section (if not overridden)
   ```

### Last-Wins Semantics

Go's `AddParseTree()` implements last-wins automatically:
- If section "guardrails" exists in both base and provider, provider wins
- If section "steps" only exists in base, base is used
- If section "custom" only exists in provider, it's available but not used (unless "main" references it)

## Testing Strategy

### Unit Tests (internal/domain/template_test.go)

```go
func TestTemplateRef_Render_NoProviderTemplate(t *testing.T) {
    // Test: TemplateRef with nil ProviderTemplate uses base directly
}

func TestTemplateRef_Render_WithProviderOverride(t *testing.T) {
    // Test: Provider override replaces base section
}

func TestTemplateRef_composeTemplate(t *testing.T) {
    // Test: Composition merges base + provider correctly
    // Test: Provider sections override base sections
    // Test: Base sections used when provider doesn't override
}
```

### Integration Tests (internal/initialize/templates_test.go)

```go
func TestTemplateManager_ParseProviderTemplates(t *testing.T) {
    // Test: Provider directories are discovered and parsed
    // Test: Missing providers directory doesn't break initialization
}

func TestTemplateManager_ProviderSlashCommand(t *testing.T) {
    // Test: Returns TemplateRef with correct base and provider templates
    // Test: Unknown provider returns TemplateRef with nil ProviderTemplate
}

func TestTemplateManager_BackwardCompatibility(t *testing.T) {
    // Test: Existing SlashCommand() method still works
    // Test: Returns same result as ProviderSlashCommand("", cmd)
}
```

### End-to-End Validation

1. Run `spectr init --provider=claude-code`
2. Verify `.claude/commands/spectr/proposal.md` contains Claude-specific guardrails
3. Verify "orchestrator pattern" and "200k context" mentioned
4. Run `spectr init --provider=codex`
5. Verify `~/.codex/prompts/spectr-proposal.md` contains Codex-specific steps
6. Verify "global prompts" and "~/.codex/skills" mentioned
7. Run `spectr init --provider=aider` (not migrated)
8. Verify `.aider/commands/spectr/proposal.md` contains generic content

## Performance Considerations

### Template Parsing
- **When:** Once during `NewTemplateManager()` (application startup)
- **Cost:** Parse ~20 template files (base + providers)
- **Impact:** Negligible (happens once, not per-request)

### Template Composition
- **When:** Once per `TemplateRef.Render()` call
- **Cost:** Clone base template + merge provider overrides
- **Impact:** Low (composition is in-memory tree operations)
- **Frequency:** Only during `spectr init` (not during normal CLI usage)

### Memory Overhead
- **Base templates:** ~10 KB (2 Markdown + 2 TOML templates)
- **Provider templates:** ~2 KB per provider (only overrides)
- **Total:** < 50 KB for all providers
- **Impact:** Negligible

## Migration Path

### For Spectr Maintainers

1. **Phase 1** (Foundation):
   - Refactor generic templates to use `{{define}}` blocks
   - Add `ProviderTemplate` field to `TemplateRef`
   - Implement composition logic
   - Add unit tests

2. **Phase 2** (Infrastructure):
   - Update `TemplateManager` to parse provider directories
   - Add `ProviderSlashCommand()` and `ProviderTOMLSlashCommand()` methods
   - Extend interface
   - Add integration tests

3. **Phase 3** (Content):
   - Create `internal/initialize/templates/providers/claude-code/` with overrides
   - Create `internal/initialize/templates/providers/codex/` with overrides
   - Create `internal/initialize/templates/providers/gemini/` (placeholder)

4. **Phase 4** (Adoption):
   - Update ClaudeProvider to use `ProviderSlashCommand()`
   - Update CodexProvider to use `ProviderSlashCommand()`
   - Update GeminiProvider to use `ProviderTOMLSlashCommand()` (future)
   - Run end-to-end validation

### For Provider Authors (Future)

To add provider-specific templates:

1. Create directory: `internal/initialize/templates/providers/{provider-id}/`
2. Add override template: `slash-proposal.md.tmpl` with `{{define "section"}}` blocks
3. Update provider's `Initializers()` method to call `tm.ProviderSlashCommand("{provider-id}", cmd)`
4. Test with `spectr init --provider={provider-id}`

## Open Questions & Future Work

### Section Documentation
- **Question:** Where to document available sections and their purpose?
- **Answer:** Create `internal/initialize/templates/SECTIONS.md` listing all sections

### Section Naming Convention
- **Question:** lowercase (`guardrails`) vs title-case (`Guardrails`)?
- **Answer:** lowercase (follows Go template conventions)

### Validation
- **Question:** When to validate that provider overrides only define known sections?
- **Answer:** At template parse time in `NewTemplateManager()` - fail fast on unknown sections

### Future Providers
- All 16 providers can eventually have custom templates
- Prioritize based on usage metrics and user feedback
- Document customization guidelines for consistency

### TOML Template Composition
- Gemini currently uses TOML format
- Same composition approach applies
- Define sections in TOML templates: `{{define "guardrails"}}[guardrails]...{{end}}`
- Not urgent - implement when Gemini-specific customizations needed

## Success Criteria

- ✅ Generic templates refactored with `{{define}}` sections
- ✅ Provider-specific templates override only targeted sections
- ✅ Template composition merges base + provider correctly
- ✅ All existing providers continue working (backward compatibility)
- ✅ Claude Code uses custom templates with orchestrator pattern references
- ✅ Codex uses custom templates with global prompts references
- ✅ `spectr init` produces provider-specific slash commands
- ✅ Providers without custom templates fall back to generic
- ✅ No runtime performance degradation
- ✅ Clear documentation for future provider customizations
