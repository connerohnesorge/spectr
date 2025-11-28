// Package config provides configuration loading and discovery for Spectr.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// ConfigFileName is the name of the configuration file
	ConfigFileName = "spectr.yaml"
	// DefaultRootDir is the default spectr directory name
	DefaultRootDir = "spectr"
)

var (
	// ErrInvalidRootDir is returned when root_dir contains invalid characters
	ErrInvalidRootDir = errors.New(
		"root_dir must be a simple directory name without path separators",
	)
)

// Config holds Spectr configuration settings
type Config struct {
	// RootDir is the name of the spectr root directory (default: "spectr")
	RootDir string `yaml:"root_dir"`

	// ProjectRoot is the absolute path to the project root directory
	// (where spectr.yaml was found, or where spectr/ directory is)
	ProjectRoot string `yaml:"-"`

	// ConfigPath is the path to the config file (empty if using defaults)
	ConfigPath string `yaml:"-"`
}

// Load discovers and loads the configuration.
// It walks up from startDir looking for spectr.yaml.
// Returns default config if no config file is found.
func Load(startDir string) (*Config, error) {
	// Convert to absolute path
	absStartDir, err := filepath.Abs(startDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	// Walk up directory tree looking for spectr.yaml
	currentDir := absStartDir
	for {
		configPath := filepath.Join(currentDir, ConfigFileName)

		// Check if config file exists
		if _, err := os.Stat(configPath); err == nil {
			// Config file found, load it
			cfg, err := loadConfigFile(configPath)
			if err != nil {
				return nil, err
			}

			cfg.ProjectRoot = currentDir
			cfg.ConfigPath = configPath

			// Validate root_dir
			if err := validateRootDir(cfg.RootDir); err != nil {
				return nil, fmt.Errorf(
					"invalid configuration in %s: %w",
					configPath,
					err,
				)
			}

			return cfg, nil
		}

		// Move up one directory
		parent := filepath.Dir(currentDir)

		// Check if we've reached the filesystem root
		if parent == currentDir {
			break
		}

		currentDir = parent
	}

	// No config file found, return default config
	return &Config{
		RootDir:     DefaultRootDir,
		ProjectRoot: absStartDir,
		ConfigPath:  "",
	}, nil
}

// loadConfigFile reads and parses the config file
func loadConfigFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		// Try to extract line/column information from YAML error
		var yamlErr *yaml.TypeError
		if errors.As(err, &yamlErr) {
			return nil, fmt.Errorf("YAML parsing error: %s", yamlErr.Error())
		}

		return nil, fmt.Errorf("invalid YAML syntax: %w", err)
	}

	// If root_dir is empty, use default
	if cfg.RootDir == "" {
		cfg.RootDir = DefaultRootDir
	}

	return &cfg, nil
}

// validateRootDir ensures root_dir is a simple directory name
func validateRootDir(rootDir string) error {
	if rootDir == "" {
		return nil // Will use default
	}

	invalidChars := []string{"/", "\\", "..", "*", "?"}
	for _, char := range invalidChars {
		if strings.Contains(rootDir, char) {
			return fmt.Errorf("%w: found '%s'", ErrInvalidRootDir, char)
		}
	}

	return nil
}

// SpectrRoot returns the absolute path to the spectr root directory
func (c *Config) SpectrRoot() string {
	return filepath.Join(c.ProjectRoot, c.RootDir)
}

// ChangesDir returns the path to the changes directory
func (c *Config) ChangesDir() string {
	return filepath.Join(c.SpectrRoot(), "changes")
}

// SpecsDir returns the path to the specs directory
func (c *Config) SpecsDir() string {
	return filepath.Join(c.SpectrRoot(), "specs")
}

// ArchiveDir returns the path to the archive directory
func (c *Config) ArchiveDir() string {
	return filepath.Join(c.ChangesDir(), "archive")
}
