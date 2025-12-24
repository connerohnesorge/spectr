package providers

import "path"

// Config contains configuration for initializers, specifying the base
// directory for the Spectr project structure.
//
// All path methods return paths relative to the filesystem root
// (projectFs or globalFs). Paths use forward slashes regardless of OS,
// as they work with afero.Fs abstractions.
type Config struct {
	// SpectrDir is the base directory for Spectr files (e.g., "spectr").
	// All other paths are derived from this base directory.
	// This is relative to the filesystem root (projectFs).
	SpectrDir string
}

// NewDefaultConfig returns a Config with default Spectr directory ("spectr").
func NewDefaultConfig() *Config {
	return &Config{
		SpectrDir: "spectr",
	}
}

// SpecsDir returns the path to the specs directory.
// Specs are stored at <SpectrDir>/specs.
//
// Example: "spectr/specs"
func (c *Config) SpecsDir() string {
	return path.Join(c.SpectrDir, "specs")
}

// ChangesDir returns the path to the changes directory.
// Change proposals are stored at <SpectrDir>/changes.
//
// Example: "spectr/changes"
func (c *Config) ChangesDir() string {
	return path.Join(c.SpectrDir, "changes")
}

// ProjectFile returns the path to the project.md file.
// The project file is stored at <SpectrDir>/project.md.
//
// Example: "spectr/project.md"
func (c *Config) ProjectFile() string {
	return path.Join(c.SpectrDir, "project.md")
}

// AgentsFile returns the path to the AGENTS.md file.
// The agents file is stored at <SpectrDir>/AGENTS.md.
//
// Example: "spectr/AGENTS.md"
func (c *Config) AgentsFile() string {
	return path.Join(c.SpectrDir, "AGENTS.md")
}
