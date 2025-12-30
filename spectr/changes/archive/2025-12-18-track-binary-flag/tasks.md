# Implementation Tasks

## 1. Implementation

- [ ] 1.1 Add `IncludeBinaries` flag to `TrackCmd` struct in `cmd/track.go`
- [ ] 1.2 Add `IncludeBinaries` field to `track.Config` struct in
  `internal/track/tracker.go`
- [ ] 1.3 Pass flag through from TrackCmd.Run to Config in `cmd/track.go`
- [ ] 1.4 Add binary detection function to `internal/track/committer.go` using
  git diff --numstat
- [ ] 1.5 Update `filterTaskFiles` function to also filter binary files when
  flag is false
- [ ] 1.6 Add warning message when binary files are skipped
- [ ] 1.7 Update `Committer` struct to accept IncludeBinaries setting
- [ ] 1.8 Pass IncludeBinaries from Config to Committer in tracker.go

## 2. Testing

- [ ] 2.1 Add unit tests for binary detection in `committer_test.go`
- [ ] 2.2 Add unit tests for filtering with IncludeBinaries=false
- [ ] 2.3 Add unit tests for including binaries with IncludeBinaries=true
- [ ] 2.4 Add integration test for track command with binary files
- [ ] 2.5 Test with various binary file types (images, executables, archives)

## 3. Validation

- [ ] 3.1 Run `go test ./...` to ensure all tests pass
- [ ] 3.2 Run `golangci-lint run` to ensure code quality
- [ ] 3.3 Manual testing with real binary files
- [ ] 3.4 Verify default behavior excludes binaries
- [ ] 3.5 Verify --include-binaries flag includes binaries
