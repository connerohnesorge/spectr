# Change: Fix projectPath Not Set in Interactive Changes and Add Display

## Why

The `projectPath` field in `interactiveModel` is not being set when `RunInteractiveChanges()` is called, unlike `RunInteractiveSpecs()` which correctly receives and sets it. This creates an inconsistency where the field is defined on the struct but never populated for changes mode. Additionally, users benefit from seeing the project path in the TUI to understand which project they are working with, especially when working across multiple Spectr projects.

## What Changes

- Update `RunInteractiveChanges()` signature to accept `projectPath string` parameter
- Pass `projectPath` from `cmd/list.go` to `RunInteractiveChanges()`
- Set `projectPath` on the `interactiveModel` in both `RunInteractiveChanges()` and `RunInteractiveArchive()`
- Add project path display to the TUI help text or header for both changes and specs modes

## Impact

- Affected specs: cli-interface
- Affected code:
  - `internal/list/interactive.go` - Update `RunInteractiveChanges()` and `RunInteractiveArchive()` signatures and model initialization
  - `cmd/list.go` - Pass `projectPath` to `RunInteractiveChanges()`
  - `cmd/archive.go` - Pass `projectPath` to `RunInteractiveArchive()` if used there
- Breaking changes: **BREAKING** - Function signature changes for `RunInteractiveChanges()` and `RunInteractiveArchive()`
- User benefit: Consistency in model initialization and visual confirmation of working project path
