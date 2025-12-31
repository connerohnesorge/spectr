# Support Codex Specification Delta

## ADDED Requirements

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
