# Change: Auto-Install Slash Commands During Init

## Why

Currently, users selecting config-based AI tools (e.g., `claude-code`, `cline`) in the init wizard must also separately select the corresponding slash command tool (e.g., `claude`, `cline`) to get the full Spectr experience. This creates confusion and duplicates tool entries in the wizard, making the selection process less intuitive.

By automatically installing slash commands when their matching config-based tool is selected, we provide a more cohesive experience similar to OpenSpec's approach, where selecting a tool gives you everything you need in one action.

## What Changes

- Config-based tools (claude-code, cline, etc.) will automatically trigger installation of their corresponding slash command files
- Redundant slash-only tool entries will be removed from the wizard registry
- Tool selection remains simple and flat, but each selection does more behind the scenes
- The executor will invoke both config and slash command configurators when applicable
- No UI changes - slash commands install silently as part of the tool configuration

**Affected tools**:
- `claude-code` → auto-installs `.claude/commands/spectr/{proposal,apply,archive}.md`
- `cline` → auto-installs `.clinerules/spectr-{proposal,apply,archive}.md`
- `cursor` → auto-installs `.cursor/commands/spectr-{proposal,apply,archive}.md`
- All 11 tools with slash command variants

## Impact

**Affected specs**:
- `cli-interface` - Tool selection and initialization workflow

**Affected code**:
- `internal/init/registry.go` - Tool registry and mappings
- `internal/init/executor.go` - Configuration orchestration
- `internal/init/configurator.go` - Potentially needs composite configurator pattern

**Breaking changes**:
- Users who previously selected only slash command tools (without config) will no longer see them in the wizard
- Migration: They should now select the config-based version which installs both

**User benefits**:
- Simpler tool selection (fewer redundant entries)
- Complete Spectr setup with one selection
- Consistent with OpenSpec's tool handling approach
