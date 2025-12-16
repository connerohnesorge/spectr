## 1. Core Interfaces and Types

- [ ] 1.1 Create `internal/initialize/providers/interfaces.go` with new `Provider`, `Initializer`, `Config`, and `Registration` types
- [ ] 1.2 Create `internal/initialize/providers/registry_new.go` with new registration API using `Registration` struct
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

## 5. Cleanup

- [ ] 5.1 Remove old `Provider` interface and `BaseProvider` from `provider.go`
- [ ] 5.2 Remove old `TemplateRenderer` interface
- [ ] 5.3 Remove `helpers.go` functions that are now in initializers
- [ ] 5.4 Remove old registry functions (keep new ones)
- [ ] 5.5 Update `constants.go` to remove unused constants

## 6. Documentation and Testing

- [ ] 6.1 Update CLI help text for `spectr init` if needed
- [ ] 6.2 Run full test suite and fix any failures
- [ ] 6.3 Manual testing: `spectr init` with various provider combinations
- [ ] 6.4 Verify git diff shows expected file changes after init
