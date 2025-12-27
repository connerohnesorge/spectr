package providers

import (
	"errors"
	"strings"
)

// Config contains initialization configuration.
type Config struct {
	SpectrDir string // e.g., "spectr" (relative to fs root)
}

// Validate checks Config fields for basic correctness.
func (c *Config) Validate() error {
	if c.SpectrDir == "" {
		return errors.New("SpectrDir must not be empty")
	}
	if strings.HasPrefix(c.SpectrDir, "/") {
		return errors.New("SpectrDir must be relative, not absolute")
	}
	if strings.Contains(c.SpectrDir, "..") {
		return errors.New("SpectrDir must not contain path traversal")
	}

	return nil
}

// SpecsDir returns the specs directory path derived from SpectrDir.
func (c *Config) SpecsDir() string {
	return c.SpectrDir + "/specs"
}

// ChangesDir returns the changes directory path derived from SpectrDir.
func (c *Config) ChangesDir() string {
	return c.SpectrDir + "/changes"
}

// ProjectFile returns the project file path derived from SpectrDir.
func (c *Config) ProjectFile() string {
	return c.SpectrDir + "/project.md"
}

// AgentsFile returns the agents file path derived from SpectrDir.
func (c *Config) AgentsFile() string {
	return c.SpectrDir + "/AGENTS.md"
}
