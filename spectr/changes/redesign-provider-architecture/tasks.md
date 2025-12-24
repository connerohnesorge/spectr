## 0. Domain Package: Shared Types to Break Import Cycles

Create `internal/domain` package with shared domain types that can be imported by both `providers` and `templates` packages without creating import cycles.

- [ ] 0.1 Create `internal/domain/template.go` with:
  - `TemplateRef` struct with `Name string` and `Template *template.Template` fields
  - `Render(ctx TemplateContext) (string, error)` method on `TemplateRef`
  - `TemplateContext` struct with `BaseDir`, `SpecsDir`, `ChangesDir`, `ProjectFile`, `AgentsFile` fields
  - `DefaultTemplateContext()` function returning default values
- [ ] 0.2 Create `internal/domain/slashcmd.go` with:
  - `SlashCommand int` type with `SlashProposal`, `SlashApply` constants
  - `String() string` method for debugging
  - `TemplateName() string` method returning the .tmpl file name
- [ ] 0.3 Add unit tests for `internal/domain/template_test.go`:
  - Test `TemplateRef.Render()` with mock template
  - Test `DefaultTemplateContext()` returns expected values
- [ ] 0.4 Add unit tests for `internal/domain/slashcmd_test.go`:
  - Test `SlashCommand.String()` returns correct names
  - Test `SlashCommand.TemplateName()` returns correct template file names

## 1. Foundation: Core Interfaces and Types

These tasks establish the new architecture without breaking the existing code.

- [ ] 1.1 Create `internal/initialize/providers/initializer.go` with new `Initializer` interface including:
  - `Init(ctx context.Context, fs afero.Fs, cfg *Config, tm *TemplateManager) (InitResult, error)`
  - `IsSetup(fs afero.Fs, cfg *Config) bool`
  - `Path() string` (for deduplication)
  - `IsGlobal() bool` (true = use globalFs, false = use projectFs)
- [ ] 1.2 Create `internal/initialize/providers/result.go` with `InitResult` struct:
  - `CreatedFiles []string`
  - `UpdatedFiles []string`
- [ ] 1.3 Create `internal/initialize/providers/config.go` with `Config` struct containing:
  - `SpectrDir string` field
  - `SpecsDir() string` method (returns SpectrDir + "/specs")
  - `ChangesDir() string` method (returns SpectrDir + "/changes")
  - `ProjectFile() string` method (returns SpectrDir + "/project.md")
  - `AgentsFile() string` method (returns SpectrDir + "/AGENTS.md")
- [ ] 1.4 Rewrite `internal/initialize/providers/provider.go` with new minimal `Provider` interface returning `[]Initializer` (replace old interface in-place)
- [ ] 1.5 Create `internal/initialize/providers/registration.go` with:
  - `Registration` struct (ID, Name, Priority, Provider)
  - `RegisterProvider(reg Registration) error` function with validation
  - `RegisterAllProviders() error` function that registers all built-in providers explicitly (no init())
  - Registry data structure and accessor functions (`Get`, `All`, `IDs`, `Count`)

## 2. Type-Safe Template System

Update TemplateManager to use domain types and add type-safe accessor methods.

- [ ] 2.1 Update `internal/initialize/templates.go` to import and use `domain.TemplateRef` and `domain.TemplateContext`
- [ ] 2.2 Add type-safe accessor methods to `TemplateManager` returning `domain.TemplateRef`:
  - `InstructionPointer() domain.TemplateRef`
  - `Agents() domain.TemplateRef`
  - `Project() domain.TemplateRef`
  - `CIWorkflow() domain.TemplateRef`
- [ ] 2.3 Add `SlashCommand(cmd domain.SlashCommand) domain.TemplateRef` method to `TemplateManager`
- [ ] 2.4 Add unit tests verifying all accessor methods return valid `domain.TemplateRef`
- [ ] 2.5 Update any existing code that uses TemplateContext to use `domain.TemplateContext`

## 3. Built-in Initializers

Create the three composable initializers that providers will use. Each must implement `Path()` and `IsGlobal()`.

