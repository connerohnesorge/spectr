# Cli Interface Specification

## Purpose

Defines the CLI framework using Kong for declarative commands, struct tags, subcommands (archive, list, validate, view), and interactive features (TUI, wizard, visual styling).

## Requirements

### Requirement: Archive Command

The system SHALL provide an `archive` command to move completed changes and apply deltas.

#### Scenario: Archive with change ID
- WHEN `spectr archive <change-id>`
- THEN archive change without prompting

#### Scenario: Interactive archive selection
- WHEN `spectr archive` (no ID)
- THEN prompt to select active change

#### Scenario: Non-interactive archiving with yes flag
- WHEN `spectr archive <change-id> --yes`
- THEN archive without confirmation

#### Scenario: Skip spec updates for tooling changes
- WHEN `spectr archive <change-id> --skip-specs`
- THEN archive without updating specs

#### Scenario: Skip validation with confirmation
- WHEN `spectr archive <change-id> --no-validate`
- THEN warn and require confirmation (unless --yes)

### Requirement: Positional Argument Support

The system SHALL support positional arguments via struct fields tagged `arg`.

#### Scenario: Optional positional argument
- WHEN command accepts optional arg
- THEN define with `arg` and `optional` tags (pointer or zero-value)

#### Scenario: Required positional argument
- WHEN command requires arg
- THEN define with `arg` tag; fail if missing

### Requirement: Built-in Help Generation

The system SHALL generate help text from struct tags.

#### Scenario: Root help display
- WHEN `--help` or no args
- THEN list subcommands, descriptions, and args

#### Scenario: Subcommand help display
- WHEN subcommand `--help`
- THEN show description, flags, types, and args

### Requirement: Error Handling and Exit Codes

The system SHALL provide clear errors and exit codes.

#### Scenario: Parse error handling
- WHEN invalid flags/args
- THEN show error, usage, exit non-zero

#### Scenario: Execution error handling
- WHEN Run method returns error
- THEN show error, exit non-zero

### Requirement: Backward-Compatible CLI Interface

The system SHALL maintain syntax and flag compatibility.

#### Scenario: Init command compatibility
- WHEN `spectr init` invoked
- THEN match previous Cobra behavior (flags, aliases, args)

#### Scenario: Help text accessibility
- WHEN `spectr --help`
- THEN document all commands/flags

### Requirement: List Command for Changes

The system SHALL provide a `list` command that enumerates active changes (IDs default).

#### Scenario: List changes with IDs only
- WHEN `spectr list`
- THEN display IDs (sorted, excluding archive)

#### Scenario: List changes with details
- WHEN `spectr list --long`
- THEN format: `{id}: {title} [deltas {count}] [tasks {completed}/{total}]`

#### Scenario: List changes as JSON
- WHEN `spectr list --json`
- THEN output JSON array (id, title, deltaCount, taskStatus)

#### Scenario: No changes found
- WHEN no active changes
- THEN display "No items found"

### Requirement: List Command for Specs

The system SHALL support a `--specs` flag to list specifications.

#### Scenario: List specs with IDs only
- WHEN `spectr list --specs`
- THEN display spec IDs

#### Scenario: List specs with details
- WHEN `spectr list --specs --long`
- THEN format: `{id}: {title} [requirements {count}]`

#### Scenario: List specs as JSON
- WHEN `spectr list --specs --json`
- THEN output JSON array (id, title, requirementCount)

#### Scenario: No specs found
- WHEN no specs exist
- THEN display "No items found"

### Requirement: Change Discovery

The system SHALL scan `spectr/changes/` for subdirectories with `proposal.md`.

#### Scenario: Find active changes
- WHEN scanning
- THEN include `spectr/changes/*/proposal.md`
- AND exclude `archive/` and hidden dirs

### Requirement: Spec Discovery

The system SHALL scan `spectr/specs/` for subdirectories with `spec.md`.

#### Scenario: Find specs
- WHEN scanning
- THEN include `spectr/specs/*/spec.md`
- AND exclude hidden dirs

### Requirement: Title Extraction

The system SHALL extract title from first level-1 heading in markdown (removing "Change:"/"Spec:" prefix).

#### Scenario: Extract title from proposal
- WHEN reading proposal
- THEN extract title from `# Change: ...`

#### Scenario: Extract title from spec
- WHEN reading spec
- THEN extract title from `# ...`

#### Scenario: Fallback to ID when title not found
- WHEN no title found
- THEN use directory name

### Requirement: Task Counting

