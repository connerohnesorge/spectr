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

#### Scenario: Provider registration

- **WHEN** the Gemini provider is registered
- **THEN** it SHALL use the new Registration struct with metadata
- **AND** registration SHALL include ID `gemini`, Name `Gemini CLI`, Priority 2
- **AND** the Provider implementation SHALL return initializers

#### Scenario: Provider returns initializers

- **WHEN** the provider's `Initializers(ctx context.Context, tm
  *TemplateManager)` method is called
- **THEN** it SHALL return a `DirectoryInitializer` for
  `.gemini/commands/spectr/`
- **AND** it SHALL return a `TOMLSlashCommandsInitializer` for TOML format slash
  commands
- **AND** it SHALL NOT return a `ConfigFileInitializer` (Gemini has no
  instruction file)

### Requirement: No Instruction File

The Gemini provider SHALL NOT create an instruction file since Gemini CLI does
not support project-level instruction files.

#### Scenario: Config file check

- **WHEN** `HasConfigFile()` is called on Gemini provider
- **THEN** it returns false

### Requirement: Gemini Slash Commands

The provider SHALL create slash commands in `.gemini/commands/spectr/` directory
using TOML format.

#### Scenario: Command directory structure

- **WHEN** the provider returns initializers
- **THEN** DirectoryInitializer SHALL create `.gemini/commands/spectr/`
  directory

#### Scenario: Command paths

- **WHEN** the `TOMLSlashCommandsInitializer` executes
- **THEN** it SHALL create `.gemini/commands/spectr/proposal.toml`
- **AND** it SHALL create `.gemini/commands/spectr/apply.toml`

#### Scenario: TOML command format

- **WHEN** slash command files are created by `TOMLSlashCommandsInitializer`
- **THEN** they SHALL use TOML format with `.toml` extension
- **AND** it SHALL include `description` field with command description
- **AND** it SHALL include `prompt` field with command content

### Requirement: Custom TOML Generation

The Gemini provider SHALL override the base Configure method to generate TOML
files instead of Markdown.

#### Scenario: TOML content structure

- **WHEN** a TOML command file is generated
- **THEN** it includes a comment header `# Spectr command for Gemini CLI`
- **AND** description is a quoted string
- **AND** prompt uses TOML multiline string syntax (triple quotes)
