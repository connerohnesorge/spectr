# Tasks

## 1. Git Package Foundation

- [x] 1.1 Create `internal/git/doc.go` with package documentation
- [x] 1.2 Create `internal/git/platform.go` with platform detection:
  - [x] Define Platform constants (GitHub, GitLab, Gitea, Bitbucket, Unknown)
  - [x] Define PlatformInfo struct with Platform, CLITool, RepoURL fields
  - [x] Implement DetectPlatform() that parses remote URL patterns
  - [x] Support both HTTPS and SSH URL formats
- [x] 1.3 Create `internal/git/platform_test.go` with comprehensive tests:
  - [x] Test GitHub URLs (https, ssh, enterprise)
  - [x] Test GitLab URLs (https, ssh, self-hosted)
  - [x] Test Gitea/Forgejo URLs
  - [x] Test Bitbucket URLs
  - [x] Test unknown/custom hosts
- [x] 1.4 Create `internal/git/worktree.go` with worktree management:
  - [x] Implement CreateWorktree() - creates temp worktree on new branch
  - [x] Implement CleanupWorktree() - removes worktree safely
  - [x] Implement GetBaseBranch() - auto-detect main/master
  - [x] Implement FetchOrigin() - fetch latest refs from origin
  - [x] Implement GetRepoRoot() - get git repository root
  - [x] Implement BranchExists() - check if remote branch exists
  - [x] Implement DeleteRemoteBranch() - delete branch from origin
- [x] 1.5 Create `internal/git/worktree_test.go` with unit tests:
  - [x] Test worktree creation with real git repo
  - [x] Test cleanup (normal and forced)
  - [x] Test base branch detection
  - [x] Test error handling (missing git, no remote, etc.)

## 2. PR Package Foundation

- [x] 2.1 Create `internal/pr/doc.go` with package documentation
- [x] 2.2 Create `internal/pr/templates.go` with message templates:
  - [x] Define CommitMessage struct and template for archive mode
  - [x] Define CommitMessage struct and template for new mode
  - [x] Define PRBody struct and template for archive mode
  - [x] Define PRBody struct and template for new mode
  - [x] Implement Render() methods using text/template
- [x] 2.3 Create `internal/pr/templates_test.go`:
  - [x] Test commit message generation for both modes
  - [x] Test PR body generation for both modes
  - [x] Test with various operation counts
- [x] 2.4 Create `internal/pr/workflow.go` with PR workflow:
  - [x] Define PRConfig struct with all workflow options
  - [x] Define PRResult struct with PR URL and metadata
  - [x] Implement ExecutePR() - unified PR workflow for archive and new modes
  - [x] Implement common helpers: validatePrereqs(), createPR(), etc.
- [x] 2.5 Create `internal/pr/workflow_test.go`:
  - [x] Test prerequisite validation
  - [x] Test workflow steps in isolation
  - [x] Test error handling at each stage

## 3. CLI Command Implementation

- [x] 3.1 Create `cmd/pr.go` with Kong command structure:
  - [x] Define PRCmd struct with Archive and New subcommands
  - [x] Define PRArchiveCmd with flags: --base, --draft, --force, --dry-run, --skip-specs
  - [x] Define PRNewCmd with flags: --base, --draft, --force, --dry-run
  - [x] Implement Run() methods that delegate to internal/pr
- [x] 3.2 Update `cmd/root.go` to add PR command to CLI struct
- [x] 3.3 Add shell completion support for change IDs in PR commands (via predictor:"changeID" tag)
- [x] 3.4 Create `cmd/pr_test.go`:
  - [x] Test command parsing
  - [x] Test flag validation
  - [x] Test help output

## 4. Platform CLI Integration

- [x] 4.1 Implement GitHub PR creation with `gh`:
  - [x] Build gh command with --title, --body-file, --base
  - [x] Handle --draft flag
  - [x] Parse PR URL from output
- [x] 4.2 Implement GitLab MR creation with `glab`:
  - [x] Build glab command with --title, --description, --target-branch
  - [x] Handle --draft flag (--draft or -d)
  - [x] Parse MR URL from output
- [x] 4.3 Implement Gitea PR creation with `tea`:
  - [x] Build tea command with --title, --description, --base
  - [x] Parse PR URL from output
- [x] 4.4 Implement Bitbucket fallback:
  - [x] Generate manual PR URL
  - [x] Display instructions to user
- [x] 4.5 Add CLI availability detection:
  - [x] Check if gh/glab/tea is installed
  - [x] Provide helpful error messages with install suggestions

## 5. Integration Testing

- [x] 5.1 Create integration test harness:
  - [x] Set up temp git repos with origin remote
  - [x] Mock or stub platform CLI calls
  - [x] Helper to create test change proposals
- [x] 5.2 Test `spectr pr archive` end-to-end:
  - [x] Test with valid change, verify worktree isolation (ALL Shell Commands Should Only Run In The Worktree)
  - [x] Test archive workflow executed in worktree
  - [x] Test commit message content
  - [x] Test cleanup on success
- [x] 5.3 Test `spectr pr new` end-to-end:
  - [x] Test with valid change, verify copy operation
  - [x] Test commit message content
  - [x] Test cleanup on success
- [x] 5.4 Test error scenarios:
  - [x] Test with missing git
  - [x] Test with no origin remote
  - [x] Test with missing platform CLI
  - [x] Test with non-existent change
  - [x] Test cleanup on failure

## 6. Documentation and Polish

- [x] 6.1 Add `spectr pr --help` documentation (via Kong help tags)
- [x] 6.2 Add `spectr pr archive --help` documentation (via Kong help tags)
- [x] 6.3 Add `spectr pr new --help` documentation (via Kong help tags)
- [x] 6.4 Update AGENTS.md with new command reference
- [x] 6.5 Add examples to documentation showing common workflows
