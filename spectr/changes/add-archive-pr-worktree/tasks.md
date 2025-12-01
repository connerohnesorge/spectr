# Tasks: Implement `--pr` Flag for Archive Command with Git Worktree Isolation

## 1. Git Package Foundation

- [ ] 1.1 Create `internal/git/` package directory
- [ ] 1.2 Create `internal/git/doc.go` with package documentation
- [ ] 1.3 Create `internal/git/types.go` with Platform enum and error types

## 2. Platform Detection

- [ ] 2.1 Create `internal/git/platform.go` with `DetectPlatform(remoteURL string) (Platform, error)`
- [ ] 2.2 Implement HTTPS URL parsing (github.com, gitlab.com, bitbucket.org, etc.)
- [ ] 2.3 Implement SSH URL parsing (git@github.com:..., git@gitlab.com:..., etc.)
- [ ] 2.4 Implement self-hosted detection (gitlab in hostname, gitea/forgejo keywords)
- [ ] 2.5 Create `internal/git/platform_test.go` with table-driven tests for all URL patterns
- [ ] 2.6 Add function to get CLI tool name for platform (`gh`, `glab`, `tea`, or empty for Bitbucket)

## 3. CLI Tool Validation

- [ ] 3.1 Create `internal/git/cli.go` with `CheckCLIInstalled(tool string) error`
- [ ] 3.2 Implement version checking for minimum requirements
- [ ] 3.3 Create `internal/git/cli_test.go` with mock tests
- [ ] 3.4 Add helpful error messages with installation URLs for each tool

## 4. Worktree Management

- [ ] 4.1 Create `internal/git/worktree.go` with worktree operations
- [ ] 4.2 Implement `CreateWorktree(baseBranch, newBranch, path string) error`
- [ ] 4.3 Implement `RemoveWorktree(path string) error` with force cleanup
- [ ] 4.4 Implement `GetBaseBranch() (string, error)` to detect main/master
- [ ] 4.5 Implement `GetRemoteURL() (string, error)` to get origin URL
- [ ] 4.6 Create `internal/git/worktree_test.go` with integration tests
- [ ] 4.7 Add worktree path generation with UUID suffix in temp directory

## 5. Git Operations

- [ ] 5.1 Create `internal/git/operations.go` with basic git commands
- [ ] 5.2 Implement `Add(path string) error` for staging files
- [ ] 5.3 Implement `Commit(message string) error` with heredoc support
- [ ] 5.4 Implement `Push(branch string) error` with upstream tracking
- [ ] 5.5 Implement `IsGitRepository() bool` check
- [ ] 5.6 Create `internal/git/operations_test.go` with mock tests

## 6. PR Creation Abstraction

- [ ] 6.1 Create `internal/git/pr.go` with platform-specific PR creation
- [ ] 6.2 Implement `CreateGitHubPR(title, body, base string) (url string, error)`
- [ ] 6.3 Implement `CreateGitLabMR(title, body, base string) (url string, error)`
- [ ] 6.4 Implement `CreateGiteaPR(title, body, base string) (url string, error)`
- [ ] 6.5 Implement `GetBitbucketPRURL(branch, base string) string` for manual creation
- [ ] 6.6 Add `--draft` flag support for GitHub and GitLab
- [ ] 6.7 Create `internal/git/pr_test.go` with mock tests
- [ ] 6.8 Implement PR URL extraction from CLI output

## 7. Archive Command Integration

- [ ] 7.1 Add `PR bool` field to `ArchiveCmd` struct in `internal/archive/cmd.go`
- [ ] 7.2 Add `Draft bool` field for draft PR support
- [ ] 7.3 Add flag validation: `--pr` requires explicit change ID
- [ ] 7.4 Add flag validation: `--pr` incompatible with `--interactive`
- [ ] 7.5 Add flag validation: `--pr` incompatible with `--no-validate`

## 8. PR Workflow Orchestration

