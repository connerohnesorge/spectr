## Context

The current provider architecture uses a monolithic pattern where each provider embeds `BaseProvider` and implements a fixed interface with methods like `ConfigFile()`, `GetProposalCommandPath()`, `HasConfigFile()`, etc. This creates several issues:

1. Adding new file types requires changing the Provider interface
2. `HasConfigFile()` and `HasSlashCommands()` are tech debt - providers explicitly define their capabilities
3. `BaseProvider` adds indirection without clear benefit
4. Configuration logic is duplicated across providers

**Stakeholders**: Developers adding new providers, contributors extending file type support

## Goals / Non-Goals

### Goals
- Providers declare their file initializers as a composable list
- Provider interface is minimal (6 methods)
- No `BaseProvider` - providers implement interface directly using helper functions
- Remove tech debt methods (`HasConfigFile()`, `HasSlashCommands()`, `ConfigFile()`, etc.)
- Single atomic migration of all 23 providers

### Non-Goals
- Changing the external CLI interface or user-facing behavior
- Adding actual Claude skill/agent implementations (follow-up proposals)
- Modifying the registry pattern or provider discovery mechanism
- Supporting runtime-configurable initializers

## Decisions

### Decision 1: FileInitializer Interface

```go
// FileInitializer creates or updates a single file for a provider.
type FileInitializer interface {
    // ID returns a unique identifier for this initializer.
    // Format: "{type}:{path}" e.g., "instruction:CLAUDE.md"
    ID() string

    // FilePath returns the relative path this initializer manages.
    // May contain ~ for home directory (expanded internally during Configure).
    FilePath() string

    // Configure creates or updates the file.
    Configure(projectPath string, tm TemplateRenderer) error

    // IsConfigured checks if the file exists and is properly configured.
    IsConfigured(projectPath string) bool
}
```

**Why**: Narrow interface focused on single-file operations. ID enables deduplication. Path expansion handled internally.

### Decision 2: Minimal Provider Interface (6 methods)

```go
type Provider interface {
    ID() string
    Name() string
    Priority() int
    Initializers() []FileInitializer
    IsConfigured(projectPath string) bool
    GetFilePaths() []string
}
```

**Why**:
- No `Configure()` method - use helper function instead
- No `HasConfigFile()` / `HasSlashCommands()` - tech debt removed
- No `ConfigFile()` / `GetProposalCommandPath()` / `GetApplyCommandPath()` - derived from initializers
- No `CommandFormat()` - implicit in initializer type

### Decision 3: Remove BaseProvider, Use Helper Functions

Instead of embedding `BaseProvider`, providers implement the interface directly and use exported helper functions:

```go
// Helper functions in providers/helpers.go
func ConfigureInitializers(
    inits []FileInitializer,
    projectPath string,
    tm TemplateRenderer,
) error

func AreInitializersConfigured(
    inits []FileInitializer,
    projectPath string,
) bool

func GetInitializerPaths(
    inits []FileInitializer,
) []string
```

**Why**: Simpler composition, no inheritance-like patterns, explicit control flow.

### Decision 4: Generic Slash Command Initializers

Slash command initializers take the command name as a parameter, not as separate types:

```go
// Generic - supports any command name
NewMarkdownSlashCommandInitializer(
    path string,        // ".claude/commands/spectr/proposal.md"
    commandName string, // "proposal", "apply", "sync", etc.
    frontmatter string,
)
```

**Why**: Flexible for future commands. No need to create new types for each command.

### Decision 5: Remove Wizard Filtering

The wizard previously used `HasConfigFile()` and `HasSlashCommands()` to filter providers. These methods are removed entirely.

**Why**: All providers shown equally. Users select what they want regardless of internal capabilities. Simpler UX.

### Decision 6: Path Expansion Handled Internally

`FilePath()` returns the template path (may include `~`). The `Configure()` method expands paths internally.

