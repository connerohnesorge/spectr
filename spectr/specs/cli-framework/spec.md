# Cli Framework Specification

## Requirements

### Requirement: Provider-Specific Template Resolution

The template system SHALL support per-provider template overrides with
fallback to generic templates.

#### Scenario: Provider with custom AGENTS.md template

- **WHEN** rendering AGENTS.md for a provider with
  `templates/{provider-id}/AGENTS.md.tmpl`
- **THEN** the system uses the provider-specific template
- **AND** the output reflects provider-specific tool names and
  patterns

#### Scenario: Provider without custom template falls back to generic

- **WHEN** rendering AGENTS.md for a provider without a custom template
- **THEN** the system uses `templates/spectr/AGENTS.md.tmpl` (generic)
- **AND** the provider still receives valid instructions

#### Scenario: Partial override (some templates custom, some generic)

- **WHEN** a provider has `templates/{provider-id}/AGENTS.md.tmpl` but no
  `slash-proposal.md.tmpl`
- **THEN** AGENTS.md uses the custom template
- **AND** slash-proposal uses the generic
  `templates/tools/slash-proposal.md.tmpl`

### Requirement: Provider-Specific Template Directory Structure

The template system SHALL organize provider templates in dedicated
directories.

#### Scenario: Template directory layout

- **WHEN** the template system initializes
- **THEN** it recognizes the structure:

  ```text
  templates/
  ├── spectr/           # Generic templates (fallback)
  │   ├── AGENTS.md.tmpl
  │   ├── instruction-pointer.md.tmpl
  │   └── project.md.tmpl
  ├── tools/            # Generic slash commands (fallback)
  │   ├── slash-proposal.md.tmpl
  │   └── slash-apply.md.tmpl
  ├── claude-code/      # Provider-specific overrides
  │   ├── AGENTS.md.tmpl
  │   └── slash-proposal.md.tmpl
  ├── crush/
  │   └── AGENTS.md.tmpl
  └── ci/
      └── spectr-ci.yml.tmpl
  ```

#### Scenario: Provider ID mapping to directory

- **WHEN** resolving templates for provider `claude-code`
- **THEN** the system looks in `templates/claude-code/` first
- **AND** falls back to generic directories if not found

### Requirement: Template Renderer Provider Context

The TemplateRenderer interface SHALL accept provider context for resolution.

#### Scenario: Render with provider ID

- **WHEN** a provider calls `RenderAgents(ctx, providerID)`
- **THEN** the renderer uses the provider ID to select the appropriate
  template
- **AND** the template context variables are still populated correctly

#### Scenario: Backward compatibility for generic rendering

- **WHEN** `RenderAgents(ctx, "")` is called with empty provider ID
- **THEN** the system uses the generic template
- **AND** existing behavior is preserved

