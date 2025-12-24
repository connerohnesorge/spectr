# Design: Provider Architecture Redesign

## Context

Spectr supports 17 AI CLI/IDE tools (Claude Code, Cursor, Cline, etc.) through a provider system. Each provider configures:
1. An instruction file (e.g., `CLAUDE.md`) with marker-based updates
2. Slash commands (e.g., `.claude/commands/spectr/proposal.md`)

The current implementation has each provider implement a 12-method interface, with most embedding `BaseProvider`. This leads to ~50 lines of boilerplate per provider when the actual variance is just configuration values.

**Current Problems:**
1. **Import cycle**: `providers.TemplateManager` cannot import `internal/initialize/templates` without creating an import cycle, forcing the use of `any` as placeholder types
2. **Silent registration failures**: Provider registration in `init()` assigns errors to the blank identifier, silently discarding failures

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
- Break import cycles with a dedicated `internal/domain` package
- Fail-fast registration with explicit error handling (no silent failures)

**Non-Goals:**
- Runtime plugin loading (all providers compiled in)
- Backwards compatibility with existing configurations
- Rollback on partial failure
- New instruction file formats (separate proposal)

## Decisions

### 0. Domain Package for Shared Types

**Decision**: Create `internal/domain` package containing shared domain types to break import cycles.

**Problem**: The current architecture has an import cycle issue:
- `providers.TemplateManager` needs to reference template types like `TemplateRef` and `SlashCommand`
- These types are defined in `internal/initialize/templates`
- But `templates` cannot import `providers` and vice versa without creating a cycle
- Currently, `any` is used as a placeholder, with the concrete adapter in `executor.go` explaining the real types

**Solution**: Extract domain objects into `internal/domain`:

```go
// internal/domain/template.go
package domain

import (
    "bytes"
    "html/template"
)

// TemplateRef is a type-safe reference to a parsed template.
// It can be safely passed between packages without creating import cycles.
type TemplateRef struct {
    Name     string              // template file name (e.g., "instruction-pointer.md.tmpl")
    Template *template.Template  // pre-parsed template
}

// Render executes the template with the given context.
func (tr TemplateRef) Render(ctx TemplateContext) (string, error) {
    var buf bytes.Buffer
    if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
        return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
    }
    return buf.String(), nil
}

// TemplateContext holds path-related template variables for dynamic directory names.
type TemplateContext struct {
    BaseDir     string // e.g., "spectr"
    SpecsDir    string // e.g., "spectr/specs"
    ChangesDir  string // e.g., "spectr/changes"
    ProjectFile string // e.g., "spectr/project.md"
    AgentsFile  string // e.g., "spectr/AGENTS.md"
}

// DefaultTemplateContext returns a TemplateContext with default values.
func DefaultTemplateContext() TemplateContext {
    return TemplateContext{
        BaseDir:     "spectr",
        SpecsDir:    "spectr/specs",
        ChangesDir:  "spectr/changes",
        ProjectFile: "spectr/project.md",
        AgentsFile:  "spectr/AGENTS.md",
    }
}
```

```go
// internal/domain/slashcmd.go
package domain

// SlashCommand represents a type-safe slash command identifier.
type SlashCommand int

const (
    SlashProposal SlashCommand = iota
    SlashApply
)

// String returns the command name for debugging.
func (s SlashCommand) String() string {
    names := []string{"proposal", "apply"}
    if int(s) < len(names) {
        return names[s]
    }
    return "unknown"
}

// TemplateName returns the template file name for this command.
func (s SlashCommand) TemplateName() string {
    names := map[SlashCommand]string{
        SlashProposal: "slash-proposal.md.tmpl",
        SlashApply:    "slash-apply.md.tmpl",
    }
    return names[s]
}
```

**Import structure after change:**
```
internal/domain          <- shared types (no dependencies on other internal packages)
    ├── template.go      <- TemplateRef, TemplateContext
    └── slashcmd.go      <- SlashCommand enum

internal/initialize/templates
    └── imports domain   <- uses domain.TemplateRef, domain.TemplateContext

internal/initialize/providers
    └── imports domain   <- uses domain.TemplateRef, domain.SlashCommand
```

**Benefits:**
- Clean separation of domain types from implementation
- No import cycles possible (domain has no internal dependencies)
- Types can be shared freely between packages
- Clear ownership of domain concepts

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

### 2. Registration API with Explicit Calls (No init())

**Decision**: Register providers explicitly from a central location, not via `init()`.

**Problem**: The current `init()` pattern has multiple issues:
```go
// CURRENT (BAD):
// 1. Error is discarded - registration failures are silent
// 2. init() is implicit - hard to test, hard to trace, order-dependent
// 3. No control over when registration happens
func init() {
    Register(NewOpencodeProvider())  // Returns error but it's ignored
}
```

