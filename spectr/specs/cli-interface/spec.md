# Cli Interface Specification

## Purpose

This specification defines the CLI framework structure using Kong for declarative command definitions with struct tags, supporting subcommands (archive, list, validate, view), flags, positional arguments, automatic method dispatch, and built-in help generation.
Additionally This specification defines interactive CLI features including navigable table interfaces for list and archive commands, cross-platform clipboard operations, initialization wizard flows, and visual styling for enhanced user experience.

## Requirements

### Requirement: Archive Command
The CLI SHALL provide an `archive` command that moves completed changes to a dated archive directory and applies delta specifications to main specs.

#### Scenario: Archive with change ID
- WHEN user runs `spectr archive <change-id>`
- THEN the system archives the specified change without prompting

#### Scenario: Interactive archive selection
- WHEN user runs `spectr archive` without specifying a change ID
- THEN the system displays a list of active changes and prompts for selection

#### Scenario: Non-interactive archiving with yes flag
- WHEN user runs `spectr archive <change-id> --yes`
- THEN the system archives without any confirmation prompts

#### Scenario: Skip spec updates for tooling changes
- WHEN user runs `spectr archive <change-id> --skip-specs`
- THEN the system archives the change without updating main specs

#### Scenario: Skip validation with confirmation
- WHEN user runs `spectr archive <change-id> --no-validate`
- THEN the system warns about skipping validation and requires confirmation unless --yes flag is also provided

### Requirement: Archive Command Flags
The archive command SHALL support flags for controlling behavior.

#### Scenario: Yes flag skips all prompts
- WHEN user provides the `-y` or `--yes` flag
- THEN the system skips all confirmation prompts for automated usage

#### Scenario: Skip specs flag bypasses spec updates
- WHEN user provides the `--skip-specs` flag
- THEN the system moves the change to archive without applying delta specs

#### Scenario: No validate flag skips validation
- WHEN user provides the `--no-validate` flag
- THEN the system skips validation but requires confirmation unless --yes is also provided

### Requirement: Struct-Based Command Definition
The CLI framework SHALL use Go struct types with struct tags to declaratively define command structure, subcommands, flags, and arguments. Provider configuration SHALL be retrieved from the `Registry` interface rather than static global maps.

#### Scenario: Root command definition
- WHEN the CLI is initialized
- THEN it SHALL use a root struct with subcommand fields tagged with `cmd` for command definitions
- AND each subcommand SHALL be a nested struct type with appropriate tags

#### Scenario: Subcommand registration
- WHEN a new subcommand is added to the CLI
- THEN it SHALL be defined as a struct field on the parent command struct
- AND it SHALL use `cmd` tag to indicate it is a subcommand
- AND it SHALL include a `help` tag describing the command purpose

#### Scenario: Tool configuration lookup
- WHEN the executor needs tool configuration
- THEN it SHALL query the `Registry` via `Get(id)` method
- AND it SHALL NOT use hardcoded global maps

### Requirement: Declarative Flag Definition
The CLI framework SHALL define flags using struct fields with Kong struct tags instead of imperative flag registration.

#### Scenario: String flag definition
- WHEN a command requires a string flag
- THEN it SHALL be defined as a struct field with `name` tag for the flag name
- AND it MAY include `short` tag for single-character shorthand
- AND it SHALL include `help` tag describing the flag purpose
- AND it MAY include `default` tag for default values

#### Scenario: Boolean flag definition
- WHEN a command requires a boolean flag
- THEN it SHALL be defined as a bool struct field with appropriate tags
- AND the flag SHALL default to false unless explicitly set

#### Scenario: Slice flag definition
- WHEN a command requires a multi-value flag
- THEN it SHALL be defined as a slice type struct field
- AND it SHALL support comma-separated values or repeated flag usage

### Requirement: Positional Argument Support
The CLI framework SHALL support positional arguments using struct fields tagged with `arg`.

#### Scenario: Optional positional argument
- WHEN a command accepts an optional positional argument
- THEN it SHALL be defined with `arg` and `optional` tags
- AND the field SHALL be a pointer type or have a zero value for "not provided"

#### Scenario: Required positional argument
- WHEN a command requires a positional argument
- THEN it SHALL be defined with `arg` tag without `optional`
- AND parsing SHALL fail if the argument is not provided

### Requirement: Automatic Method Dispatch
The CLI framework SHALL automatically invoke the appropriate command's Run method after parsing.

#### Scenario: Command execution
- WHEN a command is successfully parsed
- THEN the framework SHALL call the command struct's `Run() error` method
- AND it SHALL pass any configured context values to the Run method
- AND it SHALL handle the returned error appropriately

### Requirement: Built-in Help Generation
The CLI framework SHALL automatically generate help text from struct tags and types.

#### Scenario: Root help display
- WHEN the CLI is invoked with `--help` or no arguments
- THEN it SHALL display a list of available subcommands
- AND it SHALL show descriptions from `help` tags
- AND it SHALL indicate required vs optional arguments

#### Scenario: Subcommand help display
- WHEN a subcommand is invoked with `--help`
- THEN it SHALL display the command description
- AND it SHALL list all flags with their types and help text
- AND it SHALL show positional argument requirements

### Requirement: Error Handling and Exit Codes
The CLI framework SHALL provide appropriate error messages and exit codes for parsing and execution failures.

#### Scenario: Parse error handling
- WHEN invalid flags or arguments are provided
- THEN it SHALL display an error message
- AND it SHALL show usage information
- AND it SHALL exit with non-zero status code

#### Scenario: Execution error handling
- WHEN a command's Run method returns an error
- THEN it SHALL display the error message
- AND it SHALL exit with non-zero status code

### Requirement: Backward-Compatible CLI Interface
The CLI framework SHALL maintain the same command syntax and flag names as the previous implementation.

#### Scenario: Init command compatibility
- WHEN users invoke `spectr init` with existing flag combinations
- THEN the behavior SHALL be identical to the previous Cobra-based implementation
- AND all flag names SHALL remain unchanged
- AND short flag aliases SHALL remain unchanged
- AND positional argument handling SHALL remain unchanged

#### Scenario: Help text accessibility
- WHEN users invoke `spectr --help` or `spectr init --help`
- THEN help information SHALL be displayed (format may differ from Cobra)
- AND all commands and flags SHALL be documented

### Requirement: List Command for Changes
The system SHALL provide a `list` command that enumerates all active changes in the project, displaying their IDs by default.

#### Scenario: List changes with IDs only
- WHEN user runs `spectr list` without flags
- THEN the system displays change IDs, one per line, sorted alphabetically
- AND excludes archived changes in the `archive/` directory

#### Scenario: List changes with details
- WHEN user runs `spectr list --long`
- THEN the system displays each change with format: `{id}: {title} [deltas {count}] [tasks {completed}/{total}]`
- AND sorts output alphabetically by ID

#### Scenario: List changes as JSON
- WHEN user runs `spectr list --json`
- THEN the system outputs a JSON array of objects with fields: `id`, `title`, `deltaCount`, `taskStatus` (with `total` and `completed`)
- AND sorts the array by ID

#### Scenario: No changes found
- WHEN user runs `spectr list` and no active changes exist
- THEN the system displays "No items found"

### Requirement: List Command for Specs
The system SHALL support a `--specs` flag that switches the list command to enumerate specifications instead of changes.

#### Scenario: List specs with IDs only
- WHEN user runs `spectr list --specs` without other flags
- THEN the system displays spec IDs, one per line, sorted alphabetically
- AND only includes directories with valid `spec.md` files

#### Scenario: List specs with details
- WHEN user runs `spectr list --specs --long`
- THEN the system displays each spec with format: `{id}: {title} [requirements {count}]`
- AND sorts output alphabetically by ID

#### Scenario: List specs as JSON
- WHEN user runs `spectr list --specs --json`
- THEN the system outputs a JSON array of objects with fields: `id`, `title`, `requirementCount`
- AND sorts the array by ID

#### Scenario: No specs found
- WHEN user runs `spectr list --specs` and no specs exist
- THEN the system displays "No items found"

### Requirement: Change Discovery
The system SHALL discover active changes by scanning the `spectr/changes/` directory and identifying subdirectories that contain a `proposal.md` file, excluding the `archive/` directory.

#### Scenario: Find active changes
- WHEN the system scans for changes
- THEN it includes all subdirectories of `spectr/changes/` that contain `proposal.md`
- AND excludes the `spectr/changes/archive/` directory and its contents
- AND excludes hidden directories (starting with `.`)

### Requirement: Spec Discovery
The system SHALL discover specs by scanning the `spectr/specs/` directory and identifying subdirectories that contain a `spec.md` file.

#### Scenario: Find specs
- WHEN the system scans for specs
- THEN it includes all subdirectories of `spectr/specs/` that contain `spec.md`
- AND excludes hidden directories (starting with `.`)

### Requirement: Title Extraction
The system SHALL extract titles from proposal and spec markdown files by finding the first level-1 heading and removing the "Change:" or "Spec:" prefix if present.

#### Scenario: Extract title from proposal
- WHEN the system reads a `proposal.md` file with heading `# Change: Add Feature`
- THEN it extracts the title as "Add Feature"

#### Scenario: Extract title from spec
- WHEN the system reads a `spec.md` file with heading `# CLI Framework`
- THEN it extracts the title as "CLI Framework"

#### Scenario: Fallback to ID when title not found
- WHEN the system cannot extract a title from a markdown file
- THEN it uses the directory name (ID) as the title

### Requirement: Task Counting
The system SHALL count tasks from `tasks.jsonc` or `tasks.md` files. Legacy `tasks.json` files are silently ignored (breaking change).

#### Scenario: Count tasks from JSONC
- WHEN the system counts tasks and `tasks.jsonc` exists
- THEN it reads task status from the JSONC file
- AND strips any comments before parsing
- AND counts tasks by status field values
- AND reports `taskStatus` with total and completed counts

#### Scenario: Count tasks from Markdown
- WHEN the system counts tasks and `tasks.jsonc` does not exist but `tasks.md` exists
- THEN it identifies lines matching `- [ ]` or `- [x]` (case-insensitive)
- AND counts completed tasks by `[x]` markers
- AND reports `taskStatus` with total and completed counts

#### Scenario: Ignore legacy tasks.json
- WHEN the system counts tasks and only `tasks.json` exists (no `tasks.jsonc` or `tasks.md`)
- THEN it reports `taskStatus` as `{ total: 0, completed: 0 }`
- AND does NOT read the legacy `tasks.json` file

#### Scenario: Handle missing tasks file
- WHEN the system cannot find `tasks.jsonc` or `tasks.md` for a change
- THEN it reports `taskStatus` as `{ total: 0, completed: 0 }`
- AND continues processing without error

#### Scenario: JSONC takes precedence over Markdown
- WHEN both `tasks.jsonc` and `tasks.md` exist
- THEN the system reads from `tasks.jsonc`
- AND ignores `tasks.md`

### Requirement: Validate Command Structure
The CLI SHALL provide a validate command for checking spec and change document correctness.

#### Scenario: Validate command registration
- WHEN the CLI is initialized
- THEN it SHALL include a ValidateCmd struct field tagged with `cmd`
- AND the command SHALL be accessible via `spectr validate`
- AND help text SHALL describe validation functionality

#### Scenario: Direct item validation invocation
- WHEN user invokes `spectr validate <item-name>`
- THEN the command SHALL validate the named item (change or spec)
- AND SHALL print validation results to stdout
- AND SHALL exit with code 0 for valid, 1 for invalid

#### Scenario: Bulk validation invocation
- WHEN user invokes `spectr validate --all`
- THEN the command SHALL validate all changes and specs
- AND SHALL print summary of results
- AND SHALL display full issue details for each failed item including level, path, and message
- AND SHALL exit with code 1 if any item fails validation

#### Scenario: Interactive validation invocation
- WHEN user invokes `spectr validate` without arguments in a TTY
- THEN the command SHALL prompt for what to validate
- AND SHALL execute the user's selection

### Requirement: Validate Command Flags
The validate command SHALL support flags for controlling validation behavior and output format. Validation always treats warnings as errors.

#### Scenario: Default validation behavior (always strict)
- WHEN user runs `spectr validate <item>` without any strict flag
- THEN validation SHALL treat warnings as errors
- AND exit code SHALL be 1 if warnings or errors exist
- AND validation report SHALL show valid=false for any issues

#### Scenario: JSON output flag
- WHEN user provides `--json` flag
- THEN output SHALL be formatted as JSON
- AND SHALL include items, summary, and version fields
- AND SHALL be parseable by standard JSON tools

#### Scenario: Type disambiguation flag
- WHEN user provides `--type change` or `--type spec`
- THEN the command SHALL treat the item as the specified type
- AND SHALL skip type auto-detection
- AND SHALL error if item does not exist as that type

#### Scenario: All items flag
- WHEN user provides `--all` flag
- THEN the command SHALL validate all changes and all specs
- AND SHALL run in bulk validation mode

#### Scenario: Changes only flag
- WHEN user provides `--changes` flag
- THEN the command SHALL validate all changes only
- AND SHALL skip specs

#### Scenario: Specs only flag
- WHEN user provides `--specs` flag
- THEN the command SHALL validate all specs only
- AND SHALL skip changes

#### Scenario: Non-interactive flag
- WHEN user provides `--no-interactive` flag
- THEN the command SHALL not prompt for input
- AND SHALL print usage hint if no item specified
- AND SHALL exit with code 1

### Requirement: Validate Command Help Text
The validate command SHALL provide comprehensive help documentation.

#### Scenario: Command help display
- WHEN user invokes `spectr validate --help`
- THEN help text SHALL describe validation purpose
- AND SHALL list all available flags with descriptions
- AND SHALL show usage examples for common scenarios
- AND SHALL indicate optional vs required arguments

### Requirement: Positional Argument Support for Item Name
The validate command SHALL accept an optional positional argument for the item to validate.

#### Scenario: Optional item name argument
- WHEN validate command is defined
- THEN it SHALL have an ItemName field tagged with `arg:"" optional:""`
- AND the field type SHALL be pointer to string or string with zero value check
- AND omitting the argument SHALL be valid (triggers interactive or bulk mode)

#### Scenario: Item name provided
- WHEN user provides item name as positional argument
- THEN the command SHALL validate that specific item
- AND SHALL auto-detect whether it's a change or spec
- AND SHALL respect --type flag if provided for disambiguation

### Requirement: View Command Structure
The CLI SHALL provide a `view` command that displays a comprehensive project dashboard with summary metrics, active changes, completed changes, and specifications.

#### Scenario: View command registration
- WHEN the CLI is initialized
- THEN it SHALL include a ViewCmd struct field tagged with `cmd`
- AND the command SHALL be accessible via `spectr view`
- AND help text SHALL describe dashboard functionality

