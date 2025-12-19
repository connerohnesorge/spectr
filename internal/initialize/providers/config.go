// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
//
// This file defines the Config struct that holds configuration for spectr
// initialization, including the base spectr directory and derived path methods.
package providers

// DefaultSpectrDir is the default directory name for spectr files.
const DefaultSpectrDir = "spectr"

// Config holds configuration for spectr initialization.
//
// The Config struct provides:
//   - SpectrDir: the base directory for spectr files (e.g., "spectr")
//
// And derived path methods:
//   - SpecsDir(): returns SpectrDir + "/specs"
//   - ChangesDir(): returns SpectrDir + "/changes"
//   - ProjectFile(): returns SpectrDir + "/project.md"
//   - AgentsFile(): returns SpectrDir + "/AGENTS.md"
//
//nolint:revive // line-length-limit - struct documentation
type Config struct {
	// SpectrDir is the base directory for spectr files (relative to project).
	// Default: "spectr"
	SpectrDir string
}

// NewConfig creates a new Config with the given spectr directory.
// If spectrDir is empty, it defaults to "spectr".
func NewConfig(spectrDir string) *Config {
	dir := spectrDir
	if dir == "" {
		dir = DefaultSpectrDir
	}

	return &Config{
		SpectrDir: dir,
	}
}

// SpecsDir returns the path to the specs directory.
func (c *Config) SpecsDir() string { return c.SpectrDir + "/specs" }

// ChangesDir returns the path to the changes directory.
func (c *Config) ChangesDir() string { return c.SpectrDir + "/changes" }

// ProjectFile returns the path to the project configuration file.
func (c *Config) ProjectFile() string { return c.SpectrDir + "/project.md" }

// AgentsFile returns the path to the AGENTS.md file.
func (c *Config) AgentsFile() string { return c.SpectrDir + "/AGENTS.md" }
