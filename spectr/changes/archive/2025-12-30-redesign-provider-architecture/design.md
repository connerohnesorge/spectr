# Design: Provider Architecture Redesign

## Context

Spectr supports 15 AI CLI/IDE tools (Claude Code, Cursor, Cline, etc.) through a
provider system. Each provider configures:

1. An instruction file (e.g., `CLAUDE.md`) with marker-based updates
2. Slash commands (e.g., `.claude/commands/spectr/proposal.md`)

The current implementation has each provider implement a 12-method interface,
with most embedding `BaseProvider`. This leads to ~50 lines of boilerplate per
provider when the actual variance is just configuration values.

**Current Problems:**

1. **Import cycle**: `providers.TemplateManager` cannot import
  `internal/initialize/templates` without creating an import cycle, forcing the
  use of `any` as placeholder types
2. **Silent registration failures**: Provider registration in `init()` assigns
  errors to the blank identifier, silently discarding failures

## Scope

**Behavioral equivalence**: Architectural refactor maintaining identical
user-facing behavior. No new user-facing features. Internal implementation uses
new patterns (domain package, initializer interface, dual filesystems) to
eliminate boilerplate.

## Goals / Non-Goals

**Goals:**

- Reduce provider authoring to ~10 lines of registration code
- Enable sharing and deduplication of common initialization logic
- Improve testability of initialization steps
- Maintain support for all 15 current providers
- Use `afero.Fs` rooted at project for cleaner path handling
- Explicit change tracking via ExecutionResult return values
- Break import cycles with a dedicated `internal/domain` package
- Fail-fast registration with explicit error handling (no silent failures)

**Non-Goals:**

- Runtime plugin loading (all providers compiled in)
- Backwards compatibility with existing configurations
- Rollback on partial failure
- New instruction file formats (separate proposal)

## Decisions

### 0. Domain Package for Shared Types

**Decision**: Create `internal/domain` package containing shared domain types to
break import cycles.

**Problem**: The current architecture has an import cycle issue:

- `providers.TemplateManager` needs to reference template types like
  `TemplateRef` and `SlashCommand`
- These types are defined in `internal/initialize/templates`
- But `templates` cannot import `providers` and vice versa without creating a
  cycle
- Currently, `any` is used as a placeholder, with the concrete adapter in
  `executor.go` explaining the real types

**Solution**: Extract domain objects into `internal/domain`:

```go
// internal/domain/template.go
package domain

import (
    "html/template"
)

// TemplateRef is a type-safe reference to a parsed template.
// It serves as a lightweight typed handle that can be safely passed
// between packages without creating import cycles.
// Rendering is performed by TemplateManager, not by TemplateRef
// itself.
type TemplateRef struct {
    Name     string              // template file name (e.g., "instruction-pointer.md.tmpl")
    Template *template.Template  // pre-parsed template
}

// TemplateContext holds path-related template variables for dynamic
// directory names.
// Created via templateContextFromConfig(cfg) in the executor, not via
// defaults.
type TemplateContext struct {
    BaseDir     string // e.g., "spectr" (from cfg.SpectrDir)
    SpecsDir    string // e.g., "spectr/specs" (from cfg.SpecsDir())
    ChangesDir  string // e.g., "spectr/changes" (from cfg.ChangesDir())
    ProjectFile string // e.g., "spectr/project.md" (from cfg.ProjectFile())
    AgentsFile  string // e.g., "spectr/AGENTS.md" (from cfg.AgentsFile())
}
```text

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
```text

**Slash command template consolidation:**

The slash command templates (`slash-proposal.md.tmpl`, `slash-apply.md.tmpl`)
will be moved from `internal/initialize/templates/tools/` into the
`internal/domain` package. This achieves full consolidation by:

1. **Embedding templates in domain package:**

   ```go
   // internal/domain/templates.go
   package domain

   import "embed"

   //go:embed templates/*.tmpl
   var TemplateFS embed.FS
```text

2. **Template file locations:**
   - Move: `internal/initialize/templates/tools/slash-proposal.md.tmpl` →
     `internal/domain/templates/slash-proposal.md.tmpl`
   - Move: `internal/initialize/templates/tools/slash-apply.md.tmpl` →
     `internal/domain/templates/slash-apply.md.tmpl`
   - Create: `internal/domain/templates/slash-proposal.toml.tmpl` (for Gemini
     TOML format)
   - Create: `internal/domain/templates/slash-apply.toml.tmpl` (for Gemini TOML
     format)

3. **Template access via TemplateManager:**
   Templates are accessed through `TemplateManager.SlashCommand(cmd)`, NOT via a
   method on `SlashCommand` itself. This keeps template rendering concerns
   separate from the domain type.

**Import structure after change:**

```text
internal/domain                    <- shared types + embedded slash
                                      templates
    ├── template.go                <- TemplateRef, TemplateContext
    ├── slashcmd.go                <- SlashCommand enum (NO
                                      TemplateName method)
    ├── templates.go               <- embed.FS for slash command
                                      templates
    └── templates/                 <- embedded template files
        ├── slash-proposal.md.tmpl     <- Markdown format
        ├── slash-apply.md.tmpl        <- Markdown format
        ├── slash-proposal.toml.tmpl   <- TOML format (Gemini)
        └── slash-apply.toml.tmpl      <- TOML format (Gemini)

