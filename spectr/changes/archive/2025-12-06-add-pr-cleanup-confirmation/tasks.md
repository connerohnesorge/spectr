## 1. Implementation

- [x] 1.1 Create Bubbletea TUI confirmation menu component in `internal/tui/confirm.go` (reusable Yes/No menu with arrow key navigation)
- [x] 1.2 Create `removeChangeDirectory` helper function in `internal/pr/helpers.go`
- [x] 1.3 Add `Yes` flag to `PRProposalCmd` struct for non-interactive mode (skip prompt, keep change)
- [x] 1.4 Add cleanup confirmation flow to `cmd/pr.go` after successful PR creation (proposal mode only)
- [x] 1.5 Wire up the TUI confirmation menu with styled rendering matching spectr's existing TUI components

## 2. Testing

- [x] 2.1 Add unit tests for `removeChangeDirectory` helper
- [x] 2.2 Add unit tests for TUI confirmation menu component
- [x] 2.3 Test non-interactive mode (--yes flag skips prompt, keeps change)

## 3. Documentation

- [x] 3.1 Update CLI help text for `spectr pr proposal` to document `--yes` flag and cleanup prompt behavior
