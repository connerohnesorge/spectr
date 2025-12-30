# Change: Add View Command for Dashboard Display

## Why

Spectr currently lacks a comprehensive dashboard that provides a quick overview
of the entire project state. Users must run multiple commands (`list`,
`validate`) or manually navigate directories to understand active changes, task
progress, and available specifications. A view command that aggregates this
information into a single, visual dashboard would dramatically improve project
visibility and developer productivity. OpenSpec provides such a dashboard—Spectr
should offer equivalent functionality to give users an at-a-glance understanding
of their spec-driven development workflow.

## What Changes

- Add new `view` command to the CLI that displays a comprehensive project
  dashboard
- Show summary metrics for specifications, active changes, and completed changes
- Display active changes with visual progress bars based on task completion
- Display completed changes (all tasks done or no tasks defined)
- Display specifications with requirement counts
- Add `--json` flag for machine-readable structured output
- Use colored, formatted terminal output with box-drawing characters for visual
  appeal
- Sort active changes by completion percentage (ascending) for priority
  visibility
- Sort specifications by requirement count (descending) to highlight complexity
- Provide helpful footer with navigation hints to related commands

## Impact

- **Affected specs**: `cli-framework` (extends existing)
- **Affected code**:
  - `cmd/root.go` - Add ViewCmd struct to CLI
  - New `cmd/view.go` - View command implementation
  - New `internal/view/` - View dashboard functionality package
    - `dashboard.go` - Core dashboard rendering logic
    - `types.go` - Data structures for dashboard data
    - `formatters.go` - Text and JSON output formatting
    - `progress.go` - Progress bar rendering utilities
  - Reuse `internal/discovery/` - Item discovery utilities (from list/validate
    commands)
  - Reuse `internal/parsers/` - Markdown parsing (from list/validate commands)
  - May need additional styling package or use lipgloss for colors and
    formatting

## Benefits

- **Project visibility**: Single command provides complete project overview
- **Progress tracking**: Visual progress bars show task completion at a glance
- **Developer productivity**: Reduces cognitive load by aggregating key
  information
- **Priority guidance**: Sorting by completion helps identify what needs
  attention
- **Consistency**: Matches OpenSpec's proven dashboard UX pattern
- **Automation-friendly**: JSON output enables integration with other tools
- **Visual clarity**: Colored output and formatting make information scannable
