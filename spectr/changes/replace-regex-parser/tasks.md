## Phase 1: Build Generic Markdown Parser (internal/mdparser)

- [x] 1.1 Define Token types (TokenHeader, TokenText, TokenCodeFence, TokenCodeContent, TokenListItem, TokenBlankLine, TokenEOF)
- [x] 1.2 Define AST node interfaces and structs (Node, Document, Header, Paragraph, CodeBlock, List)
- [x] 1.3 Implement Lexer struct and state machine
- [x] 1.4 Implement lexer state functions (lexText, lexHeader, lexCodeBlock, lexList)
- [x] 1.5 Implement Parser to build AST from token stream
- [x] 1.6 Add unit tests for Lexer (tokenization correctness, state transitions)
- [x] 1.7 Add unit tests for Parser (AST construction, error recovery)
- [x] 1.8 Test edge cases: code blocks with markdown syntax, nested structures, malformed input

## Phase 2: Build Benchmark Infrastructure

- [x] 2.1 Create test corpus in testdata/benchmarks/ (small.md, medium.md, large.md, pathological.md)
- [x] 2.2 Implement benchmark for current regex implementation (BenchmarkRegexRequirementParser, BenchmarkRegexDeltaParser)
- [x] 2.3 Implement benchmark for new lexer/parser (BenchmarkLexerRequirementParser, BenchmarkLexerDeltaParser)
- [x] 2.4 Add correctness validation (ensure both parsers produce identical results)
- [x] 2.5 Run benchmarks and document results (ns/op, bytes/op, allocs/op)

## Phase 3: Decision Gate - Performance Validation

- [x] 3.1 Analyze benchmark results (compare speed, memory, allocations)
- [x] 3.2 Verify correctness: new parser handles edge cases regex cannot
- [x] 3.3 Decision: If performance regression >2x, profile and optimize (token pooling, lazy AST)
- [x] 3.4 Document trade-offs and decision to proceed or iterate
- [x] 3.5 Get approval to proceed with migration (if benchmarks acceptable)

## Phase 4: Build Spectr Extractors (internal/parsers)

- [x] 4.1 Implement extractor for Requirements (ExtractRequirements function)
- [x] 4.2 Implement extractor for Scenarios (traverse AST, validate hierarchy)
- [x] 4.3 Implement extractor for Delta specs (ADDED, MODIFIED, REMOVED, RENAMED)
- [x] 4.4 Implement extractor for RENAMED FROM/TO parsing
- [x] 4.5 Add unit tests for all extractors
- [x] 4.6 Test extractor edge cases (Requirements in code blocks should be ignored)

## Phase 5: Migrate internal/parsers/ Package

- [x] 5.1 Update requirement_parser.go to use mdparser + extractors (remove 3 regex patterns)
- [x] 5.2 Update delta_parser.go to use mdparser + extractors (remove 11 regex patterns)
- [x] 5.3 Update parsers.go utility functions (CountTasks, CountDeltas, CountRequirements)
- [x] 5.4 Run existing unit tests in internal/parsers/ and verify all pass
- [x] 5.5 Remove old regex-based code from internal/parsers/

## Phase 6: Migrate internal/validation/ Package

- [x] 6.1 Replace ExtractSections in parser.go with mdparser API
- [x] 6.2 Replace ExtractRequirements in parser.go with extractor API
- [x] 6.3 Replace ExtractScenarios in parser.go with extractor API
- [x] 6.4 Replace ContainsShallOrMust with lexer token-based check
- [x] 6.5 Replace NormalizeRequirementName (integrate into extractor or keep as utility)
- [x] 6.6 Replace parseRenamedRequirements in change_rules.go with extractor API
- [x] 6.7 Run existing unit tests in internal/validation/ and verify all pass
- [x] 6.8 Remove old regex-based parsing code from internal/validation/

## Phase 7: Migrate internal/archive/ Package

- [x] 7.1 Replace reconstructSpec blank line normalization with AST-based approach
- [x] 7.2 Replace splitSpec section splitting with mdparser section extraction
- [x] 7.3 Replace extractOrderedRequirements with extractor API
- [x] 7.4 Run existing unit tests in internal/archive/ and verify all pass
- [x] 7.5 Remove old regex-based code from internal/archive/spec_merger.go

## Phase 8: Integration Testing and Validation

- [x] 8.1 Run full test suite across all packages (go test ./...)
- [x] 8.2 Test spectr validate with real spec files (including edge cases)
- [x] 8.3 Test spectr archive with complex delta specs
- [x] 8.4 Test end-to-end workflows (create change, validate, archive)
- [x] 8.5 Verify no regressions in existing functionality
- [x] 8.6 Update documentation (if internal APIs changed)

## Phase 9: Cleanup and Documentation

- [x] 9.1 Remove all unused regex pattern constants
- [x] 9.2 Update package documentation for internal/mdparser
- [x] 9.3 Update package documentation for internal/parsers extractors
- [x] 9.4 Add examples and usage documentation
- [x] 9.5 Final code review and cleanup
