## 1. Core Implementation

- [ ] 1.1 Add `GitShowFile(ref, path string) ([]byte, error)` helper in `internal/git/` to read files from branches
- [ ] 1.2 Add `--base-branch` flag to `cmd/validate.go`
- [ ] 1.3 Create `BaseSpecResolver` interface in `internal/validation/` to abstract spec resolution
- [ ] 1.4 Implement `LocalBaseSpecResolver` (current behavior) and `GitRefBaseSpecResolver` (git show)
- [ ] 1.5 Update `validateDeltaAgainstBaseSpec` in `change_rules.go` to accept a resolver

## 2. Pre-flight Validation

- [ ] 2.1 Add `ValidateChangeAgainstBranch(changePath, targetBranch string)` function
- [ ] 2.2 Integrate pre-flight validation into `internal/pr/workflow.go` before worktree creation
- [ ] 2.3 Update error messages to explain local vs remote spec discrepancy

## 3. Testing

- [ ] 3.1 Add unit tests for `GitShowFile` with mock git commands
- [ ] 3.2 Add unit tests for `GitRefBaseSpecResolver`
- [ ] 3.3 Add integration test for `--base-branch` validation flag
- [ ] 3.4 Add integration test for pre-flight validation in PR workflow

## 4. Documentation

- [ ] 4.1 Update CLI help text for `--base-branch` flag
- [ ] 4.2 Add troubleshooting section in AGENTS.md for this validation scenario
