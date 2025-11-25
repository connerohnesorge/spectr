## 1. Core Interface and Registry

- [x] 1.1 Create `Provider` interface in `internal/init/providers/provider.go`
- [x] 1.2 Create `Registry` struct with `Register()`, `Get()`, `All()` methods
- [x] 1.3 Create `CommandFormat` type (Markdown, TOML)
- [x] 1.4 Create `BaseProvider` struct with common logic for config file + slash commands
- [x] 1.5 Add unit tests for registry and base provider

## 2. Migrate Existing Providers

- [x] 2.1 Create `providers/claude.go` (CLAUDE.md + .claude/commands/)
- [x] 2.2 Create `providers/cline.go` (CLINE.md + .clinerules/commands/)
- [x] 2.3 Create `providers/cursor.go` (.cursorrules/commands/)
- [x] 2.4 Create `providers/copilot.go` (.github/copilot/commands/)
- [x] 2.5 Create `providers/aider.go` (.aider/commands/)
- [x] 2.6 Create `providers/continue.go` (.continue/commands/)
- [x] 2.7 Create `providers/mentat.go` (.mentat/commands/)
- [x] 2.8 Create `providers/tabnine.go` (.tabnine/commands/)
- [x] 2.9 Create `providers/smol.go` (.smol/commands/)
- [x] 2.10 Create `providers/costrict.go` (COSTRICT.md + .costrict/commands/)
- [x] 2.11 Create `providers/windsurf.go` (.windsurf/commands/)
- [x] 2.12 Create `providers/codebuddy.go` (CODEBUDDY.md + .codebuddy/commands/)
- [x] 2.13 Create `providers/qwen.go` (QWEN.md + .qwen/commands/)
- [x] 2.14 Create `providers/qoder.go` (QODER.md + .qoder/commands/)
- [x] 2.15 Create `providers/kilocode.go` (.kilocode/commands/)
- [x] 2.16 Create `providers/antigravity.go` (AGENTS.md + .agent/workflows/)

## 3. Add Gemini Provider

- [x] 3.1 Create `providers/gemini.go` with TOML format support
- [x] 3.2 Implement TOML command file generation in Configure()
- [x] 3.3 Add unit tests for Gemini TOML generation

## 4. Update Executor Integration

- [x] 4.1 Update `InitExecutor` to use new `Registry`
- [x] 4.2 Simplify `configureTools()` (no separate slash provider lookup needed)
- [x] 4.3 Update wizard to display providers from registry

## 5. Cleanup

- [x] 5.1 Remove `tool_definitions.go` global maps
- [x] 5.2 Remove old `ToolRegistry` struct
- [x] 5.3 Update all existing tests to use new API
- [x] 5.4 Run `go test ./...` to verify no regressions

## 6. Documentation

- [x] 6.1 Add doc comments to `Provider` interface
- [x] 6.2 Add "How to add a new provider" example in provider.go comments
