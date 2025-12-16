## Context

The validation system currently resolves base specs using the local filesystem path:
```go
baseSpecPath := filepath.Join(spectrRoot, "specs", capability, "spec.md")
```

This works correctly for local development but fails when the target branch (e.g., `main`) has specs that differ from the local working directory. The `pr archive` command creates a git worktree from `origin/main` and copies the change there, causing validation to run against `main`'s specs instead of local specs.

**Stakeholders**: Users who have long-running feature branches or work alongside other contributors whose changes get merged first.

## Goals / Non-Goals

**Goals:**
- Catch validation errors locally that would fail in `pr archive`
- Provide clear error messages explaining the local vs remote discrepancy
- Maintain backward compatibility (no breaking changes to existing workflows)

**Non-Goals:**
- Automatic resolution of spec conflicts
- Making `pr archive` use local specs (the current behavior is correct for archive)
- Supporting arbitrary remote refs (just main/default branch initially)

## Decisions

### Decision 1: Add `--base-branch` flag to validate command

**What**: Add an optional `--base-branch` flag that causes validation to read base specs from the specified branch instead of the local filesystem.

**Why**: This lets users validate as-if they were archiving, catching issues before the `pr archive` workflow starts.

**Alternatives considered**:
- Always validate against main: Too restrictive, breaks local development workflow
- Auto-detect and validate both: Complex, doubles validation time, unclear which result to trust
- Only fix the error message: Doesn't prevent the issue, just explains it better

### Decision 2: Pre-flight validation in PR workflow

**What**: Before creating the worktree, run validation against the target branch specs to fail fast with a better error message.

**Why**: The current error appears deep in the worktree/archive workflow with a path like `/tmp/spectr-pr-*`. Pre-flight validation provides the error in the context of the user's working directory.

**Implementation**: Read base spec content using `git show <branch>:spectr/specs/<capability>/spec.md` and pass it to the validator.

### Decision 3: Use git show for branch file access

**What**: Use `git show <ref>:<path>` to read file contents from other branches without creating a worktree.

**Why**: Lightweight, doesn't require filesystem changes, widely supported git operation.

**Alternatives considered**:
- Sparse checkout: More complex setup for one-off reads
- Full clone to temp dir: Wasteful for reading a few files
- git archive: Requires extracting, more complex

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Git command failures (branch doesn't exist) | Graceful fallback to local validation with warning |
| Performance overhead of git show per spec | Cache results or batch reads; base specs are small/few |
| Confusing dual validation (local vs remote) | Clear messaging about which validation mode is active |

## Migration Plan

1. Implement `--base-branch` flag (backward compatible)
2. Add pre-flight validation to `pr archive` (transparent improvement)
3. Update documentation and error messages
4. No breaking changes - existing workflows continue to work

## Open Questions

- Should `--base-branch` default to detecting the target PR branch, or require explicit specification?
- Should we add a `validate --pre-archive <change-id>` convenience command that auto-detects the target branch?
