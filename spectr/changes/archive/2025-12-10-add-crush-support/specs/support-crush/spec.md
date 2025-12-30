# Delta Specification

## ADDED Requirements

### Requirement: Crush Provider Configuration

The provider SHALL be configured with these settings:

- ID: `crush`
- Name: `Crush`
- Priority: 16 (after existing providers)
- Config File: `CRUSH.md`
- Command Format: Markdown

#### Scenario: Provider identification

- **WHEN** the registry queries for Crush provider
- **THEN** it SHALL return provider with ID `crush`
- **AND** the provider priority is 16

#### Scenario: Provider metadata

- **WHEN** displaying provider options to users
- **THEN** the provider name is "Crush"
- **AND** it appears in the provider list after Continue (priority 15)

### Requirement: Crush Instruction File

The provider SHALL create and maintain a `CRUSH.md` instruction file in the
project root.

#### Scenario: Instruction file creation

- **WHEN** `spectr init` runs with Crush provider selected
- **THEN** the system creates `CRUSH.md` in project root
- **AND** inserts Spectr instructions between `<!-- spectr:START -->` and `<!--
  spectr:END -->` markers

#### Scenario: Instruction file updates

- **WHEN** `spectr init` runs in a project with Crush provider
- **THEN** the system updates content between markers in `CRUSH.md`
- **AND** preserves any user content outside the markers

### Requirement: Crush Slash Commands

The provider SHALL create slash commands in `.crush/commands/spectr/` directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** it creates `.crush/commands/spectr/` directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths

- **WHEN** the provider generates slash command files
- **THEN** it creates `.crush/commands/spectr/proposal.md`
- **AND** it creates `.crush/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

### Requirement: Standard Frontmatter

The provider SHALL use standard frontmatter templates for each command type.

#### Scenario: Proposal command frontmatter

- **WHEN** generating the proposal command file
- **THEN** the frontmatter description is "Scaffold a new Spectr change and
  validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** generating the apply command file
- **THEN** the frontmatter description is "Implement an approved Spectr change
  and keep tasks in sync."

### Requirement: Command Path Construction

The provider SHALL construct command paths using the standard pattern.

#### Scenario: Path construction for Crush

- **WHEN** determining command file paths
- **THEN** it uses base directory `.crush/commands`
- **AND** appends `/spectr/` subdirectory
- **AND** appends command name with `.md` extension
- **AND** results in paths like `.crush/commands/spectr/proposal.md`
