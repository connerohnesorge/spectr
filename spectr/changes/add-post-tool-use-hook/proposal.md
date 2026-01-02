# Change: Add postToolUse Hook Configuration

## Why

Currently conclaude supports hooks for `stop`, `subagentStop`, `preToolUse`,
`permissionRequest`, and `userPromptSubmit`. However, there's no way to hook
into tool use results after they complete. This prevents use cases like:

- Logging Q&A interactions from AskUserQuestion for documentation
- Capturing search results for later reference
- Auditing tool usage patterns
- Building custom integrations that react to tool outputs

A `postToolUse` hook would enable read-only observation of tool results,
allowing users to log, document, or integrate with external systems.

## What Changes

- Add `postToolUse` configuration section to `.conclaude.yaml` schema
- Define hook command configuration with tool filtering (per-tool or all tools)
- Specify environment variable interface for passing tool data to hooks
- Document read-only semantics (hooks observe but cannot modify)
- Add Spectr integration example for AskUserQuestion Q&A logging

## Impact

- Affected specs: `conclaude-hooks` (new capability)
- Affected code: conclaude configuration parser, hook executor
- Not a breaking change - purely additive configuration
