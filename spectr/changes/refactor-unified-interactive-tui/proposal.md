# Change: Refactor Unified Interactive TUI

## Why

The codebase has three separate implementations of interactive table TUIs with significant code duplication:

1. **`internal/list/interactive.go`** (~650 lines) - handles changes, specs, unified mode, archive selection
2. **`internal/validation/interactive.go`** (~450 lines) - handles validation menu and item picker

Both implementations duplicate:
- Table styling logic (`applyTableStyles`)
- String truncation logic (`truncateString`)
- Bubbletea model patterns (Init, Update, View)
- Key binding handling (q, Ctrl+C, Enter, arrow keys)
- Help text formatting patterns

This makes maintenance difficult and creates inconsistency risk when updating TUI behavior.

## What Changes

- Extract shared TUI components into a new `internal/tui` package
- Create composable building blocks: `TablePicker`, `MenuPicker`, styling utilities
- Refactor `list/interactive.go` to use shared components
- Refactor `validation/interactive.go` to use shared components
- Maintain 100% backward compatibility with existing CLI behavior
- Reduce total lines of code while improving testability

## Impact

- Affected specs: `cli-interface`, `validation`
- Affected code:
  - `internal/list/interactive.go` - major refactor to use shared components
  - `internal/list/helpers.go` - move `applyTableStyles`, `truncateString` to tui package
  - `internal/validation/interactive.go` - major refactor to use shared components
  - NEW: `internal/tui/` package with shared components
- No breaking changes to CLI interface or user-visible behavior
- Tests will need updates to test the new shared components
