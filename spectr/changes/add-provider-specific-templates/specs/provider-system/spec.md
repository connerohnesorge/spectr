# Provider System Specification Delta

## MODIFIED Requirements

### Requirement: Domain Package

The system SHALL define a `internal/domain` package containing shared domain types to break import cycles, including support for provider-specific template composition.

#### Scenario: TemplateRef in domain package

- **WHEN** code needs to reference a template
- **THEN** it SHALL use `domain.TemplateRef` from `internal/domain`
- **AND** `TemplateRef` SHALL have the following structure:

```go
type TemplateRef struct {
    Name             string             // template file name (e.g., "slash-proposal.md.tmpl")
    Template         *template.Template // base template with all sections
    ProviderTemplate *template.Template // provider overrides (nil if none)
}
```

- **AND** `TemplateRef` SHALL support template composition when `ProviderTemplate` is non-nil
- **AND** the `Render(ctx *TemplateContext)` method SHALL compose base + provider templates before rendering
- **AND** composition SHALL use Go's `template.Clone()` and `AddParseTree()` for last-wins semantics

#### Scenario: TemplateRef with no provider overrides

- **WHEN** `TemplateRef.ProviderTemplate` is nil
- **THEN** `Render()` SHALL execute the base template directly without composition
- **AND** behavior SHALL be identical to the previous implementation

#### Scenario: TemplateRef with provider overrides

- **WHEN** `TemplateRef.ProviderTemplate` is non-nil
- **THEN** `Render()` SHALL clone the base template
- **AND** SHALL merge provider template sections using `AddParseTree()`
- **AND** SHALL execute the "main" template from the composed result
- **AND** provider sections SHALL override base sections with the same name (last-wins)

#### Scenario: Template composition error handling

- **WHEN** template composition fails during `AddParseTree()`
- **THEN** `Render()` SHALL return an error with context about which template failed to merge
- **AND** the error message SHALL include the provider template name for debugging

## ADDED Requirements

### Requirement: Provider Template Directory Structure

The system SHALL support provider-specific template overrides organized in subdirectories.

#### Scenario: Provider template discovery

- **WHEN** `TemplateManager` initializes
- **THEN** it SHALL check for `templates/providers/` directory in the embedded filesystem
- **AND** SHALL iterate over subdirectories representing provider IDs
- **AND** SHALL parse template files (*.tmpl) for each provider into separate `template.Template` instances
- **AND** SHALL store provider templates in a map: `providerTemplates[providerID] = template.Template`

#### Scenario: Missing providers directory

- **WHEN** `templates/providers/` directory does not exist
- **THEN** initialization SHALL succeed with empty `providerTemplates` map
- **AND** all providers SHALL use generic templates (backward compatibility)

#### Scenario: Provider directory with no templates

- **WHEN** a provider subdirectory exists but contains no .tmpl files
- **THEN** that provider SHALL be skipped (not added to `providerTemplates` map)
- **AND** initialization SHALL continue without error

### Requirement: Provider-Aware Template Resolution

The `TemplateManager` SHALL provide methods to resolve templates with provider-specific overrides.

#### Scenario: ProviderSlashCommand method

- **WHEN** code calls `tm.ProviderSlashCommand(providerID, cmd)`
- **THEN** the method SHALL return a `domain.TemplateRef` with:
  - `Name` set to the slash command template name (e.g., "slash-proposal.md.tmpl")
  - `Template` set to the base template from `tm.templates`
  - `ProviderTemplate` set to `tm.providerTemplates[providerID]` (nil if provider not found)

#### Scenario: ProviderTOMLSlashCommand method

- **WHEN** code calls `tm.ProviderTOMLSlashCommand(providerID, cmd)`
- **THEN** the method SHALL return a `domain.TemplateRef` with:
  - `Name` set to the TOML slash command template name (e.g., "slash-proposal.toml.tmpl")
  - `Template` set to the base template
  - `ProviderTemplate` set to the provider's TOML template (nil if not found)

#### Scenario: Unknown provider ID

- **WHEN** `ProviderSlashCommand("unknown-provider", cmd)` is called
- **THEN** the method SHALL return a `domain.TemplateRef` with nil `ProviderTemplate`
- **AND** rendering SHALL fall back to the generic base template

### Requirement: Backward Compatible Template Access

The `TemplateManager` SHALL maintain existing methods for backward compatibility.