- [ ] 8.1 Create `internal/archive/pr.go` with PR workflow orchestration
- [ ] 8.2 Implement `PRWorkflow` struct with configuration and state
- [ ] 8.3 Implement `PRWorkflow.Execute(changeID string, opCounts OperationCounts, capabilities []string) error`
- [ ] 8.4 Implement pre-flight checks (git repo, origin remote, CLI tool, base branch)
- [ ] 8.5 Implement worktree creation with proper cleanup defer
- [ ] 8.6 Implement self-invocation of `spectr archive` within worktree
- [ ] 8.7 Implement staging, commit, push sequence
- [ ] 8.8 Implement PR creation with platform detection
- [ ] 8.9 Implement result display with PR URL

## 9. Commit and PR Message Generation

- [ ] 9.1 Create `internal/archive/pr_format.go` with message generators
- [ ] 9.2 Implement `GenerateCommitMessage(changeID string, archivePath string, opCounts OperationCounts) string`
- [ ] 9.3 Implement `GeneratePRTitle(changeID string) string`
- [ ] 9.4 Implement `GeneratePRBody(changeID string, archivePath string, opCounts OperationCounts, capabilities []string) string`
- [ ] 9.5 Create `internal/archive/pr_format_test.go` with message format tests

## 10. Archiver Integration

- [ ] 10.1 Modify `Archive()` function to return operation counts and capabilities
- [ ] 10.2 Add PR workflow invocation after successful archive (when --pr flag set)
- [ ] 10.3 Ensure archive success message still displays before PR workflow
- [ ] 10.4 Handle PR workflow errors without affecting archive success status
- [ ] 10.5 Display clear status messages during PR workflow stages

## 11. Error Handling

- [ ] 11.1 Define custom error types in `internal/git/errors.go`
- [ ] 11.2 Implement `ErrNotGitRepository` with remediation message
- [ ] 11.3 Implement `ErrNoRemoteOrigin` with remediation message
- [ ] 11.4 Implement `ErrCLINotInstalled` with installation URL
- [ ] 11.5 Implement `ErrBaseBranchNotFound` with remediation message
- [ ] 11.6 Implement `ErrBranchExists` with remediation message
- [ ] 11.7 Implement `ErrPushFailed` with common causes
- [ ] 11.8 Implement `ErrPRCreationFailed` with CLI output

## 12. Testing

- [ ] 12.1 Create unit tests for platform detection with all URL patterns
- [ ] 12.2 Create unit tests for commit message generation
- [ ] 12.3 Create unit tests for PR body generation
- [ ] 12.4 Create integration test for worktree lifecycle (create, use, cleanup)
- [ ] 12.5 Create mock tests for PR CLI invocation
- [ ] 12.6 Test flag validation: --pr with missing change ID
- [ ] 12.7 Test flag validation: --pr with --interactive
- [ ] 12.8 Test flag validation: --pr with --no-validate
- [ ] 12.9 Test error paths: no git repo, no remote, CLI not installed
- [ ] 12.10 Test cleanup: worktree removed on success and failure

## 13. Documentation and Help

- [ ] 13.1 Update archive command help text to document `--pr` flag
- [ ] 13.2 Add `--draft` flag help text
- [ ] 13.3 Add usage examples in command help
- [ ] 13.4 Document git and CLI tool dependencies in error messages
- [ ] 13.5 Verify help output with `spectr archive --help`

## 14. Code Quality

- [ ] 14.1 Run `golangci-lint` and fix any issues
- [ ] 14.2 Ensure all exported functions have doc comments
- [ ] 14.3 Add inline comments for complex logic
- [ ] 14.4 Verify error wrapping with `fmt.Errorf` context
- [ ] 14.5 Check for proper resource cleanup (worktrees, temp files)

## 15. End-to-End Validation

- [ ] 15.1 Test complete workflow with GitHub repository
- [ ] 15.2 Test complete workflow with GitLab repository
- [ ] 15.3 Test complete workflow with Gitea repository
- [ ] 15.4 Verify branch is created with correct name
- [ ] 15.5 Verify commit message format is correct
- [ ] 15.6 Verify PR title and body format
- [ ] 15.7 Verify worktree is cleaned up after success
- [ ] 15.8 Verify worktree is cleaned up after failure

---

**Completion Criteria**:
- All tests pass
- No linting errors
- `spectr archive <change-id> --pr` creates a PR with isolated worktree workflow
- All error paths provide clear guidance
- Works with GitHub, GitLab, and Gitea platforms
- Archive completes successfully even if PR creation fails
- Worktrees are always cleaned up
