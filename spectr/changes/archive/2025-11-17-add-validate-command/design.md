# Design: Validate Command

## Context

The validate command needs to parse markdown files, extract structured
information (sections, requirements, scenarios), apply validation rules, and
report issues with helpful guidance. The OpenSpec reference implementation
provides a proven architecture using TypeScript, which we'll adapt to Go's
idioms and type system.

### Constraints

- Must work with existing Spectr directory structure (`spectr/specs/`,
  `spectr/changes/`)
- Must integrate with Kong CLI framework
- Must support both interactive and non-interactive modes
- Must be performant for projects with many specs/changes
- Should follow Go best practices (error handling, interfaces, testability)

### Stakeholders

- Developers using Spectr for spec-driven development
- CI/CD systems requiring validation before merge
- AI assistants that need to validate generated specs

## Goals / Non-Goals

### Goals

- Implement comprehensive validation matching OpenSpec functionality
- Provide clear, actionable error messages with remediation steps
- Support bulk validation with parallel processing for performance
- Enable CI/CD integration via JSON output and exit codes
- Maintain compatibility with Spectr conventions

### Non-Goals

- Auto-fixing validation errors (out of scope for initial implementation)
- Web UI for validation results (CLI-only for now)
- Real-time validation during editing (editor integration is separate)
- Schema evolution or migration tools (validation only)

## Decisions

### Decision 1: Package Structure

**Choice**: Separate `internal/validation` and `internal/discovery` packages

**Rationale**:

- **Separation of concerns**: Discovery (finding items) vs validation (checking
  items)
- **Testability**: Each package can be unit tested independently
- **Reusability**: Discovery package can be used by future commands (show, list,
  archive)
- **Go idioms**: Flat, focused packages are more idiomatic than deep hierarchies

**Structure**:

```text
internal/
├── validation/
│   ├── validator.go        # Main Validator interface and orchestration
│   ├── types.go            # ValidationReport, ValidationIssue, etc.
│   ├── spec_rules.go       # Spec-specific validation rules
│   ├── change_rules.go     # Change delta validation rules
│   └── parser.go           # Markdown parsing utilities
├── discovery/
│   └── discovery.go        # Find changes and specs in filesystem
└── init/                   # Existing init package
```text

**Alternatives considered**:

- Single `internal/validate` package: Rejected due to mixing concerns
- Nesting under `cmd/validate/`: Rejected because internal packages should be
  reusable

### Decision 2: Markdown Parsing Library

**Choice**: Use `github.com/gomarkdown/markdown` for parsing, with custom
extraction

**Rationale**:

- **Standard library first**: Go's `bufio.Scanner` can handle simple section
  extraction
- **Goldmark alternative**: `github.com/yuin/goldmark` is more featureful but
  adds complexity
- **Custom parsing**: Spec format is well-defined, so custom regex-based
  extraction is tractable
- **Zero dependencies preference**: Start simple, add library only if needed
  during implementation

**Implementation approach**:

1. Use `bufio.Scanner` to read lines
2. Use regex to match headers (`##`, `###`, `####`)
3. Build section map while iterating
4. Extract requirements and scenarios using pattern matching

**Alternatives considered**:

- Full AST parsing with goldmark: Overkill for structured document validation
- Generic markdown parser: Too general, would require extensive post-processing

### Decision 3: Validation Architecture

**Choice**: Strategy pattern with separate validators for specs and changes

**Rationale**:

- **Type safety**: Spec validation and change validation have different rules
- **Extensibility**: Easy to add new validation rules or validators
- **Testing**: Can test spec and change validation independently
- **Clear contracts**: Explicit interfaces make behavior predictable

**Interface design**:

```go
type Validator interface {
    ValidateSpec(path string) (*ValidationReport, error)
    ValidateChange(changeDir string) (*ValidationReport, error)
}

type ValidationReport struct {
    Valid   bool
    Issues  []ValidationIssue
    Summary ValidationSummary
}

type ValidationIssue struct {
    Level   ValidationLevel // ERROR, WARNING, INFO
    Path    string          // Section/requirement path
    Message string          // Human-readable message
}
```text

**Alternatives considered**:

- Single validate method with type switching: Less type-safe, harder to test
- Separate validator types: Possible, but interface unifies the API

### Decision 4: Parallel Validation Strategy

**Choice**: Worker pool with configurable concurrency, default 6 workers

**Rationale**:

- **Performance**: Validation is I/O-bound (file reads), parallelism helps
- **Resource limits**: Worker pool prevents resource exhaustion
- **Compatibility**: Matches OpenSpec default concurrency (6)
- **Go idioms**: Goroutines and channels are idiomatic for this pattern

**Implementation**:

- Use `sync.WaitGroup` for coordination
- Channel-based work queue
- Configurable via flag or environment variable
- Graceful degradation if concurrency=1 (serial execution)

**Alternatives considered**:

- Unlimited goroutines: Could exhaust file descriptors or memory
- Serial-only execution: Too slow for large projects
- Dependency on third-party pool library: Unnecessary for this use case

