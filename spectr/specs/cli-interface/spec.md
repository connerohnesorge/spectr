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
The interactive table SHALL use clear visual styling to distinguish headers, selected rows, and borders.

#### Scenario: Visual hierarchy in table
- **WHEN** interactive mode is displayed
- **THEN** column headers are visually distinct from data rows
- **AND** selected row has contrasting background/foreground colors
- **AND** table borders are visible and styled consistently
- **AND** table fits within terminal width gracefully

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
- **AND** the system creates `.claude/commands/spectr/archive.md`
- **AND** all files are tracked in the execution result
- **AND** the completion screen shows all 4 files created

#### Scenario: Multiple tools with slash commands selected

- **WHEN** user selects both `claude-code` and `cursor` in the init wizard
- **THEN** the system creates `CLAUDE.md` and both config + slash commands for Claude
- **AND** the system creates `.cursor/commands/spectr-proposal.md` and slash commands for Cursor
- **AND** all files from both tools are created and tracked separately
- **AND** the completion screen lists all created files grouped by tool

#### Scenario: Slash command files already exist

- **WHEN** user run init and selects `claude-code`
- **AND** `.claude/commands/spectr/proposal.md` already exists
- **THEN** the existing file's content between `<!-- spectr:START -->` and `<!-- spectr:END -->` is updated
- **AND** the file's YAML frontmatter is preserved
- **AND** no error occurs
- **AND** the file is marked as "updated" rather than "created" in execution result

#### Scenario: Config-based tool without slash mapping

- **WHEN** a config-based tool has no slash command equivalent in the mapping
- **THEN** only the config file is created
- **AND** no error occurs
- **AND** the system continues with remaining tool configurations

#### Scenario: Tool mapping is explicit and centralized

- **WHEN** a developer reviews the mapping logic
- **THEN** they find the tool mapping integrated into the tool registry configuration
- **AND** the registry uses data-driven tool definitions with type-safe IDs
- **AND** the mapping can be extended for new tools through configuration

#### Scenario: ToolDefinition structure simplified

- **WHEN** a developer reviews the ToolDefinition struct in `internal/init/models.go`
- **THEN** the struct contains: ID (type-safe ToolID), Name, Type, Priority, and Configured fields
- **AND** the struct does NOT contain a ConfigPath field
- **AND** file paths are determined by configurator implementations, not the registry
- **AND** the `getToolFileInfo()` function queries configurators for actual file paths

### Requirement: Archive Command PR Flag

The `spectr archive` command SHALL accept a `--pr` flag that triggers automated pull request creation after successful archive completion, integrating git branch creation, commit, push, and platform-specific PR CLI invocation into the archive workflow.

#### Scenario: User archives with PR flag

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** the change exists and archive operation succeeds
- **AND** git repository is configured with origin remote
- **THEN** the archive workflow completes (validation, spec merging, directory move)
- **AND** a git branch `archive-my-feature` is created
- **AND** archived files and updated specs are committed
- **AND** the branch is pushed to origin
- **AND** a PR is created using the detected platform CLI tool
- **AND** the PR URL is displayed

#### Scenario: PR flag combines with other archive flags

- **WHEN** user runs `spectr archive my-feature --pr --yes`
- **THEN** confirmation prompts are skipped
- **AND** PR creation proceeds automatically after archive

- **WHEN** user runs `spectr archive my-feature --pr --skip-specs`
- **THEN** spec updates are skipped
- **AND** PR is created with only archived directory
- **AND** PR body notes specs were skipped

- **WHEN** user runs `spectr archive my-feature --pr --interactive`
- **THEN** interactive change selection is used
- **AND** after selection, PR workflow proceeds

#### Scenario: PR flag without change ID uses interactive mode

- **WHEN** user runs `spectr archive --pr` with no change ID argument
- **THEN** interactive change selection is displayed
- **AND** after user selects a change, archive proceeds with PR creation
- **AND** the selected change is archived and PR is created

#### Scenario: PR flag error when archive fails

- **WHEN** user runs `spectr archive invalid --pr`
- **AND** the change does not exist or validation fails
- **THEN** an error is displayed about the archive failure
- **AND** no git operations are performed
- **AND** no branch or PR is created
- **AND** the command exits with error code 1

#### Scenario: PR flag error when git not available

- **WHEN** user runs `spectr archive my-feature --pr`
- **AND** archive succeeds but git is not available or no git repository exists
- **THEN** an error is displayed about git requirements
- **AND** the archive operation is complete
- **AND** changes remain uncommitted
- **AND** user is guided to manually commit or configure git
- **AND** the command exits with error code 1

#### Scenario: Help text documents PR flag

- **WHEN** user runs `spectr archive --help`
- **THEN** the `--pr` flag is listed in the available flags
- **AND** the description explains: "Create pull request after successful archive"
- **AND** the help text notes the dependency on git and PR CLI tools (gh/glab/tea)
- **AND** the configurator implementations use tool-specific configuration data

#### Scenario: Type-safe tool ID usage

- **WHEN** a developer references a tool ID in code
- **THEN** they use a defined ToolID constant (e.g., `ToolClaudeCode`)
- **AND** string literals for tool IDs trigger compiler warnings or errors
- **AND** IDE autocomplete suggests available tool ID constants
- **AND** typos in tool IDs are caught at compile time
