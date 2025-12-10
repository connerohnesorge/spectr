# Change: Add Provider Search/Filter to Init Wizard

## Why

The init wizard's provider selection screen currently lists all available AI tool providers in a flat list. As the number of supported providers grows (currently 15+), users need a way to quickly find specific tools without manually scrolling through the entire list. The existing search functionality in `spectr list` interactive modes provides a proven pattern that can be adapted for the provider selection screen.

## What Changes

- Add `/` hotkey to activate search mode in the provider selection step of the init wizard
- Filter providers in real-time by matching search query against provider name (case-insensitive)
- Preserve selection state when filtering (checked providers remain checked even when filtered out)
- Allow navigation and selection within filtered results
- Exit search mode with Escape to restore full list

## Impact

- Affected specs: `cli-interface`
- Affected code: `internal/initialize/wizard.go`
- No breaking changes - existing navigation and selection behavior preserved
- Follows established search pattern from `spectr list` interactive modes
