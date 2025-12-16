# Change: Redesign Provider Architecture with Composable Initializers

## Why

The current provider system has 17 providers, each implementing a 12-method interface with significant code duplication. Most providers embed `BaseProvider` and only vary in configuration values. This redesign introduces a composable initializer architecture that:

1. Reduces boilerplate by separating provider identity from initialization logic
2. Enables shared initializers (ConfigFile, SlashCommands, Directory) to be reused and deduped
3. Improves testability by isolating initialization logic into small, focused units
4. Simplifies adding new providers to ~10 lines of registration code

## What Changes

- **BREAKING**: Remove current `Provider` interface (12 methods) and `BaseProvider` struct
- **BREAKING**: Replace with minimal `Provider` interface returning `[]Initializer`
- **BREAKING**: Provider metadata (ID, name, priority) moves to registration time
- **NEW**: `Initializer` interface with `Init()` and `IsSetup()` methods
- **NEW**: Built-in initializers: `ConfigFileInitializer`, `SlashCommandsInitializer`, `DirectoryInitializer`
- **NEW**: `Config` struct with `SpectrDir` field only
- **NEW**: Use `afero.NewBasePathFs(osFs, projectPath)` so all paths are relative to project root
- **REMOVED**: `GetFilePaths()`, `HasConfigFile()`, `HasSlashCommands()` methods
- **NEW**: Use git diff after initialization to detect changed files (no upfront declarations)
- **NEW**: Add missing instruction file support for providers:
  - Gemini → `GEMINI.md`
  - Cursor → `.cursorrules`
  - Aider → `AIDER-SPECTR.md`
  - OpenCode → `AGENTS.md`
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
