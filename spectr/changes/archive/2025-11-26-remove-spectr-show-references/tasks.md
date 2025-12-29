# Implementation Tasks

## 1. Core AI Agent Documentation Files

These files contain instructions for AI agents. Replace `spectr show` with
direct file reading instructions (NOT `spectr view`).

- [x] 1.1 Update `CLAUDE.md` - replace `spectr show` with direct file reading
  (`spectr/specs/<cap>/spec.md`, `spectr/changes/<id>/proposal.md`)
- [x] 1.2 Update `AGENTS.md` - replace `spectr show` with direct file reading
- [x] 1.3 Update `spectr/AGENTS.md` - replace `spectr show` with direct file
  reading
- [x] 1.4 Update `.github/copilot-instructions.md` - replace `spectr show` with
  direct file reading

## 2. AI Tool Command Files

These files are slash commands for AI agents. Replace `spectr show` with direct
file reading instructions.

- [x] 2.1 Update `.claude/commands/spectr/proposal.md` - replace `spectr show`
  with direct file reading
- [x] 2.2 Update `.claude/commands/spectr/apply.md` - replace `spectr show` with
  direct file reading
- [x] 2.3 Update `.claude/commands/spectr/sync.md` - replace `spectr show` with
  direct file reading
- [x] 2.4 Update `.gemini/commands/spectr/*.toml` - replace `spectr show` with
  direct file reading
- [x] 2.5 Update `.agent/workflows/*.md` - replace `spectr show` with direct
  file reading (if exists)
- [x] 2.6 Update `.gemini/commands/spectr-*.toml` - replace `spectr show` with
  direct file reading

## 3. Template Files

These are templates for `spectr init`. Replace `spectr show` with direct file
reading instructions.

- [x] 3.1 Update `internal/init/templates/spectr/AGENTS.md.tmpl` - replace
  `spectr show` with direct file reading
- [x] 3.2 Update `internal/init/templates/tools/slash-apply.md.tmpl` - replace
  `spectr show` with direct file reading
- [x] 3.3 Update `internal/init/templates/tools/slash-proposal.md.tmpl` -
  replace `spectr show` with direct file reading
- [x] 3.4 Update `internal/init/templates/tools/slash-sync.md.tmpl` - replace
  `spectr show` with direct file reading

## 4. Public User Documentation

These files are for human users. Replace `spectr show` with `spectr view` (the
correct command for users).

- [x] 4.1 Update `docs/src/content/docs/reference/cli-commands.md` - replace
  `spectr show` with `spectr view`
- [x] 4.2 Update `docs/src/content/docs/guides/archiving-workflow.md` - replace
  `spectr show` with `spectr view`
- [x] 4.3 Update `README.md` - replace `spectr show` with `spectr view` (if any
  references exist)

## 5. Validation

- [x] 5.1 Run `spectr validate remove-spectr-show-references --strict` to ensure
  change is valid
- [x] 5.2 Run `rg "spectr show" --type md --type go` to verify no remaining
  references (except archive/ history)
