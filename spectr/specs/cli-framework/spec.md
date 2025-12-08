# Cli Framework Specification

## Purpose

This specification defines the CLI framework structure using Kong for declarative command definitions with struct tags, supporting subcommands (archive, list, validate, view), flags, positional arguments, automatic method dispatch, and built-in help generation.

## Requirements

### Requirement: Archive Command
The CLI SHALL provide an `archive` command that moves completed changes to a dated archive directory and applies delta specifications to main specs.

#### Scenario: Archive with change ID
- **WHEN** user runs `spectr archive <change-id>`
- **THEN** the system archives the specified change without prompting

#### Scenario: Interactive archive selection
- **WHEN** user runs `spectr archive` without specifying a change ID
- **THEN** the system displays a list of active changes and prompts for selection

#### Scenario: Non-interactive archiving with yes flag
- **WHEN** user runs `spectr archive <change-id> --yes`
- **THEN** the system archives without any confirmation prompts

#### Scenario: Skip spec updates for tooling changes
- **WHEN** user runs `spectr archive <change-id> --skip-specs`
- **THEN** the system archives the change without updating main specs

#### Scenario: Skip validation with confirmation
- **WHEN** user runs `spectr archive <change-id> --no-validate`
- **THEN** the system warns about skipping validation and requires confirmation unless --yes flag is also provided

### Requirement: Archive Command Flags
The archive command SHALL support flags for controlling behavior.

#### Scenario: Yes flag skips all prompts
- **WHEN** user provides the `-y` or `--yes` flag
- **THEN** the system skips all confirmation prompts for automated usage

#### Scenario: Skip specs flag bypasses spec updates
- **WHEN** user provides the `--skip-specs` flag
- **THEN** the system moves the change to archive without applying delta specs

#### Scenario: No validate flag skips validation
- **WHEN** user provides the `--no-validate` flag
- **THEN** the system skips validation but requires confirmation unless --yes is also provided

### Requirement: Struct-Based Command Definition
The CLI framework SHALL use Go struct types with struct tags to declaratively define command structure, subcommands, flags, and arguments. Provider configuration SHALL be retrieved from the `Registry` interface rather than static global maps.

#### Scenario: Root command definition
- **WHEN** the CLI is initialized
- **THEN** it SHALL use a root struct with subcommand fields tagged with `cmd` for command definitions
- **AND** each subcommand SHALL be a nested struct type with appropriate tags

#### Scenario: Subcommand registration
- **WHEN** a new subcommand is added to the CLI
- **THEN** it SHALL be defined as a struct field on the parent command struct
- **AND** it SHALL use `cmd` tag to indicate it is a subcommand
- **AND** it SHALL include a `help` tag describing the command purpose

#### Scenario: Tool configuration lookup
- **WHEN** the executor needs tool configuration
- **THEN** it SHALL query the `Registry` via `Get(id)` method
- **AND** it SHALL NOT use hardcoded global maps

### Requirement: Declarative Flag Definition
The CLI framework SHALL define flags using struct fields with Kong struct tags instead of imperative flag registration.

#### Scenario: String flag definition
- **WHEN** a command requires a string flag
- **THEN** it SHALL be defined as a struct field with `name` tag for the flag name
- **AND** it MAY include `short` tag for single-character shorthand
- **AND** it SHALL include `help` tag describing the flag purpose
- **AND** it MAY include `default` tag for default values

#### Scenario: Boolean flag definition
- **WHEN** a command requires a boolean flag
- **THEN** it SHALL be defined as a bool struct field with appropriate tags
- **AND** the flag SHALL default to false unless explicitly set

#### Scenario: Slice flag definition
- **WHEN** a command requires a multi-value flag
- **THEN** it SHALL be defined as a slice type struct field
- **AND** it SHALL support comma-separated values or repeated flag usage

### Requirement: Positional Argument Support
The CLI framework SHALL support positional arguments using struct fields tagged with `arg`.

