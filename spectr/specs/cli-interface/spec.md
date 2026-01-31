# Cli Interface Specification

## Purpose

Defines the CLI framework using Kong for declarative commands, struct tags,
subcommands (archive, list, validate, view), and interactive features (TUI,
wizard, visual styling).

## Requirements

### Requirement: Archive Command

The system SHALL provide an `archive` command to move completed changes and
apply deltas.

#### Scenario: Archive with change ID

- THEN archive change without prompting

#### Scenario: Interactive archive selection

- THEN prompt to select active change

#### Scenario: Non-interactive archiving with yes flag

- THEN archive without confirmation

#### Scenario: Skip spec updates for tooling changes

- THEN archive without updating specs

#### Scenario: Skip validation with confirmation

- THEN warn and require confirmation (unless --yes)

### Requirement: Positional Argument Support

The system SHALL support positional arguments via struct fields tagged `arg`.

#### Scenario: Optional positional argument

- THEN define with `arg` and `optional` tags (pointer or zero-value)

#### Scenario: Required positional argument

- THEN define with `arg` tag; fail if missing

### Requirement: Built-in Help Generation

The system SHALL generate help text from struct tags.

#### Scenario: Root help display

- THEN list subcommands, descriptions, and args

#### Scenario: Subcommand help display

- THEN show description, flags, types, and args

### Requirement: Error Handling and Exit Codes

The system SHALL provide clear errors and exit codes.

#### Scenario: Parse error handling

- THEN show error, usage, exit non-zero

#### Scenario: Execution error handling

- THEN show error, exit non-zero

### Requirement: Backward-Compatible CLI Interface

The system SHALL maintain syntax and flag compatibility.

#### Scenario: Init command compatibility

- THEN match previous Cobra behavior (flags, aliases, args)

#### Scenario: Help text accessibility

- THEN document all commands/flags

### Requirement: List Command for Changes

The system SHALL provide a `list` command that enumerates active changes (IDs default).

#### Scenario: List changes with IDs only

- THEN display IDs (sorted, excluding archive)

#### Scenario: List changes with details

- THEN format: `{id}: {title} [deltas {count}] [tasks {completed}/{total}]`

#### Scenario: List changes as JSON

- THEN output JSON array (id, title, deltaCount, taskStatus)

#### Scenario: No changes found

- THEN display "No items found"

### Requirement: List Command for Specs

The system SHALL support a `--specs` flag to list specifications.

#### Scenario: List specs with IDs only

- THEN display spec IDs

#### Scenario: List specs with details

- THEN format: `{id}: {title} [requirements {count}]`

#### Scenario: List specs as JSON

- THEN output JSON array (id, title, requirementCount)

#### Scenario: No specs found

- THEN display "No items found"

### Requirement: Change Discovery

The system SHALL scan `spectr/changes/` for subdirectories with `proposal.md`.

#### Scenario: Find active changes

- THEN include `spectr/changes/*/proposal.md`
- AND exclude `archive/` and hidden dirs

### Requirement: Spec Discovery

The system SHALL scan `spectr/specs/` for subdirectories with `spec.md`.

#### Scenario: Find specs

- THEN include `spectr/specs/*/spec.md`
- AND exclude hidden dirs

### Requirement: Title Extraction

The system SHALL extract title from first level-1 heading in markdown
(removing "Change:"/"Spec:" prefix).

#### Scenario: Extract title from proposal

- THEN extract title from `# Change: ...`

#### Scenario: Extract title from spec

- THEN extract title from `# ...`

#### Scenario: Fallback to ID when title not found

- THEN use directory name

### Requirement: Task Counting

The system SHALL count tasks from `tasks.jsonc` or `tasks.md` (ignore `tasks.json`).

#### Scenario: Count tasks from JSONC

- THEN count by status (strip comments)

#### Scenario: Count tasks from Markdown

- THEN count `- [ ]` and `- [x]`

#### Scenario: Ignore legacy tasks.json

- THEN report 0 tasks

#### Scenario: Handle missing tasks file

- THEN report 0 tasks

#### Scenario: JSONC takes precedence over Markdown

- THEN use `tasks.jsonc`

### Requirement: Validate Command

The system SHALL check spec/change correctness.

#### Scenario: Validate command registration

- THEN register `spectr validate`

#### Scenario: Direct item validation invocation

- THEN validate item, print results, exit 0/1

#### Scenario: Bulk validation invocation

- THEN validate all, print summary/issues, exit 1 on failure

#### Scenario: Interactive validation invocation

- THEN prompt for selection

#### Scenario: Default validation behavior (always strict)

- THEN treat warnings as errors (exit 1)

#### Scenario: JSON output flag

- THEN output structured JSON

#### Scenario: Type disambiguation flag

- THEN validate as type (error if missing)

#### Scenario: Changes only flag

- THEN validate changes only

#### Scenario: Specs only flag

- THEN validate specs only

#### Scenario: Non-interactive flag

- THEN no prompt, exit 1 if no item

### Requirement: Validate Command Help Text

The system SHALL document validation purpose and flags.

#### Scenario: Command help display

- THEN show description, flags, examples

### Requirement: Positional Argument Support for Item Name

The system SHALL accept an optional item name.

#### Scenario: Optional item name argument

- THEN `arg:"" optional:""`

#### Scenario: Item name provided

- THEN validate specific item (auto-detect type)

### Requirement: View Command

The system SHALL display a project dashboard (`spectr view`).

#### Scenario: View command registration

- THEN register `spectr view`

#### Scenario: View command invocation

- THEN show dashboard (summary, active/completed changes, specs)

#### Scenario: View command with JSON output

- THEN output dashboard JSON