internal/initialize/templates      <- main template manager
    ├── imports domain             <- uses domain.TemplateRef,
                                      domain.TemplateContext,
                                      domain.TemplateFS
    └── templates/                 <- instruction and doc templates
                                      (non-slash)
        ├── AGENTS.md.tmpl
        └── instruction-pointer.md.tmpl

internal/initialize/providers
    └── imports domain             <- uses domain.TemplateRef,
                                      domain.SlashCommand
```text

**Benefits:**

- Clean separation of domain types from implementation
- No import cycles possible (domain has no internal dependencies)
- Types can be shared freely between packages
- Clear ownership of domain concepts

### 1. Provider Interface

**Decision**: Providers return a list of initializers; metadata lives at
registration.

```go
type Provider interface {
    // Initializers returns the list of initializers for this provider.
    // Receives TemplateManager to allow passing TemplateRef directly to ConfigFileInitializer.
    Initializers(ctx context.Context, tm *TemplateManager) []Initializer
}

type Initializer interface {
    // Init creates or updates files. Returns result with file changes and
    // error if initialization fails.
    // Must be idempotent (safe to run multiple times).
    // Receives both filesystems - initializer decides which to use based on
    // its configuration.
    Init(ctx context.Context, projectFs, homeFs afero.Fs, cfg *Config,
        tm *TemplateManager) (ExecutionResult, error)

    // IsSetup returns true if this initializer's artifacts already exist.
    // Receives both filesystems - initializer checks the appropriate one.
    // PURPOSE: Used by the setup wizard to show which providers are already
    // configured.
    // NOT used to skip initializers during execution - Init() always runs
    // (idempotent).
    IsSetup(projectFs, homeFs afero.Fs, cfg *Config) bool
}

type Config struct {
    SpectrDir string // e.g., "spectr" (relative to fs root)
}