#### Scenario: Optional positional argument
- **WHEN** a command accepts an optional positional argument
- **THEN** it SHALL be defined with `arg` and `optional` tags
- **AND** the field SHALL be a pointer type or have a zero value for "not provided"

#### Scenario: Required positional argument
- **WHEN** a command requires a positional argument
- **THEN** it SHALL be defined with `arg` tag without `optional`
- **AND** parsing SHALL fail if the argument is not provided

### Requirement: Automatic Method Dispatch
The CLI framework SHALL automatically invoke the appropriate command's Run method after parsing.

#### Scenario: Command execution
- **WHEN** a command is successfully parsed
- **THEN** the framework SHALL call the command struct's `Run() error` method
- **AND** it SHALL pass any configured context values to the Run method
- **AND** it SHALL handle the returned error appropriately

### Requirement: Built-in Help Generation
The CLI framework SHALL automatically generate help text from struct tags and types.

#### Scenario: Root help display
- **WHEN** the CLI is invoked with `--help` or no arguments
- **THEN** it SHALL display a list of available subcommands
- **AND** it SHALL show descriptions from `help` tags
- **AND** it SHALL indicate required vs optional arguments

#### Scenario: Subcommand help display
- **WHEN** a subcommand is invoked with `--help`
- **THEN** it SHALL display the command description
- **AND** it SHALL list all flags with their types and help text
- **AND** it SHALL show positional argument requirements

### Requirement: Error Handling and Exit Codes
The CLI framework SHALL provide appropriate error messages and exit codes for parsing and execution failures.

#### Scenario: Parse error handling
- **WHEN** invalid flags or arguments are provided
- **THEN** it SHALL display an error message
- **AND** it SHALL show usage information
- **AND** it SHALL exit with non-zero status code

#### Scenario: Execution error handling
- **WHEN** a command's Run method returns an error
- **THEN** it SHALL display the error message
- **AND** it SHALL exit with non-zero status code

### Requirement: Backward-Compatible CLI Interface
The CLI framework SHALL maintain the same command syntax and flag names as the previous implementation.

#### Scenario: Init command compatibility
- **WHEN** users invoke `spectr init` with existing flag combinations
- **THEN** the behavior SHALL be identical to the previous Cobra-based implementation
- **AND** all flag names SHALL remain unchanged
- **AND** short flag aliases SHALL remain unchanged
- **AND** positional argument handling SHALL remain unchanged

#### Scenario: Help text accessibility
- **WHEN** users invoke `spectr --help` or `spectr init --help`
- **THEN** help information SHALL be displayed (format may differ from Cobra)
- **AND** all commands and flags SHALL be documented

### Requirement: List Command for Changes
The system SHALL provide a `list` command that enumerates all active changes in the project, displaying their IDs by default.

#### Scenario: List changes with IDs only
- **WHEN** user runs `spectr list` without flags
- **THEN** the system displays change IDs, one per line, sorted alphabetically
- **AND** excludes archived changes in the `archive/` directory

#### Scenario: List changes with details
- **WHEN** user runs `spectr list --long`
- **THEN** the system displays each change with format: `{id}: {title} [deltas {count}] [tasks {completed}/{total}]`
- **AND** sorts output alphabetically by ID

#### Scenario: List changes as JSON
- **WHEN** user runs `spectr list --json`
- **THEN** the system outputs a JSON array of objects with fields: `id`, `title`, `deltaCount`, `taskStatus` (with `total` and `completed`)
- **AND** sorts the array by ID

#### Scenario: No changes found
- **WHEN** user runs `spectr list` and no active changes exist
- **THEN** the system displays "No items found"

### Requirement: List Command for Specs
The system SHALL support a `--specs` flag that switches the list command to enumerate specifications instead of changes.

#### Scenario: List specs with IDs only
- **WHEN** user runs `spectr list --specs` without other flags
- **THEN** the system displays spec IDs, one per line, sorted alphabetically
- **AND** only includes directories with valid `spec.md` files

