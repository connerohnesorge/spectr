# Lexer/Parser Capability

## ADDED Requirements

### Requirement: Parse Markdown Structure
The system SHALL parse markdown documents into a structured Abstract Syntax Tree (AST) that accurately represents headers, paragraphs, code blocks, and lists.

#### Scenario: Parsing basic structure
WHEN the parser processes a document with headers and paragraphs
THEN it produces an AST with corresponding Header and Paragraph nodes
AND the hierarchy is preserved

#### Scenario: Parsing code blocks
WHEN the parser processes a document containing code blocks
THEN the content inside the code blocks is treated as literal text
AND markdown syntax inside code blocks is NOT parsed as structure

### Requirement: Extract Spectr Elements
The system SHALL extract Spectr-specific elements (Requirements, Scenarios, Delta sections) by traversing the AST.

#### Scenario: Extracting valid requirement
WHEN the extractor traverses an AST with a "### Requirement: Foo" header
THEN it identifies a Requirement named "Foo"
AND captures the following content as the requirement body

#### Scenario: Ignoring syntax in code blocks
WHEN the extractor encounters "### Requirement: Foo" inside a code block
THEN it ignores it
AND does NOT create a Requirement entity

#### Scenario: Extracting delta sections
WHEN the extractor traverses an AST with "## ADDED Requirements", "## MODIFIED Requirements", "## REMOVED Requirements", or "## RENAMED Requirements" headers
THEN it identifies the appropriate delta operation type
AND extracts requirements within each delta section

#### Scenario: Extracting RENAMED requirements
WHEN the extractor encounters a RENAMED section with "- FROM: ### Requirement: OldName" and "- TO: ### Requirement: NewName"
THEN it captures both the old and new requirement names
AND associates them as a rename operation

### Requirement: Report Parsing Errors
The system SHALL report parsing errors with precise line and column information.

#### Scenario: Reporting malformed input
WHEN the parser encounters invalid syntax
THEN it returns an error
AND the error includes the line and column number of the failure

### Requirement: Generic Markdown Parser
The system SHALL provide a generic markdown parsing API in `internal/mdparser` that has zero knowledge of Spectr-specific semantics.

#### Scenario: Parser is domain-agnostic
WHEN the mdparser package is used to parse a markdown document
THEN it produces generic AST nodes (Document, Header, Paragraph, CodeBlock, List)
AND it does NOT recognize or specially handle "Requirement", "Scenario", or "Delta" keywords
AND it has no dependencies on Spectr domain logic

#### Scenario: Public API for external use
WHEN external code needs to parse markdown
THEN it can import internal/mdparser
AND use NewLexer(), Parse(), and AST traversal functions
AND get a complete AST representation of any markdown document

### Requirement: Two-Layer Architecture
The system SHALL separate generic markdown parsing (Layer 1) from Spectr-specific extraction (Layer 2).

#### Scenario: Layer 1 is reusable
WHEN Layer 1 (internal/mdparser) is implemented
THEN it can parse any markdown document without Spectr knowledge
AND it could be extracted as a standalone library
AND it has zero imports from internal/parsers or other Spectr packages

#### Scenario: Layer 2 consumes Layer 1
WHEN Layer 2 (internal/parsers extractors) processes a spec file
THEN it calls mdparser.Parse() to get an AST
AND it traverses the AST to find Spectr-specific elements
AND it applies Spectr business rules (requirement hierarchy, delta sections)

### Requirement: Performance Benchmarking
The system SHALL include a benchmark suite that compares the new lexer/parser performance against the current regex implementation.

#### Scenario: Benchmark test corpus
WHEN the benchmark suite runs
THEN it uses test files of varying sizes (small, medium, large, pathological)
AND each file represents realistic Spectr spec content
AND files include edge cases (code blocks with markdown, nested structures)

#### Scenario: Performance metrics
WHEN benchmarks execute
THEN they measure speed (nanoseconds per operation)
AND they measure memory usage (bytes allocated per operation)
AND they measure allocation count (heap allocations per operation)
AND they verify correctness (both parsers produce equivalent results)

#### Scenario: Performance acceptance criteria
WHEN benchmark results are analyzed
THEN correctness is prioritized over raw speed
AND performance regression of <2x is acceptable
AND if regression >2x, profiling and optimization is required before migration
AND documented trade-offs justify any performance difference

### Requirement: Conditional Replacement
The system SHALL only replace the regex implementation if benchmark validation confirms acceptable performance.

#### Scenario: Migration gate
WHEN benchmark results are reviewed
THEN migration proceeds only if all tests pass
AND regression is within acceptable limits (<2x or optimized)
AND correctness validation shows edge case improvements
AND decision is documented with benchmark data

#### Scenario: Rollback plan
WHEN migration causes unexpected issues
THEN the old regex implementation is available for rollback
AND a feature flag or build tag can switch implementations
AND both implementations are tested until confidence is high