// Validate checks Config fields for basic correctness.
func (c *Config) Validate() error {
    if c.SpectrDir == "" {
        return fmt.Errorf("SpectrDir must not be empty")
    }
    if strings.HasPrefix(c.SpectrDir, "/") {
        return fmt.Errorf("SpectrDir must be relative, not absolute")
    }
    if strings.Contains(c.SpectrDir, "..") {
        return fmt.Errorf("SpectrDir must not contain path traversal")
    }
    return nil
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
```text

**Alternatives considered:**

- Keep metadata in Provider interface (current design) - More boilerplate
- Use functional options pattern - Harder to test
- Store all paths in Config - Redundant, error-prone

### 2. Registration API with Explicit Calls (No init())

**Decision**: Register providers explicitly from a central location, not via
`init()`.

**Problem**: The current `init()` pattern has multiple issues:

```go
// CURRENT (BAD):
// 1. Error is discarded - registration failures are silent
// 2. init() is implicit - hard to test, hard to trace, order-dependent
// 3. No control over when registration happens
func init() {
    Register(NewOpencodeProvider())  // Returns error but it's ignored
}
```text

**Solution**: Explicit registration from a central `RegisterAllProviders()`
function:

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
// RegisteredProviders returns all registered providers sorted by priority
// (lowest first).
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
        {ID: "claude-code", Name: "Claude Code", Priority: 1,
            Provider: &ClaudeProvider{}},
        {ID: "gemini", Name: "Gemini CLI", Priority: 2,
            Provider: &GeminiProvider{}},
        {ID: "costrict", Name: "CoStrict", Priority: 3,
            Provider: &CostrictProvider{}},
        {ID: "qoder", Name: "Qoder", Priority: 4,
            Provider: &QoderProvider{}},
        {ID: "qwen", Name: "Qwen Code", Priority: 5,
            Provider: &QwenProvider{}},
        {ID: "antigravity", Name: "Antigravity", Priority: 6,
            Provider: &AntigravityProvider{}},
        {ID: "cline", Name: "Cline", Priority: 7,
            Provider: &ClineProvider{}},
        {ID: "cursor", Name: "Cursor", Priority: 8,
            Provider: &CursorProvider{}},
        {ID: "codex", Name: "Codex CLI", Priority: 9,
            Provider: &CodexProvider{}},
        {ID: "aider", Name: "Aider", Priority: 10,
            Provider: &AiderProvider{}},
        {ID: "windsurf", Name: "Windsurf", Priority: 11,
            Provider: &WindsurfProvider{}},
        {ID: "kilocode", Name: "Kilocode", Priority: 12,
            Provider: &KilocodeProvider{}},
        {ID: "continue", Name: "Continue", Priority: 13,
            Provider: &ContinueProvider{}},
        {ID: "crush", Name: "Crush", Priority: 14,
            Provider: &CrushProvider{}},
        {ID: "opencode", Name: "OpenCode", Priority: 15,
            Provider: &OpencodeProvider{}},
    }

    for _, reg := range providers {
        if err := RegisterProvider(reg); err != nil {
            // Note: Successfully registered providers remain registered
            // (no rollback)
            return fmt.Errorf("failed to register %s provider: %w",
                reg.ID, err)
        }
    }
    return nil
}
```text

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
```text

**Rationale**:

- **Testability**: Tests can call `RegisterProvider()` individually or skip
  registration entirely
- **Clarity**: All registrations in one place, easy to see what providers
  exist
- **Error handling**: Errors are explicit and propagated, not silently
  discarded
- **Debuggability**: Easy to set breakpoints, trace registration order
- **No implicit dependencies**: No reliance on init() ordering between
  packages

**Why not keep a compatibility shim?**

We considered adding a deprecated `Register(_ any)` function during migration:

```go
// REJECTED APPROACH:
// Deprecated: Use RegisterProvider instead
func Register(_ any) {
    // Silently swallow calls to prevent compilation errors
    // OR log a deprecation warning
}
```text

**Why we rejected this:**

1. **Hides problems**: Code appears to work but providers aren't actually
  registered
2. **Creates ambiguity**: Two registration paths exist, unclear which to
  use
3. **Technical debt**: Deprecated code that must be maintained and
  eventually removed
4. **Violates zero tech debt policy**: We have an explicit policy against
  keeping deprecated code around
5. **Delayed migration**: Developers can postpone fixing the real problem
6. **No benefit**: This is a single change that migrates all providers at
  once - there's no partial migration state where compatibility is needed

**Our approach instead:**

- Complete removal of old `Register()` function
- Any code calling it will fail to compile
- Developers get clear error: `undefined: Register`
- Forces explicit migration to new system
- No hidden behavior, no confusion
- Clean codebase from day one

### 3. Built-in Initializers

**Decision**: Provide composable initializers with type-safe template
selection. Use separate types for project vs home filesystem operations.

```go
// Creates directories in project filesystem
// (e.g., .claude/commands/spectr/)
func NewDirectoryInitializer(paths ...string) Initializer

// Creates directories in home filesystem (e.g., ~/.codex/prompts/)
func NewHomeDirectoryInitializer(paths ...string) Initializer

// Creates/updates instruction file with markers
// Takes TemplateRef directly (not a function) for simpler API
func NewConfigFileInitializer(path string, template TemplateRef) Initializer

// Creates Markdown slash commands from templates in project filesystem
// Uses map[SlashCommand]TemplateRef for early binding (templates resolved
// at construction)
func NewSlashCommandsInitializer(dir string,
    commands map[SlashCommand]TemplateRef) Initializer

// Creates Markdown slash commands from templates in home filesystem
func NewHomeSlashCommandsInitializer(dir string,
    commands map[SlashCommand]TemplateRef) Initializer

// Creates Markdown slash commands with custom prefix in project filesystem
// (for Antigravity)
// Output: {dir}/{prefix}{command}.md
// (e.g., .agent/workflows/spectr-proposal.md)
func NewPrefixedSlashCommandsInitializer(dir, prefix string,
    commands map[SlashCommand]TemplateRef) Initializer

// Creates Markdown slash commands with custom prefix in home filesystem
// (for Codex)
// Output: {dir}/{prefix}{command}.md
// (e.g., ~/.codex/prompts/spectr-proposal.md)
func NewHomePrefixedSlashCommandsInitializer(dir, prefix string,
    commands map[SlashCommand]TemplateRef) Initializer

// Creates TOML slash commands from templates in project filesystem
// (Gemini only)
func NewTOMLSlashCommandsInitializer(dir string,
    commands map[SlashCommand]TemplateRef) Initializer
```text

**Rationale for separate types**:

- Clear intent: Type name makes filesystem and format choice explicit
- No constructor parameters for format/filesystem:
  `TOMLSlashCommandsInitializer` vs `SlashCommandsInitializer`
- Type-safe deduplication: can dedupe by type + path without checking
  internal flags
- Simpler API: no `ext` parameter needed
- PrefixedSlashCommandsInitializer handles special cases (Antigravity,
  Codex) with custom naming

**Special Provider Path Exceptions**:

| Provider | Path Pattern | Init Type | Notes |
|----------|--------------|-----------|-------|
| Standard | `.tool/commands/spectr/proposal.md` | `SlashCommandsInit` | Subdir |
| Antigravity | `.agent/workflows/spectr-proposal.md` | `PrefixedSlash` | Proj |
| Codex | `~/.codex/prompts/spectr-proposal.md` | `HomePrefixedInit` | Home |

These are intentional design exceptions, not bugs. Antigravity and Codex
intentionally share `AGENTS.md` as their config file - both use the same
file with same markers. Deduplication handles this - only one
ConfigFileInitializer runs.

**Deduplication**: When multiple providers share an initializer with same
config, run once.

**ConfigFileInitializer Marker Handling**:

The ConfigFileInitializer uses `<!-- spectr:start -->` and `<!-- spectr:end -->`
markers to update instruction files. **Marker matching is case-insensitive for
reading** (matches both uppercase and lowercase), but **always writes lowercase
markers**. This ensures behavioral equivalence with files created by older
versions.

```go
// Pseudocode for marker handling logic
// Note: findMarkerCaseInsensitive returns the index and length of the
// matched marker
func findMarkerCaseInsensitive(content string, marker string)
    (index int, length int) {
    lower := strings.ToLower(content)
    lowerMarker := strings.ToLower(marker)
    idx := strings.Index(lower, lowerMarker)
    if idx == -1 {
        return -1, 0
    }
    return idx, len(marker)
}

func updateWithMarkers(content, newContent string) (string, error) {
    // Always write lowercase markers
    startMarker := "<!-- spectr:start -->"
    endMarker := "<!-- spectr:end -->"

    // Case-insensitive search for existing markers
    startIdx, _ := findMarkerCaseInsensitive(content, startMarker)

    if startIdx == -1 {
        // No start marker - check for orphaned end marker
        // (case-insensitive)
        endIdx, _ := findMarkerCaseInsensitive(content, endMarker)
        if endIdx != -1 {
            return "", fmt.Errorf(
                "orphaned end marker at position %d without start marker",
                endIdx)
        }
        // No markers exist - append new block at end with lowercase markers
        return content + "\n\n" + startMarker + "\n" + newContent +
            "\n" + endMarker, nil
    }

    // Start marker found - look for end marker AFTER the start
    // (case-insensitive)
    searchFrom := startIdx + len(startMarker)
    endIdx, _ := findMarkerCaseInsensitive(content[searchFrom:], endMarker)

    if endIdx != -1 {
        // Normal case: both markers present and properly paired
        endIdx += searchFrom // Adjust to absolute position

        // Check for nested start marker before end (case-insensitive)
        nextStartIdx, _ := findMarkerCaseInsensitive(
            content[searchFrom:endIdx], startMarker)
        if nextStartIdx != -1 {
            return "", fmt.Errorf(
                "nested start marker at position %d before end marker at %d",
                searchFrom+nextStartIdx, endIdx)
        }

        before := content[:startIdx]
        after := content[endIdx+len(endMarker):]
        // Always write lowercase markers
        return before + startMarker + "\n" + newContent + "\n" + endMarker +
            after, nil
    }

    // Start marker exists but no end marker immediately after
    // Search for trailing end marker anywhere in the file
    // (case-insensitive)
    trailingEndIdx, _ := findMarkerCaseInsensitive(content[searchFrom:],
        endMarker)
    if trailingEndIdx != -1 {
        trailingEndIdx += searchFrom
        // Found end marker after start - use it
        before := content[:startIdx]
        after := content[trailingEndIdx+len(endMarker):]
        // Always write lowercase markers
        return before + startMarker + "\n" + newContent + "\n" + endMarker +
            after, nil
    }

    // Check for multiple start markers without end (case-insensitive)
    nextStartIdx, _ := findMarkerCaseInsensitive(content[searchFrom:],
        startMarker)
    if nextStartIdx != -1 {
        return "", fmt.Errorf(
            "multiple start markers at positions %d and %d without end markers",
            startIdx, searchFrom+nextStartIdx)
    }

    // No end marker anywhere after start - orphaned start marker
    // Replace everything from start marker onward with new block
    before := content[:startIdx]
    return before + startMarker + "\n" + newContent + "\n" + endMarker, nil
}
```text

**Marker Format**: Marker matching is **case-insensitive for reading**
(matches both `<!-- spectr:START -->` and `<!-- spectr:start -->`), but
**always writes lowercase** `<!-- spectr:start -->` and `<!-- spectr:end -->`
for consistency across all providers. This ensures behavioral equivalence with
files created by older versions.

**Edge Cases Handled**:

- **Missing markers**: Insert at end of file with markers
- **Orphaned end marker**: Return error (corrupted file)
- **Nested markers**: Return error (not supported)
- **Multiple start markers**: Return error (ambiguous)
- **Orphaned start marker**: Replace from start marker to end of file
- **Normal case**: Replace content between existing markers

**Rationale**:

- **Prevents duplicate blocks**: Don't append when start marker already
  exists
- **Error on corruption**: Fail fast for orphaned end, nested, or multiple
  starts
- **Recovers gracefully**: Create markers if missing, use trailing end if
  found
- **Consistent format**: Lowercase markers across all providers

### 4. Filesystem Abstraction

**Decision**: Use two filesystem instances to support both project-relative and
home directory paths.

```go
// Project-relative filesystem (for CLAUDE.md, .claude/commands/, etc.)
projectFs := afero.NewBasePathFs(afero.NewOsFs(), projectPath)

// Home filesystem (for ~/.config/tool/commands/, etc.)
// Uses os.UserHomeDir() and converts to afero.Fs
// MUST return error if home directory cannot be determined
homeDir, err := os.UserHomeDir()
if err != nil {
    return nil, fmt.Errorf("failed to get home directory: %w", err)
}
homeFs := afero.NewBasePathFs(afero.NewOsFs(), homeDir)

// Executor provides both to initializers
type ExecutorContext struct {
    ProjectFs afero.Fs
    HomeFs    afero.Fs
    Config    *Config
    Templates *TemplateManager
}
```text

**Rationale**:

- Project paths are cleaner and easier to test
- Home paths support tools like Codex that use ~/.codex/
- Initializers receive both filesystems and decide internally which to use
  based on their type (Home* initializers use homeFs)

### 5. File Change Detection

**Decision**: Each initializer returns an `ExecutionResult` containing the files
it created/updated. The executor merges results inline.

```go
// ExecutionResult contains results from initialization
// Note: Error is returned separately, not stored in this struct
type ExecutionResult struct {
    CreatedFiles []string // All files created
    UpdatedFiles []string // All files updated
}

// Fail-fast execution: stop on first error
var allCreated, allUpdated []string

for _, init := range allInitializers {
    result, err := init.Init(ctx, projectFs, homeFs, cfg, tm)
    if err != nil {
        // Fail fast: return immediately on first error
        // Files created before error remain on disk (no rollback)
        // Return partial results so user knows what was created
        return ExecutionResult{
            CreatedFiles: allCreated,
            UpdatedFiles: allUpdated,
        }, err
    }
    allCreated = append(allCreated, result.CreatedFiles...)
    allUpdated = append(allUpdated, result.UpdatedFiles...)
}

// All succeeded
return ExecutionResult{
    CreatedFiles: allCreated,
    UpdatedFiles: allUpdated,
}, nil
```text

**Rationale**: Explicit change tracking; initializers know what they create;
works in non-git projects; more testable; single result type simplifies the
API.

## Example: Claude Code Provider

```go
// internal/initialize/providers/claude.go
package providers

import (
    "context"

    "github.com/connerohnesorge/spectr/internal/domain"
)

// ClaudeProvider configures Claude Code with CLAUDE.md and
// .claude/commands/spectr/.
// No init() - registration happens in RegisterAllProviders().
type ClaudeProvider struct{}

func (p *ClaudeProvider) Initializers(ctx context.Context,
    tm *TemplateManager) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".claude/commands/spectr"),
        // TemplateRef passed directly - simpler API, template resolved at
        // provider construction
        NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer()),
        // Early binding: TemplateRef resolved at construction time
        // SlashCommandsInitializer creates Markdown files (.md)
        NewSlashCommandsInitializer(".claude/commands/spectr",
            map[domain.SlashCommand]domain.TemplateRef{
                domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
                domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
            }),
    }
}
```text

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

func (p *GeminiProvider) Initializers(ctx context.Context,
    tm *TemplateManager) []Initializer {
    return []Initializer{
        NewDirectoryInitializer(".gemini/commands/spectr"),
        // No config file for Gemini - uses TOML slash commands only
        // Early binding: TemplateRef resolved at construction time
        NewTOMLSlashCommandsInitializer(".gemini/commands/spectr",
            map[domain.SlashCommand]domain.TemplateRef{
                domain.SlashProposal: tm.TOMLSlashCommand(
                    domain.SlashProposal),
                domain.SlashApply: tm.TOMLSlashCommand(domain.SlashApply),
            }),
    }
}
```text

