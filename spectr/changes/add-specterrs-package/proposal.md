# Change: Add internal/specterrs Package for Centralized Error Handling

## Why

Error definitions are scattered across 10+ files with only 1 sentinel error (`ErrUserCancelled`) and 1 error constant (`errEmptyPath`). The remaining ~35 `errors.New()` calls use inline strings, making errors hard to find, test, and maintain consistently.

## What Changes

- **ADDED**: New `internal/specterrs/` package with domain-organized custom error types
- **ADDED**: 21 custom error types across 8 domain files (git, archive, validation, initialize, list, environment, pr, doc)
- **MODIFIED**: All ~35 `errors.New()` usages migrated to custom types
- **REMOVED**: Sentinel error `ErrUserCancelled` from `internal/archive/types.go`
- **REMOVED**: Error constant `errEmptyPath` from `internal/initialize/filesystem.go`

## Impact

- Affected specs: No existing specs modified; added: `error-handling`
- Affected code:
  - `internal/specterrs/` (new package)
  - `internal/archive/` (11 errors migrated)
  - `internal/git/` (5 errors migrated)
  - `internal/validation/` (3 errors migrated)
  - `internal/initialize/` (3 errors migrated)
  - `internal/pr/` (2 errors migrated)
  - `internal/list/` (2 errors migrated)
  - `cmd/` (9 errors across validate.go, init.go, list.go, accept.go)