The system SHALL count tasks from `tasks.jsonc` or `tasks.md` (ignore `tasks.json`).

#### Scenario: Count tasks from JSONC
- WHEN `tasks.jsonc` exists
- THEN count by status (strip comments)

#### Scenario: Count tasks from Markdown
- WHEN `tasks.md` exists (no jsonc)
- THEN count `- [ ]` and `- [x]`

#### Scenario: Ignore legacy tasks.json
- WHEN only `tasks.json` exists
- THEN report 0 tasks

#### Scenario: Handle missing tasks file
- WHEN no tasks file
- THEN report 0 tasks

#### Scenario: JSONC takes precedence over Markdown
- WHEN both exist
- THEN use `tasks.jsonc`

### Requirement: Validate Command

The system SHALL check spec/change correctness.

#### Scenario: Validate command registration
- WHEN initialized
- THEN register `spectr validate`

#### Scenario: Direct item validation invocation
- WHEN `spectr validate <item-name>`
- THEN validate item, print results, exit 0/1

#### Scenario: Bulk validation invocation
- WHEN `spectr validate --all`
- THEN validate all, print summary/issues, exit 1 on failure

#### Scenario: Interactive validation invocation
- WHEN `spectr validate` (no args, TTY)
- THEN prompt for selection

#### Scenario: Default validation behavior (always strict)
- WHEN validating
- THEN treat warnings as errors (exit 1)

#### Scenario: JSON output flag
- WHEN `--json`
- THEN output structured JSON

#### Scenario: Type disambiguation flag
- WHEN `--type change` or `--type spec`
- THEN validate as type (error if missing)

#### Scenario: Changes only flag
- WHEN `--changes`
- THEN validate changes only

#### Scenario: Specs only flag
- WHEN `--specs`
- THEN validate specs only

#### Scenario: Non-interactive flag
- WHEN `--no-interactive`
- THEN no prompt, exit 1 if no item

### Requirement: Validate Command Help Text

The system SHALL document validation purpose and flags.

#### Scenario: Command help display
- WHEN `spectr validate --help`
- THEN show description, flags, examples

### Requirement: Positional Argument Support for Item Name

The system SHALL accept an optional item name.

#### Scenario: Optional item name argument
- WHEN defined
- THEN `arg:"" optional:""`

#### Scenario: Item name provided
- WHEN provided
- THEN validate specific item (auto-detect type)

### Requirement: View Command

The system SHALL display a project dashboard (`spectr view`).

#### Scenario: View command registration
- WHEN initialized
- THEN register `spectr view`

#### Scenario: View command invocation
- WHEN `spectr view`
- THEN show dashboard (summary, active/completed changes, specs)

#### Scenario: View command with JSON output
- WHEN `spectr view --json`
- THEN output dashboard JSON

#### Scenario: JSON structure
- WHEN JSON output
- THEN include `summary`, `activeChanges`, `completedChanges`, `specs`

#### Scenario: JSON arrays sorted consistently
- WHEN JSON output
- THEN sort active (progress, ID), completed (ID), specs (reqs, ID)

#### Scenario: JSON with no items
- WHEN empty
- THEN return empty arrays

### Requirement: Dashboard Summary Metrics

The system SHALL display aggregated stats.

#### Scenario: Display summary with all metrics
- WHEN rendering summary
- THEN show counts: specs, requirements, active/completed changes, tasks

#### Scenario: Calculate total requirements
- WHEN aggregating
- THEN sum requirements from all specs

#### Scenario: Calculate task progress
- WHEN aggregating
- THEN sum tasks/completed from active changes

### Requirement: Active Changes Display

The system SHALL show active changes with progress bars.

#### Scenario: List active changes with progress
- WHEN rendering active changes
- THEN show ID, progress bar, %, indicator (◉)
- AND sort by completion asc, ID

#### Scenario: Render progress bar
- WHEN rendering bar
- THEN 20 chars, full/light blocks, green/gray

#### Scenario: Handle zero tasks
- WHEN 0 tasks
- THEN empty bar, 0%

#### Scenario: No active changes
- WHEN none
- THEN hide section

### Requirement: Completed Changes Display

The system SHALL show completed changes.

#### Scenario: List completed changes
- WHEN rendering completed
- THEN show ID, checkmark (✓), sort ID

#### Scenario: Determine completion status
- WHEN evaluating
- THEN complete if all tasks done or 0 tasks

#### Scenario: No completed changes
- WHEN none
- THEN hide section

### Requirement: Specifications Display

The system SHALL show specs sorted by complexity.