### Decision 5: Error Reporting and Guidance

**Choice**: Multi-level issues (ERROR/WARNING/INFO) with contextual remediation
messages

**Rationale**:

- **User experience**: Helpful messages reduce friction
- **Strict mode**: Teams can enforce zero warnings if desired
- **Gradual adoption**: Warnings allow incremental improvements
- **Compatibility**: Matches OpenSpec behavior

**Message structure**:

- Issue level and path clearly indicated
- Primary error message describes what's wrong
- Remediation guidance explains how to fix it
- Examples show correct format

**Example**:

```text
✗ [ERROR] auth/spec.md: Requirement "User Authentication" must include
  at least one scenario
Next steps:
  - Add scenario using #### Scenario: format
  - Example:
    #### Scenario: Login success
    - **WHEN** valid credentials provided
    - **THEN** JWT token returned
```text

**Alternatives considered**:

- Error codes only: Not helpful for humans
- Verbose documentation links: Adds complexity, docs may drift
- Auto-fix suggestions: Out of scope, requires more complex implementation

### Decision 6: CLI Flags and Modes

**Choice**: Support multiple invocation modes matching OpenSpec UX

**Flags**:

- `--strict`: Treat warnings as errors (exit code 1)
- `--json`: Machine-readable output for CI/CD
- `--all`: Validate all changes and specs
- `--changes`: Validate only changes
- `--specs`: Validate only specs
- `--type change|spec`: Disambiguate when item name exists in both
- `--no-interactive`: Disable interactive prompts (for CI)

**Modes**:

1. **Direct item**: `spectr validate <item-name>`
2. **Bulk validation**: `spectr validate --all`
3. **Interactive selection**: `spectr validate` (TTY only)

**Rationale**:

- **Compatibility**: Matches OpenSpec for familiar UX
- **Flexibility**: Supports both human and machine use cases
- **CI-friendly**: Non-interactive mode prevents hangs in pipelines

## Risks / Trade-offs

### Risk: Markdown Parsing Complexity

**Mitigation**:

- Start with simple regex-based parsing for well-defined format
- Add proper parser library if edge cases emerge during implementation
- Comprehensive test suite with realistic spec examples

### Risk: Performance with Large Projects

**Mitigation**:

- Parallel validation with worker pool
- Configurable concurrency for tuning
- Future optimization: caching validation results based on file mtime

### Risk: Inconsistent Validation Between OpenSpec and Spectr

**Mitigation**:

- Port OpenSpec validation rules directly
- Shared test cases between implementations
- Document any intentional differences

### Trade-off: Go vs TypeScript Idioms

**Decision**: Favor Go idioms over direct TypeScript port

- Use interfaces instead of classes
- Use error returns instead of exceptions
- Use channels instead of promises
- Use struct tags for JSON marshaling

This maintains Go best practices while preserving validation semantics.

## Migration Plan

### Phase 1: Core Implementation

1. Implement `internal/validation` package
2. Implement `internal/discovery` package
3. Add `cmd/validate.go` with basic flags
4. Unit tests for validation rules

### Phase 2: CLI Integration

1. Wire validate command into Kong CLI
2. Implement interactive mode (if TTY)
3. Implement bulk validation
4. Integration tests

### Phase 3: Polish

1. Add helpful error messages and guidance
2. JSON output format
3. Performance tuning (parallel validation)
4. Documentation updates

### Rollback Plan

If validation proves problematic:

1. Command can be disabled via feature flag
2. No breaking changes to existing commands
3. Validation package is isolated, can be removed cleanly

## Open Questions

### Question 1: Should we validate referenced specs during change validation?

**Context**: When a change references a spec in `## MODIFIED Requirements`,
should we verify that spec exists?

**Options**:

- **A**: No cross-validation (simple, what OpenSpec does)
- **B**: Warn if spec doesn't exist (helpful but adds complexity)
- **C**: Error if spec doesn't exist (strict but may block valid workflows)

**Recommendation**: Start with Option A (OpenSpec behavior), add Option B later
if users request it.

### Question 2: Should validation cache results?

**Context**: Validating the same unchanged file multiple times is wasteful.

**Options**:

- **A**: No caching (simple, always correct)
- **B**: In-memory cache for single command invocation (helps bulk validation)
- **C**: Persistent cache based on file mtime (complex, high risk of stale data)

**Recommendation**: Start with Option A, add Option B if profiling shows
redundant work.

### Question 3: How to handle non-standard scenario formats?

**Context**: Users might write scenarios in different formats (bullets, numbered
lists, etc.)

**Options**:

- **A**: Strict enforcement of `#### Scenario:` format (simple, consistent)
- **B**: Accept alternatives with warnings (flexible but inconsistent)
- **C**: Auto-detect and normalize (complex, risky)

**Recommendation**: Option A (strict enforcement) matches OpenSpec and ensures
consistency. Users adapt quickly to clear requirements.
