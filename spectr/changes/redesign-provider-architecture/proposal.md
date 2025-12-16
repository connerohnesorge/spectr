# Change: Redesign Provider Architecture with Composable Initializers

## Why

The current provider system has 17 providers, each implementing a 12-method interface with significant code duplication. Most providers embed `BaseProvider` and only vary in configuration values. This redesign introduces a composable initializer architecture that:

1. Reduces boilerplate by separating provider identity from initialization logic
2. Enables shared initializers (ConfigFile, SlashCommands, Directory) to be reused and deduped
3. Improves testability by isolating initialization logic into small, focused units
4. Simplifies adding new providers to ~10 lines of registration code

## What Changes

### New Architecture
- **NEW**: Minimal `Provider` interface returning `[]Initializer`
- **NEW**: `Initializer` interface with `Init()` and `IsSetup()` methods
- **NEW**: Built-in initializers: `ConfigFileInitializer`, `SlashCommandsInitializer`, `DirectoryInitializer`
- **NEW**: `Config` struct with `SpectrDir` field only
- **NEW**: Use `afero.NewBasePathFs(osFs, projectPath)` so all paths are relative to project root
- **NEW**: Instance-only `Registry` struct (no global state) for testability
- **NEW**: Shared helper functions migrated to use `afero.Fs`
- **NEW**: Add missing instruction file support for providers:
  - Gemini → `GEMINI.md`
  - Cursor → `.cursorrules`
  - Aider → `AIDER-SPECTR.md`
  - OpenCode → `AGENTS.md`

### Breaking Changes
- **BREAKING**: Remove current `Provider` interface (12 methods)
- **BREAKING**: Remove `BaseProvider` struct and all its methods
- **BREAKING**: Remove `TemplateRenderer` interface
- **BREAKING**: Provider metadata (ID, name, priority) moves to registration time
- **BREAKING**: Remove global registry functions (`Register`, `Get`, `All`, `IDs`, `Count`, `WithConfigFile`, `WithSlashCommands`, `Reset`)

### Removed Code (Files to Delete/Modify)
- **DELETE**: `provider.go` - Old `Provider` interface (lines 107-158), `TemplateRenderer` interface (lines 164-179), `BaseProvider` struct and all methods (lines 183-479)
- **MODIFY**: `registry.go` - Remove global `registry` variable and global functions; keep only instance-based `Registry` struct
- **MODIFY**: `helpers.go` - Migrate to use `afero.Fs` instead of `os` package
- **MODIFY**: `constants.go` - Remove `StandardCommandPaths()`, `StandardFrontmatter()`, priority constants used only by old system
- **MODIFY**: All 17 provider files (`claude.go`, `gemini.go`, etc.) - Complete rewrite to new interface

### Design Decisions
- **NO DRY-RUN**: Git diff after initialization provides sufficient visibility
- **NO ROLLBACK**: Partial failures leave partial state; users can re-run or fix manually
- **MIGRATION**: Users must re-run `spectr init` (clean break, no automatic migration)

## Impact

- Affected specs: `support-*` (all 17 provider specs), `cli-interface` (init command)

## Note

Individual `support-*` delta specs are blocked by a validator limitation (same-named requirements across capabilities). See `fix-cross-capability-requirement-validation` proposal. Once that fix is implemented, the provider deltas will be added.

## Affected Code

- `internal/initialize/providers/*.go` - Complete rewrite
- `internal/initialize/executor.go` - Simplified provider configuration
- `internal/initialize/wizard.go` - Provider selection UI unchanged
- `cmd/init.go` - Minor updates for new API
