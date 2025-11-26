# Change: Add Codex CLI Provider Support

## Why
Codex CLI is OpenAI's open-source CLI for agentic coding. It uses custom prompts (slash commands) stored in `~/.codex/prompts/` with YAML frontmatter. Adding Codex as a provider expands Spectr's reach to users of this emerging AI coding tool.

## What Changes
- Add new `codex` provider to `internal/init/providers/`
- Provider uses `AGENTS.md` instruction file and global slash commands
- Commands install to global path `~/.codex/prompts/spectr/`
- **BREAKING PATTERN**: First provider using global paths rather than project-level
- Frontmatter uses `description:` and optionally `argument-hint:` fields

## Impact
- Affected specs: New `support-codex` capability
- Affected code: `internal/init/providers/codex.go` (new file)
- Affected code: `internal/init/providers/constants.go` (add priority constant)
- Affected code: `internal/init/providers/provider.go` (may need global path support)