**Template Files by Initializer Type**:

- `SlashCommandsInitializer` / `HomeSlashCommandsInitializer` /
  `PrefixedSlashCommandsInitializer` /
  `HomePrefixedSlashCommandsInitializer` → `slash-proposal.md.tmpl`,
  `slash-apply.md.tmpl`
- `TOMLSlashCommandsInitializer` → `slash-proposal.toml.tmpl`,
  `slash-apply.toml.tmpl`

**TOML Template Structure** (Gemini only):

```toml
description = "Create a Spectr change proposal"
prompt = """
{{ .Content }}
"""
```text

Minimal structure with `description` and `prompt` fields only.

**TOML/Slash Command Update Behavior**:

- All slash command initializers (TOML and Markdown) always overwrite
  existing files
- This ensures idempotent behavior: running `spectr init` multiple times
  produces consistent results
- User modifications to slash command files will be lost on re-initialization
- This is intentional: slash commands are generated from templates and should
  reflect the current template content

### 8. Initializer Ordering (Documented Guarantee)

**Decision**: Initializers are sorted by type before execution. This is a
documented guarantee.

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
    case *DirectoryInitializer, *HomeDirectoryInitializer:
        return 1
    case *ConfigFileInitializer:
        return 2
    case *SlashCommandsInitializer, *HomeSlashCommandsInitializer,
        *PrefixedSlashCommandsInitializer,
        *HomePrefixedSlashCommandsInitializer, *TOMLSlashCommandsInitializer:
        return 3
    default:
        return 99
    }
}

