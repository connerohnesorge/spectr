## 1. Theme Package Foundation

- [x] 1.1 Create `internal/theme/` package with `Theme` struct defining all color slots
- [x] 1.2 Define color slot interface: `Primary`, `Secondary`, `Success`, `Error`, `Warning`, `Muted`, `Border`, `Header`, `Selected`, `Highlight`, `GradientStart`, `GradientEnd`
- [x] 1.3 Implement `default` theme matching current hardcoded colors
- [x] 1.4 Implement `dark` theme optimized for dark terminal backgrounds
- [x] 1.5 Implement `light` theme optimized for light terminal backgrounds
- [x] 1.6 Implement `solarized` theme matching Solarized color palette
- [x] 1.7 Implement `monokai` theme matching Monokai color palette
- [x] 1.8 Add `Get(name string) (*Theme, error)` function for theme lookup
- [x] 1.9 Write unit tests for theme package

## 2. Configuration Integration

- [x] 2.1 Add `Theme string` field to `Config` struct in `internal/config/config.go`
- [x] 2.2 Default `Theme` to `"default"` when not specified
- [x] 2.3 Add validation for theme name (must match a known preset)
- [x] 2.4 Update config tests for theme field

## 3. TUI Styles Migration

- [x] 3.1 Update `internal/tui/styles.go` to accept theme parameter
- [x] 3.2 Replace hardcoded `ColorBorder`, `ColorHeader`, `ColorSelected`, `ColorHighlight`, `ColorHelp` with theme values
- [x] 3.3 Update `ApplyTableStyles`, `TitleStyle`, `HelpStyle`, `SelectedStyle`, `ChoiceStyle` to use theme
- [x] 3.4 Update tui tests

## 4. Dashboard View Migration

- [x] 4.1 Update `internal/view/formatters.go` to use theme colors
- [x] 4.2 Replace hardcoded header, active, completed, spec indicator colors
- [x] 4.3 Update `internal/view/progress.go` to use theme colors for filled/empty bar portions
- [x] 4.4 Update view tests

## 5. Init Wizard Migration

- [x] 5.1 Update `internal/init/wizard.go` to use theme colors
- [x] 5.2 Replace hardcoded `titleStyle`, `selectedStyle`, `dimmedStyle`, `cursorStyle`, `errorStyle`, `successStyle`, `infoStyle`, `subtleStyle`
- [x] 5.3 Update `applyGradient` call to use theme's `GradientStart` and `GradientEnd`
- [x] 5.4 Update init wizard tests

## 6. Integration Testing

- [x] 6.1 Add integration tests verifying theme changes apply to all TUI output
- [x] 6.2 Test backward compatibility (default theme when not specified)
- [x] 6.3 Test invalid theme name produces helpful error message

## 7. Documentation

- [x] 7.1 Update `spectr.yaml` example in README with `theme` option
- [x] 7.2 Document available preset themes and their characteristics
