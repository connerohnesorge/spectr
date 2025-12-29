## MODIFIED Requirements

### Requirement: Antigravity Slash Commands

The provider SHALL create slash commands in `.agent/workflows/` directory.

#### Scenario: Command directory structure

- **WHEN** the provider configures slash commands
- **THEN** it uses `.agent/workflows/` as base directory (not `.agent/commands/`)
- **AND** all Spectr commands reside in `.agent/workflows/` subdirectory

#### Scenario: Command file paths

- **WHEN** the provider creates slash command files
- **THEN** it creates `.agent/workflows/spectr-proposal.md`
- **AND** it creates `.agent/workflows/spectr-apply.md`

#### Scenario: Command file format

- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter at the top
- **AND** frontmatter includes a `description` field
