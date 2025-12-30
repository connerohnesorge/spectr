# Implementation Tasks

## 1. Define Template Context

- [ ] 1.1 Add `TemplateContext` struct to `internal/initialize/models.go` with
  fields:
  - `BaseDir` (default: `spectr`)
  - `SpecsDir` (default: `spectr/specs`)
  - `ChangesDir` (default: `spectr/changes`)
  - `ProjectFile` (default: `spectr/project.md`)
  - `AgentsFile` (default: `spectr/AGENTS.md`)
- [ ] 1.2 Add `DefaultTemplateContext()` function returning default values

## 2. Update Template Manager

- [ ] 2.1 Update `RenderAgents()` to accept `TemplateContext` parameter
- [ ] 2.2 Update `RenderInstructionPointer()` to accept `TemplateContext`
  parameter
- [ ] 2.3 Update `RenderSlashCommand()` to accept `TemplateContext` parameter
- [ ] 2.4 Update all render method signatures and callers in `executor.go`

## 3. Update Template Files

- [ ] 3.1 Update `templates/spectr/AGENTS.md.tmpl` - replace `spectr/` paths
  with variables
- [ ] 3.2 Update `templates/spectr/instruction-pointer.md.tmpl` - replace paths
  with variables
- [ ] 3.3 Update `templates/tools/slash-proposal.md.tmpl` - replace paths with
  variables
- [ ] 3.4 Update `templates/tools/slash-apply.md.tmpl` - replace paths with
  variables
- [ ] 3.5 Update `templates/tools/slash-sync.md.tmpl` - replace paths with
  variables

## 4. Validation

- [ ] 4.1 Run existing tests to ensure backward compatibility
- [ ] 4.2 Verify `spectr init` produces identical output with default context
- [ ] 4.3 Run `go build` to confirm compilation
