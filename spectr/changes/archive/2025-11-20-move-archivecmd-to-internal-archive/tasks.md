# Implementation Tasks

## 1. Code Reorganization
- [ ] 1.1 Create `internal/archive/cmd.go` with ArchiveCmd struct definition
- [ ] 1.2 Copy ArchiveCmd struct from `cmd/archive.go` to `internal/archive/cmd.go`
- [ ] 1.3 Implement Run() method on ArchiveCmd that calls archive.RunArchive(cmd)
- [ ] 1.4 Add package comment and export doc comments to ArchiveCmd in new location
- [ ] 1.5 Update imports in `internal/archive/cmd.go` if needed (fmt, etc.)

## 2. Update Archive Implementation
- [ ] 2.1 Change `Archive(changeID string)` signature to `Archive(cmd *ArchiveCmd)` in `internal/archive/archiver.go`
- [ ] 2.2 Update all references to flags inside Archive() to use `cmd.Yes`, `cmd.SkipSpecs`, `cmd.NoValidate`, `cmd.Interactive`, `cmd.PR`
- [ ] 2.3 Remove `NewArchiver()` function and related flag parameters as they're now in ArchiveCmd
- [ ] 2.4 Update Archive() implementation to extract changeID from `cmd.ChangeID`

## 3. Update CLI Root Command
- [ ] 3.1 Update `cmd/root.go` imports to use `archive "github.com/conneroisu/spectr/internal/archive"`
- [ ] 3.2 Change `Archive ArchiveCmd` field to `Archive archive.ArchiveCmd` in CLI struct
- [ ] 3.3 Verify all other imports in `cmd/root.go` remain correct

## 4. Update Tests
- [ ] 4.1 Update `internal/archive/archiver_test.go` to pass ArchiveCmd to Archive() method
- [ ] 4.2 Update any test fixtures that create Archiver instances
- [ ] 4.3 Verify all archive-related tests still pass

## 5. Cleanup
- [ ] 5.1 Delete `cmd/archive.go` file entirely
- [ ] 5.2 Verify no other files import from `cmd/archive.go`

## 6. Testing & Validation
- [ ] 6.1 Run `go build ./...` to ensure clean compilation
- [ ] 6.2 Run `go test ./...` to ensure all tests pass
- [ ] 6.3 Verify archive command still works: `spectr archive --help`
- [ ] 6.4 Run full lint with `nix develop -c 'lint'` to verify no violations
- [ ] 6.5 Validate with `spectr validate --strict` if applicable
