# Cli Interface Specification

## Purpose

This specification defines interactive CLI features including navigable table interfaces for list and archive commands, cross-platform clipboard operations, initialization wizard flows, and visual styling for enhanced user experience.

## Requirements

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
- **WHEN** the user runs `spectr list --interactive`
- **THEN** the behavior is identical to before this change
- **AND** only changes are displayed
- **AND** columns show: ID, Title, Deltas, Tasks

#### Scenario: Unified mode opt-in
- **WHEN** the user explicitly uses `--all --interactive`
- **THEN** the new unified behavior is enabled
- **AND** users must opt-in to the new functionality
- **AND** columns show: Type, ID, Title, Details (context-aware)

#### Scenario: Unified mode displays both types
- **WHEN** unified mode is active
- **THEN** changes show Type="CHANGE" with delta and task counts
- **AND** specs show Type="SPEC" with requirement counts
- **AND** both types are navigable and selectable in the same table

#### Scenario: Type-specific actions in unified mode
- **WHEN** user presses 'e' on a change row in unified mode
- **THEN** the action is ignored (no edit for changes)
- **AND** help text does not show 'e' option
- **WHEN** user presses 'e' on a spec row in unified mode
- **THEN** the spec opens in the editor as usual

#### Scenario: Help text uses minimal footer by default
- **WHEN** interactive mode is displayed in any mode (changes, specs, or unified)
- **THEN** the footer shows: item count, project path, and `?: help`
- **AND** the full hotkey reference is hidden until `?` is pressed

#### Scenario: Help text format for changes mode
- **WHEN** user presses `?` in changes mode (`spectr list -I`)
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | q: quit`
- **AND** pressing `?` again or navigating hides the help

#### Scenario: Help text format for specs mode
- **WHEN** user presses `?` in specs mode (`spectr list --specs -I`)
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | q: quit`
- **AND** archive hotkey is NOT shown (specs cannot be archived)

### Requirement: Clipboard Copy on Selection
When a user presses Enter on a selected row in interactive mode, the item's ID SHALL be copied to the system clipboard.

#### Scenario: Copy change ID to clipboard
- **WHEN** user selects a change row and presses Enter
- **THEN** the change ID (kebab-case identifier) is copied to clipboard
- **AND** a success message is displayed (e.g., "Copied: add-archive-command")
- **AND** the interactive mode exits

#### Scenario: Copy spec ID to clipboard
- **WHEN** user selects a spec row and presses Enter
- **THEN** the spec ID is copied to clipboard
- **AND** a success message is displayed
- **AND** the interactive mode exits

#### Scenario: Clipboard failure handling
- **WHEN** clipboard operation fails
- **THEN** display error message to user
- **AND** do not exit interactive mode
- **AND** user can retry or quit manually

### Requirement: Interactive Mode Exit Controls
Users SHALL be able to exit interactive mode using standard quit commands.

#### Scenario: Quit with q key
- **WHEN** user presses 'q'
- **THEN** interactive mode exits
- **AND** no clipboard operation occurs
- **AND** command returns successfully

#### Scenario: Quit with Ctrl+C
- **WHEN** user presses Ctrl+C
- **THEN** interactive mode exits immediately
- **AND** no clipboard operation occurs
- **AND** command returns successfully

### Requirement: Table Visual Styling
The interactive table SHALL use clear visual styling to distinguish headers, selected rows, and borders, provided by the shared `internal/tui` package.

#### Scenario: Visual hierarchy in table
- **WHEN** interactive mode is displayed
- **THEN** column headers are visually distinct from data rows
- **AND** selected row has contrasting background/foreground colors
- **AND** table borders are visible and styled consistently
- **AND** table fits within terminal width gracefully
- **AND** styling SHALL be applied via `tui.ApplyTableStyles()`

#### Scenario: Consistent styling across commands
- **WHEN** user uses `spectr list -I`, `spectr archive`, or `spectr validate` interactive modes
- **THEN** all tables SHALL use identical styling
- **AND** colors, borders, and highlights SHALL match exactly
- **AND** the shared `tui.ApplyTableStyles()` function SHALL be the single source of truth

### Requirement: Cross-Platform Clipboard Support
Clipboard operations SHALL work across Linux, macOS, and Windows platforms.

#### Scenario: Clipboard on Linux
- **WHEN** running on Linux
- **THEN** clipboard operations use X11 or Wayland clipboard APIs as appropriate
- **AND** fallback to OSC 52 escape sequences if desktop clipboard unavailable

#### Scenario: Clipboard on macOS
- **WHEN** running on macOS
- **THEN** clipboard operations use pbcopy or native clipboard APIs

#### Scenario: Clipboard on Windows
- **WHEN** running on Windows
- **THEN** clipboard operations use Windows clipboard APIs

#### Scenario: Clipboard in SSH/remote session
- **WHEN** running over SSH without X11 forwarding
- **THEN** use OSC 52 escape sequences to copy to local clipboard
- **AND** document this behavior for users

### Requirement: Initialization Next Steps Message

The `spectr init` command SHALL display a formatted "Next steps" message after successful initialization that provides users with clear, actionable guidance for getting started with Spectr.

The message SHALL include:
1. Three progressive steps with copy-paste ready prompts for AI assistants
2. Visual separators to make the message stand out
3. References to key Spectr files and documentation
4. Placeholder text that users can customize (e.g., "[YOUR FEATURE HERE]")

The init command SHALL NOT automatically create project files outside the `spectr/` directory (such as README.md). Users maintain full control over their project's root-level documentation.

#### Scenario: Interactive mode initialization succeeds

- **WHEN** a user completes initialization via the interactive TUI wizard
- **THEN** the completion screen SHALL display the next steps message
- **AND** the message SHALL appear after the list of created/updated files
- **AND** the message SHALL be visually distinct with a separator line
- **AND** the message SHALL provide three numbered steps with specific prompts

#### Scenario: Non-interactive mode initialization succeeds

- **WHEN** a user runs `spectr init --non-interactive` and initialization succeeds
- **THEN** the command output SHALL display the next steps message
- **AND** the message SHALL appear after the list of created/updated files
- **AND** the message SHALL be formatted consistently with the interactive mode
- **AND** the message SHALL include the same three progressive steps

#### Scenario: Initialization fails with errors

- **WHEN** initialization fails with errors
- **THEN** the next steps message SHALL NOT be displayed
- **AND** only error messages SHALL be shown

#### Scenario: Next steps message content

- **WHEN** the next steps message is displayed
- **THEN** step 1 SHALL guide users to populate spectr/project.md
- **AND** step 2 SHALL guide users to create their first change proposal
- **AND** step 3 SHALL guide users to learn the Spectr workflow from spectr/AGENTS.md
- **AND** each step SHALL include a complete, copy-paste ready prompt in quotes
- **AND** the message SHALL include a visual separator using dashes or similar characters

#### Scenario: Init does not create README

- **WHEN** a user runs `spectr init` on a project without a README.md
- **THEN** the init command SHALL NOT create a README.md file
- **AND** only files within the `spectr/` directory SHALL be created
- **AND** tool-specific files (e.g., CLAUDE.md, .cursor/) SHALL be created as configured

### Requirement: Flat Tool List in Initialization Wizard

The initialization wizard SHALL present all AI tool options in a single unified flat list without visual grouping by tool type. Slash-only tool entries SHALL be removed from the registry as their functionality is now provided via automatic installation when the corresponding config-based tool is selected.

#### Scenario: Display only config-based tools in wizard

