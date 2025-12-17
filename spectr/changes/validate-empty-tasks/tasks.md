## 1. Implementation
- [ ] 1.1 Add `validateTasksFile()` function to `internal/validation/change_rules.go`
- [ ] 1.2 Call `validateTasksFile()` from `ValidateChangeDeltaSpecs()` and append issues to report
- [ ] 1.3 Add test cases in `internal/validation/change_rules_test.go`

## 2. Verification
- [ ] 2.1 Run `go test ./internal/validation/...` to verify tests pass
- [ ] 2.2 Run `spectr validate --all` to verify existing changes still validate
- [ ] 2.3 Manually test with an empty tasks.md file to confirm error is reported
