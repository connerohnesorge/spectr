# Change: Filter PR Proposal Interactive List to Unmerged Changes Only

## Why

When running `spectr pr proposal` (or `spectr pr p`) interactively without a change ID, the list currently shows ALL active changes. However, some changes may already exist on the main branch (e.g., they were committed directly or merged via a different workflow). Creating a proposal PR for a change that already exists on main is pointless and confusing.

The interactive selection should only show changes that are genuinely new and not yet present on the target branch.

## What Changes

- When `spectr pr proposal` is invoked without a change ID argument, the interactive selection filters out changes that already exist on the main branch
- A new function in `internal/git` checks if a path exists on a specific git ref using `git ls-tree`
- The filtering happens before displaying the interactive list, not after selection
- If all changes are already on main, display a message indicating no unmerged proposals exist

## Impact

- Affected specs: `cli-interface`
- Affected code:
  - `cmd/pr.go`: `selectChangeInteractive()` function
  - `internal/git/branch.go`: New `PathExistsOnRef()` function
  - `internal/list/lister.go`: Potentially add filter capability
