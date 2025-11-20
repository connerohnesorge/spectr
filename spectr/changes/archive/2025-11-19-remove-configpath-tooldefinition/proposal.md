# Change: Remove ConfigPath from ToolDefinition

## Why

The `ConfigPath` field in `ToolDefinition` (internal/init/models.go:24) is misleading and serves no functional purpose. While it's set for each tool in the registry and returned by `getToolFileInfo()`, it's never actually used to determine where files are created. Each configurator hardcodes its own file path independently:

- `ClaudeCodeConfigurator` uses `"CLAUDE.md"` (configurator.go:163)
- `ClineConfigurator` uses `"CLINE.md"` (configurator.go:193)
- `CostrictConfigurator` uses `"COSTRICT.md"` (configurator.go:221)

The values stored in `ConfigPath` (like `.claude/claude.json`, `.cline/cline_mcp_settings.json`) don't match the actual files created and are never referenced when creating configuration files. This creates confusion and maintenance overhead.

## What Changes

- **BREAKING**: Remove `ConfigPath` field from `ToolDefinition` struct
- Remove `ConfigPath` assignments in `registry.go` for all tools
- Update `getToolFileInfo()` to return actual file paths from configurators instead of using ConfigPath
- Update tests that reference ConfigPath
- No user-facing behavior changes (files are still created in the same locations)

## Impact

- Affected specs: cli-interface (initialization wizard and tool configuration)
- Affected code:
  - `internal/init/models.go` - remove field from struct
  - `internal/init/registry.go` - remove ConfigPath assignments
  - `internal/init/executor.go` - update getToolFileInfo logic
  - `internal/init/registry_test.go` - remove ConfigPath assertions
  - `internal/init/init_test.go` - update tests if needed
