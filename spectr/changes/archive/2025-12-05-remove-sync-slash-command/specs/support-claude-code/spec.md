# Delta Specification

## MODIFIED Requirements

### Requirement: Claude Code Slash Commands

The provider SHALL create slash commands in `.claude/commands/spectr/`
directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** it creates `.claude/commands/spectr/` directory
- **AND** all Spectr commands are placed in this subdirectory

#### Scenario: Command paths

- **WHEN** the provider generates slash command files
- **THEN** it creates `.claude/commands/spectr/proposal.md`
- **AND** it creates `.claude/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field
