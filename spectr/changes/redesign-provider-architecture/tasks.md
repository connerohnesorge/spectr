## 1. Foundation: Core Interfaces and Types

These tasks establish the new architecture without breaking the existing code.

- [ ] 1.1 Create `internal/initialize/providers/initializer.go` with new `Initializer` interface including:
  - `Init(ctx context.Context, fs afero.Fs, cfg *Config, tm *TemplateManager) error`
  - `IsSetup(fs afero.Fs, cfg *Config) bool`
  - `Path() string` (for deduplication)
  - `IsGlobal() bool` (true = use globalFs, false = use projectFs)
- [ ] 1.2 Create `internal/initialize/providers/config.go` with `Config` struct containing:
  - `SpectrDir string` field
  - `SpecsDir() string` method (returns SpectrDir + "/specs")
  - `ChangesDir() string` method (returns SpectrDir + "/changes")
  - `ProjectFile() string` method (returns SpectrDir + "/project.md")
  - `AgentsFile() string` method (returns SpectrDir + "/AGENTS.md")
- [ ] 1.3 Create `internal/initialize/providers/provider_new.go` with new minimal `Provider` interface returning `[]Initializer`
- [ ] 1.4 Create `internal/initialize/providers/registration.go` with `Registration` struct (ID, Name, Priority, Provider) and new registration API

## 2. Built-in Initializers

Create the three composable initializers that providers will use. Each must implement `Path()` and `IsGlobal()`.

- [ ] 2.1 Create `internal/initialize/providers/initializers/directory.go` with `DirectoryInitializer`:
  - Implements `Init()`, `IsSetup()`, `Path()`, `IsGlobal()`
  - Accepts directory path(s) and isGlobal flag
  - Creates directories with `MkdirAll`
- [ ] 2.2 Create `internal/initialize/providers/initializers/configfile.go` with `ConfigFileInitializer`:
  - Implements `Init()`, `IsSetup()`, `Path()`, `IsGlobal()`
  - Uses `*TemplateManager` for marker-based instruction file updates
  - Handles both create and update scenarios
- [ ] 2.3 Create `internal/initialize/providers/initializers/slashcmds.go` with `SlashCommandsInitializer`:
  - Implements `Init()`, `IsSetup()`, `Path()`, `IsGlobal()`
  - Supports both Markdown and TOML formats
  - Receives `*TemplateManager` for template rendering
- [ ] 2.4 Add unit tests for `DirectoryInitializer` with `afero.MemMapFs`
- [ ] 2.5 Add unit tests for `ConfigFileInitializer` with `afero.MemMapFs` - test create and marker update scenarios
- [ ] 2.6 Add unit tests for `SlashCommandsInitializer` with `afero.MemMapFs` - test Markdown and TOML formats

## 3. New Registry Implementation

Replace the old registry with metadata-separated registration.

- [ ] 3.1 Create `internal/initialize/providers/registry_v2.go` with new `Registry` type using `Registration` struct for metadata
- [ ] 3.2 Implement `Register(Registration)`, `Get(id)`, `All()`, `IDs()`, `Count()` methods on new registry
- [ ] 3.3 Add priority-sorted retrieval maintaining backwards-compatible behavior
- [ ] 3.4 Add duplicate ID rejection with clear error messages
- [ ] 3.5 Add unit tests for new registry: registration, retrieval, priority sorting, duplicate rejection

## 4. Git Integration for Change Detection

Implement git-based file change detection to replace `GetFilePaths()` declarations. Require git repository.

- [ ] 4.1 Create `internal/initialize/git/detector.go` with `ChangeDetector` type
- [ ] 4.2 Implement `IsGitRepo(path string) bool` function for early validation
- [ ] 4.3 Implement `Snapshot() (string, error)` method to capture git working tree state (using `git stash create` or staging area comparison)
- [ ] 4.4 Implement `ChangedFiles(beforeSnapshot string) ([]string, error)` method using `git diff --name-only`
- [ ] 4.5 Handle edge cases: untracked files (use `git status --porcelain`), dirty working tree
- [ ] 4.6 Add clear error message for non-git repos: "spectr init requires a git repository. Run 'git init' first."
- [ ] 4.7 Add unit tests for `ChangeDetector` mocking git commands

## 5. Migrate Providers (In-Place Replacement)

Migrate each provider to the new interface, deleting old code as each migration completes.

