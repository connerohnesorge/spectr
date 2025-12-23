# Change: Redesign Provider Architecture with Composable Initializers

## Why

The current provider system has 17 providers, each implementing a 12-method interface with significant code duplication. Most providers embed `BaseProvider` and only vary in configuration values. This redesign introduces a composable initializer architecture that:

1. Reduces boilerplate by separating provider identity from initialization logic
2. Enables shared initializers (ConfigFile, SlashCommands, Directory) to be reused and deduped
3. Improves testability by isolating initialization logic into small, focused units
4. Simplifies adding new providers to ~10 lines of registration code

## Scope

**Minimal viable**: Focus on reducing boilerplate while keeping behavior identical. No new features, no new instruction files (those belong in a separate proposal).

## What Changes

- **BREAKING**: Remove current `Provider` interface (12 methods) and `BaseProvider` struct
- **BREAKING**: Replace with minimal `Provider` interface returning `[]Initializer`
- **BREAKING**: Provider metadata (ID, name, priority) moves to registration time
- **NEW**: `Initializer` interface with `Init(ctx, fs, cfg, tm)` and `IsSetup(fs, cfg)` methods
- **NEW**: Built-in initializers: `DirectoryInitializer`, `ConfigFileInitializer`, `SlashCommandsInitializer`
- **NEW**: `Config` struct with `SpectrDir` field; other paths derived (SpecsDir = SpectrDir/specs, etc.)
- **NEW**: Two filesystem instances: `projectFs` (project-relative) and `globalFs` (home directory)
- **REMOVED**: `GetFilePaths()`, `HasConfigFile()`, `HasSlashCommands()` methods
- **NEW**: Git-based change detection after initialization (replaces upfront declarations)
- **NEW**: Require git repository - fail early with clear error if not a git repo
- **MIGRATION**: Users must re-run `spectr init` (clean break, no automatic migration)

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Non-git projects | Require git, fail early | Simplifies change detection; git is assumed for spectr |
| Initializer ordering | Implicit, documented guarantee | Directory → ConfigFile → SlashCommands; simple and predictable |
| Partial failure | No rollback, report failures | Keep simple; users can re-run init |
| Template variables | Derive from SpectrDir | SpecsDir = SpectrDir/specs, etc.; single source of truth |
| Global paths | Two fs instances | projectFs and globalFs; supports ~/.config/tool/ patterns |
| Deduplication | By file path | Same path = initialize once; prevents redundant writes |
| Git check timing | Early fail | Check at init start, not lazily; clear error messaging |

## Impact

- Affected specs: `support-*` (all 17 provider specs), `cli-interface` (init command)

## Affected Code

- `internal/initialize/providers/*.go` - Complete rewrite
- `internal/initialize/providers/initializers/*.go` - New initializer implementations
- `internal/initialize/git/detector.go` - New change detection
- `internal/initialize/executor.go` - Simplified provider orchestration
- `internal/initialize/wizard.go` - Provider selection UI unchanged
- `cmd/init.go` - Minor updates for new API
