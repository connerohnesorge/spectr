# Error Handling Specification

## Requirements

### Requirement: Centralized Error Package

The system SHALL provide a centralized `internal/specterrs` package containing all custom error types used throughout the codebase.

#### Scenario: Import error types from specterrs

- **WHEN** a package needs to return a domain-specific error
- **THEN** it SHALL import from `internal/specterrs` and use the appropriate typed error

#### Scenario: Error message preservation

- **WHEN** migrating from inline `errors.New()` to custom types
- **THEN** the error messages SHALL remain identical to preserve backward compatibility

### Requirement: Domain-Based Error Organization

The system SHALL organize error types into domain-specific files within the specterrs package.

#### Scenario: Git errors in git.go

- **WHEN** defining git-related errors (repository, branch, remote operations)
- **THEN** they SHALL be defined in `internal/specterrs/git.go`

#### Scenario: Archive errors in archive.go

- **WHEN** defining archive workflow errors (cancellation, conflicts, validation)
- **THEN** they SHALL be defined in `internal/specterrs/archive.go`

#### Scenario: Validation errors in validation.go

- **WHEN** defining spec/change validation errors
- **THEN** they SHALL be defined in `internal/specterrs/validation.go`

### Requirement: Custom Error Types with Structured Fields

The system SHALL use custom struct types for all errors, with optional fields for contextual information.

#### Scenario: Error with context fields

- **WHEN** an error benefits from additional context (e.g., file path, operation name)
- **THEN** the error type SHALL include struct fields for that context

#### Scenario: Error interface implementation

- **WHEN** defining a custom error type
- **THEN** it SHALL implement the `error` interface via an `Error() string` method with a pointer receiver

#### Scenario: Error wrapping support

- **WHEN** an error type wraps an underlying error
- **THEN** it SHALL implement `Unwrap() error` to support `errors.Is()` and `errors.As()`

### Requirement: No Sentinel Errors

The system SHALL NOT define sentinel error variables (e.g., `var ErrFoo = errors.New(...)`).

#### Scenario: Existing sentinels removed

- **WHEN** the migration is complete
- **THEN** the existing `ErrUserCancelled` sentinel SHALL be removed and replaced with `UserCancelledError` type

#### Scenario: Error checking with types

- **WHEN** code needs to check for a specific error condition
- **THEN** it SHALL use `errors.As()` with a pointer to the error type instead of `errors.Is()` with a sentinel
