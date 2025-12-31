# Support Codex Specification

## Purpose

Specifies how Spectr integrates with Codex CLI using global prompt files in the
user's home directory.

## Requirements

### Requirement: Codex Provider Configuration

The provider SHALL be configured with these settings:

- ID: `codex`
- Name: `Codex CLI`
- Priority: 9
- Config File: `AGENTS.md`
- Command Format: Markdown

#### Scenario: Provider registration

- **WHEN** the Codex provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `codex`, Name `Codex CLI`, Priority 9
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers with home paths

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `HomeDirectoryInitializer` for `.codex/prompts/`
  (relative to home directory)
- **AND** it SHALL return a `ConfigFileInitializer` for `AGENTS.md` with
  TemplateRef from TemplateManager using `tm.Agents()`
- **AND** it SHALL return a `HomePrefixedSlashCommandsInitializer` with prefix
  `spectr-` for home slash commands in `.codex/prompts/`

#### Scenario: Provider metadata

- **WHEN** provider is registered
- **THEN** the provider name is "Codex CLI"
- **AND** it appears after Cursor (priority 8) and before Aider (priority 10)

#### Scenario: Instruction file

- **WHEN** the provider returns initializers
- **THEN** it includes a ConfigFileInitializer for "AGENTS.md"

#### Scenario: Create new instruction file

- **WHEN** `AGENTS.md` does not exist
- **THEN** the ConfigFileInitializer SHALL create it with instruction content
  between markers
- **AND** the markers SHALL be `<!-- spectr:start -->` and `<!-- spectr:end -->`
  (lowercase)

#### Scenario: Update existing instruction file

- **WHEN** `AGENTS.md` exists with spectr markers
- **THEN** the ConfigFileInitializer SHALL replace content between markers
- **AND** it SHALL preserve content outside markers
- **AND** the marker search SHALL be case-insensitive (matches both uppercase
  and lowercase)
- **AND** when writing, the system SHALL always use lowercase markers

### Requirement: Codex Global Slash Commands

The provider SHALL create slash commands in the home `~/.codex/prompts/`
directory.

#### Scenario: Home command directory structure

- **WHEN** the provider returns initializers
- **THEN** `HomeDirectoryInitializer` SHALL create `~/.codex/prompts/` directory
- **AND** the directory is created in user's home directory via homeFs

#### Scenario: Command paths with prefix

- **WHEN** the `HomePrefixedSlashCommandsInitializer` executes
- **THEN** it SHALL create `.codex/prompts/spectr-proposal.md` in the home
  filesystem
- **AND** it SHALL create `.codex/prompts/spectr-apply.md` in the home
  filesystem

#### Scenario: Home path handling

- **WHEN** Home* initializers execute
- **THEN** the executor provides both projectFs and homeFs filesystems
- **AND** the Home* initializers use homeFs automatically (no configuration flag
  needed)
- **AND** paths work correctly regardless of current project directory

### Requirement: Codex Command Format

The provider SHALL use Markdown format with YAML frontmatter for slash commands.

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** each file SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field

### Requirement: Global Path Support in Provider Framework

The provider framework SHALL support global paths (starting with `~/` or `/`) in
addition to project-relative paths.

#### Scenario: Global path detection

- **WHEN** a command path starts with `~/` or `/`
- **THEN** the system treats it as a global path
- **AND** does not prepend the project directory

#### Scenario: Home directory expansion

- **WHEN** a path starts with `~/`
- **THEN** the system expands `~` to the user's home directory
- **AND** uses `os.UserHomeDir()` for cross-platform compatibility

#### Scenario: IsConfigured with global paths

- **WHEN** checking if a provider with global paths is configured
- **THEN** the system checks the expanded absolute path
- **AND** does not look in the project directory

### Requirement: Codex Command Invocation

Users SHALL invoke Spectr commands in Codex using the `/spectr-<command>`
pattern.

#### Scenario: Invoking proposal command

- **WHEN** user types `/spectr-proposal` in Codex
- **THEN** Codex loads and executes the proposal prompt

#### Scenario: Invoking apply command

- **WHEN** user types `/spectr-apply` in Codex
- **THEN** Codex loads and executes the apply prompt

### Requirement: Codex Skills Directory

The provider SHALL create a `.codex/skills/` directory for agent skills.

#### Scenario: Skills directory creation

- **WHEN** the provider returns initializers
- **THEN** it SHALL include a `DirectoryInitializer` for `.codex/skills/`
- **AND** the directory SHALL be created in the project filesystem (not home)
- **AND** the directory SHALL be created before skills are installed

#### Scenario: Skills directory location

- **WHEN** the skills directory is created
- **THEN** it SHALL be located at `.codex/skills/` relative to project root
- **AND** it SHALL NOT be created in the home directory
- **AND** it SHALL use `projectFs` filesystem, not `homeFs`

### Requirement: spectr-accept-wo-spectr-bin Skill

The provider SHALL install the `spectr-accept-wo-spectr-bin` skill for accepting
changes without the spectr binary.

#### Scenario: Accept skill installation path

- **WHEN** the provider returns initializers
- **THEN** it SHALL include an `AgentSkillsInitializer` for
  `spectr-accept-wo-spectr-bin`
- **AND** the skill SHALL be installed at `.codex/skills/spectr-accept-wo-spectr-bin/`
- **AND** the skill SHALL be installed in the project filesystem (not home)

#### Scenario: Accept skill structure

- **WHEN** the skill is installed
- **THEN** it SHALL create `SKILL.md` with valid AgentSkills frontmatter
- **AND** it SHALL create `scripts/accept.sh` for task conversion
- **AND** `scripts/accept.sh` SHALL be executable

#### Scenario: Accept skill SKILL.md content

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
  `.codex/skills/spectr-validate-wo-spectr-bin/`
- **AND** the skill SHALL be installed in the project filesystem (not home)

#### Scenario: Validate skill structure

- **WHEN** the skill is installed
- **THEN** it SHALL create `SKILL.md` with valid AgentSkills frontmatter
- **AND** it SHALL create `scripts/validate.sh` for specification validation
- **AND** `scripts/validate.sh` SHALL be executable (0755 permissions)

#### Scenario: Validate skill SKILL.md content

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

#### Scenario: validate.sh functionality

- **WHEN** `scripts/validate.sh` is executed
- **THEN** it SHALL validate spec files for required sections and formatting
- **AND** it SHALL validate change deltas for
  ADDED/MODIFIED/REMOVED/RENAMED requirements
- **AND** it SHALL check that requirements have scenarios and SHALL/MUST
  statements
- **AND** it SHALL support `--spec <id>`, `--change <id>`, and `--all` modes
- **AND** it SHALL support `--json` output format
- **AND** it SHALL exit with code 0 on success, 1 on failure, 2 on usage error

### Requirement: Initializers Ordering for Codex

The provider SHALL return initializers in the correct order to ensure directories
exist before files are created.

#### Scenario: Directory initializers first

- **WHEN** the Codex provider's `Initializers()` method is called
- **THEN** `HomeDirectoryInitializer` for `.codex/prompts` SHALL be first
- **AND** `DirectoryInitializer` for `.codex/skills` SHALL be second
- **AND** directory initializers SHALL come before file-creating initializers

#### Scenario: Skills installed after directory creation

- **WHEN** the Codex provider's initializers are executed
- **THEN** `.codex/skills/` directory SHALL be created before skills are installed
- **AND** `AgentSkillsInitializer` instances SHALL come after
  `DirectoryInitializer` for `.codex/skills/`
- **AND** this ensures the target directory exists when skills are copied
