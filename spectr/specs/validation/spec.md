# Validation Specification

## Purpose

The validation system ensures that spec files and change proposals conform to Spectr conventions, providing structural correctness checks, helpful error messages, and bulk validation capabilities for CI/CD integration.

## Requirements

### Requirement: Spec File Validation
The validation system SHALL validate spec files for structural correctness and adherence to Spectr conventions.

#### Scenario: Valid spec with all required sections
- **WHEN** a spec file contains a Requirements section with properly formatted requirements and scenarios
- **THEN** validation SHALL pass with no errors
- **AND** the validation report SHALL indicate valid=true

#### Scenario: Missing Requirements section
- **WHEN** a spec file lacks a "## Requirements" section
- **THEN** validation SHALL fail with an ERROR level issue
- **AND** the error message SHALL provide example of correct structure

#### Scenario: Requirement without scenarios
- **WHEN** a requirement exists without any "#### Scenario:" subsections
- **THEN** validation SHALL report a WARNING level issue
- **AND** in strict mode validation SHALL fail (valid=false)
- **AND** the warning SHALL include example scenario format

#### Scenario: Requirement missing SHALL or MUST
- **WHEN** a requirement text does not contain "SHALL" or "MUST" keywords
- **THEN** validation SHALL report a WARNING level issue
- **AND** the message SHALL suggest using normative language

#### Scenario: Incorrect scenario format
- **WHEN** scenarios use formats other than "#### Scenario:" (e.g., bullets or bold text)
- **THEN** validation SHALL report an ERROR
- **AND** the message SHALL show the correct "#### Scenario:" header format

#### Scenario: Parsing uses markdown package
- **WHEN** the validation system parses spec or delta files
- **THEN** it SHALL use the `internal/markdown/` package for AST-based parsing
- **AND** it SHALL NOT define local regex patterns for structural markdown elements
- **AND** it SHALL use the visitor pattern and query functions from the markdown package

### Requirement: Change Delta Validation
The validation system SHALL validate change delta specs for structural correctness and delta operation validity.

#### Scenario: Valid change with deltas
- **WHEN** a change directory contains specs with proper ADDED/MODIFIED/REMOVED/RENAMED sections
- **THEN** validation SHALL pass with no errors
- **AND** each delta requirement SHALL be counted toward the total

#### Scenario: Change with no deltas
- **WHEN** a change directory has no specs/ subdirectory or no delta sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL explain that at least one delta is required
- **AND** remediation guidance SHALL explain the delta header format

