## ADDED Requirements

### Requirement: spectr-validate-wo-spectr-bin Skill

The provider SHALL install the `spectr-validate-wo-spectr-bin` skill for
validating specifications without the spectr binary.

#### Scenario: Skill installation path

- **WHEN** the provider returns initializers
- **THEN** it SHALL include an `AgentSkillsInitializer` for
  `spectr-validate-wo-spectr-bin`
- **AND** the skill SHALL be installed at
  `.claude/skills/spectr-validate-wo-spectr-bin/`

#### Scenario: Skill structure

- **WHEN** the skill is installed
- **THEN** it SHALL create `SKILL.md` with valid AgentSkills frontmatter
- **AND** it SHALL create `scripts/validate.sh` for specification validation
- **AND** `scripts/validate.sh` SHALL be executable

#### Scenario: SKILL.md content

- **WHEN** the `SKILL.md` file is created
- **THEN** the frontmatter `name` SHALL be `spectr-validate-wo-spectr-bin`
- **AND** the frontmatter `description` SHALL describe the skill's purpose
- **AND** the `compatibility` section SHALL note required tools (grep, sed, awk)
- **AND** the body SHALL contain usage instructions for the validate script

#### Scenario: validate.sh single spec validation

- **WHEN** `scripts/validate.sh` is executed with `--spec <spec-id>` argument
- **THEN** it SHALL validate `spectr/specs/<spec-id>/spec.md`
- **AND** it SHALL check for `## Requirements` section
- **AND** it SHALL check requirements contain SHALL or MUST
- **AND** it SHALL check requirements have `#### Scenario:` blocks
- **AND** it SHALL report errors with file path and line number

#### Scenario: validate.sh single change validation

- **WHEN** `scripts/validate.sh` is executed with `--change <change-id>` argument
- **THEN** it SHALL validate `spectr/changes/<change-id>/specs/*/spec.md`
- **AND** it SHALL check for delta sections (ADDED, MODIFIED, REMOVED, RENAMED)
- **AND** it SHALL validate requirements in delta sections
- **AND** it SHALL validate `tasks.md` if present

#### Scenario: validate.sh bulk validation

- **WHEN** `scripts/validate.sh` is executed with `--all` argument
- **THEN** it SHALL discover all specs in `spectr/specs/`
- **AND** it SHALL discover all changes in `spectr/changes/` (excluding archive)
- **AND** it SHALL validate each item and aggregate results
- **AND** it SHALL report summary with passed/failed counts

#### Scenario: validate.sh exit codes

- **WHEN** validation completes
- **THEN** exit code 0 SHALL indicate all validations passed
- **AND** exit code 1 SHALL indicate one or more validations failed
- **AND** exit code 2 SHALL indicate usage error (invalid arguments)

#### Scenario: validate.sh JSON output

- **WHEN** `scripts/validate.sh` is executed with `--json` flag
- **THEN** it SHALL output valid JSON
- **AND** the JSON SHALL include `items` array with per-item results
- **AND** the JSON SHALL include `summary` with totals
- **AND** human-readable output SHALL be suppressed
