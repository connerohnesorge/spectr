## Context

Spectr currently has 35+ inline error strings scattered across 10+ files. Only 2 errors are defined as constants/variables. This makes errors:
- Hard to find and audit
- Difficult to ensure consistency
- Impossible to reuse across packages
- Hard to test error handling with `errors.Is()` and `errors.As()`

## Goals / Non-Goals

**Goals:**
- Centralize all error definitions in one package
- Provide custom types with structured fields for programmatic access
- Support standard error wrapping (`Unwrap()`)
- Maintain backward-compatible error messages

**Non-Goals:**
- Creating a complex error hierarchy
- Adding error codes or internationalization
- Changing error message wording (preserve existing messages)

## Decisions

### Decision: Custom Types Only (No Sentinels)

**Rationale:** Custom types provide structured fields for rich context and enable type-based error checking with `errors.As()`. The user explicitly requested this approach.

**Alternatives considered:**
- Sentinel errors: Simpler but no structured context
- Mixed approach: More complex to maintain

### Decision: Domain-Based File Organization

**Rationale:** Groups related errors logically, scales well, and matches the existing `internal/` package structure.

**Alternatives considered:**
- Single file: Would become unwieldy with 20+ types
- Flat with prefixes: Less clear organization

### Decision: Pointer Receivers for Error() Method

All error types use pointer receivers (`*ErrorType`) for consistency and to enable optional fields.

## Package Structure

```
internal/specterrs/
├── doc.go           # Package documentation
├── git.go           # Git errors (5 types)
├── archive.go       # Archive errors (6 types)
├── validation.go    # Validation errors (3 types)
├── initialize.go    # Init errors (3 types)
├── list.go          # List/flag errors (1 type)
├── environment.go   # Environment errors (1 type)
└── pr.go            # PR workflow errors (2 types)
```

## Error Type Pattern

Each error type follows this pattern:

```go
// TypeNameError describes when X happens.
type TypeNameError struct {
    Field1 string // Contextual field
    Err    error  // For wrapping (optional)
}

func (e *TypeNameError) Error() string {
    // Human-readable message
}

func (e *TypeNameError) Unwrap() error {
    return e.Err
}
```

## Error Type Definitions

### Git Errors (`git.go`)

| Type | Fields | Message |
|------|--------|---------|
| `EmptyRemoteURLError` | - | "empty remote URL" |
| `BranchNameRequiredError` | - | "branch name is required" |
| `BaseBranchRequiredError` | - | "base branch is required" |
| `NotInGitRepositoryError` | `Path string` | "not in a git repository" |
| `BaseBranchNotFoundError` | `BranchName string` | "could not determine base branch..." |

### Archive Errors (`archive.go`)

| Type | Fields | Message |
|------|--------|---------|
| `UserCancelledError` | `Operation string` | "user cancelled selection" |
| `ArchiveCancelledError` | `Reason string` | "archive cancelled" |
| `ValidationRequiredError` | `Operation string` | "validation errors must be fixed before {operation}" |
| `DeltaConflictError` | `Section1, Section2, RequirementName string` | "requirement appears in both {s1} and {s2} sections" |

### Validation Errors (`validation.go`)

| Type | Fields | Message |
|------|--------|---------|
| `ValidationFailedError` | `ItemCount, ErrorCount, WarningCount int` | "validation failed" |
| `DeltaSpecParseError` | `SpecPath string, Line int, Err error` | "failed to parse delta spec..." |

### Initialize Errors (`initialize.go`)

| Type | Fields | Message |
|------|--------|---------|
| `EmptyPathError` | `Operation string` | "path cannot be empty" |
| `WizardModelCastError` | `ActualType string` | "failed to cast final model to WizardModel" |
| `InitializationCompletedWithErrorsError` | `ErrorCount int, Errors []error` | "initialization completed with errors" |

### List Errors (`list.go`)

| Type | Fields | Message |
|------|--------|---------|
| `IncompatibleFlagsError` | `Flag1, Flag2 string` | "cannot use {flag1} with {flag2}" |

### Environment Errors (`environment.go`)

| Type | Fields | Message |
|------|--------|---------|
| `EditorNotSetError` | `Operation string` | "EDITOR environment variable not set" |

### PR Errors (`pr.go`)

| Type | Fields | Message |
|------|--------|---------|
| `UnknownPlatformError` | `Platform, RepoURL string` | "unknown platform; please create PR manually" |
| `PRPrerequisiteError` | `Check, Details string, Err error` | "PR prerequisite failed ({check}): {details}" |

## Migration Plan

1. Create package with all error types (no code changes)
2. Migrate errors by domain, least to most coupled:
   - Environment (1 error)
   - List (2 errors)
   - Initialize (4 errors)
   - Git (5 errors)
   - PR (2 errors)
   - Validation (3 errors)
   - Archive (11 errors)
3. Run tests after each domain migration
4. Remove old constants/sentinels after migration complete

## Risks / Trade-offs

**Risk:** Breaking error checks using string matching
**Mitigation:** Preserve exact error messages; `errors.As()` works regardless

**Risk:** Increased verbosity in return statements
**Mitigation:** Optional fields allow simple `&ErrorType{}` usage