#### Scenario: List specs with details
- **WHEN** user runs `spectr list --specs --long`
- **THEN** the system displays each spec with format: `{id}: {title} [requirements {count}]`
- **AND** sorts output alphabetically by ID

#### Scenario: List specs as JSON
- **WHEN** user runs `spectr list --specs --json`
- **THEN** the system outputs a JSON array of objects with fields: `id`, `title`, `requirementCount`
- **AND** sorts the array by ID

#### Scenario: No specs found
- **WHEN** user runs `spectr list --specs` and no specs exist
- **THEN** the system displays "No items found"

### Requirement: Change Discovery
The system SHALL discover active changes by scanning the `spectr/changes/` directory and identifying subdirectories that contain a `proposal.md` file, excluding the `archive/` directory.

#### Scenario: Find active changes
- **WHEN** the system scans for changes
- **THEN** it includes all subdirectories of `spectr/changes/` that contain `proposal.md`
- **AND** excludes the `spectr/changes/archive/` directory and its contents
- **AND** excludes hidden directories (starting with `.`)

### Requirement: Spec Discovery
The system SHALL discover specs by scanning the `spectr/specs/` directory and identifying subdirectories that contain a `spec.md` file.

#### Scenario: Find specs
- **WHEN** the system scans for specs
- **THEN** it includes all subdirectories of `spectr/specs/` that contain `spec.md`
- **AND** excludes hidden directories (starting with `.`)

### Requirement: Title Extraction
The system SHALL extract titles from proposal and spec markdown files by finding the first level-1 heading and removing the "Change:" or "Spec:" prefix if present.

#### Scenario: Extract title from proposal
- **WHEN** the system reads a `proposal.md` file with heading `# Change: Add Feature`
- **THEN** it extracts the title as "Add Feature"

#### Scenario: Extract title from spec
- **WHEN** the system reads a `spec.md` file with heading `# CLI Framework`
- **THEN** it extracts the title as "CLI Framework"

#### Scenario: Fallback to ID when title not found
- **WHEN** the system cannot extract a title from a markdown file
- **THEN** it uses the directory name (ID) as the title

### Requirement: Task Counting
The system SHALL count tasks in `tasks.md` files by identifying lines matching the pattern `- [ ]` or `- [x]` (case-insensitive), with completed tasks marked by `[x]`. The system SHALL prefer `tasks.json` when present.

#### Scenario: Count tasks from JSON
- **WHEN** the system counts tasks and `tasks.json` exists
- **THEN** it reads task status from the JSON file
- **AND** counts tasks by status field values
- **AND** reports `taskStatus` with total and completed counts

#### Scenario: Count completed and total tasks
- **WHEN** the system reads a `tasks.md` file with 3 tasks, 2 marked `[x]` and 1 marked `[ ]`
- **THEN** it reports `taskStatus` as `{ total: 3, completed: 2 }`

#### Scenario: Handle missing tasks file
- **WHEN** the system cannot find or read a `tasks.md` or `tasks.json` file for a change
- **THEN** it reports `taskStatus` as `{ total: 0, completed: 0 }`
- **AND** continues processing without error

#### Scenario: JSON takes precedence over Markdown
- **WHEN** both `tasks.json` and `tasks.md` exist (should not happen normally)
- **THEN** the system reads from `tasks.json`
- **AND** ignores `tasks.md`

### Requirement: Validate Command Structure
The CLI SHALL provide a validate command for checking spec and change document correctness.

