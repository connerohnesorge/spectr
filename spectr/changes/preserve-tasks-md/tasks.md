# Implementation Tasks

## 1. Modify Accept Command

- [ ] 1.1 Remove `os.Remove(tasksMdPath)` call in `cmd/accept.go` after
  tasks.jsonc is written
- [ ] 1.2 Add comment explaining why tasks.md is preserved
- [ ] 1.3 Update success message to indicate both files now exist

## 2. Update Documentation

- [ ] 2.1 Update `cmd/accept.go` command description to mention tasks.md
  preservation
- [ ] 2.2 Add note in `spectr/AGENTS.md` explaining that both files coexist
  after accept
- [ ] 2.3 Update any workflow documentation referencing tasks.md deletion

## 3. Add Validation

- [ ] 3.1 Add optional validation in `spectr validate` to check if tasks.md and
  tasks.jsonc exist but diverge
- [ ] 3.2 Show informational warning (not error) if files exist but have
  different task counts or IDs

## 4. Verify Backward Compatibility

- [ ] 4.1 Verify `internal/discovery/changes.go` correctly prefers tasks.jsonc
  when both exist
- [ ] 4.2 Test that `spectr list` works correctly with both files present
- [ ] 4.3 Test that `spectr view` works correctly with both files present
- [ ] 4.4 Test that `spectr archive` works correctly with both files present

## 5. Testing

- [ ] 5.1 Update `cmd/accept_test.go` to verify tasks.md is NOT deleted
- [ ] 5.2 Add test case for accept with both tasks.md and tasks.jsonc present
- [ ] 5.3 Run full test suite: `go test ./...`
- [ ] 5.4 Manual integration test: accept a real change and verify both files
  exist

## 6. Edge Cases

- [ ] 6.1 Test behavior when tasks.md already exists from previous accept
  (should overwrite or preserve?)
- [ ] 6.2 Document behavior when user manually deletes tasks.jsonc after accept
- [ ] 6.3 Consider adding `--sync-from-md` flag to re-generate tasks.jsonc from
  tasks.md

## 7. Template and Documentation Updates

- [ ] 7.1 Update `internal/domain/templates/slash-proposal.md.tmpl` Step 6 to
  instruct LLM to add header warning in tasks.md
- [ ] 7.2 Update `internal/domain/templates/slash-proposal.toml.tmpl` with same
  modification
- [ ] 7.3 Update `internal/domain/templates/slash-apply.md.tmpl` to add reminder
  not to edit tasks.md
- [ ] 7.4 Update `internal/domain/templates/slash-apply.toml.tmpl` with same
  modification
- [ ] 7.5 Update `cmd/accept_writer.go` tasksJSONHeader constant to reference
  AGENTS.md
- [ ] 7.6 Add "Task File Management Workflow" section to `spectr/AGENTS.md` with
  workflow diagram and examples
