# Change: Add Post-PR Proposal Cleanup Confirmation

## Why

After running `spectr pr proposal` commands, the change proposal directory remains in the user's working directory (`spectr/changes/<change-id>/`). This can be confusing because:

1. The proposal has been pushed to a remote branch and a PR has been created
2. The local copy may diverge from the PR branch if the user makes local edits
3. Users may forget to clean up old proposals, cluttering the changes directory

A post-PR confirmation menu gives users a clear opportunity to remove the local change proposal, keeping their working directory clean.

## What Changes

- Implement the existing spec requirement "PR Proposal Local Change Cleanup Confirmation" which is specified but not yet fully implemented
- After a successful `spectr pr proposal` command displays the PR URL, show a Bubbletea TUI confirmation menu asking whether to remove the local change directory
- The menu defaults to "No, keep it" for safety
- Non-interactive mode (`--yes` flag) skips the prompt and keeps the local directory

## Impact

- Affected specs: `cli-interface` (implements existing requirement)
- Affected code: `cmd/pr.go`, potentially new TUI component in `internal/tui` or `internal/pr`
