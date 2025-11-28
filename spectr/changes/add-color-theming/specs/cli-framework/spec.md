## ADDED Requirements

### Requirement: Theme Configuration

The system SHALL support an optional `theme` setting in `spectr.yaml` that allows users to select from built-in preset color themes, applying the selected theme to all TUI output.

#### Scenario: Default theme when not specified
- **WHEN** user runs any spectr command without `theme` in `spectr.yaml`
- **THEN** the system uses the `default` theme
- **AND** all TUI output matches the current (pre-theming) color scheme
- **AND** no warning is displayed about missing theme configuration

#### Scenario: Valid theme selection
- **WHEN** `spectr.yaml` contains `theme: solarized`
- **AND** user runs any spectr command with TUI output
- **THEN** dashboard, wizard, interactive modes, and progress bars use Solarized colors
- **AND** theme is applied consistently across all visual output

#### Scenario: Invalid theme name
- **WHEN** `spectr.yaml` contains `theme: nonexistent`
- **AND** user runs any spectr command
- **THEN** display error: "invalid theme 'nonexistent', available themes: default, dark, light, solarized, monokai"
- **AND** command exits with non-zero status

### Requirement: Preset Theme Library

The system SHALL provide a set of built-in preset themes optimized for different terminal environments and user preferences.

#### Scenario: Default theme colors
- **WHEN** user selects `theme: default`
- **THEN** colors match the original hardcoded values
- **AND** headers use purple/violet tones
- **AND** success indicators use green
- **AND** error indicators use red
- **AND** muted text uses dim gray

#### Scenario: Dark theme colors
- **WHEN** user selects `theme: dark`
- **THEN** colors are optimized for dark terminal backgrounds
- **AND** high contrast is maintained for readability
- **AND** bright accent colors are used for visibility

#### Scenario: Light theme colors
- **WHEN** user selects `theme: light`
- **THEN** colors are optimized for light terminal backgrounds
- **AND** darker accent colors are used for visibility against light backgrounds
- **AND** text remains readable without eye strain

#### Scenario: Solarized theme colors
- **WHEN** user selects `theme: solarized`
- **THEN** colors match the Solarized Dark color palette
- **AND** base colors use Solarized base0, base01, base02
- **AND** accent colors use Solarized blue, green, yellow, red

#### Scenario: Monokai theme colors
- **WHEN** user selects `theme: monokai`
- **THEN** colors match the Monokai color palette
- **AND** accents use Monokai pink, green, orange, purple

### Requirement: Theme Application Scope

The theme SHALL apply to all TUI output components including dashboard view, init wizard, interactive selection modes, progress bars, and table styling.

#### Scenario: Dashboard view theming
- **WHEN** user runs `spectr view` with a non-default theme
- **THEN** section headers use theme's header color
- **AND** active change indicators use theme's warning color
- **AND** completed indicators use theme's success color
- **AND** spec indicators use theme's primary color
- **AND** progress bar filled portions use theme's success color
- **AND** progress bar empty portions use theme's muted color

#### Scenario: Init wizard theming
- **WHEN** user runs `spectr init` with a non-default theme
- **THEN** ASCII art gradient uses theme's gradient start and end colors
- **AND** titles use theme's primary color
- **AND** selected items use theme's secondary color
- **AND** cursor uses theme's highlight color
- **AND** success messages use theme's success color
- **AND** error messages use theme's error color

#### Scenario: Interactive mode theming
- **WHEN** user enters interactive mode (e.g., `spectr validate` without args)
- **THEN** table headers use theme's header color
- **AND** selected rows use theme's selected and highlight colors
- **AND** borders use theme's border color
- **AND** help text uses theme's muted color
