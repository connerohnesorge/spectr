# Tasks: Implement `spectr propose` Command

## 1. Core Implementation

- [ ] 1.1 Create `internal/propose/` package directory
- [ ] 1.2 Create `internal/propose/proposer.go` with main Proposer struct
- [ ] 1.3 Implement git platform detection from remote URL
- [ ] 1.4 Implement git operations (branch creation, staging, commit, push)
- [ ] 1.5 Implement PR CLI tool selection and invocation logic
- [ ] 1.6 Implement PR URL extraction from tool output
- [ ] 1.7 Add error handling for all git and CLI operations

## 2. Command Integration

- [ ] 2.1 Create `cmd/propose.go` with ProposeCmd struct
- [ ] 2.2 Register ProposeCmd in `cmd/root.go` CLI struct
- [ ] 2.3 Implement ProposeCmd.Run() method
- [ ] 2.4 Add command help text and flag descriptions

## 3. Validation & Safety

- [ ] 3.1 Validate that change folder exists before any git operations
- [ ] 3.2 Validate that change folder is NOT tracked by git
- [ ] 3.3 Validate git repository exists (has .git directory)
- [ ] 3.4 Validate origin remote is configured
- [ ] 3.5 Validate required PR CLI tool is installed

## 4. Testing

- [ ] 4.1 Create `internal/propose/proposer_test.go` with unit tests
- [ ] 4.2 Test git platform detection for GitHub, GitLab, Gitea URLs
- [ ] 4.3 Test git operations (mocked git calls)
- [ ] 4.4 Test PR CLI tool selection logic
- [ ] 4.5 Test error cases (missing change, already tracked, no remote)
- [ ] 4.6 Test PR URL extraction from different tool outputs
- [ ] 4.7 Create `cmd/propose_test.go` with command integration tests

## 5. Documentation & Verification

- [ ] 5.1 Verify command works end-to-end with a real change proposal
- [ ] 5.2 Test on Linux, macOS, and Windows (or use CI)
- [ ] 5.3 Ensure error messages are helpful and actionable
- [ ] 5.4 Add inline comments for complex logic

## 6. Polish

- [ ] 6.1 Run `golangci-lint` and fix any linting issues
- [ ] 6.2 Ensure all exported functions have doc comments
- [ ] 6.3 Test with edge cases (spaces in change ID, special chars in URL)
- [ ] 6.4 Verify that unrelated uncommitted changes are not included in commit

---

**Completion Criteria**:
- All tests pass
- No linting errors
- Command successfully creates a PR for a new change proposal
- All error paths are tested and provide clear guidance
