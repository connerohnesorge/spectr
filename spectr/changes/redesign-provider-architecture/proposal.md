# Change: Redesign Provider Architecture with Composable Initializers

## Why

The current provider system has 15 providers, each implementing a 12-method interface with significant code duplication. Most providers embed `BaseProvider` and only vary in configuration values. This redesign introduces a composable initializer architecture that:

1. Reduces boilerplate by separating provider identity from initialization logic
2. Enables shared initializers (ConfigFile, SlashCommands, Directory) to be reused and deduped
3. Improves testability by isolating initialization logic into small, focused units
4. Simplifies adding new providers to ~10 lines of registration code

## Scope

**Behavioral equivalence**: Architectural refactor maintaining identical user-facing behavior. No new user-facing features, no new instruction file formats (those belong in a separate proposal). Internal implementation uses new patterns (domain package, initializer interface, dual filesystems) to eliminate boilerplate.

**Zero technical debt policy**: Complete removal of old registration system. No compatibility shims, no deprecated functions that silently swallow calls. Clean break with compile-time errors to force explicit migration.

## What Changes

- **BREAKING**: Remove current `Provider` interface (12 methods) and `BaseProvider` struct
- **BREAKING**: Replace with minimal `Provider` interface returning `[]Initializer`
- **BREAKING**: Provider metadata (ID, name, priority) moves to `Registration` struct at registration time
- **BREAKING**: **COMPLETELY REMOVE** old `Register(p Provider)` function and all `init()` registration - zero tech debt policy means NO deprecated `Register(_ any)` compatibility shim
- **BREAKING**: Remove all provider `init()` functions that call `Register()`
- **CHANGED**: Marker matching is case-insensitive for reading (matches both `<!-- spectr:START -->` and `<!-- spectr:start -->`), always writes lowercase for consistency
- **NEW**: `internal/domain` package containing shared domain types (`TemplateRef`, `SlashCommand`, `TemplateContext`) to break import cycles
- **NEW**: `internal/domain` package embeds slash command templates moved from `internal/initialize/templates/tools/`:
  - `slash-proposal.md.tmpl`, `slash-apply.md.tmpl` (Markdown format)
  - `slash-proposal.toml.tmpl`, `slash-apply.toml.tmpl` (TOML format for Gemini)
- **NEW**: `Initializer` interface with `Init(ctx, projectFs, homeFs, cfg, tm)` and `IsSetup(projectFs, homeFs, cfg)` methods
- **NEW**: Built-in initializers with separate types for filesystem and format:
  - `DirectoryInitializer`, `HomeDirectoryInitializer`
  - `ConfigFileInitializer`
  - `SlashCommandsInitializer` (Markdown), `HomeSlashCommandsInitializer` (Markdown)
  - `PrefixedSlashCommandsInitializer` (Markdown with prefix for Antigravity)
  - `HomePrefixedSlashCommandsInitializer` (Markdown with prefix for Codex, home fs)
  - `TOMLSlashCommandsInitializer` (TOML for Gemini)
- **NEW**: `Config` struct with `SpectrDir` field; other paths derived (SpecsDir = SpectrDir/specs, etc.)
- **NEW**: Two filesystem instances: `projectFs` (project-relative) and `homeFs` (home directory)
- **REMOVED**: `GetFilePaths()`, `HasConfigFile()`, `HasSlashCommands()` methods
- **NEW**: `Initializer.Init()` returns `ExecutionResult` containing created/updated files (explicit change tracking)
- **CHANGED**: Provider registration uses explicit `RegisterAllProviders()` called at startup (no init() in provider files, proper error propagation)
- **MIGRATION**: Users must re-run `spectr init` (clean break, no automatic migration)

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Change detection | ExecutionResult return value | Each initializer returns files it created/updated; explicit and testable |
| Initializer ordering | Documented guarantee (implicit by type) | Directory → ConfigFile → SlashCommands; simple and predictable |
| Partial failure | Fail-fast, no rollback | Stop on first error, files remain on disk; user fixes and re-runs |
| Registration failure | Fail-fast, application exits | Stop on first error; application won't start if any provider fails to register |
| Template variables | Derive from SpectrDir | SpecsDir = SpectrDir/specs, etc.; single source of truth |
| Home paths | Separate initializer types | HomeDirectoryInitializer, HomeSlashCommandsInitializer, HomePrefixedSlashCommandsInitializer for ~/.config/tool/ patterns |
| Deduplication | By type + path, provider priority wins | Sort by type (stable), then dedupe (keep first); higher-priority provider's initializer kept |
| Template selection | TemplateRef directly | ConfigFileInitializer takes TemplateRef, not function; Provider.Initializers() receives TemplateManager |
| TOML support | Separate initializer type | TOMLSlashCommandsInitializer for Gemini; uses .toml.tmpl templates |
| Marker format | Case-insensitive read, lowercase write | Read both uppercase/lowercase for backward compatibility, always write lowercase |
| Template collision | Last-wins precedence | Later template overwrites earlier; no error |
| Filesystem root | os.UserHomeDir() to afero.Fs | Explicit Go stdlib function, converted to afero.Fs |
| Home directory failure | Fail initialization entirely | Home directory access required; initialization aborts if unavailable |
| Directory creation | Recursive (MkdirAll style) | Create all missing parents automatically |
| Slash command update | Always overwrite | Idempotent behavior; user modifications lost on re-init |

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
- `internal/initialize/providers/result.go` - New ExecutionResult type
- `internal/initialize/templates/*.go` - Updated to use domain types from `internal/domain`
- `internal/initialize/executor.go` - Simplified provider orchestration with result collection
- `internal/initialize/wizard.go` - Provider selection UI unchanged
- `cmd/init.go` - Minor updates for new API
