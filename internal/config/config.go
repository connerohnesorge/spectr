// Package config handles Spectr configuration file loading and validation.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/connerohnesorge/spectr/internal/theme"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultRootDir is the default directory name for Spectr files
	DefaultRootDir = "spectr"
	// ConfigFileName is the name of the Spectr configuration file
	ConfigFileName = "spectr.yaml"
)

// Config holds the Spectr configuration
type Config struct {
	// RootDir is the directory name where Spectr files are stored
	// (e.g., "spectr", "specs")
	RootDir string `yaml:"root_dir"`
	// ProjectRoot is the absolute path to the project root
	// (where spectr.yaml was found or where we're running from)
	ProjectRoot string `yaml:"-"`
	// Theme is the name of the color theme to use
	// (default, dark, light, solarized, monokai)
	Theme string `yaml:"theme"`
}

// Load searches for spectr.yaml starting from the current working directory,
// walking up the directory tree. If found, it parses the configuration.
// If not found, returns default configuration.
func Load() (*Config, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get current directory: %w", err)
	}

	return LoadFromPath(cwd)
}

// LoadFromPath searches for spectr.yaml starting from the given path,
// walking up the directory tree. If found, it parses the configuration.
// If not found, returns default configuration with startPath as ProjectRoot.
func LoadFromPath(startPath string) (*Config, error) {
	// Convert to absolute path
	absPath, err := filepath.Abs(startPath)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to resolve absolute path for %q: %w",
			startPath,
			err,
		)
	}

	// Walk up the directory tree looking for spectr.yaml
	currentPath := absPath
	for {
		configPath := filepath.Join(currentPath, ConfigFileName)

		// Check if config file exists
		if _, err := os.Stat(configPath); err == nil {
			// Found config file, parse it
			cfg, err := parseConfigFile(configPath)
			if err != nil {
				return nil, err
			}
			cfg.ProjectRoot = currentPath

			// Validate configuration
			if err := cfg.validate(); err != nil {
				return nil, fmt.Errorf(
					"invalid configuration in %s: %w",
					configPath,
					err,
				)
			}

			return cfg, nil
		}

		// Move up one directory
		parentPath := filepath.Dir(currentPath)
		if parentPath == currentPath {
			// Reached root directory without finding config
			break
		}
		currentPath = parentPath
	}

	// No config file found, return defaults
	return &Config{
		RootDir:     DefaultRootDir,
		ProjectRoot: absPath,
		Theme:       "default",
	}, nil
}

// parseConfigFile reads and parses a spectr.yaml file
func parseConfigFile(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		// Try to provide better error messages for YAML syntax errors
		var yamlErr *yaml.TypeError
		if errors.As(err, &yamlErr) {
			return nil, fmt.Errorf("invalid YAML syntax: %v", yamlErr.Errors)
		}

		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Apply defaults if not set
	if cfg.RootDir == "" {
		cfg.RootDir = DefaultRootDir
	}
	if cfg.Theme == "" {
		cfg.Theme = "default"
	}

	return &cfg, nil
}

// validate checks if the configuration is valid
func (c *Config) validate() error {
	if c.RootDir == "" {
		return errors.New("root_dir cannot be empty")
	}

	// Check for invalid characters in root_dir
	// Must be a simple directory name (no path separators, .., *)
	invalidChars := []string{"/", "\\", "..", "*"}
	var foundInvalid []string

	for _, char := range invalidChars {
		if strings.Contains(c.RootDir, char) {
			foundInvalid = append(foundInvalid, char)
		}
	}

	if len(foundInvalid) > 0 {
		return fmt.Errorf(
			"root_dir must be a simple directory name "+
				"(found invalid characters: %s)",
			strings.Join(foundInvalid, ", "),
		)
	}

	// Additional check: root_dir shouldn't start with . (hidden directory)
	if strings.HasPrefix(c.RootDir, ".") {
		return errors.New(
			"root_dir cannot start with '.' (hidden directories not allowed)",
		)
	}

	// Validate theme name
	if _, err := theme.Get(c.Theme); err != nil {
		available := theme.Available()

		return fmt.Errorf(
			"invalid theme '%s', available themes: %s",
			c.Theme,
			strings.Join(available, ", "),
		)
	}

	return nil
}

// RootPath returns the absolute path to the Spectr root directory
func (c *Config) RootPath() string {
	return filepath.Join(c.ProjectRoot, c.RootDir)
}

// SpecsPath returns the absolute path to the specs directory
func (c *Config) SpecsPath() string {
	return filepath.Join(c.RootPath(), "specs")
}

// ChangesPath returns the absolute path to the changes directory
func (c *Config) ChangesPath() string {
	return filepath.Join(c.RootPath(), "changes")
}
