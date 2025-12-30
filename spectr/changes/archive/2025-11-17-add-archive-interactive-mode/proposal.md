# Change: Add Interactive Table Mode to Archive Command

## Why

Currently, the archive command uses a basic numbered list selection when no
change ID is provided, which is inconsistent with the list command's polished
interactive table interface. Users have a better experience with the list -I
interface (keyboard navigation, visual styling, clipboard support), and this
same UX should be available when selecting changes to archive.

## What Changes

- Replace the basic text-based selection in archive command with the same
  interactive table interface used by `list -I`
- Add an `-I` or `--interactive` flag to the archive command to explicitly
  trigger interactive mode
- When no change ID is provided and the flag is not present, default to
  interactive mode for better UX
- Reuse existing interactive table components from the list package
- Display change metadata (ID, title, deltas, tasks) in the table to help users
  make informed decisions

## Impact

- Affected specs: cli-interface
- Affected code:
  - cmd/archive.go (add Interactive flag, update Run logic)
  - internal/archive/archiver.go (replace selectChange with interactive table
    call)
  - internal/list/interactive.go (potentially expose or adapt
    RunInteractiveChanges for archive use case)
- User benefit: More consistent and polished UX across all interactive commands
- No breaking changes: Existing `spectr archive <change-id>` usage remains
  unchanged
