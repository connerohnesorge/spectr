// Package config provides user configuration types and defaults for Spectr CLI.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ConfigPath returns the expected path for the Spectr config file.
// It follows XDG Base Directory specification:
//   - Uses $XDG_CONFIG_HOME/spectr/config.yaml if XDG_CONFIG_HOME
//     is set
//   - Falls back to ~/.config/spectr/config.yaml otherwise
//
// The path is returned even if the file doesn't exist
// (useful for display purposes).
func ConfigPath() string {
	configDir := os.Getenv("XDG_CONFIG_HOME")
	if configDir == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			// Fallback to a reasonable default if home directory
			// cannot be determined
			return filepath.Join(".config", "spectr", "config.yaml")
		}
		configDir = filepath.Join(homeDir, ".config")
	}

	return filepath.Join(configDir, "spectr", "config.yaml")
}

// Load loads the user configuration from the default config path.
// If the config file doesn't exist, it returns a Config with default
// values (no error). If the config file exists but is invalid YAML,
// it logs a warning and returns defaults.
func Load() (*Config, error) {
	return LoadFromPath(ConfigPath())
}

// LoadFromPath loads the user configuration from a specific path.
// If the file doesn't exist, it returns a Config with default values
// (no error). If the file exists but is invalid YAML, it returns an
// error.
func LoadFromPath(path string) (*Config, error) {
	// Check if file exists
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// File doesn't exist - return defaults with no error
		return defaultConfig(), nil
	}
	if err != nil {
		// Other stat error (permissions, etc.)
		return nil, fmt.Errorf("failed to check config file: %w", err)
	}

	// Read and parse the file
	cfg, err := readAndParseConfig(path)
	if err != nil {
		return nil, err
	}

	// Validate and fix invalid theme colors
	fixInvalidThemeColors(&cfg)

	// Merge with defaults (preserve defaults for unset fields)
	return mergeWithDefaults(&cfg), nil
}

// readAndParseConfig reads and parses the YAML config file.
func readAndParseConfig(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		// Invalid YAML - log warning and return defaults
		fmt.Fprintf(
			os.Stderr,
			"Warning: Invalid YAML in config file %s: %v\n",
			path,
			err,
		)
		fmt.Fprint(os.Stderr, "Using default configuration.\n")

		return *defaultConfig(), nil
	}

	return cfg, nil
}

// fixInvalidThemeColors validates theme colors and resets invalid ones.
func fixInvalidThemeColors(cfg *Config) {
	validationErrors := validateThemeWithFieldTracking(cfg.Theme)
	if len(validationErrors) == 0 {
		return
	}

	defaults := DefaultTheme()
	for _, fieldErr := range validationErrors {
		fmt.Fprintf(os.Stderr, "Warning: %v\n", fieldErr.err)
		resetThemeField(&cfg.Theme, fieldErr.fieldName, defaults)
	}
}

// resetThemeField resets a specific theme field to its default value.
func resetThemeField(theme *Theme, fieldName string, defaults Theme) {
	switch fieldName {
	case "accent":
		theme.Accent = defaults.Accent
	case "error":
		theme.Error = defaults.Error
	case "success":
		theme.Success = defaults.Success
	case "border":
		theme.Border = defaults.Border
	case "help":
		theme.Help = defaults.Help
	case "selected":
		theme.Selected = defaults.Selected
	case "highlight":
		theme.Highlight = defaults.Highlight
	case "header":
		theme.Header = defaults.Header
	}
}

// defaultConfig returns a Config with all default values.
func defaultConfig() *Config {
	theme := DefaultTheme()

	return &Config{
		Theme: theme,
	}
}

// mergeWithDefaults takes a user config and fills in any nil theme
// fields with defaults. This enables partial override - users only need
// to specify the colors they want to change.
func mergeWithDefaults(userConfig *Config) *Config {
	defaults := DefaultTheme()
	merged := userConfig.Theme

	// Override each field only if it's nil (unset)
	if merged.Accent == nil {
		merged.Accent = defaults.Accent
	}
	if merged.Error == nil {
		merged.Error = defaults.Error
	}
	if merged.Success == nil {
		merged.Success = defaults.Success
	}
	if merged.Border == nil {
		merged.Border = defaults.Border
	}
	if merged.Help == nil {
		merged.Help = defaults.Help
	}
	if merged.Selected == nil {
		merged.Selected = defaults.Selected
	}
	if merged.Highlight == nil {
		merged.Highlight = defaults.Highlight
	}
	if merged.Header == nil {
		merged.Header = defaults.Header
	}

	return &Config{
		Theme: merged,
	}
}
