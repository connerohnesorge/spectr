## ADDED Requirements

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
