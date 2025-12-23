# Design: Provider Architecture Redesign

## Context

Spectr supports 17 AI CLI/IDE tools (Claude Code, Cursor, Cline, etc.) through a provider system. Each provider configures:
1. An instruction file (e.g., `CLAUDE.md`) with marker-based updates
2. Slash commands (e.g., `.claude/commands/spectr/proposal.md`)

The current implementation has each provider implement a 12-method interface, with most embedding `BaseProvider`. This leads to ~50 lines of boilerplate per provider when the actual variance is just configuration values.

## Scope

**Minimal viable refactor**: Reduce boilerplate while keeping behavior identical. No new features.

## Goals / Non-Goals

**Goals:**
- Reduce provider authoring to ~10 lines of registration code
- Enable sharing and deduplication of common initialization logic
- Improve testability of initialization steps
- Maintain support for all 17 current providers
- Use `afero.Fs` rooted at project for cleaner path handling
- Explicit change tracking via InitResult return values

**Non-Goals:**
- Runtime plugin loading (all providers compiled in)
- Backwards compatibility with existing configurations
- Rollback on partial failure
- New instruction file formats (separate proposal)

## Decisions

### 1. Provider Interface

**Decision**: Providers return a list of initializers; metadata lives at registration.

```go
type Provider interface {
    Initializers(ctx context.Context) []Initializer
}

// InitResult contains the files created or modified by an initializer.
type InitResult struct {
    CreatedFiles []string
    UpdatedFiles []string
}

type Initializer interface {
    // Init creates or updates files. Returns result with file changes and error if initialization fails.
    // Must be idempotent (safe to run multiple times).
    Init(ctx context.Context, fs afero.Fs, cfg *Config, tm *TemplateManager) (InitResult, error)

    // IsSetup returns true if this initializer's artifacts already exist.
    IsSetup(fs afero.Fs, cfg *Config) bool

    // Path returns the file/directory path this initializer manages.
    // Used for deduplication: same path = run once.
    Path() string

    // IsGlobal returns true if this initializer uses globalFs instead of projectFs.
    IsGlobal() bool
}

type Config struct {
    SpectrDir string // e.g., "spectr" (relative to fs root)
}

// Derived paths (computed, not stored):
// - SpecsDir:    cfg.SpectrDir + "/specs"
// - ChangesDir:  cfg.SpectrDir + "/changes"
// - ProjectFile: cfg.SpectrDir + "/project.md"
// - AgentsFile:  cfg.SpectrDir + "/AGENTS.md"

func (c *Config) SpecsDir() string    { return c.SpectrDir + "/specs" }
func (c *Config) ChangesDir() string  { return c.SpectrDir + "/changes" }
func (c *Config) ProjectFile() string { return c.SpectrDir + "/project.md" }
func (c *Config) AgentsFile() string  { return c.SpectrDir + "/AGENTS.md" }
```

**Alternatives considered:**
- Keep metadata in Provider interface (current design) - More boilerplate
- Use functional options pattern - Harder to test
- Store all paths in Config - Redundant, error-prone

### 2. Registration API

**Decision**: Metadata provided at registration time.

```go
// Register a provider with its metadata
providers.Register(providers.Registration{
    ID:       "claude-code",
    Name:     "Claude Code",
    Priority: 1,
    Provider: &ClaudeProvider{},
})
```

**Rationale**: Providers don't need to know their own ID/name/priority. This is registry concern.

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

**Decision**: Use two filesystem instances to support both project-relative and global paths.

```go
// Project-relative filesystem (for CLAUDE.md, .claude/commands/, etc.)
projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectPath)

// Global filesystem (for ~/.config/tool/commands/, etc.)
globalFs := afero.NewBasePathFs(afero.NewOsFs(), os.UserHomeDir())

// Executor provides both to initializers based on IsGlobal()
type ExecutorContext struct {
    ProjectFs afero.Fs
    GlobalFs  afero.Fs
    Config    *Config
    Templates *TemplateManager
}
```