#### Scenario: Validate command registration
- **WHEN** the CLI is initialized
- **THEN** it SHALL include a ValidateCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr validate`
- **AND** help text SHALL describe validation functionality

#### Scenario: Direct item validation invocation
- **WHEN** user invokes `spectr validate <item-name>`
- **THEN** the command SHALL validate the named item (change or spec)
- **AND** SHALL print validation results to stdout
- **AND** SHALL exit with code 0 for valid, 1 for invalid

#### Scenario: Bulk validation invocation
- **WHEN** user invokes `spectr validate --all`
- **THEN** the command SHALL validate all changes and specs
- **AND** SHALL print summary of results
- **AND** SHALL display full issue details for each failed item including level, path, and message
- **AND** SHALL exit with code 1 if any item fails validation

#### Scenario: Interactive validation invocation
- **WHEN** user invokes `spectr validate` without arguments in a TTY
- **THEN** the command SHALL prompt for what to validate
- **AND** SHALL execute the user's selection

### Requirement: Validate Command Flags
The validate command SHALL support flags for controlling validation behavior and output format.

#### Scenario: Strict mode flag
- **WHEN** user provides `--strict` flag
- **THEN** validation SHALL treat warnings as errors
- **AND** exit code SHALL be 1 if warnings exist
- **AND** validation report SHALL reflect strict mode in valid field

#### Scenario: JSON output flag
- **WHEN** user provides `--json` flag
- **THEN** output SHALL be formatted as JSON
- **AND** SHALL include items, summary, and version fields
- **AND** SHALL be parseable by standard JSON tools

#### Scenario: Type disambiguation flag
- **WHEN** user provides `--type change` or `--type spec`
- **THEN** the command SHALL treat the item as the specified type
- **AND** SHALL skip type auto-detection
- **AND** SHALL error if item does not exist as that type

#### Scenario: All items flag
- **WHEN** user provides `--all` flag
- **THEN** the command SHALL validate all changes and all specs
- **AND** SHALL run in bulk validation mode

#### Scenario: Changes only flag
- **WHEN** user provides `--changes` flag
- **THEN** the command SHALL validate all changes only
- **AND** SHALL skip specs

#### Scenario: Specs only flag
- **WHEN** user provides `--specs` flag
- **THEN** the command SHALL validate all specs only
- **AND** SHALL skip changes

#### Scenario: Non-interactive flag
- **WHEN** user provides `--no-interactive` flag
- **THEN** the command SHALL not prompt for input
- **AND** SHALL print usage hint if no item specified
- **AND** SHALL exit with code 1

### Requirement: Validate Command Help Text
The validate command SHALL provide comprehensive help documentation.

#### Scenario: Command help display
- **WHEN** user invokes `spectr validate --help`
- **THEN** help text SHALL describe validation purpose
- **AND** SHALL list all available flags with descriptions
- **AND** SHALL show usage examples for common scenarios
- **AND** SHALL indicate optional vs required arguments

### Requirement: Positional Argument Support for Item Name
The validate command SHALL accept an optional positional argument for the item to validate.

#### Scenario: Optional item name argument
- **WHEN** validate command is defined
- **THEN** it SHALL have an ItemName field tagged with `arg:"" optional:""`
- **AND** the field type SHALL be pointer to string or string with zero value check
- **AND** omitting the argument SHALL be valid (triggers interactive or bulk mode)

#### Scenario: Item name provided
- **WHEN** user provides item name as positional argument
- **THEN** the command SHALL validate that specific item
- **AND** SHALL auto-detect whether it's a change or spec
- **AND** SHALL respect --type flag if provided for disambiguation

### Requirement: View Command Structure
The CLI SHALL provide a `view` command that displays a comprehensive project dashboard with summary metrics, active changes, completed changes, and specifications.

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
The view command SHALL display summary metrics aggregating key project statistics in a dedicated section at the top of the dashboard.

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
- **AND** SHALL calculate overall percentage as `(completedTasks / totalTasks) * 100`
- **AND** SHALL handle division by zero (display 0% if no tasks)

### Requirement: Active Changes Display
The view command SHALL display active changes with visual progress bars showing task completion status.

#### Scenario: List active changes with progress
- **WHEN** the dashboard displays active changes
- **THEN** each change SHALL show its ID padded to 30 characters
- **AND** SHALL show a progress bar rendered with block characters
- **AND** SHALL show completion percentage after the progress bar
- **AND** SHALL use yellow circle indicator (◉) before each change
- **AND** SHALL sort changes by completion percentage ascending, then by ID alphabetically

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
The view command SHALL display changes that have all tasks completed or no tasks defined.

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
The view command SHALL display all specifications sorted by requirement count to highlight complexity.

#### Scenario: List specifications with requirement counts
- **WHEN** the dashboard displays specifications
- **THEN** each spec SHALL show its ID padded to 30 characters
- **AND** SHALL show requirement count with format `{count} requirement(s)`
- **AND** SHALL use blue square indicator (▪) before each spec
- **AND** SHALL sort specs by requirement count descending, then by ID alphabetically

#### Scenario: Pluralize requirement label
- **WHEN** displaying requirement count
- **THEN** use "requirement" for count of 1
- **AND** use "requirements" for count != 1

#### Scenario: No specifications found
- **WHEN** no specs exist in spectr/specs/
- **THEN** the specifications section SHALL not be displayed
- **AND** the dashboard SHALL complete without error

### Requirement: Dashboard Visual Formatting
The view command SHALL use colored output, Unicode box-drawing characters, and consistent styling for visual clarity.

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
The view command SHALL support `--json` flag to output dashboard data as structured JSON for programmatic consumption.

#### Scenario: JSON structure
- **WHEN** user provides `--json` flag
- **THEN** output SHALL be a JSON object with top-level fields: `summary`, `activeChanges`, `completedChanges`, `specs`
- **AND** `summary` SHALL contain: `totalSpecs`, `totalRequirements`, `activeChanges`, `completedChanges`, `totalTasks`, `completedTasks`
- **AND** `activeChanges` SHALL be an array of objects with: `id`, `title`, `progress` (object with `total`, `completed`, `percentage`)
- **AND** `completedChanges` SHALL be an array of objects with: `id`, `title`
- **AND** `specs` SHALL be an array of objects with: `id`, `title`, `requirementCount`

#### Scenario: JSON arrays sorted consistently
- **WHEN** outputting JSON
- **THEN** `activeChanges` array SHALL be sorted by percentage ascending, then ID alphabetically
- **AND** `completedChanges` array SHALL be sorted by ID alphabetically
- **AND** `specs` array SHALL be sorted by requirementCount descending, then ID alphabetically

#### Scenario: JSON with no items
- **WHEN** outputting JSON and a category has no items
- **THEN** the corresponding array SHALL be empty `[]`
- **AND** summary counts SHALL reflect zero appropriately

### Requirement: Sorting Strategy
The view command SHALL sort dashboard items to surface the most relevant information first.

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
The view command SHALL reuse existing discovery and parsing infrastructure to avoid code duplication.

#### Scenario: Discover changes and specs
- **WHEN** building dashboard data
- **THEN** use `internal/discovery` package functions to find changes
- **AND** use `internal/discovery` package functions to find specs
- **AND** exclude archived changes from active/completed lists

#### Scenario: Parse titles and counts
- **WHEN** extracting metadata from markdown files
- **THEN** use `internal/parsers` package to parse proposal.md for titles
- **AND** use `internal/parsers` package to parse spec.md for titles and requirement counts
- **AND** use `internal/parsers` package to parse tasks.md for task counts

### Requirement: View Command Help Text
The view command SHALL provide comprehensive help documentation.

#### Scenario: Command help display
- **WHEN** user invokes `spectr view --help`
- **THEN** help text SHALL describe dashboard purpose
- **AND** SHALL list available flags (--json)
- **AND** SHALL indicate that no positional arguments are required

### Requirement: Provider Interface
The init system SHALL define a `Provider` interface that all AI CLI tool integrations implement, with one provider per tool handling both instruction files and slash commands.

#### Scenario: Provider interface methods
- **WHEN** a new provider is created
- **THEN** it SHALL implement `ID() string` returning a unique kebab-case identifier
- **AND** it SHALL implement `Name() string` returning the human-readable name
- **AND** it SHALL implement `Priority() int` returning display order
- **AND** it SHALL implement `ConfigFile() string` returning instruction file path or empty string
- **AND** it SHALL implement `GetProposalCommandPath() string` returning relative path for proposal command or empty string
- **AND** it SHALL implement `GetApplyCommandPath() string` returning relative path for apply command or empty string
- **AND** it SHALL implement `CommandFormat() CommandFormat` returning Markdown or TOML
- **AND** it SHALL implement `Configure(projectPath, spectrDir string) error` for configuration
- **AND** it SHALL implement `IsConfigured(projectPath string) bool` for status checks

#### Scenario: Single provider per tool
- **WHEN** a tool has both an instruction file and slash commands
- **THEN** one provider SHALL handle both (e.g., ClaudeProvider handles CLAUDE.md and .claude/commands/)
- **AND** there SHALL NOT be separate config and slash providers for the same tool

#### Scenario: Flexible command paths
- **WHEN** a provider returns paths from command path methods
- **THEN** each method SHALL return a relative path including directory and filename
- **AND** paths MAY have different directories for each command type
- **AND** paths MAY have different file extensions based on CommandFormat
- **AND** empty string indicates the provider does not support that command

#### Scenario: HasSlashCommands detection
- **WHEN** code calls `HasSlashCommands()` on a provider
- **THEN** it SHALL return true if ANY command path method returns a non-empty string
- **AND** it SHALL return false only if ALL command path methods return empty strings

### Requirement: Provider Registry
The init system SHALL provide a `Registry` that manages registration and lookup of providers using a registry pattern similar to `database/sql`.

#### Scenario: Register provider
- **WHEN** a provider calls `Register(provider Provider)`
- **THEN** the registry SHALL store the provider by its ID
- **AND** duplicate registration SHALL panic with a descriptive message

#### Scenario: Get provider by ID
- **WHEN** code calls `Get(id string) (Provider, bool)`
- **THEN** the registry SHALL return the provider and true if found
- **AND** SHALL return nil and false if not found

#### Scenario: List all providers
- **WHEN** code calls `All() []Provider`
- **THEN** the registry SHALL return all registered providers
- **AND** providers SHALL be sorted by Priority ascending

### Requirement: Per-Provider File Organization
The init system SHALL organize provider implementations as separate Go files under `internal/initialize/providers/`, with one file per provider.

#### Scenario: Provider file structure
- **WHEN** a provider file is created
- **THEN** it SHALL be named `{provider-id}.go` (e.g., `claude.go`, `gemini.go`)
- **AND** it SHALL contain an `init()` function that registers its provider
- **AND** it SHALL be self-contained with all provider-specific configuration

#### Scenario: Adding a new provider
- **WHEN** a developer adds a new AI CLI provider
- **THEN** they SHALL create a single file under `internal/initialize/providers/`
- **AND** the file SHALL implement the `Provider` interface
- **AND** the file SHALL call `Register()` in its `init()` function
- **AND** no other files SHALL require modification

### Requirement: Init Function Registration
The init system SHALL use Go's `init()` function pattern for automatic provider registration at startup.

#### Scenario: Auto-registration at startup
- **WHEN** the program starts
- **THEN** all provider `init()` functions SHALL execute before `main()`
- **AND** all providers SHALL be registered in the global registry
- **AND** registration order SHALL not affect functionality

### Requirement: Command Format Support
The init system SHALL support multiple command file formats through the `CommandFormat` type.

#### Scenario: Markdown command format
- **WHEN** a provider returns `FormatMarkdown` from `CommandFormat()`
- **THEN** slash command files SHALL be generated as `.md` files
- **AND** files SHALL use frontmatter and spectr markers

#### Scenario: TOML command format
- **WHEN** a provider returns `FormatTOML` from `CommandFormat()`
- **THEN** slash command files SHALL be generated as `.toml` files
- **AND** the TOML SHALL include `description` field with command description
- **AND** the TOML SHALL include `prompt` field with the command prompt content

### Requirement: Version Command Structure
The CLI SHALL provide a `version` command that displays version information including version number, git commit hash, and build date.

#### Scenario: Version command registration
- **WHEN** the CLI is initialized
- **THEN** it SHALL include a VersionCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr version`
- **AND** help text SHALL describe version display functionality