- **WHEN** user runs `spectr init` and reaches the tool selection screen
- **THEN** only config-based AI tools are displayed (e.g., `claude-code`, `cline`, `cursor`)
- **AND** slash-only tool entries (e.g., `claude`, `kilocode`) are not shown
- **AND** tools are sorted by priority
- **AND** no section headers (e.g., "Config-Based Tools", "Slash Command Tools") are shown
- **AND** each tool appears as a single checkbox item with its name

#### Scenario: Keyboard navigation across displayed tools

- **WHEN** user navigates with arrow keys (↑/↓)
- **THEN** the cursor moves through all displayed config-based tools sequentially
- **AND** navigation is continuous without group boundaries
- **AND** the first tool is selected by default on screen load

#### Scenario: Tool selection works uniformly

- **WHEN** user presses space to toggle any tool
- **THEN** the checkbox state changes (checked/unchecked)
- **AND** selection state is preserved when navigating
- **AND** both config file and slash commands will be installed when confirmed

#### Scenario: Bulk selection operations

- **WHEN** user presses 'a' to select all
- **THEN** all displayed config-based tools are checked
- **AND** WHEN user presses 'n' to select none
- **THEN** all tools are unchecked
- **AND** operations work across all displayed tools

#### Scenario: Help text clarity

- **WHEN** the tool selection screen is displayed
- **THEN** the help text shows keyboard controls (↑/↓, space, a, n, enter, q)
- **AND** the help text does NOT reference tool groupings or categories
- **AND** the screen title clearly indicates "Select AI Tools to Configure"

#### Scenario: Reduced tool count in wizard

- **WHEN** the wizard displays the tool list
- **THEN** fewer total tools are shown compared to the previous implementation
- **AND** the count reflects only config-based tools (not slash-only duplicates)
- **AND** navigation and selection work correctly with the reduced count

### Requirement: Interactive Archive Mode
The archive command SHALL provide an interactive table interface when no change ID argument is provided or when the `-I` or `--interactive` flag is used, displaying available changes in a navigable table format identical to the list command's interactive mode with project path information.

#### Scenario: User runs archive with no arguments
- **WHEN** user runs `spectr archive` with no change ID argument
- **THEN** an interactive table is displayed with columns: ID, Title, Deltas, Tasks
- **AND** the table supports arrow key navigation (↑/↓, j/k)
- **AND** the first row is selected by default
- **AND** the table uses the same visual styling as list -I
- **AND** the project path is displayed in the interface

#### Scenario: User runs archive with -I flag
- **WHEN** user runs `spectr archive -I`
- **THEN** an interactive table is displayed even if other flags are present
- **AND** the behavior is identical to running archive with no arguments
- **AND** the project path is displayed in the interface

#### Scenario: User selects change for archiving
- **WHEN** user presses Enter on a selected row in archive interactive mode
- **THEN** the change ID is captured (not copied to clipboard)
- **AND** the interactive mode exits
- **AND** the archive workflow proceeds with the selected change ID
- **AND** validation, task checking, and spec updates proceed as normal

#### Scenario: User cancels archive selection
- **WHEN** user presses 'q' or Ctrl+C in archive interactive mode
- **THEN** interactive mode exits
- **AND** archive command returns successfully without archiving anything
- **AND** a "Cancelled" message is displayed

#### Scenario: No changes available for archiving
- **WHEN** user runs `spectr archive` and no changes exist in changes/ directory
- **THEN** display "No changes available to archive" message
- **AND** exit cleanly without entering interactive mode
- **AND** command returns successfully

#### Scenario: Archive with explicit change ID bypasses interactive mode
- **WHEN** user runs `spectr archive <change-id>`
- **THEN** interactive mode is NOT triggered
- **AND** archive proceeds directly with the specified change ID
- **AND** behavior is unchanged from current implementation

### Requirement: Archive Interactive Table Display
The archive command's interactive table SHALL display the same information columns as the list command to help users make informed archiving decisions.

#### Scenario: Table columns match list command
- **WHEN** archive interactive mode is displayed
- **THEN** columns are: ID (30 chars), Title (40 chars), Deltas (10 chars), Tasks (15 chars)
- **AND** column widths match the list -I command exactly
- **AND** title text is truncated with ellipsis if longer than 38 characters
- **AND** task status shows format "completed/total" (e.g., "5/10")

#### Scenario: Visual styling consistency
- **WHEN** archive interactive table is displayed
- **THEN** the table uses identical styling to list -I
- **AND** column headers are visually distinct from data rows
- **AND** selected row has contrasting background/foreground colors
- **AND** table borders are visible and styled consistently
- **AND** help text shows navigation controls (↑/↓, j/k, enter, q)

### Requirement: Archive Selection Without Clipboard
The archive command's interactive mode SHALL NOT copy the selected change ID to the clipboard, unlike the list command, since the ID is immediately consumed by the archive workflow.

#### Scenario: Enter key captures selection
- **WHEN** user presses Enter on a selected change
- **THEN** the change ID is captured internally
- **AND** NO clipboard operation occurs
- **AND** NO "Copied: <id>" message is displayed
- **AND** the archive workflow proceeds immediately with the selected ID

#### Scenario: Workflow continuation
- **WHEN** a change is selected in interactive mode
- **THEN** the Archiver.Archive() method receives the selected change ID
- **AND** validation, task checking, and spec updates proceed as if the ID was provided as an argument
- **AND** all confirmation prompts and flags (--yes, --skip-specs) work normally

### Requirement: Validation Output Format
The validate command SHALL display validation issues in a consistent, detailed format for both single-item and bulk validation modes.

#### Scenario: Single item validation with issues
- **WHEN** user runs `spectr validate <item>` and validation finds issues
- **THEN** output SHALL display "✗ <item> has N issue(s):"
- **AND** each issue SHALL be displayed on a separate line with format "  [LEVEL] PATH: MESSAGE"
- **AND** the command SHALL exit with code 1

#### Scenario: Bulk validation with issues
- **WHEN** user runs `spectr validate --all` and validation finds issues in multiple items
- **THEN** output SHALL display "✗ <item> (<type>): N issue(s)" for each failed item
- **AND** immediately following each failed item, all issue details SHALL be displayed
- **AND** each issue SHALL use the format "  [LEVEL] PATH: MESSAGE"
- **AND** a summary line SHALL display "N passed, M failed, T total"
- **AND** the command SHALL exit with code 1

#### Scenario: Bulk validation all passing
- **WHEN** user runs `spectr validate --all` and all items are valid
- **THEN** output SHALL display "✓ <item> (<type>)" for each item
- **AND** a summary line SHALL display "N passed, 0 failed, N total"
- **AND** the command SHALL exit with code 0

#### Scenario: JSON output format
- **WHEN** user provides `--json` flag with any validation command
- **THEN** output SHALL be valid JSON
- **AND** SHALL include full issue details with level, path, and message fields
- **AND** SHALL include per-item results and summary statistics

### Requirement: Editor Hotkey in Interactive Specs List
The interactive specs list mode SHALL provide an 'e' hotkey that opens the selected spec file in the user's configured editor.

#### Scenario: User presses 'e' to edit a spec
- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses the 'e' key on a selected spec
- **THEN** the file `spectr/specs/<spec-id>/spec.md` is opened in the editor specified by $EDITOR environment variable
- **AND** the TUI waits for the editor to close
- **AND** the TUI remains active after the editor closes
- **AND** the same row remains selected

#### Scenario: User edits spec and returns to TUI
- **WHEN** user presses 'e' to open a spec
- **AND** makes changes in the editor and saves
- **AND** closes the editor
- **THEN** the TUI returns to the interactive list view
- **AND** the user can continue navigating or edit another spec
- **AND** the user can quit with 'q' or Ctrl+C as normal

