## 1. Create internal/tui package foundation
- [ ] 1.1 Create `internal/tui/styles.go` with `ApplyTableStyles()` function
- [ ] 1.2 Create `internal/tui/helpers.go` with `TruncateString()` and `CopyToClipboard()` functions
- [ ] 1.3 Create `internal/tui/types.go` with shared types (`KeyAction`, `ActionResult`)
- [ ] 1.4 Write tests for helpers in `internal/tui/helpers_test.go`

## 2. Implement TablePicker component
- [ ] 2.1 Create `internal/tui/table.go` with `TablePicker` struct
- [ ] 2.2 Implement `Init()`, `Update()`, `View()` for `TablePicker`
- [ ] 2.3 Add configurable action registration (`WithAction`)
- [ ] 2.4 Add help text generation from registered actions
- [ ] 2.5 Write tests for `TablePicker` in `internal/tui/table_test.go`

## 3. Implement MenuPicker component
- [ ] 3.1 Create `internal/tui/menu.go` with `MenuPicker` struct
- [ ] 3.2 Implement `Init()`, `Update()`, `View()` for `MenuPicker`
- [ ] 3.3 Add selection callback support
- [ ] 3.4 Write tests for `MenuPicker` in `internal/tui/menu_test.go`

## 4. Refactor list/interactive.go
- [ ] 4.1 Update imports to use `internal/tui` package
- [ ] 4.2 Replace inline `applyTableStyles` calls with `tui.ApplyTableStyles`
- [ ] 4.3 Replace inline `truncateString` calls with `tui.TruncateString`
- [ ] 4.4 Refactor `RunInteractiveChanges` to use `TablePicker`
- [ ] 4.5 Refactor `RunInteractiveSpecs` to use `TablePicker`
- [ ] 4.6 Refactor `RunInteractiveAll` to use `TablePicker`
- [ ] 4.7 Refactor `RunInteractiveArchive` to use `TablePicker`
- [ ] 4.8 Update tests in `internal/list/interactive_test.go`

## 5. Refactor validation/interactive.go
- [ ] 5.1 Update imports to use `internal/tui` package
- [ ] 5.2 Replace inline `applyTableStyles` calls with `tui.ApplyTableStyles`
- [ ] 5.3 Replace inline `truncateString` calls with `tui.TruncateString`
- [ ] 5.4 Refactor `menuModel` to use `MenuPicker`
- [ ] 5.5 Refactor `itemPickerModel` to use `TablePicker`
- [ ] 5.6 Update tests in `internal/validation/interactive_test.go`

## 6. Cleanup and verification
- [ ] 6.1 Remove orphaned helper functions from `internal/list/helpers.go`
- [ ] 6.2 Remove duplicate implementations from `internal/validation/interactive.go`
- [ ] 6.3 Run `go test ./...` and ensure all tests pass
- [ ] 6.4 Run `golangci-lint run` and fix any linting issues
- [ ] 6.5 Manually test `spectr list -I`, `spectr list --specs -I`, `spectr list --all -I`
- [ ] 6.6 Manually test `spectr archive` interactive mode
- [ ] 6.7 Manually test `spectr validate` interactive mode
