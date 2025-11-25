# Change: Polish List Interactive TUI

## Why

The current interactive TUI help text is verbose and cramped on a single line, making it harder to read. Additionally, users frequently need to archive changes after reviewing them in the list, requiring them to exit and run a separate command.

## What Changes

- Condense navigation hint from `↑/↓ or j/k: navigate` to `↑/↓/j/k: navigate`
- Shorten `e: edit proposal` to `e: edit` (context is already clear from mode)
- Split help text into two lines: controls on line 1, project path on line 2
- Add `a` hotkey to archive the selected change directly from interactive mode

## Impact

- Affected specs: cli-interface
- Affected code: internal/list/interactive.go
