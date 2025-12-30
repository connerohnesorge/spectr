# Change: Remove incorrect `spectr show` references from documentation

## Why

The documentation across multiple files references a `spectr show` command that
does **not exist**. The CLI only has these commands: `init`, `list`, `validate`,
`archive`, `view`, `version`. Users attempting to use `spectr show` receive an
error: "unexpected argument show". This misleads users and AI assistants who
rely on the documentation.

Additionally, AI agent documentation incorrectly suggests using CLI commands to
read specs/changes. **AI agents should read files directly** using their native
file reading capabilities (e.g., `Read` tool, `cat`), not CLI commands like
`spectr view`. The `spectr view` command is for human users only.

## What Changes

- **BREAKING**: Remove all references to `spectr show` command from
  documentation
- Remove `spectr show [item]`, `spectr show <id>`, `spectr show <spec>`, and
  related examples
- **For user documentation**: Reference `spectr view` as the correct command for
  viewing specs/changes
- **For AI agent documentation**: Replace CLI viewing commands with direct file
  reading instructions
  - Specs: Read `spectr/specs/<capability>/spec.md` directly
  - Changes: Read `spectr/changes/<change-id>/proposal.md` directly
- Update troubleshooting sections to remove `spectr show` debug suggestions
- Update Quick Reference sections appropriately for the audience (users vs AI
  agents)

## Impact

- Affected specs: `documentation`
- Affected files:
  - `CLAUDE.md` (root)
  - `AGENTS.md` (root)
  - `spectr/AGENTS.md`
  - `.github/copilot-instructions.md`
  - `.claude/commands/spectr/*.md`
  - `.gemini/commands/spectr/*.toml`
  - `.agent/workflows/*.md`
  - `internal/init/templates/spectr/AGENTS.md.tmpl`
  - `internal/init/templates/tools/*.md.tmpl`
  - `docs/src/content/docs/reference/cli-commands.md`
  - `docs/src/content/docs/guides/archiving-workflow.md`
  - `README.md`
