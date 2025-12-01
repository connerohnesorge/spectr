# Design: `--pr` Flag for Archive Command with Git Worktree Isolation

## Context

The `spectr archive` command performs multiple operations:
1. Validates the change proposal
2. Checks task completion
3. Merges delta specs into main specs (unless `--skip-specs`)
4. Moves the change directory to `spectr/changes/archive/YYYY-MM-DD-<change-id>/`

After these operations complete, users typically want to commit and create a PR for team review. The previous `--pr` implementation operated directly on the user's working directory, which risked including uncommitted changes and caused confusion about the repository state.

This design uses **git worktrees** to provide complete isolation, executing the archive in a separate working tree on a fresh branch.

## Goals

- **Primary**: Provide complete isolation using git worktrees - never modify the user's main working directory
- **Primary**: Automate the branch → archive → commit → push → PR workflow atomically
- **Secondary**: Support multiple git hosting platforms (GitHub, GitLab, Gitea, Forgejo, Bitbucket)
- **Secondary**: Clean up worktrees automatically, even on failure

## Non-Goals

- Modify the user's current working directory or checkout state
- Support git operations without the `git` CLI (no libgit2)
- Auto-merge or handle PR review workflows
- Support running `--pr` without an actual archive operation

## Decisions

### 1. Worktree-Based Isolation

**Decision**: Use `git worktree` to create an isolated environment for the archive operation.

**Workflow**:
```bash
# 1. Create worktree on new branch
git worktree add /tmp/spectr-archive-<uuid> -b archive-<change-id> origin/main

# 2. Execute archive within worktree
cd /tmp/spectr-archive-<uuid>
spectr archive <change-id> --yes

# 3. Stage and commit
git add spectr/
git commit -m "[message]"

# 4. Push and create PR
git push -u origin archive-<change-id>
gh pr create ...

# 5. Cleanup worktree
cd -
git worktree remove /tmp/spectr-archive-<uuid>
```

**Rationale**:
- User's working directory is completely untouched
- No risk of including uncommitted changes
- Branch is based on `origin/main`, not local state
- Failed operations don't pollute the repository
- Worktrees are lightweight and fast (shared objects)

**Requirements**:
- Git >= 2.5 (worktree support)
- Clean base branch on remote (typically `main` or `master`)

### 2. Git Hosting Platform Detection

**Decision**: Detect platform from `origin` remote URL and select appropriate CLI tool.

**Detection Algorithm**:
```
URL Pattern                    → Platform    → CLI Tool
─────────────────────────────────────────────────────────
github.com                     → GitHub      → gh
gitlab.com OR has "gitlab"     → GitLab      → glab
gitea OR forgejo               → Gitea       → tea
bitbucket.org OR bitbucket     → Bitbucket   → (manual URL)
ssh://git@<custom>:...         → Unknown     → Error with guidance
```

**Implementation**:
```go
func DetectPlatform(remoteURL string) (Platform, error) {
    // Parse URL (handle both HTTPS and SSH formats)
    // Match against known patterns
    // Return platform enum and CLI tool name
}
```

**Rationale**:
- Single source of truth for platform detection
- Extensible for future platforms
- Clear error messages for unsupported platforms

### 3. Branch Naming Convention

**Decision**: Create branch with name `archive-<change-id>`.

**Examples**:
- `archive-add-user-auth`
- `archive-refactor-init-package-rename`

**Rationale**:
- Clearly indicates branch purpose
- Follows Spectr's kebab-case convention
- Unlikely to conflict with feature branches

**Conflict Handling**:
- If branch exists remotely: Error with message to delete or use different name
- If branch exists locally: Auto-handled by worktree (creates new branch)

### 4. Worktree Location and Naming

**Decision**: Create worktrees in system temp directory with UUID suffix.

**Pattern**: `{os.TempDir()}/spectr-archive-<uuid>/`

