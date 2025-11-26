## 1. Implementation

- [ ] 1.1 Add `searchMode` bool and `searchQuery` string fields to `interactiveModel` struct
- [ ] 1.2 Add `textinput` component from bubbles for search input
- [ ] 1.3 Implement '/' key handler to toggle search mode
- [ ] 1.4 Implement Escape key handler to exit search mode
- [ ] 1.5 Add filtering logic to filter rows by ID and title matching search query
- [ ] 1.6 Update `View()` to render search input when search mode is active
- [ ] 1.7 Update help text to show '/' hotkey for all modes

## 2. Testing

- [ ] 2.1 Add unit tests for search mode toggle behavior
- [ ] 2.2 Add unit tests for row filtering logic
- [ ] 2.3 Add integration tests using teatest for search workflow
- [ ] 2.4 Run existing tests to ensure no regressions

## 3. Validation

- [ ] 3.1 Run `go test ./internal/list/...` to verify all tests pass
- [ ] 3.2 Run `spectr validate add-search-hotkey --strict` to verify spec validity
