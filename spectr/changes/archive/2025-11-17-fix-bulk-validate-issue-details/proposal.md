# Change: Fix Bulk Validation to Show Issue Details

## Why

When users run `spectr validate --all --strict`, they see that items failed
validation and the count of issues (e.g., "4 issue(s)"), but the actual issue
details are not displayed. This forces users to run individual validation
commands for each failed item to see what's wrong, which is inefficient and
frustrating.

The issue exists because `printBulkHumanResults()` in
`cmd/validate_print.go:63-98` only prints the issue count but doesn't iterate
through and display the actual issues like `printHumanReport()` does.

## What Changes

- Modify the bulk validation human-readable output to display full issue details
  (level, path, message) for each failed item
- Ensure the output format is consistent with single-item validation output
- Maintain the summary line showing pass/fail counts

## Impact

- **Affected specs**: cli-framework (Validate Command Structure requirement),
  cli-interface (validate command output format)
- **Affected code**: `cmd/validate_print.go` (printBulkHumanResults function)
- **User experience**: Users will see full validation details immediately
  instead of needing to re-run validation per item
