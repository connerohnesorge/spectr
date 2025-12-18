// Package providers defines the core interfaces for the provider architecture.
//
// This file contains the Config struct which holds initialization configuration
// and provides methods for computing derived paths.
package providers

import "path"

// Config contains initialization configuration for spectr providers.
//
// Config holds the base spectr directory path and provides methods for
// computing derived paths (specs, changes, project file, agents file).
// All paths are relative to the filesystem root (either project root for
// project files, or home directory for global files).
//
// # Example Usage
//
//	cfg := &Config{SpectrDir: "spectr"}
//	specsDir := cfg.SpecsDir()       // "spectr/specs"
//	changesDir := cfg.ChangesDir()   // "spectr/changes"
//	projectFile := cfg.ProjectFile() // "spectr/project.md"
//	agentsFile := cfg.AgentsFile()   // "spectr/AGENTS.md"
//
// # Design Principles
//
// 1. **Single Source of Truth**: Only SpectrDir is stored; all other paths
// are computed to avoid redundancy and ensure consistency.
//
// 2. **Path Joining**: Uses path.Join for clean path concatenation that works
// correctly with the afero filesystem abstraction.
type Config struct {
	// SpectrDir is the base directory for spectr files (e.g., "spectr").
	// All other paths are derived from this value.
	// This is relative to the filesystem root (project directory).
	SpectrDir string
}

// SpecsDir returns the path to the specs directory.
// This is where capability specifications are stored.
//
// Example: "spectr/specs"
func (c *Config) SpecsDir() string {
	return path.Join(c.SpectrDir, "specs")
}

// ChangesDir returns the path to the changes directory.
// This is where change proposals are stored.
//
// Example: "spectr/changes"
func (c *Config) ChangesDir() string {
	return path.Join(c.SpectrDir, "changes")
}

// ProjectFile returns the path to the project configuration file.
// This file contains project-level spectr settings.
//
// Example: "spectr/project.md"
func (c *Config) ProjectFile() string {
	return path.Join(c.SpectrDir, "project.md")
}

// AgentsFile returns the path to the agents file.
// This file contains instructions for AI assistants working in the project.
//
// Example: "spectr/AGENTS.md"
func (c *Config) AgentsFile() string {
	return path.Join(c.SpectrDir, "AGENTS.md")
}

// TemplateContext holds path-related template variables for dynamic directory names.
// This struct provides all the path information needed for rendering templates
// in instruction files and slash commands.
//
// # Usage
//
// Create a TemplateContext from a Config using NewTemplateContext, or use
// DefaultTemplateContext for the standard "spectr" base directory.
//
// Example:
//
//	cfg := &Config{SpectrDir: "spectr"}
//	ctx := NewTemplateContext(cfg)
type TemplateContext struct {
	// BaseDir is the base directory for spectr files (default: "spectr")
	BaseDir string
	// SpecsDir is the directory for spec files (default: "spectr/specs")
	SpecsDir string
	// ChangesDir is the directory for change proposals (default: "spectr/changes")
	ChangesDir string
	// ProjectFile is the path to the project configuration file (default: "spectr/project.md")
	ProjectFile string
	// AgentsFile is the path to the agents file (default: "spectr/AGENTS.md")
	AgentsFile string
}

// NewTemplateContext creates a TemplateContext from a Config.
// This ensures the template context paths are consistent with the configuration.
func NewTemplateContext(cfg *Config) TemplateContext {
	return TemplateContext{
		BaseDir:     cfg.SpectrDir,
		SpecsDir:    cfg.SpecsDir(),
		ChangesDir:  cfg.ChangesDir(),
		ProjectFile: cfg.ProjectFile(),
		AgentsFile:  cfg.AgentsFile(),
	}
}

// DefaultTemplateContext returns a TemplateContext with default values.
// Uses "spectr" as the base directory.
func DefaultTemplateContext() TemplateContext {
	cfg := &Config{SpectrDir: "spectr"}
	return NewTemplateContext(cfg)
}