#### Scenario: JSON structure

- THEN include `summary`, `activeChanges`, `completedChanges`, `specs`

#### Scenario: JSON arrays sorted consistently

- THEN sort active (progress, ID), completed (ID), specs (reqs, ID)

#### Scenario: JSON with no items

- THEN return empty arrays

### Requirement: Dashboard Summary Metrics

The system SHALL display aggregated stats.

#### Scenario: Display summary with all metrics

- THEN show counts: specs, requirements, active/completed changes, tasks

#### Scenario: Calculate total requirements

- THEN sum requirements from all specs

#### Scenario: Calculate task progress

- THEN sum tasks/completed from active changes

### Requirement: Active Changes Display

The system SHALL show active changes with progress bars.

#### Scenario: List active changes with progress

- THEN show ID, progress bar, %, indicator (◉)
- AND sort by completion asc, ID

#### Scenario: Render progress bar

- THEN 20 chars, full/light blocks, green/gray

#### Scenario: Handle zero tasks

- THEN empty bar, 0%

#### Scenario: No active changes

- THEN hide section

### Requirement: Completed Changes Display

The system SHALL show completed changes.

#### Scenario: List completed changes

- THEN show ID, checkmark (✓), sort ID

#### Scenario: Determine completion status

- THEN complete if all tasks done or 0 tasks

#### Scenario: No completed changes

- THEN hide section

### Requirement: Specifications Display

The system SHALL show specs sorted by complexity.

#### Scenario: List specifications with requirement counts

- THEN show ID, count, indicator (▪), sort count desc

#### Scenario: Pluralize requirement label

- THEN handle singular/plural

#### Scenario: No specifications found

- THEN hide section

### Requirement: Dashboard Visual Formatting

The system SHALL use colors and box-drawing chars.

#### Scenario: Render dashboard header

- THEN title, double-line separator

#### Scenario: Render section headers

- THEN bold cyan, single-line separator

#### Scenario: Render footer

- THEN double-line separator, hints

#### Scenario: Color scheme consistency

- THEN cyan (headers), yellow (active), green (done), blue (specs)

### Requirement: Sorting Strategy

The system SHALL sort for relevance.

#### Scenario: Sort active changes by priority

- THEN % asc, ID

#### Scenario: Sort specs by complexity

- THEN count desc, ID

#### Scenario: Sort completed changes alphabetically

- THEN ID

### Requirement: View Command Help Text

The system SHALL document the view command.

#### Scenario: View command help display

- THEN show description, `--json`

### Requirement: Provider Interface

The system SHALL provide a `Provider` interface for AI tools (instruction
files + slash commands).

#### Scenario: Provider interface methods

- THEN implement ID, Name, Priority, ConfigFile, paths, Configure, IsConfigured

#### Scenario: Single provider per tool

- THEN one provider handles both

#### Scenario: Flexible command paths

- THEN relative path or empty

#### Scenario: HasSlashCommands detection

- THEN true if any path non-empty

### Requirement: Command Format Support

The system SHALL support Markdown/TOML formats.

#### Scenario: Markdown command format

- THEN create `.md` with frontmatter

#### Scenario: TOML command format

- THEN create `.toml` with description/prompt

### Requirement: Version Command Structure

The system SHALL display version/commit/date.

#### Scenario: Version command registration

- THEN register `spectr version`

#### Scenario: Version command invocation

- THEN show version info

#### Scenario: Version command with short flag

- THEN show version only

#### Scenario: Version command with JSON flag

- THEN output JSON

### Requirement: Version Variable Injection

The system SHALL inject version variables via ldflags.

#### Scenario: Goreleaser/Nix/Dev injection

- THEN set version/commit/date (default dev/unknown)

### Requirement: Completion Command Structure

The system SHALL generate shell completions.

#### Scenario: Completion command registration

- THEN register `spectr completion`

#### Scenario: Bash/Zsh/Fish completion output

- THEN output script

### Requirement: Custom Predictors for Dynamic Arguments

The system SHALL suggest IDs.

#### Scenario: Change/Spec ID completion

- THEN suggest active changes/specs

#### Scenario: Item type completion

- THEN suggest `change`, `spec`

### Requirement: Accept Command Structure

The CLI SHALL provide an `accept` command that converts `tasks.md` to
`tasks.jsonc` format with header comments for stable agent manipulation during
implementation.

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

- WHEN user runs `spectr accept <change-id>` on a change that already has
  tasks.jsonc
- THEN the system displays a message indicating change is already accepted
- AND the system exits with code 0 (success, idempotent)

#### Scenario: Accept change without tasks.md

- WHEN user runs `spectr accept <change-id>` on a change without tasks.md
- THEN the system displays an error indicating no tasks.md found
- AND the system exits with code 1

### Requirement: Tasks JSON Schema

The system SHALL generate a versioned `tasks.jsonc` with optional hierarchical
structure for large task files, automatically splitting files over 100 lines
while preserving section boundaries.

#### Scenario: JSONC file structure

- **WHEN** generating tasks.jsonc
- **THEN** the file SHALL start with header comments documenting status values
  and transitions
- **AND** the file SHALL contain a root object with `version` and `tasks` fields
- **AND** `version` SHALL be integer 2 for hierarchical files, integer 1 for
  flat files
- **AND** `tasks` SHALL be an array of task objects

#### Scenario: Version 2 hierarchical structure

- **WHEN** tasks.md exceeds 100 lines during accept
- **THEN** the system SHALL generate version 2 tasks.jsonc with split files
- **AND** root file SHALL contain reference tasks with `children` field
- **AND** child files SHALL be named `tasks-N.jsonc` where N is sequential
- **AND** root file SHALL include `includes` field with glob pattern
  `["tasks-*.jsonc"]`