#### Scenario: EDITOR environment variable not set
- **WHEN** user presses 'e' to edit a spec
- **AND** $EDITOR environment variable is not set
- **THEN** display an error message "EDITOR environment variable not set"
- **AND** the TUI remains in interactive mode
- **AND** the user can continue navigating or quit

#### Scenario: Spec file does not exist
- **WHEN** user presses 'e' to edit a spec
- **AND** the spec file at `spectr/specs/<spec-id>/spec.md` does not exist
- **THEN** display an error message "Spec file not found: <path>"
- **AND** the TUI remains in interactive mode
- **AND** the user can continue navigating or quit

#### Scenario: Editor launch fails
- **WHEN** user presses 'e' to edit a spec
- **AND** the editor process fails to launch (e.g., editor binary not found, permission error)
- **THEN** display an error message with the underlying error details
- **AND** the TUI remains in interactive mode
- **AND** the user can retry or quit

#### Scenario: Help text shows editor hotkey
- **WHEN** interactive specs mode is displayed
- **THEN** the help text includes "e: edit spec" or similar guidance
- **AND** the help text shows all available keys including navigation, enter, e, and quit keys

### Requirement: Editor Hotkey Scope
The 'e' hotkey for opening files in $EDITOR SHALL only be available in specs list mode, not in changes list mode.

#### Scenario: Editor hotkey not available for changes
- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses 'e' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'e: edit' option
- **AND** only standard navigation and clipboard actions are available

#### Scenario: Rationale for specs-only scope
- **WHEN** user reviews this specification
- **THEN** they understand that changes have multiple files (proposal.md, tasks.md, design.md, delta specs)
- **AND** pressing 'e' on a change would be ambiguous (which file to open?)
- **AND** specs have a single canonical file (spec.md) making 'e' unambiguous
- **AND** this design decision can be revisited in a future change if multi-file editing is needed

### Requirement: Project Path Display in Interactive Mode
The interactive table interfaces SHALL display the project root path to provide users with context about which project they are working with.

#### Scenario: Project path shown in changes interactive mode
- **WHEN** user runs `spectr list -I` for changes
- **THEN** the project root path is displayed in the help text or table header
- **AND** the path is the absolute path to the project directory

#### Scenario: Project path shown in specs interactive mode
- **WHEN** user runs `spectr list --specs -I`
- **THEN** the project root path is displayed in the help text or table header
- **AND** the path is the absolute path to the project directory

#### Scenario: Project path shown in archive interactive mode
- **WHEN** user runs `spectr archive` without arguments
- **THEN** the project root path is displayed in the help text or table header
- **AND** the path is the absolute path to the project directory

#### Scenario: Project path properly initialized for changes
- **WHEN** `RunInteractiveChanges()` is invoked
- **THEN** the `projectPath` parameter is passed from the calling command
- **AND** the `projectPath` field on `interactiveModel` is set during initialization

#### Scenario: Project path properly initialized for archive
- **WHEN** `RunInteractiveArchive()` is invoked
- **THEN** the `projectPath` parameter is passed from the calling command
- **AND** the `projectPath` field on `interactiveModel` is set during initialization

### Requirement: Unified Item List Display
The system SHALL display changes and specifications together in a single interactive table when invoked with appropriate flags, allowing users to browse both item types simultaneously with clear visual differentiation.

#### Scenario: User opens unified interactive list
- **WHEN** the user runs `spectr list --interactive --all` from a directory with both changes and specs
- **THEN** a table appears showing both changes and specs rows
- **AND** each row indicates its type (change or spec)
- **AND** the table maintains correct ordering and alignment

#### Scenario: Unified list shows correct columns
- **WHEN** the unified interactive mode is active
- **THEN** the table displays: Type, ID, Title, and Type-Specific Details columns
- **AND** "Type-Specific Details" shows "Deltas/Tasks" for changes
- **AND** "Type-Specific Details" shows "Requirements" for specs

#### Scenario: User navigates mixed items
- **WHEN** the user navigates with arrow keys through a mixed list
- **THEN** the cursor moves smoothly between change and spec rows
- **AND** help text remains accurate and updated
- **AND** the selected row is clearly highlighted

### Requirement: Type-Aware Item Selection
The system SHALL track whether a selected item is a change or spec and provide type-appropriate actions (e.g., edit only works for specs).

#### Scenario: Selecting a spec in unified mode
- **WHEN** the user presses Enter on a spec row
- **THEN** the spec ID is copied to clipboard
- **AND** a success message displays the ID and type indicator
- **AND** the user is returned to the interactive session or exited cleanly

#### Scenario: Selecting a change in unified mode
- **WHEN** the user presses Enter on a change row
- **THEN** the change ID is copied to clipboard
- **AND** a success message displays the ID and type indicator
- **AND** no edit action is attempted

#### Scenario: Edit action restricted to specs
- **WHEN** the user presses 'e' on a change row in unified mode
- **THEN** the action is ignored or a helpful message appears
- **AND** the interactive session continues

### Requirement: Backward-Compatible Separate Modes
The system SHALL maintain existing interactive modes for changes-only and specs-only when `--all` flag is not provided.

#### Scenario: Changes-only mode still works
- **WHEN** the user runs `spectr list --interactive` without `--all`
- **THEN** only changes are displayed
- **AND** behavior is identical to the previous implementation
- **AND** edit functionality works as before for changes

#### Scenario: Specs-only mode still works
- **WHEN** the user runs `spectr list --specs --interactive` without `--all`
- **THEN** only specs are displayed
- **AND** behavior is identical to the previous implementation
- **AND** edit functionality works as before for specs

### Requirement: Enhanced List Command Flags
The system SHALL support new flag combinations to control listing behavior while maintaining validation for mutually exclusive options.

#### Scenario: Flag validation for unified mode
- **WHEN** the user attempts `spectr list --interactive --all --json`
- **THEN** an error message is returned: "cannot use --interactive with --json"
- **AND** the command exits without running

#### Scenario: All flag with separate type flags
- **WHEN** the user provides `--all` with `--specs`
- **THEN** `--all` takes precedence and unified mode is used
- **AND** a warning may be shown (optional) about the redundant flag

#### Scenario: All flag in non-interactive mode
- **WHEN** the user runs `spectr list --all` without `--interactive`
- **THEN** both changes and specs are listed in text format
- **AND** each item shows its type in the output

### Requirement: Automatic Slash Command Installation

When a config-based AI tool is selected during initialization, the system SHALL automatically install the corresponding slash command files for that tool without requiring separate user selection.

Config-based tools include those that create instruction files (e.g., `claude-code` creates `CLAUDE.md`). Slash command files are the workflow command files (e.g., `.claude/commands/spectr/proposal.md`).

The `ToolDefinition` model SHALL NOT include a `ConfigPath` field, as actual file paths are determined by individual configurators. The registry maintains tool metadata (ID, Name, Type, Priority) but delegates file path resolution to configurator implementations. Tool IDs SHALL use a type-safe constant approach to prevent typos and enable compile-time validation.

This automatic installation provides users with complete Spectr integration in a single selection, eliminating the need for redundant tool entries in the wizard.

#### Scenario: Claude Code auto-installs slash commands

- **WHEN** user selects `claude-code` in the init wizard
- **THEN** the system creates `CLAUDE.md` in the project root
- **AND** the system creates `.claude/commands/spectr/proposal.md`
- **AND** the system creates `.claude/commands/spectr/apply.md`
- **AND** all files are tracked in the execution result
- **AND** the completion screen shows all 3 files created

#### Scenario: Multiple tools with slash commands selected

- **WHEN** user selects both `claude-code` and `cursor` in the init wizard
- **THEN** the system creates `CLAUDE.md` and both config + slash commands for Claude
- **AND** the system creates `.cursor/commands/spectr/proposal.md` and slash commands for Cursor
- **AND** all files from both tools are created and tracked separately
- **AND** the completion screen lists all created files grouped by tool