#### Scenario: View command invocation
- WHEN user runs `spectr view` without flags
- THEN the system displays a dashboard with colored terminal output
- AND includes summary metrics section
- AND includes active changes section with progress bars
- AND includes completed changes section
- AND includes specifications section with requirement counts
- AND includes footer with navigation hints

#### Scenario: View command with JSON output
- WHEN user runs `spectr view --json`
- THEN the system outputs dashboard data as JSON
- AND includes summary, activeChanges, completedChanges, and specs fields
- AND SHALL be parseable by standard JSON tools

### Requirement: Dashboard Summary Metrics
The view command SHALL display summary metrics aggregating key project statistics in a dedicated section at the top of the dashboard.

#### Scenario: Display summary with all metrics
- WHEN the dashboard is rendered
- THEN the summary section SHALL include total number of specifications
- AND SHALL include total number of requirements across all specs
- AND SHALL include number of active changes (in progress)
- AND SHALL include number of completed changes
- AND SHALL include total task count across all active changes
- AND SHALL include completed task count across all active changes

#### Scenario: Calculate total requirements
- WHEN aggregating specification requirements
- THEN the system SHALL sum requirement counts from all specs
- AND SHALL parse each spec.md file to count requirements
- AND SHALL handle specs with zero requirements gracefully

#### Scenario: Calculate task progress
- WHEN aggregating task progress
- THEN the system SHALL sum all tasks from all active changes
- AND SHALL count completed tasks (marked `[x]`)
- AND SHALL calculate overall percentage as `(completedTasks / totalTasks) * 100`
- AND SHALL handle division by zero (display 0% if no tasks)

### Requirement: Active Changes Display
The view command SHALL display active changes with visual progress bars showing task completion status.

#### Scenario: List active changes with progress
- WHEN the dashboard displays active changes
- THEN each change SHALL show its ID padded to 30 characters
- AND SHALL show a progress bar rendered with block characters
- AND SHALL show completion percentage after the progress bar
- AND SHALL use yellow circle indicator (◉) before each change
- AND SHALL sort changes by completion percentage ascending, then by ID alphabetically

#### Scenario: Render progress bar
- WHEN rendering a progress bar for a change
- THEN the bar SHALL have fixed width of 20 characters
- AND filled portion SHALL use full block character (█) in green
- AND empty portion SHALL use light block character (░) in dim gray
- AND filled width SHALL be `round((completed / total) * 20)`
- AND format SHALL be `[████████████░░░░░░░░]`

#### Scenario: Handle zero tasks
- WHEN a change has zero total tasks in tasks.md
- THEN the progress bar SHALL render as empty `[░░░░░░░░░░░░░░░░░░░░]`
- AND percentage SHALL display as `0%`
- AND the change SHALL still appear in active changes section

#### Scenario: No active changes
- WHEN no active changes exist (all completed or none exist)
- THEN the active changes section SHALL not be displayed
- AND the dashboard SHALL proceed to display other sections

### Requirement: Completed Changes Display
The view command SHALL display changes that have all tasks completed or no tasks defined.

#### Scenario: List completed changes
- WHEN the dashboard displays completed changes
- THEN each change SHALL show its ID
- AND SHALL use green checkmark indicator (✓) before each change
- AND SHALL sort changes alphabetically by ID

#### Scenario: Determine completion status
- WHEN evaluating if a change is completed
- THEN a change is completed if tasks.md has all tasks marked `[x]`
- OR if tasks.md has zero total tasks
- AND changes with partial completion remain in active changes

#### Scenario: No completed changes
- WHEN no completed changes exist
- THEN the completed changes section SHALL not be displayed
- AND the dashboard SHALL proceed to display other sections

### Requirement: Specifications Display
The view command SHALL display all specifications sorted by requirement count to highlight complexity.

#### Scenario: List specifications with requirement counts
- WHEN the dashboard displays specifications
- THEN each spec SHALL show its ID padded to 30 characters
- AND SHALL show requirement count with format `{count} requirement(s)`
- AND SHALL use blue square indicator (▪) before each spec
- AND SHALL sort specs by requirement count descending, then by ID alphabetically

#### Scenario: Pluralize requirement label
- WHEN displaying requirement count
- THEN use "requirement" for count of 1
- AND use "requirements" for count != 1

#### Scenario: No specifications found
- WHEN no specs exist in spectr/specs/
- THEN the specifications section SHALL not be displayed
- AND the dashboard SHALL complete without error

### Requirement: Dashboard Visual Formatting
The view command SHALL use colored output, Unicode box-drawing characters, and consistent styling for visual clarity.

#### Scenario: Render dashboard header
- WHEN the dashboard is displayed
- THEN it SHALL start with bold title "Spectr Dashboard" (or similar)
- AND SHALL use double-line separator (═) below the title with width 60
- AND SHALL use consistent spacing between sections

#### Scenario: Render section headers
- WHEN displaying a section (Summary, Active Changes, etc.)
- THEN the section name SHALL be bold and cyan
- AND SHALL use single-line separator (─) below the header with width 60

#### Scenario: Render footer
- WHEN the dashboard completes rendering
- THEN it SHALL display a closing double-line separator (═) with width 60
- AND SHALL display a dim hint referencing related commands
- AND hint SHALL mention `spectr list --changes` and `spectr list --specs`

#### Scenario: Color scheme consistency
- WHEN applying colors to dashboard elements
- THEN use cyan for section headers
- AND use yellow for active change indicators
- AND use green for completed indicators and filled progress bars
- AND use blue for spec indicators
- AND use dim gray for empty progress bars and footer hints

### Requirement: JSON Output Format
The view command SHALL support `--json` flag to output dashboard data as structured JSON for programmatic consumption.

#### Scenario: JSON structure
- WHEN user provides `--json` flag
- THEN output SHALL be a JSON object with top-level fields: `summary`, `activeChanges`, `completedChanges`, `specs`
- AND `summary` SHALL contain: `totalSpecs`, `totalRequirements`, `activeChanges`, `completedChanges`, `totalTasks`, `completedTasks`
- AND `activeChanges` SHALL be an array of objects with: `id`, `title`, `progress` (object with `total`, `completed`, `percentage`)
- AND `completedChanges` SHALL be an array of objects with: `id`, `title`
- AND `specs` SHALL be an array of objects with: `id`, `title`, `requirementCount`

#### Scenario: JSON arrays sorted consistently
- WHEN outputting JSON
- THEN `activeChanges` array SHALL be sorted by percentage ascending, then ID alphabetically
- AND `completedChanges` array SHALL be sorted by ID alphabetically
- AND `specs` array SHALL be sorted by requirementCount descending, then ID alphabetically

#### Scenario: JSON with no items
- WHEN outputting JSON and a category has no items
- THEN the corresponding array SHALL be empty `[]`
- AND summary counts SHALL reflect zero appropriately

### Requirement: Sorting Strategy
The view command SHALL sort dashboard items to surface the most relevant information first.

#### Scenario: Sort active changes by priority
- WHEN sorting active changes
- THEN calculate completion percentage as `(completed / total) * 100`
- AND sort by percentage ascending (least complete first)
- AND for ties, sort alphabetically by ID

#### Scenario: Sort specs by complexity
- WHEN sorting specifications
- THEN sort by requirement count descending (most requirements first)
- AND for ties, sort alphabetically by ID

#### Scenario: Sort completed changes alphabetically
- WHEN sorting completed changes
- THEN sort by ID alphabetically

### Requirement: Data Reuse from Discovery and Parsers
The view command SHALL reuse existing discovery and parsing infrastructure to avoid code duplication.

#### Scenario: Discover changes and specs
- WHEN building dashboard data
- THEN use `internal/discovery` package functions to find changes
- AND use `internal/discovery` package functions to find specs
- AND exclude archived changes from active/completed lists

#### Scenario: Parse titles and counts
- WHEN extracting metadata from markdown files
- THEN use `internal/parsers` package to parse proposal.md for titles
- AND use `internal/parsers` package to parse spec.md for titles and requirement counts
- AND use `internal/parsers` package to parse tasks.md for task counts

### Requirement: View Command Help Text
The view command SHALL provide comprehensive help documentation.

#### Scenario: Command help display
- WHEN user invokes `spectr view --help`
- THEN help text SHALL describe dashboard purpose
- AND SHALL list available flags (--json)
- AND SHALL indicate that no positional arguments are required

### Requirement: Provider Interface
The init system SHALL define a `Provider` interface that all AI CLI tool integrations implement, with one provider per tool handling both instruction files and slash commands.

#### Scenario: Provider interface methods
- WHEN a new provider is created
- THEN it SHALL implement `ID() string` returning a unique kebab-case identifier
- AND it SHALL implement `Name() string` returning the human-readable name
- AND it SHALL implement `Priority() int` returning display order
- AND it SHALL implement `ConfigFile() string` returning instruction file path or empty string
- AND it SHALL implement `GetProposalCommandPath() string` returning relative path for proposal command or empty string
- AND it SHALL implement `GetApplyCommandPath() string` returning relative path for apply command or empty string
- AND it SHALL implement `CommandFormat() CommandFormat` returning Markdown or TOML
- AND it SHALL implement `Configure(projectPath, spectrDir string) error` for configuration
- AND it SHALL implement `IsConfigured(projectPath string) bool` for status checks

#### Scenario: Single provider per tool
- WHEN a tool has both an instruction file and slash commands
- THEN one provider SHALL handle both (e.g., ClaudeProvider handles CLAUDE.md and .claude/commands/)
- AND there SHALL NOT be separate config and slash providers for the same tool

#### Scenario: Flexible command paths
- WHEN a provider returns paths from command path methods
- THEN each method SHALL return a relative path including directory and filename
- AND paths MAY have different directories for each command type
- AND paths MAY have different file extensions based on CommandFormat
- AND empty string indicates the provider does not support that command

#### Scenario: HasSlashCommands detection
- WHEN code calls `HasSlashCommands()` on a provider
- THEN it SHALL return true if ANY command path method returns a non-empty string
- AND it SHALL return false only if ALL command path methods return empty strings

### Requirement: Provider Registry
The init system SHALL provide a `Registry` that manages registration and lookup of providers using a registry pattern similar to `database/sql`.

#### Scenario: Register provider
- WHEN a provider calls `Register(provider Provider)`
- THEN the registry SHALL store the provider by its ID
- AND duplicate registration SHALL panic with a descriptive message

#### Scenario: Get provider by ID
- WHEN code calls `Get(id string) (Provider, bool)`
- THEN the registry SHALL return the provider and true if found
- AND SHALL return nil and false if not found

#### Scenario: List all providers
- WHEN code calls `All() []Provider`
- THEN the registry SHALL return all registered providers
- AND providers SHALL be sorted by Priority ascending

### Requirement: Per-Provider File Organization
The init system SHALL organize provider implementations as separate Go files under `internal/initialize/providers/`, with one file per provider.

#### Scenario: Provider file structure
- WHEN a provider file is created
- THEN it SHALL be named `{provider-id}.go` (e.g., `claude.go`, `gemini.go`)
- AND it SHALL contain an `init()` function that registers its provider
- AND it SHALL be self-contained with all provider-specific configuration

#### Scenario: Adding a new provider
- WHEN a developer adds a new AI CLI provider
- THEN they SHALL create a single file under `internal/initialize/providers/`
- AND the file SHALL implement the `Provider` interface
- AND the file SHALL call `Register()` in its `init()` function
- AND no other files SHALL require modification

### Requirement: Init Function Registration
The init system SHALL use Go's `init()` function pattern for automatic provider registration at startup.

#### Scenario: Auto-registration at startup
- WHEN the program starts
- THEN all provider `init()` functions SHALL execute before `main()`
- AND all providers SHALL be registered in the global registry
- AND registration order SHALL not affect functionality

### Requirement: Command Format Support
The init system SHALL support multiple command file formats through the `CommandFormat` type.

#### Scenario: Markdown command format
- WHEN a provider returns `FormatMarkdown` from `CommandFormat()`
- THEN slash command files SHALL be created as `.md` files
- AND files SHALL use frontmatter and spectr markers

#### Scenario: TOML command format
- WHEN a provider returns `FormatTOML` from `CommandFormat()`
- THEN slash command files SHALL be created as `.toml` files
- AND the TOML SHALL include `description` field with command description
- AND the TOML SHALL include `prompt` field with the command prompt content

### Requirement: Version Command Structure
The CLI SHALL provide a `version` command that displays version information including version number, git commit hash, and build date.

#### Scenario: Version command registration
- WHEN the CLI is initialized
- THEN it SHALL include a VersionCmd struct field tagged with `cmd`
- AND the command SHALL be accessible via `spectr version`
- AND help text SHALL describe version display functionality

#### Scenario: Version command invocation
- WHEN user runs `spectr version` without flags
- THEN the system displays version in format: `spectr version {version} (commit: {commit}, built: {date})`
- AND version SHALL be the semantic version (e.g., `0.1.0` or `dev`)
- AND commit SHALL be the git commit hash (short or full) or `unknown`
- AND date SHALL be the build date in ISO 8601 format or `unknown`

#### Scenario: Version command with short flag
- WHEN user runs `spectr version --short`
- THEN the system displays only the version number (e.g., `0.1.0`)
- AND no other information is displayed

#### Scenario: Version command with JSON flag
- WHEN user runs `spectr version --json`
- THEN the system outputs version data as JSON
- AND JSON SHALL include fields: `version`, `commit`, `date`
- AND SHALL be parseable by standard JSON tools

### Requirement: Version Variable Injection
The version information SHALL be injectable at build time via Go ldflags, supporting both goreleaser releases and nix flake builds.

#### Scenario: Goreleaser version injection
- WHEN goreleaser builds the binary
- THEN version SHALL be set from git tag via ldflags
- AND commit SHALL be set from git commit hash via ldflags
- AND date SHALL be set from build timestamp via ldflags

#### Scenario: Nix flake version injection
- WHEN nix builds the binary via flake.nix
- THEN version SHALL be set from the flake package version attribute via ldflags
- AND commit and date MAY be `unknown` if not available in nix build context

#### Scenario: Development build defaults
- WHEN binary is built without ldflags (e.g., `go build`)
- THEN version SHALL default to `dev`
- AND commit SHALL default to `unknown`
- AND date SHALL default to `unknown`

### Requirement: Version Package Location
The version variables SHALL be defined in a dedicated `internal/version` package for clean separation and easy ldflags targeting.

#### Scenario: Package structure
- WHEN the version package is imported
- THEN it SHALL expose `Version`, `Commit`, and `Date` string variables
- AND variables SHALL have default values for development builds
- AND the ldflags path SHALL be `github.com/connerohnesorge/spectr/internal/version`