#### Scenario: Version 1 flat structure preserved

- **WHEN** tasks.md is 100 lines or fewer
- **THEN** the system SHALL generate version 1 tasks.jsonc with flat structure
- **AND** no split files SHALL be created

#### Scenario: Header comment content

- **WHEN** generating root tasks.jsonc
- **THEN** the header SHALL use `//` line comment syntax
- **AND** the header SHALL indicate the file is machine-generated by `spectr
  accept`
- **AND** the header SHALL document the three valid status values: "pending",
  "in_progress", "completed"
- **AND** the header SHALL document valid status transitions: pending ->
  in_progress -> completed
- **AND** the header SHALL explain that agents should mark a task "in_progress"
  when starting work
- **AND** the header SHALL explain that agents should mark a task "completed"
  only after verification
- **AND** the header SHALL note that skipping directly from "pending" to
  "completed" is allowed for trivial tasks

#### Scenario: Child file header content

- **WHEN** generating child tasks-N.jsonc files
- **THEN** the header SHALL include the full header from root file
- **AND** the header SHALL add "Generated by: spectr accept `change-id`"
- **AND** the header SHALL add "Parent change: `change-id`"
- **AND** the header SHALL add "Parent task: `parent-task-id`"

#### Scenario: Task object structure

- **WHEN** a task is serialized to JSONC
- **THEN** it SHALL have `id` field containing the task identifier (e.g., "1.1")
- **AND** it SHALL have `section` field containing the section header (e.g.,
  "Implementation")
- **AND** it SHALL have `description` field containing the full task text
- **AND** it SHALL have `status` field with value "pending", "in_progress", or
  "completed"
- **AND** it MAY have `children` field with value "$ref:`path`" for reference
  tasks

#### Scenario: Status value mapping from Markdown

- **WHEN** converting tasks from markdown
- **THEN** map `- [ ]` to status "pending"
- **AND** map `- [x]` to status "completed"

#### Scenario: Reference task structure

- **WHEN** a task references child tasks in a split file
- **THEN** it SHALL have all standard fields (id, section, description, status)
- **AND** it SHALL have `children` field with value "$ref:tasks-N.jsonc"
- **AND** the description SHALL be the section header or summary
- **AND** the status SHALL reflect aggregate status of children (pending if any
  pending, in_progress if any in_progress, completed if all completed)

#### Scenario: Child file structure

- **WHEN** generating a child tasks-N.jsonc file
- **THEN** it SHALL have `version` field set to 2
- **AND** it SHALL have `parent` field containing the parent task ID
- **AND** it SHALL have `tasks` array with child task objects
- **AND** child task IDs SHALL use dot notation (e.g., "1.1", "1.2" for parent
  "1")

#### Scenario: Includes field for discovery

- **WHEN** generating version 2 root tasks.jsonc
- **THEN** it SHALL include `includes` field with value `["tasks-*.jsonc"]`
- **AND** agents MAY use this glob pattern to discover all split files

### Requirement: Accept Command Flags

The system SHALL provide flags to control behavior.

#### Scenario: Dry-run flag

- THEN preview JSON

#### Scenario: Interactive change selection

- THEN prompt list

### Requirement: List Command Alias

The system SHALL alias `ls` to `list`.

#### Scenario: User runs spectr ls shorthand

- THEN same as `list` (flags work)

### Requirement: Item Name Path Normalization

The system SHALL normalize paths to IDs.

#### Scenario: Path/ID normalization

- THEN extract ID, infer type

### Requirement: Interactive List Mode

The system SHALL provide a unified TUI for changes/specs.

#### Scenario: Default behavior unchanged

- THEN show changes only

#### Scenario: Unified mode opt-in

- THEN show changes and specs (Type column)

#### Scenario: Type-specific actions

- THEN only specs editable

#### Scenario: Help text formatting

- THEN minimal footer, `?` for full help

### Requirement: Clipboard Copy on Selection

The system SHALL copy ID on Enter.

#### Scenario: Copy ID

- THEN copy ID to clipboard (exit)

### Requirement: Interactive Mode Exit Controls

The system SHALL provide standard quit controls.

#### Scenario: Quit

- THEN exit

### Requirement: Table Visual Styling

The system SHALL provide consistent TUI styling.

#### Scenario: Visual hierarchy

- THEN styled headers, selection, borders (`tui.ApplyTableStyles`)

### Requirement: Cross-Platform Clipboard Support

The system SHALL support Linux/macOS/Windows/SSH.

#### Scenario: Clipboard support

- THEN use native API or OSC 52 fallback

### Requirement: Initialization Next Steps Message

The system SHALL guide the user after init.

#### Scenario: Next steps display

- THEN show 3 steps (project.md, proposal, AGENTS.md)

#### Scenario: Init does not create README

- THEN do not create README.md

### Requirement: Flat Tool List in Initialization Wizard

The system SHALL provide a unified tool list.

#### Scenario: Display only config-based tools

- THEN show config tools only (auto-install slash commands)

#### Scenario: Navigation/Selection

- THEN navigate flat list, space toggle, a/n bulk

### Requirement: Interactive Archive Mode

The system SHALL provide a table interface for archive.

#### Scenario: Archive no args / -I

- THEN show interactive table (same columns as list)

#### Scenario: Selection behavior

- THEN capture ID (no copy), proceed to archive

### Requirement: Archive Interactive Table Display

The system SHALL match list columns in archive mode.

#### Scenario: Table columns

- THEN ID, Title, Deltas, Tasks (consistent style)

### Requirement: Archive Selection Without Clipboard

The system SHALL capture selection internally only.

#### Scenario: Enter key

