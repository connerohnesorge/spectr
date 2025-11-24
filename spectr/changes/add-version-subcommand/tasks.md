# Implementation Tasks

## 1. Core Implementation
- [ ] 1.1 Create `VERSION` file at project root with JSON format (version, commit, date fields)
- [ ] 1.2 Create version package/file that embeds VERSION using Go's embed package
- [ ] 1.3 Create `cmd/version.go` implementing VersionCmd that reads from embedded file
- [ ] 1.4 Implement parsing logic for VERSION JSON file with error handling
- [ ] 1.5 Add VersionCmd field to CLI struct in `cmd/root.go`
- [ ] 1.6 Implement `--short` flag for version-only output
- [ ] 1.7 Implement `--json` flag for JSON-formatted output
- [ ] 1.8 Handle missing/empty VERSION file gracefully (show 'dev' version)

## 2. Build Integration
- [ ] 2.1 Update `.goreleaser.yaml` to generate VERSION file in before hooks
- [ ] 2.2 Update `flake.nix` to generate VERSION file during Nix build phase
- [ ] 2.3 Test local build by creating VERSION file and running `go build`
- [ ] 2.4 Verify GoReleaser builds generate VERSION file correctly
- [ ] 2.5 Verify Nix builds generate VERSION file correctly
- [ ] 2.6 Document VERSION file format and generation in code comments

## 3. Testing
- [ ] 3.1 Create `cmd/version_test.go` with table-driven tests
- [ ] 3.2 Test VERSION file parsing (valid JSON, invalid JSON, missing file)
- [ ] 3.3 Test default output format reads from embedded file
- [ ] 3.4 Test `--short` flag output
- [ ] 3.5 Test `--json` flag output and JSON structure
- [ ] 3.6 Test behavior when VERSION file is missing or has empty values
- [ ] 3.7 Add test for Run() method existence (following existing pattern)

## 4. Documentation
- [ ] 4.1 Add doc comments to VersionCmd struct
- [ ] 4.2 Add doc comments to Run() method
- [ ] 4.3 Add usage examples in code comments
- [ ] 4.4 Ensure `spectr version --help` displays clear help text
- [ ] 4.5 Document VERSION file JSON format in code comments
- [ ] 4.6 Document embed usage and file location requirements

## 5. Validation
- [ ] 5.1 Run `go build` to verify compilation
- [ ] 5.2 Run `spectr version` to test default output
- [ ] 5.3 Run `spectr version --short` to verify short output
- [ ] 5.4 Run `spectr version --json` to verify JSON output
- [ ] 5.5 Run `go test ./cmd/...` to verify all tests pass
- [ ] 5.6 Run `golangci-lint run` to verify no linting issues
- [ ] 5.7 Verify VERSION file is created correctly in builds
- [ ] 5.8 Test that embedded file is accessible at runtime
