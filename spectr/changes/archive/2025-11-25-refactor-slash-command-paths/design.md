## Context

The Provider interface currently uses `SlashDir()` to return a single directory path for slash commands. The BaseProvider then constructs file paths using a hardcoded pattern: `{slashDir}/spectr-{cmd}.md`. This creates coupling between the interface contract and file naming conventions.

Gemini CLI requires TOML files instead of markdown, forcing it to override 5+ methods just to change file extensions. Future providers may need even more flexibility (different directories per command, custom naming schemes, etc.).

## Goals / Non-Goals

**Goals:**

- Replace single `SlashDir()` with three explicit path methods
- Each method returns a complete relative path including filename
- Simplify provider implementations that need custom paths
- Maintain backward-compatible behavior for standard providers

**Non-Goals:**

- Changing command content generation (templates unchanged)
- Adding new command types
- Modifying instruction file handling (`ConfigFile()`)

## Decisions

### Decision: Three separate path methods

Replace `SlashDir() string` with:

```go
GetProposalCommandPath() string  // e.g., ".claude/commands/spectr-proposal.md"
GetArchiveCommandPath() string   // e.g., ".claude/commands/spectr-archive.md"
GetApplyCommandPath() string     // e.g., ".claude/commands/spectr-apply.md"
```

**Why separate methods instead of `GetCommandPath(cmd string)`:**

- Type safety: compile-time verification of command names
- Explicit: each command path is independently configurable
- No magic strings: callers don't need to know valid command names

**Why relative paths:**

- Consistent with `ConfigFile()` which returns relative paths
- Callers (Configure, IsConfigured, GetFilePaths) already have `projectPath`
- Easier to test (no absolute paths in expected values)

### Decision: Update BaseProvider with helper fields

Instead of a single `slashDir` field, use three path fields:

```go
type BaseProvider struct {
    id              string
    name            string
    priority        int
    configFile      string
    proposalPath    string  // e.g., ".claude/commands/spectr-proposal.md"
    archivePath     string  // e.g., ".claude/commands/spectr-archive.md"
    applyPath       string  // e.g., ".claude/commands/spectr-apply.md"
    commandFormat   CommandFormat
    frontmatter     map[string]string
}
```

**Alternative considered: Keep slashDir + add filename overrides**
Rejected because it adds complexity without solving the core problem (Gemini would still need directory + extension overrides).

### Decision: HasSlashCommands checks any path

`HasSlashCommands()` returns true if ANY of the three paths is non-empty:

```go
func (p *BaseProvider) HasSlashCommands() bool {
    return p.proposalPath != "" || p.archivePath != "" || p.applyPath != ""
}
```

This allows providers to support only some commands if needed.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Breaking change to interface | All providers are internal; update in single PR |
| More boilerplate per provider | Add helper function to generate standard paths |
| Forgotten path updates | Tests verify all registered providers have valid paths |

## Migration Plan

1. Update Provider interface (add 3 methods, remove SlashDir)
2. Update BaseProvider (replace slashDir with 3 path fields)
3. Add helper: `StandardCommandPaths(dir, ext string) (proposal, archive, apply string)`
4. Update each provider to use helper or custom paths
5. Simplify GeminiProvider (remove method overrides, use TOML paths directly)
6. Update tests to verify new methods
7. Remove deprecated `getSlashCommandPath` private method

## Open Questions

None - design decisions confirmed via user input.
