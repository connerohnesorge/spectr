# Provider System Specification Delta

## MODIFIED Requirements

### Requirement: Domain Package

The system SHALL define a `internal/domain` package containing shared domain types, including support for provider-specific template composition.

#### Scenario: TemplateRef supports composition

- **WHEN** code creates a `domain.TemplateRef`
- **THEN** it SHALL have fields: `Name`, `Template`, and `ProviderTemplate`
- **AND** `ProviderTemplate` SHALL be optional (nil if no provider customization)
- **AND** the `Render(ctx *TemplateContext)` method SHALL compose base + provider templates before rendering
- **AND** if `ProviderTemplate` is nil, rendering SHALL proceed without composition

#### Scenario: Template section composition

- **WHEN** `TemplateRef.Render()` is called with non-nil `ProviderTemplate`
- **THEN** it SHALL clone the base template
- **AND** SHALL merge provider template sections using Go's `AddParseTree()`
- **AND** SHALL execute the "main" template from the composed result
- **AND** provider sections SHALL override base sections (last-wins semantics)

#### Scenario: Backward compatible rendering

- **WHEN** `TemplateRef.Render()` is called with nil `ProviderTemplate`
- **THEN** behavior SHALL be identical to previous implementation
- **AND** template SHALL render correctly with all context variables

## ADDED Requirements

### Requirement: Provider Template Directory Structure

The system SHALL organize provider-specific template overrides in a structured directory hierarchy.

#### Scenario: Provider template discovery

- **WHEN** `TemplateManager` initializes
- **THEN** it SHALL check for `templates/providers/` directory in embedded filesystem
- **AND** SHALL iterate over provider subdirectories
- **AND** SHALL parse template files for each provider

#### Scenario: Missing provider directory

- **WHEN** `templates/providers/` directory does not exist
- **THEN** initialization SHALL succeed with empty provider templates
- **AND** all providers SHALL use generic templates (backward compatibility)

#### Scenario: Provider template validation

- **WHEN** provider templates are parsed during initialization
- **THEN** the system SHALL validate that provider templates only define known sections
- **AND** if unknown sections are found, initialization SHALL fail with a descriptive error
- **AND** known sections are: `guardrails`, `steps`, `reference`, `main`, `base_guardrails`, `base_steps`, `base_reference`

### Requirement: Provider-Aware Template Resolution

The `TemplateManager` SHALL provide methods to resolve templates with provider-specific overrides.

#### Scenario: ProviderSlashCommand method

- **WHEN** code calls `tm.ProviderSlashCommand(providerID, cmd)`
- **THEN** it SHALL return a `domain.TemplateRef` with:
  - `Name` set to the slash command template name (e.g., "slash-proposal.md.tmpl")
  - `Template` set to the base template
  - `ProviderTemplate` set to provider's template (nil if provider not found)

#### Scenario: ProviderTOMLSlashCommand method

- **WHEN** code calls `tm.ProviderTOMLSlashCommand(providerID, cmd)`
- **THEN** it SHALL return a `domain.TemplateRef` with:
  - `Name` set to the TOML template name (e.g., "slash-proposal.toml.tmpl")
  - `Template` set to the base TOML template
  - `ProviderTemplate` set to provider's TOML template (nil if not found)

#### Scenario: Unknown provider ID

- **WHEN** `ProviderSlashCommand("unknown-provider", cmd)` is called
- **THEN** the method SHALL return a `domain.TemplateRef` with nil `ProviderTemplate`
- **AND** rendering SHALL use only the base template (automatic fallback)

### Requirement: Backward Compatible Template Access

Existing `TemplateManager` methods SHALL continue working without changes.

#### Scenario: SlashCommand method unchanged

- **WHEN** code calls `tm.SlashCommand(cmd)` (existing method)
- **THEN** it SHALL internally call `tm.ProviderSlashCommand("", cmd)`
- **AND** SHALL return a `domain.TemplateRef` with nil `ProviderTemplate`
- **AND** behavior SHALL be identical to before this change

#### Scenario: TOMLSlashCommand method unchanged

- **WHEN** code calls `tm.TOMLSlashCommand(cmd)` (existing method)
- **THEN** it SHALL internally call `tm.ProviderTOMLSlashCommand("", cmd)`
- **AND** SHALL return a `domain.TemplateRef` with nil `ProviderTemplate`

### Requirement: Template Section Definitions

Generic templates SHALL define named sections for composition.

#### Scenario: Markdown slash command sections

- **WHEN** parsing generic markdown templates (slash-proposal.md.tmpl, slash-apply.md.tmpl)
- **THEN** the template SHALL define sections: `guardrails`, `steps`, `reference`, `main`
- **AND** each section SHALL use `{{define "name"}}...{{end}}` syntax
- **AND** the `main` section SHALL reference other sections via `{{template "name" .}}`

#### Scenario: Base sections for provider reuse

- **WHEN** a provider template needs to include base content
- **THEN** the provider SHALL call `{{template "base_guardrails" .}}` or similar
- **AND** base sections (with `base_` prefix) SHALL contain the original generic content
- **AND** this allows providers to extend rather than completely replace sections

#### Scenario: Provider section override

- **WHEN** a provider defines a section (e.g., `{{define "guardrails"}}`)
- **THEN** composition SHALL merge the provider section with the base template
- **AND** provider sections SHALL have access to the same `TemplateContext` variables
- **AND** other sections from base template SHALL remain available if not overridden

### Requirement: TemplateManager Interface Extension

The `TemplateManager` interface SHALL be extended with provider-aware methods.

#### Scenario: New interface methods

- **WHEN** the interface is defined
- **THEN** it SHALL include:
  - `ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
  - `ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
- **AND** existing methods SHALL remain: `InstructionPointer()`, `Agents()`, `SlashCommand()`, `TOMLSlashCommand()`, `SkillFS()`

#### Scenario: Concrete implementation

- **WHEN** `internal/initialize.TemplateManager` implements the interface
- **THEN** it SHALL implement all new methods with provider template resolution
- **AND** all existing methods SHALL work identically to before

### Requirement: Provider Opt-In

Providers SHALL opt-in to using custom templates by calling new methods.

#### Scenario: Provider using custom templates

- **WHEN** a provider's `Initializers()` calls `tm.ProviderSlashCommand(providerID, cmd)`
- **THEN** provider-specific templates SHALL be used if they exist
- **AND** automatic fallback to generic templates if provider-specific templates don't exist
- **AND** no error handling required (fallback is transparent)

#### Scenario: Provider not migrated

- **WHEN** a provider calls `tm.SlashCommand(cmd)` (old method)
- **THEN** generic templates SHALL be used (nil `ProviderTemplate`)
- **AND** behavior SHALL be identical to current implementation

#### Scenario: Gradual migration

- **WHEN** some providers use new methods and others use old methods
- **THEN** all providers SHALL initialize successfully
- **AND** each SHALL receive correct templates based on method called
- **AND** no coordination between providers is required

## Cross-References

This capability:
- **Extends**: Domain Package (adds composition support to TemplateRef)
- **Extends**: CLI Interface (template-based commands now support provider customization)
- **Used by**: support-claude-code provider (claude-code custom templates)
- **Used by**: support-codex provider (codex custom templates)
- **Optional for**: All other providers (can opt-in at any time)
