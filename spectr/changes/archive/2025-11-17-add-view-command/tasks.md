# Implementation Tasks

## 1. Foundation

- [x] 1.1 Create `internal/view/` package structure
- [x] 1.2 Define data types in `internal/view/types.go` for DashboardData,
  SummaryMetrics, ChangeProgress, etc.
- [x] 1.3 Verify `internal/discovery/` package exists (from list/validate
  commands)
- [x] 1.4 Verify `internal/parsers/` package exists (from list/validate
  commands)

## 2. Dashboard Data Collection

- [x] 2.1 Implement `dashboard.CollectData()` to gather all dashboard
  information
- [x] 2.2 Use `discovery.GetActiveChanges()` to find all changes
- [x] 2.3 Use `discovery.GetSpecs()` to find all specifications
- [x] 2.4 Parse each change's `proposal.md` for title using
  `parsers.ExtractTitle()`
- [x] 2.5 Parse each change's `tasks.md` for task counts using
  `parsers.CountTasks()`
- [x] 2.6 Parse each spec's `spec.md` for title and requirement count using
  parsers
- [x] 2.7 Categorize changes into active vs completed based on task completion
- [x] 2.8 Calculate summary metrics (total specs, total requirements, total
  tasks, etc.)

## 3. Sorting Logic

- [x] 3.1 Implement sorting for active changes by completion percentage
  (ascending), then ID (alphabetical)
- [x] 3.2 Implement sorting for completed changes by ID (alphabetical)
- [x] 3.3 Implement sorting for specs by requirement count (descending), then ID
  (alphabetical)
- [x] 3.4 Add utility functions for percentage calculation with zero-division
  handling

## 4. Progress Bar Rendering

- [x] 4.1 Implement `progress.RenderBar()` to create progress bar strings
- [x] 4.2 Use fixed width of 20 characters for progress bars
- [x] 4.3 Calculate filled width as `round((completed / total) * 20)`
- [x] 4.4 Use filled block character (█) for completed portion
- [x] 4.5 Use empty block character (░) for remaining portion
- [x] 4.6 Apply green color to filled portion using lipgloss
- [x] 4.7 Apply dim gray to empty portion using lipgloss
- [x] 4.8 Handle edge case of zero total tasks (render empty bar)
- [x] 4.9 Calculate and format percentage display

## 5. Text Output Formatting

- [x] 5.1 Implement `formatters.FormatDashboardText()` for human-readable
  terminal output
- [x] 5.2 Create dashboard header with title "Spectr Dashboard" and double-line
  separator (═)
- [x] 5.3 Format summary section with bullet points and metrics
- [x] 5.4 Format active changes section with headers, progress bars, and
  percentages
- [x] 5.5 Format completed changes section with checkmarks
- [x] 5.6 Format specifications section with requirement counts
- [x] 5.7 Add single-line separators (─) between sections
- [x] 5.8 Create footer with double-line separator and navigation hints
- [x] 5.9 Apply color scheme using lipgloss (cyan headers, yellow/green/blue
  indicators)
- [x] 5.10 Implement conditional section display (hide empty sections)

## 6. JSON Output Formatting

- [x] 6.1 Implement `formatters.FormatDashboardJSON()` for machine-readable
  output
- [x] 6.2 Create JSON structure with `summary`, `activeChanges`,
  `completedChanges`, `specs` fields
- [x] 6.3 Format summary object with all metrics (totalSpecs, totalRequirements,
  etc.)
- [x] 6.4 Format activeChanges array with id, title, progress object
- [x] 6.5 Format completedChanges array with id, title
- [x] 6.6 Format specs array with id, title, requirementCount
- [x] 6.7 Ensure arrays are sorted consistently with text output
- [x] 6.8 Use proper JSON encoding with indentation for readability

## 7. Command Integration

- [x] 7.1 Add `ViewCmd` struct to `cmd/root.go` with Kong tags
- [x] 7.2 Create `cmd/view.go` with `ViewCmd.Run()` method
- [x] 7.3 Add `--json` flag to ViewCmd for JSON output
- [x] 7.4 Implement command execution flow (collect data → format → output)
- [x] 7.5 Add error handling for missing spectr directory
- [x] 7.6 Add error handling for discovery/parsing failures
- [x] 7.7 Set appropriate exit codes for errors

## 8. Testing

- [x] 8.1 Write unit tests for `dashboard.CollectData()` with mock data
- [x] 8.2 Write unit tests for sorting functions
- [x] 8.3 Write unit tests for `progress.RenderBar()` with various completion
  levels
- [x] 8.4 Write unit tests for percentage calculation including edge cases
- [x] 8.5 Write unit tests for JSON formatting with various data structures
- [x] 8.6 Write integration test for `spectr view` (default text output)
- [x] 8.7 Write integration test for `spectr view --json`
- [x] 8.8 Test with empty project (no changes, no specs)
- [x] 8.9 Test with only active changes, no completed changes
- [x] 8.10 Test with only completed changes, no active changes
- [x] 8.11 Test with changes having zero tasks
- [x] 8.12 Test color output respects NO_COLOR environment variable

## 9. Documentation and Polish

- [x] 9.1 Add comprehensive help text via Kong struct tags
- [x] 9.2 Test help output with `spectr view --help`
- [x] 9.3 Verify terminal output matches design specification (colors, symbols,
  spacing)
- [x] 9.4 Verify JSON output schema matches design specification
- [x] 9.5 Test on multiple terminal emulators (iTerm2, GNOME Terminal, Windows
  Terminal)
- [x] 9.6 Ensure Unicode box-drawing characters render correctly
- [x] 9.7 Add comments and documentation to all public functions
- [x] 9.8 Verify consistency with OpenSpec reference behavior
