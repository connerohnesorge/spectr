# Spec Delta

## ADDED Requirements

### Requirement: Kimi Provider Implementation

The system SHALL provide a Kimi provider that initializes Kimi CLI configuration
directories and files.

#### Scenario: Initialize Kimi directories

- **WHEN** `spectr init` runs in a project directory
- **THEN** the Kimi provider creates `.kimi/skills` directory
- **AND** creates `.kimi/commands` directory

#### Scenario: Create Kimi instruction file

- **WHEN** the Kimi provider initializes
- **THEN** it creates `AGENTS.md` file in the project root
- **AND** the file contains Kimi-specific instructions for working with Spectr

### Requirement: Kimi Slash Commands

The system SHALL create slash command files for Kimi CLI to support
Spectr workflows.

#### Scenario: Create spectr-proposal command

- **WHEN** the Kimi provider initializes
- **THEN** it creates `.kimi/commands/spectr-proposal.md`
- **AND** the file contains instructions for creating Spectr proposals

#### Scenario: Create spectr-apply command

- **WHEN** the Kimi provider initializes
- **THEN** it creates `.kimi/commands/spectr-apply.md`
- **AND** the file contains instructions for applying Spectr changes

### Requirement: Kimi Agent Skills

The system SHALL install Spectr agent skills in the Kimi skills
directory.

#### Scenario: Install spectr-accept-wo-spectr-bin skill

- **WHEN** the Kimi provider initializes
- **THEN** it installs the `spectr-accept-wo-spectr-bin` skill to
  `.kimi/skills/spectr-accept-wo-spectr-bin/`
- **AND** the skill includes SKILL.md and any executable scripts
- **AND** scripts maintain executable permissions after installation

#### Scenario: Install spectr-validate-wo-spectr-bin skill

- **WHEN** the Kimi provider initializes
- **THEN** it installs the `spectr-validate-wo-spectr-bin` skill to
  `.kimi/skills/spectr-validate-wo-spectr-bin/`
- **AND** the skill includes SKILL.md and any executable scripts
- **AND** scripts maintain executable permissions after installation

### Requirement: Kimi Provider Registration

The system SHALL register the Kimi provider in the provider registry with
appropriate priority.

#### Scenario: Provider is discoverable

- **WHEN** the provider registry is initialized
- **THEN** the Kimi provider is registered with ID "kimi"
- **AND** has name "Kimi"
- **AND** has appropriate priority for initialization order
