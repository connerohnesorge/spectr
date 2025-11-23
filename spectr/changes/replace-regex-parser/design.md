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
