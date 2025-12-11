## 1. Setup

- [ ] 1.1 Add blackfriday v2 dependency to go.mod
- [ ] 1.2 Create `internal/markdown/` package directory structure

## 2. Core Markdown Package

- [ ] 2.1 Implement `internal/markdown/types.go` with Header, Section, Task types
- [ ] 2.2 Implement `internal/markdown/parser.go` with Parse() function returning AST
- [ ] 2.3 Implement `internal/markdown/headers.go` with header extraction (H1-H4)
- [ ] 2.4 Implement `internal/markdown/sections.go` with section content extraction
- [ ] 2.5 Implement `internal/markdown/tasks.go` with task checkbox parsing
- [ ] 2.6 Write unit tests for all markdown package functions

## 3. Replace Parsers Package

- [ ] 3.1 Update `internal/parsers/parsers.go` to use markdown package for task counting
- [ ] 3.2 Update `internal/parsers/parsers.go` to use markdown package for delta counting
- [ ] 3.3 Update `internal/parsers/parsers.go` to use markdown package for requirement counting
- [ ] 3.4 Update `internal/parsers/requirement_parser.go` to use markdown package
- [ ] 3.5 Update `internal/parsers/delta_parser.go` to use markdown package
- [ ] 3.6 Remove unused regex imports from parsers package
- [ ] 3.7 Verify all existing parsers tests pass

## 4. Replace Validation Parser

- [ ] 4.1 Update `internal/validation/parser.go` ExtractSections() to use markdown package
- [ ] 4.2 Update `internal/validation/parser.go` ExtractRequirements() to use markdown package
- [ ] 4.3 Update `internal/validation/parser.go` ExtractScenarios() to use markdown package
- [ ] 4.4 Keep ContainsShallOrMust() and NormalizeRequirementName() as-is (not markdown parsing)
- [ ] 4.5 Remove unused regex imports from validation/parser.go
- [ ] 4.6 Verify all existing validation tests pass

## 5. Replace Archive Spec Merger

- [ ] 5.1 Update `internal/archive/spec_merger.go` splitSpec() to use markdown package
- [ ] 5.2 Update `internal/archive/spec_merger.go` extractOrderedRequirements() to use markdown package
- [ ] 5.3 Keep reconstructSpec() newline normalization regex (utility, not markdown)
- [ ] 5.4 Remove unused regex imports from spec_merger.go
- [ ] 5.5 Verify all existing archive tests pass

## 6. Replace Accept Command Parser

- [ ] 6.1 Update `cmd/accept.go` task parsing to use markdown package
- [ ] 6.2 Remove package-level regex variables from accept.go
- [ ] 6.3 Verify accept command tests pass

## 7. Validation and Cleanup

- [ ] 7.1 Run full test suite (`go test ./...`)
- [ ] 7.2 Run linter (`golangci-lint run`)
- [ ] 7.3 Test with example spec files in `examples/` directory
- [ ] 7.4 Test with actual project specs in `spectr/specs/`
- [ ] 7.5 Update any documentation referencing regex parsing
