# Tasks: Provider Architecture Redesign

## Overview

This document defines the implementation tasks for redesigning the provider architecture with composable initializers. Tasks are organized by phase to ensure proper dependency ordering.

**Key Principles:**
- Zero technical debt: Complete removal of old code, no compatibility shims
- Clean break: No backwards compatibility, users must re-run `spectr init`
- 15 providers total (compacted priorities 1-15)

---

## Phase 0: Domain Package - Shared Types

Create `internal/domain` package to break import cycles with shared types.

- [ ] 0.1 Create `internal/domain/template.go` with `TemplateRef` struct (public fields `Name`, `Template`), `Render()` method, `TemplateContext` struct, and `DefaultTemplateContext()` function

- [ ] 0.2 Create `internal/domain/slashcmd.go` with `SlashCommand` typed int, `SlashProposal`/`SlashApply` constants, and `String()` method (NO `TemplateName()`)

- [ ] 0.3 Create `internal/domain/template_test.go` with tests for TemplateRef.Render() and DefaultTemplateContext()

- [ ] 0.4 Create `internal/domain/slashcmd_test.go` with tests for SlashCommand.String()

- [ ] 0.5 Create `internal/domain/templates.go` with `//go:embed templates/*.tmpl` and exported `TemplateFS embed.FS`

- [ ] 0.6 Create `internal/domain/templates/` directory for embedded slash command templates

- [ ] 0.7 Move `internal/initialize/templates/tools/slash-proposal.md.tmpl` to `internal/domain/templates/slash-proposal.md.tmpl`

- [ ] 0.8 Move `internal/initialize/templates/tools/slash-apply.md.tmpl` to `internal/domain/templates/slash-apply.md.tmpl`

- [ ] 0.9 Verify `internal/domain/slashcmd.go` has ONLY `String()` method, no `TemplateName()`

- [ ] 0.10 Delete empty `internal/initialize/templates/tools/` directory after moving templates

---

## Phase 1: Foundation - Core Interfaces and Types

Create the new provider system types and interfaces.

- [ ] 1.1 Create `internal/initialize/providers/initializer.go` with `Initializer` interface: `Init(ctx, projectFs, globalFs, cfg, tm) (InitResult, error)` and `IsSetup(projectFs, globalFs, cfg) bool`

- [ ] 1.2 Create `internal/initialize/providers/result.go` with `InitResult` struct (`CreatedFiles`, `UpdatedFiles`), `ExecutionResult` struct (adds `Errors`), and `aggregateResults()` function

- [ ] 1.3 Create `internal/initialize/providers/config.go` with `Config` struct (`SpectrDir`) and derived path methods: `SpecsDir()`, `ChangesDir()`, `ProjectFile()`, `AgentsFile()`

- [ ] 1.4 Rewrite `internal/initialize/providers/provider.go` with minimal `Provider` interface returning `[]Initializer`; DELETE old 12-method interface, BaseProvider, TemplateRenderer, old TemplateContext

- [ ] 1.5 Create `internal/initialize/providers/registration.go` with `Registration` struct: `ID`, `Name`, `Priority`, `Provider`

---

## Phase 2: Type-Safe Template System

Update TemplateManager to use domain types and provide type-safe accessors.

- [ ] 2.1 Update `internal/initialize/templates.go` to import `internal/domain`, merge templates from `templateFS` and `domain.TemplateFS` in `NewTemplateManager()`

- [ ] 2.2 Add type-safe accessor methods to TemplateManager: `InstructionPointer() domain.TemplateRef`, `Agents() domain.TemplateRef`

- [ ] 2.3 Add `SlashCommand(cmd domain.SlashCommand) domain.TemplateRef` method to TemplateManager

- [ ] 2.4 Add unit tests in `internal/initialize/templates_test.go` verifying all accessors return valid `domain.TemplateRef` and Render() works

- [ ] 2.5 Update any existing code using TemplateContext to use `domain.TemplateContext`

---

## Phase 3: Built-in Initializers

Create the three reusable initializer implementations.

- [ ] 3.1 Create `internal/initialize/providers/initializers/directory.go` with `DirectoryInitializer`: creates directories, internal `global` flag, implements `dedupeKey()`

- [ ] 3.2 Create `internal/initialize/providers/initializers/configfile.go` with `ConfigFileInitializer`: marker-based updates, orphaned marker handling with `strings.LastIndex`, prevents duplicate blocks

- [ ] 3.3 Create `internal/initialize/providers/initializers/slashcmds.go` with `SlashCommandsInitializer`: creates proposal/apply files, supports Markdown and TOML formats, internal `global` flag

