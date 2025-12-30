# Delta Specification

## MODIFIED Requirements

### Requirement: Windsurf Slash Commands

The provider SHALL create slash commands in `.windsurf/commands/spectr/`
directory.

#### Scenario: Command paths

- **WHEN** the provider configures slash commands
- **THEN** it SHALL create `.windsurf/commands/spectr/proposal.md`
- **AND** it SHALL create `.windsurf/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with YAML frontmatter
- **AND** frontmatter SHALL include a `description` field

#### Scenario: Proposal command frontmatter

- **WHEN** the proposal command is created
- **THEN** the frontmatter description SHALL be "Scaffold a new Spectr change
  and validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** the apply command is created
- **THEN** the frontmatter description SHALL be "Implement an approved Spectr
  change and keep tasks in sync."
