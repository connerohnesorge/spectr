# Tasks: Add Codex CLI Provider Support

## 1. Extend Path Resolution in BaseProvider

- [x] 1.1 Add `expandPath(path string) string` helper function in `providers/helpers.go`
- [x] 1.2 Add `isGlobalPath(path string) bool` helper to detect `~/` or `/` prefixed paths
- [x] 1.3 Modify `BaseProvider.configureSlashCommands()` to use expanded paths for global providers
- [x] 1.4 Modify `BaseProvider.IsConfigured()` to handle global paths correctly
- [x] 1.5 Add unit tests for `expandPath()` and `isGlobalPath()` helpers

## 2. Add Codex Provider Implementation

- [x] 2.1 Add `PriorityCodex = 10` constant to `constants.go`
- [x] 2.2 Create `internal/init/providers/codex.go` with CodexProvider struct
- [x] 2.3 Configure paths: `~/.codex/prompts/spectr/{proposal,sync,apply}.md`
- [x] 2.4 Use standard markdown frontmatter with `description:` field
- [x] 2.5 Register provider in `init()` function

## 3. Add Specification

- [ ] 3.1 Create `spectr/specs/support-codex/spec.md` after implementation verified
- [ ] 3.2 Archive this change

## 4. Testing

- [x] 4.1 Add unit tests for CodexProvider in `providers/codex_test.go`
- [x] 4.2 Test global path expansion on actual filesystem
- [x] 4.3 Verify `spectr init` with codex provider selection
- [x] 4.4 Verify `spectr init` updates global codex prompts
- [x] 4.5 Run full test suite: `go test ./...`