- [ ] 3.1 Create `internal/initialize/providers/initializers/directory.go` with `DirectoryInitializer`:
  - Implements `Init()`, `IsSetup()`, `Path()`, `IsGlobal()`
  - Accepts directory path(s) and isGlobal flag
  - Creates directories with `MkdirAll`
- [ ] 3.2 Create `internal/initialize/providers/initializers/configfile.go` with `ConfigFileInitializer`:
  - Implements `Init()`, `IsSetup()`, `Path()`, `IsGlobal()`
  - Receives `TemplateGetter func(*TemplateManager) domain.TemplateRef` (compile-time checked)
  - Handles both create and update scenarios with marker-based updates
- [ ] 3.3 Create `internal/initialize/providers/initializers/slashcmds.go` with `SlashCommandsInitializer`:
  - Implements `Init()`, `IsSetup()`, `Path()`, `IsGlobal()`
  - Receives `[]domain.SlashCommand` (compile-time checked command types from domain package)
  - Supports both Markdown and TOML output formats
- [ ] 3.4 Add unit tests for `DirectoryInitializer` with `afero.MemMapFs`
- [ ] 3.5 Add unit tests for `ConfigFileInitializer` with `afero.MemMapFs` - test create and marker update scenarios
- [ ] 3.6 Add unit tests for `SlashCommandsInitializer` with `afero.MemMapFs` - test Markdown and TOML formats

## 4. New Registry Implementation (No init())

Replace the old registry with explicit registration - no init() functions in provider files. Replace in-place, do not create V2 files.

- [ ] 4.1 Rewrite `internal/initialize/providers/registry.go` with new `Registry` type using `Registration` struct for metadata
- [ ] 4.2 Implement `RegisterProvider(reg Registration) error` with validation (returns error, not panic)
- [ ] 4.3 Implement `Get(id)`, `All()`, `IDs()`, `Count()` methods on registry
- [ ] 4.4 Add priority-sorted retrieval maintaining backwards-compatible behavior
- [ ] 4.5 Add duplicate ID rejection with clear error messages
- [ ] 4.6 Create `RegisterAllProviders() error` function that explicitly registers all 17 providers in one place
- [ ] 4.7 Rewrite `registry_test.go` with tests for new registry: registration, retrieval, priority sorting, duplicate rejection
- [ ] 4.8 Add unit tests for `RegisterAllProviders()` verifying all providers are registered correctly

## 5. Migrate Providers (In-Place Replacement, No init())

Migrate each provider to the new interface. Each provider file should ONLY contain the struct and Initializers() method - NO init() function. Registration happens in `RegisterAllProviders()`.

- [ ] 5.1 Migrate `claude.go` to new Provider interface (reference implementation) - delete old struct, `NewClaudeProvider()`, and `init()`
- [ ] 5.2 Migrate `gemini.go` to new Provider interface (TOML format example) - delete old struct, `Configure()` override, and `init()`
- [ ] 5.3 Migrate `cursor.go` to new Provider interface - delete old code and `init()`
- [ ] 5.4 Migrate `cline.go` to new Provider interface - delete old code and `init()`
- [ ] 5.5 Migrate `aider.go` to new Provider interface - delete old code and `init()`
- [ ] 5.6 Migrate `codex.go` to new Provider interface - delete old code and `init()`
- [ ] 5.7 Migrate `costrict.go` to new Provider interface - delete old code and `init()`
- [ ] 5.8 Migrate `qoder.go` to new Provider interface - delete old code and `init()`
- [ ] 5.9 Migrate `codebuddy.go` to new Provider interface - delete old code and `init()`
- [ ] 5.10 Migrate `qwen.go` to new Provider interface - delete old code and `init()`
- [ ] 5.11 Migrate `antigravity.go` to new Provider interface - delete old code and `init()`
- [ ] 5.12 Migrate `windsurf.go` to new Provider interface - delete old code and `init()`
- [ ] 5.13 Migrate `kilocode.go` to new Provider interface - delete old code and `init()`
- [ ] 5.14 Migrate `continue.go` to new Provider interface - delete old code and `init()`
- [ ] 5.15 Migrate `crush.go` to new Provider interface - delete old code and `init()`
- [ ] 5.16 Migrate `opencode.go` to new Provider interface - delete old code and `init()`

