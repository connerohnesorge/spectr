# Tasks

## 1. Update Markdown Templates

- [ ] 1.1 Add YAML frontmatter to `internal/domain/templates/slash-proposal.md.tmpl`
  - Include `description`, `allowed-tools`, `agent`, `model`, `subtask` fields
  - Keep existing template content after frontmatter
- [ ] 1.2 Add YAML frontmatter to `internal/domain/templates/slash-apply.md.tmpl`
  - Include `description`, `allowed-tools`, `agent`, `model`, `subtask` fields
  - Keep existing template content after frontmatter

## 2. Update TOML Templates

- [ ] 2.1 Review `internal/domain/templates/slash-proposal.toml.tmpl` for consistency
  - TOML already has `description` field, verify it matches markdown version
- [ ] 2.2 Review `internal/domain/templates/slash-apply.toml.tmpl` for consistency
  - TOML already has `description` field, verify it matches markdown version

## 3. Testing

- [ ] 3.1 Run existing template tests to ensure no regressions
  - `go test ./internal/initialize/...`
  - `go test ./internal/domain/...`
- [ ] 3.2 Manually verify generated files contain frontmatter
  - Run `spectr init` in a test directory
  - Check generated slash command files have frontmatter
- [ ] 3.3 Test with Claude Code
  - Verify `/help` shows spectr commands with descriptions
  - Test that `SlashCommand` tool can invoke commands

## 4. Validation

- [ ] 4.1 Run `spectr validate add-slash-command-frontmatter`
- [ ] 4.2 Run linter: `nix develop -c lint`
- [ ] 4.3 Run full test suite: `nix develop -c tests`
