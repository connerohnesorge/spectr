# Implementation Tasks

## 1. Foundation
- [ ] 1.1 Add goldmark dependency to go.mod
- [ ] 1.2 Create internal/parsers/ast_utils.go with utilities: ParseMarkdown(), FindHeading(), SegmentToLineNumber(), SegmentToLocation() (line/column/byte offset), ExtractTextContent()
- [ ] 1.3 Write tests for AST utility functions, covering line/column/byte offset calculations

## 2. Simple Parsers
- [x] 2.1 Migrate ExtractTitle() to use goldmark AST (find first H1 node)
- [x] 2.2 Migrate CountRequirements(), CountDeltas(), CountTasks() to AST traversal
- [x] 2.3 Update parsers_test.go to test goldmark implementations
- [ ] 2.4 Run parallel tests comparing regex vs goldmark output
- [ ] 2.5 Test title extraction with code blocks containing heading-like text
- [ ] 2.6 Test task counting with code blocks containing checkbox-like text

## 3. Requirements Parsing
- [ ] 3.1 Update RequirementBlock struct with ASTNode and Location (SourceLocation with LineNumber, Column, ByteOffset) fields
- [ ] 3.2 Migrate ParseRequirements() to walk AST for H3 "Requirement:" nodes
- [ ] 3.3 Migrate ParseScenarios() to walk AST for H4 "Scenario:" nodes
- [ ] 3.4 Update requirement_parser_test.go for goldmark
- [ ] 3.5 Test requirement extraction with edge cases (code blocks, nested lists)

## 4. Delta Section Parsing
- [ ] 4.1 Implement AST walker to find H2 delta section headers (ADDED, MODIFIED, etc.)
- [ ] 4.2 Parse ADDED and MODIFIED sections using AST requirement collection
- [ ] 4.3 Parse REMOVED section extracting requirement names
- [ ] 4.4 Parse RENAMED section with list post-processing for FROM:/TO: format
- [ ] 4.5 Update DeltaPlan struct with SourceLocation metadata (line/column/byte offset)
- [ ] 4.6 Update delta_parser_test.go for goldmark
- [ ] 4.7 Add tests for malformed markdown detection (unclosed code blocks, malformed headers)
- [ ] 4.8 Implement error reporting with line/column numbers from AST node SourceLocation

## 5. Integration and Testing
- [ ] 5.1 Update internal/archive/spec_merger.go to use new parser structures
- [ ] 5.2 Update internal/validation/change_rules.go to use goldmark parsers
- [ ] 5.3 Update internal/validation/delta_validators.go for new structures
- [ ] 5.4 Update validation error messages to include line/column context from SourceLocation
- [ ] 5.5 Update internal/list/lister.go and internal/view/dashboard.go
- [ ] 5.6 Add integration tests for malformed markdown error reporting with line numbers
- [ ] 5.7 Run full test suite (go test ./...)
- [ ] 5.8 Performance benchmark goldmark vs regex on real spec files
- [ ] 5.9 Remove old regex parsing code from parsers.go, requirement_parser.go, delta_parser.go
- [ ] 5.10 Update any documentation referencing old parsing approach
