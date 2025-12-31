# Design: Provider-Specific Template Overrides

## Architecture

### Template Composition with {{define}} Blocks

Use Go's `text/template` `{{define}}` directive for section-based composition:

**Base Template Structure** (slash-proposal.md.tmpl):
```markdown
{{define "base_guardrails"}}
- Favor straightforward, minimal implementations
- Keep changes tightly scoped to the requested outcome
- Refer to spectr/AGENTS.md and spectr/project.md for conventions
{{end}}

{{define "guardrails"}}
## Guardrails
{{template "base_guardrails" .}}
{{end}}

{{define "steps"}}
## Steps
1. Review spectr/project.md...
2. Choose a unique change-id...
{{end}}

{{define "reference"}}
## Reference
- Read delta specs at spectr/changes/<id>/specs/...
{{end}}

{{define "main"}}
# Proposal Creation Guide
{{template "guardrails" .}}
{{template "steps" .}}
{{template "reference" .}}
{{end}}
```

**Provider Override** (providers/claude-code/slash-proposal.md.tmpl):
```markdown
{{define "guardrails"}}
## Guardrails (Claude Code Edition)
- **Use the orchestrator pattern**: Delegate to coder/tester/stuck subagents
- **Leverage 200k context**: Comprehensive analysis and planning
{{template "base_guardrails" .}}
{{end}}
```

### Why This Approach

- **Native feature**: Go's standard `text/template` package
- **Type-safe**: Compile-time template validation
- **Clean composition**: `Clone()` + `AddParseTree()` with last-wins semantics
- **Provider flexibility**: Providers can reuse base sections via `{{template "base_X" .}}`
- **Backward compatible**: Minimal changes to existing templates

## Directory Structure

```
internal/
  domain/
    templates/
      slash-proposal.md.tmpl          # Base template (refactored with {{define}})
      slash-apply.md.tmpl             # Base template (refactored with {{define}})
      slash-proposal.toml.tmpl        # Base TOML (unchanged for now)
      slash-apply.toml.tmpl           # Base TOML (unchanged for now)

  initialize/
    templates/
      providers/                      # NEW: Provider overrides
        claude-code/
          slash-proposal.md.tmpl      # Override "guardrails" section
          slash-apply.md.tmpl         # Override "guardrails" section
        codex/
          slash-proposal.md.tmpl      # Override "steps" and "reference"
          slash-apply.md.tmpl         # Override "reference"
        gemini/
          .gitkeep                    # Placeholder for future TOML templates
```

## Code Implementation Details

### 1. Update TemplateRef (internal/domain/template.go)

```go
type TemplateRef struct {
    Name             string             // "slash-proposal.md.tmpl"
    Template         *template.Template // Base template with all sections
    ProviderTemplate *template.Template // Provider overrides (nil if none)
}

func (tr TemplateRef) Render(ctx *TemplateContext) (string, error) {
    // No provider overrides - use base directly
    if tr.ProviderTemplate == nil {
        var buf bytes.Buffer
        if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
            return "", fmt.Errorf("failed to render %s: %w", tr.Name, err)
        }
        return buf.String(), nil
    }

    // Compose base + provider
    composed, err := tr.composeTemplate()
    if err != nil {
        return "", err
    }

    // Execute "main" (composition entry point)
    var buf bytes.Buffer
    if err := composed.ExecuteTemplate(&buf, "main", ctx); err != nil {
        return "", fmt.Errorf("failed to render composed template: %w", err)
    }
    return buf.String(), nil
}

func (tr TemplateRef) composeTemplate() (*template.Template, error) {
    // Clone base template
    composed := template.Must(tr.Template.Clone())

    // Merge provider overrides (last-wins semantics)
    for _, t := range tr.ProviderTemplate.Templates() {
        if _, err := composed.AddParseTree(t.Name(), t.Tree); err != nil {
            return nil, fmt.Errorf("failed to merge provider template: %w", err)
        }
    }

    return composed, nil
}
```

**Key Design:**
- `ProviderTemplate` is optional (nil = no composition)
- Composition is lazy (happens in `Render()`, not during manager initialization)
- Clone base template to avoid mutations
- `AddParseTree()` implements last-wins automatically

### 2. Update TemplateManager (internal/initialize/templates.go)

