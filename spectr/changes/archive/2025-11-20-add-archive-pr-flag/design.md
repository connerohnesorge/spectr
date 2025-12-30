# Design: `--pr` Flag for Archive Command

## Context

The `spectr archive` command performs multiple operations:

1. Validates the change proposal
2. Checks task completion
3. Merges delta specs into main specs (unless `--skip-specs`)
4. Moves the change directory to
  `spectr/changes/archive/YYYY-MM-DD-<change-id>/`

After these operations complete, users typically want to commit these changes
and create a PR for team review. Currently this requires manual git operations.
This design adds a `--pr` flag to automate PR creation after successful archive
completion.

The `spectr propose` command (from `add-propose-command`) provides a pattern for
git platform detection and PR creation that can be reused or adapted here.

## Goals

- **Primary**: Automate the branch → commit → push → PR workflow for archive
  operations
- **Secondary**: Reuse git detection and PR creation logic from `spectr propose`
  to maintain consistency
- **Secondary**: Only create PR after archive operation fully succeeds (atomic
  behavior)
- **Secondary**: Commit both the archived change and all updated specs in one
  atomic commit

## Non-Goals

- Modify git configuration or user settings
- Support creating PRs for failed archive operations
- Auto-merge or handle PR review workflows
- Support running `--pr` in isolation (must be combined with archive operation)

## Decisions

### 1. PR Workflow Activation Point

**Decision**: The `--pr` flag SHALL only execute after the entire archive
operation completes successfully.

**Rationale**:

- Archive operations must be atomic; partial archives should not be committed
- If validation fails, spec merging fails, or directory move fails, no git
  operations occur
- This ensures the PR only contains valid, complete archive results

**Implementation**:

- PR workflow is the final step in `Archiver.Archive()` method
- Only executes if all prior steps (validation, spec merging, move) succeed
- Git operations happen after the "Successfully archived" message

### 2. Code Reuse from Propose Command

**Decision**: Extract git platform detection and PR creation logic into a shared
package that both `propose` and `archive` commands can use.

**Options**:

- **Option A**: Create `internal/git/` package with shared PR logic
- **Option B**: Create `internal/pr/` package with PR-specific operations
- **Option C**: Duplicate logic in archive package

**Choice**: Option A - `internal/git/` package

**Rationale**:

- Both commands need identical platform detection (GitHub/GitLab/Gitea)
- Both commands need identical PR CLI invocation logic
- Shared package reduces duplication and maintenance burden
- `git` is a clear domain name encompassing both detection and operations
- Allows future commands to reuse PR creation logic

**Dependency Note**: This change depends on or should coordinate with
`add-propose-command` to avoid duplicating git/PR logic.

### 3. Branch Naming Convention

**Decision**: Create branch with name `archive-<change-id>`.

**Rationale**:

- Mirrors the `add-<change-id>` pattern from propose command
- Clearly indicates the branch purpose (archiving a completed change)
- Unlikely to conflict with proposal or feature branches
- Consistent with Spectr's kebab-case convention

**Examples**:

- `archive-user-auth` for archiving the `user-auth` change
- `archive-add-cli-commands` for archiving the `add-cli-commands` change

### 4. Commit Strategy

**Decision**: Stage ALL archive-related changes in a single atomic commit.

**Files to stage**:

1. Removal of original change directory (`spectr/changes/<change-id>/`) -
  detected automatically by git as a rename/delete
2. Addition of archived directory
  (`spectr/changes/archive/YYYY-MM-DD-<change-id>/`)
3. All updated spec files in `spectr/specs/` (unless `--skip-specs` was used)

**Implementation**:

```bash
git checkout -b archive-<change-id>
git add spectr/changes/archive/YYYY-MM-DD-<change-id>/
git add spectr/specs/  # All updated specs
# Original change dir removal is detected automatically
git commit -m "[commit message - see below]"
```text

**Rationale**:

- Archive operations should be atomic in version control
- Reviewers need to see both the archived change AND the spec updates together
- Staging `spectr/specs/` captures all modified specs without manual enumeration
- Git automatically detects directory moves/renames

### 5. Commit Message Format

**Decision**: Use structured commit message with archive summary and operation
counts.

**Format**:

```text
Archive: <change-id>

Completed change '<change-id>' archived to changes/archive/YYYY-MM-DD-<change-id>/

Spec operations applied:
+ <N> added
~ <N> modified
- <N> removed
→ <N> renamed

Change-Id: <change-id>
```text

**Rationale**:

- First line follows conventional commit style with "Archive:" prefix
- Body includes context about what was archived and where
- Spec operation summary provides reviewers with at-a-glance understanding of
  spec changes
- `Change-Id` trailer aids linking tools and automation
- If `--skip-specs` was used, omit the "Spec operations" section

**Example**:

```text
Archive: user-authentication

Completed change 'user-authentication' archived to changes/archive/2025-11-20-user-authentication/

Spec operations applied:
+ 3 added
~ 1 modified
- 0 removed
→ 0 renamed

Change-Id: user-authentication
```text

### 6. PR Title and Body Format

**Decision**: Use consistent PR title/body format that mirrors archive operation
results.

**PR Title**:

```text
Archive: <change-id>
```text

**PR Body Template**:

```markdown
## Archive Summary

Archived completed change: `<change-id>`

Location: `spectr/changes/archive/YYYY-MM-DD-<change-id>/`

## Spec Updates

<if specs were updated>
Spec operations applied:
- **+ <N> added**
- **~ <N> modified**
- **- <N> removed**
- **→ <N> renamed**

Updated capabilities:
- `<capability-1>`
- `<capability-2>`
<else if --skip-specs>
Spec updates skipped (--skip-specs flag used)
<endif>

## Review Notes

This PR archives a completed change and updates specifications to reflect the
implemented functionality. Please review:

1. Archived change structure and completeness
2. Spec delta accuracy and correctness
3. Merged spec content

---
Generated by `spectr archive --pr`
```text

