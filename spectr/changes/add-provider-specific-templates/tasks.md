# Implementation Tasks: Provider-Specific Templates

## Phase 1: Foundation (No Breaking Changes)

### 1.1 Refactor Generic Templates with {{define}} Blocks

- [ ] Read current `internal/domain/templates/slash-proposal.md.tmpl`
- [ ] Wrap content in `{{define "guardrails"}}...{{end}}` block
- [ ] Wrap content in `{{define "steps"}}...{{end}}` block
- [ ] Wrap content in `{{define "reference"}}...{{end}}` block
- [ ] Add `{{define "main"}}` block that calls `{{template "guardrails" .}}`, `{{template "steps" .}}`, `{{template "reference" .}}`
- [ ] Test that template still renders correctly with existing tests
- [ ] Read current `internal/domain/templates/slash-apply.md.tmpl`
- [ ] Apply same `{{define}}` block structure to slash-apply.md.tmpl
- [ ] Test that template still renders correctly

### 1.2 Update TemplateRef for Composition

- [ ] Read `internal/domain/template.go`
- [ ] Add `ProviderTemplate *template.Template` field to `TemplateRef` struct
- [ ] Update `Render()` method to check if `ProviderTemplate` is nil
- [ ] If nil, render base template directly (existing behavior)
- [ ] If non-nil, call new `composeTemplate()` helper
- [ ] Implement `composeTemplate()` private method using `template.Clone()` and `AddParseTree()`
- [ ] Execute "main" template from composed result
- [ ] Add error handling for composition failures

### 1.3 Add Unit Tests for Composition

