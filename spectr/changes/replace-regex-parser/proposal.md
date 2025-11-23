# Replace Regex Parser with Lexer/Parser Architecture

## Goal
Replace the current regex-based markdown parsing implementation with a robust lexer/parser architecture inspired by the Go compiler and Rob Pike's "Lexical Scanning in Go" talk. This change aims to eliminate brittleness, improve error reporting, and correctly handle edge cases like nested structures and code blocks.

## Context
The current implementation relies heavily on `regexp` and line-by-line scanning to parse Spectr specifications (requirements, scenarios, deltas). This approach has several limitations:
- **Brittleness**: Code blocks containing markdown syntax (e.g., `### Requirement:`) are often incorrectly parsed as actual requirements.
- **No Context**: The parser lacks awareness of state (e.g., "inside a code block"), making it difficult to handle nested structures.
- **Poor Error Reporting**: Errors are often generic or missing precise line/column information.
- **Maintenance**: Complex regex patterns are hard to read and maintain.

## Solution
We will implement a hand-written lexer and parser modeled after the current Go compiler (`cmd/compile/internal/syntax`), using a **two-layer architecture** that separates generic markdown parsing from Spectr-specific domain logic:

### Layer 1: Generic Markdown Parser (`internal/mdparser`)
A reusable, general-purpose markdown parser that has no knowledge of Spectr requirements, scenarios, or deltas:
1.  **Lexer**: A state-machine-based lexer that tokenizes the input stream. It will handle state transitions (e.g., entering/exiting code blocks) and emit tokens, similar to the Go compiler's scanner.
2.  **Parser**: A recursive descent parser that consumes tokens and builds an Abstract Syntax Tree (AST) representing the document structure (headings, paragraphs, lists, code blocks, etc.). This mirrors the Go compiler's parser design for handling structure and error recovery.

### Layer 2: Spectr-Specific Extraction (`internal/parsers`)
Domain-specific logic that walks the generic AST to extract Spectr structures:
3.  **Extractor**: A semantic analysis layer that walks the generic AST to identify and extract Spectr-specific structures (Requirements, Scenarios, Deltas) and validates them according to Spectr rules.

This separation ensures:
- **Reusability**: The `internal/mdparser` package is standalone and could be extracted for use in other projects
- **Maintainability**: Generic markdown parsing logic is isolated from Spectr-specific business rules
- **Testability**: Each layer can be tested independently with appropriate fixtures

## Benchmarking
Before replacing the existing regex implementation, we will validate that the new parser meets performance requirements:

### Test Corpus
Benchmarks will run against:
- **Small files**: 1-5 requirements, <100 lines (typical change deltas)
- **Medium files**: 10-20 requirements, 200-500 lines (typical capability specs)
- **Large files**: 50+ requirements, 1000+ lines (stress testing)
- **Edge cases**: Deep nesting, code blocks with markdown syntax, escaped characters

### Metrics
For each test case, we will measure:
- **Speed**: nanoseconds per operation (ns/op)
- **Memory**: bytes allocated per operation (bytes/op)
- **Allocations**: number of allocations per operation (allocs/op)
- **Correctness**: AST output matches expected structure (100% required)

### Acceptance Criteria
- **Correctness**: MUST achieve 100% correctness on all test cases (non-negotiable)
- **Performance**: SHOULD not regress more than 2x compared to regex implementation
  - If performance regression exceeds 2x, we will profile and optimize before replacement
  - Trade-off: correctness and maintainability take priority over raw speed
  - Rationale: parsing is not in hot path; file sizes are typically <1000 lines

### Validation Gate
The regex implementation will only be replaced if:
1. All correctness tests pass (100% requirement)
2. Performance is within acceptable bounds (tracked, but not blocking if benefits outweigh costs)
3. Benchmark results are documented in the change's design.md or tasks.md

## Affected Files
This change will introduce a new package and replace regex usage in 6 existing files across 3 packages:

### NEW: internal/mdparser/ (generic markdown parser)
A new standalone package providing generic markdown parsing capabilities:
- **lexer.go**: State-machine tokenizer for markdown streams
- **parser.go**: Recursive descent parser producing generic AST
- **ast.go**: AST node types (Heading, Paragraph, List, CodeBlock, etc.)
- **lexer_test.go**, **parser_test.go**: Comprehensive test coverage
- **benchmark_test.go**: Performance benchmarks vs regex implementation

**Key characteristic**: This package has NO dependencies on Spectr-specific types or business logic. It's a reusable markdown parser.

### UPDATED: internal/parsers/ (3 files - Spectr-specific extraction)
These files will be updated to use the generic `internal/mdparser` AST instead of regex:
- **requirement_parser.go** (3 regex patterns → AST walker): Parses requirement blocks and scenarios from spec files
- **delta_parser.go** (11 regex patterns → AST walker): Parses delta specs (ADDED, MODIFIED, REMOVED, RENAMED sections)
- **parsers.go** (3 regex patterns → AST walker): Utility functions for counting tasks, deltas, and requirements

**Note**: These files contain Spectr-specific business logic (e.g., "what is a Requirement", "what is a Scenario"). They will walk the generic AST from `internal/mdparser`.

### UPDATED: internal/validation/ (2 files - CRITICAL)
- **parser.go** (5 regex patterns → shared AST walker, 260 lines): Duplicate parsing implementations for ExtractSections, ExtractRequirements, ExtractScenarios, ContainsShallOrMust, and NormalizeRequirementName
- **change_rules.go** (2 regex patterns → shared AST walker): Parses RENAMED delta sections in parseRenamedRequirements

**Critical fix**: The validation package currently contains duplicate parsing logic that suffers from the same brittleness issues. This duplication will be eliminated by using the shared `internal/mdparser` and `internal/parsers` logic.

### UPDATED: internal/archive/ (1 file)
- **spec_merger.go** (4 regex patterns → AST walker): Merges delta specs into base specs during archiving (reconstructSpec, splitSpec, extractOrderedRequirements)

### Architecture Summary
```
internal/mdparser/          # Generic, reusable markdown parser
├── lexer.go                # Tokenization (state machine)
├── parser.go               # AST construction (recursive descent)
└── ast.go                  # Generic AST nodes

internal/parsers/           # Spectr-specific extraction layer
├── requirement_parser.go   # Walks AST → Requirements, Scenarios
├── delta_parser.go         # Walks AST → Deltas (ADDED/MODIFIED/etc)
└── parsers.go              # Walks AST → counts, utilities

internal/validation/        # Uses shared parsers (no duplication)
internal/archive/           # Uses shared parsers (no duplication)
```

## Impact
- **Reliability**: Correctly handles all valid markdown and Spectr-specific constructs across all parsing contexts (parsing, validation, archiving).
- **Debuggability**: Provides precise error messages with line and column numbers throughout the system.
- **Extensibility**: Easier to add new syntax or features in the future by modifying the grammar and AST.
- **Code Quality**: Eliminates duplicate parsing logic between internal/parsers and internal/validation packages.
- **Consistency**: Ensures validation and archiving use the same robust parsing logic as the main parsers.
- **Reusability**: Creates a standalone markdown parser (`internal/mdparser`) with no Spectr dependencies that could be extracted for use in other projects.
- **Performance Validation**: Benchmarks ensure no unacceptable regression; correctness and maintainability take priority, with performance tracked and optimized if needed.
