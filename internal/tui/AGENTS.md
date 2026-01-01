# TUI Package

Bubble Tea components and interactive flows. Charmbracelet-based.

## OVERVIEW
Internal TUI utilities for interactive CLI commands. Uses Bubble Tea framework, Bubbles widgets, Lipgloss styling. Provides reusable components: prompts, menus, dashboards.

## STRUCTURE
```go

internal/tui/
├── menu.go              # Interactive menu selection
├── styles.go            # Lipgloss styles/constants
├── helpers.go           # TUI utility functions
└── *_test.go            # teatest-based tests
```

## WHERE TO LOOK

| Task | Location | Notes |
|------|----------|-------|
| Menu selection | menu.go | Bubble Tea model |
| Colors/styles | styles.go | Lipgloss theme |
| Helper utilities | helpers.go | Common patterns |

## CONVENTIONS
- **Bubble Tea**: Model-Update-View pattern
- **Lipgloss**: Use predefined styles, avoid inline styling
- **teatest**: Use charmbracelet/x/exp/teatest for tests
- **Cleanup models**: Always quit with tea.Quit

## KEY COMPONENTS
- **Menu**: Interactive selection list with keyboard navigation
- **Styles**: Colors, borders, spacing constants
- **Helpers**: Spinner functions, text wrapping, common patterns

## ANTI-PATTERNS
- **NO blocking I/O**: TUI must not block event loop
- **DON'T mix CLI and TUI**: Choose one mode per command
- **NO hardcoded colors**: Use styles.go constants

## COMMON PATTERNS
```go
// Create TUI program
p := tea.NewProgram(initialModel, tea.WithAltScreen())
if _, err := p.Run(); err != nil {
    log.Fatal(err)
}

// Menu selection
m := menu.NewModel(items, menu.Width(50))
selected := menu.Run(m)
```

## TESTING
- Use `teatest.NewProgram()` for component tests
- Verify final model state after key presses
- Test rendering with teatest.String()

## FLOW
1. Initialize Bubble Tea model
2. Define update() function for Msg handling
3. Define view() function for rendering
4. Run program, return selected value