**Solution**: Explicit registration from a central `RegisterAllProviders()` function:

```go
// RegisterProvider registers a provider and returns an error if registration fails.
func RegisterProvider(reg Registration) error {
    if reg.ID == "" {
        return fmt.Errorf("provider ID is required")
    }
    if reg.Provider == nil {
        return fmt.Errorf("provider implementation is required")
    }
    if _, exists := registry[reg.ID]; exists {
        return fmt.Errorf("provider %q already registered", reg.ID)
    }
    registry[reg.ID] = reg
    return nil
}

// Registration contains provider metadata and implementation.
type Registration struct {
    ID       string   // Unique identifier (kebab-case, e.g., "claude-code")
    Name     string   // Human-readable name (e.g., "Claude Code")
    Priority int      // Display order (lower = higher priority)
    Provider Provider // Implementation
}

// RegisterAllProviders registers all built-in providers.
// Called once at application startup (e.g., from main() or cmd/root.go).
// Returns error if any registration fails.
func RegisterAllProviders() error {
    providers := []Registration{
        {ID: "claude-code", Name: "Claude Code", Priority: 1, Provider: &ClaudeProvider{}},
        {ID: "gemini", Name: "Gemini CLI", Priority: 2, Provider: &GeminiProvider{}},
        {ID: "cursor", Name: "Cursor", Priority: 3, Provider: &CursorProvider{}},
        {ID: "cline", Name: "Cline", Priority: 4, Provider: &ClineProvider{}},
        {ID: "aider", Name: "Aider", Priority: 5, Provider: &AiderProvider{}},
        {ID: "opencode", Name: "OpenCode", Priority: 6, Provider: &OpencodeProvider{}},
        // ... all other providers
    }

    for _, reg := range providers {
        if err := RegisterProvider(reg); err != nil {
            return fmt.Errorf("failed to register %s provider: %w", reg.ID, err)
        }
    }
    return nil
}
```

**Application startup (cmd/root.go or main.go):**
```go
func init() {
    // Single init() that calls explicit registration with error handling
    if err := providers.RegisterAllProviders(); err != nil {
        panic(err)
    }
}

// Or even better - called explicitly from main():
func main() {
    if err := providers.RegisterAllProviders(); err != nil {
        log.Fatalf("failed to register providers: %v", err)
    }
    // ... rest of application
}
```

**Rationale**:
- **Testability**: Tests can call `RegisterProvider()` individually or skip registration entirely
- **Clarity**: All registrations in one place, easy to see what providers exist
- **Error handling**: Errors are explicit and propagated, not silently discarded
- **Debuggability**: Easy to set breakpoints, trace registration order
- **No implicit dependencies**: No reliance on init() ordering between packages

### 3. Built-in Initializers

**Decision**: Provide three composable initializers with type-safe template selection:

```go
// Creates directories (e.g., .claude/commands/spectr/)
func NewDirectoryInitializer(paths ...string) Initializer

// Creates/updates instruction file with markers
// Uses TemplateGetter for compile-time checked template selection
type TemplateGetter func(*TemplateManager) TemplateRef
func NewConfigFileInitializer(path string, getTemplate TemplateGetter) Initializer

// Creates slash commands from templates
// Uses []SlashCommand for compile-time checked command selection
func NewSlashCommandsInitializer(dir, ext string, commands []SlashCommand) Initializer
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
// internal/initialize/providers/claude.go
package providers

import (
    "context"

    "github.com/connerohnesorge/spectr/internal/domain"
)

// ClaudeProvider configures Claude Code with CLAUDE.md and .claude/commands/spectr/.
// No init() - registration happens in RegisterAllProviders().
type ClaudeProvider struct{}

func (p *ClaudeProvider) Initializers(ctx context.Context) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".claude/commands/spectr"),
        // Type-safe: (*TemplateManager).InstructionPointer is a method reference
        // Compiler catches typos - no runtime surprises
        NewConfigFileInitializer("CLAUDE.md", (*TemplateManager).InstructionPointer),
        // Type-safe: domain.SlashProposal/domain.SlashApply are typed constants from domain package
        NewSlashCommandsInitializer(".claude/commands/spectr", ".md", []domain.SlashCommand{
            domain.SlashProposal,
            domain.SlashApply,
        }),
    }
}
```

## Example: Gemini Provider (TOML format)

