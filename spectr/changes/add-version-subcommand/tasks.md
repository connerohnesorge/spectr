## 1. Implementation

- [x] 1.1 Create `internal/version/version.go` with version, commit, and date variables
- [x] 1.2 Create `cmd/version.go` with VersionCmd struct and Run method
- [x] 1.3 Add VersionCmd field to CLI struct in `cmd/root.go`
- [x] 1.4 Update `.goreleaser.yaml` to inject ldflags for version, commit, and date
- [x] 1.5 Update `flake.nix` to inject ldflags for version in buildGoModule

## 2. Testing

- [x] 2.1 Verify `spectr version` displays version information
- [x] 2.2 Verify `spectr version --json` outputs valid JSON
- [x] 2.3 Verify `spectr version --short` outputs only version number
- [x] 2.4 Verify goreleaser builds inject version correctly (via dry run)
- [x] 2.5 Verify nix build injects version correctly
