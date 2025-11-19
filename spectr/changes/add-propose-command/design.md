# Design: `spectr propose` Command

## Context

The `spectr propose <id>` command automates the workflow of creating a PR from a newly scaffolded change proposal. Users scaffold proposals with `spectr init` or manually, validate them with `spectr validate`, and then need to push them to a remote repository and open a PR.

Currently, this requires manual git operations and platform-specific PR CLI invocations. This design automates these steps while preserving user control and providing clear error feedback.

## Goals

- **Primary**: Automate the branch → commit → push → PR workflow for change proposals
- **Secondary**: Detect git hosting platform automatically to use the correct PR CLI
- **Secondary**: Only work on uncommitted changes (proposed, not yet in version control)
- **Secondary**: Ensure all staged files are ONLY from the target change folder

## Non-Goals

- Modify git configuration or user settings
- Support multiple simultaneous branches or PRs from one proposal
- Auto-merge or handle PR review workflows
- Support git hosting platforms not detected (manual PR creation still available)

## Decisions

### 1. Uncommitted Change Requirement

**Decision**: The command SHALL only accept change proposals that are NOT currently committed to any branch.

**Rationale**:
- The purpose is to help users create PRs for NEW proposals during active development
- If a proposal is already committed to a branch, users can use standard git/gh/glab/tea workflows
- This prevents accidental duplicate branches or PRs

**Implementation**:
- Use `git ls-files` to check if `spectr/changes/<id>/` is tracked
- If tracked, return an error with guidance to commit and PR the existing branch instead

### 2. Git Platform Detection

**Decision**: Auto-detect the git hosting platform from the `origin` remote URL.

**Detection Logic**:
```
github.com → GitHub (use `gh`)
gitlab.com → GitLab (use `glab`)
gitea/forgejo in URL → Gitea (use `tea`)
```

**Rationale**:
- Most developers use one platform per repository
- Parsing `git config --get remote.origin.url` is reliable and requires no external dependencies
- Supports both HTTPS and SSH URLs
- Gracefully handles self-hosted GitLab, Gitea, and Forgejo instances

**Fallback**: If no platform is detected, return an error with the detected remote URL and ask user to run `gh pr create`, `glab mr create`, or `tea pr create` manually.

### 3. Branch Naming Convention

**Decision**: Create branch with name `add-<change-id>`.

**Rationale**:
- Consistent with Spectr change ID convention (verb-led, kebab-case)
- Clearly indicates the branch purpose (adding a proposal)
- Unlikely to conflict with existing branches
- PR titles will reference the change ID, making history clear

### 4. Commit Strategy

**Decision**: Stage ONLY the `spectr/changes/<id>/` folder, commit in one step.

**Implementation**:
```bash
git checkout -b add-<id>
git add spectr/changes/<id>/
git commit -m "Propose: Add <change-id> change proposal

This proposal introduces [brief description from proposal.md].

Change-Id: <change-id>"
```

**Rationale**:
- Isolates changes to exactly what user intended
- If user has other uncommitted changes, they remain uncommitted
- Commit message includes context for reviewers
- Git trailer `Change-Id:` aids with linking tools

### 5. PR Creation

**Decision**: Delegate PR creation to the appropriate CLI tool with minimal options.

**For GitHub** (`gh pr create`):
```bash
gh pr create --title "Propose: <change-id>" --body "..."
```

**For GitLab** (`glab mr create`):
```bash
glab mr create --title "Propose: <change-id>" --body "..."
```

**For Gitea** (`tea pr create`):
```bash
tea pr create --title "Propose: <change-id>" --body "..."
```

**PR Body Template**:
- First line: purpose from `proposal.md` ("Why" section)
- Blank line
- "What Changes" section from `proposal.md`
- Blank line
- "This proposal is ready for review." (optional)

**Rationale**:
- Minimal flags reduce complexity and platform-specific branching
- Body includes context from the proposal for reviewers
- Each tool has its own conventions; let them handle formatting
- User can edit the PR title/body after creation if needed

### 6. Error Handling

**Decision**: Fail fast with descriptive error messages at each step.

**Error Cases**:
- Change folder doesn't exist → "Change proposal '<id>' not found in spectr/changes/"
- Change folder already tracked → "Change '<id>' is already committed. Use git/gh/glab/tea directly to create a PR."
- Not in a git repository → "Not in a git repository. Initialize git with 'git init'."
- Remote not found → "No 'origin' remote configured. Run 'git remote add origin <url>' first."
- Branch creation fails → "Failed to create branch 'add-<id>': [git error]"
- Platform not detected → "Could not detect git hosting platform. Remote URL: [url]. Please create PR manually using gh, glab, or tea."
- PR CLI tool not installed → "[gh/glab/tea] not found. Install from [url]."
- PR creation fails → "Failed to create PR: [tool output]"

**Rationale**:
- Users get specific guidance on what went wrong and how to fix it
- Preserves git state on failure (branch is created but not pushed if PR fails)

### 7. Cross-Platform Support

**Decision**: Use Go's `exec.Command` to invoke git and PR CLI tools; ensure compatibility with Windows, macOS, and Linux.

**Rationale**:
- `exec.Command` abstracts OS-specific shell differences
- All three platforms have git available
- `gh`, `glab`, `tea` are cross-platform and work identically

### 8. PR URL Output

**Decision**: Parse the PR CLI output to extract and display the PR URL.

**Rationale**:
- Provides immediate feedback to the user
- Can be piped or captured for automation
- Each tool outputs the URL in a different format; parsing handles this

## Alternatives Considered

### Alternative 1: Require user to provide PR title/body
**Rejected**: Too verbose for typical case; users can edit PR after creation.

### Alternative 2: Support any branch name (flag for custom branch)
**Rejected**: Adds complexity; `add-<id>` is sufficient and clear.

### Alternative 3: Push to tracking branch automatically
**Rejected**: Requires `--set-upstream` logic; current approach delegates to the PR CLI tools which handle this.

### Alternative 4: Validate the proposal before creating PR
**Rejected**: User already ran `spectr validate`; redundant here. If needed, they can run it again.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|-----------|
| Git operation fails mid-workflow | Branch created but not pushed; user confusion | Clear error message; user can clean up or continue manually |
| PR CLI tool not installed | Command fails | Error message with installation URL for the detected tool |
| User has unstaged changes in the same folder | Unexpected behavior | Explicitly check if folder is uncommitted; fail with guidance |
| PR title/body too long or invalid | PR creation fails | Use simple, guaranteed-valid format; user can edit in UI |
| Auto-detection fails (weird URL format) | Command fails to detect platform | Graceful fallback; user manually runs correct tool |

## Migration Plan

This is a new feature; no migration required. Existing workflows continue to work unchanged.

## Open Questions

1. Should we support --dry-run flag to preview branch/commit without pushing?
2. Should we parse proposal.md ourselves or trust the validation is done?
3. Do we need to support custom remote names (not just "origin")?

---

**Decision Status**: Ready for implementation after approval.
