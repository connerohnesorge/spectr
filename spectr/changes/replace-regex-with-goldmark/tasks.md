# Implementation Tasks

## 1. Foundation
- [ ] 1.1 Add goldmark dependency to go.mod (`github.com/yuin/goldmark`)
- [ ] 1.2 Create `internal/parsers/ast_utils.go` with basic AST walker utilities
- [ ] 1.3 Add helper functions: `walkAST()`, `findHeadings()`, `extractTextContent()`
- [ ] 1.4 Write unit tests for AST utilities

## 2. Simple Parsers (Titles and Counts)
- [ ] 2.1 Implement AST-based `ExtractTitleAST()` in parsers.go
- [ ] 2.2 Implement AST-based `CountRequirementsAST()` in parsers.go
- [ ] 2.3 Implement AST-based `CountTasksAST()` in parsers.go
- [ ] 2.4 Write equivalence tests (regex output == AST output for all test cases)
- [ ] 2.5 Add edge case tests (code blocks with headers, nested lists)

## 3. Requirement Parsing
- [ ] 3.1 Enhance `RequirementBlock` struct with `SourcePos` field (line, column)
- [ ] 3.2 Implement AST-based `ParseRequirementsAST()` in requirement_parser.go
- [ ] 3.3 Implement AST-based `ParseScenariosAST()` in requirement_parser.go
- [ ] 3.4 Write equivalence tests comparing regex vs AST output
- [ ] 3.5 Add tests for malformed markdown (missing closing headers, code blocks)

## 4. Delta Section Parsing
- [ ] 4.1 Enhance `DeltaPlan` struct with source location metadata
- [ ] 4.2 Implement `parseDeltaSectionAST()` for ADDED/MODIFIED/REMOVED sections
- [ ] 4.3 Create AST walker for extracting section content by heading
- [ ] 4.4 Write equivalence tests for delta parsing
- [ ] 4.5 Add tests for edge cases (empty sections, malformed deltas)

## 5. Rename Parsing (Complex)
- [ ] 5.1 Implement AST + post-processing for RENAMED section parsing
- [ ] 5.2 Extract bullet list items with `FROM:` and `TO:` patterns
- [ ] 5.3 Parse code-fenced requirement names from bullets
- [ ] 5.4 Write tests for well-formed and malformed rename pairs
- [ ] 5.5 Add validation for missing FROM or TO in pairs

## 6. Integration and Switchover
- [ ] 6.1 Replace `ExtractTitle()` calls with `ExtractTitleAST()` across codebase
- [ ] 6.2 Replace `ParseRequirements()` calls with `ParseRequirementsAST()`
- [ ] 6.3 Replace `ParseDeltaSpec()` calls with `ParseDeltaSpecAST()`
- [ ] 6.4 Update validation package to use AST-based error locations
- [ ] 6.5 Update archive package to leverage AST metadata
- [ ] 6.6 Remove old regex-based functions once all consumers migrated

## 7. Testing and Validation
- [ ] 7.1 Run full test suite and ensure all tests pass
- [ ] 7.2 Run `go run main.go validate --all --strict` on entire project
- [ ] 7.3 Performance benchmarks: compare regex vs AST parsing time
- [ ] 7.4 Manual testing: archive a change, validate deltas, view dashboard
- [ ] 7.5 Update test fixtures if needed for improved error messages

## 8. Documentation and Cleanup
- [ ] 8.1 Update godoc comments to reflect AST-based implementation
- [ ] 8.2 Remove unused regex patterns and helper functions
- [ ] 8.3 Add code examples in comments showing AST traversal
- [ ] 8.4 Run `golangci-lint` and fix any linting issues
- [ ] 8.5 Update CHANGELOG or commit messages with migration notes
