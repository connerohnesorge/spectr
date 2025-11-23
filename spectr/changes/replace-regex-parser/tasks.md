## 1. Build Core Lexer/Parser Infrastructure

- [ ] 1.1 Define `Token` types and `Lexer` structure <!-- id: 0 -->
- [ ] 1.2 Implement `Lexer` state functions (Text, Header, CodeBlock) <!-- id: 1 -->
- [ ] 1.3 Define AST nodes (`Document`, `Header`, `Block`, etc.) <!-- id: 2 -->
- [ ] 1.4 Implement `Parser` to build AST from tokens <!-- id: 3 -->
- [ ] 1.5 Implement `Extractor` for Requirements and Scenarios <!-- id: 4 -->
- [ ] 1.6 Implement `Extractor` for Delta specs (ADDED, MODIFIED, etc.) <!-- id: 5 -->
- [ ] 1.7 Add unit tests for Lexer, Parser, and Extractor <!-- id: 6 -->

## 2. Replace internal/parsers/ Package

- [ ] 2.1 Replace usage in `internal/parsers/requirement_parser.go` (3 regex patterns) <!-- id: 7 -->
- [ ] 2.2 Replace usage in `internal/parsers/delta_parser.go` (11 regex patterns) <!-- id: 8 -->
- [ ] 2.3 Replace utility functions in `internal/parsers/parsers.go` (CountTasks, CountDeltas, CountRequirements) <!-- id: 10 -->

## 3. Replace internal/validation/ Package (CRITICAL)

- [ ] 3.1 Replace `ExtractSections` in `internal/validation/parser.go` (## section headers) <!-- id: 11 -->
- [ ] 3.2 Replace `ExtractRequirements` in `internal/validation/parser.go` (### Requirement: headers) <!-- id: 12 -->
- [ ] 3.3 Replace `ExtractScenarios` in `internal/validation/parser.go` (#### Scenario: headers) <!-- id: 13 -->
- [ ] 3.4 Replace `ContainsShallOrMust` in `internal/validation/parser.go` (use lexer tokens) <!-- id: 14 -->
- [ ] 3.5 Replace `NormalizeRequirementName` in `internal/validation/parser.go` (whitespace normalization) <!-- id: 15 -->
- [ ] 3.6 Replace `parseRenamedRequirements` in `internal/validation/change_rules.go` (FROM/TO parsing) <!-- id: 16 -->

## 4. Replace internal/archive/ Package

- [ ] 4.1 Replace `reconstructSpec` in `internal/archive/spec_merger.go` (blank line normalization) <!-- id: 17 -->
- [ ] 4.2 Replace `splitSpec` in `internal/archive/spec_merger.go` (section splitting) <!-- id: 18 -->
- [ ] 4.3 Replace `extractOrderedRequirements` in `internal/archive/spec_merger.go` (requirement extraction) <!-- id: 19 -->

## 5. Testing and Validation

- [ ] 5.1 Verify all unit tests pass across all packages <!-- id: 20 -->
- [ ] 5.2 Run integration tests with edge cases (markdown in code blocks) <!-- id: 21 -->
- [ ] 5.3 Verify `spectr validate` works correctly with new parser <!-- id: 9 -->
- [ ] 5.4 Test archive workflow with complex delta specs <!-- id: 22 -->
