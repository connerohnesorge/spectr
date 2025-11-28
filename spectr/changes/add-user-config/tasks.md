# Tasks: Add User Configuration with Color Theming

## 1. Core Configuration Infrastructure

- [ ] 1.1 Create `internal/config/` package with user config types
- [ ] 1.2 Implement XDG-compliant config file discovery (`$XDG_CONFIG_HOME/spectr/config.yaml` or `~/.config/spectr/config.yaml`)
- [ ] 1.3 Add YAML parsing for user config with sensible defaults
- [ ] 1.4 Write unit tests for config loading and defaults

## 2. Theme Color System

- [ ] 2.1 Define `Theme` struct with overridable color fields (accent, error, success, border, help, selected, highlight, header)
- [ ] 2.2 Set default theme values matching current hardcoded colors
- [ ] 2.3 Add color validation (hex format, ANSI 256 codes)
- [ ] 2.4 Write unit tests for theme parsing and validation

## 3. Style Refactoring

- [ ] 3.1 Update `internal/tui/styles.go` to use centralized theme config
- [ ] 3.2 Update `internal/view/formatters.go` to use centralized theme config
- [ ] 3.3 Create style factory functions that consume theme values
- [ ] 3.4 Write unit tests verifying style generation from config

## 4. Integration

- [ ] 4.1 Initialize config loader at application startup (in `cmd/` entry points)
- [ ] 4.2 Pass theme to TUI components via existing model structures
- [ ] 4.3 Add `spectr config` subcommand to display current config and path
- [ ] 4.4 Write integration tests for end-to-end config loading

## 5. Documentation & Examples

- [ ] 5.1 Add example config file with all color options documented
- [ ] 5.2 Update CLI help text to mention user config location
