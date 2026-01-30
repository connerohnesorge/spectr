# Amp Support Specification

## Purpose

Defines Spectr support for Amp (ampcode.com), a production-grade AI coding assistant based on Claude Code that uses agent skills for extensibility.

## ADDED Requirements

### Requirement: Amp Provider

The system SHALL provide an `AmpProvider` that generates Amp-compatible agent skills in `.agents/skills/`.

#### Scenario: Provider returns skill initializers

- **WHEN** `AmpProvider.Initializers(ctx, tm)` is called
- **THEN** it SHALL return initializers for:
  - `.agents/skills/spectr-proposal/` directory creation
  - `.agents/skills/spectr-apply/` directory creation
  - `.agents/skills/spectr-proposal/SKILL.md` file creation
  - `.agents/skills/spectr-apply/SKILL.md` file creation
  - `.agents/skills/spectr-accept-wo-spectr-bin/` embedded skill
  - `.agents/skills/spectr-validate-wo-spectr-bin/` embedded skill

#### Scenario: Skill directory structure

- **WHEN** Amp initializers execute
- **THEN** they SHALL create the following structure:
  ```
  .agents/skills/
  ├── spectr-proposal/
  │   └── SKILL.md
  ├── spectr-apply/
  │   └── SKILL.md
  ├── spectr-accept-wo-spectr-bin/
  │   ├── SKILL.md
  │   └── scripts/accept.sh
  └── spectr-validate-wo-spectr-bin/
      ├── SKILL.md
      └── scripts/validate.sh
  ```

### Requirement: Amp Skill Frontmatter

Agent skills for Amp SHALL use YAML frontmatter with `name` and `description` fields.

#### Scenario: Proposal skill frontmatter

- **WHEN** the spectr-proposal skill is generated
- **THEN** the SKILL.md SHALL contain frontmatter:
  ```yaml
  ---
  name: spectr-proposal
  description: Create a Spectr change proposal with delta specs and tasks
  ---
  ```
- **AND** the frontmatter SHALL be followed by instructional content

#### Scenario: Apply skill frontmatter

- **WHEN** the spectr-apply skill is generated
- **THEN** the SKILL.md SHALL contain frontmatter:
  ```yaml
  ---
  name: spectr-apply
  description: Apply a Spectr change proposal by converting tasks.md to tasks.jsonc
  ---
  ```
- **AND** the frontmatter SHALL be followed by instructional content

#### Scenario: Accept-without-binary skill frontmatter

- **WHEN** the spectr-accept-wo-spectr-bin skill is generated
- **THEN** the SKILL.md SHALL contain frontmatter:
  ```yaml
  ---
  name: spectr-accept-wo-spectr-bin
  description: Accept Spectr change proposals without requiring the spectr binary
  ---
  ```
- **AND** the skill SHALL include embedded `scripts/accept.sh`

#### Scenario: Validate-without-binary skill frontmatter

- **WHEN** the spectr-validate-wo-spectr-bin skill is generated
- **THEN** the SKILL.md SHALL contain frontmatter:
  ```yaml
  ---
  name: spectr-validate-wo-spectr-bin
  description: Validate Spectr specifications without requiring the spectr binary
  ---
  ```
- **AND** the skill SHALL include embedded `scripts/validate.sh`

### Requirement: Skill Templates

The system SHALL provide skill templates in `internal/domain/templates/` for Amp skills.

#### Scenario: Skill template files

- **WHEN** skill templates are embedded
- **THEN** the following templates SHALL exist:
  - `skill-proposal.md.tmpl` - Proposal creation skill
  - `skill-apply.md.tmpl` - Proposal application skill

#### Scenario: Template variable substitution

- **WHEN** skill templates are rendered
- **THEN** they SHALL support the following template variables:
  - `{{.BaseDir}}` - Spectr base directory (e.g., "spectr")
  - `{{.SpecsDir}}` - Specs directory (e.g., "spectr/specs")
  - `{{.ChangesDir}}` - Changes directory (e.g., "spectr/changes")
  - `{{.ProjectFile}}` - Project file path (e.g., "spectr/project.md")
  - `{{.AgentsFile}}` - Agents file path (e.g., "spectr/AGENTS.md")

#### Scenario: Template content structure

