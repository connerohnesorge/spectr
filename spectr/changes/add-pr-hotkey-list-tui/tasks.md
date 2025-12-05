## 1. Implementation

- [ ] 1.1 Add `PRRequested` field to `ActionResult` struct in `internal/tui/types.go`
- [ ] 1.2 Register `P` (Shift+P) hotkey action in the changes list TUI that sets `PRRequested: true` and exits
- [ ] 1.3 Update `renderQuitView()` in `internal/tui/table.go` to handle PR requested state with appropriate message
- [ ] 1.4 Update list command handler to detect `PRRequested` result and invoke `spectr pr` workflow

## 2. Testing

- [ ] 2.1 Add unit tests for the new `P` hotkey action in `internal/tui/table_test.go`
- [ ] 2.2 Verify the hotkey only appears in changes mode (not specs mode)
- [ ] 2.3 Test that `P` hotkey is ignored when search mode is active (if applicable)

## 3. Documentation

- [ ] 3.1 Create `assets/vhs/pr-hotkey.tape` demonstrating the Shift+P hotkey utility
- [ ] 3.2 Generate `assets/gifs/pr-hotkey.gif` from the tape file
- [ ] 3.3 Create example project in `examples/pr-hotkey/` for the VHS demo

## 4. Validation

- [ ] 4.1 Run `spectr validate add-pr-hotkey-list-tui --strict` and fix any issues
- [ ] 4.2 Run `go test ./...` to ensure all tests pass
- [ ] 4.3 Manually test the hotkey in the TUI
