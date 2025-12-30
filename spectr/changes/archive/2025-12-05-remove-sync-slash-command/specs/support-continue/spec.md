# Delta Specification

## MODIFIED Requirements

### Requirement: Continue Slash Commands

The provider SHALL create slash commands in `.continue/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** commands SHALL be placed in `.continue/commands/spectr/` directory

#### Scenario: Command paths

- **WHEN** the provider creates slash command files
- **THEN** it SHALL create `.continue/commands/spectr/proposal.md`
- **AND** it SHALL create `.continue/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** files SHALL include YAML frontmatter
- **AND** frontmatter SHALL include `description` field
