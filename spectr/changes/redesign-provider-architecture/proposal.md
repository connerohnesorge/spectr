# Change: Redesign Provider Architecture with Composable Initializers

## Why

The current provider system has 17 providers, each implementing a 12-method interface with significant code duplication. Most providers embed `BaseProvider` and only vary in configuration values. This redesign introduces a composable initializer architecture that:

1. Reduces boilerplate by separating provider identity from initialization logic
2. Enables shared initializers (ConfigFile, SlashCommands, Directory) to be reused and deduped
3. Improves testability by isolating initialization logic into small, focused units
4. Simplifies adding new providers to ~10 lines of registration code

## Scope

**Minimal viable**: Focus on reducing boilerplate while keeping behavior identical. No new features, no new instruction files (those belong in a separate proposal).

**Zero technical debt policy**: Complete removal of old registration system. No compatibility shims, no deprecated functions that silently swallow calls. Clean break with compile-time errors to force explicit migration.

## What Changes

- **BREAKING**: Remove current `Provider` interface (12 methods) and `BaseProvider` struct
- **BREAKING**: Replace with minimal `Provider` interface returning `[]Initializer`
- **BREAKING**: Provider metadata (ID, name, priority) moves to registration time
- **BREAKING**: **COMPLETELY REMOVE** old `Register(p Provider)` function and all `init()` registration - zero tech debt policy means NO deprecated `Register(_ any)` compatibility shim
- **BREAKING**: Remove all provider `init()` functions that call `Register()`
- **NEW**: `internal/domain` package containing shared domain types (`TemplateRef`, `SlashCommand`, `TemplateContext`) to break import cycles
- **NEW**: `internal/domain` package embeds slash command templates (`slash-proposal.md.tmpl`, `slash-apply.md.tmpl`) moved from `internal/initialize/templates/tools/`
- **NEW**: `Initializer` interface with `Init(ctx, fs, cfg, tm)` and `IsSetup(fs, cfg)` methods
- **NEW**: Built-in initializers: `DirectoryInitializer`, `ConfigFileInitializer`, `SlashCommandsInitializer`
- **NEW**: `Config` struct with `SpectrDir` field; other paths derived (SpecsDir = SpectrDir/specs, etc.)
- **NEW**: Two filesystem instances: `projectFs` (project-relative) and `globalFs` (home directory)
- **REMOVED**: `GetFilePaths()`, `HasConfigFile()`, `HasSlashCommands()` methods
- **NEW**: `Initializer.Init()` returns `InitResult` containing created/updated files (explicit change tracking)
- **CHANGED**: Provider registration uses explicit `RegisterAllProviders()` called at startup (no init() in provider files, proper error propagation)
- **MIGRATION**: Users must re-run `spectr init` (clean break, no automatic migration)

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Change detection | InitResult return value | Each initializer returns files it created/updated; explicit and testable |
| Initializer ordering | Implicit, documented guarantee | Directory → ConfigFile → SlashCommands; simple and predictable |
| Partial failure | No rollback, report failures | Keep simple; users can re-run init |
| Template variables | Derive from SpectrDir | SpecsDir = SpectrDir/specs, etc.; single source of truth |
| Global paths | Two fs instances | projectFs and globalFs; supports ~/.config/tool/ patterns |
| Deduplication | By file path | Same path = initialize once; prevents redundant writes |

## Impact

- Affected specs: `support-*` (all 17 provider specs), `cli-interface` (init command)

## Affected Code

- `internal/domain/*.go` - New domain package with shared types (`TemplateRef`, `SlashCommand`, `TemplateContext`)
- `internal/domain/templates/*.tmpl` - Slash command templates moved from `internal/initialize/templates/tools/`
- `internal/domain/templates.go` - Embed directive for domain templates
- `internal/initialize/providers/*.go` - Complete rewrite with explicit registration error handling
- `internal/initialize/providers/initializers/*.go` - New initializer implementations
- `internal/initialize/providers/result.go` - New InitResult type
- `internal/initialize/templates/*.go` - Updated to use domain types from `internal/domain`
- `internal/initialize/executor.go` - Simplified provider orchestration with result collection
- `internal/initialize/wizard.go` - Provider selection UI unchanged
- `cmd/init.go` - Minor updates for new API
