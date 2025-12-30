# Implementation Tasks

## 1. Create Helper Function

- [x] 1.1 Create a `formatNextStepsMessage()` helper function in
  `internal/init/executor.go`
- [x] 1.2 Function should return a formatted string with the next steps message
- [x] 1.3 Include visual separators (dashes) for readability
- [x] 1.4 Include three numbered steps with copy-paste ready prompts

## 2. Update Non-Interactive Mode

- [x] 2.1 Modify `runNonInteractiveInit()` in `cmd/root.go` (around line 126)
- [x] 2.2 Call the helper function after displaying created/updated files
- [x] 2.3 Only display next steps if there are no errors
- [x] 2.4 Ensure consistent spacing and formatting

## 3. Update Interactive Mode

- [x] 3.1 Modify `renderComplete()` in `internal/init/wizard.go` (around line
  394)
- [x] 3.2 Add next steps message to the success output
- [x] 3.3 Integrate with existing lipgloss styles for visual consistency
- [x] 3.4 Position after created/updated files list, before the quit instruction

## 4. Testing

- [x] 4.1 Test interactive mode initialization with tool selection
- [x] 4.2 Test interactive mode initialization with no tools
- [x] 4.3 Test non-interactive mode with `--non-interactive` flag
- [x] 4.4 Test non-interactive mode with `--tools` flag
- [x] 4.5 Verify message does not appear when errors occur
- [x] 4.6 Verify message formatting and readability in terminal

## 5. Documentation

- [x] 5.1 Update README.md if needed to reflect new onboarding flow
- [x] 5.2 Ensure the message text matches Spectr conventions
