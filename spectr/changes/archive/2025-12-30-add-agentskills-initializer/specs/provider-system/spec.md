## ADDED Requirements

### Requirement: AgentSkills Initializer

The system SHALL provide a built-in `AgentSkillsInitializer` for copying
embedded skill directories to target paths.

```go
// AgentSkillsInitializer copies an embedded skill directory to a target path.
type AgentSkillsInitializer struct {
    skillName string // name of the skill (matches embedded directory name)
    targetDir string // target directory path (e.g., ".claude/skills/my-skill")
}

// NewAgentSkillsInitializer creates an AgentSkillsInitializer for the given
// skill name and target directory.
func NewAgentSkillsInitializer(
    skillName, targetDir string,
    tm TemplateManager,
) *AgentSkillsInitializer
```

#### Scenario: AgentSkillsInitializer construction

- **WHEN** an AgentSkillsInitializer is created via
  `NewAgentSkillsInitializer(skillName, targetDir, tm)`
- **THEN** it SHALL receive a skill name matching an embedded skill directory
- **AND** it SHALL receive a target directory path for skill installation
- **AND** it SHALL receive a TemplateManager to access embedded skill files
- **AND** the initializer SHALL use `projectFs` for all file operations

#### Scenario: Copy skill directory

- **WHEN** `Init()` is called on an AgentSkillsInitializer
- **THEN** it SHALL recursively copy all files from the embedded skill
  directory
- **AND** it SHALL preserve the directory structure (e.g., `scripts/accept.sh`)
- **AND** it SHALL preserve file permissions (executable scripts remain
  executable)
- **AND** it SHALL create the target directory if it does not exist

#### Scenario: Skill not found

- **WHEN** `Init()` is called with a skill name that does not exist
- **THEN** it SHALL return an error indicating the skill was not found
- **AND** the error message SHALL include the requested skill name

#### Scenario: Idempotent execution

- **WHEN** `Init()` is called multiple times
- **THEN** existing files SHALL be overwritten with embedded content
- **AND** the result SHALL be equivalent to a fresh installation

#### Scenario: IsSetup check for skills

- **WHEN** `IsSetup()` is called on an AgentSkillsInitializer
- **THEN** it SHALL return `true` if the `SKILL.md` file exists in the target
  directory
- **AND** it SHALL return `false` if `SKILL.md` is missing

#### Scenario: Deduplication key for skills

- **WHEN** `dedupeKey()` is called on an AgentSkillsInitializer
- **THEN** it SHALL return `AgentSkillsInitializer:<targetDir>`
- **AND** the path SHALL be normalized with `filepath.Clean`

### Requirement: TemplateManager Skill Access

The TemplateManager SHALL provide access to embedded skill directories.

```go
// SkillFS returns an fs.FS rooted at the skill directory for the given skill
// name. Returns an error if the skill does not exist.
func (tm *TemplateManager) SkillFS(skillName string) (fs.FS, error)
```

#### Scenario: Retrieve skill filesystem

- **WHEN** `SkillFS(skillName)` is called with a valid skill name
- **THEN** it SHALL return an `fs.FS` interface for the skill directory
- **AND** the filesystem SHALL contain all files under
  `templates/skills/<skillName>/`
- **AND** file paths SHALL be relative to the skill root (e.g., `SKILL.md`,
  `scripts/accept.sh`)

#### Scenario: Skill not found error

- **WHEN** `SkillFS(skillName)` is called with an unknown skill name
- **THEN** it SHALL return an error
- **AND** the error message SHALL indicate the skill was not found

### Requirement: Embedded Skill Templates

The system SHALL embed skill directories under
`internal/initialize/templates/skills/`.

#### Scenario: Skill directory structure

- **WHEN** a skill is embedded
- **THEN** it SHALL be located at `templates/skills/<skill-name>/`
- **AND** it SHALL contain at minimum a `SKILL.md` file
- **AND** optional directories (`scripts/`, `references/`, `assets/`) MAY be
  included

#### Scenario: Embedded skill globbing

- **WHEN** templates are parsed
- **THEN** the embed directive SHALL use `//go:embed templates/skills/**/*`
- **AND** all skill files SHALL be available via `SkillFS()`