#### Scenario: SlashCommand method unchanged

- **WHEN** code calls `tm.SlashCommand(cmd)` (existing method)
- **THEN** the method SHALL internally call `tm.ProviderSlashCommand("", cmd)`
- **AND** SHALL return a `domain.TemplateRef` with nil `ProviderTemplate`
- **AND** behavior SHALL be identical to previous implementation (uses generic template)

#### Scenario: TOMLSlashCommand method unchanged

- **WHEN** code calls `tm.TOMLSlashCommand(cmd)` (existing method)
- **THEN** the method SHALL internally call `tm.ProviderTOMLSlashCommand("", cmd)`
- **AND** SHALL return a `domain.TemplateRef` with nil `ProviderTemplate`
- **AND** behavior SHALL be identical to previous implementation

### Requirement: Template Section Definitions

Generic templates SHALL define named sections using Go's `{{define}}` directive.

#### Scenario: Markdown slash command sections

- **WHEN** parsing `slash-proposal.md.tmpl` or `slash-apply.md.tmpl`
- **THEN** the template SHALL define the following sections:
  - `{{define "guardrails"}}` - Best practices and constraints
  - `{{define "steps"}}` - Ordered workflow instructions
  - `{{define "reference"}}` - Supporting documentation pointers
  - `{{define "main"}}` - Composition entry point
- **AND** the "main" section SHALL reference other sections via `{{template "guardrails" .}}`, etc.

#### Scenario: Provider template section override

- **WHEN** a provider template defines a section (e.g., `{{define "guardrails"}}`)
- **THEN** composition SHALL replace the base section with the provider section
- **AND** the provider section SHALL have access to the same `TemplateContext` variables
- **AND** undefined sections SHALL fall back to base definitions

### Requirement: TemplateManager Interface Extension

The `TemplateManager` interface in `internal/initialize/providers/provider.go` SHALL be extended with provider-aware methods.

#### Scenario: New interface methods

- **WHEN** the `TemplateManager` interface is defined
- **THEN** it SHALL include the following new methods:
  - `ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
  - `ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
- **AND** existing methods SHALL remain unchanged:
  - `InstructionPointer() domain.TemplateRef`
  - `Agents() domain.TemplateRef`
  - `SlashCommand(cmd domain.SlashCommand) domain.TemplateRef`
  - `TOMLSlashCommand(cmd domain.SlashCommand) domain.TemplateRef`
  - `SkillFS(skillName string) (fs.FS, error)`

#### Scenario: Interface implementation by concrete TemplateManager

- **WHEN** `internal/initialize.TemplateManager` implements the interface
- **THEN** it SHALL implement all new methods with provider template resolution
- **AND** it SHALL maintain backward compatibility for all existing methods

### Requirement: Provider Opt-In

Providers SHALL opt-in to using custom templates by updating their `Initializers()` method.

#### Scenario: Provider using custom templates

- **WHEN** a provider's `Initializers()` method calls `tm.ProviderSlashCommand(providerID, cmd)`
- **THEN** the provider SHALL receive provider-specific templates if they exist
- **AND** SHALL automatically fall back to generic templates if provider-specific templates don't exist
- **AND** no error handling is required (fallback is transparent)

#### Scenario: Provider not migrated

- **WHEN** a provider continues using `tm.SlashCommand(cmd)` (old method)
- **THEN** the provider SHALL receive generic templates (nil `ProviderTemplate`)
- **AND** behavior SHALL be identical to before this change

#### Scenario: Gradual migration

- **WHEN** some providers use `ProviderSlashCommand()` and others use `SlashCommand()`
- **THEN** all providers SHALL initialize successfully
- **AND** each provider SHALL receive the correct templates based on the method called
- **AND** no coordination between providers is required

## Cross-References

This capability builds on:
- **Domain Package** (existing) - Extends `TemplateRef` with composition support
- **Initializer Interface** (existing) - No changes to initializer pattern
- **Provider Interface** (existing) - No changes to `Provider` interface (providers opt-in via method calls)

This capability is used by:
- **support-claude-code** - First provider to use `ProviderSlashCommand("claude-code", ...)`
- **support-codex** - Second provider to use `ProviderSlashCommand("codex", ...)`
- **support-gemini** - Future use of `ProviderTOMLSlashCommand("gemini", ...)` (placeholder)