### Requirement: Completion Command Structure
The CLI SHALL provide a `completion` subcommand that outputs shell completion scripts for supported shells using the kong-completion library.

#### Scenario: Completion command registration
- WHEN the CLI is initialized
- THEN it SHALL include a Completion field using `kongcompletion.Completion` type
- AND the command SHALL be accessible via `spectr completion`
- AND help text SHALL describe shell completion functionality

#### Scenario: Bash completion output
- WHEN user runs `spectr completion bash`
- THEN the system outputs a valid bash completion script
- AND the script can be sourced directly or added to bash-completion.d

#### Scenario: Zsh completion output
- WHEN user runs `spectr completion zsh`
- THEN the system outputs a valid zsh completion script
- AND the script can be sourced or placed in $fpath

#### Scenario: Fish completion output
- WHEN user runs `spectr completion fish`
- THEN the system outputs a valid fish completion script
- AND the script can be saved to fish completions directory

### Requirement: Custom Predictors for Dynamic Arguments
The completion system SHALL provide context-aware suggestions for arguments that accept dynamic values like change IDs or spec IDs.

#### Scenario: Change ID completion
- WHEN user types `spectr archive <TAB>` or `spectr validate <TAB>`
- AND the argument expects a change ID
- THEN completion suggests all active change IDs from `spectr/changes/`
- AND excludes archived changes

#### Scenario: Spec ID completion
- WHEN user types `spectr validate --type spec <TAB>`
- AND the argument expects a spec ID
- THEN completion suggests all spec IDs from `spectr/specs/`

#### Scenario: Item type completion
- WHEN user types `spectr validate --type <TAB>`
- THEN completion suggests `change` and `spec`

### Requirement: Kong-Completion Integration Pattern
The CLI initialization SHALL follow the kong-completion pattern where Kong is initialized, completions are registered, and then arguments are parsed.

#### Scenario: Initialization order
- WHEN the program starts
- THEN `kong.Must()` is called first to create the Kong application
- AND `kongcompletion.Register()` is called before parsing
- AND `app.Parse()` is called after completion registration
- AND this order ensures completions work correctly

#### Scenario: Predictor registration
- WHEN custom predictors are defined
- THEN they SHALL be registered via `kongcompletion.WithPredictor()`
- AND struct fields SHALL reference predictors using `predictor:"name"` tag

### Requirement: Accept Command Structure
The CLI SHALL provide an `accept` command that converts `tasks.md` to `tasks.jsonc` format with header comments for stable agent manipulation during implementation.

#### Scenario: Accept command registration
- WHEN the CLI is initialized
- THEN it SHALL include an AcceptCmd struct field tagged with `cmd`
- AND the command SHALL be accessible via `spectr accept`
- AND the command help text SHALL reference tasks.jsonc output

#### Scenario: Accept with change ID
- WHEN user runs `spectr accept <change-id>`
- THEN the system validates the change exists in `spectr/changes/<change-id>/`
- AND the system parses `tasks.md` into structured format
- AND the system writes `tasks.jsonc` with proper schema and header comments
- AND the system removes `tasks.md` to prevent drift

#### Scenario: Accept with validation
- WHEN user runs `spectr accept <change-id>`
- THEN the system validates the change before conversion
- AND if validation fails, the system aborts
- AND the system displays validation errors

#### Scenario: Accept dry-run mode
- WHEN user runs `spectr accept <change-id> --dry-run`
- THEN the system displays what would be converted
- AND the system does NOT write tasks.jsonc
- AND the system does NOT remove tasks.md

#### Scenario: Accept already accepted change
- WHEN user runs `spectr accept <change-id>` on a change that already has tasks.jsonc
- THEN the system displays a message indicating change is already accepted
- AND the system exits with code 0 (success, idempotent)

#### Scenario: Accept change without tasks.md
- WHEN user runs `spectr accept <change-id>` on a change without tasks.md
- THEN the system displays an error indicating no tasks.md found
- AND the system exits with code 1

### Requirement: Tasks JSON Schema
The accept command SHALL generate `tasks.jsonc` files conforming to a versioned schema with structured task objects and header comments.

#### Scenario: JSONC file structure
- WHEN the accept command creates tasks.jsonc
- THEN the file SHALL start with header comments documenting task status values
- AND the file SHALL contain a root object with `version` and `tasks` fields
- AND `version` SHALL be integer 1 for this schema version
- AND `tasks` SHALL be an array of task objects

#### Scenario: Header comment content
- WHEN the accept command creates tasks.jsonc
- THEN the header SHALL use `//` line comment syntax
- AND the header SHALL indicate the file is machine-generated by `spectr accept`
- AND the header SHALL document the three valid status values: "pending", "in_progress", "completed"
- AND the header SHALL document valid status transitions: pending → in_progress → completed
- AND the header SHALL explain that agents should mark a task "in_progress" when starting work
- AND the header SHALL explain that agents should mark a task "completed" only after verification
- AND the header SHALL note that skipping directly from "pending" to "completed" is allowed for trivial tasks

#### Scenario: Task object structure
- WHEN a task is serialized to JSONC
- THEN it SHALL have `id` field containing the task identifier (e.g., "1.1")
- AND it SHALL have `section` field containing the section header (e.g., "Implementation")
- AND it SHALL have `description` field containing the full task text
- AND it SHALL have `status` field with value "pending", "in_progress", or "completed"

#### Scenario: Status value mapping from Markdown
- WHEN converting tasks.md to tasks.jsonc
- THEN `- [ ]` SHALL map to status "pending"
- AND `- [x]` (case-insensitive) SHALL map to status "completed"

### Requirement: Accept Command Flags
The accept command SHALL support flags for controlling behavior.

#### Scenario: Dry-run flag
- WHEN user provides the `--dry-run` flag
- THEN the system previews the conversion without writing files
- AND displays the JSON that would be generated

#### Scenario: Interactive change selection
- WHEN user runs `spectr accept` without specifying a change ID
- THEN the system displays a list of active changes with tasks.md files
- AND prompts for selection using existing TUI components

### Requirement: List Command Alias
The `spectr list` command SHALL support `ls` as a shorthand alias, allowing users to invoke `spectr ls` as equivalent to `spectr list`.

#### Scenario: User runs spectr ls shorthand
- WHEN user runs `spectr ls`
- THEN the system displays the list of changes identically to `spectr list`
- AND all flags (`--specs`, `--all`, `--long`, `--json`, `--interactive`) work with the alias

#### Scenario: User runs spectr ls with flags
- WHEN user runs `spectr ls --specs --long`
- THEN the command behaves identically to `spectr list --specs --long`
- AND specs are displayed in long format

#### Scenario: Help text shows list alias
- WHEN user runs `spectr --help`
- THEN the help text displays `list` with its `ls` alias
- AND the alias is shown in parentheses or as comma-separated alternatives

### Requirement: Item Name Path Normalization
Commands accepting item names (validate, archive, accept) SHALL normalize path arguments to extract the item ID and infer the item type from the path structure.

#### Scenario: Path with spectr/changes prefix
- WHEN user runs a command with argument `spectr/changes/my-change`
- THEN the system SHALL extract `my-change` as the item ID
- AND SHALL infer the item type as "change"

#### Scenario: Path with spectr/changes prefix and trailing content
- WHEN user runs a command with argument `spectr/changes/my-change/specs/foo/spec.md`
- THEN the system SHALL extract `my-change` as the item ID
- AND SHALL infer the item type as "change"

#### Scenario: Path with spectr/specs prefix
- WHEN user runs a command with argument `spectr/specs/my-spec`
- THEN the system SHALL extract `my-spec` as the item ID
- AND SHALL infer the item type as "spec"

#### Scenario: Path with spectr/specs prefix and spec.md file
- WHEN user runs a command with argument `spectr/specs/my-spec/spec.md`
- THEN the system SHALL extract `my-spec` as the item ID
- AND SHALL infer the item type as "spec"

#### Scenario: Simple ID without path prefix
- WHEN user runs a command with argument `my-change`
- THEN the system SHALL use `my-change` as-is for lookup
- AND SHALL use existing auto-detection logic for item type

#### Scenario: Absolute path normalization
- WHEN user runs a command with argument `/home/user/project/spectr/changes/my-change`
- THEN the system SHALL extract `my-change` as the item ID
- AND SHALL infer the item type as "change"

#### Scenario: Inferred type precedence
- WHEN user provides a path argument that contains `spectr/changes/` or `spectr/specs/`
- THEN the inferred type from path SHALL be used for validation
- AND SHALL NOT trigger "exists as both change and spec" ambiguity errors

### Requirement: Interactive List Mode
The interactive list mode in `spectr list` is extended to support unified display of changes and specifications alongside existing separate modes.

#### Previous behavior
The system displays either changes OR specs in interactive mode based on the `--specs` flag. Columns and behavior are specific to each item type.

#### New behavior
- When `--all` is provided with `--interactive`, both changes and specs are shown together with unified columns
- When neither `--all` nor `--specs` are provided, changes-only mode is default (backward compatible)
- When `--specs` is provided without `--all`, specs-only mode is used (backward compatible)
- Each item type is clearly labeled in the Type column (CHANGE or SPEC)
- Type-aware actions apply based on selected item (edit only for specs)

#### Scenario: Default behavior unchanged
- WHEN the user runs `spectr list --interactive`
- THEN the behavior is identical to before this change
- AND only changes are displayed
- AND columns show: ID, Title, Deltas, Tasks

#### Scenario: Unified mode opt-in
- WHEN the user explicitly uses `--all --interactive`
- THEN the new unified behavior is enabled
- AND users must opt-in to the new functionality
- AND columns show: Type, ID, Title, Details (context-aware)

#### Scenario: Unified mode displays both types
- WHEN unified mode is active
- THEN changes show Type="CHANGE" with delta and task counts
- AND specs show Type="SPEC" with requirement counts
- AND both types are navigable and selectable in the same table

#### Scenario: Type-specific actions in unified mode
- WHEN user presses 'e' on a change row in unified mode
- THEN the action is ignored (no edit for changes)
- AND help text does not show 'e' option
- WHEN user presses 'e' on a spec row in unified mode
- THEN the spec opens in the editor as usual

#### Scenario: Help text uses minimal footer by default
- WHEN interactive mode is displayed in any mode (changes, specs, or unified)
- THEN the footer shows: item count, project path, and `?: help`
- AND the full hotkey reference is hidden until `?` is pressed

#### Scenario: Help text format for changes mode
- WHEN user presses `?` in changes mode (`spectr list -I`)
- THEN the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | q: quit`
- AND pressing `?` again or navigating hides the help

#### Scenario: Help text format for specs mode
- WHEN user presses `?` in specs mode (`spectr list --specs -I`)
- THEN the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | q: quit`
- AND archive hotkey is NOT shown (specs cannot be archived)

### Requirement: Clipboard Copy on Selection
When a user presses Enter on a selected row in interactive mode, the item's ID SHALL be copied to the system clipboard.

#### Scenario: Copy change ID to clipboard
- WHEN user selects a change row and presses Enter
- THEN the change ID (kebab-case identifier) is copied to clipboard
- AND a success message is displayed (e.g., "Copied: add-archive-command")
- AND the interactive mode exits

#### Scenario: Copy spec ID to clipboard
- WHEN user selects a spec row and presses Enter
- THEN the spec ID is copied to clipboard
- AND a success message is displayed
- AND the interactive mode exits

#### Scenario: Clipboard failure handling
- WHEN clipboard operation fails
- THEN display error message to user
- AND do not exit interactive mode
- AND user can retry or quit manually

### Requirement: Interactive Mode Exit Controls
Users SHALL be able to exit interactive mode using standard quit commands.

#### Scenario: Quit with q key
- WHEN user presses 'q'
- THEN interactive mode exits
- AND no clipboard operation occurs
- AND command returns successfully

#### Scenario: Quit with Ctrl+C
- WHEN user presses Ctrl+C
- THEN interactive mode exits immediately
- AND no clipboard operation occurs
- AND command returns successfully

### Requirement: Table Visual Styling
The interactive table SHALL use clear visual styling to distinguish headers, selected rows, and borders, provided by the shared `internal/tui` package.

#### Scenario: Visual hierarchy in table
- WHEN interactive mode is displayed
- THEN column headers are visually distinct from data rows
- AND selected row has contrasting background/foreground colors
- AND table borders are visible and styled consistently
- AND table fits within terminal width gracefully
- AND styling SHALL be applied via `tui.ApplyTableStyles()`

#### Scenario: Consistent styling across commands
- WHEN user uses `spectr list -I`, `spectr archive`, or `spectr validate` interactive modes
- THEN all tables SHALL use identical styling
- AND colors, borders, and highlights SHALL match exactly
- AND the shared `tui.ApplyTableStyles()` function SHALL be the single source of truth

### Requirement: Cross-Platform Clipboard Support
Clipboard operations SHALL work across Linux, macOS, and Windows platforms.

#### Scenario: Clipboard on Linux
- WHEN running on Linux
- THEN clipboard operations use X11 or Wayland clipboard APIs as appropriate
- AND fallback to OSC 52 escape sequences if desktop clipboard unavailable

#### Scenario: Clipboard on macOS
- WHEN running on macOS
- THEN clipboard operations use pbcopy or native clipboard APIs

#### Scenario: Clipboard on Windows
- WHEN running on Windows
- THEN clipboard operations use Windows clipboard APIs

#### Scenario: Clipboard in SSH/remote session
- WHEN running over SSH without X11 forwarding
- THEN use OSC 52 escape sequences to copy to local clipboard
- AND document this behavior for users

### Requirement: Initialization Next Steps Message

The `spectr init` command SHALL display a formatted "Next steps" message after successful initialization that provides users with clear, actionable guidance for getting started with Spectr.

The message SHALL include:
1. Three progressive steps with copy-paste ready prompts for AI assistants
2. Visual separators to make the message stand out
3. References to key Spectr files and documentation
4. Placeholder text that users can customize (e.g., "[YOUR FEATURE HERE]")

The init command SHALL NOT automatically create project files outside the `spectr/` directory (such as README.md). Users maintain full control over their project's root-level documentation.

#### Scenario: Interactive mode initialization succeeds

- WHEN a user completes initialization via the interactive TUI wizard
- THEN the completion screen SHALL display the next steps message
- AND the message SHALL appear after the list of created/updated files
- AND the message SHALL be visually distinct with a separator line
- AND the message SHALL provide three numbered steps with specific prompts

#### Scenario: Non-interactive mode initialization succeeds

- WHEN a user runs `spectr init --non-interactive` and initialization succeeds
- THEN the command output SHALL display the next steps message
- AND the message SHALL appear after the list of created/updated files
- AND the message SHALL be formatted consistently with the interactive mode
- AND the message SHALL include the same three progressive steps

