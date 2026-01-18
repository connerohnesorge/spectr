# Tasks: Fix JSONC File Validation

## 1. Foundation

- [ ] **1.1** Research and select Go JSONC parser library (check go.mod for existing options, research `tidwall/gjson` or similar)
- [ ] **1.2** Add JSONC parser dependency to go.mod if needed (via `go get`)
- [ ] **1.3** Create `cmd/accept_validator.go` with package documentation and imports

## 2. Core Validation Functions

- [ ] **2.1** Implement `validateJSONCOutput(jsonData []byte) error` function that:
  - Strips JSONC comments using existing `parsers.StripJSONComments`
  - Attempts to unmarshal into `map[string]interface{}`
  - Returns detailed error if parsing fails
- [ ] **2.2** Implement `validateTasksJSONC(tasksFile TasksFile, jsonData []byte) error` for round-trip validation:
  - Parse JSONC back into TasksFile struct
  - Compare all fields deeply (ID, Section, Description, Status, Children for each task)
  - Return error with details if any field differs
- [ ] **2.3** Implement `validateWithExternalParser(jsonData []byte) error` if external library selected:
  - Use external JSONC parser to validate format
  - Return detailed error on failure
- [ ] **2.4** Add validation calls to `writeTasksJSONC()` in `cmd/accept_writer.go`:
  - Call `validateJSONCOutput()` before writing file
  - Call `validateTasksJSONC()` for round-trip check
  - Return error if validation fails (don't write file)
- [ ] **2.5** Add validation calls to `writeRootTasksJSONC()` in `cmd/accept_writer.go`
- [ ] **2.6** Add validation calls to `writeChildTasksJSONC()` in `cmd/accept_writer.go`

## 3. Property-Based Testing

- [ ] **3.1** Create `cmd/accept_validator_test.go` with test infrastructure
- [ ] **3.2** Implement `TestJSONCValidation_SpecialCharacters` with test cases for:
  - Backslash `\`
  - Quote `"`
  - Newline `\n`
  - Tab `\t`
  - Carriage return `\r`
  - Backspace `\b`
  - Form feed `\f`
- [ ] **3.3** Implement `TestJSONCValidation_Unicode` with test cases for:
  - Emoji (🚀, 💻, ✅)
  - Non-ASCII characters (你好, مرحبا, Здравствуй)
  - Mixed unicode and ASCII
- [ ] **3.4** Implement `TestJSONCValidation_JSONMetaCharacters` with test cases for:
  - Curly braces in descriptions: `{`, `}`
  - Square brackets: `[`, `]`
  - Colons and commas: `:`, `,`
- [ ] **3.5** Implement `TestJSONCValidation_CommentLikeStrings` with test cases for:
  - Single-line comment syntax: `//`
  - Multi-line comment syntax: `/* */`
  - Mixed: `Update // comment-like text`
- [ ] **3.6** Implement `TestJSONCValidation_EdgeCases` with test cases for:
  - Very long descriptions (>1000 chars)
  - Empty descriptions
  - Whitespace-only descriptions
  - Descriptions ending with backslash

## 4. Round-Trip Testing

- [ ] **4.1** Implement `TestRoundTripConversion_AllFields` that:
  - Creates TasksFile with all field types populated
  - Marshals to JSONC
  - Validates JSONC
  - Unmarshals back to TasksFile
  - Compares original and result deeply
- [ ] **4.2** Implement `TestRoundTripConversion_Version2Hierarchical` for hierarchical format:
  - Test root file with Children references
  - Test child file with Parent field
  - Verify Includes array preserved
- [ ] **4.3** Implement `TestRoundTripConversion_RealWorldData` using archived tasks.jsonc files:
  - Read existing tasks.jsonc from `spectr/changes/archive/*/tasks.jsonc`
  - Validate they pass round-trip test
  - Report any files that fail

## 5. Integration with Existing Code

- [ ] **5.1** Update `writeAndCleanup()` in `cmd/accept.go` to handle validation errors gracefully
- [ ] **5.2** Add helpful error messages when validation fails:
  - Show task ID and description excerpt
  - Suggest common fixes (escape backslashes, etc.)
  - Include character position if available
- [ ] **5.3** Ensure dry-run mode (`--dry-run`) also validates JSONC before printing

## 6. Comprehensive Testing

- [ ] **6.1** Run all new tests: `go test -v ./cmd -run Validation`
- [ ] **6.2** Run full test suite: `go test ./...`
- [ ] **6.3** Test with real tasks.md files:
  - Create test tasks.md with special characters
  - Run `spectr accept` and verify JSONC is valid
  - Edit tasks.jsonc manually with invalid escape sequences and verify error detection
- [ ] **6.4** Run linting: `nix develop -c lint` or `golangci-lint run`

## 7. Documentation and Validation

- [ ] **7.1** Add godoc comments to all new validation functions
- [ ] **7.2** Update CLAUDE.md if new testing patterns introduced
- [ ] **7.3** Run `spectr validate fix-jsonc-validation --strict` to validate this change
- [ ] **7.4** Verify all scenarios in spec delta are covered by tests (cross-reference)

## 8. Final Verification

- [ ] **8.1** Test accept command with archived changes: `spectr accept <archived-change-id>` (should work with existing files)
- [ ] **8.2** Create a tasks.md with intentionally problematic descriptions and verify errors are caught
- [ ] **8.3** Verify performance impact is minimal (validation should be fast)
- [ ] **8.4** Run `nix build` to ensure build succeeds
