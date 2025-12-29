# Design: `spectr pr` Subcommand with Git Worktree Isolation

## Context

Spectr manages change proposals through a structured workflow: create proposal,
implement, archive. After archiving (or when sharing work-in-progress), users
typically want to create a PR for team review. The current workflow requires
multiple manual git commands:

```bash
git checkout -b archive-<change-id>
spectr archive <change-id> --yes
git add spectr/
git commit -m "..."
git push -u origin archive-<change-id>
gh pr create ...
```text

This is tedious and error-prone. Worse, it pollutes the user's working directory
if they have uncommitted changes.

This design introduces a new `spectr pr` command namespace with two subcommands
that use **git worktrees** for complete isolation.

## Goals

- **Primary**: Provide complete isolation using git worktrees - never modify
  user's main working directory
- **Primary**: Support both "archive and PR" and "copy and PR" workflows
- **Primary**: Automate the branch → operation → commit → push → PR workflow
  atomically
- **Secondary**: Support multiple git hosting platforms (GitHub, GitLab, Gitea,
  Forgejo, Bitbucket)
- **Secondary**: Clean up worktrees automatically, even on failure

## Non-Goals

- Modify the user's current working directory or checkout state
- Support git operations without the `git` CLI (no libgit2)
- Auto-merge or handle PR review workflows
- Replace the existing `spectr archive` command (this is additive)

## Command Structure

### Option 1: Flat Subcommands (Recommended)

```text
spectr pr archive <change-id> [flags]
spectr pr new <change-id> [flags]
```text

**Rationale**: Clear, verb-led subcommands. Matches existing CLI patterns in
Spectr.

### Option 2: Single Command with Mode Flag

```text
spectr pr <change-id> --mode=archive|new
```text

**Rejected**: Less discoverable, requires flag for common operation.

## Decisions

### 1. Worktree-Based Isolation

**Decision**: Use `git worktree` to create an isolated environment for all PR
operations.

**Workflow (archive)**:

```bash
# 1. Create worktree on new branch
git worktree add /tmp/spectr-pr-<uuid> -b spectr/<change-id> origin/main

# 2. Execute archive within worktree
cd /tmp/spectr-pr-<uuid>
spectr archive <change-id> --yes

# 3. Stage and commit
git add spectr/
git commit -m "[message]"

# 4. Push and create PR
git push -u origin spectr/<change-id>
gh pr create ...

# 5. Cleanup worktree
cd -
git worktree remove /tmp/spectr-pr-<uuid>
```text

**Workflow (new)**:

```bash
# 1. Create worktree on new branch
git worktree add /tmp/spectr-pr-<uuid> -b spectr/<change-id> origin/main

# 2. Copy change to worktree
cp -r spectr/changes/<change-id> /tmp/spectr-pr-<uuid>/spectr/changes/

# 3. Stage and commit
git add spectr/
git commit -m "[message]"

# 4. Push and create PR
git push -u origin spectr/<change-id>
gh pr create ...

# 5. Cleanup worktree
cd -
git worktree remove /tmp/spectr-pr-<uuid>
```text

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

**Decision**: Detect platform from `origin` remote URL and select appropriate
CLI tool.

**Detection Algorithm**:

```text
URL Pattern                    → Platform    → CLI Tool
─────────────────────────────────────────────────────────
github.com                     → GitHub      → gh
gitlab.com OR has "gitlab"     → GitLab      → glab
gitea OR forgejo               → Gitea       → tea
bitbucket.org OR bitbucket     → Bitbucket   → (manual URL)
ssh://git@<custom>:...         → Unknown     → Error with guidance
```text

**Implementation**:

```go
type Platform string

const (
    PlatformGitHub    Platform = "github"
    PlatformGitLab    Platform = "gitlab"
    PlatformGitea     Platform = "gitea"
    PlatformBitbucket Platform = "bitbucket"
    PlatformUnknown   Platform = "unknown"
)

type PlatformInfo struct {
    Platform Platform
    CLITool  string   // "gh", "glab", "tea", ""
    RepoURL  string   // For generating manual PR URLs
}

func DetectPlatform(remoteURL string) (PlatformInfo, error)
```text

**Rationale**:

- Single source of truth for platform detection
- Extensible for future platforms
- Clear error messages for unsupported platforms

### 3. Branch Naming Convention

**Decision**: Create branch with name `spectr/<change-id>`.

**Examples**:

- `spectr/add-user-auth`
- `spectr/refactor-init-package-rename`

