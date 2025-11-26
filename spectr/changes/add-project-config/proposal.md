# Change: Add YAML Project Configuration

## Why

Users want to customize their Spectr installation, particularly renaming the root folder from `spectr/` to better fit their workflow (e.g., `specs/`, `.spectr/`, `openspec/`). Currently the `spectr` directory name is hardcoded throughout the codebase in ~50 locations, making customization impossible.

## What Changes

- Add optional `spectr.yaml` configuration file support at the project root
- Allow customization of the root directory name via `root_dir` setting
- Configuration is optional - defaults to `spectr/` when no config file exists
- All internal path resolution uses a centralized config loader
- Discovery walks up the directory tree to find `spectr.yaml` or falls back to `spectr/` directory

## Impact

- Affected specs: `cli-framework` (new capability for configuration)
- Affected code:
  - `internal/discovery/` - add config discovery and loading
  - `internal/validation/helpers.go` - use config instead of hardcoded constant
  - `internal/init/` - optionally generate config file
  - `internal/archive/` - use config for paths
  - `internal/list/` - use config for paths
  - `internal/view/` - use config for paths
  - All commands that reference spectr paths
