## ADDED Requirements

### Requirement: spectr-validate-w-spectr-bin Skill

The provider SHALL install the `spectr-validate-w-spectr-bin` skill for validating specifications and change proposals using the installed spectr binary.

#### Scenario: Validate (binary) skill installation path
- **WHEN** the provider returns initializers
- **THEN** it SHALL include an `AgentSkillsInitializer` for `spectr-validate-w-spectr-bin`
- **AND** the skill SHALL be installed at `.codex/skills/spectr-validate-w-spectr-bin/`
- **AND** the skill SHALL be installed in the project filesystem (not home)

#### Scenario: Validate (binary) skill structure
- **WHEN** the skill is installed
- **THEN** it SHALL create `SKILL.md` with valid AgentSkills frontmatter
- **AND** it SHALL create `scripts/validate.sh` that wraps the spectr binary
- **AND** `scripts/validate.sh` SHALL be executable

#### Scenario: Validate (binary) SKILL.md content
- **WHEN** the `SKILL.md` file is created
- **THEN** the frontmatter `name` SHALL be `spectr-validate-w-spectr-bin`
- **AND** the frontmatter `description` SHALL state it uses the spectr binary
- **AND** the `compatibility` section SHALL list `spectr` as a requirement
- **AND** the body SHALL contain usage instructions matching `spectr validate`

#### Scenario: validate.sh (binary) functionality
- **WHEN** `scripts/validate.sh` is executed
- **THEN** it SHALL forward all arguments to `spectr validate`
- **AND** it SHALL execute `spectr validate "$@"`
