# Implementation Tasks

## 1. Remove Slash Command File

- [x] 1.1 Delete `.claude/commands/spectr/sync.md`

## 2. Remove Documentation

- [x] 2.1 Delete `WHY_SYNC.md` from project root

## 3. Update Templates

- [x] 3.1 Delete `internal/initialize/templates/tools/slash-sync.md.tmpl`
- [x] 3.2 Update `internal/initialize/templates/spectr/AGENTS.md.tmpl` - remove
  "Stage 3: Syncing Specs" section and sync-related references

## 4. Update Provider Interface and Constants

- [x] 4.1 Remove `GetSyncCommandPath()` from provider interface in
  `internal/initialize/providers/provider.go`
- [x] 4.2 Remove `syncPath` field from `baseProvider` struct
- [x] 4.3 Remove `FrontmatterSync` constant from
  `internal/initialize/providers/constants.go`
- [x] 4.4 Update `StandardCommandPaths()` to return only proposal and apply
  paths
- [x] 4.5 Update `PrefixedCommandPaths()` to return only proposal and apply
  paths
- [x] 4.6 Update `CommandFrontmatter` map to remove sync entry

## 5. Update All Provider Implementations

- [x] 5.1 Update `internal/initialize/providers/claude.go`
- [x] 5.2 Update `internal/initialize/providers/aider.go`
- [x] 5.3 Update `internal/initialize/providers/antigravity.go`
- [x] 5.4 Update `internal/initialize/providers/cline.go`
- [x] 5.5 Update `internal/initialize/providers/codebuddy.go`
- [x] 5.6 Update `internal/initialize/providers/codex.go`
- [x] 5.7 Update `internal/initialize/providers/continue.go`
- [x] 5.8 Update `internal/initialize/providers/costrict.go`
- [x] 5.9 Update `internal/initialize/providers/cursor.go`
- [x] 5.10 Update `internal/initialize/providers/gemini.go`
- [x] 5.11 Update `internal/initialize/providers/kilocode.go`
- [x] 5.12 Update `internal/initialize/providers/qoder.go`
- [x] 5.13 Update `internal/initialize/providers/qwen.go`
- [x] 5.14 Update `internal/initialize/providers/tabnine.go`
- [x] 5.15 Update `internal/initialize/providers/windsurf.go`

## 6. Update Tests

- [x] 6.1 Update `internal/initialize/providers/provider_test.go` - remove
  sync-related tests
- [x] 6.2 Update `internal/initialize/templates_test.go` - remove sync command
  tests
- [x] 6.3 Update `internal/initialize/wizard_test.go` - remove sync.md
  expectations

## 7. Validation

- [x] 7.1 Run `go build ./...` to verify compilation
- [x] 7.2 Run `go test ./...` to verify tests pass
- [x] 7.3 Run `spectr validate remove-sync-slash-command --strict` to verify
  proposal
