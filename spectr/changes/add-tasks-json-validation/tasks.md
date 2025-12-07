## 1. Core Implementation
- [ ] 1.1 Add `IsValid()` method to `TaskStatusValue` type in `internal/parsers/types.go`
- [ ] 1.2 Add `ValidateTasksJson()` function in `internal/parsers/parsers.go` that validates JSON structure and status values
- [ ] 1.3 Add `ValidTaskStatusValues()` helper function returning valid status slice

## 2. Integration
- [ ] 2.1 Add tasks.json validation to change validation in `internal/validation/change_rules.go`
- [ ] 2.2 Ensure validation errors include task ID and descriptive message

## 3. Testing
- [ ] 3.1 Add unit tests for `IsValid()` method with valid and invalid status values
- [ ] 3.2 Add tests for `ValidateTasksJson()` with valid tasks.json
- [ ] 3.3 Add tests for `ValidateTasksJson()` with invalid status values (e.g., "done")
- [ ] 3.4 Add tests for malformed JSON handling
- [ ] 3.5 Add integration test verifying `spectr validate` catches invalid status values
