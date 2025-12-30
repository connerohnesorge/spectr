# Change: Intelligent Partial ID Matching for `spectr archive <id>`

## Why

Users currently must type the exact full change ID when running `spectr archive
<id>`. Change IDs like `refactor-unified-interactive-tui` are long and
error-prone to type. Partial matching would let users type `refactor` or
`unified` and have the system resolve to the correct change, improving UX
without sacrificing precision.

## What Changes

- Add partial ID matching to `spectr archive` when a non-exact ID is provided
- Match algorithm: prefix match first, then substring match
- Require unique match (error if multiple changes match the partial)
- Display matched ID to user before proceeding with confirmation
- No changes to interactive mode (already works well for discovery)
- Add VHS tape demo showing partial ID matching in action

## Impact

- Affected specs: `cli-interface` (archive command behavior)
- Affected code: `internal/archive/archiver.go`, `internal/discovery/changes.go`
- Affected docs: `assets/vhs/partial-match.tape` (new),
  `assets/gifs/partial-match.gif` (generated)
