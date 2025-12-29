# Design: View Command Dashboard

## Context

The view command provides a comprehensive dashboard for Spectr projects, aggregating information from changes and specs into a single, visually appealing overview. This design is based on OpenSpec's proven view command implementation, adapted for Go and Spectr's architecture.

## Goals

- Provide at-a-glance project status with minimal cognitive load
- Display all key metrics (specs, changes, requirements, tasks) in one view
- Use visual elements (progress bars, colors, symbols) for quick comprehension
- Support both human-readable terminal output and machine-readable JSON
- Reuse existing discovery and parsing infrastructure from list/validate commands
- Maintain consistency with OpenSpec's dashboard layout and information hierarchy

## Non-Goals

- Interactive TUI with navigation (this is explicitly static output)
- Real-time updates or live progress tracking
- Detailed validation or error reporting (use validate command for that)
- Filtering or searching within the dashboard view

## Decisions

### Decision: Use lipgloss for Styling

**Rationale**: Spectr already uses `charmbracelet/lipgloss` for the init wizard. Leveraging it for the view command ensures consistency and avoids adding new dependencies.

**Alternatives considered**:

- Plain ANSI escape codes: More control but harder to maintain, less readable
- Custom styling package: Unnecessary when lipgloss is already available

### Decision: Reuse Discovery and Parsing from Other Commands

**Rationale**: The list and validate commands already implement change/spec discovery and markdown parsing. Rather than duplicating code, the view command will depend on shared `internal/discovery` and `internal/parsers` packages.

**Alternatives considered**:

- Duplicate logic in view package: Violates DRY, increases maintenance burden
- Inline all logic in view command: Makes code harder to test and reuse

### Decision: Static Output with Box-Drawing Characters

**Rationale**: Following OpenSpec's pattern, the dashboard uses static colored output with Unicode box-drawing characters (═, ─, etc.) for visual structure. This is simple, portable, and works in all modern terminals.

**Alternatives considered**:

- Interactive TUI with bubbletea: User explicitly requested static output
- ASCII-only output: Less visually appealing, doesn't leverage modern terminal capabilities

### Decision: JSON Output for Automation

**Rationale**: Adding a `--json` flag enables scripting, CI/CD integration, and programmatic consumption of dashboard data. The JSON structure mirrors the visual output's information hierarchy.

**JSON Schema**:

```json
{
  "summary": {
    "totalSpecs": 5,
    "totalRequirements": 42,
    "activeChanges": 3,
    "completedChanges": 2,
    "totalTasks": 15,
    "completedTasks": 8
  },
  "activeChanges": [
    {
      "id": "add-view-command",
      "title": "Add View Command for Dashboard Display",
      "progress": {
        "total": 8,
        "completed": 3,
        "percentage": 37
      }
    }
  ],
  "completedChanges": [
    {
      "id": "add-list-command",
      "title": "Add List Command"
    }
  ],
  "specs": [
    {
      "id": "cli-framework",
      "title": "CLI Framework",
      "requirementCount": 15
    }
  ]
}
```

### Decision: Progress Bar Rendering Algorithm

**Rationale**: Visual progress bars provide immediate intuition about task completion. The algorithm uses filled (█) and empty (░) Unicode block characters with a fixed width of 20 characters.

**Implementation**:

- Calculate percentage: `completed / total`
- Fill width: `round(percentage * 20)`
- Render: `[████████████░░░░░░░░]`
- Color filled portion green, empty portion dim gray
- Show percentage as `37%` after the bar

**Edge case**: If total is 0, render empty bar `[░░░░░░░░░░░░░░░░░░░░]` with dim styling

### Decision: Sorting Strategy

**Rationale**: The dashboard sorts items to surface the most relevant information:

1. **Active changes**: Sort by completion percentage (ascending), then by ID alphabetically
   - **Why**: Changes with lower completion need more attention and should appear first
   - **Tie-breaker**: Alphabetical by ID ensures deterministic ordering

