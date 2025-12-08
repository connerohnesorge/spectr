# Change: Add Responsive Table Column Trimming

## Why

The interactive TUI tables in `spectr list -I` and other list commands use fixed column widths that total ~95-108 characters. On smaller terminals (80 columns, split panes, or mobile SSH sessions), the table content overflows and wraps awkwardly, making it difficult to read and navigate.

## What Changes

- Detect terminal width at TUI initialization using `tea.WindowSizeMsg`
- Define column priority levels for each table view (ID is highest priority, less important columns like Tasks/Deltas are lowest)
- Progressively hide or narrow columns based on available width:
  - **Full width (110+)**: Show all columns at default widths
  - **Medium width (90-109)**: Reduce Title column width, truncate more aggressively
  - **Narrow width (70-89)**: Hide lowest-priority columns (Tasks/Details), narrow remaining columns
  - **Minimal width (<70)**: Show only ID and Title columns with aggressive truncation
- Handle dynamic terminal resize events during TUI session
- Apply responsive behavior consistently across all interactive table views (changes, specs, unified, archive)

## Impact

- Affected specs: cli-interface
- Affected code: internal/list/interactive.go, internal/tui/table.go, internal/tui/types.go
