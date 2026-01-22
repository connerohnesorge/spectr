# Design: Frontmatter Override System

## Problem Statement

Claude Code now supports `context: fork` in slash command frontmatter to run
commands in forked sub-agent contexts. We need to customize the proposal
command's frontmatter for Claude Code by:

- Adding `context: fork`
- Removing `agent: plan` (not supported by Claude Code slash commands)

However, duplicating entire template files for provider-specific variations is unmaintainable.

## Goals

1. **Zero duplication**: Single template source for each command
2. **Provider flexibility**: Each provider can customize frontmatter as needed
3. **Type safety**: Compile-time guarantees for override operations
4. **Backward compatibility**: Existing providers work unchanged
5. **TOML support**: Works for both Markdown (Claude, Windsurf) and
   TOML (Gemini) formats

## Non-Goals

- Complex frontmatter queries (e.g., JSONPath selectors)
- Conditional overrides based on user config
- Frontmatter validation (handled elsewhere)
- Per-file customization (only per-provider)

## Architecture Overview

### Data Flow

```text
┌─────────────────────────────────────────────────────────────┐
│ 1. Provider.Initializers() calls tm.SlashCommand()         │
│    with optional FrontmatterOverride                        │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. TemplateManager returns TemplateRef with overrides       │
│    TemplateRef { Name, Template, Overrides }                │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 3. SlashCommandsInitializer.Init() creates TemplateContext  │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 4. TemplateRef.Render(ctx) executes Go template            │
│    Result: "---\nkey: val\n---\nBody content"              │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼ (if Overrides != nil)
┌─────────────────────────────────────────────────────────────┐
│ 5. ParseFrontmatter() extracts YAML from rendered content   │
│    Returns: map[string]interface{}, body string             │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 6. ApplyFrontmatterOverrides() merges Set/Remove ops        │
│    Set: merge into map, Remove: delete keys                 │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 7. RenderFrontmatter() serializes back to YAML              │
│    Returns: "---\nmodified: val\n---\nBody"                 │
└───────────────────────────┬─────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────┐
│ 8. Write final content to .claude/commands/spectr/*.md      │
└─────────────────────────────────────────────────────────────┘
```

## Component Design

### 1. Domain Types (`internal/domain/frontmatter.go`)

#### FrontmatterOverride

```go
// FrontmatterOverride specifies modifications to slash command frontmatter.
// Set operations are applied first, then Remove operations.
type FrontmatterOverride struct {
    // Set contains fields to add or modify.
    // Values must be YAML-serializable (string, bool, int, []string, map, etc.)
    Set map[string]interface{}

    // Remove contains field names to delete.
    // Applied after Set to allow field replacement.
    Remove []string
}
```

**Design choice**: Separate Set/Remove maps rather than using nil values
because:

- Clear intent: `Set{"context": "fork"}` vs
  `Set{"context": "fork", "agent": nil}`
- Type safety: nil has different meanings in Go (missing vs null)
- Explicit ordering: Remove always happens after Set

#### ParseFrontmatter

**Note**: This function is no longer needed in the current design since
templates don't contain frontmatter. Frontmatter is defined in Go code and
assembled dynamically.

If needed for other purposes (e.g., reading existing files), the signature
would be:

```go
// ParseFrontmatter extracts YAML frontmatter from markdown content.
// Returns the frontmatter as a map, the body content, and any error.
func ParseFrontmatter(content string) (map[string]interface{}, string, error)
```

But for our use case (generating files), we only need `RenderFrontmatter()`.

#### ApplyFrontmatterOverrides

```go
// ApplyFrontmatterOverrides applies Set and Remove operations to a frontmatter map.
// Returns a new map; does not mutate the input.
//
// Operation order:
//   1. Copy base map
//   2. Apply all Set operations (merge/overwrite)
//   3. Apply all Remove operations (delete keys)
func ApplyFrontmatterOverrides(
    base map[string]interface{},
    overrides *FrontmatterOverride,
) map[string]interface{}
```

**Design choices**:

- Returns new map (immutable operation) for safety
- Nil overrides returns copy of base (no-op)
- Set always overwrites existing values
- Remove ignores non-existent keys (idempotent)

**Open question**: Should Remove support nested paths like "metadata.author"?

#### RenderFrontmatter

```go
// RenderFrontmatter serializes a frontmatter map to YAML and combines with body.
// Returns markdown with YAML frontmatter block.
//
// Output format:
//   ---
//   key: value
//   ---
//   Body content
func RenderFrontmatter(fm map[string]interface{}, body string) (string, error)
```

