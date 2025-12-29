# Change: Add Next Steps Message to Init Command

## Why

After running `spectr init`, users are left without clear guidance on what to do next. The current output only shows created files but doesn't explain how to actually start using Spectr with their AI coding assistant. This creates a poor first-run experience and may lead to confusion about the Spectr workflow.

## What Changes

Add a helpful "Next steps" message at the end of the initialization process that provides:

- Clear, copy-paste ready prompts for AI assistants
- Three progressive steps: populate project context, create first change, learn workflow
- References to key files (spectr/project.md, spectr/AGENTS.md)
- Formatted output with visual separators for clarity

This change affects both interactive (TUI wizard) and non-interactive (CLI) initialization modes.

## Impact

- Affected specs: `cli-interface` (initialization output)
- Affected code:
  - `cmd/root.go` (non-interactive mode output at lines 126-153)
  - `internal/init/wizard.go` (interactive mode completion screen at lines 394-463)
- User benefit: Improved onboarding, clearer path to productivity
- No breaking changes
