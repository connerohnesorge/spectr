// Package config provides configuration loading for Spectr projects.
package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	// ConfigFileName is the name of the Spectr configuration file.
	ConfigFileName = ".spectr.yaml"
	// DefaultDir is the default directory name for Spectr files.
	DefaultDir = "spectr"
)

// Config holds the Spectr configuration loaded from .spectr.yaml.
type Config struct {
	Dir string `yaml:"dir"`
}

// Validate checks Config fields for correctness.
func (c *Config) Validate() error {
	if c.Dir == "" {
		return errors.New("dir must not be empty")
	}

	if strings.HasPrefix(c.Dir, "/") {
		return errors.New("dir must be relative, not absolute")
	}

	if strings.Contains(c.Dir, "..") {
		return errors.New("dir must not contain path traversal")
	}

	return nil
}

// SpecsDir returns the path to the specs directory.
func (c *Config) SpecsDir() string {
	return c.Dir + "/specs"
}

// ChangesDir returns the path to the changes directory.
func (c *Config) ChangesDir() string {
	return c.Dir + "/changes"
}

// ProjectFile returns the path to the project.md file.
func (c *Config) ProjectFile() string {
	return c.Dir + "/project.md"
}

// AgentsFile returns the path to the AGENTS.md file.
func (c *Config) AgentsFile() string {
	return c.Dir + "/AGENTS.md"
}

// Load reads configuration from .spectr.yaml in the given project root.
// Returns default configuration if file doesn't exist.
func Load(projectRoot string) (*Config, error) {
	configPath := filepath.Join(projectRoot, ConfigFileName)

	// Default config
	cfg := &Config{Dir: DefaultDir}

	// Check if config file exists
	data, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		return cfg, nil // Return default, no error
	}

	if err != nil {
		return nil, err
	}

	// Parse YAML
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	// Apply default if empty
	if cfg.Dir == "" {
		cfg.Dir = DefaultDir
	}

	// Validate
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Default returns a Config with default values.
func Default() *Config {
	return &Config{Dir: DefaultDir}
}
