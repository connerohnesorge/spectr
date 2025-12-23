package providers

// Config holds configuration for provider initialization.
// All paths are derived from SpectrDir to maintain a single source of truth.
type Config struct {
	// SpectrDir is the base directory for spectr files (e.g., "spectr").
	// All other paths are derived from this value.
	SpectrDir string
}

// SpecsDir returns the directory for spec files.
// Derived from SpectrDir: <SpectrDir>/specs
func (c *Config) SpecsDir() string {
	return c.SpectrDir + "/specs"
}

// ChangesDir returns the directory for change proposals.
// Derived from SpectrDir: <SpectrDir>/changes
func (c *Config) ChangesDir() string {
	return c.SpectrDir + "/changes"
}

// ProjectFile returns the path to the project configuration file.
// Derived from SpectrDir: <SpectrDir>/project.md
func (c *Config) ProjectFile() string {
	return c.SpectrDir + "/project.md"
}

// AgentsFile returns the path to the agents file.
// Derived from SpectrDir: <SpectrDir>/AGENTS.md
func (c *Config) AgentsFile() string {
	return c.SpectrDir + "/AGENTS.md"
}
