# Provider System Specification (Delta)

## ADDED Requirements

### Requirement: TemplateManager Skill Template Accessors

The TemplateManager SHALL provide methods for accessing agent skill templates.

```go
// ProposalSkill returns the skill-proposal.md.tmpl template reference.
// Used by providers that generate Amp-style agent skills.
func (tm *TemplateManager) ProposalSkill() domain.TemplateRef

// ApplySkill returns the skill-apply.md.tmpl template reference.
// Used by providers that generate Amp-style agent skills.
func (tm *TemplateManager) ApplySkill() domain.TemplateRef
```

#### Scenario: ProposalSkill accessor

- **WHEN** `TemplateManager.ProposalSkill()` is called
- **THEN** it SHALL return a `domain.TemplateRef` for `skill-proposal.md.tmpl`
- **AND** the template SHALL be pre-parsed and ready for rendering

#### Scenario: ApplySkill accessor

- **WHEN** `TemplateManager.ApplySkill()` is called
- **THEN** it SHALL return a `domain.TemplateRef` for `skill-apply.md.tmpl`
- **AND** the template SHALL be pre-parsed and ready for rendering

#### Scenario: Skill template rendering

- **WHEN** a skill template is rendered via `TemplateManager.Render(templateName, data)`
- **THEN** it SHALL accept a `domain.TemplateContext` as the data parameter
- **AND** it SHALL substitute all template variables (BaseDir, SpecsDir, ChangesDir, etc.)
- **AND** it SHALL return the fully rendered skill content with YAML frontmatter

### Requirement: SkillFileInitializer

The system SHALL provide a `SkillFileInitializer` for creating individual SKILL.md files from templates.

```go
// SkillFileInitializer creates a SKILL.md file from a template.
type SkillFileInitializer struct {
    targetPath string           // target file path (e.g., ".agents/skills/spectr-proposal/SKILL.md")
    template   domain.TemplateRef // template to render for skill content
}

// NewSkillFileInitializer creates a SkillFileInitializer for the given path and template.
func NewSkillFileInitializer(
    targetPath string, template domain.TemplateRef) *SkillFileInitializer
```

#### Scenario: SkillFileInitializer construction

- **WHEN** a SkillFileInitializer is created via `NewSkillFileInitializer(targetPath, template)`
- **THEN** it SHALL receive a target file path ending in `SKILL.md`
- **AND** it SHALL receive a TemplateRef directly (not a function)
- **AND** the TemplateRef SHALL be resolved at provider construction time
- **AND** the initializer SHALL use `projectFs` for file operations

#### Scenario: Create new skill file

- **WHEN** the skill file does not exist
- **THEN** the initializer SHALL create the parent directory if needed
- **AND** SHALL create the SKILL.md file with rendered content
- **AND** SHALL return the file path in `ExecutionResult.CreatedFiles`

#### Scenario: Update existing skill file

- **WHEN** the skill file already exists
- **THEN** the initializer SHALL overwrite it with rendered content
- **AND** SHALL return the file path in `ExecutionResult.UpdatedFiles`
- **AND** this ensures idempotent execution

#### Scenario: IsSetup check for skill file

- **WHEN** `IsSetup()` is called on a SkillFileInitializer
- **THEN** it SHALL return `true` if the SKILL.md file exists at the target path
- **AND** SHALL return `false` if the file does not exist

#### Scenario: Deduplication key for skill file

- **WHEN** `dedupeKey()` is called on a SkillFileInitializer
- **THEN** it SHALL return `SkillFileInitializer:<targetPath>`
- **AND** the path SHALL be normalized with `filepath.Clean`
- **AND** multiple initializers with the same target path SHALL deduplicate

#### Scenario: Initializer ordering for skill files

- **WHEN** skill file initializers are executed
- **THEN** they SHALL run after `DirectoryInitializer` (to ensure parent directories exist)
- **AND** SHALL run after `ConfigFileInitializer` (instruction pointers first)
- **AND** SHALL run before or alongside other file creation initializers

### Requirement: Embedded Skill Template Location

Agent skill templates SHALL be embedded in `internal/domain/templates/`.

#### Scenario: Skill template embedding

- **WHEN** the template embed directive is processed
- **THEN** it SHALL include `skill-proposal.md.tmpl` and `skill-apply.md.tmpl`
- **AND** templates SHALL be located in `internal/domain/templates/`
- **AND** templates SHALL be accessible via `TemplateManager.ProposalSkill()` and `TemplateManager.ApplySkill()`

#### Scenario: Skill template parsing

- **WHEN** templates are loaded into TemplateManager
- **THEN** skill templates SHALL be parsed with the same template engine as other templates
- **AND** SHALL support the same template context variables (BaseDir, SpecsDir, etc.)
- **AND** parsing errors SHALL be reported during initialization