#### Scenario: List specifications with requirement counts
- WHEN rendering specs
- THEN show ID, count, indicator (▪), sort count desc

#### Scenario: Pluralize requirement label
- WHEN displaying
- THEN handle singular/plural

#### Scenario: No specifications found
- WHEN none
- THEN hide section

### Requirement: Dashboard Visual Formatting

The system SHALL use colors and box-drawing chars.

#### Scenario: Render dashboard header
- WHEN rendering
- THEN title, double-line separator

#### Scenario: Render section headers
- WHEN rendering section
- THEN bold cyan, single-line separator

#### Scenario: Render footer
- WHEN rendering footer
- THEN double-line separator, hints

#### Scenario: Color scheme consistency
- WHEN coloring
- THEN cyan (headers), yellow (active), green (done), blue (specs)

### Requirement: Sorting Strategy

The system SHALL sort for relevance.

#### Scenario: Sort active changes by priority
- WHEN sorting active
- THEN % asc, ID

#### Scenario: Sort specs by complexity
- WHEN sorting specs
- THEN count desc, ID

#### Scenario: Sort completed changes alphabetically
- WHEN sorting completed
- THEN ID

### Requirement: View Command Help Text

The system SHALL document the view command.

#### Scenario: View command help display
- WHEN `spectr view --help`
- THEN show description, `--json`

### Requirement: Provider Interface

The system SHALL provide a `Provider` interface for AI tools (instruction files + slash commands).

#### Scenario: Provider interface methods
- WHEN provider created
- THEN implement ID, Name, Priority, ConfigFile, paths, Configure, IsConfigured

#### Scenario: Single provider per tool
- WHEN tool has config & slash
- THEN one provider handles both

#### Scenario: Flexible command paths
- WHEN returning paths
- THEN relative path or empty

#### Scenario: HasSlashCommands detection
- WHEN checking
- THEN true if any path non-empty

### Requirement: Command Format Support

The system SHALL support Markdown/TOML formats.

#### Scenario: Markdown command format
- WHEN `FormatMarkdown`
- THEN create `.md` with frontmatter

#### Scenario: TOML command format
- WHEN `FormatTOML`
- THEN create `.toml` with description/prompt

### Requirement: Version Command Structure

The system SHALL display version/commit/date.

#### Scenario: Version command registration
- WHEN initialized
- THEN register `spectr version`

#### Scenario: Version command invocation
- WHEN `spectr version`
- THEN show version info

#### Scenario: Version command with short flag
- WHEN `spectr version --short`
- THEN show version only

#### Scenario: Version command with JSON flag
- WHEN `spectr version --json`
- THEN output JSON

### Requirement: Version Variable Injection

The system SHALL inject version variables via ldflags.

#### Scenario: Goreleaser/Nix/Dev injection
- WHEN building
- THEN set version/commit/date (default dev/unknown)

### Requirement: Completion Command Structure

The system SHALL generate shell completions.

#### Scenario: Completion command registration
- WHEN initialized
- THEN register `spectr completion`

#### Scenario: Bash/Zsh/Fish completion output
- WHEN `spectr completion <shell>`
- THEN output script

### Requirement: Custom Predictors for Dynamic Arguments

The system SHALL suggest IDs.

#### Scenario: Change/Spec ID completion
- WHEN tabbing ID arg
- THEN suggest active changes/specs

#### Scenario: Item type completion
- WHEN tabbing `--type`
- THEN suggest `change`, `spec`

### Requirement: Accept Command Structure

The system SHALL convert `tasks.md` to `tasks.jsonc`.

#### Scenario: Accept command registration
- WHEN initialized
- THEN register `spectr accept`

#### Scenario: Accept with change ID
- WHEN `spectr accept <id>`
- THEN validate, parse md, write jsonc (preserve md)

#### Scenario: Accept with validation
- WHEN running
- THEN abort if invalid

#### Scenario: Accept dry-run mode
- WHEN `--dry-run`
- THEN preview only

#### Scenario: Accept already accepted change
- WHEN jsonc exists
- THEN regenerate from md

#### Scenario: Accept change without tasks.md
- WHEN no md
- THEN error

#### Scenario: Both tasks.md and tasks.jsonc exist
- WHEN both exist
- THEN prefer jsonc runtime, md reference

### Requirement: Tasks JSON Schema

The system SHALL generate a versioned `tasks.jsonc`.

#### Scenario: JSONC file structure
- WHEN generating
- THEN header comments, version, tasks array

