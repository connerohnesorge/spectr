## MODIFIED Requirements

### Requirement: CoStrict Slash Commands

The provider SHALL create slash commands in `.costrict/commands/spectr/` directory.

#### Scenario: Command paths

- **WHEN** the provider configures slash commands
- **THEN** it creates `.costrict/commands/spectr/proposal.md`
- **AND** it creates `.costrict/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
