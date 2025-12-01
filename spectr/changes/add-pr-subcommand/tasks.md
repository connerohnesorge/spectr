# Tasks

## 1. Git Package Foundation

- [ ] 1.1 Create `internal/git/doc.go` with package documentation
- [ ] 1.2 Create `internal/git/platform.go` with platform detection:
  - [ ] Define Platform constants (GitHub, GitLab, Gitea, Bitbucket, Unknown)
  - [ ] Define PlatformInfo struct with Platform, CLITool, RepoURL fields
  - [ ] Implement DetectPlatform() that parses remote URL patterns
  - [ ] Support both HTTPS and SSH URL formats
- [ ] 1.3 Create `internal/git/platform_test.go` with comprehensive tests:
  - [ ] Test GitHub URLs (https, ssh, enterprise)
  - [ ] Test GitLab URLs (https, ssh, self-hosted)
  - [ ] Test Gitea/Forgejo URLs
  - [ ] Test Bitbucket URLs
  - [ ] Test unknown/custom hosts
- [ ] 1.4 Create `internal/git/worktree.go` with worktree management:
  - [ ] Implement CreateWorktree() - creates temp worktree on new branch
  - [ ] Implement CleanupWorktree() - removes worktree safely
  - [ ] Implement ExecuteInWorktree() - runs command in worktree context
  - [ ] Implement GetBaseBranch() - auto-detect main/master
- [ ] 1.5 Create `internal/git/worktree_test.go` with unit tests:
  - [ ] Test worktree creation with real git repo
  - [ ] Test cleanup (normal and forced)
  - [ ] Test base branch detection
  - [ ] Test error handling (missing git, no remote, etc.)

## 2. PR Package Foundation

- [ ] 2.1 Create `internal/pr/doc.go` with package documentation
- [ ] 2.2 Create `internal/pr/templates.go` with message templates:
  - [ ] Define CommitMessage struct and template for archive mode
  - [ ] Define CommitMessage struct and template for new mode
  - [ ] Define PRBody struct and template for archive mode
  - [ ] Define PRBody struct and template for new mode
  - [ ] Implement Render() methods using text/template
- [ ] 2.3 Create `internal/pr/templates_test.go`:
  - [ ] Test commit message generation for both modes
  - [ ] Test PR body generation for both modes
  - [ ] Test with various operation counts
- [ ] 2.4 Create `internal/pr/workflow.go` with PR workflow:
  - [ ] Define PRConfig struct with all workflow options
  - [ ] Define PRResult struct with PR URL and metadata
  - [ ] Implement ExecuteArchivePR() - full archive PR workflow
  - [ ] Implement ExecuteNewPR() - full new PR workflow
  - [ ] Implement common helpers: validatePrereqs(), createPR(), etc.
- [ ] 2.5 Create `internal/pr/workflow_test.go`:
  - [ ] Test prerequisite validation
  - [ ] Test workflow steps in isolation
  - [ ] Test error handling at each stage

## 3. CLI Command Implementation

- [ ] 3.1 Create `cmd/pr.go` with Kong command structure:
  - [ ] Define PRCmd struct with Archive and New subcommands
  - [ ] Define PRArchiveCmd with flags: --base, --draft, --force, --dry-run, --skip-specs
  - [ ] Define PRNewCmd with flags: --base, --draft, --force, --dry-run
  - [ ] Implement Run() methods that delegate to internal/pr
- [ ] 3.2 Update `cmd/root.go` to add PR command to CLI struct
- [ ] 3.3 Add shell completion support for change IDs in PR commands
- [ ] 3.4 Create `cmd/pr_test.go`:
  - [ ] Test command parsing
  - [ ] Test flag validation
  - [ ] Test help output

## 4. Platform CLI Integration

- [ ] 4.1 Implement GitHub PR creation with `gh`:
  - [ ] Build gh command with --title, --body-file, --base
  - [ ] Handle --draft flag
  - [ ] Parse PR URL from output
- [ ] 4.2 Implement GitLab MR creation with `glab`:
  - [ ] Build glab command with --title, --description, --target-branch
  - [ ] Handle --draft flag (--draft or -d)
  - [ ] Parse MR URL from output
- [ ] 4.3 Implement Gitea PR creation with `tea`:
  - [ ] Build tea command with --title, --description, --base
  - [ ] Parse PR URL from output
- [ ] 4.4 Implement Bitbucket fallback:
  - [ ] Generate manual PR URL
  - [ ] Display instructions to user
- [ ] 4.5 Add CLI availability detection:
  - [ ] Check if gh/glab/tea is installed
  - [ ] Check if CLI is authenticated
  - [ ] Provide helpful error messages

## 5. Integration Testing

- [ ] 5.1 Create integration test harness:
  - [ ] Set up temp git repos with origin remote
  - [ ] Mock or stub platform CLI calls
  - [ ] Helper to create test change proposals
- [ ] 5.2 Test `spectr pr archive` end-to-end:
  - [ ] Test with valid change, verify worktree isolation
  - [ ] Test archive workflow executed in worktree
  - [ ] Test commit message content
  - [ ] Test cleanup on success
- [ ] 5.3 Test `spectr pr new` end-to-end:
  - [ ] Test with valid change, verify copy operation
  - [ ] Test commit message content
  - [ ] Test cleanup on success
- [ ] 5.4 Test error scenarios:
  - [ ] Test with missing git
  - [ ] Test with no origin remote
  - [ ] Test with missing platform CLI
  - [ ] Test with non-existent change
  - [ ] Test cleanup on failure

## 6. Documentation and Polish

- [ ] 6.1 Add `spectr pr --help` documentation
- [ ] 6.2 Add `spectr pr archive --help` documentation
- [ ] 6.3 Add `spectr pr new --help` documentation
- [ ] 6.4 Update AGENTS.md with new command reference
- [ ] 6.5 Add examples to documentation showing common workflows
