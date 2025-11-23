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

## Impact
- **Reliability**: Correctly handles all valid markdown and Spectr-specific constructs.
- **Debuggability**: Provides precise error messages with line and column numbers.
- **Extensibility**: Easier to add new syntax or features in the future by modifying the grammar and AST.
