# Parsers Package

Parses requirements and delta operations from markdown AST. Bridge between markdown/ and domain/.

## OVERVIEW
Converts markdown/ AST nodes into domain types (Requirement, Scenario, DeltaOp). Supports incremental parsing updates.

## STRUCTURE
```go
internal/parsers/
├── parsers.go           # Main parsing entry points
├── delta_parser.go      # Delta operation parsing
├── parsers_test.go      # Table-driven tests
└── testdata/           # Fixture markdown files
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Parse requirements | ParseRequirements() | Extract from markdown AST |
| Parse deltas | ParseDelta() | ADDED/MODIFIED/REMOVED/RENAMED |
| Extract scenarios | ParseScenarios() | WHEN/THEN/AND bullets |

## CONVENTIONS
- **AST-based**: Parse from markdown/ nodes, not raw text
- **Strict errors**: Invalid structure returns ParseError
- **Table tests**: All parsers use t.Run() subtests with fixtures

## UNIQUE PATTERNS
- **Delta detection**: Identify operation headers (`## ADDED Requirements`) by name matching
- **Scenario parsing**: Extract `- **WHEN**`, `- **THEN**`, `- **AND**` bullets
- **Requirement identity**: Extract name from `### Requirement: <name>` header

## ANTI-PATTERNS
- **NO text parsing**: Use markdown/ AST, not regex on raw text
- **DON'T ignore sections**: All sections (Purpose, Requirements) must be parsed

## KEY FUNCTIONS
- `ParseRequirements(node markdown.Node) ([]*Requirement, []ParseError)` - Extract requirements
- `ParseScenarios(node markdown.Node) ([]*Scenario, []ParseError)` - Extract scenarios from requirement
- `ParseDelta(node markdown.Node) (*Delta, []ParseError)` - Parse delta operations

## DELTA TYPES
- **DeltaOp**: ADDED, MODIFIED, REMOVED, RENAMED
- **Delta**: Contains multiple DeltaOps per capability
- **Requirement**: Name, Description, Scenarios
- **Scenario**: Name, WHEN/THEN/AND steps

## FLOW
1. Parse Markdown source with `markdown.Parse()`
2. Find section headers with markdown.Find()
3. For each requirement section:
   - Extract requirement name from header
   - Parse scenario bullets (WHEN/THEN/AND)
4. For delta specs:
   - Identify operation headers (ADDED/MODIFIED/REMOVED/RENAMED)
   - Parse requirements under each operation
5. Return domain objects or ParseErrors