- THEN proceed with ID (no clipboard)

### Requirement: Validation Output Format

The system SHALL provide consistent issue reporting.

#### Scenario: Single/Bulk output

- THEN show issues (Level, Path, Message), summary

#### Scenario: JSON output

- THEN structured issue data

### Requirement: Editor Hotkey in Interactive Specs List

The system SHALL allow editing specs with 'e'.

#### Scenario: Edit spec

- THEN open `$EDITOR` (wait), return to TUI

### Requirement: Editor Hotkey Scope

The system SHALL limit editor hotkey to specs only.

#### Scenario: No edit for changes

- THEN ignore

### Requirement: Project Path Display in Interactive Mode

The system SHALL show context by displaying the project path.

#### Scenario: Path display

- THEN show project root path

### Requirement: Unified Item List Display

The system SHALL provide a mixed table for changes/specs.

#### Scenario: Unified display

- THEN show Type, ID, Title, Details

### Requirement: Type-Aware Item Selection

The system SHALL handle types correctly.

#### Scenario: Selection

- THEN copy ID (both types)

#### Scenario: Edit restriction

- THEN specs only

### Requirement: Backward-Compatible Separate Modes

The system SHALL preserve existing modes.

#### Scenario: Separate modes

- WHEN `--specs -I` -> Specs

### Requirement: Enhanced List Command Flags

The system SHALL validate flags.

#### Scenario: Validation

- THEN error

### Requirement: Automatic Slash Command Installation

The system SHALL install slash commands with config.

#### Scenario: Auto-install

- THEN install slash commands (CLAUDE.md + .claude/...)

### Requirement: Archive Hotkey in Interactive Changes Mode

The system SHALL allow archiving with 'a'.

#### Scenario: Archive action

- THEN exit and archive selected change

### Requirement: Shared TUI Component Library

The system SHALL use `internal/tui`.

#### Scenario: Components

- THEN use `TablePicker`, `MenuPicker`, `TruncateString`, `CopyToClipboard`

### Requirement: Search Hotkey in Interactive Lists

The system SHALL filter lists with '/'.

#### Scenario: Search mode

- THEN input field, filter rows by ID/Title

### Requirement: Help Toggle Hotkey

The system SHALL toggle help with '?'.

#### Scenario: Help toggle

- THEN toggle full/minimal help

### Requirement: Partial Change ID Resolution for Archive Command

The system SHALL resolve prefixes/substrings.

#### Scenario: Resolution

- THEN resolve if unique (prefix > substring), else error

### Requirement: Configured Provider Detection in Init Wizard

The system SHALL detect existing config.

#### Scenario: Detection

- THEN mark/select existing providers

### Requirement: Instruction File Pointer Template

The system SHALL use pointers in root files.

#### Scenario: Pointer content

- THEN point to `spectr/AGENTS.md`

### Requirement: PR Archive Subcommand Alias

The system SHALL alias `a` to `archive`.

#### Scenario: Alias

- THEN same as `spectr pr archive`

### Requirement: PR Branch Naming Convention

The system SHALL use consistent branch naming.

#### Scenario: Naming

- WHEN proposal -> `spectr/proposal/<id>`

### Requirement: PR Command Structure

The system SHALL support `spectr pr <subcommand>`.

#### Scenario: Subcommands

- THEN execute workflow

### Requirement: PR Archive Subcommand

The system SHALL create an archive PR (isolated worktree).

#### Scenario: Workflow

- THEN worktree, archive --yes, commit, push, create PR

### Requirement: PR Proposal Subcommand

The system SHALL create a proposal PR (isolated worktree).

#### Scenario: Proposal workflow

- THEN worktree, copy change, commit, push, create PR (no archive)

### Requirement: PR Common Flags

The system SHALL provide shared flags.

#### Scenario: Flags

- THEN apply behavior

### Requirement: Git Platform Detection

The system SHALL detect GitHub/GitLab/Gitea/Bitbucket.

#### Scenario: Platform detection

- THEN detect platform, use CLI (gh/glab/tea)

### Requirement: Platform CLI Availability

The system SHALL check if CLI is installed.

#### Scenario: Check

- THEN error if CLI missing/unauthenticated

### Requirement: Git Worktree Isolation

The system SHALL isolate operations.

#### Scenario: Isolation

- THEN use temp worktree, clean up after

### Requirement: PR Commit Message Format

The system SHALL use conventional commits.

#### Scenario: Format

- THEN `spectr(archive/proposal): <id>`

### Requirement: PR Body Content

The system SHALL provide a useful description.

#### Scenario: Content

- THEN summary, checklist, location

### Requirement: PR Branch Naming

The system SHALL use branch naming pattern `spectr/<id>`.

#### Scenario: Branch naming

- THEN `spectr/<id>`

### Requirement: PR Error Handling

The system SHALL provide clear errors.

#### Scenario: Error handling

- THEN display specific error

### Requirement: Partial Change ID Resolution for PR Commands

The system SHALL resolve IDs.

#### Scenario: ID resolution

- THEN resolve unique match

### Requirement: PR Proposal Interactive Selection Filters Unmerged Changes

The system SHALL filter already merged changes.

#### Scenario: Filtering

- THEN hide changes already on main

### Requirement: Template Path Variables

The system SHALL support dynamic paths.

#### Scenario: Variables

- THEN use `{{ .BaseDir }}` etc

### Requirement: Copy Populate Context Prompt in Init Next Steps

The system SHALL copy prompt with 'c'.

#### Scenario: Copy prompt

- THEN copy prompt to clipboard

### Requirement: PR Hotkey in Interactive Changes List Mode

The system SHALL enable PR workflow with 'P'.

