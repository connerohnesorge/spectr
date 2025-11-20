# Design: Provider Composition Architecture

## Context

The Spectr init system currently has 7 memory file-based tools (Claude Code, Cline, Qoder, CodeBuddy, Qwen, Costrict, Antigravity) and 16+ slash command-based tools. The existing `Configurator` interface mixes these concerns, and a partial refactor to `ToolProvider` pattern is in progress with new provider interfaces alongside the legacy interface.

This change completes the migration by:
1. Creating distinct provider types for each concern
2. Composing them into unified tool providers
3. Eliminating the legacy Configurator interface

## Goals

- **Clear separation**: Memory files vs slash commands as distinct interfaces
- **Composability**: Tools can mix and match providers via embedding
- **Consistency**: All memory file tools update `spectr/AGENTS.md` with Spectr instructions
- **Simplicity**: Remove legacy technical debt (Configurator interface)

## Non-Goals

- Changing file paths or content of generated files (backward compatible)
- Modifying template system or rendering logic
- Changing registry tool definitions or IDs
- Backward compatibility with Configurator interface (breaking change approved)

## Decisions

### Decision 1: Keep AgentsFileProvider Specific

**Rationale**: AgentsFileProvider is currently used by Antigravity to manage the repo root AGENTS.md file. We keep it focused on this single responsibility rather than making it generic.

**Why not generic?** Each memory file (CLAUDE.md, CLINE.md, etc.) may need tool-specific logic in the future. Separate providers maintain flexibility.

**Alternatives considered**:
- Generic configurable provider accepting filename - Rejected: Less type-safe, harder to customize per-tool
- Rename to GenericMemoryFileProvider - Rejected: Implies multi-purpose usage we don't need

### Decision 2: Create SpectrAgentsUpdater Provider

**Purpose**: Ensures all memory file-based tools update `spectr/AGENTS.md` (in the spectr/ directory) with generic Spectr usage instructions.

**Why separate provider?** This is a cross-cutting concern that applies to all memory file tools. By making it a reusable provider, we avoid code duplication and ensure consistency.

**Implementation**:
- Implements `MemoryFileProvider` interface
- Updates `spectr/AGENTS.md` using marker-based updates (same as other memory files)
- Template: `internal/init/templates/spectr/AGENTS.md.tmpl` (already exists)

### Decision 3: Use Embedded Fields for Composition

**Chosen pattern**:
```go
type ClaudeCodeToolProvider struct {
    *ClaudeMemoryFileProvider      // CLAUDE.md in repo root
    *ClaudeSlashCommandProvider    // .claude/commands/spectr/
    *SpectrAgentsUpdater           // spectr/AGENTS.md
}
```

**Benefits**:
- Automatic interface implementation (no wrapper methods needed)
- Clean composition syntax
- Go idiomatic pattern for "has-a" relationships

**Alternatives considered**:
- Explicit named fields with getters - Rejected: More boilerplate, less idiomatic
- Interface-based composition - Rejected: Runtime overhead, less type-safe

### Decision 4: Naming Convention for Providers

**Memory File Providers**: `{Tool}MemoryFileProvider`
- ClaudeMemoryFileProvider, ClineMemoryFileProvider, etc.

**Slash Command Providers**: `{Tool}SlashCommandProvider`
- ClaudeSlashCommandProvider, ClineSlashCommandProvider, etc.

**Composite Providers**: `{Tool}ToolProvider`
- ClaudeCodeToolProvider, ClineToolProvider, etc.

**Rationale**: Clear, consistent, and indicates the type of provider at a glance.

### Decision 5: Complete Migration (No Backward Compatibility)

**Decision**: Remove `Configurator` interface entirely in this change.

**Why now?** The partial refactor already created new interfaces. Maintaining both patterns creates confusion and tech debt. User confirmed backward compatibility is not required.

**Migration strategy**:
1. Create all new providers
2. Update registry to use ToolProvider exclusively
3. Update executor to use ToolProvider exclusively
4. Delete legacy Configurator interface and implementations
5. All in single atomic change to avoid inconsistent state

## Provider Interface Hierarchy

