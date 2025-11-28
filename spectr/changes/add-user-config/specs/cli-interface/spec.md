## ADDED Requirements

### Requirement: User Configuration File

The system SHALL support an optional user configuration file at the XDG-compliant location `$XDG_CONFIG_HOME/spectr/config.yaml` (defaulting to `~/.config/spectr/config.yaml` when `$XDG_CONFIG_HOME` is unset).

#### Scenario: Config file exists and is valid
- **WHEN** a valid config file exists at the user config path
- **THEN** the system loads and applies the configuration settings

#### Scenario: Config file does not exist
- **WHEN** no config file exists at the user config path
- **THEN** the system uses default values for all settings

#### Scenario: Config file is invalid YAML
- **WHEN** the config file contains invalid YAML syntax
- **THEN** the system logs a warning and uses default values

### Requirement: Color Theme Overrides

The system SHALL allow users to override TUI colors via the `theme` section of the user config file, supporting both hex color codes (#RRGGBB) and ANSI 256 color codes (0-255).

#### Scenario: Custom accent color
- **WHEN** the config file specifies `theme.accent: "#FF5733"`
- **THEN** the TUI uses the specified hex color for accent elements

#### Scenario: ANSI 256 color code
- **WHEN** the config file specifies `theme.border: "240"`
- **THEN** the TUI uses ANSI color code 240 for border elements

#### Scenario: Invalid color format
- **WHEN** the config file specifies an invalid color value (e.g., `theme.accent: "not-a-color"`)
- **THEN** the system logs a warning and uses the default color for that field

#### Scenario: Partial theme override
- **WHEN** the config file only overrides some theme colors (e.g., only `accent` and `error`)
- **THEN** the system uses the specified overrides and defaults for unspecified colors

### Requirement: Config Command

The system SHALL provide a `spectr config` subcommand to display the current configuration and its source location.

#### Scenario: Display config when file exists
- **WHEN** the user runs `spectr config`
- **AND** a config file exists
- **THEN** the system displays the loaded configuration and the file path

#### Scenario: Display config when no file exists
- **WHEN** the user runs `spectr config`
- **AND** no config file exists
- **THEN** the system displays the default configuration and indicates no config file was found

### Requirement: Theme Color Fields

The theme configuration SHALL support the following color fields: `accent`, `error`, `success`, `border`, `help`, `selected`, `highlight`, and `header`.

#### Scenario: All theme fields are configurable
- **WHEN** the user specifies all eight color fields in the config
- **THEN** the TUI applies each color to its corresponding UI element

#### Scenario: Default theme colors match current behavior
- **WHEN** no theme overrides are specified
- **THEN** the TUI colors match the current hardcoded defaults (border: 240, header: 99, selected: 229, highlight: 57, help: 240)
