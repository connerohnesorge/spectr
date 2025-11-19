# Design: Dependency Syntax for Change Proposals

## Context

Spectr change proposals currently operate in isolation without formal dependency tracking. When changes build upon other changes or require specific capabilities to exist, these relationships are documented only in prose (if at all). This creates several problems:

1. **Unclear ordering**: No way to determine which changes should be implemented first
2. **Missing context**: Reviewers don't see dependencies without reading full proposals
3. **Validation gaps**: Tools can't verify that prerequisites exist
4. **Orchestration limits**: Future tooling can't automate dependency resolution

The goal is to add a lightweight, inline syntax that makes dependencies explicit, parseable, and validatable.

### Constraints
- Must be easy to read and write (human-friendly)
- Must be parseable by simple regex (no complex grammar)
- Must integrate with existing validation system
- Must not break existing proposals (backward compatible)
- Must work well with AI assistants (clear syntax)

### Stakeholders
- Developers writing change proposals
- AI assistants generating proposals
- Code reviewers evaluating dependencies
- Future tooling for dependency resolution and ordering

## Goals / Non-Goals

### Goals
- Enable explicit declaration of change dependencies via inline syntax
- Enable explicit declaration of spec requirements via inline syntax
- Validate that referenced changes and specs actually exist
- Provide clear error messages when dependencies are missing
- Support soft enforcement (warnings) during development, strict enforcement (errors) at archive time
- Document syntax clearly for both humans and AI assistants

### Non-Goals
- **Circular dependency detection** - Deferred to future enhancement (would require graph traversal)
- **Automatic dependency ordering** - Future enhancement for orchestration
- **Dependency version constraints** - Changes are immutable once archived, no versioning needed
- **Cross-repository dependencies** - Out of scope, Spectr is single-repo focused
- **Automatic dependency resolution** - Tool should validate, not auto-fix

## Decisions

### Decision 1: Inline Syntax Using @ Directives
**Choice**: Use `@depends(change-id)` and `@requires(spec:capability-id)` inline in proposal.md

**Rationale**:
- **Inline placement**: Keeps dependencies close to the explanation, not in separate metadata
- **@ prefix**: Familiar from social media mentions, git blame annotations, docstring tags
- **Explicit type markers**: `@depends` for changes, `@requires` for specs - clear distinction
- **Kebab-case IDs**: Matches existing Spectr conventions for change and spec naming
- **Namespace separation**: `spec:` prefix in @requires distinguishes specs from changes
- **Simple parsing**: Regex-friendly syntax, no complex grammar needed

**Example usage in proposal.md**:
```markdown
## Why
This change builds on the work from @depends(add-validation-command) to add
dependency tracking. It requires the @requires(spec:validation) capability
to be implemented first.

The syntax enables explicit dependency declarations that can be validated
automatically, improving orchestration and preventing missing prerequisites.
```

**Alternatives considered**:
- **YAML frontmatter**: More structured but harder to read in prose, requires YAML parser
- **Separate dependencies.md file**: Separates deps from context, easier to miss
- **List-based syntax in dedicated section**: More verbose, separates from explanation
- **Comment-based syntax**: Easy to overlook, not prominent enough

### Decision 2: Soft Enforcement During Development, Strict at Archive
**Choice**: Generate WARNING for missing dependencies in normal validation, ERROR in strict mode / archive

**Rationale**:
- **Development flexibility**: Allows working on dependent changes in parallel
- **Archive safety**: Ensures dependencies exist before marking change complete
- **Incremental workflow**: Can create proposals that reference not-yet-created dependencies
- **CI/CD compatibility**: Archive validation (strict mode) prevents incomplete deployments

**Implementation**:
- Normal `spectr validate` generates WARNING level issues for missing deps
- `spectr validate --strict` generates ERROR level issues (used during archive)
- Validation checks both `spectr/changes/` and `spectr/changes/archive/` for changes
- Validation checks `spectr/specs/` for required specs

**Alternatives considered**:
- **Always strict**: Too rigid, prevents parallel work on dependent changes
- **Always soft**: Allows archiving changes with missing dependencies (unsafe)
- **Manual flag control**: Confusing UX, archive should always be strict

### Decision 3: No Circular Dependency Detection (This Change)
**Choice**: Document circular dependencies as possible, but don't validate for them yet

**Rationale**:
- **Complexity**: Requires graph traversal, cycle detection algorithms
- **Rare case**: Most dependencies are linear or DAG-like
- **Future work**: Can add later without breaking existing syntax
- **YAGNI principle**: Don't build it until we need it

**Mitigation**:
- Document that circular dependencies are possible but not recommended
- Add TODO comment in validation code for future enhancement
- Design.md explains why detection is deferred

**Future implementation notes**:
- Build dependency graph from all active changes
- Use DFS or Tarjan's algorithm to detect cycles
- Report cycles as ERROR level issues
- Could be added in separate change proposal

### Decision 4: Validation Integration via Parser + Rules
**Choice**: Create separate parser package and validation rules file, integrate into existing validator

