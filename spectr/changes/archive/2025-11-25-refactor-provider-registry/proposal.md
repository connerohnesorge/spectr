# Change: Refactor tool registry to interface-driven provider architecture

## Why

The current `internal/init/` tool registry is tightly coupled with hardcoded
tool configurations in a single `tool_definitions.go` file. This makes it
difficult to add new AI CLI providers (e.g., Gemini CLI) and violates Go's
interface-driven design principles. The current implementation uses global maps
and requires modifying multiple files to add a single provider.

## What Changes

- **Provider interface**: Define a `Provider` interface that all CLI/IDE
  integrations implement
- **Per-provider files**: Each provider (Claude, Gemini, Cline, Cursor, etc.)
  gets its own Go file containing its configuration
- **Registration pattern**: Use a registry pattern with `Register()` functions
  called from `init()` in each provider file
- **Support TOML-based commands**: Add support for Gemini-style TOML command
  definitions alongside Markdown
- **Remove global maps**: Replace `toolConfigs` and `slashToolConfigs` maps with
  the interface-driven registry

## Impact

- Affected specs: `cli-framework`
- Affected code:
  - `internal/init/registry.go` (major refactor)
  - `internal/init/tool_definitions.go` (split into provider files)
  - `internal/init/configurator.go` (minor updates)
  - `internal/init/executor.go` (use new Provider interface)
  - New files: `internal/init/providers/*.go` (one per provider)
