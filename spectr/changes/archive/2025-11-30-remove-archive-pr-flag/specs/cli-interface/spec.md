## REMOVED Requirements

### Requirement: Archive Command PR Flag

**Reason**: The `--pr` flag adds complexity (worktree management, platform detection, multi-CLI support) for functionality readily available through standard git workflows. Removing simplifies the codebase and maintenance burden.

**Migration**: Users can achieve the same result by running standard git commands after archive:

1. `spectr archive <change-id>`
2. `git add spectr/`
3. `git commit -m "Archive: <change-id>"`
4. `gh pr create` (or `glab mr create`, `tea pr create`)