#### Scenario: Header comment content
- WHEN generating
- THEN document status values/transitions, machine-generated warning

#### Scenario: Task object structure
- WHEN serializing
- THEN id, section, description, status

#### Scenario: Status value mapping from Markdown
- WHEN converting
- THEN map `[ ]` -> pending, `[x]` -> completed

### Requirement: Accept Command Flags

The system SHALL provide flags to control behavior.

#### Scenario: Dry-run flag
- WHEN `--dry-run`
- THEN preview JSON

#### Scenario: Interactive change selection
- WHEN no ID
- THEN prompt list

### Requirement: List Command Alias

The system SHALL alias `ls` to `list`.

#### Scenario: User runs spectr ls shorthand
- WHEN `spectr ls`
- THEN same as `list` (flags work)

### Requirement: Item Name Path Normalization

The system SHALL normalize paths to IDs.

#### Scenario: Path/ID normalization
- WHEN path provided (e.g. `spectr/changes/id`)
- THEN extract ID, infer type

### Requirement: Interactive List Mode

The system SHALL provide a unified TUI for changes/specs.

#### Scenario: Default behavior unchanged
- WHEN `spectr list -I`
- THEN show changes only

#### Scenario: Unified mode opt-in
- WHEN `spectr list -I --all`
- THEN show changes and specs (Type column)

#### Scenario: Type-specific actions
- WHEN unified
- THEN only specs editable

#### Scenario: Help text formatting
- WHEN interactive
- THEN minimal footer, `?` for full help

### Requirement: Clipboard Copy on Selection

The system SHALL copy ID on Enter.

#### Scenario: Copy ID
- WHEN Enter pressed
- THEN copy ID to clipboard (exit)

### Requirement: Interactive Mode Exit Controls

The system SHALL provide standard quit controls.

#### Scenario: Quit
- WHEN `q` or `Ctrl+C`
- THEN exit

### Requirement: Table Visual Styling

The system SHALL provide consistent TUI styling.

#### Scenario: Visual hierarchy
- WHEN displaying
- THEN styled headers, selection, borders (`tui.ApplyTableStyles`)

### Requirement: Cross-Platform Clipboard Support

The system SHALL support Linux/macOS/Windows/SSH.

#### Scenario: Clipboard support
- WHEN copying
- THEN use native API or OSC 52 fallback

### Requirement: Initialization Next Steps Message

The system SHALL guide the user after init.

#### Scenario: Next steps display
- WHEN init succeeds
- THEN show 3 steps (project.md, proposal, AGENTS.md)

#### Scenario: Init does not create README
- WHEN init
- THEN do not create README.md

### Requirement: Flat Tool List in Initialization Wizard

The system SHALL provide a unified tool list.

#### Scenario: Display only config-based tools
- WHEN selecting tools
- THEN show config tools only (auto-install slash commands)

#### Scenario: Navigation/Selection
- WHEN interacting
- THEN navigate flat list, space toggle, a/n bulk

### Requirement: Interactive Archive Mode

The system SHALL provide a table interface for archive.

#### Scenario: Archive no args / -I
- WHEN `spectr archive [-I]`
- THEN show interactive table (same columns as list)

#### Scenario: Selection behavior
- WHEN selecting
- THEN capture ID (no copy), proceed to archive

### Requirement: Archive Interactive Table Display

The system SHALL match list columns in archive mode.

#### Scenario: Table columns
- WHEN displaying
- THEN ID, Title, Deltas, Tasks (consistent style)

### Requirement: Archive Selection Without Clipboard

The system SHALL capture selection internally only.

#### Scenario: Enter key
- WHEN Enter
- THEN proceed with ID (no clipboard)

### Requirement: Validation Output Format

The system SHALL provide consistent issue reporting.

#### Scenario: Single/Bulk output
- WHEN validating
- THEN show issues (Level, Path, Message), summary

#### Scenario: JSON output
- WHEN `--json`
- THEN structured issue data

### Requirement: Editor Hotkey in Interactive Specs List

The system SHALL allow editing specs with 'e'.

#### Scenario: Edit spec
- WHEN `e` pressed (specs mode)
- THEN open `$EDITOR` (wait), return to TUI

### Requirement: Editor Hotkey Scope

The system SHALL limit editor hotkey to specs only.

#### Scenario: No edit for changes
- WHEN `e` pressed (changes mode)
- THEN ignore

### Requirement: Project Path Display in Interactive Mode

The system SHALL show context by displaying the project path.

#### Scenario: Path display
- WHEN interactive
- THEN show project root path

