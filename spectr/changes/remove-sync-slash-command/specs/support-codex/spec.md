## MODIFIED Requirements

### Requirement: Codex Global Slash Commands
The provider SHALL create slash commands in the global `~/.codex/prompts/` directory.

#### Scenario: Global command directory structure
- **WHEN** the provider configures slash commands
- **THEN** it creates `~/.codex/prompts/spectr/` directory
- **AND** the directory is created in user's home directory, not project directory

#### Scenario: Command paths
- **WHEN** the provider generates slash command files
- **THEN** it creates `~/.codex/prompts/spectr-proposal.md`
- **AND** it creates `~/.codex/prompts/spectr-apply.md`

#### Scenario: Global path expansion
- **WHEN** resolving command paths
- **THEN** the `~` prefix is expanded to user's home directory
- **AND** paths work correctly regardless of current project directory

### Requirement: Codex Command Format
The provider SHALL use Markdown format with YAML frontmatter for slash commands.

#### Scenario: Command format
- **WHEN** slash command files are created
- **THEN** they use Markdown format with `.md` extension
- **AND** each file includes YAML frontmatter
- **AND** frontmatter includes `description` field

#### Scenario: Proposal command frontmatter
- **WHEN** generating the proposal command file
- **THEN** the frontmatter description is "Scaffold a new Spectr change and validate strictly."

#### Scenario: Apply command frontmatter
- **WHEN** generating the apply command file
- **THEN** the frontmatter description is "Implement an approved Spectr change and keep tasks in sync."

### Requirement: Codex Command Invocation
Users SHALL invoke Spectr commands in Codex using the `/spectr-<command>` pattern.

#### Scenario: Invoking proposal command
- **WHEN** user types `/spectr-proposal` in Codex
- **THEN** Codex loads and executes the proposal prompt

#### Scenario: Invoking apply command
- **WHEN** user types `/spectr-apply` in Codex
- **THEN** Codex loads and executes the apply prompt
