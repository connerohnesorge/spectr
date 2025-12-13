//nolint:revive // unused-receiver acceptable for interface compliance
package providers

func init() {
	Register(&GeminiProvider{})
}

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses .gemini/commands/ for TOML-based slash commands
// (no instruction file).
type GeminiProvider struct{}

// ID returns the unique identifier for this provider.
func (p *GeminiProvider) ID() string { return "gemini" }

// Name returns the human-readable name for display.
func (p *GeminiProvider) Name() string { return "Gemini CLI" }

// Priority returns the display order (lower = higher priority).
func (p *GeminiProvider) Priority() int { return PriorityGemini }

// Initializers returns the file initializers for this provider.
func (p *GeminiProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(".gemini/commands", ".toml")

	return []FileInitializer{
		NewTOMLSlashCommandInitializer(
			proposalPath,
			"proposal",
			"Scaffold a new Spectr change and validate strictly.",
		),
		NewTOMLSlashCommandInitializer(
			applyPath,
			"apply",
			"Implement an approved Spectr change and keep tasks in sync.",
		),
	}
}

// IsConfigured checks if all files for this provider exist.
func (p *GeminiProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

// GetFilePaths returns the file paths managed by this provider.
func (p *GeminiProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
