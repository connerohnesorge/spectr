# Change: Improve validate --all output formatting

## Why

The current `spectr validate --all` output is difficult to scan when there are
multiple failures. Error messages run together without visual separation, full
absolute paths are repeated for each issue, and there's no clear distinction
between passed and failed items.

## What Changes

- Add blank line spacing between failed items for visual separation
- Use relative paths from spectr/ root instead of absolute paths
- Group errors by file when multiple issues exist in the same file
- Add color coding for error/warning levels using lipgloss (red for errors,
  yellow for warnings)
- Enhance summary to show error and warning counts separately
- Add item type icons (change/spec) for faster visual scanning

## Impact

- Affected specs: `validation`
- Affected code: `internal/validation/formatters.go`
