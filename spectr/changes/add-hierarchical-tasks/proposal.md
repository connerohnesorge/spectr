# Change: Add Hierarchical tasks.jsonc Structure

## Why

Large change proposals (like `redesign-provider-architecture` with 60+ tasks) produce `tasks.jsonc` files that are too large for AI agents to read effectively in a single Read call. Additionally, there's currently no way to associate tasks with specific delta specs—all tasks live in a flat root file.

This change introduces hierarchical `tasks.jsonc` files that:

1. Allow delta specs (`spectr/changes/<id>/specs/<capability>/`) to have their own task files
2. Enable root tasks to reference child tasks in delta spec files
3. Provide a summary view in the root file so agents can drill down into specific capabilities
4. Auto-split `tasks.md` sections into capability-specific task files during `spectr accept`

## What Changes

- **ADDED**: Support for `tasks.jsonc` files inside delta spec directories (`specs/<capability>/tasks.jsonc`)
- **ADDED**: Task reference syntax allowing root tasks to point to child tasks (`"children": "$ref:specs/support-aider/tasks.jsonc"`)
- **ADDED**: Auto-discovery of delta spec task files via glob pattern
- **ADDED**: Summary counts in root `tasks.jsonc` for quick progress overview
- **MODIFIED**: `spectr accept` command to auto-split tasks by capability when section names match
- **MODIFIED**: Task ID schema to support hierarchical IDs (e.g., `5.1.1`, `5.1.2` for children under `5.1`)
- **ADDED**: `spectr tasks` command with `--flatten` flag to merge hierarchical tasks into single view

## Impact

- **Affected specs**: `cli-interface`
- **Affected code**:
  - `cmd/accept.go` - task splitting logic
  - `cmd/accept_writer.go` - hierarchical file generation
  - `internal/parsers/types.go` - new Task fields
  - New `cmd/tasks.go` - tasks command implementation
- **Breaking changes**: None - existing flat `tasks.jsonc` files remain valid
