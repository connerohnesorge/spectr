// Package providers implements the interface-driven provider architecture for
// AI CLI/IDE/Orchestrator tools.
package providers

// FileInitializer creates or updates a single file for a provider.
//
// Each initializer manages exactly one file and encapsulates all logic for
// creating, updating, and checking that file's configuration state.
//
// Implementations include:
//   - InstructionFileInitializer: Manages instruction files (e.g., CLAUDE.md)
//   - MarkdownSlashCommandInitializer: Manages markdown slash commands
//   - TOMLSlashCommandInitializer: Manages TOML slash commands (Gemini)
//
// ID Format:
//
//	"{type}:{path}" where type indicates the initializer kind:
//	  - "instruction:CLAUDE.md"
//	  - "markdown-cmd:.claude/commands/spectr/proposal.md"
//	  - "toml-cmd:.gemini/commands/spectr/proposal.toml"
type FileInitializer interface {
	// ID returns a unique identifier for this initializer.
	// Format: "{type}:{path}" e.g., "instruction:CLAUDE.md"
	//
	// The ID is used for deduplication when multiple providers
	// share the same file path.
	ID() string

	// FilePath returns the relative path this initializer manages.
	// May contain ~ for home directory (expanded internally during Configure).
	//
	// Examples:
	//   - "CLAUDE.md" (project-relative)
	//   - "~/.codex/CODEX.md" (global, home-relative)
	//   - ".claude/commands/spectr/proposal.md" (project-relative)
	FilePath() string

	// Configure creates or updates the file.
	//
	// Path expansion (~ to home directory) is handled internally.
	// For project-relative paths, projectPath is joined with FilePath().
	// For global paths (starting with ~), the path is expanded independently.
	//
	// The TemplateRenderer is used to render file content from templates.
	//
	// Returns an error if file creation/update fails.
	Configure(projectPath string, tm TemplateRenderer) error

	// IsConfigured checks if the file exists and is properly configured.
	//
	// Path expansion follows the same rules as Configure().
	// Returns true if the file exists at the expected location.
	IsConfigured(projectPath string) bool
}
