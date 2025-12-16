## 1. Implementation
- [ ] 1.1 Implement `cleanupLocalChange` function in `internal/pr/worktree.go` to remove local change directory
- [ ] 1.2 Call `cleanupLocalChange` after successful PR creation in `executeWorkflow` (for both archive and remove modes)
- [ ] 1.3 Display warning message before cleanup: "Cleaning up local change directory: spectr/changes/<id>/"

## 2. Testing
- [ ] 2.1 Add unit test for `cleanupLocalChange` function
- [ ] 2.2 Add test for PR archive with local cleanup
- [ ] 2.3 Add test for PR remove with local cleanup

## 3. Validation
- [ ] 3.1 Run `spectr validate fix-leftover-tasks-json --strict`
- [ ] 3.2 Manual testing: create change with tasks.json, run `spectr pr rm`, verify local cleanup
- [ ] 3.3 Verify existing tests pass