#### Scenario: PR action

- THEN enter PR workflow

### Requirement: VHS Demo for PR Hotkey

The system SHALL include a demo asset.

#### Scenario: Demo

- THEN `pr-hotkey.tape` exists

### Requirement: PR Proposal Local Change Cleanup Confirmation

The system SHALL prompt to remove local changes.

#### Scenario: Prompt

- THEN prompt remove local (default No)

### Requirement: CI Workflow Setup Option in Init Wizard Review Step

The system SHALL setup GitHub Actions.

#### Scenario: CI option

- THEN checkbox for `.github/workflows/spectr-ci.yml`

### Requirement: PR Remove Subcommand

The system SHALL remove change PR.

#### Scenario: Remove workflow

- THEN worktree, remove dir, commit, PR, clean local

### Requirement: Remove PR Branch Naming

The system SHALL use branch naming pattern `spectr/remove/<id>`.

#### Scenario: Remove branch naming

- THEN `spectr/remove/<id>`

### Requirement: Remove PR Commit Message Format

The system SHALL use commit message format `spectr(remove): <id>`.

#### Scenario: Remove commit format

- THEN structured message

### Requirement: Remove PR Body Content

The system SHALL explain removal in PR body.

#### Scenario: Remove PR content

- THEN summary, removed path

### Requirement: Responsive Table Column Layout

The system SHALL adapt columns to width.

#### Scenario: Responsive columns

- THEN hide/narrow columns based on width

### Requirement: Dynamic Terminal Resize Handling

The system SHALL handle terminal resize.

#### Scenario: Resize

- THEN recalculate layout

### Requirement: Column Priority System

The system SHALL prioritize columns.

#### Scenario: Priority

- THEN ID > Title > Deltas/Reqs > Tasks

### Requirement: Provider Search in Init Wizard

The system SHALL search tools with '/'.

#### Scenario: Search

- THEN filter tools

### Requirement: Stdout Output Mode for Interactive List

The system SHALL output ID to stdout.

#### Scenario: Stdout mode

- THEN print ID (no clipboard)

### Requirement: JSONC Comment Parsing

The system SHALL strip comments.

#### Scenario: Parsing

- THEN strip `//` and `/* */`

### Requirement: TTY Error Hint

The system SHALL provide hints for non-TTY.

#### Scenario: Hint

- THEN suggest `--non-interactive`

### Requirement: File Coexistence Documentation

The system SHALL document tasks.md/jsonc coexistence.

#### Scenario: Docs

- THEN mention coexistence

### Requirement: Slash Command Template Updates

The system SHALL provide instructions for tasks.

#### Scenario: Templates

- THEN instruct on tasks.md/jsonc usage

### Requirement: Project Configuration File

The system SHALL support an optional `spectr.yaml` configuration file at the
project root.

#### Scenario: Config file present

- **WHEN** `spectr.yaml` exists in project root
- **THEN** load and parse the configuration

#### Scenario: Config file absent

- **WHEN** `spectr.yaml` does not exist
- **THEN** proceed with default behavior (no appended tasks)

#### Scenario: Config file malformed

- **WHEN** `spectr.yaml` contains invalid YAML
- **THEN** display error message and exit non-zero

### Requirement: Append Tasks Configuration

The system SHALL support an `append_tasks` section in `spectr.yaml` with a
configurable section name and list of tasks.

#### Scenario: Valid append_tasks configuration

- **WHEN** config contains `append_tasks.section` and `append_tasks.tasks`
- **THEN** parse section name as string and tasks as list of strings

#### Scenario: Missing section name

- **WHEN** `append_tasks.tasks` exists but `append_tasks.section` is missing
- **THEN** use default section name "Automated Tasks"

#### Scenario: Empty tasks list

- **WHEN** `append_tasks.tasks` is empty or missing
- **THEN** do not append any tasks

### Requirement: Auto-Append Tasks on Accept

The system SHALL append configured tasks to `tasks.jsonc` during
`spectr accept`.

#### Scenario: Append tasks with configured section

- **WHEN** `spectr accept <id>` runs with valid `append_tasks` config
- **THEN** append tasks to `tasks.jsonc` under the configured section name
- **AND** generate sequential task IDs continuing from the last task

#### Scenario: Task ID generation for appended tasks

- **WHEN** appending tasks after existing tasks (e.g., last ID was 3.2)
- **THEN** start appended tasks at next section number (e.g., 4.1, 4.2)

#### Scenario: No config present during accept

- **WHEN** `spectr accept <id>` runs without `spectr.yaml`
- **THEN** produce identical output to current behavior (no appended tasks)

### Requirement: Pre-Command Task Sync Hook

The system SHALL synchronize task statuses from `tasks.jsonc` to `tasks.md`
before every command execution using Kong's BeforeRun hook pattern.

#### Scenario: Sync runs before command execution

- **WHEN** any spectr subcommand is invoked
- **THEN** the system SHALL sync all active changes' task statuses before
  executing the command

#### Scenario: Sync updates only status markers

- **WHEN** syncing tasks.md from tasks.jsonc
- **THEN** the system SHALL update only checkbox markers (`[ ]` to `[x]` or
  vice versa)
- **AND** preserve all other markdown content (comments, links, formatting)

#### Scenario: Sync matches tasks by ID

- **WHEN** matching tasks between files
- **THEN** the system SHALL match by task ID (e.g., `1.1`, `2.3`)
- **AND** handle flexible ID formats (decimal, dot-suffixed, number-only)

#### Scenario: Sync handles missing files gracefully

- **WHEN** tasks.jsonc does not exist for a change
- **THEN** the system SHALL skip sync for that change silently

#### Scenario: Sync handles missing tasks.md gracefully

