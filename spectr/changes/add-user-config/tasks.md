# Tasks: Add User Configuration with Color Theming

## 1. Core Configuration Infrastructure

- [x] 1.1 Create `internal/config/` package with user config types
- [x] 1.2 Implement XDG-compliant config file discovery (`$XDG_CONFIG_HOME/spectr/config.yaml` or `~/.config/spectr/config.yaml`)
- [x] 1.3 Add YAML parsing for user config with sensible defaults
- [x] 1.4 Write unit tests for config loading and defaults

## 2. Theme Color System

- [x] 2.1 Define `Theme` struct with overridable color fields (accent, error, success, border, help, selected, highlight, header)
- [x] 2.2 Set default theme values matching current hardcoded colors
- [x] 2.3 Add color validation (hex format, ANSI 256 codes)
- [x] 2.4 Write unit tests for theme parsing and validation

## 3. Style Refactoring

- [x] 3.1 Update `internal/tui/styles.go` to use centralized theme config
- [x] 3.2 Update `internal/view/formatters.go` to use centralized theme config
- [x] 3.3 Create style factory functions that consume theme values
- [x] 3.4 Write unit tests verifying style generation from config

## 4. Integration

- [x] 4.1 Initialize config loader at application startup (in `cmd/` entry points)
- [x] 4.2 Pass theme to TUI components via existing model structures
- [x] 4.3 Add `spectr config` subcommand to display current config and path
- [x] 4.4 Write integration tests for end-to-end config loading

## 5. Documentation & Examples

- [x] 5.1 Add example config file with all color options documented
- [x] 5.2 Update CLI help text to mention user config location
