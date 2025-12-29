# Implementation Tasks

## 1. Dependencies and Setup

- [x] 1.1 Add clipboard library to go.mod (`atotto/clipboard` or use existing `termenv` for OSC 52)
- [x] 1.2 Verify `github.com/charmbracelet/bubbles/table` is available (already in go.mod)
- [x] 1.3 Run `go mod tidy` to ensure all dependencies are resolved

## 2. Add Interactive Flag to List Command

- [x] 2.1 Add `Interactive bool` field to `ListCmd` struct in `cmd/list.go` with appropriate tags
- [x] 2.2 Update flag documentation/help text for `-I`/`--interactive` flag
- [x] 2.3 Add conditional logic in `Run()` to check for interactive mode
- [x] 2.4 Ensure interactive flag is mutually validated with JSON flag (cannot use both)

## 3. Create Interactive Table Model

- [x] 3.1 Create `internal/list/interactive.go` with bubbletea model struct
- [x] 3.2 Implement `Init()` method for table initialization
- [x] 3.3 Implement `Update(msg tea.Msg)` method with keyboard handlers:
  - [x] 3.3.1 Handle arrow keys and j/k for navigation
  - [x] 3.3.2 Handle Enter key for clipboard copy and exit
  - [x] 3.3.3 Handle 'q' and Ctrl+C for quit without copy
  - [x] 3.3.4 Handle Esc key for focus toggle (optional)
- [x] 3.4 Implement `View()` method to render table with borders and styling
- [x] 3.5 Add lipgloss styling for headers, selected rows, and borders

## 4. Implement Table Data Conversion

- [x] 4.1 Create function `buildChangesTable(changes []ChangeInfo) table.Model` in `interactive.go`
- [x] 4.2 Define columns: ID, Title, Deltas, Tasks (format: "completed/total")
- [x] 4.3 Convert `[]ChangeInfo` to `[]table.Row` format
- [x] 4.4 Create function `buildSpecsTable(specs []SpecInfo) table.Model`
- [x] 4.5 Define columns for specs: ID, Title, Requirements
- [x] 4.6 Convert `[]SpecInfo` to `[]table.Row` format
- [x] 4.7 Set appropriate column widths based on terminal size

## 5. Implement Clipboard Integration

- [x] 5.1 Create helper function `copyToClipboard(text string) error`
- [x] 5.2 Implement cross-platform clipboard logic:
  - [x] 5.2.1 Use `atotto/clipboard` for standard desktop environments, OR
  - [x] 5.2.2 Use termenv OSC 52 sequences for SSH/remote sessions
- [x] 5.3 Add clipboard error handling and user-friendly error messages
- [x] 5.4 Display success message with copied ID after clipboard operation

## 6. Wire Up Interactive Mode in List Command

- [x] 6.1 In `listChanges()`, check if `Interactive` flag is set
- [x] 6.2 If interactive, call `runInteractiveChanges(changes)` instead of formatting
- [x] 6.3 In `listSpecs()`, check if `Interactive` flag is set
- [x] 6.4 If interactive, call `runInteractiveSpecs(specs)` instead of formatting
- [x] 6.5 Handle empty list case (show message and exit without starting bubbletea)

## 7. Create Integration Functions

- [x] 7.1 Create `runInteractiveChanges(changes []ChangeInfo) error` function
- [x] 7.2 Build table model, create bubbletea program, run and handle result
- [x] 7.3 Extract selected row ID from final model state
- [x] 7.4 Copy ID to clipboard and display confirmation
- [x] 7.5 Create `runInteractiveSpecs(specs []SpecInfo) error` function
- [x] 7.6 Implement same flow for specs

## 8. Testing

- [x] 8.1 Create `internal/list/interactive_test.go`
- [x] 8.2 Write unit tests for table data conversion functions
- [x] 8.3 Write tests for clipboard helper (mock clipboard operations)
- [x] 8.4 Write integration test for flag parsing
- [x] 8.5 Manual testing: Run `spectr list -I` with various data states
- [x] 8.6 Manual testing: Verify clipboard on Linux, macOS, Windows
- [x] 8.7 Manual testing: Test SSH session with OSC 52
- [x] 8.8 Manual testing: Test empty list case
- [x] 8.9 Manual testing: Test quit with 'q' and Ctrl+C

## 9. Documentation

- [x] 9.1 Update `--help` text for list command to mention `-I` flag
- [x] 9.2 Add usage examples in code comments
- [x] 9.3 Document keyboard shortcuts (Enter, q, Ctrl+C, arrows/j/k)
- [x] 9.4 Document clipboard behavior for SSH users (OSC 52 note)

## 10. Edge Cases and Polish

- [x] 10.1 Handle very long titles (truncate with ellipsis in table)
- [x] 10.2 Handle narrow terminal widths gracefully
- [x] 10.3 Ensure table height adjusts to terminal size
- [x] 10.4 Test with single-item lists
- [x] 10.5 Test with 100+ items (scrolling behavior)
- [x] 10.6 Add visual indicator for clipboard copy success
- [x] 10.7 Ensure proper cleanup on exit signals
