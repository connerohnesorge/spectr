## MODIFIED Requirements

### Requirement: Cursor Slash Commands
The provider SHALL create slash commands in `.cursorrules/commands/spectr/` directory.

#### Scenario: Command paths
- **WHEN** the provider configures slash commands
- **THEN** it creates `.cursorrules/commands/spectr/proposal.md`
- **AND** it creates `.cursorrules/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with YAML frontmatter
- **AND** frontmatter includes `description` field
