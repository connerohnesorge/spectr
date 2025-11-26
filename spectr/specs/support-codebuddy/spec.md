# Support CodeBuddy Specification

## Purpose
Documents the CodeBuddy provider integration for Spectr, enabling CodeBuddy AI coding assistant to work with Spectr's spec-driven development workflow.

## Requirements

### Requirement: CodeBuddy Provider Configuration
The provider SHALL be configured with these settings:
- ID: `codebuddy`
- Name: `CodeBuddy`
- Priority: 5
- Config File: `CODEBUDDY.md`
- Command Format: Markdown

#### Scenario: Provider identification
- **WHEN** the registry queries for CodeBuddy provider
- **THEN** it SHALL return provider with ID `codebuddy`
- **AND** the provider name SHALL be `CodeBuddy`

#### Scenario: Provider priority
- **WHEN** multiple providers are available
- **THEN** CodeBuddy SHALL have priority 5
- **AND** it SHALL be sorted accordingly in provider selection UI

### Requirement: CodeBuddy Instruction File
The provider SHALL create and maintain a `CODEBUDDY.md` instruction file in the project root.

#### Scenario: Instruction file creation
- **WHEN** `spectr init` runs with CodeBuddy provider selected
- **THEN** the system SHALL create `CODEBUDDY.md` in project root
- **AND** it SHALL insert Spectr instructions between `<!-- spectr:START -->` and `<!-- spectr:END -->` markers

#### Scenario: Instruction file update
- **WHEN** `spectr update` runs in a project with CodeBuddy provider
- **THEN** the system SHALL update the managed block in `CODEBUDDY.md`
- **AND** it SHALL preserve any user content outside the markers

### Requirement: CodeBuddy Slash Commands
The provider SHALL create slash commands in `.codebuddy/commands/spectr/` directory.

#### Scenario: Command directory structure
- **WHEN** the provider initializes slash commands
- **THEN** it SHALL create `.codebuddy/commands/spectr/` directory
- **AND** all Spectr commands SHALL be placed in this subdirectory

#### Scenario: Command file paths
- **WHEN** the provider configures slash commands
- **THEN** it SHALL create `.codebuddy/commands/spectr/proposal.md`
- **AND** it SHALL create `.codebuddy/commands/spectr/sync.md`
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

#### Scenario: Sync command frontmatter
- **WHEN** the sync command file is created
- **THEN** it SHALL include frontmatter with description: "Detect spec drift from code and update specs interactively."

### Requirement: CodeBuddy Command Content
The provider SHALL use standard Spectr command instructions for each slash command.

#### Scenario: Command instruction templates
- **WHEN** slash command files are generated
- **THEN** they SHALL include the appropriate template content from `internal/init/commands/`
- **AND** the content SHALL guide CodeBuddy through the Spectr workflow

### Requirement: CodeBuddy Provider Registration
The provider SHALL be automatically registered with the provider registry on package initialization.

#### Scenario: Automatic registration
- **WHEN** the providers package is imported
- **THEN** the CodeBuddy provider SHALL be registered via init() function
- **AND** it SHALL be available for selection during `spectr init`
