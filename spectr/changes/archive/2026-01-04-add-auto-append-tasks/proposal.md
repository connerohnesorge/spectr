# Change: Add Auto-Append Tasks on Accept

## Why

Teams often have project-specific workflow tasks that should be tracked on every
change proposal (e.g., "Update changelog", "Run linter and tests", "Notify
stakeholders"). Currently, these must be manually added to each `tasks.md` file,
which is error-prone and inconsistent.

This change introduces a project-level configuration file (`spectr.yaml`) that
allows teams to define tasks that are automatically appended during
`spectr accept`, ensuring consistent project workflows across all changes.

## What Changes

- **ADDED**: Support for `spectr.yaml` configuration file at project root
- **ADDED**: `append_tasks` configuration section with:
  - `section`: Configurable section name for appended tasks (e.g., "Project
    Workflow")
  - `tasks`: List of task descriptions to append
- **MODIFIED**: `spectr accept` command to read config and append defined tasks
  to `tasks.jsonc` output
- **ADDED**: Config file discovery and parsing (`internal/config/` package)

## Impact

- **Affected specs**: `cli-interface`
- **Affected code**:
  - New `internal/config/config.go` - YAML config parsing
  - `cmd/accept.go` - integrate config loading and task appending
  - `cmd/accept_writer.go` - modify JSONC generation to include appended tasks
- **Breaking changes**: None - `spectr.yaml` is optional; existing behavior
  unchanged when config is absent