- **WHEN** tasks.md does not exist but tasks.jsonc does
- **THEN** the system SHALL skip sync for that change silently

#### Scenario: tasks.jsonc is source of truth for status

- **WHEN** a task status differs between files
- **THEN** the system SHALL use the status from tasks.jsonc
- **AND** update tasks.md to match

### Requirement: Global No-Sync Flag

The system SHALL provide a `--no-sync` global flag to disable automatic task
synchronization.

#### Scenario: Skip sync with flag

- **WHEN** `spectr --no-sync <command>` is invoked
- **THEN** the system SHALL skip the pre-command sync subroutine

#### Scenario: Flag applies to all subcommands

- **WHEN** `--no-sync` is provided
- **THEN** the system SHALL skip sync regardless of which subcommand follows

### Requirement: Silent Sync by Default

The system SHALL perform sync operations silently unless verbose mode is
enabled or errors occur.

#### Scenario: No output on successful sync

- **WHEN** sync completes successfully
- **THEN** the system SHALL produce no output

#### Scenario: Error output on sync failure

- **WHEN** sync encounters an error
- **THEN** the system SHALL print the error to stderr
- **AND** continue with command execution (non-blocking)

#### Scenario: Verbose output with flag

- **WHEN** `--verbose` flag is set and sync makes changes
- **THEN** the system SHALL print "Synced N task statuses in [change-id]"

### Requirement: Global Verbose Flag

The system SHALL provide a `--verbose` global flag to enable detailed output
including sync operations.

#### Scenario: Verbose flag registration

- **WHEN** CLI is initialized
- **THEN** the system SHALL register `--verbose` as a global flag

#### Scenario: Verbose sync output

- **WHEN** `spectr --verbose <command>` syncs tasks
- **THEN** the system SHALL print sync details for each affected change

### Requirement: Active Changes Only Sync

The system SHALL only synchronize task files in active (non-archived) changes.

#### Scenario: Exclude archived changes

- **WHEN** discovering changes to sync
- **THEN** the system SHALL exclude `spectr/changes/archive/` subdirectories

#### Scenario: Include all active changes

- **WHEN** discovering changes to sync
- **THEN** the system SHALL include all directories in `spectr/changes/` with
  `tasks.jsonc`

### Requirement: Status Mapping for Sync

The system SHALL map task statuses between jsonc and markdown formats
correctly.

#### Scenario: Pending and in_progress map to unchecked

- **WHEN** status is `pending` or `in_progress` in tasks.jsonc
- **THEN** the system SHALL write `[ ]` in tasks.md

#### Scenario: Completed maps to checked

- **WHEN** status is `completed` in tasks.jsonc
- **THEN** the system SHALL write `[x]` in tasks.md

### Requirement: Proposal Frontmatter Support

The system SHALL support YAML frontmatter in `proposal.md` files for declaring
metadata including dependencies.

#### Scenario: Parse valid frontmatter

- **WHEN** proposal.md contains YAML frontmatter block (`---` delimiters)

- **THEN** extract and parse frontmatter into ProposalMetadata

- **AND** make metadata available for validation and commands

#### Scenario: Handle missing frontmatter

- **WHEN** proposal.md has no frontmatter block

- **THEN** treat as valid with empty metadata (backward compatible)

#### Scenario: Handle malformed frontmatter

- **WHEN** proposal.md has invalid YAML in frontmatter

- **THEN** emit validation error with line number and details

### Requirement: Chained Proposal Dependencies

The system SHALL support `requires` and `enables` fields in proposal frontmatter
to declare dependencies between proposals.

#### Scenario: Declare required dependencies

- **WHEN** frontmatter contains `requires` list with id and optional reason

- **THEN** parse as list of Dependency objects

- **AND** validate that referenced change IDs exist or are archived

#### Scenario: Declare enabled proposals

- **WHEN** frontmatter contains `enables` list with id and optional reason

- **THEN** parse as informational metadata (not enforced)

#### Scenario: Self-reference detection

- **WHEN** proposal requires itself

- **THEN** emit validation error "Proposal cannot require itself"

### Requirement: Dependency Validation in Validate Command

The system SHALL check proposal dependencies during `spectr validate`.

#### Scenario: Warn on unmet dependencies

- **WHEN** `spectr validate <id>` runs on proposal with requires

- **AND** a required proposal is not archived

- **THEN** emit warning "Dependency '[id]' is not yet archived"

#### Scenario: Error on circular dependencies

- **WHEN** `spectr validate <id>` detects a dependency cycle

- **THEN** emit error "Circular dependency detected: A → B → ... → A"

- **AND** exit non-zero

#### Scenario: Pass validation when dependencies met

- **WHEN** all required proposals are archived

- **AND** no circular dependencies exist

- **THEN** validation passes (no dependency warnings/errors)

### Requirement: Dependency Enforcement in Accept Command

The system SHALL enforce `requires` dependencies before accepting a proposal.

#### Scenario: Block accept with unmet dependencies

- **WHEN** `spectr accept <id>` runs on proposal with requires

- **AND** any required proposal is not archived

- **THEN** exit non-zero with error listing unmet dependencies

- **AND** do not create tasks.jsonc

#### Scenario: Allow accept when dependencies met

- **WHEN** `spectr accept <id>` runs on proposal with requires

- **AND** all required proposals are archived

- **THEN** proceed with normal accept flow

#### Scenario: Accept proposal without dependencies

- **WHEN** `spectr accept <id>` runs on proposal without requires

- **THEN** proceed with normal accept flow (backward compatible)

### Requirement: Graph Command

The system SHALL provide a `spectr graph` command to visualize proposal
dependencies.

