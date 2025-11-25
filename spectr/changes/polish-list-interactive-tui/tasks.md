## 1. Implementation

- [ ] 1.1 Update help text format in `RunInteractiveChanges()` to use two-line format
- [ ] 1.2 Update help text format in `RunInteractiveSpecs()` to use two-line format
- [ ] 1.3 Update help text format in `RunInteractiveAll()` to use two-line format
- [ ] 1.4 Update help text format in `RunInteractiveArchive()` to use two-line format
- [ ] 1.5 Update help text format in `rebuildUnifiedTable()` to use two-line format
- [ ] 1.6 Add `a` key handler in `Update()` to trigger archive workflow
- [ ] 1.7 Add `archiveMsg` type and `handleArchive()` method for archive integration
- [ ] 1.8 Update `interactiveModel` struct with archive-related fields if needed

## 2. Testing

- [ ] 2.1 Update existing `interactive_test.go` tests for new help text format
- [ ] 2.2 Add tests for `a` key archive hotkey behavior
- [ ] 2.3 Verify archive hotkey is only available in change modes (not spec-only mode)

## 3. Validation

- [ ] 3.1 Run `go build ./...` to verify compilation
- [ ] 3.2 Run `go test ./internal/list/...` to verify tests pass
- [ ] 3.3 Manual test: `spectr list -I` shows new help text format
- [ ] 3.4 Manual test: Press `a` in list mode triggers archive workflow
