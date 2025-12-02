## ADDED Requirements

### Requirement: Template Path Variables

The template rendering system SHALL support dynamic path variables for all directory and file references, allowing template content to be decoupled from specific path names while maintaining backward-compatible defaults.

The `TemplateContext` struct SHALL provide the following fields with default values:
- `BaseDir`: The root Spectr directory name (default: `spectr`)
- `SpecsDir`: The specifications directory path (default: `spectr/specs`)
- `ChangesDir`: The changes directory path (default: `spectr/changes`)
- `ProjectFile`: The project configuration file path (default: `spectr/project.md`)
- `AgentsFile`: The agents instruction file path (default: `spectr/AGENTS.md`)

#### Scenario: Templates use path variables instead of hardcoded strings

- **WHEN** a template file contains path references
- **THEN** the path SHALL be expressed using Go template syntax (e.g., `{{ .BaseDir }}`, `{{ .SpecsDir }}`)
- **AND** hardcoded `spectr/` strings SHALL NOT appear in template files for path references
- **AND** the rendered output SHALL contain the actual path values from the context

#### Scenario: Default context produces backward-compatible output

- **WHEN** `DefaultTemplateContext()` is used for rendering
- **THEN** the rendered output SHALL be identical to the previous hardcoded output
- **AND** all path references SHALL resolve to `spectr/`, `spectr/specs/`, `spectr/changes/`, etc.

#### Scenario: Template manager methods accept context parameter

- **WHEN** `RenderAgents()`, `RenderInstructionPointer()`, or `RenderSlashCommand()` is called
- **THEN** the method SHALL accept a `TemplateContext` parameter
- **AND** the context values SHALL be available within the template

#### Scenario: All template files use consistent variable names

- **WHEN** any template file references a Spectr path
- **THEN** it SHALL use the standardized variable names (`BaseDir`, `SpecsDir`, `ChangesDir`, `ProjectFile`, `AgentsFile`)
- **AND** variable names SHALL be consistent across all template files