**Rationale**:

- Provides reviewers with clear context about what changed
- Lists updated capabilities for targeted spec review
- Acknowledges when specs were skipped
- Footer attribution aids debugging and process understanding

### 7. Flag Compatibility

**Decision**: The `--pr` flag SHALL work with all existing archive flags except
when operations would conflict.

**Compatible combinations**:

- `--pr --yes` - Auto-confirm archive prompts, then create PR
- `--pr --skip-specs` - Archive without spec updates, then create PR (PR body
  notes no specs updated)
- `--pr --yes --skip-specs` - Fully automated archive + PR without spec updates
- `--pr --interactive` - Select change interactively, then archive + PR

**Incompatible combinations**:

- `--pr --no-validate` - **ALLOWED BUT WARNED**: Archiver already warns about
  skipping validation; PR proceeds if user confirms

**Rationale**:

- `--pr` is orthogonal to archive workflow flags
- Users may legitimately want to skip specs or use auto-confirm with PR creation
- Validation skip already has safeguards; adding PR doesn't change that

### 8. Error Handling and Recovery

**Decision**: Fail fast with descriptive errors at each step; leave archive
committed but don't abandon it.

**Error Cases and Handling**:

| Error | Behavior | State After Error |
|-------|----------|-------------------|
| Archive fails | Exit before git ops | No changes; retry |
| Not in git repo | Error message | Archive done; uncommitted |
| No origin remote | Error message | Archive done; uncommitted |
| Branch creation fails | Error message | Archive done; uncommitted |
| Platform detection fails | Error message | Archive done; pushed |
| PR CLI not installed | Error message | Archive done; pushed |
| Commit fails | Error message | Archive done; branch created |
| Push fails | Error message | Archive done; committed |
| PR creation fails | Error message | Archive done; pushed |

**Rationale**:

- Archive operation is the primary goal; PR is a convenience feature
- If archive succeeds but PR fails, user still has a valid archive and can
  create PR manually
- Clear error messages guide users on how to recover
- Branch/commit/push are preserved even if PR creation fails, reducing rework

### 9. Git Platform Detection

**Decision**: Reuse platform detection logic from `spectr propose` command (once
implemented).

**Detection Logic** (from propose design):

```text
github.com → GitHub (use `gh`)
gitlab.com → GitLab (use `glab`)
gitea/forgejo in URL → Gitea (use `tea`)
```text

**Rationale**:

- Consistent behavior between propose and archive commands
- Well-tested logic that handles HTTPS and SSH URLs
- Supports self-hosted instances of GitLab, Gitea, Forgejo

### 10. Cross-Platform Support

**Decision**: Use Go's `exec.Command` for all git and PR CLI operations; ensure
Windows, macOS, Linux compatibility.

**Rationale**:

- `exec.Command` abstracts OS-specific shell differences
- All platforms have git available
- `gh`, `glab`, `tea` are cross-platform binaries
- Mirrors approach from propose command for consistency

## Alternatives Considered

### Alternative 1: Add `spectr pr` command instead of `--pr` flag

**Rejected**:

- Requires users to know which files to stage after archiving
- More error-prone (users might forget to stage updated specs)
- Less convenient than single-command workflow
- Harder to make atomic (archive + PR creation)

### Alternative 2: Always create PR after archive (no flag)

**Rejected**:

- Too opinionated; some users may not want PRs for every archive
- Some users may use different git workflows (trunk-based, direct commits to
  main)
- Breaking change to existing archive behavior

### Alternative 3: Prompt user "Create PR?" after successful archive

**Rejected**:

- Inconsistent with other flags (users prefer declarative flags over interactive
  prompts for automation)
- Harder to use in CI/CD or scripts
- Can be approximated with `--pr --yes` for automation or `--pr` for interactive
  workflows

### Alternative 4: Separate PR from archive (manual `git add`)

**Rejected**:

- Defeats purpose of automation
- Error-prone (users might forget files)
- Archive operation already knows what files changed

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Archive succeeds, git fails | Archive but no PR | Clear errors guide recovery |
| Uncommitted changes in specs/ | Unrelated changes | Check warns; prompt user |
| PR CLI incompatibility | PR fails | Document versions |
| Network failure | Partial completion | Steps independent; retry |
| Branch name conflicts | Branch creation fails | Error suggests solution |
| Commit message too long | Git warning/failure | Keep message concise |

## Migration Plan

This is a new feature; no migration required. Existing `spectr archive`
workflows continue to work identically when `--pr` flag is not used.

## Open Questions

1. Should we add `--dry-run` flag to preview PR without creating it?
2. Should we support custom branch names via `--branch` flag?
3. Should we add `--pr-title` and `--pr-body` flags for custom PR content?
4. How should we handle the case where `add-propose-command` is not yet
  implemented? (Duplicate git logic temporarily, or wait?)
5. Should `--pr` require `--yes` flag for full automation, or allow interactive
  archive + automatic PR?

**Recommendations**:

1. No dry-run initially; users can inspect git state before pushing if needed
2. No custom branch names initially; keep simple and consistent
3. No custom PR content initially; default template is sufficient and editable
  after creation
4. If `propose` is not implemented, create `internal/git/` package as part of
  this change
5. Allow `--pr` without `--yes`; they serve different purposes (auto-confirm vs
  auto-PR)

---

**Decision Status**: Ready for implementation after approval and coordination
with `add-propose-command`.
