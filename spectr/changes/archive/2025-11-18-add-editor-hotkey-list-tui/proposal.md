# Change: Add 'e' Hotkey to List TUI for Opening Specs in $EDITOR

## Why

Currently, users in the interactive list mode can navigate and copy IDs to the
clipboard, but they must manually open spec files in their editor afterward.
This adds friction to the common workflow of "list specs → review a spec → make
changes." Adding an 'e' hotkey that directly opens the selected spec in $EDITOR
streamlines this workflow and makes the TUI more productive.

## What Changes

- Add 'e' key handler in the interactive list TUI for specs mode
- When 'e' is pressed on a selected spec, open `spectr/specs/<spec-id>/spec.md`
  in the user's $EDITOR
- Update help text to show the new 'e' hotkey option
- Handle cases where $EDITOR is not set (fallback to sensible defaults or show
  error)
- Ensure the TUI remains active after the editor is closed (unless user quits)

## Impact

- Affected specs: `cli-interface`
- Affected code:
  - `internal/list/interactive.go` - Add 'e' key handling in Update() and new
    handleEdit() method
  - `internal/list/interactive.go` - Update helpText to include "e: edit spec"
  - May need helper function to launch $EDITOR with proper file path
- Breaking changes: None - this is additive functionality
- User benefit: Faster workflow for reviewing and editing specs
