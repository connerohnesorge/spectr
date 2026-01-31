# Tasks

## 1. Core Implementation

- [ ] 1.1. Add `LineNumberMode` type with constants `LineNumberOff`,
      `LineNumberRelative`, `LineNumberHybrid` in `internal/list/interactive.go`
- [ ] 1.2. Add `lineNumberMode` field to `interactiveModel` struct
- [ ] 1.3. Implement `renderLineNumbers()` function that returns a styled string
      column based on current mode and cursor position
- [ ] 1.4. Update `View()` method to prepend line number column when mode is not
      off
- [ ] 1.5. Add `#` key handler in `Update()` to cycle through line number modes
- [ ] 1.6. Update footer to show `ln: rel` or `ln: hyb` when line numbers are
      active

## 2. Styling

- [ ] 2.1. Add `LineNumberStyle()` function in `internal/tui/styles.go` returning
      dimmed style (gray, right-aligned)
- [ ] 2.2. Add `CurrentLineNumberStyle()` function for the cursor row (brighter,
      maybe bold)

## 3. Testing

- [ ] 3.1. Add unit tests for `LineNumberMode` cycling logic
- [ ] 3.2. Add unit tests for relative number calculation at different cursor
      positions (top, middle, bottom of list)
- [ ] 3.3. Add unit tests verifying hybrid mode shows absolute for current row

## 4. Documentation

- [ ] 4.1. Update help text to include `#: line numbers` in the TUI help