#### Scenario: Slash command files already exist

- **WHEN** user run init and selects `claude-code`
- **AND** `.claude/commands/spectr/proposal.md` already exists
- **THEN** the existing file's content between `<!-- spectr:START -->` and `<!-- spectr:END -->` is updated
- **AND** the file's YAML frontmatter is preserved
- **AND** no error occurs
- **AND** the file is marked as "updated" rather than "created" in execution result

### Requirement: Archive Hotkey in Interactive Changes Mode
The interactive changes list mode SHALL provide an 'a' hotkey that archives the currently selected change, invoking the same workflow as `spectr archive <change-id>`.

#### Scenario: User presses 'a' to archive a change
- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses the 'a' key on a selected change
- **THEN** the interactive mode exits
- **AND** the archive workflow begins for the selected change ID
- **AND** validation, task checking, and spec updates proceed as if the ID was provided as an argument
- **AND** all confirmation prompts and flags work normally

#### Scenario: Archive hotkey not available in specs mode
- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses 'a' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'a: archive' option

#### Scenario: Archive hotkey not available in unified mode
- **WHEN** user is in unified interactive mode (`spectr list --all -I`)
- **AND** user presses 'a' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'a: archive' option
- **AND** this avoids confusion when a spec row is selected

#### Scenario: Archive workflow integration
- **WHEN** the archive hotkey triggers the archive workflow
- **THEN** the workflow uses the same code path as `spectr archive <id>`
- **AND** the selected change ID is passed to the archive workflow
- **AND** success or failure is reported after the workflow completes

#### Scenario: Help text shows archive hotkey in changes mode
- **WHEN** interactive changes mode is displayed
- **THEN** the help text includes `a: archive` in the controls line
- **AND** the hotkey appears after `e: edit` and before `q: quit`

### Requirement: Shared TUI Component Library

The CLI SHALL use a shared `internal/tui` package for interactive TUI components, providing consistent styling, behavior, and composable building blocks across all interactive modes.

#### Scenario: TablePicker used for item selection
- **WHEN** any command needs an interactive table-based selection (list, archive, validation item picker)
- **THEN** the command SHALL use the `TablePicker` component from `internal/tui`
- **AND** the table SHALL use consistent styling from `tui.ApplyTableStyles()`
- **AND** navigation keys (↑/↓, j/k) SHALL work identically across all usages
- **AND** quit keys (q, Ctrl+C) SHALL work identically across all usages

#### Scenario: MenuPicker used for option selection
- **WHEN** any command needs an interactive menu selection (validation mode menu)
- **THEN** the command SHALL use the `MenuPicker` component from `internal/tui`
- **AND** the menu SHALL use consistent styling
- **AND** navigation and selection behavior SHALL match the TablePicker patterns

#### Scenario: Consistent string truncation
- **WHEN** any TUI component needs to truncate text for display
- **THEN** it SHALL use `tui.TruncateString()` with consistent ellipsis handling
- **AND** truncation SHALL add "..." suffix when text exceeds max length
- **AND** very short max lengths (≤3) SHALL truncate without ellipsis

#### Scenario: Consistent clipboard operations
- **WHEN** any TUI component needs to copy text to clipboard
- **THEN** it SHALL use `tui.CopyToClipboard()` from the shared package
- **AND** the function SHALL try native clipboard first
- **AND** the function SHALL fall back to OSC 52 for remote sessions

#### Scenario: Action registration pattern
- **WHEN** a command configures a TablePicker with custom actions
- **THEN** actions SHALL be registered via `WithAction(key, label, handler)`
- **AND** the help text SHALL automatically include all registered actions
- **AND** unregistered keys SHALL be ignored (no error)

#### Scenario: Domain logic remains in consuming packages
- **WHEN** the tui package is used by list or validation
- **THEN** domain-specific logic (archive workflow, validation execution) SHALL remain in consuming packages
- **AND** the tui package SHALL only provide UI primitives
- **AND** business logic SHALL not be coupled to the tui package

### Requirement: Search Hotkey in Interactive Lists
The interactive list modes SHALL provide a '/' hotkey that activates a text search mode, allowing users to filter the displayed list by typing a search query that matches against item IDs and titles.

#### Scenario: User presses '/' to enter search mode
- **WHEN** user is in any interactive list mode (changes, specs, or unified)
- **AND** user presses the '/' key
- **THEN** search mode is activated
- **AND** a text input field is displayed below or above the table
- **AND** the cursor is placed in the text input field
- **AND** the user can type a search query

#### Scenario: Search filters rows in real-time
- **WHEN** search mode is active
- **AND** user types characters into the search input
- **THEN** the table rows are filtered in real-time
- **AND** only rows where ID or title contains the search query (case-insensitive) are displayed
- **AND** the first matching row is automatically selected

#### Scenario: Search with no matches shows empty table
- **WHEN** search mode is active
- **AND** user types a query that matches no items
- **THEN** the table displays no rows
- **AND** a message indicates no matches found

#### Scenario: User presses Escape to exit search mode
- **WHEN** search mode is active
- **AND** user presses the Escape key
- **THEN** search mode is deactivated
- **AND** the search query is cleared
- **AND** all items are displayed again in the table
- **AND** the text input field is hidden

#### Scenario: User presses '/' again to clear search
- **WHEN** search mode is active
- **AND** the search query is not empty
- **AND** user presses '/' key
- **THEN** the search input gains focus (normal text input behavior)

- **WHEN** search mode is active
- **AND** the search query is empty
- **AND** user presses '/' key
- **THEN** search mode is deactivated
- **AND** all items are displayed again

#### Scenario: Navigation works while searching
- **WHEN** search mode is active
- **AND** filtered results are displayed
- **THEN** arrow key navigation (up/down, j/k) moves through filtered rows
- **AND** Enter key copies the selected filtered item's ID
- **AND** other hotkeys (e, a, t) work on the selected filtered item

#### Scenario: Help text shows search hotkey
- **WHEN** interactive mode is displayed in any mode
- **THEN** the help text includes '/: search' in the controls line
- **AND** the search hotkey is shown for all modes (changes, specs, unified)

#### Scenario: Search mode visual indicator
- **WHEN** search mode is active
- **THEN** the search input field is visually distinct
- **AND** the current search query is visible
- **AND** the help text updates to show 'Esc: exit search'

### Requirement: Help Toggle Hotkey
The interactive TUI modes SHALL hide hotkey hints by default and reveal them only when the user presses `?`, reducing visual clutter while maintaining discoverability.

#### Scenario: Default view shows minimal footer
- **WHEN** user enters any interactive TUI mode (list, archive, validate)
- **THEN** the footer displays only: item count, project path, and `?: help`
- **AND** the full hotkey reference is NOT shown
- **AND** navigation and all other hotkeys remain functional

#### Scenario: User presses '?' to reveal help
- **WHEN** user presses `?` while in interactive mode
- **THEN** the full hotkey reference is displayed in the footer area
- **AND** the reference includes all available hotkeys for the current mode
- **AND** the view updates immediately

#### Scenario: User dismisses help by pressing '?' again
- **WHEN** user presses `?` while help is visible
- **THEN** the help is hidden
- **AND** the minimal footer is restored

#### Scenario: Help auto-hides on navigation
- **WHEN** user presses a navigation key (↑/↓/j/k) while help is visible
- **THEN** the help is automatically hidden
- **AND** the navigation action is performed
- **AND** the minimal footer is restored

