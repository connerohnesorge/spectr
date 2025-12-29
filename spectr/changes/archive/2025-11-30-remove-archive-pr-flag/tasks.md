# Implementation Tasks

## 1. Remove PR Flag from Archive Command

- [x] 1.1 Remove `PR bool` field from `ArchiveCmd` struct in
  `internal/archive/cmd.go`
- [x] 1.2 Remove PR workflow logic from `Archive()` function in
  `internal/archive/archiver.go` (lines 130-147)
- [x] 1.3 Remove `PRContext` struct and `createPR` function call

## 2. Delete PR-Related Archive Files

- [x] 2.1 Delete `internal/archive/pr.go`
- [x] 2.2 Delete `internal/archive/pr_format.go`
- [x] 2.3 Delete `internal/archive/pr_test.go`

## 3. Delete Git Package

- [x] 3.1 Delete `internal/git/pr.go`
- [x] 3.2 Delete `internal/git/platform.go`
- [x] 3.3 Delete `internal/git/platform_test.go`
- [x] 3.4 Delete `internal/git/operations.go`
- [x] 3.5 Remove `internal/git/` directory

## 4. Update Dependencies

- [x] 4.1 Remove `github.com/google/uuid` from `go.mod`
- [x] 4.2 Run `go mod tidy` to clean up dependencies

## 5. Verify Changes

- [x] 5.1 Run `go build .` to ensure clean compilation
- [x] 5.2 Run `go test ./...` to verify all tests pass
- [x] 5.3 Run `go run . archive --help` to verify `--pr` flag is removed
- [x] 5.4 Run `go run . validate --strict` to validate specs
