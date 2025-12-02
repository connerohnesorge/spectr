# Cli Interface Specification

## ADDED Requirements

### Requirement: Uncommitted Filter Hotkey in Interactive Changes Mode

The interactive changes list mode SHALL provide an 'h' hotkey that toggles a filter to show only changes that have uncommitted git modifications AND fully completed tasks.md files. This helps users identify changes that are ready to be committed or archived.

#### Scenario: User presses 'h' to enable uncommitted filter
- **WHEN** user is in interactive changes mode (`spectr list -I`)
- **AND** user presses the 'h' key
- **THEN** the table is filtered to show only changes where:
  - The change directory has uncommitted git modifications (untracked or modified files)
  - AND all tasks in tasks.md are marked complete (Completed == Total, with Total > 0)
- **AND** the footer updates to indicate filter is active (e.g., "filter: uncommitted+complete")
- **AND** the item count updates to reflect filtered results

#### Scenario: User presses 'h' again to disable filter
- **WHEN** uncommitted filter is active
- **AND** user presses the 'h' key
- **THEN** the filter is disabled
- **AND** all changes are displayed again
- **AND** the footer returns to normal display
- **AND** the cursor position is preserved if the previously selected item is still visible

#### Scenario: No changes match filter criteria
- **WHEN** user presses 'h' to enable filter
- **AND** no changes have both uncommitted modifications and complete tasks
- **THEN** the table displays no rows
- **AND** a message indicates "No uncommitted changes with complete tasks"
- **AND** user can press 'h' again to show all changes

#### Scenario: Git status detection for change directory
- **WHEN** the filter checks for uncommitted modifications
- **THEN** it runs `git status --porcelain` for the change directory
- **AND** considers a change "uncommitted" if any files in `spectr/changes/<id>/` appear in git status output
- **AND** includes both untracked files and modified files

#### Scenario: Complete tasks detection
- **WHEN** the filter checks for complete tasks
- **THEN** a change is considered "complete" if TaskStatus.Completed == TaskStatus.Total
- **AND** TaskStatus.Total must be greater than 0 (empty tasks.md does not count as complete)

#### Scenario: Help text shows uncommitted filter hotkey
- **WHEN** user presses '?' to show help in changes mode
- **THEN** the help text includes 'h: uncommitted' or 'h: filter uncommitted'
- **AND** the hotkey appears in the controls line

#### Scenario: Filter not available in specs mode
- **WHEN** user is in interactive specs mode (`spectr list --specs -I`)
- **AND** user presses 'h' key
- **THEN** the key press is ignored (no action taken)
- **AND** the help text does NOT show 'h' option
- **AND** specs do not have git status or tasks concepts

#### Scenario: Filter available in unified mode for changes only
- **WHEN** user is in unified interactive mode (`spectr list --all -I`)
- **AND** user presses 'h' key
- **THEN** the filter applies only to change items
- **AND** specs are hidden when filter is active (since they cannot match criteria)
- **AND** the filter can be toggled on/off as in changes mode

#### Scenario: Search and uncommitted filter work together
- **WHEN** uncommitted filter is active
- **AND** user activates search mode with '/'
- **THEN** search operates on the already-filtered (uncommitted+complete) changes
- **AND** both filters combine to narrow results
- **AND** clearing search restores the uncommitted filter view