#### Scenario: Help content matches mode
- **WHEN** help is displayed in changes mode
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | q: quit`
- **WHEN** help is displayed in specs mode
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | q: quit`
- **WHEN** help is displayed in unified mode
- **THEN** the help shows: `↑/↓/j/k: navigate | Enter: copy ID | e: edit | a: archive | t: filter | q: quit`

### Requirement: Partial Change ID Resolution for Archive Command

The `spectr archive` command SHALL support intelligent partial ID matching when a non-exact change ID is provided as an argument. The resolution algorithm SHALL prioritize prefix matches over substring matches and require a unique match to proceed.

#### Scenario: Exact ID match takes precedence

- **WHEN** user runs `spectr archive add-feature`
- **AND** a change with ID `add-feature` exists
- **THEN** the archive proceeds with `add-feature`
- **AND** no resolution message is displayed

#### Scenario: Unique prefix match resolves successfully

- **WHEN** user runs `spectr archive refactor`
- **AND** only one change ID starts with `refactor` (e.g., `refactor-unified-interactive-tui`)
- **THEN** a message is displayed: "Resolved 'refactor' -> 'refactor-unified-interactive-tui'"
- **AND** the archive proceeds with the resolved ID

#### Scenario: Unique substring match resolves successfully

- **WHEN** user runs `spectr archive unified`
- **AND** no change ID starts with `unified`
- **AND** only one change ID contains `unified` (e.g., `refactor-unified-interactive-tui`)
- **THEN** a message is displayed: "Resolved 'unified' -> 'refactor-unified-interactive-tui'"
- **AND** the archive proceeds with the resolved ID

#### Scenario: Multiple prefix matches cause error

- **WHEN** user runs `spectr archive add`
- **AND** multiple change IDs start with `add` (e.g., `add-feature`, `add-hotkey`)
- **THEN** an error is displayed: "Ambiguous ID 'add' matches multiple changes: add-feature, add-hotkey"
- **AND** the command exits with error code 1
- **AND** no archive operation is performed

#### Scenario: Multiple substring matches cause error

- **WHEN** user runs `spectr archive search`
- **AND** no change ID starts with `search`
- **AND** multiple change IDs contain `search` (e.g., `add-search-hotkey`, `update-search-ui`)
- **THEN** an error is displayed: "Ambiguous ID 'search' matches multiple changes: add-search-hotkey, update-search-ui"
- **AND** the command exits with error code 1
- **AND** no archive operation is performed

#### Scenario: No match found

- **WHEN** user runs `spectr archive nonexistent`
- **AND** no change ID matches `nonexistent` (neither prefix nor substring)
- **THEN** an error is displayed: "No change found matching 'nonexistent'"
- **AND** the command exits with error code 1
- **AND** no archive operation is performed

#### Scenario: Case-insensitive matching

- **WHEN** user runs `spectr archive REFACTOR`
- **AND** a change ID `refactor-unified-interactive-tui` exists
- **THEN** the partial match succeeds (case-insensitive)
- **AND** the archive proceeds with the resolved ID

#### Scenario: Prefix match preferred over substring match

- **WHEN** user runs `spectr archive add`
- **AND** change ID `add-feature` exists (prefix match)
- **AND** change ID `update-add-button` exists (substring match only)
- **THEN** the prefix match `add-feature` is selected
- **AND** the substring-only match is ignored in preference calculation

### Requirement: Configured Provider Detection in Init Wizard

The initialization wizard SHALL detect which AI tool providers are already configured for the project and display this status in the tool selection screen. Already-configured providers SHALL be visually distinguished and pre-selected by default.

#### Scenario: Display configured indicator for already-configured providers

- **WHEN** user runs `spectr init` on a project with `CLAUDE.md` already present
- **AND** user reaches the tool selection screen
- **THEN** the Claude Code entry displays a "configured" indicator (e.g., dimmed text or badge)
- **AND** the indicator is visually distinct from the selection checkbox
- **AND** other unconfigured providers do NOT show the configured indicator

#### Scenario: Pre-select already-configured providers

- **WHEN** user runs `spectr init` on a project with some providers already configured
- **AND** user reaches the tool selection screen
- **THEN** already-configured providers have their checkboxes pre-selected
- **AND** users can deselect them if they don't want to update the configuration
- **AND** unconfigured providers remain unselected by default

#### Scenario: Help text explains configured indicator

- **WHEN** user is on the tool selection screen
- **THEN** the help text or screen description explains what the "configured" indicator means
- **AND** the explanation clarifies that selecting a configured provider will update its files

#### Scenario: No configured providers

- **WHEN** user runs `spectr init` on a fresh project with no providers configured
- **AND** user reaches the tool selection screen
- **THEN** no providers show the configured indicator
- **AND** no providers are pre-selected
- **AND** the screen functions as before this change

#### Scenario: All providers configured

- **WHEN** user runs `spectr init` on a project with all available providers configured
- **AND** user reaches the tool selection screen
- **THEN** all providers show the configured indicator
- **AND** all providers are pre-selected
- **AND** user can deselect providers they don't want to update

#### Scenario: Configured detection uses provider's IsConfigured method

- **WHEN** the wizard initializes
- **THEN** it calls `IsConfigured(projectPath)` on each provider
- **AND** the result is cached for the wizard session (not re-checked on each render)
- **AND** providers with global paths (like Codex) are correctly detected

### Requirement: Instruction File Pointer Template

The system SHALL use a short pointer template when injecting Spectr instructions into root-level instruction files (e.g., `CLAUDE.md`, `AGENTS.md` at project root), directing AI assistants to read `spectr/AGENTS.md` for full instructions rather than duplicating the entire content.

#### Scenario: Init creates instruction file with pointer

- **WHEN** user runs `spectr init` and selects an AI tool (e.g., Claude Code)
- **THEN** the root-level instruction file (e.g., `CLAUDE.md`) contains a short pointer between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers
- **AND** the pointer directs AI assistants to read `spectr/AGENTS.md` when handling proposals, specs, or changes
- **AND** the full instructions remain only in `spectr/AGENTS.md`

#### Scenario: Update refreshes instruction file with pointer

- **WHEN** user runs `spectr init` on an already-initialized project
- **THEN** the root-level instruction files are updated with the short pointer content
- **AND** the `spectr/AGENTS.md` file retains the full instructions

#### Scenario: Pointer content is concise

- **WHEN** the instruction pointer template is rendered
- **THEN** the output is less than 20 lines
- **AND** the output explains when to read `spectr/AGENTS.md` (proposals, specs, changes, planning)
- **AND** the output does NOT duplicate the full workflow instructions

### Requirement: PR Archive Subcommand Alias
The `spectr pr archive` subcommand SHALL support `a` as a shorthand alias, allowing users to invoke `spectr pr a <id>` as equivalent to `spectr pr archive <id>`.

#### Scenario: User runs spectr pr a shorthand
- **WHEN** user runs `spectr pr a <change-id>`
- **THEN** the system executes the archive PR workflow identically to `spectr pr archive`
- **AND** all flags (`--base`, `--draft`, `--force`, `--dry-run`, `--skip-specs`) work with the alias

#### Scenario: User runs spectr pr a with flags
- **WHEN** user runs `spectr pr a my-change --draft --force`
- **THEN** the command behaves identically to `spectr pr archive my-change --draft --force`
- **AND** a draft PR is created after deleting any existing branch

#### Scenario: Help text shows archive alias
- **WHEN** user runs `spectr pr --help`
- **THEN** the help text displays `archive` with its `a` alias
- **AND** the alias is shown in parentheses or as comma-separated alternatives

### Requirement: PR Branch Naming Convention
The system SHALL use a mode-specific branch naming convention for PR branches that distinguishes between archive and proposal branches based on the subcommand used.