## 6. Executor Integration

Update executor to use new architecture with dual filesystem, ordering, deduplication, and InitResult collection.

- [ ] 6.0 Update `cmd/root.go` or application entry point to call `providers.RegisterAllProviders()` at startup with error handling
- [ ] 6.1 Create dual filesystem in `executor.go`:
  - `projectFs := afero.NewBasePathFs(osFs, projectPath)` for project-relative paths
  - `globalFs := afero.NewBasePathFs(osFs, os.UserHomeDir())` for global paths
- [ ] 6.2 Update `executor.go` to use new registry API (`Registration` based retrieval)
- [ ] 6.3 Implement initializer collection from selected providers
- [ ] 6.4 Implement initializer deduplication by `Path()` - same path = run once
- [ ] 6.5 Implement initializer sorting by type (guaranteed order):
  - 1. `DirectoryInitializer`
  - 2. `ConfigFileInitializer`
  - 3. `SlashCommandsInitializer`
- [ ] 6.6 Update `configureProviders()` to:
  - Select `projectFs` or `globalFs` based on `initializer.IsGlobal()`
  - Call `Init(ctx, fs, cfg, templateManager)` on each initializer
  - Collect `InitResult` from each initializer
- [ ] 6.7 Aggregate `InitResult` values into `ExecutionResult`
- [ ] 6.8 Handle partial failures: report which initializers failed, continue with rest

## 7. Remove Old Provider Code (In-Place Cleanup)

Remove deprecated code that has been replaced by the new architecture.

- [ ] 7.1 Remove old `Provider` interface (12-method version) from `provider.go`
- [ ] 7.2 Remove `BaseProvider` struct from `provider.go`
- [ ] 7.3 Remove `TemplateRenderer` interface from `provider.go` (now using `*TemplateManager` directly)
- [ ] 7.4 Remove `TemplateContext` and `DefaultTemplateContext()` from `provider.go` (now in `internal/domain`)
- [ ] 7.5 Delete `helpers.go` - `EnsureDir`, `FileExists`, `UpdateFileWithMarkers` now in initializers or use `afero.Fs`
- [ ] 7.6 Remove old global registry functions from `registry.go` (keep only new `Registry` type)
- [ ] 7.7 Clean up `constants.go` - remove `StandardFrontmatter()`, `StandardCommandPaths()`, `PrefixedCommandPaths()` (moved to initializers)
- [ ] 7.8 Remove priority constants from `constants.go` (priorities now in registration calls)

## 8. Test Cleanup

Update tests to match new architecture - rewrite in-place, no V2 files.

- [ ] 8.1 Rewrite `provider_test.go` with tests for new `Provider` interface (replace old tests, not new file)
- [ ] 8.2 Add tests verifying all 17 providers return expected initializers
- [ ] 8.3 Add tests verifying provider registration metadata (ID, Name, Priority)
- [ ] 8.4 Rewrite `registry_test.go` with tests for new `Registry` type (replace old tests, not new file)
- [ ] 8.5 Add integration test for full initialization flow using `afero.MemMapFs`
- [ ] 8.6 Add integration test verifying InitResult accumulation

## 9. Final Verification

Ensure everything works end-to-end.

- [ ] 9.1 Run `go build ./...` to verify no compilation errors
- [ ] 9.2 Run `go test ./...` to verify all tests pass
- [ ] 9.3 Manual test: `spectr init` with Claude Code provider
- [ ] 9.4 Manual test: `spectr init` with multiple providers (verify deduplication)
- [ ] 9.5 Manual test: `spectr init` with Gemini provider (verify TOML format)
- [ ] 9.6 Manual test: Provider with global path (verify globalFs usage)
- [ ] 9.7 Verify InitResult reports correct created/updated files
- [ ] 9.8 Verify initializer ordering (directories created before files)
- [ ] 9.9 Update CLI help text for `spectr init` if needed
