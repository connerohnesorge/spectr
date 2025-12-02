## 1. Update Agent Prompt Files

- [x] 1.1 Update `spectr/AGENTS.md` to remove `spectr list` references and replace with directory reading instructions
- [x] 1.2 Update `.agent/workflows/spectr-proposal.md`
- [x] 1.3 Update `.agent/workflows/spectr-sync.md`
- [x] 1.4 Update `.agent/workflows/spectr-apply.md`

## 2. Update Claude Command Files

- [x] 2.1 Update `.claude/commands/spectr/proposal.md`
- [x] 2.2 Update `.claude/commands/spectr/sync.md`
- [x] 2.3 Update `.claude/commands/spectr/apply.md`

## 3. Update Gemini Command Files

- [x] 3.1 Update `.gemini/commands/spectr/proposal.toml`
- [x] 3.2 Update `.gemini/commands/spectr/sync.toml`
- [x] 3.3 Update `.gemini/commands/spectr/apply.toml`

## 4. Update Templates

- [x] 4.1 Update `internal/initialize/templates/spectr/AGENTS.md.tmpl`
- [x] 4.2 Update `internal/initialize/templates/tools/slash-apply.md.tmpl`
- [x] 4.3 Update `internal/initialize/templates/tools/slash-sync.md.tmpl`
- [x] 4.4 Update `internal/initialize/templates/tools/slash-proposal.md.tmpl`

## 5. Validation

- [x] 5.1 Run `spectr validate remove-spectr-list-from-agent-prompts --strict`
- [x] 5.2 Verify no remaining `spectr list` references in agent prompt files (grep check)
- [x] 5.3 Confirm user-facing docs still have `spectr list` references (README.md, docs/)