**Examples**:
- `/tmp/spectr-archive-a1b2c3d4/` (Linux/macOS)
- `C:\Users\...\AppData\Local\Temp\spectr-archive-a1b2c3d4\` (Windows)

**Rationale**:
- Temp directory is cleaned up by OS eventually
- UUID prevents conflicts between concurrent operations
- Predictable pattern aids debugging

### 5. Base Branch Selection

**Decision**: Base the archive branch on `origin/main` (or `origin/master` as fallback).

**Detection Order**:
1. Check if `origin/main` exists → Use `origin/main`
2. Check if `origin/master` exists → Use `origin/master`
3. Error: "Could not determine base branch. Remote has neither 'main' nor 'master'."

**Rationale**:
- Archives should be based on the current truth, not local state
- `main` is the modern default; `master` is legacy fallback
- Clear error if neither exists

### 6. Archive Execution in Worktree

**Decision**: Execute `spectr archive <change-id> --yes --skip-specs` within the worktree, then run spec merge manually.

**Workflow Detail**:
```bash
# In worktree:
cd /tmp/spectr-archive-<uuid>

# Run archive with --yes to auto-confirm (worktree is isolated)
# Use --skip-specs initially, then manually handle spec updates
spectr archive <change-id> --yes
```

**Rationale**:
- The archive operation modifies files in the worktree, not user's directory
- `--yes` is safe because worktree is isolated
- All prompts are bypassed for automation

**Self-Invocation Pattern**:
The `--pr` workflow invokes `spectr archive` as a subprocess within the worktree. This ensures:
- Same binary version is used
- All archive logic is reused
- No code duplication

### 7. Files to Stage

**Decision**: Stage the entire `spectr/` directory rather than individual files.

**Command**: `git add spectr/`

**Rationale**:
- Captures all archive-related changes (archived directory, updated specs)
- Simple and predictable
- Git handles deletions and moves automatically
- No risk of missing files

### 8. Commit Message Format

**Decision**: Use structured commit message with archive metadata.

**Template**:
```
archive(<change-id>): Archive completed change

Archived to: spectr/changes/archive/YYYY-MM-DD-<change-id>/

Spec operations applied:
+ {added} added
~ {modified} modified
- {removed} removed
→ {renamed} renamed

Generated by: spectr archive --pr
```

**Example**:
```
archive(refactor-init-package-rename): Archive completed change

Archived to: spectr/changes/archive/2025-12-01-refactor-init-package-rename/

Spec operations applied:
+ 0 added
~ 3 modified
- 0 removed
→ 0 renamed

Generated by: spectr archive --pr
```

**Rationale**:
- Conventional commit style with `archive()` scope
- Clear summary of what was archived
- Operation counts help reviewers understand scope
- Attribution aids debugging

### 9. PR Title and Body

**Decision**: Generate PR with structured title and Markdown body.

**PR Title**: `archive(<change-id>): Archive completed change`

**PR Body Template**:
```markdown
## Summary

Archived completed change: `<change-id>`

**Location**: `spectr/changes/archive/YYYY-MM-DD-<change-id>/`

## Spec Updates

| Operation | Count |
|-----------|-------|
| Added     | {N}   |
| Modified  | {N}   |
| Removed   | {N}   |
| Renamed   | {N}   |

**Updated capabilities**:
{list of capability names}

## Review Checklist

- [ ] Archived change structure is complete
- [ ] Spec deltas are accurate
- [ ] Merged spec content is correct

---
*Generated by `spectr archive --pr`*
```

**Rationale**:
- Consistent with commit message format
- Markdown table for quick scan
- Review checklist guides reviewers
- Attribution in footer

### 10. Platform-Specific PR Creation

**Decision**: Use platform CLI tools with consistent arguments.

**GitHub (`gh`)**:
```bash
gh pr create \
  --title "archive(<change-id>): Archive completed change" \
  --body-file /tmp/pr-body.md \
  --base main
```

**GitLab (`glab`)**:
```bash
glab mr create \
  --title "archive(<change-id>): Archive completed change" \
  --description "$(cat /tmp/pr-body.md)" \
  --target-branch main
