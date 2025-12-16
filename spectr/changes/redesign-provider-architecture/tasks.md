## 1. Foundation: Core Interfaces and Types

These tasks establish the new architecture without breaking the existing code.

- [ ] 1.1 Create `internal/initialize/providers/initializer.go` with new `Initializer` interface including `Init(ctx, fs, cfg, templateManager)` and `IsSetup(fs, cfg)` methods
- [ ] 1.2 Create `internal/initialize/providers/config.go` with `Config` struct containing `SpectrDir` field
- [ ] 1.3 Create `internal/initialize/providers/provider_new.go` with new minimal `Provider` interface returning `[]Initializer`
- [ ] 1.4 Create `internal/initialize/providers/registration.go` with `Registration` struct (ID, Name, Priority, Provider) and new registration API

## 2. Built-in Initializers

Create the three composable initializers that providers will use.

- [ ] 2.1 Create `internal/initialize/providers/initializers/directory.go` with `DirectoryInitializer` that accepts `afero.Fs` and creates directory paths
- [ ] 2.2 Create `internal/initialize/providers/initializers/configfile.go` with `ConfigFileInitializer` that accepts `afero.Fs` and `*TemplateManager` for marker-based instruction file updates
- [ ] 2.3 Create `internal/initialize/providers/initializers/slashcmds.go` with `SlashCommandsInitializer` supporting both Markdown and TOML formats, receiving `*TemplateManager` for template rendering
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

## 4. Git Diff Integration for Change Detection

Implement git-based file change detection to replace `GetFilePaths()` declarations.

- [ ] 4.1 Create `internal/initialize/git/detector.go` with `ChangeDetector` type
- [ ] 4.2 Implement `Snapshot() (string, error)` method to capture git working tree state (using `git stash create` or staging area comparison)
- [ ] 4.3 Implement `ChangedFiles(beforeSnapshot string) ([]string, error)` method using `git diff --name-only`
- [ ] 4.4 Handle edge cases: untracked files, not a git repo, dirty working tree
- [ ] 4.5 Add unit tests for `ChangeDetector` mocking git commands

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
- [ ] 5.12 Migrate `tabnine.go` to new Provider interface - delete old code
- [ ] 5.13 Migrate `windsurf.go` to new Provider interface - delete old code
- [ ] 5.14 Migrate `kilocode.go` to new Provider interface - delete old code
- [ ] 5.15 Migrate `continue.go` to new Provider interface - delete old code
- [ ] 5.16 Migrate `crush.go` to new Provider interface - delete old code
- [ ] 5.17 Migrate `opencode.go` to new Provider interface - delete old code

## 6. Executor Integration

Update executor to use new architecture with `afero.Fs` and initializer deduplication.

- [ ] 6.1 Update `executor.go` to create `afero.NewBasePathFs(osFs, projectPath)` filesystem
- [ ] 6.2 Update `executor.go` to use new registry API (`Registration` based retrieval)
- [ ] 6.3 Implement initializer collection from selected providers
- [ ] 6.4 Implement initializer deduplication logic based on type + config key
- [ ] 6.5 Update `configureProviders()` to run initializers with `Init(ctx, fs, cfg, templateManager)` signature
- [ ] 6.6 Integrate `git.ChangeDetector` for file change reporting (replace `GetFilePaths()` usage)
- [ ] 6.7 Update `ExecutionResult` to use git-detected changed files instead of declared paths

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
- [ ] 9.3 Manual test: `spectr init` with Claude Code provider
- [ ] 9.4 Manual test: `spectr init` with multiple providers (verify deduplication)
- [ ] 9.5 Manual test: `spectr init` with Gemini provider (verify TOML format)
- [ ] 9.6 Verify git diff shows expected file changes after init
- [ ] 9.7 Update CLI help text for `spectr init` if needed
