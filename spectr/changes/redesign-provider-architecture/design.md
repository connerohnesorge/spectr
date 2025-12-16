# Design: Provider Architecture Redesign

## Context

Spectr supports 17 AI CLI/IDE tools (Claude Code, Cursor, Cline, etc.) through a provider system. Each provider configures:
1. An instruction file (e.g., `CLAUDE.md`) with marker-based updates
2. Slash commands (e.g., `.claude/commands/spectr/proposal.md`)

The current implementation has each provider implement a 12-method interface, with most embedding `BaseProvider`. This leads to ~50 lines of boilerplate per provider when the actual variance is just configuration values.

## Goals / Non-Goals

**Goals:**
- Reduce provider authoring to ~10 lines of registration code
- Enable sharing and deduplication of common initialization logic
- Improve testability of initialization steps
- Maintain support for all 17 current providers
- Use `afero.Fs` rooted at project for cleaner path handling

**Non-Goals:**
- Runtime plugin loading (all providers compiled in)
- Backwards compatibility with existing configurations
- Ordered initializer dependencies

## Decisions

### 1. Provider Interface

**Decision**: Providers return a list of initializers; metadata lives at registration.

```go
type Provider interface {
    Initializers(ctx context.Context) []Initializer
}

type Initializer interface {
    Init(ctx context.Context, fs afero.Fs, cfg *Config) error
    IsSetup(fs afero.Fs, cfg *Config) bool
}

type Config struct {
    SpectrDir string // e.g., "spectr" (relative to fs root)
}
```

**Alternatives considered:**
- Keep metadata in Provider interface (current design) - More boilerplate
- Use functional options pattern - Harder to test

### 2. Registration API

**Decision**: Instance-only registry with metadata at registration time. No global state.

```go
// Create a registry instance
reg := providers.NewRegistry()

// Register a provider with its metadata
reg.Register(providers.Registration{
    ID:       "claude-code",
    Name:     "Claude Code",
    Priority: 1,
    Provider: &ClaudeProvider{},
})

// Get all providers sorted by priority
all := reg.All()

// Get provider by ID
claude := reg.Get("claude-code")
```

**Rationale**:
- Providers don't need to know their own ID/name/priority. This is registry concern.
- Instance-based registry improves testability (no shared global state between tests).
- No `init()` magic - explicit registration in application setup.

**Removed**: Global `Register()`, `Get()`, `All()`, `IDs()`, `Count()`, `WithConfigFile()`, `WithSlashCommands()`, `Reset()` functions.

### 3. Built-in Initializers

**Decision**: Provide three composable initializers:

```go
// Creates directories (e.g., .claude/commands/spectr/)
func NewDirectoryInitializer(paths ...string) Initializer

// Creates/updates instruction file with markers
func NewConfigFileInitializer(path string, template string) Initializer

// Creates slash commands from templates
func NewSlashCommandsInitializer(dir, ext string, format CommandFormat) Initializer
```

**Deduplication**: When multiple providers share an initializer with same config, run once.

### 4. Filesystem Abstraction

**Decision**: Use `afero.NewBasePathFs(osFs, projectPath)` so all paths are project-relative.

```go
// Instead of: "/home/user/project/CLAUDE.md"
// Use:        "CLAUDE.md"
fs := afero.NewBasePathFs(afero.NewOsFs(), projectPath)
```

**Rationale**: Cleaner paths, easier testing with `afero.MemMapFs`.

### 5. File Change Detection

**Decision**: Use git diff after initialization instead of upfront declarations.

```go
// Before init
beforeCommit := git.Stash()

// Run initializers
for _, init := range allInitializers {
    init.Init(ctx, fs, cfg)
}

// After init
changedFiles := git.DiffFiles(beforeCommit)
```

**Rationale**: Simpler provider interface; git already tracks file changes.

## Example: Claude Code Provider

```go
type ClaudeProvider struct{}

func (p *ClaudeProvider) Initializers(ctx context.Context) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".claude/commands/spectr"),
        NewConfigFileInitializer("CLAUDE.md", InstructionTemplate),
        NewSlashCommandsInitializer(".claude/commands/spectr", ".md", FormatMarkdown),
    }
}

func init() {
    providers.Register(providers.Registration{
        ID:       "claude-code",
        Name:     "Claude Code",
        Priority: 1,
        Provider: &ClaudeProvider{},
    })
}
```

## Example: Gemini Provider (TOML format)

```go
type GeminiProvider struct{}

func (p *GeminiProvider) Initializers(ctx context.Context) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".gemini/commands/spectr"),
        // No config file for Gemini
        NewSlashCommandsInitializer(".gemini/commands/spectr", ".toml", FormatTOML),
    }
}
```

## Initializer Deduplication

When initializers are collected from multiple providers:

```go
func dedupeInitializers(all []Initializer) []Initializer {
    seen := make(map[string]bool)
    var result []Initializer
    for _, init := range all {
        key := initializerKey(init) // Based on type + config
        if !seen[key] {
            seen[key] = true
            result = append(result, init)
        }
    }
    return result
}
```

### 6. Shared Helper Functions

**Decision**: Keep shared helpers, migrate to use `afero.Fs`.

```go
// helpers.go - Updated signatures
func FileExists(fs afero.Fs, path string) bool
func EnsureDir(fs afero.Fs, path string) error
func UpdateFileWithMarkers(fs afero.Fs, path, content, start, end string) error
```

**Rationale**:
- Avoids code duplication across initializers
- `afero.Fs` abstraction enables testing with `afero.MemMapFs`
- Marker-based updates are complex enough to warrant shared implementation

**Removed**: `expandPath()`, `isGlobalPath()` - No longer needed with project-relative paths.

### 7. Error Handling

**Decision**: No rollback on partial failure.

- If 2 of 5 initializers succeed and the 3rd fails, the partial state remains
- Users can re-run `spectr init` which is idempotent
- Git provides natural rollback via `git checkout` if needed

**Rationale**: Simpler implementation, git already provides recovery mechanism.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking change for users | Clear migration docs; `spectr init` re-run required |
| Loss of file path metadata | Git diff after init shows all changes |
| Initializer key collisions | Use deterministic key generation (type + sorted config) |
| Partial failure state | Initializers are idempotent; re-run is safe |

## Migration Plan

1. Implement new provider system alongside existing
2. Migrate providers one-by-one to new system
3. Remove old provider code:
   - Delete old `Provider` interface and `BaseProvider` from `provider.go`
   - Delete `TemplateRenderer` interface
   - Remove global registry functions from `registry.go`
   - Migrate `helpers.go` to use `afero.Fs`
   - Remove unused constants from `constants.go`
4. Update docs to explain re-initialization requirement
5. No automatic migration - clean break

## Resolved Questions

- **Dry-run**: Not needed. Git diff after initialization provides sufficient visibility.
- **Rollback**: Not needed. Partial state is acceptable; re-run is idempotent.
- **Helper functions**: Keep shared helpers, migrate to use `afero.Fs`.
- **Registry**: Instance-only, no global state.
