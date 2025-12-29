# Implementation Tasks

## 1. Core Implementation

- [x] 1.1 Add `ModeRemove` constant to `internal/pr/templates.go`
- [x] 1.2 Add `PRRemoveCmd` struct to `cmd/pr.go` with appropriate flags
- [x] 1.3 Register `rm` subcommand in `PRCmd` struct with aliases `r` and
  `remove`
- [x] 1.4 Implement `Run()` method for `PRRemoveCmd`

## 2. Workflow Support

- [x] 2.1 Add remove branch prefix handling in `prepareWorkflowContext()`
  (`spectr/remove/<change-id>`)
- [x] 2.2 Create `removeChangeInWorktree()` function in `internal/pr/helpers.go`
- [x] 2.3 Add remove case to `executeOperation()` switch statement in
  `workflow.go`
- [x] 2.4 Update `validatePrerequisites()` to accept `ModeRemove` as valid mode

## 3. Templates

- [x] 3.1 Add remove commit message template to `internal/pr/templates.go`
- [x] 3.2 Add remove PR body template to `internal/pr/templates.go`
- [x] 3.3 Update `RenderCommitMessage()` to handle remove mode
- [x] 3.4 Update `RenderPRBody()` to handle remove mode
- [x] 3.5 Update `GetPRTitle()` to handle remove mode

## 4. Dry Run Support

- [x] 4.1 Add remove mode handling to `executeDryRun()` in
  `internal/pr/dryrun.go`

## 5. Testing

- [x] 5.1 Add unit tests for `PRRemoveCmd.Run()` in `cmd/pr_test.go`
- [x] 5.2 Add unit tests for remove templates in `internal/pr/templates_test.go`
- [x] 5.3 Add unit tests for `removeChangeInWorktree()` in
  `internal/pr/helpers_test.go`
- [x] 5.4 Add integration test for remove workflow in
  `internal/pr/integration_test.go`

## 6. Validation

- [x] 6.1 Run `spectr validate add-pr-rm-subcommand --strict`
- [x] 6.2 Run `go test ./...` to verify all tests pass
- [x] 6.3 Run `golangci-lint run` to verify linting passes
- [x] 6.4 Manual testing of `spectr pr rm` with a test change
