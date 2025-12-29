# Support Gemini Specification

## Purpose

Documents the Gemini CLI provider integration for Spectr.

## Requirements

### Requirement: Gemini Provider Configuration

The provider SHALL be configured with these settings:

- ID: `gemini`
- Name: `Gemini CLI`
- Priority: 2
- Config File: (none - Gemini has no instruction file)
- Command Format: TOML

#### Scenario: Provider identification

- **WHEN** the registry queries for Gemini provider
- **THEN** it SHALL return provider with ID `gemini`

### Requirement: No Instruction File

The Gemini provider SHALL NOT create an instruction file since Gemini CLI does not support project-level instruction files.

#### Scenario: Config file check

- **WHEN** `HasConfigFile()` is called on Gemini provider
- **THEN** it returns false

### Requirement: Gemini Slash Commands

The provider SHALL create slash commands in `.gemini/commands/spectr/` directory using TOML format.

#### Scenario: Command paths

- **WHEN** the provider configures slash commands
- **THEN** it creates `.gemini/commands/spectr/proposal.toml`
- **AND** it creates `.gemini/commands/spectr/apply.toml`

#### Scenario: TOML command format

- **WHEN** slash command files are created
- **THEN** they use TOML format
- **AND** include `description` field with command description
- **AND** include `prompt` field with command content

### Requirement: Custom TOML Generation

The Gemini provider SHALL override the base Configure method to generate TOML files instead of Markdown.

#### Scenario: TOML content structure

- **WHEN** a TOML command file is generated
- **THEN** it includes a comment header `# Spectr command for Gemini CLI`
- **AND** description is a quoted string
- **AND** prompt uses TOML multiline string syntax (triple quotes)
