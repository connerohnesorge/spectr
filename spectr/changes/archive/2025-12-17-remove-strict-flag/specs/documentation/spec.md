# Documentation Delta Spec

## MODIFIED Requirements

### Requirement: Complete Command Reference

The system SHALL document all CLI commands with flags, examples, and expected output. Documentation SHALL only reference commands that actually exist in the CLI. Documentation SHALL distinguish between instructions for human users and AI agents.

#### Scenario: User learns init command usage

- **WHEN** a user reads the init command documentation
- **THEN** they SHALL see all available flags (`--path`, `--tools`, `--non-interactive`) with explanations and examples

#### Scenario: User learns list command options

- **WHEN** a user reads the list command documentation
- **THEN** they SHALL understand the `--specs`, `--json`, and `--long` flags with example outputs

#### Scenario: User learns validate command options

- **WHEN** a user reads the validate command documentation
- **THEN** they SHALL understand that validation always treats warnings as errors
- **AND** they SHALL see available flags (`--json`, `--type`, `--all`, `--changes`, `--specs`, `--no-interactive`) with explanations

#### Scenario: User learns archive command

- **WHEN** a user reads the archive command documentation
- **THEN** they SHALL understand the archiving workflow and `--skip-specs` flag usage

#### Scenario: User learns view command options

- **WHEN** a user reads the view command documentation
- **THEN** they SHALL understand the `--json` flag for dashboard output

#### Scenario: Documentation accuracy

- **WHEN** a user or AI assistant reads any documentation file
- **THEN** all referenced CLI commands SHALL exist and work as documented
- **AND** no nonexistent commands (such as `spectr show`) SHALL be referenced

#### Scenario: AI agent documentation uses direct file reading

- **WHEN** an AI agent reads documentation for viewing specs or changes
- **THEN** the documentation SHALL instruct agents to read files directly (e.g., `spectr/specs/<capability>/spec.md` or `spectr/changes/<change-id>/proposal.md`)
- **AND** the documentation SHALL NOT instruct agents to use CLI commands like `spectr view` for reading content
- **AND** the `spectr view` command SHALL only be documented for human users
