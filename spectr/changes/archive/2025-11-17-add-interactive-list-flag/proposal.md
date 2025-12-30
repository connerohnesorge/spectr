# Change: Add Interactive List Flag with Clipboard Support

## Why

Users frequently need to reference change or spec IDs for other commands (e.g.,
`spectr show`, `spectr validate`). Currently, they must manually type or
copy-paste IDs from text output, which is error-prone and inefficient. An
interactive table interface would allow users to browse items visually and copy
IDs directly to the clipboard with a single keypress.

## What Changes

- Add `-I` / `--interactive` flag to the `list` command
- Implement bubbletea-based table interface for browsing changes and specs
- Display ID, Title, and relevant metadata (deltas/tasks for changes,
  requirements for specs) in table columns
- Enable clipboard copying of the selected item's ID when user presses Enter
- Support standard table navigation (arrow keys, j/k, page up/down)
- Provide visual feedback when ID is copied to clipboard
- Add clipboard dependency for cross-platform clipboard access

## Impact

- **Affected specs**: `cli-interface` (new capability for interactive UI
  components)
- **Affected code**:
  - `cmd/list.go` - Add `-I` flag and wire up interactive mode
  - `internal/list/interactive.go` (new) - Bubbletea table model and logic
  - `go.mod` - Add `github.com/charmbracelet/bubbles` (already present) and
    clipboard library
  - `internal/list/types.go` - May need to expose more fields for table display
- **Dependencies**: Add clipboard library (e.g., `atotto/clipboard` or use OSC
  52 via termenv which is already available)
- **User experience**: Adds optional interactive mode; existing text/JSON output
  unchanged
