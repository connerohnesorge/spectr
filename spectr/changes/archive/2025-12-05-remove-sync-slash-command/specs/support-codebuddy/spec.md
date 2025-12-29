## MODIFIED Requirements

### Requirement: CodeBuddy Slash Commands

The provider SHALL create slash commands in `.codebuddy/commands/spectr/` directory.

#### Scenario: Command directory structure

- **WHEN** the provider initializes slash commands
- **THEN** it SHALL create `.codebuddy/commands/spectr/` directory
- **AND** all Spectr commands SHALL be placed in this subdirectory

#### Scenario: Command file paths

- **WHEN** the provider configures slash commands
- **THEN** it SHALL create `.codebuddy/commands/spectr/proposal.md`
- **AND** it SHALL create `.codebuddy/commands/spectr/apply.md`

#### Scenario: Command format

- **WHEN** slash command files are created
- **THEN** they SHALL use Markdown format with `.md` extension
- **AND** each file SHALL include YAML frontmatter
- **AND** frontmatter SHALL include a `description` field

#### Scenario: Proposal command frontmatter

- **WHEN** the proposal command file is created
- **THEN** it SHALL include frontmatter with description: "Scaffold a new Spectr change and validate strictly."

#### Scenario: Apply command frontmatter

- **WHEN** the apply command file is created
- **THEN** it SHALL include frontmatter with description: "Implement an approved Spectr change and keep tasks in sync."