// Note: Order within the same type category preserves provider priority
// order.
// Higher-priority providers (lower priority number) appear first within
// each type category.
// This ensures that when deduplication keeps the first occurrence, it's
// from the highest-priority provider.
```text

**Rationale**: Directories must exist before files can be written. This
ordering is implicit but guaranteed.

### 9. Initializer Deduplication

**Decision**: Deduplicate by initializer identity. Same initializer
configuration = run once. Keep **first** occurrence when duplicates are found.

**Provider Priority Handling**: Providers are iterated in priority order
(lowest priority number first). Initializers are then sorted by type using a
stable sort, which preserves provider order within each type category. When
deduplication keeps the first occurrence, it's always from the
highest-priority provider. For example, if Claude (priority 1) and Cline
(priority 7) both create CLAUDE.md, Claude's initializer is kept because
Claude appears first in the collection. The separate Home* types make
filesystem boundaries explicit:

```go
// deduplicatable is an optional interface for initializers that support
// deduplication.
// Initializers that implement this interface can be deduplicated based on
// their key.
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

// Execution flow:
// 1. Collect initializers from all providers (providers iterated by priority
//    order)
// 2. Sort by type priority (stable sort preserves provider order within type)
// 3. Deduplicate (keep first occurrence = highest-priority provider wins)
// 4. Execute in order (fail-fast on error)