#### Scenario: Delta sections present but empty
- **WHEN** delta sections exist (## ADDED Requirements) but contain no requirement entries
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which sections are empty
- **AND** guidance SHALL explain requirement block format

#### Scenario: ADDED requirement without scenario
- **WHEN** an ADDED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate which requirement is missing scenarios

#### Scenario: MODIFIED requirement without scenario
- **WHEN** a MODIFIED requirement lacks a "#### Scenario:" block
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL require at least one scenario for MODIFIED requirements

#### Scenario: Duplicate requirement in same section
- **WHEN** two requirements with the same normalized name appear in the same delta section
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL identify the duplicate requirement name

#### Scenario: Cross-section conflicts
- **WHEN** a requirement appears in both ADDED and MODIFIED sections
- **THEN** validation SHALL fail with an ERROR
- **AND** the message SHALL indicate the conflicting requirement and sections

#### Scenario: RENAMED requirement validation
- **WHEN** a RENAMED section contains well-formed "FROM: X TO: Y" pairs
- **THEN** validation SHALL accept the renames using AST-based parsing
- **AND** SHALL check for duplicate FROM or TO entries
- **AND** SHALL error if MODIFIED references the old name instead of new name

#### Scenario: Delta parsing uses markdown package
- **WHEN** the validation system parses delta spec files
- **THEN** it SHALL use the `internal/markdown/` package for AST-based parsing
- **AND** delta type detection SHALL use NodeSection with appropriate Level and Title checks
- **AND** requirement extraction SHALL use NodeRequirement nodes from the AST
- **AND** scenario extraction SHALL use NodeScenario nodes from the AST
- **AND** section content extraction SHALL use query functions like Find and FindFirst

### Requirement: Validation Report Structure
The validation system SHALL produce structured validation reports containing issue details and summary statistics.

#### Scenario: Report with errors and warnings
- **WHEN** validation encounters both ERROR and WARNING level issues
- **THEN** the report SHALL list all issues with level, path, and message
- **AND** the summary SHALL count errors, warnings, and info separately
- **AND** valid SHALL be false if any errors exist

#### Scenario: Report in strict mode
- **WHEN** validation runs in strict mode
- **THEN** the report SHALL treat warnings as failures
- **AND** valid SHALL be false if errors OR warnings exist
- **AND** exit code SHALL be non-zero for warnings in strict mode

#### Scenario: JSON output format
- **WHEN** validation is invoked with --json flag
- **THEN** the output SHALL be valid JSON
- **AND** SHALL include items array with per-item results
- **AND** SHALL include summary with totals and byType breakdowns
- **AND** SHALL include version field for format compatibility

### Requirement: Bulk Validation with Concurrency
The validation system SHALL support validating multiple items in parallel for performance.

#### Scenario: Parallel validation of multiple items
- **WHEN** bulk validation is invoked with multiple specs and changes
- **THEN** validation SHALL process items concurrently using a worker pool
- **AND** concurrency SHALL be configurable via flag or environment variable
- **AND** default concurrency SHALL be 6 workers

#### Scenario: Validation queue management
- **WHEN** the number of items exceeds worker pool size
- **THEN** items SHALL be queued and processed as workers become available
- **AND** progress indicators SHALL update as items complete (if not JSON mode)
- **AND** results SHALL be collected and sorted by item ID

#### Scenario: Error handling in parallel validation
- **WHEN** validation of one item fails with an error (not validation issue, but runtime error)
- **THEN** the error SHALL be captured in the results for that item
- **AND** validation of other items SHALL continue
- **AND** the final exit code SHALL indicate failure

### Requirement: Item Discovery
The validation system SHALL discover specs and changes within the project directory structure.

#### Scenario: Discover active changes
- **WHEN** the system scans the spectr/changes/ directory
- **THEN** it SHALL return all subdirectories except "archive"
- **AND** each subdirectory name SHALL be a change ID

#### Scenario: Discover specs
- **WHEN** the system scans the spectr/specs/ directory
- **THEN** it SHALL return all subdirectories containing a spec.md file
- **AND** each subdirectory name SHALL be a spec ID

#### Scenario: Handle missing directories
- **WHEN** spectr/changes/ or spectr/specs/ does not exist
- **THEN** discovery SHALL return empty list for that category
- **AND** SHALL NOT error on missing directories

### Requirement: Interactive Validation Mode
The validation system SHALL support interactive selection when invoked without arguments in a TTY, using a bubbletea-based TUI with menu-driven navigation and item picker.

#### Scenario: Interactive mode main menu
- **WHEN** validate command is invoked without arguments in an interactive terminal
- **THEN** it SHALL display a menu with options: "Validate All", "Validate All Changes", "Validate All Specs", "Pick Specific Item", "Quit"
- **AND** user SHALL be able to navigate options using arrow keys or j/k
- **AND** user SHALL select an option by pressing Enter
- **AND** selected option SHALL be executed immediately

#### Scenario: Pick specific item with search
- **WHEN** user selects "Pick Specific Item" from the main menu
- **THEN** a searchable list of all changes and specs SHALL be displayed
- **AND** items SHALL be sorted alphabetically with type indicator (change/spec)
- **AND** user SHALL navigate the list with arrow keys or j/k
- **AND** pressing Enter on an item SHALL validate that specific item
- **AND** pressing q or Ctrl+C SHALL return to main menu

#### Scenario: Non-interactive environment detection
- **WHEN** validate command is invoked without arguments in non-interactive environment (CI/CD)
- **THEN** it SHALL print usage hints for non-interactive invocation
- **AND** SHALL exit with code 1
- **AND** SHALL NOT hang waiting for input

#### Scenario: Interactive validation execution
- **WHEN** user selects a validation option in interactive mode
- **THEN** validation SHALL execute using existing validation logic
- **AND** results SHALL be displayed in human-readable format (not JSON)
- **AND** user SHALL see success/failure summary
- **AND** for failures, detailed issues SHALL be shown
- **AND** user SHALL be returned to main menu after viewing results (or exit on quit)

#### Scenario: Consistent styling with other TUIs
- **WHEN** interactive validation TUI is displayed
- **THEN** it SHALL use lipgloss styling consistent with internal/list/interactive.go
- **AND** it SHALL use the same color scheme and formatting patterns
- **AND** help text SHALL be displayed showing available key bindings
- **AND** selected items SHALL be highlighted with cursor style

### Requirement: Helpful Error Messages
The validation system SHALL provide actionable error messages with remediation guidance.

#### Scenario: Error with remediation steps
- **WHEN** validation fails due to missing required content
- **THEN** the error message SHALL explain what is wrong
- **AND** SHALL provide "Next steps" section with concrete actions
- **AND** SHALL include format examples when applicable

#### Scenario: Ambiguous item name
- **WHEN** user provides an item name that matches both a change and a spec
- **THEN** validation SHALL report the ambiguity
- **AND** SHALL suggest using --type flag to disambiguate
- **AND** SHALL show available type options (change, spec)

#### Scenario: Item not found with suggestions
- **WHEN** user provides an item name that does not exist
- **THEN** validation SHALL report item not found
- **AND** SHALL provide nearest match suggestions based on string similarity
- **AND** SHALL limit suggestions to 5 most similar items

### Requirement: Exit Code Conventions
The validation system SHALL use exit codes to indicate success or failure for scripting and CI/CD.

#### Scenario: Successful validation
- **WHEN** all validated items pass without errors (or warnings in strict mode)
- **THEN** the command SHALL exit with code 0

#### Scenario: Validation failures
- **WHEN** any validated item has errors (or warnings in strict mode)
- **THEN** the command SHALL exit with code 1

#### Scenario: Runtime errors
- **WHEN** the command encounters runtime errors (file not found, parse errors)
- **THEN** the command SHALL exit with code 1
- **AND** SHALL print error details to stderr

### Requirement: Helper Functions in Internal Package
The validation system SHALL organize helper functions in internal/validation/ package following clean architecture patterns, with cmd/ serving as a thin command layer.

#### Scenario: Helper functions accessible to internal packages
- **WHEN** validation logic needs to determine item types or format results
- **THEN** helper functions SHALL be available in internal/validation/helpers.go
- **AND** item collection functions SHALL be in internal/validation/items.go
- **AND** formatting functions SHALL be in internal/validation/formatters.go
- **AND** all helpers SHALL be unit tested in their respective test files

#### Scenario: Type determination logic reusable
- **WHEN** any validation component needs to determine if an item is a change or spec
- **THEN** it SHALL use DetermineItemType() from internal/validation/helpers.go
- **AND** the function SHALL return itemTypeInfo with isChange, isSpec, and itemType fields
- **AND** it SHALL handle ambiguous cases (item exists as both change and spec)
- **AND** it SHALL respect explicit --type flag when provided

#### Scenario: Validation item collection reusable
- **WHEN** bulk validation needs to collect items to validate
- **THEN** it SHALL use GetAllItems(), GetChangeItems(), or GetSpecItems() from internal/validation/items.go
- **AND** functions SHALL return []ValidationItem with name, itemType, and path
- **AND** functions SHALL handle missing directories gracefully
- **AND** functions SHALL leverage internal/discovery for ID enumeration

#### Scenario: Result formatting separated from command logic
- **WHEN** validation results need to be displayed
- **THEN** formatting functions SHALL be in internal/validation/formatters.go
- **AND** FormatJSONReport() and FormatHumanReport() SHALL handle single item results
- **AND** FormatBulkJSONResults() and FormatBulkHumanResults() SHALL handle multiple items
- **AND** formatters SHALL accept report data and return formatted strings
- **AND** formatters SHALL not directly write to stdout (return strings instead)

### Requirement: Interactive TUI Architecture
The validation interactive mode SHALL follow the bubbletea model-update-view pattern used in other project TUIs, ensuring consistency and maintainability.

#### Scenario: Bubbletea model structure
- **WHEN** interactive validation TUI is initialized
- **THEN** it SHALL define a model struct implementing tea.Model interface
- **AND** model SHALL contain state for current screen, selected option, validation results
- **AND** model SHALL have Init() method returning initial commands
- **AND** model SHALL have Update(msg tea.Msg) method handling events
- **AND** model SHALL have View() string method rendering current state

#### Scenario: Key binding handling
- **WHEN** user presses keys in interactive validation TUI
- **THEN** arrow up/down and j/k SHALL navigate menu items
- **AND** Enter SHALL select the highlighted option
- **AND** q or Ctrl+C SHALL quit the TUI
- **AND** Esc SHALL go back to previous screen (when in item picker)
- **AND** unrecognized keys SHALL be ignored

#### Scenario: Integration with cmd layer
- **WHEN** cmd/validate.go needs to launch interactive mode
- **THEN** it SHALL call RunInteractiveValidation() from internal/validation/interactive.go
- **AND** function SHALL accept projectPath, validator, and JSON flag
- **AND** function SHALL return error on failure
- **AND** function SHALL handle TTY detection internally
- **AND** cmd layer SHALL only handle program initialization and error printing

### Requirement: Bulk Validation Human Output Formatting
The validation system SHALL produce bulk validation human-readable output with improved spacing, relative paths, file grouping, and color-coded error levels for easier scanning.

#### Scenario: Visual separation between failed items
- **WHEN** bulk validation encounters multiple failed items in human output mode
- **THEN** output SHALL include a blank line between each failed item's error listing
- **AND** passed items SHALL be listed without blank lines between them
- **AND** failed items SHALL be visually distinct from passed items

#### Scenario: Relative path display
- **WHEN** validation issues include file paths in human output mode
- **THEN** paths SHALL be displayed relative to the spectr/ directory
- **AND** paths SHALL NOT include the project root or spectr/ prefix
- **AND** example: `changes/foo/specs/bar/spec.md` instead of `/home/user/project/spectr/changes/foo/specs/bar/spec.md`

#### Scenario: Grouping issues by file
- **WHEN** multiple issues exist in the same file in human output mode
- **THEN** the file path SHALL be displayed once as a header
- **AND** issues for that file SHALL be indented below the header
- **AND** this grouping SHALL reduce visual clutter from repeated paths

#### Scenario: Color-coded error levels
- **WHEN** bulk validation output is displayed in a TTY
- **THEN** [ERROR] labels SHALL be styled in red using lipgloss
- **AND** [WARNING] labels SHALL be styled in yellow using lipgloss
- **AND** [INFO] labels SHALL remain unstyled
- **AND** when not in a TTY, labels SHALL display without color codes

#### Scenario: Enhanced summary line
- **WHEN** bulk validation completes with failures
- **THEN** summary SHALL show "X passed, Y failed (E errors, W warnings), Z total"
- **AND** summary SHALL only show error/warning breakdown if failures exist
- **AND** example: "22 passed, 2 failed (5 errors, 1 warning), 24 total"

#### Scenario: Item type indicators
- **WHEN** bulk validation results are displayed in human output mode
- **THEN** each item SHALL show a type indicator alongside its name
- **AND** changes SHALL display "(change)" indicator
- **AND** specs SHALL display "(spec)" indicator

### Requirement: Markdown Parsing Package
The system SHALL provide a dedicated `internal/markdown/` package that provides AST-based markdown parsing with a token-based lexer and parser for Spectr-specific patterns.

#### Scenario: Two-phase parsing architecture
- **WHEN** the markdown package parses a document
- **THEN** it SHALL use a lexer to tokenize input into fine-grained tokens
- **AND** it SHALL use a parser to consume tokens and build an immutable AST
- **AND** the separation SHALL provide clear concerns and testable components

#### Scenario: Spectr-specific node types
- **WHEN** the markdown package parses Spectr documents
- **THEN** it SHALL recognize NodeRequirement nodes for `### Requirement: Name` headers
- **AND** it SHALL recognize NodeScenario nodes for `#### Scenario: Name` headers
- **AND** it SHALL recognize NodeSection nodes for standard markdown headers with Level and Title
- **AND** it SHALL recognize NodeListItem nodes with Checked state for task checkboxes

#### Scenario: Query functions for AST traversal
- **WHEN** code needs to find elements in the AST
- **THEN** it SHALL have access to `markdown.Find(root, predicate)` for finding all matching nodes
- **AND** it SHALL have access to `markdown.FindFirst(root, predicate)` for finding the first match
- **AND** it SHALL have access to predicate combinators like `IsType`, `HasName`, `And`, `Or`

#### Scenario: Visitor pattern for AST processing
- **WHEN** code needs to process AST nodes
- **THEN** it SHALL have access to `markdown.Walk(root, visitor)` for traversing the tree
- **AND** visitors SHALL implement type-specific methods like `VisitRequirement`, `VisitScenario`
- **AND** a BaseVisitor struct SHALL provide default implementations

#### Scenario: Position information
- **WHEN** nodes are created from parsing
- **THEN** they SHALL track byte offsets via `Span() (start, end int)`
- **AND** a LineIndex utility SHALL convert byte offsets to line/column numbers
- **AND** position information SHALL support error reporting with location context

#### Scenario: Immutable AST nodes
- **WHEN** the parser creates AST nodes
- **THEN** nodes SHALL be immutable after creation
- **AND** nodes SHALL provide getter methods for type-specific data
- **AND** transforms SHALL create new nodes rather than mutating existing ones
