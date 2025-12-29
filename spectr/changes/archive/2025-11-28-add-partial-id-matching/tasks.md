# Implementation Tasks

## 1. Implementation

- [x] 1.1 Add `ResolveChangeID(partialID string, projectRoot string) (string,
  error)` function to `internal/discovery/changes.go`
- [x] 1.2 Implement prefix matching logic (case-insensitive)
- [x] 1.3 Implement substring fallback when no prefix match
- [x] 1.4 Return appropriate errors: no match, multiple matches
- [x] 1.5 Integrate resolver into `Archive()` function in
  `internal/archive/archiver.go`
- [x] 1.6 Display resolved ID when partial match succeeds (e.g., "Resolved
  'unified' -> 'refactor-unified-interactive-tui'")

## 2. Testing

- [x] 2.1 Unit tests for `ResolveChangeID` with various scenarios (exact,
  prefix, substring, ambiguous, none)
- [x] 2.2 Integration test for archive command with partial ID

## 3. Documentation

- [x] 3.1 Create VHS tape `assets/vhs/partial-match.tape` demonstrating partial
  ID matching
- [x] 3.2 Generate gif with `vhs assets/vhs/partial-match.tape`
- [x] 3.3 Update docs to reference the new feature and gif

## 4. Validation

- [x] 4.1 Run `go test ./...` to ensure all tests pass
- [x] 4.2 Run `spectr validate add-partial-id-matching --strict`
