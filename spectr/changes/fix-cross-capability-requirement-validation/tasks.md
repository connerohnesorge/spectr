## 1. Update Cross-File Duplicate Detection

- [ ] 1.1 In `internal/validation/delta_validators.go`, modify the `addedReqs`, `modifiedReqs`, `removedReqs` maps to use `(capability, name)` composite keys instead of just `name`
- [ ] 1.2 Extract capability name from `specPath` (e.g., `spectr/changes/foo/specs/support-aider/spec.md` → `support-aider`)
- [ ] 1.3 Create helper function `extractCapabilityFromPath(specPath string) string`
- [ ] 1.4 Update key generation in `validateAddedRequirements()` to use `capability::normalized_name`
- [ ] 1.5 Update key generation in `validateModifiedRequirements()` to use `capability::normalized_name`
- [ ] 1.6 Update key generation in `validateRemovedRequirements()` to use `capability::normalized_name`
- [ ] 1.7 Update key generation in `validateRenamedRequirements()` for both FROM and TO names

## 2. Testing

- [ ] 2.1 Add test case for same-named REMOVED requirements across different capabilities
- [ ] 2.2 Add test case for same-named MODIFIED requirements across different capabilities
- [ ] 2.3 Add test case for same-named ADDED requirements across different capabilities (new specs)
- [ ] 2.4 Verify existing tests still pass (duplicate detection within same file)
- [ ] 2.5 Run `go test ./internal/validation/...` to verify all tests pass

## 3. Validation

- [ ] 3.1 Validate the `redesign-provider-architecture` change with the fix applied
- [ ] 3.2 Verify no regression in existing validation behavior
