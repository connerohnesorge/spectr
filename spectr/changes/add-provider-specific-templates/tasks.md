# Implementation Tasks: Provider-Specific Template Overrides

## Phase 1: Foundation (Refactor Generic Templates)

### 1.1 Refactor slash-proposal.md.tmpl

- [ ] Read current content of `internal/domain/templates/slash-proposal.md.tmpl`
- [ ] Create section: `{{define "base_guardrails"}}...{{end}}` with current guardrails content
- [ ] Create section: `{{define "guardrails"}}## Guardrails\n{{template "base_guardrails" .}}\n{{end}}`
- [ ] Create section: `{{define "base_steps"}}...{{end}}` with current steps content
- [ ] Create section: `{{define "steps"}}## Steps\n{{template "base_steps" .}}\n{{end}}`
- [ ] Create section: `{{define "base_reference"}}...{{end}}` with current reference content
- [ ] Create section: `{{define "reference"}}## Reference\n{{template "base_reference" .}}\n{{end}}`
- [ ] Create section: `{{define "main"}}# Proposal Creation Guide\n{{template "guardrails" .}}\n{{template "steps" .}}\n{{template "reference" .}}\n{{end}}`
- [ ] Verify template variables ({{ .AgentsFile }}, {{ .BaseDir }}, etc.) are preserved
- [ ] Test: `go test ./internal/domain/...` - verify existing tests still pass

### 1.2 Refactor slash-apply.md.tmpl

- [ ] Read current content of `internal/domain/templates/slash-apply.md.tmpl`
- [ ] Create section: `{{define "base_guardrails"}}...{{end}}` with current guardrails content
- [ ] Create section: `{{define "guardrails"}}## Guardrails\n{{template "base_guardrails" .}}\n{{end}}`
- [ ] Create section: `{{define "base_steps"}}...{{end}}` with current steps content
- [ ] Create section: `{{define "steps"}}## Steps\n{{template "base_steps" .}}\n{{end}}`
- [ ] Create section: `{{define "base_reference"}}...{{end}}` with current reference content
- [ ] Create section: `{{define "reference"}}## Reference\n{{template "base_reference" .}}\n{{end}}`
- [ ] Create section: `{{define "main"}}# Apply Guide\n{{template "guardrails" .}}\n{{template "steps" .}}\n{{template "reference" .}}\n{{end}}`
- [ ] Test: `go test ./internal/domain/...` - verify existing tests still pass

## Phase 2: Update TemplateRef

### 2.1 Add ProviderTemplate field

- [ ] Read `internal/domain/template.go`
- [ ] Add field to `TemplateRef` struct: `ProviderTemplate *template.Template`
- [ ] Update struct doc comment to explain the new field

### 2.2 Implement composition logic

- [ ] Update `Render()` method to check if `ProviderTemplate` is nil
- [ ] If nil: execute base template directly (existing behavior)
- [ ] If non-nil: call new `composeTemplate()` helper function
- [ ] Implement `composeTemplate()` private method:
  - Clone base template: `template.Must(tr.Template.Clone())`
  - Merge provider templates: iterate and call `AddParseTree()`
  - Return composed template
- [ ] In `Render()`, execute "main" template from composed result
- [ ] Update error messages to include template name for debugging

### 2.3 Add unit tests for composition

- [ ] Create/update `internal/domain/template_test.go`
- [ ] Test: `TestTemplateRef_Render_NoProvider` - nil ProviderTemplate uses base directly
- [ ] Test: `TestTemplateRef_Render_WithProvider` - provider section overrides base
- [ ] Test: `TestTemplateRef_Render_PartialOverride` - provider overrides one section, base used for others
- [ ] Test: `TestTemplateRef_composeTemplate` - composition merges correctly
- [ ] Test: Error handling for failed composition
- [ ] Run: `go test ./internal/domain/... -v`
- [ ] Fix any failing tests

## Phase 3: Update TemplateManager

### 3.1 Add provider template parsing