```go
type TemplateManager struct {
    templates         *template.Template
    providerTemplates map[string]*template.Template  // NEW
}

func NewTemplateManager() (*TemplateManager, error) {
    // Parse base templates (existing logic)
    mainTmpl, err := template.ParseFS(templateFS, "templates/**/*.tmpl")
    // ... error handling ...

    // Parse domain templates (existing logic)
    domainTmpl, err := template.ParseFS(domain.TemplateFS, "templates/*.tmpl")
    // ... merge into mainTmpl ...

    // NEW: Parse provider templates
    providerTmpls := make(map[string]*template.Template)

    entries, err := fs.ReadDir(templateFS, "templates/providers")
    if err != nil {
        // No providers directory yet - OK (backward compatibility)
        return &TemplateManager{
            templates:         mainTmpl,
            providerTemplates: providerTmpls,
        }, nil
    }

    for _, entry := range entries {
        if !entry.IsDir() {
            continue
        }

        providerID := entry.Name()
        pattern := fmt.Sprintf("templates/providers/%s/*.tmpl", providerID)

        // Parse provider templates
        providerTmpl, err := template.ParseFS(templateFS, pattern)
        if err != nil {
            // Provider directory exists but parsing failed - validate at init time
            return nil, fmt.Errorf(
                "failed to parse provider %s templates: %w",
                providerID,
                err,
            )
        }

        // Validate: provider templates only define known sections
        if providerTmpl != nil {
            if err := validateProviderTemplate(mainTmpl, providerTmpl); err != nil {
                return nil, fmt.Errorf(
                    "provider %s template validation failed: %w",
                    providerID,
                    err,
                )
            }
            providerTmpls[providerID] = providerTmpl
        }
    }

    return &TemplateManager{
        templates:         mainTmpl,
        providerTemplates: providerTmpls,
    }, nil
}

// validateProviderTemplate checks that provider only overrides known sections
func validateProviderTemplate(base, provider *template.Template) error {
    // Provider can override: guardrails, steps, reference, main
    // Provider can call: base_guardrails, base_steps, base_reference
    knownSections := map[string]bool{
        "guardrails":      true,
        "steps":           true,
        "reference":       true,
        "main":            true,
        "base_guardrails": true,
        "base_steps":      true,
        "base_reference":  true,
    }

    for _, t := range provider.Templates() {
        if !knownSections[t.Name()] {
            return fmt.Errorf("provider template defines unknown section: %s", t.Name())
        }
    }
    return nil
}

// ProviderSlashCommand returns base + provider templates
func (tm *TemplateManager) ProviderSlashCommand(
    providerID string,
    cmd domain.SlashCommand,
) domain.TemplateRef {
    names := map[domain.SlashCommand]string{
        domain.SlashProposal: "slash-proposal.md.tmpl",
        domain.SlashApply:    "slash-apply.md.tmpl",
    }

    return domain.TemplateRef{
        Name:             names[cmd],
        Template:         tm.templates,
        ProviderTemplate: tm.providerTemplates[providerID], // nil if not found
    }
}

// ProviderTOMLSlashCommand returns base + provider TOML templates
func (tm *TemplateManager) ProviderTOMLSlashCommand(
    providerID string,
    cmd domain.SlashCommand,
) domain.TemplateRef {
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

// SlashCommand (backward compatible - calls ProviderSlashCommand with empty ID)
func (tm *TemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
    return tm.ProviderSlashCommand("", cmd)
}
```

**Key Design:**
- Parse-time validation (fail fast on unknown sections)
- Graceful degradation (missing providers directory is OK)
- Provider ID lookup returns nil if not found (automatic fallback)
- Base sections can be called with "base_" prefix (for provider reuse)

### 3. Update TemplateManager Interface (internal/initialize/providers/provider.go)

```go
type TemplateManager interface {
    // Existing methods (unchanged)
    InstructionPointer() domain.TemplateRef
    Agents() domain.TemplateRef
    SlashCommand(cmd domain.SlashCommand) domain.TemplateRef
    TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef
    SkillFS(skillName string) (fs.FS, error)

    // NEW: Provider-aware methods
    ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef
    ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef
}
```

### 4. Provider Implementation Changes

**Claude Code** (internal/initialize/providers/claude.go):
```go
func (*ClaudeProvider) Initializers(_ context.Context, tm TemplateManager) []Initializer {
    return []Initializer{
        // ... existing initializers ...
        NewSlashCommandsInitializer(".claude/commands/spectr",
            map[domain.SlashCommand]domain.TemplateRef{
                // Use provider-specific templates
                domain.SlashProposal: tm.ProviderSlashCommand("claude-code", domain.SlashProposal),
                domain.SlashApply:    tm.ProviderSlashCommand("claude-code", domain.SlashApply),
            },
        ),
        // ... rest of initializers ...
    }
}
```

