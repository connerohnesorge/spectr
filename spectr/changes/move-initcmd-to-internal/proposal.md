# Change: Move InitCmd to internal/init package

## Why

Currently, `InitCmd` is defined in `cmd/init.go` alongside its `Run()` method and helper functions. Moving the command struct to `internal/init/` follows the project's clean architecture pattern where command definitions are business logic, not CLI framework concerns. This enables the init package to directly accept the `InitCmd` struct as an argument, improving code organization and reducing dependencies between layers. It also establishes a consistent pattern for command-related types across the codebase.

## What Changes

- Move `InitCmd` struct definition from `cmd/init.go` to `internal/init/models.go`
- Export `InitCmd` with proper doc comments in `internal/init/`
- Update `cmd/init.go` to import `InitCmd` from `internal/init/`
- Update `cmd/root.go` to import and reference `InitCmd` from `internal/init/`
- Refactor init package functions to accept `InitCmd` struct as parameter instead of individual fields
- Maintain 100% behavioral compatibility with existing CLI interface

## Impact

- **Affected specs**: cli-framework, cli-interface
- **Affected code**: cmd/init.go, cmd/root.go, internal/init/models.go, and internal/init functions
- **Breaking changes**: None (public CLI interface unchanged)
- **Migration**: All internal; no user-facing changes

## Related Work

This refactoring is part of the `refactor-init-package` change (currently 82/107 tasks). It enables cleaner separation of concerns and sets a pattern for organizing command types within the internal package structure.
