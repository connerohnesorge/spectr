# Change: Redesign Provider Architecture with Composable Initializers

## Why

The current provider system has 15 providers, each implementing a 12-method interface with significant code duplication. Most providers embed `BaseProvider` and only vary in configuration values. This redesign introduces a composable initializer architecture that:

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
- **BREAKING**: Provider metadata (ID, name, priority) moves to `Registration` struct at registration time
- **BREAKING**: **COMPLETELY REMOVE** old `Register(p Provider)` function and all `init()` registration - zero tech debt policy means NO deprecated `Register(_ any)` compatibility shim
- **BREAKING**: Remove all provider `init()` functions that call `Register()`
- **BREAKING**: All markdown markers standardized to `<!-- spectr:start -->` and `<!-- spectr:end -->` (lowercase) for consistency
- **NEW**: `internal/domain` package containing shared domain types (`TemplateRef`, `SlashCommand`, `TemplateContext`) to break import cycles
- **NEW**: `internal/domain` package embeds slash command templates moved from `internal/initialize/templates/tools/`:
  - `slash-proposal.md.tmpl`, `slash-apply.md.tmpl` (Markdown format)
  - `slash-proposal.toml.tmpl`, `slash-apply.toml.tmpl` (TOML format for Gemini)
- **NEW**: `Initializer` interface with `Init(ctx, projectFs, globalFs, cfg, tm)` and `IsSetup(projectFs, globalFs, cfg)` methods
- **NEW**: Built-in initializers with separate types for filesystem and format:
  - `DirectoryInitializer`, `GlobalDirectoryInitializer`
  - `ConfigFileInitializer`
  - `SlashCommandsInitializer` (Markdown), `GlobalSlashCommandsInitializer` (Markdown)
  - `TOMLSlashCommandsInitializer` (TOML for Gemini)
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
| Initializer ordering | Documented guarantee (implicit by type) | Directory → ConfigFile → SlashCommands; simple and predictable |
| Partial failure | Fail-fast, no rollback | Stop on first error, files remain on disk; user fixes and re-runs |
| Registration failure | Partial registrations kept | Successfully registered providers remain; no rollback |
| Template variables | Derive from SpectrDir | SpecsDir = SpectrDir/specs, etc.; single source of truth |
| Global paths | Separate initializer types | GlobalDirectoryInitializer, GlobalSlashCommandsInitializer for ~/.config/tool/ patterns |
| Deduplication | By type + path, then sort | Dedupe first by key, then sort by type priority; execute in order |
| Template selection | TemplateRef directly | ConfigFileInitializer takes TemplateRef, not function; Provider.Initializers() receives TemplateManager |
| TOML support | Separate initializer type | TOMLSlashCommandsInitializer for Gemini; uses .toml.tmpl templates |
| Marker format | Lowercase `<!-- spectr:start/end -->` | ALL markdown markers use lowercase for consistency |
| Template collision | Last-wins precedence | Later template overwrites earlier; no error |
| Filesystem root | os.UserHomeDir() to afero.Fs | Explicit Go stdlib function, converted to afero.Fs |
| Directory creation | Recursive (MkdirAll style) | Create all missing parents automatically |

## Impact

- Affected specs: `support-*` (all 15 provider specs), `cli-interface` (init command)

## Affected Code

- `internal/domain/*.go` - New domain package with shared types (`TemplateRef`, `SlashCommand`, `TemplateContext`)
- `internal/domain/templates/*.tmpl` - Slash command templates:
  - `slash-proposal.md.tmpl`, `slash-apply.md.tmpl` (moved from `internal/initialize/templates/tools/`)
  - `slash-proposal.toml.tmpl`, `slash-apply.toml.tmpl` (new for Gemini TOML format)
- `internal/domain/templates.go` - Embed directive for domain templates
- `internal/initialize/providers/*.go` - Complete rewrite with explicit registration error handling
- `internal/initialize/providers/initializers/*.go` - New initializer implementations
- `internal/initialize/providers/result.go` - New InitResult type
- `internal/initialize/templates/*.go` - Updated to use domain types from `internal/domain`
- `internal/initialize/executor.go` - Simplified provider orchestration with result collection
- `internal/initialize/wizard.go` - Provider selection UI unchanged
- `cmd/init.go` - Minor updates for new API
