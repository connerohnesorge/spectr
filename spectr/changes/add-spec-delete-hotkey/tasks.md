## 1. Implementation

- [ ] 1.1 Add `handleDelete()` method to `interactiveModel` in `internal/list/interactive.go`
- [ ] 1.2 Add 'd' key case to the `Update()` switch statement in `internal/list/interactive.go`
- [ ] 1.3 Implement confirmation prompt before deletion (inline TUI state or hchoose/confirm)
- [ ] 1.4 Implement `deleteSpecFolder()` function using `os.RemoveAll()` on `spectr/specs/<spec-id>/`
- [ ] 1.5 Handle deletion result: remove row from table, update row count, reset cursor if needed
- [ ] 1.6 Update help text for specs mode to include `d: delete`
- [ ] 1.7 Update help text for unified mode to include `d: delete` (spec rows only)
- [ ] 1.8 Ensure 'd' key is ignored when a change row is selected in unified mode

## 2. Testing

- [ ] 2.1 Add unit tests for `handleDelete()` function
- [ ] 2.2 Add integration test for deletion flow in specs mode
- [ ] 2.3 Add test for deletion being ignored on change rows in unified mode
- [ ] 2.4 Add test for cancellation flow when user declines confirmation

## 3. Validation

- [ ] 3.1 Run `go test ./...` to verify all tests pass
- [ ] 3.2 Run `golangci-lint run` to verify code quality
- [ ] 3.3 Manual testing of delete flow in TUI
