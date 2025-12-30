# Change: Add Template Path Variables for Dynamic Directory Names

## Why

The templates in `internal/initialize/templates/` hardcode the `spectr/`
directory path (e.g., `spectr/specs/`, `spectr/changes/`, `spectr/AGENTS.md`).
This prevents users from customizing the base directory name and couples
templates tightly to a specific naming convention. Introducing template
variables allows for future configurability and cleaner separation between
template content and path structure.

## What Changes

- Add a `TemplateContext` struct with path-related fields: `BaseDir`,
  `SpecsDir`, `ChangesDir`, `ProjectFile`, `AgentsFile`
- Update all template files (`.tmpl`) to use Go template variables (e.g., `{{
  .BaseDir }}`, `{{ .ChangesDir }}`) instead of hardcoded `spectr/` paths
- Update `TemplateManager` methods to accept and pass the context to templates
- Default values maintain backward compatibility (`spectr`, `spectr/specs`,
  `spectr/changes`, etc.)

## Impact

- Affected specs: `cli-interface` (Initialization and template rendering)
- Affected code:
  - `internal/initialize/models.go` - Add `TemplateContext` struct
  - `internal/initialize/templates.go` - Update render methods to accept context
  - `internal/initialize/templates/*.tmpl` - Replace hardcoded paths with
    variables
  - `internal/initialize/executor.go` - Pass default context to template
    rendering
