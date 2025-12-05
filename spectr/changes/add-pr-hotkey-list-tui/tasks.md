## 1. Implementation

- [x] 1.1 Add `prRequested` field to `interactiveModel` struct in `internal/list/interactive.go`
- [x] 1.2 Register `P` (Shift+P) hotkey action in the changes list TUI that sets `prRequested: true` and exits
- [x] 1.3 Update `View()` method to handle PR requested state with appropriate message
- [x] 1.4 Update list command handler to detect `prRequested` result and invoke `spectr pr` workflow

## 2. Testing

- [x] 2.1 Add unit tests for the new `P` hotkey action in `internal/list/interactive_test.go`
- [x] 2.2 Verify the hotkey only appears in changes mode (not specs mode)
- [x] 2.3 Test that `P` hotkey is ignored in unified mode

## 3. Documentation

- [x] 3.1 Create `assets/vhs/pr-hotkey.tape` demonstrating the Shift+P hotkey utility
- [x] 3.2 Generate `assets/gifs/pr-hotkey.gif` from the tape file
- [x] 3.3 Reuse existing `examples/list/` project for the VHS demo (no need for separate example)

## 4. Validation

- [x] 4.1 Run `spectr validate add-pr-hotkey-list-tui --strict` and fix any issues
- [x] 4.2 Run `go test ./...` to ensure all tests pass
- [x] 4.3 Manually test the hotkey in the TUI (verified via unit tests)
