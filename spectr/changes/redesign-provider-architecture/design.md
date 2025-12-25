# Design: Provider Architecture Redesign

## Context

Spectr supports 15 AI CLI/IDE tools (Claude Code, Cursor, Cline, etc.) through a provider system. Each provider configures:
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
- Maintain support for all 15 current providers
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
```

**Slash command template consolidation:**

The slash command templates (`slash-proposal.md.tmpl`, `slash-apply.md.tmpl`) will be moved from `internal/initialize/templates/tools/` into the `internal/domain` package. This achieves full consolidation by:

1. **Embedding templates in domain package:**
   ```go
   // internal/domain/templates.go
   package domain

   import "embed"

   //go:embed templates/*.tmpl
   var TemplateFS embed.FS
   ```

2. **Template file locations:**
   - Move: `internal/initialize/templates/tools/slash-proposal.md.tmpl` → `internal/domain/templates/slash-proposal.md.tmpl`
   - Move: `internal/initialize/templates/tools/slash-apply.md.tmpl` → `internal/domain/templates/slash-apply.md.tmpl`
   - Create: `internal/domain/templates/slash-proposal.toml.tmpl` (for Gemini TOML format)
   - Create: `internal/domain/templates/slash-apply.toml.tmpl` (for Gemini TOML format)

3. **Template access via TemplateManager:**
   Templates are accessed through `TemplateManager.SlashCommand(cmd)`, NOT via a method on `SlashCommand` itself. This keeps template rendering concerns separate from the domain type.

**Import structure after change:**
```
internal/domain                    <- shared types + embedded slash templates
    ├── template.go                <- TemplateRef, TemplateContext
    ├── slashcmd.go                <- SlashCommand enum (NO TemplateName method)
    ├── templates.go               <- embed.FS for slash command templates
    └── templates/                 <- embedded template files
        ├── slash-proposal.md.tmpl     <- Markdown format
        ├── slash-apply.md.tmpl        <- Markdown format
        ├── slash-proposal.toml.tmpl   <- TOML format (Gemini)
        └── slash-apply.toml.tmpl      <- TOML format (Gemini)

internal/initialize/templates      <- main template manager
    ├── imports domain             <- uses domain.TemplateRef, domain.TemplateContext, domain.TemplateFS
    └── templates/                 <- instruction and doc templates (non-slash)
        ├── AGENTS.md.tmpl
        └── instruction-pointer.md.tmpl

internal/initialize/providers
    └── imports domain             <- uses domain.TemplateRef, domain.SlashCommand
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
    // Initializers returns the list of initializers for this provider.
    // Receives TemplateManager to allow passing TemplateRef directly to ConfigFileInitializer.
    Initializers(ctx context.Context, tm *TemplateManager) []Initializer
}

// InitResult contains the files created or modified by an initializer.
type InitResult struct {
    CreatedFiles []string
    UpdatedFiles []string
}