- [ ] 3.4 Add unit tests for `DirectoryInitializer` in `directory_test.go` using `afero.MemMapFs`: creates dirs, IsSetup checks, project and global filesystem modes

- [ ] 3.5 Add unit tests for `ConfigFileInitializer` in `configfile_test.go` using `afero.MemMapFs`: new file, update between markers, orphaned start with trailing end, orphaned start with no end, no duplicate blocks

- [ ] 3.6 Add unit tests for `SlashCommandsInitializer` in `slashcmds_test.go` using `afero.MemMapFs`: Markdown format, TOML format, IsSetup checks

---

## Phase 4: New Registry Implementation

Replace old registry with explicit registration (no init()).

- [ ] 4.1 Rewrite `internal/initialize/providers/registry.go` with new Registry using `Registration` struct, package-level `registry` map

- [ ] 4.2 Implement `RegisterProvider(reg Registration) error` with validation: non-empty ID, non-nil Provider, reject duplicates with error (not panic)

- [ ] 4.3 Implement registry query methods: `RegisteredProviders() []Registration`, `Get(id) (Registration, bool)`, `Count() int`

- [ ] 4.4 Ensure `RegisteredProviders()` returns providers sorted by Priority (lower first)

- [ ] 4.5 Ensure `RegisterProvider` returns clear error for duplicate IDs: `fmt.Errorf("provider %q already registered", reg.ID)`

- [ ] 4.6 Create `RegisterAllProviders() error` that explicitly registers all 15 providers with priorities 1-15: claude-code, gemini, costrict, qoder, qwen, antigravity, cline, cursor, codex, aider, windsurf, kilocode, continue, crush, opencode

- [ ] 4.7 Rewrite `registry_test.go` with tests for: valid registration, empty ID rejection, nil Provider rejection, duplicate ID rejection, priority sorting, Get() correctness

- [ ] 4.8 Add tests for `RegisterAllProviders()`: all 15 providers registered, no errors, priorities sequential 1-15, IDs correct

---

## Phase 5: Migrate Providers

Migrate all 15 providers to new interface. Each migration DELETEs old BaseProvider, Configure(), init() and implements new Provider returning []Initializer.

- [ ] 5.1 Migrate `claude.go` (Priority 1): Config file `CLAUDE.md`, commands `.claude/commands/spectr/`

- [ ] 5.2 Migrate `gemini.go` (Priority 2, TOML format): No config file, commands `.gemini/commands/spectr/` with TOML

- [ ] 5.3 Migrate `costrict.go` (Priority 3): Config file `COSTRICT.md`, commands `.costrict/commands/spectr/`

- [ ] 5.4 Migrate `qoder.go` (Priority 4): Config file `QODER.md`, commands `.qoder/commands/spectr/`

- [ ] 5.5 Migrate `qwen.go` (Priority 5): Config file `QWEN.md`, commands `.qwen/commands/spectr/`

- [ ] 5.6 Migrate `antigravity.go` (Priority 6, non-standard paths): Config file `AGENTS.md`, commands `.agent/workflows/` with prefixed names `spectr-proposal.md`, `spectr-apply.md`

- [ ] 5.7 Migrate `cline.go` (Priority 7): Config file `CLINE.md`, commands `.clinerules/commands/spectr/`

- [ ] 5.8 Migrate `cursor.go` (Priority 8): No config file, commands `.cursorrules/commands/spectr/`

- [ ] 5.9 Migrate `codex.go` (Priority 9, global paths): Config file `AGENTS.md`, commands `~/.codex/prompts/` with prefixed names, uses globalFs

- [ ] 5.10 Migrate `aider.go` (Priority 10): No config file, commands `.aider/commands/spectr/`

- [ ] 5.11 Migrate `windsurf.go` (Priority 11): No config file, commands `.windsurf/commands/spectr/`

- [ ] 5.12 Migrate `kilocode.go` (Priority 12): No config file, commands `.kilocode/commands/spectr/`

- [ ] 5.13 Migrate `continue.go` (Priority 13): No config file, commands `.continue/commands/spectr/`

- [ ] 5.14 Migrate `crush.go` (Priority 14): Config file `CRUSH.md`, commands `.crush/commands/spectr/`

- [ ] 5.15 Migrate `opencode.go` (Priority 15): No config file, commands `.opencode/commands/spectr/`

---

## Phase 6: Executor Integration

Update the executor to use the new provider system.

