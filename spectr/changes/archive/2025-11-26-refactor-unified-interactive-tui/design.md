# Design Document

## Context

The Spectr CLI uses Bubbletea (charmbracelet/bubbletea) for interactive TUI
components. Currently, there are two main interactive implementations:

1. **List Interactive** (`internal/list/interactive.go`): Provides interactive
  selection for changes, specs, and unified views. Features include clipboard
  copy, editor opening, archive triggering, and type filtering.

2. **Validation Interactive** (`internal/validation/interactive.go`): Provides a
  menu-driven validation workflow with options to validate all items, specific
  types, or pick individual items.

Both share similar patterns but are completely independent, leading to ~1100
lines of partially duplicated code.

## Goals / Non-Goals

**Goals:**

- Create a shared `internal/tui` package with composable components
- Reduce code duplication by 40-50%
- Maintain exact same user-facing behavior
- Improve consistency of styling and key bindings
- Make adding new interactive features easier

**Non-Goals:**

- Adding new features to the interactive modes
- Changing the visual appearance
- Supporting additional TUI libraries
- Creating a general-purpose TUI framework

## Decisions

### Decision: Create `internal/tui` package structure

The new package will contain:

```text
internal/tui/
├── styles.go      # Shared lipgloss styles, applyTableStyles
├── helpers.go     # truncateString, copyToClipboard
├── table.go       # TablePicker component for item selection
├── menu.go        # MenuPicker component for option selection
└── types.go       # Shared types and interfaces
```text

**Rationale:** This structure separates concerns and allows each consumer (list,
validation) to compose only what they need.

### Decision: TablePicker as primary building block

The `TablePicker` will be a configurable table-based selector that supports:

- Configurable columns
- Row data as `[]table.Row`
- Configurable key actions (map of key -> callback)
- Standard navigation (up/down, j/k)
- Standard quit (q, Ctrl+C)
- Help text generation from registered actions

**Alternatives considered:**

- Embedding bubbletea models directly - rejected as it still requires
  duplication of Update logic
- Using interfaces for shared behavior - rejected as too abstract for the
  concrete use cases

### Decision: Keep domain logic in consuming packages

The `list` and `validation` packages will remain responsible for:

- Building their specific data structures
- Defining domain-specific actions (archive, edit, validate)
- Handling domain-specific messages

The `tui` package only provides UI primitives.

**Rationale:** Keeps the TUI package focused and prevents coupling to business
logic.

### Decision: Action registration pattern

Actions will be registered via a builder pattern:

```go
picker := tui.NewTablePicker(columns, rows).
    WithProjectPath(projectPath).
    WithAction("e", "edit", editHandler).
    WithAction("a", "archive", archiveHandler).
    WithStandardNav().
    WithStandardQuit()
```text

**Rationale:** Allows each consumer to compose exactly the actions they need
without inheritance or conditionals.

## Risks / Trade-offs

**Risk:** Over-abstraction

- **Mitigation:** Start with only clearly shared code. If something is used in
  only one place, keep it there.

**Risk:** Breaking existing behavior

- **Mitigation:** Write comprehensive tests before refactoring. Run `go test
  ./...` after each change.

**Risk:** Increased complexity for simple changes

- **Mitigation:** Keep the API simple. If adding a feature requires touching
  `internal/tui`, that's fine - it should be easy.

## Migration Plan

1. Create `internal/tui` package with basic types and helpers (non-breaking)
2. Move `applyTableStyles` and `truncateString` to tui package, update imports
3. Implement `TablePicker` component
4. Refactor `list/interactive.go` to use `TablePicker`
5. Implement `MenuPicker` component
6. Refactor `validation/interactive.go` to use shared components
7. Remove orphaned code from original files
8. Update and add tests

Each step can be tested independently. Rollback is safe at any point.

## Open Questions

- Should `copyToClipboard` move to tui package or stay in list? (Currently
  leaning: move to tui as it's a UI concern)
- Should the tui package expose its own test helpers for consumers? (Currently
  leaning: yes, teatest patterns are complex)
