# Support Gemini Specification

## Requirements

### Requirement: Gemini Skills Directory

The system SHALL create a `.gemini/skills/` directory in the project root for
storing Gemini CLI agent skills.

#### Scenario: Initialize skills directory

- **WHEN** `spectr init` runs in a project directory with Gemini provider
- **THEN** the system creates `.gemini/skills/` directory
- **AND** the directory has appropriate permissions (0755)

### Requirement: Gemini Agent Skills Installation

The system SHALL install Spectr agent skills in the Gemini skills directory
following the [Agent Skills specification][agentskills].

[agentskills]: https://agentskills.io

#### Scenario: Install spectr-accept-wo-spectr-bin skill

- **WHEN** the Gemini provider initializes
- **THEN** it installs the `spectr-accept-wo-spectr-bin` skill to
  `.gemini/skills/spectr-accept-wo-spectr-bin/`
- **AND** the skill includes SKILL.md with valid frontmatter (name, description)
- **AND** the skill includes scripts/accept.sh
- **AND** scripts maintain executable permissions (0755) after installation

#### Scenario: Install spectr-validate-wo-spectr-bin skill

- **WHEN** the Gemini provider initializes
- **THEN** it installs the `spectr-validate-wo-spectr-bin` skill to
  `.gemini/skills/spectr-validate-wo-spectr-bin/`
- **AND** the skill includes SKILL.md with valid frontmatter (name, description)
- **AND** the skill includes scripts/validate.sh
- **AND** scripts maintain executable permissions (0755) after installation

#### Scenario: Skills are idempotent

- **WHEN** `spectr init` runs multiple times
- **THEN** skills are overwritten with latest templates
- **AND** no duplicate files are created
- **AND** the operation completes without error

### Requirement: Gemini Instruction File

The system SHALL create a GEMINI.md instruction file in the project root to
provide workspace-wide guidance for Gemini CLI.

#### Scenario: Create GEMINI.md with Spectr instructions

- **WHEN** the Gemini provider initializes
- **THEN** it creates `GEMINI.md` file in the project root
- **AND** the file contains Spectr workflow guidance
- **AND** the file includes managed markers for automatic updates

#### Scenario: GEMINI.md preserves user content

- **WHEN** `spectr init` runs and GEMINI.md already exists with user content
  outside managed markers
- **THEN** the user content is preserved
- **AND** only content within managed markers is updated

### Requirement: Gemini TOML Commands Backward Compatibility

The system SHALL maintain existing TOML-based slash commands alongside the new
agent skills.

#### Scenario: Preserve existing slash commands

- **WHEN** the Gemini provider initializes
- **THEN** `.gemini/commands/spectr/` directory is maintained
- **AND** TOML slash command files (proposal.toml, apply.toml) are preserved
- **AND** both slash commands and skills are available for use

### Requirement: Gemini Provider Skill Registration

The system SHALL register the Gemini provider's skill initializers in the
correct initialization order.

#### Scenario: Skills initialize after directories

- **WHEN** the Gemini provider initializes
- **THEN** the `.gemini/skills/` directory is created first
- **AND** then the skills are copied into the directory
- **AND** the initialization order prevents errors

#### Scenario: Provider initializers are complete

- **WHEN** the Gemini provider's Initializers() method is called
- **THEN** it returns initializers for:
  - `.gemini/commands/spectr/` directory (existing)
  - `.gemini/skills/` directory (new)
  - TOML slash commands (existing)
  - GEMINI.md instruction file (new)
  - spectr-accept-wo-spectr-bin skill (new)
  - spectr-validate-wo-spectr-bin skill (new)
