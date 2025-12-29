# Implementation Tasks: Add Validate Command

## 1. Foundation - Internal Packages

### 1.1 Create validation types package

- [x] 1.1.1 Create `internal/validation/types.go`
- [x] 1.1.2 Define `ValidationLevel` enum (ERROR, WARNING, INFO)
- [x] 1.1.3 Define `ValidationIssue` struct (level, path, message)
- [x] 1.1.4 Define `ValidationReport` struct (valid, issues, summary)
- [x] 1.1.5 Define `ValidationSummary` struct (errors, warnings, info counts)
- [x] 1.1.6 Add JSON struct tags for all exported types
- [x] 1.1.7 Write unit tests for type marshaling/unmarshaling

### 1.2 Create markdown parser utilities

- [x] 1.2.1 Create `internal/validation/parser.go`
- [x] 1.2.2 Implement `ExtractSections(content string) map[string]string` for ## headers
- [x] 1.2.3 Implement `ExtractRequirements(content string) []Requirement` for ### headers
- [x] 1.2.4 Implement `ExtractScenarios(requirementBlock string) []string` for #### headers
- [x] 1.2.5 Implement `ContainsShallOrMust(text string) bool` regex checker
- [x] 1.2.6 Implement `NormalizeRequirementName(name string) string` for duplicate detection
- [x] 1.2.7 Write unit tests with various markdown formats

### 1.3 Create discovery package

- [x] 1.3.1 Create `internal/discovery/discovery.go`
- [x] 1.3.2 Implement `GetActiveChangeIDs(projectRoot string) ([]string, error)`
- [x] 1.3.3 Implement `GetSpecIDs(projectRoot string) ([]string, error)`
- [x] 1.3.4 Handle missing directories gracefully (return empty slice, not error)
- [x] 1.3.5 Filter out "archive" directory from changes
- [x] 1.3.6 Write unit tests with mock filesystem

## 2. Core Validation Logic

### 2.1 Create spec validator

- [x] 2.1.1 Create `internal/validation/spec_rules.go`
- [x] 2.1.2 Implement `ValidateSpecFile(path string, strictMode bool) (*ValidationReport, error)`
- [x] 2.1.3 Check for "## Purpose" section presence (ERROR if missing)
- [x] 2.1.4 Check for "## Requirements" section presence (ERROR if missing)
- [x] 2.1.5 Check Purpose section length (WARNING if < 50 chars)
- [x] 2.1.6 Extract all requirements and validate each one
- [x] 2.1.7 Check each requirement has SHALL or MUST (WARNING if missing)
- [x] 2.1.8 Check each requirement has at least one scenario (WARNING if missing)
- [x] 2.1.9 Check scenario format uses "#### Scenario:" (ERROR if wrong format)
- [x] 2.1.10 Aggregate all issues into ValidationReport
- [x] 2.1.11 Set valid=false if errors exist (or warnings in strict mode)
- [x] 2.1.12 Write comprehensive unit tests

### 2.2 Create change delta validator