**Rationale**:

- Clearly indicates branch purpose with `spectr/` prefix
- Follows Spectr's kebab-case convention
- Grouped under `spectr/` namespace to avoid conflicts

**Conflict Handling**:

- If branch exists remotely: Error with message to delete or use `--force`
- Add `--force` flag to delete existing branch and recreate

### 4. Worktree Location and Naming

**Decision**: Create worktrees in system temp directory with UUID suffix.

**Pattern**: `{os.TempDir()}/spectr-pr-<uuid>/`

**Examples**:

- `/tmp/spectr-pr-a1b2c3d4/` (Linux/macOS)
- `C:\Users\...\AppData\Local\Temp\spectr-pr-a1b2c3d4\` (Windows)

**Rationale**:

- Temp directory is cleaned up by OS eventually
- UUID prevents conflicts between concurrent operations
- Predictable pattern aids debugging

### 5. Base Branch Selection

**Decision**: Base the PR branch on `origin/main` (or `origin/master` as
fallback), with optional `--base` flag.

**Detection Order**:

1. If `--base <branch>` provided → Use specified branch
2. Check if `origin/main` exists → Use `origin/main`
3. Check if `origin/master` exists → Use `origin/master`
4. Error: "Could not determine base branch"

**Rationale**:

- PRs should be based on the current remote truth, not local state
- `main` is the modern default; `master` is legacy fallback
- `--base` allows targeting feature branches

### 6. Change Operation in Worktree

**Decision**: Execute appropriate change operation within the worktree.

**Archive Mode**:

```bash
cd /tmp/spectr-pr-<uuid>
spectr archive <change-id> --yes
```text

**New Mode**:

```bash
cd /tmp/spectr-pr-<uuid>
mkdir -p spectr/changes
cp -r <original>/spectr/changes/<change-id> spectr/changes/
```text

**Rationale**:

- Archive mode: Full archive workflow including spec merging
- New mode: Just copy the change proposal for review without archiving
- Both operate in isolated worktree

**Self-Invocation Pattern**:
For archive mode, the `spectr pr archive` command invokes `spectr archive` as a
subprocess. This ensures:

- Same binary version is used
- All archive logic is reused
- No code duplication

### 7. Files to Stage

**Decision**: Stage the entire `spectr/` directory.

**Command**: `git add spectr/`

**Rationale**:

- Captures all change-related modifications
- For archive: includes archived directory and updated specs
- For new: includes just the copied change
- Simple and predictable
- Git handles additions correctly

### 8. Commit Message Format

**Decision**: Use structured commit message with operation metadata.

**Archive Template**:

```text
spectr(archive): <change-id>

Archived to: spectr/changes/archive/YYYY-MM-DD-<change-id>/

Spec operations applied:
+ {added} added
~ {modified} modified
- {removed} removed
→ {renamed} renamed

Generated by: spectr pr archive
```text

**New Template**:

```text
spectr(proposal): <change-id>

Proposal for review: spectr/changes/<change-id>/

Generated by: spectr pr new
```text

**Rationale**:

- Conventional commit style with `spectr()` scope
- Clear summary of what operation was performed
- Operation counts (for archive) help reviewers understand scope
- Attribution aids debugging

### 9. PR Title and Body

**Decision**: Generate PR with structured title and Markdown body.

**PR Title**:

- Archive: `spectr(archive): <change-id>`
- New: `spectr(proposal): <change-id>`

**PR Body Template (Archive)**:

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
*Generated by `spectr pr archive`*
```text

**PR Body Template (New)**:

```markdown
## Summary

Proposal for review: `<change-id>`

**Location**: `spectr/changes/<change-id>/`

## Files

- `proposal.md` - Change overview
- `tasks.md` - Implementation checklist
- `specs/` - Delta specifications

## Review Checklist

- [ ] Proposal addresses the stated problem
- [ ] Delta specs are properly formatted
- [ ] Tasks are clear and actionable

---
*Generated by `spectr pr new`*
```text

### 10. Platform-Specific PR Creation

**Decision**: Use platform CLI tools with consistent arguments.

**GitHub (`gh`)**:

```bash
gh pr create \
  --title "<title>" \
  --body-file /tmp/pr-body.md \
  --base main
```text

**GitLab (`glab`)**:

```bash
glab mr create \
  --title "<title>" \
  --description "$(cat /tmp/pr-body.md)" \
  --target-branch main
```text

**Gitea (`tea`)**:

