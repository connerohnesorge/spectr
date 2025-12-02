## 1. Infrastructure

- [x] 1.1 Add `ModeDelete` constant to `internal/pr/templates.go`
- [x] 1.2 Add spec ID resolution to `internal/discovery/specs.go` (similar to `ResolveChangeID`)

## 2. Delete Workflow Implementation

- [x] 2.1 Add `executeDeleteInWorktree()` function to `internal/pr/helpers.go`
- [x] 2.2 Add delete case to `executeOperation()` in `internal/pr/workflow.go`
- [x] 2.3 Add delete-specific commit message template to `internal/pr/templates.go`
- [x] 2.4 Add delete-specific PR body template to `internal/pr/templates.go`

## 3. CLI Integration

- [x] 3.1 Add `PRDeleteCmd` struct to `cmd/pr.go` with flags: `Base`, `Draft`, `Force`, `DryRun`
- [x] 3.2 Add `Delete` field to `PRCmd` struct with `cmd:"" aliases:"d"` tag
- [x] 3.3 Implement `Run()` method on `PRDeleteCmd` with spec ID resolution
- [x] 3.4 Update `validatePrerequisites()` to handle delete mode (validate spec exists instead of change)

## 4. Testing

- [x] 4.1 Add unit tests for `ResolveSpecID()` function
- [x] 4.2 Add unit tests for delete templates
- [x] 4.3 Add workflow tests for delete mode
- [ ] 4.4 Manual testing: `spectr pr delete <spec-id>` and `spectr pr d <spec-id>`

## 5. Validation

- [x] 5.1 Run `spectr validate add-pr-delete-subcommand --strict`
- [x] 5.2 Run existing tests to ensure no regressions
- [x] 5.3 Run linter (`golangci-lint run`)
