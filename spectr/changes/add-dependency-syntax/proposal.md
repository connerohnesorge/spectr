# Change: Add Dependency Syntax for Change Proposals

## Why
Spectr currently lacks explicit dependency tracking between change proposals and their related changes or required specs. When a change builds on another change or requires specific capabilities to be in place, there is no formal way to declare these relationships. This leads to unclear implementation order, missed dependencies, and potential conflicts when changes are implemented out of sequence. By adding inline dependency syntax, we enable clear declaration of dependencies, better validation, and improved change orchestration.

## What Changes
- Add inline syntax for declaring change dependencies using `@depends(...)` syntax
- Add inline syntax for declaring spec requirements using `@requires(spec:...)` syntax
- Extend validation system to parse and validate dependency references
- Implement soft enforcement during development (warnings) and strict enforcement at archive time (errors)
- Update AGENTS.md documentation with dependency syntax examples and best practices
- Add validation rules to check that referenced changes and specs exist

## Impact
- **Affected specs**: `validation` (MODIFIED - adds dependency validation requirements)
- **Affected code**:
  - New `internal/parsers/dependency_parser.go` - Parse @depends() and @requires() syntax
  - New `internal/validation/dependency_rules.go` - Validate dependency references
  - Modified `internal/validation/types.go` - Add dependency-related types
  - Modified `internal/validation/validator.go` - Integrate dependency validation
  - Modified `spectr/AGENTS.md` - Document dependency syntax for AI assistants

## Benefits
- **Explicit dependencies**: Clear declaration of which changes depend on other changes
- **Validation**: Automatic checking that dependencies exist and are valid
- **Documentation**: Self-documenting proposals that show relationships
- **Orchestration**: Future tooling can use dependency data for ordering and planning
- **AI-friendly**: Clear syntax that AI assistants can easily parse and validate
