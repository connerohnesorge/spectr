// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
package providers

// FileInitializer creates or updates a single file for a provider.
// Each initializer is responsible for managing one specific file type
// (instruction files, slash commands, etc.).
type FileInitializer interface {
	// ID returns a unique identifier for this initializer.
	// Format: "{type}:{path}" e.g., "instruction:CLAUDE.md",
	// "markdown-cmd:.claude/commands/spectr/proposal.md"
	ID() string

	// FilePath returns the relative path this initializer manages.
	// May contain ~ for home directory (expanded internally during Configure).
	FilePath() string

	// Configure creates or updates the file.
	// projectPath is the root project directory.
	// tm is the template renderer for generating content.
	Configure(projectPath string, tm TemplateRenderer) error

	// IsConfigured checks if the file exists and is properly configured.
	// Returns true if the file exists at the expected location.
	IsConfigured(projectPath string) bool
}
