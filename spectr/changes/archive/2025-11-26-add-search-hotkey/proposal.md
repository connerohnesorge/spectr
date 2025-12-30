# Change: Add Search Hotkey to Interactive Lists

## Why

Users browsing long lists of changes or specs need a quick way to filter items
by text search. Currently, the only way to find a specific item is to scroll
through the entire list, which becomes cumbersome with many items.

## What Changes

- Add a '/' hotkey that activates a text search mode in interactive lists
- When search mode is active, display a text input field where users can type
  search queries
- Filter the table rows in real-time as the user types, matching against ID and
  title columns
- Allow exiting search mode with Escape to return to the full list
- Persist the search query until explicitly cleared (pressing '/' again when
  empty clears search)

## Impact

- Affected specs: cli-interface
- Affected code: internal/list/interactive.go, internal/list/interactive_test.go
