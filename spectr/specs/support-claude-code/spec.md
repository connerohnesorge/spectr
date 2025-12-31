# Support Claude Code Specification

## Purpose

Documents the Claude Code provider integration for Spectr, enabling Spectr to
work seamlessly with Claude Code through instruction file management and slash
command generation.

## Requirements

### Requirement: Claude Code Provider Configuration

The provider SHALL be configured with these settings:

- ID: `claude-code`
- Name: `Claude Code`
- Priority: 1 (highest)
- Config File: `CLAUDE.md`
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Claude Code provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `claude-code`, Name `Claude Code`,
  Priority 1
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.claude/commands/spectr/`
- **AND** it SHALL return a `ConfigFileInitializer` for `CLAUDE.md` with
  TemplateRef from TemplateManager
- **AND** it SHALL return a `SlashCommandsInitializer` for slash commands in
  Markdown format

#### Scenario: Provider metadata

- **WHEN** the provider is registered
- **THEN** the provider name SHALL be "Claude Code"
- **AND** the provider priority SHALL be 1 (highest priority)
- **AND** the provider ID SHALL be "claude-code"

### Requirement: Claude Code Instruction File

The provider SHALL create and maintain a `CLAUDE.md` instruction file in the
project root.

#### Scenario: Instruction file creation

- **WHEN** `spectr init` runs with Claude Code provider selected
- **THEN** the ConfigFileInitializer creates `CLAUDE.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:start -->` and `<!--
  spectr:end -->` markers

#### Scenario: Instruction file updates

- **WHEN** `spectr init` runs in a project with Claude Code provider
- **THEN** the ConfigFileInitializer updates content between markers in
  `CLAUDE.md`
- **AND** preserves any user content outside the markers

### Requirement: Claude Code Slash Commands

The provider SHALL create slash commands in `.claude/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.claude/commands/spectr/`
  directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths

- **WHEN** the SlashCommandsInitializer executes
- **THEN** it SHALL create `.claude/commands/spectr/proposal.md`
- **AND** it SHALL create `.claude/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** each file SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field

### Requirement: Standard Frontmatter

The provider SHALL use standard frontmatter templates for each command type.

#### Scenario: Proposal command frontmatter

- **WHEN** generating the proposal command file
- **THEN** the frontmatter description is "Scaffold a new Spectr change and
  validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** generating the apply command file
- **THEN** the frontmatter description is "Implement an approved Spectr change
  and keep tasks in sync."

#### Scenario: Sync command frontmatter

- **WHEN** generating the sync command file
- **THEN** the frontmatter description is "Detect spec drift from code and
  update specs interactively."

### Requirement: Command Path Construction

The provider SHALL construct command paths using the standard pattern.

#### Scenario: Path construction for Claude Code

- **WHEN** determining command file paths
- **THEN** it uses base directory `.claude/commands`
- **AND** appends `/spectr/` subdirectory
- **AND** appends command name with `.md` extension
- **AND** results in paths like `.claude/commands/spectr/proposal.md`

### Requirement: Claude Code Skills Directory

The provider SHALL create a `.claude/skills/` directory for agent skills.

#### Scenario: Skills directory creation

- **WHEN** the provider returns initializers
- **THEN** it SHALL include a `DirectoryInitializer` for `.claude/skills/`
- **AND** the directory SHALL be created before skills are installed

### Requirement: spectr-accept-wo-spectr-bin Skill

The provider SHALL install the `spectr-accept-wo-spectr-bin` skill for
accepting changes without the spectr binary.

#### Scenario: Skill installation path

- **WHEN** the provider returns initializers
- **THEN** it SHALL include an `AgentSkillsInitializer` for
  `spectr-accept-wo-spectr-bin`
- **AND** the skill SHALL be installed at
  `.claude/skills/spectr-accept-wo-spectr-bin/`

#### Scenario: Skill structure

- **WHEN** the skill is installed
- **THEN** it SHALL create `SKILL.md` with valid AgentSkills frontmatter
- **AND** it SHALL create `scripts/accept.sh` for task conversion
- **AND** `scripts/accept.sh` SHALL be executable

#### Scenario: SKILL.md content

- **WHEN** the `SKILL.md` file is created
- **THEN** the frontmatter `name` SHALL be `spectr-accept-wo-spectr-bin`
- **AND** the frontmatter `description` SHALL describe the skill's purpose
- **AND** the `compatibility` field SHALL note `jq` as a requirement
- **AND** the body SHALL contain usage instructions for the accept script