#### Scenario: Version command invocation
- **WHEN** user runs `spectr version` without flags
- **THEN** the system displays version in format: `spectr version {version} (commit: {commit}, built: {date})`
- **AND** version SHALL be the semantic version (e.g., `0.1.0` or `dev`)
- **AND** commit SHALL be the git commit hash (short or full) or `unknown`
- **AND** date SHALL be the build date in ISO 8601 format or `unknown`

#### Scenario: Version command with short flag
- **WHEN** user runs `spectr version --short`
- **THEN** the system displays only the version number (e.g., `0.1.0`)
- **AND** no other information is displayed

#### Scenario: Version command with JSON flag
- **WHEN** user runs `spectr version --json`
- **THEN** the system outputs version data as JSON
- **AND** JSON SHALL include fields: `version`, `commit`, `date`
- **AND** SHALL be parseable by standard JSON tools

### Requirement: Version Variable Injection
The version information SHALL be injectable at build time via Go ldflags, supporting both goreleaser releases and nix flake builds.

#### Scenario: Goreleaser version injection
- **WHEN** goreleaser builds the binary
- **THEN** version SHALL be set from git tag via ldflags
- **AND** commit SHALL be set from git commit hash via ldflags
- **AND** date SHALL be set from build timestamp via ldflags

#### Scenario: Nix flake version injection
- **WHEN** nix builds the binary via flake.nix
- **THEN** version SHALL be set from the flake package version attribute via ldflags
- **AND** commit and date MAY be `unknown` if not available in nix build context

