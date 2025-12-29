# Delta Specification

## MODIFIED Requirements

### Requirement: Aider Slash Commands

The provider SHALL create slash commands in `.aider/commands/spectr/` directory.

#### Scenario: Command paths

- **WHEN** the provider configures slash commands
- **THEN** it creates `.aider/commands/spectr/proposal.md`
- **AND** it creates `.aider/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
