# Tasks

## Implementation

- [ ] Add `HasUncommittedChanges(changeDir string) bool` function in `internal/discovery/changes.go` that uses `git status --porcelain` to detect uncommitted files in a change directory
- [ ] Add `uncommittedFilter bool` field to `interactiveModel` struct in `internal/list/interactive.go`
- [ ] Add `uncommittedChanges map[string]bool` field to cache git status results for performance
- [ ] Implement `handleUncommittedFilter()` method to toggle the filter state
- [ ] Add 'h' key handler in the `Update()` method that calls `handleUncommittedFilter()`
- [ ] Implement `applyUncommittedFilter()` method that filters `allRows` based on:
  - Git uncommitted status (from cached map)
  - Complete tasks (TaskStatus.Completed == TaskStatus.Total && Total > 0)
- [ ] Update help text in `RunInteractiveChanges()` to include `h: uncommitted`
- [ ] Update help text in `RunInteractiveAll()` to include `h: uncommitted`
- [ ] Update minimal footer to show "filter: uncommitted+complete" when filter is active
- [ ] Add handling for empty filter results (display "No uncommitted changes with complete tasks")
- [ ] Ignore 'h' key press in specs-only mode (`itemTypeSpec`)
- [ ] Ensure search filter and uncommitted filter work together correctly

## Testing

- [ ] Add unit test for `HasUncommittedChanges()` function with mock git output
- [ ] Add test for uncommitted filter toggle behavior in interactive model
- [ ] Add test for filter criteria (uncommitted AND complete tasks)
- [ ] Add test for empty filter results message
- [ ] Add test that 'h' is ignored in specs mode
- [ ] Add test for combined search + uncommitted filter
- [ ] Run `go test ./internal/list/...` to verify all tests pass
- [ ] Run `go test ./internal/discovery/...` to verify discovery tests pass

## Validation

- [ ] Run `spectr validate add-uncommitted-filter-hotkey --strict` to verify spec compliance
- [ ] Manual test: Run `spectr list -I`, press 'h' with changes that have uncommitted modifications
- [ ] Manual test: Verify filter toggles off when pressing 'h' again
- [ ] Manual test: Verify help text shows 'h' hotkey after pressing '?'
- [ ] Manual test: Verify 'h' is ignored in `spectr list --specs -I`
