# Change: Refactor Providers to Composable File Initializers

## Why

The current provider architecture has accumulated tech debt:

1. **Monolithic interface**: 10+ methods including `Configure()`, `ConfigFile()`, `GetProposalCommandPath()`, `HasConfigFile()`, `HasSlashCommands()`, `CommandFormat()`, etc.
2. **BaseProvider indirection**: Providers embed `BaseProvider` which adds complexity without clear benefit
3. **Rigid file types**: Adding new file types (Claude skills, MCP configs) requires changing the Provider interface
4. **Redundant detection methods**: `HasConfigFile()` and `HasSlashCommands()` duplicate information that's already explicit in provider composition

Moving to composable initializers with a minimal interface enables extensibility while removing tech debt.

## What Changes

- **Minimal Provider interface**: Reduced from 10+ methods to 6: `ID()`, `Name()`, `Priority()`, `Initializers()`, `IsConfigured()`, `GetFilePaths()`
- **Remove BaseProvider**: Providers implement interface directly using helper functions
- **Remove Configure() from interface**: Configuration via `ConfigureInitializers()` helper function
- **Remove tech debt methods**: `HasConfigFile()`, `HasSlashCommands()`, `ConfigFile()`, `GetProposalCommandPath()`, `GetApplyCommandPath()`, `CommandFormat()` all removed
- **Remove registry filter methods**: `WithConfigFile()` and `WithSlashCommands()` removed from registry
- **New FileInitializer interface**: Atomic unit for single-file operations (4 methods: ID, FilePath, Configure, IsConfigured)
- **Standard initializers**: `InstructionFileInitializer`, `MarkdownSlashCommandInitializer`, `TOMLSlashCommandInitializer`
- **Helper functions**: `ConfigureInitializers()`, `AreInitializersConfigured()`, `GetInitializerPaths()`
- **Initializer independence**: No dependencies between initializers; order in list implies sequence
- **Fail-fast error handling**: `ConfigureInitializers()` stops on first error, no rollback
- **Wizard simplification**: Remove filtering by capabilities (all providers shown equally)

## Impact

- Affected specs: `cli-interface` (Provider Interface, Per-Provider File Organization, Command Format Support)
- Affected code:
  - `internal/initialize/providers/provider.go` - Minimal 6-method interface, remove BaseProvider
  - `internal/initialize/providers/registry.go` - Remove `WithConfigFile()` and `WithSlashCommands()` filters
  - `internal/initialize/providers/initializer.go` - New FileInitializer interface
  - `internal/initialize/providers/*_initializer.go` - New initializer implementations
  - All 23 provider files - Rewrite to new pattern (implement 6 methods directly)
  - `internal/initialize/executor.go` - Use `ConfigureInitializers()` helper instead of `provider.Configure()`
  - `internal/initialize/wizard.go` - Remove capability-based filtering
  - `internal/initialize/providers/*_test.go` - Update tests for new interface
- **BREAKING**: Internal Provider interface changes significantly (no public API impact)

## Scope

This change establishes the initializer architecture only. Actual implementations of new file types (Claude skills, agents, MCP servers) are separate follow-up proposals that will leverage this architecture.

## Migration Approach

Single atomic PR converting all 23 providers at once. No mixed states or phased migration.

## Future Capabilities Enabled

- **Claude Skill**: Custom `.claude/skills/spectr.md` for enhanced CLI guidance (follow-up proposal)
- **Spectr Agent**: Agent definition for multi-step workflows (follow-up proposal)
- **MCP Servers**: Per-provider MCP server configurations (follow-up proposal)
- **Custom Rules**: Provider-specific rule files (follow-up proposal)
