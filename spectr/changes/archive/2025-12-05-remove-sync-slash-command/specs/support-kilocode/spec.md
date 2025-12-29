# Delta Specification

## MODIFIED Requirements

### Requirement: Kilocode Slash Commands

The provider SHALL create slash commands in `.kilocode/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** commands are placed in `.kilocode/commands/spectr/` subdirectory

#### Scenario: Command paths

- **WHEN** the provider configures slash commands
- **THEN** it creates `.kilocode/commands/spectr/proposal.md`
- **AND** it creates `.kilocode/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field

#### Scenario: Proposal command frontmatter

- **WHEN** proposal command is created
- **THEN** frontmatter contains description "Scaffold a new Spectr change and
  validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** apply command is created
- **THEN** frontmatter contains description "Implement an approved Spectr change
  and keep tasks in sync."
