# Change: Add Unified Interactive List Mode

## Why

Currently, the interactive mode only works in separate contexts: `spectr list --interactive` shows changes, and `spectr list --specs --interactive` shows specs. Users cannot see both changes and specs together in a single interactive session. This makes it harder to understand the full project landscape and navigate between related items. A unified interactive list mode would allow users to browse and interact with all item types (changes and specs) in a single navigable interface.

## What Changes

- Add a unified interactive mode that displays both changes and specs together
- Implement item type filtering/toggling within the interactive interface
- Support viewing mixed lists with clear visual indicators for item type
- Maintain support for separate changes-only and specs-only modes
- Extend interactive navigation to support both item types seamlessly

## Impact

- **Affected specs**: `cli-interface` (adds new interaction patterns)
- **Affected code**: `internal/list/interactive.go`, `cmd/list.go` (extends interactive mode)
- **New capability**: Unified item browsing across changes and specs
- **Backward compatible**: Existing `--interactive` behavior preserved with `--all` flag for unified mode

## Benefits

- Single command to browse entire project landscape
- Better navigation and discovery of related changes and specs
- Reduced context switching between different list types
- Improved user experience for project exploration

## Migration Notes

- Default behavior unchanged (separate lists)
- Users opt-in to unified mode with new `--all` flag
- No breaking changes to existing workflows
