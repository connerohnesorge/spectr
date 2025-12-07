## 1. Core Infrastructure

- [ ] 1.1 Add `terminalWidth` field to `interactiveModel` struct to track current terminal width
- [ ] 1.2 Handle `tea.WindowSizeMsg` in `Update()` to capture initial and resize terminal dimensions
- [ ] 1.3 Define column priority constants and width breakpoint thresholds

## 2. Responsive Column Logic

- [ ] 2.1 Create `calculateResponsiveColumns()` function that returns appropriate columns based on terminal width
- [ ] 2.2 Implement priority-based column hiding for narrow terminals
- [ ] 2.3 Implement dynamic truncation thresholds based on available width

## 3. Apply to All Table Views

- [ ] 3.1 Update `RunInteractiveChanges()` to use responsive columns
- [ ] 3.2 Update `RunInteractiveSpecs()` to use responsive columns
- [ ] 3.3 Update `RunInteractiveAll()` (unified mode) to use responsive columns
- [ ] 3.4 Update `RunInteractiveArchive()` to use responsive columns

## 4. Dynamic Resize Handling

- [ ] 4.1 Rebuild table with new column configuration when terminal is resized
- [ ] 4.2 Preserve cursor position and selection state during resize
- [ ] 4.3 Update footer/help text to reflect current visible columns

## 5. Testing and Validation

- [ ] 5.1 Add unit tests for `calculateResponsiveColumns()` at various widths
- [ ] 5.2 Test column visibility at each breakpoint threshold
- [ ] 5.3 Manual testing on different terminal sizes (80, 100, 120+ columns)
