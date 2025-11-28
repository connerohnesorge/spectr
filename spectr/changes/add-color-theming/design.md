## Context

Spectr's TUI output currently uses hardcoded ANSI color codes spread across four packages. Users cannot customize colors to match their terminal preferences or accessibility needs without modifying source code. The recently added `spectr.yaml` configuration system provides an ideal extension point for theme configuration.

**Current color locations:**
- `internal/tui/styles.go`: ANSI 256 codes (`240`, `99`, `229`, `57`)
- `internal/view/formatters.go`: Basic ANSI codes (`6`, `3`, `2`, `4`, `240`)
- `internal/view/progress.go`: Basic ANSI codes (`2`, `240`)
- `internal/init/wizard.go`: ANSI 256 codes (`99`, `170`, `240`, `212`, `196`, `42`, `86`, `241`, `205`)

## Goals / Non-Goals

**Goals:**
- Users can select from preset themes via `spectr.yaml`
- All TUI output respects the selected theme
- Zero configuration produces identical output to today (backward compatible)
- Clean separation between theme definitions and TUI rendering logic

**Non-Goals:**
- Custom user-defined color palettes (future enhancement)
- Per-command theme overrides
- Terminal capability detection (assume 256 color support)
- True color (24-bit) support

## Decisions

### Decision: Centralized Theme Package

Create `internal/theme/` package with:
```go
type Theme struct {
    // Semantic color slots
    Primary        lipgloss.Color  // Main accent (headers, titles)
    Secondary      lipgloss.Color  // Secondary accent (cursors, selections)
    Success        lipgloss.Color  // Success states, checkmarks
    Error          lipgloss.Color  // Errors, warnings
    Warning        lipgloss.Color  // Caution indicators
    Muted          lipgloss.Color  // Dim/subtle text
    Border         lipgloss.Color  // Table borders, separators
    Header         lipgloss.Color  // Section headers
    Selected       lipgloss.Color  // Selected item foreground
    Highlight      lipgloss.Color  // Selected item background
    GradientStart  lipgloss.Color  // ASCII art gradient start
    GradientEnd    lipgloss.Color  // ASCII art gradient end
}
```

**Alternatives considered:**
- CSS-like variables: Adds complexity without benefit for preset-only system
- Per-component configuration: Fragments configuration, harder to maintain consistency

### Decision: Preset Themes Only (Initially)

Support 5 built-in themes:
- `default` - Current colors, optimized for dark terminals
- `dark` - High contrast on dark backgrounds
- `light` - Optimized for light terminal backgrounds
- `solarized` - Matches Solarized color palette
- `monokai` - Matches Monokai color palette

**Alternatives considered:**
- Full custom palettes: Overly complex for initial release, can add later
- Single accent color derivation: Limits flexibility, harder to get right

### Decision: Global Theme Instance

Theme loaded once at config load time, accessed via `theme.Current()` global:
```go
var current *Theme

func Load(name string) error {
    t, err := Get(name)
    if err != nil { return err }
    current = t
    return nil
}

func Current() *Theme {
    if current == nil { return defaultTheme }
    return current
}
```

**Alternatives considered:**
- Pass theme through context: Requires changing many function signatures
- Inject at component construction: Adds complexity to TUI model creation

## Risks / Trade-offs

**Risk:** Global mutable state for theme
- Mitigation: Theme is effectively immutable after config load; no runtime changes expected

**Risk:** Color choices may not work well in all terminals
- Mitigation: Use ANSI 256 codes (widely supported), provide presets known to work well

**Trade-off:** Preset-only limits customization
- Accepted: Full customization can be added later; presets cover 90% of use cases

## Migration Plan

1. Create `internal/theme/` package with all theme definitions
2. Update `internal/config/` to load theme setting
3. Migrate each TUI package one at a time:
   - `internal/tui/styles.go`
   - `internal/view/`
   - `internal/init/wizard.go`
4. Each migration is backward-compatible (default theme matches current colors)

**Rollback:** Remove `theme` field from config; revert to hardcoded colors

## Open Questions

None - preset theme approach is well-defined and self-contained.
