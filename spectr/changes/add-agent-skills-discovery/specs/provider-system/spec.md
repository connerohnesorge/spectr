# Provider System Specification Delta

## ADDED Requirements

### Requirement: Skill Discovery

The system SHALL provide discovery capabilities for embedded Agent Skills following the [agentskills.io](https://agentskills.io) specification.

#### Scenario: List embedded skills with metadata

- **WHEN** embedded skills are listed via `ListEmbeddedSkills(skillFS)`
- **THEN** the system SHALL return a slice of `SkillMetadata` containing skill information
- **AND** the system SHALL parse YAML frontmatter from each skill's `SKILL.md` file
- **AND** missing optional frontmatter fields SHALL default to empty slices
- **AND** skills SHALL be sorted alphabetically by name
- **AND** the function SHALL return an error if skill parsing fails

#### Scenario: Parse skill frontmatter

- **WHEN** a `SKILL.md` file is parsed via `ParseSkillMetadata(content)`
- **THEN** the system SHALL extract the following fields from YAML frontmatter:
  - `name` (string, required) - skill identifier
  - `description` (string, required) - short description
  - `compatibility.requirements` (string array, optional) - required dependencies
  - `compatibility.optional` (string array, optional) - optional dependencies
  - `compatibility.platforms` (string array, optional) - supported platforms
- **AND** the frontmatter SHALL be delimited by `---` markers at the start and end
- **AND** the markers SHALL appear before any other content in the file
- **AND** invalid YAML SHALL return an error with parser context
- **AND** missing required fields SHALL return an error specifying which field is missing

#### Scenario: Discover skills from embedded filesystem

- **WHEN** skills are discovered from an embedded filesystem
- **THEN** the system SHALL scan the root directory of the provided `skillFS`
- **AND** each subdirectory containing a `SKILL.md` file SHALL be recognized as a skill
- **AND** subdirectories without `SKILL.md` SHALL be skipped with a warning logged
- **AND** the system SHALL read and parse each `SKILL.md` file using `ParseSkillMetadata`
- **AND** parsing errors for individual skills SHALL not stop discovery of other skills
- **AND** skills with parse errors SHALL be excluded from results with errors logged

#### Scenario: Check skill installation status

- **WHEN** checking if a skill is installed via `IsSkillInstalled(projectPath, skillName)`
- **THEN** the system SHALL check if `SKILL.md` exists in `.claude/skills/<skillName>/` relative to the project root
- **AND** the check SHALL use the provided project path as the base directory
- **AND** the function SHALL return `true` if the file exists
- **AND** the function SHALL return `false` if the file does not exist
- **AND** the function SHALL return an error only for unexpected filesystem errors (not for file-not-found)

### Requirement: Skill Information Enrichment

The system SHALL enrich skill metadata with installation status for user-facing display.

#### Scenario: Convert metadata to skill info

- **WHEN** `SkillMetadata` is converted to `SkillInfo` via `ListSkills()`
- **THEN** the system SHALL copy all metadata fields (name, description, requirements, optional, platforms)
- **AND** the system SHALL check installation status using `IsSkillInstalled()`
- **AND** the `installed` field SHALL be set to `true` if the skill is installed in the project
- **AND** the `installed` field SHALL be set to `false` if the skill is not installed
- **AND** the `installed` field SHALL be set to `false` if installation check returns an error

#### Scenario: List skills with installation status

- **WHEN** skills are listed via `Lister.ListSkills(skillFS)`
- **THEN** the system SHALL discover all embedded skills using `ListEmbeddedSkills()`
- **AND** the system SHALL enrich each skill with installation status
- **AND** the system SHALL return a slice of `SkillInfo` objects
- **AND** the results SHALL be sorted alphabetically by name (preserved from discovery)
- **AND** the function SHALL return an error if skill discovery fails

### Requirement: CLI Skill Listing

The system SHALL provide a CLI command to list available agent skills with multiple output formats.

#### Scenario: List skills in text format

- **WHEN** `spectr list --skills` is executed without format flags
- **THEN** the system SHALL output skill names, one per line
- **AND** skills SHALL be sorted alphabetically by name
- **AND** output SHALL contain only skill names with no additional formatting
- **AND** the command SHALL succeed with exit code 0

#### Scenario: List skills in long format

- **WHEN** `spectr list --skills --long` is executed
- **THEN** the system SHALL output skill information in human-readable multi-line format
- **AND** each skill SHALL be separated by a blank line
- **AND** output for each skill SHALL include:
  - Skill name on the first line
  - Description indented with 2 spaces
  - Requirements (if present) indented with 2 spaces, comma-separated
  - Optional dependencies (if present) indented with 2 spaces, comma-separated
  - Platforms (if present) indented with 2 spaces, comma-separated
  - Installation status ("Installed" or "Not Installed") indented with 2 spaces
- **AND** the command SHALL succeed with exit code 0

#### Scenario: List skills in JSON format

- **WHEN** `spectr list --skills --json` is executed
- **THEN** the system SHALL output a JSON array of skill objects
- **AND** each object SHALL have the following fields:
  - `name` (string) - skill identifier
  - `description` (string) - skill description
  - `installed` (boolean) - installation status
  - `requirements` (string array) - required dependencies (may be empty)
  - `optional` (string array) - optional dependencies (may be empty)
  - `platforms` (string array) - supported platforms (may be empty)
- **AND** the JSON SHALL be pretty-printed with 2-space indentation
- **AND** the JSON SHALL be valid and parseable
- **AND** the command SHALL succeed with exit code 0

#### Scenario: Skills flag mutual exclusivity with specs

- **WHEN** `spectr list --skills` is combined with `--specs` flag
- **THEN** the system SHALL return an `IncompatibleFlagsError`
- **AND** the error message SHALL indicate that `--skills` and `--specs` are incompatible
- **AND** the command SHALL exit with a non-zero exit code

#### Scenario: Skills flag mutual exclusivity with all

- **WHEN** `spectr list --skills` is combined with `--all` flag
- **THEN** the system SHALL return an `IncompatibleFlagsError`
- **AND** the error message SHALL indicate that `--skills` and `--all` are incompatible
- **AND** the command SHALL exit with a non-zero exit code

#### Scenario: Skills flag mutual exclusivity with interactive

- **WHEN** `spectr list --skills` is combined with `--interactive` flag
- **THEN** the system SHALL return an `IncompatibleFlagsError`
- **AND** the error message SHALL indicate that `--skills` and `--interactive` are incompatible
- **AND** the command SHALL exit with a non-zero exit code

#### Scenario: Skills flag compatibility with output formats

- **WHEN** `spectr list --skills` is combined with `--long` or `--json` flags
- **THEN** the system SHALL successfully execute the command
- **AND** output SHALL be formatted according to the specified format flag
- **AND** the command SHALL succeed with exit code 0

## Cross-References

This change extends the existing `Requirement: AgentSkills Initializer` (provider-system spec, line 714) by adding discovery capabilities for skills that can be installed via `AgentSkillsInitializer`.

The discovery functionality uses the existing `TemplateManager.SkillFS()` method (provider-system spec, lines 785-786) to access embedded skill directories.

## Design Rationale

### Embedded-Only Discovery

The initial implementation discovers only embedded skills (those bundled with the Spectr binary). This keeps the implementation simple while covering the primary use case. Future enhancements can add discovery of installed skills from other sources (.codex/skills/, custom paths).

### Repository-Focused Installation Check

Installation status checks only `.claude/skills/` in the project root, not global paths like `~/.claude/skills/`. This aligns with Spectr's project-centric philosophy and avoids complexity around user home directory detection.

### YAML Frontmatter Format

The frontmatter format follows industry standards (Hugo, Jekyll, GitHub Actions) and the [agentskills.io](https://agentskills.io) specification, making skills portable across different AI coding tools.

### CLI Integration via List Command

Adding `--skills` to the existing `spectr list` command maintains consistency with `--specs` and `--all` flags, rather than creating a separate top-level command. This reduces cognitive load for users familiar with the list command.

### Read-Only Operations

Discovery operations are purely informational with no side effects. Installing or removing skills is intentionally excluded to keep the initial implementation focused and to avoid complexity around error handling, rollback, and validation. These operations can be added in future iterations.
