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
We will implement a hand-written lexer and parser modeled after the current Go compiler (`cmd/compile/internal/syntax`):
1.  **Lexer**: A state-machine-based lexer that tokenizes the input stream. It will handle state transitions (e.g., entering/exiting code blocks) and emit tokens, similar to the Go compiler's scanner.
2.  **Parser**: A recursive descent parser that consumes tokens and builds an Abstract Syntax Tree (AST) representing the document structure. This mirrors the Go compiler's parser design for handling structure and error recovery.
3.  **Extractor**: A semantic analysis layer that walks the AST to extract Spectr-specific structures (Requirements, Scenarios, Deltas) and validates them.

## Affected Files
This change will replace regex usage in 6 files across 3 packages:

### internal/parsers/ (3 files)
- **requirement_parser.go** (3 regex patterns): Parses requirement blocks and scenarios from spec files
- **delta_parser.go** (11 regex patterns): Parses delta specs (ADDED, MODIFIED, REMOVED, RENAMED sections)
- **parsers.go** (3 regex patterns): Utility functions for counting tasks, deltas, and requirements

### internal/validation/ (2 files) - CRITICAL
- **parser.go** (5 regex patterns, 260 lines): Duplicate parsing implementations for ExtractSections, ExtractRequirements, ExtractScenarios, ContainsShallOrMust, and NormalizeRequirementName
- **change_rules.go** (2 regex patterns): Parses RENAMED delta sections in parseRenamedRequirements

### internal/archive/ (1 file)
- **spec_merger.go** (4 regex patterns): Merges delta specs into base specs during archiving (reconstructSpec, splitSpec, extractOrderedRequirements)

**Note:** The validation package has the most critical gap - it contains duplicate parsing logic that currently suffers from the same brittleness issues as the main parsers. This duplication will be eliminated by having validation use the new shared lexer/parser.

## Impact
- **Reliability**: Correctly handles all valid markdown and Spectr-specific constructs across all parsing contexts (parsing, validation, archiving).
- **Debuggability**: Provides precise error messages with line and column numbers throughout the system.
- **Extensibility**: Easier to add new syntax or features in the future by modifying the grammar and AST.
- **Code Quality**: Eliminates duplicate parsing logic between internal/parsers and internal/validation packages.
- **Consistency**: Ensures validation and archiving use the same robust parsing logic as the main parsers.
