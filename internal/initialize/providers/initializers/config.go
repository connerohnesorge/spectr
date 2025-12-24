// Package initializers provides composable initialization components
// for the provider architecture.
package initializers

import (
	"path/filepath"
	"text/template"
)

// Marker constants for managing config file updates.
const (
	SpectrStartMarker = "<!-- spectr:START -->"
	SpectrEndMarker   = "<!-- spectr:END -->"
)

// File permission constants.
const (
	// DirPerm is the permission mode for created directories (rwxr-xr-x).
	DirPerm = 0755

	// FilePerm is the permission mode for created files (rw-r--r--).
	FilePerm = 0644

	// Newline is the newline character used in file content.
	Newline = "\n"
)

// TemplateProvider provides access to templates for rendering.
// This interface is implemented by the TemplateManager to avoid import cycles.
type TemplateProvider interface {
	// GetTemplates returns the underlying template set for rendering.
	GetTemplates() *template.Template
}

// Config contains provider configuration settings.
// All paths are relative to the filesystem root (either project or global).
//
// The SpectrDir is the single source of truth; all other paths are derived
// from it using helper methods. This ensures consistency and prevents
// path-related bugs.
type Config struct {
	// SpectrDir is the base directory for spectr files (e.g., "spectr")
	// All other paths are derived from this value.
	SpectrDir string
}

// SpecsDir returns the path to the specs directory.
// Example: "spectr/specs"
func (c *Config) SpecsDir() string {
	return filepath.Join(c.SpectrDir, "specs")
}

// ChangesDir returns the path to the changes directory.
// Example: "spectr/changes"
func (c *Config) ChangesDir() string {
	return filepath.Join(c.SpectrDir, "changes")
}

// ProjectFile returns the path to the project.md file.
// Example: "spectr/project.md"
func (c *Config) ProjectFile() string {
	return filepath.Join(c.SpectrDir, "project.md")
}

// AgentsFile returns the path to the AGENTS.md file.
// Example: "spectr/AGENTS.md"
func (c *Config) AgentsFile() string {
	return filepath.Join(c.SpectrDir, "AGENTS.md")
}

// NewConfig creates a new Config with the given base directory.
// The baseDir should be relative to the filesystem root (e.g., "spectr").
func NewConfig(baseDir string) *Config {
	return &Config{
		SpectrDir: baseDir,
	}
}

// DefaultConfig returns a Config with default values ("spectr").
func DefaultConfig() *Config {
	return NewConfig("spectr")
}
