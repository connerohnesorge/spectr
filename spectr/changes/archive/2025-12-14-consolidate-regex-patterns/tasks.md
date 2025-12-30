# Implementation Tasks

## 1. Setup Package Structure

- [ ] 1.1 Create `internal/regex/` package directory
- [ ] 1.2 Create `internal/regex/doc.go` with package documentation
- [ ] 1.3 Create `internal/regex/headers.go` with H2, H3, H4 patterns
- [ ] 1.4 Add header helper functions (MatchH2SectionHeader,
  MatchH2DeltaSection, MatchH3Requirement, MatchH4Scenario)
- [ ] 1.5 Create `internal/regex/tasks.go` with TaskCheckbox, NumberedTask,
  NumberedSection patterns
- [ ] 1.6 Add task helper functions (MatchTaskCheckbox, MatchNumberedTask,
  MatchNumberedSection)
- [ ] 1.7 Create `internal/regex/renames.go` with RenamedFrom, RenamedTo, and
  Alt variants
- [ ] 1.8 Add rename helper functions (MatchRenamedFrom, MatchRenamedFromAlt,
  MatchRenamedTo, MatchRenamedToAlt)
- [ ] 1.9 Create `internal/regex/sections.go` with FindSectionContent,
  FindDeltaSectionContent, FindRequirementsSection
- [ ] 1.10 Write unit tests for headers.go in `internal/regex/headers_test.go`
- [ ] 1.11 Write unit tests for tasks.go in `internal/regex/tasks_test.go`
- [ ] 1.12 Write unit tests for renames.go in `internal/regex/renames_test.go`
- [ ] 1.13 Write unit tests for sections.go in `internal/regex/sections_test.go`

## 2. Migrate Parsers Package

- [ ] 2.1 Update `internal/parsers/parsers.go` to use `regex.TaskCheckbox` and
  `regex.MatchTaskCheckbox`
- [ ] 2.2 Update `internal/parsers/parsers.go` to use `regex.H2DeltaSection` and
  `regex.MatchH2DeltaSection`
- [ ] 2.3 Update `internal/parsers/parsers.go` to use `regex.H3Requirement`
- [ ] 2.4 Update `internal/parsers/requirement_parser.go` to use regex package
  for requirement and scenario patterns
- [ ] 2.5 Update `internal/parsers/delta_parser.go` to use
  `regex.FindSectionContent` for section extraction
- [ ] 2.6 Update `internal/parsers/delta_parser.go` to use regex package for
  requirement parsing
- [ ] 2.7 Update `internal/parsers/delta_parser.go` to use regex package for
  RENAMED parsing
- [ ] 2.8 Remove unused local regex imports from parsers package
- [ ] 2.9 Verify all existing parsers tests pass

## 3. Migrate Validation Package

- [ ] 3.1 Update `internal/validation/parser.go` ExtractSections() to use
  `regex.MatchH2SectionHeader`
- [ ] 3.2 Update `internal/validation/parser.go` ExtractRequirements() to use
  `regex.MatchH3Requirement`
- [ ] 3.3 Update `internal/validation/parser.go` ExtractScenarios() to use
  `regex.MatchH4Scenario`
- [ ] 3.4 Keep ContainsShallOrMust() inline (validation-specific, not
  structural)
- [ ] 3.5 Keep NormalizeRequirementName() space regex inline (utility)
- [ ] 3.6 Update `internal/validation/change_rules.go` RENAMED parsing to use
  `regex.MatchRenamedFromAlt` and `regex.MatchRenamedToAlt`
- [ ] 3.7 Remove unused regex imports from validation package
- [ ] 3.8 Verify all existing validation tests pass

## 4. Migrate Archive Package

- [ ] 4.1 Update `internal/archive/spec_merger.go` splitSpec() to use
  `regex.FindRequirementsSection`
- [ ] 4.2 Update `internal/archive/spec_merger.go` extractOrderedRequirements()
  to use `regex.H3Requirement`
- [ ] 4.3 Keep newline normalization regex inline (utility, single use)
- [ ] 4.4 Remove unused regex imports from spec_merger.go
- [ ] 4.5 Verify all existing archive tests pass

## 5. Migrate Accept Command

- [ ] 5.1 Update `cmd/accept.go` to use `regex.NumberedSection` and
  `regex.MatchNumberedSection`
- [ ] 5.2 Update `cmd/accept.go` to use `regex.NumberedTask` and
  `regex.MatchNumberedTask`
- [ ] 5.3 Remove package-level regex variables from accept.go
- [ ] 5.4 Verify accept command tests pass

## 6. Validation and Cleanup

- [ ] 6.1 Run full test suite (`go test ./...`)
- [ ] 6.2 Run linter (`golangci-lint run`)
- [ ] 6.3 Test with example spec files in `examples/` directory
- [ ] 6.4 Test with actual project specs in `spectr/specs/`
- [ ] 6.5 Verify no duplicate regex patterns remain in migrated files
- [ ] 6.6 Update package doc comments if needed
- [ ] 6.7 Verify output is byte-for-byte identical to pre-migration behavior
