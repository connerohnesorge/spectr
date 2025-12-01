## 1. Update Agent Prompt Files

- [ ] 1.1 Update `spectr/AGENTS.md` to remove `spectr list` references and replace with directory reading instructions
- [ ] 1.2 Update `.agent/workflows/spectr-proposal.md`
- [ ] 1.3 Update `.agent/workflows/spectr-sync.md`
- [ ] 1.4 Update `.agent/workflows/spectr-apply.md`

## 2. Update Claude Command Files

- [ ] 2.1 Update `.claude/commands/spectr/proposal.md`
- [ ] 2.2 Update `.claude/commands/spectr/sync.md`
- [ ] 2.3 Update `.claude/commands/spectr/apply.md`

## 3. Update Gemini Command Files

- [ ] 3.1 Update `.gemini/commands/spectr/proposal.toml`
- [ ] 3.2 Update `.gemini/commands/spectr/sync.toml`
- [ ] 3.3 Update `.gemini/commands/spectr/apply.toml`

## 4. Update Templates

- [ ] 4.1 Update `internal/initialize/templates/spectr/AGENTS.md.tmpl`
- [ ] 4.2 Update `internal/initialize/templates/tools/slash-apply.md.tmpl`
- [ ] 4.3 Update `internal/initialize/templates/tools/slash-sync.md.tmpl`

## 5. Validation

- [ ] 5.1 Run `spectr validate remove-spectr-list-from-agent-prompts --strict`
- [ ] 5.2 Verify no remaining `spectr list` references in agent prompt files (grep check)
- [ ] 5.3 Confirm user-facing docs still have `spectr list` references (README.md, docs/)
