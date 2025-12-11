## 1. Implementation

- [ ] 1.1 Add search state fields to `WizardModel` (searchMode, searchQuery, filteredProviders)
- [ ] 1.2 Implement search mode activation with `/` key in `handleSelectKeys`
- [ ] 1.3 Add text input handling for search query when in search mode
- [ ] 1.4 Implement provider filtering logic matching against provider name (case-insensitive)
- [ ] 1.5 Update `renderSelect` to show search input field when active
- [ ] 1.6 Update `renderProviderGroup` to render filtered providers instead of all providers
- [ ] 1.7 Implement search mode exit with Escape key (restore full list)
- [ ] 1.8 Adjust cursor position when filtering (select first match or maintain valid position)
- [ ] 1.9 Preserve selection state across filter changes (don't deselect hidden items)
- [ ] 1.10 Update help text to include `/: search` when not in search mode

## 2. Testing

- [ ] 2.1 Add unit tests for provider filtering logic
- [ ] 2.2 Add tests for cursor adjustment when filtered list changes
- [ ] 2.3 Add tests for selection preservation during filtering
- [ ] 2.4 Manual testing of search flow in init wizard TUI

## 3. Documentation

- [ ] 3.1 Update spectr/AGENTS.md if init wizard help references are affected
