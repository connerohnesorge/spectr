## 1. Core Implementation

- [ ] 1.1 Add `showHelp` boolean field to `TablePicker` struct in `internal/tui/table.go`
- [ ] 1.2 Add `?` key handler in `TablePicker.Update()` to toggle `showHelp` state
- [ ] 1.3 Create `generateMinimalFooter()` method that shows only item count, project path, and `?: help` hint
- [ ] 1.4 Modify `generateHelpText()` to be the full help text (existing implementation)
- [ ] 1.5 Update `View()` to show minimal footer by default, full help when `showHelp` is true
- [ ] 1.6 Auto-hide help on navigation keys (↑/↓/j/k) to return to minimal view

## 2. Interactive Model Updates

- [ ] 2.1 Add `showHelp` field to `interactiveModel` in `internal/list/interactive.go`
- [ ] 2.2 Add `?` key handler in `interactiveModel.Update()` to toggle help display
- [ ] 2.3 Update `View()` to conditionally show minimal or full help text
- [ ] 2.4 Update `rebuildUnifiedTable()` to preserve help toggle state
- [ ] 2.5 Auto-hide help on navigation keys to avoid cluttering view

## 3. Testing

- [ ] 3.1 Add test for `?` key toggling help visibility in TablePicker
- [ ] 3.2 Add test for minimal footer content (item count, project path, `?: help`)
- [ ] 3.3 Add test for full help content when help is shown
- [ ] 3.4 Add test for auto-hide on navigation keys
- [ ] 3.5 Add test for help toggle in interactiveModel

## 4. Validation

- [ ] 4.1 Run `go test ./...` to verify all tests pass
- [ ] 4.2 Manual test: verify `spectr list -I` shows minimal footer
- [ ] 4.3 Manual test: verify pressing `?` reveals full hotkey list
- [ ] 4.4 Manual test: verify pressing `?` again or navigating hides help
