# CLI Framework Specification Delta

## ADDED Requirements

### Requirement: View Command Structure

The CLI SHALL provide a `view` command that displays a comprehensive project
dashboard with summary metrics, active changes, completed changes, and
specifications.

#### Scenario: View command registration

- **WHEN** the CLI is initialized
- **THEN** it SHALL include a ViewCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr view`
- **AND** help text SHALL describe dashboard functionality

#### Scenario: View command invocation

- **WHEN** user runs `spectr view` without flags
- **THEN** the system displays a dashboard with colored terminal output
- **AND** includes summary metrics section
- **AND** includes active changes section with progress bars
- **AND** includes completed changes section
- **AND** includes specifications section with requirement counts
- **AND** includes footer with navigation hints

#### Scenario: View command with JSON output

- **WHEN** user runs `spectr view --json`
- **THEN** the system outputs dashboard data as JSON
- **AND** includes summary, activeChanges, completedChanges, and specs fields
- **AND** SHALL be parseable by standard JSON tools

### Requirement: Dashboard Summary Metrics

The view command SHALL display summary metrics aggregating key project
statistics in a dedicated section at the top of the dashboard.

#### Scenario: Display summary with all metrics

- **WHEN** the dashboard is rendered
- **THEN** the summary section SHALL include total number of specifications
- **AND** SHALL include total number of requirements across all specs
- **AND** SHALL include number of active changes (in progress)
- **AND** SHALL include number of completed changes
- **AND** SHALL include total task count across all active changes
- **AND** SHALL include completed task count across all active changes

#### Scenario: Calculate total requirements

- **WHEN** aggregating specification requirements
- **THEN** the system SHALL sum requirement counts from all specs
- **AND** SHALL parse each spec.md file to count requirements
- **AND** SHALL handle specs with zero requirements gracefully

#### Scenario: Calculate task progress

- **WHEN** aggregating task progress
- **THEN** the system SHALL sum all tasks from all active changes
- **AND** SHALL count completed tasks (marked `[x]`)
- **AND** SHALL calculate overall percentage as `(completedTasks / totalTasks) *
  100`
- **AND** SHALL handle division by zero (display 0% if no tasks)

### Requirement: Active Changes Display

The view command SHALL display active changes with visual progress bars showing
task completion status.

#### Scenario: List active changes with progress

- **WHEN** the dashboard displays active changes
- **THEN** each change SHALL show its ID padded to 30 characters
- **AND** SHALL show a progress bar rendered with block characters
- **AND** SHALL show completion percentage after the progress bar
- **AND** SHALL use yellow circle indicator (◉) before each change
- **AND** SHALL sort changes by completion percentage ascending, then by ID
  alphabetically

#### Scenario: Render progress bar

- **WHEN** rendering a progress bar for a change
- **THEN** the bar SHALL have fixed width of 20 characters
- **AND** filled portion SHALL use full block character (█) in green
- **AND** empty portion SHALL use light block character (░) in dim gray
- **AND** filled width SHALL be `round((completed / total) * 20)`
- **AND** format SHALL be `[████████████░░░░░░░░]`

#### Scenario: Handle zero tasks

- **WHEN** a change has zero total tasks in tasks.md
- **THEN** the progress bar SHALL render as empty `[░░░░░░░░░░░░░░░░░░░░]`
- **AND** percentage SHALL display as `0%`
- **AND** the change SHALL still appear in active changes section

#### Scenario: No active changes

- **WHEN** no active changes exist (all completed or none exist)
- **THEN** the active changes section SHALL not be displayed
- **AND** the dashboard SHALL proceed to display other sections

### Requirement: Completed Changes Display

The view command SHALL display changes that have all tasks completed or no tasks
defined.

#### Scenario: List completed changes

- **WHEN** the dashboard displays completed changes
- **THEN** each change SHALL show its ID
- **AND** SHALL use green checkmark indicator (✓) before each change
- **AND** SHALL sort changes alphabetically by ID

#### Scenario: Determine completion status

- **WHEN** evaluating if a change is completed
- **THEN** a change is completed if tasks.md has all tasks marked `[x]`
- **OR** if tasks.md has zero total tasks
- **AND** changes with partial completion remain in active changes

#### Scenario: No completed changes

- **WHEN** no completed changes exist
- **THEN** the completed changes section SHALL not be displayed
- **AND** the dashboard SHALL proceed to display other sections

### Requirement: Specifications Display

The view command SHALL display all specifications sorted by requirement count to
highlight complexity.

#### Scenario: List specifications with requirement counts

- **WHEN** the dashboard displays specifications
- **THEN** each spec SHALL show its ID padded to 30 characters
- **AND** SHALL show requirement count with format `{count} requirement(s)`
- **AND** SHALL use blue square indicator (▪) before each spec
- **AND** SHALL sort specs by requirement count descending, then by ID
  alphabetically

#### Scenario: Pluralize requirement label

- **WHEN** displaying requirement count
- **THEN** use "requirement" for count of 1
- **AND** use "requirements" for count != 1

#### Scenario: No specifications found

- **WHEN** no specs exist in spectr/specs/
- **THEN** the specifications section SHALL not be displayed
- **AND** the dashboard SHALL complete without error

### Requirement: Dashboard Visual Formatting

The view command SHALL use colored output, Unicode box-drawing characters, and
consistent styling for visual clarity.

#### Scenario: Render dashboard header

- **WHEN** the dashboard is displayed
- **THEN** it SHALL start with bold title "Spectr Dashboard" (or similar)
- **AND** SHALL use double-line separator (═) below the title with width 60
- **AND** SHALL use consistent spacing between sections

#### Scenario: Render section headers

- **WHEN** displaying a section (Summary, Active Changes, etc.)
- **THEN** the section name SHALL be bold and cyan
- **AND** SHALL use single-line separator (─) below the header with width 60

#### Scenario: Render footer

- **WHEN** the dashboard completes rendering
- **THEN** it SHALL display a closing double-line separator (═) with width 60
- **AND** SHALL display a dim hint referencing related commands
- **AND** hint SHALL mention `spectr list --changes` and `spectr list --specs`

#### Scenario: Color scheme consistency

- **WHEN** applying colors to dashboard elements
- **THEN** use cyan for section headers
- **AND** use yellow for active change indicators
- **AND** use green for completed indicators and filled progress bars
- **AND** use blue for spec indicators
- **AND** use dim gray for empty progress bars and footer hints

### Requirement: JSON Output Format

The view command SHALL support `--json` flag to output dashboard data as
structured JSON for programmatic consumption.

#### Scenario: JSON structure

- **WHEN** user provides `--json` flag
- **THEN** output SHALL be a JSON object with top-level fields: `summary`,
  `activeChanges`, `completedChanges`, `specs`
- **AND** `summary` SHALL contain: `totalSpecs`, `totalRequirements`,
  `activeChanges`, `completedChanges`, `totalTasks`, `completedTasks`
- **AND** `activeChanges` SHALL be an array of objects with: `id`, `title`,
  `progress` (object with `total`, `completed`, `percentage`)
- **AND** `completedChanges` SHALL be an array of objects with: `id`, `title`
- **AND** `specs` SHALL be an array of objects with: `id`, `title`,
  `requirementCount`

#### Scenario: JSON arrays sorted consistently

- **WHEN** outputting JSON
- **THEN** `activeChanges` array SHALL be sorted by percentage ascending, then
  ID alphabetically
- **AND** `completedChanges` array SHALL be sorted by ID alphabetically
- **AND** `specs` array SHALL be sorted by requirementCount descending, then ID
  alphabetically

#### Scenario: JSON with no items

- **WHEN** outputting JSON and a category has no items
- **THEN** the corresponding array SHALL be empty `[]`
- **AND** summary counts SHALL reflect zero appropriately

### Requirement: Sorting Strategy

The view command SHALL sort dashboard items to surface the most relevant
information first.

#### Scenario: Sort active changes by priority

- **WHEN** sorting active changes
- **THEN** calculate completion percentage as `(completed / total) * 100`
- **AND** sort by percentage ascending (least complete first)
- **AND** for ties, sort alphabetically by ID

#### Scenario: Sort specs by complexity

- **WHEN** sorting specifications
- **THEN** sort by requirement count descending (most requirements first)
- **AND** for ties, sort alphabetically by ID

#### Scenario: Sort completed changes alphabetically

- **WHEN** sorting completed changes
- **THEN** sort by ID alphabetically

### Requirement: Data Reuse from Discovery and Parsers

The view command SHALL reuse existing discovery and parsing infrastructure to
avoid code duplication.

#### Scenario: Discover changes and specs

- **WHEN** building dashboard data
- **THEN** use `internal/discovery` package functions to find changes
- **AND** use `internal/discovery` package functions to find specs
- **AND** exclude archived changes from active/completed lists

#### Scenario: Parse titles and counts

- **WHEN** extracting metadata from markdown files
- **THEN** use `internal/parsers` package to parse proposal.md for titles
- **AND** use `internal/parsers` package to parse spec.md for titles and
  requirement counts
- **AND** use `internal/parsers` package to parse tasks.md for task counts

### Requirement: View Command Help Text

The view command SHALL provide comprehensive help documentation.

#### Scenario: Command help display

- **WHEN** user invokes `spectr view --help`
- **THEN** help text SHALL describe dashboard purpose
- **AND** SHALL list available flags (--json)
- **AND** SHALL indicate that no positional arguments are required
