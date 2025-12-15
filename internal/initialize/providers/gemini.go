package providers

// init registers the Gemini CLI provider with the global registry.
func init() {
	Register(&GeminiProvider{})
}

// GeminiProvider implements the Provider interface for Gemini CLI.
// Gemini uses .gemini/commands/ for TOML-based slash commands.
// It does not use a separate instruction file.
type GeminiProvider struct{}

// ID returns the unique identifier for the Gemini CLI provider.
func (*GeminiProvider) ID() string { return "gemini" }

// Name returns the display name for Gemini CLI.
func (*GeminiProvider) Name() string { return "Gemini CLI" }

// Priority returns the display order for Gemini CLI.
func (*GeminiProvider) Priority() int { return PriorityGemini }

// Initializers returns the file initializers for the Gemini CLI provider.
func (*GeminiProvider) Initializers() []FileInitializer {
	proposalPath, applyPath := StandardCommandPaths(
		".gemini/commands",
		".toml",
	)

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

func (p *GeminiProvider) IsConfigured(projectPath string) bool {
	return AreInitializersConfigured(p.Initializers(), projectPath)
}

func (p *GeminiProvider) GetFilePaths() []string {
	return GetInitializerPaths(p.Initializers())
}