- [ ] Create `internal/domain/template_test.go` (if it doesn't exist)
- [ ] Add test: `TestTemplateRef_Render_NoProviderTemplate` (nil ProviderTemplate uses base)
- [ ] Add test: `TestTemplateRef_Render_WithProviderOverride` (provider section replaces base)
- [ ] Add test: `TestTemplateRef_composeTemplate` (composition merges correctly)
- [ ] Add test: Verify base sections used when provider doesn't override
- [ ] Add test: Verify last-wins semantics (provider overrides base)
- [ ] Run tests: `go test ./internal/domain/...`
- [ ] Fix any failing tests

## Phase 2: Template Manager Infrastructure

### 2.1 Add Provider Template Parsing

- [ ] Read `internal/initialize/templates.go`
- [ ] Add `providerTemplates map[string]*template.Template` field to `TemplateManager` struct
- [ ] In `NewTemplateManager()`, check if `templates/providers/` directory exists
- [ ] If exists, iterate over subdirectories with `fs.ReadDir(templateFS, "templates/providers")`
- [ ] For each provider directory, parse templates with pattern `templates/providers/{id}/*.tmpl`
- [ ] Store in map: `providerTemplates[providerID] = parsed template`
- [ ] Handle errors gracefully (skip provider if parsing fails, don't fail initialization)
- [ ] Return `TemplateManager` with populated `providerTemplates` map

### 2.2 Implement Provider-Aware Methods

- [ ] In `internal/initialize/templates.go`, implement `ProviderSlashCommand(providerID string, cmd domain.SlashCommand)`
- [ ] Method returns `domain.TemplateRef` with base template + provider template from map
- [ ] If `providerTemplates[providerID]` is nil, `ProviderTemplate` field is nil (automatic fallback)
- [ ] Implement `ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand)` with same pattern
- [ ] Update existing `SlashCommand()` to call `ProviderSlashCommand("", cmd)` internally
- [ ] Update existing `TOMLSlashCommand()` to call `ProviderTOMLSlashCommand("", cmd)` internally
- [ ] Verify backward compatibility (old methods work identically)

### 2.3 Update TemplateManager Interface

- [ ] Read `internal/initialize/providers/provider.go`
- [ ] Add `ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef` to interface
- [ ] Add `ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef` to interface
- [ ] Verify existing methods remain in interface (backward compatibility)
- [ ] Verify interface is satisfied by concrete `TemplateManager` (compile check)

### 2.4 Add Integration Tests

- [ ] Create or update `internal/initialize/templates_test.go`
- [ ] Add test: `TestTemplateManager_ParseProviderTemplates` (discovers and parses provider directories)
- [ ] Add test: `TestTemplateManager_ProviderSlashCommand` (returns correct TemplateRef)
- [ ] Add test: `TestTemplateManager_ProviderSlashCommand_UnknownProvider` (returns nil ProviderTemplate)
- [ ] Add test: `TestTemplateManager_BackwardCompatibility` (old methods still work)
- [ ] Add test: Missing providers directory doesn't break initialization
- [ ] Run tests: `go test ./internal/initialize/...`
- [ ] Fix any failing tests

## Phase 3: Provider-Specific Templates

### 3.1 Create Claude Code Templates

- [ ] Create directory: `internal/initialize/templates/providers/claude-code/`
- [ ] Create `slash-proposal.md.tmpl` with `{{define "guardrails"}}` override
- [ ] Include: "Use the orchestrator pattern: delegate to coder/tester/stuck subagents"
- [ ] Include: "Leverage Claude's 200k context window for comprehensive proposals"
- [ ] Include: "Reference skills: Check `.claude/skills/` for available AgentSkills"
- [ ] Include generic guardrails content (don't remove, add to them)
- [ ] Create `slash-apply.md.tmpl` with `{{define "guardrails"}}` override
- [ ] Include: "Delegate implementation: Use `coder` subagent for tasks"
- [ ] Include: "Verify subagent work: Run `git diff --stat` after coder completes"
- [ ] Include: "Test every change: Use `tester` subagent after EVERY coder completion"
- [ ] Include generic guardrails content

### 3.2 Create Codex Templates

- [ ] Create directory: `internal/initialize/templates/providers/codex/`
- [ ] Create `slash-proposal.md.tmpl` with `{{define "steps"}}` override
- [ ] Add step 2: "Check global prompts: Review `~/.codex/prompts/` for existing Spectr-related prompts"
- [ ] Add step 7: "Reference AgentSkills: If applicable, reference skills in `~/.codex/skills/`"
- [ ] Include all other generic steps (don't remove existing steps)
- [ ] Create `slash-apply.md.tmpl` with `{{define "reference"}}` override
- [ ] Add reference: "Check `~/.codex/prompts/` for global Spectr prompts"
- [ ] Add reference: "Check `~/.codex/skills/` for available AgentSkills"
- [ ] Include generic reference content

### 3.3 Create Gemini Placeholder

- [ ] Create directory: `internal/initialize/templates/providers/gemini/`
- [ ] Create empty `.gitkeep` file (placeholder for future TOML templates)
- [ ] Document in README or comment: "TOML-specific templates for Gemini provider (future work)"

### 3.4 Test Template Composition

- [ ] Parse claude-code templates manually in test
- [ ] Compose with base template
- [ ] Verify "guardrails" section contains "orchestrator pattern"
- [ ] Verify "steps" section contains generic steps (not overridden)
- [ ] Parse codex templates manually in test
- [ ] Compose with base template
- [ ] Verify "steps" section contains "~/.codex/prompts/"
- [ ] Verify "guardrails" section contains generic content (not overridden)

## Phase 4: Provider Migration

### 4.1 Update ClaudeProvider

- [ ] Read `internal/initialize/providers/claude.go`
- [ ] Find `Initializers()` method
- [ ] Locate `NewSlashCommandsInitializer` call
- [ ] Change `tm.SlashCommand(domain.SlashProposal)` to `tm.ProviderSlashCommand("claude-code", domain.SlashProposal)`
- [ ] Change `tm.SlashCommand(domain.SlashApply)` to `tm.ProviderSlashCommand("claude-code", domain.SlashApply)`
- [ ] Leave all other initializers unchanged
- [ ] Verify code compiles: `go build ./...`

### 4.2 Update CodexProvider

- [ ] Read `internal/initialize/providers/codex.go`
- [ ] Find `Initializers()` method
- [ ] Locate `NewHomeSlashCommandsInitializer` call
- [ ] Change `tm.SlashCommand(domain.SlashProposal)` to `tm.ProviderSlashCommand("codex", domain.SlashProposal)`
- [ ] Change `tm.SlashCommand(domain.SlashApply)` to `tm.ProviderSlashCommand("codex", domain.SlashApply)`
- [ ] Leave all other initializers unchanged
- [ ] Verify code compiles: `go build ./...`

### 4.3 Verify Other Providers Unchanged

- [ ] Read `internal/initialize/providers/aider.go` (or any other provider)
- [ ] Verify it still uses `tm.SlashCommand(cmd)` (old method)
- [ ] Verify code compiles
- [ ] Confirm no changes needed for non-migrated providers

## Phase 5: Validation & Testing

### 5.1 End-to-End Validation - Claude Code

- [ ] Build project: `go build -o spectr ./cmd/spectr`
- [ ] Create temp directory for testing
- [ ] Run: `./spectr init --provider=claude-code` in temp directory
- [ ] Verify `.claude/commands/spectr/proposal.md` exists
- [ ] Read file and verify it contains "orchestrator pattern"
- [ ] Read file and verify it contains "coder/tester/stuck subagents"
- [ ] Read file and verify it contains "200k context window"
- [ ] Verify `.claude/commands/spectr/apply.md` exists
- [ ] Read file and verify it contains "Delegate implementation"
- [ ] Verify all template variables ({{ .BaseDir }}, etc.) are correctly rendered

### 5.2 End-to-End Validation - Codex

- [ ] Create temp directory for testing
- [ ] Run: `./spectr init --provider=codex` in temp directory
- [ ] Verify `~/.codex/prompts/spectr-proposal.md` exists (or project-local equivalent)
- [ ] Read file and verify it contains "~/.codex/prompts/"
- [ ] Read file and verify it contains "~/.codex/skills/"
- [ ] Verify `~/.codex/prompts/spectr-apply.md` exists
- [ ] Verify all template variables are correctly rendered

### 5.3 End-to-End Validation - Aider (Non-Migrated)

- [ ] Create temp directory for testing
- [ ] Run: `./spectr init --provider=aider` in temp directory
- [ ] Verify `.aider/commands/spectr/proposal.md` exists
- [ ] Read file and verify it contains ONLY generic content
- [ ] Verify it does NOT contain "orchestrator pattern" or "~/.codex/prompts/"
- [ ] Confirm behavior is identical to before this change

### 5.4 Run Spectr Validate

- [ ] Run: `./spectr validate add-provider-specific-templates`
- [ ] Verify validation passes with zero errors
- [ ] If errors exist, read error messages
- [ ] Fix issues in delta specs (ensure all requirements have `#### Scenario:` sections)
- [ ] Re-run validation until it passes

### 5.5 Run Full Test Suite

- [ ] Run: `go test ./...`
- [ ] Verify all tests pass
- [ ] If tests fail, investigate failures
- [ ] Fix failing tests (likely related to template composition)
- [ ] Re-run until all tests pass

### 5.6 Run Linter

- [ ] Run: `golangci-lint run`
- [ ] Verify no new linting errors
- [ ] Fix any linting issues (e.g., missing doc comments, unused variables)
- [ ] Re-run linter until clean

## Phase 6: Documentation & Cleanup

### 6.1 Document Template Sections

- [ ] Create `internal/initialize/templates/SECTIONS.md`
- [ ] Document available sections: "guardrails", "steps", "reference", "main"
- [ ] Describe purpose of each section
- [ ] Provide example of provider override
- [ ] Add guidelines: what to customize (tool-specific) vs preserve (shared best practices)

### 6.2 Update Provider Documentation

- [ ] Update comment in `internal/initialize/providers/provider.go` package doc
- [ ] Add example showing how to create provider-specific templates
- [ ] Add example showing how to use `ProviderSlashCommand()` in `Initializers()`
- [ ] Reference `internal/initialize/templates/SECTIONS.md` for section definitions

### 6.3 Final Verification

- [ ] Re-run all tests: `go test ./...`
- [ ] Re-run validation: `./spectr validate add-provider-specific-templates`
- [ ] Test all 3 scenarios (claude-code, codex, aider) one final time
- [ ] Verify git status shows only expected changes
- [ ] Review all modified files for completeness

## Completion Checklist

- [ ] All tests passing
- [ ] Validation passing (zero errors)
- [ ] Linter clean
- [ ] Claude Code uses custom templates
- [ ] Codex uses custom templates
- [ ] Non-migrated providers (aider) use generic templates
- [ ] Backward compatibility verified
- [ ] Documentation complete
- [ ] Ready for PR/commit
