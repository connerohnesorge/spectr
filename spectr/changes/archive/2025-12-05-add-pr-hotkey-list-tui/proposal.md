# Change: Add Shift+P Hotkey for PR Mode in Interactive List TUI

## Why

Users working in the interactive list TUI (`spectr list -I`) currently need to exit the TUI, copy the change ID, and then run `spectr pr` separately. Adding a Shift+P hotkey provides a streamlined workflow to enter "PR mode" directly from the list, reducing context switching and improving the development experience.

## What Changes

- Add `P` (Shift+P) hotkey to the interactive changes list mode
- When pressed on a selected change, exit the TUI and invoke `spectr pr` workflow with that change ID
- Display the hotkey in the help text when `?` is pressed
- Add a VHS tape demo showing the Shift+P hotkey utility

## Impact

- Affected specs: `cli-interface`
- Affected code: `internal/tui/`, `internal/list/`, `cmd/list.go`
- User-visible changes: New `P: pr` hotkey available in changes list mode