- [ ] 5.1 Migrate `claude.go` to new Provider interface (reference implementation) - delete old `ClaudeProvider` struct and `NewClaudeProvider()` after migration
- [ ] 5.2 Migrate `gemini.go` to new Provider interface (TOML format example) - delete old `GeminiProvider` struct and `Configure()` override after migration
- [ ] 5.3 Migrate `cursor.go` to new Provider interface - delete old code
- [ ] 5.4 Migrate `cline.go` to new Provider interface - delete old code
- [ ] 5.5 Migrate `aider.go` to new Provider interface - delete old code
- [ ] 5.6 Migrate `codex.go` to new Provider interface - delete old code
- [ ] 5.7 Migrate `costrict.go` to new Provider interface - delete old code
- [ ] 5.8 Migrate `qoder.go` to new Provider interface - delete old code
- [ ] 5.9 Migrate `codebuddy.go` to new Provider interface - delete old code
- [ ] 5.10 Migrate `qwen.go` to new Provider interface - delete old code
- [ ] 5.11 Migrate `antigravity.go` to new Provider interface - delete old code
- [ ] 5.13 Migrate `windsurf.go` to new Provider interface - delete old code
- [ ] 5.14 Migrate `kilocode.go` to new Provider interface - delete old code
- [ ] 5.15 Migrate `continue.go` to new Provider interface - delete old code
- [ ] 5.16 Migrate `crush.go` to new Provider interface - delete old code
- [ ] 5.17 Migrate `opencode.go` to new Provider interface - delete old code

## 6. Executor Integration

Update executor to use new architecture with dual filesystem, ordering, and deduplication.

- [ ] 6.1 Create dual filesystem in `executor.go`:
  - `projectFs := afero.NewBasePathFs(osFs, projectPath)` for project-relative paths
  - `globalFs := afero.NewBasePathFs(osFs, os.UserHomeDir())` for global paths
- [ ] 6.2 Add git repo check at start of `Init()`:
  - Call `git.IsGitRepo(projectPath)`
  - Return clear error if not a git repo
- [ ] 6.3 Update `executor.go` to use new registry API (`Registration` based retrieval)
- [ ] 6.4 Implement initializer collection from selected providers
- [ ] 6.5 Implement initializer deduplication by `Path()` - same path = run once
- [ ] 6.6 Implement initializer sorting by type (guaranteed order):
  - 1. `DirectoryInitializer`
  - 2. `ConfigFileInitializer`
  - 3. `SlashCommandsInitializer`
- [ ] 6.7 Update `configureProviders()` to:
  - Select `projectFs` or `globalFs` based on `initializer.IsGlobal()`
  - Call `Init(ctx, fs, cfg, templateManager)` on each initializer
- [ ] 6.8 Integrate `git.ChangeDetector` for file change reporting (replace `GetFilePaths()` usage)
- [ ] 6.9 Update `ExecutionResult` to use git-detected changed files instead of declared paths
- [ ] 6.10 Handle partial failures: report which initializers failed, continue with rest

## 7. Remove Old Provider Code (In-Place Cleanup)

Remove deprecated code that has been replaced by the new architecture.

- [ ] 7.1 Remove old `Provider` interface (12-method version) from `provider.go`
- [ ] 7.2 Remove `BaseProvider` struct from `provider.go`
- [ ] 7.3 Remove `TemplateRenderer` interface from `provider.go` (now using `*TemplateManager` directly)
- [ ] 7.4 Remove `TemplateContext` and `DefaultTemplateContext()` from `provider.go` (move to config if needed)
- [ ] 7.5 Delete `helpers.go` - `EnsureDir`, `FileExists`, `UpdateFileWithMarkers` now in initializers or use `afero.Fs`
- [ ] 7.6 Remove old global registry functions from `registry.go` (keep only new `Registry` type)
- [ ] 7.7 Clean up `constants.go` - remove `StandardFrontmatter()`, `StandardCommandPaths()`, `PrefixedCommandPaths()` (moved to initializers)
- [ ] 7.8 Remove priority constants from `constants.go` (priorities now in registration calls)

## 8. Test Cleanup

Update tests to match new architecture - explicit cleanup tasks for old test files.

- [ ] 8.1 Delete `provider_test.go` (~666 lines) - old interface tests no longer applicable
- [ ] 8.2 Create `provider_new_test.go` with tests for new `Provider` interface
- [ ] 8.3 Add tests verifying all 17 providers return expected initializers
- [ ] 8.4 Add tests verifying provider registration metadata (ID, Name, Priority)
- [ ] 8.5 Delete `registry_test.go` (~212 lines) - old registry tests no longer applicable
- [ ] 8.6 Create `registry_v2_test.go` with tests for new `Registry` type
- [ ] 8.7 Add integration test for full initialization flow using `afero.MemMapFs`
- [ ] 8.8 Add integration test verifying git diff change detection

## 9. Final Verification

Ensure everything works end-to-end.

- [ ] 9.1 Run `go build ./...` to verify no compilation errors
- [ ] 9.2 Run `go test ./...` to verify all tests pass
- [ ] 9.3 Manual test: `spectr init` in non-git directory (verify early fail with clear error)
- [ ] 9.4 Manual test: `spectr init` with Claude Code provider
- [ ] 9.5 Manual test: `spectr init` with multiple providers (verify deduplication)
- [ ] 9.6 Manual test: `spectr init` with Gemini provider (verify TOML format)
- [ ] 9.7 Manual test: Provider with global path (verify globalFs usage)
- [ ] 9.8 Verify git diff shows expected file changes after init
- [ ] 9.9 Verify initializer ordering (directories created before files)
- [ ] 9.10 Update CLI help text for `spectr init` if needed
