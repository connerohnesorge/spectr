# Change: Require tasks.json for Archive

## Why

The current `spectr archive` command has auto-accept behavior that converts `tasks.md` to `tasks.json` automatically during archive. This is problematic because:

1. **Workflow enforcement**: The `spectr accept` command exists to formally "accept" a proposal for implementation. Auto-accepting during archive bypasses this intentional workflow gate.
2. **Implementation tracking**: Changes should be accepted with `spectr accept` before implementation begins, not at archive time. Auto-accept allows skipping the formal acceptance step entirely.
3. **Clearer separation of concerns**: Archive should only archive completed changes, not also handle format conversion.

## What Changes

- Remove the "Auto-Accept on Archive" requirement from archive-workflow spec
- Add new requirement that archive SHALL require `tasks.json` to exist
- Archive SHALL display an actionable error message directing users to run `spectr accept <change-id>` first
- Archive SHALL block when neither tasks.md nor tasks.json exists (no tasks at all)

## Impact

- Affected specs: `archive-workflow`
- Affected code:
  - `internal/archive/archiver.go` (remove auto-accept, add tasks.json check)