```go
// internal/initialize/providers/gemini.go
package providers

import (
    "context"

    "github.com/connerohnesorge/spectr/internal/domain"
)

// GeminiProvider configures Gemini CLI with TOML slash commands only.
// No init() - registration happens in RegisterAllProviders().
type GeminiProvider struct{}

func (p *GeminiProvider) Initializers(ctx context.Context) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".gemini/commands/spectr"),
        // No config file for Gemini - uses TOML slash commands only
        NewSlashCommandsInitializer(".gemini/commands/spectr", ".toml", []domain.SlashCommand{
            domain.SlashProposal,
            domain.SlashApply,
        }),
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

### 7. Type-Safe Template Selection

**Decision**: Use typed template references instead of raw strings for compile-time safety.

**Problem**: The original design had raw string inputs like:
```go
// Unsafe: typo "instrction-pointer" would fail at runtime
NewConfigFileInitializer("CLAUDE.md", "instruction-pointer")
```

**Solution**: Define typed template accessors on TemplateManager:

```go
// TemplateRef is a type-safe reference to a parsed template
type TemplateRef struct {
    name     string              // template file name (e.g., "instruction-pointer.md.tmpl")
    template *template.Template  // pre-parsed template
}

// Render executes the template with the given context
func (tr TemplateRef) Render(ctx TemplateContext) (string, error) {
    var buf bytes.Buffer
    if err := tr.template.ExecuteTemplate(&buf, tr.name, ctx); err != nil {
        return "", fmt.Errorf("failed to render template %s: %w", tr.name, err)
    }
    return buf.String(), nil
}

// TemplateManager exposes type-safe accessors for each template
type TemplateManager struct {
    templates *template.Template
}

// InstructionPointer returns the instruction-pointer.md.tmpl template reference
func (tm *TemplateManager) InstructionPointer() TemplateRef {
    return TemplateRef{
        name:     "instruction-pointer.md.tmpl",
        template: tm.templates,
    }
}

// Agents returns the AGENTS.md.tmpl template reference
func (tm *TemplateManager) Agents() TemplateRef {
    return TemplateRef{
        name:     "AGENTS.md.tmpl",
        template: tm.templates,
    }
}

// SlashCommand type for compile-time checked slash command types
type SlashCommand int

const (
    SlashProposal SlashCommand = iota
    SlashApply
)

// SlashCommand returns a template reference for the given slash command type
func (tm *TemplateManager) SlashCommand(cmd SlashCommand) TemplateRef {
    names := map[SlashCommand]string{
        SlashProposal: "slash-proposal.md.tmpl",
        SlashApply:    "slash-apply.md.tmpl",
    }
    return TemplateRef{
        name:     names[cmd],
        template: tm.templates,
    }
}
```

**Updated initializer constructors**:

```go
// ConfigFileInitializer now receives a TemplateRef getter function
// The getter is called with TemplateManager at Init() time
type TemplateGetter func(*TemplateManager) TemplateRef

func NewConfigFileInitializer(path string, getTemplate TemplateGetter) Initializer {
    return &ConfigFileInitializer{
        path:        path,
        getTemplate: getTemplate,
    }
}

// Usage - compile-time checked:
NewConfigFileInitializer("CLAUDE.md", (*TemplateManager).InstructionPointer)

// SlashCommandsInitializer receives slice of SlashCommand types
func NewSlashCommandsInitializer(dir, ext string, commands []SlashCommand) Initializer {
    return &SlashCommandsInitializer{
        dir:      dir,
        ext:      ext,
        commands: commands,
    }
}

// Usage - compile-time checked:
NewSlashCommandsInitializer(".claude/commands/spectr", ".md", []SlashCommand{
    SlashProposal,
    SlashApply,
})
```

**Benefits**:
- Compile-time errors for invalid template names (typos caught by compiler)
- IDE autocomplete for available templates
- Refactoring-safe: renaming a template accessor updates all usages
- Clear documentation of available templates through method signatures

**Alternatives considered:**
- String constants (`const TemplateInstructionPointer = "instruction-pointer"`) - Still strings, can be used incorrectly
- Template name validation at startup - Runtime error, not compile-time
- Passing `*template.Template` directly - Leaks implementation detail, harder to mock in tests

## Resolved Questions

| Question | Decision |
|----------|----------|
| Change detection approach? | InitResult return value from each Initializer |
| Initializer ordering? | Implicit ordering by type (Directory → ConfigFile → SlashCommands) |
| Rollback on partial failure? | No rollback; report failures, users re-run init |
| Template variable location? | Derived from SpectrDir via methods |
| Global paths support? | Two fs instances (projectFs, globalFs) |
| Deduplication key? | By file path (simple and effective) |
| Template selection type safety? | TemplateRef accessors with method references; compile-time checked |
| Import cycle resolution? | New `internal/domain` package with shared types (`TemplateRef`, `SlashCommand`, `TemplateContext`) |
| Registration error handling? | No init() - explicit `RegisterAllProviders()` called at startup with proper error propagation |

## Future Considerations (Out of Scope)

- `spectr init --dry-run` to preview changes without applying
- New instruction file support for Gemini, Cursor, Aider, OpenCode (separate proposal)
