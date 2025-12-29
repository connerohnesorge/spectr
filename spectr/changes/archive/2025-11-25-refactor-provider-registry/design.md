## Context

The current tool registry (`internal/init/`) manages AI CLI tool configurations through:

- `tool_definitions.go`: Global maps (`toolConfigs`, `slashToolConfigs`) with hardcoded configuration
- `registry.go`: `ToolRegistry` struct wrapping a map of `ToolDefinition` pointers
- `configurator.go`: `GenericConfigurator` that reads from global maps
- Separate "config" and "slash" tool entries with a mapping between them

Adding a new provider (e.g., Gemini CLI with TOML-based commands) requires:

1. Adding constants to `tool_definitions.go`
2. Adding entries to multiple global maps
3. Potentially modifying `configurator.go` for format differences

This proposal introduces a Go-idiomatic interface-driven pattern with **one interface per tool**.

## Goals / Non-Goals

**Goals:**

- Define a single `Provider` interface per tool (not separate config/slash)
- Create one Go file per provider under `internal/init/providers/`
- Each provider handles both its instruction file AND slash commands
- Support heterogeneous command formats (markdown, TOML)
- Make adding new providers a single-file addition

**Non-Goals:**

- Changing the CLI user experience (same commands, same flags)
- Supporting runtime provider discovery (compile-time registration is sufficient)
- Supporting user-defined providers (out of scope for this change)

## Decisions

### Decision: Single Provider Interface Per Tool

```go
// Provider represents an AI CLI tool (Claude Code, Gemini, Cline, etc.)
type Provider interface {
    // ID returns the unique provider identifier (kebab-case)
    ID() string
    // Name returns the human-readable provider name
    Name() string
    // Priority returns display order (lower = higher priority)
    Priority() int

    // ConfigFile returns the instruction file path (e.g., "CLAUDE.md"), empty if none
    ConfigFile() string
    // SlashDir returns the slash commands directory (e.g., ".claude/commands"), empty if none
    SlashDir() string
    // CommandFormat returns Markdown or TOML for slash command files
    CommandFormat() CommandFormat

    // Configure applies all configuration (instruction file + slash commands)
    Configure(projectPath, spectrDir string) error
    // IsConfigured checks if the provider is fully configured
    IsConfigured(projectPath string) bool
}

type CommandFormat int
const (
    FormatMarkdown CommandFormat = iota
    FormatTOML
)
```

**Rationale:** One provider = one tool. Claude Code handles both CLAUDE.md and .claude/commands/. No separate "slash provider" needed - simpler and more cohesive.

### Decision: Per-Provider Files

Create `internal/init/providers/` directory with:

- `provider.go` - Interface definition and base helpers
- `registry.go` - Global registry and registration functions
- `claude.go` - Claude Code (CLAUDE.md + .claude/commands/)
- `gemini.go` - Gemini CLI (~/.gemini/commands/ with TOML)
- `cline.go` - Cline (CLINE.md + .clinerules/commands/)
- `cursor.go`, `copilot.go`, `aider.go`, etc.

**Rationale:** Each provider file is self-contained. Adding Gemini support means adding one file.

### Decision: Registration via init()

```go
// In providers/claude.go
func init() {
    Register(&ClaudeProvider{})
}
```

**Rationale:** Standard Go pattern (see `database/sql`, `image` package).

### Decision: TOML Command Format Support

Gemini uses TOML for custom commands:

```toml
# ~/.gemini/commands/spectr-proposal.toml
description = "Scaffold a new Spectr change and validate strictly."
prompt = """
...prompt content...
"""
```

The `GeminiProvider` returns `FormatTOML` from `CommandFormat()` and generates TOML files in `Configure()`.

**Rationale:** Format-specific logic stays in provider implementation.

### Alternatives Considered

1. **Separate config/slash providers**: Rejected - adds unnecessary complexity
2. **Keep global maps, add Gemini entries**: Rejected - doesn't scale
3. **Factory pattern with type switches**: Rejected - not as extensible

## Risks / Trade-offs

- **Risk:** Breaking existing tests that reference `toolConfigs` directly
  - Mitigation: Update tests to use new `Registry` API

- **Trade-off:** More files to maintain vs. cleaner separation
  - Accepted: Clarity and extensibility benefits outweigh file count

## Migration Plan

1. Create `Provider` interface in `internal/init/providers/provider.go`
2. Create `Registry` in `internal/init/providers/registry.go`
3. Implement existing providers one at a time, starting with Claude
4. Update `executor.go` to use new registry
5. Remove deprecated code (`tool_definitions.go` globals, `ToolRegistry`)
6. Add Gemini provider
7. Run full test suite

## Open Questions

- Should TOML template rendering reuse `TemplateManager` or have its own implementation?
