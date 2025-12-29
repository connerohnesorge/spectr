# Change: Add PR Proposal Cleanup Confirmation

## Why

After running `spectr pr proposal*` commands, the change proposal directory
remains in the user's working directory (`spectr/changes/<change-id>/`). This
can be confusing because:

1. The proposal has been pushed to a remote branch and a PR has been created
2. The local copy may diverge from the PR branch if the user makes local edits
3. Users may forget to clean up old proposals, cluttering the changes directory

A post-PR confirmation menu would give users a clear opportunity to remove the
local change proposal, keeping their working directory clean.

## What Changes

- After a successful `spectr pr proposal*` command completes and displays the PR
  URL, the system prompts the user with a Bubbletea TUI confirmation menu
  (consistent with other spectr interactive modes)
- The menu asks: "Remove local change proposal from spectr/changes/?"
- Options: Yes (remove), No (keep), with No as the default for safety
- Menu uses arrow key navigation and styled rendering matching spectr's existing
  TUI components
- When `--yes` flag is provided, the prompt is skipped and the local change is
  kept (safe default - no data loss in CI/automation)
- This only applies to `pr proposal` commands, not `pr archive` (which already
  handles cleanup via archiving)

## Impact

- Affected specs: `cli-interface`
- Affected code: `cmd/pr.go`, `internal/pr/workflow.go`