#### Scenario: Graph command registration

- **WHEN** CLI is initialized

- **THEN** register `spectr graph` command

#### Scenario: Display ASCII dependency tree

- **WHEN** `spectr graph [id]`

- **THEN** display ASCII tree showing requires/enables relationships

- **AND** indicate archived status with checkmark (✓) or pending (⧖)

#### Scenario: Display single proposal graph

- **WHEN** `spectr graph <id>`

- **THEN** show dependency tree rooted at specified proposal

#### Scenario: Display full dependency graph

- **WHEN** `spectr graph` (no args)

- **THEN** show all proposals and their relationships

#### Scenario: DOT format output

- **WHEN** `spectr graph --dot`

- **THEN** output Graphviz DOT format for external visualization

#### Scenario: JSON format output

- **WHEN** `spectr graph --json`

- **THEN** output JSON with nodes and edges arrays

### Requirement: Archive Status Detection

The system SHALL detect whether a change proposal has been archived.

#### Scenario: Detect archived change

- **WHEN** checking if change-id is archived

- **THEN** search `spectr/changes/archive/` for matching directory

- **AND** handle date prefix format (YYYY-MM-DD-[id])

#### Scenario: Detect active change

- **WHEN** change-id exists in `spectr/changes/` (not archive)

- **THEN** report as not archived

#### Scenario: Handle unknown change ID

- **WHEN** change-id does not exist anywhere

- **THEN** report as unknown/not found

### Requirement: Hierarchical Tasks Structure

The tasks.jsonc format SHALL support hierarchical task organization where root
tasks can reference child tasks in delta spec directories.

#### Scenario: Root file with summary

- WHEN a tasks.jsonc uses hierarchical structure (version 2)
- THEN the root object MAY contain a `summary` field
- AND `summary` SHALL contain `total`, `completed`, `in_progress`, and
  `pending` counts
- AND counts SHALL reflect aggregated status across all referenced child
  files

#### Scenario: Child task references

- WHEN a task has `children` field with value
  "$ref:specs/`<capability>`/tasks.jsonc"
- THEN the referenced file SHALL contain tasks that are children of this
  task
- AND child task IDs SHALL use the parent ID as prefix (e.g., parent "5" has
  children "5.1", "5.2")

#### Scenario: Auto-discovery via includes

- WHEN a tasks.jsonc contains `includes` field with glob patterns
- THEN the system SHALL discover matching task files automatically
- AND discovered files SHALL be included after explicit children references
- AND duplicate files (by path) SHALL be processed only once

#### Scenario: Child file structure

- WHEN a delta spec has its own tasks.jsonc (`specs/<capability>/tasks.jsonc`)
- THEN it SHALL contain a `parent` field with the parent task ID
- AND it SHALL contain a `tasks` array with child tasks
- AND it SHALL use version 2 format
- AND the header comments MAY be abbreviated (status values only)

#### Scenario: Hierarchical ID schema

- WHEN tasks use hierarchical IDs
- THEN root-level tasks SHALL use single numbers (e.g., "1", "2", "5")
- AND first-level children SHALL use parent.child format (e.g., "5.1",
  "5.2")
- AND deeper nesting SHALL extend the pattern (e.g., "5.1.1")
- AND maximum nesting depth SHALL be 2 levels (root -> capability)

### Requirement: Auto-Split on Accept

The `spectr accept` command SHALL automatically split tasks into
capability-specific files when section names match delta spec directories.

#### Scenario: Section matches delta spec

- WHEN accepting a change with section "5. Support Aider"
- AND delta spec directory `specs/support-aider/` exists
- THEN tasks from that section SHALL be written to
  `specs/support-aider/tasks.jsonc`
- AND root tasks.jsonc SHALL contain a reference task with `children`
  field

#### Scenario: Section name normalization

- WHEN matching section names to delta spec directories
- THEN section names SHALL be normalized to kebab-case
- AND leading numbers and punctuation SHALL be stripped (e.g.,
  "5. Support Aider" -> "support-aider")
- AND matching SHALL be case-insensitive

#### Scenario: No matching delta spec

- WHEN a section does not match any delta spec directory
- THEN tasks from that section SHALL remain in the root tasks.jsonc
- AND they SHALL be written as flat tasks (no children reference)

#### Scenario: Mixed structure

- WHEN some sections match delta specs and others do not
- THEN the root tasks.jsonc SHALL contain:
  - Flat tasks for non-matching sections
  - Reference tasks with children for matching sections
- AND the appropriate child files SHALL be created

#### Scenario: Auto-discovery pattern

- WHEN hierarchical tasks.jsonc is generated
- THEN root file SHALL include `"includes": ["specs/*/tasks.jsonc"]` for
  auto-discovery
- AND this enables agents to find child files without explicit
  references

### Requirement: Tasks Command

The CLI SHALL provide a `tasks` command for viewing task status across
hierarchical task files.

#### Scenario: Tasks command registration

- WHEN the CLI is initialized
- THEN it SHALL include a TasksCmd struct field tagged with `cmd`
- AND the command SHALL be accessible via `spectr tasks`

#### Scenario: Summary view (default)

- WHEN user runs `spectr tasks <change-id>`
- THEN the system displays section-by-section progress summary
- AND each section shows completed/total counts
- AND total progress percentage is displayed at the end

#### Scenario: Flatten view

- WHEN user runs `spectr tasks <change-id> --flatten`
- THEN the system displays all tasks from root and child files
- AND tasks are shown with full hierarchical IDs
- AND output format matches flat tasks.jsonc structure

#### Scenario: JSON output

- WHEN user runs `spectr tasks <change-id> --json`
- THEN the system outputs merged task data as JSON
- AND includes all tasks from root and child files
- AND includes aggregated summary counts