#### Scenario: Initialization fails with errors

- WHEN initialization fails with errors
- THEN the next steps message SHALL NOT be displayed
- AND only error messages SHALL be shown

#### Scenario: Next steps message content

- WHEN the next steps message is displayed
- THEN step 1 SHALL guide users to populate spectr/project.md
- AND step 2 SHALL guide users to create their first change proposal
- AND step 3 SHALL guide users to learn the Spectr workflow from spectr/AGENTS.md
- AND each step SHALL include a complete, copy-paste ready prompt in quotes
- AND the message SHALL include a visual separator using dashes or similar characters

#### Scenario: Init does not create README

- WHEN a user runs `spectr init` on a project without a README.md
- THEN the init command SHALL NOT create a README.md file
- AND only files within the `spectr/` directory SHALL be created
- AND tool-specific files (e.g., CLAUDE.md, .cursor/) SHALL be created as configured

### Requirement: Flat Tool List in Initialization Wizard

The initialization wizard SHALL present all AI tool options in a single unified flat list without visual grouping by tool type. Slash-only tool entries SHALL be removed from the registry as their functionality is now provided via automatic installation when the corresponding config-based tool is selected.

#### Scenario: Display only config-based tools in wizard

- WHEN user runs `spectr init` and reaches the tool selection screen
- THEN only config-based AI tools are displayed (e.g., `claude-code`, `cline`, `cursor`)
- AND slash-only tool entries (e.g., `claude`, `kilocode`) are not shown
- AND tools are sorted by priority
- AND no section headers (e.g., "Config-Based Tools", "Slash Command Tools") are shown
- AND each tool appears as a single checkbox item with its name

#### Scenario: Keyboard navigation across displayed tools

- WHEN user navigates with arrow keys (↑/↓)
- THEN the cursor moves through all displayed config-based tools sequentially
- AND navigation is continuous without group boundaries
- AND the first tool is selected by default on screen load

#### Scenario: Tool selection works uniformly

- WHEN user presses space to toggle any tool
- THEN the checkbox state changes (checked/unchecked)
- AND selection state is preserved when navigating
- AND both config file and slash commands will be installed when confirmed

#### Scenario: Bulk selection operations

- WHEN user presses 'a' to select all
- THEN all displayed config-based tools are checked
- AND WHEN user presses 'n' to select none
- THEN all tools are unchecked
- AND operations work across all displayed tools

#### Scenario: Help text clarity

- WHEN the tool selection screen is displayed
- THEN the help text shows keyboard controls (↑/↓, space, a, n, enter, q)
- AND the help text does NOT reference tool groupings or categories
- AND the screen title clearly indicates "Select AI Tools to Configure"

#### Scenario: Reduced tool count in wizard

- WHEN the wizard displays the tool list
- THEN fewer total tools are shown compared to the previous implementation
- AND the count reflects only config-based tools (not slash-only duplicates)
- AND navigation and selection work correctly with the reduced count

### Requirement: Interactive Archive Mode
The archive command SHALL provide an interactive table interface when no change ID argument is provided or when the `-I` or `--interactive` flag is used, displaying available changes in a navigable table format identical to the list command's interactive mode with project path information.

#### Scenario: User runs archive with no arguments
- WHEN user runs `spectr archive` with no change ID argument
- THEN an interactive table is displayed with columns: ID, Title, Deltas, Tasks
- AND the table supports arrow key navigation (↑/↓, j/k)
- AND the first row is selected by default
- AND the table uses the same visual styling as list -I
- AND the project path is displayed in the interface

#### Scenario: User runs archive with -I flag
- WHEN user runs `spectr archive -I`
- THEN an interactive table is displayed even if other flags are present
- AND the behavior is identical to running archive with no arguments
- AND the project path is displayed in the interface

#### Scenario: User selects change for archiving
- WHEN user presses Enter on a selected row in archive interactive mode
- THEN the change ID is captured (not copied to clipboard)
- AND the interactive mode exits
- AND the archive workflow proceeds with the selected change ID
- AND validation, task checking, and spec updates proceed as normal

#### Scenario: User cancels archive selection
- WHEN user presses 'q' or Ctrl+C in archive interactive mode
- THEN interactive mode exits
- AND archive command returns successfully without archiving anything
- AND a "Cancelled" message is displayed

#### Scenario: No changes available for archiving
- WHEN user runs `spectr archive` and no changes exist in changes/ directory
- THEN display "No changes available to archive" message
- AND exit cleanly without entering interactive mode
- AND command returns successfully

#### Scenario: Archive with explicit change ID bypasses interactive mode
- WHEN user runs `spectr archive <change-id>`
- THEN interactive mode is NOT triggered
- AND archive proceeds directly with the specified change ID
- AND behavior is unchanged from current implementation

### Requirement: Archive Interactive Table Display
The archive command's interactive table SHALL display the same information columns as the list command to help users make informed archiving decisions.

#### Scenario: Table columns match list command
- WHEN archive interactive mode is displayed
- THEN columns are: ID (30 chars), Title (40 chars), Deltas (10 chars), Tasks (15 chars)
- AND column widths match the list -I command exactly
- AND title text is truncated with ellipsis if longer than 38 characters
- AND task status shows format "completed/total" (e.g., "5/10")

#### Scenario: Visual styling consistency
- WHEN archive interactive table is displayed
- THEN the table uses identical styling to list -I
- AND column headers are visually distinct from data rows
- AND selected row has contrasting background/foreground colors
- AND table borders are visible and styled consistently
- AND help text shows navigation controls (↑/↓, j/k, enter, q)

### Requirement: Archive Selection Without Clipboard
The archive command's interactive mode SHALL NOT copy the selected change ID to the clipboard, unlike the list command, since the ID is immediately consumed by the archive workflow.

#### Scenario: Enter key captures selection
- WHEN user presses Enter on a selected change
- THEN the change ID is captured internally
- AND NO clipboard operation occurs
- AND NO "Copied: <id>" message is displayed
- AND the archive workflow proceeds immediately with the selected ID

#### Scenario: Workflow continuation
- WHEN a change is selected in interactive mode
- THEN the Archiver.Archive() method receives the selected change ID
- AND validation, task checking, and spec updates proceed as if the ID was provided as an argument
- AND all confirmation prompts and flags (--yes, --skip-specs) work normally

### Requirement: Validation Output Format
The validate command SHALL display validation issues in a consistent, detailed format for both single-item and bulk validation modes.

#### Scenario: Single item validation with issues
- WHEN user runs `spectr validate <item>` and validation finds issues
- THEN output SHALL display "✗ <item> has N issue(s):"
- AND each issue SHALL be displayed on a separate line with format "  [LEVEL] PATH: MESSAGE"
- AND the command SHALL exit with code 1

#### Scenario: Bulk validation with issues
- WHEN user runs `spectr validate --all` and validation finds issues in multiple items
- THEN output SHALL display "✗ <item> (<type>): N issue(s)" for each failed item
- AND immediately following each failed item, all issue details SHALL be displayed
- AND each issue SHALL use the format "  [LEVEL] PATH: MESSAGE"
- AND a summary line SHALL display "N passed, M failed, T total"
- AND the command SHALL exit with code 1

#### Scenario: Bulk validation all passing
- WHEN user runs `spectr validate --all` and all items are valid
- THEN output SHALL display "✓ <item> (<type>)" for each item
- AND a summary line SHALL display "N passed, 0 failed, N total"
- AND the command SHALL exit with code 0

#### Scenario: JSON output format
- WHEN user provides `--json` flag with any validation command
- THEN output SHALL be valid JSON
- AND SHALL include full issue details with level, path, and message fields
- AND SHALL include per-item results and summary statistics

### Requirement: Editor Hotkey in Interactive Specs List
The interactive specs list mode SHALL provide an 'e' hotkey that opens the selected spec file in the user's configured editor.

#### Scenario: User presses 'e' to edit a spec
- WHEN user is in interactive specs mode (`spectr list --specs -I`)
- AND user presses the 'e' key on a selected spec
- THEN the file `spectr/specs/<spec-id>/spec.md` is opened in the editor specified by $EDITOR environment variable
- AND the TUI waits for the editor to close
- AND the TUI remains active after the editor closes
- AND the same row remains selected

#### Scenario: User edits spec and returns to TUI
- WHEN user presses 'e' to open a spec
- AND makes changes in the editor and saves
- AND closes the editor
- THEN the TUI returns to the interactive list view
- AND the user can continue navigating or edit another spec
- AND the user can quit with 'q' or Ctrl+C as normal

#### Scenario: EDITOR environment variable not set
- WHEN user presses 'e' to edit a spec
- AND $EDITOR environment variable is not set
- THEN display an error message "EDITOR environment variable not set"
- AND the TUI remains in interactive mode
- AND the user can continue navigating or quit

#### Scenario: Spec file does not exist
- WHEN user presses 'e' to edit a spec
- AND the spec file at `spectr/specs/<spec-id>/spec.md` does not exist
- THEN display an error message "Spec file not found: <path>"
- AND the TUI remains in interactive mode
- AND the user can continue navigating or quit

#### Scenario: Editor launch fails
- WHEN user presses 'e' to edit a spec
- AND the editor process fails to launch (e.g., editor binary not found, permission error)
- THEN display an error message with the underlying error details
- AND the TUI remains in interactive mode
- AND the user can retry or quit

#### Scenario: Help text shows editor hotkey
- WHEN interactive specs mode is displayed
- THEN the help text includes "e: edit spec" or similar guidance
- AND the help text shows all available keys including navigation, enter, e, and quit keys

### Requirement: Editor Hotkey Scope
The 'e' hotkey for opening files in $EDITOR SHALL only be available in specs list mode, not in changes list mode.

#### Scenario: Editor hotkey not available for changes
- WHEN user is in interactive changes mode (`spectr list -I`)
- AND user presses 'e' key
- THEN the key press is ignored (no action taken)
- AND the help text does NOT show 'e: edit' option
- AND only standard navigation and clipboard actions are available

### Requirement: Project Path Display in Interactive Mode
The interactive table interfaces SHALL display the project root path to provide users with context about which project they are working with.

#### Scenario: Project path shown in changes interactive mode
- WHEN user runs `spectr list -I` for changes
- THEN the project root path is displayed in the help text or table header
- AND the path is the absolute path to the project directory

#### Scenario: Project path shown in specs interactive mode
- WHEN user runs `spectr list --specs -I`
- THEN the project root path is displayed in the help text or table header
- AND the path is the absolute path to the project directory

#### Scenario: Project path shown in archive interactive mode
- WHEN user runs `spectr archive` without arguments
- THEN the project root path is displayed in the help text or table header
- AND the path is the absolute path to the project directory

#### Scenario: Project path properly initialized for changes
- WHEN `RunInteractiveChanges()` is invoked
- THEN the `projectPath` parameter is passed from the calling command
- AND the `projectPath` field on `interactiveModel` is set during initialization

#### Scenario: Project path properly initialized for archive
- WHEN `RunInteractiveArchive()` is invoked
- THEN the `projectPath` parameter is passed from the calling command
- AND the `projectPath` field on `interactiveModel` is set during initialization

### Requirement: Unified Item List Display
The system SHALL display changes and specifications together in a single interactive table when invoked with appropriate flags, allowing users to browse both item types simultaneously with clear visual differentiation.

#### Scenario: User opens unified interactive list
- WHEN the user runs `spectr list --interactive --all` from a directory with both changes and specs
- THEN a table appears showing both changes and specs rows
- AND each row indicates its type (change or spec)
- AND the table maintains correct ordering and alignment

#### Scenario: Unified list shows correct columns
- WHEN the unified interactive mode is active
- THEN the table displays: Type, ID, Title, and Type-Specific Details columns
- AND "Type-Specific Details" shows "Deltas/Tasks" for changes
- AND "Type-Specific Details" shows "Requirements" for specs

#### Scenario: User navigates mixed items
- WHEN the user navigates with arrow keys through a mixed list
- THEN the cursor moves smoothly between change and spec rows
- AND help text remains accurate and updated
- AND the selected row is clearly highlighted

### Requirement: Type-Aware Item Selection
The system SHALL track whether a selected item is a change or spec and provide type-appropriate actions (e.g., edit only works for specs).

#### Scenario: Selecting a spec in unified mode
- WHEN the user presses Enter on a spec row
- THEN the spec ID is copied to clipboard
- AND a success message displays the ID and type indicator
- AND the user is returned to the interactive session or exited cleanly

#### Scenario: Selecting a change in unified mode
- WHEN the user presses Enter on a change row
- THEN the change ID is copied to clipboard
- AND a success message displays the ID and type indicator
- AND no edit action is attempted

#### Scenario: Edit action restricted to specs
- WHEN the user presses 'e' on a change row in unified mode
- THEN the action is ignored or a helpful message appears
- AND the interactive session continues

### Requirement: Backward-Compatible Separate Modes
The system SHALL maintain existing interactive modes for changes-only and specs-only when `--all` flag is not provided.

#### Scenario: Changes-only mode still works
- WHEN the user runs `spectr list --interactive` without `--all`
- THEN only changes are displayed
- AND behavior is identical to the previous implementation
- AND edit functionality works as before for changes

#### Scenario: Specs-only mode still works
- WHEN the user runs `spectr list --specs --interactive` without `--all`
- THEN only specs are displayed
- AND behavior is identical to the previous implementation
- AND edit functionality works as before for specs

### Requirement: Enhanced List Command Flags
The system SHALL support new flag combinations to control listing behavior while maintaining validation for mutually exclusive options.

#### Scenario: Flag validation for unified mode
- WHEN the user attempts `spectr list --interactive --all --json`
- THEN an error message is returned: "cannot use --interactive with --json"
- AND the command exits without running

#### Scenario: All flag with separate type flags
- WHEN the user provides `--all` with `--specs`
- THEN `--all` takes precedence and unified mode is used
- AND a warning may be shown (optional) about the redundant flag

#### Scenario: All flag in non-interactive mode
- WHEN the user runs `spectr list --all` without `--interactive`
- THEN both changes and specs are listed in text format
- AND each item shows its type in the output

### Requirement: Automatic Slash Command Installation

When a config-based AI tool is selected during initialization, the system SHALL automatically install the corresponding slash command files for that tool without requiring separate user selection.

Config-based tools include those that create instruction files (e.g., `claude-code` creates `CLAUDE.md`). Slash command files are the workflow command files (e.g., `.claude/commands/spectr/proposal.md`).

The `ToolDefinition` model SHALL NOT include a `ConfigPath` field, as actual file paths are determined by individual configurators. The registry maintains tool metadata (ID, Name, Type, Priority) but delegates file path resolution to configurator implementations. Tool IDs SHALL use a type-safe constant approach to prevent typos and enable compile-time validation.

