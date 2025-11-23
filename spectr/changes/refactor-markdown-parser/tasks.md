# Implementation Tasks

## 1. Foundation Setup

- [ ] 1.1 Create `internal/parsers/markdown` package directory
- [ ] 1.2 Add package documentation (`doc.go`)
- [ ] 1.3 Define token types and Token struct (`token.go`)
- [ ] 1.4 Write token type stringer tests

## 2. Lexer Implementation

- [ ] 2.1 Implement Lexer struct with position tracking (`lexer.go`)
- [ ] 2.2 Implement state function type and basic state machine
- [ ] 2.3 Implement `lexNormal` state (heading detection, text accumulation)
- [ ] 2.4 Implement `lexHeading` state (count hashes, extract text)
- [ ] 2.5 Implement `lexCodeFence` state (triple backtick handling)
- [ ] 2.6 Implement `lexList` state (detect list items)
- [ ] 2.7 Implement helper methods (next, backup, peek, emit, accept)
- [ ] 2.8 Add position tracking (line, column) for all tokens

## 3. Lexer Testing

- [ ] 3.1 Write test for heading tokenization (all levels)
- [ ] 3.2 Write test for code fence tokenization
- [ ] 3.3 Write test for code fence with markdown inside (critical edge case)
- [ ] 3.4 Write test for list tokenization
- [ ] 3.5 Write test for mixed content (headings, text, code, lists)
- [ ] 3.6 Write test for position tracking accuracy
- [ ] 3.7 Write test for blank line detection
- [ ] 3.8 Write test for EOF handling
- [ ] 3.9 Write table-driven tests for edge cases

## 4. AST Definition

- [ ] 4.1 Define NodeType enum (`ast.go`)
- [ ] 4.2 Define Node interface with Type, Position, Children methods
- [ ] 4.3 Implement Document struct
- [ ] 4.4 Implement Heading struct with level, text, content, position
- [ ] 4.5 Implement CodeBlock struct with language, code, position
- [ ] 4.6 Implement List and ListItem structs
- [ ] 4.7 Implement Paragraph and Text structs
- [ ] 4.8 Add Position struct with line and column

## 5. Parser Implementation

- [ ] 5.1 Implement Parser struct with lexer, current, peek tokens (`parser.go`)
- [ ] 5.2 Implement token advancement methods (advance, expect, match)
- [ ] 5.3 Implement Parse() method returning Document
- [ ] 5.4 Implement parseNode() dispatcher
- [ ] 5.5 Implement parseHeading() with content collection
- [ ] 5.6 Implement parseCodeBlock() preserving content verbatim
- [ ] 5.7 Implement parseList() and parseListItem()
- [ ] 5.8 Implement parseParagraph() for text blocks
- [ ] 5.9 Implement heading hierarchy (collect content until next same/higher level)
- [ ] 5.10 Add error recovery and ParseError type

## 6. Parser Testing

- [ ] 6.1 Write test for parsing headings with hierarchy
- [ ] 6.2 Write test for parsing code blocks
- [ ] 6.3 Write test for code block NOT parsed as heading (critical)
- [ ] 6.4 Write test for nested lists
- [ ] 6.5 Write test for complete document structure
- [ ] 6.6 Write test for malformed input (error handling)
- [ ] 6.7 Write test for position preservation in AST
- [ ] 6.8 Write round-trip test (parse → reconstruct → compare)

## 7. Extractor Implementation

- [ ] 7.1 Implement RequirementExtractor struct (`extractor.go`)
- [ ] 7.2 Implement ExtractRequirements() walking AST for ### Requirement:
- [ ] 7.3 Implement scenario extraction from requirement content
- [ ] 7.4 Implement SectionExtractor for ## headers → content map
- [ ] 7.5 Implement DeltaExtractor for ADDED/MODIFIED/REMOVED/RENAMED
- [ ] 7.6 Implement markdown reconstruction from AST nodes
- [ ] 7.7 Add helper functions for AST traversal and filtering

## 8. Extractor Testing

- [ ] 8.1 Write test for requirement extraction from simple spec
- [ ] 8.2 Write test for requirement extraction with code blocks
- [ ] 8.3 Write test for scenario extraction
- [ ] 8.4 Write test for section extraction
- [ ] 8.5 Write test for delta extraction (all four types)
- [ ] 8.6 Write test for edge case: requirement with code block containing "Scenario:"
- [ ] 8.7 Verify extracted structures match current parser output