#### Scenario: Archive branch name format
- **WHEN** user runs `spectr pr archive <change-id>`
- **THEN** the branch is named `spectr/archive/<change-id>`

#### Scenario: Proposal branch name format
- **WHEN** user runs `spectr pr new <change-id>`
- **THEN** the branch is named `spectr/proposal/<change-id>`

#### Scenario: Branch name with special characters
- **WHEN** change ID contains only valid kebab-case characters
- **THEN** the branch name is valid for git

#### Scenario: Branch names clearly indicate PR purpose
- **WHEN** a developer views the branch list
- **THEN** they can distinguish archive PRs from proposal PRs by the branch prefix
- **AND** `spectr/archive/*` indicates a completed change being archived
- **AND** `spectr/proposal/*` indicates a change proposal for review

#### Scenario: Force flag for existing archive branch
- **WHEN** user runs `spectr pr archive <change-id> --force`
- **AND** branch `spectr/archive/<change-id>` already exists on remote
- **THEN** the existing branch is deleted and recreated
- **AND** the PR workflow proceeds normally

#### Scenario: Force flag for existing proposal branch
- **WHEN** user runs `spectr pr new <change-id> --force`
- **AND** branch `spectr/proposal/<change-id>` already exists on remote
- **THEN** the existing branch is deleted and recreated
- **AND** the PR workflow proceeds normally

#### Scenario: Archive branch conflict without force
- **WHEN** user runs `spectr pr archive <change-id>`
- **AND** branch `spectr/archive/<change-id>` already exists on remote
- **AND** `--force` flag is NOT provided
- **THEN** an error is displayed: "branch 'spectr/archive/<change-id>' already exists on remote; use --force to delete"
- **AND** the command exits with code 1

#### Scenario: Proposal branch conflict without force
- **WHEN** user runs `spectr pr new <change-id>`
- **AND** branch `spectr/proposal/<change-id>` already exists on remote
- **AND** `--force` flag is NOT provided
- **THEN** an error is displayed: "branch 'spectr/proposal/<change-id>' already exists on remote; use --force to delete"
- **AND** the command exits with code 1

### Requirement: PR Command Structure
The system SHALL provide a `spectr pr` command with `archive` and `proposal` subcommands for creating pull requests from Spectr changes using git worktree isolation.

#### Scenario: User runs spectr pr without subcommand
- **WHEN** user runs `spectr pr` without a subcommand
- **THEN** help text is displayed showing available subcommands (archive, proposal)
- **AND** the command exits with code 0

#### Scenario: User runs spectr pr archive
- **WHEN** user runs `spectr pr archive <change-id>`
- **THEN** the system executes the archive PR workflow
- **AND** a PR is created with the archived change

#### Scenario: User runs spectr pr proposal
- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the system executes the proposal PR workflow
- **AND** a PR is created with the change proposal copied (not archived)

### Requirement: PR Archive Subcommand
The `spectr pr archive` subcommand SHALL create a pull request containing an archived Spectr change, executing the archive workflow in an isolated git worktree.

#### Scenario: Archive PR workflow execution
- **WHEN** user runs `spectr pr archive <change-id>`
- **THEN** the system creates a temporary git worktree on branch `spectr/<change-id>`
- **AND** executes `spectr archive <change-id> --yes` within the worktree
- **AND** stages all changes in `spectr/` directory
- **AND** commits with structured message including archive metadata
- **AND** pushes the branch to origin
- **AND** creates a PR using the detected platform's CLI
- **AND** cleans up the temporary worktree
- **AND** displays the PR URL on success

#### Scenario: Archive PR with skip-specs flag
- **WHEN** user runs `spectr pr archive <change-id> --skip-specs`
- **THEN** the `--skip-specs` flag is passed to the underlying archive command
- **AND** spec merging is skipped during the archive operation

#### Scenario: Archive PR preserves user working directory
- **WHEN** user runs `spectr pr archive <change-id>`
- **AND** user has uncommitted changes in their working directory
- **THEN** the user's working directory is NOT modified
- **AND** the archive operation executes only within the isolated worktree
- **AND** uncommitted changes are NOT included in the PR

### Requirement: PR Proposal Subcommand
The `spectr pr proposal` subcommand SHALL create a pull request containing a Spectr change proposal for review, copying the change to an isolated git worktree without archiving. This command replaces the deprecated `spectr pr new` command.

The renaming from `new` to `proposal` aligns CLI terminology with the `/spectr:proposal` slash command naming convention, creating consistent vocabulary across CLI and IDE integrations.

#### Scenario: Proposal PR workflow execution
- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the system creates a temporary git worktree on branch `spectr/<change-id>`
- **AND** copies the change directory from source to worktree
- **AND** stages all changes in `spectr/` directory
- **AND** commits with structured message for proposal review
- **AND** pushes the branch to origin
- **AND** creates a PR using the detected platform's CLI
- **AND** cleans up the temporary worktree
- **AND** displays the PR URL on success

#### Scenario: Proposal PR does not archive
- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the original change remains in `spectr/changes/<change-id>/`
- **AND** the change is NOT moved to archive
- **AND** spec merging does NOT occur

#### Scenario: Proposal PR validates change first
- **WHEN** user runs `spectr pr proposal <change-id>`
- **THEN** the system runs validation on the change
- **AND** warnings are displayed if validation issues exist
- **AND** the PR workflow continues (validation does not block)

#### Scenario: User runs spectr pr without subcommand
- **WHEN** user runs `spectr pr` without a subcommand
- **THEN** help text is displayed showing available subcommands (archive, proposal)
- **AND** the command exits with code 0

#### Scenario: Unique prefix match for PR proposal command
- **WHEN** user runs `spectr pr proposal refactor`
- **AND** only one change ID starts with `refactor`
- **THEN** a resolution message is displayed
- **AND** the PR workflow proceeds with the resolved ID

### Requirement: PR Common Flags
Both `spectr pr archive` and `spectr pr proposal` subcommands SHALL support common flags for controlling PR creation behavior.

#### Scenario: Base branch flag
- **WHEN** user provides `--base <branch>` flag
- **THEN** the PR targets the specified branch instead of auto-detected main/master

#### Scenario: Auto-detect base branch
- **WHEN** user does NOT provide `--base` flag
- **AND** `origin/main` exists
- **THEN** the PR targets `main`

#### Scenario: Fallback to master
- **WHEN** user does NOT provide `--base` flag
- **AND** `origin/main` does NOT exist
- **AND** `origin/master` exists
- **THEN** the PR targets `master`

#### Scenario: Draft PR flag
- **WHEN** user provides `--draft` flag
- **THEN** the PR is created as a draft PR on platforms that support it

#### Scenario: Force flag for existing branch
- **WHEN** user provides `--force` flag
- **AND** branch `spectr/<change-id>` already exists on remote
- **THEN** the existing branch is deleted and recreated
- **AND** the PR workflow proceeds normally

#### Scenario: Branch conflict without force
- **WHEN** branch `spectr/<change-id>` already exists on remote
- **AND** `--force` flag is NOT provided
- **THEN** an error is displayed: "Branch 'spectr/<change-id>' already exists. Use --force to overwrite."
- **AND** the command exits with code 1

#### Scenario: Dry run flag
- **WHEN** user provides `--dry-run` flag
- **THEN** the system displays what would be done without executing
- **AND** no git operations are performed
- **AND** no PR is created

### Requirement: Git Platform Detection
The system SHALL automatically detect the git hosting platform from the origin remote URL and use the appropriate CLI tool for PR creation.

#### Scenario: Detect GitHub platform
- **WHEN** origin remote URL contains `github.com`
- **THEN** platform is detected as GitHub
- **AND** `gh` CLI is used for PR creation