This automatic installation provides users with complete Spectr integration in a single selection, eliminating the need for redundant tool entries in the wizard.

#### Scenario: Claude Code auto-installs slash commands

- WHEN user selects `claude-code` in the init wizard
- THEN the system creates `CLAUDE.md` in the project root
- AND the system creates `.claude/commands/spectr/proposal.md`
- AND the system creates `.claude/commands/spectr/apply.md`
- AND all files are tracked in the execution result
- AND the completion screen shows all 3 files created

#### Scenario: Multiple tools with slash commands selected

- WHEN user selects both `claude-code` and `cursor` in the init wizard
- THEN the system creates `CLAUDE.md` and both config + slash commands for Claude
- AND the system creates `.cursor/commands/spectr/proposal.md` and slash commands for Cursor
- AND all files from both tools are created and tracked separately
- AND the completion screen lists all created files grouped by tool

#### Scenario: Slash command files already exist

- WHEN user run init and selects `claude-code`
- AND `.claude/commands/spectr/proposal.md` already exists
- THEN the existing file's content between `<!-- spectr:START -->` and `<!-- spectr:END -->` is updated
- AND the file's YAML frontmatter is preserved
- AND no error occurs
- AND the file is marked as "updated" rather than "created" in execution result

### Requirement: Archive Hotkey in Interactive Changes Mode
The interactive changes list mode SHALL provide an 'a' hotkey that archives the currently selected change, invoking the same workflow as `spectr archive <change-id>`.

#### Scenario: User presses 'a' to archive a change
- WHEN user is in interactive changes mode (`spectr list -I`)
- AND user presses the 'a' key on a selected change
- THEN the interactive mode exits
- AND the archive workflow begins for the selected change ID
- AND validation, task checking, and spec updates proceed as if the ID was provided as an argument
- AND all confirmation prompts and flags work normally

#### Scenario: Archive hotkey not available in specs mode
- WHEN user is in interactive specs mode (`spectr list --specs -I`)
- AND user presses 'a' key
- THEN the key press is ignored (no action taken)
- AND the help text does NOT show 'a: archive' option

#### Scenario: Archive hotkey not available in unified mode
- WHEN user is in unified interactive mode (`spectr list --all -I`)
- AND user presses 'a' key
- THEN the key press is ignored (no action taken)
- AND the help text does NOT show 'a: archive' option
- AND this avoids confusion when a spec row is selected

#### Scenario: Archive workflow integration
- WHEN the archive hotkey triggers the archive workflow
- THEN the workflow uses the same code path as `spectr archive <id>`
- AND the selected change ID is passed to the archive workflow
- AND success or failure is reported after the workflow completes

#### Scenario: Help text shows archive hotkey in changes mode
- WHEN interactive changes mode is displayed
- THEN the help text includes `a: archive` in the controls line
- AND the hotkey appears after `e: edit` and before `q: quit`

### Requirement: Shared TUI Component Library

The CLI SHALL use a shared `internal/tui` package for interactive TUI components, providing consistent styling, behavior, and composable building blocks across all interactive modes.

#### Scenario: TablePicker used for item selection
- WHEN any command needs an interactive table-based selection (list, archive, validation item picker)
- THEN the command SHALL use the `TablePicker` component from `internal/tui`
- AND the table SHALL use consistent styling from `tui.ApplyTableStyles()`
- AND navigation keys (↑/↓, j/k) SHALL work identically across all usages
- AND quit keys (q, Ctrl+C) SHALL work identically across all usages

#### Scenario: MenuPicker used for option selection
- WHEN any command needs an interactive menu selection (validation mode menu)
- THEN the command SHALL use the `MenuPicker` component from `internal/tui`
- AND the menu SHALL use consistent styling
- AND navigation and selection behavior SHALL match the TablePicker patterns

#### Scenario: Consistent string truncation
- WHEN any TUI component needs to truncate text for display
- THEN it SHALL use `tui.TruncateString()` with consistent ellipsis handling
- AND truncation SHALL add "..." suffix when text exceeds max length
- AND very short max lengths (≤3) SHALL truncate without ellipsis

#### Scenario: Consistent clipboard operations
- WHEN any TUI component needs to copy text to clipboard
- THEN it SHALL use `tui.CopyToClipboard()` from the shared package
- AND the function SHALL try native clipboard first
- AND the function SHALL fall back to OSC 52 for remote sessions

#### Scenario: Action registration pattern
- WHEN a command configures a TablePicker with custom actions
- THEN actions SHALL be registered via `WithAction(key, label, handler)`
- AND the help text SHALL automatically include all registered actions
- AND unregistered keys SHALL be ignored (no error)

#### Scenario: Domain logic remains in consuming packages
- WHEN the tui package is used by list or validation
- THEN domain-specific logic (archive workflow, validation execution) SHALL remain in consuming packages
- AND the tui package SHALL only provide UI primitives
- AND business logic SHALL not be coupled to the tui package

### Requirement: Search Hotkey in Interactive Lists
The interactive list modes SHALL provide a '/' hotkey that activates a text search mode, allowing users to filter the displayed list by typing a search query that matches against item IDs and titles.

#### Scenario: User presses '/' to enter search mode
- WHEN user is in any interactive list mode (changes, specs, or unified)
- AND user presses the '/' key
- THEN search mode is activated
- AND a text input field is displayed below or above the table
- AND the cursor is placed in the text input field
- AND the user can type a search query

#### Scenario: Search filters rows in real-time
- WHEN search mode is active
- AND user types characters into the search input
- THEN the table rows are filtered in real-time
- AND only rows where ID or title contains the search query (case-insensitive) are displayed
- AND the first matching row is automatically selected

#### Scenario: Search with no matches shows empty table
- WHEN search mode is active
- AND user types a query that matches no items
- THEN the table displays no rows
- AND a message indicates no matches found

#### Scenario: User presses Escape to exit search mode
- WHEN search mode is active
- AND user presses the Escape key
- THEN search mode is deactivated
- AND the search query is cleared
- AND all items are displayed again in the table
- AND the text input field is hidden

#### Scenario: User presses '/' again to clear search
- WHEN search mode is active
- AND the search query is not empty
- AND user presses '/' key
- THEN the search input gains focus (normal text input behavior)

- WHEN search mode is active
- AND the search query is empty
- AND user presses '/' key
- THEN search mode is deactivated
- AND all items are displayed again

#### Scenario: Navigation works while searching
- WHEN search mode is active
- AND filtered results are displayed
- THEN arrow key navigation (up/down, j/k) moves through filtered rows
- AND Enter key copies the selected filtered item's ID
- AND other hotkeys (e, a, t) work on the selected filtered item

#### Scenario: Help text shows search hotkey
- WHEN interactive mode is displayed in any mode
- THEN the help text includes '/: search' in the controls line
- AND the search hotkey is shown for all modes (changes, specs, unified)

#### Scenario: Search mode visual indicator
- WHEN search mode is active
- THEN the search input field is visually distinct
- AND the current search query is visible
- AND the help text updates to show 'Esc: exit search'

### Requirement: Help Toggle Hotkey
The interactive TUI modes SHALL hide hotkey hints by default and reveal them only when the user presses `?`, reducing visual clutter while maintaining discoverability.

#### Scenario: Default view shows minimal footer
- WHEN user enters any interactive TUI mode (list, archive, validate)
- THEN the footer displays only: item count, project path, and `?: help`
- AND the full hotkey reference is NOT shown
- AND navigation and all other hotkeys remain functional

#### Scenario: User presses '?' to reveal help
- WHEN user presses `?` while in interactive mode
- THEN the full hotkey reference is displayed in the footer area
- AND the reference includes all available hotkeys for the current mode
- AND the view updates immediately

#### Scenario: User dismisses help by pressing '?' again
- WHEN user presses `?` while help is visible
- THEN the help is hidden
- AND the minimal footer is restored

#### Scenario: Help auto-hides on navigation
- WHEN user presses a navigation key (↑/↓/j/k) while help is visible
- THEN the help is automatically hidden
- AND the navigation action is performed
- AND the minimal footer is restored

