# Tasks: Add Frontmatter Override System

## Phase 1: Define Base Frontmatter in Go

- [ ] Create `internal/domain/frontmatter.go`
- [ ] Add `FrontmatterOverride` struct with `Set` and `Remove` fields
- [ ] Define `BaseSlashCommandFrontmatter` map with frontmatter for proposal and apply commands
- [ ] Implement `GetBaseFrontmatter(cmd domain.SlashCommand) map[string]interface{}` with deep copy
- [ ] Implement `copyMap(src map[string]interface{}) map[string]interface{}` helper for deep copying
- [ ] Write tests in `internal/domain/frontmatter_test.go` for:
  - GetBaseFrontmatter returns correct values for each command
  - GetBaseFrontmatter returns copy (mutations don't affect base)
  - GetBaseFrontmatter returns empty map for unknown commands

## Phase 2: Frontmatter Manipulation Utilities

- [ ] Implement `ApplyFrontmatterOverrides(base, overrides) map[string]interface{}` in `frontmatter.go`
  - Copy base map
  - Apply Set operations (merge/overwrite)
  - Apply Remove operations (delete keys)
  - Handle nil overrides (return copy of base)
- [ ] Implement `RenderFrontmatter(fm map[string]interface{}, body string) (string, error)` using `gopkg.in/yaml.v3`
  - Marshal fm to YAML
  - Wrap in `---` delimiters
  - Append body content
- [ ] Write tests for ApplyFrontmatterOverrides:
  - Set operations add/modify fields
  - Remove operations delete fields
  - Set + Remove (Remove wins)
  - Nil overrides returns copy
  - Empty overrides returns copy
- [ ] Write tests for RenderFrontmatter:
  - Renders valid YAML with delimiters
  - Handles empty map (empty frontmatter section)
  - Preserves body content exactly
  - Returns error on YAML marshal failure

## Phase 3: Extend TemplateRef

- [ ] Add `Command domain.SlashCommand` field to `TemplateRef` in `internal/domain/template.go`
- [ ] Add `Overrides *FrontmatterOverride` field to `TemplateRef`
- [ ] Update `TemplateRef.Render(ctx)` to:
  - Render template body first (execute Go template)
  - Get base frontmatter via `GetBaseFrontmatter(tr.Command)`
  - Apply overrides if present
  - Call `RenderFrontmatter(fm, body)` to assemble final content
- [ ] Write tests in `internal/domain/template_test.go`:
  - Render with no overrides produces base frontmatter + body
  - Render with Set overrides adds/modifies fields
  - Render with Remove overrides deletes fields
  - Render with nil overrides behaves like no overrides

## Phase 4: Remove Frontmatter from Templates

- [ ] Edit `internal/domain/templates/slash-proposal.md.tmpl` - remove `---` frontmatter block, keep only body
- [ ] Edit `internal/domain/templates/slash-apply.md.tmpl` - remove `---` frontmatter block, keep only body
- [ ] Update any tests that depend on template content format

## Phase 5: TemplateManager Extension

- [ ] Add `SlashCommandWithOverrides(cmd, overrides)` method to `TemplateManager` in `internal/initialize/templates.go`
  - Call `tm.SlashCommand(cmd)` to get base TemplateRef
  - Set `ref.Command = cmd`
  - Set `ref.Overrides = overrides`
  - Return modified ref
- [ ] Update existing `SlashCommand(cmd)` method to set `Command` field
- [ ] Add tests in `internal/initialize/templates_test.go`:
  - SlashCommandWithOverrides returns TemplateRef with Command and Overrides set
  - Rendering produces correct frontmatter

## Phase 6: Update ClaudeProvider

- [ ] Modify `ClaudeProvider.Initializers()` in `internal/initialize/providers/claude.go`:
  - For proposal: use `tm.SlashCommandWithOverrides(domain.SlashProposal, &domain.FrontmatterOverride{Set: map[string]interface{}{"context": "fork"}, Remove: []string{"agent"}})`
  - For apply: use `tm.SlashCommand(domain.SlashApply)` (no overrides)
- [ ] Create `internal/initialize/providers/claude_test.go` if it doesn't exist
- [ ] Add test verifying:
  - ClaudeProvider returns correct initializers
  - Proposal TemplateRef has overrides set
  - Apply TemplateRef has no overrides

## Phase 7: Integration Testing

- [ ] Run `nix develop -c go build` to ensure it compiles
- [ ] Run `nix develop -c tests` to verify all tests pass
- [ ] Run `nix develop -c lint` to verify linting passes
- [ ] Create temp directory and run `./spectr init --provider=claude-code`
- [ ] Verify `.claude/commands/spectr/proposal.md` contains:
  - `context: fork` in frontmatter
  - NO `agent:` field in frontmatter
  - Body content unchanged
- [ ] Verify `.claude/commands/spectr/apply.md` matches expected base content

## Phase 8: Validation & Cleanup

- [ ] Run `spectr validate add-frontmatter-override-system`
- [ ] Fix any validation errors in spec or implementation
- [ ] Add/update doc comments for all exported types and functions
- [ ] Run `nix fmt` to format code
- [ ] Final test run: `nix develop -c tests && nix develop -c lint`
