# Design: Lexer/Parser Architecture

## Architecture

The parsing logic will be split into three distinct phases: Lexing, Parsing, and Extraction.

### 1. Lexer
The lexer will be modeled after the Go compiler's scanner (`cmd/compile/internal/syntax/scanner.go`). It will scan the input text and emit a stream of `Token`s.

**Key Components:**
-   `TokenType`: Enum for token types (e.g., `TokenHeader`, `TokenText`, `TokenCodeBlock`, `TokenEOF`).
-   `Token`: Struct containing `Type`, `Value` (string), and `Position` (line/col).
-   `Lexer`: Struct holding the input, current position, and a channel/slice of tokens.
-   `StateFn`: Function type `func(*Lexer) StateFn` (adapting the Rob Pike approach to fit the Go compiler's structural style where appropriate).

**States:**
-   `lexText`: Default state, consumes plain text.
-   `lexHeader`: Handles markdown headers (`#`, `##`, etc.).
-   `lexCodeBlock`: Handles code fences (```` ``` ````) and content within them.
-   `lexList`: Handles list items.

### 2. Parser
The parser will consume the stream of tokens and build an Abstract Syntax Tree (AST), following the recursive descent pattern used in the Go compiler (`cmd/compile/internal/syntax/parser.go`).

**AST Nodes:**
-   `Node`: Interface for all AST nodes.
-   `Document`: Root node.
-   `Header`: Represents a markdown header.
-   `Paragraph`: Represents a block of text.
-   `CodeBlock`: Represents a code block.
-   `List`: Represents a list.

The parser will be a recursive descent parser that dispatches based on the current token type, ensuring robust error handling and structure recovery.

### 3. Extractor
The extractor will traverse the AST to identify and validate Spectr-specific structures. It effectively replaces the current "parsing" logic but operates on a structured tree instead of raw text.

**Logic:**
-   Find `Header` nodes with specific text (e.g., "Requirement:", "Scenario:").
-   Validate hierarchy (e.g., Scenario must be under a Requirement).
-   Extract content from `Paragraph` and `CodeBlock` nodes associated with the headers.

## Code Sharing Architecture

The new lexer/parser will be implemented as a shared component that eliminates duplicate parsing logic across packages.

### Current Duplication Problem

Currently, there are three separate parsing implementations:
1. **internal/parsers/**: Primary parsing for requirements and delta specs
2. **internal/validation/parser.go**: Duplicate implementations of ExtractSections, ExtractRequirements, ExtractScenarios
3. **internal/archive/spec_merger.go**: Additional regex-based parsing for spec merging

This duplication means:
- Bug fixes must be applied in multiple places
- Inconsistent behavior between parsing contexts
- Same brittleness issues affect all packages

### Unified Architecture

The new design will establish a single source of truth:

```
internal/parser/           (or internal/parsers/lexer/)
├── lexer.go              → Tokenization
├── parser.go             → AST construction
└── extractor.go          → Spectr element extraction

internal/parsers/
├── requirement_parser.go → Uses shared parser
└── delta_parser.go       → Uses shared parser

internal/validation/
├── parser.go             → Uses shared parser (no more local regex)
└── change_rules.go       → Uses shared parser for RENAMED parsing

internal/archive/
└── spec_merger.go        → Uses shared parser for merging
```

### Integration Strategy

**Phase 1: Build Core (Tasks 1.1-1.7)**
- Implement lexer, parser, and extractor with comprehensive test coverage
- Ensure API supports all current use cases

**Phase 2: Replace Parsers (Tasks 2.1-2.3)**
- Update internal/parsers/ to use new shared implementation
- Verify existing tests pass

**Phase 3: Refactor Validation (Tasks 3.1-3.6)**
- Replace ExtractSections with shared parser API
- Replace ExtractRequirements with shared parser API
- Replace ExtractScenarios with shared parser API
- Update ContainsShallOrMust to use lexer tokens
- Update NormalizeRequirementName (or integrate into shared parser)
- Replace parseRenamedRequirements in change_rules.go

**Phase 4: Update Archive (Tasks 4.1-4.3)**
- Replace reconstructSpec normalization logic
- Replace splitSpec section splitting
- Replace extractOrderedRequirements with shared parser

**Phase 5: Validate (Tasks 5.1-5.4)**
- Comprehensive testing across all packages
- Edge case validation (markdown in code blocks)
- Integration testing with full workflows

### Benefits of Shared Architecture

- **Single Source of Truth**: One parser implementation for all contexts
- **Consistent Behavior**: Same parsing logic in validation, archiving, and extraction
- **Easier Maintenance**: Bug fixes and enhancements in one place
- **Better Testing**: Comprehensive tests on shared parser benefit all consumers
- **Reduced Code**: Eliminates ~300+ lines of duplicate parsing logic

## Trade-offs

### Pros
-   **Robustness**: Explicitly handles state, preventing misinterpretation of syntax inside code blocks.
-   **Precision**: Tracks line and column numbers for every token, enabling high-quality error messages.
-   **Separation of Concerns**: Lexing (syntax) is separated from Parsing (structure) and Extraction (semantics).

### Cons
-   **Complexity**: More code than the current regex approach.
-   **Performance**: Likely slower than raw regex for simple cases (though potentially faster for complex ones due to single pass), but performance is not the primary bottleneck for Spectr.

## Migration Strategy
The new parser will be implemented in a new package (e.g., `internal/parser/new` or just `internal/parser` if replacing). We will switch over consumers one by one or via a feature flag if needed, but given the scope, a direct replacement in a single release is feasible.
