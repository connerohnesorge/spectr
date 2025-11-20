// Package providerkit provides shared utilities, types, and interfaces for
// implementing AI tool providers in Spectr.
//
// This package serves as the foundation for all provider implementations,
// offering:
//   - Provider interface definition (alias for Configurator)
//   - Marker-based file update utilities
//   - Template rendering capabilities
//   - Filesystem helper functions
//   - Base slash command provider implementation
//
// Provider implementations should use this package to ensure consistent
// behavior across all AI tools supported by Spectr.
package providerkit

// Provider interface defines the contract that all AI tool providers
// must implement. This interface is an alias for the Configurator interface
// to maintain backward compatibility while improving naming clarity.
//
// Providers are responsible for:
//   - Configuring AI tools in a project (creating/updating instruction files)
//   - Checking if a tool is already configured
//   - Providing a human-readable name for the tool
//
// Implementation types:
//   - Config-based providers: Create single instruction files (CLAUDE.md)
//   - Slash command providers: Create multiple command files in directories
type Provider interface {
	// Configure configures an AI tool for the given project path.
	// It creates or updates the necessary configuration files.
	//
	// Parameters:
	//   - projectPath: Root directory of the project being configured
	//   - spectrDir: Path to Spectr directory (usually projectPath/spectr)
	//
	// Returns an error if configuration fails.
	Configure(projectPath, spectrDir string) error

	// IsConfigured checks if the AI tool is already configured for the
	// given project.
	//
	// Parameters:
	//   - projectPath: The root directory of the project to check
	//
	// Returns true if the tool is configured, false otherwise.
	IsConfigured(projectPath string) bool

	// GetName returns the human-readable name of the AI tool.
	// This name is displayed to users during initialization and in
	// CLI output.
	//
	// Returns a string like "Claude Code", "Cline", etc.
	GetName() string
}

// Configurator is an alias for Provider, maintained for backward
// compatibility with existing code that references the Configurator interface.
//
// New code should prefer using Provider, but both names refer to the
// same interface.
type Configurator = Provider
