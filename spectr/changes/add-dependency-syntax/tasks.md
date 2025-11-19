# Implementation Tasks: Add Dependency Syntax

## 1. Parser Implementation

### 1.1 Create dependency parser
- [ ] 1.1.1 Create `internal/parsers/dependency_parser.go`
- [ ] 1.1.2 Define `DependencyReference` struct (type, id, location)
- [ ] 1.1.3 Define `DependencyType` enum (CHANGE, SPEC)
- [ ] 1.1.4 Implement `ParseDependencies(content string) []DependencyReference`
- [ ] 1.1.5 Implement regex pattern for `@depends(change-id)` syntax
- [ ] 1.1.6 Implement regex pattern for `@requires(spec:capability-id)` syntax
- [ ] 1.1.7 Extract line numbers for better error reporting
- [ ] 1.1.8 Handle edge cases (malformed syntax, duplicates, whitespace)
- [ ] 1.1.9 Write unit tests for parser with various inputs

## 2. Validation Types Extension

### 2.1 Extend validation types
- [ ] 2.1.1 Open `internal/validation/types.go`
- [ ] 2.1.2 Add `DependencyReference` type (or import from parsers)
- [ ] 2.1.3 Add `DependencyValidationIssue` type
- [ ] 2.1.4 Add dependency-specific issue messages
- [ ] 2.1.5 Write unit tests for new types

## 3. Dependency Validation Rules

### 3.1 Create dependency validator
- [ ] 3.1.1 Create `internal/validation/dependency_rules.go`
- [ ] 3.1.2 Implement `ValidateDependencies(deps []DependencyReference, projectPath string, strictMode bool) []ValidationIssue`
- [ ] 3.1.3 For @depends() references, verify change exists in spectr/changes/
- [ ] 3.1.4 For @requires() references, verify spec exists in spectr/specs/
- [ ] 3.1.5 Generate WARNING level issues for missing dependencies in development mode
- [ ] 3.1.6 Generate ERROR level issues for missing dependencies when archiving (strict mode)
- [ ] 3.1.7 Check for self-references (change depending on itself)
- [ ] 3.1.8 Check for duplicate dependency declarations
- [ ] 3.1.9 Provide helpful error messages with remediation guidance
- [ ] 3.1.10 Write comprehensive unit tests

### 3.2 Circular dependency detection (future enhancement)
- [ ] 3.2.1 Document that circular dependency detection is not implemented in this change
- [ ] 3.2.2 Add TODO comment for future circular dependency validation
- [ ] 3.2.3 Document in design.md why circular detection is deferred

## 4. Integration with Validator

### 4.1 Integrate dependency validation
- [ ] 4.1.1 Open `internal/validation/validator.go`
- [ ] 4.1.2 Import dependency parser and validation functions
- [ ] 4.1.3 In `ValidateChange()`, parse dependencies from proposal.md
- [ ] 4.1.4 Call `ValidateDependencies()` and append issues to report
- [ ] 4.1.5 Ensure dependency validation respects strictMode flag
- [ ] 4.1.6 Write integration tests for end-to-end dependency validation

## 5. Documentation Updates

### 5.1 Update AGENTS.md
- [ ] 5.1.1 Open `spectr/AGENTS.md`
- [ ] 5.1.2 Add section on dependency syntax under "Creating Change Proposals"
- [ ] 5.1.3 Document `@depends(change-id)` syntax with examples
- [ ] 5.1.4 Document `@requires(spec:capability-id)` syntax with examples
- [ ] 5.1.5 Explain when to use dependencies (builds on another change)
- [ ] 5.1.6 Explain when to use requirements (needs existing capability)
- [ ] 5.1.7 Document validation behavior (warnings vs errors)
- [ ] 5.1.8 Add troubleshooting section for common dependency errors
- [ ] 5.1.9 Add examples to proposal.md template section

## 6. Testing

### 6.1 Parser tests
- [ ] 6.1.1 Test parsing valid @depends(change-id) references
- [ ] 6.1.2 Test parsing valid @requires(spec:capability-id) references
- [ ] 6.1.3 Test parsing multiple dependencies in same file
- [ ] 6.1.4 Test parsing dependencies with whitespace variations
- [ ] 6.1.5 Test parsing malformed syntax (missing parens, etc.)
- [ ] 6.1.6 Test parsing empty dependency references
- [ ] 6.1.7 Test line number extraction accuracy

### 6.2 Validation tests
- [ ] 6.2.1 Test validation with existing change dependency
- [ ] 6.2.2 Test validation with missing change dependency (WARNING)
- [ ] 6.2.3 Test validation with existing spec requirement
- [ ] 6.2.4 Test validation with missing spec requirement (WARNING)
- [ ] 6.2.5 Test validation in strict mode (errors instead of warnings)
- [ ] 6.2.6 Test validation with self-reference
- [ ] 6.2.7 Test validation with duplicate dependencies
- [ ] 6.2.8 Test validation with multiple mixed dependencies

### 6.3 Integration tests
- [ ] 6.3.1 Create test fixture with valid dependencies
- [ ] 6.3.2 Create test fixture with invalid dependencies
- [ ] 6.3.3 Test end-to-end validation with `spectr validate`
- [ ] 6.3.4 Test that validation report includes dependency issues
- [ ] 6.3.5 Test JSON output includes dependency validation results
- [ ] 6.3.6 Test exit codes for dependency validation failures

## 7. Validation and Cleanup

### 7.1 Self-validation
- [ ] 7.1.1 Run `spectr validate add-dependency-syntax --strict`
- [ ] 7.1.2 Fix any validation issues found
- [ ] 7.1.3 Ensure this change proposal itself uses dependency syntax if applicable

### 7.2 Code quality
- [ ] 7.2.1 Run `go fmt` on all new and modified files
- [ ] 7.2.2 Run `go vet` and fix any warnings
- [ ] 7.2.3 Run `go test ./...` and ensure all tests pass
- [ ] 7.2.4 Check test coverage for new code (aim for >80%)

### 7.3 Final verification
- [ ] 7.3.1 Test on fresh clone of repository
- [ ] 7.3.2 Verify `go build` succeeds
- [ ] 7.3.3 Run smoke tests with real proposals
- [ ] 7.3.4 Verify documentation is clear and complete
