# Technical Design: Auto-Install Slash Commands

## Context

The init system currently has two separate tool types:
1. **Config-based tools** (6 tools): Create single instruction files like `CLAUDE.md`
2. **Slash command tools** (11 tools): Create 3 command files in tool-specific directories

This separation creates UX friction - users must understand the distinction and select both types to get complete Spectr integration. OpenSpec solved this by having each tool selection install everything relevant to that tool.

## Goals / Non-Goals

**Goals**:
- Automatic slash command installation when config-based tool is selected
- Remove redundant tool entries from wizard (keep one entry per tool)
- Maintain backward compatibility with existing slash command files
- Simple implementation leveraging existing configurator infrastructure

**Non-Goals**:
- Changing the configurator implementation details
- Modifying template content or frontmatter
- Adding new UI elements or wizard screens
- Supporting partial installations (user always gets both config + slash commands)

## Decisions

### Decision 1: Sequential Invocation vs Composite Configurator

**Chosen**: Sequential invocation in executor

**Rationale**:
- Simpler to implement - reuse existing `getConfigurator()` logic
- No new configurator types needed
- Clearer separation of concerns (executor orchestrates, configurators execute)
- Easier to test and debug

**Alternative considered**: Composite configurator pattern
- Would require new `CompositeConfigurator` wrapper class
- Adds complexity for minimal benefit
- Harder to maintain with 11 tool pairs

**Implementation**:
```go
// In executor.go configureTools():
for _, toolID := range selectedToolIDs {
    // 1. Configure the main tool (config file)
    configurator := e.getConfigurator(toolID)
    if configurator != nil {
        configurator.Configure(projectPath, spectrDir)
    }

    // 2. Check if tool has slash command equivalent
    if slashToolID, hasSlash := getSlashToolMapping(toolID); hasSlash {
        slashConfig := e.getConfigurator(slashToolID)
        if slashConfig != nil {
            slashConfig.Configure(projectPath, spectrDir)
        }
    }
}
```

### Decision 2: Mapping Structure

**Chosen**: Simple map in registry.go

**Rationale**:
- Explicit and easy to understand
- Centralized in one location
- Type-safe with compile-time checking
- Easy to extend for new tools

**Alternative considered**: Convention-based (derive slash ID from config ID)
- Would work for claude-code→claude pattern
- Breaks for tools like cline-config→cline, cursor-config→cursor
- Implicit behavior harder to debug

**Implementation**:
```go
// internal/init/registry.go
var configToSlashMapping = map[string]string{
    "claude-code":      "claude",
    "cline":            "cline",
    "cursor":           "cursor",
    "costrict-config":  "costrict",
    "qoder-config":     "qoder",
    "codebuddy":        "codebuddy",
    "qwen":             "qwen",
    // ... remaining mappings
}

func getSlashToolMapping(configToolID string) (string, bool) {
    slashID, exists := configToSlashMapping[configToolID]
    return slashID, exists
}
```

### Decision 3: Registry Cleanup Strategy

**Chosen**: Remove slash-only tools entirely from registry

**Rationale**:
- Eliminates duplication and confusion
- Forces users toward the "correct" selection (config-based)
- Simplifies wizard display (fewer items)
- Matches user's stated preference: "no change - just auto-install silently"

**Alternative considered**: Keep both, mark slash-only as deprecated
- Would maintain backward compat for users who somehow prefer slash-only
- Adds complexity and keeps redundant options visible
- Delays the UX improvement

**Migration impact**:
- Users with existing projects won't be affected (files already created)
- New users get simpler, better experience
- `spectr update` could theoretically detect and migrate, but not necessary

### Decision 4: Tool Naming - Keep Current Names

**Chosen**: No changes to tool display names

**Rationale**:
- User selected "No change - just auto-install silently"
- Keeps wizard clean and familiar
- Slash command installation is an implementation detail

**Alternative considered**: Add suffix like "Claude Code (config + commands)"
- More explicit but verbose
- Would require wizard UI changes
- Contradicts user's preference

## Risks / Trade-offs

### Risk 1: Breaking change for slash-only users

**Risk**: Users who previously selected only slash command tools (not config) won't see them anymore

**Mitigation**:
- Document in changelog/release notes
- These users can now select the config-based tool to get both
- Edge case - most users likely selected both anyway

**Trade-off accepted**: Better UX for 95% of users worth the breaking change

### Risk 2: Tool count mismatch in wizard

**Risk**: Wizard currently shows "17 AI tools" - this will become ~6 after cleanup

**Mitigation**:
- Update any hardcoded tool counts in code/docs
- Test wizard display with new count
- Verify navigation still works correctly

**Trade-off accepted**: Accurate count is better than inflated count

### Risk 3: Confusion about what gets installed

**Risk**: Users might not know slash commands are being created

**Mitigation**:
- Completion screen shows all created files (both config and slash)
- `ExecutionResult` tracks both file types
- Users will see the files listed in success output

**Trade-off accepted**: Silent installation is simpler than adding UI complexity

## Migration Plan

### For New Projects
1. Run `spectr init`
2. Select config-based tools (e.g., `claude-code`)
3. Both config file and slash commands are created automatically
4. No additional steps needed

### For Existing Projects
1. No action required - existing files are preserved
2. If they run `spectr init` again, configurators respect existing files
3. `spectr update` behavior unchanged (only updates instruction content)

### Rollback Strategy
If this change causes issues:
1. Revert registry.go to restore slash-only tool entries
2. Revert executor.go to remove auto-invocation logic
3. No data loss - files remain on disk
4. Users can manually delete unwanted slash command files

## Open Questions

None - user clarifications received:
- ✅ Automatic installation confirmed
- ✅ All matching pairs confirmed
- ✅ Silent installation (no UI changes) confirmed