```bash
tea pr create \
  --title "<title>" \
  --description "$(cat /tmp/pr-body.md)" \
  --base main
```text

**Bitbucket**:
No official CLI; output manual URL:

```text
PR creation not automated for Bitbucket.
Create manually at: https://bitbucket.org/<org>/<repo>/pull-requests/new?source=spectr/<change-id>&dest=main
```text

### 11. Error Handling Strategy

**Decision**: Fail fast with descriptive errors; always cleanup worktree.

**Error Hierarchy**:

```text
Level 1: Pre-flight checks (before any git ops)
├── Not in git repository
├── No origin remote
├── Required CLI tool not installed
├── Base branch not found
└── Change does not exist

Level 2: Worktree operations
├── Worktree creation failed
├── Change operation failed (archive/copy)
├── Commit failed
└── Push failed

Level 3: PR creation
└── PR CLI invocation failed
```text

**Cleanup Guarantee**:

```go
defer func() {
    if worktreePath != "" {
        cleanupWorktree(worktreePath)
    }
}()
```text

**Error Messages**:

- Include what failed and why
- Suggest remediation steps
- Include state information (e.g., "Branch was created and pushed, but PR
  creation failed")

### 12. Flag Design

**Common Flags**:

- `--base <branch>` - Target branch for PR (default: auto-detect main/master)
- `--draft` - Create as draft PR
- `--force` - Delete existing remote branch if present
- `--dry-run` - Show what would be done without executing

**Archive-Specific Flags**:

- `--skip-specs` - Pass through to `spectr archive` to skip spec merging

**Rationale**:

- `--draft` is commonly needed for WIP PRs
- `--force` handles branch conflicts
- `--dry-run` aids debugging and validation
- `--skip-specs` provides flexibility for archive workflow

## Package Structure

```text
internal/
├── git/                   # NEW: Git operations package
│   ├── platform.go        # Platform detection
│   ├── platform_test.go
│   ├── worktree.go        # Worktree management
│   ├── worktree_test.go
│   └── doc.go
├── pr/                    # NEW: PR workflow package
│   ├── workflow.go        # PR workflow orchestration
│   ├── workflow_test.go
│   ├── templates.go       # Commit/PR message templates
│   ├── templates_test.go
│   └── doc.go
└── archive/               # Existing
    └── ...

cmd/
├── root.go                # Add PR command
└── pr.go                  # NEW: PR command with subcommands
```text

## Alternatives Considered

### Alternative 1: Flag on Archive (`spectr archive --pr`)

**Considered**: Add `--pr` flag to existing archive command.

**Trade-offs**:

- Pro: Single command for archive+PR
- Con: No support for "new" (non-archive) PR workflow
- Con: Overloads archive command with git concerns
- Con: Flag combinations become complex

**Decision**: Rejected in favor of dedicated namespace for cleaner separation.

### Alternative 2: Single `spectr pr` Without Subcommands

**Considered**: `spectr pr <change-id>` with `--archive` flag.

**Trade-offs**:

- Pro: Simpler command structure
- Con: "new" becomes the default, archive requires flag
- Con: Less discoverable via help

**Decision**: Rejected. Subcommands are more explicit and match existing
patterns.

### Alternative 3: Stash-Based Isolation

**Considered**: Stash user changes, operate, restore.

**Trade-offs**:

- Pro: Simpler than worktrees
- Con: Risk of conflicts on unstash
- Con: Modifies user's working directory (even temporarily)
- Con: Slower than worktrees

**Decision**: Rejected. Worktrees provide true isolation without touching user's
state.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Git < 2.5 no worktrees | Command fails | Check git version; error clearly |
| Concurrent `spectr pr` | Branch conflict | UUID in path; branch from ID |
| Worktree cleanup fails | Temp pollution | Log warning; manual cleanup docs |
| Network failure on push | Orphaned branch | Provide recovery instructions |
| PR CLI auth missing | PR create fails | Check auth; provide login steps |
| User in worktree | Inception | Detect and error with message |

## Open Questions

1. **Should we support `--reviewer` flag to add reviewers?**
   - Recommendation: Not initially; platform CLIs handle this differently

2. **Should we auto-open the PR URL in browser?**
   - Recommendation: Yes, with `--no-browser` to disable

3. **Should `spectr pr new` validate the change first?**
   - Recommendation: Yes, run `spectr validate <change-id>` and warn on failures

---

**Decision Status**: Ready for implementation.