#### Scenario: Development build defaults
- **WHEN** binary is built without ldflags (e.g., `go build`)
- **THEN** version SHALL default to `dev`
- **AND** commit SHALL default to `unknown`
- **AND** date SHALL default to `unknown`

### Requirement: Version Package Location
The version variables SHALL be defined in a dedicated `internal/version` package for clean separation and easy ldflags targeting.

#### Scenario: Package structure
- **WHEN** the version package is imported
- **THEN** it SHALL expose `Version`, `Commit`, and `Date` string variables
- **AND** variables SHALL have default values for development builds
- **AND** the ldflags path SHALL be `github.com/connerohnesorge/spectr/internal/version`

### Requirement: Completion Command Structure
The CLI SHALL provide a `completion` subcommand that outputs shell completion scripts for supported shells using the kong-completion library.

#### Scenario: Completion command registration
- **WHEN** the CLI is initialized
- **THEN** it SHALL include a Completion field using `kongcompletion.Completion` type
- **AND** the command SHALL be accessible via `spectr completion`
- **AND** help text SHALL describe shell completion functionality

#### Scenario: Bash completion output
- **WHEN** user runs `spectr completion bash`
- **THEN** the system outputs a valid bash completion script
- **AND** the script can be sourced directly or added to bash-completion.d