type Initializer interface {
    // Init creates or updates files. Returns result with file changes and error if initialization fails.
    // Must be idempotent (safe to run multiple times).
    // Receives both filesystems - initializer decides which to use based on its configuration.
    Init(ctx context.Context, projectFs, globalFs afero.Fs, cfg *Config, tm *TemplateManager) (InitResult, error)

    // IsSetup returns true if this initializer's artifacts already exist.
    // Receives both filesystems - initializer checks the appropriate one.
    IsSetup(projectFs, globalFs afero.Fs, cfg *Config) bool
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
// RegisteredProviders returns all registered providers sorted by priority (lowest first).
func RegisteredProviders() []Registration {
    result := make([]Registration, 0, len(registry))
    for _, reg := range registry {
        result = append(result, reg)
    }
    sort.Slice(result, func(i, j int) bool {
        return result[i].Priority < result[j].Priority
    })
    return result
}

func RegisterAllProviders() error {
    providers := []Registration{
        {ID: "claude-code", Name: "Claude Code", Priority: 1, Provider: &ClaudeProvider{}},
        {ID: "gemini", Name: "Gemini CLI", Priority: 2, Provider: &GeminiProvider{}},
        {ID: "costrict", Name: "Costrict", Priority: 3, Provider: &CostrictProvider{}},
        {ID: "qoder", Name: "Qoder", Priority: 4, Provider: &QoderProvider{}},
        {ID: "qwen", Name: "Qwen", Priority: 5, Provider: &QwenProvider{}},
        {ID: "antigravity", Name: "Antigravity", Priority: 6, Provider: &AntigravityProvider{}},
        {ID: "cline", Name: "Cline", Priority: 7, Provider: &ClineProvider{}},
        {ID: "cursor", Name: "Cursor", Priority: 8, Provider: &CursorProvider{}},
        {ID: "codex", Name: "Codex CLI", Priority: 9, Provider: &CodexProvider{}},
        {ID: "aider", Name: "Aider", Priority: 10, Provider: &AiderProvider{}},
        {ID: "windsurf", Name: "Windsurf", Priority: 11, Provider: &WindsurfProvider{}},
        {ID: "kilocode", Name: "Kilocode", Priority: 12, Provider: &KilocodeProvider{}},
        {ID: "continue", Name: "Continue", Priority: 13, Provider: &ContinueProvider{}},
        {ID: "crush", Name: "Crush", Priority: 14, Provider: &CrushProvider{}},
        {ID: "opencode", Name: "OpenCode", Priority: 15, Provider: &OpencodeProvider{}},
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

**Why not keep a compatibility shim?**

We considered adding a deprecated `Register(_ any)` function during migration:

```go
// REJECTED APPROACH:
// Deprecated: Use RegisterProvider instead
func Register(_ any) {
    // Silently swallow calls to prevent compilation errors
    // OR log a deprecation warning
}
```

**Why we rejected this:**

1. **Hides problems**: Code appears to work but providers aren't actually registered
2. **Creates ambiguity**: Two registration paths exist, unclear which to use
3. **Technical debt**: Deprecated code that must be maintained and eventually removed
4. **Violates zero tech debt policy**: We have an explicit policy against keeping deprecated code around
5. **Delayed migration**: Developers can postpone fixing the real problem
6. **No benefit**: This is a single change that migrates all providers at once - there's no partial migration state where compatibility is needed

**Our approach instead:**
- Complete removal of old `Register()` function
- Any code calling it will fail to compile
- Developers get clear error: `undefined: Register`
- Forces explicit migration to new system
- No hidden behavior, no confusion
- Clean codebase from day one

### 3. Built-in Initializers

**Decision**: Provide composable initializers with type-safe template selection. Use separate types for local vs global filesystem operations.

```go
// Creates directories in project filesystem (e.g., .claude/commands/spectr/)
func NewDirectoryInitializer(paths ...string) Initializer

// Creates directories in global filesystem (e.g., ~/.codex/prompts/)
func NewGlobalDirectoryInitializer(paths ...string) Initializer

// Creates/updates instruction file with markers
// Takes TemplateRef directly (not a function) for simpler API
func NewConfigFileInitializer(path string, template TemplateRef) Initializer

// Creates Markdown slash commands from templates in project filesystem
// Uses []SlashCommand for compile-time checked command selection
func NewSlashCommandsInitializer(dir string, commands []SlashCommand) Initializer

// Creates Markdown slash commands from templates in global filesystem
func NewGlobalSlashCommandsInitializer(dir string, commands []SlashCommand) Initializer

// Creates TOML slash commands from templates in project filesystem (Gemini only)
func NewTOMLSlashCommandsInitializer(dir string, commands []SlashCommand) Initializer
```

**Rationale for separate types**:
- Clear intent: Type name makes filesystem and format choice explicit
- No constructor parameters for format/filesystem: `TOMLSlashCommandsInitializer` vs `SlashCommandsInitializer`
- Type-safe deduplication: can dedupe by type + path without checking internal flags
- Simpler API: no `ext` parameter needed

**Deduplication**: When multiple providers share an initializer with same config, run once.

**ConfigFileInitializer Marker Handling**:

The ConfigFileInitializer uses `<!-- spectr:START -->` and `<!-- spectr:END -->` markers to update instruction files. It must handle orphaned markers correctly:

```go
// Pseudocode for marker handling logic
func updateWithMarkers(content, newContent, startMarker, endMarker string) string {
    startIdx := strings.Index(content, startMarker)

    if startIdx == -1 {
        // No markers exist - append new block at end
        return content + "\n\n" + startMarker + newContent + endMarker
    }

    // Start marker found - look for end marker AFTER the start
    searchFrom := startIdx + len(startMarker)
    endIdx := strings.Index(content[searchFrom:], endMarker)

    if endIdx != -1 {
        // Normal case: both markers present
        endIdx += searchFrom // Adjust to absolute position
        before := content[:startIdx]
        after := content[endIdx+len(endMarker):]
        return before + startMarker + newContent + endMarker + after
    }

    // Start marker exists but no end marker immediately after
    // Search for trailing end marker anywhere in the file
    trailingEndIdx := strings.LastIndex(content, endMarker)

    if trailingEndIdx > startIdx {
        // Found end marker after start - use it
        before := content[:startIdx]
        after := content[trailingEndIdx+len(endMarker):]
        return before + startMarker + newContent + endMarker + after
    }

    // No end marker anywhere after start - orphaned start marker
    // Replace everything from start marker onward with new block
    before := content[:startIdx]
    return before + startMarker + newContent + endMarker
}
```

**Rationale**:
- **Prevents duplicate blocks**: Don't append when start marker already exists
- **Handles corrupted files**: Recovers from missing end markers
- **Uses trailing end marker**: If end marker exists elsewhere, use it
- **Clean replacement**: Orphaned start markers are replaced, not duplicated

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
- Initializers receive both filesystems and decide internally which to use based on their configuration

### 5. File Change Detection

**Decision**: Each initializer returns an `InitResult` containing the files it created/updated.

```go
// ExecutionResult aggregates results from all initializers
type ExecutionResult struct {
    CreatedFiles []string  // All files created across all initializers
    UpdatedFiles []string  // All files updated across all initializers
    Errors       []error   // Any errors encountered during initialization
}

// aggregateResults combines multiple InitResult values into a single ExecutionResult
func aggregateResults(results []InitResult, errors []error) ExecutionResult {
    var created, updated []string
    for _, r := range results {
        created = append(created, r.CreatedFiles...)
        updated = append(updated, r.UpdatedFiles...)
    }
    return ExecutionResult{
        CreatedFiles: created,
        UpdatedFiles: updated,
        Errors:       errors,
    }
}

// Fail-fast execution: stop on first error
var allResults []InitResult

for _, init := range allInitializers {
    result, err := init.Init(ctx, projectFs, globalFs, cfg, tm)
    if err != nil {
        // Fail fast: return immediately on first error
        return ExecutionResult{
            CreatedFiles: collectCreatedFiles(allResults),
            UpdatedFiles: collectUpdatedFiles(allResults),
            Errors:       []error{err},
        }, err
    }
    allResults = append(allResults, result)
}

// All succeeded
return aggregateResults(allResults, nil), nil
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

func (p *ClaudeProvider) Initializers(ctx context.Context, tm *TemplateManager) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".claude/commands/spectr"),
        // TemplateRef passed directly - simpler API, template resolved at provider construction
        NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer()),
        // Type-safe: domain.SlashProposal/domain.SlashApply are typed constants from domain package
        // SlashCommandsInitializer creates Markdown files (.md)
        NewSlashCommandsInitializer(".claude/commands/spectr", []domain.SlashCommand{
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

func (p *GeminiProvider) Initializers(ctx context.Context, tm *TemplateManager) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".gemini/commands/spectr"),
        // No config file for Gemini - uses TOML slash commands only
        // TOMLSlashCommandsInitializer uses slash-*.toml.tmpl templates
        NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", []domain.SlashCommand{
            domain.SlashProposal,
            domain.SlashApply,
        }),
    }
}
```

**Template Files by Initializer Type**:
- `SlashCommandsInitializer` / `GlobalSlashCommandsInitializer` → `slash-proposal.md.tmpl`, `slash-apply.md.tmpl`
- `TOMLSlashCommandsInitializer` → `slash-proposal.toml.tmpl`, `slash-apply.toml.tmpl`

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
    case *DirectoryInitializer, *GlobalDirectoryInitializer:
        return 1
    case *ConfigFileInitializer:
        return 2
    case *SlashCommandsInitializer, *GlobalSlashCommandsInitializer, *TOMLSlashCommandsInitializer:
        return 3
    default:
        return 99
    }
}
```

**Rationale**: Directories must exist before files can be written. This ordering is implicit but guaranteed.

### 9. Initializer Deduplication

**Decision**: Deduplicate by initializer identity. Same initializer configuration = run once.

When initializers are collected from multiple providers, deduplication is based on the initializer type and path. The separate Global* types make this explicit:

```go
// deduplicatable is an optional interface for initializers that support deduplication.
// Initializers that implement this interface can be deduplicated based on their key.
type deduplicatable interface {
    dedupeKey() string
}

func dedupeInitializers(all []Initializer) []Initializer {
    seen := make(map[string]bool)
    var result []Initializer
    for _, init := range all {
        // Check if initializer supports deduplication
        if d, ok := init.(deduplicatable); ok {
            key := d.dedupeKey()
            if seen[key] {
                continue // Skip duplicate
            }
            seen[key] = true
        }
        result = append(result, init)
    }
    return result
}

// Example dedupeKey implementations (type name encodes filesystem and format):
// DirectoryInitializer:           "DirectoryInitializer:.claude/commands/spectr"
// GlobalDirectoryInitializer:     "GlobalDirectoryInitializer:.codex/prompts"
// ConfigFileInitializer:          "ConfigFileInitializer:CLAUDE.md"
// SlashCommandsInitializer:       "SlashCommandsInitializer:.claude/commands/spectr"
// GlobalSlashCommandsInitializer: "GlobalSlashCommandsInitializer:.codex/prompts"
// TOMLSlashCommandsInitializer:   "TOMLSlashCommandsInitializer:.gemini/commands/spectr"
```

**Example**: If Claude Code and Cline both return `ConfigFileInitializer{path: "CLAUDE.md"}`, only one runs.

**Rationale**:
- Type-based keys: `GlobalDirectoryInitializer` and `DirectoryInitializer` have different type names, preventing accidental deduplication across filesystem boundaries
- Optional interface: Initializers that don't need deduplication don't have to implement it
- Simple and covers the common case (multiple providers sharing same file)

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking change for users | Clear migration docs; `spectr init` re-run required |
| Path-based dedup misses edge cases | Simple first; can enhance if needed |
| No rollback on failure | Clear error reporting; users can re-run init |

## Migration Plan

### Zero Technical Debt Policy

This migration follows a **zero technical debt** policy - no compatibility shims, no deprecated functions kept around. Clean break, complete removal of old code.

### Migration Steps

1. Implement new provider system (new files, new types)
2. Migrate all 15 providers in-place (delete old code, write new code in same file)
3. **Completely remove** old registration system - no `Register(p Provider)` function kept around
4. Update docs to explain re-initialization requirement
5. No rollback needed - old configs continue to work

### No Compatibility Shims

**We will NOT do this:**
```go
// BAD: Deprecated compatibility shim that silently swallows calls
func Register(_ any) {
    // No-op to prevent compilation errors during migration
}
```

This approach:
- Hides migration problems instead of surfacing them
- Creates technical debt
- Makes it unclear which registration method to use
- Allows old code to "work" but not actually register providers

**Instead, we do this:**
- Remove `func Register(p Provider)` entirely from registry.go
- Remove all `init()` functions from provider files
- Implement `RegisterAllProviders()` in one central location
- Any attempt to call old `Register()` will cause a **compile-time error**
- Developers must explicitly migrate to the new registration system

### Migration is All-or-Nothing

The migration happens in a single change. There is no "partial migration" state where some providers use init() and others don't. The tasks are ordered to ensure:

1. New system is fully implemented first (sections 0-4)
2. All 15 providers are migrated together (section 5)
3. Old code is completely removed (section 7)

This ensures no ambiguity about which registration approach to use.

### Enforcement: No Compatibility Shims Allowed

If during implementation someone is tempted to add:

```go
// DON'T DO THIS - violates zero tech debt policy
func Register(_ any) {
    log.Printf("WARNING: Register() is deprecated, use RegisterProvider()")
    // ... caller info ...
}
```

**STOP.** This violates the zero technical debt policy. Instead:

1. **Delete the old Register() function entirely** (task 7.6)
2. **Delete all init() functions** that call it (tasks 5.1-5.16)
3. **Implement RegisterAllProviders()** in one location (task 4.6)
4. Let the **compiler enforce migration** with clear errors

The only valid reason to keep deprecated code would be if external code depends on it. Since all providers are internal to this codebase, there is no valid reason to keep compatibility shims.

### 6. TemplateManager Integration

**Decision**: Initializers receive `*TemplateManager` instead of implementing `TemplateRenderer`.

```go
type Initializer interface {
    Init(ctx context.Context, projectFs, globalFs afero.Fs, cfg *Config, tm *TemplateManager) (InitResult, error)
    IsSetup(projectFs, globalFs afero.Fs, cfg *Config) bool
}
```

**TemplateManager loads templates from both locations:**

```go
// NewTemplateManager creates a new template manager with all embedded templates loaded.
// It merges templates from:
// 1. internal/initialize/templates (main templates: AGENTS.md, instruction-pointer.md)
// 2. internal/domain (slash command templates: slash-proposal.md, slash-apply.md)
func NewTemplateManager() (*TemplateManager, error) {
    // Parse main templates
    mainTmpl, err := template.ParseFS(
        templateFS,  // internal/initialize/templates embed.FS
        "templates/**/*.tmpl",
    )
    if err != nil {
        return nil, fmt.Errorf("failed to parse main templates: %w", err)
    }

    // Parse and merge domain templates (slash commands)
    domainTmpl, err := template.ParseFS(
        domain.TemplateFS,  // internal/domain embed.FS
        "templates/*.tmpl",
    )
    if err != nil {
        return nil, fmt.Errorf("failed to parse domain templates: %w", err)
    }

    // Merge: add domain templates to main template set
    for _, t := range domainTmpl.Templates() {
        if _, err := mainTmpl.AddParseTree(t.Name(), t.Tree); err != nil {
            return nil, fmt.Errorf("failed to merge template %s: %w", t.Name(), err)
        }
    }

    return &TemplateManager{
        templates: mainTmpl,
    }, nil
}
```

**Rationale**:
- Reuses existing `TemplateManager` from `internal/initialize/templates.go`
- Merges templates from both `internal/initialize` and `internal/domain` packages
- Avoids duplicating template rendering logic in each initializer
- Simpler than the old `TemplateRenderer` interface pattern
- Domain package remains dependency-free (only exposes embed.FS, doesn't use it)

**Alternatives considered:**
- Each initializer implements own template rendering - More code duplication
- Pass templates as strings - Less flexible, harder to maintain
- Keep all templates in `internal/initialize` - Creates import cycle for domain types

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
    Name     string              // template file name (e.g., "instruction-pointer.md.tmpl")
    Template *template.Template  // pre-parsed template
}

// Render executes the template with the given context
func (tr TemplateRef) Render(ctx TemplateContext) (string, error) {
    var buf bytes.Buffer
    if err := tr.Template.ExecuteTemplate(&buf, tr.Name, ctx); err != nil {
        return "", fmt.Errorf("failed to render template %s: %w", tr.Name, err)
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
        Name:     "instruction-pointer.md.tmpl",
        Template: tm.templates,
    }
}

// Agents returns the AGENTS.md.tmpl template reference
func (tm *TemplateManager) Agents() TemplateRef {
    return TemplateRef{
        Name:     "AGENTS.md.tmpl",
        Template: tm.templates,
    }
}

// SlashCommand returns a Markdown template reference for the given slash command type.
// Used by SlashCommandsInitializer and GlobalSlashCommandsInitializer.
func (tm *TemplateManager) SlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
    names := map[domain.SlashCommand]string{
        domain.SlashProposal: "slash-proposal.md.tmpl",
        domain.SlashApply:    "slash-apply.md.tmpl",
    }
    return domain.TemplateRef{
        Name:     names[cmd],
        Template: tm.templates,
    }
}

// TOMLSlashCommand returns a TOML template reference for the given slash command type.
// Used by TOMLSlashCommandsInitializer (Gemini only).
func (tm *TemplateManager) TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef {
    names := map[domain.SlashCommand]string{
        domain.SlashProposal: "slash-proposal.toml.tmpl",
        domain.SlashApply:    "slash-apply.toml.tmpl",
    }
    return domain.TemplateRef{
        Name:     names[cmd],
        Template: tm.templates,
    }
}
```

**Updated initializer constructors**:

```go
// ConfigFileInitializer receives TemplateRef directly
// TemplateRef is resolved at provider construction time when Initializers() is called
func NewConfigFileInitializer(path string, template TemplateRef) Initializer {
    return &ConfigFileInitializer{
        path:     path,
        template: template,
    }
}

// Usage - compile-time checked, TemplateRef passed directly:
// Provider.Initializers(ctx, tm) receives TemplateManager
NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer())

// SlashCommandsInitializer receives slice of SlashCommand types (creates .md files)
func NewSlashCommandsInitializer(dir string, commands []SlashCommand) Initializer {
    return &SlashCommandsInitializer{
        dir:      dir,
        commands: commands,
    }
}

// TOMLSlashCommandsInitializer receives slice of SlashCommand types (creates .toml files)
func NewTOMLSlashCommandsInitializer(dir string, commands []SlashCommand) Initializer {
    return &TOMLSlashCommandsInitializer{
        dir:      dir,
        commands: commands,
    }
}

// Usage - compile-time checked, type determines format:
NewSlashCommandsInitializer(".claude/commands/spectr", []SlashCommand{
    SlashProposal,
    SlashApply,
})

NewTOMLSlashCommandsInitializer(".gemini/commands/spectr", []SlashCommand{
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
| Partial failure handling? | Fail-fast: stop on first error, return partial results |
| Template variable location? | Derived from SpectrDir via methods |
| Global paths support? | Separate types: `GlobalDirectoryInitializer`, `GlobalSlashCommandsInitializer` |
| Deduplication key? | By type name + path via optional `deduplicatable` interface |
| Template selection type safety? | TemplateRef passed directly to ConfigFileInitializer; Provider.Initializers() receives TemplateManager |
| Import cycle resolution? | New `internal/domain` package with shared types (`TemplateRef`, `SlashCommand`, `TemplateContext`) |
| Registration error handling? | No init() - explicit `RegisterAllProviders()` called in `cmd/root.go` with proper error propagation |
| Initializer interface methods? | Two methods only: `Init()` and `IsSetup()` - both receive both filesystems |
| TemplateRef field visibility? | Public fields (`Name`, `Template`) for external package access |
| TOML template support? | Separate `TOMLSlashCommandsInitializer` type with dedicated `TOMLSlashCommand()` accessor |
| Frontmatter structure? | Minimal: only `description` field required |

## Future Considerations (Out of Scope)

- `spectr init --dry-run` to preview changes without applying
- New instruction file support for Gemini, Cursor, Aider, OpenCode (separate proposal)
