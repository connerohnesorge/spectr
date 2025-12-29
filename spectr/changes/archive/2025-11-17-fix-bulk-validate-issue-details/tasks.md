# Implementation Tasks

## 1. Update Bulk Validation Output

- [x] 1.1 Modify `printBulkHumanResults` in `cmd/validate_print.go` to display
  full issue details for failed items
- [x] 1.2 Iterate through `result.Report.Issues` and print each issue with
  level, path, and message
- [x] 1.3 Ensure output format matches single-item validation format for
  consistency
- [x] 1.4 Preserve the summary line showing pass/fail/total counts

## 2. Testing

- [x] 2.1 Test `spectr validate --all --strict` with items containing various
  issue types
- [x] 2.2 Verify issue details are displayed in human-readable format
- [x] 2.3 Verify JSON output remains unchanged (already working correctly)
- [x] 2.4 Test with mix of passing and failing items
- [x] 2.5 Verify summary line displays correctly