#### Scenario: Zsh completion output
- **WHEN** user runs `spectr completion zsh`
- **THEN** the system outputs a valid zsh completion script
- **AND** the script can be sourced or placed in $fpath

#### Scenario: Fish completion output
- **WHEN** user runs `spectr completion fish`
- **THEN** the system outputs a valid fish completion script
- **AND** the script can be saved to fish completions directory

### Requirement: Custom Predictors for Dynamic Arguments
The completion system SHALL provide context-aware suggestions for arguments that accept dynamic values like change IDs or spec IDs.

#### Scenario: Change ID completion
- **WHEN** user types `spectr archive <TAB>` or `spectr validate <TAB>`
- **AND** the argument expects a change ID
- **THEN** completion suggests all active change IDs from `spectr/changes/`
- **AND** excludes archived changes

#### Scenario: Spec ID completion
- **WHEN** user types `spectr validate --type spec <TAB>`
- **AND** the argument expects a spec ID
- **THEN** completion suggests all spec IDs from `spectr/specs/`

#### Scenario: Item type completion
- **WHEN** user types `spectr validate --type <TAB>`
- **THEN** completion suggests `change` and `spec`

### Requirement: Kong-Completion Integration Pattern
The CLI initialization SHALL follow the kong-completion pattern where Kong is initialized, completions are registered, and then arguments are parsed.