#### Scenario: accept.sh functionality

- **WHEN** `scripts/accept.sh` is executed with a change-id argument
- **THEN** it SHALL read `spectr/changes/<id>/tasks.md`
- **AND** it SHALL parse markdown task lists into structured format
- **AND** it SHALL output `spectr/changes/<id>/tasks.jsonc`
- **AND** the output SHALL be valid JSONC format

### Requirement: spectr-validate-wo-spectr-bin Skill

The provider SHALL install the `spectr-validate-wo-spectr-bin` skill for
validating specifications and change proposals without the spectr binary.

#### Scenario: Validate skill installation path

- **WHEN** the provider returns initializers
- **THEN** it SHALL include an `AgentSkillsInitializer` for
  `spectr-validate-wo-spectr-bin`
- **AND** the skill SHALL be installed at
  `.claude/skills/spectr-validate-wo-spectr-bin/`

#### Scenario: Validate skill structure

- **WHEN** the skill is installed
- **THEN** it SHALL create `SKILL.md` with valid AgentSkills frontmatter
- **AND** it SHALL create `scripts/validate.sh` for specification validation
- **AND** `scripts/validate.sh` SHALL be executable (0755 permissions)

#### Scenario: Validate SKILL.md content

- **WHEN** the `SKILL.md` file is created
- **THEN** the frontmatter `name` SHALL be `spectr-validate-wo-spectr-bin`
- **AND** the frontmatter `description` SHALL describe the skill's purpose
- **AND** the `compatibility` section SHALL list bash 4.0+, grep, sed, find as
  requirements
- **AND** the `compatibility` section SHALL list jq as optional for JSON output
- **AND** the body SHALL contain usage instructions with examples
- **AND** the body SHALL document validation rules matching `spectr validate`
  behavior
- **AND** the body SHALL document exit codes (0=success, 1=failure, 2=usage)
- **AND** the body SHALL document limitations (no pre-merge validation, no
  cross-capability duplicates, sequential processing)

#### Scenario: validate.sh basic structure

- **WHEN** `scripts/validate.sh` is executed
- **THEN** it SHALL use `#!/usr/bin/env bash` shebang
- **AND** it SHALL set `set -euo pipefail` for error handling
- **AND** it SHALL support `SPECTR_DIR` environment variable (default: `spectr`)
- **AND** it SHALL define regex patterns matching `internal/markdown/` matchers

#### Scenario: validate.sh spec file validation

- **WHEN** `scripts/validate.sh` validates a spec file
- **THEN** it SHALL check for `## Requirements` section (ERROR if missing)
- **AND** it SHALL check requirements contain SHALL or MUST (ERROR if missing)
- **AND** it SHALL check requirements have `#### Scenario:` blocks (ERROR if
  missing)
- **AND** it SHALL detect malformed scenario formatting (ERROR for wrong header
  levels, bullets, bold)
- **AND** it SHALL report errors with file path and line number

#### Scenario: validate.sh change delta validation

- **WHEN** `scripts/validate.sh` validates a change directory
- **THEN** it SHALL verify `specs/` directory exists
- **AND** it SHALL find all `spec.md` files under `specs/`
- **AND** it SHALL check for at least one delta section (ADDED, MODIFIED,
  REMOVED, RENAMED)
- **AND** it SHALL check delta sections are not empty (ERROR if no
  requirements)
- **AND** it SHALL validate ADDED requirements have scenarios (ERROR if missing)
- **AND** it SHALL validate ADDED requirements have SHALL/MUST (ERROR if
  missing)
- **AND** it SHALL validate MODIFIED requirements have scenarios (ERROR if
  missing)
- **AND** it SHALL validate MODIFIED requirements have SHALL/MUST (ERROR if
  missing)
- **AND** it SHALL skip scenario/SHALL validation for REMOVED requirements
- **AND** it SHALL skip normal validation for RENAMED section (different format)
- **AND** it SHALL report ERROR if total delta count is zero

#### Scenario: validate.sh tasks file validation

