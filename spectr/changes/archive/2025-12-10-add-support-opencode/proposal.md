# Change: Add OpenCode Provider Support

## Why
OpenCode is an emerging AI coding assistant (https://opencode.ai) with a command system similar to Claude Code. Users of OpenCode should be able to use Spectr for spec-driven development with proper instruction file integration and slash commands.

## What Changes
- Add new provider `opencode` to the provider registry
- Create instruction file integration (insert Spectr instructions into user's config)
- Generate slash commands in `.opencode/command/spectr/` directory using Markdown format with YAML frontmatter
- Add `PriorityOpencode` constant to provider priorities

## Impact
- Affected specs: `support-opencode` (new capability)
- Affected code:
  - `internal/initialize/providers/constants.go` (add priority constant)
  - `internal/initialize/providers/opencode.go` (new provider implementation)
  - `internal/initialize/providers/opencode_test.go` (new provider tests)
