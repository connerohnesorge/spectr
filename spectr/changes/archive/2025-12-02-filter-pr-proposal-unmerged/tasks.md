# Implementation Tasks

## 1. Implementation

- [x] 1.1 Add `PathExistsOnRef(ref, path string) (bool, error)` function to
  `internal/git/branch.go`
- [x] 1.2 Add `FilterChangesNotOnRef(changes []ChangeInfo, ref string)
  ([]ChangeInfo, error)` function to filter changes
- [x] 1.3 Modify `cmd/pr.go` `selectChangeInteractive()` to filter changes when
  called from proposal command
- [x] 1.4 Add context parameter or separate function to distinguish proposal vs
  archive interactive selection
- [x] 1.5 Display helpful message when no unmerged proposals exist

## 2. Testing

- [x] 2.1 Add unit tests for `PathExistsOnRef` function
- [x] 2.2 Add unit tests for filtering logic
- [x] 2.3 Add integration test verifying interactive list only shows unmerged
  changes
- [x] 2.4 Test edge case: all changes already on main
- [x] 2.5 Test edge case: no changes exist

## 3. Validation

- [x] 3.1 Run `spectr validate filter-pr-proposal-unmerged --strict`
- [x] 3.2 Verify existing archive interactive behavior is unchanged
- [x] 3.3 Manual testing with real repository
