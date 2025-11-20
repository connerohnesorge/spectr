# Design: Move ArchiveCmd to Internal Archive Package

## Context

The current CLI structure separates command definitions (`cmd/`) from their implementation logic (`internal/`). While this provides clean separation of concerns, it creates a pattern where related code is spread across two locations:

- `cmd/archive.go` - Command struct definition
- `internal/archive/archiver.go` - Implementation logic

This is a pure refactoring with no functional changes.

## Goals

1. **Improve package cohesion** - Keep all archive-related code (struct + logic) in the same package
2. **Reduce cmd/ complexity** - cmd/ should contain only the CLI root and command registrations
3. **Establish clearer pattern** - Each internal package owns its command struct definition
4. **Minimize file count** - Reduce unnecessary files in cmd/ directory

## Non-Goals

- Changing command behavior or flags
- Modifying CLI interface
- Refactoring other commands (this sets pattern for future changes)

## Decisions

### Decision 1: Place ArchiveCmd in internal/archive/cmd.go

**What**: Create `internal/archive/cmd.go` for the ArchiveCmd struct and Run() method.

**Rationale**:
- Keeps all archive functionality in one package
- New file name clearly indicates it's a command definition
- Follows Go package naming conventions

**Alternatives considered**:
- Place in archiver.go directly - Could cause archiver.go to be too large with mixed concerns
- Place in types.go - cmd.go is clearer and more maintainable

### Decision 2: Archive() accepts full ArchiveCmd struct

**What**: Change `Archive(changeID string)` signature to `Archive(cmd *ArchiveCmd)` so the method receives the entire command struct with all flags.

**Rationale**:
- Encapsulates all command-specific logic in one place
- Each flag is self-contained within the struct
- Simplifies method signature (one parameter instead of multiple flags)
- Future changes to flags only require updating ArchiveCmd struct
- Makes the intent clearer - the entire command determines behavior

**Alternatives considered**:
- Keep current `Archive(changeID string)` signature - Loses encapsulation of flags
- Pass individual flags as parameters - Creates long method signature, harder to extend
- Keep NewArchiver() helper - Still separates command from implementation

### Decision 3: Use package-qualified import in root.go

**What**: Import archive package and reference `archive.ArchiveCmd` instead of `ArchiveCmd`.

**Rationale**:
- Explicit about where the type comes from
- Avoids polluting cmd package namespace
- Makes dependency relationships clear

**Alternative**: Could do type aliasing, but explicit qualification is clearer.

## Migration Path

1. Create `internal/archive/cmd.go` with ArchiveCmd definition
2. Update Archive() signature from `Archive(changeID string)` to `Archive(cmd *ArchiveCmd)`
3. Remove NewArchiver() helper function
4. Replace all flag parameters with references to cmd struct fields
5. Update archiver implementation to work with cmd struct
6. Update `cmd/root.go` imports and CLI struct
7. Update all tests to pass ArchiveCmd to Archive()
8. Delete `cmd/archive.go`
9. Verify build, tests, and lint all pass

No end-user visible changes - CLI interface remains identical.

## Risks & Mitigations

| Risk | Mitigation |
|------|-----------|
| Import cycles | `cmd/` imports `internal/archive/` (no reverse imports), safe |
| Other imports of cmd.ArchiveCmd | Check codebase - only used in cmd/root.go CLI struct |
| Test references | Update test imports if any tests reference cmd.ArchiveCmd |

## Future Pattern

This establishes a pattern for other commands:
- `InitCmd` could move to `internal/init/cmd.go`
- `ValidateCmd` could move to `internal/validation/cmd.go`
- `ListCmd` could move to `internal/list/cmd.go`

However, this change focuses only on ArchiveCmd to keep scope minimal.
