# Implementation Tasks

## 1. Core Implementation
- [ ] 1.1 Add version variables to `main.go` (version, commit, date, goVersion)
- [ ] 1.2 Create `cmd/version.go` implementing VersionCmd struct with flags
- [ ] 1.3 Implement `Run()` method to display version information
- [ ] 1.4 Add VersionCmd field to CLI struct in `cmd/root.go`
- [ ] 1.5 Implement `--short` flag for version-only output
- [ ] 1.6 Implement `--json` flag for JSON-formatted output
- [ ] 1.7 Handle "dev" version gracefully when ldflags not provided

## 2. Build Integration
- [ ] 2.1 Update `.goreleaser.yaml` to inject version metadata via ldflags
- [ ] 2.2 Test local build with `go build -ldflags` to verify injection
- [ ] 2.3 Verify GoReleaser builds include correct version information
- [ ] 2.4 Document build process in code comments

## 3. Testing
- [ ] 3.1 Create `cmd/version_test.go` with table-driven tests
- [ ] 3.2 Test default output format (human-readable)
- [ ] 3.3 Test `--short` flag output
- [ ] 3.4 Test `--json` flag output and JSON structure
- [ ] 3.5 Test behavior when version variables are empty/default
- [ ] 3.6 Add test for Run() method existence (following existing pattern)

## 4. Documentation
- [ ] 4.1 Add doc comments to VersionCmd struct
- [ ] 4.2 Add doc comments to Run() method
- [ ] 4.3 Add usage examples in code comments
- [ ] 4.4 Ensure `spectr version --help` displays clear help text

## 5. Validation
- [ ] 5.1 Run `go build` to verify compilation
- [ ] 5.2 Run `spectr version` to test default output
- [ ] 5.3 Run `spectr version --short` to verify short output
- [ ] 5.4 Run `spectr version --json` to verify JSON output
- [ ] 5.5 Run `go test ./cmd/...` to verify all tests pass
- [ ] 5.6 Run `golangci-lint run` to verify no linting issues
