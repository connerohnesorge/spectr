# Implementation Tasks

## 1. Core Implementation
- [x] 1.1 Create `VERSION` file at project root with JSON format (version, commit, date fields)
- [x] 1.2 Create version package/file that embeds VERSION using Go's embed package
- [x] 1.3 Create `cmd/version.go` implementing VersionCmd that reads from embedded file
- [x] 1.4 Implement parsing logic for VERSION JSON file with error handling
- [x] 1.5 Add VersionCmd field to CLI struct in `cmd/root.go`
- [x] 1.6 Implement `--short` flag for version-only output
- [x] 1.7 Implement `--json` flag for JSON-formatted output
- [x] 1.8 Handle missing/empty VERSION file gracefully (show 'dev' version)

## 2. Build Integration
- [x] 2.1 Update `.goreleaser.yaml` to generate VERSION file in before hooks
- [x] 2.2 Update `flake.nix` to generate VERSION file during Nix build phase
- [x] 2.3 Test local build by creating VERSION file and running `go build`
- [x] 2.4 Verify GoReleaser builds generate VERSION file correctly
- [x] 2.5 Verify Nix builds generate VERSION file correctly
- [x] 2.6 Document VERSION file format and generation in code comments

## 3. Testing
- [x] 3.1 Create `cmd/version_test.go` with table-driven tests
- [x] 3.2 Test VERSION file parsing (valid JSON, invalid JSON, missing file)
- [x] 3.3 Test default output format reads from embedded file
- [x] 3.4 Test `--short` flag output
- [x] 3.5 Test `--json` flag output and JSON structure
- [x] 3.6 Test behavior when VERSION file is missing or has empty values
- [x] 3.7 Add test for Run() method existence (following existing pattern)

## 4. Documentation
- [x] 4.1 Add doc comments to VersionCmd struct
- [x] 4.2 Add doc comments to Run() method
- [x] 4.3 Add usage examples in code comments
- [x] 4.4 Ensure `spectr version --help` displays clear help text
- [x] 4.5 Document VERSION file JSON format in code comments
- [x] 4.6 Document embed usage and file location requirements

## 5. Validation
- [x] 5.1 Run `go build` to verify compilation
- [x] 5.2 Run `spectr version` to test default output
- [x] 5.3 Run `spectr version --short` to verify short output
- [x] 5.4 Run `spectr version --json` to verify JSON output
- [x] 5.5 Run `go test ./cmd/...` to verify all tests pass
- [x] 5.6 Run `golangci-lint run` to verify no linting issues
- [x] 5.7 Verify VERSION file is created correctly in builds
- [x] 5.8 Test that embedded file is accessible at runtime
