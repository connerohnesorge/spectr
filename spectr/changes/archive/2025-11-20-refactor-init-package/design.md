# Design: Refactor internal/init Package

## Context

The `internal/init` package manages Spectr initialization, creating directory structures, templates, and configuring AI tool integrations. Currently handles 7 AI tools with both config-file and slash-command integration methods. The package has grown to ~2,000 lines with significant duplication as new tools were added incrementally without revisiting the architecture.

### Current Pain Points
- Adding a new AI tool requires touching 5+ files and writing ~100 lines of boilerplate
- No compile-time safety for tool IDs (all strings, prone to typos)
- Three different switch statements must stay in sync manually
- Template rendering inconsistent (some use TemplateManager, some use fmt.Sprintf)
- Constants duplicated across files with inconsistent naming

## Goals / Non-Goals

### Goals
- Reduce code duplication by 60%+ (from ~1,400 to ~550 lines in affected files)
- Make adding new tools declarative (data-driven) instead of imperative (code-driven)
- Improve type safety with tool ID constants
- Consolidate all tool configuration into single registry
- Extract reusable utilities (marker file updates, template rendering)
- Maintain 100% backward compatibility (all existing tests pass unchanged)

### Non-Goals
- Not changing public API or CLI interface
- Not modifying template content or file formats
- Not refactoring wizard.go TUI logic (separate concern)
- Not adding new tool integrations in this change
- Not changing test coverage requirements

## Decisions

### Decision 1: Data-Driven Tool Registry

**Choice**: Single tool registry with declarative configuration structs

**Why**: Eliminates 18 configurator structs and 15 factory functions

**Implementation**:
```go
// Before: 18 separate structs, each with 3 methods
type ClaudeCodeConfigurator struct{}
func (*ClaudeCodeConfigurator) Configure(...) error { /* 15 lines */ }
func (*ClaudeCodeConfigurator) IsConfigured(...) bool { /* 3 lines */ }
func (*ClaudeCodeConfigurator) GetName() string { return "Claude Code" }

// After: Single struct + data
type ToolConfig struct {
    ID           ToolID
    Name         string
    Type         ToolType
    ConfigFile   string  // for config-based tools
    SlashPaths   map[string]string // for slash command tools
    Frontmatter  map[string]string // for slash commands
}

type GenericConfigurator struct {
    config ToolConfig
}
```

**Alternatives Considered**:
- Keep individual structs, extract common logic: Still requires maintaining all structs
- Interface-based plugin system: Over-engineered for fixed set of 7 tools
- Code generation: Adds build complexity without runtime benefits

### Decision 2: Type-Safe Tool IDs

**Choice**: Use string-based const type for tool IDs

**Why**: Compile-time safety without heavy ceremony

**Implementation**:
```go
type ToolID string

const (
    ToolClaudeCode     ToolID = "claude-code"
    ToolCline          ToolID = "cline"
    ToolCostrict       ToolID = "costrict-"
    // ... etc
)
```

**Alternatives Considered**:
- Enums with iota: Breaks string serialization needs
- Plain strings: Current state, no type safety
- Complex type system: Over-engineered for simple need

### Decision 3: Unified Template Rendering

**Choice**: All templates go through TemplateManager, remove inline fmt.Sprintf templates

**Why**: Consistency and easier to test/maintain

**Current Inconsistency**:
- `RenderAgents()` uses TemplateManager with embedded template
- `RenderSpec()` uses inline fmt.Sprintf (templates.go:56-68)
- `RenderProposal()` uses inline fmt.Sprintf (templates.go:71-86)

**After**: All use TemplateManager with .tmpl files

**Alternatives Considered**:
- Keep mixed approach: Inconsistent, harder to maintain
- Move all to fmt.Sprintf: Loses template file benefits (syntax highlighting, reuse)

### Decision 4: Extract Marker Update Utility

**Choice**: Move `UpdateFileWithMarkers` to dedicated `marker_utils.go`

**Why**: Reusable across package, clear responsibility

**Current**: Mixed into configurator.go with all configurator logic
**After**: Standalone utility with focused tests

### Decision 5: Consolidate Constants

**Choice**: Single constants.go with unified naming

**Current Issues**:
- `filePerm` and `filePerms` both exist
- `dirPerm` and `dirPerms` both exist
- Scattered across constants.go and filesystem.go

**After**: Single canonical constant name per concept, all in constants.go

## File Organization

### Before
```
internal/init/
├── configurator.go     (875 lines - configurators + marker logic)
├── registry.go         (146 lines - registry + mapping)
├── executor.go         (510 lines - execution + giant switches)
├── templates.go        (106 lines - mixed template approaches)
├── filesystem.go       (212 lines)
├── models.go           (78 lines)
├── constants.go        (16 lines - incomplete)
├── wizard.go           (564 lines)
```

### After
```
internal/init/
├── tool_definitions.go (250 lines - all tool configs)
├── configurator.go     (120 lines - single generic configurator)
├── registry.go         (100 lines - simple registry)
├── executor.go         (350 lines - no switches, registry lookups)
├── templates.go        (80 lines - unified template rendering)
├── marker_utils.go     (80 lines - extracted marker logic)
├── filesystem.go       (200 lines - minor cleanup)
├── models.go           (90 lines - add ToolConfig struct)
├── constants.go        (50 lines - all constants, including ToolIDs)
├── wizard.go           (564 lines - unchanged)
```

## Migration Strategy

### Phase 1: Add New Code (Non-Breaking)
1. Create `tool_definitions.go` with all tool configs
2. Create `marker_utils.go` with extracted logic
3. Update `constants.go` with ToolID type and all constants
4. Add `GenericConfigurator` alongside existing configurators

### Phase 2: Refactor Internals (Breaking Changes Contained)
1. Update `registry.go` to use new tool definitions
2. Update `executor.go` to use registry lookups instead of switches
3. Update `templates.go` to unified approach
4. Update `configurator.go` - remove old configurators, keep only generic one

### Phase 3: Cleanup
1. Remove old configurator implementations
2. Remove old constants from filesystem.go
3. Update all tests to use new constants
4. Verify all integration tests pass

## Testing Strategy

- Run full test suite after each phase
- All existing tests must pass unchanged (behavior preservation)
- Add new tests for extracted utilities (marker_utils.go)
- Verify tool registration with table-driven tests
- Integration test: init with each tool, verify correct files created

## Risks / Trade-offs

### Risk: Breaking Backward Compatibility
**Mitigation**: All tests pass unchanged; public API (executor.Execute) unchanged

### Risk: Introducing Bugs During Refactor
**Mitigation**: Phased approach; run tests after each file change; preserve exact behavior

### Trade-off: More Indirection
**Cost**: Tool config now in registry instead of explicit code
**Benefit**: 60% less code, easier to add tools, fewer files to touch
**Decision**: Worth it - data-driven is clearer for this use case

### Trade-off: All Templates Must Use TemplateManager
**Cost**: Can't quickly add inline template with fmt.Sprintf
**Benefit**: Consistency, testability, template syntax highlighting
**Decision**: Worth it - only 2 inline templates currently, easy to migrate

## Open Questions

1. ~~Should ToolID be exported or package-private?~~
   **Decision**: Exported - used by registry initialization

2. ~~Keep backward compat for old function names (NewClaudeSlashConfigurator)?~~
   **Decision**: No - internal package only, no public API

3. ~~Move templates to separate templates/ subpackage?~~
   **Decision**: No - current embed pattern works well, no benefit to separate package
