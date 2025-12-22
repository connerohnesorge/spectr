// Package types defines the shared types for the initialize package.
package types

import (
	"path/filepath"
)

// Config holds the configuration for provider initialization.
type Config struct {
	// SpectrDir is the path to the spectr directory (default: "spectr")
	SpectrDir string
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() *Config {
	return &Config{
		SpectrDir: "spectr",
	}
}

// SpecsDir returns the path to the specs directory.
func (c *Config) SpecsDir() string {
	return filepath.Join(c.SpectrDir, "specs")
}

// ChangesDir returns the path to the changes directory.
func (c *Config) ChangesDir() string {
	return filepath.Join(c.SpectrDir, "changes")
}

// ProjectFile returns the path to the project configuration file.
func (c *Config) ProjectFile() string {
	return filepath.Join(
		c.SpectrDir,
		"project.md",
	)
}

// AgentsFile returns the path to the agents file.
func (c *Config) AgentsFile() string {
	return filepath.Join(c.SpectrDir, "AGENTS.md")
}