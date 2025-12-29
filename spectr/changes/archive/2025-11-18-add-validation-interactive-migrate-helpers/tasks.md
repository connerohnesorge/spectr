## 1. Preparation and Analysis

- [x] 1.1 Review existing helper files to understand all exported functions and types
- [x] 1.2 Review existing TUI patterns in internal/list/interactive.go for consistency
- [x] 1.3 Identify all import dependencies for migrated code
- [x] 1.4 Plan package structure for internal/validation/ additions

## 2. Migrate Helper Files to Internal Package

- [x] 2.1 Create internal/validation/helpers.go with functions from cmd/validate_helpers.go
  - [x] 2.1.1 Migrate determineItemType and itemTypeInfo
  - [x] 2.1.2 Migrate validateItemByType
  - [x] 2.1.3 Migrate validateSingleItem
  - [x] 2.1.4 Migrate contains helper
  - [x] 2.1.5 Export functions that cmd/ will need (uppercase names)
  - [x] 2.1.6 Keep internal-only functions unexported (lowercase names)
- [x] 2.2 Create internal/validation/items.go with functions from cmd/validate_items.go
  - [x] 2.2.1 Migrate createValidationItems
  - [x] 2.2.2 Migrate getAllItems, getChangeItems, getSpecItems
  - [x] 2.2.3 Migrate validationItem type
  - [x] 2.2.4 Export necessary types and functions
- [x] 2.3 Create internal/validation/formatters.go with functions from cmd/validate_print.go
  - [x] 2.3.1 Migrate printJSONReport, printHumanReport
  - [x] 2.3.2 Migrate printBulkJSONResults, printBulkHumanResults
  - [x] 2.3.3 Migrate bulkResult type
  - [x] 2.3.4 Convert from ValidateCmd methods to package functions
- [x] 2.4 Update cmd/validate.go to import and use new internal packages
  - [x] 2.4.1 Add imports for internal/validation helpers, items, formatters
  - [x] 2.4.2 Update function calls to use package-level functions
  - [x] 2.4.3 Remove method receivers where helpers are now package functions
- [x] 2.5 Delete migrated files from cmd/
  - [x] 2.5.1 Delete cmd/validate_helpers.go
  - [x] 2.5.2 Delete cmd/validate_items.go
  - [x] 2.5.3 Delete cmd/validate_print.go

## 3. Implement Interactive Validation TUI

- [x] 3.1 Create internal/validation/interactive.go following list/interactive.go patterns
  - [x] 3.1.1 Define interactiveValidationModel struct with bubbletea.Model interface
  - [x] 3.1.2 Implement Init() method
  - [x] 3.1.3 Implement Update(msg tea.Msg) method with key handlers
  - [x] 3.1.4 Implement View() method with styled output
  - [x] 3.1.5 Add menu options: All, All changes, All specs, Pick specific item
  - [x] 3.1.6 Handle selection and execute validation accordingly
- [x] 3.2 Add selection list for "Pick specific item" mode
  - [x] 3.2.1 Show combined list of changes and specs
  - [x] 3.2.2 Allow navigation with arrow keys / j/k
  - [x] 3.2.3 Execute validation on Enter
- [x] 3.3 Integrate interactive mode into cmd/validate.go
  - [x] 3.3.1 Detect TTY vs non-interactive environment
  - [x] 3.3.2 Call interactive TUI when no arguments and interactive
  - [x] 3.3.3 Show usage error when no arguments and non-interactive
  - [x] 3.3.4 Remove getInteractiveModeError() function
- [x] 3.4 Apply lipgloss styling consistent with other TUIs
  - [x] 3.4.1 Use project's existing style definitions
  - [x] 3.4.2 Add help text with key bindings
  - [x] 3.4.3 Show validation progress during execution

## 4. Testing

- [x] 4.1 Run existing validation tests to ensure no regressions
  - [x] 4.1.1 Run go test ./internal/validation/...
  - [x] 4.1.2 Verify all existing tests pass
- [x] 4.2 Create tests for migrated helpers
  - [x] 4.2.1 Add internal/validation/helpers_test.go
  - [x] 4.2.2 Add internal/validation/items_test.go
  - [x] 4.2.3 Add internal/validation/formatters_test.go
  - [x] 4.2.4 Ensure test coverage matches or exceeds original
- [x] 4.3 Create tests for interactive TUI
  - [x] 4.3.1 Add internal/validation/interactive_test.go
  - [x] 4.3.2 Test model initialization
  - [x] 4.3.3 Test key event handling
  - [x] 4.3.4 Test selection logic
- [x] 4.4 Manual testing of interactive mode
  - [x] 4.4.1 Test in TTY: spectr validate (no args)
  - [x] 4.4.2 Test selecting "All" option
  - [x] 4.4.3 Test selecting "All changes" option
  - [x] 4.4.4 Test selecting "All specs" option
  - [x] 4.4.5 Test selecting specific item
  - [x] 4.4.6 Test quit/cancel behavior
  - [x] 4.4.7 Test in non-TTY environment (CI simulation)

## 5. Documentation and Validation

- [x] 5.1 Run linter and fix any issues
  - [x] 5.1.1 Run golangci-lint run
  - [x] 5.1.2 Fix any linting errors
- [x] 5.2 Validate the change proposal
  - [x] 5.2.1 Run spectr validate add-validation-interactive-migrate-helpers --strict
  - [x] 5.2.2 Fix any validation issues
- [x] 5.3 Run full test suite
  - [x] 5.3.1 Run go test ./...
  - [x] 5.3.2 Ensure all tests pass
  - [x] 5.3.3 Check test coverage hasn't decreased
- [x] 5.4 Build and smoke test the binary
  - [x] 5.4.1 Run go build
  - [x] 5.4.2 Test validation command in various modes
  - [x] 5.4.3 Verify interactive mode works as expected