#### Scenario: Detect GitLab platform
- **WHEN** origin remote URL contains `gitlab.com` or `gitlab`
- **THEN** platform is detected as GitLab
- **AND** `glab` CLI is used for MR creation

#### Scenario: Detect Gitea platform
- **WHEN** origin remote URL contains `gitea` or `forgejo`
- **THEN** platform is detected as Gitea
- **AND** `tea` CLI is used for PR creation

#### Scenario: Detect Bitbucket platform
- **WHEN** origin remote URL contains `bitbucket.org` or `bitbucket`
- **THEN** platform is detected as Bitbucket
- **AND** manual PR URL is provided (no CLI automation)

#### Scenario: Unknown platform error
- **WHEN** origin remote URL does not match any known platform
- **THEN** an error is displayed with the detected URL
- **AND** guidance is provided for manual PR creation

#### Scenario: SSH URL format support
- **WHEN** origin remote uses SSH format (e.g., `git@github.com:org/repo.git`)
- **THEN** platform detection correctly identifies the host

#### Scenario: HTTPS URL format support
- **WHEN** origin remote uses HTTPS format (e.g., `https://github.com/org/repo.git`)
- **THEN** platform detection correctly identifies the host

### Requirement: Platform CLI Availability
The system SHALL verify that the required platform CLI tool is installed and authenticated before attempting PR creation.

#### Scenario: CLI not installed
- **WHEN** the required CLI tool (gh/glab/tea) is not installed
- **THEN** an error is displayed: "<tool> CLI is required for <platform> PR creation. Install: <install-url>"
- **AND** the command exits with code 1

#### Scenario: CLI not authenticated
- **WHEN** the required CLI tool is installed but not authenticated
- **THEN** an error is displayed with authentication instructions
- **AND** the command exits with code 1

### Requirement: Git Worktree Isolation
The PR commands SHALL use git worktrees to provide complete isolation from the user's working directory.

#### Scenario: Worktree created in temp directory
- **WHEN** PR workflow starts
- **THEN** a worktree is created in the system temp directory
- **AND** the worktree path includes a UUID to prevent conflicts

#### Scenario: Worktree based on origin branch
- **WHEN** worktree is created
- **THEN** it is based on the remote base branch (origin/main or origin/master)
- **AND** it does NOT include any local uncommitted changes

#### Scenario: Worktree cleanup on success
- **WHEN** PR workflow completes successfully
- **THEN** the temporary worktree is removed
- **AND** no temporary files remain

#### Scenario: Worktree cleanup on failure
- **WHEN** PR workflow fails at any stage
- **THEN** the temporary worktree is still removed
- **AND** an appropriate error message is displayed

#### Scenario: Git version requirement
- **WHEN** git version is less than 2.5
- **THEN** an error is displayed: "Git >= 2.5 required for worktree support. Current version: <version>"
- **AND** the command exits with code 1

### Requirement: PR Commit Message Format
The system SHALL generate structured commit messages that follow conventional commit format and include relevant metadata.

#### Scenario: Archive commit message format
- **WHEN** `spectr pr archive` commits changes
- **THEN** the commit message includes:
  - Title: `spectr(archive): <change-id>`
  - Archive location path
  - Spec operation counts (added, modified, removed, renamed)
  - Attribution: "Generated by: spectr pr archive"

#### Scenario: Proposal commit message format
- **WHEN** `spectr pr proposal` commits changes
- **THEN** the commit message includes:
  - Title: `spectr(proposal): <change-id>`
  - Proposal location path
  - Attribution: "Generated by: spectr pr proposal"

### Requirement: PR Body Content
The system SHALL generate PR body content that helps reviewers understand the change.

#### Scenario: Archive PR body content
- **WHEN** PR is created for archive
- **THEN** the PR body includes:
  - Summary section with change ID and archive location
  - Spec updates table with operation counts
  - List of updated capabilities
  - Review checklist

#### Scenario: Proposal PR body content
- **WHEN** PR is created for proposal
- **THEN** the PR body includes:
  - Summary section with change ID and location
  - List of included files (proposal.md, tasks.md, specs/)
  - Review checklist

### Requirement: PR Branch Naming
The system SHALL use a consistent branch naming convention for PR branches.

#### Scenario: Branch name format
- **WHEN** PR workflow creates a branch
- **THEN** the branch is named `spectr/<change-id>`

#### Scenario: Branch name with special characters
- **WHEN** change ID contains only valid kebab-case characters
- **THEN** the branch name is valid for git

### Requirement: PR Error Handling
The system SHALL provide clear error messages and guidance when PR creation fails.

#### Scenario: Not in git repository
- **WHEN** user runs `spectr pr` from outside a git repository
- **THEN** an error is displayed: "Not in a git repository"
- **AND** the command exits with code 1

#### Scenario: No origin remote
- **WHEN** user runs `spectr pr` and no origin remote exists
- **THEN** an error is displayed: "No 'origin' remote configured"
- **AND** guidance is provided to add a remote

#### Scenario: Change does not exist
- **WHEN** user runs `spectr pr <subcommand> <change-id>`
- **AND** the change does not exist
- **THEN** an error is displayed: "Change '<change-id>' not found"
- **AND** the command exits with code 1

#### Scenario: Push failure
- **WHEN** git push fails (e.g., network error)
- **THEN** an error is displayed with the git error message
- **AND** guidance is provided for manual recovery
- **AND** worktree is still cleaned up

#### Scenario: PR creation failure with pushed branch
- **WHEN** push succeeds but PR creation fails
- **THEN** an error is displayed with the PR CLI error
- **AND** a message indicates: "Branch was pushed. Create PR manually or retry."
- **AND** worktree is still cleaned up

### Requirement: Partial Change ID Resolution for PR Commands
The `spectr pr` subcommands SHALL support intelligent partial ID matching consistent with the archive command's resolution algorithm.

#### Scenario: Exact ID match for PR commands
- **WHEN** user runs `spectr pr archive exact-change-id`
- **AND** a change with ID `exact-change-id` exists
- **THEN** the PR workflow proceeds with that change

#### Scenario: Unique prefix match for PR commands
- **WHEN** user runs `spectr pr proposal refactor`
- **AND** only one change ID starts with `refactor`
- **THEN** a resolution message is displayed
- **AND** the PR workflow proceeds with the resolved ID

### Requirement: PR Proposal Interactive Selection Filters Unmerged Changes

The `spectr pr proposal` command's interactive selection mode SHALL only display changes that do not already exist on the target branch (main/master), ensuring users only see changes that genuinely need proposal PRs.

#### Scenario: Interactive list excludes changes on main

- **WHEN** user runs `spectr pr proposal` without a change ID argument
- **AND** some changes in `spectr/changes/` already exist on `origin/main`
- **THEN** only changes NOT present on `origin/main` are displayed in the interactive list
- **AND** changes that exist on main are filtered out before display

#### Scenario: All changes already on main

- **WHEN** user runs `spectr pr proposal` without a change ID argument
- **AND** all active changes already exist on `origin/main`
- **THEN** a message is displayed: "No unmerged proposals found. All changes already exist on main."
- **AND** the command exits gracefully without entering interactive mode

#### Scenario: No changes exist at all

- **WHEN** user runs `spectr pr proposal` without a change ID argument
- **AND** no changes exist in `spectr/changes/`
- **THEN** a message is displayed: "No changes found."
- **AND** the command exits gracefully

#### Scenario: Explicit change ID bypasses filter

- **WHEN** user runs `spectr pr proposal <change-id>` with an explicit argument
- **THEN** the filter is NOT applied
- **AND** the command proceeds with the specified change ID
- **AND** existing behavior is preserved for direct invocation

