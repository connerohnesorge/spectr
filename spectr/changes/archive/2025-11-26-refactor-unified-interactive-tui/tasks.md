## 1. Create internal/tui package foundation

- [x] 1.1 Create `internal/tui/styles.go` with `ApplyTableStyles()` function
- [x] 1.2 Create `internal/tui/helpers.go` with `TruncateString()` and `CopyToClipboard()` functions
- [x] 1.3 Create `internal/tui/types.go` with shared types (`KeyAction`, `ActionResult`)
- [x] 1.4 Write tests for helpers in `internal/tui/helpers_test.go`

## 2. Implement TablePicker component

- [x] 2.1 Create `internal/tui/table.go` with `TablePicker` struct
- [x] 2.2 Implement `Init()`, `Update()`, `View()` for `TablePicker`
- [x] 2.3 Add configurable action registration (`WithAction`)
- [x] 2.4 Add help text generation from registered actions
- [x] 2.5 Write tests for `TablePicker` in `internal/tui/table_test.go`

## 3. Implement MenuPicker component

- [x] 3.1 Create `internal/tui/menu.go` with `MenuPicker` struct
- [x] 3.2 Implement `Init()`, `Update()`, `View()` for `MenuPicker`
- [x] 3.3 Add selection callback support
- [x] 3.4 Write tests for `MenuPicker` in `internal/tui/menu_test.go`

## 4. Refactor list/interactive.go

- [x] 4.1 Update imports to use `internal/tui` package
- [x] 4.2 Replace inline `applyTableStyles` calls with `tui.ApplyTableStyles`
- [x] 4.3 Replace inline `truncateString` calls with `tui.TruncateString`
- [x] 4.4 Refactor `RunInteractiveChanges` to use `TablePicker`
- [x] 4.5 Refactor `RunInteractiveSpecs` to use `TablePicker`
- [x] 4.6 Refactor `RunInteractiveAll` to use `TablePicker`
- [x] 4.7 Refactor `RunInteractiveArchive` to use `TablePicker`
- [x] 4.8 Update tests in `internal/list/interactive_test.go`

## 5. Refactor validation/interactive.go

- [x] 5.1 Update imports to use `internal/tui` package
- [x] 5.2 Replace inline `applyTableStyles` calls with `tui.ApplyTableStyles`
- [x] 5.3 Replace inline `truncateString` calls with `tui.TruncateString`
- [x] 5.4 Refactor `menuModel` to use `MenuPicker`
- [x] 5.5 Refactor `itemPickerModel` to use `TablePicker`
- [x] 5.6 Update tests in `internal/validation/interactive_test.go`

## 6. Cleanup and verification

- [x] 6.1 Remove orphaned helper functions from `internal/list/helpers.go`
- [x] 6.2 Remove duplicate implementations from `internal/validation/interactive.go`
- [x] 6.3 Run `go test ./...` and ensure all tests pass
- [x] 6.4 Run `golangci-lint run` and fix any linting issues
- [x] 6.5 Manually test `spectr list -I`, `spectr list --specs -I`, `spectr list --all -I`
- [x] 6.6 Manually test `spectr archive` interactive mode
- [x] 6.7 Manually test `spectr validate` interactive mode