#### Scenario: Help content matches mode
- WHEN help is displayed in changes mode
- THEN the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | q: quit`
- WHEN help is displayed in specs mode
- THEN the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | q: quit`
- WHEN help is displayed in unified mode
- THEN the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | t: filter | q: quit`

### Requirement: Partial Change ID Resolution for Archive Command

The `spectr archive` command SHALL support intelligent partial ID matching when a non-exact change ID is provided as an argument. The resolution algorithm SHALL prioritize prefix matches over substring matches and require a unique match to proceed.

#### Scenario: Exact ID match takes precedence

- WHEN user runs `spectr archive add-feature`
- AND a change with ID `add-feature` exists
- THEN the archive proceeds with `add-feature`
- AND no resolution message is displayed

#### Scenario: Unique prefix match resolves successfully

- WHEN user runs `spectr archive refactor`
- AND only one change ID starts with `refactor` (e.g., `refactor-unified-interactive-tui`)
- THEN a message is displayed: "Resolved 'refactor' -> 'refactor-unified-interactive-tui'"
- AND the archive proceeds with the resolved ID

#### Scenario: Unique substring match resolves successfully

- WHEN user runs `spectr archive unified`
- AND no change ID starts with `unified`
- AND only one change ID contains `unified` (e.g., `refactor-unified-interactive-tui`)
- THEN a message is displayed: "Resolved 'unified' -> 'refactor-unified-interactive-tui'"
- AND the archive proceeds with the resolved ID

#### Scenario: Multiple prefix matches cause error

- WHEN user runs `spectr archive add`
- AND multiple change IDs start with `add` (e.g., `add-feature`, `add-hotkey`)
- THEN an error is displayed: "Ambiguous ID 'add' matches multiple changes: add-feature, add-hotkey"
- AND the command exits with error code 1
- AND no archive operation is performed

#### Scenario: Multiple substring matches cause error

- WHEN user runs `spectr archive search`
- AND no change ID starts with `search`
- AND multiple change IDs contain `search` (e.g., `add-search-hotkey`, `update-search-ui`)
- THEN an error is displayed: "Ambiguous ID 'search' matches multiple changes: add-search-hotkey, update-search-ui"
- AND the command exits with error code 1
- AND no archive operation is performed

#### Scenario: No match found

- WHEN user runs `spectr archive nonexistent`
- AND no change ID matches `nonexistent` (neither prefix nor substring)
- THEN an error is displayed: "No change found matching 'nonexistent'"
- AND the command exits with error code 1
- AND no archive operation is performed

#### Scenario: Case-insensitive matching

- WHEN user runs `spectr archive REFACTOR`
- AND a change ID `refactor-unified-interactive-tui` exists
- THEN the partial match succeeds (case-insensitive)
- AND the archive proceeds with the resolved ID

#### Scenario: Prefix match preferred over substring match

- WHEN user runs `spectr archive add`
- AND change ID `add-feature` exists (prefix match)
- AND change ID `update-add-button` exists (substring match only)
- THEN the prefix match `add-feature` is selected
- AND the substring-only match is ignored in preference calculation

### Requirement: Configured Provider Detection in Init Wizard

The initialization wizard SHALL detect which AI tool providers are already configured for the project and display this status in the tool selection screen. Already-configured providers SHALL be visually distinguished and pre-selected by default.

#### Scenario: Display configured indicator for already-configured providers

- WHEN user runs `spectr init` on a project with `CLAUDE.md` already present
- AND user reaches the tool selection screen
- THEN the Claude Code entry displays a "configured" indicator (e.g., dimmed text or badge)
- AND the indicator is visually distinct from the selection checkbox
- AND other unconfigured providers do NOT show the configured indicator

#### Scenario: Pre-select already-configured providers

- WHEN user runs `spectr init` on a project with some providers already configured
- AND user reaches the tool selection screen
- THEN already-configured providers have their checkboxes pre-selected
- AND users can deselect them if they don't want to update the configuration
- AND unconfigured providers remain unselected by default

#### Scenario: Help text explains configured indicator

- WHEN user is on the tool selection screen
- THEN the help text or screen description explains what the "configured" indicator means
- AND the explanation clarifies that selecting a configured provider will update its files

#### Scenario: No configured providers

- WHEN user runs `spectr init` on a fresh project with no providers configured
- AND user reaches the tool selection screen
- THEN no providers show the configured indicator
- AND no providers are pre-selected
- AND the screen functions as before this change

#### Scenario: All providers configured

- WHEN user runs `spectr init` on a project with all available providers configured
- AND user reaches the tool selection screen
- THEN all providers show the configured indicator
- AND all providers are pre-selected
- AND user can deselect providers they don't want to update

#### Scenario: Configured detection uses provider's IsConfigured method

- WHEN the wizard initializes
- THEN it calls `IsConfigured(projectPath)` on each provider
- AND the result is cached for the wizard session (not re-checked on each render)
- AND providers with global paths (like Codex) are correctly detected

### Requirement: Instruction File Pointer Template

The system SHALL use a short pointer template when injecting Spectr instructions into root-level instruction files (e.g., `CLAUDE.md`, `AGENTS.md` at project root), directing AI assistants to read `spectr/AGENTS.md` for full instructions rather than duplicating the entire content.

#### Scenario: Init creates instruction file with pointer

- WHEN user runs `spectr init` and selects an AI tool (e.g., Claude Code)
- THEN the root-level instruction file (e.g., `CLAUDE.md`) contains a short pointer between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers
- AND the pointer directs AI assistants to read `spectr/AGENTS.md` when handling proposals, specs, or changes
- AND the full instructions remain only in `spectr/AGENTS.md`

#### Scenario: Update refreshes instruction file with pointer

- WHEN user runs `spectr init` on an already-initialized project
- THEN the root-level instruction files are updated with the short pointer content
- AND the `spectr/AGENTS.md` file retains the full instructions

#### Scenario: Pointer content is concise

- WHEN the instruction pointer template is rendered
- THEN the output is less than 20 lines
- AND the output explains when to read `spectr/AGENTS.md` (proposals, specs, changes, planning)
- AND the output does NOT duplicate the full workflow instructions

### Requirement: PR Archive Subcommand Alias
The `spectr pr archive` subcommand SHALL support `a` as a shorthand alias, allowing users to invoke `spectr pr a <id>` as equivalent to `spectr pr archive <id>`.

#### Scenario: User runs spectr pr a shorthand
- WHEN user runs `spectr pr a <change-id>`
- THEN the system executes the archive PR workflow identically to `spectr pr archive`
- AND all flags (`--base`, `--draft`, `--force`, `--dry-run`, `--skip-specs`) work with the alias

#### Scenario: User runs spectr pr a with flags
- WHEN user runs `spectr pr a my-change --draft --force`
- THEN the command behaves identically to `spectr pr archive my-change --draft --force`
- AND a draft PR is created after deleting any existing branch

#### Scenario: Help text shows archive alias
- WHEN user runs `spectr pr --help`
- THEN the help text displays `archive` with its `a` alias
- AND the alias is shown in parentheses or as comma-separated alternatives

### Requirement: PR Branch Naming Convention
The system SHALL use a mode-specific branch naming convention for PR branches that distinguishes between archive and proposal branches based on the subcommand used.

#### Scenario: Archive branch name format
- WHEN user runs `spectr pr archive <change-id>`
- THEN the branch is named `spectr/archive/<change-id>`

#### Scenario: Proposal branch name format
- WHEN user runs `spectr pr new <change-id>`
- THEN the branch is named `spectr/proposal/<change-id>`

#### Scenario: Branch name with special characters
- WHEN change ID contains only valid kebab-case characters
- THEN the branch name is valid for git

#### Scenario: Branch names clearly indicate PR purpose
- WHEN a developer views the branch list
- THEN they can distinguish archive PRs from proposal PRs by the branch prefix
- AND `spectr/archive/*` indicates a completed change being archived
- AND `spectr/proposal/*` indicates a change proposal for review

#### Scenario: Force flag for existing archive branch
- WHEN user runs `spectr pr archive <change-id> --force`
- AND branch `spectr/archive/<change-id>` already exists on remote
- THEN the existing branch is deleted and recreated
- AND the PR workflow proceeds normally

#### Scenario: Force flag for existing proposal branch
- WHEN user runs `spectr pr new <change-id> --force`
- AND branch `spectr/proposal/<change-id>` already exists on remote
- THEN the existing branch is deleted and recreated
- AND the PR workflow proceeds normally

#### Scenario: Archive branch conflict without force
- WHEN user runs `spectr pr archive <change-id>`
- AND branch `spectr/archive/<change-id>` already exists on remote
- AND `--force` flag is NOT provided
- THEN an error is displayed: "branch 'spectr/archive/<change-id>' already exists on remote; use --force to delete"
- AND the command exits with code 1

#### Scenario: Proposal branch conflict without force
- WHEN user runs `spectr pr new <change-id>`
- AND branch `spectr/proposal/<change-id>` already exists on remote
- AND `--force` flag is NOT provided
- THEN an error is displayed: "branch 'spectr/proposal/<change-id>' already exists on remote; use --force to delete"
- AND the command exits with code 1

### Requirement: PR Command Structure
The system SHALL provide a `spectr pr` command with `archive` and `proposal` subcommands for creating pull requests from Spectr changes using git worktree isolation.

#### Scenario: User runs spectr pr without subcommand
- WHEN user runs `spectr pr` without a subcommand
- THEN help text is displayed showing available subcommands (archive, proposal)
- AND the command exits with code 0

#### Scenario: User runs spectr pr archive
- WHEN user runs `spectr pr archive <change-id>`
- THEN the system executes the archive PR workflow
- AND a PR is created with the archived change

#### Scenario: User runs spectr pr proposal
- WHEN user runs `spectr pr proposal <change-id>`
- THEN the system executes the proposal PR workflow
- AND a PR is created with the change proposal copied (not archived)

### Requirement: PR Archive Subcommand
The `spectr pr archive` subcommand SHALL create a pull request containing an archived Spectr change, executing the archive workflow in an isolated git worktree.

#### Scenario: Archive PR workflow execution
- WHEN user runs `spectr pr archive <change-id>`
- THEN the system creates a temporary git worktree on branch `spectr/<change-id>`
- AND executes `spectr archive <change-id> --yes` within the worktree
- AND stages all changes in `spectr/` directory
- AND commits with structured message including archive metadata
- AND pushes the branch to origin
- AND creates a PR using the detected platform's CLI
- AND cleans up the temporary worktree
- AND displays the PR URL on success

#### Scenario: Archive PR with skip-specs flag
- WHEN user runs `spectr pr archive <change-id> --skip-specs`
- THEN the `--skip-specs` flag is passed to the underlying archive command
- AND spec merging is skipped during the archive operation

#### Scenario: Archive PR preserves user working directory
- WHEN user runs `spectr pr archive <change-id>`
- AND user has uncommitted changes in their working directory
- THEN the user's working directory is NOT modified
- AND the archive operation executes only within the isolated worktree
- AND uncommitted changes are NOT included in the PR

### Requirement: PR Proposal Subcommand
The `spectr pr proposal` subcommand SHALL create a pull request containing a Spectr change proposal for review, copying the change to an isolated git worktree without archiving. This command replaces the deprecated `spectr pr new` command.

The renaming from `new` to `proposal` aligns CLI terminology with the `/spectr:proposal` slash command naming convention, creating consistent vocabulary across CLI and IDE integrations.

#### Scenario: Proposal PR workflow execution
- WHEN user runs `spectr pr proposal <change-id>`
- THEN the system creates a temporary git worktree on branch `spectr/<change-id>`
- AND copies the change directory from source to worktree
- AND stages all changes in `spectr/` directory
- AND commits with structured message for proposal review
- AND pushes the branch to origin
- AND creates a PR using the detected platform's CLI
- AND cleans up the temporary worktree
- AND displays the PR URL on success

#### Scenario: Proposal PR does not archive
- WHEN user runs `spectr pr proposal <change-id>`
- THEN the original change remains in `spectr/changes/<change-id>/`
- AND the change is NOT moved to archive
- AND spec merging does NOT occur

#### Scenario: Proposal PR validates change first
- WHEN user runs `spectr pr proposal <change-id>`
- THEN the system runs validation on the change
- AND warnings are displayed if validation issues exist
- AND the PR workflow continues (validation does not block)

#### Scenario: User runs spectr pr without subcommand
- WHEN user runs `spectr pr` without a subcommand
- THEN help text is displayed showing available subcommands (archive, proposal)
- AND the command exits with code 0

#### Scenario: Unique prefix match for PR proposal command
- WHEN user runs `spectr pr proposal refactor`
- AND only one change ID starts with `refactor`
- THEN a resolution message is displayed
- AND the PR workflow proceeds with the resolved ID

### Requirement: PR Common Flags
Both `spectr pr archive` and `spectr pr proposal` subcommands SHALL support common flags for controlling PR creation behavior.

#### Scenario: Base branch flag
- WHEN user provides `--base <branch>` flag
- THEN the PR targets the specified branch instead of auto-detected main/master

#### Scenario: Auto-detect base branch
- WHEN user does NOT provide `--base` flag
- AND `origin/main` exists
- THEN the PR targets `main`

#### Scenario: Fallback to master
- WHEN user does NOT provide `--base` flag
- AND `origin/main` does NOT exist
- AND `origin/master` exists
- THEN the PR targets `master`

#### Scenario: Draft PR flag
- WHEN user provides `--draft` flag
- THEN the PR is created as a draft PR on platforms that support it

#### Scenario: Force flag for existing branch
- WHEN user provides `--force` flag
- AND branch `spectr/<change-id>` already exists on remote
- THEN the existing branch is deleted and recreated
- AND the PR workflow proceeds normally

#### Scenario: Branch conflict without force
- WHEN branch `spectr/<change-id>` already exists on remote
- AND `--force` flag is NOT provided
- THEN an error is displayed: "Branch 'spectr/<change-id>' already exists. Use --force to overwrite."
- AND the command exits with code 1

#### Scenario: Dry run flag
- WHEN user provides `--dry-run` flag
- THEN the system displays what would be done without executing
- AND no git operations are performed
- AND no PR is created

### Requirement: Git Platform Detection
The system SHALL automatically detect the git hosting platform from the origin remote URL and use the appropriate CLI tool for PR creation.

#### Scenario: Detect GitHub platform
- WHEN origin remote URL contains `github.com`
- THEN platform is detected as GitHub
- AND `gh` CLI is used for PR creation

#### Scenario: Detect GitLab platform
- WHEN origin remote URL contains `gitlab.com` or `gitlab`
- THEN platform is detected as GitLab
- AND `glab` CLI is used for MR creation

#### Scenario: Detect Gitea platform
- WHEN origin remote URL contains `gitea` or `forgejo`
- THEN platform is detected as Gitea
- AND `tea` CLI is used for PR creation

#### Scenario: Detect Bitbucket platform
- WHEN origin remote URL contains `bitbucket.org` or `bitbucket`
- THEN platform is detected as Bitbucket
- AND manual PR URL is provided (no CLI automation)

#### Scenario: Unknown platform error
- WHEN origin remote URL does not match any known platform
- THEN an error is displayed with the detected URL
- AND guidance is provided for manual PR creation

#### Scenario: SSH URL format support
- WHEN origin remote uses SSH format (e.g., `git@github.com:org/repo.git`)
- THEN platform detection correctly identifies the host

#### Scenario: HTTPS URL format support
- WHEN origin remote uses HTTPS format (e.g., `https://github.com/org/repo.git`)
- THEN platform detection correctly identifies the host

### Requirement: Platform CLI Availability
The system SHALL verify that the required platform CLI tool is installed and authenticated before attempting PR creation.

#### Scenario: CLI not installed
- WHEN the required CLI tool (gh/glab/tea) is not installed
- THEN an error is displayed: "<tool> CLI is required for <platform> PR creation. Install: <install-url>"
- AND the command exits with code 1

#### Scenario: CLI not authenticated
- WHEN the required CLI tool is installed but not authenticated
- THEN an error is displayed with authentication instructions
- AND the command exits with code 1

### Requirement: Git Worktree Isolation
The PR commands SHALL use git worktrees to provide complete isolation from the user's working directory.

#### Scenario: Worktree created in temp directory
- WHEN PR workflow starts
- THEN a worktree is created in the system temp directory
- AND the worktree path includes a UUID to prevent conflicts

#### Scenario: Worktree based on origin branch
- WHEN worktree is created
- THEN it is based on the remote base branch (origin/main or origin/master)
- AND it does NOT include any local uncommitted changes

#### Scenario: Worktree cleanup on success
- WHEN PR workflow completes successfully
- THEN the temporary worktree is removed
- AND no temporary files remain

#### Scenario: Worktree cleanup on failure
- WHEN PR workflow fails at any stage
- THEN the temporary worktree is still removed
- AND an appropriate error message is displayed

#### Scenario: Git version requirement
- WHEN git version is less than 2.5
- THEN an error is displayed: "Git >= 2.5 required for worktree support. Current version: <version>"
- AND the command exits with code 1

### Requirement: PR Commit Message Format
The system SHALL generate structured commit messages that follow conventional commit format and include relevant metadata.

#### Scenario: Archive commit message format
- WHEN `spectr pr archive` commits changes
- THEN the commit message includes:
  - Title: `spectr(archive): <change-id>`
  - Archive location path
  - Spec operation counts (added, modified, removed, renamed)
  - Attribution: "Generated by: spectr pr archive"

#### Scenario: Proposal commit message format
- WHEN `spectr pr proposal` commits changes
- THEN the commit message includes:
  - Title: `spectr(proposal): <change-id>`
  - Proposal location path
  - Attribution: "Generated by: spectr pr proposal"

### Requirement: PR Body Content
The system SHALL generate PR body content that helps reviewers understand the change.

#### Scenario: Archive PR body content
- WHEN PR is created for archive
- THEN the PR body includes:
  - Summary section with change ID and archive location
  - Spec updates table with operation counts
  - List of updated capabilities
  - Review checklist

#### Scenario: Proposal PR body content
- WHEN PR is created for proposal
- THEN the PR body includes:
  - Summary section with change ID and location
  - List of included files (proposal.md, tasks.md, specs/)
  - Review checklist

### Requirement: PR Branch Naming
The system SHALL use a consistent branch naming convention for PR branches.

#### Scenario: Branch name format
- WHEN PR workflow creates a branch
- THEN the branch is named `spectr/<change-id>`

#### Scenario: Branch name with special characters
- WHEN change ID contains only valid kebab-case characters
- THEN the branch name is valid for git

### Requirement: PR Error Handling
The system SHALL provide clear error messages and guidance when PR creation fails.

#### Scenario: Not in git repository
- WHEN user runs `spectr pr` from outside a git repository
- THEN an error is displayed: "Not in a git repository"
- AND the command exits with code 1

#### Scenario: No origin remote
- WHEN user runs `spectr pr` and no origin remote exists
- THEN an error is displayed: "No 'origin' remote configured"
- AND guidance is provided to add a remote

#### Scenario: Change does not exist
- WHEN user runs `spectr pr <subcommand> <change-id>`
- AND the change does not exist
- THEN an error is displayed: "Change '<change-id>' not found"
- AND the command exits with code 1

#### Scenario: Push failure
- WHEN git push fails (e.g., network error)
- THEN an error is displayed with the git error message
- AND guidance is provided for manual recovery
- AND worktree is still cleaned up

#### Scenario: PR creation failure with pushed branch
- WHEN push succeeds but PR creation fails
- THEN an error is displayed with the PR CLI error
- AND a message indicates: "Branch was pushed. Create PR manually or retry."
- AND worktree is still cleaned up

### Requirement: Partial Change ID Resolution for PR Commands
The `spectr pr` subcommands SHALL support intelligent partial ID matching consistent with the archive command's resolution algorithm.

#### Scenario: Exact ID match for PR commands
- WHEN user runs `spectr pr archive exact-change-id`
- AND a change with ID `exact-change-id` exists
- THEN the PR workflow proceeds with that change

#### Scenario: Unique prefix match for PR commands
- WHEN user runs `spectr pr proposal refactor`
- AND only one change ID starts with `refactor`
- THEN a resolution message is displayed
- AND the PR workflow proceeds with the resolved ID

### Requirement: PR Proposal Interactive Selection Filters Unmerged Changes

The `spectr pr proposal` command's interactive selection mode SHALL only display changes that do not already exist on the target branch (main/master), ensuring users only see changes that genuinely need proposal PRs.

#### Scenario: Interactive list excludes changes on main

- WHEN user runs `spectr pr proposal` without a change ID argument
- AND some changes in `spectr/changes/` already exist on `origin/main`
- THEN only changes NOT present on `origin/main` are displayed in the interactive list
- AND changes that exist on main are filtered out before display

#### Scenario: All changes already on main

- WHEN user runs `spectr pr proposal` without a change ID argument
- AND all active changes already exist on `origin/main`
- THEN a message is displayed: "No unmerged proposals found. All changes already exist on main."
- AND the command exits gracefully without entering interactive mode

#### Scenario: No changes exist at all

- WHEN user runs `spectr pr proposal` without a change ID argument
- AND no changes exist in `spectr/changes/`
- THEN a message is displayed: "No changes found."
- AND the command exits gracefully

#### Scenario: Explicit change ID bypasses filter

- WHEN user runs `spectr pr proposal <change-id>` with an explicit argument
- THEN the filter is NOT applied
- AND the command proceeds with the specified change ID
- AND existing behavior is preserved for direct invocation

#### Scenario: Archive command unaffected

- WHEN user runs `spectr pr archive` without a change ID argument
- THEN all active changes are displayed in the interactive list
- AND no filtering based on main branch presence is applied
- AND existing archive behavior is preserved

#### Scenario: Detection uses git ls-tree

- WHEN the system checks if a change exists on main
- THEN it uses `git ls-tree` to check if `spectr/changes/<change-id>` path exists on `origin/main`
- AND the check is performed before displaying the interactive list
- AND fetch is performed first to ensure refs are current

#### Scenario: Custom base branch respected

- WHEN user runs `spectr pr proposal --base develop` without a change ID
- THEN the filter checks against `origin/develop` instead of `origin/main`
- AND only changes not present on `origin/develop` are displayed

### Requirement: Template Path Variables

The template rendering system SHALL support dynamic path variables for all directory and file references, allowing template content to be decoupled from specific path names while maintaining backward-compatible defaults.

The `TemplateContext` struct SHALL provide the following fields with default values:
- `BaseDir`: The root Spectr directory name (default: `spectr`)
- `SpecsDir`: The specifications directory path (default: `spectr/specs`)
- `ChangesDir`: The changes directory path (default: `spectr/changes`)
- `ProjectFile`: The project configuration file path (default: `spectr/project.md`)
- `AgentsFile`: The agents instruction file path (default: `spectr/AGENTS.md`)

#### Scenario: Templates use path variables instead of hardcoded strings

- WHEN a template file contains path references
- THEN the path SHALL be expressed using Go template syntax (e.g., `{{ .BaseDir }}`, `{{ .SpecsDir }}`)
- AND hardcoded `spectr/` strings SHALL NOT appear in template files for path references
- AND the rendered output SHALL contain the actual path values from the context

#### Scenario: Default context produces backward-compatible output

- WHEN `DefaultTemplateContext()` is used for rendering
- THEN the rendered output SHALL be identical to the previous hardcoded output
- AND all path references SHALL resolve to `spectr/`, `spectr/specs/`, `spectr/changes/`, etc.

#### Scenario: Template manager methods accept context parameter

- WHEN `RenderAgents()`, `RenderInstructionPointer()`, or `RenderSlashCommand()` is called
- THEN the method SHALL accept a `TemplateContext` parameter
- AND the context values SHALL be available within the template

#### Scenario: All template files use consistent variable names

- WHEN any template file references a Spectr path
- THEN it SHALL use the standardized variable names (`BaseDir`, `SpecsDir`, `ChangesDir`, `ProjectFile`, `AgentsFile`)
- AND variable names SHALL be consistent across all template files

### Requirement: Copy Populate Context Prompt in Init Next Steps

The Next Steps completion screen in the interactive initialization wizard SHALL provide a hotkey to copy the "populate project context" prompt (step 1) to the system clipboard.

#### Scenario: Copy prompt with 'c' hotkey

- WHEN the user is on the Next Steps completion screen after successful initialization
- AND the user presses the 'c' key
- THEN the raw prompt text (without surrounding quotes) "Review spectr/project.md and help me fill in our project's tech stack, conventions, and description. Ask me questions to understand the codebase." is copied to the clipboard
- AND the wizard exits immediately and returns to the shell
- AND no success message is displayed (silent exit, consistent with list mode Enter behavior)

#### Scenario: Clipboard copy failure handling

- WHEN the user presses 'c' to copy the prompt
- AND the clipboard operation fails
- THEN an error message is displayed (e.g., "Failed to copy to clipboard: [error]")
- AND the wizard does NOT exit
- AND the user can retry the copy operation or press 'q' to quit manually

#### Scenario: Help text shows copy hotkey

- WHEN the Next Steps completion screen is displayed after successful initialization
- THEN the footer help text SHALL include the 'c' hotkey
- AND the help text format is: "c: copy step 1 | q: quit" or "c: copy prompt | q: quit"
- AND the hotkey is clearly described

#### Scenario: Copy hotkey only on success screen

- WHEN initialization fails and the error screen is displayed
- THEN the 'c' hotkey is NOT active
- AND the help text does NOT mention the copy hotkey
- AND only quit controls are available

#### Scenario: Clipboard uses OSC 52 fallback

- WHEN the user presses 'c' in an SSH/remote session without native clipboard access
- THEN the copy operation uses OSC 52 escape sequences as fallback
- AND the operation is considered successful (OSC 52 does not report errors)
- AND the success message is displayed

### Requirement: PR Hotkey in Interactive Changes List Mode

The interactive changes list mode SHALL provide a `P` (Shift+P) hotkey that exits the TUI and enters the `spectr pr` workflow for the selected change, allowing users to create pull requests without manually copying the change ID.

#### Scenario: User presses Shift+P to enter PR mode

- WHEN user is in interactive changes mode (`spectr list -I`)
- AND user presses the `P` key (Shift+P) on a selected change
- THEN the interactive mode exits
- AND the system enters PR mode for the selected change ID
- AND the user is prompted to select PR type (archive or proposal)

#### Scenario: PR hotkey not available in specs mode

- WHEN user is in interactive specs mode (`spectr list --specs -I`)
- AND user presses `P` key
- THEN the key press is ignored (no action taken)
- AND the help text does NOT show `P: pr` option

#### Scenario: PR hotkey not available in unified mode

- WHEN user is in unified interactive mode (`spectr list --all -I`)
- AND user presses `P` key
- THEN the key press is ignored (no action taken)
- AND the help text does NOT show `P: pr` option
- AND this avoids confusion when a spec row is selected

#### Scenario: Help text shows PR hotkey in changes mode

- WHEN user presses `?` in changes mode
- THEN the help text includes `P: pr` in the controls line
- AND the hotkey appears alongside other change-specific hotkeys (e, a)

#### Scenario: PR workflow integration

- WHEN the PR hotkey triggers the PR workflow
- THEN the workflow uses the same code path as `spectr pr`
- AND the selected change ID is passed to the PR workflow
- AND the user can select between archive and proposal modes

### Requirement: VHS Demo for PR Hotkey

The system SHALL provide a VHS tape demonstrating the Shift+P hotkey utility in the interactive list TUI.

#### Scenario: Developer finds PR hotkey demo

- WHEN a developer reviews the VHS tape files in `assets/vhs/`
- THEN they SHALL find `pr-hotkey.tape` demonstrating the PR hotkey workflow

#### Scenario: User sees PR hotkey demo

- WHEN a user views the PR hotkey demo GIF
- THEN they SHALL see the interactive list being invoked
- AND they SHALL see the `P` key being pressed
- AND they SHALL see the PR mode being entered for the selected change

### Requirement: PR Proposal Local Change Cleanup Confirmation

After a successful `spectr pr proposal` command that creates a pull request, the system SHALL prompt the user with a Bubbletea TUI confirmation menu asking whether to remove the local change proposal directory from `spectr/changes/`.

This prompt helps users maintain a clean working directory by offering an opportunity to remove proposals that have been submitted for review, while defaulting to "No" for safety. The menu uses arrow key navigation and styled rendering consistent with other spectr interactive modes.

#### Scenario: Successful PR proposal triggers cleanup prompt

- WHEN user runs `spectr pr proposal <change-id>` and PR creation succeeds
- AND the PR URL is displayed to the user
- THEN the system displays a Bubbletea TUI menu: "Remove local change proposal from spectr/changes/?"
- AND the menu shows two options: "No, keep it" and "Yes, remove it"
- AND the default selection is "No, keep it" (cursor starts on this option)
- AND the menu supports arrow key navigation (↑/↓) and Enter to confirm

#### Scenario: User confirms cleanup via TUI

- WHEN the cleanup TUI menu is displayed
- AND user navigates to "Yes, remove it" and presses Enter
- THEN the system removes the directory `spectr/changes/<change-id>/`
- AND displays confirmation: "Removed local change: <change-id>"

#### Scenario: User declines cleanup via TUI

- WHEN the cleanup TUI menu is displayed
- AND user presses Enter on the default "No, keep it" option
- THEN the system keeps the local change directory
- AND displays: "Local change kept: spectr/changes/<change-id>/"

#### Scenario: User cancels cleanup menu

- WHEN the cleanup TUI menu is displayed
- AND user presses 'q' or Ctrl+C
- THEN the system keeps the local change directory (same as declining)
- AND the command exits successfully

#### Scenario: Non-interactive mode skips prompt

- WHEN user runs `spectr pr proposal <change-id> --yes`
- AND PR creation succeeds
- THEN the cleanup prompt is NOT displayed
- AND the local change directory is kept (safe default)
- AND the command exits successfully

#### Scenario: Cleanup for archive mode

- WHEN user runs `spectr pr archive <change-id>`
- AND PR creation succeeds
- THEN the system displays: "Cleaning up local change directory: spectr/changes/<change-id>/"
- AND the local change directory is removed
- AND the change is archived in the worktree (pulled when PR merges)

#### Scenario: PR creation fails

- WHEN user runs `spectr pr proposal <change-id>`
- AND PR creation fails at any step
- THEN the cleanup prompt is NOT displayed
- AND the local change directory remains intact

#### Scenario: Cleanup removal error handling

- WHEN the user confirms cleanup
- AND removal of the change directory fails (e.g., permission error)
- THEN display a warning: "Warning: Failed to remove change directory: <error>"
- AND the command still exits successfully (non-fatal error)

### Requirement: CI Workflow Setup Option in Init Wizard Review Step
The initialization wizard's Review step SHALL include an optional checkbox to create a GitHub Actions workflow file (`.github/workflows/spectr-ci.yml`) for automated Spectr validation during CI/CD. This option is presented alongside the tool selection summary, keeping the wizard flow quick without adding a separate step.

#### Scenario: CI option displayed in Review step
- WHEN user completes tool selection and proceeds to the Review step
- THEN a "Spectr CI Validation" checkbox option is displayed after the tool summary
- AND the option appears before the creation plan section
- AND a description explains: "Validate specs automatically on push and pull requests"

#### Scenario: CI option detects existing workflow
- WHEN user runs `spectr init` on a project that already has `.github/workflows/spectr-ci.yml`
- AND user reaches the Review step
- THEN the "Spectr CI Validation" option shows a "(configured)" indicator
- AND the option is pre-selected by default
- AND selecting it will update the existing workflow file

#### Scenario: CI option not pre-selected on fresh projects
- WHEN user runs `spectr init` on a project without `.github/workflows/spectr-ci.yml`
- AND user reaches the Review step
- THEN the "Spectr CI Validation" option is NOT pre-selected by default
- AND the user must explicitly select it to enable CI workflow creation

#### Scenario: User toggles CI option in Review step
- WHEN user is on the Review step
- AND user presses Space
- THEN the CI workflow checkbox toggles between selected and unselected
- AND the creation plan updates to reflect the change
- AND the visual state updates immediately

#### Scenario: CI workflow created when selected
- WHEN user selects the "Spectr CI Validation" option in Review
- AND user presses Enter to proceed with initialization
- THEN the system creates `.github/workflows/` directory if it doesn't exist
- AND the system creates `.github/workflows/spectr-ci.yml` with the Spectr validation workflow
- AND the workflow file is tracked in the execution result as created or updated

#### Scenario: CI workflow not created when unselected
- WHEN user does NOT select the "Spectr CI Validation" option
- AND user proceeds with initialization
- THEN no `.github/workflows/spectr-ci.yml` file is created
- AND any existing `.github/workflows/spectr-ci.yml` file is left unchanged

#### Scenario: CI workflow content uses pinned action version
- WHEN the CI workflow file is created
- THEN the workflow contains a single `spectr-validate` job
- AND the workflow uses `connerohnesorge/spectr-action@v0.0.2` (pinned version)
- AND the workflow triggers on push to `main` branch only
- AND the workflow triggers on pull requests to all branches
- AND the workflow uses `fetch-depth: 0` for full git history
- AND the workflow includes concurrency management to cancel in-progress runs
- AND the workflow runs on `ubuntu-latest`

#### Scenario: Creation plan shows CI workflow when enabled
- WHEN user has selected the "Spectr CI Validation" option
- THEN the creation plan section shows `.github/workflows/spectr-ci.yml`
- AND the file is listed with the tool configurations

#### Scenario: Creation plan hides CI workflow when disabled
- WHEN user has NOT selected the "Spectr CI Validation" option
- THEN the creation plan does NOT mention `.github/workflows/spectr-ci.yml`

#### Scenario: Completion screen shows CI workflow file
- WHEN the CI workflow file is successfully created
- THEN the completion screen lists `.github/workflows/spectr-ci.yml` in created or updated files
- AND the file path is displayed with the appropriate icon

#### Scenario: Non-interactive mode supports CI workflow flag
- WHEN user runs `spectr init --non-interactive --ci-workflow`
- THEN the CI workflow file is created without TUI interaction
- AND the workflow file content matches the interactive mode output

#### Scenario: Non-interactive mode without CI flag skips workflow
- WHEN user runs `spectr init --non-interactive` without `--ci-workflow`
- THEN no CI workflow file is created
- AND existing workflow files are not modified

#### Scenario: Review step help text includes Space for toggle
- WHEN user is on the Review step
- THEN the help text shows: "Space: Toggle CI  Enter: Initialize  Backspace: Back  q: Quit"

### Requirement: PR Remove Subcommand

The `spectr pr rm` subcommand SHALL create a pull request that removes a change directory from the repository, using the same git worktree isolation as other PR subcommands.

The command supports aliases `r` and `remove` for convenience.

#### Scenario: User runs spectr pr rm with change ID

- WHEN user runs `spectr pr rm <change-id>`
- THEN the system creates a temporary git worktree on branch `spectr/remove/<change-id>`
- AND removes the change directory from `spectr/changes/<change-id>` in the worktree
- AND stages the deletion
- AND commits with a structured message indicating removal
- AND pushes the branch to origin
- AND creates a PR using the detected platform's CLI
- AND cleans up the temporary worktree
- AND displays the PR URL on success
- AND removes the local change directory after successful PR creation

#### Scenario: User runs spectr pr rm without change ID

- WHEN user runs `spectr pr rm` without a change ID argument
- THEN an interactive table is displayed showing available changes
- AND user can navigate and select a change
- AND the remove workflow proceeds with the selected change ID

#### Scenario: User runs spectr pr r shorthand

- WHEN user runs `spectr pr r <change-id>`
- THEN the system executes the remove PR workflow identically to `spectr pr rm`
- AND all flags work with the alias

#### Scenario: Remove PR with draft flag

- WHEN user runs `spectr pr rm <change-id> --draft`
- THEN the PR is created as a draft PR on platforms that support it

#### Scenario: Remove PR with force flag

- WHEN user runs `spectr pr rm <change-id> --force`
- AND branch `spectr/remove/<change-id>` already exists on remote
- THEN the existing branch is deleted and recreated
- AND the PR workflow proceeds normally

#### Scenario: Remove branch conflict without force

- WHEN user runs `spectr pr rm <change-id>`
- AND branch `spectr/remove/<change-id>` already exists on remote
- AND `--force` flag is NOT provided
- THEN an error is displayed: "branch 'spectr/remove/<change-id>' already exists on remote; use --force to delete"
- AND the command exits with code 1

#### Scenario: Remove PR with dry-run flag

- WHEN user runs `spectr pr rm <change-id> --dry-run`
- THEN the system displays what would be done without executing
- AND no git operations are performed
- AND no PR is created
- AND no local cleanup is performed

#### Scenario: Remove PR with base branch flag

- WHEN user runs `spectr pr rm <change-id> --base develop`
- THEN the PR targets the `develop` branch instead of auto-detected main/master

#### Scenario: Change does not exist

- WHEN user runs `spectr pr rm <change-id>`
- AND the change does not exist in `spectr/changes/`
- THEN an error is displayed: "change '<change-id>' not found in spectr/changes/"
- AND the command exits with code 1

#### Scenario: Remove cleans up local change directory

- WHEN user runs `spectr pr rm <change-id>`
- AND PR creation succeeds
- THEN the system displays: "Cleaning up local change directory: spectr/changes/<change-id>/"
- AND the local change directory is removed including all files (tracked and untracked)

#### Scenario: Partial ID resolution for remove command

- WHEN user runs `spectr pr rm refactor`
- AND only one change ID starts with `refactor`
- THEN a resolution message is displayed
- AND the PR workflow proceeds with the resolved ID

### Requirement: Remove PR Branch Naming

The `spectr pr rm` command SHALL use the branch naming pattern `spectr/remove/<change-id>` to clearly indicate the PR's purpose.

#### Scenario: Remove branch name format

- WHEN user runs `spectr pr rm <change-id>`
- THEN the branch is named `spectr/remove/<change-id>`

#### Scenario: Remove branch distinguishable from archive and proposal

- WHEN a developer views the branch list
- THEN they can distinguish remove PRs from archive and proposal PRs by the branch prefix
- AND `spectr/remove/*` indicates a change removal PR
- AND `spectr/archive/*` indicates a completed change being archived
- AND `spectr/proposal/*` indicates a change proposal for review

### Requirement: Remove PR Commit Message Format

The `spectr pr rm` command SHALL generate a structured commit message that clearly indicates the removal.

#### Scenario: Remove commit message content

- WHEN `spectr pr rm` commits changes
- THEN the commit message includes:
  - Title: `spectr(remove): <change-id>`
  - Removed path: `spectr/changes/<change-id>`
  - Attribution: "Generated by: spectr pr rm"

### Requirement: Remove PR Body Content

The `spectr pr rm` command SHALL generate PR body content that explains the removal.

#### Scenario: Remove PR body content

- WHEN PR is created for removal
- THEN the PR body includes:
  - Summary section with change ID and removal context
  - The removed path
  - Review checklist for confirming removal is intentional

### Requirement: Responsive Table Column Layout
The interactive TUI table views SHALL detect terminal width and dynamically adjust column visibility and widths to ensure readable display across different screen sizes.

#### Scenario: Full width terminal displays all columns
- WHEN user runs `spectr list -I` on a terminal with 110+ columns
- THEN all columns are displayed at their default widths
- AND for changes view: ID (30), Title (40), Deltas (10), Tasks (15) are shown
- AND for specs view: ID (35), Title (45), Requirements (15) are shown
- AND for unified view: ID (30), Type (8), Title (40), Details (20) are shown

#### Scenario: Medium width terminal narrows Title column
- WHEN user runs `spectr list -I` on a terminal with 90-109 columns
- THEN the Title column width is reduced proportionally
- AND title truncation threshold is reduced to match narrower column
- AND all columns remain visible

#### Scenario: Narrow width terminal hides low-priority columns
- WHEN user runs `spectr list -I` on a terminal with 70-89 columns
- THEN the lowest-priority columns are hidden
- AND for changes view: Tasks column is hidden, Deltas may be narrowed
- AND for specs view: Requirements column width is reduced
- AND for unified view: Details column is hidden
- AND remaining columns are adjusted to fit available width

#### Scenario: Minimal width terminal shows essential columns only
- WHEN user runs `spectr list -I` on a terminal with fewer than 70 columns
- THEN only ID and Title columns are displayed
- AND title truncation is aggressive to fit available space
- AND help text indicates some columns are hidden

### Requirement: Dynamic Terminal Resize Handling
The interactive TUI SHALL respond to terminal resize events by recalculating and rebuilding the table layout without losing user state.

#### Scenario: Terminal resized wider during session
- WHEN user is in interactive mode and terminal width increases
- THEN table columns are recalculated for the new width
- AND previously hidden columns may become visible
- AND cursor position is preserved on the same item
- AND search filter state is preserved if active

#### Scenario: Terminal resized narrower during session
- WHEN user is in interactive mode and terminal width decreases
- THEN table columns are recalculated for the new width
- AND low-priority columns are hidden as needed
- AND cursor position is preserved on the same item
- AND the view does not overflow horizontally

#### Scenario: Resize does not interrupt search mode
- WHEN user is in search mode and terminal is resized
- THEN search input remains active
- AND filtered results are preserved
- AND table layout adapts to new width

### Requirement: Column Priority System
Each table view SHALL define column priorities to determine which columns are shown at each width breakpoint.

#### Scenario: Changes view column priorities
- WHEN calculating responsive columns for changes view
- THEN ID has highest priority (always shown)
- AND Title has second priority (always shown, width adjustable)
- AND Deltas has third priority (hidden below 80 columns)
- AND Tasks has lowest priority (hidden below 90 columns)

#### Scenario: Specs view column priorities
- WHEN calculating responsive columns for specs view
- THEN ID has highest priority (always shown)
- AND Title has second priority (always shown, width adjustable)
- AND Requirements has lowest priority (width reduced or hidden below 70 columns)

#### Scenario: Unified view column priorities
- WHEN calculating responsive columns for unified view
- THEN ID has highest priority (always shown)
- AND Type has second priority (always shown at fixed 8-character width)
- AND Title has third priority (width adjustable)
- AND Details has lowest priority (hidden below 90 columns)

### Requirement: Provider Search in Init Wizard

The initialization wizard's tool selection step SHALL provide a `/` hotkey that activates a text search mode, allowing users to filter the displayed provider list by typing a search query that matches against provider names.

#### Scenario: User presses '/' to enter search mode

- WHEN user is on the tool selection step of the init wizard (`StepSelect`)
- AND user presses the '/' key
- THEN search mode is activated
- AND a text input field is displayed below the provider list
- AND the cursor is placed in the text input field
- AND the user can type a search query

#### Scenario: Search filters providers in real-time

- WHEN search mode is active
- AND user types characters into the search input
- THEN the provider list is filtered in real-time
- AND only providers whose name contains the search query (case-insensitive) are displayed
- AND the cursor moves to the first matching provider if current selection is filtered out

#### Scenario: Search with no matches shows empty list

- WHEN search mode is active
- AND user types a query that matches no providers
- THEN the provider list displays no items
- AND a message indicates no matches found (e.g., "No providers match 'xyz'")

#### Scenario: User presses Escape to exit search mode

- WHEN search mode is active
- AND user presses the Escape key
- THEN search mode is deactivated
- AND the search query is cleared
- AND all providers are displayed again in the list
- AND the text input field is hidden

#### Scenario: Selection preserved during filtering

- WHEN search mode is active
- AND user has previously selected providers (checked checkboxes)
- AND user types a query that filters out some selected providers
- THEN the selection state of filtered-out providers is preserved
- AND when search is cleared, previously selected providers remain selected

#### Scenario: Navigation works while searching

- WHEN search mode is active
- AND filtered results are displayed
- THEN arrow key navigation (up/down, j/k) moves through filtered rows
- AND space key toggles selection on the currently highlighted filtered provider
- AND Enter key proceeds to the Review step with all selections (including filtered-out ones)

#### Scenario: Help text shows search hotkey

- WHEN the tool selection step is displayed and search mode is NOT active
- THEN the help text includes '/: search' in the controls line
- AND the search hotkey is shown alongside existing controls (navigate, toggle, all, none, enter, quit)

#### Scenario: Search mode visual indicator

- WHEN search mode is active
- THEN the search input field is visually distinct (styled text input)
- AND the current search query is visible in the input field
- AND the help text updates to show 'Esc: exit search' instead of '/: search'

#### Scenario: Cursor adjustment on filter change

- WHEN search mode is active
- AND the user types additional characters that reduce the filtered list
- AND the current cursor position is beyond the new list length
- THEN the cursor is adjusted to the last valid position in the filtered list
- AND the cursor does not go out of bounds

### Requirement: Stdout Output Mode for Interactive List
The `spectr list` command SHALL support a `--stdout` flag that, when combined with interactive mode (`-I`), outputs the selected item ID to stdout instead of copying it to the system clipboard.

#### Scenario: User runs list with --stdout and -I flags
- WHEN user runs `spectr list -I --stdout`
- AND user navigates to a row and presses Enter
- THEN the selected ID is printed to stdout (just the ID, no formatting)
- AND no clipboard operation is performed
- AND the command exits with code 0

#### Scenario: Stdout mode with changes
- WHEN user runs `spectr list -I --stdout` (changes mode)
- AND user selects a change and presses Enter
- THEN only the change ID is printed to stdout (e.g., `add-feature`)
- AND no "Copied:" prefix or other formatting is included

#### Scenario: Stdout mode with specs
- WHEN user runs `spectr list --specs -I --stdout`
- AND user selects a spec and presses Enter
- THEN only the spec ID is printed to stdout (e.g., `cli-interface`)
- AND no "Copied:" prefix or other formatting is included

#### Scenario: Stdout mode with unified view
- WHEN user runs `spectr list --all -I --stdout`
- AND user selects an item and presses Enter
- THEN only the item ID is printed to stdout
- AND no "Copied:" prefix or other formatting is included

#### Scenario: Stdout flag requires interactive mode
- WHEN user runs `spectr list --stdout` without `-I`
- THEN an error is displayed: "cannot use --stdout without --interactive (-I)"
- AND the command exits with code 1

#### Scenario: Stdout flag mutually exclusive with JSON
- WHEN user runs `spectr list -I --stdout --json`
- THEN an error is displayed: "cannot use --stdout with --json"
- AND the command exits with code 1

#### Scenario: Stdout mode enables piping
- WHEN user runs `spectr list -I --stdout | xargs echo`
- AND user selects an item
- THEN the pipeline receives the clean ID string
- AND the downstream command processes the ID correctly

#### Scenario: User quits without selection in stdout mode
- WHEN user runs `spectr list -I --stdout`
- AND user presses 'q' or Ctrl+C without selecting
- THEN nothing is printed to stdout
- AND the command exits with code 0

#### Scenario: Stdout mode help text
- WHEN user runs `spectr list --help`
- THEN the help text shows `--stdout` flag
- AND the description explains it outputs to stdout instead of clipboard
- AND the help indicates it requires `-I` flag

### Requirement: JSONC Comment Parsing
The system SHALL support reading JSONC files by stripping comments before JSON parsing.

#### Scenario: Strip line comments
- WHEN reading a JSONC file containing `//` line comments
- THEN the parser SHALL remove all text from `//` to end of line
- AND the resulting JSON SHALL parse correctly

#### Scenario: Strip block comments
- WHEN reading a JSONC file containing `/* */` block comments
- THEN the parser SHALL remove all text between `/*` and `*/`
- AND the resulting JSON SHALL parse correctly

#### Scenario: Preserve comments in strings
- WHEN reading a JSONC file with `//` or `/*` inside a JSON string value
- THEN the parser SHALL NOT treat the text as a comment
- AND the string value SHALL be preserved intact