#### Scenario: Archive command unaffected

- **WHEN** user runs `spectr pr archive` without a change ID argument
- **THEN** all active changes are displayed in the interactive list
- **AND** no filtering based on main branch presence is applied
- **AND** existing archive behavior is preserved

#### Scenario: Detection uses git ls-tree

- **WHEN** the system checks if a change exists on main
- **THEN** it uses `git ls-tree` to check if `spectr/changes/<change-id>` path exists on `origin/main`
- **AND** the check is performed before displaying the interactive list
- **AND** fetch is performed first to ensure refs are current

#### Scenario: Custom base branch respected

- **WHEN** user runs `spectr pr proposal --base develop` without a change ID
- **THEN** the filter checks against `origin/develop` instead of `origin/main`
- **AND** only changes not present on `origin/develop` are displayed

### Requirement: Template Path Variables

The template rendering system SHALL support dynamic path variables for all directory and file references, allowing template content to be decoupled from specific path names while maintaining backward-compatible defaults.

The `TemplateContext` struct SHALL provide the following fields with default values:
- `BaseDir`: The root Spectr directory name (default: `spectr`)
- `SpecsDir`: The specifications directory path (default: `spectr/specs`)
- `ChangesDir`: The changes directory path (default: `spectr/changes`)
- `ProjectFile`: The project configuration file path (default: `spectr/project.md`)
- `AgentsFile`: The agents instruction file path (default: `spectr/AGENTS.md`)

#### Scenario: Templates use path variables instead of hardcoded strings

- **WHEN** a template file contains path references
- **THEN** the path SHALL be expressed using Go template syntax (e.g., `{{ .BaseDir }}`, `{{ .SpecsDir }}`)
- **AND** hardcoded `spectr/` strings SHALL NOT appear in template files for path references
- **AND** the rendered output SHALL contain the actual path values from the context

#### Scenario: Default context produces backward-compatible output

- **WHEN** `DefaultTemplateContext()` is used for rendering
- **THEN** the rendered output SHALL be identical to the previous hardcoded output
- **AND** all path references SHALL resolve to `spectr/`, `spectr/specs/`, `spectr/changes/`, etc.

#### Scenario: Template manager methods accept context parameter

- **WHEN** `RenderAgents()`, `RenderInstructionPointer()`, or `RenderSlashCommand()` is called
- **THEN** the method SHALL accept a `TemplateContext` parameter
- **AND** the context values SHALL be available within the template

#### Scenario: All template files use consistent variable names

- **WHEN** any template file references a Spectr path
- **THEN** it SHALL use the standardized variable names (`BaseDir`, `SpecsDir`, `ChangesDir`, `ProjectFile`, `AgentsFile`)
- **AND** variable names SHALL be consistent across all template files

### Requirement: PR Hotkey in Interactive Changes List Mode

The interactive changes list mode SHALL provide a `P` (Shift+P) hotkey that exits the TUI and enters the `spectr pr` workflow for the selected change, allowing users to create pull requests without manually copying the change ID.

#### Scenario: User presses Shift+P to enter PR mode

- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses the `P` key (Shift+P) on a selected change
- **THEN** the interactive mode exits
- **AND** the system enters PR mode for the selected change ID
- **AND** the user is prompted to select PR type (archive or proposal)

#### Scenario: PR hotkey not available in specs mode

- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses `P` key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show `P: pr` option

#### Scenario: PR hotkey not available in unified mode

- **WHEN** user is in unified interactive mode (`spectr list --all -I`)
- **AND** user presses `P` key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show `P: pr` option
- **AND** this avoids confusion when a spec row is selected

#### Scenario: Help text shows PR hotkey in changes mode

- **WHEN** user presses `?` in changes mode
- **THEN** the help text includes `P: pr` in the controls line
- **AND** the hotkey appears alongside other change-specific hotkeys (e, a)

#### Scenario: PR workflow integration

- **WHEN** the PR hotkey triggers the PR workflow
- **THEN** the workflow uses the same code path as `spectr pr`
- **AND** the selected change ID is passed to the PR workflow
- **AND** the user can select between archive and proposal modes

### Requirement: VHS Demo for PR Hotkey

The system SHALL provide a VHS tape demonstrating the Shift+P hotkey utility in the interactive list TUI.

#### Scenario: Developer finds PR hotkey demo

- **WHEN** a developer reviews the VHS tape files in `assets/vhs/`
- **THEN** they SHALL find `pr-hotkey.tape` demonstrating the PR hotkey workflow

#### Scenario: User sees PR hotkey demo

- **WHEN** a user views the PR hotkey demo GIF
- **THEN** they SHALL see the interactive list being invoked
- **AND** they SHALL see the `P` key being pressed
- **AND** they SHALL see the PR mode being entered for the selected change

### Requirement: PR Proposal Local Change Cleanup Confirmation

After a successful `spectr pr proposal` command that creates a pull request, the system SHALL prompt the user with a Bubbletea TUI confirmation menu asking whether to remove the local change proposal directory from `spectr/changes/`.

This prompt helps users maintain a clean working directory by offering an opportunity to remove proposals that have been submitted for review, while defaulting to "No" for safety. The menu uses arrow key navigation and styled rendering consistent with other spectr interactive modes.

#### Scenario: Successful PR proposal triggers cleanup prompt

- **WHEN** user runs `spectr pr proposal <change-id>` and PR creation succeeds
- **AND** the PR URL is displayed to the user
- **THEN** the system displays a Bubbletea TUI menu: "Remove local change proposal from spectr/changes/?"
- **AND** the menu shows two options: "No, keep it" and "Yes, remove it"
- **AND** the default selection is "No, keep it" (cursor starts on this option)
- **AND** the menu supports arrow key navigation (↑/↓) and Enter to confirm

#### Scenario: User confirms cleanup via TUI

- **WHEN** the cleanup TUI menu is displayed
- **AND** user navigates to "Yes, remove it" and presses Enter
- **THEN** the system removes the directory `spectr/changes/<change-id>/`
- **AND** displays confirmation: "Removed local change: <change-id>"

#### Scenario: User declines cleanup via TUI

- **WHEN** the cleanup TUI menu is displayed
- **AND** user presses Enter on the default "No, keep it" option
- **THEN** the system keeps the local change directory
- **AND** displays: "Local change kept: spectr/changes/<change-id>/"

#### Scenario: User cancels cleanup menu

- **WHEN** the cleanup TUI menu is displayed
- **AND** user presses 'q' or Ctrl+C
- **THEN** the system keeps the local change directory (same as declining)
- **AND** the command exits successfully

#### Scenario: Non-interactive mode skips prompt

- **WHEN** user runs `spectr pr proposal <change-id> --yes`
- **AND** PR creation succeeds
- **THEN** the cleanup prompt is NOT displayed
- **AND** the local change directory is kept (safe default)
- **AND** the command exits successfully

#### Scenario: Cleanup prompt only for proposal mode

- **WHEN** user runs `spectr pr archive <change-id>`
- **AND** PR creation succeeds
- **THEN** the cleanup prompt is NOT displayed
- **AND** the change is already moved to archive by the archive workflow

#### Scenario: PR creation fails

- **WHEN** user runs `spectr pr proposal <change-id>`
- **AND** PR creation fails at any step
- **THEN** the cleanup prompt is NOT displayed
- **AND** the local change directory remains intact

#### Scenario: Cleanup removal error handling

- **WHEN** the user confirms cleanup
- **AND** removal of the change directory fails (e.g., permission error)
- **THEN** display a warning: "Warning: Failed to remove change directory: <error>"
- **AND** the command still exits successfully (non-fatal error)
