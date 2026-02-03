# Change: Add Multi-Repo Project Nesting Support

## Why

Spectr currently requires running from the exact project root where `spectr/`
lives. In mono-repo setups with nested git repositories (each with their own
`spectr/` directory), users cannot run spectr from subdirectories or aggregate
results across multiple spectr roots. This limits usability for organizations
using mono-repo governance patterns.

**GitHub Issue**: #363

## What Changes

- **Discovery overhaul**: Walk up from cwd to find all `spectr/` directories,
  stopping at `.git` boundaries (each git repo is isolated)
- **Environment override**: Add `SPECTR_ROOT` env var for explicit root
  selection (single path, takes precedence over discovery)
- **Aggregated output**: Commands like `list`, `validate`, `view` aggregate
  results from all discovered roots, prefixing items with source path
  (e.g., `[project] add-feature`)
- **TUI path copying**: Change Enter key behavior to copy path relative to cwd
  instead of just the ID, enabling correct proposal application in nested
  contexts
- **No inheritance**: Each `spectr/` directory remains completely independent;
  no cross-project spec dependencies

## Impact

- Affected specs: `cli-interface` (discovery, list, view, TUI behavior)
- Affected code:
  - `internal/discovery/` - new multi-root discovery logic
  - `cmd/root.go` - env var handling, multi-root iteration
  - `cmd/list.go` - aggregated output with prefixes
  - `cmd/view.go` - aggregated dashboard
  - `cmd/validate.go` - validate across roots
  - `internal/tui/` - relative path copying
- **Backward compatible**: Single-root projects work exactly as before