**Design choices**:

- Always includes `---` delimiters even for empty map
- Uses `yaml.v3` for consistent formatting
- Preserves map key order if possible (YAML 1.2 behavior)

### 2. Base Frontmatter Definition (`internal/domain/frontmatter.go`)

Frontmatter is defined as Go data structures, not embedded in templates:

```go
// BaseSlashCommandFrontmatter defines default frontmatter for each slash command.
// Templates (.tmpl files) contain only body content; frontmatter is data.
var BaseSlashCommandFrontmatter = map[domain.SlashCommand]map[string]interface{}{
    domain.SlashProposal: {
        "description":   "Proposal Creation Guide (project)",
        "allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
        "agent":         "plan",
        "subtask":       false,
    },
    domain.SlashApply: {
        "description":   "Change Proposal Application/Acceptance Process (project)",
        "allowed-tools": "Read, Glob, Grep, Write, Edit, Bash(spectr:*)",
        "subtask":       false,
    },
}

// GetBaseFrontmatter returns a copy of the base frontmatter for a command.
// Returns empty map if command not found.
func GetBaseFrontmatter(cmd domain.SlashCommand) map[string]interface{} {
    base, ok := BaseSlashCommandFrontmatter[cmd]
    if !ok {
        return make(map[string]interface{})
    }

    // Return deep copy to prevent mutation
    return copyMap(base)
}
```

**Design choices**:

- Frontmatter is data (maps), not text in templates
- Templates (.tmpl files) contain only body content (no `---` blocks)
- Base frontmatter is versioned in Go code (type-safe, tracked in git)
- Deep copy prevents accidental mutation of base maps

### 3. TemplateRef Extension (`internal/domain/template.go`)

```go
type TemplateRef struct {
    Name      string                  // Template file name
    Template  *template.Template      // Pre-parsed Go template
    Command   domain.SlashCommand     // Slash command type (for frontmatter lookup)
    Overrides *FrontmatterOverride    // Optional frontmatter modifications
}

func (tr TemplateRef) Render(ctx *TemplateContext) (string, error) {
    // 1. Execute Go template to get body content
    var buf bytes.Buffer
    if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
        return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
    }
    body := buf.String()

    // 2. Get base frontmatter for this command
    fm := GetBaseFrontmatter(tr.Command)

    // 3. Apply overrides if present
    if tr.Overrides != nil {
        fm = ApplyFrontmatterOverrides(fm, tr.Overrides)
    }

    // 4. Render frontmatter as YAML + body content
    return RenderFrontmatter(fm, body)
}
```

**Design choices**:

- Added `Command` field to TemplateRef for frontmatter lookup
- Template body rendered first (resolves {{.BaseDir}} etc.)
- Frontmatter assembled dynamically: base + overrides
- Standard markdown `---` fences via RenderFrontmatter()
- Nil overrides = use base frontmatter unchanged

**Key insight**: Templates are pure content, frontmatter is pure data. This
separates concerns and makes both easier to maintain.

### 4. TemplateManager Extension (`internal/initialize/templates.go`)

```go
// SlashCommandWithOverrides returns a Markdown template with frontmatter overrides.
// Used when providers need to customize slash command frontmatter.
func (tm *TemplateManager) SlashCommandWithOverrides(
    cmd domain.SlashCommand,
    overrides *FrontmatterOverride,
) domain.TemplateRef {
    ref := tm.SlashCommand(cmd)  // Get base template
    ref.Command = cmd             // Set command for frontmatter lookup
    ref.Overrides = overrides     // Add overrides
    return ref
}

// TOMLSlashCommandWithOverrides returns a TOML template with frontmatter overrides.
// TOML uses different format but same override logic.
func (tm *TemplateManager) TOMLSlashCommandWithOverrides(
    cmd domain.SlashCommand,
    overrides *FrontmatterOverride,
) domain.TemplateRef {
    ref := tm.TOMLSlashCommand(cmd)
    ref.Overrides = overrides
    return ref
}
```

**Design choices**:

- Separate methods for Markdown and TOML (clear intent)
- Returns modified copy, doesn't mutate cached TemplateRef
- Nil overrides behaves like SlashCommand() (backward compatible)

**Note on TOML**: Per interview decision, we're only supporting Markdown YAML
frontmatter initially. TOML support can be added later if Gemini needs it.

### 5. Template File Changes

All `.tmpl` files will have frontmatter removed:

**Before** (`internal/domain/templates/slash-proposal.md.tmpl`):

```yaml
---
description: Proposal Creation Guide (project)
allowed-tools: Read, Glob, Grep, Write, Edit, Bash(spectr:*)
agent: plan
subtask: false
---

# Proposal Creation Guide

## Guardrails
...
```

**After**:

```markdown
# Proposal Creation Guide

## Guardrails
...
```

Frontmatter moved to `BaseSlashCommandFrontmatter` map in Go code.

### 6. Provider Usage Example

```go
// ClaudeProvider.Initializers() - BEFORE
return []Initializer{
    NewSlashCommandsInitializer(
        ".claude/commands/spectr",
        map[domain.SlashCommand]domain.TemplateRef{
            domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
            domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
        },
    ),
}

// ClaudeProvider.Initializers() - AFTER
return []Initializer{
    NewSlashCommandsInitializer(
        ".claude/commands/spectr",
        map[domain.SlashCommand]domain.TemplateRef{
            domain.SlashProposal: tm.SlashCommandWithOverrides(
                domain.SlashProposal,
                &domain.FrontmatterOverride{
                    Set:    map[string]interface{}{"context": "fork"},
                    Remove: []string{"agent"},
                },
            ),
            domain.SlashApply: tm.SlashCommand(domain.SlashApply), // No overrides
        },
    ),
}
```

## Trade-offs & Alternatives

### Alternative 1: Provider-Specific Templates

**Approach**: Create `slash-proposal-claude.md.tmpl` with hardcoded frontmatter.

**Pros**:

- Simple: no parsing/merging logic
- Explicit: frontmatter visible in template

**Cons**:

- Duplication: must maintain N copies of each template
- Brittle: body changes require updating all variants
- Discovery: hard to find all customizations

**Rejected**: Violates DRY principle and makes maintenance difficult.

### Alternative 2: Frontmatter in Provider Config

**Approach**: Providers define frontmatter in Go structs, completely override
template frontmatter.

**Pros**:

- Full control: provider owns entire frontmatter
- Type-safe: Go structs validate at compile time

**Cons**:

- Duplication: every provider must specify all frontmatter fields
- Fragile: default changes break all providers
- Verbose: simple overrides require full struct definitions

**Rejected**: Too much boilerplate for simple customizations.

### Alternative 3: Post-Processing Hook

**Approach**: SlashCommandsInitializer calls a hook after Render() to modify
content.

**Pros**:

- Flexible: can modify any part of content
- Separation: frontmatter logic separate from template system

**Cons**:

- Generic: hook signature `func(string) string` loses type safety
- Complex: providers must implement their own parsing
- Error-prone: string manipulation is fragile

**Rejected**: Too generic, loses benefits of structured frontmatter.

## Design Decisions (from Interview)

All open questions have been resolved:

### 1. TOML Support Scope

**Decision**: Markdown YAML only initially. Gemini doesn't need overrides yet,
we can add TOML support later if needed.

### 2. Nested Field Removal

**Decision**: Top-level keys only. `Remove: []string{"metadata"}` removes
entire section. Simpler code, and we don't have nested frontmatter currently.

### 3. Error Handling Philosophy

**Decision**: Strict mode. Return error if YAML rendering fails. Fail fast to
catch bugs early during template development.

### 4. Frontmatter Source

**Decision**: Define base frontmatter in Go code (`BaseSlashCommandFrontmatter`
map), not in .tmpl files. Templates contain only body content. This separates
data from presentation.

### 5. Override Precedence

**Decision**: Set first, then Remove. If same field in both, Remove wins. Clear
precedence rule: `Set{"context": "fork"}, Remove{"context"}` results in field
removed.

## Implementation Plan

1. **Phase 1**: Core frontmatter utilities (parse, apply, render)
2. **Phase 2**: Extend TemplateRef to support overrides
3. **Phase 3**: Add TemplateManager helper methods
4. **Phase 4**: Update ClaudeProvider to use overrides
5. **Phase 5**: Tests and validation

## Success Criteria

- [ ] Generated `.claude/commands/spectr/proposal.md` has `context: fork`
- [ ] Generated `.claude/commands/spectr/proposal.md` does NOT have `agent:` field
- [ ] Generated `.claude/commands/spectr/apply.md` unchanged from base template
- [ ] All existing tests pass (backward compatibility)
- [ ] New tests cover parse/override/render logic
- [ ] `spectr init --provider=claude-code` works end-to-end
