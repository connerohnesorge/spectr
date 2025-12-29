# Design: Add Codex CLI Provider Support

## Context

Codex CLI is OpenAI's open-source agentic coding tool. Unlike all existing Spectr providers (Claude Code, Gemini, Aider, etc.), Codex stores custom prompts in a **global** location (`~/.codex/prompts/`) rather than project-local paths (`.claude/commands/`, `.gemini/commands/`, etc.).

This is the first provider requiring global path support, which introduces new architectural considerations.

## Goals

- Add Codex CLI as a supported provider
- Support global path installation (`~/.codex/prompts/spectr/`)
- Maintain backward compatibility with existing project-local providers
- Keep the provider interface clean and extensible

## Non-Goals

- Adding project-local support for Codex (Codex doesn't support it)
- Modifying how existing providers work
- Adding complex global/local path negotiation logic

## Decisions

### Decision 1: Extend BaseProvider with Global Path Support

**What**: Add a `globalPath` flag or method to distinguish global vs project-local providers.

**Why**: The current `Configure(projectPath, spectrDir string)` assumes project-local paths. For Codex, we need to resolve `~/.codex/prompts/` independent of `projectPath`.

**Alternatives considered**:

1. **Override Configure() in CodexProvider** - Would work but duplicates path resolution logic
2. **Add IsGlobal() method to Provider interface** - Clean but adds interface complexity for one provider
3. **Use absolute paths in command paths** - Simple: if path starts with `~/` or `/`, treat as global

**Decision**: Option 3 - Use absolute paths. If `proposalPath` starts with `~` or `/`, don't prepend `projectPath`. This is minimally invasive and self-documenting.

### Decision 2: Use `spectr/` Subdirectory in Global Prompts

**What**: Install to `~/.codex/prompts/spectr/` not `~/.codex/prompts/`.

**Why**: Prevents naming conflicts with user's existing prompts and clearly namespaces Spectr commands.

**Pattern**: `/prompts:spectr/proposal`, `/prompts:spectr/apply`, `/prompts:spectr/sync`

### Decision 3: Home Directory Expansion

**What**: Expand `~` to user's home directory at runtime.

**Why**: Storing `~/.codex/prompts/spectr/proposal.md` as literal path allows clean configuration while resolving to actual path at execution time.

**Implementation**: Use `os.UserHomeDir()` when path starts with `~/`.

### Decision 4: IsConfigured() for Global Paths

**What**: `IsConfigured(projectPath)` must handle global paths correctly.

**Why**: Currently checks `filepath.Join(projectPath, relPath)`. For global paths, should check absolute path instead.

### Decision 5: Priority Placement

**What**: Assign Codex priority 10 (between Cursor at 9 and Aider at 11).

**Why**: Places Codex among similar CLI-focused tools.

## Implementation Approach

### Phase 1: Extend Path Resolution

1. Add `expandPath(path string) string` helper to resolve `~` to home dir
2. Modify `BaseProvider.configureSlashCommands()` to check if path is absolute/global
3. Modify `BaseProvider.IsConfigured()` similarly
4. Modify `BaseProvider.GetFilePaths()` to return expanded paths for display

### Phase 2: Add Codex Provider

1. Create `internal/init/providers/codex.go`
2. Add `PriorityCodex = 10` to constants.go
3. Use global paths: `~/.codex/prompts/spectr/{proposal,sync,apply}.md`
4. Standard markdown frontmatter with `description:` field

### Phase 3: Update Tests

1. Add unit tests for `expandPath()`
2. Add tests for global path handling in BaseProvider
3. Add tests for CodexProvider

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Global install affects all projects | Use `spectr/` subdirectory to namespace; document behavior |
| `~` expansion varies by platform | Use `os.UserHomeDir()` which handles cross-platform |
| User may not have `~/.codex/` dir | Create parent directories as needed (existing behavior) |

## Migration Plan

No migration needed - this is additive. Existing providers unchanged.

## Open Questions

None remaining after user clarification.