### Requirement: Unified Item List Display

The system SHALL provide a mixed table for changes/specs.

#### Scenario: Unified display
- WHEN `--all -I`
- THEN show Type, ID, Title, Details

### Requirement: Type-Aware Item Selection

The system SHALL handle types correctly.

#### Scenario: Selection
- WHEN Enter
- THEN copy ID (both types)

#### Scenario: Edit restriction
- WHEN `e`
- THEN specs only

### Requirement: Backward-Compatible Separate Modes

The system SHALL preserve existing modes.

#### Scenario: Separate modes
- WHEN `-I` (no `--all`) -> Changes
- WHEN `--specs -I` -> Specs

### Requirement: Enhanced List Command Flags

The system SHALL validate flags.

#### Scenario: Validation
- WHEN conflicting flags (e.g. `-I` + `--json`)
- THEN error

### Requirement: Automatic Slash Command Installation

The system SHALL install slash commands with config.

#### Scenario: Auto-install
- WHEN config tool selected
- THEN install slash commands (CLAUDE.md + .claude/...)

### Requirement: Archive Hotkey in Interactive Changes Mode

The system SHALL allow archiving with 'a'.

#### Scenario: Archive action
- WHEN `a` pressed (changes mode)
- THEN exit and archive selected change

### Requirement: Shared TUI Component Library

The system SHALL use `internal/tui`.

#### Scenario: Components
- WHEN building UI
- THEN use `TablePicker`, `MenuPicker`, `TruncateString`, `CopyToClipboard`

### Requirement: Search Hotkey in Interactive Lists

The system SHALL filter lists with '/'.

#### Scenario: Search mode
- WHEN `/` pressed
- THEN input field, filter rows by ID/Title

### Requirement: Help Toggle Hotkey

The system SHALL toggle help with '?'.

#### Scenario: Help toggle
- WHEN `?` pressed
- THEN toggle full/minimal help

### Requirement: Partial Change ID Resolution for Archive Command

The system SHALL resolve prefixes/substrings.

#### Scenario: Resolution
- WHEN partial ID provided
- THEN resolve if unique (prefix > substring), else error

### Requirement: Configured Provider Detection in Init Wizard

The system SHALL detect existing config.

#### Scenario: Detection
- WHEN initializing
- THEN mark/select existing providers

### Requirement: Instruction File Pointer Template

The system SHALL use pointers in root files.

#### Scenario: Pointer content
- WHEN creating `CLAUDE.md` etc
- THEN point to `spectr/AGENTS.md`

### Requirement: PR Archive Subcommand Alias

The system SHALL alias `a` to `archive`.

#### Scenario: Alias
- WHEN `spectr pr a`
- THEN same as `spectr pr archive`

### Requirement: PR Branch Naming Convention

The system SHALL use consistent branch naming.

#### Scenario: Naming
- WHEN archive -> `spectr/archive/<id>`
- WHEN proposal -> `spectr/proposal/<id>`

### Requirement: PR Command Structure

The system SHALL support `spectr pr <subcommand>`.

#### Scenario: Subcommands
- WHEN `archive`/`proposal`
- THEN execute workflow

### Requirement: PR Archive Subcommand

The system SHALL create an archive PR (isolated worktree).

#### Scenario: Workflow
- WHEN `spectr pr archive <id>`
- THEN worktree, archive --yes, commit, push, create PR

### Requirement: PR Proposal Subcommand

The system SHALL create a proposal PR (isolated worktree).

#### Scenario: Workflow
- WHEN `spectr pr proposal <id>`
- THEN worktree, copy change, commit, push, create PR (no archive)

### Requirement: PR Common Flags

The system SHALL provide shared flags.

#### Scenario: Flags
- WHEN `--base`, `--draft`, `--force`, `--dry-run`
- THEN apply behavior

### Requirement: Git Platform Detection

The system SHALL detect GitHub/GitLab/Gitea/Bitbucket.

#### Scenario: Detection
- WHEN remote URL
- THEN detect platform, use CLI (gh/glab/tea)

### Requirement: Platform CLI Availability

The system SHALL check if CLI is installed.

#### Scenario: Check
- WHEN running
- THEN error if CLI missing/unauthenticated

### Requirement: Git Worktree Isolation

The system SHALL isolate operations.

#### Scenario: Isolation
- WHEN running
- THEN use temp worktree, clean up after

### Requirement: PR Commit Message Format

The system SHALL use conventional commits.

#### Scenario: Format
- WHEN committing
- THEN `spectr(archive/proposal): <id>`

