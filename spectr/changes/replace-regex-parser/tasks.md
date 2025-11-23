## Phase 1: Build Generic Markdown Parser (internal/mdparser)

- [ ] 1.1 Define Token types (TokenHeader, TokenText, TokenCodeFence, TokenCodeContent, TokenListItem, TokenBlankLine, TokenEOF)
- [ ] 1.2 Define AST node interfaces and structs (Node, Document, Header, Paragraph, CodeBlock, List)
- [ ] 1.3 Implement Lexer struct and state machine
- [ ] 1.4 Implement lexer state functions (lexText, lexHeader, lexCodeBlock, lexList)
- [ ] 1.5 Implement Parser to build AST from token stream
- [ ] 1.6 Add unit tests for Lexer (tokenization correctness, state transitions)
- [ ] 1.7 Add unit tests for Parser (AST construction, error recovery)
- [ ] 1.8 Test edge cases: code blocks with markdown syntax, nested structures, malformed input

## Phase 2: Build Benchmark Infrastructure

- [ ] 2.1 Create test corpus in testdata/benchmarks/ (small.md, medium.md, large.md, pathological.md)
- [ ] 2.2 Implement benchmark for current regex implementation (BenchmarkRegexRequirementParser, BenchmarkRegexDeltaParser)
- [ ] 2.3 Implement benchmark for new lexer/parser (BenchmarkLexerRequirementParser, BenchmarkLexerDeltaParser)
- [ ] 2.4 Add correctness validation (ensure both parsers produce identical results)
- [ ] 2.5 Run benchmarks and document results (ns/op, bytes/op, allocs/op)

## Phase 3: Decision Gate - Performance Validation

- [ ] 3.1 Analyze benchmark results (compare speed, memory, allocations)
- [ ] 3.2 Verify correctness: new parser handles edge cases regex cannot
- [ ] 3.3 Decision: If performance regression >2x, profile and optimize (token pooling, lazy AST)
- [ ] 3.4 Document trade-offs and decision to proceed or iterate
- [ ] 3.5 Get approval to proceed with migration (if benchmarks acceptable)

## Phase 4: Build Spectr Extractors (internal/parsers)

- [ ] 4.1 Implement extractor for Requirements (ExtractRequirements function)
- [ ] 4.2 Implement extractor for Scenarios (traverse AST, validate hierarchy)
- [ ] 4.3 Implement extractor for Delta specs (ADDED, MODIFIED, REMOVED, RENAMED)
- [ ] 4.4 Implement extractor for RENAMED FROM/TO parsing
- [ ] 4.5 Add unit tests for all extractors
- [ ] 4.6 Test extractor edge cases (Requirements in code blocks should be ignored)

## Phase 5: Migrate internal/parsers/ Package

- [ ] 5.1 Update requirement_parser.go to use mdparser + extractors (remove 3 regex patterns)
- [ ] 5.2 Update delta_parser.go to use mdparser + extractors (remove 11 regex patterns)
- [ ] 5.3 Update parsers.go utility functions (CountTasks, CountDeltas, CountRequirements)
- [ ] 5.4 Run existing unit tests in internal/parsers/ and verify all pass
- [ ] 5.5 Remove old regex-based code from internal/parsers/

## Phase 6: Migrate internal/validation/ Package

- [ ] 6.1 Replace ExtractSections in parser.go with mdparser API
- [ ] 6.2 Replace ExtractRequirements in parser.go with extractor API
- [ ] 6.3 Replace ExtractScenarios in parser.go with extractor API
- [ ] 6.4 Replace ContainsShallOrMust with lexer token-based check
- [ ] 6.5 Replace NormalizeRequirementName (integrate into extractor or keep as utility)
- [ ] 6.6 Replace parseRenamedRequirements in change_rules.go with extractor API
- [ ] 6.7 Run existing unit tests in internal/validation/ and verify all pass
- [ ] 6.8 Remove old regex-based parsing code from internal/validation/

## Phase 7: Migrate internal/archive/ Package

- [ ] 7.1 Replace reconstructSpec blank line normalization with AST-based approach
- [ ] 7.2 Replace splitSpec section splitting with mdparser section extraction
- [ ] 7.3 Replace extractOrderedRequirements with extractor API
- [ ] 7.4 Run existing unit tests in internal/archive/ and verify all pass
- [ ] 7.5 Remove old regex-based code from internal/archive/spec_merger.go

## Phase 8: Integration Testing and Validation

- [ ] 8.1 Run full test suite across all packages (go test ./...)
- [ ] 8.2 Test spectr validate with real spec files (including edge cases)
- [ ] 8.3 Test spectr archive with complex delta specs
- [ ] 8.4 Test end-to-end workflows (create change, validate, archive)
- [ ] 8.5 Verify no regressions in existing functionality
- [ ] 8.6 Update documentation (if internal APIs changed)

## Phase 9: Cleanup and Documentation

- [ ] 9.1 Remove all unused regex pattern constants
- [ ] 9.2 Update package documentation for internal/mdparser
- [ ] 9.3 Update package documentation for internal/parsers extractors
- [ ] 9.4 Add examples and usage documentation
- [ ] 9.5 Final code review and cleanup
