# Change: Replace SlashDir with Per-Command Path Methods

## Why

The current `SlashDir()` method returns a single directory for all slash commands, forcing a rigid file structure (`{slashDir}/spectr-{cmd}.md`). This prevents providers from having:

- Different locations for each command type
- Custom file naming conventions (not just `spectr-{cmd}`)
- Mixed formats within a provider (e.g., one command as TOML, another as markdown)
- Commands stored outside a single directory

Gemini already works around this by overriding multiple methods (`Configure`, `configureSlashCommands`, `getSlashCommandPath`, `IsConfigured`, `GetFilePaths`). The new design eliminates this workaround pattern.

## What Changes

- **BREAKING**: Remove `SlashDir() string` from Provider interface
- **BREAKING**: Remove `slashDir` field from BaseProvider struct
- Add `GetProposalCommandPath() string` returning relative path for proposal command
- Add `GetArchiveCommandPath() string` returning relative path for archive command
- Add `GetApplyCommandPath() string` returning relative path for apply command
- Update `HasSlashCommands()` to return true if ANY command path is non-empty
- Update all provider implementations to use new path methods
- Simplify GeminiProvider by removing method overrides (uses TOML paths directly)

## Impact

- Affected specs: `cli-framework` (Provider Interface requirement)
- Affected code:
  - `internal/init/providers/provider.go` (interface + BaseProvider)
  - `internal/init/providers/*.go` (all provider implementations)
  - `internal/init/providers/*_test.go` (tests)
