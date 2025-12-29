# Change: Add Archive Command for Completing and Applying Changes

## Why

After implementing and deploying a Spectr change, users currently have no automated way to:

1. Move the change to an archive with a dated timestamp
2. Apply delta specs to the main specification files
3. Validate that the change is complete and correct before archiving
4. Track which changes have been deployed

OpenSpec provides a robust archive command (`openspec archive`) that handles all these tasks systematically. Spectr should have equivalent functionality to complete the change lifecycle and maintain spec-change consistency.

## What Changes

- Add new `archive` command that moves completed changes to `spectr/changes/archive/YYYY-MM-DD-<change-id>/`
- Implement delta spec application that merges ADDED/MODIFIED/REMOVED/RENAMED requirements into main specs
- Add pre-archive validation of change delta specs and proposal
- Support interactive change selection when no change-id provided
- Add task completion checking with warnings for incomplete tasks
- Implement spec update workflow:
  - Parse delta operations from `changes/<id>/specs/*/spec.md`
  - Apply operations to corresponding files in `spectr/specs/*/spec.md`
  - Create new spec files if they don't exist
  - Validate rebuilt specs before writing
  - Display operation counts (added, modified, removed, renamed)
- Add command-line flags:
  - `-y, --yes` - Skip confirmation prompts for non-interactive usage
  - `--skip-specs` - Skip spec updates (for tooling/infrastructure-only changes)
  - `--no-validate` - Skip validation (requires confirmation, not recommended)
- Prevent duplicate requirement names within delta sections
- Prevent cross-section conflicts (e.g., same requirement in ADDED and MODIFIED)
- Handle RENAMED operations correctly with MODIFIED operations
- Archive naming format: `YYYY-MM-DD-<change-id>` for chronological sorting

## Impact

- **Affected specs**: `cli-framework` (extends existing command set)
- **New spec**: `archive-workflow` - Defines archive command behavior and spec merging
- **Affected code**:
  - `cmd/root.go` - Add ArchiveCmd struct to CLI
  - New `cmd/archive.go` - Archive command implementation
  - New `internal/archive/` - Archive functionality package
    - `archiver.go` - Core archive logic
    - `spec_merger.go` - Delta spec parsing and merging
    - `validator.go` - Pre-archive validation
    - `types.go` - Data structures for archive operations
  - Extend `internal/parsers/` - Add delta parsing
    - `delta_parser.go` - Parse ADDED/MODIFIED/REMOVED/RENAMED sections
    - `requirement_parser.go` - Parse requirement blocks with scenarios
  - Extend `internal/discovery/` - Add spec finding
    - `specs.go` - Find main specs that need updates

## Benefits

- **Complete workflow**: Users can now scaffold → implement → validate → archive changes
- **Spec consistency**: Automated merging prevents manual errors
- **Audit trail**: Dated archives provide deployment history
- **Safe defaults**: Validation prevents invalid specs from being archived
- **Flexibility**: Flags support both interactive and CI/CD usage
- **Alignment**: Matches OpenSpec's proven archive command semantics

## Migration Notes

This is a new capability. Existing changes can be archived using:

```bash
spectr archive <change-id>
```

For changes that only modify tooling/documentation without spec updates:

```bash
spectr archive <change-id> --skip-specs
```