```

**Gitea (`tea`)**:
```bash
tea pr create \
  --title "archive(<change-id>): Archive completed change" \
  --description "$(cat /tmp/pr-body.md)" \
  --base main
```

**Bitbucket**:
No official CLI; output manual URL:
```
PR creation not automated for Bitbucket.
Create manually at: https://bitbucket.org/<org>/<repo>/pull-requests/new?source=archive-<change-id>&dest=main
```

### 11. Error Handling Strategy

**Decision**: Fail fast with descriptive errors; always cleanup worktree.

**Error Hierarchy**:
```
Level 1: Pre-flight checks (before any git ops)
├── Not in git repository
├── No origin remote
├── Required CLI tool not installed
└── Base branch not found

Level 2: Worktree operations
├── Worktree creation failed
├── Archive execution failed
├── Commit failed
└── Push failed

Level 3: PR creation
└── PR CLI invocation failed
```

**Cleanup Guarantee**:
```go
defer func() {
    if worktreePath != "" {
        cleanupWorktree(worktreePath)
    }
}()
```

**Error Messages**:
- Include what failed and why
- Suggest remediation steps
- Include state information (e.g., "Branch was created and pushed, but PR creation failed")

### 12. Flag Compatibility

**Decision**: `--pr` works with `--skip-specs` but not with `--interactive` or `--no-validate`.

**Compatible**:
- `spectr archive <id> --pr` - Standard PR workflow
- `spectr archive <id> --pr --skip-specs` - PR without spec updates

**Incompatible**:
- `spectr archive --pr` (no change ID) - Error: `--pr` requires explicit change ID
- `spectr archive <id> --pr --interactive` - Error: `--pr` cannot be used with interactive mode
- `spectr archive <id> --pr --no-validate` - Error: `--pr` requires validation

**Rationale**:
- Interactive mode doesn't make sense for automated worktree workflow
- Validation is required to ensure only valid archives are PR'd
- Explicit change ID prevents ambiguity in automation

## Package Structure

```
internal/
├── archive/
│   ├── archiver.go       # Existing archive logic
│   ├── cmd.go             # Add PR flag
│   └── pr.go              # NEW: PR workflow orchestration
└── git/                   # NEW: Git operations package
    ├── platform.go        # Platform detection
    ├── platform_test.go
    ├── worktree.go        # Worktree management
    ├── worktree_test.go
    ├── pr.go              # PR CLI abstraction
    └── pr_test.go
```

## Alternatives Considered

### Alternative 1: Direct Branch Creation (Previous Approach)

**Rejected**: Risk of including uncommitted changes; modifies user's working directory.

### Alternative 2: Stash User Changes Before Archive

**Rejected**: Complex state management; risk of conflicts when unstashing.

### Alternative 3: Shallow Clone Instead of Worktree

**Rejected**: Slower than worktrees; requires network access; doesn't share objects.

### Alternative 4: Use go-git Library

**Rejected**: Heavy dependency; CLI approach is simpler and more debuggable.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Git < 2.5 doesn't support worktrees | Command fails | Check git version on startup; provide clear error |
| Concurrent archive --pr operations | Branch name conflict | UUID in worktree path; branch name from change-id |
| Worktree cleanup fails | Temp directory pollution | Log warning; user can manually clean |
| Network failure during push | Orphaned local branch | Provide recovery instructions |
| PR CLI authentication | PR creation fails | Check CLI auth status; provide login instructions |

## Open Questions

1. **Should we support `--base` flag to specify target branch?**
   - Recommendation: Not initially; default to main/master is sufficient

2. **Should we support `--draft` flag for draft PRs?**
   - Recommendation: Yes, add in initial implementation

3. **Should we delete the remote branch on PR merge?**
   - Recommendation: No; let platform handle via PR settings

4. **Should we support custom PR templates from `.github/PULL_REQUEST_TEMPLATE.md`?**
   - Recommendation: Not initially; let platform CLI handle this

---

**Decision Status**: Ready for implementation.
