# Change: Replace Regex Parsing with Goldmark AST

## Why
Current regex-based line-by-line parsing is brittle for edge cases (code blocks, nested markdown structures), difficult to maintain and extend, and doesn't properly handle standard markdown syntax. This creates fragile parsing that breaks on valid markdown documents.

## What Changes
- Replace all regex line-matching with goldmark AST traversal
- Migrate parsers.go, requirement_parser.go, and delta_parser.go to goldmark-based implementation
- Add custom goldmark AST walker for delta sections (ADDED/MODIFIED/REMOVED/RENAMED)
- Improve internal data structures (RequirementBlock, DeltaPlan) with AST metadata for better accuracy
- Add goldmark dependency (github.com/yuin/goldmark) to go.mod
- **INTERNAL BREAKING**: Parser function signatures may change to accommodate AST-based approach
- Enhance error reporting with line numbers and context from AST nodes

## Impact
- Affected specs:
  - validation (core parsing logic)
  - cli-framework (parser function interfaces)
- Affected code:
  - internal/parsers/* (3 core files: parsers.go, requirement_parser.go, delta_parser.go)
  - internal/parsers/*_test.go (3 test files)
  - internal/archive/spec_merger.go (requirement matching)
  - internal/validation/change_rules.go (delta validation)
- internal/validation/delta_validators.go (scenario validation)
- internal/list/lister.go (requirement counting)
- internal/view/dashboard.go (spec display)
- Dependencies: Add github.com/yuin/goldmark to go.mod
- Testing: All existing parser tests (parsers_test.go, requirement_parser_test.go, delta_parser_test.go) must pass with goldmark implementation
- Performance: AST parsing may have different performance characteristics (measure before/after)