- [x] 2.2.1 Create `internal/validation/change_rules.go`
- [x] 2.2.2 Implement `ValidateChangeDeltaSpecs(changeDir string, strictMode bool) (*ValidationReport, error)`
- [x] 2.2.3 Scan changeDir/specs/ for capability subdirectories
- [x] 2.2.4 For each spec.md, parse delta sections (## ADDED/MODIFIED/REMOVED/RENAMED Requirements)
- [x] 2.2.5 Track total delta count across all files
- [x] 2.2.6 ERROR if total deltas == 0
- [x] 2.2.7 ERROR if delta sections present but no requirement blocks found
- [x] 2.2.8 For ADDED requirements: check SHALL/MUST and scenarios (ERROR if missing)
- [x] 2.2.9 For MODIFIED requirements: check SHALL/MUST and scenarios (ERROR if missing)
- [x] 2.2.10 For REMOVED requirements: names only, no scenario check
- [x] 2.2.11 For RENAMED requirements: validate FROM/TO pairs
- [x] 2.2.12 Check for duplicate requirement names within each section
- [x] 2.2.13 Check for cross-section conflicts (same name in ADDED and MODIFIED)
- [x] 2.2.14 Aggregate issues and build ValidationReport
- [x] 2.2.15 Write comprehensive unit tests with fixture files

### 2.3 Create main validator orchestrator

- [x] 2.3.1 Create `internal/validation/validator.go`
- [x] 2.3.2 Define `Validator` struct with strictMode field
- [x] 2.3.3 Implement `NewValidator(strictMode bool) *Validator` constructor
- [x] 2.3.4 Implement `ValidateSpec(path string) (*ValidationReport, error)` wrapper
- [x] 2.3.5 Implement `ValidateChange(changeDir string) (*ValidationReport, error)` wrapper
- [x] 2.3.6 Add helper method `CreateReport(issues []ValidationIssue) *ValidationReport`
- [x] 2.3.7 Implement strict mode logic (warnings count as errors)
- [x] 2.3.8 Write integration tests

## 3. CLI Command Implementation

### 3.1 Create validate command file

- [x] 3.1.1 Create `cmd/validate.go`
- [x] 3.1.2 Define `ValidateCmd` struct with Kong tags
- [x] 3.1.3 Add `ItemName *string` field with `arg:"" optional:"" help:"..."`
- [x] 3.1.4 Add `Strict bool` field with `name:"strict" help:"..."`
- [x] 3.1.5 Add `JSON bool` field with `name:"json" help:"..."`
- [x] 3.1.6 Add `All bool` field with `name:"all" help:"..."`
- [x] 3.1.7 Add `Changes bool` field with `name:"changes" help:"..."`
- [x] 3.1.8 Add `Specs bool` field with `name:"specs" help:"..."`
- [x] 3.1.9 Add `Type *string` field with `name:"type" help:"..."`
- [x] 3.1.10 Add `NoInteractive bool` field with `name:"no-interactive" help:"..."`

### 3.2 Implement Run method - routing logic

- [x] 3.2.1 Implement `(c *ValidateCmd) Run() error`
- [x] 3.2.2 Check if bulk flags (--all, --changes, --specs) are set
- [x] 3.2.3 If bulk flags set, call `runBulkValidation()` and return
- [x] 3.2.4 If no ItemName and interactive mode, call `runInteractiveMode()` and return
- [x] 3.2.5 If no ItemName and non-interactive, print usage hint and exit 1
- [x] 3.2.6 If ItemName provided, call `runDirectValidation()` and return

### 3.3 Implement direct validation

- [x] 3.3.1 Create `runDirectValidation(itemName string) error` method
- [x] 3.3.2 Use discovery to get list of changes and specs
- [x] 3.3.3 Check if itemName exists in changes list
- [x] 3.3.4 Check if itemName exists in specs list
- [x] 3.3.5 If in both and no --type flag, error with disambiguation message
- [x] 3.3.6 If in neither, error with "not found" and nearest match suggestions
- [x] 3.3.7 Determine type (from --type flag or auto-detection)
- [x] 3.3.8 Call appropriate validator (spec or change)
- [x] 3.3.9 Print report (JSON or human-readable)
- [x] 3.3.10 Set exit code based on validation result

### 3.4 Implement bulk validation

- [x] 3.4.1 Create `runBulkValidation() error` method
- [x] 3.4.2 Determine scope from flags (all, changes-only, specs-only)
- [x] 3.4.3 Use discovery to get relevant item lists
- [x] 3.4.4 Create worker pool with concurrency=6 (or env var)
- [x] 3.4.5 Build task queue (one task per item)
- [x] 3.4.6 Execute tasks in parallel, collecting results
- [x] 3.4.7 Sort results by item ID
- [x] 3.4.8 Calculate summary statistics (totals, by-type)
- [x] 3.4.9 Print results (JSON or human-readable)
- [x] 3.4.10 Set exit code (1 if any failures)

### 3.5 Implement interactive mode

- [x] 3.5.1 Create `runInteractiveMode() error` method
- [x] 3.5.2 Check if stdin is TTY, error if not
- [x] 3.5.3 Use interactive prompt library (e.g., promptui or survey)
- [x] 3.5.4 Show options: "All", "All changes", "All specs", "Pick specific item"
- [x] 3.5.5 If "Pick specific item", show second prompt with all changes and specs
- [x] 3.5.6 Execute selected validation mode
- [x] 3.5.7 Handle user cancellation gracefully

### 3.6 Implement output formatting

- [x] 3.6.1 Create `printReport(itemType, itemName string, report *ValidationReport)` method
- [x] 3.6.2 For JSON mode, marshal report to JSON and print
- [x] 3.6.3 For human mode, print "✓ valid" or "✗ has issues"
- [x] 3.6.4 For each issue, print level icon, path, and message
- [x] 3.6.5 After issues, print "Next steps" guidance based on issue types
- [x] 3.6.6 Create `printBulkResults(results []BulkItemResult)` for bulk output
- [x] 3.6.7 Print per-item status with checkmark/x icons
- [x] 3.6.8 Print summary line with totals

## 4. CLI Integration

### 4.1 Wire validate command into CLI

- [x] 4.1.1 Open `cmd/root.go`
- [x] 4.1.2 Add `Validate ValidateCmd` field to `CLI` struct
- [x] 4.1.3 Add appropriate struct tag: `cmd:"" help:"Validate specs and changes"`
- [x] 4.1.4 Build and test `spectr validate --help`

## 5. Testing

### 5.1 Unit tests for validation rules

- [x] 5.1.1 Test spec validation with valid spec
- [x] 5.1.2 Test spec validation with missing Purpose
- [x] 5.1.3 Test spec validation with missing Requirements
- [x] 5.1.4 Test spec validation with requirement lacking scenarios
- [x] 5.1.5 Test spec validation with requirement lacking SHALL/MUST
- [x] 5.1.6 Test spec validation with incorrect scenario format

### 5.2 Unit tests for change validation

- [x] 5.2.1 Test change validation with valid deltas
- [x] 5.2.2 Test change validation with no deltas
- [x] 5.2.3 Test change validation with empty delta sections
- [x] 5.2.4 Test change validation with missing scenarios in ADDED
- [x] 5.2.5 Test change validation with duplicate requirements
- [x] 5.2.6 Test change validation with cross-section conflicts

### 5.3 Integration tests

- [x] 5.3.1 Create test fixtures in testdata/ directory
- [x] 5.3.2 Test end-to-end: validate valid spec file
- [x] 5.3.3 Test end-to-end: validate invalid spec file
- [x] 5.3.4 Test end-to-end: validate valid change
- [x] 5.3.5 Test end-to-end: validate invalid change
- [x] 5.3.6 Test bulk validation with mixed results
- [x] 5.3.7 Test JSON output parsing

### 5.4 CLI tests

- [x] 5.4.1 Test command registration in root CLI
- [x] 5.4.2 Test flag parsing (all combinations)
- [x] 5.4.3 Test exit codes for success and failure
- [x] 5.4.4 Test error messages for ambiguous items
- [x] 5.4.5 Test error messages for not found items

## 6. Documentation and Polish

### 6.1 Update documentation

- [x] 6.1.1 Update README with validate command usage
- [x] 6.1.2 Add examples to spectr/AGENTS.md if applicable
- [x] 6.1.3 Document validation rules in project documentation

### 6.2 Error message polish

- [x] 6.2.1 Review all error messages for clarity
- [x] 6.2.2 Ensure remediation guidance is helpful and accurate
- [x] 6.2.3 Add examples to error messages where useful

### 6.3 Performance testing

- [x] 6.3.1 Create large test project (50+ specs and changes)
- [x] 6.3.2 Benchmark bulk validation performance
- [x] 6.3.3 Tune concurrency if needed
- [x] 6.3.4 Profile for memory usage

## 7. Validation and Cleanup

### 7.1 Self-validation

- [x] 7.1.1 Use the new validate command on this change proposal
- [x] 7.1.2 Fix any issues found
- [x] 7.1.3 Ensure all specs in project validate successfully

### 7.2 Code review prep

- [x] 7.2.1 Run `go fmt` on all new code
- [x] 7.2.2 Run `go vet` and fix warnings
- [x] 7.2.3 Run `go test ./...` and ensure all tests pass
- [x] 7.2.4 Check test coverage (aim for >80%)

### 7.3 Final verification

- [x] 7.3.1 Test on fresh clone of repository
- [x] 7.3.2 Verify `go build` succeeds
- [x] 7.3.3 Verify `spectr validate --help` displays correctly
- [x] 7.3.4 Run smoke tests on real project
