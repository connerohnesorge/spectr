# Agent Instructions Specification (Delta)

## ADDED Requirements

### Requirement: Amp User-Invocable Skills

Amp agents SHALL support user-invocable skills via slash command syntax (e.g., `/spectr:proposal`).

#### Scenario: User invokes skill with slash syntax

- **WHEN** a user types `/spectr:proposal` in Amp
- **THEN** Amp SHALL locate the `spectr-proposal` skill in `.agents/skills/`
- **AND** SHALL load the `SKILL.md` content into agent context
- **AND** the agent SHALL follow the skill instructions to create a proposal

#### Scenario: Agent autonomously loads skill

- **WHEN** an agent determines it needs to create a Spectr proposal
- **THEN** the agent MAY autonomously load the `spectr-proposal` skill
- **AND** SHALL inject the skill content into context
- **AND** SHALL follow the skill instructions

#### Scenario: Skill content injection

- **WHEN** a skill is loaded by Amp
- **THEN** the entire `SKILL.md` content (excluding frontmatter metadata) SHALL be injected into agent context
- **AND** the agent SHALL treat it as instructional guidance
- **AND** the agent SHALL follow the skill's step-by-step instructions

### Requirement: Amp Skill Discovery Locations

Amp SHALL discover agent skills from multiple locations in priority order.

#### Scenario: Project-local skill priority

- **WHEN** Amp discovers skills in a project
- **THEN** it SHALL search `.agents/skills/` first (highest priority)
- **AND** SHALL search `.claude/skills/` second (compatibility)
- **AND** SHALL search `~/.config/agents/skills/` last (global skills)

#### Scenario: Skill name collision resolution

- **WHEN** skills with the same name exist in multiple locations
- **THEN** Amp SHALL use the skill from the highest-priority location
- **AND** SHALL NOT merge or combine skills from different locations
- **AND** SHALL log which skill location was chosen (if logging is enabled)

### Requirement: Agent Skill Lazy Loading

Amp agents SHALL lazy-load skill content on-demand rather than preloading all skills.

#### Scenario: Context efficiency with lazy loading

- **WHEN** an agent session starts in Amp
- **THEN** skill content SHALL NOT be preloaded into context
- **AND** only skill names and descriptions SHALL be available initially
- **AND** full skill content SHALL be loaded only when explicitly invoked

#### Scenario: On-demand skill loading

- **WHEN** an agent or user invokes a skill (e.g., `/spectr:proposal`)
- **THEN** Amp SHALL read the `SKILL.md` file at that moment
- **AND** SHALL inject the content into the active context
- **AND** the content SHALL remain available for the duration of that task

### Requirement: Amp SKILL.md Frontmatter Requirements

Agent skills in Amp SHALL use minimal YAML frontmatter with required fields.

#### Scenario: Required frontmatter fields

- **WHEN** Amp loads a SKILL.md file
- **THEN** the file SHALL contain YAML frontmatter delimited by `---`
- **AND** the frontmatter SHALL include a `name` field (lowercase, kebab-case)
- **AND** the frontmatter SHALL include a `description` field (1-2 sentence summary)
- **AND** additional fields MAY be present but are optional

#### Scenario: Frontmatter parsing errors

- **WHEN** a SKILL.md file has malformed or missing frontmatter
- **THEN** Amp SHALL report an error indicating the skill could not be loaded
- **AND** SHALL NOT attempt to use the skill
- **AND** other skills SHALL remain available

### Requirement: Amp Skills for Spectr

Spectr-generated skills for Amp SHALL provide clear, actionable instructions for common workflows.

#### Scenario: Proposal skill instructions

- **WHEN** the `spectr-proposal` skill is loaded in Amp
- **THEN** it SHALL guide the agent through:
  - Reviewing `spectr/project.md` and existing specs
  - Choosing a unique verb-led change ID
  - Scaffolding proposal files (proposal.md, tasks.md, design.md)
  - Mapping changes to capabilities with delta specs
  - Validating with `spectr validate <id>`

#### Scenario: Apply skill instructions

- **WHEN** the `spectr-apply` skill is loaded in Amp
- **THEN** it SHALL guide the agent through:
  - Validating the proposal with `spectr validate <id>`
  - Running `spectr accept <id>` to convert tasks.md to tasks.jsonc
  - Verifying tasks.jsonc was created
  - Explaining next steps for implementation
