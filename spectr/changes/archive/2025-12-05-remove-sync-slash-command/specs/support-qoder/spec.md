## MODIFIED Requirements

### Requirement: Qoder Slash Commands
The provider SHALL create slash commands in `.qoder/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider configures slash commands
- **THEN** it creates `.qoder/commands/spectr/` directory
- **AND** all Spectr commands are placed under this directory

#### Scenario: Standard command paths
- **WHEN** the provider generates command file paths
- **THEN** it creates `.qoder/commands/spectr/proposal.md`
- **AND** it creates `.qoder/commands/spectr/apply.md`

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter at the top
- **AND** frontmatter includes a `description` field

#### Scenario: Proposal command frontmatter
- **WHEN** the proposal command file is created
- **THEN** frontmatter description SHALL be "Scaffold a new Spectr change and validate strictly."

#### Scenario: Apply command frontmatter
- **WHEN** the apply command file is created
- **THEN** frontmatter description SHALL be "Implement an approved Spectr change and keep tasks in sync."
