# Implementation Tasks

## 1. Update Function Signatures

- [x] 1.1 Update `RunInteractiveChanges()` in `internal/list/interactive.go` to accept `projectPath string` parameter
- [x] 1.2 Update `RunInteractiveArchive()` in `internal/list/interactive.go` to accept `projectPath string` parameter
- [x] 1.3 Set `projectPath` field on `interactiveModel` in `RunInteractiveChanges()` during initialization
- [x] 1.4 Set `projectPath` field on `interactiveModel` in `RunInteractiveArchive()` during initialization

## 2. Update Callers

- [x] 2.1 Update `listChanges()` in `cmd/list.go` to pass `projectPath` to `RunInteractiveChanges()`
- [x] 2.2 Update `runInteractiveArchiveForArchiver()` in `internal/archive/interactive_bridge.go` to accept and pass `projectPath`
- [x] 2.3 Update call to `runInteractiveArchiveForArchiver()` in `internal/archive/archiver.go` to pass `projectRoot`

## 3. Add Project Path Display

- [x] 3.1 Update help text in `RunInteractiveChanges()` to include project path display
- [x] 3.2 Update help text in `RunInteractiveArchive()` to include project path display
- [x] 3.3 Update help text in `RunInteractiveSpecs()` to include project path display for consistency

## 4. Update Tests

- [x] 4.1 Update test calls to `RunInteractiveChanges()` in `internal/list/interactive_test.go` to pass `projectPath`
- [x] 4.2 Update test calls to `RunInteractiveArchive()` in `internal/list/interactive_test.go` to pass `projectPath`
- [x] 4.3 Verify that all tests pass with the updated signatures

## 5. Verification

- [x] 5.1 Run `go build` to verify compilation
- [x] 5.2 Run `go test ./...` to verify all tests pass
- [x] 5.3 Run `spectr list -I` manually to verify project path is displayed
- [x] 5.4 Run `spectr archive` (no args) manually to verify project path is displayed
- [x] 5.5 Run `spectr list --specs -I` manually to verify project path is displayed
