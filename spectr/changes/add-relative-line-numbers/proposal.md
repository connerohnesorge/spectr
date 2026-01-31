# Change: Add Relative Line Numbers to List TUI

## Why

The list TUI currently supports vim-style count prefixes for navigation (e.g.,
`9j` to move down 9 rows). However, users must mentally calculate how many rows
away their target is. Adding relative line numbers - a feature popularized by
Vim's `set relativenumber` - displays each row's distance from the currently
selected row, making count-prefix navigation intuitive and efficient.

**GitHub Issue**: #364

## What Changes

- **Line number column**: Add an optional leading column displaying line
  numbers in the interactive list TUI
- **Relative mode**: By default, show relative distances (1, 2, 3...) from the
  cursor row; the cursor row shows absolute position
- **Hybrid mode**: Current row shows absolute number, other rows show relative
  distances (matching Vim's `set number relativenumber`)
- **Toggle hotkey**: Add `#` key to cycle through display modes: off (default)
  -> relative -> hybrid -> off
- **Footer indicator**: Show current line number mode in footer when active

## Impact

- Affected specs: `cli-interface` (interactive list mode requirements)
- Affected code:
  - `internal/list/interactive.go` - line number rendering logic
  - `internal/tui/styles.go` - line number column styling (dimmed)
- **Backward compatible**: Default behavior unchanged (no line numbers shown);
  feature is opt-in via `#` toggle