```
ToolProvider (interface)
├── GetName() string
├── GetMemoryFileProvider() MemoryFileProvider    // Can be nil
└── GetSlashCommandProvider() SlashCommandProvider // Can be nil

MemoryFileProvider (interface)
├── ConfigureMemoryFile(projectPath string) error
└── IsMemoryFileConfigured(projectPath string) bool

SlashCommandProvider (interface)
├── ConfigureSlashCommands(projectPath string) error
└── AreSlashCommandsConfigured(projectPath string) bool
```

## Implementation Structure

### Memory File Providers (7 new + 1 existing)

Each manages a specific memory file in repo root:
- `ClaudeMemoryFileProvider` → `CLAUDE.md`
- `ClineMemoryFileProvider` → `CLINE.md`
- `QoderMemoryFileProvider` → `QODER.md`
- `CodeBuddyMemoryFileProvider` → `CODEBUDDY.md`
- `QwenMemoryFileProvider` → `QWEN.md`
- `CostrictMemoryFileProvider` → `COSTRICT.md`
- `AgentsFileProvider` → `AGENTS.md` (existing, for Antigravity)

### SpectrAgentsUpdater Provider

Special provider used by all memory file tools:
- Updates `spectr/AGENTS.md` with generic Spectr instructions
- Uses same marker-based update pattern
- Composed into all memory file tool providers

### Slash Command Providers (16+ existing)

Rename existing implementations to follow pattern:
- `ClaudeSlashCommandProvider`
- `ClineSlashCommandProvider`
- `CursorSlashCommandProvider`
- etc.

### Composite Tool Providers (7 new)

Combine providers for each tool:

```go
type ClaudeCodeToolProvider struct {
    *ClaudeMemoryFileProvider
    *ClaudeSlashCommandProvider
    *SpectrAgentsUpdater
}

type ClineToolProvider struct {
    *ClineMemoryFileProvider
    *ClineSlashCommandProvider
    *SpectrAgentsUpdater
}

// ... etc for all 7 memory file tools
```

## Configuration Flow

### Before (Legacy Configurator)
```
Executor → getConfigurator(toolID) → Configurator.Configure()
                                    ↓
                            Single monolithic method handles everything
```

### After (Composite Providers)
```
Executor → getToolProvider(toolID) → ToolProvider
                                    ↓
                    Executor calls each provider in sequence:
                    1. MemoryFileProvider.ConfigureMemoryFile()
                    2. SlashCommandProvider.ConfigureSlashCommands()
                    3. SpectrAgentsUpdater.ConfigureMemoryFile()
```

## File Organization

- `internal/init/providers.go` - Base providers + memory file providers + SpectrAgentsUpdater
- `internal/init/slash_providers.go` - All slash command provider implementations (existing)
- `internal/init/composite_providers.go` - NEW: Composite tool providers
- `internal/init/interfaces.go` - Provider interfaces (Configurator removed)
- `internal/init/registry.go` - Updated to use ToolProvider
- `internal/init/executor.go` - Updated to use ToolProvider

## Risks / Trade-offs

### Risk: Breaking Existing Integrations
**Mitigation**: File paths and content remain identical. Only internal implementation changes.

### Trade-off: More Types
**Benefit**: Type safety and clarity
**Cost**: More structs to maintain (~20 provider types)
**Verdict**: Worth it for improved separation of concerns

### Risk: Embedded Field Confusion
**Mitigation**: Clear naming conventions and documentation. Embedded fields are idiomatic Go.

## Migration Plan

### Phase 1: Create New Providers
1. Add 7 memory file providers to `providers.go`
2. Add SpectrAgentsUpdater to `providers.go`
3. Rename slash command providers in `slash_providers.go`
4. Create composite providers in new `composite_providers.go`

### Phase 2: Update Infrastructure
1. Update registry to use ToolProvider pattern
2. Update executor `getToolProvider()` method
3. Update executor configuration flow to call providers in sequence

### Phase 3: Remove Legacy Code
1. Delete legacy Configurator interface from `interfaces.go`
2. Delete all configurator implementations from `configurator.go`
3. Delete or rename `configurator.go` if empty

### Phase 4: Validation
1. Run all tests
2. Test init command with each tool
3. Verify file generation is identical to before

## Open Questions

None - all design questions resolved during planning phase.
