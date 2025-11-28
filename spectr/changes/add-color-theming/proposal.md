# Change: Add Preset Color Theming Support

## Why

Users want to customize the visual appearance of Spectr's TUI output to match their terminal preferences or improve accessibility. Currently, all colors are hardcoded across multiple packages (`internal/tui/styles.go`, `internal/view/formatters.go`, `internal/view/progress.go`, `internal/init/wizard.go`), making customization impossible without code changes.

## What Changes

- Extend `spectr.yaml` configuration to support a `theme` setting with preset theme names
- Create a centralized `internal/theme` package that defines color palettes for each preset
- Add 5 built-in themes: `default`, `dark`, `light`, `solarized`, `monokai`
- Refactor all TUI color usage to read from the active theme instead of hardcoded values
- Theme applies to: dashboard view, init wizard, interactive modes, progress bars, and table styling

## Impact

- Affected specs: `cli-framework` (new capability for theme configuration)
- Affected code:
  - `internal/config/config.go` - add `Theme` field to Config struct
  - `internal/theme/` (new) - theme definitions and color palette types
  - `internal/tui/styles.go` - use theme colors instead of constants
  - `internal/view/formatters.go` - use theme colors for dashboard styling
  - `internal/view/progress.go` - use theme colors for progress bars
  - `internal/init/wizard.go` - use theme colors and gradient endpoints
