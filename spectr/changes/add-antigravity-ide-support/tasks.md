# Implementation Tasks: Add Antigravity IDE Support

## 1. Update Tool Registry
- [ ] 1.1 Add Antigravity tool definition to `NewRegistry()` in `internal/init/registry.go`
  - Priority: 7 (after Qwen)
  - ID: `antigravity`
  - Name: `Antigravity`
  - Type: `ToolTypeConfig`
  - ConfigPath: `.antigravity/config.json`
  - Configured: false
- [ ] 1.2 Add Antigravity-to-slash mapping in `configToSlashMapping` variable
  - Key: `antigravity`
  - Value: `antigravity-slash` (following pattern of other tools)

## 2. Implement AntigravityConfigurator
- [ ] 2.1 Create `AntigravityConfigurator` struct in `internal/init/configurator.go`
  - Follow exact pattern of existing 6 configurators (ClaudeCodeConfigurator, ClineConfigurator, etc.)
- [ ] 2.2 Implement `Configure()` method
  - Renders AGENTS.md template using `NewTemplateManager()`
  - Updates `AGENTS.md` file using `UpdateFileWithMarkers()`
  - Uses `SpectrStartMarker` and `SpectrEndMarker` markers
  - Creates `.agent/workflows/spectr-*.md` slash command files
- [ ] 2.3 Implement `IsConfigured()` method
  - Checks for `AGENTS.md` file in project path
- [ ] 2.4 Implement `GetName()` method
  - Returns `"Antigravity"`

## 3. Register Configurator
- [ ] 3.1 Add Antigravity configurator to wizard/executor
  - No changes needed to `internal/init/wizard.go` (auto-populated from registry)
  - No changes needed to `internal/init/executor.go` (existing dispatch logic handles new configurator)

## 4. Testing
- [ ] 4.1 Add test for Antigravity tool registration in `internal/init/registry_test.go`
  - Verify tool is found with `GetTool("antigravity")`
  - Verify mapping exists in `configToSlashMapping`
- [ ] 4.2 Add test for AntigravityConfigurator in `internal/init/configurator_test.go`
  - Test `Configure()` creates `AGENTS.md` with markers
  - Test `IsConfigured()` returns true when file exists
  - Test `GetName()` returns correct name
  - Test file update behavior when markers already exist

## 5. Validation
- [ ] 5.1 Run `spectr validate add-antigravity-ide-support --strict`
  - All validation rules pass
  - No missing requirements or scenarios
  - Spec delta syntax is correct
- [ ] 5.2 Manual testing
  - Run `spectr init` and verify Antigravity appears in tool list
  - Select Antigravity and verify files are created correctly
  - Verify `ANTIGRAVITY.md` contains spectr instructions between markers
  - Verify slash commands are auto-installed in `.antigravity/commands/spectr/`

## Notes
- Antigravity becomes the 7th config-based tool in the registry
- All tests must pass before archiving
- Follow existing code patterns for consistency
- No breaking changes to public APIs