**Rationale**:
- Project paths are cleaner and easier to test
- Global paths support tools like Aider that use ~/.config/
- Initializers declare via `IsGlobal()` which fs to use

### 5. File Change Detection

**Decision**: Each initializer returns an `InitResult` containing the files it created/updated.

```go
// Collect results from all initializers
var allResults []InitResult

for _, init := range allInitializers {
    result, err := init.Init(ctx, fs, cfg, tm)
    if err != nil {
        errors = append(errors, err)
        continue
    }
    allResults = append(allResults, result)
}

// Aggregate into ExecutionResult
executionResult := aggregateResults(allResults)
```

**Rationale**: Explicit change tracking; initializers know what they create; works in non-git projects; more testable.

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

### 8. Initializer Ordering (Documented Guarantee)

**Decision**: Initializers are sorted by type before execution. This is a documented guarantee.

```go
// Execution order (guaranteed):
// 1. DirectoryInitializer   - Create directories first
// 2. ConfigFileInitializer  - Then config files (may need directories)
// 3. SlashCommandsInitializer - Then slash commands (may need directories)

func sortInitializers(all []Initializer) []Initializer {
    sort.SliceStable(all, func(i, j int) bool {
        return initializerPriority(all[i]) < initializerPriority(all[j])
    })
    return all
}

func initializerPriority(init Initializer) int {
    switch init.(type) {
    case *DirectoryInitializer:
        return 1
    case *ConfigFileInitializer:
        return 2
    case *SlashCommandsInitializer:
        return 3
    default:
        return 99
    }
}
```

**Rationale**: Directories must exist before files can be written. This ordering is implicit but guaranteed.

### 9. Initializer Deduplication

**Decision**: Deduplicate by file path. Same path = run once.

When initializers are collected from multiple providers:

```go
func dedupeInitializers(all []Initializer) []Initializer {
    seen := make(map[string]bool)
    var result []Initializer
    for _, init := range all {
        key := init.Path() // Simple: just the path
        if !seen[key] {
            seen[key] = true
            result = append(result, init)
        }
    }
    return result
}
```

**Example**: If Claude Code and Cline both return `ConfigFileInitializer{path: "CLAUDE.md"}`, only one runs.

**Rationale**: Path-based deduplication is simple and covers the common case (multiple providers sharing same file).

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking change for users | Clear migration docs; `spectr init` re-run required |
| Path-based dedup misses edge cases | Simple first; can enhance if needed |
| No rollback on failure | Clear error reporting; users can re-run init |

## Migration Plan

1. Implement new provider system alongside existing
2. Migrate providers one-by-one to new system
3. Remove old provider code
4. Update docs to explain re-initialization requirement
5. No rollback needed - old configs continue to work

### 6. TemplateManager Integration

**Decision**: Initializers receive `*TemplateManager` instead of implementing `TemplateRenderer`.

```go
type Initializer interface {
    Init(ctx context.Context, fs afero.Fs, cfg *Config, tm *TemplateManager) (InitResult, error)
    IsSetup(fs afero.Fs, cfg *Config) bool
}
```

**Rationale**:
- Reuses existing `TemplateManager` from `internal/initialize/templates.go`
- Avoids duplicating template rendering logic in each initializer
- Simpler than the old `TemplateRenderer` interface pattern

**Alternatives considered:**
- Each initializer implements own template rendering - More code duplication
- Pass templates as strings - Less flexible, harder to maintain

## Resolved Questions

| Question | Decision |
|----------|----------|
| Change detection approach? | InitResult return value from each Initializer |
| Initializer ordering? | Implicit ordering by type (Directory → ConfigFile → SlashCommands) |
| Rollback on partial failure? | No rollback; report failures, users re-run init |
| Template variable location? | Derived from SpectrDir via methods |
| Global paths support? | Two fs instances (projectFs, globalFs) |
| Deduplication key? | By file path (simple and effective) |

## Future Considerations (Out of Scope)

- `spectr init --dry-run` to preview changes without applying
- New instruction file support for Gemini, Cursor, Aider, OpenCode (separate proposal)
