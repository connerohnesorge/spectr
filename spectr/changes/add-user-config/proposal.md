# Change: Add User Configuration with Color Theming

## Why

Users cannot customize the TUI appearance to match their terminal preferences or accessibility needs. All colors are currently hardcoded across multiple files (`internal/tui/styles.go`, `internal/view/formatters.go`), making personalization impossible without modifying source code.

## What Changes

- Add user configuration file support at `~/.config/spectr/config.yaml`
- Allow color overrides for key TUI elements (accent, error, success, border, help, selected, highlight)
- Configuration is optional - defaults to current hardcoded colors when no config exists
- Discovery follows XDG Base Directory spec (`$XDG_CONFIG_HOME/spectr/config.yaml` or `~/.config/spectr/config.yaml`)
- Centralize style definitions to consume configuration values

## Impact

- Affected specs: `cli-interface` (new user configuration capability)
- Affected code:
  - `internal/config/` - new package for user config discovery and loading
  - `internal/tui/styles.go` - refactor to consume config values
  - `internal/view/formatters.go` - refactor to consume config values
  - `internal/init/gradient.go` - optionally use config colors
  - All TUI components that reference color constants
