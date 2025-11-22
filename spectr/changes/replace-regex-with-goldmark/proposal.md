# Change: Replace Regex Parsing with Goldmark AST

## Why

The current markdown parsing implementation relies heavily on regex patterns and line-by-line scanning, which has several limitations:

1. **Fragility**: Regex patterns are brittle and fail on edge cases like code blocks containing headers, nested structures, or unusual whitespace
2. **Limited Context**: Line-by-line scanning loses document structure, making it difficult to provide accurate error locations and context
3. **Maintainability**: Regex patterns are hard to read, test, and extend as spec format evolves
4. **Performance**: Multiple passes over files for different parsing tasks (titles, requirements, scenarios, deltas)

Goldmark is the canonical Go markdown parser, providing a robust AST (Abstract Syntax Tree) that:
- Handles all CommonMark edge cases correctly
- Preserves source positions for accurate error reporting
- Enables single-pass parsing with complete document structure
- Supports extensions for custom syntax (delta sections)
- Maintained by the Go community with comprehensive test coverage

## What Changes

### Core Parsing (Complete Replacement)
- **Remove**: All regex-based line scanning in `internal/parsers/parsers.go`, `requirement_parser.go`, `delta_parser.go`
- **Add**: Goldmark dependency (`github.com/yuin/goldmark v1.7.0+`)
- **Add**: AST walker utilities for extracting titles, requirements, scenarios
- **Add**: Custom goldmark walker for delta sections (`ADDED`, `MODIFIED`, `REMOVED`, `RENAMED`)
- **Add**: AST post-processing for rename format (`FROM:`/`TO:` bullets)

### Data Structures (Improved Design)
- **Enhance** `RequirementBlock`: Add `SourcePos` field with line/column information from AST
- **Enhance** `DeltaPlan`: Add source location metadata for better error reporting
- **Add**: `ASTParser` interface for extensibility
- **Preserve**: Backward compatibility for all public APIs used by archive/validation/list/view

### Affected Capabilities
- **validation**: Parsing logic changes, error reporting gains AST context (line:column references)
- **cli-framework**: All parsing functions migrate to goldmark (titles, tasks, requirements, deltas)
- **archive-workflow**: Requirement merging uses AST-aware comparison for better conflict detection

## Impact

### Breaking Changes
**None** - All public APIs remain compatible. Internal parsing implementation is fully encapsulated.

### Affected Code
- `internal/parsers/parsers.go` - Replace `ExtractTitle`, `CountRequirements` with AST-based implementations
- `internal/parsers/requirement_parser.go` - Replace `ParseRequirements`, `ParseScenarios` with AST traversal
- `internal/parsers/delta_parser.go` - Replace regex section extraction with AST walking
- `internal/validation/*` - Update to use enhanced error locations from AST
- `internal/archive/spec_merger.go` - May leverage AST metadata for better merging
- `go.mod` - Add goldmark dependency

### Testing Strategy
- Parallel test execution (regex vs goldmark) during migration to validate equivalence
- New test cases for goldmark-specific edge cases (code blocks, nested lists, unusual markdown)
- Performance benchmarks on real spec files to ensure no regression
- Integration tests with existing validation and archive workflows

### Migration Risk
**Low** - Changes are internal implementation details. Comprehensive test coverage exists. Phased rollout allows validation at each step.

### Benefits
1. **Reliability**: Handles all CommonMark edge cases correctly out of the box
2. **Better Errors**: Precise line:column error locations from AST source positions
3. **Extensibility**: Easy to add new parsing features using AST traversal
4. **Performance**: Single-pass parsing replaces multiple regex scans
5. **Maintainability**: Clear AST traversal code replaces cryptic regex patterns
6. **Future-Proof**: Built on canonical Go markdown parser with active maintenance
