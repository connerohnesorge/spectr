# Change: Implement Interactive Validation TUI and Migrate Helpers to Internal

## Why

The validation command currently returns a "not yet implemented" error when
invoked without arguments in interactive mode, despite the validation spec
defining this requirement. This creates a poor user experience and doesn't
leverage the project's existing TUI capabilities.

Additionally, the cmd/ package contains helper files (validate_helpers.go,
validate_items.go, validate_print.go) that violate the clean architecture
pattern. Per project conventions, cmd/ should be a thin layer delegating to
internal packages. These helpers contain substantial business logic that belongs
in internal/validation/ for better testability, reusability, and architectural
consistency.

## What Changes

- Implement interactive TUI for validation using bubbletea/lipgloss (following
  internal/list/interactive.go patterns)
- Add interactive mode that prompts users with options: All, All changes, All
  specs, Pick specific item
- Migrate validate_helpers.go functions to internal/validation/helpers.go
- Migrate validate_items.go functions to internal/validation/items.go
- Migrate validate_print.go functions to internal/validation/formatters.go
- Update cmd/validate.go imports to reference new internal locations
- Preserve all existing functionality and test coverage
- Add tests for new interactive TUI components

## Impact

- Affected specs: `validation`
- Affected code:
  - `cmd/validate.go` - Update imports, remove helper method receivers
  - `cmd/validate_helpers.go` - **REMOVED** (migrated to
    internal/validation/helpers.go)
  - `cmd/validate_items.go` - **REMOVED** (migrated to
    internal/validation/items.go)
  - `cmd/validate_print.go` - **REMOVED** (migrated to
    internal/validation/formatters.go)
  - `internal/validation/` - **NEW** files: helpers.go, items.go, formatters.go,
    interactive.go
  - `internal/validation/` - **NEW** tests: interactive_test.go
- Breaking changes: None - this is internal refactoring + new functionality
- User benefit: Interactive validation mode for improved UX, better code
  organization