#### Scenario: Interactive change selection

- WHEN user runs `spectr tasks` without specifying a change ID
- THEN the system displays a list of active changes with tasks.jsonc files
- AND prompts for selection using existing TUI components

#### Scenario: Missing tasks file

- WHEN user runs `spectr tasks <change-id>` on a change without tasks.jsonc
- THEN the system displays an error indicating no tasks.jsonc found
- AND suggests running `spectr accept` first
- AND the system exits with code 1

### Requirement: Child File Validation

The validation system SHALL validate hierarchical task file references.

#### Scenario: Valid child reference

- WHEN validating a tasks.jsonc with children references
- AND all referenced files exist and have valid format
- THEN validation SHALL pass

#### Scenario: Missing child file

- WHEN validating a tasks.jsonc with children reference to non-existent file
- THEN validation SHALL report an error
- AND error message SHALL include the missing file path

#### Scenario: Invalid child format

- WHEN a referenced child tasks.jsonc has invalid format
- THEN validation SHALL report an error
- AND error message SHALL identify the specific file and issue

#### Scenario: Orphaned child files

- WHEN a delta spec has tasks.jsonc but is not referenced by root
- AND root has auto-discovery includes pattern
- THEN validation SHALL pass (file is discovered)

#### Scenario: Circular references

- WHEN a child tasks.jsonc references another child file
- THEN validation SHALL report an error
- AND error message SHALL indicate circular reference detected

### Requirement: Automatic Task File Splitting

The system SHALL automatically split large tasks.md files during `spectr accept`
based on a 100-line threshold, using section boundaries for clean breaks.

#### Scenario: Detect large tasks.md file

- **WHEN** running `spectr accept <change-id>`
- **THEN** the system SHALL count lines in tasks.md
- **AND** if line count exceeds 100, enable splitting mode

#### Scenario: Split by section boundaries

- **WHEN** splitting is enabled
- **THEN** the system SHALL identify top-level sections (lines starting with `##
  N.`)
- **AND** each section SHALL be evaluated for splitting into its own file
- **AND** the root tasks.jsonc SHALL contain reference tasks for each split
  section

#### Scenario: Split large sections by subsections

- **WHEN** a single section exceeds 100 lines
- **THEN** the system SHALL identify subsections by task ID prefixes (e.g.,
  "1.1", "1.2")
- **AND** group related tasks with common ID prefixes together
- **AND** split subsection groups into separate files if needed

#### Scenario: Preserve small sections in root

- **WHEN** splitting is enabled
- **THEN** sections with fewer than 10 tasks MAY remain in root file
- **AND** the system SHALL prioritize splitting largest sections first

#### Scenario: Sequential file naming

- **WHEN** generating split files
- **THEN** name files sequentially: tasks-1.jsonc, tasks-2.jsonc, etc.
- **AND** the sequence SHALL follow section order from tasks.md

#### Scenario: Regenerate from tasks.md

- **WHEN** running `spectr accept` on a change with existing tasks.jsonc
- **THEN** the system SHALL delete all existing tasks-*.jsonc files
- **AND** regenerate all files from tasks.md
- **AND** preserve status values from existing tasks.jsonc where task IDs match

### Requirement: Hierarchical Task ID Schema

The system SHALL support hierarchical task IDs using dot notation to represent
parent-child relationships.

#### Scenario: Root task IDs

- **WHEN** a task is in the root tasks.jsonc
- **THEN** its ID SHALL be a single integer or decimal (e.g., "1", "2", "5")

#### Scenario: Child task IDs

- **WHEN** a task is in a child file with parent ID "5"
- **THEN** its ID SHALL be prefixed with parent ID (e.g., "5.1", "5.2")

#### Scenario: Nested child task IDs

- **WHEN** a child task has its own children (subsections within a split
  section)
- **THEN** IDs SHALL use additional dot notation (e.g., "5.1.1", "5.1.2")

#### Scenario: ID uniqueness within change

- **WHEN** generating tasks across multiple files
- **THEN** all task IDs SHALL be unique within the change
- **AND** parent-child relationships SHALL be explicit via dot notation

### Requirement: Task Status Aggregation

The system SHALL compute aggregate status for reference tasks based on their
children.

#### Scenario: All children pending

- **WHEN** all child tasks have status "pending"
- **THEN** the parent reference task SHALL have status "pending"

#### Scenario: Any child in progress

- **WHEN** at least one child task has status "in_progress"
- **THEN** the parent reference task SHALL have status "in_progress"

#### Scenario: All children completed

- **WHEN** all child tasks have status "completed"
- **THEN** the parent reference task SHALL have status "completed"

#### Scenario: Mixed pending and completed

- **WHEN** some children are "completed" and others are "pending" (none
  "in_progress")
- **THEN** the parent reference task SHALL have status "in_progress"

### Requirement: Backwards Compatibility

The system SHALL maintain full compatibility with version 1 flat tasks.jsonc
files.

#### Scenario: Read version 1 files

- **WHEN** reading a tasks.jsonc with `version: 1`
- **THEN** the system SHALL parse it as a flat task list
- **AND** no split file lookup SHALL be performed

#### Scenario: Generate version 1 for small files

- **WHEN** tasks.md is 100 lines or fewer
- **THEN** generate version 1 tasks.jsonc
- **AND** use existing flat structure without hierarchical fields

#### Scenario: Upgrade from version 1 to version 2

- **WHEN** running accept on a change with version 1 tasks.jsonc where tasks.md
  now exceeds 100 lines
- **THEN** the system SHALL regenerate as version 2 with split files
- **AND** preserve existing task statuses where IDs match