- **WHEN** a skill template is rendered
- **THEN** it SHALL include:
  - YAML frontmatter with `name` and `description`
  - Clear usage instructions for agents
  - Examples of skill invocation
  - References to Spectr conventions and files

### Requirement: User-Invocable Skills

Skills SHALL support user invocation via slash command syntax.

#### Scenario: User invokes proposal skill

- **WHEN** a user types `/spectr:proposal` in Amp
- **THEN** Amp SHALL load the spectr-proposal skill
- **AND** the skill content SHALL be injected into agent context
- **AND** the agent SHALL follow the skill instructions to create a proposal

#### Scenario: User invokes apply skill

- **WHEN** a user types `/spectr:apply` in Amp
- **THEN** Amp SHALL load the spectr-apply skill
- **AND** the skill content SHALL be injected into agent context
- **AND** the agent SHALL follow the skill instructions to apply the proposal

### Requirement: Amp Provider Registration

The Amp provider SHALL be registered in the global provider registry.

#### Scenario: Provider registration

- **WHEN** `RegisterAllProviders()` is called
- **THEN** it SHALL register the Amp provider with:
  - ID: `amp`
  - Name: `Amp`
  - Priority: `15` (after Claude Code at 10, before Gemini at 20)
  - Provider: `&AmpProvider{}`

### Requirement: Embedded Skill Support

The Amp provider SHALL use `AgentSkillsInitializer` to copy embedded skills.

#### Scenario: Copy accept skill

- **WHEN** Amp provider initializers execute
- **THEN** `AgentSkillsInitializer` SHALL copy `spectr-accept-wo-spectr-bin` skill
- **AND** SHALL preserve directory structure
- **AND** SHALL preserve file permissions (executable scripts)

#### Scenario: Copy validate skill

- **WHEN** Amp provider initializers execute
- **THEN** `AgentSkillsInitializer` SHALL copy `spectr-validate-wo-spectr-bin` skill
- **AND** SHALL preserve directory structure
- **AND** SHALL preserve file permissions (executable scripts)

### Requirement: TemplateManager Skill Methods

The TemplateManager SHALL provide methods to access Amp skill templates.

#### Scenario: Skill template accessors

- **WHEN** `TemplateManager` is queried for Amp skill templates
- **THEN** it SHALL provide:
  - `ProposalSkill() domain.TemplateRef` - Returns skill-proposal.md.tmpl
  - `ApplySkill() domain.TemplateRef` - Returns skill-apply.md.tmpl
- **AND** these methods SHALL return `domain.TemplateRef` with pre-parsed templates

#### Scenario: Template rendering for skills

- **WHEN** a skill template is rendered via `tm.Render()`
- **THEN** it SHALL accept a `domain.TemplateContext` with path variables
- **AND** it SHALL substitute all template variables
- **AND** it SHALL return the fully rendered skill content

### Requirement: Compatibility with Claude Code

Amp skills SHALL coexist with Claude Code skills without conflicts.

#### Scenario: Dual provider selection

- **WHEN** a user selects both Claude Code and Amp during `spectr init`
- **THEN** Claude Code skills SHALL be generated in `.claude/skills/`
- **AND** Amp skills SHALL be generated in `.agents/skills/`
- **AND** both sets of skills SHALL function independently

#### Scenario: Skill deduplication

- **WHEN** embedded skills are copied for both providers
- **THEN** deduplication SHALL occur based on target directory
- **AND** `.claude/skills/spectr-accept-wo-spectr-bin/` SHALL be separate from `.agents/skills/spectr-accept-wo-spectr-bin/`
- **AND** both SHALL be created if both providers are selected

### Requirement: Skill Discovery by Amp

Amp SHALL automatically discover skills in `.agents/skills/` directories.

#### Scenario: Project-local skill discovery

- **WHEN** Amp loads a project with `.agents/skills/` directory
- **THEN** Amp SHALL discover all skills with `SKILL.md` files
- **AND** skills SHALL be available for agent loading
- **AND** skills SHALL be available for user invocation

#### Scenario: Backward compatibility with .claude/skills

- **WHEN** Amp loads a project with `.claude/skills/` directory
- **THEN** Amp SHALL discover skills in `.claude/skills/` as well
- **AND** `.agents/skills/` SHALL take precedence over `.claude/skills/`
- **AND** skills from both locations SHALL be available