// Path normalization: Paths are normalized before generating deduplication
// keys.
// - Remove trailing slashes
// - Normalize path separators
// - Clean path (remove redundant `.` and `..`)
func (d *DirectoryInitializer) dedupeKey() string {
    return fmt.Sprintf("DirectoryInitializer:%s", filepath.Clean(d.path))
}

// Example dedupeKey implementations (type name encodes filesystem and
// format):
// DirectoryInitializer:                 "DirectoryInitializer:.claude/commands/spectr"
// HomeDirectoryInitializer:             "HomeDirectoryInitializer:.codex/prompts"
// ConfigFileInitializer:                "ConfigFileInitializer:CLAUDE.md"
// SlashCommandsInitializer:             "SlashCommandsInitializer:.claude/commands/spectr"
// HomeSlashCommandsInitializer:         "HomeSlashCommandsInitializer:.codex/prompts"
// PrefixedSlashCommandsInitializer:     "PrefixedSlashCommandsInitializer:.agent/workflows:spectr-"
// HomePrefixedSlashCommandsInitializer: "HomePrefixedSlashCommandsInitializer:.codex/prompts:spectr-"
// TOMLSlashCommandsInitializer:         "TOMLSlashCommandsInitializer:.gemini/commands/spectr"
```text

**Example**: If Claude Code and Cline both return `ConfigFileInitializer{path:
"CLAUDE.md"}`, only one runs.

**Rationale**:

- Type-based keys: `HomeDirectoryInitializer` and `DirectoryInitializer` have
  different type names, preventing accidental deduplication across filesystem
  boundaries
- Optional interface: Initializers that don't need deduplication don't have to
  implement it
- Simple and covers the common case (multiple providers sharing same file)

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking change for users | Clear migration docs; `spectr init` re-run |
| Path-based dedup misses edge cases | Simple first; can enhance if needed |
| No rollback on failure | Clear error reporting; users can re-run init |

## Migration Plan

### Zero Technical Debt Policy

This migration follows a **zero technical debt** policy - no compatibility
shims, no deprecated functions kept around. Clean break, complete removal of old
code.

### Migration Steps

1. Implement new provider system (new files, new types)
2. Migrate all 15 providers in-place (delete old code, write new code in same
  file)
3. **Completely remove** old registration system - no `Register(p Provider)`
  function kept around
4. Update docs to explain re-initialization requirement
5. No rollback needed - old configs continue to work

### No Compatibility Shims

**We will NOT do this:**

```go
// BAD: Deprecated compatibility shim that silently swallows calls
func Register(_ any) {
    // No-op to prevent compilation errors during migration
}
```text

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

The migration happens in a single change. There is no "partial migration" state
where some providers use init() and others don't. The tasks are ordered to
ensure:

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
```text

