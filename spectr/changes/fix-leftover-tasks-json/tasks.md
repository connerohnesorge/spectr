## 1. Implementation
- [ ] 1.1 Add `NoCleanup` field to `PRConfig` struct in `internal/pr/workflow.go`
- [ ] 1.2 Add `--no-cleanup` flag to `PRArchiveCmd` in `cmd/pr.go`
- [ ] 1.3 Add `--no-cleanup` flag to `PRRemoveCmd` in `cmd/pr.go`
- [ ] 1.4 Implement `cleanupLocalChange` function in `internal/pr/worktree.go` to remove local change directory
- [ ] 1.5 Call `cleanupLocalChange` after successful PR creation in `executeWorkflow` (for both archive and remove modes)
- [ ] 1.6 Display warning message before cleanup: "Cleaning up local change directory: spectr/changes/<id>/"
- [ ] 1.7 Skip cleanup when `--no-cleanup` flag is provided and display message: "Skipping local cleanup (--no-cleanup)"

## 2. Testing
- [ ] 2.1 Add unit test for `cleanupLocalChange` function
- [ ] 2.2 Add test for PR archive with local cleanup
- [ ] 2.3 Add test for PR remove with local cleanup
- [ ] 2.4 Add test for PR archive with `--no-cleanup` flag
- [ ] 2.5 Add test for PR remove with `--no-cleanup` flag

## 3. Validation
- [ ] 3.1 Run `spectr validate fix-leftover-tasks-json --strict`
- [ ] 3.2 Manual testing: create change with tasks.json, run `spectr pr rm`, verify local cleanup
- [ ] 3.3 Verify existing tests pass
