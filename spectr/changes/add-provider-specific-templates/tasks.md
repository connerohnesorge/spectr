# 1. Template Resolution Infrastructure

## Tasks

- [ ] 1.1 Update `TemplateRenderer` interface to accept provider ID
  parameter
- [ ] 1.2 Implement provider-first template lookup in `TemplateManager`
- [ ] 1.3 Add helper function to resolve template path with fallback
  logic
- [ ] 1.4 Write unit tests for template resolution with/without provider
  overrides

## 2. Provider Integration

- [ ] 2.1 Update `BaseProvider.Configure()` to pass provider ID to template
  renderer
- [ ] 2.2 Update all `TemplateRenderer` method signatures
- [ ] 2.3 Ensure backward compatibility (empty provider ID = generic)
- [ ] 2.4 Write integration tests for provider-specific template
  rendering

## 3. Initial Provider Templates

- [ ] 3.1 Create `templates/claude-code/` directory structure
- [ ] 3.2 Create Claude Code-specific `AGENTS.md.tmpl` with tool-specific
  references
- [ ] 3.3 Create `templates/crush/` directory structure  
- [ ] 3.4 Create Crush-specific `AGENTS.md.tmpl` with agent delegation
  patterns
- [ ] 3.5 Verify `spectr init` produces provider-appropriate output

## 4. Validation & Documentation

- [ ] 4.1 Run full test suite (`go test ./...`)
- [ ] 4.2 Test `spectr init` with claude-code and crush providers
- [ ] 4.3 Verify fallback works for providers without custom
  templates
- [ ] 4.4 Update provider documentation comments if needed
