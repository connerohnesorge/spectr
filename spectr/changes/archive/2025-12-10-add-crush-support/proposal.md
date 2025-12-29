# Change: Add Crush AI Assistant Support

## Why

Add support for Crush AI assistant in Spectr to enable developers to use Crush's enhanced capabilities for spec-driven development workflows. Crush is a glamorous CLI AI coding agent from Charmbracelet that supports multiple LLMs, MCP integration, and session-based work.

## What Changes

- Add Crush provider to the provider registry (`internal/initialize/providers/crush.go`)
- Add `PriorityCrush` constant to `internal/initialize/providers/constants.go`
- Create Crush-specific slash commands in `.crush/commands/spectr/` directory
- Provider configures `CRUSH.md` instruction file for the Spectr marker injection
- **BREAKING**: None - this is additive functionality

## Impact

- Affected specs: support-crush (new capability)
- Affected code:
  - `internal/initialize/providers/crush.go` (new file)
  - `internal/initialize/providers/constants.go` (add priority constant)
- New files created during `spectr init`:
  - `CRUSH.md` (instruction file with spectr markers)
  - `.crush/commands/spectr/proposal.md`
  - `.crush/commands/spectr/apply.md`