**STOP.** This violates the zero technical debt policy. Instead:

1. **Delete the old Register() function entirely** (task 7.6)
2. **Delete all init() functions** that call it (tasks 5.1-5.16)
3. **Implement RegisterAllProviders()** in one location (task 4.6)
4. Let the **compiler enforce migration** with clear errors

The only valid reason to keep deprecated code would be if external code depends
on it. Since all providers are internal to this codebase, there is no valid
reason to keep compatibility shims.

### 6. TemplateManager Integration

**Decision**: Initializers receive `*TemplateManager` instead of
implementing `TemplateRenderer`.

```go
type Initializer interface {
    Init(ctx context.Context, projectFs, homeFs afero.Fs, cfg *Config,
        tm *TemplateManager) (ExecutionResult, error)
    IsSetup(projectFs, homeFs afero.Fs, cfg *Config) bool
}
```text

**TemplateManager loads templates from both locations:**

```go
// NewTemplateManager creates a new template manager with all embedded
// templates loaded.
// It merges templates from:
// 1. internal/initialize/templates (main templates: AGENTS.md,
//    instruction-pointer.md)
// 2. internal/domain (slash command templates: slash-proposal.md,
//    slash-apply.md)
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
    // If duplicate template names exist, last-wins precedence applies
    for _, t := range domainTmpl.Templates() {
        if _, err := mainTmpl.AddParseTree(t.Name(), t.Tree); err != nil {
            return nil, fmt.Errorf("failed to merge template %s: %w",
                t.Name(), err)
        }
    }
    // Note: Template name collisions result in last-wins (later template
    // overwrites earlier)

    return &TemplateManager{
        templates: mainTmpl,
    }, nil
}
```text

**Rationale**:

- Reuses existing `TemplateManager` from `internal/initialize/templates.go`
- Merges templates from both `internal/initialize` and `internal/domain`
  packages
