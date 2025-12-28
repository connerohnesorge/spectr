package providers

import (
	"errors"
	"strings"
)

// Config holds configuration for provider initialization.
// All paths are derived from SpectrDir to ensure consistency.
type Config struct {
	// SpectrDir is the base directory for spectr files,
	// relative to the project root.
	// Example: "spectr"
	// Must not be empty, absolute, or contain path traversal.
	SpectrDir string
}

// Validate checks Config fields for basic correctness.
// Returns an error if any validation rule fails.
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

// SpecsDir returns the directory for spec files.
// Derived from SpectrDir: {SpectrDir}/specs
func (c *Config) SpecsDir() string {
	return c.SpectrDir + "/specs"
}

// ChangesDir returns the directory for change proposals.
// Derived from SpectrDir: {SpectrDir}/changes
func (c *Config) ChangesDir() string {
	return c.SpectrDir + "/changes"
}

// ProjectFile returns the path to the project configuration file.
// Derived from SpectrDir: {SpectrDir}/project.md
func (c *Config) ProjectFile() string {
	return c.SpectrDir + "/project.md"
}

// AgentsFile returns the path to the agents file.
// Derived from SpectrDir: {SpectrDir}/AGENTS.md
func (c *Config) AgentsFile() string {
	return c.SpectrDir + "/AGENTS.md"
}