#### Scenario: Initialization order
- **WHEN** the program starts
- **THEN** `kong.Must()` is called first to create the Kong application
- **AND** `kongcompletion.Register()` is called before parsing
- **AND** `app.Parse()` is called after completion registration
- **AND** this order ensures completions work correctly

#### Scenario: Predictor registration
- **WHEN** custom predictors are defined
- **THEN** they SHALL be registered via `kongcompletion.WithPredictor()`
- **AND** struct fields SHALL reference predictors using `predictor:"name"` tag

### Requirement: Accept Command Structure
The CLI SHALL provide an `accept` command that converts `tasks.md` to `tasks.json` format for stable agent manipulation during implementation.

#### Scenario: Accept command registration
- **WHEN** the CLI is initialized
- **THEN** it SHALL include an AcceptCmd struct field tagged with `cmd`
- **AND** the command SHALL be accessible via `spectr accept`
- **AND** help text SHALL describe task format conversion functionality

#### Scenario: Accept with change ID
- **WHEN** user runs `spectr accept <change-id>`
- **THEN** the system validates the change exists in `spectr/changes/<change-id>/`
- **AND** the system parses `tasks.md` into structured format
- **AND** the system writes `tasks.json` with proper schema
- **AND** the system removes `tasks.md` to prevent drift

#### Scenario: Accept with validation
- **WHEN** user runs `spectr accept <change-id>`
- **THEN** the system validates the change before conversion
- **AND** the system blocks acceptance if validation fails
- **AND** the system displays validation errors

#### Scenario: Accept dry-run mode
- **WHEN** user runs `spectr accept <change-id> --dry-run`
- **THEN** the system displays what would be converted
- **AND** the system does NOT write tasks.json
- **AND** the system does NOT remove tasks.md

#### Scenario: Accept already accepted change
- **WHEN** user runs `spectr accept <change-id>` on a change that already has tasks.json
- **THEN** the system displays a message indicating change is already accepted
- **AND** the system exits with code 0 (success, idempotent)

#### Scenario: Accept change without tasks.md
- **WHEN** user runs `spectr accept <change-id>` on a change without tasks.md
- **THEN** the system displays an error indicating no tasks.md found
- **AND** the system exits with code 1

### Requirement: Tasks JSON Schema
The accept command SHALL generate `tasks.json` files conforming to a versioned schema with structured task objects.

#### Scenario: JSON file structure
- **WHEN** the accept command creates tasks.json
- **THEN** the file SHALL contain a root object with `version` and `tasks` fields
- **AND** `version` SHALL be integer 1 for this schema version
- **AND** `tasks` SHALL be an array of task objects

#### Scenario: Task object structure
- **WHEN** a task is serialized to JSON
- **THEN** it SHALL have `id` field containing the task identifier (e.g., "1.1")
- **AND** it SHALL have `section` field containing the section header (e.g., "Implementation")
- **AND** it SHALL have `description` field containing the full task text
- **AND** it SHALL have `status` field with value "pending", "in_progress", or "completed"

#### Scenario: Status value mapping from Markdown
- **WHEN** converting tasks.md to tasks.json
- **THEN** `- [ ]` SHALL map to status "pending"
- **AND** `- [x]` (case-insensitive) SHALL map to status "completed"

### Requirement: Accept Command Flags
The accept command SHALL support flags for controlling behavior.

#### Scenario: Dry-run flag
- **WHEN** user provides the `--dry-run` flag
- **THEN** the system previews the conversion without writing files
- **AND** displays the JSON that would be generated

#### Scenario: Interactive change selection
- **WHEN** user runs `spectr accept` without specifying a change ID
- **THEN** the system displays a list of active changes with tasks.md files
- **AND** prompts for selection using existing TUI components
