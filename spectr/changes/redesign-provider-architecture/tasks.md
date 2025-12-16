## 1. Core Interfaces and Types

- [ ] 1.1 Create `internal/initialize/providers/interfaces.go` with new `Provider`, `Initializer`, `Config`, and `Registration` types
- [ ] 1.2 Create `internal/initialize/providers/registry_new.go` with instance-only `Registry` struct (no global state)
- [ ] 1.3 Add tests for new registry: registration, retrieval, priority sorting, duplicate rejection

## 2. Built-in Initializers

- [ ] 2.1 Create `internal/initialize/providers/initializers/directory.go` with `DirectoryInitializer`
- [ ] 2.2 Create `internal/initialize/providers/initializers/configfile.go` with `ConfigFileInitializer` (marker-based updates)
- [ ] 2.3 Create `internal/initialize/providers/initializers/slashcmds.go` with `SlashCommandsInitializer` (Markdown and TOML formats)
- [ ] 2.4 Add unit tests for each initializer with `afero.MemMapFs`
- [ ] 2.5 Implement initializer deduplication logic in executor

## 3. Migrate Providers

- [ ] 3.1 Migrate `claude.go` to new Provider interface (reference implementation)
- [ ] 3.2 Migrate `gemini.go` to new Provider interface (TOML format example)
- [ ] 3.3 Migrate remaining 15 providers to new interface (batch)
- [ ] 3.4 Add tests verifying all providers return expected initializers

## 4. Executor Integration

- [ ] 4.1 Update `executor.go` to use `afero.NewBasePathFs(osFs, projectPath)`
- [ ] 4.2 Update `executor.go` to collect and dedupe initializers from selected providers
- [ ] 4.3 Update `executor.go` to run initializers with new `Init(ctx, fs, cfg)` signature
- [ ] 4.4 Add integration tests for full initialization flow

## 5. Cleanup - Remove Old Provider System

### 5.1 provider.go - Remove old interfaces and structs
- [ ] 5.1.1 Delete old `Provider` interface (12 methods: `ID`, `Name`, `Priority`, `ConfigFile`, `GetProposalCommandPath`, `GetApplyCommandPath`, `CommandFormat`, `Configure`, `IsConfigured`, `GetFilePaths`, `HasConfigFile`, `HasSlashCommands`)
- [ ] 5.1.2 Delete `TemplateRenderer` interface (`RenderAgents`, `RenderInstructionPointer`, `RenderSlashCommand`)
- [ ] 5.1.3 Delete `BaseProvider` struct and all methods (~300 lines)
- [ ] 5.1.4 Delete `TemplateContext` struct and `DefaultTemplateContext()` function
- [ ] 5.1.5 Keep only `CommandFormat` type and constants (`FormatMarkdown`, `FormatTOML`) - used by new system

### 5.2 registry.go - Remove global state
- [ ] 5.2.1 Delete global `registry` variable and `registryLock` mutex
- [ ] 5.2.2 Delete global functions: `Register()`, `Get()`, `All()`, `IDs()`, `Count()`, `WithConfigFile()`, `WithSlashCommands()`, `Reset()`
- [ ] 5.2.3 Delete `NewRegistryFromGlobal()` function (no global to copy from)
- [ ] 5.2.4 Keep instance-based `Registry` struct and its methods

### 5.3 helpers.go - Migrate to afero.Fs
- [ ] 5.3.1 Update `FileExists(path string)` to `FileExists(fs afero.Fs, path string)`
- [ ] 5.3.2 Update `EnsureDir(path string)` to `EnsureDir(fs afero.Fs, path string)`
- [ ] 5.3.3 Update `UpdateFileWithMarkers()` to accept `afero.Fs` as first parameter
- [ ] 5.3.4 Update `updateSlashCommandBody()` to accept `afero.Fs` as first parameter
- [ ] 5.3.5 Delete `expandPath()` function - no longer needed with project-relative paths
- [ ] 5.3.6 Delete `isGlobalPath()` function - no longer needed with project-relative paths
- [ ] 5.3.7 Delete `findMarkerIndex()` if absorbed into marker functions

### 5.4 constants.go - Remove unused constants
- [ ] 5.4.1 Delete `StandardCommandPaths()` function
- [ ] 5.4.2 Delete `StandardFrontmatter()` function
- [ ] 5.4.3 Delete provider priority constants (`PriorityClaudeCode`, etc.) if only used by old BaseProvider
- [ ] 5.4.4 Keep marker constants (`SpectrStartMarker`, `SpectrEndMarker`) - still used
- [ ] 5.4.5 Keep file permission constants (`dirPerm`, `filePerm`) - still used

### 5.5 Provider files - Verify migration complete
- [ ] 5.5.1 Verify all 17 provider files use new interface (no `BaseProvider` embedding)
- [ ] 5.5.2 Remove `init()` auto-registration from all provider files
- [ ] 5.5.3 Update `NewXxxProvider()` functions to return new `Provider` implementation

## 6. Documentation and Testing

- [ ] 6.1 Update CLI help text for `spectr init` if needed
- [ ] 6.2 Run full test suite and fix any failures
- [ ] 6.3 Manual testing: `spectr init` with various provider combinations
- [ ] 6.4 Verify git diff shows expected file changes after init