- Avoids duplicating template rendering logic in each initializer
- Simpler than the old `TemplateRenderer` interface pattern
- Domain package remains dependency-free (only exposes embed.FS, doesn't use it)

**Alternatives considered:**

- Each initializer implements own template rendering - More code duplication
- Pass templates as strings - Less flexible, harder to maintain
- Keep all templates in `internal/initialize` - Creates import cycle for
  domain types

### 7. Type-Safe Template Selection

**Decision**: Use typed template references instead of raw strings for
compile-time safety.

**Problem**: The original design had raw string inputs like:

```go
// Unsafe: typo "instrction-pointer" would fail at runtime
NewConfigFileInitializer("CLAUDE.md", "instruction-pointer")
```text

**Solution**: Define typed template accessors on TemplateManager:

```go
// TemplateRef is a type-safe reference to a parsed template.
// It's a lightweight handle without rendering logic.
// Rendering is performed by TemplateManager.
type TemplateRef struct {
    Name     string              // template file name (e.g., "instruction-pointer.md.tmpl")
    Template *template.Template  // pre-parsed template
}

// TemplateManager exposes type-safe accessors for each template.
// Rendering is performed via tm.Render(templateRef.Name, ctx).
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

// SlashCommand returns a Markdown template reference for the given slash
// command type.
// Used by SlashCommandsInitializer, HomeSlashCommandsInitializer, and
// PrefixedSlashCommandsInitializer.
func (tm *TemplateManager) SlashCommand(cmd domain.SlashCommand)
    domain.TemplateRef {
    names := map[domain.SlashCommand]string{
        domain.SlashProposal: "slash-proposal.md.tmpl",
        domain.SlashApply:    "slash-apply.md.tmpl",
    }
    return domain.TemplateRef{
        Name:     names[cmd],
        Template: tm.templates,
    }
}

// TOMLSlashCommand returns a TOML template reference for the given slash
// command type.
// Used by TOMLSlashCommandsInitializer (Gemini only).
func (tm *TemplateManager) TOMLSlashCommand(cmd domain.SlashCommand)
    domain.TemplateRef {
    names := map[domain.SlashCommand]string{
        domain.SlashProposal: "slash-proposal.toml.tmpl",
        domain.SlashApply:    "slash-apply.toml.tmpl",
    }
    return domain.TemplateRef{
        Name:     names[cmd],
        Template: tm.templates,
    }
}
```text

**Updated initializer constructors**:

```go
// ConfigFileInitializer receives TemplateRef directly
// TemplateRef is resolved at provider construction time when
// Initializers() is called
func NewConfigFileInitializer(path string, template TemplateRef) Initializer {
    return &ConfigFileInitializer{
        path:     path,
        template: template,
    }
}

// Usage - compile-time checked, TemplateRef passed directly:
// Provider.Initializers(ctx, tm) receives TemplateManager
NewConfigFileInitializer("CLAUDE.md", tm.InstructionPointer())

// SlashCommandsInitializer receives map of SlashCommand to TemplateRef
// (early binding)
// TemplateRef is resolved at provider construction time when
// Initializers() is called
func NewSlashCommandsInitializer(dir string,
    commands map[SlashCommand]TemplateRef) Initializer {
    return &SlashCommandsInitializer{
        dir:      dir,
        commands: commands,
    }
}

// TOMLSlashCommandsInitializer receives map of SlashCommand to TemplateRef
// (early binding)
func NewTOMLSlashCommandsInitializer(dir string,
    commands map[SlashCommand]TemplateRef) Initializer {
    return &TOMLSlashCommandsInitializer{
        dir:      dir,
        commands: commands,
    }
}

// Usage - compile-time checked, TemplateRef passed directly:
// Provider.Initializers(ctx, tm) receives TemplateManager
NewSlashCommandsInitializer(".claude/commands/spectr",
    map[domain.SlashCommand]domain.TemplateRef{
        domain.SlashProposal: tm.SlashCommand(domain.SlashProposal),
        domain.SlashApply:    tm.SlashCommand(domain.SlashApply),
    })

NewTOMLSlashCommandsInitializer(".gemini/commands/spectr",
    map[domain.SlashCommand]domain.TemplateRef{
        domain.SlashProposal: tm.TOMLSlashCommand(domain.SlashProposal),
        domain.SlashApply:    tm.TOMLSlashCommand(domain.SlashApply),
    })
```text

**Benefits**:

- Compile-time errors for invalid template names (typos caught by compiler)
- IDE autocomplete for available templates
- Refactoring-safe: renaming a template accessor updates all usages
- Clear documentation of available templates through method signatures

**Alternatives considered:**

- String constants (`const TemplateInstructionPointer =
  "instruction-pointer"`) - Still strings, can be used incorrectly
- Template name validation at startup - Runtime error, not compile-time
- Passing `*template.Template` directly - Leaks implementation detail,
  harder to mock in tests

### 10. TemplateContext Creation

**Decision**: TemplateContext is derived from Config.SpectrDir in the executor.

```go
func templateContextFromConfig(cfg *Config) domain.TemplateContext {
    return domain.TemplateContext{
        BaseDir:     cfg.SpectrDir,
        SpecsDir:    cfg.SpecsDir(),
        ChangesDir:  cfg.ChangesDir(),
        ProjectFile: cfg.ProjectFile(),
        AgentsFile:  cfg.AgentsFile(),
    }
}
```text

This ensures all template variables are derived from a single source of
truth (Config.SpectrDir).

## Resolved Questions

| Question | Decision |
|----------|----------|
| Change detection approach? | ExecutionResult return value |
| Initializer ordering? | Implicit by type (Directory → ConfigFile → Slash) |
| Partial failure handling? | Fail-fast: stop on first error, partial results |
| Template variable location? | Derived from SpectrDir via methods |
| Home paths support? | Separate types: `Home*Initializer` variants |
| Deduplication key? | By type name + path via `deduplicatable` interface |
| Template selection type safety? | TemplateRef passed directly to initializers |
| Import cycle resolution? | New `internal/domain` package with shared types |
| Registration error handling? | Explicit `RegisterAllProviders()` in cmd/init |
| Initializer methods? | Two: `Init()` and `IsSetup()`, both receive fs |
| TemplateRef field visibility? | Public fields for external package access |
| TOML template support? | Separate `TOMLSlashCommandsInitializer` type |
| Frontmatter structure? | Minimal: only `description` field |
| Marker search algorithm? | Case-insensitive read; lowercase write |
| os.UserHomeDir() failure? | Fail initialization entirely |
| Directory already exists? | Silent success; don't report in UpdatedFiles |
| Re-run behavior? | Always re-run; idempotent initializers |
| IsSetup() purpose? | For setup wizard UI only; Init() always runs |
| ConfigFileInitializer idempotency? | Content between markers replaced |
| TemplateContext creation? | Derived from Config.SpectrDir in executor |
| Template collision? | Last-wins precedence, silent |
| SlashCommand filename? | Use `SlashCommand.String()` + extension |
| Provider priority constraints? | Unique positive integers; gaps allowed |
| Slash command update behavior? | Always overwrite; user mods lost on re-init |

## Future Considerations (Out of Scope)

- `spectr init --dry-run` to preview changes without applying
- New instruction file support for Gemini, Cursor, Aider, OpenCode (separate
  proposal)