### Requirement: PR Body Content

The system SHALL provide a useful description.

#### Scenario: Content
- WHEN creating PR
- THEN summary, checklist, location

### Requirement: PR Branch Naming

The system SHALL use branch naming pattern `spectr/<id>`.

#### Scenario: Naming
- WHEN creating branch
- THEN `spectr/<id>`

### Requirement: PR Error Handling

The system SHALL provide clear errors.

#### Scenario: Errors
- WHEN no git/remote/change/push fail
- THEN display specific error

### Requirement: Partial Change ID Resolution for PR Commands

The system SHALL resolve IDs.

#### Scenario: Resolution
- WHEN partial ID
- THEN resolve unique match

### Requirement: PR Proposal Interactive Selection Filters Unmerged Changes

The system SHALL filter already merged changes.

#### Scenario: Filtering
- WHEN selecting proposal
- THEN hide changes already on main

### Requirement: Template Path Variables

The system SHALL support dynamic paths.

#### Scenario: Variables
- WHEN templating
- THEN use `{{ .BaseDir }}` etc

### Requirement: Copy Populate Context Prompt in Init Next Steps

The system SHALL copy prompt with 'c'.

#### Scenario: Copy prompt
- WHEN 'c' pressed (success screen)
- THEN copy prompt to clipboard

### Requirement: PR Hotkey in Interactive Changes List Mode

The system SHALL enable PR workflow with 'P'.

#### Scenario: PR action
- WHEN `Shift+P` pressed
- THEN enter PR workflow

### Requirement: VHS Demo for PR Hotkey

The system SHALL include a demo asset.

#### Scenario: Demo
- WHEN viewing assets
- THEN `pr-hotkey.tape` exists

### Requirement: PR Proposal Local Change Cleanup Confirmation

The system SHALL prompt to remove local changes.

#### Scenario: Prompt
- WHEN proposal PR success
- THEN prompt remove local (default No)

### Requirement: CI Workflow Setup Option in Init Wizard Review Step

The system SHALL setup GitHub Actions.

#### Scenario: CI option
- WHEN init review
- THEN checkbox for `.github/workflows/spectr-ci.yml`

### Requirement: PR Remove Subcommand

The system SHALL remove change PR.

#### Scenario: Remove workflow
- WHEN `spectr pr rm <id>`
- THEN worktree, remove dir, commit, PR, clean local

### Requirement: Remove PR Branch Naming

The system SHALL use branch naming pattern `spectr/remove/<id>`.

#### Scenario: Naming
- WHEN remove PR
- THEN `spectr/remove/<id>`

### Requirement: Remove PR Commit Message Format

The system SHALL use commit message format `spectr(remove): <id>`.

#### Scenario: Format
- WHEN committing
- THEN structured message

### Requirement: Remove PR Body Content

The system SHALL explain removal in PR body.

#### Scenario: Content
- WHEN PR body
- THEN summary, removed path

### Requirement: Responsive Table Column Layout

The system SHALL adapt columns to width.

#### Scenario: Responsive columns
- WHEN displaying
- THEN hide/narrow columns based on width

### Requirement: Dynamic Terminal Resize Handling

The system SHALL handle terminal resize.

#### Scenario: Resize
- WHEN resized
- THEN recalculate layout

### Requirement: Column Priority System

The system SHALL prioritize columns.

#### Scenario: Priority
- WHEN calculating
- THEN ID > Title > Deltas/Reqs > Tasks

### Requirement: Provider Search in Init Wizard

The system SHALL search tools with '/'.

#### Scenario: Search
- WHEN `/` in wizard
- THEN filter tools

### Requirement: Stdout Output Mode for Interactive List

The system SHALL output ID to stdout.

#### Scenario: Stdout mode
- WHEN `-I --stdout`
- THEN print ID (no clipboard)

### Requirement: JSONC Comment Parsing

The system SHALL strip comments.

#### Scenario: Parsing
- WHEN reading JSONC
- THEN strip `//` and `/* */`

### Requirement: TTY Error Hint

The system SHALL provide hints for non-TTY.

#### Scenario: Hint
- WHEN TTY error
- THEN suggest `--non-interactive`

### Requirement: File Coexistence Documentation

The system SHALL document tasks.md/jsonc coexistence.

#### Scenario: Docs
- WHEN help/success
- THEN mention coexistence

### Requirement: Slash Command Template Updates

The system SHALL provide instructions for tasks.

#### Scenario: Templates
- WHEN proposal/apply
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
