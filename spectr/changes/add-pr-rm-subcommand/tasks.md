## 1. Core Implementation

- [ ] 1.1 Add `ModeRemove` constant to `internal/pr/helpers.go`
- [ ] 1.2 Add `PRRemoveCmd` struct to `cmd/pr.go` with appropriate flags
- [ ] 1.3 Register `rm` subcommand in `PRCmd` struct with aliases `r` and `remove`
- [ ] 1.4 Implement `Run()` method for `PRRemoveCmd`

## 2. Workflow Support

- [ ] 2.1 Add remove branch prefix handling in `prepareWorkflowContext()` (`spectr/remove/<change-id>`)
- [ ] 2.2 Create `removeChangeInWorktree()` function in `internal/pr/helpers.go`
- [ ] 2.3 Add remove case to `executeOperation()` switch statement in `workflow.go`
- [ ] 2.4 Update `validatePrerequisites()` to accept `ModeRemove` as valid mode

## 3. Templates

- [ ] 3.1 Add remove commit message template to `internal/pr/templates.go`
- [ ] 3.2 Add remove PR body template to `internal/pr/templates.go`
- [ ] 3.3 Update `RenderCommitMessage()` to handle remove mode
- [ ] 3.4 Update `RenderPRBody()` to handle remove mode
- [ ] 3.5 Update `GetPRTitle()` to handle remove mode

## 4. Dry Run Support

- [ ] 4.1 Add remove mode handling to `executeDryRun()` in `internal/pr/dryrun.go`

## 5. Testing

- [ ] 5.1 Add unit tests for `PRRemoveCmd.Run()` in `cmd/pr_test.go`
- [ ] 5.2 Add unit tests for remove templates in `internal/pr/templates_test.go`
- [ ] 5.3 Add unit tests for `removeChangeInWorktree()` in `internal/pr/helpers_test.go`
- [ ] 5.4 Add integration test for remove workflow in `internal/pr/integration_test.go`

## 6. Validation

- [ ] 6.1 Run `spectr validate add-pr-rm-subcommand --strict`
- [ ] 6.2 Run `go test ./...` to verify all tests pass
- [ ] 6.3 Run `golangci-lint run` to verify linting passes
- [ ] 6.4 Manual testing of `spectr pr rm` with a test change