2. **Specs**: Sort by requirement count (descending), then by ID alphabetically
   - **Why**: Specs with more requirements represent larger capabilities and are more complex
   - **Tie-breaker**: Alphabetical by ID ensures deterministic ordering

3. **Completed changes**: Sort alphabetically by ID
   - **Why**: These are done; alphabetical ordering is sufficient for browsing

## Data Flow

```
┌─────────────────┐
│  View Command   │
│   (cmd/view.go) │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────┐
│  Dashboard Builder                  │
│  (internal/view/dashboard.go)       │
│                                     │
│  1. Discover changes via discovery  │
│  2. Discover specs via discovery    │
│  3. Parse metadata via parsers      │
│  4. Calculate progress metrics      │
│  5. Sort items per strategy         │
│  6. Format output (text or JSON)    │
└────────┬────────────────────────────┘
         │
         ├──────────────┬──────────────┐
         ▼              ▼              ▼
   ┌──────────┐  ┌──────────┐  ┌──────────┐
   │Discovery │  │ Parsers  │  │Formatters│
   └──────────┘  └──────────┘  └──────────┘
```

## Output Format Specification

### Human-Readable Terminal Output

```
OpenSpec Dashboard

════════════════════════════════════════════════════════════

Summary:
  ● Specifications: 5 specs, 42 requirements
  ● Active Changes: 3 in progress
  ● Completed Changes: 2
  ● Task Progress: 8/15 (53% complete)

Active Changes
────────────────────────────────────────────────────────────
  ◉ add-view-command              [████████░░░░░░░░░░░░] 37%
  ◉ add-validate-command          [████████████░░░░░░░░] 60%
  ◉ migrate-to-kong-cli           [████████████████████] 100%

Completed Changes
────────────────────────────────────────────────────────────
  ✓ add-list-command

Specifications
────────────────────────────────────────────────────────────
  ▪ cli-framework                 15 requirements
  ▪ validation                    12 requirements

════════════════════════════════════════════════════════════

Use spectr list --changes or spectr list --specs for detailed views
```

**Color Scheme** (using lipgloss):

- Section headers: Bold, cyan
- Summary labels: Cyan bullet (●)
- Active change indicator: Yellow circle (◉)
- Completed change indicator: Green checkmark (✓)
- Spec indicator: Blue square (▪)
- Progress bar filled: Green
- Progress bar empty: Dim gray
- Percentage: Dim
- Separator lines: Default terminal color
- Footer hint: Dim

## Dependencies

**Go Packages**:

- `github.com/charmbracelet/lipgloss` (already in go.mod) - Styling and colors
- Standard library `os`, `path/filepath`, `encoding/json`, `sort`

**Internal Packages** (to be created or reused):

- `internal/discovery` - Change and spec discovery (shared with list/validate)
- `internal/parsers` - Markdown parsing for titles, tasks, requirements (shared with list/validate)
- `internal/view` - New package for dashboard-specific logic

## Risks & Trade-offs

### Risk: Terminal Compatibility

**Mitigation**: Use widely-supported Unicode box-drawing characters. Test on common terminals (iTerm2, Terminal.app, Windows Terminal, GNOME Terminal). Provide `NO_COLOR` environment variable support.

### Risk: Performance with Large Projects

**Mitigation**: Discovery and parsing are inherently I/O-bound. For projects with hundreds of changes/specs, add progress spinner during data collection phase.

### Trade-off: Static vs Interactive

**Decision**: Static output means no navigation or drill-down. Users must run `show` or `list` commands to see details. This is acceptable because the view command is explicitly a dashboard overview, not a detailed explorer.

### Trade-off: Completion Detection Logic

**Question**: How to determine if a change is "completed"?
**Decision**: A change is completed if all tasks in `tasks.md` are marked `[x]` OR if `tasks.md` has zero tasks (edge case: change with proposal but no tasks yet). This matches OpenSpec's behavior and is simple to implement.

## Migration Plan

N/A - This is a new feature with no breaking changes.

## Open Questions

None - All design decisions have been made based on user requirements and OpenSpec reference implementation.