- **WHEN** `scripts/validate.sh` validates a change with tasks.md
- **THEN** it SHALL check if tasks.md exists (skip if not present)
- **AND** it SHALL count task items using pattern `- [ ]` or `- [x]`
- **AND** it SHALL support task items with optional IDs (`- [ ] 1.1 Task`)
- **AND** it SHALL support uppercase X (`- [X] Task`)
- **AND** it SHALL report ERROR if tasks.md exists but has zero tasks
- **AND** it SHALL provide helpful error message with expected format

#### Scenario: validate.sh single spec mode

- **WHEN** `scripts/validate.sh` is executed with `--spec <spec-id>`
- **THEN** it SHALL validate `spectr/specs/<spec-id>/spec.md`
- **AND** it SHALL verify the spec directory exists before validation
- **AND** it SHALL report ERROR if spec directory or file not found
- **AND** it SHALL output validation results for that single spec

#### Scenario: validate.sh single change mode

- **WHEN** `scripts/validate.sh` is executed with `--change <change-id>`
- **THEN** it SHALL validate `spectr/changes/<change-id>/`
- **AND** it SHALL verify the change directory exists before validation
- **AND** it SHALL validate all delta specs in the change
- **AND** it SHALL validate tasks.md if present
- **AND** it SHALL output validation results for that single change

#### Scenario: validate.sh bulk validation mode

- **WHEN** `scripts/validate.sh` is executed with `--all`
- **THEN** it SHALL discover all specs in `spectr/specs/`
- **AND** it SHALL discover all changes in `spectr/changes/` (excluding archive)
- **AND** it SHALL validate each discovered spec
- **AND** it SHALL validate each discovered change
- **AND** it SHALL aggregate results across all items
- **AND** it SHALL output summary with passed/failed counts

#### Scenario: validate.sh spec discovery

- **WHEN** the script discovers specs
- **THEN** it SHALL find all directories under `spectr/specs/` containing
  `spec.md`
- **AND** it SHALL return empty list if `spectr/specs/` doesn't exist
- **AND** it SHALL sort spec IDs alphabetically
- **AND** it SHALL match discovery behavior of
  `internal/discovery/GetSpecIDs()`

#### Scenario: validate.sh change discovery

- **WHEN** the script discovers changes
- **THEN** it SHALL find all directories under `spectr/changes/` except
  `archive`
- **AND** it SHALL return empty list if `spectr/changes/` doesn't exist
- **AND** it SHALL sort change IDs alphabetically
- **AND** it SHALL match discovery behavior of
  `internal/discovery/GetChangeIDs()`

#### Scenario: validate.sh human-readable output

- **WHEN** validation completes without `--json` flag
- **THEN** it SHALL output results in human-readable format
- **AND** it SHALL group issues by file path
- **AND** it SHALL print file path as header followed by indented issues
- **AND** it SHALL format issues as `[LEVEL] line N: message`
- **AND** it SHALL color-code error levels if stdout is TTY (ERROR in red)
- **AND** it SHALL not use ANSI codes if stdout is not TTY
- **AND** it SHALL print summary line `X passed, Y failed (E errors), Z total`
- **AND** it SHALL add blank line separators between failed items

#### Scenario: validate.sh JSON output

- **WHEN** validation completes with `--json` flag
- **THEN** it SHALL output valid JSON to stdout
- **AND** the JSON SHALL have structure matching `spectr validate --json`
- **AND** the JSON SHALL include `version` field (value: 1)
- **AND** the JSON SHALL include `items` array with per-item results
- **AND** each item SHALL have `name`, `type`, `valid`, and `issues` fields
- **AND** the JSON SHALL include `summary` object with `total`, `passed`,
  `failed`, `errors`, `warnings` fields
- **AND** it SHALL use jq for JSON generation if available
- **AND** it SHALL warn and fallback to human output if jq unavailable

#### Scenario: validate.sh exit codes

- **WHEN** validation completes
- **THEN** exit code 0 SHALL indicate all validations passed (no errors)
- **AND** exit code 1 SHALL indicate one or more validations failed
- **AND** exit code 2 SHALL indicate usage error (invalid arguments)
- **AND** the exit code SHALL be usable in CI pipelines for pass/fail
  determination

#### Scenario: validate.sh argument parsing

