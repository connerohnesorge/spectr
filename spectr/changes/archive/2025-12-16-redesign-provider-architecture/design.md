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

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking change for users | Clear migration docs; `spectr init` re-run required |
| Loss of file path metadata | Use git diff; add `spectr init --dry-run` if needed |
| Initializer key collisions | Use deterministic key generation (type + sorted config) |

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
    Init(ctx context.Context, fs afero.Fs, cfg *Config, tm *TemplateManager) error
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

### 7. Git-Based Change Detection

**Decision**: Use git diff after initialization instead of upfront `GetFilePaths()` declarations.

```go
// internal/initialize/git/detector.go
type ChangeDetector struct {
    repoPath string
}

func (d *ChangeDetector) Snapshot() (string, error) {
    // Create git stash or record HEAD state
}

func (d *ChangeDetector) ChangedFiles(before string) ([]string, error) {
    // git diff --name-only before...HEAD
}
```

**Implementation details:**
- Use `git stash create` to capture pre-init state without modifying working tree
- Use `git diff --name-only` to detect changed files after init
- Handle edge cases: untracked files (use `git status --porcelain`), not a git repo, dirty working tree

**Rationale**:
- Eliminates need for providers to declare file paths upfront
- Git is source of truth for file changes anyway
- Simpler provider interface (no `GetFilePaths()` method)

**Trade-offs:**
- Requires git repo (graceful degradation for non-git projects)
- Slightly more complex executor flow

## Open Questions

- Should `spectr init --dry-run` be added to preview changes without applying?
- Should initializers support rollback on partial failure?
