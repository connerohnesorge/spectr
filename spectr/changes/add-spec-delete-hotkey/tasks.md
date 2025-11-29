## 1. Implementation

- [x] 1.1 Add `handleDelete()` method to `interactiveModel` in `internal/list/interactive.go`
- [x] 1.2 Add 'd' key case to the `Update()` switch statement in `internal/list/interactive.go`
- [x] 1.3 Implement confirmation prompt before deletion (inline TUI state or hchoose/confirm)
- [x] 1.4 Implement `deleteSpecFolder()` function using `os.RemoveAll()` on `spectr/specs/<spec-id>/`
- [x] 1.5 Handle deletion result: remove row from table, update row count, reset cursor if needed
- [x] 1.6 Update help text for specs mode to include `d: delete`
- [x] 1.7 Update help text for unified mode to include `d: delete (specs)` (spec rows only)
- [x] 1.8 Ensure 'd' key is ignored when a change row is selected in unified mode

## 2. Testing

- [ ] 2.1 Add unit tests for `handleDelete()` function
- [ ] 2.2 Add integration test for deletion flow in specs mode
- [ ] 2.3 Add test for deletion being ignored on change rows in unified mode
- [ ] 2.4 Add test for cancellation flow when user declines confirmation

## 3. Validation

- [x] 3.1 Run `go test ./...` to verify all tests pass
- [x] 3.2 Run `golangci-lint run` to verify code quality
- [ ] 3.3 Manual testing of delete flow in TUI