#### Scenario: Parse plain JSON
- WHEN reading a JSON file without comments
- THEN the parser SHALL parse it successfully
- AND no comment stripping side effects SHALL occur

### Requirement: Track Command
The CLI SHALL provide a `track` command that watches task status changes and automatically commits related changes.

#### Scenario: Track with change ID
- **WHEN** user runs `spectr track <change-id>`
- **THEN** the system watches tasks.json for the specified change
- **AND** displays current task status (X/Y completed)
- **AND** runs until all tasks are complete or interrupted

#### Scenario: Interactive track selection
- **WHEN** user runs `spectr track` without specifying a change ID
- **THEN** the system displays a list of active changes with tasks.json
- **AND** prompts for selection

#### Scenario: Auto-commit on task completion
- **WHEN** a task status changes to "completed" in tasks.json
- **THEN** the system detects modified files via git status
- **AND** stages all modified files except tasks.json, tasks.jsonc, tasks.md
- **AND** creates a commit with message format: `spectr(<change-id>): complete task <task-id>`
- **AND** includes footer: `[Automated by spectr track]`

#### Scenario: Auto-commit on task start
- **WHEN** a task status changes to "in_progress" in tasks.json
- **THEN** the system detects modified files via git status
- **AND** stages all modified files except tasks.json, tasks.jsonc, tasks.md
- **AND** creates a commit with message format: `spectr(<change-id>): start task <task-id>`
- **AND** includes footer: `[Automated by spectr track]`

#### Scenario: No files to commit warning
- **WHEN** a task status changes but no files have been modified (excluding task files)
- **THEN** the system prints a warning: "No files to commit for task <task-id>"
- **AND** continues watching for more task changes

#### Scenario: Git commit failure stops tracking
- **WHEN** a git commit operation fails (e.g., merge conflict, hook rejection)
- **THEN** the system displays the git error message
- **AND** stops tracking immediately
- **AND** exits with non-zero status code

#### Scenario: Graceful interruption
- **WHEN** user presses Ctrl+C during tracking
- **THEN** the system stops watching and exits cleanly
- **AND** displays "Tracking stopped" message

#### Scenario: All tasks already complete
- **WHEN** user runs `spectr track <change-id>` and all tasks are already completed
- **THEN** the system displays a message indicating all tasks are complete
- **AND** exits without starting the watch loop

### Requirement: Track Command Flags
The track command SHALL support flags for controlling behavior.

#### Scenario: No-interactive flag disables prompts
- **WHEN** user provides the `--no-interactive` flag
- **AND** no change-id is provided
- **THEN** the system displays usage error instead of prompting for selection