- [ ] Read `internal/initialize/templates.go`
- [ ] Add field to `TemplateManager` struct: `providerTemplates map[string]*template.Template`
- [ ] In `NewTemplateManager()`, after parsing domain templates, add provider parsing:
  - Call `fs.ReadDir(templateFS, "templates/providers")`
  - If error (directory doesn't exist), continue with empty map (backward compatible)
  - For each directory entry that is a directory:
    - Set `providerID = entry.Name()`
    - Parse templates: `template.ParseFS(templateFS, "templates/providers/{id}/*.tmpl")`
    - If parsing fails, return error with context
    - If successful, validate with `validateProviderTemplate()`
    - Store in map: `providerTmpls[providerID] = parsed`

### 3.2 Add validation function

- [ ] Implement `validateProviderTemplate(base, provider *template.Template) error`
- [ ] Check each template in provider against known sections
- [ ] Known sections: `guardrails`, `steps`, `reference`, `main`, `base_guardrails`, `base_steps`, `base_reference`
- [ ] Return error if unknown section found
- [ ] Error message should be clear and actionable

### 3.3 Implement provider-aware methods

- [ ] Implement `ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
  - Build template name map (slash-proposal.md.tmpl, slash-apply.md.tmpl)
  - Return TemplateRef with base template + provider template from map
- [ ] Implement `ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
  - Same pattern but with TOML template names

### 3.4 Update backward compatible methods

- [ ] Update `SlashCommand()` to call `ProviderSlashCommand("", cmd)` internally
- [ ] Update `TOMLSlashCommand()` to call `ProviderTOMLSlashCommand("", cmd)` internally
- [ ] Verify behavior is identical (providers should see nil ProviderTemplate)

### 3.5 Add integration tests

- [ ] Create/update `internal/initialize/templates_test.go`
- [ ] Test: `TestTemplateManager_ParseProviderTemplates` - discovers and parses provider directories
- [ ] Test: `TestTemplateManager_MissingProvidersDirectory` - initialization succeeds with no providers dir
- [ ] Test: `TestTemplateManager_ProviderSlashCommand` - returns correct TemplateRef
- [ ] Test: `TestTemplateManager_ProviderSlashCommand_Unknown` - unknown provider returns nil ProviderTemplate
- [ ] Test: `TestTemplateManager_Backward Compatibility` - old methods still work
- [ ] Test: `TestTemplateManager_ValidationError` - invalid provider template fails at init
- [ ] Run: `go test ./internal/initialize/... -v`
- [ ] Fix any failing tests

## Phase 4: Update Interfaces

### 4.1 Extend TemplateManager interface

- [ ] Read `internal/initialize/providers/provider.go`
- [ ] Add methods to interface:
  - `ProviderSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
  - `ProviderTOMLSlashCommand(providerID string, cmd domain.SlashCommand) domain.TemplateRef`
- [ ] Verify all existing methods remain in interface
- [ ] Add doc comments explaining provider-aware template selection

## Phase 5: Create Provider Templates

### 5.1 Create claude-code templates

- [ ] Create directory: `internal/initialize/templates/providers/claude-code/`
- [ ] Create `slash-proposal.md.tmpl`:
  - Define `guardrails` section
  - Include: "Use the orchestrator pattern: delegate to coder/tester/stuck"
  - Include: "Leverage 200k context window"
  - Include: "Reference AgentSkills in .claude/skills/"
  - Call `{{template "base_guardrails" .}}` to include generic guardrails
- [ ] Create `slash-apply.md.tmpl`:
  - Define `guardrails` section
  - Include: "Delegate implementation to coder subagent"
  - Include: "Verify subagent work: run git diff --stat"
  - Include: "Test every change: use tester subagent"
  - Call `{{template "base_guardrails" .}}`

### 5.2 Create codex templates

- [ ] Create directory: `internal/initialize/templates/providers/codex/`
- [ ] Create `slash-proposal.md.tmpl`:
  - Define `steps` section
  - Include all generic steps
  - Add codex-specific step: "Check ~/.codex/prompts/ for existing Spectr prompts"
  - Add codex-specific step: "Reference AgentSkills in ~/.codex/skills/"
- [ ] Create `slash-apply.md.tmpl`:
  - Define `reference` section
  - Add: "Check ~/.codex/prompts/ for global Spectr prompts"
  - Add: "Check ~/.codex/skills/ for available AgentSkills"
  - Call `{{template "base_reference" .}}` to include generic references

### 5.3 Create gemini placeholder

- [ ] Create directory: `internal/initialize/templates/providers/gemini/`
- [ ] Create `.gitkeep` file
- [ ] Leave as placeholder for future TOML template customization

## Phase 6: Migrate Providers

### 6.1 Update ClaudeProvider

- [ ] Read `internal/initialize/providers/claude.go`
- [ ] Find `Initializers()` method
- [ ] Find `NewSlashCommandsInitializer` call with `tm.SlashCommand()` calls
- [ ] Change `tm.SlashCommand(domain.SlashProposal)` to `tm.ProviderSlashCommand("claude-code", domain.SlashProposal)`
- [ ] Change `tm.SlashCommand(domain.SlashApply)` to `tm.ProviderSlashCommand("claude-code", domain.SlashApply)`
- [ ] Leave all other code unchanged
- [ ] Verify: `go build ./...` - compiles without errors

### 6.2 Update CodexProvider

- [ ] Read `internal/initialize/providers/codex.go`
- [ ] Find `Initializers()` method
- [ ] Find `NewHomeSlashCommandsInitializer` call (or similar) with `tm.SlashCommand()` calls
- [ ] Change `tm.SlashCommand(domain.SlashProposal)` to `tm.ProviderSlashCommand("codex", domain.SlashProposal)`
- [ ] Change `tm.SlashCommand(domain.SlashApply)` to `tm.ProviderSlashCommand("codex", domain.SlashApply)`
- [ ] Leave all other code unchanged
- [ ] Verify: `go build ./...` - compiles without errors

### 6.3 Verify other providers unchanged

- [ ] Pick a non-migrated provider (e.g., aider)
- [ ] Verify it still uses `tm.SlashCommand()` and `tm.TOMLSlashCommand()`
- [ ] Verify: `go build ./...` - compiles without errors
- [ ] Confirm provider will use generic templates (no custom templates for aider)

## Phase 7: Testing & Validation

### 7.1 Run full test suite

- [ ] Run: `go test ./... -v`
- [ ] Fix any failing tests
- [ ] Run: `go test ./... -cover` - verify coverage is adequate
- [ ] Re-run until all tests pass

### 7.2 Run linter

- [ ] Run: `golangci-lint run ./...`
- [ ] Fix any linting errors (unused variables, missing doc comments, etc.)
- [ ] Re-run linter until clean

### 7.3 Build binary

- [ ] Run: `go build -o spectr ./cmd/spectr` (or equivalent for your project)
- [ ] Verify build succeeds without errors

### 7.4 End-to-End Test - Claude Code

- [ ] Create temp directory: `mkdir /tmp/spectr-test-claude`
- [ ] Initialize: `./spectr init --provider=claude-code` in temp dir
- [ ] Verify `.claude/commands/spectr/proposal.md` exists
- [ ] Check file contains: "orchestrator pattern"
- [ ] Check file contains: "coder/tester/stuck"
- [ ] Check file contains: "200k context"
- [ ] Check file contains: "Proposal Creation Guide" (main title)
- [ ] Verify `.claude/commands/spectr/apply.md` exists
- [ ] Check file contains: "Delegate implementation"
- [ ] Check file contains: "Verify subagent work"
- [ ] Verify all template variables rendered (no {{ .BaseDir }} in output)

### 7.5 End-to-End Test - Codex

- [ ] Create temp directory: `mkdir /tmp/spectr-test-codex`
- [ ] Initialize: `./spectr init --provider=codex` in temp dir
- [ ] Verify `~/.codex/prompts/spectr-proposal.md` exists (or equivalent project-local path)
- [ ] Check file contains: "~/.codex/prompts/"
- [ ] Check file contains: "~/.codex/skills/"
- [ ] Check file contains: "Proposal Creation Guide"
- [ ] Verify apply command file exists
- [ ] Check file contains generic reference content

### 7.6 End-to-End Test - Non-Migrated Provider (Aider)

- [ ] Create temp directory: `mkdir /tmp/spectr-test-aider`
- [ ] Initialize: `./spectr init --provider=aider` in temp dir
- [ ] Verify `.aider/commands/spectr/proposal.md` exists
- [ ] Check file does NOT contain: "orchestrator pattern"
- [ ] Check file does NOT contain: "~/.codex/prompts/"
- [ ] Check file contains generic content only
- [ ] Verify behavior is identical to before this change

### 7.7 Validate with spectr

- [ ] Run: `spectr validate add-provider-specific-templates`
- [ ] Verify validation passes: `✓ add-provider-specific-templates valid`
- [ ] If errors exist, read error messages and fix delta specs
- [ ] Re-run until validation passes

## Completion Checklist

- [ ] All unit tests passing
- [ ] All integration tests passing
- [ ] Linter clean (no errors or warnings)
- [ ] Build successful
- [ ] End-to-end tests pass (all 3 providers)
- [ ] Spectr validation passes
- [ ] Design.md documentation complete
- [ ] Code ready for review and merge
