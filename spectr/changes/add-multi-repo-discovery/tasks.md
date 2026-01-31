# Tasks: Add Multi-Repo Project Nesting Support

## 1. Core Discovery Infrastructure

- [ ] 1.1 Create `SpectrRoot` type in `internal/discovery/roots.go` with fields:
  `Path` (absolute), `RelativeTo` (relative to cwd), `GitRoot` (parent .git dir)
- [ ] 1.2 Implement `FindSpectrRoots(cwd string) ([]SpectrRoot, error)` that walks
  up from cwd, finds all `spectr/` directories, stops at `.git` boundaries
- [ ] 1.3 Add git boundary detection: check for `.git` directory at each level,
  don't traverse beyond git root
- [ ] 1.4 Add `SPECTR_ROOT` env var support: if set, return single root from env
  var instead of discovery (validate path exists)
- [ ] 1.5 Add unit tests for `FindSpectrRoots` covering: single root, multiple
  roots, git boundary stopping, env var override, no roots found

## 2. Command Integration

- [ ] 2.1 Create helper `GetDiscoveredRoots()` in `cmd/` that wraps
  `FindSpectrRoots` with cwd and returns `[]SpectrRoot`
- [ ] 2.2 Update `cmd/root.go` `AfterApply()` to iterate over all discovered
  roots for sync operations
- [ ] 2.3 Update `cmd/list.go` to aggregate changes/specs from all roots,
  prefix each item with `[relative-path]`
- [ ] 2.4 Update `cmd/view.go` dashboard to show aggregated stats from all
  roots, with per-root breakdown
- [ ] 2.5 Update `cmd/validate.go` to validate items across all discovered
  roots

## 3. TUI Path Copying

- [ ] 3.1 Modify `internal/tui/list.go` (or equivalent) to track the full path
  of each item, not just ID
- [ ] 3.2 Change Enter key handler to copy path relative to cwd instead of ID
- [ ] 3.3 Update clipboard copy to use `filepath.Rel(cwd, itemPath)` for the
  copied value
- [ ] 3.4 Ensure path works for both changes (`spectr/changes/<id>/proposal.md`)
  and specs (`spectr/specs/<id>/spec.md`)

## 4. Output Formatting

- [ ] 4.1 Add `FormatItemWithRoot(root SpectrRoot, id string) string` helper
  that returns `[relative-path] id` format
- [ ] 4.2 Update list command table to include a "Root" or "Project" column
  when multiple roots detected
- [ ] 4.3 For single-root scenarios, omit the prefix to maintain current
  behavior (backward compatible)
- [ ] 4.4 Add tests for output formatting with single and multiple roots

## 5. Documentation and Polish

- [ ] 5.1 Update `spectr/AGENTS.md` to document multi-repo discovery behavior
- [ ] 5.2 Add `SPECTR_ROOT` env var documentation to help text and README
- [ ] 5.3 Add integration test: create nested git repos with spectr dirs,
  verify discovery and aggregation
