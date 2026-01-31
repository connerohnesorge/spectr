# Proposal: Add Frontmatter Override System for Provider-Specific Slash Commands

## Problem

Claude Code now supports `context: fork` in slash command frontmatter to
run commands in a forked sub-agent context. We need to add this frontmatter
to the proposal slash command for Claude Code, but:

1. We don't want to duplicate entire templates for each provider variation
2. Hardcoding frontmatter in templates is inflexible
3. Other providers (or Claude Code's other commands) may need different
   overrides

The current system embeds frontmatter directly in `.tmpl` files, requiring
duplication for provider-specific variations. This causes:

## Solution

Implement an intelligent frontmatter override system that:

1. **Looks up** base frontmatter from the frontmatter package
2. **Merges** provider-specific overrides intelligently
3. **Renders** the merged frontmatter back to YAML
4. **Inserts** the final frontmatter into the generated file

This allows:

- Base frontmatter stored centrally in the package (not in templates)

- Providers to override specific fields (add `context: fork`, remove `agent`)

- Clean, maintainable template reuse without duplication or parsing from templates

## Architecture

```text
Frontmatter Package (base frontmatter map)
  ↓ lookup base frontmatter
Base Frontmatter (map[string]interface{})
  ↓ + Provider overrides (map[string]interface{})
Merged Frontmatter (map[string]interface{})
  ↓ render to YAML + template body
Final Slash Command (.md file)

```

### Key Components

1. **FrontmatterOverride** struct in `internal/domain`:
   - `Set map[string]interface{}` - fields to add/modify
   - `Remove []string` - fields to delete

2. **TemplateManager.SlashCommandWithOverrides()** method:
   - Takes `cmd domain.SlashCommand` and `overrides *FrontmatterOverride`
   - Returns modified `domain.TemplateRef`

3. **ClaudeProvider** updated to pass overrides:
   - Proposal command: add `context: fork`, remove `agent`
   - Apply command: no overrides (unchanged)

## Implementation

### Phase 1: Core Frontmatter System

- Create base frontmatter map in `internal/frontmatter` package
- Implement `FrontmatterOverride` type in `internal/domain` with Set and Remove
  fields
- Add YAML rendering utilities to convert merged map back to YAML
- Add tests for frontmatter merge logic and override application

### Phase 2: TemplateManager Extension

- Add `SlashCommandWithOverrides()` method to lookup frontmatter
- Apply provider-specific overrides (Set then Remove)
- Render merged frontmatter map to YAML
- Combine rendered YAML with template body to create final slash command file

### Phase 3: ClaudeProvider Integration

- Update `ClaudeProvider.Initializers()` to use overrides

- Add `context: fork` to proposal command

- Remove `agent: plan` from proposal command

### Phase 4: Validation & Testing

- Test with `spectr init --provider=claude-code`

- Verify generated `.claude/commands/spectr/proposal.md` has correct frontmatter

- Run existing tests to ensure no regression

## Benefits

- **No template duplication**: Single source of truth for command content
- **Flexible**: Easy to add provider-specific customizations
- **Maintainable**: Frontmatter changes don't require duplicating templates
- **Type-safe**: Go structs for overrides prevent typos
- **Future-proof**: Other providers can use same mechanism

## Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| YAML parsing errors | Comprehensive error handling and tests |
| Breaking existing providers | Defaults keep current behavior, overrides opt-in |
| Complex merge logic | Overrides win, then removes execute |

## Open Questions

None - approach is straightforward and well-defined.
