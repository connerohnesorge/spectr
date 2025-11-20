# Change: Move ArchiveCmd struct to internal/archive package

## Why

The `ArchiveCmd` struct is currently defined in `cmd/archive.go`, which separates the command definition from its implementation logic in `internal/archive/archiver.go`. Moving the struct into the `internal/archive/` package improves package cohesion by keeping all archive-related code together, reduces files in the thin `cmd/` layer, and establishes a clearer pattern where each internal package owns its CLI command definition. The `Run()` method can then delegate to `Archive()` by passing itself as an argument, keeping command-specific logic encapsulated.

## What Changes

- Move `ArchiveCmd` struct definition from `cmd/archive.go` to `internal/archive/cmd.go`
- Update `ArchiveCmd.Run()` method to call `Archive(cmd *ArchiveCmd)` passing itself as argument
- Modify `Archive()` signature from `Archive(changeID string)` to `Archive(cmd *ArchiveCmd)` to accept the full command struct
- Update `cmd/root.go` to import `ArchiveCmd` from the `archive` package using package-qualified name
- Remove `cmd/archive.go` entirely
- No changes to CLI behavior, flags, or user-facing interface
- No changes to spec requirements

## Impact

- **Affected code**: `cmd/root.go`, `cmd/archive.go` (deleted), `internal/archive/cmd.go` (new), `internal/archive/archiver.go` (Archive method signature)
- **User impact**: None - CLI interface remains identical
- **Architecture**: Improves package organization and encapsulation by keeping command logic with implementation
- **Testing**: Existing tests need updates to pass the ArchiveCmd struct to Archive()