- **WHEN** the script is invoked with arguments
- **THEN** it SHALL support `--spec <spec-id>` flag
- **AND** it SHALL support `--change <change-id>` flag
- **AND** it SHALL support `--all` flag
- **AND** it SHALL support `--json` flag (combinable with validation mode)
- **AND** it SHALL support `-h` and `--help` flags
- **AND** it SHALL require exactly one mode flag (--spec, --change, or --all)
- **AND** it SHALL error with code 2 if mode flag missing
- **AND** it SHALL error with code 2 if unknown flag provided
- **AND** it SHALL print usage message on error

#### Scenario: validate.sh SPECTR_DIR environment variable

- **WHEN** `SPECTR_DIR` environment variable is set
- **THEN** the script SHALL use that value as spectr directory location
- **AND** it SHALL default to `spectr` if variable not set
- **AND** it SHALL derive specs directory as `$SPECTR_DIR/specs`
- **AND** it SHALL derive changes directory as `$SPECTR_DIR/changes`

#### Scenario: validate.sh regex patterns

- **WHEN** the script parses markdown files
- **THEN** it SHALL use pattern `^##[[:space:]]+Requirements[[:space:]]*$` for
  Requirements section
- **AND** it SHALL use pattern
  `^###[[:space:]]+Requirement:[[:space:]]+(.+)$` for requirement headers
- **AND** it SHALL use pattern `^####[[:space:]]+Scenario:` for scenario headers
- **AND** it SHALL use pattern
  `^##[[:space:]]+(ADDED|MODIFIED|REMOVED|RENAMED)[[:space:]]+Requirements` for
  delta sections
- **AND** it SHALL use pattern `^[[:space:]]*-[[:space:]]+\[([ xX])\]` for task
  items
- **AND** patterns SHALL match behavior of `internal/markdown/` matchers

#### Scenario: validate.sh malformed scenario detection

- **WHEN** the script detects malformed scenarios
- **THEN** it SHALL detect `### Scenario:` (3 hashtags) as malformed
- **AND** it SHALL detect `##### Scenario:` (5 hashtags) as malformed
- **AND** it SHALL detect `###### Scenario:` (6 hashtags) as malformed
- **AND** it SHALL detect `**Scenario:` (bold) as malformed
- **AND** it SHALL detect `- **Scenario:` (bullet + bold) as malformed
- **AND** it SHALL report ERROR with line number for malformed scenarios

#### Scenario: validate.sh line-by-line parsing

- **WHEN** the script validates files
- **THEN** it SHALL process files line-by-line using bash `read` loop
- **AND** it SHALL track current line number for error reporting
- **AND** it SHALL maintain state for current section and requirement
- **AND** it SHALL flush requirement validation on section boundaries
- **AND** it SHALL flush requirement validation at end of file

#### Scenario: validate.sh requirement state tracking

- **WHEN** the script processes requirements
- **THEN** it SHALL track whether requirement has scenario (`has_scenario`
  flag)
- **AND** it SHALL track whether requirement has SHALL/MUST (`has_shall_must`
  flag)
- **AND** it SHALL reset state flags when entering new requirement
- **AND** it SHALL validate accumulated state when flushing requirement

#### Scenario: validate.sh issue collection

- **WHEN** the script detects validation issues
- **THEN** it SHALL collect issues in global `ISSUES` array
- **AND** each issue SHALL have level, path, line number, and message
- **AND** it SHALL count errors in `ISSUE_COUNTS` associative array
- **AND** it SHALL use `add_issue()` function for standardized issue creation

#### Scenario: validate.sh color setup

- **WHEN** the script initializes
- **THEN** it SHALL detect if stdout is TTY using `[[ -t 1 ]]`
- **AND** it SHALL define ANSI color codes if TTY (RED, YELLOW, GREEN, NC)
- **AND** it SHALL define empty color codes if not TTY
- **AND** color codes SHALL be used in human-readable output formatting

#### Scenario: validate.sh relative path formatting

- **WHEN** the script formats file paths for output
- **THEN** it SHALL remove `spectr/` prefix from paths
- **AND** output SHALL show `changes/foo/specs/bar/spec.md` not
  `spectr/changes/foo/specs/bar/spec.md`
- **AND** paths SHALL match format from `internal/validation/formatters.go`

#### Scenario: validate.sh usage message

- **WHEN** the script prints usage with `-h`, `--help`, or on error
- **THEN** it SHALL display all available flags (--spec, --change, --all,
  --json)
- **AND** it SHALL provide usage examples for each mode
- **AND** it SHALL document exit codes
- **AND** it SHALL document SPECTR_DIR environment variable
