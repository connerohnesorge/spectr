# Implementation Tasks

## 1. Implementation

- [x] 1.1 Add lipgloss styling constants to formatters.go for error/warning
  colors
- [x] 1.2 Create helper function to convert absolute paths to spectr-relative
  paths
- [x] 1.3 Update PrintBulkHumanResults to add blank lines between failed items
- [x] 1.4 Update PrintBulkHumanResults to group issues by file path
- [x] 1.5 Apply color styling to [ERROR] and [WARNING] labels
- [x] 1.6 Enhance summary line to show "X errors, Y warnings" breakdown
- [x] 1.7 Add type indicators (change/spec) to item names

## 2. Testing

- [x] 2.1 Add unit tests for relative path conversion
- [x] 2.2 Add unit tests for colored output formatting
- [x] 2.3 Manual testing with spectr validate --all on real project
- [x] 2.4 Verify output renders correctly in terminal without TTY

## 3. Validation

- [x] 3.1 Run spectr validate improve-validate-output --strict
- [x] 3.2 Run go test ./internal/validation/...
- [x] 3.3 Run golangci-lint