**Codex** (internal/initialize/providers/codex.go):
```go
func (*CodexProvider) Initializers(_ context.Context, tm TemplateManager) []Initializer {
    return []Initializer{
        // ... existing initializers ...
        NewHomeSlashCommandsInitializer("spectr",
            map[domain.SlashCommand]domain.TemplateRef{
                domain.SlashProposal: tm.ProviderSlashCommand("codex", domain.SlashProposal),
                domain.SlashApply:    tm.ProviderSlashCommand("codex", domain.SlashApply),
            },
        ),
        // ... rest of initializers ...
    }
}
```

**Other Providers**: No changes needed (continue using `tm.SlashCommand()`)

## Template Section Reference

This section documents all available sections for template authors:

### Base Sections (defined in generic template)

- **`base_guardrails`**: Core best practices (don't override, embed with `{{template "base_guardrails" .}}`)
- **`base_steps`**: Core workflow steps (don't override, embed with `{{template "base_steps" .}}`)
- **`base_reference`**: Core reference links (don't override, embed with `{{template "base_reference" .}}`)

### Composite Sections (can be overridden)

- **`guardrails`**: Wraps `base_guardrails`, providers add tool-specific guidance
- **`steps`**: Numbered workflow steps, providers can insert tool-specific steps
- **`reference`**: Documentation links and references, providers add tool-specific resources

### Entry Point

- **`main`**: Composition root, calls `{{template "guardrails" .}}`, `{{template "steps" .}}`, `{{template "reference" .}}`

### Provider Override Pattern

Providers override composite sections and call base versions:

```markdown
{{define "guardrails"}}
## Guardrails (Provider Edition)
- Tool-specific guidance
- More guidance
{{template "base_guardrails" .}}
{{end}}
```

## Backward Compatibility

### Phase 1: Foundation (no breaking changes)
- Add `ProviderTemplate` field (optional, defaults to nil)
- Refactor templates with `{{define}}` blocks
- Existing `SlashCommand()` still works unchanged

### Phase 2: Infrastructure (backward compatible)
- Parse provider directories (missing directory is OK)
- Implement new methods (`ProviderSlashCommand`, `ProviderTOMLSlashCommand`)
- Old methods unchanged (call new methods internally)

### Phase 3: Provider Migration (opt-in)
- Update 2-3 providers to use new methods
- Other 13 providers unchanged (continue using old methods)
- Both old and new methods work simultaneously

### Phase 4: Long-term (optional future)
- Document migration path for remaining providers
- No timeline for full migration (providers can opt-in at their own pace)

## Testing Strategy

### Unit Tests (internal/domain/template_test.go)
- TemplateRef with nil ProviderTemplate renders base directly
- TemplateRef with provider template composes correctly
- Provider section overrides base section (last-wins)
- Base section used when provider doesn't override

### Integration Tests (internal/initialize/templates_test.go)
- Provider directories are discovered and parsed
- Unknown sections in provider template cause parse error
- Composition produces correct merged template
- All 16 providers initialize successfully
- Both old and new methods work simultaneously

### End-to-End Validation
- `spectr init --provider=claude-code` produces Claude-specific templates
- `spectr init --provider=codex` produces Codex-specific templates
- `spectr init --provider=aider` produces generic templates (unchanged)
- All template variables correctly rendered

## Performance

- **Template parsing**: Once at application startup (negligible cost)
- **Template composition**: Once per `Render()` call (in-memory tree operations, low cost)
- **Frequency**: Only during `spectr init` (not in hot paths)
- **Memory overhead**: < 50 KB for all providers combined

## Future Extensions

### Adding New Provider Templates

To add custom templates for a new provider:

1. Create `internal/initialize/templates/providers/{provider-id}/` directory
2. Add `slash-proposal.md.tmpl` with `{{define}}` section overrides
3. Update provider's `Initializers()` to call `tm.ProviderSlashCommand("{id}", cmd)`
4. Test with `spectr init --provider={id}`

### TOML Template Customization

Gemini currently uses TOML format. Same composition approach applies:

```toml
{{define "guardrails"}}
[guardrails]
items = [
    "Tool-specific guidance",
    # ... rest of content
]
{{end}}
```

Currently placeholder (`gemini/.gitkeep`), implement when Gemini-specific customizations are needed.

## Open Design Decisions

**Section Documentation**: Kept in design.md (single source of truth, updated with template changes)

**Section Validation**: Parse-time validation in `NewTemplateManager()` - fail fast on unknown sections

**Provider Reuse Pattern**: Providers can call `{{template "base_X" .}}` to include base section content