```go
func (i *InstructionFileInitializer) FilePath() string {
    return "CLAUDE.md"  // or "~/.codex/CODEX.md"
}

func (i *InstructionFileInitializer) Configure(projectPath string, tm TemplateRenderer) error {
    fullPath := expandPath(i.path, projectPath)  // handles ~ expansion
    // ... create/update file
}
```

**Why**: Consistent with current behavior. Initializers encapsulate path logic.

### Decision 7: Standalone Initializer IDs

IDs are standalone without provider context:

```go
ID() = "instruction:CLAUDE.md"
ID() = "markdown-cmd:.claude/commands/spectr/proposal.md"
ID() = "toml-cmd:.gemini/commands/spectr/proposal.toml"
```

**Why**: Path provides uniqueness. Simpler format.

### Decision 8: Executor Changes

The executor currently calls `provider.Configure()`. This changes to:

```go
// Before
err := provider.Configure(projectPath, spectrDir, tm)

// After
err := ConfigureInitializers(provider.Initializers(), projectPath, tm)
```

**Why**: Configuration logic centralized in helper function, not interface method.

## Example Provider Implementation

```go
// claude.go
package providers

func init() {
    Register(&ClaudeProvider{})
}

type ClaudeProvider struct{}

func (p *ClaudeProvider) ID() string       { return "claude-code" }
func (p *ClaudeProvider) Name() string     { return "Claude Code" }
func (p *ClaudeProvider) Priority() int    { return PriorityClaudeCode }

func (p *ClaudeProvider) Initializers() []FileInitializer {
    return []FileInitializer{
        NewInstructionFileInitializer("CLAUDE.md"),
        NewMarkdownSlashCommandInitializer(
            ".claude/commands/spectr/proposal.md",
            "proposal",
            StandardProposalFrontmatter,
        ),
        NewMarkdownSlashCommandInitializer(
            ".claude/commands/spectr/apply.md",
            "apply",
            StandardApplyFrontmatter,
        ),
    }
}

func (p *ClaudeProvider) IsConfigured(projectPath string) bool {
    return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *ClaudeProvider) GetFilePaths() []string {
    return GetInitializerPaths(p.Initializers())
}
```

## Risks / Trade-offs

### Risk: More Boilerplate Per Provider
Each provider now implements 6 methods directly instead of embedding BaseProvider.

**Mitigation**: The methods are simple one-liners. Helper functions do the heavy lifting.

### Risk: Large Single PR
All 23 providers migrated at once.

**Accepted**: Atomic change avoids mixed states. PR can be reviewed file-by-file.

### Trade-off: No Filtering in Wizard
Users see all providers regardless of capabilities.

**Accepted**: Simpler UX. Users know what tools they use.

## Migration Plan

### Single Atomic Migration
1. Add `FileInitializer` interface and implementations
2. Add helper functions
3. Update Provider interface (remove old methods)
4. Convert all 23 providers in same PR
5. Update executor to use helper function
6. Remove old BaseProvider and related code
7. Update wizard to remove filtering logic

### Rollback
Revert entire PR if issues arise. No partial states.

## Resolved Questions

1. **Should initializers support dependencies?**
   **Decision: No.** Order in list implies sequence. No enforcement.

2. **How should errors aggregate?**
   **Decision: Fail-fast.** Stop on first error.

3. **Should this include Claude skill/agent?**
   **Decision: No.** Architecture only. Follow-up proposals.

4. **Where should initializers live?**
   **Decision: Mixed.** Generic in own files, provider-specific in provider files.

5. **Type detection for HasConfigFile/HasSlashCommands?**
   **Decision: Remove these methods entirely.** Tech debt.

6. **BaseProvider?**
   **Decision: Remove.** Use helper functions.

7. **Slash command names?**
   **Decision: Generic.** Parameter, not separate types.

8. **Path expansion?**
   **Decision: Internal.** Initializer handles it.

9. **Wizard filtering?**
   **Decision: Remove.** All providers shown equally.

10. **Configure() on interface?**
    **Decision: No.** Helper function only.

11. **Initializer ID format?**
    **Decision: Standalone.** "{type}:{path}"

12. **Migration strategy?**
    **Decision: All at once.** Single atomic PR.
