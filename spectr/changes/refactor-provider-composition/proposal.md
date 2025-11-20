# Change: Refactor Provider Composition Architecture

## Why

The current configurator system mixes two distinct concerns: memory files (like CLAUDE.md, AGENTS.md) that are included in every agent invocation, and slash commands (like .claude/commands/spectr/*.md) that are invoked conditionally. This makes it difficult to understand which tools do what, and prevents flexible composition of providers.

Additionally, there's no systematic way to ensure that memory file-based tools update `spectr/AGENTS.md` with Spectr usage instructions, leading to inconsistent tool integration.

## What Changes

- Create specific memory file providers: `ClaudeMemoryFileProvider`, `ClineMemoryFileProvider`, `QoderMemoryFileProvider`, `CodeBuddyMemoryFileProvider`, `QwenMemoryFileProvider`, `CostrictMemoryFileProvider`
- Keep `AgentsFileProvider` focused solely on managing AGENTS.md for Antigravity
- Create `SpectrAgentsUpdater` provider to ensure all memory file tools update `spectr/AGENTS.md` with Spectr instructions
- Rename slash command providers to follow consistent pattern: `ClaudeSlashCommandProvider`, `ClineSlashCommandProvider`, etc.
- Create composite tool providers using embedded fields composition pattern
- **BREAKING**: Remove legacy `Configurator` interface entirely
- Migrate registry and executor to use `ToolProvider` interface exclusively
- Update all 7 memory file configurators to new pattern
- Update all 16+ slash command configurators to new pattern

## Impact

### Affected Specs
- `cli-interface` - Modified "Automatic Slash Command Installation" requirement to use ToolProvider composition pattern
- `cli-interface` - Added requirements for provider composition architecture, SpectrAgentsUpdater, provider interfaces, and embedded field pattern

### Affected Code
- `internal/init/providers.go` - New memory file provider implementations (~300 lines added)
- `internal/init/configurator.go` - Remove all legacy configurator implementations (~875 lines removed)
- `internal/init/composite_providers.go` - New file for composite tool providers (~200 lines added)
- `internal/init/interfaces.go` - Remove legacy Configurator interface (~10 lines removed)
- `internal/init/registry.go` - Update to use ToolProvider pattern (~50 lines modified)
- `internal/init/executor.go` - Update to use ToolProvider pattern (~100 lines modified)

### Migration Path
All existing tool configurations will automatically migrate to the new provider pattern. No user action required as the file paths and content remain identical.
