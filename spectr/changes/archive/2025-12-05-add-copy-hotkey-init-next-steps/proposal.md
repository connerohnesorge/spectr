# Change: Add Copy Hotkey to Init Next Steps Screen

## Why

After running `spectr init`, users see a Next Steps message with three helpful prompts for their AI assistant. The first step provides a prompt to populate the project context:

> "Review spectr/project.md and help me fill in our project's tech stack, conventions, and description. Ask me questions to understand the codebase."

Currently, users must manually select and copy this text from the terminal. This creates friction in the onboarding experience, especially since this is one of the first interactions users have with Spectr. The interactive wizard already demonstrates hotkey functionality (e.g., Enter to copy IDs in list mode), so users would expect similar convenience here.

## What Changes

Add a 'c' hotkey to the Next Steps completion screen in the interactive init wizard that copies the "populate project context" prompt (step 1) to the clipboard and exits. This allows users to quickly paste the prompt into their AI coding assistant without manual text selection. The behavior mirrors the 'Enter' key in `spectr list` interactive mode (copy ID and exit).

The change affects:

- Interactive mode completion screen keyboard handlers
- Screen footer help text to indicate the 'c' hotkey availability
- Clipboard integration using the existing `internal/tui.CopyToClipboard` helper

## Impact

- Affected specs: `cli-interface` (initialization wizard)
- Affected code:
  - `internal/initialize/wizard.go` (handleCompleteKeys handler at lines 252-259, renderComplete at lines 466-487)
  - `internal/initialize/constants.go` (add keyCopy constant)
- User benefit: Faster onboarding, reduced friction when starting with Spectr
- No breaking changes: Purely additive feature
- Consistent with existing clipboard patterns in `spectr list` interactive mode