**Architecture**:
```
internal/parsers/dependency_parser.go
  - ParseDependencies(content string) []DependencyReference
  - Regex-based extraction
  - Returns typed references (CHANGE or SPEC)

internal/validation/dependency_rules.go
  - ValidateDependencies(deps []DependencyReference, projectPath string, strictMode bool) []ValidationIssue
  - Checks existence via discovery package
  - Generates appropriate WARNING or ERROR issues

internal/validation/validator.go (modified)
  - In ValidateChange(), parse dependencies from proposal.md
  - Call ValidateDependencies() and merge issues into report
```

**Rationale**:
- **Separation of concerns**: Parsing logic separate from validation logic
- **Reusability**: Parser can be used by future tools (show, list, etc.)
- **Testability**: Each component can be unit tested independently
- **Go idioms**: Flat package structure, clear interfaces

**Alternatives considered**:
- **Single validation file**: Mixes parsing and validation, harder to test
- **Inline parsing in validator**: Tight coupling, not reusable
- **Complex AST**: Overkill for simple regex-based syntax

### Decision 5: Line Number Tracking for Better Errors
**Choice**: Parser tracks line numbers where dependencies are declared

**Rationale**:
- **Better error messages**: "proposal.md:12: Missing dependency @depends(foo)"
- **Debugging aid**: Helps locate issues in long proposals
- **Low cost**: Simple to implement during regex matching

**Implementation**:
```go
type DependencyReference struct {
    Type     DependencyType  // CHANGE or SPEC
    ID       string          // change-id or capability-id
    Location string          // "proposal.md:12"
    Line     int             // 12
}
```

## Risks / Trade-offs

### Risk: Syntax Ambiguity in Prose
**Problem**: `@depends(foo)` might appear in code examples or quoted text, causing false positives

**Mitigation**:
- Document that dependency declarations should be in "Why" or "What Changes" sections
- If false positives become an issue, add comment-based override: `<!-- no-dep-check -->`
- Parser can be enhanced to skip code blocks (fenced with ```) in future

### Risk: Incomplete Dependency Coverage
**Problem**: Not all dependencies might be declared (human error)

**Mitigation**:
- Validation can only check declared dependencies, not infer missing ones
- Review process should catch missing dependencies
- Over time, best practices will emerge around when to declare deps
- AI assistants can be prompted to identify and declare dependencies

### Trade-off: Flexibility vs Strictness
**Decision**: Prioritize flexibility during development, strictness at archive time

**Rationale**: Allows iterative workflow (create multiple related proposals), but ensures completeness before marking done

**Impact**: Teams need to understand the dual enforcement model (warnings vs errors)

### Trade-off: Simple Syntax vs Rich Metadata
**Decision**: Keep syntax minimal (just ID reference), no extra metadata

**What we're NOT including**:
- Dependency versions (changes are immutable)
- Optional vs required dependencies (all are required)
- Dependency reasons/descriptions (prose explains why)

**Rationale**: YAGNI - add complexity only when proven necessary

## Migration Plan

### Phase 1: Implementation (This Change)
1. Implement dependency parser in `internal/parsers/dependency_parser.go`
2. Create validation rules in `internal/validation/dependency_rules.go`
3. Integrate into `internal/validation/validator.go`
4. Update `spectr/AGENTS.md` with syntax documentation
5. Write comprehensive tests

### Phase 2: Adoption
1. Update existing active change proposals to add dependencies (optional)
2. Document best practices in project documentation
3. AI assistants automatically include dependencies in new proposals
4. Review process checks for dependency declarations

### Phase 3: Future Enhancements (Separate Changes)
1. Add circular dependency detection
2. Add `spectr deps` command to visualize dependency graphs
3. Add automatic ordering suggestions based on dependencies
4. Add dependency validation in `spectr archive` command

### Rollback Plan
If dependency syntax proves problematic:
- Parser and validation rules are isolated, can be removed
- Syntax is optional, existing proposals without deps continue to work
- Can disable dependency validation via feature flag if needed

## Open Questions

### Question 1: Should @depends support archived changes?
**Context**: When a change depends on something already archived, should we validate against archive too?

**Options**:
- **A**: Only check active changes (simple, encourages cleanup)
- **B**: Check both active and archive (complete, but more complex)
- **C**: Error if depending on archived change (strict, but limits reuse)

**Recommendation**: Option B - Check both `spectr/changes/` and `spectr/changes/archive/`. Reason: Changes can legitimately build on archived work, validation should confirm the dependency existed at some point.

**Implementation**: Discovery function searches both directories.

### Question 2: Should we support @depends on multiple changes in single declaration?
**Context**: Some changes might depend on multiple other changes.

**Options**:
- **A**: Multiple separate declarations: `@depends(foo) @depends(bar)`
- **B**: Comma-separated: `@depends(foo, bar)`
- **C**: Both supported

**Recommendation**: Option A (multiple separate declarations) for initial implementation. Reason: Simpler parsing, clear one-dependency-per-declaration, easier error reporting. Can add Option C later if users request it.

### Question 3: Should validation check dependency order in archives?
**Context**: If change-A depends on change-B, should change-B's archive date be before change-A's?

**Options**:
- **A**: No timestamp validation (simple, allows retroactive deps)
- **B**: Warn if dependency archived after dependent (helps catch errors)
- **C**: Error if dependency archived after dependent (strict)

**Recommendation**: Defer to future enhancement. Start with Option A. Archive timestamp validation adds complexity and may not be necessary if teams follow proper workflow.
