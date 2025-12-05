# Tasks

## Implementation

- [ ] Add `keyCopy = "c"` constant to `internal/initialize/constants.go`
- [ ] Extract the populate context prompt text as a constant (e.g., `populateContextPrompt`) in `internal/initialize/executor.go` or `constants.go`
  - Text should be raw content WITHOUT surrounding quotes: "Review spectr/project.md and help me fill in our project's tech stack, conventions, and description. Ask me questions to understand the codebase."
- [ ] Update `handleCompleteKeys()` in `internal/initialize/wizard.go` to handle the 'c' key press
  - Import `internal/tui` package for `CopyToClipboard()`
  - Call `tui.CopyToClipboard(populateContextPrompt)` when 'c' is pressed
  - On success: immediately return tea.Quit (silent exit, no message)
  - On error: display error message and return without exiting (allow retry)
  - Note: Success behavior matches list mode's Enter key (copy and exit silently)
- [ ] Update `renderComplete()` in `internal/initialize/wizard.go` to include 'c' hotkey in help text
  - Only show on success screen (not error screen)
  - Format: "c: copy prompt | q: quit" or similar using `subtleStyle`
- [ ] Write unit tests for the new keyboard handler behavior in `internal/initialize/wizard_test.go`
  - Test 'c' key triggers clipboard operation and returns tea.Quit
  - Test error handling for clipboard failures (does not exit)
  - Test help text includes 'c' on success screen
  - Test help text excludes 'c' on error screen
  - Test copied text excludes surrounding quotes (raw content only)

## Validation

- [ ] Run `spectr validate add-copy-hotkey-init-next-steps --strict` and fix any issues
- [ ] Run existing tests: `go test ./internal/initialize/...`
- [ ] Manual testing:
  - Run `spectr init` in interactive mode
  - Complete initialization successfully
  - Press 'c' on Next Steps screen
  - Verify prompt is copied to clipboard
  - Paste in another application to confirm
  - Test in SSH session to verify OSC 52 fallback works