- [ ] 6.1 Update `cmd/root.go` or entry point to call `providers.RegisterAllProviders()` with error handling at startup

- [ ] 6.2 Create dual filesystem in `executor.go`: `projectFs` rooted at project, `globalFs` rooted at home directory

- [ ] 6.3 Update `executor.go` to use `RegisteredProviders()` for sorted provider list

- [ ] 6.4 Implement initializer collection from selected providers in `executor.go`

- [ ] 6.5 Implement initializer deduplication by `dedupeKey()` in `executor.go`

- [ ] 6.6 Implement initializer sorting by type in `executor.go`: DirectoryInitializer (1), ConfigFileInitializer (2), SlashCommandsInitializer (3)

- [ ] 6.7 Update `configureProviders()` to pass both filesystems, collect InitResult, handle errors without stopping

- [ ] 6.8 Use `aggregateResults(allResults, errors)` to combine all InitResult into ExecutionResult

- [ ] 6.9 Handle partial failures: log errors, continue with remaining initializers, include all errors in ExecutionResult.Errors

---

## Phase 7: Remove Old Code

Complete removal of deprecated code. Zero technical debt.

- [ ] 7.1 Remove old 12-method Provider interface from `provider.go` (already replaced in 1.4)

- [ ] 7.2 Remove `BaseProvider` struct and all its methods from `provider.go`

- [ ] 7.3 Remove `TemplateRenderer` interface from `provider.go`

- [ ] 7.4 Remove old `TemplateContext` and `DefaultTemplateContext()` from `provider.go` (now in domain)

- [ ] 7.5 DELETE entire `helpers.go` file: `EnsureDir`, `FileExists`, `UpdateFileWithMarkers` now in initializers

- [ ] 7.6 COMPLETELY DELETE old registry functions: `Register(p Provider)`, old Get/All/IDs/Count, `WithConfigFile()`, `WithSlashCommands()`, `Reset()` - NO compatibility shims

- [ ] 7.7 Clean up `constants.go`: DELETE `StandardFrontmatter()`, `StandardCommandPaths()`, `PrefixedCommandPaths()`

- [ ] 7.8 Remove priority constants from `constants.go` (priorities now in RegisterAllProviders)

---

## Phase 8: Test Cleanup

Update tests to match new architecture.

- [ ] 8.1 Rewrite `provider_test.go` with tests for new Provider interface: each provider returns expected initializers

- [ ] 8.2 Test all 15 providers return expected initializers: correct count, correct types, correct paths

- [ ] 8.3 Test provider registration metadata: ID, Name, Priority (1-15) all correct

- [ ] 8.4 Ensure `registry_test.go` covers new Registry (already done in 4.7/4.8)

- [ ] 8.5 Add integration test in `executor_test.go` for full initialization flow using `afero.MemMapFs`

- [ ] 8.6 Add integration test for InitResult accumulation: CreatedFiles, UpdatedFiles, Errors all correct

---

## Phase 9: Final Verification

Comprehensive testing before completion.

- [ ] 9.1 Run `go build ./...` to verify no compilation errors

- [ ] 9.2 Run `go test ./...` to verify all tests pass

- [ ] 9.3 Manual test: `spectr init` with Claude Code provider - verify CLAUDE.md and .claude/commands/spectr/ created

- [ ] 9.4 Manual test: `spectr init` with multiple providers - verify deduplication works, no duplicate marker blocks

- [ ] 9.5 Manual test: `spectr init` with Gemini provider - verify TOML format in .gemini/commands/spectr/

- [ ] 9.6 Manual test: `spectr init` with Codex provider - verify global paths in ~/.codex/prompts/

- [ ] 9.7 Verify InitResult reports correct created/updated files

- [ ] 9.8 Verify initializer ordering: directories first, then config files, then slash commands

- [ ] 9.9 Update CLI help text for `spectr init` if needed

---

## Summary

| Phase | Tasks | Description |
|-------|-------|-------------|
| 0 | 0.1-0.10 | Domain package with shared types |
| 1 | 1.1-1.5 | Core interfaces and types |
| 2 | 2.1-2.5 | Type-safe template system |
| 3 | 3.1-3.6 | Built-in initializers |
| 4 | 4.1-4.8 | New registry (no init()) |
| 5 | 5.1-5.15 | Migrate all 15 providers |
| 6 | 6.1-6.9 | Executor integration |
| 7 | 7.1-7.8 | Remove old code |
| 8 | 8.1-8.6 | Test cleanup |
| 9 | 9.1-9.9 | Final verification |

**Total: 68 tasks**
