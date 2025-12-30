# Change: Add --stdout flag to spectr list -I

## Why

When using `spectr list -I` in scripts or automated workflows, users need a way
to output the selected ID to stdout instead of copying it to the clipboard. This
enables piping the selection to other commands and integration with shell
scripts without requiring clipboard access.

## What Changes

- Add `--stdout` flag to the `spectr list` command
- When `--stdout` is combined with `-I` (interactive mode), pressing Enter
  prints the selected ID to stdout instead of copying to clipboard
- The flag provides a clean, scriptable output (just the ID, no "Copied:" prefix
  or other formatting)
- Mutually exclusive with non-interactive modes (only works with `-I`)

## Impact

- Affected specs: `cli-interface`
- Affected code:
  - `cmd/list.go` - Add Stdout flag to ListCmd struct
  - `internal/list/interactive.go` - Add stdout output mode to interactiveModel
  - `cmd/list_test.go` - Add tests for --stdout flag
  - `internal/list/interactive_test.go` - Add tests for stdout mode behavior
