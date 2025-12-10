## 1. Implementation

- [ ] 1.1 Add `Stdout` boolean flag to `ListCmd` struct in `cmd/list.go` with help text
- [ ] 1.2 Add validation that `--stdout` requires `-I` (interactive mode)
- [ ] 1.3 Add `stdoutMode` field to `interactiveModel` in `internal/list/interactive.go`
- [ ] 1.4 Modify `handleEnter()` to print to stdout when `stdoutMode` is true instead of copying to clipboard
- [ ] 1.5 Update `View()` to show clean output (just the ID) when exiting in stdout mode
- [ ] 1.6 Pass `stdoutMode` parameter to `RunInteractiveChanges()` and `RunInteractiveSpecs()` and `RunInteractiveAll()`
- [ ] 1.7 Update `listChanges()`, `listSpecs()`, and `listAll()` in `cmd/list.go` to pass stdout mode

## 2. Testing

- [ ] 2.1 Add unit tests in `cmd/list_test.go` for `--stdout` flag parsing and validation
- [ ] 2.2 Add unit test for `--stdout` without `-I` returns error
- [ ] 2.3 Add unit test for `--stdout` with `--json` returns error (mutually exclusive)
- [ ] 2.4 Add unit tests in `internal/list/interactive_test.go` for stdout mode behavior
- [ ] 2.5 Add test that stdout mode prints ID without formatting prefix
- [ ] 2.6 Add test that stdout mode does not attempt clipboard copy
- [ ] 2.7 Run existing tests to ensure no regressions: `go test ./cmd/... ./internal/list/...`

## 3. Validation

- [ ] 3.1 Run `spectr validate add-stdout-flag-list --strict` to ensure proposal is valid
- [ ] 3.2 Manually test `spectr list -I --stdout` outputs selected ID to stdout
- [ ] 3.3 Manually test piping works: `spectr list -I --stdout | xargs echo "Selected:"`
