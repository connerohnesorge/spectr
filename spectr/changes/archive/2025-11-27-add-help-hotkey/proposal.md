# Change: Show Hotkeys Only on '?' Keypress

## Why

The current TUI displays all hotkeys at all times in the help text, which takes up screen real estate and adds visual clutter. Users familiar with the interface don't need to see `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | q: quit` on every screen. A common pattern in terminal applications (like vim, less, htop) is to show a minimal footer and reveal the full help on demand via `?`.

## What Changes

- Remove inline hotkey display from the default TUI view
- Show only a minimal footer with item count, project path, and a hint: `?: help`
- Add `?` hotkey that toggles display of the full hotkey reference
- When help is shown, display all available hotkeys in the footer area
- Pressing `?` again (or any navigation key) hides the help and returns to minimal footer

## Impact

- Affected specs: cli-interface
- Affected code: internal/tui/table.go, internal/list/interactive.go, internal/list/interactive_test.go
