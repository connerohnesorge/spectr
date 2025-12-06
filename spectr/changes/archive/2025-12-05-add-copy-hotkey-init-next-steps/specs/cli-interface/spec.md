# Cli Interface Specification Delta

## ADDED Requirements

### Requirement: Copy Populate Context Prompt in Init Next Steps

The Next Steps completion screen in the interactive initialization wizard SHALL provide a hotkey to copy the "populate project context" prompt (step 1) to the system clipboard.

#### Scenario: Copy prompt with 'c' hotkey

- **WHEN** the user is on the Next Steps completion screen after successful initialization
- **AND** the user presses the 'c' key
- **THEN** the raw prompt text (without surrounding quotes) "Review spectr/project.md and help me fill in our project's tech stack, conventions, and description. Ask me questions to understand the codebase." is copied to the clipboard
- **AND** the wizard exits immediately and returns to the shell
- **AND** no success message is displayed (silent exit, consistent with list mode Enter behavior)

#### Scenario: Clipboard copy failure handling

- **WHEN** the user presses 'c' to copy the prompt
- **AND** the clipboard operation fails
- **THEN** an error message is displayed (e.g., "Failed to copy to clipboard: [error]")
- **AND** the wizard does NOT exit
- **AND** the user can retry the copy operation or press 'q' to quit manually

#### Scenario: Help text shows copy hotkey

- **WHEN** the Next Steps completion screen is displayed after successful initialization
- **THEN** the footer help text SHALL include the 'c' hotkey
- **AND** the help text format is: "c: copy step 1 | q: quit" or "c: copy prompt | q: quit"
- **AND** the hotkey is clearly described

#### Scenario: Copy hotkey only on success screen

- **WHEN** initialization fails and the error screen is displayed
- **THEN** the 'c' hotkey is NOT active
- **AND** the help text does NOT mention the copy hotkey
- **AND** only quit controls are available

#### Scenario: Clipboard uses OSC 52 fallback

- **WHEN** the user presses 'c' in an SSH/remote session without native clipboard access
- **THEN** the copy operation uses OSC 52 escape sequences as fallback
- **AND** the operation is considered successful (OSC 52 does not report errors)
- **AND** the success message is displayed