## 9. Integration with Existing Parsers Package

- [ ] 9.1 Update `internal/parsers/parsers.go` to use markdown parser
- [ ] 9.2 Replace ExtractTitle() implementation
- [ ] 9.3 Keep TaskStatus and CountTasks() unchanged (not markdown parsing)
- [ ] 9.4 Update CountDeltas() to use new parser
- [ ] 9.5 Update CountRequirements() to use new parser
- [ ] 9.6 Update `internal/parsers/requirement_parser.go` to use new parser
- [ ] 9.7 Replace ParseRequirements() implementation
- [ ] 9.8 Replace ParseScenarios() implementation
- [ ] 9.9 Update `internal/parsers/delta_parser.go` to use new parser
- [ ] 9.10 Replace ParseDeltaSpec() implementation
- [ ] 9.11 Ensure all public APIs remain unchanged

## 10. Validation Package Integration

- [ ] 10.1 Update `internal/validation/parser.go` to use new parser
- [ ] 10.2 Replace ExtractSections() implementation
- [ ] 10.3 Replace ExtractRequirements() implementation
- [ ] 10.4 Replace ExtractScenarios() implementation
- [ ] 10.5 Verify improved error messages with position info

## 11. Test Suite Validation

- [ ] 11.1 Run all existing parser tests - must pass
- [ ] 11.2 Run all validation tests - must pass
- [ ] 11.3 Run all archive tests - must pass
- [ ] 11.4 Add new tests for previously broken edge cases
- [ ] 11.5 Verify test coverage >90% for markdown package

## 12. Edge Case Testing

- [ ] 12.1 Test spec with code block containing "### Requirement:"
- [ ] 12.2 Test spec with code block containing "#### Scenario:"
- [ ] 12.3 Test spec with nested lists containing code blocks
- [ ] 12.4 Test spec with indented code blocks (4 spaces)
- [ ] 12.5 Test spec with mixed code fence languages
- [ ] 12.6 Test spec with escaped markdown characters
- [ ] 12.7 Test spec with HTML comments
- [ ] 12.8 Test spec with blank lines in various contexts
- [ ] 12.9 Create integration test using actual failing spec examples

## 13. Performance Benchmarking

- [ ] 13.1 Create benchmark for new lexer
- [ ] 13.2 Create benchmark for new parser
- [ ] 13.3 Create benchmark for old regex-based parser
- [ ] 13.4 Compare new vs old parser performance
- [ ] 13.5 Verify <10% performance regression
- [ ] 13.6 Profile memory allocations
- [ ] 13.7 Optimize hot paths if needed
- [ ] 13.8 Test with large generated specs (100, 500, 1000 requirements)

## 14. Integration Testing

- [ ] 14.1 Test parsing all existing specs in spectr/specs/
- [ ] 14.2 Test parsing all existing changes in spectr/changes/
- [ ] 14.3 Compare output with old parser (should match or improve)
- [ ] 14.4 Document any differences as improvements
- [ ] 14.5 Verify no regressions in validation behavior
- [ ] 14.6 Test with CI specs (if any)

## 15. Documentation

- [ ] 15.1 Write comprehensive package documentation (`doc.go`)
- [ ] 15.2 Add godoc comments for all exported types and functions
- [ ] 15.3 Add code examples in godoc
- [ ] 15.4 Document lexer state machine with diagram
- [ ] 15.5 Document AST structure with examples
- [ ] 15.6 Add README.md in markdown package explaining architecture
- [ ] 15.7 Update CHANGELOG.md with parser improvements

## 16. Cleanup and Finalization

- [ ] 16.1 Remove old regex-based parsing code (after validation)
- [ ] 16.2 Update internal documentation references
- [ ] 16.3 Run golangci-lint and fix any issues
- [ ] 16.4 Run go fmt on all new code
- [ ] 16.5 Review code for simplification opportunities
- [ ] 16.6 Verify no commented-out code remains
- [ ] 16.7 Final test suite run (all tests passing)
- [ ] 16.8 Mark all tasks as complete

## 17. Validation and Approval

- [ ] 17.1 Run `spectr validate refactor-markdown-parser --strict`
- [ ] 17.2 Fix any validation issues
- [ ] 17.3 Request proposal approval from maintainers
- [ ] 17.4 Address any feedback
- [ ] 17.5 Get final approval before implementation
