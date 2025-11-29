# Change: Add spec folder delete hotkey in TUI view

## Why

Users need the ability to remove entire spec folders directly from the interactive TUI without manually running shell commands. The 'd' hotkey provides a quick and consistent way to delete specs while the user is already navigating the list, improving workflow efficiency.

## What Changes

- Add 'd' hotkey in interactive specs mode (`spectr list --specs -I`) and unified mode (`spectr list --all -I`)
- When pressed on a spec row, prompt for confirmation before deleting the entire `spectr/specs/<spec-id>/` folder
- Display success/error message after deletion
- Refresh the table view after successful deletion
- The hotkey is ignored for change rows in unified mode (changes have their own lifecycle via archive)

## Impact

- Affected specs: `cli-interface`
- Affected code:
  - `internal/list/interactive.go` - Add delete handler and confirmation logic
  - `internal/tui/types.go` - Add `DeleteRequested` field to `ActionResult` (optional)
  - Help text updates in all relevant modes
